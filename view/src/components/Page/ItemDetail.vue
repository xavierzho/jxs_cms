<template>
<div>
  <el-table v-loading="loading" :data="data" border stripe :max-height="maxHeight" style="width: 100%" :span-method="spanMethod">
    <slot name="table-top-column"></slot>
    <el-table-column prop="item_id" :label="$t('inquire.item.item_id')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.item_id}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column prop="item_name" :label="$t('inquire.item.item_name')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.item_name}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column v-if="showLevelName" prop="level_name" :label="$t('inquire.item.level_name')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.level_name}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column :label="$t('inquire.item.cover_thumb')" min-width="80px" align="center">
      <template v-slot="detail">
        <el-image v-if="detail.row.item_id !== '0'" :src="detail.row.cover_thumb" style="width: 75px; height: 75px"></el-image>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column prop="show_price" :label="$t('inquire.item.show_price')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.show_price | localeNum2f}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column prop="inner_price" :label="$t('inquire.item.inner_price')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.inner_price | localeNum2f}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <el-table-column prop="recycling_price" :label="$t('inquire.item.recycling_price')" min-width="90px" align="center">
      <template v-slot="detail">
        <template v-if="detail.row.item_id !== '0'" >{{detail.row.recycling_price | localeNum2f}}</template>
        <template v-else>-</template>
      </template>
    </el-table-column>
    <slot name="table-tail-column"></slot>
  </el-table>
</div>
</template>

<script>
export default {
  name: "ItemDetail",

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
