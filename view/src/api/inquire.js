import request from '@/utils/request'
import { saveFile } from '@/utils/file-saver'

export default {
  // item
  getInquireItemLogTypeOptions() {
    return request.get('inquire/item/log-type/options', )
  },
  getInquireItemLog(params) {
    return request.get('inquire/item/log', {params})
  },
  exportInquireItemLog(params) {
    return saveFile('inquire/item/log/export', params, "user_item_log")
  },
  getInquireItemDetail(params) {
    return request.get('inquire/item/detail', {params})
  },
  exportInquireItemDetailLog(params) {
    return saveFile('inquire/item/detail/export', params, "user_item_detail_log")
  },

  // gacha
  getInquireGachaTypeOptions() {
    return request.get('inquire/gacha/type/options', )
  },
  getInquireGachaRevenue(params) {
    return request.get('inquire/gacha/revenue', {params})
  },
  getInquireGachaDetail(params) {
    return request.get('inquire/gacha/detail', {params})
  },
  exportInquireGachaDetail(params) {
    let fileName = `gacha_detail-${params.gacha_id}`
    if (params.box_out_no){
      fileName += `-${params.box_out_no}`
    }
    return saveFile('inquire/gacha/detail/export', params, fileName)
  },

  // balance
  getInquireBalance(params) {
    return request.get('inquire/balance', {params})
  },
  exportInquireBalance(params) {
    return saveFile('inquire/balance/export', params, "user_balance_log")
  },
  getInquireBalanceSourceTypeOptions(params) {
    return request.get('inquire/balance/source-type/options', {params})
  },
  getInquireBalanceChannelTypeOptions(params) {
    return request.get('inquire/balance/channel-type/options', {params})
  },

  getInquireBalancePaySourceTypeOptions(params) {
    return request.get('inquire/balance/pay-source-type/options', {params})
  },
  addComment(params){
    return request.post('inquire/balance/comment', params)
  },
  deleteComment(params){
    return request.delete('inquire/balance/comment', {params})
  },

  // coupon
  getInquireCoupon(params) {
    return request.get('inquire/coupon', {params})
  },
  exportInquireCoupon(params) {
    return saveFile('inquire/coupon/export', params, "user_coupon_log")
  },
  getInquireCouponTypeOptions(params) {
    return request.get('inquire/coupon/type/options', {params})
  },
  getInquireCouponActionTypeOptions(params) {
    return request.get('inquire/coupon/action/options', {params})
  },

  // invite
  getInquireInviteRec(params) {
    return request.get('inquire/invite-rec', {params})
  },
  exportInquireInviteRec(params) {
    return saveFile('inquire/invite-rec/export', params, "invite_rec")
  },
}
