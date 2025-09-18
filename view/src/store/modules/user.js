import {
  getPagePermission,
  setPagePermission,
  getPermission,
  getToken,
  getUserInfo,
  removeAllCookies,
  removeMenu,
  removePagePermission,
  removePermission,
  removeToken,
  setPermission,
  setToken,
  setUserInfo,
} from "@/utils/cache"
import { resetRouter } from '@/router'
import api from '@/api'
import i18n from "@/utils/i18n";

const getDefaultState = () => {
  return {
    token: getToken() || null,
    info: getUserInfo() || null,
    permission: getPermission() || null,
    pagePermission: getPagePermission() || {},
  }
}

const state = getDefaultState()

const actions = {
  // user login
  login({ commit, dispatch }, userInfo) {
    return new Promise((resolve, reject) => {
      api.login(userInfo).then(res => {
        commit('SET_TOKEN', res.data.token)
        let user = {
          id: res.data.id,
          name: res.data.name,
          email: res.data.email,
          isAdmin: res.data.is_admin,
        }
        commit('SET_USER', user)
        commit('SET_PERMISSION', res.data.permission)
        setToken(res.data.token)
        setUserInfo(user)
        setPermission(res.data.permission)
        resolve()
      }).catch(err => {
        reject(err)
      })
    })
  },
  // user logout
  logout({ commit, dispatch, state }) {
    return new Promise((resolve, reject) => {
      // 暂时只清除客户端Token
      resetRouter()
      removeAllCookies()
      removePermission()
      removeMenu()
      removePagePermission()
      commit('RESET_STATE')
      dispatch('tagsView/delAllVisitedViews', null, { root: true })
      dispatch('tagsView/delAllCachedViews', null, { root: true })
      resolve()
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      removeToken() // must remove  token  first
      commit('RESET_STATE')
      resolve()
    })
  },
  // 页面权限
  getPagePermissions({ commit }) {
    return new Promise(((resolve, reject) => {
      api.getPagePermission().then(res => {
        commit('SET_PAGE_PERMISSIONS', res.data)
        setPagePermission(res.data)
      }).catch(() => {
        reject()
      })
    }))
  },
}

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState())
  },
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_USER: (state, user) => {
    state.info = user
  },
  SET_PERMISSION: (state, permissions) => {
    state.permission = permissions
  },
  SET_PAGE_PERMISSIONS(state, pagePermission) {
    state.pagePermission = pagePermission
  },
}

function getDictFromPagePermission(pagePermissionArr) {
  let dict = {}

  Object.keys(pagePermissionArr).forEach(key => {
    Object.keys(pagePermissionArr[key]).forEach(key2 => {
      let item = pagePermissionArr[key][key2]
      // 修改后端传回数据
      dict[item.name] = i18n.locale === 'en' ? item.display_name : item.description
    })
  })
  return dict
}

const getters = {
  userInfo(state){
    return state.info
  },
  pagePermissionDict(state) {
    return getDictFromPagePermission(state.pagePermission)
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}

