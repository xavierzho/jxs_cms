import Vue from 'vue'
import Router from 'vue-router'
import Layout from '@/layout'
import store from '@/store'
import BaseView from '@/views/Report/BaseView.vue'

Vue.use(Router)

export const constantRoutes = [
  {
    path: '/redirect',
    name: 'redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index'),
      },
    ],
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true,
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true,
  },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        meta: {
          title: 'Dashboard',
          noCache: true,
          affix: true,
        },
        component: () => import('@/views/Dashboard/index'),
      },
    ],
  },
]

const dynamicRouter = [
  {
    path: "/advertising",
    component: Layout,
    redirect: '/advertising/analysis',
    children: [
      {
        path: 'analysis',
        name: 'advertising_analysis',
        meta: { permission: 'advertising_analysis_view' },
        component: () => import('@/views/Advertising/Analysis.vue'),
      },
    ]
  },
  {
    path: '/report',
    component: Layout,
    redirect: '/report/revenue',
    children: [
      {
        path: 'revenue',
        name: 'report_revenue',
        meta: { permission: 'report_revenue_view' },
        component: () => import('@/views/Report/Revenue.vue'),
      },
      {
        path: 'cohort',
        name: 'report_cohort',
        meta: { permission: 'report_cohort_view' },
        component: () => import('@/views/Report/Cohort.vue'),
      },
      {
        path: 'realtime',
        name: 'report_realtime',
        meta: { permission: 'report_realtime_view' },
        component: () => import('@/views/Report/Realtime.vue'),
      },
      {
        path: 'market',
        name: 'report_market',
        meta: { permission: 'report_market_view' },
        component: () => import('@/views/Report/Market.vue'),
      },
      {
        path: 'bet',
        name: 'report_bet',
        meta: { permission: 'report_bet_view' },
        component: () => import('@/views/Report/Bet.vue'),
      },
      {
        path: 'invite',
        name: 'report_invite',
        meta: { permission: 'report_invite_view' },
        component: () => import('@/views/Report/Invite.vue'),
      },
      {
        path: 'order',
        name: 'report_order',
        meta: { permission: 'report_order_view' },
        component: () => import('@/views/Report/Order.vue'),
      },
      {
        path: 'recall',
        name: 'report_recall',
        meta: { permission: 'report_recall_view' },
        component: () => import('@/views/Report/Recall.vue'),
      },
    ],
  },
  {
    path: '/inquire',
    component: Layout,
    redirect: '/inquire/item',
    children: [
      {
        path: 'item',
        name: 'inquire_item',
        meta: { permission: 'inquire_item_view' },
        component: () => import('@/views/Inquire/Item.vue'),
      },
      {
        path: 'revenue-item',
        name: 'inquire_revenue_item',
        meta: { permission: 'inquire_revenue_item_view' },
        component: () => import('@/views/Inquire/RevenueItem.vue'),
      },
      {
        path: 'bet-item',
        name: 'inquire_bet_item',
        meta: { permission: 'inquire_bet_item_view' },
        component: () => import('@/views/Inquire/BetItem.vue'),
      },
      {
        path: 'gacha',
        name: 'inquire_gacha',
        meta: { permission: 'inquire_gacha_view' },
        component: () => import('@/views/Inquire/Gacha.vue'),
      },
      {
        path: 'balance',
        name: 'inquire_balance',
        meta: { permission: 'inquire_balance_view' },
        component: () => import('@/views/Inquire/Balance.vue'),
      },
      {
        path: 'coupon',
        name: 'inquire_coupon',
        meta: { permission: 'inquire_coupon_view' },
        component: () => import('@/views/Inquire/Coupon.vue'),
      },
      {
        path: 'invite',
        name: 'inquire_invite',
        meta: { permission: 'inquire_invite_view' },
        component: () => import('@/views/Inquire/Invite.vue'),
      },
      {
        path: 'recall',
        name: 'inquire_recall',
        meta: { permission: 'inquire_recall_view' },
        component: () => import('@/views/Inquire/Recall.vue'),
      },
      {
        path: 'task',
        name: 'inquire_task',
        meta: { permission: 'inquire_task_view' },
        component: () => import('@/views/Inquire/Task.vue'),
      },
    ]
  },
  {
    path: '/activity',
    component: Layout,
    redirect: '/activity/cost-award',
    children: [
      {
        path: 'cost-award',
        name: 'cost_award',
        meta: { permission: 'activity_cost_award_view' },
        component: () => import('@/views/Activity/CostAward.vue'),
      },
      {
        path: 'cost-award-log',
        name: 'cost_award_log',
        meta: { permission: 'activity_cost_award_log_view' },
        component: () => import('@/views/Activity/CostAwardLog.vue'),
      },
      {
        path: 'turntable',
        name: 'activity_turntable',
        meta: { permission: 'activity_turntable_view' },
        component: () => import('@/views/Activity/Turntable.vue'),
      },
      {
        path: 'step-by-step',
        name: 'activity_step_by_step',
        meta: { permission: 'activity_step_by_step_view' },
        component: () => import('@/views/Activity/StepByStep.vue'),
      },
      {
        path: 'sign-in',
        name: 'activity_sign_in',
        meta: { permission: 'activity_sign_in_view' },
        component: () => import('@/views/Activity/SignIn.vue'),
      },
      {
        path: 'team-pk',
        name: 'activity_team_pk',
        meta: { permission: 'activity_team_pk_view' },
        component: () => import('@/views/Activity/TeamPK.vue'),
      },
      {
        path: 'redemption-code',
        name: 'activity_redemption_code',
        meta: { permission: 'activity_redemption_code_view' },
        component: () => import('@/views/Activity/RedemptionCode.vue'),
      },
    ]
  },
  {
    path: '/management',
    component: Layout,
    redirect: '/management/user',
    children: [
      {
        path: 'user',
        name: 'management_user',
        meta: { permission: 'management_user_view' },
        component: () => import('@/views/User/UserList'),
      },
      {
        path: 'role',
        name: 'management_role',
        meta: { permission: 'management_role_view' },
        component: () => import('@/views/User/RoleList'),
      },
    ],
  },
]

const filterRouter = (routerArr) => {
  let arr = []
  if (!store.state.user.permission) {
    return []
  }
  for (let i = 0; i < routerArr.length; i++) {
    // 判断是否存在权限
    if (routerArr[i].meta && routerArr[i].meta.permission && !store.state.user.permission.includes(routerArr[i].meta.permission)) {
      continue
    }
    let currentMenu = {
      name: routerArr[i].name,
      path: routerArr[i].path,
      component: routerArr[i].component,
      meta: routerArr[i].meta,
    }
    if (routerArr[i].hasOwnProperty('redirect')) {
      currentMenu.redirect = routerArr[i].redirect
    }
    if (routerArr[i].hasOwnProperty('children') && routerArr[i].children.length > 0) {
      let children = filterRouter(routerArr[i].children)
      if (children.length === 0) {
        continue
      } else {
        currentMenu.children = children
      }
    }
    arr.push(currentMenu)
  }
  return arr
}

const createRouter = (routerArr = []) => {
  let routes = []
  routes.push(...constantRoutes)
  if (routerArr.length > 0) {
    routes.push(...routerArr)
  }

  // // 404 page must be placed at the end !!!
  routes.push({ path: '*', redirect: '/404', hidden: true })
  return new Router({
    scrollBehavior: () => ({ y: 0 }),
    mode: 'history',
    routes: routes,
  })
}

const router = createRouter(dynamicRouter)

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter(filterRouter(dynamicRouter))
  router.matcher = newRouter.matcher // reset router
}

export function resetRouterPromise() {
  return new Promise(resolve => {
    const newRouter = createRouter(filterRouter(dynamicRouter))
    router.matcher = newRouter.matcher // reset router
    resolve()
  })
}

export default router
