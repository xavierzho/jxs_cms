<template>
  <div>
    <el-form :inline="true" :model="searchForm">
      <el-form-item v-if="$checkPermission('report_revenue_filter')">
        <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
      </el-form-item>
    </el-form>

    <el-tabs type="border-card" v-loading="loading" v-model="tabActiveName">
      <el-tab-pane :label="$t('report.revenue.type.active')" name="active">
        <chart-table ref="active" :loading="loading" :tableData="tableData.active" :chartOption="chartOption.active">
          <template v-slot:table-column>
            <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
<!--            <el-table-column prop="activate_cnt" :label="$t('report.revenue.active.activate_cnt')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{ data.row.activate_cnt | localeNum }}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="activate_cnt_new" :label="$t('report.revenue.active.activate_cnt_new')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{ data.row.activate_cnt_new | localeNum }}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="register_cnt" :label="$t('report.revenue.active.register_cnt')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.register_cnt | localeNum }}</template>
            </el-table-column>
<!--            <el-table-column prop="register_cnt_rate" :label="$t('report.revenue.active.register_cnt_rate')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.register_cnt_rate | rate}}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="active_cnt" :label="$t('report.revenue.active.active_cnt')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.active_cnt | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="active_cnt_new" :label="$t('report.revenue.active.active_cnt_new')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.active_cnt_new | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="active_cnt_old" :label="$t('report.revenue.active.active_cnt_old')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.active_cnt_old | localeNum }}</template>
            </el-table-column>
<!--            <el-table-column prop="max_online_cnt" :label="$t('report.revenue.active.max_online_cnt')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{ data.row.max_online_cnt | localeNum }}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="validated_cnt_7" :label="$t('report.revenue.active.validated_cnt_7')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.validated_cnt_7 | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="validated_cnt_15" :label="$t('report.revenue.active.validated_cnt_15')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.validated_cnt_15 | localeNum }}</template>
            </el-table-column>
<!--            <el-table-column prop="per_online_time" :label="$t('report.revenue.active.per_online_time')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.per_online_time | localeNum2f}}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="wastage_rate_1" :label="$t('report.revenue.active.wastage_rate_1')" min-width="90px" align="center">
              <template v-slot="data">
                <template v-if="data.row.date >= wastageCalDate">-</template>
                <template v-else>{{data.row.wastage_rate_1 | rate}}</template>
              </template>
            </el-table-column>
            <el-table-column prop="wastage_rate_3" :label="$t('report.revenue.active.wastage_rate_3')" min-width="90px" align="center">
              <template v-slot="data">
                <template v-if="data.row.date >= wastageCalDate">-</template>
                <template v-else>{{data.row.wastage_rate_3 | rate}}</template>
              </template>
            </el-table-column>
          </template>
        </chart-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.revenue.type.pating')" name="pating">
        <chart-table ref="pating" :loading="loading" :tableData="tableData.pating" :chartOption="chartOption.pating">
          <template v-slot:table-column>
            <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
            <el-table-column prop="user_cnt" :label="$t('report.revenue.pating.user_cnt')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="rate" :label="$t('report.revenue.pating.rate')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.rate | rate}}</template>
            </el-table-column>
            <el-table-column prop="user_cnt_new" :label="$t('report.revenue.pating.user_cnt_new')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt_new | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="rate_new" :label="$t('report.revenue.pating.rate_new')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.rate_new | rate}}</template>
            </el-table-column>
          </template>
        </chart-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.revenue.type.pay')" name="pay">
        <chart-table ref="pay" :loading="loading" :tableData="tableData.pay" :chartOption="chartOption.pay">
          <template v-slot:table-column>
            <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
            <el-table-column prop="amount" :label="$t('report.revenue.pay.amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="amount_new" :label="$t('report.revenue.pay.amount_new')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.amount_new | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="amount_old" :label="$t('report.revenue.pay.amount_old')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.amount_old | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="user_cnt" :label="$t('report.revenue.pay.user_cnt')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="user_cnt_new" :label="$t('report.revenue.pay.user_cnt_new')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt_new | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="user_cnt_old" :label="$t('report.revenue.pay.user_cnt_old')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt_old | localeNum }}</template>
            </el-table-column>
<!--            <el-table-column prop="user_cnt_first" :label="$t('report.revenue.pay.user_cnt_first')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{ data.row.user_cnt_first | localeNum }}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="arpu" :label="$t('report.revenue.pay.arpu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.arpu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="arppu" :label="$t('report.revenue.pay.arppu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.arppu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="arppu_new" :label="$t('report.revenue.pay.arppu_new')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.arppu_new | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="arppu_old" :label="$t('report.revenue.pay.arppu_old')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.arppu_old | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="permeability" :label="$t('report.revenue.pay.permeability')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.permeability | rate}}</template>
            </el-table-column>
            <el-table-column prop="permeability_new" :label="$t('report.revenue.pay.permeability_new')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.permeability_new | rate}}</template>
            </el-table-column>
            <el-table-column prop="permeability_old" :label="$t('report.revenue.pay.permeability_old')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.permeability_old | rate}}</template>
            </el-table-column>
          </template>
        </chart-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.revenue.type.draw')" name="draw">
        <chart-table ref="draw" :loading="loading" :tableData="tableData.draw" :chartOption="chartOption.draw">
          <template v-slot:table-column>
            <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
            <el-table-column prop="amount" :label="$t('report.revenue.draw.amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="user_cnt" :label="$t('report.revenue.draw.user_cnt')" min-width="90px" align="center">
              <template v-slot="data">{{ data.row.user_cnt | localeNum }}</template>
            </el-table-column>
            <el-table-column prop="per_amount" :label="$t('report.revenue.draw.per_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.per_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="rate" :label="$t('report.revenue.draw.rate')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.rate | rate}}</template>
            </el-table-column>
<!--            <el-table-column prop="tax" :label="$t('report.revenue.draw.tax')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.tax | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="tax_new" :label="$t('report.revenue.draw.tax_new')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.tax_new | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="tax_old" :label="$t('report.revenue.draw.tax_old')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.tax_old | localeNum2f}}</template>-->
<!--            </el-table-column>-->
          </template>
        </chart-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.revenue.type.summary')" name="summary">
        <chart-table ref="summary" :loading="loading" :tableData="tableData.summary" :chartOption="chartOption.summary">
          <template v-slot:table-column>
            <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
            <el-table-column prop="wallet_balance" :label="$t('report.revenue.summary.wallet_balance')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.wallet_balance | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="gold_balance" :label="$t('report.revenue.summary.gold_balance')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.gold_balance | localeNum2f}}</template>
            </el-table-column>
<!--            <el-table-column prop="merchant_balance" :label="$t('report.revenue.summary.merchant_balance')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.merchant_balance | localeNum2f}}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="pay_amount" :label="$t('report.revenue.summary.pay_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.pay_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="pay_amount_bet" :label="$t('report.revenue.summary.pay_amount_bet')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.pay_amount_bet | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="recharge_amount" :label="$t('report.revenue.summary.recharge_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="discount_amount" :label="$t('report.revenue.summary.discount_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.discount_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_amount" :label="$t('report.revenue.summary.saving_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="recharge_refund_amount" :label="$t('report.revenue.summary.recharge_refund_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_refund_amount | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_refund_amount" :label="$t('report.revenue.summary.saving_refund_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_refund_amount | localeNum2f}}</template>
            </el-table-column>

<!--            <el-table-column prop="recharge_amount_wechat" :label="$t('report.revenue.summary.recharge_amount_wechat')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.recharge_amount_wechat | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="discount_amount_wechat" :label="$t('report.revenue.summary.discount_amount_wechat')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.discount_amount_wechat | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="saving_amount_wechat" :label="$t('report.revenue.summary.saving_amount_wechat')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.saving_amount_wechat | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="recharge_refund_amount_wechat" :label="$t('report.revenue.summary.recharge_refund_amount_wechat')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.recharge_refund_amount_wechat | localeNum2f}}</template>-->
<!--            </el-table-column>-->
<!--            <el-table-column prop="saving_refund_amount_wechat" :label="$t('report.revenue.summary.saving_refund_amount_wechat')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.saving_refund_amount_wechat | localeNum2f}}</template>-->
<!--            </el-table-column>-->

            <el-table-column prop="recharge_amount_ali" :label="$t('report.revenue.summary.recharge_amount_ali')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_amount_ali | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="discount_amount_ali" :label="$t('report.revenue.summary.discount_amount_ali')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.discount_amount_ali | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_amount_ali" :label="$t('report.revenue.summary.saving_amount_ali')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_amount_ali | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="recharge_refund_amount_ali" :label="$t('report.revenue.summary.recharge_refund_amount_ali')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_refund_amount_ali | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_refund_amount_ali" :label="$t('report.revenue.summary.saving_refund_amount_ali')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_refund_amount_ali | localeNum2f}}</template>
            </el-table-column>

            <el-table-column prop="recharge_amount_huifu" :label="$t('report.revenue.summary.recharge_amount_huifu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_amount_huifu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="discount_amount_huifu" :label="$t('report.revenue.summary.discount_amount_huifu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.discount_amount_huifu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_amount_huifu" :label="$t('report.revenue.summary.saving_amount_huifu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_amount_huifu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="recharge_refund_amount_huifu" :label="$t('report.revenue.summary.recharge_refund_amount_huifu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.recharge_refund_amount_huifu | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="saving_refund_amount_huifu" :label="$t('report.revenue.summary.saving_refund_amount_huifu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.saving_refund_amount_huifu | localeNum2f}}</template>
            </el-table-column>

            <el-table-column prop="draw_amount" :label="$t('report.revenue.summary.draw_amount')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.draw_amount | localeNum2f}}</template>
            </el-table-column>
<!--            <el-table-column prop="tax_amount" :label="$t('report.revenue.summary.tax_amount')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.tax_amount | localeNum2f}}</template>-->
<!--            </el-table-column>-->
            <el-table-column prop="revenue" :label="$t('report.revenue.summary.revenue')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.revenue | localeNum2f}}</template>
            </el-table-column>
            <el-table-column prop="revenue_rate" :label="$t('report.revenue.summary.revenue_rate')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.revenue_rate | rate}}</template>
            </el-table-column>
            <el-table-column prop="revenue_arpu" :label="$t('report.revenue.summary.revenue_arpu')" min-width="90px" align="center">
              <template v-slot="data">{{data.row.revenue_arpu | localeNum2f}}</template>
            </el-table-column>
<!--            <el-table-column prop="refund_amount" :label="$t('report.revenue.summary.refund_amount')" min-width="90px" align="center">-->
<!--              <template v-slot="data">{{data.row.refund_amount | localeNum2f}}</template>-->
<!--            </el-table-column>-->
          </template>
        </chart-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import api from '@/api'
import moment from 'moment'
import DatePicker from '@/components/DatePicker.vue'
import ChartTable from '@/components/Table/ChartTable.vue'
export default {
  name: "Revenue",
  components: {
    DatePicker, ChartTable
  },
  data() {
    return {
      wastageCalDate: moment().subtract(6, 'day').format('YYYY-MM-DD'),
      loading: false,
      tabActiveName: "active",
      searchForm: {
        date_range: [],
        data_type: "",
      },
      tableData: {
        active: [],
        pating: [],
        pay: [],
        draw: [],
        summary: [],
      },
      selected: {},
      chartOption:{
        active: {},
        pating: {},
        pay: {},
        draw: {},
        summary: {},
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
      moment().subtract(1, 'months').format('YYYY-MM-DD'),
      moment().format('YYYY-MM-DD'),
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
      api.getReportRevenue(this.searchForm).then(res => {
        this.tableData[this.tabActiveName] = res.data || []
      }).catch(()=>{
        this.tableData[this.tabActiveName] = []
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
          if (key === 'date') {
            return
          }
          if (!_seriesData.hasOwnProperty(key)){
            _seriesData[key] = []
          }
          _seriesData[key].push(dataItem[key] || 0)
        })
      })

      Object.keys(_seriesData).forEach(seriesName => {
        let _series = {
          name: this.$t('report.revenue.' + this.tabActiveName + '.' + seriesName),
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
        this.selected[_series.name] = ['active_cnt', 'user_cnt', 'amount', 'revenue'].includes(seriesName);
      })
      this.setChartOptions(seriesData)
    },
    setChartOptions(seriesData) {
      let title = this.$t('report.revenue.type.' + this.tabActiveName)
      this.chartOption[this.tabActiveName] = {
        title: {
          text: title,
          left: 100,
        },
        tooltip: {
          trigger: 'axis',
        },
        legend: {
          type: 'scroll',
          orient: 'vertical',
          left: 10,
          top: 20,
          bottom: 20,
          selected: this.selected,
        },
        xAxis: {
          type: 'category',
          data: this.timeRangeArr,
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
  },
}
</script>

<style scoped>

</style>
