<template>
    <div>
      <search-table
        :loading="loading"
        :table-data="tableData" show-summary :summary-method="summaryMethod "
        :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
        :height="'75vh'"
        @fetch="fetch"
        >
    <template v-slot:table-search>
          <el-form inline v-model="searchForm">
            <el-form-item :label="$t('common.DateTime')">
              <date-picker v-model="searchForm.date_range" type="daterange" value-format="yyyy-MM-dd" :clearable="false"></date-picker>
            </el-form-item>
                <!-- <el-form-item>
                  <span slot="label">{{$t('user.UserType')}}
                    <el-tooltip effect="dark" :content="$t('user.UserType')" placement="top"></el-tooltip>
                  </span>
                  <el-select v-model="searchForm.user_type"  clearable>
                    <el-option v-for="item in inviteUserType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
                  </el-select>
                </el-form-item> -->
                <el-form-item :label="$t('user.UserID')">
                  <el-input v-model="searchForm.user_id" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('user.Name')">
                  <el-input v-model="searchForm.user_name" clearable></el-input>
                </el-form-item>
                <el-form-item :label="$t('user.Tel')">
                  <input-number v-model="searchForm.tel" :range="[0,99999999999]" clearable></input-number>
                </el-form-item>
    
          <!-- <el-form-item :label="$t('user.IsRewards')" @click="() => {fetchDaily({page: 1});}" ><el-switch v-model="searchForm.is_rewards" /></el-form-item> -->
          </el-form>
         
          <el-form inline v-model="searchForm">
            <el-form-item :label="$t('activity.turntable.type')">
            <el-select v-model="searchForm.type"  clearable>       <!-- multiple  多选-->
              <el-option v-for="item in prizeType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
            </el-form-item>

            <el-form-item>
            <span slot="label">{{$t('activity.turntable.name')}}
            </span>
            <el-input v-model="searchForm.name" clearable></el-input>
            </el-form-item>
   

          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
          </el-form-item>
        </el-form>

    </template>
            <template v-slot:table-column>
                <el-table-column prop="date" :label="$t('common.Date')" min-width="70px" align="center"></el-table-column>
                <el-table-column prop="user_id" :label="$t('activity.turntable.user_id')" min-width="50px" align="center">
                  <template v-slot="data">{{ data.row.user_id }}</template>
                </el-table-column>
                <el-table-column prop="user_name" :label="$t('activity.turntable.user_name')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.user_name}}</template>
                </el-table-column>
                <el-table-column prop="name" :label="$t('activity.turntable.name')" min-width="70px" align="center">
                  <template v-slot="data">{{ data.row.name }}</template>
                </el-table-column>
                <el-table-column prop="period" :label="$t('activity.turntable.period')" min-width="50px" align="center">
                  <template v-slot="data">{{ data.row.period }}</template>
                </el-table-column>

                <el-table-column prop="point_type_name" :label="$t('activity.turntable.point_type_name')" min-width="50px" align="center">
                  <template v-slot="data">{{ data.row.point_type_name }}</template>
                </el-table-column>

                <el-table-column prop="point" :label="$t('activity.turntable.point')" min-width="50px" align="center">
                  <template v-slot="data">{{ data.row.point }}</template>
                </el-table-column>

                <el-table-column prop="type" :label="$t('activity.turntable.type')" min-width="50px" align="center">
                  <template v-slot="data">{{ data.row.type_name }}</template>
                </el-table-column>

                <el-table-column prop="item_id" :label="$t('activity.turntable.item_id')" min-width="60px" align="center">
                  <template v-slot="data">{{ data.row.item_id }}</template>
                </el-table-column>

                <el-table-column prop="item_name" :label="$t('activity.turntable.item_name')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.item_name }}</template>
                </el-table-column>

                <el-table-column prop="prize_value" :label="$t('activity.turntable.prize_value')" min-width="90px" align="center">
                  <template v-slot="data">{{ data.row.prize_value }}</template>
                </el-table-column>
            </template>
    </search-table>
    </div>
  </template>
  
  <script>
  import DatePicker from "@/components/DatePicker.vue";
  import moment from "moment/moment";
  import api from "@/api";
  import SearchTable from "@/components/Table/SearchTable.vue";
  import InputNumber from "@/components/Input/InputNumber.vue";
  import {mapGetters} from "vuex";
  
  export default {
    name: "Invite",
    components: {InputNumber,SearchTable, DatePicker},
    data(){
      return {
        loading: false,
        searchForm: {
          page: 1,
          page_size: 50,
          date_range: [],
          user_id: '',
          user_name: '',
          tel: '',
          name:'',
          type:'',
        },
        total: 0,
        amount: 0,
        totalamount: 0,
        difference: 0,
        tableData: [],
        summary: {},
      }
    },
  
    computed: {
      ...mapGetters({
        prizeType: 'option/prizeType',
      }),
    },
    created() {
      this.searchForm.date_range = [
        moment().subtract(1, 'month').format('YYYY-MM-DD'),
        moment().format('YYYY-MM-DD'),
      ]
      this.fetch()
    },
    methods:{
  
      fetch(params = {}) {
        this.loading = true
        this.searchForm = Object.assign({}, this.searchForm, params)
        api.getActivityTurntable(this.searchForm).then(res => {
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
        api.exportActivityTurntable(this.searchForm).
        finally(() => {
          this.loading = false
        })
      },
  
  
      summaryMethod({ columns, data }){
        return [
          'Total',
          undefined, undefined, undefined, undefined,undefined, undefined, undefined, undefined,undefined,
          Number(this.summary['total_amount'] || 0).toLocaleString(),
        ]
      },
    },
  }
  </script>
  
  <style scoped>
  
  </style>
  