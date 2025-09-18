<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="130px">
      <el-form-item :label="$t('user.Name')" prop="name">
        <el-input v-model="form.name" :placeholder="$t('user.Name')" clearable></el-input>
      </el-form-item>
      <el-form-item :label="$t('user.UserName')" prop="user_name" v-if="mode==='create'">
        <el-input v-model="form.user_name" :placeholder="$t('user.UserName')" clearable></el-input>
      </el-form-item>
      <el-form-item :label="$t('user.Email')" prop="email">
        <el-input v-model="form.email" :placeholder="$t('user.Email')" clearable></el-input>
      </el-form-item>
      <template v-if="mode === 'self-update'">
        <el-form-item :label="$t('user.NewPassword')" prop="new_password">
          <el-input :placeholder="$t('user.NewPassword')" v-model="form.new_password" show-password></el-input>
        </el-form-item>
      </template>
      <template v-else>
        <el-form-item :label="$t('user.IsEnabled')" prop="is_lock">
          <el-switch v-model="form.is_lock" active-color="#13ce66"
                     :active-value="false"
                     :inactive-value="true"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('user.Password')" prop="password">
          <el-input :placeholder="$t('user.Password')" v-model="form.password" show-password></el-input>
        </el-form-item>
      </template>
      <el-form-item :label="$t('user.Role')" prop="role_id_list" v-if="mode!=='self-update'">
        <el-checkbox-group v-model="form.role_id_list">
          <el-checkbox :label="option.id" v-for="option in roleOptions" :key="option.id">{{ option.name }}
          </el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import api from '@/api'

export default {
  name: 'UserForm',
  model: {
    prop: 'form',
    event: 'change',
  },
  props: {
    mode: {
      type: String,
      default: 'create',
    },
    form: {
      type: Object,
      default: {},
    },
  },
  data() {
    return {}
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

      if (this.form.hasOwnProperty('password') && this.form.password !== '') {
        rule = Object.assign({}, rule, {
          password: [
            { min: 6, max: 20, message: this.$t('user.validMsg.PasswordLenLimit'), trigger: 'blur' },
          ],
        })
      }
      if (this.mode === 'create') {
        rule = Object.assign({}, rule, {
          user_name: [
            { required: true, message: this.$t('user.validMsg.UsernameRequired'), trigger: 'blur' },
            { min: 2, max: 100, message: this.$t('user.validMsg.UsernameLenLimit'), trigger: 'blur' },
          ],
          password: [
            { required: true, message: this.$t('user.validMsg.PasswordRequired'), trigger: 'blur' },
            { min: 6, max: 20, message: this.$t('user.validMsg.PasswordLenLimit'), trigger: 'blur' },
          ],
        })
      }
      if (this.mode === 'self-update') {
        if (this.form.hasOwnProperty('new_password') && this.form.new_password !== '') {
          rule = Object.assign({}, rule, {
            new_password: [
              { min: 6, max: 20, message: this.$t('user.validMsg.PasswordLenLimit'), trigger: 'blur' },
            ],
          })
        }
      }
      return rule
    },
    ...mapState({
      roleOptions: (state) => state.option.role,
    }),
  },
  created() {
    if (this.roleOptions.length === 0) {
      this.$store.dispatch('option/update', 'role')
    }
  },
  methods: {
    submitUserForm() {
      return new Promise((resolve, reject) => {
        this.$refs.form.validate(valid => {
          if (valid) {
            let apiName = 'createUser'
            switch (this.mode) {
              case 'create':
                apiName = 'createUser'
                break
              case 'edit':
                apiName = 'updateUser'
                break
              case 'self-update':
                apiName = 'updateUserSelf'
                break
            }
            this.btnLoading = true
            api[apiName](this.form).then(res => {
              this.$message.success('Success')
              resolve(res)
              this.resetForm()
            }).catch(err => {
              reject(err)
            })
          } else {
            resolve()
          }
        })
      })
    },
    resetForm() {
      if (this.form.hasOwnProperty('id')) {
        delete this.form.id
      }
      this.$refs.form.resetFields()
      this.form.user_name = ''
    },
  },
}
</script>

<style scoped>

</style>
