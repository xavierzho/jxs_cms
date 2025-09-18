<template>
  <div>
    <search-table
      ref="searchTable"
      :loading="loading"
      :tableData="tableData"
      :total="total"
      :page="searchForm.page"
      :page-size="searchForm.page_size"
      height="calc(100% - 160px)"
      @fetch="fetch"
      style="height: 90vh"
    >
      <template v-slot:table-search>
        <el-form :inline="true" :model="searchForm">
          <el-form-item>
            <el-input :placeholder="$t('user.UserID')" v-model="searchForm.user_id" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <el-select v-model="searchForm.module" clearable :placeholder="$t('user.Module')">
              <el-option :label="$t('option.OperationLogModule[0]')" value="router"></el-option>
              <el-option :label="$t('option.OperationLogModule[1]')" value="roles"></el-option>
              <el-option :label="$t('option.OperationLogModule[2]')" value="users"></el-option>
              <el-option :label="$t('option.OperationLogModule[3]')" value="balance_log"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item v-if="searchForm.module!==''">
            <el-input :placeholder="$t('user.ModuleID')" v-model="searchForm.module_id" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetch({page:1})">{{ $t('common.Search') }}</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-slot:table-column>
        <el-table-column type="expand">
          <template slot-scope="props">
            <div class="flex-center" style="padding: 0 50px">
              <el-table
                border
                :data="[jsonToObj(props.row.request)]"
                style="width: 80%"
              >
                <el-table-column v-for="key in Object.keys(jsonToObj(props.row.request))" :key="key" :prop="key" :label="key"
                ></el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="ID"></el-table-column>
        <el-table-column prop="User.user_name" :label="$t('user.Operator')"></el-table-column>
        <el-table-column prop="operation" :label="$t('user.Operation')"></el-table-column>
        <el-table-column prop="module_name" :label="$t('user.Module')"></el-table-column>
        <el-table-column prop="module_id" :label="$t('user.ModuleID')"></el-table-column>
        <el-table-column prop="created_at" :label="$t('user.CreateTime')"></el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import SearchTable from '@/components/Table/SearchTable'
import api from '@/api'

export default {
  name: 'operation_log_list',
  components: {
    SearchTable,
  },
  data() {
    return {
      tableData: [],
      loading: false,
      total: 0,
      searchForm: {
        user_id: '',
        module: '',
        module_id: '',
        page: 1,
        page_size: 20,
      },
    }
  },
  watch: {
    'searchForm.module'(val) {
      if (val === '') {
        this.searchForm.module_id = ''
      }
    },
  },
  created() {
    this.fetch({ page: 1 })
  },
  methods: {
    fetch(params = {}) {
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.loading = true
      api.getLogList(this.searchForm).then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total
      }).catch(() => {
        this.tableData = []
        this.total = 0
      }).finally(() => {
        this.loading = false
      })
    },
    jsonToObj(json) {
      if (json === '') {
        return {}
      }
      try {
        let obj = JSON.parse(json)
        let newObj = {}
        Object.keys(obj).forEach(key => {
          newObj[key] = obj[key].join(',')
        })
        return newObj
      } catch (e) {
        return {}
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.flex-center {
  display: flex;
  justify-content: center;
}
</style>
