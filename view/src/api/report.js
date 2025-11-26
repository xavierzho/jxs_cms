import request from '@/utils/request'
import {saveFile} from "@/utils/file-saver";

export default {
  getReportRevenue(params) {
    return request.get('report/revenue', {params})
  },
  getReportCohort(params) {
    return request.get('report/cohort', {params})
  },
  getReportRealtime(params) {
    return request.get('report/realtime', {params})
  },
  getReportMarket(params) {
    return request.get('report/market', {params})
  },
  exportReportMarket(params) {
    return saveFile('report/market/export', params, "market_report")
  },
  getReportBet(params) {
    return request.get('report/bet', {params})
  },
  getReportInviteBet(params) {
    return request.get('report/invite-bet', {params})
  },
  exportReportInviteBet(params) {
    return saveFile('report/invite-bet/export', params, "invite_bet")
  },
  getReportInviteBetDaily(params) {
    return request.get('report/invite-bet/daily', {params})
  },
  exportReportInviteBetDaily(params) {
    return saveFile('report/invite-bet/daily/export', params, "invite_bet_daily")
  },
  getReportDashboard(params) {
    return request.get('report/dashboard', {params})
  },
  //发货报表
  getReportDeliveryOrderDaily(params) {
    return request.get('report/delivery_order', {params})
  },
  exportReportDeliveryOrderDaily(params) {
    return saveFile('report/delivery_order/export', params, "report_delivery_order_daily")
  },
  //召回报表
  getReportRecall(params) {
    return request.get('report/recall', {params})
  },
  exportReportRecall(params) {
    return saveFile('report/recall/export', params, "recall")
  },
  getReportRecallDaily(params) {
    return request.get('report/recall/daily', {params})
  },
  exportReportRecallDaily(params) {
    return saveFile('report/recall/daily/export', params, "recall_daily")
  },
}
