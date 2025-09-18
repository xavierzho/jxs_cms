const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.user.token,
  name: state => state.user.info ? state.user.info.name : '',
  menus: state => state.app.menus,
  lang: state => state.app.lang,
}
export default getters
