<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
      </el-form-item>
    </el-form>

    <el-table ref="tableData" :data="tableData" v-loading="loading" border style="width: 100%;margin-bottom: 20px">
      <el-table-column prop="newUserCntSummary" :label="$t('report.realtime.summary.newUserCntSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.newUserCntSummary | localeNum }}</template>
      </el-table-column>
<!--      <el-table-column prop="maxOnlineUserCnt" :label="$t('report.realtime.summary.maxOnlineUserCnt')" min-width="90px" align="center">-->
<!--        <template v-slot="data">{{ data.row.maxOnlineUserCnt | localeNum }}</template>-->
<!--      </el-table-column>-->
      <el-table-column prop="payUserCntSummary" :label="$t('report.realtime.summary.payUserCntSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.payUserCntSummary | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="payAmountSummary" :label="$t('report.realtime.summary.payAmountSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.payAmountSummary | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="rechargeUserCntSummary" :label="$t('report.realtime.summary.rechargeUserCntSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.rechargeUserCntSummary | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="rechargeAmountSummary" :label="$t('report.realtime.summary.rechargeAmountSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.rechargeAmountSummary | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="drawAmountSummary" :label="$t('report.realtime.summary.drawAmountSummary')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.drawAmountSummary | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="payAmountARPPU" :label="$t('report.realtime.summary.payAmountARPPU')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.payAmountARPPU | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="newUserCntSummaryToday" :label="$t('report.realtime.summary.newUserCntSummaryToday')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.newUserCntSummaryToday | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="payAmountSummaryToday" :label="$t('report.realtime.summary.payAmountSummaryToday')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.payAmountSummaryToday | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="rechargeAmountSummaryToday" :label="$t('report.realtime.summary.rechargeAmountSummaryToday')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.rechargeAmountSummaryToday | localeNum2f }}</template>
      </el-table-column>
    </el-table>

    <chart ref="active" :options="chartOption.active" :seriesData="chartOption.active.series" width="100%" height="300px"></chart>
    <chart ref="pating" :options="chartOption.pating" :seriesData="chartOption.pating.series" width="100%" height="300px"></chart>
    <chart ref="payAmountSummary" :options="chartOption.payAmountSummary" :seriesData="chartOption.payAmountSummary.series" width="100%" height="300px"></chart>
    <chart ref="rechargeAmountSummary" :options="chartOption.rechargeAmountSummary" :seriesData="chartOption.rechargeAmountSummary.series" width="100%" height="300px"></chart>
    <chart ref="newUserCnt" :options="chartOption.newUserCnt" :seriesData="chartOption.newUserCnt.series" width="100%" height="300px"></chart>
    <chart ref="newUserCntSummary" :options="chartOption.newUserCntSummary" :seriesData="chartOption.newUserCntSummary.series" width="100%" height="300px"></chart>
  </div>
</template>

<script>
import Chart from "@/components/Chart/Chart.vue";
import moment from "moment/moment";
import api from "@/api";
import DatePicker from "@/components/DatePicker.vue";

export default {
  name: "Realtime",
  components: {DatePicker, Chart},
  data() {
    return {
      loading: false,
      tableData: [],
      tData: {
      },
      yData: {
      },
      chartTypeMap: {
        active: ['active'],
        pating: ['pating', 'patingFirstPrize', 'patingGashapon', 'patingChao', 'patingHole', 'patingMarket', 'patingDeliver', 'patingShop'],
        payAmountSummary: ['payAmountSummary'],
        rechargeAmountSummary: ['rechargeAmountSummary'],
        newUserCnt: ['newUserCnt'],
        newUserCntSummary: ['newUserCntSummary'],
      },
      chartOption:{
        active: {},
        pating: {},
        payAmountSummary: {},
        rechargeAmountSummary: {},
        newUserCnt: {},
        newUserCntSummary: {},
      },
    }
  },
  created() {
    this.fetch()
  },
  computed: {
    timeRangeArr() {
      let timeArr = []
      let start = moment('00:10:00', 'HH:mm:ss')
      let end = moment('23:50:00', 'HH:mm:ss')
      while (start <= end) {
        timeArr.push(start.format('HH:mm:ss'))
        start = start.add(10, 'minute')
      }

      timeArr.push("24:00:00")
      return timeArr
    },
  },
  methods: {
    fetch(params = {}) {
      this.loading = true
      api.getReportRealtime(params).then(res => {
        this.tData = res.data["tData"] || {}
        this.yData = res.data["yData"] || {}
        this.tableData = [res.data["summaryData"]] || []
      }).catch(() => {
        this.tData = {}
        this.yData = {}
        this.tableData = []
      }).finally(() => {
        this.parseData()
        this.rebuildCharts()
        this.loading = false
      })
    },
    parseData(){
      Object.keys(this.chartTypeMap).forEach(chartType => {
        let seriesData = []
        let _seriesData = {}
        Object.values(this.chartTypeMap[chartType]).forEach(seriesName =>{
          _seriesData[this.$t('common.Today')+this.$t('report.realtime.'+chartType+'.'+seriesName)] = this.tData[seriesName]
          _seriesData[this.$t('common.Yesterday')+this.$t('report.realtime.'+chartType+'.'+seriesName)] = this.yData[seriesName]
        })

        Object.keys(_seriesData).forEach(seriesName => {
          let _series = {
            name: seriesName,
            type: 'line',
            smooth: true,
            connectNulls: true,
            emphasis: {
              focus: 'self',
            },
            showSymbol: false,
            data: _seriesData[seriesName],
          }
          seriesData.push(_series)
        })
        this.setChartOptions(chartType, seriesData)
      })
    },
    setChartOptions(chartType, seriesData){
      let title = this.$t('report.realtime.type.' + chartType)
      this.chartOption[chartType] = {
        title: {
          text: title,
          left: 100,
        },
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          type: 'scroll',
          top: 20,
          left: 100,
        },
        xAxis: {
          type: 'category',
          data: this.timeRangeArr,
        },
        yAxis: {
            type: 'value',
            position: 'left',
            minInterval: 1,
        },
        series: seriesData,
      }
    },
    rebuildCharts(){
      this.$nextTick(()=>{
        this.$refs.active.rebuildCharts()
        this.$refs.pating.rebuildCharts()
        this.$refs.payAmountSummary.rebuildCharts()
        this.$refs.newUserCnt.rebuildCharts()
        this.$refs.newUserCntSummary.rebuildCharts()
      })
    },
  }
}
</script>

<style scoped>

</style>
