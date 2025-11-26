<template>
  <div>
    <el-table v-loading="loading" :data="data" border stripe :max-height="maxHeight" style="width: 100%"
              :span-method="spanMethod">
      <slot name="table-top-column"></slot>

      <el-table-column prop="award_type" :label="$t('reward.award_type')" min-width="90px" align="center">
        <template v-slot="detail">
          {{ detail.row.award_type }}
        </template>
      </el-table-column>

      <el-table-column prop="award_name" :label="$t('reward.award_name')" min-width="90px" align="center">
        <template v-slot="detail">
          {{ detail.row.award_name }}
        </template>
      </el-table-column>

      <el-table-column prop="reward_value_item" :label="$t('reward.reward_value_item')" min-width="90px" align="center">
        <template v-slot="detail">
          {{ detail.row.reward_value_item }}
        </template>
      </el-table-column>

      <el-table-column prop="reward_value_cost_award_point" :label="$t('reward.reward_value_cost_award_point')" min-width="90px" align="center">
        <template v-slot="detail">
          {{ detail.row.reward_value_cost_award_point }}
        </template>
      </el-table-column>

      <el-table-column prop="award_num" :label="$t('reward.award_num')" min-width="90px" align="center">
        <template v-slot="detail">
          {{ detail.row.award_num }}
        </template>
      </el-table-column>

      <slot name="table-tail-column"></slot>
    </el-table>
  </div>
</template>

<script>
export default {
  name: "AwardDetail",

  props: {
    loading: {
      type: Boolean,
      default: false,
    },

    data: {
      type: Array,
    },

    showLevelName: {
      type: Boolean,
      default: false,
    },

    maxHeight: {
      type: String,
      default: '600px',
    },

    spanMethod: {
      type: Function,
    },
  },
}
</script>

<style scoped>

</style>
