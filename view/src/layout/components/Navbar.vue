<template>
  <div class="navbar">
    <hamburger :is-active="sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar"/>

    <breadcrumb class="breadcrumb-container"/>

    <div class="right-menu">
      <div class="right-menu-item">
        <el-tooltip class="item" effect="dark" content="Tag View Switch" placement="bottom">
          <el-switch
            v-model="openTagView"
            active-color="#13ce66"
            inactive-color="#e4e4e4"
          >
          </el-switch>
        </el-tooltip>
      </div>
      <div class="clock right-menu-item">{{ nowTime }}</div>
<!--      <el-dropdown class="right-menu-item" trigger="click" @command="changeLang" v-perm="'lang_update'">-->
<!--        <span class="el-dropdown-link user-select-none">-->
<!--          {{ $t('common.Language') }}<i class="el-icon-arrow-down el-icon&#45;&#45;right"></i>-->
<!--        </span>-->
<!--        <el-dropdown-menu slot="dropdown">-->
<!--          <el-dropdown-item command="en">English</el-dropdown-item>-->
<!--          <el-dropdown-item command="zh">简体中文</el-dropdown-item>-->
<!--        </el-dropdown-menu>-->
<!--      </el-dropdown>-->

      <el-dropdown class="avatar-container right-menu-item" trigger="click" @command="handleUser">
        <span class="username user-select-none">
          {{ name }}<i class="el-icon-arrow-down el-icon--right"></i>
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="home">
            Home
          </el-dropdown-item>
          <el-dropdown-item command="self-update">
            Edit
          </el-dropdown-item>
          <el-dropdown-item command="logout" divided>
            <span style="display:block;">Log Out</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
    <el-dialog
      title="Modify Personal Info"
      :visible.sync="dialog.visible"
      width="30%"
    >
      <user-form mode="self-update" v-loading="dialog.loading" v-model="dialog.form" ref="userForm"></user-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialog.visible = false">{{ $t('common.Cancel') }}</el-button>
        <el-button type="primary" :loading="dialog.btnLoading" @click="submitUserForm">{{ $t('common.Confirm') }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters, mapState } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import UserForm from '@/components/Form/UserForm'
import { checkPermission } from '@/utils/auth'
import api from '@/api'
import moment from 'moment'

export default {
  components: {
    Breadcrumb, Hamburger, UserForm,
  },
  data() {
    return {
      dialog: {
        loading: false,
        btnLoading: false,
        visible: false,
        form: {
          name: '',
          email: '',
          new_password: '',
        },
        rules: {},
      },
      polling: null,
      nowTime: moment().utcOffset('+05:30').format('YYYY-MM-DD HH:mm:ss'),
    }
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'name',
    ]),
    ...mapState({
      tagsView: state => state.settings.tagsView,
    }),
    openTagView: {
      get: function () {
        return this.tagsView
      },
      set: function (newValue) {
        if (newValue === false){
          this.$store.dispatch('tagsView/delAllVisitedViews')
          this.$store.dispatch('tagsView/delAllCachedViews')
        }
        this.$store.dispatch('settings/setTagsView', newValue)
      },
    },
  },
  created() {
    if (typeof this.$store.state.app.lang === 'undefined') {
      this.$store.dispatch('app/toggleLang', this.$i18n.locale)
    } else {
      if (checkPermission('lang_update')) {
        this.$i18n.locale = this.$store.state.app.lang
      } else {
        this.changeLang('zh')
      }
    }
    this.initTime()
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.polling)
    this.polling = null
    next()
  },
  methods: {
    initTime() {
      this.polling = setInterval(() => {
        this.nowTime = moment().utcOffset('+05:30').format('YYYY-MM-DD HH:mm:ss')
      }, 500)
    },
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      await this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    },
    changeLang(lang) {
      this.$store.dispatch('app/toggleLang', lang)
      this.$i18n.locale = lang
    },
    handleUser(command) {
      switch (command) {
        case 'logout':
          this.logout()
          break
        case 'self-update':
          this.openUserDialog()
          break
        case 'home':
          this.$router.push('/')
          break
      }
    },
    openUserDialog() {
      // 获取用户信息
      this.dialog.loading = true
      this.dialog.visible = true
      api.getUserDetail().then(res => {
        Object.keys(this.dialog.form).forEach(key => {
          if (res.data.hasOwnProperty(key)) {
            this.dialog.form[key] = res.data[key]
          }
        })
      }).finally(() => {
        this.dialog.loading = false
      })
    },
    submitUserForm() {
      this.dialog.btnLoading = true
      this.$refs.userForm.submitUserForm().then(res => {
        let user = {
          id: res.data.id,
          name: res.data.name,
          email: res.data.email,
        }
        this.$store.commit('user/SET_USER', user)
        this.dialog.visible = false
      }).finally(() => {
        this.dialog.btnLoading = false
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.user-select-none {
  user-select: none;
}

.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, .08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .clock {
      cursor: default !important;
      user-select: none;
    }

    .right-menu-item {
      cursor: pointer;
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 14px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .username {
        cursor: pointer;
        width: 40px;
        height: 40px;
        border-radius: 10px;
      }

      .el-icon-caret-bottom {
        cursor: pointer;
        position: absolute;
        right: -20px;
        top: 25px;
        font-size: 12px;
      }
    }
  }
}
</style>
