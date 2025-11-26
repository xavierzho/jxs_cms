<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData"
      :page="searchFormLog.page" :page-size="searchFormLog.pageSize" :total="total"
      :height="'64vh'"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <div v-perm="'inquire_bet_item_filter'">
          <el-form inline v-model="searchFormLog">
            <el-form-item :label="$t('common.DateTime')">
              <date-picker v-model="searchFormLog.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss" :clearable="false"></date-picker>
            </el-form-item>
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

          <el-form inline v-model="searchFormLog">
            <el-form-item :label="$t('inquire.item.log_type')">
              <el-select v-model="searchFormLog.log_type_list" multiple clearable>
                <el-option v-for="item in inquireItemLogType" :label="item.label" :value="item.value" :key="item.value" v-if="item.value >= 100 && item.value <= 199" ></el-option>
              </el-select>
            </el-form-item>
            <el-form-item>
            <span slot="label">{{$t('inquire.item.gacha_name')}}
              <el-tooltip effect="dark" :content="$t('inquire.item.filterOptionTip_gachaName')" placement="top"><i class='el-icon-info' /></el-tooltip>
            </span>
              <el-input v-model="searchFormLog.gacha_name" clearable></el-input>
            </el-form-item>
          </el-form>

          <el-form inline v-model="searchFormLog">
            <el-form-item :label="$t('inquire.item.update_amount')">
              <input-number-range v-model="searchFormLog.update_amount_range" clearable :precision="2"></input-number-range>
            </el-form-item>
            <el-form-item :label="$t('inquire.item.show_price')">
              <input-number-range v-model="searchFormLog.show_price_range" clearable :range="[0, 2**31-1]" :precision="2"></input-number-range>
            </el-form-item>
            <el-form-item :label="$t('inquire.item.inner_price')">
              <input-number-range v-model="searchFormLog.inner_price_range" clearable :range="[0, 2**31-1]" :precision="2"></input-number-range>
            </el-form-item>
          </el-form>
        </div>
        <el-form inline v-model="searchFormLog">
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="date_time" :label="$t('common.DateTime')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"></el-table-column>
<!--        <el-table-column prop="log_type_str" :label="$t('inquire.item.log_type')" min-width="90px" align="center"></el-table-column>-->
<!--        <el-table-column prop="log_type_name" :label="$t('inquire.item.gacha_name')" min-width="90px" align="center"></el-table-column>-->
        <el-table-column prop="no" :label="$t('inquire.item.no')" min-width="40px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.no | localeNum}}</template>
          </template>
        </el-table-column>
<!--        <el-table-column prop="box_out_no" :label="$t('inquire.item.box_out_no')" min-width="90px" align="center"></el-table-column>-->
        <el-table-column prop="item_id" :label="$t('inquire.item.item_id')" min-width="90px" align="center"/>
        <el-table-column prop="item_name" :label="$t('inquire.item.item_name')" min-width="90px" align="center"/>
        <el-table-column prop="level_name" :label="$t('inquire.item.level_name')" min-width="90px" align="center"/>
        <el-table-column :label="$t('inquire.item.cover_thumb')" min-width="80px" align="center">
          <template v-slot="detail">
            <el-image :src="detail.row.cover_thumb" style="width: 75px; height: 75px"></el-image>
          </template>
        </el-table-column>
        <el-table-column prop="show_price" :label="$t('inquire.item.show_price')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.show_price | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="inner_price" :label="$t('inquire.item.inner_price')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.inner_price | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="recycling_price" :label="$t('inquire.item.recycling_price')" min-width="90px" align="center">
          <template v-slot="detail">
            <template>{{detail.row.recycling_price | localeNum2f}}</template>
          </template>
        </el-table-column>
        <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
          <template v-slot="detail">{{detail.row.nums | localeNum}}</template>
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
import {isInvalidNumber} from "@/utils/validate";
import {checkPermission} from "@/utils/auth";

export default {
  name: "BetItem",
  components: {SearchTable, ItemDetail, InputNumber, InputNumberRange, DatePicker},
  data(){
    return {
      loading: false,
      searchFormLog: {
        page: 1,
        page_size: 50,
        date_time_range: [],
        user_ids: '',
        user_name: '',
        tel: '',
        // is_admin: false,
        channel: null,
        log_type_list: ['101','102','103','104','105','106'],
        gacha_name: '',
        update_amount_range: [null, -0.01], // 没填则不传, 用边界补齐 两个元素
        show_price_range: [],
        inner_price_range: [1000, 999999],
      },
      tableData: [],
      total: 0,
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
      userChannel: 'option/userChannel',
      inquireItemLogType: 'option/inquireItemLogType',
    }),
  },
  created() {
    if (this.inquireItemLogType.length === 0 && checkPermission('inquire_bet_item_filter')) {
      this.$store.dispatch('option/update', 'inquireItemLogType')
    }
    let yDay = moment(moment().subtract(1, 'day').format('YYYY-MM-DD'))
    this.searchFormLog.date_time_range = [
      yDay.format('YYYY-MM-DD HH:mm:ss'),
      yDay.add(1, "day").add(-1, "second").format('YYYY-MM-DD HH:mm:ss'),
    ]
    this.fetch()
  },
  methods:{
    fetch(params = {}){
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      searchFormLog.show_price_range = this.setRangeBorder(searchFormLog.show_price_range)
      searchFormLog.inner_price_range = this.setRangeBorder(searchFormLog.inner_price_range)
      api.getInquireBetItemDetailList(searchFormLog).
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
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      searchFormLog.show_price_range = this.setRangeBorder(searchFormLog.show_price_range)
      searchFormLog.inner_price_range = this.setRangeBorder(searchFormLog.inner_price_range)
      api.exportInquireBetItemDetailList(searchFormLog).
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

</style>
