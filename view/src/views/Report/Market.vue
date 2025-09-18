<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod"
      :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
      :height="'75vh'"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <el-form :inline="true" :model="searchForm">
          <el-form-item>
            <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_cnt" :label="$t('report.market.user_cnt')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.user_cnt | localeNum }}</template>
        </el-table-column>
        <el-table-column prop="order_cnt" :label="$t('report.market.order_cnt')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.order_cnt | localeNum }}</template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('report.market.amount')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.amount | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="amount_0" :label="$t('report.market.amount_0')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.amount_0 | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="amount_1" :label="$t('report.market.amount_1')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.amount_1 | localeNum2f }}</template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import SearchTable from "@/components/Table/SearchTable.vue";

export default {
  name: "Market",
  components: {SearchTable, DatePicker},
  data(){
    return {
      loading: false,
      searchForm: {
        page: 1,
        page_size: 50,
        date_range: [],
      },
      total: 0,
      tableData: [],
      summary: {},
    }
  },
  created() {
    this.searchForm.date_range = [
      moment().subtract(1, 'month').format('YYYY-MM-DD'),
      moment().format('YYYY-MM-DD'),
    ]
    this.fetch()
  },
  methods:{
    fetch(params = {}) {
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      api.getReportMarket(this.searchForm).then(res => {
        this.tableData = res.data || []
        this.summary = res.headers.summary || {}
        this.total = this.summary['total'] || 0
      }).catch(()=>{
        this.tableData = []
        this.summary = {}
        this.total = 0
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExport(params = {}){
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      api.exportReportMarket(this.searchForm).
      finally(() => {
        this.loading = false
      })
    },
    summaryMethod({ columns, data }){
      return [
        'Total',
        undefined, Number(this.summary['order_cnt']).toLocaleString(),
        Number(this.summary['amount']).toLocaleString(), Number(this.summary['amount_0']).toLocaleString(), Number(this.summary['amount_1']).toLocaleString(),
      ]
    },
  },
}
</script>

<style scoped>

</style>
