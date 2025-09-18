<template>
  <div v-if="$checkPermission('report_dashboard_view')">
    <!--form-->
    <el-form inline>
      <el-form-item>
        <el-button v-loading="loading" type="primary" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
      </el-form-item>
    </el-form>

    <!--board-->
    <div v-if="!loading" style="display: flex; flex-direction: column;">
      <div class="group">
        <dashboard-value class-name="recharge_amount" :title="$t('report.dashboard.recharge_amount')" :value="[todayData.recharge_amount,yesterdayData.recharge_amount]" value-type="value2f"/>
        <dashboard-value class-name="recharge_amount_wechat" :title="$t('report.dashboard.recharge_amount_wechat')" :value="[todayData.recharge_amount_wechat,yesterdayData.recharge_amount_wechat]" value-type="value2f"/>
        <dashboard-value class-name="recharge_amount_ali" :title="$t('report.dashboard.recharge_amount_ali')" :value="[todayData.recharge_amount_ali,yesterdayData.recharge_amount_ali]" value-type="value2f"/>
        <dashboard-value class-name="draw_amount" :title="$t('report.dashboard.draw_amount')" :value="[todayData.draw_amount,yesterdayData.draw_amount]" value-type="value2f"/>
        <dashboard-value class-name="new_user_cnt" :title="$t('report.dashboard.new_user_cnt')" :value="[todayData.new_user_cnt,yesterdayData.new_user_cnt]"/>
      </div>

      <div class="group">
        <dashboard-value class-name="active_user_cnt" :title="$t('report.dashboard.active_user_cnt')" :value="[todayData.active_user_cnt,yesterdayData.active_user_cnt]"/>
        <dashboard-value class-name="recharge_user_cnt" :title="$t('report.dashboard.recharge_user_cnt')" :value="[todayData.recharge_user_cnt,yesterdayData.recharge_user_cnt]" :value-classify="[todayData.recharge_user_cnt_new, todayData.recharge_user_cnt-todayData.recharge_user_cnt_new]"/>
        <dashboard-value class-name="pating_rate_new" :title="$t('report.dashboard.pating_rate_new')" :value="[todayData.pating_rate_new,yesterdayData.pating_rate_new]" value-type="percentage"/>
        <dashboard-value class-name="new_user_cnt_summary" :title="$t('report.dashboard.new_user_cnt_summary')" :value="[summary.new_user_cnt,0]" is-summary/>
        <dashboard-value class-name="recharge_user_cnt_summary" :title="$t('report.dashboard.recharge_user_cnt_summary')" :value="[summary.recharge_user_cnt,0]" is-summary/>
      </div>

      <div class="group">
        <dashboard-value class-name="recharge_amount_summary" :title="$t('report.dashboard.recharge_amount_summary')" :value="[summary.recharge_amount,0]" value-type="value2f" is-summary/>
        <dashboard-value class-name="draw_amount_summary" :title="$t('report.dashboard.draw_amount_summary')" :value="[summary.draw_amount,0]" value-type="value2f" is-summary/>
      </div>
    </div>

    <!--table-->
    <el-table ref="table" v-loading="loading" :data="tableData" border stripe style="width: 100%">
      <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
      <el-table-column prop="recharge_amount" :label="$t('report.dashboard.recharge_amount')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.recharge_amount | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="recharge_amount_wechat" :label="$t('report.dashboard.recharge_amount_wechat')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.recharge_amount_wechat | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="recharge_amount_ali" :label="$t('report.dashboard.recharge_amount_ali')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.recharge_amount_ali | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="draw_amount" :label="$t('report.dashboard.draw_amount')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.draw_amount | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="new_user_cnt" :label="$t('report.dashboard.new_user_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.new_user_cnt | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="active_user_cnt" :label="$t('report.dashboard.active_user_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.active_user_cnt | localeNum }}</template>
      </el-table-column>
      <el-table-column prop="recharge_user_cnt" :label="$t('report.dashboard.recharge_user_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{ data.row.recharge_user_cnt | localeNum }}</template>
      </el-table-column>
    </el-table>
  </div>
  <div v-else class="dashboard-container">
    <div class="dashboard-text">Hello {{ name }}</div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { checkPermission } from '@/utils/auth'
import api from "@/api";
import DashboardValue from "@/components/Page/DashboardValue.vue";

export default {
  name: 'Dashboard',
  components: {DashboardValue},

  data(){
    return {
      loading: false,
      summary: {},
      tableData: [],
      todayData: [],
      yesterdayData: [],
    }
  },

  computed: {
    ...mapGetters([
      'name'
    ])
  },

  created() {
    if (checkPermission('report_dashboard_view')){
      this.fetch()
    }
  },

  methods:{
    fetch(params = {}){
      this.loading = true
      api.getReportDashboard(params).
        then(res => {
          this.tableData = res.data || []
          this.todayData = res.data[0] || []
          this.yesterdayData = res.data[1] || []
          this.summary = res.headers.summary
        }).catch(()=>{
          this.tableData = []
          this.todayData = []
          this.yesterdayData = []
          this.summary = {}
        }).finally(() => {
          this.loading = false
        })
    },
  },
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}

.group{
  padding: 5px 5px 20px;
  display: flex;
  overflow: auto;
}
</style>
