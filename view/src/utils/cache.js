import Cookies from 'js-cookie'

const TokenKey = 'chaoshe_cms_token'
const UserInfoKey = 'chaoshe_cms_user_info'
const MenuKey = 'chaoshe_cms_menu'
const LangKey = 'chaoshe_cms_lang'
const PermissionKey = 'chaoshe_cms_permission'
const EventNamesKey = 'chaoshe_cms_event_names'
const PagePermissionKey = 'chaoshe_cms_page_permission'
const TagViewKey = 'chaoshe_cms_tag_view'
const expireDay = 1

// 所有的cookies操作
export function removeAllCookies() {
  Cookies.remove(TokenKey)
  Cookies.remove(UserInfoKey)
}

// 细分cookies操作
export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token, { expires: expireDay })
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getUserInfo() {
  let jsonUserInfo = Cookies.get(UserInfoKey)
  if (typeof jsonUserInfo === 'undefined') {
    return null
  }
  return JSON.parse(jsonUserInfo)
}

export function setUserInfo(userInfo) {
  return Cookies.set(UserInfoKey, JSON.stringify(userInfo), { expires: expireDay })
}

export function removeUserInfo() {
  return Cookies.remove(UserInfoKey)
}

export function getMenu() {
  let jsonStr = localStorage.getItem(MenuKey)
  if (typeof jsonStr === 'undefined') {
    return null
  }
  return JSON.parse(jsonStr)
}

export function setMenu(menu) {
  localStorage.setItem(MenuKey, JSON.stringify(menu))
}

export function removeMenu() {
  localStorage.removeItem(MenuKey)
}

export function getLang() {
  return Cookies.get(LangKey)
}

export function setLang(lang) {
  return Cookies.set(LangKey, lang, { expires: expireDay })
}

export function getPermission() {
  let jsonStr = localStorage.getItem(PermissionKey)
  if (typeof jsonStr === 'undefined') {
    return null
  }
  return JSON.parse(jsonStr)
}

export function setPermission(permission) {
  localStorage.setItem(PermissionKey, JSON.stringify(permission))
}

export function removePermission() {
  localStorage.removeItem(PermissionKey)
}

export function getTagViewSetting() {
  let jsonStr = localStorage.getItem(TagViewKey)
  if (typeof jsonStr === 'undefined') {
    return null
  }
  return JSON.parse(jsonStr)
}

export function setTagViewSetting(value) {
  localStorage.setItem(TagViewKey, JSON.stringify(value))
}

export function removeTagViewSetting() {
  localStorage.removeItem(TagViewKey)
}

export function getShowEventNames() {
  let jsonStr = localStorage.getItem(EventNamesKey)
  if (typeof jsonStr === 'undefined') {
    return null
  }
  return JSON.parse(jsonStr)
}

export function setShowEventNames(eventNames) {
  localStorage.setItem(EventNamesKey, JSON.stringify(eventNames))
}

export function removeShowEventNames() {
  localStorage.removeItem(EventNamesKey)
}

export function getPagePermission() {
  let jsonStr = localStorage.getItem(PagePermissionKey)
  if (typeof jsonStr === 'undefined') {
    return null
  }
  return JSON.parse(jsonStr)
}

export function setPagePermission(pagePermission) {
  localStorage.setItem(PagePermissionKey, JSON.stringify(pagePermission))
}

export function removePagePermission() {
  localStorage.removeItem(PagePermissionKey)
}
