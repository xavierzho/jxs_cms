<template>
  <div>
    <div class="input-number-range" :class="{ 'is-disabled': disabled }">
      <div class="flex">
        <div class="from">
          <el-input
            ref="input_from"
            v-model="userInputForm"
            :disabled="disabled"
            :clearable="clearable"
            :placeholder="$t('common.MinValue')"
            @blur="handleBlurFrom"
            @focus="handleFocusFrom"
            @input="handleInputFrom"
            @change="handleInputChangeFrom"
          ></el-input>
        </div>
        <div class="center">
          <span>-</span>
        </div>
        <div class="to">
          <el-input
            ref="input_to"
            v-model="userInputTo"
            :disabled="disabled"
            :clearable="clearable"
            :placeholder="$t('common.MaxValue')"
            @blur="handleBlurTo"
            @focus="handleFocusTo"
            @input="handleInputTo"
            @change="handleInputChangeTo"
          ></el-input>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'InputNumberRange',

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

  data() {
    return {
      userInputForm: null,
      userInputTo: null
    }
  },

  watch: {
    value: {
      immediate: true,
      handler(value) {
        /** 初始化范围 */
        if (value instanceof Array && this.precision !== undefined) {
          this.userInputForm = value[0]
          this.userInputTo = value[1]
        }
      }
    }
  },

  methods: {
    handleBlurFrom(event) {
      this.$emit('blurfrom', event)
    },

    handleFocusFrom(event) {
      this.$emit('focusfrom', event)
    },

    handleBlurTo(event) {
      this.$emit('blurto', event)
    },

    handleFocusTo(event) {
      this.$emit('focusto', event)
    },

    handleInputFrom(value) {
      this.$emit('inputfrom', value)
    },

    handleInputTo(value) {
      this.$emit('inputto', value)
    },

    // from输入框change事件
    handleInputChangeFrom(value) {
      // 如果是非数字空返回null
      if (isNaN(value) || value === '') {
        this.$emit('input', [null, this.userInputTo])
        this.$emit('changefrom', this.userInputForm)
        return
      }

      // 初始化数字精度
      this.userInputForm = this.formatValue(value)

      // 如果from > to 将from值替换成to
      if (typeof this.userInputTo === 'number') {
        this.userInputForm =
          parseFloat(this.userInputForm) <= parseFloat(this.userInputTo)
            ? this.userInputForm
            : this.userInputTo
      }
      this.$emit('input', [this.userInputForm, this.userInputTo])
      this.$emit('changefrom', this.userInputForm)
    },

    // to输入框change事件
    handleInputChangeTo(value) {
      // 如果是非数字空返回null
      if (isNaN(value) || value === '') {
        this.$emit('input', [this.userInputForm, null])
        this.$emit('changeto', this.userInputTo)
        return
      }

      // 初始化数字精度
      this.userInputTo = this.formatValue(value)

      // 如果to < tfrom 将to值替换成from
      if (typeof this.userInputForm === 'number') {
        this.userInputTo =
          parseFloat(this.userInputTo) >= parseFloat(this.userInputForm)
            ? this.userInputTo
            : this.userInputForm
      }
      this.$emit('input', [this.userInputForm, this.userInputTo])
      this.$emit('changeto', this.userInputTo)
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
<style lang="scss" scoped>
// 取消element原有的input框样式
::v-deep .el-input--mini .el-input__inner {
  border: 0px;
  margin: 0;
  padding: 0 15px;
  background-color: transparent;
}
.input-number-range {
  background-color: #fff;
  //border: 1px solid #dcdfe6;
  //border-radius: 4px;
}
.flex {
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: center;
  align-items: center;
  .center {
    margin-top: 1px;
  }
}
.is-disabled {
  background-color: #eef0f6;
  border-color: #e4e7ed;
  color: #c0c4cc;
  cursor: not-allowed;
}
</style>
