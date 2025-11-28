<template>
  <div>
    <search-table
      :loading="loading"
      :table-data="tableData" show-summary :summary-method="summaryMethod"
      :page="searchFormLog.page" :page-size="searchFormLog.pageSize" :total="total"
      :height="'64vh'"
      @fetch="fetch"
      @expandChange="expandChange"
    >
      <template v-slot:table-search>
        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('common.DateTime')">
            <date-picker v-model="searchFormLog.date_time_range" type="datetimerange" value-format="yyyy-MM-dd HH:mm:ss"
                         :clearable="false"></date-picker>
          </el-form-item>
          <el-form-item>
            <span slot="label">{{ $t('user.UserID') }}
              <el-tooltip effect="dark" :content="$t('user.filterOptionTip_userID')" placement="top"><i
                class='el-icon-info'/></el-tooltip>
            </span>
            <el-input v-model="searchFormLog.user_ids" clearable></el-input>
          </el-form-item>

          <el-form-item :label="$t('user.Name')">
            <el-input v-model="searchFormLog.user_name" clearable></el-input>
          </el-form-item>
          <el-form-item :label="$t('user.Tel')">
            <input-number v-model="searchFormLog.tel" :range="[0,99999999999]" clearable></input-number>
          </el-form-item>
          <el-form-item :label="$t('user.UserType')">
            <el-select v-model="searchFormLog.role" clearable multiple>
              <el-option v-for="(item, index) in userType" :key="index" :label="item.label"
                         :value="item.value"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('user.UserChannel')">
            <el-select v-model="searchFormLog.channel" clearable>
              <el-option v-for="(item, index) in userChannel" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
          </el-form-item>
        </el-form>

        <el-form inline v-model="searchFormLog">
          <el-form-item :label="$t('inquire.item.log_type')">
            <el-select v-model="searchFormLog.log_type_list" multiple clearable>
              <el-option v-for="item in inquireItemLogType" :label="item.label" :value="item.value"
                         :key="item.value"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <span slot="label">{{ $t('inquire.item.gacha_name') }}
              <el-tooltip effect="dark" :content="$t('inquire.item.filterOptionTip_gachaName')" placement="top"><i
                class='el-icon-info'/></el-tooltip>
            </span>
            <el-input v-model="searchFormLog.gacha_name" clearable></el-input>
          </el-form-item>
        </el-form>

        <el-form inline v-model="searchFormLog">
          <el-form-item>
            <span slot="label">{{ $t('inquire.item.update_amount') }}
              <el-tooltip effect="dark" :content="$t('inquire.item.filterOptionTip_updateAmount')" placement="top"><i
                class='el-icon-info'/></el-tooltip>
            </span>
            <input-number-range v-model="searchFormLog.update_amount_range" clearable
                                :precision="2"></input-number-range>
          </el-form-item>
          <el-form-item :label="$t('inquire.item.show_price')">
            <input-number-range v-model="searchFormLog.show_price_range" clearable :range="[0, 2**31-1]"
                                :precision="2"></input-number-range>
          </el-form-item>
          <el-form-item :label="$t('inquire.item.inner_price')">
            <input-number-range v-model="searchFormLog.inner_price_range" clearable :range="[0, 2**31-1]"
                                :precision="2"></input-number-range>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" v-loading="loading" @click="fetch({page: 1})">{{
                $t('common.Search')
              }}
            </el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExport({})">{{ $t('common.Export') }}</el-button>
            <el-button type="warning" v-loading="loading" @click="fetchExportDetail({})">
              {{ $t('inquire.item.export_detail') }}
            </el-button>
          </el-form-item>
        </el-form>
      </template>
      <template v-slot:table-column>
        <el-table-column prop="date_time" :label="$t('common.DateTime')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="user_name" :label="$t('user.Name')" min-width="90px" align="center"></el-table-column>
        <el-table-column prop="log_type_str" :label="$t('inquire.item.log_type')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="log_type_name" :label="$t('inquire.item.gacha_name')" min-width="90px"
                         align="center"></el-table-column>
        <el-table-column prop="bet_nums" :label="$t('inquire.item.bet_nums')" min-width="90px" align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type<=199 && data.row.level_type === 1">{{ data.row.bet_nums | localeNum }}
            </template>
            <template v-else-if="data.row.log_type===300">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else-if="data.row.log_type===701">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else-if="data.row.log_type===100002">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else-if="data.row.log_type===100004">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else-if="data.row.log_type===100005">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else-if="data.row.log_type===999999">{{ data.row.bet_nums | localeNum }}</template>
            <template v-else>-</template>
          </template>
        </el-table-column>
        <el-table-column prop="level_type_str" :label="$t('inquire.item.level_type')" min-width="90px" align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type<=199">{{ data.row.level_type_str }}</template>
            <template v-else>-</template>
          </template>
        </el-table-column>
        <el-table-column prop="update_amount" :label="$t('inquire.item.update_amount')" min-width="90px" align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type===100002 && data.row.log_type_name === '兑换'">-</template>
            <template v-else-if="data.row.log_type===100003">-</template>
            <template v-else-if="data.row.log_type===100004">-</template>
            <template v-else-if="data.row.log_type===100005">-</template>
            <template v-else-if="data.row.log_type===100006">-</template>
            <template v-else-if="data.row.log_type===100007">-</template>
            <template v-else-if="data.row.log_type===100011">-</template>
            <template v-else-if="data.row.log_type===999999">-</template>
            <template v-else>{{ data.row.update_amount | localeNum2f }}</template>
          </template>
        </el-table-column>
        <el-table-column prop="show_price" :label="$t('inquire.item.show_price')" min-width="90px" align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type!==200 && data.row.log_type!==100004">
              {{ data.row.show_price | localeNum2f }}
            </template>
            <template v-else>-</template>
          </template>
        </el-table-column>
        <el-table-column prop="inner_price" :label="$t('inquire.item.inner_price')" min-width="90px" align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type!==200 && data.row.log_type!==100004">
              {{ data.row.inner_price | localeNum2f }}
            </template>
            <template v-else>-</template>
          </template>
        </el-table-column>
        <el-table-column prop="recycling_price" :label="$t('inquire.item.recycling_price')" min-width="90px"
                         align="center">
          <template v-slot="data">
            <template v-if="data.row.log_type!==200 && data.row.log_type!==100004">
              {{ data.row.recycling_price | localeNum2f }}
            </template>
            <template v-else>-</template>
          </template>
        </el-table-column>
        <el-table-column type="expand">
          <template v-slot="record">
            <template v-if="(record.row.log_type===200)&& !loading">
              <el-row :gutter="20">
                <el-col :span="12">
                  <item-detail :loading="loading" :data="record.row.detail[0]"
                               :span-method="spanMethodDetailMarket(record.row.detail[0])">
                    <template v-slot:table-top-column>
                      <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px"
                                       align="center"></el-table-column>
                      <el-table-column prop="user_name" :label="$t('inquire.item.creator')" min-width="90px"
                                       align="center"></el-table-column>
                      <el-table-column prop="amount" :label="$t('inquire.item.ask_price')" min-width="90px"
                                       align="center">
                        <template v-slot="detail">{{ detail.row.amount | localeNum2f }}</template>
                      </el-table-column>
                    </template>
                    <template v-slot:table-tail-column>
                      <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                        <template v-slot="detail">
                          <template v-if="detail.row.item_id !== '0'">{{ detail.row.nums | localeNum }}</template>
                          <template v-else>-</template>
                        </template>
                      </el-table-column>
                    </template>
                  </item-detail>
                </el-col>
                <el-col :span="12">
                  <item-detail :loading="loading" :data="record.row.detail[1]"
                               :span-method="spanMethodDetailMarket(record.row.detail[1])">
                    <template v-slot:table-top-column>
                      <el-table-column prop="user_id" :label="$t('user.UserID')" min-width="90px"
                                       align="center"></el-table-column>
                      <el-table-column prop="user_name" :label="$t('inquire.item.offerer')" min-width="90px"
                                       align="center"></el-table-column>
                      <el-table-column prop="amount" :label="$t('inquire.item.offer_price')" min-width="90px"
                                       align="center">
                        <template v-slot="detail">{{ detail.row.amount | localeNum2f }}</template>
                      </el-table-column>
                    </template>
                    <template v-slot:table-tail-column>
                      <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                        <template v-slot="detail">
                          <template v-if="detail.row.item_id !== '0'">{{ detail.row.nums | localeNum }}</template>
                          <template v-else>-</template>
                        </template>
                      </el-table-column>
                    </template>
                  </item-detail>
                </el-col>
              </el-row>
            </template>
            <template
              v-else-if="(record.row.log_type===100002 || record.row.log_type===100003 || record.row.log_type===100005 || record.row.log_type===100006 ||record.row.log_type===100007 || record.row.log_type===100010 || record.row.log_type===100011 ||record.row.log_type===200000) && !loading">
              <item-detail :loading="loading" :data="record.row.detail">
                <template v-slot:table-tail-column>
                  <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                    <template v-slot="detail">{{ detail.row.nums * (record.row.bet_nums || 1) | localeNum }}</template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
            <template v-else-if="(record.row.log_type===300) && !loading">
              <item-detail :loading="loading" :data="record.row.detail">
                <template v-slot:table-tail-column>
                  <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                    <template v-slot="detail">{{ detail.row.nums | localeNum }}</template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
            <template v-else-if="(record.row.log_type===701) && !loading">
              <item-detail :loading="loading" :data="record.row.detail">
                <template v-slot:table-tail-column>
                  <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                    <template v-slot="detail">{{ detail.row.nums | localeNum }}</template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
            <template v-else-if="(record.row.log_type===100004) && !loading">
              <el-col :span="12">
                <item-detail :loading="loading" :data="record.row.detail[0]">
                  <template v-slot:table-tail-column>
                    <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                      <template v-slot="detail">
                        <template v-if="detail.row.item_id !== '0'">{{ detail.row.nums | localeNum }}</template>
                        <template v-else>-</template>
                      </template>
                    </el-table-column>
                  </template>
                </item-detail>
              </el-col>
              <el-col :span="12">
                <item-detail :loading="loading" :data="record.row.detail[1]">
                  <template v-slot:table-tail-column>
                    <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                      <template v-slot="detail">
                        <template v-if="detail.row.item_id !== '0'">{{ detail.row.nums | localeNum }}</template>
                        <template v-else>-</template>
                      </template>
                    </el-table-column>
                  </template>
                </item-detail>
              </el-col>
            </template>
            <template v-else-if="(record.row.log_type===999999) && !loading">
              <item-detail :loading="loading" :data="record.row.detail" show-level-name>
                <template v-slot:table-tail-column>
                  <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                    <template v-slot="detail">{{ detail.row.nums | localeNum }}</template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
            <template v-else-if="!loading">
              <item-detail :loading="loading" :data="record.row.detail" show-level-name
                           :span-method="spanMethodDetailBet(record.row.detail)">
                <template v-slot:table-top-column>
                  <!--                  <el-table-column prop="gacha_name" :label="$t('inquire.item.gacha_name')" min-width="90px" align="center"></el-table-column>-->
                  <el-table-column prop="box_out_no" :label="$t('inquire.item.box_out_no')" min-width="90px"
                                   align="center"></el-table-column>
                  <template v-if="record.row.log_type > 100 && record.row.log_type < 199">
                    <el-table-column prop="no" :label="$t('inquire.item.no')" min-width="40px" align="center">
                      <template v-slot="detail">
                        <template v-if="detail.row.no !== 0">{{ detail.row.no | localeNum }}</template>
                        <template v-else>-</template>
                      </template>
                    </el-table-column>
                  </template>
                </template>
                <template v-slot:table-tail-column>
                  <el-table-column prop="nums" :label="$t('inquire.item.nums')" min-width="90px" align="center">
                    <template v-slot="detail">{{ detail.row.nums | localeNum }}</template>
                  </el-table-column>
                </template>
              </item-detail>
            </template>
          </template>
        </el-table-column>
      </template>
    </search-table>
  </div>
</template>

<script>
import DatePicker from "@/components/DatePicker.vue";
import moment from "moment/moment";
import api from "@/api";
import {mapGetters,} from "vuex";
import InputNumber from "@/components/Input/InputNumber.vue";
import InputNumberRange from "@/components/Input/InputNumberRange.vue"
import ItemDetail from "@/components/Page/ItemDetail.vue";
import SearchTable from "@/components/Table/SearchTable.vue";
import {isInvalidNumber} from '@/utils/validate';

export default {
  name: "Item",
  components: {SearchTable, ItemDetail, InputNumber, InputNumberRange, DatePicker},
  data() {
    return {
      loading: false,
      searchFormLog: {
        page: 1,
        page_size: 50,
        date_time_range: [],
        user_ids: '',
        user_name: '',
        tel: '',
        role: [],
        channel: null,
        log_type_list: ['101', '102', '103', '104', '105','106'],
        gacha_name: '',
        update_amount_range: [], // 没填则不传, 用边界补齐 两个元素
        show_price_range: [],
        inner_price_range: [],
      },
      searchFormDetail: {
        id: "",
        log_type: 0,
        level_type: 0,
      },
      tableData: [],
      total: 0,
      summary: {},
    }
  },
  computed: {
    ...mapGetters({
      userType: 'option/userType',
      userChannel: 'option/userChannel',
      inquireItemLogType: 'option/inquireItemLogType',
    }),
  },
  created() {
    if (this.inquireItemLogType.length === 0) {
      this.$store.dispatch('option/update', 'inquireItemLogType')
    }
    this.searchFormLog.date_time_range = [
      moment().subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss'),
      moment().format('YYYY-MM-DD HH:mm:ss'),
    ]
    this.fetch()
  },
  methods: {
    fetch(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      searchFormLog.show_price_range = this.setRangeBorder(searchFormLog.show_price_range)
      searchFormLog.inner_price_range = this.setRangeBorder(searchFormLog.inner_price_range)
      api.getInquireItemLog(searchFormLog).then(res => {
        this.tableData = res.data || []
        this.total = res.headers.total || 0
        this.summary = res.headers.summary || {}
      }).catch(() => {
        this.tableData = []
        this.total = 0
        this.summary = {}
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExport(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      searchFormLog.show_price_range = this.setRangeBorder(searchFormLog.show_price_range)
      searchFormLog.inner_price_range = this.setRangeBorder(searchFormLog.inner_price_range)
      api.exportInquireItemLog(searchFormLog).finally(() => {
        this.loading = false
      })
    },
    fetchDetail(params = {}, row) {
      this.loading = true
      this.searchFormDetail = Object.assign({}, this.searchFormDetail, params)
      api.getInquireItemDetail(this.searchFormDetail).then(res => {
        row.detail = res.data || []
      }).catch(() => {
        row.detail = []
      }).finally(() => {
        this.loading = false
      })
    },
    fetchExportDetail(params = {}) {
      this.loading = true
      this.searchFormLog = Object.assign(this.searchFormLog, params)
      let searchFormLog = Object.assign({}, this.searchFormLog)
      searchFormLog.update_amount_range = this.setRangeBorder(searchFormLog.update_amount_range)
      searchFormLog.show_price_range = this.setRangeBorder(searchFormLog.show_price_range)
      searchFormLog.inner_price_range = this.setRangeBorder(searchFormLog.inner_price_range)
      api.exportInquireItemDetailLog(searchFormLog).finally(() => {
        this.loading = false
      })
    },
    setRangeBorder(amountRange) {
      if (amountRange.length === 2) {
        if (isInvalidNumber(amountRange[0]) && isInvalidNumber(amountRange[1])) {
          return []
        } else if (isInvalidNumber(amountRange[0]) && !isInvalidNumber(amountRange[1])) {
          return [-(2 ** 30), amountRange[1]]
        } else if (!isInvalidNumber(amountRange[0]) && isInvalidNumber(amountRange[1])) {
          return [amountRange[0], (2 ** 30)]
        }
      }
      return amountRange
    },
    expandChange(row) {
      if (row.detail === undefined) {
        this.fetchDetail({id: row.id, log_type: row.log_type, level_type: row.level_type,}, row)
      }
    },
    spanMethodDetailBet(detail) {
      return ({row, column, rowIndex, columnIndex}) => {
        if ([0].includes(columnIndex)) {
          if (rowIndex === 0) {
            return {rowspan: detail.length, colspan: 1}
          } else {
            return {rowspan: 0, colspan: 0}
          }
        }
      }
    },
    spanMethodDetailMarket(detail) {
      return ({row, column, rowIndex, columnIndex}) => {
        if ([0, 1, 2].includes(columnIndex)) {
          if (rowIndex === 0) {
            return {rowspan: detail.length, colspan: 1}
          } else {
            return {rowspan: 0, colspan: 0}
          }
        }
      }
    },
    summaryMethod({columns, data}) {
      return [
        'Total', Number(this.summary["user_cnt"]).toLocaleString(), undefined, undefined, undefined,
        Number(this.summary["bet_nums"]).toLocaleString(), undefined,
        Number(this.summary["update_amount"]).toLocaleString(), Number(this.summary["show_price"]).toLocaleString(), Number(this.summary["inner_price"]).toLocaleString(), Number(this.summary["recycling_price"]).toLocaleString(),
        undefined,
      ]
    },
    changeIsAmin(value) {
      if (value === "") {
        delete (this.searchFormLog.is_admin)
      } else {
        this.searchFormLog.is_admin = value
      }
    },
  },
}
</script>

<style scoped>

</style>
