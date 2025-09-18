<template>
  <div class="infinite-list-wrapper" style="overflow-x: hidden;overflow-y: auto;height: calc(100vh - 150px);">
    <div v-infinite-scroll="load" infinite-scroll-disabled="disabled">
      <slot></slot>
    </div>
    <div class="flex-center">
      <p v-if="loading" style="color: #409EFF"><i class="el-icon-loading"></i> Loading...</p>
      <p v-if="noMore" style="color: #409EFF">No More</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Waterfall',
  props: {
    total: {
      type: Number,
      default: 0,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    noMore: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {}
  },
  computed: {
    disabled() {
      return this.loading || this.noMore
    },
  },
  methods: {
    load() {
      // 防止 一开始就加载
      this.$emit('update:loading', true)
      if (this.total !== 0) {
        setTimeout(() => {
          this.$emit('fetch')
        }, 1000)
      }
    },
  },
}
</script>

<style scoped>

</style>
