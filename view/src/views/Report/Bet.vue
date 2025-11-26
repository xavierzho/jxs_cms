<template>
  <div>
    <el-form :inline="true" :model="searchForm">
      <el-form-item v-if="$checkPermission('report_bet_filter')">
        <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
      </el-form-item>
    </el-form>

    <el-tabs type="border-card" v-loading="loading" v-model="tabActiveName">
      <el-tab-pane :label="$t('report.bet.type.FirstPrize')" name="FirstPrize">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.FirstPrize"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.bet.type.Gashapon')" name="Gashapon">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.Gashapon"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.bet.type.Chao')" name="Chao">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.Chao"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.bet.type.Hole')" name="Hole">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.Hole"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.bet.type.ChaoShe')" name="ChaoShe">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.ChaoShe"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
      <el-tab-pane :label="$t('report.bet.type.ShareBill')" name="ShareBill">
        <bet-table
          ref="betTable"
          :loading="loading"
          :tableData="tableData.ShareBill"
          :show-summary="false"
          :summary-method="summaryMethod"
        >
        </bet-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import Decimal from "decimal.js";
import BetTable from "@/components/Page/BetTable.vue";

export default {
  name: "Bet",
  components: {BetTable, DatePicker},
  data(){
    return {
      loading: false,
      tabActiveName: "FirstPrize",
      searchForm: {
        date_range: [],
      },
      tableData: {
        FirstPrize: [],
        Gashapon: [],
        Chao: [],
        Hole: [],
        ChaoShe: [],
        ShareBill: [],
      }
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
      moment().format('YYYY-MM-DD'),
    ]
    this.fetch()
  },
  methods: {
    fetch(params = {}){
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.searchForm.data_type = this.tabActiveName
      api.getReportBet(this.searchForm).then(res => {
        this.tableData[this.tabActiveName] = res.data || []
      }).catch(()=>{
        this.tableData[this.tabActiveName] = []
      }).finally(() => {
        this.loading = false
      })
    },
    summaryMethod({ columns, data }) {
      let row = Array(columns.length).fill(0)
      let sumValue = {
        "user_cnt": new Decimal(0),
        "bet_nums": new Decimal(0),
        "box_cnt_remaining": new Decimal(0),
        "box_cnt_new": new Decimal(0),
        "box_cnt_close": new Decimal(0),
        "amount": new Decimal(0),
        "amount_balance": new Decimal(0),
        "amount_wechat": new Decimal(0),
        "amount_ali": new Decimal(0),
      }
      data.forEach(item =>{
        Object.keys(item).forEach(key => {
          if (sumValue.hasOwnProperty(key)){
            sumValue[key] = sumValue[key].add(item[key] || 0)
          }
        })
      })

      columns.forEach((el, index) => {
        switch (el.property) {
          case "date":
            row[index] = "TOTAL"
            break;
          default:
            row[index] = Number(sumValue[el.property].toFixed(2, Decimal.ROUND_HALF_UP)).toLocaleString()
        }
      })

      return row
    },
  },
}
</script>

<style scoped>

</style>
