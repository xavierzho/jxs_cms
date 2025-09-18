package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

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
	"data_backend/pkg/token"
	"data_backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const URR_KEY_EXPIRED_DAY = time.Hour * 24 * 3 // UserRoleRKey的redis过期时间 // TODO 感觉没必要

type UserSvc struct {
	ctx      *gin.Context
	engine   *gorm.DB
	rdb      *redisdb.RedisClient
	logger   *logger.Logger
	roleSvc  *RoleSvc
	logSvc   *OperationLogSvc
	dao      *dao.UserDao
	newAlarm func(log *logger.Logger) message.Alarm
}

func NewUserSvc(ctx *gin.Context, engine *gorm.DB, rdb *redisdb.RedisClient, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *UserSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".UserSvc")))
	return &UserSvc{
		ctx:      ctx,
		engine:   engine,
		rdb:      rdb,
		logger:   log,
		dao:      dao.NewUserDao(engine, log),
		newAlarm: newAlarm,
	}
}

func (svc *UserSvc) Login(params *form.LoginRequest) (*form.UserLoginInfo, *errcode.Error) {
	// 查找用户
	user, err := svc.dao.First([]database.QueryWhere{
		{Prefix: "user_name = ?", Value: []interface{}{params.UserName}},
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, iErrcode.UserNotExist
	} else if err != nil {
		return nil, iErrcode.LoginFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}

	// 校验
	if !util.ComparePassword(user.Password, params.Password) {
		return nil, iErrcode.IncorrectPassword
	}
	if user.IsLock == 1 {
		return nil, iErrcode.UserIsLock
	}

	// 生成Token
	tokenStr, err := token.GenerateToken(user.UserName, user.ID)
	if err != nil {
		svc.logger.Errorf("Login, token.GenerateToken: %v", err)
	}
	_, err = svc.rdb.Set(svc.ctx, token.GetRKeyByUserID(user.ID), tokenStr, token.JWTSetting.Expire).Result()
	if err != nil {
		svc.logger.Errorf("Login, rdb.Set: %v", err)
	}

	// isAdmin
	isAdmin := user.IsAdmin()

	// updateInfo
	lastLogonTime := time.Now()
	user.LastLogonTime = &lastLogonTime
	err = svc.dao.Update(user)
	if err != nil {
		return nil, iErrcode.LoginFail.WithDetails(errcode.UpdateFail.WithDetails(err.Error()).Error())
	}

	permList, e := svc.GetPermNameList(user.ID)
	if e != nil {
		return nil, iErrcode.LoginFail.WithDetails(e.Error())
	}

	return &form.UserLoginInfo{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Token:      tokenStr,
		Permission: permList,
		IsAdmin:    isAdmin,
	}, nil
}

func (svc *UserSvc) Create(params *form.UserCreateRequest) *errcode.Error {
	svc.roleSvc = NewRoleSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)
	svc.logSvc = NewOperationLogSvc(svc.ctx, svc.engine, svc.logger, svc.newAlarm)

	// 查找用户
	flag, err := svc.checkExist(0, params.UserName, params.Email)
	if err != nil {
		return errcode.CreateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}
	if flag {
		return iErrcode.UserExist.WithDetails("duplicate name or email address")
	}

	var isLock uint8 = 0
	if params.IsLock {
		isLock = 1
	}
	roleList, e := svc.roleSvc.All(database.QueryWhereGroup{{Prefix: "id in ?", Value: []interface{}{params.RoleIDList}}})
	if e != nil {
		return errcode.CreateFail.WithDetails(e.Error())
	}
	user := &dao.User{
		UserName: params.UserName,
		Name:     params.Name,
		Email:    params.Email,
		IsLock:   isLock,
		Password: util.GeneratePassword(params.Password),
		Role:     roleList,
	}
	err = svc.dao.Create(user)
	if err != nil {
		return errcode.CreateFail.WithDetails(err.Error())
	}

	go svc.logSvc.Create(&dao.OperationLog{ModuleName: user.TableName(), Operation: "Create", ModuleID: convert.GetString(user.ID)})

	return nil
}

func (svc *UserSvc) checkExist(id uint32, userName, email string) (bool, error) {
	_, err := svc.dao.First(database.QueryWhereGroup{
		{Prefix: "id <> ? and (user_name = ? or email = ?)", Value: []interface{}{id, userName, email}},
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// 获取当前登录用户的信息
// * 注意分别处理返回的不同异常
func (svc *UserSvc) CurrentUser() (user *dao.User, e *errcode.Error) {
	if userID, ok := svc.ctx.Get(app.USER_ID_KEY); !ok {
		return nil, errcode.UnauthorizedTokenError.WithDetails("userID doesn't exist in the ctx")
	} else {
		user, err := svc.dao.First([]database.QueryWhere{
			{Prefix: "id = ?", Value: []interface{}{userID.(string)}},
		})
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}

		return user, nil
	}
}

func (svc *UserSvc) List(params *form.UserListRequest) ([]form.User, int64, *errcode.Error) {
	queryParams := params.Parse()
	users, total, err := svc.dao.List(queryParams, app.GetPager(svc.ctx))
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err := form.UserFormat(users, app.GetOrderBy(svc.ctx))
	if err != nil {
		svc.logger.Errorf("List, UserFormat: %v", err)
		return nil, 0, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, total, nil
}

func (svc *UserSvc) Update(id uint32, params *form.UserUpdateRequest) (e *errcode.Error) {
	if err := params.Valid(); err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}
	svc.roleSvc = NewRoleSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)
	svc.logSvc = NewOperationLogSvc(svc.ctx, svc.engine, svc.logger, svc.newAlarm)

	// 查找用户
	user, err := svc.dao.First([]database.QueryWhere{
		{Prefix: "id = ?", Value: []interface{}{id}},
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return iErrcode.UserNotExist
	} else if err != nil {
		return errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}

	flag, err := svc.checkExist(id, "", params.Email)
	if err != nil {
		return errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}
	if flag {
		return iErrcode.UserExist.WithDetails("duplicate email address")
	}

	var isLock uint8 = 0
	if params.IsLock {
		isLock = 1
	}
	roleList, e := svc.roleSvc.All(database.QueryWhereGroup{{Prefix: "id in ?", Value: []interface{}{params.RoleIDList}}})
	if e != nil {
		return errcode.UpdateFail.WithDetails(e.Error())
	}
	user.Name = params.Name
	user.Email = params.Email
	user.IsLock = isLock
	user.Role = roleList
	// 若更改了密码 则重新生成密码
	if params.Password != "" {
		if !util.ComparePassword(user.Password, params.Password) {
			user.Password = util.GeneratePassword(params.Password)
		}
	}

	err = svc.dao.UpdateAndAssociationReplace(user)
	if err != nil {
		return errcode.UpdateFail.WithDetails(err.Error())
	}

	go svc.logSvc.Create(&dao.OperationLog{ModuleName: user.TableName(), Operation: "Update", ModuleID: convert.GetString(user.ID)})
	go svc.delUserRoleCache(user.ID)

	return nil
}

func (svc *UserSvc) UpdateSelf(params *form.UserUpdateSelfRequest) (user *dao.User, e *errcode.Error) {
	if err := params.Valid(); err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}
	svc.logSvc = NewOperationLogSvc(svc.ctx, svc.engine, svc.logger, svc.newAlarm)

	user, e = svc.CurrentUser()
	if e.Is(errcode.UnauthorizedTokenError) {
		return nil, e
	} else if e != nil {
		return nil, errcode.UpdateFail.WithDetails(e.Error())
	}

	flag, err := svc.checkExist(user.ID, "", params.Email)
	if err != nil {
		return nil, errcode.UpdateFail.WithDetails(errcode.QueryFail.WithDetails(err.Error()).Error())
	}
	if flag {
		return nil, iErrcode.UserExist.WithDetails("duplicate email address")
	}

	user.Name = params.Name
	user.Email = params.Email
	user.Password = util.GeneratePassword(params.NewPassword)
	err = svc.dao.Update(user)
	if err != nil {
		return nil, errcode.UpdateFail.WithDetails(err.Error())
	}

	go svc.logSvc.Create(&dao.OperationLog{ModuleName: user.TableName(), Operation: "UpdateSelf", ModuleID: convert.GetString(user.ID)})

	return user, nil
}

// 优先从缓存获取
func (svc *UserSvc) getRoleIDList(userID uint32) (roleIDList []uint32, err error) {
	userRoleRKey := dao.User{Model: dao.Model{ID: userID}}.UserRoleRKey()
	roleIDListStr, err := svc.rdb.Get(svc.ctx, userRoleRKey).Result()
	if err != nil {
		if err != redis.Nil {
			svc.logger.Errorf("GetRoleByID rdb.Get: %v", err)
		}
		roleIDList, err := svc.dao.GetRoleIDList(database.QueryWhereGroup{{Prefix: "u.id = ?", Value: []interface{}{userID}}})
		if err != nil {
			return nil, err
		}

		go svc.cachedUserRole(userRoleRKey, roleIDList)

		return roleIDList, nil
	}

	roleIDStrList := strings.Split(roleIDListStr, ",")
	for _, item := range roleIDStrList {
		roleIDList = append(roleIDList, convert.GetUint32(item))
	}

	return roleIDList, nil
}

func (svc *UserSvc) cachedUserRole(userRoleRKey string, roleIDList []uint32) {
	var roleIDStrList = make([]string, 0, len(roleIDList))
	for _, item := range roleIDList {
		roleIDStrList = append(roleIDStrList, strconv.FormatInt(int64(item), 10))
	}

	_, err := svc.rdb.Set(svc.ctx, userRoleRKey, strings.Join(roleIDStrList, ","), URR_KEY_EXPIRED_DAY).Result()
	if err != nil {
		svc.logger.Errorf("GetRoleByID rdb.Set: %v", err)
	}
}

func (svc *UserSvc) delUserRoleCache(userID uint32) (err error) {
	_, err = svc.rdb.Del(svc.ctx, dao.User{Model: dao.Model{ID: userID}}.UserRoleRKey()).Result()
	if err != nil {
		svc.logger.Errorf("DelRoleCache: %v", err)
		return err
	}
	return nil
}

func (svc *UserSvc) GetPermNameList(userID uint32) (permList []string, e *errcode.Error) {
	svc.roleSvc = NewRoleSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)

	// 获取所有的角色
	roleIDList, err := svc.getRoleIDList(userID)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	// 获取角色权限的并集
	permList, e = svc.roleSvc.GetPermNameList(roleIDList)
	if e != nil {
		return nil, e
	}

	return permList, nil
}

func (svc *UserSvc) GetPermIDList(userID uint32) ([]uint32, *errcode.Error) {
	permIDList, err := svc.dao.GetPermIDList(database.QueryWhereGroup{
		{Prefix: "u.id = ?", Value: []interface{}{userID}},
	})
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return permIDList, nil
}

// 中间件调用， 需要修改为运行时ctx
// 不可用结构体指针 避免竞态
func (svc UserSvc) CheckPerm(ctx *gin.Context, userID uint32, permission []string) (bool, error) {
	if svc.ctx != ctx {
		svc.ctx = ctx
		svc.logger.WithContext(ctx)
	}

	// 获取该用户所有的权限
	allPermission, e := svc.GetPermNameList(userID)
	if e != nil {
		return false, e
	}

	allow := util.PermissionCheckAll(permission, allPermission)
	return allow, nil
}

// 不可用结构体指针 避免竞态
func (svc UserSvc) CheckPermOr(ctx *gin.Context, userID uint32, permission []string) (bool, error) {
	if svc.ctx != ctx {
		svc.ctx = ctx
		svc.logger.WithContext(ctx)
	}

	// 获取该用户所有的权限
	allPermission, e := svc.GetPermNameList(userID)
	if e != nil {
		return false, e
	}

	allow := util.PermissionCheckOr(permission, allPermission)
	return allow, nil
}

func (svc *UserSvc) CanShowPhoneNum() bool {
	if userID := svc.ctx.Value(app.USER_ID_KEY); userID != nil {
		hasPerm, err := svc.CheckPerm(svc.ctx, convert.GetUint32(userID), []string{dao.PERMISSION_SHOW_SENSITIVE_INFO})
		if err == nil && hasPerm {
			return true
		}
	}

	return false
}

func (svc *UserSvc) PagePermission() (result map[string][]map[string]interface{}, e *errcode.Error) {
	permDao := dao.NewPermissionDao(svc.engine, svc.logger)

	user, e := svc.CurrentUser()
	if e != nil {
		return nil, e
	}
	var isAdmin = user.IsAdmin()

	// 当前用户所有权限
	userPermissions, e := svc.GetPermNameList(user.ID)
	if e != nil {
		return nil, e
	}
	permissions, err := permDao.Options()
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	result = form.PagePermissionFormat(permissions, userPermissions, isAdmin)

	return result, nil
}

func (svc *UserSvc) Options() ([]map[string]interface{}, *errcode.Error) {
	data, err := svc.dao.Options()
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, nil
}
