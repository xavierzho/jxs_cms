<template>
    <div>
      <search-table
        :loading="loading"
        :table-data="tableData" show-summary :summary-method="summaryMethod "
        :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
        :height="'75vh'"
        @fetch="fetch"
        >
    <template v-slot:table-search>
          <el-form inline v-model="searchForm">
            <el-form-item :label="$t('common.DateTime')">
              <date-picker v-model="searchForm.date_range" type="daterange" value-format="yyyy-MM-dd" :clearable="false"></date-picker>
            </el-form-item>
                <el-form-item :label="$t('user.UserID')">
                  <el-input v-model="searchForm.user_id" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('user.Name')">
                  <el-input v-model="searchForm.user_name" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('user.Tel')">
                  <input-number v-model="searchForm.tel" :range="[0,99999999999]" clearable></input-number>
                </el-form-item>
              <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
              <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
  
          </el-form>
    </template>
            <template v-slot:table-column>
                <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
                <el-table-column prop="user_id" :label="$t('report.order.user_id')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.user_id }}</template>
                </el-table-column>
                <el-table-column prop="user_name" :label="$t('report.order.user_name')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.user_name}}</template>
                </el-table-column>
                <el-table-column prop="show_price" :label="$t('report.order.show_price')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.show_price }}</template>
                </el-table-column>
                <el-table-column prop="inner_price" :label="$t('report.order.inner_price')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.inner_price }}</template>
                </el-table-column>
                <el-table-column prop="recycling_price" :label="$t('report.order.recycling_price')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.recycling_price }}</template>
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
  import InputNumber from "@/components/Input/InputNumber.vue";
  import {mapGetters} from "vuex";
  
  export default {
    name: "Invite",
    components: {InputNumber,SearchTable, DatePicker},
    data(){
      return {
        loading: false,
        searchForm: {
          page: 1,
          page_size: 50,
          date_range: [],
          user_id: '',
          user_name: '',
          tel: '',
        },
        total: 0,
        amount: 0,
        totalamount: 0,
        difference: 0,
        tableData: [],
        summary: {},
      }
    },
  
    computed: {
      ...mapGetters({
        inviteUserType: 'option/inviteUserType',
      }),
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
        api.getReportDeliveryOrderDaily(this.searchForm).then(res => {
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
        api.exportReportDeliveryOrderDaily(this.searchForm).
        finally(() => {
          this.loading = false
        })
      },
  
  
      summaryMethod({ columns, data }){
        return [
          'Total',
          undefined, undefined, 
          Number(this.summary['show_price'] || 0).toLocaleString(),
          Number(this.summary['inner_price'] || 0).toLocaleString(),
          Number(this.summary['recycling_price'] || 0).toLocaleString(),
        ]
      },
    },
  }
  </script>
  
  <style scoped>
  
  </style>
  