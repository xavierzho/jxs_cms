package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"data_backend/internal/dao"
	"data_backend/internal/global"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"gorm.io/gorm"
)

const (
	update_file_path = "data/updatePermissions.json" // 更新文件地址
	copy_file_path   = "data/permissions.json"       // 副本文件地址
)

type jsonPermission struct {
	Permission []*dao.Permission `json:"permission"`
}

func (jp jsonPermission) writeJsonFile() error {
	file, err := os.OpenFile(filepath.Join(global.StoragePath, copy_file_path), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	permissionJson, err := json.Marshal(jp)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	json.Indent(&out, permissionJson, "", "\t")
	_, err = out.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

type PermissionSvc struct {
	logger *logger.Logger
	alarm  message.Alarm
	dao    *dao.PermissionDao
}

func NewPermissionSvc(ctx context.Context, engine *gorm.DB, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *PermissionSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".PermissionSvc")))
	return &PermissionSvc{
		logger: log,
		alarm:  newAlarm(log),
		dao:    dao.NewPermissionDao(engine, log),
	}
}

func (svc *PermissionSvc) All(queryParams []database.QueryWhere) (data []*dao.Permission, e *errcode.Error) {
	data, err := svc.dao.All(queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, nil
}

// 更新权限表
// 增加权限：updatePermissions文件中增加
// 修改权限：数据库中修改(不能修改name)
// 删除权限：数据库删除（updatePermissions文件中不能再包含）
func (svc *PermissionSvc) Refresh() (e *errcode.Error) {
	var dbPermission []*dao.Permission                       // 数据库权限列表
	var dbPermissionMap = make(map[string]*dao.Permission)   // 数据库权限map
	var addPermissionMap = make(map[string]*dao.Permission)  // 更新文件map
	var copyPermissionMap = make(map[string]*dao.Permission) // 副本文件map
	var isEqual = true                                       // 副本文件是否与数据库一致
	var updateFlag bool                                      // 是否更新了数据库
	var err error

	// 读取数库权限
	if dbPermission, err = svc.dao.All(nil); err != nil {
		return errcode.QueryFail.WithDetails(err.Error())
	}
	for _, permission := range dbPermission {
		dbPermissionMap[permission.Name] = permission
	}

	// 新增权限
	if addPermissionMap, _ = svc.parseJson(update_file_path); len(addPermissionMap) > 0 {
		if updateFlag, err = svc.add(dbPermissionMap, addPermissionMap); err != nil {
			return errcode.CreateFail.WithDetails(err.Error())
		}
		if updateFlag { // 更新后重新获取
			if dbPermission, err = svc.dao.All(nil); err != nil {
				return errcode.QueryFail.WithDetails(err.Error())
			}
		}
	}

	// 副本文件中权限
	copyPermissionMap, _ = svc.parseJson(copy_file_path)
	// 两个文件权限数不等 更新
	if len(dbPermissionMap) != len(copyPermissionMap) {
		isEqual = false
	}
	if !updateFlag && isEqual {
		// 对比副本中权限是否与数库不一致
		for name := range dbPermissionMap {
			item := copyPermissionMap[name]
			dbItem := dbPermissionMap[name]
			if *item != *dbItem {
				isEqual = false
				break
			}
		}
	}

	// 将全部数保存到本地
	if (updateFlag || !isEqual) && len(dbPermission) > 0 {
		err = jsonPermission{Permission: dbPermission}.writeJsonFile()
		if err != nil {
			errMsg := fmt.Sprintf("写入本地权限文件失败: %v", err)
			svc.alarm.AlertErrorMsg(errMsg, message.CMS_ID)
			return errcode.ExecuteFail.WithDetails(errMsg)
		}
	}

	return nil
}

func (svc *PermissionSvc) parseJson(path string) (map[string]*dao.Permission, error) {
	bytes, err := os.ReadFile(filepath.Join(global.StoragePath, path))
	if err != nil {
		svc.alarm.AlertErrorMsg(fmt.Sprintf("读取本地权限文件失败: %v", err), message.CMS_ID)
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, nil
	}
	p := &jsonPermission{
		Permission: []*dao.Permission{},
	}
	err = json.Unmarshal(bytes, p)
	if err != nil {
		svc.alarm.AlertErrorMsg(fmt.Sprintf("解析本地权限数据失败: %v", err), message.CMS_ID)
		return nil, err
	}
	permissionMap := make(map[string]*dao.Permission)
	for _, permission := range p.Permission {
		permissionMap[permission.Name] = permission
	}

	return permissionMap, nil
}

func (svc *PermissionSvc) add(dbPermissionMap, addPermissions map[string]*dao.Permission) (bool, error) {
	var createList []*dao.Permission
	for key, item := range addPermissions {
		if dbItem, ok := dbPermissionMap[key]; ok && item != dbItem { // 不处理更新的情况
		} else if !ok {
			createList = append(createList, item)
		}
	}

	if len(createList) > 0 {
		if err := svc.dao.Create(createList); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (svc *PermissionSvc) Options() ([]*dao.Permission, *errcode.Error) {
	data, err := svc.dao.Options()
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, nil
}
