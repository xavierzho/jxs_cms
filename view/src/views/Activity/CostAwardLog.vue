<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod" :cellClassName="cellClassName"
      :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
      :height="'75vh'"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchForm">
          <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchForm.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss" :clearable="false"></date-picker>
          </el-form-item>

          <el-form-item :label="$t('user.UserID')">
            <input-number v-model="searchForm.user_id" clearable></input-number>
          </el-form-item>
          <el-form-item :label="$t('user.Name')">
            <el-input v-model="searchForm.user_name" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('user.Tel')">
            <input-number v-model="searchForm.tel" :range="[0,99999999999]" clearable></input-number>
          </el-form-item>
          <el-form-item :label="$t('user.IsAdmin')">
            <el-select v-model="searchForm.is_admin" clearable @change="changeIsAmin">
              <el-option v-for="(item, index) in userType" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('user.UserChannel')">
            <el-select v-model="searchForm.channel" clearable>
              <el-option v-for="(item, index) in userChannel" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('activity.costAwardLog.log_type')">
            <el-select v-model="searchForm.log_type" multiple clearable>
              <el-option v-for="item in activityCostAwardLogType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="created_at" :label="$t('common.DateTime')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="log_type_str" :label="$t('activity.costAwardLog.log_type')" min-width="90px" align="center"></el-table-column>
        <el-table-column :label="$t('activity.costAwardLog.point')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.before_point | localeNum2f}} -> {{data.row.after_point | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="update_point" :label="$t('activity.costAwardLog.update_point')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.update_point | localeNum2f }}</template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import SearchTable from "@/components/Table/SearchTable.vue";
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import {mapGetters} from "vuex";
import InputNumber from "@/components/Input/InputNumber.vue";

export default {
  name: "CostAwardLog",
  components: {InputNumber, SearchTable, DatePicker},
  data(){
    return {
      loading: false,
      searchForm: {
        page: 1,
        page_size: 50,
        date_time_range: [],
        user_id: '',
        user_name: '',
        tel: '',
        is_admin: false,
        channel: null,
        log_type: [],
      },
      total: 0,
      tableData: [],
      summary: {},
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
      userChannel: 'option/userChannel',
      activityCostAwardLogType: 'option/activityCostAwardLogType',
    }),
  },
  created() {
    if (this.activityCostAwardLogType.length === 0) {
      this.$store.dispatch('option/update', 'activityCostAwardLogType')
    }
    this.searchForm.date_time_range = [
      moment().subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss'),
      moment().format('YYYY-MM-DD HH:mm:ss'),
    ]
    this.fetch()
  },
  methods: {
    fetch(params = {}){
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      api.getActivityCostAwardLog(this.searchForm).then(res => {
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
      api.exportActivityCostAwardLog(this.searchForm).
      finally(() => {
        this.loading = false
      })
    },
    cellClassName(cell){
      if (cell.column.property !== 'update_point') {
        return ''
      }
      if (cell.row.update_point > 0){
        return 'value-green'
      }else if(cell.row.update_point < 0){
        return 'value-red'
      }
      return ''
    },
    summaryMethod({ columns, data }){
      return [
        'Total',
        Number(this.summary['user_cnt']).toLocaleString(), undefined, undefined,
        undefined, Number(this.summary['update_point']).toLocaleString(),
      ]
    },
    changeIsAmin(value){
      if (value === ""){
        delete(this.searchForm.is_admin)
      }else{
        this.searchForm.is_admin = value
      }
    },
  },
}
</script>

<style lang="scss" scoped>
::v-deep .value-red{
  color: #F56C6C;
}

::v-deep .value-green{
  color: #67C23A;
}

</style>
