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
            <el-input :placeholder="$t('user.RoleName')" v-model="searchForm.name" clearable></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetch({page:1})">{{ $t('common.Search') }}</el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="success" @click="showUserDialog('create')" v-perm="'management_role_create'">{{ $t('common.Create') }}</el-button>
          </el-form-item>
        </el-form>
      </template>

      <template v-slot:table-column>
        <el-table-column prop="id" label="ID"></el-table-column>
        <el-table-column prop="name" :label="$t('user.RoleName')"></el-table-column>
        <el-table-column fixed="right" :label="$t('common.Operation')" align="center" width="120px">
          <template slot-scope="scope">
            <el-button type="warning" @click="showUserDialog('edit',scope.row)" v-perm="'management_role_update'">{{ $t('common.Edit') }}</el-button>
          </template>
        </el-table-column>
      </template>
    </search-table>

    <!-- create edit role drawer -->
    <el-drawer
      :title="dialog.mode==='create'?'Create':'Edit'"
      :visible.sync="dialog.visible"
      direction="rtl"
      size="50%"

    >
      <el-form :model="dialog.form" :rules="rules" ref="form" label-width="100px">
        <el-form-item :label="$t('user.RoleName')" prop="name">
          <el-input v-model="dialog.form.name" :placeholder="$t('user.RoleName')" clearable></el-input>
        </el-form-item>
        <el-form-item :label="$t('user.Permission')" >
          <el-tree
          :data="pageList.data"
          show-checkbox
          node-key="node_name"
          :default-checked-keys="selectedNode"
          :props="pageList.defaultProps"
          @check-change="selectNode"
          ref="tree"
          >
          </el-tree>
        </el-form-item>
        <el-form-item>
          <el-button @click="checkedAll">{{$t('common.CheckedAll')}}</el-button>
          <el-button @click="unselectAll">{{$t('common.UnselectAll')}}</el-button>
          <el-button type="primary" :loading="dialog.btnLoading" @click="submitForm">{{$t('common.Confirm')}}</el-button>
          <el-button @click="cancel">{{$t('common.Cancel')}}</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script>
import DatePicker from '@/components/DatePicker'
import SearchTable from '@/components/Table/SearchTable'
import { getPermission, getMenu, removeMenu, setPermission } from '@/utils/cache'
import api from '@/api'
import { mapState, mapGetters } from 'vuex'
import { isArrEqual } from '@/utils/array'

export default {
  name: 'role_list',
  components: {
    SearchTable, DatePicker
  },
  data() {
    return {
      tableData: [],
      loading: false,
      total: 0,
      searchForm: {
        name: '',
        page: 1,
        page_size: 20
      },
      dialog: {
        mode: 'create',
        visible: false,
        btnLoading: false,
        form: {
          name: '',
          permission_id_list: []
        },
        tree: false,
      },
      putForm: {
          permissionIds: {}
      },
      pageList: {
          data:[
          ],
          defaultProps: {
              children: 'children',
              label: 'label'
          },
      },
      selectedNode: [],
      userPermission: {},
      loginUserPermission: [],
    }
  },
  computed: {
    rules() {
      return {
        name: [
          { required: true, message: this.$t('user.validMsg.RoleNameRequired'), trigger: 'blur' },
          { min: 2, max: 100, message: this.$t('user.validMsg.RoleNameLenLimit'), trigger: 'blur' }
        ]
      }
    },
    ...mapState({
      permissionOptions: (state) => {
        return state.option.permission.map(el => {
          return {
            key: el.id,
            label: state.app.lang === 'zh' ? el.description : el.display_name
          }
        })
      },
      pagePermission: (state) => {
        return state.user.pagePermission
      },
    }),
    ...mapGetters({
      menusDict: 'app/menusDict',
      pagePermissionDict: 'user/pagePermissionDict'
    }),
  },
  created() {
    if ((getMenu() || []).length === 0) {
      this.$store.dispatch('app/getMenus')
    }

    if (Object.keys(this.pagePermission).length === 0){
      this.$store.dispatch('user/getPagePermissions')
    }

    this.fetch({ page: 1 })
    if (this.permissionOptions.length === 0) {
      this.$store.dispatch('option/update', 'permission')
    }
  },
  methods: {
    fetch(params = {}) {
      this.searchForm = Object.assign({}, this.searchForm, params)
      this.loading = true
      api.getRoleList(this.searchForm).then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
      }).catch(() => {
        this.tableData = []
        this.total = 0
      }).finally(() => {
        this.loading = false
      })
    },
    clearUserForm() {
      if (this.dialog.form.hasOwnProperty('id')) {
        delete this.dialog.form.id
      }
      this.$nextTick(() => {
        this.$refs.form.resetFields()
      })
    },
    async showUserDialog(mode, row = {}) {
      this.dialog.mode = mode
      this.dialog.visible = true
      this.clearUserForm()
      if (mode === 'edit') {
        this.$nextTick(() => {
          Object.keys(this.dialog.form).forEach(key => {
            if (row.hasOwnProperty(key) && row[key] !== null) {
              this.dialog.form[key] = row[key]
            }
          })
          this.dialog.form.permission_id_list = row.permission ? row.permission.map(el => el.id) : []
          this.dialog.form.id = row.id
        })
      }

      if ((getMenu() || []).length === 0) {
        await this.$store.dispatch('app/getMenus')
      }

      if (Object.keys(this.pagePermission).length === 0){
        await this.$store.dispatch('user/getPagePermissions')
      }
      this.loginUserPermission = getPermission()
      this.pageList.data = getMenu()
      this.createPageTree("tree", this.pageList.data)

      this.userPermission = (mode === 'edit') ? row.permission : []
      this.putForm.permissionIds = {}
      this.selectedNode = this.showSelect(this.pageList.data, [])
      if (this.$refs.tree){
        this.$refs.tree.setCheckedKeys(this.selectedNode)
      }
    },
    submitForm() {
      this.$refs.form.validate(valid => {
        if (valid) {
          let apiName = 'createRole'
          if (this.dialog.mode === 'edit') {
            apiName = 'updateRole'
          }
          this.dialog.btnLoading = true
          this.dialog.form.permission_id_list = Object.keys(this.putForm.permissionIds)
          api[apiName](this.dialog.form).then(res =>  {
            this.$store.dispatch('option/update', 'role')
            this.$message.success('Success')
            this.fetch()
            this.dialog.visible = false
            if (this.dialog.mode === 'edit') {
              // 更新登录用户权限
              setPermission(res.data)
              // 更新菜单页面
              if (!isArrEqual(this.loginUserPermission, res.data)){
                this.$store.dispatch('app/getMenus')
              }
            }
          }).finally(() => {
            this.dialog.btnLoading = false
          })
        }
      })
    },
    checkedAll(){
        this.$refs.tree.setCheckedNodes(this.pageList.data);
    },
    unselectAll(){
      this.$refs.tree.setCheckedKeys([])
    },
    cancel(){
      // this.$refs.tree.setCheckedKeys([]) // 异步的；会调动selectNode；似乎会将多个setCheckedKeys操作合并到一起完成；
      this.putForm.permissionIds = {}
      this.selectedNode = this.showSelect(this.pageList.data, [])
      this.$refs.tree.setCheckedKeys(this.selectedNode)
      this.dialog.visible = false
    },
    // 只有展开的节点变化时会调用该函数
    selectNode(node, selectStatus, sonSelectStatus){
      // 需要迭代地去选中子树 未展开时 无法自动迭代全部树
      if (node.is_permission){
        if (selectStatus){
            this.putForm.permissionIds[node.id] = true
        }else{
            delete this.putForm.permissionIds[node.id]
        }
      }else if (selectStatus && node.children){
        node.children.forEach(item=>{
            this.selectNode(item, true, false)
        })
      } else if (!selectStatus && !sonSelectStatus && node.children){
        // 子节点未展开不会调用该函数！！！
        // 当取消操作时，父节点变化但按该逻辑无法传递到未被选中的子节点
        // 所以取消要进行初始化，而不能直接setCheckedKeys(this.selectedNode)
        node.children.forEach(item=>{
            this.selectNode(item, false, false)
        })
      }
    },
    createPageTree(preNode, data){
      if (data===undefined){return}
      data.forEach((el, _)=>{
        // 设置名字
        if (el.is_permission){ // 权限
          el.label = this.pagePermissionDict[el.name]
          el.node_name = preNode + el.name + el.id
        } else{
          el.label = this.menusDict[el.name]
          el.node_name = preNode + el.name
        }

        if (el.is_permission===undefined && el.children === undefined && this.pagePermission[el.name] !== undefined){
            let children = []
            this.pagePermission[el.name].forEach((permission, _)=> {
              if (permission.is_admin || this.loginUserPermission.indexOf(permission.name) !== -1){
                children.push(permission)
              }
            })
            el.children = children
        }
        this.createPageTree(el.name, el.children)
      } )
      return
    },
    showSelect(data, selectedNode){
      if (data===undefined){return selectedNode}
      data.forEach((el, _)=>{
        if (el.is_permission){
          this.userPermission.forEach(item=>{
            if (item.id === el.id){
              selectedNode.push(el.node_name)
              this.putForm.permissionIds[el.id] = true
            }
          })
        }
        selectedNode = this.showSelect(el.children, selectedNode)
      } )
      return selectedNode
    },
  }
}
</script>

<style  scoped>
/*.child-pd-10 > * {*/
/*  margin: 3px 3px;*/
/*}*/
/deep/ .el-transfer-panel__body{
  height: 500px;
}
/deep/ .el-transfer-panel__list{
  height: 450px;
}
</style>
