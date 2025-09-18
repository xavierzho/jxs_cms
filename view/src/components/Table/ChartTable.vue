<template>
  <div>
    <chart ref="chart" :options="chartOption" :seriesData="chartOption.series" width="100%" height="30vh"></chart>
    <el-table ref="tableData" v-loading="loading" :data="tableData"
              height="50vh" border style="width: 100%;margin-bottom: 20px"
              :cell-class-name="cellClass"
              :show-summary="showSummary"
              :summary-method="summaryMethod"
              @selection-change="handleSelectionChange"
              @sort-change="handleSortChange"
              @expand-change="expandChange"
    >
      <slot name="table-column"></slot>
    </el-table>
  </div>
</template>

<script>
import Chart from "@/components/Chart/Chart.vue";

export default {
  name: "ChartTable",
  components: {Chart},
  props: {
    tableData: {
      type: Array,
      default: () => {
        return []
      },
    },
    loading: {
      type: Boolean,
      default: () => {
        return false
      },
    },
    chartOption: {
      type: Object,
      default: () => {
        return {series: []}
      }
    },
    showSummary: {
      type: Boolean,
      default: false,
    },
    summaryMethod: Function,
  },
  updated() {
    this.$nextTick(() => {
      this.$refs.tableData.doLayout()
    })
  },
  methods:{
    rebuildCharts() {
      this.$refs.chart.rebuildCharts()
    },
    cellClass(cell) {
      if (cell.column.property === 'date') {
        return ''
      }
      let value = cell.row[cell.column.property]
      if (value === null || value === '' || typeof value === 'undefined') {
        return ''
      }
      if (value === 0 || value === '0' || value === '0%') {
        return ''
      }
      return 'table-cell-color'
    },
    handleSelectionChange(val) {
      this.$emit('selectionChange', val)
    },
    handleSortChange(val) {
      if (val.order !== null) {
        val.order = val.order === 'ascending' ? 'ASC' : 'DESC'
      }
      this.$emit('sortChange', val)
    },
    expandChange(row, expandedRows) {
      this.$emit('expandChange', row, expandedRows)
    },
  },
}
</script>

<style scoped>

</style>
