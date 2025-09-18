<template>
  <div id="tags-view-container" class="tags-view-container" style="margin: 20px">
    <tabs v-model="currentTab" type="card" closable @tab-remove="removeTab" @tab-right-click="rightClickTab">
      <tab-pane
        ref="tag"
        v-for="(tag, index) in visitedViews"
        :key="tag.name"
        :label="menusDict[tag.title]"
        :name="tag.name"
        :closable="!tag.meta.affix"
      >
      </tab-pane>
    </tabs>
    <slot></slot>
    <ul v-show="visible" :style="{left:left+'px',top:top+'px'}" class="contextmenu">
      <li @click="refreshSelectedTag(selectedTag)">{{ $t('common.Refresh') }}</li>
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)">{{ $t('common.Close') }}</li>
      <li @click="closeOthersTags">{{ $t('common.CloseOthers') }}</li>
      <li @click="closeAllTags(selectedTag)">{{ $t('common.CloseAll') }}</li>
    </ul>
  </div>
</template>

<script>
import Tabs from '@/components/Tab/Tabs'
import TabPane from "@/components/Tab/TabPane";
import path from 'path'
import {mapGetters} from 'vuex'

export default {
  components: { Tabs, TabPane },
  data() {
    return {
      visible: false,
      top: 0,
      left: 0,
      selectedTag: {},
      affixTags: [],
    }
  },
  computed: {
    ...mapGetters({
      menusDict: 'app/menusDict',
    }),
    visitedViews() {
      return this.$store.state.tagsView.visitedViews
    },
    routes() {
      return this.$store.state.permission.routes
    },
    currentTab: {
      get: function () {
        return this.$route.name
      },
      set: function (newRouteName, old) {
        const tags = this.visitedViews
        for (const tag of tags) {
          // 切换tab变化时 只有变化了tab才能进行跳转
          if (tag.name === newRouteName && this.currentTab !== newRouteName) {
            this.$router.push({ path: tag.path, query: tag.query, fullPath: tag.fullPath })
            break
          }
        }
      },
    },
  },
  watch: {
    $route() {
      this.addTags()
      this.moveToCurrentTag()
    },
    visible(value) {
      if (value) {
        document.body.addEventListener('click', this.closeMenu)
      } else {
        document.body.removeEventListener('click', this.closeMenu)
      }
    },
  },
  mounted() {
    this.initTags()
    this.addTags()
  },
  methods: {
    removeTab(targetName) {
      const tags = this.visitedViews
      for (const tag of tags) {
        if (tag.name === targetName) {
          this.closeSelectedTag(tag)
          break
        }
      }
    },
    isActive(route) {
      return route.path === this.$route.path
    },
    isAffix(tag) {
      return tag.meta && tag.meta.affix
    },
    filterAffixTags(routes, basePath = '/') {
      let tags = []
      routes.forEach(route => {
        if (route.meta && route.meta.affix) {
          const tagPath = path.resolve(basePath, route.path)
          tags.push({
            fullPath: tagPath,
            path: tagPath,
            name: route.name,
            meta: { ...route.meta },
          })
        }
        if (route.children) {
          const tempTags = this.filterAffixTags(route.children, route.path)
          if (tempTags.length >= 1) {
            tags = [...tags, ...tempTags]
          }
        }
      })
      return tags
    },
    initTags() {
      const affixTags = this.affixTags = this.filterAffixTags(this.$router.matcher.getRoutes())
      for (const tag of affixTags) {
        // Must have tag name
        if (tag.name) {
          this.$store.dispatch('tagsView/addVisitedView', tag)
        }
      }
    },
    addTags() {
      const { name } = this.$route
      if (name) {
        this.$store.dispatch('tagsView/addView', this.$route)
      }
      return false
    },
    moveToCurrentTag() {
      const tags = this.visitedViews
      this.$nextTick(() => {
        for (const tag of tags) {
          if (tag.path === this.$route.path) {
            // when query is different then update
            if (tag.fullPath !== this.$route.fullPath) {
              this.$store.dispatch('tagsView/updateVisitedView', this.$route)
            }
            break
          }
        }
      })
    },
    refreshSelectedTag(view) {
      this.$store.dispatch('tagsView/delCachedView', view).then(() => {
        const { fullPath } = view
        this.$nextTick(() => {
          this.$router.replace({
            path: '/redirect' + fullPath,
          })
        })
      })
    },
    closeSelectedTag(view) {
      this.$store.dispatch('tagsView/delView', view).then(({ visitedViews }) => {
        if (this.isActive(view)) {
          this.toLastView(visitedViews, view)
        }
      })
    },
    closeOthersTags() {
      if (this.selectedTag.name !== this.currentTab) {
        this.$router.push(this.selectedTag)
      }
      this.$store.dispatch('tagsView/delOthersViews', this.selectedTag).then(() => {
        this.moveToCurrentTag()
      })
    },
    closeAllTags(view) {
      this.$store.dispatch('tagsView/delAllViews').then(({ visitedViews }) => {
        if (this.affixTags.some(tag => tag.path === view.path)) {
          return
        }
        this.toLastView(visitedViews, view)
      })
    },
    toLastView(visitedViews, view) {
      const latestView = visitedViews.slice(-1)[0]
      if (latestView) {
        this.$router.push(latestView.fullPath)
      } else {
        // now the default is to redirect to the home page if there is no tags-view,
        // you can adjust it according to your needs.
        if (view.name === 'Dashboard') {
          // to reload home page
          this.$router.replace({ path: '/redirect' + view.fullPath })
        } else {
          this.$router.push('/')
        }
      }
    },
    openMenu(event, tag) {
      const menuMinWidth = 105
      const offsetLeft = this.$el.getBoundingClientRect().left // container margin left
      const offsetWidth = this.$el.offsetWidth // container width
      const maxLeft = offsetWidth - menuMinWidth // left boundary
      const left = event.clientX - offsetLeft + 15 // 15: margin right

      if (left > maxLeft) {
        this.left = maxLeft
      } else {
        this.left = left
      }

      this.top = event.clientY
      this.visible = true
      this.selectedTag = tag
    },
    closeMenu() {
      this.visible = false
    },
    rightClickTab(tab, event) {
      const menuMinWidth = 105
      const offsetLeft = this.$el.getBoundingClientRect().left // container margin left
      const offsetWidth = this.$el.offsetWidth // container width
      const maxLeft = offsetWidth - menuMinWidth // left boundary
      const left = event.clientX - offsetLeft + 15 // 15: margin right

      if (left > maxLeft) {
        this.left = maxLeft
      } else {
        this.left = left
      }

      this.top = event.clientY
      this.visible = true
      this.selectedTag = this.getTagByTab(tab)
    },
    getTagByTab(tab) {
      const tags = this.visitedViews

      return tags[tab.index]
    },
  },
}
</script>

<style lang="scss" scoped>
.tags-view-container {
  //height: 34px;
  //width: 100%;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);

  .tags-view-wrapper {
    .tags-view-item {
      user-select: none;
      display: inline-block;
      position: relative;
      cursor: pointer;
      height: 26px;
      line-height: 26px;
      border: 1px solid #d8dce5;
      color: #495060;
      background: #fff;
      padding: 0 8px;
      font-size: 12px;
      margin-left: 5px;
      margin-top: 4px;

      &:first-of-type {
        margin-left: 15px;
      }

      &:last-of-type {
        margin-right: 15px;
      }

      &.active {
        background-color: #42b983;
        color: #fff;
        border-color: #42b983;

        &::before {
          content: '';
          background: #fff;
          display: inline-block;
          width: 8px;
          height: 8px;
          border-radius: 50%;
          position: relative;
          margin-right: 2px;
        }
      }
    }
  }

  .contextmenu {
    margin: 0;
    background: #fff;
    z-index: 3000;
    position: absolute;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, .3);

    li {
      margin: 0;
      padding: 7px 16px;
      cursor: pointer;

      &:hover {
        background: #eee;
      }
    }
  }
}
</style>

<style lang="scss">
//reset element css of el-icon-close
.tags-view-wrapper {
  .tags-view-item {
    .el-icon-close {
      width: 16px;
      height: 16px;
      vertical-align: 2px;
      border-radius: 50%;
      text-align: center;
      transition: all .3s cubic-bezier(.645, .045, .355, 1);
      transform-origin: 100% 50%;

      &:before {
        transform: scale(.6);
        display: inline-block;
        vertical-align: -3px;
      }

      &:hover {
        background-color: #b4bccc;
        color: #fff;
      }
    }
  }
}
</style>
