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

  pointType: [
    {"value": 0, "label": i18n.t("option.pointType[0]")},
    {"value": 2, "label": i18n.t("option.pointType[2]")},
    {"value": 10, "label": i18n.t("option.pointType[10]")},
    {"value": 11, "label": i18n.t("option.pointType[11]")},
    {"value": 100, "label": i18n.t("option.pointType[100]")},
  ],

  awardType: [
    {"value": 10, "label": i18n.t("option.awardType[10]")},
    {"value": 20, "label": i18n.t("option.awardType[20]")},
    {"value": 30, "label": i18n.t("option.awardType[30]")},
  ],

  recallUserType: [
    {"value": "user", "label": i18n.t("option.recallUserType[0]")},
    {"value": "parent_user", "label": i18n.t("option.recallUserType[1]")},
  ],

  userChannel: [
    {"value": 1, "label": i18n.t("option.userChannel[1]")},
    {"value": 2, "label": i18n.t("option.userChannel[2]")},
    {"value": 4, "label": i18n.t("option.userChannel[4]")},
    {"value": 8, "label": i18n.t("option.userChannel[8]")},
    {"value": 16, "label": i18n.t("option.userChannel[16]")},
    {"value": 32, "label": i18n.t("option.userChannel[32]")},
    {"value": 64, "label": i18n.t("option.userChannel[64]")},
    {"value": 128, "label": i18n.t("option.userChannel[128]")},
  ],

  inquireItemLogType: [],
  inquireGachaType: [],
  inquireBalanceSourceType: [],
  inquireBalanceChannelType: [],
  inquireBalancePaySourceType:[],
  inquireCouponType: [],
  inquireCouponActionType: [],
  activityCostAwardLogType: [],
  inquireTaskType: [],
  inquireTaskKey: [],
  inquireBalanceType:[],
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
    // getters[payload.name]()
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
  recallUserType(state){
    return state.recallUserType
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
  pointType(state){
    return state.pointType
  },
  awardType(state){
    return state.awardType
  },
  inquireTaskType(state){
    return state.inquireTaskType
  },
  inquireTaskKey(state){
    return state.inquireTaskKey
  },
  inquireBalanceType(){
    return state.inquireBalanceType
  },
  userChannel(state){
    return state.userChannel
  },
}



export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters,
}
