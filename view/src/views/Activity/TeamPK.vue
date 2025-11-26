<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData"
      :page="searchFormLog.page" :page-size="searchFormLog.pageSize" :total="total"
      :height="'70vh'" :row-style="{height: '75px'}"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchFormLog.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss" :clearable="false"></date-picker>
          </el-form-item>
        </el-form>

        <el-form inline v-model="searchFormLog">
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="team_id" :label="$t('activity.teamPK.team_id')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="team_name" :label="$t('activity.teamPK.team_name')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="team_amount" :label="$t('activity.teamPK.team_amount')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.team_amount | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="team_rate" :label="$t('activity.teamPK.team_rate')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="team_point" :label="$t('activity.teamPK.team_point')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.team_point | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="team_no" :label="$t('activity.teamPK.team_no')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="team_gap" :label="$t('activity.teamPK.team_gap')" min-width="90px" align="center"></el-table-column>
        <el-table-column type="expand">
          <template v-slot="record">
            <el-table v-loading="loading" :data="record.row.user" border stripe>
              <el-table-column prop="user_id" :label="$t('activity.teamPK.user_id')" min-width="90px" align="center"></el-table-column>
              <el-table-column prop="user_name" :label="$t('activity.teamPK.user_name')" min-width="90px" align="center"></el-table-column>
              <el-table-column prop="user_amount" :label="$t('activity.teamPK.user_amount')" min-width="90px" align="center">
                <template v-slot="data">{{data.row.user_amount | localeNum2f}}</template>
              </el-table-column>
              <el-table-column prop="user_rate" :label="$t('activity.teamPK.user_rate')" min-width="90px" align="center"></el-table-column>
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
import moment from "moment";
import api from "@/api";

export default {
  name: "TeamPK",
  components: {SearchTable, DatePicker},
  data(){
    return {
      loading: false,
      searchFormLog: {
        page: 1,
        page_size: 20,
        date_time_range: [],
      },
      tableData: [],
      total: 0,
    }
  },
  computed: {
  },
  created() {
    this.searchFormLog.date_time_range = [
      moment().subtract(1, 'day').format('YYYY-MM-DD HH:mm:ss'),
      moment().format('YYYY-MM-DD HH:mm:ss'),
    ]
    this.fetch()
  },
  methods:{
    fetch(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      api.getActivityTeamPK(this.searchFormLog).
      then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
      }).catch(()=>{
        this.tableData = []
        this.total = 0
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExport(params = {}){
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      api.exportActivityTeamPK(this.searchFormLog).
      finally(() => {
        this.loading = false
      })
    },
  },
}
</script>

<style scoped>

</style>
