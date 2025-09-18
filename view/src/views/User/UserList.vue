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
      @sortChange="sortChange"
      style="height: 90vh"
    >
      <template v-slot:table-search>
        <el-form :inline="true" :model="searchForm">
          <el-form-item>
            <el-input :placeholder="$t('user.Name')" v-model="searchForm.name" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <el-input :placeholder="$t('user.UserName')" v-model="searchForm.user_name" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetch({page:1})">{{ $t('common.Search') }}</el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="success" @click="showUserDialog('create')" v-perm="'management_user_create'">{{ $t('common.Create') }}</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-slot:table-column>
        <el-table-column prop="id" label="ID" width="50px"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.UserName')"></el-table-column>
        <el-table-column prop="name" :label="$t('user.Name')" min-width="120px"></el-table-column>
        <el-table-column prop="email" :label="$t('user.Email')" min-width="140px"></el-table-column>
        <el-table-column prop="roles" :label="$t('user.Role')" width="140px" sortable="custom">
          <template slot-scope="scope">
            <div class="child-pd-10">
              <el-tag type="danger" v-for="role in scope.row.roles" :key="role.id">{{ role.name }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="is_lock" :label="$t('user.IsEnabled')" width="90px">
          <template v-slot="scope">
            <el-tag :type=" scope.row.is_lock === 1 ? 'info':'success'">
              {{ scope.row.is_lock === 1 ? 'false' : 'true' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_logon_time" :label="$t('user.LastLogonTime')" width="140px"></el-table-column>
        <el-table-column prop="created_at" :label="$t('user.CreatedAt')" width="140px" sortable="custom"></el-table-column>
        <el-table-column prop="updated_at" :label="$t('user.UpdatedAt')" width="140px"></el-table-column>
        <el-table-column fixed="right" :label="$t('common.Operation')" width="80px">
          <template slot-scope="scope">
            <el-button type="warning" @click="showUserDialog('edit',scope.row)" v-perm="'management_user_update'">{{ $t('common.Edit') }}</el-button>
          </template>
        </el-table-column>
      </template>
    </search-table>

    <!-- 创建编辑弹窗 -->
    <el-dialog
      :title="dialog.mode==='create'?'Create':'Edit'"
      :visible.sync="dialog.visible"
      width="30%"
    >
      <user-form :mode="dialog.mode" v-model="dialog.form" ref="userForm"></user-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialog.visible = false">{{ $t('common.Cancel') }}</el-button>
        <el-button type="primary" :loading="dialog.btnLoading" @click="submitUserForm">{{ $t('common.Confirm') }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import DatePicker from '@/components/DatePicker'
import SearchTable from '@/components/Table/SearchTable'
import UserForm from '@/components/Form/UserForm'
import { mapState } from 'vuex'
import api from '@/api'

export default {
  name: 'user_list',
  components: {
    SearchTable, DatePicker, UserForm,
  },
  data() {
    return {
      tableData: [],
      loading: false,
      total: 0,
      searchForm: {
        name: '',
        user_name: '',
        order_field: '',
        order: '',
        page: 1,
        page_size: 20,
      },
      dialog: {
        mode: 'create',
        visible: false,
        btnLoading: false,
        form: {
          name: '',
          user_name: '',
          email: '',
          password: '',
          is_lock: false,
          role_id_list: [],
        },
      },
    }
  },
  computed: {
    rules() {
      let rule = {
        name: [
          { required: true, message: this.$t('user.validMsg.NameRequired'), trigger: 'blur' },
          { min: 2, max: 100, message: this.$t('user.validMsg.NameLenLimit'), trigger: 'blur' },
        ],
        email: [
          { required: true, message: this.$t('user.validMsg.EmailRequired'), trigger: 'blur' },
          { min: 2, max: 100, message: this.$t('user.validMsg.EmailLenLimit'), trigger: 'blur' },
        ],
      }
      if (this.dialog.form.password !== '') {
        rule = Object.assign({}, rule, {
          password: [
            { min: 6, max: 20, message: this.$t('user.validMsg.PasswordLenLimit'), trigger: 'blur' },
          ],
        })
      }
      if (this.dialog.mode === 'create') {
        rule = Object.assign({}, rule, {
          user_name: [
            { required: true, message: this.$t('user.validMsg.UsernameRequired'), trigger: 'blur' },
            { min: 2, max: 100, message: this.$t('user.validMsg.UsernameRequired'), trigger: 'blur' },
          ],
          password: [
            { required: true, message: this.$t('user.validMsg.PasswordRequired'), trigger: 'blur' },
            { min: 6, max: 20, message: this.$t('user.validMsg.PasswordLenLimit'), trigger: 'blur' },
          ],
        })
      }
      return rule
    },
    ...mapState({
      roleOptions: (state) => state.option.role,
    }),
  },
  created() {
    this.fetch({ page: 1 })
    if (this.roleOptions.length === 0) {
      this.$store.dispatch('option/update', 'role')
    }
  },
  methods: {
    fetch(params = {}) {
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.loading = true
      api.getUserList(this.searchForm).then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
      }).catch(() => {
        this.tableData = []
        this.total = 0
      }).finally(() => {
        this.loading = false
      })
    },
    showUserDialog(mode, row = {}) {
      this.dialog.mode = mode
      this.dialog.visible = true
      this.$nextTick(() => {
        this.$refs.userForm.resetForm()
        if (mode === 'edit') {
          Object.keys(this.dialog.form).forEach(key => {
            if (row.hasOwnProperty(key) && row[key] !== null) {
              this.dialog.form[key] = row[key]
            }
          })
          this.dialog.form.is_lock = row.is_lock === 1
          this.dialog.form.role_id_list = row.roles ? row.roles.map(el => el.id) : []
          this.dialog.form.id = row.id
        }
      })
    },
    submitUserForm() {
      this.dialog.btnLoading = true
      this.$refs.userForm.submitUserForm().then(() => {
        this.dialog.visible = false
      }).finally(() => {
        this.fetch()
        this.dialog.btnLoading = false
      })
    },
    sortChange(val) {
      if (val.order == null) {
        return this.fetch({ order_field: '', order: '', page: 1 })
      }
      this.fetch({ order_field: val.prop, order: val.order, page: 1 })
    },
  },
}
</script>

<style lang="scss" scoped>
.child-pd-10 > * {
  margin: 3px 3px;
}
</style>
