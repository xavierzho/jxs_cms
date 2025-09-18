import router from './index'
import store from '../store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/cache' // get token from cookie
import getPageTitle from '@/utils/get-page-title'
import i18n from '@/utils/i18n'

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login'] // no redirect whitelist

router.beforeEach(async(to, from, next) => {
  // start progress bar
  NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in
  const hasToken = getToken()

  if (hasToken) {
    // 判断是否已获取菜单
    if (store.state.app.menus.length === 0) {
      await store.dispatch('app/getMenus')
    }

    if (to.path === '/login') {
      // if is logged in, redirect to the home page
      next({ path: '/' })
      NProgress.done()
    } else {
      // 确定该路由是否需要权限限制
      if (to.meta.permission && !store.state.user.permission.includes(to.meta.permission)) {
        Message({
          message: i18n.t('common.PermissionDenied'),
          type: 'warning',
          duration: 5 * 1000,
        })
        next('/dashboard')
      }
      next()
    }
  } else {
    /* has no token */
    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      next()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
