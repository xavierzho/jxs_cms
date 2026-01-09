<template>
  <div>
    <el-table ref="table" v-loading="loading" :data="tableData" height="80vh" border stripe style="width: 100%"
              :show-summary="showSummary" :summary-method="summaryMethod"
    >
      <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
      <el-table-column prop="user_cnt" :label="$t('report.bet.user_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.user_cnt | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="bet_nums" :label="$t('report.bet.bet_nums')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.bet_nums | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="box_cnt_remaining" :label="$t('report.bet.box_cnt_remaining')" min-width="90px" align="center">
        <template v-slot="data">{{formatBoxCntRemaining(data.row.date, data.row.box_cnt_remaining)}}</template>
      </el-table-column>
      <el-table-column prop="box_cnt_new" :label="$t('report.bet.box_cnt_new')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.box_cnt_new | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="box_cnt_close" :label="$t('report.bet.box_cnt_close')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.box_cnt_close | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="amount" :label="$t('option.balance_type.0')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_balance" :label="$t('option.balance_type.2')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_balance | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_jidou" :label="$t('option.balance_type.3')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_jidou | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_points" :label="$t('option.balance_type.10')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_points | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_loyalty" :label="$t('option.balance_type.11')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_loyalty | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_wechat" :label="$t('report.bet.amount_wechat')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_wechat | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_ali" :label="$t('report.bet.amount_ali')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_ali | localeNum2f }}</template>
      </el-table-column>
      <el-table-column prop="amount_huifu" :label="$t('report.bet.amount_huifu')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.amount_huifu | localeNum2f }}</template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>

import moment from "moment";

export default {
  name: "BetTable",
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
    tableData: {
      type: Array,
      default() {
        return []
      },
    },
    height: {
      type: String,
      default() {
        return '600px'
      },
    },
    showSummary: {
      type: Boolean,
      default: false,
    },
    summaryMethod: Function,
  },
  data() {
    return {}
  },
  updated() {
    this.$nextTick(() => {
      this.$refs['table'].doLayout()
    })
  },
  methods: {
    formatBoxCntRemaining(date, boxCntRemaining) {
      return date === moment().format('YYYY-MM-DD') ? Number(boxCntRemaining).toLocaleString() : '-'
    },
  },
}
</script>

<style scoped>

</style>
