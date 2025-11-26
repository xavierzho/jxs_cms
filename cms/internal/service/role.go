package service

import (
	"context"
	"errors"
	"fmt"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/internal/form"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"
	"data_backend/pkg/redisdb"
	"data_backend/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleSvc struct {
	ctx      *gin.Context
	engine   *gorm.DB
	rdb      *redisdb.RedisClient
	logger   *logger.Logger
	userSvc  *UserSvc
	dao      *dao.RoleDao
	newAlarm func(log *logger.Logger) message.Alarm
}

func NewRoleSvc(ctx *gin.Context, engine *gorm.DB, rdb *redisdb.RedisClient, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *RoleSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".RoleSvc")))
	return &RoleSvc{
		ctx:      ctx,
		engine:   engine,
		rdb:      rdb,
		logger:   log,
		dao:      dao.NewRoleDao(engine, log),
		newAlarm: newAlarm,
	}
}

func (svc *RoleSvc) Create(params *form.RoleCreateRequest) *errcode.Error {
	permDao := dao.NewPermissionDao(svc.engine, svc.logger)

	if flag, err := svc.checkExist(0, params.Name); err != nil {
		return errcode.CreateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	} else if flag {
		return iErrcode.RoleExist.WithDetails("duplicate role name")
	}

	permList, err := permDao.All(database.QueryWhereGroup{{Prefix: "id", Value: []interface{}{params.PermissionIDList}}})
	if err != nil {
		return errcode.CreateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}

	data := &dao.Role{Name: params.Name, Permission: permList}
	if err := svc.dao.Create(data); err != nil {
		return errcode.CreateFail.WithDetails(err.Error())
	}

	go svc.cachedRolePerm(data.ID)

	return nil
}

func (svc *RoleSvc) checkExist(id uint32, name string) (bool, error) {
	_, err := svc.dao.First(database.QueryWhereGroup{
		{Prefix: "id <> ? and name = ?", Value: []interface{}{id, name}},
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (svc *RoleSvc) List(params *form.RoleListRequest) (data []*dao.Role, count int64, e *errcode.Error) {
	queryParams := params.Parse()

	data, count, err := svc.dao.List(queryParams, app.GetPager(svc.ctx))
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, count, nil
}

func (svc *RoleSvc) All(queryParams database.QueryWhereGroup) ([]*dao.Role, *errcode.Error) {
	data, err := svc.dao.All(queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, nil
}

func (svc *RoleSvc) Update(id uint32, params *form.RoleUpdateRequest) (permList []string, e *errcode.Error) {
	svc.userSvc = NewUserSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)
	permDao := dao.NewPermissionDao(svc.engine, svc.logger)

	// 执行用户
	operator, e := svc.userSvc.CurrentUser()
	if e.Is(errcode.UnauthorizedTokenError) {
		return nil, e
	} else if e != nil {
		return nil, errcode.UpdateFail.WithDetails(e.Error())
	}

	// 查询要修改的角色
	role, err := svc.dao.First([]database.QueryWhere{{Prefix: "id", Value: []interface{}{id}}})
	if err != nil {
		return nil, errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}
	// 若新改的名字存在 则返回错误
	flag, err := svc.checkExist(role.ID, params.Name)
	if err != nil {
		return nil, errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}
	if flag {
		return nil, iErrcode.RoleExist.WithDetails("duplicate role name")
	}
	role.Name = params.Name
	var rolePermission []uint32
	for _, item := range role.Permission {
		rolePermission = append(rolePermission, item.ID)
	}

	if diff := util.Uint32SliceSymmetricDifferenceSet(rolePermission, params.PermissionIDList); len(diff) == 0 {
		err = svc.dao.Update(role)
		if err != nil {
			return nil, errcode.UpdateFail.WithDetails(err.Error())
		}
	} else {
		operatorPermission, e := svc.userSvc.GetPermIDList(operator.ID)
		if e != nil {
			return nil, errcode.UpdateFail.WithDetails(e.Error())
		}
		var isAdmin = operator.IsAdmin()
		// 目标用户最终的权限
		var modify []uint32
		if !isAdmin {
			canModify := util.Uint32SliceIntersectionSet(operatorPermission, params.PermissionIDList) // 非管理员账户仅可修改自己有的权限
			unModify := util.Uint32SliceDifferenceSet(rolePermission, operatorPermission)             // 被修改用户所拥有的执行者没有的权限不能被修改
			modify = util.Uint32SliceUnionSet(unModify, canModify)
		} else {
			modify = params.PermissionIDList // 管理员账户可以修改全部权限
		}
		newPermList, err := permDao.All(database.QueryWhereGroup{{Prefix: "id", Value: []interface{}{modify}}})
		if err != nil {
			return nil, errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
		}

		role.Permission = newPermList
		err = svc.dao.UpdateAndAssociationReplace(role)
		if err != nil {
			return nil, errcode.UpdateFail.WithDetails(err.Error())
		}

		go svc.delRolePermCache(role.ID)
	}

	// 当修改的是自己的权限时, 需要重新获取权限; 合并处理了 // TODO 去掉这部分
	permissions, e := svc.userSvc.GetPermNameList(operator.ID)
	if e != nil {
		return nil, errcode.UpdateFail.WithDetails(e.Error())
	}

	return permissions, nil
}

// GetPermNameList 优先从缓存获取
func (svc *RoleSvc) GetPermNameList(roleIDList []uint32) (result []string, e *errcode.Error) {
	// 生成redis的缓存key
	roleIDRKeyList := make([]string, 0, len(roleIDList))
	rParams := make([]interface{}, 0, len(roleIDList)+1)
	rParams = append(rParams, "SUNION")
	for _, roleID := range roleIDList {
		roleIDRKey := dao.Role{Model: dao.Model{ID: convert.GetUint32(roleID)}}.RolePermRKey()
		roleIDRKeyList = append(roleIDRKeyList, roleIDRKey)
		rParams = append(rParams, roleIDRKey)
	}

	// 判断缓存是否存在
	existNum, err := svc.rdb.Exists(svc.ctx, roleIDRKeyList...).Result()
	if err != nil {
		svc.logger.Errorf("GetPermNameList rdb.Exists: %v", err)
		return svc.getPermNameList(roleIDList)
	}
	if int(existNum) != len(roleIDList) {
		go svc.cachedRolePerm(roleIDList...)
		return svc.getPermNameList(roleIDList)
	}

	// 根据role进行查询权限
	result, err = svc.rdb.Do(svc.ctx, rParams...).StringSlice()
	if err != nil {
		svc.logger.Errorf("GetPermNameList rdb.Do: %v", err)
		return svc.getPermNameList(roleIDList)
	}
	if len(result) == 0 {
		return svc.getPermNameList(roleIDList)
	}

	return result, nil
}

// 从数据库获取角色的全部权限
func (svc *RoleSvc) getPermNameList(roleIDList []uint32) ([]string, *errcode.Error) {
	permNameList, err := svc.dao.GetPermNameList(database.QueryWhereGroup{
		{Prefix: "r.id", Value: []interface{}{roleIDList}},
	})
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return permNameList, nil
}

// TODO 通过 参数进行更新 而不是在内部查询
// 设置 roleIDList 中角色对应的权限到 redis
func (svc *RoleSvc) cachedRolePerm(roleIDList ...uint32) (addCount int, err error) {
	if err = svc.rdb.Lock(svc.ctx, "SetRolePerm"); err != nil {
		svc.logger.Errorf("cachedRolePerm Lock: %v", err)
		return 0, err
	}
	defer func() {
		if unlockErr := svc.rdb.Unlock(svc.ctx, "SetRolePerm"); unlockErr != nil {
			svc.logger.Errorf("cachedRolePerm Unlock: %v", unlockErr)
			err = fmt.Errorf("cachedRolePerm: %v; unlockErr: %v", err, unlockErr)
		}
	}()

	var queryParams database.QueryWhereGroup
	if len(roleIDList) > 0 {
		queryParams = append(queryParams, database.QueryWhere{Prefix: "id", Value: []interface{}{roleIDList}})
	}
	roles, err := svc.dao.All(queryParams)
	if err != nil {
		return 0, err
	}

	rolePermMap := make(map[string][]interface{})
	rolePermKeys := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDRKey := role.RolePermRKey()
		rolePermKeys = append(rolePermKeys, roleIDRKey)
		for _, permission := range role.Permission {
			rolePermMap[roleIDRKey] = append(rolePermMap[roleIDRKey], permission.Name)
		}
	}

	// 删除缓存
	_, err = svc.rdb.Del(svc.ctx, rolePermKeys...).Result()
	if err != nil {
		svc.logger.Errorf("cachedRolePerm Del %+v: %v", rolePermKeys, err)
		return 0, err
	}
	// 插入缓存
	for rKey, rolePerm := range rolePermMap {
		_, err := svc.rdb.SAdd(svc.ctx, rKey, rolePerm...).Result()
		if err != nil {
			svc.logger.Errorf("cachedRolePerm SAdd %v: %v", rKey, err)
			return 0, err
		}
		addCount++
	}

	return addCount, nil
}

func (svc *RoleSvc) delRolePermCache(roleID uint32) error {
	roleIDRKey := dao.Role{Model: dao.Model{ID: roleID}}.RolePermRKey()
	if _, err := svc.rdb.Del(svc.ctx, roleIDRKey).Result(); err != nil {
		svc.logger.Errorf("DelRolePermCache: %v", err)
		return err
	}
	return nil
}

func (svc *RoleSvc) GetPermIDList(roleID uint32) (permIDList []uint32, e *errcode.Error) {
	permIDList, err := svc.dao.GetPermIDList(database.QueryWhereGroup{
		{Prefix: "r.id", Value: []interface{}{roleID}},
	})
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return permIDList, nil
}

func (svc *RoleSvc) Options() (data []*dao.Role, e *errcode.Error) {
	data, err := svc.dao.Options()
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, nil
}
