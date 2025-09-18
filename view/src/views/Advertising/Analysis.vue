<template>
  <div>
    <!--form-->
    <el-form :inline="true" :model="searchForm">
      <el-form-item>
        <date-picker v-model="searchForm.date_range" type="daterange" :clearable="false"></date-picker>
      </el-form-item>
      <el-form-item>
        <el-select v-model="searchForm.dim" :placeholder="$t('common.dim')">
          <el-option v-for="(item, index) in dim" :key="index" :label="item.label" :value="item.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="fetch({page:1})">{{ $t('common.Search') }}</el-button>
        <el-button type="success" :loading="loading" v-perm="'advertising_analysis_create'"
                   @click="showDialog('create', {date: searchForm.date_range[1]})">
          {{ $t('common.Create') }}
        </el-button>
      </el-form-item>
    </el-form>

    <!--chart-->
    <chart ref="chart" :options="chartOption" :seriesData="chartOption.series" width="100%" height="350px"></chart>

    <!--table-->
    <el-table
      ref="tableData"
      :data="tableData"
      height="480"
      border
      v-loading="loading"
      style="width: 100%;margin-bottom: 20px"
      :show-summary="true" :summary-method="summary"
    >
      <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
      <el-table-column prop="cost" :label="$t('advertising.cost')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.cost | localeNum}}</template>
      </el-table-column>
      <el-table-column prop="exposure_cnt" :label="$t('advertising.exposure_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.exposure_cnt | localeNum}}</template>
      </el-table-column>
      <el-table-column prop="exposure_user" :label="$t('advertising.exposure_user')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.exposure_user | localeNum}}</template>
      </el-table-column>
      <el-table-column prop="exposure_per_time" :label="$t('advertising.exposure_per_time')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.exposure_per_time | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="click_cnt" :label="$t('advertising.click_cnt')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.click_cnt | localeNum}}</template>
      </el-table-column>
      <el-table-column prop="click_user" :label="$t('advertising.click_user')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.click_user | localeNum}}</template>
      </el-table-column>
      <el-table-column prop="click_per_cost" :label="$t('advertising.click_per_cost')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.click_per_cost | localeNum2f}}</template>
      </el-table-column>
      <el-table-column prop="click_rate" :label="$t('advertising.click_rate')" min-width="90px" align="center">
        <template v-slot="data">{{data.row.click_rate | localeNum2f}}%</template>
      </el-table-column>
      <el-table-column v-if="$checkPermission('advertising_analysis_update')" fixed="right" >
        <template v-slot="data">
          <el-button type="warning" :loading="loading" :disabled="data.row.date.length<=7"
                     @click="showDialog('edit', data.row)">
            {{ $t('common.Edit') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!--dialog-->
    <el-dialog
      :visible.sync="dialog.flag!==''"
      width="60%"
      :before-close="closeDialog"
    >
      <el-form ref="form" :model="dialog.form" :rules="dialog.rule" label-width="200px">
        <el-form-item prop="date" :label="$t('common.Date')">
          <date-picker v-if="dialog.flag==='create'" v-model="dialog.form.date" type="date" :clearable="false"></date-picker>
          <template v-if="dialog.flag==='edit'">{{dialog.form.date}}</template>
        </el-form-item>
        <el-form-item prop="cost" :label="$t('advertising.cost')">
          <el-input-number v-model="dialog.form.cost"></el-input-number>
        </el-form-item>
        <el-form-item prop="exposure_cnt" :label="$t('advertising.exposure_cnt')">
          <el-input-number v-model="dialog.form.exposure_cnt"></el-input-number>
        </el-form-item>
        <el-form-item prop="exposure_user" :label="$t('advertising.exposure_user')">
          <el-input-number v-model="dialog.form.exposure_user"></el-input-number>
        </el-form-item>
        <el-form-item :label="$t('advertising.exposure_per_time')">
          {{(dialog.form.exposure_cnt / dialog.form.exposure_user).toFixed(2) | localeNum}}
        </el-form-item>
        <el-form-item prop="click_cnt" :label="$t('advertising.click_cnt')">
          <el-input-number v-model="dialog.form.click_cnt"></el-input-number>
        </el-form-item>
        <el-form-item prop="click_user" :label="$t('advertising.click_user')">
          <el-input-number v-model="dialog.form.click_user"></el-input-number>
        </el-form-item>
        <el-form-item :label="$t('advertising.click_per_cost')">
          {{(dialog.form.cost / dialog.form.click_cnt).toFixed(2) | localeNum}}
        </el-form-item>
        <el-form-item :label="$t('advertising.click_rate')">
          {{(100 * dialog.form.click_cnt / dialog.form.exposure_cnt).toFixed(2) | rate}}
        </el-form-item>
        <el-form-item>
          <el-button type="warning" @click="formConfirm()">{{$t('common.Confirm')}}</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>

  </div>
</template>

<script>
import api from '@/api'
import moment from 'moment'
import DatePicker from '@/components/DatePicker'
import Chart from '@/components/Chart/Chart'
import {mapGetters} from 'vuex'
import {confirmBox} from '@/utils/box'
import i18n from "@/utils/i18n";

export default {
  name: "Analysis",
  components: {
    DatePicker, Chart
  },
  data(){
    return {
      confirmBox: confirmBox,
      loading: false,
      searchForm: {
        date_range: [],
        dim: "daily",
      },
      chartOption: {},
      tableData: [],
      dialog: {
        flag: "",
        form: {
          date: "",
          cost: 0,
          exposure_cnt: 0,
          exposure_user: 0,
          click_cnt: 0,
          click_user: 0,
        },
        rule: {
          cost: [
            {required: true, message: this.$t('advertising.validMsg.cost'), trigger: 'change'},
            {validator: this.validCost, trigger: 'change'}
          ],
          exposure_cnt: [
            {required: true, message: this.$t('advertising.validMsg.exposure_cnt'), trigger: 'change'},
            {validator: this.validExposureCnt, trigger: 'change'}
          ],
          exposure_user: [
            {required: true, message: this.$t('advertising.validMsg.exposure_user'), trigger: 'change'},
            {validator: this.validExposureUser, trigger: 'change'}
          ],
          click_cnt: [
            {required: true, message: this.$t('advertising.validMsg.click_cnt'), trigger: 'change'},
            {validator: this.validClickCnt, trigger: 'change'}
          ],
          click_user: [
            {required: true, message: this.$t('advertising.validMsg.click_user'), trigger: 'change'},
            {validator: this.validClickUser, trigger: 'change'}
          ],
        },
      },
    }
  },
  computed: {
    ...mapGetters({
      dim: "option/adDim",
    }),
    timeRangeArr() {
      let timeArr = []
      if (this.searchForm.dim === "daily"){
        let start = moment(this.searchForm.date_range[0], 'YYYY-MM-DD')
        let end = moment(this.searchForm.date_range[1], 'YYYY-MM-DD')
        while (start <= end) {
          timeArr.push(start.format('YYYY-MM-DD'))
          start = start.add(1, 'day')
        }
      }else{
        let start = moment(this.searchForm.date_range[0], 'YYYY-MM')
        let end = moment(this.searchForm.date_range[1], 'YYYY-MM')
        while (start <= end) {
          timeArr.push(start.format('YYYY-MM'))
          start = start.add(1, 'month')
        }
      }
      return timeArr
    },
  },
  created() {
    this.searchForm.date_range = [
      moment().subtract(1, 'month').format('YYYY-MM-DD'),
      moment().subtract(1, 'day').format('YYYY-MM-DD'),
    ]
  },
  mounted() {
    this.fetch({})
  },
  methods:{
    fetch(params={}){
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.loading = true
      api.getAdvertisingList(this.searchForm).then(res =>{
        this.tableData = res.data || []
      }).catch(()=>{
        this.tableData = []
      }).finally(()=>{
        this.parseData(this.tableData)
        this.rebuildCharts()
        this.loading = false
      })
    },
    async createItem(form){
      this.loading = true
      await api.createAdvertisingList(form).then().catch().finally(()=>{this.loading = false})
    },
    async updateItem(form){
      this.loading = true
      await api.updateAdvertisingList(form.date, form).then().catch().finally(()=>{this.loading = false})
    },

    showDialog(flag, form){
      this.dialog.form.date = form.date || ""
      this.dialog.form.cost = form.cost || 0
      this.dialog.form.exposure_cnt = form.exposure_cnt || 0
      this.dialog.form.exposure_user = form.exposure_user || 0
      this.dialog.form.click_cnt = form.click_cnt || 0
      this.dialog.form.click_user = form.click_user || 0

      this.dialog.flag = flag
      this.$nextTick(() => {  this.$refs['form'].clearValidate();});
    },
    closeDialog(){
      this.showDialog("", {})
    },
    formConfirm(){
      this.$refs["form"].validate((valid) => {
        if (valid){
          this.confirmBox().then(async ()=>{
            if (this.dialog.flag === 'create'){
              await this.createItem(this.dialog.form)
            }else{
              await this.updateItem(this.dialog.form)
            }

            this.closeDialog()
            this.fetch()
          }).catch(()=>{
          }).finally(()=>{
          })
        }
      });
    },
    validCost(rule, value, callback){
      if (value < 0) {
        callback(new Error(i18n.t("advertising.validMsg.cost")))
        return
      }
      callback()
    },
    validExposureCnt(rule, value, callback){
      if (!Number.isInteger(Number(value)) || value < 0) {
        callback(new Error(i18n.t("advertising.validMsg.exposure_cnt")))
        return
      }
      if (this.dialog.form.exposure_cnt < this.dialog.form.exposure_user){
        callback(new Error(i18n.t("advertising.validMsg.exposure_cnt_ge_user")))
        return
      }
      if (this.dialog.form.exposure_cnt < this.dialog.form.click_cnt){
        callback(new Error(i18n.t("advertising.validMsg.exposure_cnt_ge_click_cnt")))
        return
      }
      callback()
    },
    validExposureUser(rule, value, callback){
      if (!Number.isInteger(Number(value)) || value < 0) {
        callback(new Error(i18n.t("advertising.validMsg.exposure_user")))
        return
      }
      if (this.dialog.form.exposure_cnt < this.dialog.form.exposure_user){
        callback(new Error(i18n.t("advertising.validMsg.exposure_cnt_ge_user")))
        return
      }
      if (this.dialog.form.exposure_user < this.dialog.form.click_user){
        callback(new Error(i18n.t("advertising.validMsg.exposure_user_ge_click_user")))
        return
      }
      callback()
    },
    validClickCnt(rule, value, callback){
      if (!Number.isInteger(Number(value)) || value < 0) {
        callback(new Error(i18n.t("advertising.validMsg.click_cnt")))
        return
      }
      if (this.dialog.form.exposure_cnt < this.dialog.form.click_cnt){
        callback(new Error(i18n.t("advertising.validMsg.exposure_cnt_ge_click_cnt")))
        return
      }
      if (this.dialog.form.click_cnt < this.dialog.form.click_user){
        callback(new Error(i18n.t("advertising.validMsg.click_cnt_ge_user")))
        return
      }
      callback()
    },
    validClickUser(rule, value, callback){
      if (!Number.isInteger(Number(value)) || value < 0) {
        callback(new Error(i18n.t("advertising.validMsg.click_user")))
        return
      }
      if (this.dialog.form.click_cnt < this.dialog.form.click_user){
        callback(new Error(i18n.t("advertising.validMsg.click_cnt_ge_user")))
        return
      }
      if (this.dialog.form.exposure_user < this.dialog.form.click_user){
        callback(new Error(i18n.t("advertising.validMsg.exposure_user_ge_click_user")))
        return
      }
      callback()
    },

    rebuildCharts(){
      this.$nextTick(()=>{
        this.chartOption = {}

        this.$refs["chart"].rebuildCharts()
      })
    },
    parseData(dataList) {
      let seriesData = []
      let _seriesData = {
        cost: [],
        exposure_cnt: [],
        exposure_user: [],
        exposure_per_time: [],
        click_cnt: [],
        click_user: [],
        click_per_cost: [],
        click_rate: [],
      }
      Object.values(dataList).forEach(dataItem => { // 每条数据
        Object.keys(_seriesData).forEach(key =>{
            _seriesData[key].push(dataItem[key]||0)
        })
      })

      Object.keys(_seriesData).forEach(seriesName => {
        let _series = {
          name: this.$t('advertising.'+seriesName),
          type: 'line',
          smooth: true,
          connectNulls: true,
          emphasis: {
            focus: 'self',
          },
          showSymbol: false,
          data: _seriesData[seriesName].reverse(),
        }
        if (seriesName === 'click_rate'){
          _series['yAxisIndex'] = 1
        }
        seriesData.push(_series)
      })
      this.setChartOptions(seriesData)
    },
    setChartOptions(seriesData) {
      let title = this.$t('advertising.analysis')
      this.chartOption = {
        title: {
          text: title,
          left: 100,
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross'
          }
        },
        legend: {},
        xAxis: {
          type: 'category',
          data: this.timeRangeArr,
        },
        yAxis: [
          {
            type: 'value',
            position: 'left',
            minInterval: 1,
          },
          {
            type: 'value',
            name: this.$t('advertising.click_rate'),
            position: 'right',
            axisLabel: {
              formatter: '{value}%'
            }
          },
        ],
        series: seriesData,
      }
      let selected = {
        [this.$t('advertising.exposure_cnt')]: true,
        [this.$t('advertising.click_cnt')]: true,
      }
      for (let i=0; i<seriesData.length; i++){
        if (seriesData[i].name !== this.$t('advertising.exposure_cnt') &&
          seriesData[i].name !== this.$t('advertising.click_cnt')){
          selected[seriesData[i].name] = false
        }
      }
      this.chartOption["legend"]["selected"] = selected
    },
    summary({ columns, data }){
      let row = Array(columns.length)
      let sumValue = {}
      for (let i = 0; i < data.length; i++) {
        for (let key in data[i]){
          if (key !== "date"){
            sumValue[key] = data[i][key] + (sumValue[key] || 0)
          }
        }
      }

      sumValue["exposure_per_time"] = (sumValue["exposure_cnt"] / sumValue["exposure_user"]).toFixed(2)
      sumValue["click_per_cost"] = (sumValue["cost"] / sumValue["click_cnt"]).toFixed(2)
      sumValue["click_rate"] = (100 * sumValue["click_cnt"] / sumValue["exposure_cnt"]).toFixed(2)
      columns.forEach((el, index) => {
        switch (el.property){
          case "date":
            row[index] = "TOTAL"
            break
          case "click_rate":
            if (isNaN(sumValue[el.property]) || !isFinite(sumValue[el.property])){
              row[index] = '0%'
            }else{
              row[index] = sumValue[el.property]+'%'
            }
            break
          case undefined:
            break
          default:
            if (isNaN(sumValue[el.property]) || !isFinite(sumValue[el.property])){
              row[index] = '0'
            }else{
              row[index] = sumValue[el.property]
            }
            break
        }
      })

      return row
    },
  },
}
</script>

<style scoped>
.el-table{
  overflow:visible !important;
}
</style>
