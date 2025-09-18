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
}
