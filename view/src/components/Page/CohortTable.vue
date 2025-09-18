<template>
  <chart-table ref="tab" :loading="loading" :tableData="tableData" :chartOption="chartOption"
               :show-summary="true" :summary-method="summaryMethod"
  >
    <template v-slot:table-column>
      <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
      <el-table-column prop="total" :label="$t('report.cohort.total')" min-width="90px" align="center">
        <template v-slot="scope">
          {{ scope.row.total | localeNum }}
        </template>
      </el-table-column>
      <el-table-column prop="first_day" :label="$t('report.cohort.first_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.first_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.first_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="second_day" :label="$t('report.cohort.second_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.second_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.second_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="third_day" :label="$t('report.cohort.third_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.third_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.third_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="fourth_day" :label="$t('report.cohort.fourth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.fourth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.fourth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="fifth_day" :label="$t('report.cohort.fifth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.fifth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.fifth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="sixth_day" :label="$t('report.cohort.sixth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.sixth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.sixth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="seventh_day" :label="$t('report.cohort.seventh_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.seventh_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.seventh_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="fourteenth_day" :label="$t('report.cohort.fourteenth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.fourteenth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.fourteenth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="thirtieth_day" :label="$t('report.cohort.thirtieth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.thirtieth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.thirtieth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="sixtieth_day" :label="$t('report.cohort.sixtieth_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.sixtieth_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.sixtieth_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="ninety_day" :label="$t('report.cohort.ninety_day')" min-width="90px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.ninety_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.ninety_day | localeNum }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="no_180_day" :label="$t('report.cohort.no_180_day')" min-width="100px" align="center">
        <template v-slot="scope">
          <div class="table-line-height">{{ scope.row.no_180_day_rate | rate }}</div>
          <div class="table-line-height">{{ scope.row.no_180_day | localeNum }}</div>
        </template>
      </el-table-column>
    </template>
  </chart-table>
</template>

<script>
import ChartTable from "@/components/Table/ChartTable.vue";

export default {
  name: 'CohortTable',
  components: {ChartTable},
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
    summaryMethod: Function,
  },
  methods:{
    rebuildCharts() {
      this.$refs.tab.rebuildCharts()
    },
  },
}
</script>

<style scoped>

</style>
