import request from '@/utils/request'
import {saveFile} from "@/utils/file-saver";

export default {
  getActivityCostAward(params) {
    return request.get('activity/cost-award', {params})
  },
  exportActivityCostAward(params) {
    return saveFile('activity/cost-award/export', params, "cost_award")
  },
  getActivityCostAwardLogTypeOptions(params) {
    return request.get('activity/cost-award-log/log-type/options', {params})
  },
  getActivityCostAwardLog(params) {
    return request.get('activity/cost-award-log', {params})
  },
  exportActivityCostAwardLog(params) {
    return saveFile('activity/cost-award-log/export', params, "cost_award_log")
  },

  // turntable
  getActivityTurntable(params) {
    return request.get('activity/turntable', {params})
  },
  exportActivityTurntable(params) {
    return saveFile('activity/turntable/export', params, "turntable")
  },
  getActivityStepByStep(params) {
    return request.get('activity/step-by-step', {params})
  },
  exportActivityStepByStep(params) {
    return saveFile('activity/step-by-step/export', params, "step_by_step")
  },
  getActivityStepByStepDetail(params) {
    return request.get('activity/step-by-step/detail', {params})
  },
  exportActivityStepByStepDetail(params) {
    return saveFile('activity/step-by-step/detail/export', params, "step_by_step_detail")
  },
  // signIn
  getActivitySignIn(params) {
    return request.get('activity/sign-in', {params})
  },
  exportActivitySignIn(params) {
    return saveFile('activity/sign-in/export', params, "sign-in")
  },
  // teamPK
  getActivityTeamPK(params) {
    return request.get('activity/team-pk', {params})
  },
  exportActivityTeamPK(params) {
    return saveFile('activity/team-pk/export', params, "team_pk")
  },
  // redemptionCode
  getActivityRedemptionCode(params) {
    return request.get('activity/redemption-code', {params})
  },
  getActivityRedemptionCodeAwardDetail(params) {
    return request.get('activity/redemption-code/award-detail', {params})
  },
  exportActivityRedemptionCodeLog(params) {
    return saveFile('activity/redemption-code/export', params, "activity_redemption_code_log")
  },
  exportActivityRedemptionCodeAwardDetail(params) {
    return saveFile('activity/redemption-code/detail/export', params, "activity_redemption_code_award_detail")
  },
}
