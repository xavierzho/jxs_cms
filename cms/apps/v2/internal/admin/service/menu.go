package service

import "data_backend/internal/dao"

/*
树节点:
	Title: 大驼峰 // 简短
	Name: 标题 小蛇形格式 // 前缀标题+标题 小蛇形格式(用于权限中标识pages) //详细
	Path: 标题 中划线格式 // 简短
	Permission: 前缀+标题+后缀 小蛇形格式 // 详细

菜单结构尽可能与代码结构统一
route中的路径每一段尽可能简化
*/
func MenuList() []*dao.Menu {
	menuList := []*dao.Menu{
		{Title: "Dashboard", Name: "dashboard", Path: "/dashboard"},
		{
			Title: "Report",
			Name:  "report",
			Path:  "/report",
			Children: []*dao.Menu{
				{Title: "Cohort", Name: "report_cohort", Path: "cohort", Permission: "report_cohort_view"},
				{Title: "Revenue", Name: "report_revenue", Path: "revenue", Permission: "report_revenue_view"},
				{Title: "Realtime", Name: "report_realtime", Path: "realtime", Permission: "report_realtime_view"},
				{Title: "Market", Name: "report_market", Path: "market", Permission: "report_market_view"},
				{Title: "Bet", Name: "report_bet", Path: "bet", Permission: "report_bet_view"},
				{Title: "Invite", Name: "report_invite", Path: "invite", Permission: "report_invite_view"},
				{Title: "Order", Name: "report_order", Path: "order", Permission: "report_order_view"},
			},
		},
		{
			Title: "Inquire",
			Name:  "inquire",
			Path:  "/inquire",
			Children: []*dao.Menu{
				{Title: "Item", Name: "inquire_item", Path: "item", Permission: "inquire_item_view"},
				{Title: "Gacha", Name: "inquire_gacha", Path: "gacha", Permission: "inquire_gacha_view"},
				{Title: "Balance", Name: "inquire_balance", Path: "balance", Permission: "inquire_balance_view"},
				{Title: "Coupon", Name: "inquire_coupon", Path: "coupon", Permission: "inquire_coupon_view"},
				{Title: "Invite", Name: "inquire_invite", Path: "invite", Permission: "inquire_invite_view"},
			},
		},
		{
			Title: "Activity",
			Name:  "activity",
			Path:  "/activity",
			Children: []*dao.Menu{
				{Title: "CostAward", Name: "cost_award", Path: "cost-award", Permission: "activity_cost_award_view"},
				{Title: "CostAwardLog", Name: "cost_award_log", Path: "cost-award-log", Permission: "activity_cost_award_log_view"},
				{Title: "Turntable", Name: "activity_turntable", Path: "turntable", Permission: "activity_turntable_view"},
			},
		},
		{
			Title: "Management",
			Name:  "management",
			Path:  "/management",
			Children: []*dao.Menu{
				{Title: "User", Name: "management_user", Path: "user", Permission: "management_user_view"},
				{Title: "Role", Name: "management_role", Path: "role", Permission: "management_role_view"},
				{Title: "Log", Name: "management_log", Path: "log", Permission: "management_log_view"},
			},
		},
	}

	return menuList
}
