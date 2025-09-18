import i18nMessage from '@/i18n'
import VueI18n from 'vue-i18n'
import Vue from 'vue'

Vue.use(VueI18n)

export default new VueI18n({
  locale: 'zh', // 设置地区
  messages: i18nMessage // 设置地区信息
})
