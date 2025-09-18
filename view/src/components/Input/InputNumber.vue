<template>
  <div>
    <el-input
      ref="input"
      v-model="innerValue"
      :disabled="disabled"
      :clearable="clearable"
      @blur="handleBlur"
      @focus="handleFocus"
      @input="handleInput"
      @change="handleInputChange"
    ></el-input>
  </div>
</template>

<script>
export default {
  name: 'InputNumber',

  props: {
    // 初始化范围
    value: { required: true },

    range: {
      type: Array,
      default: ()=>{return [-(2**31), 2**31-1]},
      validator(val) {
        return val.length === 2
      }
    },

    // 是否禁用
    disabled: {
      type: Boolean,
      default: false
    },

    clearable: {
      type: Boolean,
      default: false
    },

    // 精度参数
    precision: {
      type: Number,
      default: 0,
      validator(val) {
        return val >= 0 && val === parseInt(val, 10)
      }
    }
  },

  data(){
    return {
      innerValue: null
    }
  },

  watch: {
    value: {
      immediate: true,
      handler(value) {
        this.innerValue = value
      }
    }
  },

  methods: {
    handleBlur(event) {
      this.$emit('blur', event)
    },

    handleFocus(event) {
      this.$emit('focus', event)
    },

    handleInput(value) {
      this.$emit('input', value)
    },

    // from输入框change事件
    handleInputChange(value) {
      // 如果是非数字空返回null
      if (isNaN(value) || value === '') {
        this.$emit('input', null)
        return
      }
      this.$emit('input', this.formatValue(value))
    },

    formatValue(value){
      if (value < this.range[0]){
        return this.range[0]
      }else if (value > this.range[1]){
        return this.range[1]
      }

      return this.setPrecisionValue(value)
    },

    // 设置成精度数字
    setPrecisionValue(value) {
      if (this.precision !== undefined) {
        return this.toPrecision(value, this.precision)
      }
      return null
    },
    // 根据精度保留数字
    toPrecision(num, precision) {
      if (precision === undefined) precision = 0
      return parseFloat(
        Math.round(num * Math.pow(10, precision)) / Math.pow(10, precision)
      )
    },
  }
}
</script>
