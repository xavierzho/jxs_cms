<template>
  <el-date-picker
    v-model="datePickerModel"
    :type="type"
    range-separator="-"
    :start-placeholder="startPlaceholder"
    :end-placeholder="endPlaceholder"
    :value-format="valueFormat"
    :picker-options="pickerOptions"
    :unlink-panels="true"
    :clearable="clearable"
  >
  </el-date-picker>
</template>

<script>
import moment from 'moment'

export default {
  name: 'DatePicker',
  model: {
    prop: 'value',
    event: 'change',
  },
  props: {
    value: {
      type: [String, Array],
      default: '',
    },
    type: {
      type: String,
      default: 'date',
    },
    startPlaceholder: {
      type: String,
      default: 'start',
    },
    endPlaceholder: {
      type: String,
      default: 'end',
    },
    valueFormat: {
      type: String,
      default: 'yyyy-MM-dd',
    },
    clearable: {
      type: Boolean,
      default: true,
    },
    utcOffset: {
      type: String,
      default: '',
    },
    optionRangeType: {
      type: String,
      default: 'before',
    },
  },
  computed: {
    datePickerModel: {
      get: function () {
        return this.value
      },
      set: function (newValue) {
        if (newValue === null) {
          newValue = []
        }
        this.$emit('change', newValue)
      },
    },
    pickerOptions() {
      let shortcuts = [
        {
          text: this.$t('option.DatePickerOptions[0]'),
          onClick: (picker) => {
            const end = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            const start = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            if (this.optionRangeType === 'before') {
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            } else {
              end.setTime(end.getTime() + 3600 * 1000 * 24 * 7)
            }
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: this.$t('option.DatePickerOptions[1]'),
          onClick: (picker) => {
            const end = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            const start = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            if (this.optionRangeType === 'before') {
              start.setMonth(start.getMonth() - 1)
            } else {
              end.setMonth(end.getMonth() + 1)
            }
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: this.$t('option.DatePickerOptions[2]'),
          onClick: (picker) => {
            const end = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            const start = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            if (this.optionRangeType === 'before') {
              start.setMonth(start.getMonth() - 3)
            } else {
              end.setMonth(end.getMonth() + 3)
            }
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: this.$t('option.DatePickerOptions[3]'),
          onClick: (picker) => {
            const end = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            const start = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            if (this.optionRangeType === 'before') {
              start.setMonth(start.getMonth() - 6)
            } else {
              end.setMonth(end.getMonth() + 6)
            }
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: this.$t('option.DatePickerOptions[4]'),
          onClick: (picker) => {
            const end = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            const start = new Date(moment().utcOffset(this.utcOffset).format('YYYY-MM-DD HH:mm:ss'))
            if (this.optionRangeType === 'before') {
              start.setFullYear(start.getFullYear() - 1)
            } else {
              end.setFullYear(end.getFullYear() + 1)
            }
            picker.$emit('pick', [start, end])
          },
        },
      ]
      return this.type.includes('range') ? {
        shortcuts: shortcuts,
      } : {}
    },
  },
  methods: {},
}
</script>

<style scoped>

</style>
