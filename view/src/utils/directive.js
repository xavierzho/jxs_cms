import { checkPermission } from '@/utils/auth'
import Vue from 'vue'

export function bindingCheckPermission(el, binding) {
  const { value } = binding

  let hasPermission = checkPermission(value)
  if (!hasPermission) {
    Vue.nextTick(() => {
      el.parentNode && el.parentNode.removeChild(el)
    })
  }
}

