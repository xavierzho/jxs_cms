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
        <el-form inline v-model="searchForm">
          <el-form-item>
            <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
          </el-form-item>
          <el-form-item>
            <el-select v-model="searchForm.balance_type"  clearable>
              <el-option label="积分" :value="10"></el-option>
              <el-option label="吉祥值" :value="11"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="get_user_cnt" :label="$t('activity.costAward.get_user_cnt')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.get_user_cnt | localeNum }}</template>
        </el-table-column>
        <el-table-column prop="get_point" :label="$t('activity.costAward.get_point')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.get_point | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="accept_user_cnt" :label="$t('activity.costAward.accept_user_cnt')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.accept_user_cnt | localeNum }}</template>
        </el-table-column>
        <el-table-column prop="accept_point" :label="$t('activity.costAward.accept_point')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.accept_point | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="award_amount" :label="$t('activity.costAward.award_amount')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.award_amount | localeNum2f }}</template>
        </el-table-column>
        <!--财务要求展示价在成本价前-->
        <el-table-column prop="award_item_show_price" :label="$t('activity.costAward.award_item_show_price')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.award_item_show_price | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="award_item_inner_price" :label="$t('activity.costAward.award_item_inner_price')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.award_item_inner_price | localeNum2f }}</template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment";
import api from "@/api";
import SearchTable from "@/components/Table/SearchTable.vue";
import Decimal from "decimal.js";

export default {
  name: "CostAward",
  components: {SearchTable, DatePicker},
  data(){
    return {
      loading: false,
      searchForm: {
        page: 1,
        page_size: 50,
        balance_type:10,
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
  methods: {
    fetch(params = {}){
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      api.getActivityCostAward(this.searchForm).then(res => {
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
      api.exportActivityCostAward(this.searchForm).
      finally(() => {
        this.loading = false
      })
    },
    summaryMethod({ columns, data }){
      return [
        'Total',
        undefined, Number(this.summary['get_point']).toLocaleString(), undefined, Number(this.summary['accept_point']).toLocaleString(),
        Number(this.summary['award_amount']).toLocaleString(), Number(this.summary['award_item_show_price']).toLocaleString(), Number(this.summary['award_item_inner_price']).toLocaleString(),
      ]
    },
  },
}
</script>

<style scoped>

</style>
