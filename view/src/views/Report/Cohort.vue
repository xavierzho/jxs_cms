<template>
  <div>
    <el-form :inline="true" :model="searchForm">
      <el-form-item v-if="$checkPermission('report_cohort_filter')">
        <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
      </el-form-item>
    </el-form>

    <el-tabs type="border-card" v-loading="loading" v-model="tabActiveName">
      <el-tab-pane :label="$t('report.cohort.type.new_user_active')" name="new_user_active">
        <cohort-table ref="new_user_active" :loading="loading"
                      :tableData="tableData.new_user_active" :chartOption="chartOption.new_user_active"
                      :summary-method="summaryMethod"
        >
        </cohort-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.cohort.type.new_user_validated')" name="new_user_validated">
        <cohort-validated-user-table ref="new_user_validated" :loading="loading"
                                     :tableData="tableData.new_user_validated" :chartOption="chartOption.new_user_validated"
                                     :summary-method="summaryMethod"
        >
        </cohort-validated-user-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.cohort.type.new_user_consume')" name="new_user_consume">
        <cohort-table ref="new_user_consume" :loading="loading"
                      :tableData="tableData.new_user_consume" :chartOption="chartOption.new_user_consume"
                      :summary-method="summaryMethod"
        >
        </cohort-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.cohort.type.pating_user_active')" name="pating_user_active">
        <cohort-table ref="pating_user_active" :loading="loading"
                      :tableData="tableData.pating_user_active" :chartOption="chartOption.pating_user_active"
                      :summary-method="summaryMethod"
        >
        </cohort-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.cohort.type.consume_user_active')" name="consume_user_active">
        <cohort-table ref="consume_user_active" :loading="loading"
                      :tableData="tableData.consume_user_active" :chartOption="chartOption.consume_user_active"
                      :summary-method="summaryMethod"
        >
        </cohort-table>
      </el-tab-pane>
<!--      <el-tab-pane :label="$t('report.cohort.type.invited_new_user_active')" name="invited_new_user_active">-->
<!--        <cohort-table ref="invited_new_user_active" :loading="loading"-->
<!--                      :tableData="tableData.invited_new_user_active" :chartOption="chartOption.invited_new_user_active"-->
<!--                      :summary-method="summaryMethod"-->
<!--        >-->
<!--        </cohort-table>-->
<!--      </el-tab-pane>-->
<!--      <el-tab-pane :label="$t('report.cohort.type.invited_new_user_consume')" name="invited_new_user_consume">-->
<!--        <cohort-table ref="invited_new_user_consume" :loading="loading"-->
<!--                      :tableData="tableData.invited_new_user_consume" :chartOption="chartOption.invited_new_user_consume"-->
<!--                      :summary-method="summaryMethod"-->
<!--        >-->
<!--        </cohort-table>-->
<!--      </el-tab-pane>-->
    </el-tabs>
  </div>
</template>

<script>
import api from '@/api'
import moment from 'moment'
import DatePicker from '@/components/DatePicker.vue'
import CohortTable from "@/components/Page/CohortTable.vue";
import CohortValidatedUserTable from "@/components/Page/CohortValidatedUserTable.vue";
import Decimal from "decimal.js";
export default {
  name: "Cohort",
  components: {
    CohortValidatedUserTable,
    CohortTable,
    DatePicker
  },
  data() {
    return {
      loading: false,
      tabActiveName: "new_user_active",
      searchForm: {
        date_range: [],
        data_type: "",
      },
      tableData: {
        new_user_active: [],
        new_user_validated: [],
        new_user_consume: [],
        pating_user_active: [],
        consume_user_active: [],
        invited_new_user_active: [],
        invited_new_user_consume: [],
      },
      count:{
        new_user_active: 0,
        new_user_validated: 0,
        new_user_consume: 0,
        pating_user_active: 0,
        consume_user_active: 0,
        invited_new_user_active: 0,
        invited_new_user_consume: 0,
      },
      legendTotal: [this.$t('report.cohort.total')],
      legendData: [
        this.$t('report.cohort.first_day'),this.$t('report.cohort.second_day'),this.$t('report.cohort.third_day'),
        this.$t('report.cohort.fourth_day'),this.$t('report.cohort.fifth_day'),this.$t('report.cohort.sixth_day'),
        this.$t('report.cohort.seventh_day'),this.$t('report.cohort.fourteenth_day'),this.$t('report.cohort.thirtieth_day'),
        this.$t('report.cohort.sixtieth_day'),this.$t('report.cohort.ninety_day'),this.$t('report.cohort.no_180_day'),
      ],
      legendRateData: [
        this.$t('report.cohort.first_day_rate'),this.$t('report.cohort.second_day_rate'),this.$t('report.cohort.third_day_rate'),
        this.$t('report.cohort.fourth_day_rate'),this.$t('report.cohort.fifth_day_rate'),this.$t('report.cohort.sixth_day_rate'),
        this.$t('report.cohort.seventh_day_rate'),this.$t('report.cohort.fourteenth_day_rate'),this.$t('report.cohort.thirtieth_day_rate'),
        this.$t('report.cohort.sixtieth_day_rate'),this.$t('report.cohort.ninety_day_rate'),this.$t('report.cohort.no_180_day_rate'),
      ],
      validatedLegendTotal: [this.$t('report.cohort.new_user_validated.total')],
      validatedLegendData: [
        this.$t('report.cohort.new_user_validated.first_day'),this.$t('report.cohort.new_user_validated.second_day'),
        this.$t('report.cohort.new_user_validated.third_day'),this.$t('report.cohort.new_user_validated.fourth_day'),
        this.$t('report.cohort.new_user_validated.fifth_day'),this.$t('report.cohort.new_user_validated.sixth_day'),
        this.$t('report.cohort.new_user_validated.seventh_day'),this.$t('report.cohort.new_user_validated.fourteenth_day'),
        this.$t('report.cohort.new_user_validated.thirtieth_day'),this.$t('report.cohort.new_user_validated.sixtieth_day'),
        this.$t('report.cohort.new_user_validated.ninety_day'),this.$t('report.cohort.new_user_validated.no_180_day'),
      ],
      validatedLegendRateData: [
        this.$t('report.cohort.new_user_validated.first_day_rate'),this.$t('report.cohort.new_user_validated.second_day_rate'),
        this.$t('report.cohort.new_user_validated.third_day_rate'),this.$t('report.cohort.new_user_validated.fourth_day_rate'),
        this.$t('report.cohort.new_user_validated.fifth_day_rate'),this.$t('report.cohort.new_user_validated.sixth_day_rate'),
        this.$t('report.cohort.new_user_validated.seventh_day_rate'),this.$t('report.cohort.new_user_validated.fourteenth_day_rate'),
        this.$t('report.cohort.new_user_validated.thirtieth_day_rate'),this.$t('report.cohort.new_user_validated.sixtieth_day_rate'),
        this.$t('report.cohort.new_user_validated.ninety_day_rate'),this.$t('report.cohort.new_user_validated.no_180_day_rate'),
      ],
      selected: {
        [this.$t('report.cohort.total')]: true,
        [this.$t('report.cohort.new_user_validated.total')]: true,

        [this.$t('report.cohort.first_day')]: false,
        [this.$t('report.cohort.second_day')]: false,
        [this.$t('report.cohort.third_day')]: false,
        [this.$t('report.cohort.fourth_day')]: false,
        [this.$t('report.cohort.fifth_day')]: false,
        [this.$t('report.cohort.sixth_day')]: false,
        [this.$t('report.cohort.seventh_day')]: false,
        [this.$t('report.cohort.fourteenth_day')]: false,
        [this.$t('report.cohort.thirtieth_day')]: false,
        [this.$t('report.cohort.sixtieth_day')]: false,
        [this.$t('report.cohort.ninety_day')]: false,
        [this.$t('report.cohort.no_180_day')]: false,

        [this.$t('report.cohort.new_user_validated.first_day')]: false,
        [this.$t('report.cohort.new_user_validated.second_day')]: false,
        [this.$t('report.cohort.new_user_validated.third_day')]: false,
        [this.$t('report.cohort.new_user_validated.fourth_day')]: false,
        [this.$t('report.cohort.new_user_validated.fifth_day')]: false,
        [this.$t('report.cohort.new_user_validated.sixth_day')]: false,
        [this.$t('report.cohort.new_user_validated.seventh_day')]: false,
        [this.$t('report.cohort.new_user_validated.fourteenth_day')]: false,
        [this.$t('report.cohort.new_user_validated.thirtieth_day')]: false,
        [this.$t('report.cohort.new_user_validated.sixtieth_day')]: false,
        [this.$t('report.cohort.new_user_validated.ninety_day')]: false,
        [this.$t('report.cohort.new_user_validated.no_180_day')]: false,

        [this.$t('report.cohort.first_day_rate')]: false,
        [this.$t('report.cohort.second_day_rate')]: false,
        [this.$t('report.cohort.third_day_rate')]: false,
        [this.$t('report.cohort.fourth_day_rate')]: false,
        [this.$t('report.cohort.fifth_day_rate')]: false,
        [this.$t('report.cohort.sixth_day_rate')]: false,
        [this.$t('report.cohort.seventh_day_rate')]: false,
        [this.$t('report.cohort.fourteenth_day_rate')]: false,
        [this.$t('report.cohort.thirtieth_day_rate')]: false,
        [this.$t('report.cohort.sixtieth_day_rate')]: false,
        [this.$t('report.cohort.ninety_day_rate')]: false,
        [this.$t('report.cohort.no_180_day_rate')]: false,

        [this.$t('report.cohort.new_user_validated.first_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.second_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.third_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.fourth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.fifth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.sixth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.seventh_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.fourteenth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.thirtieth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.sixtieth_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.ninety_day_rate')]: false,
        [this.$t('report.cohort.new_user_validated.no_180_day_rate')]: false,
      },
      chartOption:{
        new_user_active: {},
        new_user_validated: {},
        new_user_consume: {},
        pating_user_active: {},
        consume_user_active: {},
        invited_new_user_active: {},
        invited_new_user_consume: {},
      },
    }
  },
  watch: {
    tabActiveName(name) {
      if (this.tableData[name].length === 0) {
        this.fetch()
      }
    },
  },
  created() {
    this.searchForm.date_range = [
      moment().subtract(1, 'month').format('YYYY-MM-DD'),
      moment().subtract(1, 'day').format('YYYY-MM-DD'),
    ]
    this.fetch()
  },
  computed: {
    timeRangeArr() {
      let timeArr = []
      let start = moment(this.searchForm.date_range[0], 'YYYY-MM-DD')
      let end = moment(this.searchForm.date_range[1], 'YYYY-MM-DD')
      while (start <= end) {
        timeArr.push(start.format('YYYY-MM-DD'))
        start = start.add(1, 'day')
      }
      return timeArr
    },
  },
  methods: {
    fetch(params = {}) {
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.searchForm.data_type = this.tabActiveName
      api.getReportCohort(this.searchForm).then(res => {
        this.tableData[this.tabActiveName] = res.data || []
        this.count[this.tabActiveName] = res.headers.user_cnt || 0
      }).catch(()=>{
        this.tableData[this.tabActiveName] = []
        this.count[this.tabActiveName] = 0
      }).finally(() => {
        this.parseData(this.tableData[this.tabActiveName])
        this.$refs[this.tabActiveName].rebuildCharts()
        this.loading = false
      })
    },
    parseData(dataList) {
      let seriesData = []
      let _seriesData = {}
      Object.values(dataList).forEach(dataItem => { // 每条数据
        Object.keys(dataItem).forEach(key => {
          if (key === 'date' || key === 'data_type' || key === 'created_at' || key === 'updated_at') {
            return
          }

          if (!_seriesData.hasOwnProperty(key)){
            _seriesData[key] = []
          }
          _seriesData[key].push(dataItem[key] || 0)
        })
      })
      Object.keys(_seriesData).forEach(seriesName => {
        let name = this.$t('report.cohort.' + seriesName)
        if (this.tabActiveName === 'new_user_validated'){
          name = this.$t('report.cohort.new_user_validated.' + seriesName)
        }
        let _series = {
          name: name,
          type: 'line',
          smooth: true,
          connectNulls: true,
          emphasis: {
            focus: 'self',
          },
          showSymbol: false,
          data: _seriesData[seriesName].reverse(),
        }
        if (seriesName.includes('rate')){
          _series['yAxisIndex'] = 1
        }
        seriesData.push(_series)
      })
      this.setChartOptions(seriesData)
    },
    setChartOptions(seriesData) {
      let title = this.$t('report.cohort.type.' + this.tabActiveName)
      let timeRangeArr = this.timeRangeArr
      let legendTotal = this.legendTotal
      let legendData = this.legendData
      let legendRateData = this.legendRateData
      if (this.tabActiveName === 'new_user_validated'){
        timeRangeArr = timeRangeArr.slice(0, this.tableData['new_user_validated'].length)
        legendTotal = this.validatedLegendTotal
        legendData = this.validatedLegendData
        legendRateData = this.validatedLegendRateData
      }
      this.chartOption[this.tabActiveName] = {
        title: {
          text: title,
          left: 100,
        },
        tooltip: {
          trigger: 'axis',
        },
        legend: [
          {
            data: legendTotal,
            type: 'scroll',
            top: 10,
            selected: this.selected,
          },
          {
            data: legendData,
            type: 'scroll',
            orient: 'vertical',
            left: 10,
            top: 20,
            bottom: 20,
            selected: this.selected,
          },
          {
            data: legendRateData,
            type: 'scroll',
            orient: 'vertical',
            right: 10,
            top: 20,
            bottom: 20,
            selected: this.selected,
          },
        ],
        xAxis: {
          type: 'category',
          data: timeRangeArr,
        },
        yAxis: [
          {
            type: 'value',
            position: 'left',
            minInterval: 1,
          },
          {
            type: 'value',
            position: 'right',
            axisLabel: {
              formatter: '{value}%'
            }
          },
        ],
        series: seriesData,
      }
    },
    summaryMethod({ columns, data }) {
      let row = Array(columns.length).fill("0%")
      let totalRegisterNum = new Decimal(0)
      let offset = {
        first_day: new Decimal(0),
        second_day: new Decimal(0),
        third_day: new Decimal(0),
        fourth_day: new Decimal(0),
        fifth_day: new Decimal(0),
        sixth_day: new Decimal(0),
        seventh_day: new Decimal(0),
        fourteenth_day: new Decimal(0),
        thirtieth_day: new Decimal(0),
        sixtieth_day: new Decimal(0),
        ninety_day: new Decimal(0),
        no_180_day: new Decimal(0),
      }
      let baseNum = {
        first_day: new Decimal(0),
        second_day: new Decimal(0),
        third_day: new Decimal(0),
        fourth_day: new Decimal(0),
        fifth_day: new Decimal(0),
        sixth_day: new Decimal(0),
        seventh_day: new Decimal(0),
        fourteenth_day: new Decimal(0),
        thirtieth_day: new Decimal(0),
        sixtieth_day: new Decimal(0),
        ninety_day: new Decimal(0),
        no_180_day: new Decimal(0),
      }

      let order = {
        first_day: 0,
        second_day: 1,
        third_day: 2,
        fourth_day: 3,
        fifth_day: 4,
        sixth_day: 5,
        seventh_day: 6,
        fourteenth_day: 13,
        thirtieth_day: 29,
        sixtieth_day: 59,
        ninety_day: 89,
        no_180_day: 179,
      }
      let nowDate = moment(moment().format('YYYY-MM-DD'))
      for (let i = 0; i < data.length; i++) {
        if (data[i].hasOwnProperty('total')) {
          totalRegisterNum = totalRegisterNum.add(data[i].total || 0)
          // 若该日期在今日之后 判断该行的注册人数是否为偏差值
          Object.keys(order).forEach(key => {
            if (moment(data[i].date).add(order[key], 'days').isAfter(nowDate)) {
              offset[key] = offset[key].add(new Decimal(data[i].total || 0))
            }
          })
        }
        Object.keys(order).forEach(key => {
          if (data[i].hasOwnProperty(key)) {
            baseNum[key] = baseNum[key].add(new Decimal(data[i][key] || 0))
          }
        })
      }
      for (let i = 0; i < columns.length; i++) {
        if (offset.hasOwnProperty(columns[i].property)) {
          if (!totalRegisterNum.sub(offset[columns[i].property]).equals(0)) {
            let value = baseNum[columns[i].property].div(totalRegisterNum.sub(offset[columns[i].property])).mul(100).toFixed(2)
            row[i] = Number(value).toLocaleString() + '%'
          } else {
            row[i] = '0%'
          }
        }
      }
      row[0] = 'TOTAL'
      row[1] = Number(this.count[this.tabActiveName]).toLocaleString()

      return row
    },
  },
}
</script>

