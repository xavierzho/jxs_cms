import axios from 'axios'
import { Message, Notification } from 'element-ui'
import store from '@/store'
import router from '@/router'
import qs from 'qs'
import i18n from '@/utils/i18n'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 10000, // request timeout,
  transformRequest: [function (data, headers) {
    // 上传格式则不进行序列化 序列化会去掉分隔符
    if (headers['Content-Type'] === 'multipart/form-data') {
      return data
    }
    // 对 data 进行任意转换处理
    return qs.stringify(data, { arrayFormat: 'brackets' })
  }],
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent
    if (store.getters.lang) {
      config.headers['locale'] = store.getters.lang
    }
    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['X-Token'] = store.getters.token
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  },
)

// response interceptor
service.interceptors.response.use(
  response => {
    switch (response.config.responseType){
      case 'arraybuffer': // 用于获取文件名
        return response
      default:
        return response.data
    }
  },
  async error => {
    console.log('err ' + error) // for debug
    if (error.response) {
      switch (error.response.data.code) {
        case 40100000:
        case 40100001:
        case 40100002:
        case 40100003:
          Message({
            message: i18n.t('common.LoginExpired'),
            type: 'error',
            duration: 5 * 1000,
          })
          await store.dispatch('user/logout')
          router.push('/login')
          break
        case 40300002:
          Message({
            message: i18n.t('common.PermissionDenied'),
            type: 'error',
            duration: 5 * 1000,
          })
          break
        default:
          Message({
            message: error.response.data.msg,
            type: 'error',
            duration: 5 * 1000,
          })
          if (error.response.data && error.response.data.details.length > 0) {
            error.response.data.details.forEach((detail, index) => {
              setTimeout(() => {
                Notification({
                  title: 'Warning',
                  message: detail,
                  type: 'warning',
                  duration: 3 * 1000,
                })
              }, index * 1000)
            })
          }
      }
    }

    return Promise.reject(error)
  },
)

export default service
