package admin

import (
	"data_backend/apps/v2/internal/admin/api"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/pkg/cronjob"
	"data_backend/pkg/queue"

	"github.com/gin-gonic/gin"
)

func AddJobList() error {
	return local.JobWorker.AddJobList(
		map[string][]cronjob.CronCommonJob{},
	)
}

func AddQueueJob() error {
	return local.QueueWorker.AddQueueJob(
		[]*queue.QueueJob{},
	)
}

func InitRouter(r *gin.RouterGroup) (err error) {
	userApi := api.NewUserApi()
	r.POST("login", userApi.Login)

	r = r.Group("")
	r.Use(local.JWT.JWT())

	// 菜单
	{
		rg := r.Group("menu")
		menuApi := api.NewMenuApi()
		rg.GET("", menuApi.All)
	}

	// management
	{
		rg := r.Group("management")
		// 角色接口
		{
			rg := rg.Group("role")
			roleApi := api.NewRoleApi()
			rg.POST("create", local.PermissionGate.CheckPerm("management_role_create"), roleApi.Create)
			rg.GET("", local.PermissionGate.CheckPerm("management_role_view"), roleApi.List)
			rg.PUT("update/:id", local.PermissionGate.CheckPerm("management_role_update"), roleApi.Update)
			rg.GET("options", roleApi.Options)
		}

		// 用户接口
		{
			rg := rg.Group("user")
			rg.POST("create", local.PermissionGate.CheckPerm("management_user_create"), userApi.Create)
			rg.GET("detail", userApi.Detail)
			rg.GET("", local.PermissionGate.CheckPerm("management_user_view"), userApi.List)
			rg.PUT("update/:id", local.PermissionGate.CheckPerm("management_user_update"), userApi.Update)
			rg.PUT("update-self", userApi.UpdateSelf)
			rg.GET("page-permission", local.PermissionGate.CheckPermOr("management_role_create", "management_role_update"), userApi.PagePermission)
			rg.GET("options", userApi.Options)
		}

		// 权限
		{
			rg := rg.Group("permission")
			permissionApi := api.NewPermissionApi()
			rg.GET("options", permissionApi.Options)
		}
	}
	return nil
}

func AppendMigrateModel() {
	local.MigrateModelList = append(local.MigrateModelList, []any{}...)
}
