import request from '@/utils/request'

export default {
  createAdvertisingList(params) {
    return request.post('advertising/analysis', params)
  },
  getAdvertisingList(params) {
    return request.get('advertising/analysis', { params })
  },
  updateAdvertisingList(date, params) {
    return request.put('advertising/analysis/'+date, params)
  },
}
