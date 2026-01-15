<template>
  <div :class="'component '+ className">
    <div class="title">{{title}}</div>
    <div class="value">
      <div class="value-today">
        <div>{{format(today)}}</div>
        <div class="value-today-classify">{{classifyStr}}</div>
      </div>
      <div v-if="!isSummary" :class="'value-diff '+valueColorClass(diff)">{{valueDiffStr(diff)}}</div>
      <div v-if="!isSummary" class="value-yesterday">昨天&nbsp {{format(yesterday)}}</div>
    </div>
  </div>
</template>

<script>
import Decimal from "decimal.js";

export default {
  name: "DashboardValue",

  props:{
    className:{
      type: String,
      required: true,
    },

    title: {
      type: String,
      required: true,
    },

    value: {
      type: Array,
      required: true,
      validator(val) {
        return val.length === 2
      }
    },

    valueClassify: {
      type: Array,
      validator(val) {
        return val.length === 2
      }
    },

    valueType: {
      type: String,
      default: 'value',
    },

    isSummary: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      today: 0,
      yesterday: 0,
      diff: 0,
      classifyStr: '',
    }
  },

  watch:{
    value: {
      immediate: true,
      handler(value) {
        this.today = value[0]
        this.yesterday = value[1]
        this.diff = value[0] - value[1]
      }
    },
    valueClassify: {
      immediate: true,
      handler(value) {
        if (value instanceof Array){
          this.classifyStr = '新 ' + this.format(value[0]) + ' / 旧 ' + this.format(value[1])
        }
      }
    },
  },

  methods:{
    format(value){
      if (this.valueType === 'value'){
        return Number(value).toLocaleString()
      }else if(this.valueType === 'value2f'){
        return Number(new Decimal(value).toFixed(2, Decimal.ROUND_HALF_UP)).toLocaleString()
      }else if(this.valueType === 'percentage'){
        return Number(new Decimal(value).toFixed(2, Decimal.ROUND_HALF_UP)).toLocaleString()+'%'
      }
      return value
    },
    valueColorClass(value){
      if (value > 0){
        return 'value-green'
      }else if(value < 0){
        return 'value-red'
      }else{
        return ''
      }
    },
    valueDiffStr(value){
      let valueStr = this.format(value)
      if (value>0){
        return '+'+valueStr
      }else if (value === 0){
        return ''
      }
      return valueStr
    }
  },

}
</script>

<style lang="scss" scoped>
.component{
  position: relative;
  width: 300px;
  height: 100px;
  margin-left: 40px;
  padding: 8px;
  border-radius: 4px;
  box-shadow: 0 0 3px gray;
  display: flex;
  flex-direction: column;
  flex-grow: 0;
  flex-shrink: 0;
}
.component:first-child{
  margin-left: 0;
}

.title{
  color: #909399;
}

.value{
  display: flex;
  align-items: baseline;
}

.value-today{
  font-size: 3rem;
}

.value-today-classify{
  font-size: 0.7rem;
  color: #909399;
}

.value-diff{
  position: relative;
  left: 0.5rem;
}

.value-yesterday {
  position: absolute;
  right: 0;
  bottom: 5px;
  width: 35%;
  font-size: 0.9rem;
  color: #909399;
}

.value-red{
  color: #F56C6C;
}

.value-green{
  color: #67C23A;
}

.recharge_amount{
  border-bottom: #409EFF 2px solid;
}

.recharge_amount_huifu{
  border-bottom: #6190E8 2px solid;
}

.recharge_amount_ali{
  border-bottom: #10A3E7 2px solid;
}

.draw_amount{
  border-bottom: #E6A23C 2px solid;
}

.new_user_cnt{
  border-bottom: #67C23A 2px solid;
}

.active_user_cnt{
  border-bottom: #E1F3D8 2px solid;
}

.recharge_user_cnt{
  border-bottom: #409EFF 2px solid;
}

.pating_rate_new{
  border-bottom: #F0F9EB 2px solid;
}

.new_user_cnt_summary{
  border-bottom: #67C23A 2px solid;
}

.recharge_user_cnt_summary{
  border-bottom: #409EFF 2px solid;
}

.recharge_amount_summary{
  border-bottom: #409EFF 2px solid;
}

.draw_amount_summary{
  border-bottom: #E6A23C 2px solid;
}

</style>
