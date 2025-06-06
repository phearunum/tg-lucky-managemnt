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
          <div type="border-card" class="demo-tabs" style="min-height: 70vh; margin-top: -10px">
            <el-table v-loading="loading" :data="userList">
              <el-table-column prop="name" label="Account" align="left" width="180" style="text-align: left"> </el-table-column>
              <el-table-column :label="proxy.$t('Group ID')" align="left" sortable key="group_id" width="180" prop="group_id" />
              <el-table-column :label="proxy.$t('lables.token')" align="left" key="token" prop="token" width="350" />
              <el-table-column label="Win" align="right" key="total_win" prop="total_win" />
              <!-- <el-table-column :label="proxy.$t('Type')" align="left" key="bot_type" width="100" prop="bot_type" /> -->
              <el-table-column :label="proxy.$t('status')" align="left" key="status" width="100" prop="status" />
              <el-table-column :label="proxy.$t('message')" align="left" key="message" width="250" prop="message" />
              <el-table-column :label="proxy.$t('Image')" align="right" key="image" prop="image" width="170">
                <template #default="scope">
                  <el-image :src="scope.row.image" />
                </template>
              </el-table-column>

              <el-table-column :label="proxy.$t('action')" align="right" width="100">
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
      <el-dialog :title="title" v-model="open" width="600px" append-to-body :close-on-click-modal="false" draggable>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="150px">
          <el-row>
            <el-col :lg="24">
              <el-form-item :label="proxy.$t('lables.name')" prop="name">
                <el-input v-model="form.name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Group ID" prop="group_id">
                <el-input v-model="form.group_id" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Token" prop="token">
                <el-input v-model="form.token" type="textarea" :rows="2" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Message" prop="message">
                <el-input v-model="form.message" type="textarea" :rows="2" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Image" prop="image">
                <el-input v-model="form.image" type="textarea" :rows="2" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item label="Total Win" prop="total_win">
                <el-input-number v-model="form.total_win" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
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
import { getLuckySetting, updateLuckySetting, saveLuckySetting } from '@/api/lucky/luckysetting.service'
import { Picture as IconPicture } from '@element-plus/icons-vue'
import { onMounted } from 'vue'
import { getToken, getUserInfo, setToken, removeToken } from '@/utils/auth'

const { proxy } = getCurrentInstance()
const userList = ref([])
const optionMember = ref([])
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
    status_value: 'expired',
    status_name: 'expired'
  },
  {
    status_value: 'available',
    status_name: 'available'
  },
  {
    status_value: 'collected',
    status_name: 'collected'
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
  getLuckySetting(data.queryParams).then((res) => {
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
  console.log(row)
  form.value = row
  open.value = true
  title.value = 'Update Record'
}
const memberOption = (filter) => {
  data.queryParams.query = filter
  getTelegramMember(data.queryParams).then((res) => {
    loading.value = false
    optionMember.value = res.data
  })
}
const submitForm = () => {
  const { userId, roleId, companyId, username } = JSON.parse(getUserInfo())
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      form.value.create_by = username
      if (form.value.id !== undefined && form.value.id > 0) {
        //    form.value.type = 'edit'
        updateLuckySetting(form.value).then((response) => {
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

        saveLuckySetting(form.value).then((response) => {
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
<style scoped>
::v-deep(.el-table th) {
  font-size: 0.9rem;
  text-align: center;
  line-height: 1;
  white-space: normal;
  word-break: break-word;
}
::v-deep(.el-table td) {
  text-align: center;
  align-content: center;
  padding: 0px;
}
</style>
