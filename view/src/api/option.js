import request from '@/utils/request'

export default {
  getRoleOptions() {
    return request.get('management/role/options')
  },
  getUserOptions() {
    return request.get('management/user/options')
  },
  getPermissionOptions() {
    return request.get('management/permission/options')
  },
}
