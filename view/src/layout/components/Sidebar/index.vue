<template>
  <div :class="{'has-logo':showLogo}">
    <logo v-if="showLogo" :collapse="isCollapse"/>
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :background-color="variables.menuBg"
        :text-color="variables.menuText"
        :unique-opened="false"
        :active-text-color="variables.menuActiveText"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item v-for="route in menus" :key="route.name" :item="route" :base-path="route.path"/>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Logo from './Logo'
import SidebarItem from './SidebarItem'
import variables from '@/styles/variables.scss'
import { resetRouter } from '@/router'

export default {
  components: { SidebarItem, Logo },
  computed: {
    menus() {
      // 重置
      resetRouter()
      let menus = JSON.parse(JSON.stringify(this.$store.getters.menus))
      menus = this.filterMenus(menus)
      return menus
    },
    ...mapGetters([
      'sidebar',
    ]),
    routes() {
      return this.$router.options.routes
    },
    activeMenu() {
      const route = this.$route
      const { meta, path } = route
      // if set path, the sidebar will highlight the path you set
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    },
    showLogo() {
      return this.$store.state.settings.sidebarLogo
    },
    variables() {
      return variables
    },
    isCollapse() {
      return !this.sidebar.opened
    },
  },
  methods: {
    filterMenus(menus) {
        for (let i = menus.length-1; i >= 0; i--){
          if (!menus[i].show){
            menus.splice(i, 1)
          }else if (menus[i].children){
            menus[i].children = this.filterMenus(menus[i].children)
          }
        }

        return menus
    }
  },
}
</script>
