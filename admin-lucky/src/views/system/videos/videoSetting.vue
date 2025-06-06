<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-col :span="4" :xs="24">
        <div class="head-container"></div>
      </el-col>
      <el-col :lg="24" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" plain icon="plus" @click="handleAdd" v-hasPermi="['system:server:add']">{{ $t('btn.add') }}</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-input v-model="queryParams.q" placeholder="Search" clearable prefix-icon="el-icon-search" style="margin-bottom: 20px" />
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
            <el-table-column :label="proxy.$t('server.stream.name')" align="center" sortable key="server_name" prop="server_name" />
            <el-table-column :label="proxy.$t('server.stream.rtmp')" align="center" sortable key="rtmp" width="180" prop="rtmp" />
            <el-table-column :label="proxy.$t('server.stream.size')" align="center" sortable key="video_size" prop="video_size" />
            <el-table-column :label="proxy.$t('server.stream.type')" align="center" sortable key="video_type" prop="video_type" />
            <el-table-column :label="proxy.$t('server.stream.duration')" align="center" sortable key="duration" prop="duration" />
            <el-table-column :label="proxy.$t('server.stream.before_start')" align="center" sortable key="before_start" prop="before_start" />
            <el-table-column :label="proxy.$t('server.stream.after_end')" align="center" sortable key="after_end" prop="after_end" />
            <el-table-column :label="proxy.$t('server.stream.bucket_name')" align="center" sortable key="bucket_name" prop="bucket_name" />
            <!-- <el-table-column
              :label="proxy.$t('server.stream.service_account')"
              align="center"
              sortable
              key="service_account"
              prop="service_account" /> -->
            <el-table-column :label="proxy.$t('server.stream.output_path')" align="center" sortable key="output_path" prop="output_path" />
            <el-table-column
              :label="proxy.$t('server.stream.delete_local_store')"
              align="center"
              sortable
              key="delete_local_store"
              prop="delete_local_store" />

            <el-table-column
              :label="proxy.$t('create on')"
              align="center"
              sortable
              prop="create_at"
              v-if="columns.showColumn('create_at')"
              width="200" />

            <el-table-column fixed="right" :label="proxy.$t('action')" align="center" width="150" class-name="small-padding fixed-width">
              <template #default="scope">
                <el-button
                  size="small"
                  text
                  icon="edit"
                  :title="$t('btn.edit')"
                  @click.stop="handleUpdate(scope.row)"
                  v-hasPermi="['system:server:update']">
                </el-button>
                <el-button
                  v-if="scope.row.id > 1"
                  size="small"
                  text
                  icon="delete"
                  :title="$t('btn.delete')"
                  @click.stop="handleDelete(scope.row)"
                  v-hasPermi="['system:server:delete']">
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <br />
        <el-pagination
          v-model:current-page="queryParams.num_page"
          v-model:page-size="queryParams.limit"
          :page-sizes="[100, 200, 300, 400]"
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
              <el-form-item :label="proxy.$t('server.stream.access_domain')" prop="access_domain">
                <el-input v-model="form.access_domain" :placeholder="proxy.$t('system.Enter')" />
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
              <el-form-item label="Prefix" prop="prefix">
                <el-input v-model="form.prefix" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>

              <el-form-item :label="proxy.$t('server.stream.delete_local_store')" prop="delete_local_store">
                <el-select v-model="form.delete_local_store" placeholder="Select" style="width: 130px">
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
    </el-row>
  </div>
</template>

<script setup name="user">
import { videoSettingList, videoSettingByID, addVideoSetting, updateVideoSetting, deleteVideoSetting } from '@/api/system/videosnap'
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
    status_name: proxy.$t('common.no')
  },
  {
    status_value: true,
    status_name: proxy.$t('common.yes')
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
    limit: 10, //page
    q: ''
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
const loadList = (_page = 1) => {
  loading.value = true
  videoSettingList(data.queryParams).then((res) => {
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
      return deleteVideoSetting(ids)
    })
    .then((response) => {
      if (response.statusCode == 200) {
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
