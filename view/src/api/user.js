import request from '@/utils/request'

export default {
  // user
  createUser(params) {
    return request.post('management/user/create', params)
  },
  getUserDetail(params) {
    return request.get('management/user/detail', { params })
  },
  getUserList(params) {
    return request.get('management/user', { params })
  },
  updateUser({ id, ...params }) {
    return request.put('management/user/update/' + id, params)
  },
  updateUserSelf(params) {
    return request.put('management/user/update-self', params)
  },
  getPagePermission() {
    return request.get('management/user/page-permission')
  },
  // role
  createRole(params) {
    return request.post('management/role/create', params)
  },
  getRoleList(params) {
    return request.get('management/role', { params })
  },
  updateRole({ id, ...params }) {
    return request.put('management/role/update/' + id, params)
  },
}
