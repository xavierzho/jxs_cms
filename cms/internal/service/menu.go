package service

import (
	"context"
	"strings"

	"data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/i18n"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"
	"data_backend/pkg/redisdb"
	"data_backend/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuSvc struct {
	ctx      *gin.Context
	engine   *gorm.DB
	rdb      *redisdb.RedisClient
	logger   *logger.Logger
	userSvc  *UserSvc       // 使用时创建
	permSvc  *PermissionSvc // 使用时创建
	newAlarm func(log *logger.Logger) message.Alarm
}

func NewMenuSvc(ctx *gin.Context, engine *gorm.DB, rdb *redisdb.RedisClient, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *MenuSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".MenuSvc")))
	return &MenuSvc{
		ctx:      ctx,
		engine:   engine,
		rdb:      rdb,
		logger:   log,
		newAlarm: newAlarm,
	}
}

func (svc *MenuSvc) All(menus []*dao.Menu) ([]*dao.Menu, *errcode.Error) {
	svc.userSvc = NewUserSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)
	svc.permSvc = NewPermissionSvc(svc.ctx, svc.engine, svc.logger, svc.newAlarm)

	// 获取当前用户及权限
	user, e := svc.userSvc.CurrentUser()
	if e.Is(errcode.UnauthorizedTokenError) {
		return nil, e
	} else if e != nil {
		return nil, iErrcode.InitMenuFail.WithDetails(e.Error())
	}

	permissions, e := svc.userSvc.GetPermNameList(user.ID)
	if e != nil {
		return nil, iErrcode.InitMenuFail.WithDetails(e.Error())
	}

	// 若包含管理员角色 则默认全部菜单(为了添加新权限)
	var isAdmin = user.IsAdmin()

	return svc.filter(menus, permissions, isAdmin, []string{}), nil
}

// * 数据整体用于权限树显示, show 表示是否在侧边栏显示该级菜单
func (svc *MenuSvc) filter(menus []*dao.Menu, allPermission []string, isAdmin bool, titleList []string) []*dao.Menu {
	newMenu := make([]*dao.Menu, 0, len(menus))
	for _, menu := range menus {
		allow := util.PermissionCheckAll([]string{menu.Permission}, allPermission)
		if allow || isAdmin { // * admin 需要能看到全部菜单
			currentMenu := &dao.Menu{
				Name:       menu.Name,
				Title:      global.I18n.T(svc.ctx.Request.Context(), "menu", strings.Join(append(titleList, menu.Title), i18n.NESTED_SEPARATOR)),
				Permission: menu.Permission,
				Path:       menu.Path,
				Show:       allow,
			}
			// 若该菜单有子页面 则该菜单为菜单节点
			if len(menu.Children) > 0 {
				childMenu := svc.filter(menu.Children, allPermission, isAdmin, append(titleList, menu.Title))
				// 菜单节点若没有子菜单则不显示菜单节点
				if len(childMenu) == 0 {
					continue
				}
				currentMenu.Children = childMenu
				if isAdmin {
					currentMenu.Show = false
					// 当角色为admin时需要遍历子菜单判断是否显示该一级菜单
					for _, item := range childMenu {
						if item.Show {
							currentMenu.Show = true
							break
						}
					}
				}
			}
			newMenu = append(newMenu, currentMenu)
		}
	}
	return newMenu
}
