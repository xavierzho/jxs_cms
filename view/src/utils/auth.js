import store from '@/store'

export function checkPermission(value){
  const permissions = store.state.user.permission || []
  const perm = value
  let hasPermission = false
  if (value && typeof value === 'string') {
     hasPermission = permissions.some(permission => {
      return value === permission
    })
  } else if (value && value instanceof Array) {
    // 传入数组
    if (value.length > 0) {
       hasPermission = permissions.some(permission => {
        return perm.includes(permission)
      })
    }
  }

  return hasPermission
}
