<template>
  <div>
    <page-header/>
    <el-card :body-style="{ padding: '50px 50px',minHeight: 'calc(100vh - 250px)' }">
      <!--    主数据-->
      <el-form
          ref="ruleFormRef"
          style="max-width: 500px"
          :model="createPlanForm"
          :rules="rules"
          label-width="auto"
          class="demo-ruleForm"
          :size="formSize"
          status-icon
      >
        <el-form-item label="执行模式" prop="executionMode">
          <el-radio-group v-model="createPlanForm.executionMode" size="default">
            <el-radio-button label="ES模式" value="es"/>
            <el-radio-button label="Agent模式" value="agent"/>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="项目" prop="project">
          <el-select v-model="createPlanForm.project" placeholder="请选择项目" @focus="fetchProjectOptions">
            <el-option
                v-for="item in projectOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="执行机" prop="machine">
          <el-select v-model="createPlanForm.machine" placeholder="请选择执行机" @focus="fetchEsOptions">
            <el-option
                v-for="item in esOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="录制持续时间" prop="duration">
          <el-input-number v-model="createPlanForm.duration" :min="1" :max="10">
            <template #suffix>
              <span>min</span>
            </template>
          </el-input-number>
        </el-form-item>
      </el-form>
      <el-divider/>
      <el-tabs v-model="tabsValue" class="demo-tabs">
        <!--      匹配规则tab-->
        <el-tab-pane label="匹配规则" name="match">
          <div class="mb-4">
            <el-button @click="addMatchRule" size="small" type="primary">新增</el-button>
          </div>
          <el-table :data="matchTableData" style="width: 100%">
            <el-table-column type="selection" width="55"/>
            <el-table-column label="匹配模式" width="180">
              <template #default="scope">
                <el-input size="small" v-model="scope.row.matchMode"></el-input>
              </template>
            </el-table-column>

            <el-table-column label="匹配值" width="180">
              <template #default="scope">
                <el-input size="small" v-model="scope.row.matchValue"></el-input>
              </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" min-width="120">
              <template #default="scope">
                <el-button
                    link
                    type="primary"
                    size="small"
                    @click.prevent="deleteMatchRule(scope.$index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

        </el-tab-pane>
        <!--      <el-tab-pane label="Config" name="second">Config</el-tab-pane>-->
        <!--      <el-tab-pane label="Role" name="third">Role</el-tab-pane>-->
        <!--      <el-tab-pane label="Task" name="fourth">Task</el-tab-pane>-->
      </el-tabs>


      <template #footer>
        <div style="text-align: center">
          <el-button type="primary" @click="submitForm(ruleFormRef)">
            确认
          </el-button>
          <el-button @click="resetForm(ruleFormRef)">返回</el-button>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue'
import type {ComponentSize, FormInstance, FormRules} from 'element-plus'
import PageHeader from "@/components/PageHeader.vue";
import {CreatePlanForm, MatchTableColumn, Option} from './types'
import {getCommonOptionProject, getCommonOptionEs} from '@/api/basicdata/common'
import {useRouter} from "vue-router";

const formSize = ref<ComponentSize>('default')
const ruleFormRef = ref<FormInstance>()
const router = useRouter()
const createPlanForm = reactive<CreatePlanForm>({
  executionMode: 'es',
  project: '',
  machine: '',
  duration: 1,
})
const matchTableData = reactive<MatchTableColumn[]>([
  { id:"11111",
    matchMode: "全文",
    matchValue: "demo"
  },
])
//定义项目选择list
const projectOptions = reactive<Option[]>([])
const esOptions = reactive<Option[]>([])
const tabsValue = ref('match')
const rules = reactive<FormRules<CreatePlanForm>>({
  executionMode: [
    {required: true, message: '执行模式必选', trigger: 'blur'},
  ],
  project: [
    {required: true, message: '项目必选', trigger: 'blur',},
  ],
  machine: [{required: true, message: '执行机必选', trigger: 'blur',},
  ],
  duration: [{required: true, message: '持续时间必填', trigger: 'change',},
  ],

})
const  addMatchRule =()=>{
  // 新增空数据的函数
    matchTableData.push({
      matchMode: "",
      matchValue: ""
    });
}
const deleteMatchRule = (index: number) => {
  matchTableData.splice(index, 1); // 删除指定索引的元素
};
//获取项目options
const fetchProjectOptions = async () => {
  await getCommonOptionProject().then((res: any) => {
    console.log(res.data.data)
    projectOptions.splice(0, projectOptions.length, ...res.data.data)

  })
}
//获取esoptions
const fetchEsOptions = async () => {
  await getCommonOptionEs().then((res: any) => {
    console.log(res.data.data)
    esOptions.splice(0, esOptions.length, ...res.data.data)

  })
}
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log('submit!')
      router.push("/record/plan")
    } else {
      console.log('error submit!', fields)
    }
  })
}

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.resetFields()
  router.push("/record/plan")
}

const options = Array.from({length: 10000}).map((_, idx) => ({
  value: `${idx + 1}`,
  label: `${idx + 1}`,
}))
</script>
