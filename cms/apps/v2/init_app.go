package v2

import (
	"fmt"

	costAward "data_backend/apps/v2/internal/activity/cost_award"
	redemptionCode "data_backend/apps/v2/internal/activity/redemption_code"
	signIn "data_backend/apps/v2/internal/activity/sign_in"
	stepByStep "data_backend/apps/v2/internal/activity/step_by_step"
	teamPK "data_backend/apps/v2/internal/activity/team_pk"
	"data_backend/apps/v2/internal/activity/turntable"
	"data_backend/apps/v2/internal/admin"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/balance"
	"data_backend/apps/v2/internal/inquire/coupon"
	"data_backend/apps/v2/internal/inquire/gacha"
	iInvite "data_backend/apps/v2/internal/inquire/invite"
	"data_backend/apps/v2/internal/inquire/item"
	iRecall "data_backend/apps/v2/internal/inquire/recall"
	iTask "data_backend/apps/v2/internal/inquire/task"
	"data_backend/apps/v2/internal/marketing"
	"data_backend/apps/v2/internal/report/bet"
	"data_backend/apps/v2/internal/report/cohort"
	"data_backend/apps/v2/internal/report/dashboard"
	reportInvite "data_backend/apps/v2/internal/report/invite"
	"data_backend/apps/v2/internal/report/market"
	"data_backend/apps/v2/internal/report/order"
	"data_backend/apps/v2/internal/report/realtime"
	"data_backend/apps/v2/internal/report/recall"
	"data_backend/apps/v2/internal/report/revenue"
	iDao "data_backend/internal/dao"
	"data_backend/internal/global"
	iService "data_backend/internal/service"
	"data_backend/pkg/queue"
	"data_backend/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitSetting(config *setting.Config) (err error) {
	err = config.ReadSection("Database", &local.DatabaseSetting)
	if err != nil {
		return fmt.Errorf("ReadSection Database: %w", err)
	}
	err = config.ReadSection("Redis", &local.RedisSetting)
	if err != nil {
		return fmt.Errorf("ReadSection Redis: %w", err)
	}

	return nil
}

func InitObject() (err error) {
	// 日志
	if err = local.SetupLogger(); err != nil {
		return fmt.Errorf("SetupLogger: %w", err)
	}
	if err = local.SetupAlarm(); err != nil {
		return fmt.Errorf("SetupAlarm: %w", err)
	}
	if err = local.SetupDBEngine(); err != nil {
		return fmt.Errorf("SetupDBEngine: %w", err)
	}
	if global.ServerSetting.RunMode != global.RUN_MODE_MIGRATE {
		if err = local.SetupRedis(); err != nil {
			return fmt.Errorf("SetupRedis: %w", err)
		}
	}

	// 迁移模式仅初始化日志 数据库
	if global.ServerSetting.RunMode == global.RUN_MODE_MIGRATE {
		return
	}

	// 刷新权限表
	permSvc := iService.NewPermissionSvc(local.Ctx, local.CMSDB, local.Logger, local.NewAlarm)
	e := permSvc.Refresh()
	if e != nil {
		return e
	}

	if err = local.SetupMiddlewareObject(); err != nil {
		return fmt.Errorf("SetupMiddlewareObject: %w", err)
	}

	// debug 模式不继续
	if global.ServerSetting.RunMode == global.RUN_MODE_DEBUG {
		return nil
	}

	// 定时任务
	if err = local.SetupJobs(); err != nil {
		return fmt.Errorf("SetupJobs: %w", err)
	}
	if err = startJobs(); err != nil {
		return fmt.Errorf("startJobs: %w", err)
	}

	// 队列
	if err = local.SetupQueue(); err != nil {
		return fmt.Errorf("SetupQueue: %w", err)
	}
	if err = startQueue(); err != nil {
		return fmt.Errorf("startQueue: %w", err)
	}

	return nil
}

// 执行任务
func startJobs() (err error) {
	// admin
	{
		if err = admin.AddJobList(); err != nil {
			return err
		}
	}

	// report
	{
		if err = revenue.AddJobList(); err != nil {
			return err
		}

		if err = cohort.AddJobList(); err != nil {
			return err
		}

		if err = bet.AddJobList(); err != nil {
			return err
		}

		if err = market.AddJobList(); err != nil {
			return err
		}

		if err = realtime.AddJobList(); err != nil {
			return err
		}

		if err = dashboard.AddJobList(); err != nil {
			return err
		}

		if err = reportInvite.AddJobList(); err != nil {
			return err
		}
		if err = order.AddJobList(); err != nil {
			return err
		}
		if err = recall.AddJobList(); err != nil {
			return err
		}
		if err = marketing.AddJobList(); err != nil {
			return err
		}
	}

	// inquire
	{
		if err = item.AddJobList(); err != nil {
			return err
		}

		if err = gacha.AddJobList(); err != nil {
			return err
		}

		if err = balance.AddJobList(); err != nil {
			return err
		}

		if err = coupon.AddJobList(); err != nil {
			return err
		}

		if err = iInvite.AddJobList(); err != nil { //inquire/invite
			return err
		}
		if err = iRecall.AddJobList(); err != nil { //inquire/recall
			return err
		}
		if err = iTask.AddJobList(); err != nil { //inquire/recall
			return err
		}

	}

	// activity
	{
		if err = costAward.AddJobList(); err != nil {
			return err
		}
		if err = turntable.AddJobList(); err != nil { //activity/turntable
			return err
		}
		if err = stepByStep.AddJobList(); err != nil {
			return err
		}
		if err = signIn.AddJobList(); err != nil {
			return err
		}
		if err = teamPK.AddJobList(); err != nil {
			return err
		}
		if err = redemptionCode.AddJobList(); err != nil {
			return err
		}
	}

	if err = local.JobWorker.StartJob(local.CronChain, local.QueueCronChain); err != nil {
		return err
	}
	if err = local.JobWorker.StartFrequentlyJob(local.QueueCronChain); err != nil {
		return err
	}

	return nil
}

// 加入队列
func startQueue() (err error) {
	// 将队列job加入进去
	err = local.QueueWorker.AddQueueJob([]*queue.QueueJob{
		{Name: local.JobWorker.QueueRKey(), Retry: false, Run: local.JobWorker.QueueJobRun},
	})
	if err != nil {
		return err
	}

	// admin
	{
		if err = admin.AddQueueJob(); err != nil {
			return err
		}
	}

	// report
	{
		if err = revenue.AddQueueJob(); err != nil {
			return err
		}

		if err = cohort.AddQueueJob(); err != nil {
			return err
		}

		if err = bet.AddQueueJob(); err != nil {
			return err
		}

		if err = market.AddQueueJob(); err != nil {
			return err
		}

		if err = realtime.AddQueueJob(); err != nil {
			return err
		}

		if err = dashboard.AddQueueJob(); err != nil {
			return err
		}

		if err = reportInvite.AddQueueJob(); err != nil {
			return err
		}
		if err = order.AddQueueJob(); err != nil {
			return err
		}
		if err = recall.AddJobList(); err != nil {
			return err
		}
		if err = marketing.AddQueueJob(); err != nil {
			return err
		}

	}

	// inquire
	{
		if err = item.AddQueueJob(); err != nil {
			return err
		}

		if err = gacha.AddQueueJob(); err != nil {
			return err
		}

		if err = balance.AddQueueJob(); err != nil {
			return err
		}

		if err = coupon.AddQueueJob(); err != nil {
			return err
		}
		if err = iInvite.AddQueueJob(); err != nil { //inquire/invite
			return err
		}
		if err = iRecall.AddQueueJob(); err != nil { //inquire/recall
			return err
		}
		if err = iTask.AddQueueJob(); err != nil { //inquire/recall
			return err
		}
	}

	// activity
	{
		if err = costAward.AddQueueJob(); err != nil {
			return err
		}
		if err = turntable.AddQueueJob(); err != nil {
			return err
		}
		if err = stepByStep.AddQueueJob(); err != nil {
			return err
		}
		if err = signIn.AddQueueJob(); err != nil {
			return err
		}
		if err = teamPK.AddQueueJob(); err != nil {
			return err
		}
		if err = redemptionCode.AddQueueJob(); err != nil {
			return err
		}
	}

	local.QueueWorker.Start()
	return nil
}

func InitRouter(rg *gin.RouterGroup) (err error) {
	// admin
	{
		rg := rg.Group("")
		if err = admin.InitRouter(rg); err != nil {
			return fmt.Errorf("admin.InitRouter: %v", err)
		}
	}

	// report
	{
		rg := rg.Group("report")
		rg.Use(local.JWT.JWT())
		if err = revenue.InitRouter(rg); err != nil {
			return fmt.Errorf("revenue.InitRouter: %v", err)
		}

		if err = cohort.InitRouter(rg); err != nil {
			return fmt.Errorf("cohort.InitRouter: %v", err)
		}

		if err = bet.InitRouter(rg); err != nil {
			return fmt.Errorf("bet.InitRouter: %v", err)
		}

		if err = market.InitRouter(rg); err != nil {
			return fmt.Errorf("market.InitRouter: %v", err)
		}

		if err = realtime.InitRouter(rg); err != nil {
			return fmt.Errorf("realtime.InitRouter: %v", err)
		}

		if err = dashboard.InitRouter(rg); err != nil {
			return fmt.Errorf("dashboard.InitRouter: %v", err)
		}

		if err = reportInvite.InitRouter(rg); err != nil {
			return fmt.Errorf("invite.InitRouter: %v", err) //report/invite
		}

		if err = order.InitRouter(rg); err != nil {
			return fmt.Errorf("order.InitRouter: %v", err) //report/order
		}

		if err = recall.InitRouter(rg); err != nil {
			return fmt.Errorf("recall.InitRouter: %v", err) //report/recall
		}

	}

	// inquire
	{
		rg := rg.Group("inquire")
		rg.Use(local.JWT.JWT())
		if err = item.InitRouter(rg); err != nil {
			return fmt.Errorf("item.InitRouter: %v", err)
		}

		if err = gacha.InitRouter(rg); err != nil {
			return fmt.Errorf("gacha.InitRouter: %v", err)
		}

		if err = balance.InitRouter(rg); err != nil {
			return fmt.Errorf("balance.InitRouter: %v", err)
		}

		if err = coupon.InitRouter(rg); err != nil {
			return fmt.Errorf("coupon.InitRouter: %v", err)
		}
		if err = iInvite.InitRouter(rg); err != nil {
			return fmt.Errorf("invite.InitRouter: %v", err) //inquire/invite
		}
		if err = iRecall.InitRouter(rg); err != nil { //inquire/recall
			return fmt.Errorf("recall.InitRouter: %v", err)
		}
		if err = iTask.InitRouter(rg); err != nil { //inquire/recall
			return fmt.Errorf("task.InitRouter: %v", err)
		}
	}

	// marketing
	{
		// marketing has its own group "marketing"
		if err = marketing.InitRouter(rg); err != nil {
			return fmt.Errorf("marketing.InitRouter: %v", err)
		}
	}

	// activity
	{
		rg := rg.Group("activity")
		rg.Use(local.JWT.JWT())
		if err = costAward.InitRouter(rg); err != nil {
			return fmt.Errorf("costAward.InitRouter: %v", err)
		}
		if err = turntable.InitRouter(rg); err != nil {
			return fmt.Errorf("turntable.InitRouter: %v", err) //activity/turntable
		}
		if err = stepByStep.InitRouter(rg); err != nil {
			return fmt.Errorf("stepByStep.InitRouter: %v", err)
		}
		if err = signIn.InitRouter(rg); err != nil {
			return fmt.Errorf("signIn.InitRouter: %v", err)
		}
		if err = teamPK.InitRouter(rg); err != nil {
			return fmt.Errorf("teamPK.InitRouter: %v", err)
		}
		if err = redemptionCode.InitRouter(rg); err != nil {
			return fmt.Errorf("redemptionCode.InitRouter: %v", err)
		}
	}

	return nil
}

// 迁移app
func MigrateModel() (err error) {
	// migrate common
	if err = migrateModel(); err != nil {
		return err
	}

	// admin
	admin.AppendMigrateModel()

	// report
	{
		revenue.AppendMigrateModel()
		cohort.AppendMigrateModel()
		bet.AppendMigrateModel()
		market.AppendMigrateModel()
		realtime.AppendMigrateModel()
		dashboard.AppendMigrateModel()
		reportInvite.AppendMigrateModel()
		order.AppendMigrateModel()
		recall.AppendMigrateModel()
	}

	// inquire
	{
		item.AppendMigrateModel()
		gacha.AppendMigrateModel()
		balance.AppendMigrateModel()
		coupon.AppendMigrateModel()
		iInvite.AppendMigrateModel()
		iRecall.AppendMigrateModel()
		iTask.AppendMigrateModel()
	}

	// marketing
	{
		marketing.AppendMigrateModel()
	}

	// activity
	{
		costAward.AppendMigrateModel()
		turntable.AppendMigrateModel()
		stepByStep.AppendMigrateModel()
		signIn.AppendMigrateModel()
		teamPK.AppendMigrateModel()
		redemptionCode.AppendMigrateModel()
	}

	_ = local.MigrateModel()

	return nil
}

func migrateModel() (err error) {
	modelArr := []interface{}{
		&iDao.Permission{}, &iDao.Role{}, &iDao.User{},
	}

	// 第一次启动 进行初始化
	if !local.CMSDB.Migrator().HasTable((&iDao.User{}).TableName()) {
		// 先迁移其他表后再进行初始化
		err = local.CMSDB.AutoMigrate(&iDao.Permission{})
		if err != nil {
			return err
		}

		// 刷新权限表
		permSvc := iService.NewPermissionSvc(local.Ctx, local.CMSDB, local.Logger, local.NewAlarm)
		e := permSvc.Refresh()
		if e != nil {
			return e
		}

		err = iDao.InitFirstUser(local.CMSDB, local.Logger)
		if err != nil {
			return err
		}
	}

	err = local.CMSDB.AutoMigrate(modelArr...)
	if err != nil {
		return err
	}

	return nil
}
