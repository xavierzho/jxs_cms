<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod"
      :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
      :height="'75vh'"
      @fetch="fetch"
      @expandChange="expandChange"
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

          <el-form-item :label="$t('activity.point_type')">
            <el-select v-model="searchForm.point_type"  clearable>       <!-- multiple  多选-->
              <el-option v-for="item in pointType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExportDetail({})">{{ $t('common.ExportDetail') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="created_at" :label="$t('common.DateTime')" min-width="90px" align="center"/>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"/>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"/>
        <el-table-column prop="period" :label="$t('activity.stepByStep.period')" min-width="90px" align="center"/>
        <el-table-column prop="point_type" :label="$t('activity.point_type')" min-width="90px" align="center">
          <template v-slot="data">{{ $t(`option.pointType[${data.row.point_type}]`) }}</template>
        </el-table-column>
        <el-table-column prop="point" :label="$t('activity.point')" min-width="90px" align="center"/>
        <el-table-column prop="step_no" :label="$t('activity.stepByStep.step_no')" min-width="90px" align="center"/>
        <el-table-column prop="cell_no" :label="$t('activity.stepByStep.cell_no')" min-width="90px" align="center"/>
        <el-table-column prop="inner_price" :label="$t('item.inner_price')" min-width="90px" align="center">
          <template v-slot="data">{{ data.row.inner_price | localeNum2f }}</template>
        </el-table-column>
        <el-table-column type="expand">
          <template v-slot="record">
            <el-table v-loading="loading" :data="record.row.detail" border stripe style="width: 100%">
              <el-table-column prop="award_type" :label="$t('activity.award.type')" min-width="90px" align="center">
                <template v-slot="data">{{ $t(`option.awardType[${data.row.award_type}]`) }}</template>
              </el-table-column>
              <el-table-column prop="award_value" :label="$t('activity.award.value')" min-width="90px" align="center"/>
              <el-table-column prop="award_name" :label="$t('activity.award.name')" min-width="90px" align="center"/>
              <el-table-column prop="award_params" :label="$t('activity.award.params')" min-width="90px" align="center"/>
              <el-table-column prop="award_num" :label="$t('activity.award.num')" min-width="90px" align="center">
                <template v-slot="data">{{ data.row.award_num | localeNum }}</template>
              </el-table-column>
              <el-table-column prop="inner_price" :label="$t('item.inner_price')" min-width="90px" align="center">
                <template v-slot="data">{{ data.row.inner_price | localeNum2f }}</template>
              </el-table-column>
            </el-table>
          </template>
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
  name: "StepByStep",
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
        point_type: '',
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
      pointType: 'option/pointType',
    }),
  },
  created() {
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
      api.getActivityStepByStep(this.searchForm).then(res => {
        this.tableData = res.data || []
        this.summary = res.headers.summary || {}
        this.total = res.headers.total || 0
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
      api.exportActivityStepByStep(this.searchForm).
      finally(() => {
        this.loading = false
      })
    },
    fetchDetail(params = {}, row){
      this.loading = true
      api.getActivityStepByStepDetail(params).
      then(res => {
        row.detail = res.data || []
      }).catch(()=>{
        row.detail = []
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExportDetail(params = {}){
      this.loading = true
      this.searchForm = Object.assign({}, this.searchForm, params)
      api.exportActivityStepByStepDetail(this.searchForm).
      finally(() => {
        this.loading = false
      })
    },

    changeIsAmin(value){
      if (value === ""){
        delete(this.searchForm.is_admin)
      }else{
        this.searchForm.is_admin = value
      }
    },
    expandChange(row){
      if (row.detail === undefined){
        this.fetchDetail({cell_config_id: row.id}, row)
      }
    },

    summaryMethod({ columns, data }){
      return [
        'Total',
        undefined, undefined, undefined, undefined,
        Number(this.summary['point']).toLocaleString(),
        undefined, undefined,
        Number(this.summary['inner_price']).toLocaleString(),
      ]
    },
  },
}
</script>

<style lang="scss" scoped>
</style>
