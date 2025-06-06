<template>
<el-form :model="queryParams" ref="queryForm" :inline="true" label-width="68px">
     <el-form-item label="dictionary name" prop="dictType">
       <el-select v-model="queryParams.dictType">
         <el-option v-for="item in typeOptions" :key="item.dictId" :label="item.dictName" :value="item.dictType" />
       </el-select>
     </el-form-item>
     <el-form-item label="status" prop="status">
       <el-select v-model="queryParams.status" placeholder="data status" clearable>
         <el-option v-for="dict in statusOptions" :key="dict.dictValue" :label="dict.dictLabel" :value="dict.dictValue" />
       </el-select>
     </el-form-item>
     <el-form-item>
       <el-button type="primary" icon="search" @click="handleQuery">Search</el-button>
       <el-button icon="refresh" @click="resetQuery">Reset</el-button>
     </el-form-item>
   </el-form>

  <el-row :gutter="10" class="mb8">
    <el-col :span="1.5">
      <el-button type="primary" plain icon="plus" @click="handleAdd" v-hasPermi="['system:dict:add']">AddNew</el-button>
    </el-col>
  </el-row>
  <el-table :data="dataList">
    <el-table-column type="selection" width="55" align="center" />
     <el-table-column label="dictionary code" align="center" prop="dictCode" />
     <el-table-column label="dictionary label" align="center" prop="dictLabel">
      <template #default="scope">
        <span v-if="scope.row.listClass == '' || scope.row.listClass == 'default'" :class="scope.row.cssClass">{{ scope.row.dictLabel }}</span>
        <el-tag v-else :type="scope.row.listClass == 'primary' ? '' : scope.row.listClass" :class="scope.row.cssClass"
          >{{ scope.row.dictLabel }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="dictionary key value" align="center" prop="dictValue" sortable />
     <el-table-column label="dictionary sort" align="center" prop="dictSort" sortable />
     <el-table-column label="status" align="center" prop="status">
      <template #default="scope">
        <dict-tag :options="statusOptions" :value="scope.row.status" />
      </template>
    </el-table-column>
    <el-table-column label="remarks" align="center" prop="remark" :show-overflow-tooltip="true" />
     <el-table-column label="operation" align="center" class-name="small-padding fixed-width" width="130px">
      <template #default="scope">
        <el-button text size="small" icon="edit" @click="handleUpdate(scope.row)" v-hasPermi="['system:dict:edit']">edit</el-button>
        <el-button text size="small" icon="delete" @click="handleDelete(scope.row)" v-hasPermi="['system:dict:remove']">delete </el-button>
      </template>
    </el-table-column>
  </el-table>
  <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

  <!-- 添加或修改参数配置对话框 -->
  <el-dialog :title="title" v-model="open" draggable width="500px" append-to-body>
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="字典类型">
        <el-input v-model="form.dictType" :disabled="true" />
      </el-form-item>
      <el-form-item label="数据标签" prop="dictLabel">
        <el-input v-model="form.dictLabel" placeholder="请输入数据标签" />
      </el-form-item>
      <el-form-item label="数据键值" prop="dictValue">
        <el-input v-model="form.dictValue" placeholder="请输入数据键值" />
      </el-form-item>
      <el-form-item label="样式属性" prop="cssClass">
        <!-- <el-input v-model="form.cssClass" placeholder="请输入样式属性" /> -->
        <el-select v-model="form.cssClass" clearable="">
          <el-option
            v-for="dict in cssClassOptions"
            :class="dict.value"
            :key="dict.value"
            :label="dict.label + '(' + dict.value + ')'"
            :value="dict.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="显示排序" prop="dictSort">
        <el-input-number v-model="form.dictSort" controls-position="right" :min="0" />
      </el-form-item>
      <el-form-item label="回显样式" prop="listClass">
        <el-select v-model="form.listClass">
          <el-option
            v-for="item in listClassOptions"
            :key="item.value"
            :label="item.label + '(' + item.value + ')'"
            :value="item.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-radio-group v-model="form.status">
          <el-radio v-for="dict in statusOptions" :key="dict.dictValue" :label="dict.dictValue">{{ dict.dictLabel }}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="form.remark" type="textarea" placeholder="请输入内容"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button text @click="cancel">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup name="dictData">
import { listData, getData, delData, addData, updateData } from '@/api/system/dict/data'
import { listType, getType } from '@/api/system/dict/type'
const { proxy } = getCurrentInstance()
const props = defineProps({
  dictId: {
    type: Number,
    default: 0
  }
})
watch(
  () => props.dictId,
  (newVal, oldValue) => {
    if (newVal) {
      getTypeInfo(newVal)
      getTypeList()
      proxy.getDicts('sys_normal_disable').then((response) => {
        statusOptions.value = response.data
      })
    }
  },
  {
    immediate: true,
    deep: true
  }
)

const ids = ref()
const loading = ref(false)
// total number
const total = ref(0)
// dictionary form data
const dataList = ref([])
// Default dictionary type
const defaultDictType = ref('')
// popup layer title
const title = ref('')
// Whether to display the popup layer
const open = ref(false)
// data label echo style
const listClassOptions = ref([
   {
     value: 'default',
     label: 'default'
   },
   {
     value: 'primary',
     label: 'main'
   },
   {
     value: 'success',
     label: 'success'
   },
   {
     value: 'info',
     label: 'information'
   },
   {
     value: 'warning',
     label: 'Warning'
   },
   {
     value: 'danger',
     label: 'Dangerous'
   }
])

const cssClassOptions = ref([
   {
     value: 'text-primary',
     label: 'main'
   },
   {
     value: 'text-success',
     label: 'success'
   },
   {
     value: 'text-info',
     label: 'information'
   },
   {
     value: 'text-warning',
     label: 'Warning'
   },
   {
     value: 'text-danger',
     label: 'Dangerous'
   },
   {
     value: 'text-orange',
     label: 'orange red'
   },
   {
     value: 'text-hotpink',
     label: 'pink'
   },
   {
     value: 'text-green',
     label: 'green'
   },
   {
     value: 'text-greenyellow',
     label: 'yellow green'
   },
   {
     value: 'text-purple',
     label: 'purple'
   }
])
// state data dictionary
const statusOptions = ref([])
// type data dictionary
const typeOptions = ref([])
// query parameters
const queryParams = reactive({
   pageNum: 1,
   pageSize: 10,
   dictName: undefined,
   dictType: undefined,
   status: undefined
})
// form parameters

const formRef = ref()
const state = reactive({
  form: {},
  rules: {
    dictLabel: [{ required: true, message: 'Data label cannot be empty', trigger: 'blur' }],
     dictValue: [{ required: true, message: 'Data key value cannot be empty', trigger: 'blur' }],
     dictSort: [{ required: true, message: 'Data order cannot be empty', trigger: 'blur' }]
  }
})

const { form, rules } = toRefs(state)
/** Query dictionary type details */
function getTypeInfo(dictId) {
  getType(dictId).then((response) => {
    queryParams.dictType = response.data.dictType
    defaultDictType.value = response.data.dictType
    getList()
  })
}
/** Query list of dictionary types */
function getTypeList() {
  listType().then((response) => {
    typeOptions.value = response.data.result
  })
}

/** Query dictionary data list */
function getList() {
  loading.value = true
  listData(queryParams).then((response) => {
    dataList.value = response.data.result
    total.value = response.data.totalNum
    loading.value = false
  })
}

// cancel button
function cancel() {
  open.value = false
  reset()
}
// form reset
function reset() {
  form.value = {
    dictCode: undefined,
    dictLabel: undefined,
    dictValue: undefined,
    dictSort: 0,
    status: '0',
    remark: undefined
  }
  proxy.resetForm('formRef')
}

/** Search button action */
function handleQuery() {
   queryParams.pageNum = 1
   getList()
}

/** Reset button action */
function resetQuery() {
   proxy. resetForm('queryForm')
   queryParams.dictType = defaultDictType.value
   handleQuery()
}
/** Add button operation */
function handleAdd() {
   reset()
   open.value = true
   title.value = 'Add dictionary data'
   form.value.dictType = queryParams.dictType
}
// multi-select box selected data
// function handleSelectionChange(selection) {
// this.ids = selection.map((item) => item.dictCode)
// this.single = selection.length != 1
// this.multiple = !selection.length
// }
/** Modify the button operation */
function handleUpdate(row) {
  reset()
  const dictCode = row.dictCode || ids.value
  getData(dictCode).then((response) => {
    form.value = response.data
    open.value = true
    title.value = '修改字典数据'
  })
}
/** Submit button */
function submitForm() {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.dictCode != undefined) {
        updateData(form.value).then((response) => {
          proxy.$modal.msgSuccess('修改成功')
          open.value = false
          getList()
        })
      } else {
        addData(form.value).then((response) => {
          proxy.$modal.msgSuccess('新增成功')
          open.value = false
          getList()
        })
      }
    }
  })
}

/** delete button action */
function handleDelete(row) {
  const dictCodes = row.dictCode || ids.value
  proxy
    .$confirm('是否确认删除字典编码为"' + dictCodes + '"的数据项?', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    .then(function () {
      return delData(dictCodes)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
}
</script>
