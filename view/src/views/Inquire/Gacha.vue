<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod" :cellClassName="revenueClass"
      :page="searchFormRevenue.page" :page-size="searchFormRevenue.pageSize" :total="total"
      :height="'73vh'"
      @fetch="fetch"
      @expandChange="expandChange"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchFormRevenue">
          <el-form-item :label="$t('inquire.gacha.gacha_type')">
            <el-select v-model="searchFormRevenue.type_list" multiple clearable>
              <el-option v-for="item in inquireGachaType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('inquire.gacha.gacha_name')">
            <el-input v-model="searchFormRevenue.gacha_name" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('inquire.gacha.bet_rate')">
            <input-number-range v-model="searchFormRevenue.bet_rate" clearable :range="[0, 100]"></input-number-range>
          </el-form-item>
          <el-form-item :label="$t('inquire.gacha.revenue_range')">
            <input-number-range v-model="searchFormRevenue.revenue_range" clearable></input-number-range>
          </el-form-item>
          <el-form-item :label="$t('inquire.gacha.revenue_rate_range')">
            <input-number-range v-model="searchFormRevenue.revenue_rate_range" clearable></input-number-range>
          </el-form-item>
          <el-form-item :label="$t('inquire.gacha.is_box_dim')">
            <el-switch v-model="searchFormRevenue.is_box_dim"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-slot:table-column>
        <el-table-column prop="gacha_type_str" :label="$t('inquire.gacha.gacha_type')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="gacha_name" :label="$t('inquire.gacha.gacha_name')" min-width="90px" align="center"></el-table-column>
        <el-table-column v-if="searchFormRevenue.is_box_dim" prop="box_out_no" :label="$t('inquire.gacha.box_out_no')" min-width="90px" align="center"></el-table-column>
        <el-table-column :label="$t('inquire.gacha.bet_nums')" min-width="120px" align="center">
          <template v-slot="data">
            <el-progress :percentage="data.row.bet_rate" :format="formatBetRate(data.row.bet_nums, data.row.total_nums)"></el-progress>
          </template>
        </el-table-column>
        <el-table-column prop="price" :label="$t('inquire.gacha.price')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.price | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="discount_price" :label="$t('inquire.gacha.discount_price')" min-width="90px" align="center">
          <template v-slot="data">
            <div v-if="data.row.discount_price!=='0'">{{data.row.discount_price | localeNum2f}}</div>
            <div v-else>-</div>
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('inquire.gacha.amount')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.amount | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="amount_left" :label="$t('inquire.gacha.amount_left')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.amount_left | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="inner_price_bet" :label="$t('inquire.gacha.inner_price_bet')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.inner_price_bet | localeNum2f}}<br/>{{data.row.inner_price_bet_extra | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="inner_price_left" :label="$t('inquire.gacha.inner_price_left')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.inner_price_left | localeNum2f}}<br/>{{data.row.inner_price_left_extra | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="revenue" :label="$t('inquire.gacha.revenue')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.revenue | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="revenue_rate" :label="$t('inquire.gacha.revenue_rate')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.revenue_rate | rate}}</template>
        </el-table-column>
        <el-table-column min-width="50px" align="center">
          <template v-slot="data">
            <el-button type="warning" v-loading="loading" @click="fetchExport({gacha_id: data.row.gacha_id, box_out_no: data.row.box_out_no,})">{{ $t('common.Export') }}</el-button>
          </template>
        </el-table-column>
        <el-table-column type="expand">
          <template v-slot="data">
            <template v-if="!loading">
              <item-detail :loading="loading" :data="data.row.detail" show-level-name>
                <template v-slot:table-tail-column>
                  <el-table-column :label="$t('inquire.gacha.bet_nums')" min-width="120px" align="center">
                    <template v-slot="data">
                      <el-progress :percentage="data.row.bet_rate" :format="formatBetRate(data.row.bet_nums, data.row.total_nums)"></el-progress>
                    </template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
          </template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import SearchTable from "@/components/Table/SearchTable.vue";
import ItemDetail from "@/components/Page/ItemDetail.vue";
import InputNumber from "@/components/Input/InputNumber.vue";
import InputNumberRange from "@/components/Input/InputNumberRange.vue";
import DatePicker from "@/components/DatePicker.vue";
import {mapGetters} from "vuex";
import api from "@/api";
import { isInvalidNumber } from '@/utils/validate';

export default {
  name: "Gacha",
  components: {SearchTable, ItemDetail, InputNumber, InputNumberRange, DatePicker},
  data(){
    return {
      loading: false,
      searchFormRevenue: {
        page: 1,
        page_size: 50,
        is_box_dim: false,
        type_list: [],
        gacha_name: '',
        bet_rate: [], // 没填则不传, 用边界补齐 两个元素
        revenue_range: [],
        revenue_rate_range: [],
      },
      searchFormDetail: {
        gacha_id: 0,
        box_out_no: 0,
      },
      tableData: [],
      total: 0,
      summary: {},
    }
  },
  computed: {
    ...mapGetters({
      inquireGachaType: 'option/inquireGachaType',
    }),
  },
  created() {
    if (this.inquireGachaType.length === 0) {
      this.$store.dispatch('option/update', 'inquireGachaType')
    }
    this.fetch()
  },
  watch: {
    'searchFormRevenue.is_box_dim': {
      handler(){
        this.fetch({page: 1})
      },
    },
  },
  methods:{
    fetch(params = {}) {
      this.loading = true
      this.searchFormRevenue = Object.assign(this.searchFormRevenue, params)
      let searchFormRevenue = Object.assign({}, this.searchFormRevenue)
      searchFormRevenue.bet_rate = this.setRangeBorder(searchFormRevenue.bet_rate)
      searchFormRevenue.revenue_range = this.setRangeBorder(searchFormRevenue.revenue_range)
      searchFormRevenue.revenue_rate_range = this.setRangeBorder(searchFormRevenue.revenue_rate_range)
      api.getInquireGachaRevenue(searchFormRevenue).
      then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
        this.summary = res.headers.summary || {}
      }).catch(()=>{
        this.tableData = []
        this.total = 0
        this.summary = {}
      }).finally(() => {
        this.loading = false
      })
    },
    fetchDetail(params = {}, row){
      this.loading = true
      this.searchFormDetail = Object.assign({}, this.searchFormDetail, params)
      api.getInquireGachaDetail(this.searchFormDetail).
      then(res => {
        row.detail = res.data || []
      }).catch(()=>{
        row.detail = []
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExport(params = {}){
      this.loading = true
      api.exportInquireGachaDetail(params).
      finally(() => {
        this.loading = false
      })
    },
    setRangeBorder(amountRange){
      if (amountRange.length === 2){
        if (isInvalidNumber(amountRange[0]) && isInvalidNumber(amountRange[1])){
          return []
        }else if (isInvalidNumber(amountRange[0]) && !isInvalidNumber(amountRange[1])){
          return [-(2**30), amountRange[1]]
        }else if (!isInvalidNumber(amountRange[0]) && isInvalidNumber(amountRange[1])){
          return [amountRange[0], (2**30)]
        }
      }
      return amountRange
    },
    formatBetRate(bet_nums, total_nums){
      return (bet_rate) => {
        return Number(bet_rate).toLocaleString() + '% | ' + Number(bet_nums).toLocaleString() + ' / ' + Number(total_nums).toLocaleString()
      }
    },
    expandChange(row){
      if (row.detail === undefined){
        this.fetchDetail({gacha_id: row.gacha_id, box_out_no: row.box_out_no,}, row)
      }
    },
    summaryMethod({ columns, data }){
      let row = [
        Number(this.summary.bet_nums).toLocaleString() + ' / ' + Number(this.summary.total_nums).toLocaleString(),
        undefined, undefined,
        Number(this.summary.amount).toLocaleString(), Number(this.summary.amount_left).toLocaleString(),
        <div><p>{Number(this.summary.inner_price_bet).toLocaleString()}</p><p>{Number(this.summary.inner_price_bet_extra).toLocaleString()}</p></div>,
        <div><p>{Number(this.summary.inner_price_left).toLocaleString()}</p><p>{Number(this.summary.inner_price_left_extra).toLocaleString()}</p></div>,
        Number(this.summary.revenue).toLocaleString(), Number(this.summary.revenue_rate).toLocaleString()+'%',
      ]
      if (this.searchFormRevenue.is_box_dim){
        let head = ['Total', undefined, undefined,]
        head.push(...row)
        return head
      }
      let head = ['Total', undefined,]
      head.push(...row)
      return head
    },
    revenueClass(cell){
      if (cell.column.property !== 'revenue' && cell.column.property !== 'revenue_rate') {
        return ''
      }
      if (cell.row.revenue_rate < -100){
        return 'revenue_warning_2'
      }else if (cell.row.revenue_rate < -50){
        return 'revenue_warning_1'
      }else if (cell.row.revenue_rate < -0){
        return 'revenue_warning_0'
      }else if (cell.row.revenue_rate < 50){
        return 'revenue_good_0'
      }else if (cell.row.revenue_rate < 100){
        return 'revenue_good_1'
      }else{
        return 'revenue_good_2'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
::v-deep .revenue_warning_0{
  background-color: rgb(255,220,160) !important;
}
::v-deep .revenue_warning_1{
  background-color: rgb(255,200,200) !important;
}
::v-deep .revenue_warning_2{
  background-color: rgb(255,100,100) !important;
}
::v-deep .revenue_good_0{
  background-color: rgb(200,255,200) !important;
}
::v-deep .revenue_good_1{
  background-color: rgb(100,255,100) !important;
}
::v-deep .revenue_good_2{
  background-color: rgb(50,255,50) !important;
}
</style>
