import Cookies from 'js-cookie'
import api from '@/api'
import { getLang, getMenu, setLang, setMenu } from '@/utils/cache'
import { resetRouter } from '@/router'

const state = {
  sidebar: {
    opened: Cookies.get('sidebarStatus') ? !!+Cookies.get('sidebarStatus') : true,
    withoutAnimation: false,
  },
  device: 'desktop',
  menus: getMenu() || [],
  lang: getLang() || 'zh',
}

const actions = {
  toggleSideBar({ commit }) {
    commit('TOGGLE_SIDEBAR')
  },
  closeSideBar({ commit }, { withoutAnimation }) {
    commit('CLOSE_SIDEBAR', withoutAnimation)
  },
  toggleDevice({ commit }, device) {
    commit('TOGGLE_DEVICE', device)
  },
  toggleLang({ commit, dispatch }, lang) {
    commit('TOGGLE_LANG', lang)
    dispatch('getMenus')
    setLang(lang)
  },
  getMenus({ commit }) {
    return new Promise(((resolve, reject) => {
      api.getMenuList().then(res => {
        commit('SET_MENUS', res.data)
        // 缓存相应的菜单
        setMenu(res.data||[])
        resetRouter()
        resolve()
      }).catch(() => {
        reject()
      })
    }))
  },
}

const mutations = {
  TOGGLE_SIDEBAR: state => {
    state.sidebar.opened = !state.sidebar.opened
    state.sidebar.withoutAnimation = false
    if (state.sidebar.opened) {
      Cookies.set('sidebarStatus', 1)
    } else {
      Cookies.set('sidebarStatus', 0)
    }
  },
  CLOSE_SIDEBAR: (state, withoutAnimation) => {
    Cookies.set('sidebarStatus', 0)
    state.sidebar.opened = false
    state.sidebar.withoutAnimation = withoutAnimation
  },
  TOGGLE_DEVICE: (state, device) => {
    state.device = device
  },
  TOGGLE_LANG: (state, lang) => {
    state.lang = lang
  },
  SET_MENUS(state, menus) {
    state.menus = menus
  },
}

function getDictFromMenus(menuArr) {
  let dict = {}
  for (let i = 0; i < menuArr.length; i++) {
    if (dict.hasOwnProperty(menuArr[i].name)) {
      console.error(`This menu ${menuArr[i].name} has been registered`)
    } else {
      dict[menuArr[i].name] = menuArr[i].title
    }
    if (menuArr[i].hasOwnProperty('children')) {
      let childDict = getDictFromMenus(menuArr[i].children)
      dict = Object.assign(dict, childDict)
    }
  }
  return dict
}

const getters = {
  menusDict(state) {
    return getDictFromMenus(state.menus)
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
}
