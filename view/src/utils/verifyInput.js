import i18n from "@/utils/i18n";

function checkPositiveInteger(rule, value, callback)  {
  if (!Number.isInteger(Number(value)) || value < 0) {
    callback(new Error(i18n.t("verify.VerifyTip.PositiveIntegerTip")))
    return
  }
  callback()
}

function checkoutPositiveNumber(rule, value, callback) {
  if (value < 0) {
    callback(new Error(i18n.t("verify.VerifyTip.PositiveNumberTip")))
    return
  }
  callback()
}

function checkoutBlankString(rule, value, callback) {
  if (value===""){
    callback(new Error(i18n.t("verify.VerifyTip.SelectOptionTip")))
    return
  }
  callback()
}

function checkoutDecimalDigits(value, digits){
  let valueStr = value.toString()
  if(valueStr.indexOf('.')<0){
    return false
  }
  let re = /0{1,}$/
  valueStr = valueStr.replace(re, '')
  return (valueStr.length - valueStr.indexOf('.') - 1) > digits;
}

export {
  checkPositiveInteger,
  checkoutPositiveNumber,
  checkoutBlankString,
  checkoutDecimalDigits,
}
