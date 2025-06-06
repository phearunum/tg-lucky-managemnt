<template>
  <div class="app-container">
    <el-card shadow="never">
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
              <el-input v-model="queryParams.query" placeholder="filter by account , chatid .." clearable style="margin-bottom: 20px" />
            </el-col>
            <el-col :span="1.5">
              <el-button type="primary" icon="search" @click="loadList(1)">{{ $t('btn.search') }}</el-button>
            </el-col>
          </el-row>
          <div type="border-card" class="demo-tabs" style="min-height: 70vh; margin-top: 10px">
            <el-table v-loading="loading" :data="userList">
              <el-table-column :label="proxy.$t('Account Name')" align="left" sortable key="account_name" prop="account_name" />
              <el-table-column :label="proxy.$t('Chat ID')" align="center" sortable key="chat_id" width="180" prop="chat_id" />
              <el-table-column :label="proxy.$t('Leave')" align="center" sortable key="start_time" prop="start_time" />
              <el-table-column :label="proxy.$t('Back')" align="center" sortable key="end_time" prop="end_time" />

              <el-table-column :label="proxy.$t('Total')" align="left" sortable key="total_time" prop="total_time">
                <template #default="scope">
                  <el-tag type="primary" style="min-width: 60px">{{ formatTime(parseFloat(scope?.row?.total_time).toFixed(0)) }} </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('Request')" align="left" sortable key="request_type" prop="request_type">
                <template #default="scope">
                  <el-tag type="primary" style="width: 60px">{{ formatTime(scope?.row?.allow_time) }} </el-tag> =
                  {{ scope?.row?.request_type }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('message')" align="right" sortable key="message" prop="message" />
              <el-table-column :label="proxy.$t('bot_name')" align="right" sortable key="bot_name" prop="bot_name" />
              <!-- 
                  <el-table-column :label="proxy.$t('Start')" align="center" sortable key="created_at" prop="created_at">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.created_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('End')" align="center" sortable key="created_at" prop="created_at">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.updated_at) }}
                </template>
              </el-table-column> 
             -->
            </el-table>
          </div>
          <br />
          <el-pagination
            v-model:current-page="queryParams.page"
            v-model:page-size="queryParams.limit"
            :page-sizes="[15, 20, 30, 40, 50, 60, 70, 80, 90, 100, 120, 200, 300, 500, 1000]"
            :small="true"
            :disabled="disabled"
            :background="background"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleChangeSize"
            @current-change="handlePageChangeSize" />
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup name="user">
import { getTelegram } from '@/api/telegram/telegram'
import * as XLSX from 'xlsx'
import { onMounted } from 'vue'
const { proxy } = getCurrentInstance()
const userList = ref([])
const loading = ref(true)
const total = ref(0)
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
  // Create a new workbook
  const wb = XLSX.utils.book_new()
  const date = new Date()
  // Create a worksheet from userList data
  const ws = XLSX.utils.json_to_sheet(userList.value)

  // Append the worksheet to the workbook
  XLSX.utils.book_append_sheet(wb, ws, 'User Data')

  // Generate a file name
  const fileName = `break_time-${date}.xlsx`

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
  getTelegram(data.queryParams).then((res) => {
    loading.value = false
    userList.value = res.data
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
    catName: undefined,
    catCode: undefined
  }
}

onMounted(() => {
  loadList()
})
</script>
