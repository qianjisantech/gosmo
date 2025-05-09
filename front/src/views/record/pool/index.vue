<template>
  <div>
    <page-header/>

    <div class="table-search-wrapper">
      <a-card :bordered="false">
        <a-form :model="queryForm" @finish="onFinish">
          <a-row>
            <a-col flex="0 1 400px">
              <a-form-item name="planNo" label="索引" style="max-width: 300px;">
                <a-input v-model:value="queryForm.planNo" placeholder="请输入名称" allowClear></a-input>
              </a-form-item>
            </a-col>
            <a-col flex="0 1 400px">
              <a-form-item name="host" label="域名" style="max-width: 300px;">
                <a-input v-model:value="queryForm.host" placeholder="请输入域名" allowClear></a-input>
              </a-form-item>
            </a-col>
            <a-col flex="0 1 400px">
              <a-form-item name="url" label="url" style="max-width: 300px;">
                <a-input v-model:value="queryForm.url" placeholder="请输入url" allowClear></a-input>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row>
            <a-col>
              <a-button type="primary" html-type="submit" :loading="tableLoading">查询</a-button>
              <a-button style="margin: 0 8px" @click="resetFields">重置</a-button>
            </a-col>
          </a-row>
        </a-form>
      </a-card>
    </div>

    <div class="table-wrapper">
      <a-card :bordered="false" :headStyle="{ borderBottom: 'none', padding: '20px 24px' }"
              :bodyStyle="{ padding: '0 24px', minHeight: 'calc(100vh - 400px)' }">
        <template #extra>
          <el-space></el-space>
        </template>
        <a-table
            :rowKey="record => record.id"
            :columns="columns"
            :data-source="dataSource"
            :row-selection="rowSelection"
            :loading="tableLoading"
            @change="handleTableChange"
            :scroll="{ x: 500, y: 'calc(100vh - 500px)' }"
            :pagination="pagination"
            :style="{ minHeight: '500px' }"
        >
          <template #expandedRowRender="{ record }">
            <a-descriptions title="请求信息" bordered>
              <a-descriptions-item label="请求体" :span="10">
                <pre>{{ record.request.body || '无请求体' }}</pre>
              </a-descriptions-item>
              <a-descriptions-item label="请求头" :span="2">
                <a-table
                    :columns="headerColumns"
                    :data-source="formatHeaders(record.request.headers)"
                    :pagination="false"
                    bordered
                    size="small"
                >
                  <template #bodyCell="{ column, record: headerRecord }">
                    <template v-if="column.dataIndex === 'value'">
                      <span v-if="!isSensitiveKey(headerRecord.key)" class="header-value">
                        {{ headerRecord.value }}
                      </span>
                    </template>
                  </template>
                </a-table>
              </a-descriptions-item>
            </a-descriptions>

            <a-descriptions title="响应信息" bordered style="margin-top: 16px">
              <a-descriptions-item label="响应体" :span="10">
                <pre>{{ record.response.body || '无响应体' }}</pre>
              </a-descriptions-item>
              <a-descriptions-item label="响应头" :span="2">
                <a-table
                    :columns="headerColumns"
                    :data-source="formatHeaders(record.response.headers)"
                    :pagination="false"
                    bordered
                    size="small"
                >
                  <template #bodyCell="{ column, record: headerRecord }">
                    <template v-if="column.dataIndex === 'value'">
                      <span v-if="!isSensitiveKey(headerRecord.key)" class="header-value">
                        {{ headerRecord.value }}
                      </span>
                    </template>
                  </template>
                </a-table>
              </a-descriptions-item>
            </a-descriptions>
          </template>

          <template #bodyCell="{ column, record, index }">
            <template v-if="column.dataIndex === 'operation'">
              <div style="display: flex; gap: 15px;">
                <a v-on:click="modalHandle('view', index)">查看</a>
                <a-popconfirm title="确定删除吗？" ok-text="确定" cancel-text="取消" @confirm="deleteRow(record)">
                  <a>删除</a>
                </a-popconfirm>
              </div>
            </template>
          </template>
        </a-table>
      </a-card>
    </div>

    <div>
      <a-modal v-model:open="openModal" @ok="handleModalSumbit" :width="800" :destroyOnClose="true"
               :confirmLoading="modalSubmitLoading" style="top: 30px">
        <template #title>
          <span>{{ modalTitle === 'create' ? '新建岗位' : (modalTitle === 'view' ? '查看流量' : '') }}</span>
        </template>
      </a-modal>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, unref, onMounted } from 'vue';
import PageHeader from '@/components/PageHeader.vue';
import { Table, message, Modal } from 'ant-design-vue';
import { TrafficPoolPage } from '@/api/traffic/pool'
import type { TableColumnsType } from 'ant-design-vue';
import { QueryForm } from './types'
import { paginationRequest } from '@/common/pagination'
import { TableDataType } from "@/common/types";

const tableLoading = ref(false);
const openModal = ref(false);
const modalTitle = ref('');
const modalSubmitLoading = ref(false);
const dataSource = ref<TableDataType[]>([]);
const selectedRowKeys = ref<TableDataType['id'][]>([]);

const pagination = reactive({
  current: 1,
  pageSize: 10,
  defaultPageSize: 10,
  showSizeChanger: true,
  total: 0,
  showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条 / 总共 ${total} 条`
})

const headerColumns = [
  {
    title: 'Key',
    dataIndex: 'key',
    width: '30%',
    align: 'center',
    ellipsis: true
  },
  {
    title: 'Value',
    dataIndex: 'value',
    align: 'left',
    ellipsis: true
  }
]

const columns: TableColumnsType = [
  {
    title: '索引',
    dataIndex: 'index',
    align: 'center',
  },
  {
    title: 'id',
    dataIndex: 'id',
    align: 'center'
  },
  {
    title: '录制时间',
    dataIndex: 'timestamp',
    align: 'center'
  },
  {
    title: '域名',
    dataIndex: 'host',
    align: 'center',
  },
  {
    title: '请求路径',
    dataIndex: 'url',
    align: 'center',
  },
  {
    title: '请求方式',
    dataIndex: 'method',
    align: 'center',
  },
  {
    title: '响应状态',
    dataIndex: 'status',
    align: 'center',
  },
  {
    title: '操作',
    dataIndex: 'operation',
    align: 'center',
    fixed: 'right',
    width: 150
  }
];

const queryForm = reactive<QueryForm>({
  ...paginationRequest,
  planNo: '',
  url: "",
  key: ""
})

// 标记敏感key（如Authorization）
const isSensitiveKey = (key) => {
  const sensitiveKeys = ['authorization', 'cookie', 'token']
  return sensitiveKeys.includes(key.toLowerCase())
}

// 格式化headers为表格数据
const formatHeaders = (headers) => {
  if (!headers) return []
  return Object.entries(headers).map(([key, value]) => ({
    key,
    value: typeof value === 'string' ? value : JSON.stringify(value)
  }))
}

const loadingData = () => {
  tableLoading.value = true;
  TrafficPoolPage(queryForm).then((res) => {
    const result = res.data;
    dataSource.value = result.data
    pagination.total = result.total || 0;
    tableLoading.value = false;
  }).catch(error => {
    console.log(error);
    tableLoading.value = false;
  })
}

onMounted(() => loadingData());

const onFinish = () => {
  pagination.current = 1;
  loadingData();
};

const resetFields = () => {
  Object.keys(queryForm).forEach((key: string) => {
    delete queryForm[key];
  });
  loadingData();
}

const handleTableChange = (values: any) => {
  pagination.current = values.current;
  pagination.pageSize = values.pageSize;
  loadingData();
}

const onSelectChange = (selectingRowKeys: TableDataType['id'][]) => {
  selectedRowKeys.value = selectingRowKeys;
}

const rowSelection = computed(() => {
  return {
    selectedRowKeys: unref(selectedRowKeys),
    onChange: onSelectChange,
    hideDefaultSelections: true,
    selections: [
      Table.SELECTION_ALL,
      Table.SELECTION_INVERT,
      Table.SELECTION_NONE
    ]
  }
});

const modalHandle = (modalType: string, index?: number) => {
  modalTitle.value = modalType;
  openModal.value = true;
}

const handleModalSumbit = () => {
  modalSubmitLoading.value = true;
  modalSubmitLoading.value = false;
  openModal.value = false;
}

const deleteRow = (row: TableDataType) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这条记录吗？',
    onOk() {
      // 这里调用删除API
      message.success('删除成功');
      loadingData();
    }
  });
}
</script>

<style lang="scss" scoped>
.table-search-wrapper {
  margin-block-end: 16px;
}

.header-value {
  word-break: break-word;
  font-family: monospace;
}

.sensitive-value {
  color: #ff4d4f;
  font-style: italic;
}

pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  margin: 0;
  padding: 8px;
  background: #f5f5f5;
  border-radius: 2px;
  max-height: 300px;
  overflow: auto;
}

:deep(.ant-descriptions-item-label) {
  font-weight: bold;
  width: 100px;
}
</style>