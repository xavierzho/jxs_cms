<template>
<div>
    <search-table
    :loading="loading"
    :table-data="tableData" show-summary :summary-method="!searchForm.is_daily ? summaryMethod : summaryMethodDaily"
    :page="searchForm.page" :page-size="searchForm.page_size" :total="total"
    :height="'75vh'"
    @fetch="fetch" 
    >

    <template v-slot:table-search>
        <el-form inline v-model="searchForm">
        <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchForm.date_range" type="daterange" value-format="yyyy-MM-dd" :clearable="false"></date-picker>
        </el-form-item>
        <template v-if="!searchForm.is_daily">
            <el-form-item :label="$t('user.UserType')">
            <el-select v-model="searchForm.user_type" clearable>
                <el-option v-for="item in recallUserType" :label="item.label" :value="item.value" :key="item.value" ></el-option>
            </el-select>
            </el-form-item>
            <el-form-item :label="$t('user.UserID')">
            <el-input v-model="searchForm.user_id" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('user.Name')">
            <el-input v-model="searchForm.user_name" clearable></el-input>
            </el-form-item>
            <el-form-item :label="$t('user.Tel')">
            <input-number v-model="searchForm.tel" :range="[0,99999999999]" clearable></input-number>
            </el-form-item>
        </template>

        <el-form-item :label="$t('report.recall.is_daily')">
            <el-switch v-model="searchForm.is_daily" @change="() => handleSwitchChange(searchForm.is_daily)"/>
        </el-form-item>

        <el-form-item v-if="!searchForm.is_daily">
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
        </el-form-item>
        <el-form-item v-else>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{ $t('common.Search') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchDailyExport({})">{{ $t('common.Export') }}</el-button>
        </el-form-item>
        </el-form>
    </template>

    <template v-slot:table-column>
        <template v-if="!searchForm.is_daily">
        <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('report.recall.user_id')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.user_id }}</template>
        </el-table-column>
        <el-table-column prop="user_name" :label="$t('report.recall.user_name')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.user_name }}</template>
        </el-table-column>
        <el-table-column prop="parent_user_id" :label="$t('report.recall.parent_user_id')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.parent_user_id }}</template>
        </el-table-column>
        <el-table-column prop="parent_user_name" :label="$t('report.recall.parent_user_name')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.parent_user_name }}</template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('report.recall.amount')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.amount | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="point" :label="$t('report.recall.point')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.point | localeNum2f }}</template>
        </el-table-column>
        </template>
        <template v-else>
        <el-table-column prop="date" :label="$t('common.Date')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="total_amount" :label="$t('report.recall.daily.total_amount')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.total_amount | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('report.recall.amount')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.amount | localeNum2f }}</template>
        </el-table-column>
        <el-table-column prop="difference" :label="$t('report.recall.daily.difference')" min-width="90px" align="center">
            <template v-slot="data">{{ data.row.difference | localeNum2f }}</template>
        </el-table-column>
        </template>
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
components: {InputNumber, SearchTable, DatePicker},
data(){
    return {
    loading: false,
    searchForm: {
        page: 1,
        page_size: 50,
        date_range: [],
        user_type: '',
        user_id: '',
        user_name: '',
        tel: '',
        is_daily: false,
    },
    total: 0,
    tableData: [],
    summary: {},
    }
},

computed: {
    ...mapGetters({
    recallUserType: 'option/recallUserType',
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
    handleSwitchChange(isDaily) {
        this.fetch();
    },
    fetch(params = {}) {
    this.loading = true
    this.searchForm = Object.assign({}, this.searchForm, params)
    if (this.searchForm.is_daily) {
        api.getReportRecallDaily(this.searchForm).then(res => {
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
        }else {
        api.getReportRecall(this.searchForm).then(res => {
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

        }
    
    },
    fetchExport(params = {}){
    this.loading = true
    this.searchForm = Object.assign({}, this.searchForm, params)
    api.exportReportRecall(this.searchForm).
    finally(() => {
        this.loading = false
    })
    },
    fetchDailyExport(params = {}){
    this.loading = true
    this.searchForm = Object.assign({}, this.searchForm, params)
    api.exportReportRecallDaily(this.searchForm).
    finally(() => {
        this.loading = false
    })
    },
    summaryMethod({ columns, data }){
    return [
        'Total',
        undefined, undefined, undefined, undefined,
        Number(this.summary['amount'] || 0).toLocaleString(),Number(this.summary['point'] || 0).toLocaleString(),
    ]
    },
    summaryMethodDaily({ columns, data }){
    return [
        'Total',
        Number(this.summary['total_amount'] || 0).toLocaleString(), Number(this.summary['amount'] || 0).toLocaleString(),
        Number(this.summary['difference'] || 0).toLocaleString(),
    ]
    },
},
}
</script>

<style scoped>

</style>
