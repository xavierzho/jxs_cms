import Vue from 'vue'
import App from './App'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import store from './store'
import router from './router'
import 'normalize.css/normalize.css' // A modern alternative to CSS resets
import '@/styles/index.scss' // global css
import '@/icons' // icon
import '@/router/guard' // permission control
import i18n from '@/utils/i18n'
import filter from '@/utils/filter'
import { bindingCheckPermission } from '@/utils/directive'
import { checkPermission } from '@/utils/auth'
import echarts from '@/utils/echart'

Vue.prototype.$echarts = echarts

// set ElementUI lang to EN
Vue.use(ElementUI, {
  // locale,
  i18n: (key, value) => i18n.t(key, value),
  size: 'small',
})

filter.forEach(el => {
  Vue.filter(el.name, el.func)
})

// 权限绑定
Vue.directive('perm', bindingCheckPermission)

Vue.prototype.$checkPermission = checkPermission

Vue.config.productionTip = false

new Vue({
  el: '#app',
  i18n,
  router,
  store,
  render: h => h(App),
})
