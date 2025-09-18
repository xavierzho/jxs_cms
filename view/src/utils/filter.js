import moment from 'moment'
import Decimal from "decimal.js";

export default [
  {
    name: 'indiaNum',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined') {
        return '0'
      }
      return Number(value).toLocaleString('en-IN')
    },
  },
  {
    name: 'localeNum',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined' || isNaN(value)) {
        return '0'
      }
      return Number(value).toLocaleString()
    },
  },
  {
    name: 'localeNum2f',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined' || isNaN(value)) {
        return '0'
      }
      return Number(new Decimal(value).toFixed(2, Decimal.ROUND_HALF_UP)).toLocaleString()
    },
  },
  {
    name: 'rate',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined' || isNaN(value)) {
        return '0%'
      }
      return value + '%'
    },
  },
  {
    name: 'rate2f',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined' || isNaN(value)) {
        return '0.00%'
      }
      return value
    },
  },
  {
    name: 'date',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined') {
        return '-'
      }
      return moment(value).format('YYYY-MM-DD')
    },
  },
  {
    name: 'Duration',
    func: function (value) {
      if (value === null || value === '' || typeof value === 'undefined') {
        return '0s'
      }
      return value
    },
  }
]
