<template>
  <section class="app-main" :style="{'padding-top': needTagsView?'0px':'20px !important','min-height':needTagsView?'calc(100vh - 145px)':'calc(100vh - 50px)'}">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key"/>
      </keep-alive>
    </transition>
  </section>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'AppMain',
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    },
    ...mapState({
      needTagsView: state => state.settings.tagsView,
    }),
  },
}
</script>

<style lang="scss" scoped>
.app-main {
  /*50 = navbar  */
  padding: 0 20px 20px 20px;
  //min-height: calc(100vh - 130px);
  width: 100%;
  position: relative;
  overflow: hidden;
}

.fixed-header + .app-main {
  padding-top: 50px;
}

.hasTagsView {
  .app-main {
    /* 84 = navbar + tags-view = 50 + 34 */
    min-height: calc(100vh - 84px);
  }

  .fixed-header + .app-main {
    padding-top: 84px;
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
