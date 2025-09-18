<template>
  <div id="chart" :style="{width:width, height:height}"></div>
</template>

<script>
export default {
  name: 'Chart',
  props: {
    options: {
      type: Object,
      required: true,
    },
    width: {
      type: String,
      default: '400px',
    },
    height: {
      type: String,
      default: '400px',
    },
    seriesData: {
      type: Array,
      default: () => {
        return []
      },
    },
  },
  data() {
    return {
      chart: null,
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initCharts()
    })
  },
  watch: {
    options: {
      handler(options) {
        if (this.chart !== null){
          this.chart.setOption(this.options)
        }
      },
      deep: true,
    },
  },
  methods: {
    initCharts() {
      this.chart = this.$echarts.init(this.$el)
      this.setOptions()
    },
    rebuildCharts() {
      this.chart.dispose()
      this.initCharts()
    },
    setOptions(params = {}) {
      let chartOptions = Object.assign({}, this.options, params)
      this.chart.setOption(chartOptions)
      this.resize()
    },
    resize() {
      this.$nextTick(() => {
        this.chart.resize()
      })
    },
  },
}
</script>

<style scoped>

</style>
