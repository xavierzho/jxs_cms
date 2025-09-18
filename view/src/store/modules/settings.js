import defaultSettings from '@/settings'
import { getShowEventNames, getTagViewSetting, setShowEventNames, setTagViewSetting } from '@/utils/cache'

const { showSettings, tagsView, fixedHeader, sidebarLogo } = defaultSettings

const state = {
  showSettings: showSettings,
  tagsView: getTagViewSetting() || false,
  fixedHeader: fixedHeader,
  sidebarLogo: sidebarLogo,
  eventNames: getShowEventNames() || [],
}

const mutations = {
  CHANGE_SETTING: (state, { key, value }) => {
    // eslint-disable-next-line no-prototype-builtins
    if (state.hasOwnProperty(key)) {
      state[key] = value
    }
  },
  SET_EVENT_NAMES: (state, eventNames) => {
    state.eventNames = eventNames
    setShowEventNames(eventNames)
  },
  SET_TAG_VIEW: (state, value) => {
    state.tagsView = value
    setTagViewSetting(value)
  },
}

const actions = {
  changeSetting({ commit }, data) {
    commit('CHANGE_SETTING', data)
  },
  setShowEventNames({ commit }, eventNames) {
    commit('SET_EVENT_NAMES', eventNames)
  },
  setTagsView({ commit }, value) {
    if (typeof value === 'boolean') {
      commit('SET_TAG_VIEW', value)
    }
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
}

