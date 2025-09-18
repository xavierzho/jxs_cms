import i18n from "@/utils/i18n";

function confirmBox() {
  return this.$confirm(i18n.t('common.ConfirmOperation'), i18n.t("common.Tips"), {
    confirmButtonText: i18n.t('common.Confirm'),
    cancelButtonText: i18n.t('common.Cancel'),
    type: 'warning'
  })
}

export {confirmBox}


