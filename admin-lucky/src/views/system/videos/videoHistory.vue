<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-col :span="4" :xs="24">
        <div class="head-container"></div>
      </el-col>
      <el-col :lg="24" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-input
              v-model="queryParams.query"
              placeholder="Enter GameNo .."
              clearable
              prefix-icon="el-icon-search"
              style="margin-bottom: 20px" />
          </el-col>
          <el-col :span="1.5">
            <el-button type="primary" icon="search" @click="loadList(1)">{{ $t('btn.search') }}</el-button>
          </el-col>
        </el-row>
        <div type="border-card" class="demo-tabs" style="min-height: 70vh; margin-top: 10px">
          <el-table
            v-loading="loading"
            :data="userList"
            v-hasPermi="['setting:video:view']"
            :header-cell-style="{
              background: 'var(--el-fill-color-light)',
              color: 'var(--el-text-color-primary)'
            }">
            <el-table-column :label="proxy.$t('server.stream.id')" align="center" sortable key="id" width="100" prop="id" />
            <el-table-column :label="proxy.$t('server.stream.game_no')" align="center" sortable key="period" prop="period" min-width="200" />
            <el-table-column :label="proxy.$t('server.stream.process_id')" align="center" sortable key="process_id" prop="process_id" width="120" />
            <el-table-column :label="proxy.$t('server.stream.rtmp')" align="center" sortable key="rtmp" prop="rtmp" />
            <el-table-column :label="proxy.$t('server.stream.stream_key')" align="center" sortable key="rtmpurl" prop="rtmpurl" width="250" />
            <el-table-column :label="proxy.$t('server.stream.image_url')" align="center" sortable key="image_url" prop="image_url" width="120">
              <template #default="scope">
                <a :href="scope.row.image_url" target="_blank">
                  <el-icon style="font-size: 1.5rem"><Picture /></el-icon>
                </a>
              </template>
            </el-table-column>

            <el-table-column :label="proxy.$t('server.stream.video_url')" align="center" sortable key="video_url" prop="video_url" width="100">
              <template #default="scope">
                <a :href="scope.row.video_url" target="_blank">
                  <el-icon style="font-size: 1.5rem"><VideoPlay /></el-icon>
                </a>
              </template>
            </el-table-column>
            <el-table-column :label="proxy.$t('common.status')" align="center" sortable key="status" prop="status" width="80">
              <template #default="scope">
                <el-tag class="ml-2" v-if="scope.row.status == true" type="success" title="Active"
                  ><el-icon><SuccessFilled /></el-icon
                ></el-tag>
                <el-tag class="ml-2" v-if="scope.row.status == false" type="danger" title="Disabled"
                  ><el-icon><CircleCloseFilled /></el-icon
                ></el-tag>
              </template>
            </el-table-column>
            <!-- <el-table-column
                :label="proxy.$t('server.stream.service_account')"
                align="center"
                sortable
                key="service_account"
                prop="service_account" /> -->

            <el-table-column
              :label="proxy.$t('create on')"
              align="center"
              sortable
              prop="created_at"
              v-if="columns.showColumn('create_at')"
              width="150">
              <template #default="scope">
                <span>{{ FormatDate24(scope?.row?.created_at) }}</span>
              </template>
            </el-table-column>
            <el-table-column
              :label="proxy.$t('create on')"
              align="center"
              sortable
              prop="updated_at"
              v-if="columns.showColumn('create_at')"
              width="150">
              <template #default="scope">
                <span>{{ FormatDate24(scope?.row?.updated_at) }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <br />
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.limit"
          :page-sizes="[40, 50, 60, 80, 90, 100, 150, 200, 250]"
          :small="true"
          :disabled="disabled"
          :background="background"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleChangeSize"
          @current-change="handlePageChangeSize" />
      </el-col>

      <el-dialog :title="title" v-model="open" width="600px" append-to-body :close-on-click-modal="false">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
          <el-row>
            <el-col :lg="24">
              <el-form-item :label="proxy.$t('server.stream.name')" prop="server_name">
                <el-input v-model="form.server_name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.rtmp')" prop="rtmp">
                <el-input v-model="form.rtmp" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.duration')" prop="duration">
                <el-input-number v-model="form.duration" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
              </el-form-item>

              <el-form-item :label="proxy.$t('server.stream.size')" prop="video_size">
                <el-input v-model="form.video_size" :placeholder="proxy.$t('system.Enter')" style="width: 120px" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.type')" prop="video_type">
                <el-input v-model="form.video_type" :placeholder="proxy.$t('system.Enter')" style="width: 120px" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.before_start')" prop="before_start">
                <el-input-number v-model="form.before_start" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.after_end')" prop="after_end">
                <el-input-number v-model="form.after_end" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
              </el-form-item>

              <el-form-item :label="proxy.$t('server.stream.bucket_name')" prop="bucket_name">
                <el-input v-model="form.bucket_name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.service_account')" prop="service_account">
                <el-input v-model="form.service_account" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.stream.output_path')" prop="output_path">
                <el-input v-model="form.output_path" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <template #footer>
          <el-button text @click="open = false">{{ $t('btn.cancel') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('btn.submit') }}</el-button>
        </template>
      </el-dialog>
    </el-row>
  </div>
</template>

<script setup name="user">
import { videoRecordsList, videoRecordsByID } from '@/api/system/videohistory'
import { videoSettingList } from '@/api/system/videosnap'
import { onMounted } from 'vue'
const { proxy } = getCurrentInstance()
const userList = ref([])
const ids = ref([])
const loading = ref(true)
const total = ref(0)
const title = ref('')
const open = ref(false)
const current = ref(1)
const pageSize = ref(10)
const columns = ref([])
const statusOptions = ref()
statusOptions.value = [
  {
    status_value: false,
    status_name: proxy.$t('system.Disable')
  },
  {
    status_value: true,
    status_name: proxy.$t('system.Active')
  }
]
const unPaidOption = ref()
unPaidOption.value = [
  {
    status_value: false,
    status_name: proxy.$t('common.no')
  },
  {
    status_value: true,
    status_name: proxy.$t('common.yes')
  }
]
const data = reactive({
  queryParams: {
    page: 0, //num_page
    limit: 20, //page
    query: ''
  },
  form: {}
})

const { queryParams, form } = toRefs(data)
const FormatDate24 = (dateString) => {
  if (!dateString) return 'N/A'

  // Trim to remove any hidden characters
  const cleanedDateString = dateString.trim()
  console.log('Parsing:', JSON.stringify(cleanedDateString)) // Debugging

  // Try parsing the date
  const date = new Date(cleanedDateString)
  if (isNaN(date.getTime())) {
    console.error('Invalid Date Detected:', cleanedDateString)
    return 'Invalid Date'
  }

  // Format output
  const formattedDate = date.toLocaleDateString('en-CA') // YYYY-MM-DD
  const formattedTime = date
    .toLocaleTimeString('en-CA', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false // 24-hour format
    })
    .replace(/\s/g, '') // Remove unexpected spaces

  return `${formattedDate}, ${formattedTime}`
}

const loadList = (_page = 1) => {
  loading.value = true
  videoRecordsList(data.queryParams).then((res) => {
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
const handleAdd = () => {
  rsetForm()
  open.value = true
  title.value = proxy.$t('system.Add New')
}
function handleUpdate(row) {
  const id = row.id || ids.value
  videoSettingByID(id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = 'Update Record'
  })
}
const submitForm = () => {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.id !== undefined && form.value.id > 0) {
        //    form.value.type = 'edit'
        updateVideoSetting(form.value).then((response) => {
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

        addVideoSetting(form.value).then((response) => {
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

const handleDelete = (row) => {
  const ids = row.id || ids.value
  proxy
    .$confirm(proxy.$t('system.Do you want to delete') + ` = ${ids} ?`, proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('yes'),
      cancelButtonText: proxy.$t('no'),
      type: 'warning'
    })
    .then(() => {
      return deleteleaveType(ids)
    })
    .then((response) => {
      if (response.status == 'OK') {
        proxy.$modal.msgSuccess(proxy.$t('system.Success'))
        loadList()
      } else {
        proxy.$modal.msgError(proxy.$t('system.Failed'))
      }
    })
}

onMounted(() => {
  loadList()
})
</script>
