<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData"
      :height="'64vh'"
      :total="tableData.length"
      class="hide-pagination"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <div v-perm="'inquire_revenue_item_filter'">
          <el-form inline>
            <el-form-item>
            <span slot="label">{{ $t('user.UserID') }}
              <el-tooltip effect="dark" :content="$t('user.filterOptionTip_userID')" placement="top"><i
                class='el-icon-info'/></el-tooltip>
            </span>
              <el-input v-model="searchFormLog.user_ids" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('user.Name')">
              <el-input v-model="searchFormLog.user_name" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('user.Tel')">
              <input-number v-model="searchFormLog.tel" :range="[0,99999999999]" clearable></input-number>
            </el-form-item>
            <el-form-item :label="$t('user.IsAdmin')">
              <el-select v-model="searchFormLog.is_admin" clearable @change="changeIsAmin">
                <el-option v-for="(item, index) in userType" :key="index" :label="item.label" :value="item.value"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('user.UserChannel')">
              <el-select v-model="searchFormLog.channel" clearable>
                <el-option v-for="(item, index) in userChannel" :key="index" :label="item.label" :value="item.value"></el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <el-form inline v-model="searchFormLog">
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
<!--            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>-->
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="update_amount_7_day" :label="$t('inquire.item.update_amount_7_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.update_amount_7_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="recycling_price_7_day" :label="$t('inquire.item.recycling_price_7_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.recycling_price_7_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="revenue_7_day" :label="$t('inquire.item.revenue_7_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.revenue_7_day | localeNum2f}}</template>
          </template>
        </el-table-column>

        <el-table-column prop="update_amount_15_day" :label="$t('inquire.item.update_amount_15_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.update_amount_15_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="recycling_price_15_day" :label="$t('inquire.item.recycling_price_15_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.recycling_price_15_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="revenue_15_day" :label="$t('inquire.item.revenue_15_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.revenue_15_day | localeNum2f}}</template>
          </template>
        </el-table-column>

        <el-table-column prop="update_amount_30_day" :label="$t('inquire.item.update_amount_30_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.update_amount_30_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="recycling_price_30_day" :label="$t('inquire.item.recycling_price_30_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.recycling_price_30_day | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="revenue_30_day" :label="$t('inquire.item.revenue_30_day')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{ detail.row.revenue_30_day | localeNum2f}}</template>
          </template>
        </el-table-column>

      </template>
    </search-table>
  </div>
</template>

<script>
import api from "@/api";
import {mapGetters,} from "vuex";
import InputNumber from "@/components/Input/InputNumber.vue";
import SearchTable from "@/components/Table/SearchTable.vue";

export default {
  name: "RevenueItem",
  components: {SearchTable, InputNumber},
  data(){
    return {
      loading: false,
      searchFormLog: {
        user_ids: '',
        user_name: '',
        tel: '',
        is_admin: false,
        channel: null,
      },
      tableData: [],
      total:1
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
      userChannel: 'option/userChannel',
    }),
  },
  created() {
    this.fetch()
  },
  methods:{
    fetch(params = {}){
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)

      api.getInquireRevenueItemDetailList(searchFormLog).
      then(res => {
        this.tableData = res.data || []
      }).catch(()=>{
        this.tableData = []
      }).finally(() => {
        this.loading = false
      })
    },
    changeIsAmin(value){
      if (value === ""){
        delete(this.searchFormLog.is_admin)
      }else{
        this.searchFormLog.is_admin = value
      }
    },
  },
}
</script>

<style scoped>
/* 隐藏分页组件 */
::v-deep .hide-pagination .pagebox {
  display: none !important;
}
</style>
