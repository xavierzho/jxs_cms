<template>
    <div>
      <search-table
        :loading="loading"
        :table-data="tableData"
        :page="searchForm.page" :page-size="searchForm.pageSize" :total="total"
        :height="'75vh'"
        @fetch="fetch"
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

            <el-form-item>
              <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
              <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
            </el-form-item>
          </el-form>
        </template>
        <template v-slot:table-column>
          <el-table-column prop="created_at" :label="$t('common.DateTime')" min-width="90px" align="center"></el-table-column>
          <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.user_id }}</template>
          </el-table-column>
          <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.user_name}}</template>
          </el-table-column>
          <el-table-column prop="sign_in_type" :label="$t('activity.signIn.sign_in_type')" min-width="90px" align="center">
            <template v-slot="data">
                {{ data.row.sign_in_type}}
            </template>
          </el-table-column>
          <el-table-column prop="day_no" :label="$t('activity.signIn.day_no')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.day_no }}</template>
          </el-table-column>
          <el-table-column prop="type" :label="$t('activity.signIn.type')" min-width="90px" align="center">
            <template v-slot="data">
                {{ data.row.type }}
            </template>
          </el-table-column>
          <el-table-column prop="value" :label="$t('activity.signIn.value')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.value }}</template>
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
    name: "SignIn",
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
        },
        total: 0,
        tableData: [],
      }
    },
    computed: {
      ...mapGetters({
        userType: 'option/userType',
        userChannel: 'option/userChannel',
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
        api.getActivitySignIn(this.searchForm).then(res => {
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
        this.searchForm = Object.assign({}, this.searchForm, params)
        api.exportActivitySignIn(this.searchForm).
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
    },
  }
  </script>

  <style lang="scss" scoped>
  </style>
