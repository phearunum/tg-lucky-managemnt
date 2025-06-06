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
              <el-button type="primary" plain icon="plus" @click="handleAdd">{{ $t('btn.add') }}</el-button>
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
              <el-table-column :label="proxy.$t('lables.name')" align="left" sortable key="name" prop="name" />
              <el-table-column :label="proxy.$t('lables.allow')" align="center" sortable key="time_rest" width="180" prop="time_rest">
                <template #default="scope">
                  <el-tag type="primary" style="width: 60px">{{ formatTime(scope?.row?.time_rest) }} </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('lables.type')" align="center" sortable key="button_type" prop="button_type" />
              <el-table-column label="Group ID" align="center" sortable key="bot_name" prop="bot_name" />
              <el-table-column :label="proxy.$t('lables.sort')" align="center" sortable key="order_no" prop="order_no" />
              <el-table-column :label="proxy.$t('lables.status')" align="center" sortable key="status" prop="status" />

              <el-table-column :label="proxy.$t('lables.created_on')" align="center" sortable key="created_at" prop="created_at">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.created_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('lables.updated_on')" align="center" sortable key="updated_at" prop="updated_at">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.updated_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('action')" align="right" class-name="small-padding fixed-width">
                <template #default="scope">
                  <el-button type="primary" icon="edit" :title="$t('btn.edit')" @click.stop="handleUpdate(scope.row)"> </el-button>
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
      <el-dialog :title="title" v-model="open" width="600px" append-to-body :close-on-click-modal="false" draggable>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
          <el-row>
            <el-col :lg="24">
              <el-form-item :label="proxy.$t('server.stream.name')" prop="name">
                <el-input v-model="form.id" style="display: none" />
                <el-input v-model="form.name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Group ID" prop="bot_name">
                <el-input v-model="form.bot_name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.allow')" prop="time_rest">
                <el-input-number v-model="form.time_rest" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.sort')" prop="order_no">
                <el-input-number v-model="form.order_no" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
              </el-form-item>

              <el-form-item :label="proxy.$t('lables.type')" prop="button_type">
                <el-select v-model="form.button_type" placeholder="Select" style="width: 150px">
                  <el-option v-for="item in trackOptions" :key="item.status_value" :label="item.status_name" :value="item.status_value">
                  </el-option>
                </el-select>
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
    </el-card>
  </div>
</template>

<script setup name="user">
import { getTelegramSetting, saveTelegramSetting, updateTelegramSetting } from '@/api/telegram/telegram'
import { onMounted } from 'vue'
const { proxy } = getCurrentInstance()
const userList = ref([])
const loading = ref(true)
const open = ref(false)
const total = ref(0)
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
    status_name: 'Normal'
  },
  {
    status_value: 'track',
    status_name: 'Track'
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
  getTelegramSetting(data.queryParams).then((res) => {
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
      if (form.value.id !== undefined && form.value.id > 0) {
        //    form.value.type = 'edit'
        updateTelegramSetting(form.value).then((response) => {
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

        saveTelegramSetting(form.value).then((response) => {
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

onMounted(() => {
  loadList()
})
</script>
