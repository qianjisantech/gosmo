<template>
  <div>
    <page-header/>

    <div class="table-search-wrapper">
      <a-card :bordered="false">
        <a-form :model="queryForm" @finish="onFinish">
          <a-row>
            <a-col flex="0 1 400px">
              <a-form-item name="name" label="计划编号" style="max-width: 300px;">
                <a-input v-model:value="queryForm.planNo" placeholder="请输入计划编号" allowClear></a-input>
              </a-form-item>
            </a-col>
            <a-col flex="0 1 400px">
              <a-form-item name="executionModeCode" label="执行模式" style="max-width: 300px;">
                <a-select v-model:value="queryForm.executionModeCode" placeholder="全部" allowClear>
                  <a-select-option value="1">ES模式</a-select-option>
                  <a-select-option value="2">Agent模式</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col flex="0 1 400px">
              <a-form-item name="executionStatus" label="执行状态" style="max-width: 300px;">
                <a-select v-model:value="queryForm.executionStatus" placeholder="全部" allowClear>
                  <a-select-option value="1">未开始</a-select-option>
                  <a-select-option value="2">运行中</a-select-option>
                  <a-select-option value="3">已完成</a-select-option>
                  <a-select-option value="4">运行异常</a-select-option>
                  <a-select-option value="5">已暂停</a-select-option>
                </a-select>
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
          <a-button type="primary" :icon="h(PlusOutlined)" @click="modalHandle('create')" style="margin-right: 10px;">
            新建
          </a-button>
<!--          <a-dropdown>-->
<!--            <template #overlay>-->
<!--              <a-menu @click="handleMoreClick">-->
<!--                <a-menu-item key="1"><span style="margin-right: 10px;"><CheckOutlined/></span><span>批量启用</span>-->
<!--                </a-menu-item>-->
<!--                <a-menu-item key="2"><span style="margin-right: 10px;"><StopOutlined/></span><span>批量停用</span>-->
<!--                </a-menu-item>-->
<!--              </a-menu>-->
<!--            </template>-->
<!--            <a-button>更多-->
<!--              <DownOutlined/>-->
<!--            </a-button>-->
<!--          </a-dropdown>-->
        </template>
        <a-table :rowKey="record => record.id" :columns="columns" :data-source="dataSource"
                 :row-selection="rowSelection" :loading="tableLoading"
                 @change="handleTableChange" :scroll="{ x: 500, y: 'calc(100vh - 500px)' }" :pagination="pagination"
                 :style="{ minHeight: '500px' }">
          <template v-slot:bodyCell="{ column, record, index }">
            <template v-if="column.dataIndex === 'index'">
              <span>{{ index + 1 }}</span>
            </template>
            <template v-if="column.dataIndex === 'executionModeName'">
              <span>
                <el-tag effect="light"
                        :type="record.executionModeCode==='1' ? 'primary' : 'success'">{{
                    record.executionModeName
                  }}</el-tag>
              </span>
            </template>
            <template v-if="column.dataIndex === 'executionStatus'">
              <span>
       <a-badge :color="getModeStatusBadgeColor(record.executionStatus)"/>{{
                  getModeStatusStatusText(record.executionStatus)
                }}

              </span>
            </template>
            <template v-if="column.dataIndex === 'operation'">
              <div style="display: flex; gap: 15px;">
                <a v-on:click="modalHandle('view', index)">查看</a>
                <a v-on:click="modalHandle('update', index)">编辑</a>
                <a-popconfirm title="确定删除吗？" ok-text="确定" cancel-text="取消" @confirm="deleteRow(record)">
                  <a v-show="record.executionStatus==='1'">删除</a>
                </a-popconfirm>
                <a v-show="record.executionStatus!=='1'" v-on:click="modalHandle('pool', index)">流量池</a>
              </div>
            </template>
          </template>
        </a-table>
      </a-card>
    </div>

  </div>
</template>

<script lang="ts" setup>
import {ref, reactive, computed, unref, onMounted, h} from 'vue';
import PageHeader from '@/components/PageHeader.vue';
import {Table, message, Modal} from 'ant-design-vue';
import {
  getPositionPage,
  createPosition,
  updatePosition,
  deletePosition,
  batchEnablePosition,
  batchDisablePosition
} from '@/api/sys/position'
import {getPlanPage} from '@/api/record/plan'

import {cloneDeep, isEmpty} from '@/utils/util';
import {PlusOutlined, DownOutlined, CheckOutlined, StopOutlined} from '@ant-design/icons-vue';
import type {TableColumnsType, MenuProps} from 'ant-design-vue';
import type {SearchDataType, TableDataType} from '@/common/types'
import {useRouter} from "vue-router";
import  {paginationRequest} from '@/common/pagination'

const router = useRouter()
const createForm = ref();
const updateForm = ref();
const tableLoading = ref(false);
const openModal = ref(false);
const modalTitle = ref('');
const modalSubmitLoading = ref(false);
const detailStateLoading = ref(false);
const dataSource = ref<TableDataType[]>([]);
const selectedRowKeys = ref<TableDataType['id'][]>([]);

const pagination = reactive({
  current: 1,
  pageSize: 10,
  defaultPageSize: 10,
  showSizeChanger: true,
  total: dataSource.value.length,
  showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条 / 总共 ${total} 条`
})
const createState = reactive({
  name: '',
  order: 1,
  description: ''
})
const updateState = reactive({
  id: undefined,
  name: '',
  order: 1,
  available: true,
  description: ''
})

const columns: TableColumnsType = [
  {
    title: '序号',
    dataIndex: 'index',
    align: 'center',
    width: 80
  },
  {
    title: '计划编号',
    dataIndex: 'planNo',
    align: 'center'
  },
  {
    title: '项目',
    dataIndex: 'project',
    align: 'center'
  },
  {
    title: '执行模式',
    dataIndex: 'executionModeName',
    align: 'center'
  },
  {
    title: '执行状态',
    dataIndex: 'executionStatus',
    align: 'center'
  },
  {
    title: '备注',
    dataIndex: 'remark',
    align: 'center',
    ellipsis: true,
    width: 500
  },
  {
    title: '操作',
    dataIndex: 'operation',
    align: 'center',
    fixed: 'right',
    width: 200
  }
];
const queryForm= reactive({
  ...paginationRequest,
  planNo:'',
  executionModeCode:'',
  executionStatus:undefined
})
const getModeStatusBadgeColor = (status: string) => {
  switch (status) {
    case '1':
      return 'gray';   // 未开始，灰色
    case '2':
      return 'blue';   // 运行中，蓝色
    case '3':
      return 'green';  // 已完成，绿色
    case  '4':
      return 'red'
    case  '5':
      return 'orange'
    default:
      return 'red';    // 未知状态，红色

  }
}
const getModeStatusStatusText = (status: string) => {
  switch (status) {
    case '1':
      return '未开始';
    case '2':
      return '运行中';
    case '3':
      return '已完成';
    case '4':
      return '运行异常';
    case '5':
      return '已暂停';
    default:
      return '未知状态';
  }
}
const loadingData = () => {
  tableLoading.value = true;

  getPlanPage(queryForm).then(response => {
    const result = response.data;
    dataSource.value = result.data;
    pagination.total = result.total;
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
  pagination.current = 1;
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

  if (modalType === 'create') {
    router.push("/record/plan/create")
  }
  if (modalType === 'update' && index !== undefined) {
    router.push("/record/plan/update")
  }
  if (modalType === 'pool' && index !== undefined) {
    router.push(
        {
          path: '/record/plan/pool'
          , query: {planNo: index}
        })
  }
}
const handleModalSumbit = () => {
  modalSubmitLoading.value = true;

  if (modalTitle.value === 'view') {
    modalSubmitLoading.value = false;
    openModal.value = false;

  } else if (modalTitle.value === 'create') {
    createForm.value.validate().then(() => {
      const createBody = cloneDeep(createState);
      Object.keys(createBody).forEach(key => {
        if (isEmpty(createBody[key])) {
          delete createBody[key];
        }
      })

      createPosition(createBody).then(response => {
        modalSubmitLoading.value = false;
        openModal.value = false;
        Object.keys(createState).forEach(key => delete createState[key]);
        createState.order = 1;
        message.success(response.data.message);
        loadingData();

      }).catch(error => {
        modalSubmitLoading.value = false;
        console.log(error)
      })

    }).catch(error => {
      modalSubmitLoading.value = false;
      console.log(error)
    })

  } else {
    updateForm.value.validate().then(() => {
      updatePosition(updateState).then(response => {
        modalSubmitLoading.value = false;
        openModal.value = false;
        message.success(response.data.message);
        loadingData();
      })

    }).catch(error => {
      modalSubmitLoading.value = false;
      console.log(error)
    })
  }
}

const deleteRow = (row: TableDataType) => {
  deletePosition({id: row.id}).then(response => {
    message.success(response.data.message);
    loadingData();
  }).catch(error => {
    console.log(error)
  })
}

const handleMoreClick: MenuProps['onClick'] = e => {
  if (!selectedRowKeys.value || !(selectedRowKeys.value.length > 0)) {
    message.warning('请先勾选数据');
    return;
  }

  Modal.confirm({
    title: '提示',
    content: e.key == 1 ? '是否确定启用选择项？' : '是否确定停用选择项？',
    onOk() {
      const body = {ids: selectedRowKeys.value};
      const batchApi = e.key == 1 ? batchEnablePosition(body) : batchDisablePosition(body);
      batchApi.then(response => {
        message.success(response.data.message);
        selectedRowKeys.value = [];
        loadingData();
      }).catch(error => {
        console.log(error);
      })
    }
  });
}

</script>

<style lang="scss" scoped>
.table-search-wrapper {
  margin-block-end: 16px;
}
</style>