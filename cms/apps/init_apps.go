package apps

import (
	v2 "data_backend/apps/v2"
	"data_backend/internal/app"
	"data_backend/internal/global"
	"data_backend/internal/middleware/contextTimeout"
	"data_backend/internal/middleware/cors"
	"data_backend/internal/middleware/recovery"
	"data_backend/internal/middleware/requestInfo"
	"data_backend/internal/middleware/translations"
	"data_backend/pkg/setting"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitSetting(config *setting.Config) (err error) {
	if err = v2.InitSetting(config); err != nil {
		return fmt.Errorf("v2 InitSetting: %w", err)
	}

	return
}

func InitObject() (err error) {
	if err = v2.InitObject(); err != nil {
		return fmt.Errorf("v2 InitObject: %w", err)
	}

	return
}

// 路由
func initRouter() (*gin.Engine, error) {
	gin.SetMode(global.ServerSetting.RunMode.String())
	r := gin.New()
	if global.ServerSetting.RunMode == global.RUN_MODE_DEBUG {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
		r.GET("/debug/vars", app.Expvar)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.Use(requestInfo.RequestInfo())
	} else {
		r.Use(recovery.NewRecovery(global.Logger, global.Alarm).Recovery())
		r.Use(requestInfo.RequestInfo())
	}

	// 装载公共中间件
	r.Use(contextTimeout.ContextTimeout(global.APPSetting.DefaultContextTimeout))
	r.Use(translations.Translations())
	r.Use(cors.Cors())

	// 加载其他路由
	otherRouter(r)

	return r, nil
}

func otherRouter(r *gin.Engine) {
	// ping
	r.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
}

func InitAPPsRouter() (r *gin.Engine, err error) {
	// 初始化路由
	if r, err = initRouter(); err != nil {
		return nil, fmt.Errorf("apps.initRouter: %w", err)
	}
	if err = v2.InitRouter(r.Group("api/v2")); err != nil {
		return nil, fmt.Errorf("v2 InitialRouter: %w", err)
	}

	return r, nil
}

// MigrateModel 迁移app
func MigrateModel() (err error) {
	if err = v2.MigrateModel(); err != nil {
		return fmt.Errorf("v2 MigrateModel: %w", err)
	}

	return nil
}
