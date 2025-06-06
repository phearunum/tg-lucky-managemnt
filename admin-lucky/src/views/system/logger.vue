<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-col :span="4" :xs="24">
        <div class="head-container"></div>
      </el-col>
      <el-col :lg="24" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-input v-model="queryParams.query" style="max-width: 150px" placeholder="Username">
              <template #prepend>Filter</template>
            </el-input>
          </el-col>
          <el-col :span="1.5">
            <el-input v-model="queryParams.limit" style="max-width: 150px" placeholder="Please input">
              <template #prepend>Row</template>
            </el-input>
          </el-col>
          <el-col :span="1.5">
            <el-date-picker v-model="queryParams.fdate" format="YYYY/MM/DD" type="date" placeholder="From date" />
          </el-col>
          <el-col :span="1.5">
            <el-date-picker v-model="queryParams.tdate" type="date" format="YYYY/MM/DD" placeholder="End data" />
          </el-col>
          <el-col :span="1.5">
            <el-button type="primary" icon="search" @click="loadList">{{ $t('btn.search') }}</el-button>
          </el-col>
          <el-col :span="1.5" align="right">
            <!--
               <el-pagination
              v-model:current-page="queryParams.page"
              v-model:page-size="queryParams.limit"
              :disabled="disabled"
              :background="background"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              :small="true"
              @size-change="handleChangeSize"
              @current-change="handlePageChangeSize" />
            -->
          </el-col>
        </el-row>
        <el-tabs type="border-card" class="demo-tabs" style="min-height: 75vh; margin-top: 10px">
          <el-tab-pane label="Logger">
            <el-table v-loading="loading" :data="userList">
              <el-table-column type="expand" style="background-color: black">
                <template #default="props">
                  <div style="background-color: black; color: white; padding: 6px">
                    <p m="t-0 b-2">Logger Old Data =>[ {{ props.row.additionalInfo }} ]</p>
                    <p m="t-0 b-2">Logger New Data => [{{ props.row.data }}]</p>
                  </div>
                </template>
              </el-table-column>

              <el-table-column label="Logger" key="api" prop="api" width="100" />
              <el-table-column label="Action" key="action" prop="action" width="150" />
              <el-table-column label="Date" key="timestamp" prop="timestamp" width="200" />
              <el-table-column label="By" key="userId" prop="userId" width="50" />

              <el-table-column label="Other" key="additionalInfo" prop="additionalInfo" align="left">
                <template #default="props">
                  <span>{{ truncateString(JSON.stringify(props.row.additionalInfo, null, 2)) }}</span>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
        <br />
      </el-col>
    </el-row>
  </div>
</template>

<script setup name="daily-meal">
import { getLogger } from '@/api/system/logger'
const { proxy } = getCurrentInstance()
const userList = ref([])
const ids = ref([])
const loading = ref(true)
const total = ref(0)
const title = ref('')
const open = ref(false)
const current = ref(1)
const pageSize = ref(10)
const deptOptions = ref([])
const empOption = ref([])
const ExcelData = ref([])
const statusOptions = ref()
statusOptions.value = [
  {
    status_value: false,
    status_name: 'Disabled'
  },
  {
    status_value: true,
    status_name: 'Active'
  }
]
function truncateString(input, maxLength = 150) {
  const str = String(input) // Ensure the input is a string
  return str.length <= maxLength ? str : str.slice(0, maxLength) + '...'
}
const data = reactive({
  queryParams: {
    page: 1, //num_page
    limit: 100, //page
    query: undefined,
    type: '',
    fdate: new Date().toJSON().slice(0, 10),
    tdate: new Date().toJSON().slice(0, 10)
  },
  form: {}
})

const { queryParams, form } = toRefs(data)
const router = useRouter()
function addNew(Id = 0) {
  router.push({ path: `/products/categoryadd` })
}

const loadList = (_page = 0) => {
  loading.value = true
  getLogger(queryParams.value).then((res) => {
    loading.value = false
    userList.value = res.data
    total.value = res.meta.total
  })
}
const handleChangeSize = (val) => {
  current.value = val
  proxy.$nextTick(() => {
    loadList()
  })
}
const handlePageChangeSize = (val) => {
  console.log(val)
  pageSize.value = val
  proxy.$nextTick(() => {
    loadList(val)
  })
}
const handleAdd = () => {
  open.value = true
  title.value = proxy.$t('system.Add New')
}
function handleUpdate(row) {
  const id = row.id || ids.value
  getCategory(id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = 'Update Record'
  })
}
const submitForm = () => {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.id !== undefined && form.value.id > 0) {
        form.value.type = 'edit'
        updateCategory(form.value).then((response) => {
          open.value = false
          if (response.status == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            loadList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
        })
      } else {
        form.value.type = 'add'
        addCategory(form.value).then((response) => {
          open.value = false
          if (response.status == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            loadList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
        })
      }
    }
  })
}

onMounted(() => {
  loadList()
})
</script>
