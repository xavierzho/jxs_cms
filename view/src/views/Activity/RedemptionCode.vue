<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod"
      :page="searchFormLog.page" :page-size="searchFormLog.pageSize" :total="total"
      :height="'64vh'"
      @fetch="fetch"
      @expandChange="expandChange"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchFormLog.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss"
                         :clearable="false"></date-picker>
          </el-form-item>
          <el-form-item :label="$t('user.UserID')">
            <input-number v-model="searchFormLog.user_id" clearable></input-number>
          </el-form-item>
          <el-form-item :label="$t('user.Name')">
            <el-input v-model="searchFormLog.user_name" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('user.Tel')">
            <input-number v-model="searchFormLog.tel" :range="[0,99999999999]" clearable></input-number>
          </el-form-item>
          <el-form-item :label="$t('user.IsAdmin')">
            <el-select v-model="searchFormLog.is_admin" clearable @change="changeIsAmin">
              <el-option v-for="(item, index) in userType" :key="index" :label="item.label"
                         :value="item.value"></el-option>
            </el-select>
          </el-form-item>
        </el-form>

        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('activity.redemptionCode.name')">
            <el-input v-model="searchFormLog.name" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('activity.redemptionCode.code')">
            <el-input v-model="searchFormLog.code" clearable></el-input>
          </el-form-item>
        </el-form>

        <el-form inline v-model="searchFormLog">
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{
                $t('common.Search')
              }}
            </el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExportDetail({})">
              {{ $t('common.ExportDetail') }}
            </el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="date_time" :label="$t('common.DateTime')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="name" :label="$t('activity.redemptionCode.name')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="code" :label="$t('activity.redemptionCode.code')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="reward_value_item" :label="$t('reward.reward_value_item')" min-width="90px"
                         align="center">
          <template v-slot="data">
            {{ data.row.reward_value_item | localeNum2f }}
          </template>
        </el-table-column>

        <el-table-column prop="reward_value_cost_award_point" :label="$t('reward.reward_value_cost_award_point')"
                         min-width="90px" align="center">
          <template v-slot="data">
            {{ data.row.reward_value_cost_award_point | localeNum2f }}
          </template>
        </el-table-column>

        <el-table-column type="expand">
          <template v-slot="record">
            <award-detail :loading="loading" :data="record.row.detail" show-level-name>
            </award-detail>
          </template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import {mapGetters,} from "vuex";
import InputNumber from "@/components/Input/InputNumber.vue";
import InputNumberRange from "@/components/Input/InputNumberRange.vue"
import ItemDetail from "@/components/Page/ItemDetail.vue";
import SearchTable from "@/components/Table/SearchTable.vue";
import AwardDetail from "@/components/Page/AwardDetail.vue";

export default {
  name: "Task",
  components: {AwardDetail, SearchTable, ItemDetail, InputNumber, InputNumberRange, DatePicker},
  data() {
    return {
      loading: false,
      searchFormLog: {
        page: 1,
        page_size: 50,
        date_time_range: [],
        user_id: '',
        user_name: '',
        tel: '',
        is_admin: false,
        name: '',
        code: '',
      },
      searchFormDetail: {
        log_id: 0,
      },
      tableData: [],
      total: 0,
      summary: {},
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
    }),
  },
  created() {
    this.searchFormLog.date_time_range = [
      moment().subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss'),
      moment().format('YYYY-MM-DD HH:mm:ss'),
    ]
    this.fetch()
  },
  methods: {
    fetch(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      api.getActivityRedemptionCode(searchFormLog).then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
        this.summary = res.headers.summary || {}
      }).catch(() => {
        this.tableData = []
        this.total = 0
        this.summary = {}
      }).finally(() => {
        this.loading = false
      })
    },

    fetchDetail(params = {}, row) {
      this.loading = true
      this.searchFormDetail = Object.assign({}, this.searchFormDetail, params)
      api.getActivityRedemptionCodeAwardDetail(this.searchFormDetail).then(res => {
        row.detail = res.data || []
      }).catch(() => {
        row.detail = []
      }).finally(() => {
        this.loading = false
      })
    },
    expandChange(row) {
      if (row.detail === undefined) {
        this.fetchDetail({log_id: row.log_id,}, row)
      }
    },
    summaryMethod({columns, data}) {
      return [
        'Total', undefined, undefined, undefined, undefined, Number(this.summary["reward_value_item"]).toLocaleString(),
        Number(this.summary["reward_value_cost_award_point"]).toLocaleString(),
      ]
    },
    changeIsAmin(value) {
      if (value === "") {
        delete (this.searchFormLog.is_admin)
      } else {
        this.searchFormLog.is_admin = value
      }
    },
    fetchExport(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      api.exportActivityRedemptionCodeLog(searchFormLog).finally(() => {
        this.loading = false
      })
    },
    fetchExportDetail(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      api.exportActivityRedemptionCodeAwardDetail(searchFormLog).finally(() => {
        this.loading = false
      })
    },
  },
}
</script>

<style scoped>

</style>
