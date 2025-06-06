<template>
  <div class="app-container">
    <div>
      <el-tabs type="border-card" style="height: 80vh">
        <el-tab-pane label="Add Phone">
          <el-row :gutter="24">
            <el-col :span="4" :xs="24">
              <div class="head-container"></div>
            </el-col>
            <el-col :lg="24" :xm="24">
              <el-row :gutter="10" class="mb8">
                <el-col :span="1.5">
                  <el-button type="primary" plain icon="plus" @click="handleAdd">{{ $t('btn.add') }}</el-button>
                </el-col>
                <el-col :span="1.5">
                  <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="Delete">
                    {{ $t('btn.delete') }}
                  </el-button>
                </el-col>

                <el-col :span="1.5">
                  <el-input v-model="queryParams.query" placeholder="filter by account , chatid .." clearable style="margin-bottom: 20px" />
                </el-col>
                <el-col :span="1.5">
                  <el-button type="primary" icon="search" @click="loadList(1)">{{ $t('btn.search') }}</el-button>
                </el-col>
              </el-row>
              <div type="border-card" class="demo-tabs" style="min-height: 70vh; margin-top: 10px">
                <el-table v-loading="loading" :data="userList" @selection-change="handleSelectionChange">
                  <el-table-column type="selection" width="50" align="center" />
                  <el-table-column label="Number" align="left" sortable key="phone" prop="phone" />

                  <el-table-column label="Bot name" align="center" sortable key="bot_name" prop="bot_name" />

                  <el-table-column :label="proxy.$t('lables.status')" align="center" sortable key="status" prop="status" />

                  <el-table-column :label="proxy.$t('lables.created_on')" align="center" key="created_at" prop="created_at">
                    <template #default="scope">
                      {{ FormatDate24(scope?.row?.created_at) }}
                    </template>
                  </el-table-column>
                  <el-table-column :label="proxy.$t('lables.updated_on')" align="center" key="updated_at" prop="updated_at">
                    <template #default="scope">
                      {{ FormatDate24(scope?.row?.updated_at) }}
                    </template>
                  </el-table-column>
                  <el-table-column :label="proxy.$t('action')" align="right" class-name="small-padding fixed-width">
                    <template #default="scope">
                      <el-button type="primary" size="small" icon="edit" :title="$t('btn.edit')" @click.stop="handleUpdate(scope.row)"> </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <br />
              <el-pagination
                v-model:current-page="queryParams.page"
                v-model:page-size="queryParams.limit"
                :page-sizes="[15, 20, 30, 40, 50, 60, 70, 80, 90, 100, 120]"
                :small="true"
                :disabled="disabled"
                :background="background"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="handleChangeSize"
                @current-change="handlePageChangeSize" />
            </el-col>
          </el-row>
        </el-tab-pane>
        <el-tab-pane label="History" tab-click="loadListHistory">
          <el-row :gutter="24">
            <el-col :span="4" :xs="24">
              <div class="head-container"></div>
            </el-col>
            <el-col :lg="24" :xm="24">
              <el-row :gutter="10" class="mb8">
                <el-col :span="1.5">
                  <el-button type="primary" plain icon="UploadFilled" @click="exportToExcel">{{ $t('btn.export') }}</el-button>
                </el-col>
                <el-col :span="1.5"
                  ><el-date-picker
                    v-model="queryParams.start_date"
                    type="datetime"
                    placeholder="From Date"
                    format="YYYY-MM-DD HH:mm:ss"
                    value-format="YYYY-MM-DD HH:mm:ss"
                /></el-col>
                <el-col :span="1.5"
                  ><el-date-picker
                    v-model="queryParams.end_date"
                    type="datetime"
                    placeholder="To Date"
                    format="YYYY-MM-DD HH:mm:ss"
                    value-format="YYYY-MM-DD HH:mm:ss"
                /></el-col>
                <el-col :span="1.5">
                  <el-input v-model="queryParams.query" placeholder="filter phone .." clearable style="margin-bottom: 20px" />
                </el-col>
                <el-col :span="1.5">
                  <el-button type="primary" icon="search" @click="loadListHistory(1)">{{ $t('btn.search') }}</el-button>
                </el-col>
              </el-row>
              <div type="border-card" class="demo-tabs" style="min-height: 70vh; margin-top: 10px">
                <el-table v-loading="loading" :data="userListHistory">
                  <el-table-column label="Number" align="left" sortable key="phone" prop="phone" />

                  <el-table-column label="Bot name" align="center" sortable key="bot_name" prop="bot_name" />
                  <el-table-column label="Request" align="center" sortable key="requester" prop="requester" />
                  <el-table-column label="Request" align="center" sortable key="request_date" prop="request_date" />
                  <el-table-column :label="proxy.$t('lables.status')" align="center" key="status" prop="status" />

                  <el-table-column :label="proxy.$t('lables.created_on')" align="center" key="created_at" prop="created_at">
                    <template #default="scope">
                      {{ FormatDate24(scope?.row?.created_at) }}
                    </template>
                  </el-table-column>
                  <el-table-column :label="proxy.$t('lables.updated_on')" align="center" key="updated_at" prop="updated_at">
                    <template #default="scope">
                      {{ FormatDate24(scope?.row?.updated_at) }}
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <br />
              <el-pagination
                v-model:current-page="queryParams.page"
                v-model:page-size="queryParams.limit"
                :page-sizes="[15, 20, 30, 40, 50, 60, 70, 80, 90, 100, 120]"
                :small="true"
                :disabled="disabled"
                :background="background"
                layout="total, sizes, prev, pager, next, jumper"
                :total="total"
                @size-change="handleChangeSize"
                @current-change="handlePageChangeSize" />
            </el-col>
          </el-row>
        </el-tab-pane>
      </el-tabs>

      <el-dialog :title="title" v-model="open" width="600px" append-to-body :close-on-click-modal="false" draggable>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
          <el-row>
            <el-col :lg="24">
              <el-form-item label="Phone" prop="phone">
                <el-input v-model="form.phone" :rows="8" type="textarea" placeholder="049485495,012375868,...." />
              </el-form-item>
              <el-form-item label="Telegram Bot:" prop="bot_name">
                <el-input v-model="form.bot_name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.status')" prop="status">
                <el-select v-model="form.status" placeholder="Select" style="width: 150px">
                  <el-option v-for="item in statusOptions" :key="item.status_value" :label="item.status_name" :value="item.status_value">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <template #footer>
          <el-button text @click="open = false">{{ $t('btn.cancel') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('btn.submit') }}</el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup name="user">
import { getPhoneRecords, savePhone, updatePhone, deletePhone } from '@/api/telegram/telegram'
import { onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as XLSX from 'xlsx'
const { proxy } = getCurrentInstance()
const userList = ref([])
const userListHistory = ref([])
const loading = ref(true)
const open = ref(false)
const total = ref(0)
const ids = ref([])
const multiple = ref(true)
const title = ref('')
const current = ref(1)
const pageSize = ref(15)
const data = reactive({
  queryParams: {
    page: 0, //num_page
    limit: 15, //page
    start_date: '',
    end_date: '',
    query: ''
  },
  form: {}
})
const statusOptions = ref()
statusOptions.value = [
  {
    status_value: 'no',
    status_name: proxy.$t('common.no')
  },
  {
    status_value: 'yes',
    status_name: proxy.$t('common.yes')
  }
]
const trackOptions = ref()
trackOptions.value = [
  {
    status_value: 'normal',
    status_name: proxy.$t('common.no')
  },
  {
    status_value: 'track',
    status_name: proxy.$t('common.yes')
  }
]
const { queryParams, form } = toRefs(data)
const FormatDate24 = (dateString) => {
  const date = new Date(dateString)
  const formattedDate = date.toLocaleDateString('en-CA') // Format to YYYY-MM-DD
  const formattedTime = date.toLocaleTimeString('en-CA', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',

    hour12: false // Include AM/PM
  })
  return `${formattedDate}, ${formattedTime}` // Combine both parts
}
const exportToExcel = () => {
  const date = new Date()
  // Create a new workbook
  const wb = XLSX.utils.book_new()

  // Create a worksheet from userList data
  const ws = XLSX.utils.json_to_sheet(userList.value)

  // Append the worksheet to the workbook
  XLSX.utils.book_append_sheet(wb, ws, 'User Data')

  // Generate a file name
  const fileName = `request-history-${date}.xlsx`

  // Write the workbook and trigger a download
  XLSX.writeFile(wb, fileName)
}
const formatTime = (minutes) => {
  if (minutes < 60) {
    return `${minutes}min`
  } else {
    const hours = Math.floor(minutes / 60)
    const remainingMinutes = minutes % 60
    return `${hours}h ${remainingMinutes} min`
  }
}
const loadList = (_page = 1) => {
  loading.value = true
  data.queryParams.status = 'yes'
  getPhoneRecords(data.queryParams).then((res) => {
    loading.value = false
    userList.value = res.data
    total.value = res.meta.total
  })
}
const loadListHistory = (_page = 1) => {
  loading.value = true
  data.queryParams.status = 'used'
  getPhoneRecords(data.queryParams).then((res) => {
    loading.value = false
    userListHistory.value = res.data
    total.value = res.meta.total
  })
}
const handleChangeSize = (val) => {
  current.value = val
  proxy.$nextTick(() => {
    loadList(val)
  })
}
const handlePageChangeSize = (val) => {
  pageSize.value = val
  proxy.$nextTick(() => {
    loadList()
  })
}
const rsetForm = () => {
  form.value = {
    name: '',
    id: undefined
  }
}
const handleAdd = () => {
  rsetForm()
  open.value = true
  title.value = proxy.$t('system.Add New')
}
function handleUpdate(row) {
  console.log(row)
  form.value = row
  open.value = true
  title.value = 'Update Record'
}
const submitForm = () => {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      const phoneNumbers = form.value.phone.split(',')
      const phoneListEntries = phoneNumbers.map((phone) => ({
        phone: phone.trim(), // Remove any extra spaces
        bot_name: form.value.bot_name, // Replace with actual bot name if necessary
        status: form.value.status
      }))
      if (form.value.id !== undefined && form.value.id > 0) {
        //    form.value.type = 'edit'
        updatePhone(form.value).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            loadList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
        })
      } else {
        //    form.value.type = 'add'

        savePhone(phoneListEntries).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
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
function handleSelectionChange(selection) {
  ids.value = selection.map((item) => item.id)
  multiple.value = !selection.length
}
const Delete = (id) => {
  ElMessageBox.confirm('proxy will permanently delete the records. Continue?', 'Warning', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel',
    type: 'warning'
  })
    .then(() => {
      deletePhone({ id: ids.value })
      ElMessage({
        type: 'success',
        message: 'Delete completed'
      })
      loadList()
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: 'Delete canceled'
      })
    })
}
onMounted(() => {
  loadList()
  loadListHistory()
})
</script>
