export default {
  analysis: "广告分析",

  cost: "花费",
  exposure_cnt: "曝光次数",
  exposure_user: "曝光人数",
  exposure_per_time: "人均曝光次数",
  click_cnt: "点击次数",
  click_user: "点击人数",
  click_per_cost: "点击均价",
  click_rate: "点击率",

  validMsg: {
    cost: "请输入广告总花费[0, +∞)",
    exposure_cnt: "请输入曝光次数[0, +∞)",
    exposure_user: "请输入曝光人数[0, +∞)",
    click_cnt: "请输入点击次数[0, +∞)",
    click_user: "请输入点击人数[0, +∞)",

    exposure_cnt_ge_user: "曝光次数应大于等于曝光人数",
    exposure_cnt_ge_click_cnt: "曝光次数应大于等于点击次数",
    click_cnt_ge_user: "点击次数应大于等于人数",
    exposure_user_ge_click_user: "曝光人数应大于等于点击人数",
  },
}
