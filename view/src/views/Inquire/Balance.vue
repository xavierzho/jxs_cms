<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod" :cellClassName="cellClassName"
      :page="searchFormLog.page" :page-size="searchFormLog.pageSize" :total="total"
      :height="'70vh'" :row-style="{height: '75px'}"
      @fetch="fetch"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('common.dateTimeType')">
            <el-select v-model="searchFormLog.date_time_type">
              <el-option v-for="(item, index) in dateTimeType" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchFormLog.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss" :clearable="false"></date-picker>
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
          <el-form-item :label="$t('user.UserType')">
            <el-select v-model="searchFormLog.role" clearable multiple>
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
          <el-form-item :label="$t('inquire.balance.source_type')">
            <el-select v-model="searchFormLog.source_type" multiple clearable>
              <el-option v-for="item in inquireBalanceSourceType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <span slot="label">{{$t('inquire.balance.item_name')}}
              <el-tooltip effect="dark" :content="$t('inquire.balance.filterOptionTipBet')" placement="top"><i class='el-icon-info' /></el-tooltip>
            </span>
            <el-input v-model="searchFormLog.item_name" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <span slot="label">{{$t('inquire.balance.channel_type')}}
              <el-tooltip effect="dark" :content="$t('inquire.balance.filterOptionTipRecharge')" placement="top"><i class='el-icon-info' /></el-tooltip>
            </span>
            <el-select v-model="searchFormLog.channel_type" multiple clearable>
              <el-option v-for="item in inquireBalanceChannelType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <span slot="label">{{$t('inquire.balance.pay_source_type')}}
               <el-tooltip effect="dark" :content="$t('inquire.balance.filterOptionTipRecharge')" placement="top"><i class='el-icon-info' /></el-tooltip>
            </span>
            <el-select v-model="searchFormLog.pay_source_type" multiple clearable>
              <el-option v-for="item in inquireBalancePaySourceType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('inquire.balance.balance_type')">
            <el-select v-model="searchFormLog.balance_type" multiple clearable>
              <el-option v-for="item in inquireBalanceType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
          </el-form-item>

          <el-form-item :label="$t('inquire.balance.update_amount')">
            <input-number-range v-model="searchFormLog.update_amount_range" clearable></input-number-range>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="id" :label="$t('inquire.balance.id')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdDateTime')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="finish_at" :label="$t('common.finishDateTime')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="source_type_str" :label="$t('inquire.balance.source_type')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="item_name" :label="$t('inquire.balance.item_name')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.item_name}}</template>
        </el-table-column>
        <el-table-column prop="platform_order_id" :label="$t('inquire.balance.platform_order_id')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.platform_order_id}}</template>
        </el-table-column>
        <el-table-column prop="pay_source_type_str" :label="$t('inquire.balance.pay_source_type')" min-width="90px" align="center"></el-table-column>
        <el-table-column :label="$t('inquire.balance.amount')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.before_balance | localeNum2f}} -> {{data.row.after_balance | localeNum2f}}</template>
        </el-table-column>
        <el-table-column prop="update_amount" :label="$t('inquire.balance.update_amount')" min-width="90px" align="center">
          <template v-slot="data">{{formatUpdateAmount(data.row.update_amount)}}</template>
        </el-table-column>
        <el-table-column prop="balance_type_name" :label="$t('inquire.balance.balance_type')" min-width="90px" align="center">
          <template v-slot="data">{{data.row.balance_type_name}} </template>
        </el-table-column>
        <el-table-column prop="comment" :label="$t('inquire.balance.comment')" min-width="90px" align="center">
          <template v-slot="data">
            <a @click="editComment(data.row.id, data.row.comment)">
              <el-row v-if="data.row.comment.length>0">
                <el-col style="border-radius: 5px; box-shadow: 0 2px 12px 0 rgba(0,0,0,.1)">
                  <el-tooltip effect="light" :content="data.row.comment[0].Comment" placement="top-start">
                    <div style="padding: 5px; text-align: left">{{formatComment(data.row.comment[0].Comment)}}</div>
                  </el-tooltip>
                </el-col>
                <el-col v-if="data.row.comment.length>1" style="font-size: 2rem">...</el-col>
              </el-row>
              <el-row v-else><el-col style="font-size: 2rem">+</el-col></el-row>
            </a>
          </template>
        </el-table-column>
      </template>
    </search-table>

    <el-dialog :visible.sync="commentDialog.visible" :title="$t('inquire.balance.edit_comment')" @close="editCommentClose">
      <el-row>
        <el-col>
          <el-timeline>
            <el-timeline-item
              v-for="comment in commentDialog.list" :key="comment.ID"
              :timestamp="comment.CreatedAt" placement="top">
              <el-card>
                <p style="font-size: 0.75rem; color: gray">{{$t('user.UserID')}}: {{comment.UserID}}</p>
                <p style="margin:10px 0; font-size: 1.5rem">{{comment.Comment}}</p>
                <p v-if="userInfo.isAdmin || comment.UserID === userInfo.id"
                   style="display: flex; justify-content: right">
                  <el-button type="warning" :loading="loading" @click="deleteComment({id: commentDialog.id, comment_id: comment.ID})">{{$t('common.Delete')}}</el-button>
                </p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-col>
      </el-row>
      <el-row style="display: flex; align-items: center">
        <el-col :span="2">{{$t('inquire.balance.add_comment')}}:</el-col>
        <el-col :span="10"><el-input v-model="commentDialog.newComment"></el-input></el-col>
        <el-col :span="6" :offset="6" style="display: flex; justify-content: right">
          <el-button type="primary" :loading="loading" @click="addComment({id: commentDialog.id, comment: commentDialog.newComment})">{{$t('common.Create')}}</el-button>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import {mapGetters,} from "vuex";
import InputNumber from "@/components/Input/InputNumber.vue"
import InputNumberRange from "@/components/Input/InputNumberRange.vue"
import SearchTable from "@/components/Table/SearchTable.vue";
import { isInvalidNumber } from '@/utils/validate'
import Decimal from "decimal.js";

export default {
  name: "Balance",
  components: {SearchTable, InputNumber, InputNumberRange, DatePicker},
  data(){
    return {
      loading: false,
      searchFormLog: {
        page: 1,
        page_size: 50,
        date_time_type: 'created',
        date_time_range: [],
        user_id: '',
        user_name: '',
        tel: '',
        role: [],
        channel: null,
        source_type: [],
        channel_type: [],
        pay_source_type: [],
        item_name: '',
        update_amount_range: [], // 没填则不传, 用边界补齐 两个元素
        balance_type: [],
      },
      total: 0,
      tableData: [],
      summary: {},
      commentDialog: {
        visible: false,
        id: 0,
        list: [],
        newComment: "",
      },
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
      userChannel: 'option/userChannel',
      inquireBalanceSourceType: 'option/inquireBalanceSourceType',
      inquireBalanceChannelType: 'option/inquireBalanceChannelType',
      inquireBalancePaySourceType:'option/inquireBalancePaySourceType',
      userInfo: 'user/userInfo',
      dateTimeType: "option/dateTimeType",
      inquireBalanceType: 'option/inquireBalanceType',
    }),
  },
  created() {
    if (this.inquireBalanceSourceType.length === 0) {
      this.$store.dispatch('option/update', 'inquireBalanceSourceType')
    }
    if (this.inquireBalanceChannelType.length === 0) {
      this.$store.dispatch('option/update', 'inquireBalanceChannelType')
    }
    if (this.inquireBalancePaySourceType.length === 0) {
      this.$store.dispatch('option/update', 'inquireBalancePaySourceType')
    }
    if (this.inquireBalanceType.length === 0) {
      this.$store.dispatch('option/update', 'inquireBalanceType')
    }
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
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      api.getInquireBalance(searchFormLog).
      then(res => {
        this.tableData = res.data || []
        this.summary = res.headers.summary || {}
        this.total = this.summary['cnt'] || 0
        this.formatData()
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
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      api.exportInquireBalance(searchFormLog).
      finally(() => {
        this.loading = false
      })
    },
    addComment(params={}){
      this.loading = true
      api.addComment(params).then(() =>{
        this.editCommentClose()
        this.fetch()
      }).finally(()=>{
        this.loading = false
      })
    },
    deleteComment(params={}){
      this.loading = true
      api.deleteComment(params).then(() =>{
        this.editCommentClose()
        this.fetch()
      }).finally(()=>{
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
    formatUpdateAmount(value){
      let valueStr = Number(new Decimal(value).toFixed(2, Decimal.ROUND_HALF_UP)).toLocaleString()
      if (value > 0){
        return "+" + valueStr
      }
      return valueStr
    },
    formatData(){
      for (let i=0; i<this.tableData.length; i++){
        this.tableData[i]['comment'] =this.tableData[i]['comment'] || []
      }
    },
    formatComment(comment){
      return comment.slice(0, 100)
    },
    cellClassName(cell){
      if (cell.column.property !== 'update_amount') {
        return ''
      }
      if (cell.row.update_amount > 0){
        return 'value-green'
      }else if(cell.row.update_amount < 0){
        return 'value-red'
      }
      return ''
    },
    summaryMethod({ columns, data }){
      return [
        'Total',
        undefined, Number(this.summary['user_cnt']).toLocaleString(),
        undefined, undefined, undefined, undefined, undefined, undefined,
        Number(this.summary['update_amount']).toLocaleString(), undefined,
      ]
    },
    editComment(id, commentList){
      this.commentDialog.id = id
      this.commentDialog.list = commentList
      this.commentDialog.newComment = ""
      this.commentDialog.visible = true
    },
    editCommentClose(){
      this.commentDialog.id = 0
      this.commentDialog.list = []
      this.commentDialog.newComment = ""
      this.commentDialog.visible = false
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

<style lang="scss" scoped>
::v-deep .value-red{
  color: #F56C6C;
}

::v-deep .value-green{
  color: #67C23A;
}

</style>
