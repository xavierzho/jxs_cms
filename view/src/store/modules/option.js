import api from '@/api'
import i18n from '@/utils/i18n'

const state = {
  role: [],
  user: [],
  permission: [],

  adDim: [
    {value: "daily", label: i18n.t("option.dim[0]")},
    {value: "month", label: i18n.t("option.dim[1]")},
  ],

  dateTimeType: [
    {value: "created", label: i18n.t("option.dateTimeType[0]")},
    {value: "finish", label: i18n.t("option.dateTimeType[1]")},
  ],

  userType: [
    {"value": true, "label": i18n.t("option.YesOrNo[0]")},
    {"value": false, "label": i18n.t("option.YesOrNo[1]")},
  ],

  inviteUserType: [
    {"value": "user", "label": i18n.t("option.inviteUserType[0]")},
    {"value": "parent_user", "label": i18n.t("option.inviteUserType[1]")},
  ],

  prizeType: [
    {"value": 1, "label": i18n.t("option.prizeType[0]")},
    {"value": 10, "label": i18n.t("option.prizeType[1]")},
    {"value": 20, "label": i18n.t("option.prizeType[2]")},
  ],

  inquireItemLogType: [],
  inquireGachaType: [],
  inquireBalanceSourceType: [],
  inquireBalanceChannelType: [],
  inquireBalancePaySourceType:[],
  inquireCouponType: [],
  inquireCouponActionType: [],
  activityCostAwardLogType: [],
}

const actions = {
  update({ commit }, optionName) {
    // 强制更新选项数据
    let upperCaseOption = optionName.slice(0, 1).toUpperCase().concat(optionName.slice(1))
    let apiName = `get${upperCaseOption}Options`
    api[apiName]().then(res => {
      commit('SET_OPTION', {
        name: optionName,
        data: res.data,
      })
    })
  },
}

const mutations = {
  SET_OPTION(state, payload) {
    state[payload.name] = payload.data
  },
}

const getters = {
  user(state) {
    let arr = JSON.parse(JSON.stringify(state.user)) || []
    return arr.map((el, index) => {
      el.label = el.name
      el.value = el.id
      el.key = index
      return el
    })
  },
  adDim(state){
    return state.adDim
  },
  dateTimeType(state){
    return state.dateTimeType
  },
  userType(){
    return state.userType
  },
  inviteUserType(state){
    return state.inviteUserType
  },
  inquireItemLogType(state){
    return state.inquireItemLogType
  },
  inquireGachaType(state){
    return state.inquireGachaType
  },
  inquireBalanceSourceType(){
    return state.inquireBalanceSourceType
  },
  inquireBalanceChannelType(){
    return state.inquireBalanceChannelType
  },
  inquireBalancePaySourceType(){
    return state.inquireBalancePaySourceType
  },
  activityCostAwardLogType(){
    return state.activityCostAwardLogType
  },
  inquireCouponType(){
    return state.inquireCouponType
  },
  inquireCouponActionType(){
    return state.inquireCouponActionType
  },
  prizeType(state){
    return state.prizeType
  },
}



export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
}
