<template>
  <div>
    <el-table
      ref="table"
      v-loading="loading"
      :data="tableData"
      :height="height"
      border
      stripe
      style="width: 100%"
      :show-summary="showSummary"
      :summary-method="summaryMethod"
      @selection-change="handleSelectionChange"
      @sort-change="handleSortChange"
      @expand-change="expandChange"
    >
      <!-- Table字段 -->
      <slot name="table-column"></slot>
    </el-table>
    <div class="pagebox">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="page"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :disabled="loading"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ListTable',
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
    tableData: {
      type: Array,
      default() {
        return []
      },
    },
    total: {
      type: Number,
      required: true,
    },
    page: {
      type: Number,
      default: 1,
    },
    pageSize: {
      type: Number,
      default: 50,
    },
    height: {
      type: String,
      default() {
        return '600px'
      },
    },
    showSummary: {
      type: Boolean,
      default: false,
    },
    summaryMethod: Function,
  },
  data() {
    return {}
  },
  updated() {
    this.doLayout()
  },
  methods: {
    doLayout(){
      this.$nextTick(() => {
        this.$refs['table'].doLayout()
      })
    },
    handleSizeChange(pageSize) {
      this.$emit('fetch', {
        page_size: pageSize,
        page: 1,
      })
    },
    handleCurrentChange(currentPage) {
      this.$emit('fetch', {
        page: currentPage,
      })
      this.$refs.table.bodyWrapper.scrollTop = 0
    },
    handleSelectionChange(val) {
      this.$emit('selectionChange', val)
    },
    handleSortChange(val) {
      if (val.order !== null) {
        val.order = val.order === 'ascending' ? 'ASC' : 'DESC'
      }
      this.$emit('sortChange', val)
    },
    expandChange(row, expandedRows) {
      this.$emit('expandChange', row, expandedRows)
    },
  },
}
</script>

<style lang="scss" scoped>
.pagebox {
  margin-top: 20px;
}
</style>
