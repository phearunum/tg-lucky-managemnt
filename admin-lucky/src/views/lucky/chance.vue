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
              <el-table-column prop="profile" label="Account" align="left" width="180" style="text-align: left">
                <template #default="scope">
                  <el-row>
                    <el-col :span="8">
                      <el-image
                        :src="scope?.row?.profile"
                        style="border-radius: 30%; width: 35px; height: 35px; margin-top: 5px; background-color: #f2f2f2; margin-left: -20px">
                        <template #error>
                          <center>
                            <el-icon><icon-picture /></el-icon>
                          </center>
                        </template>
                      </el-image>
                    </el-col>
                    <el-col :span="16" style="line-height: 10px; text-align: left; align-items: center; padding-top: 8px">
                      <span>
                        {{ scope?.row?.account_name }}<br />
                        <span style="font-size: 10px; line-height: 20px"> {{ scope?.row?.employee_id }}</span>
                      </span>
                    </el-col>
                  </el-row>
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('Chat ID')" align="left" sortable key="chat_id" width="250" prop="chat_id" />
              <el-table-column :label="proxy.$t('lables.total_chance')" align="left" key="total_chance" prop="total_chance">
                <template #default="scope">
                  <el-tag type="primary" style="min-width: 60px">
                    {{ scope.row?.total_chance }} <el-icon style="color: orangered"><Orange /></el-icon>
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('lables.status')" align="left" key="status" prop="status">
                <template #default="scope">
                  <el-tag type="primary" style="min-width: 75px" v-if="scope.row?.status == 'available'">
                    {{ scope.row?.status }} <el-icon><CircleCheck /></el-icon
                  ></el-tag>
                  <el-tag type="danger" v-else style="min-width: 75px">
                    {{ scope.row?.status }} <el-icon><CircleClose /></el-icon>
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('Expired')" align="left" sortable key="expired" width="100" prop="expired" />
              <el-table-column :label="proxy.$t('created_at')" align="right" key="created_at" prop="created_at" width="170">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.created_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('updated_at')" align="right" key="updated_at" prop="updated_at" width="170">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.updated_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('action')" align="right" width="100">
                <template #default="scope">
                  <el-button
                    type="primary"
                    v-if="scope.row?.expired == 'no'"
                    size="small"
                    icon="edit"
                    :title="$t('btn.edit')"
                    @click.stop="handleUpdate(scope.row)">
                  </el-button>
                  <el-button type="primary" v-else size="small" disabled icon="edit" :title="$t('btn.edit')" @click.stop="handleUpdate(scope.row)">
                  </el-button>
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
              <el-form-item :label="proxy.$t('Member')" prop="chat_id">
                <el-select
                  v-model="form.chat_id"
                  filterable
                  remote
                  reserve-keyword
                  :placeholder="proxy.$t('system.Enter')"
                  remote-show-suffix
                  :remote-method="memberOption"
                  @change="memberOption"
                  :loading="loading"
                  style="width: 240px">
                  <el-option v-for="item in optionMember" :key="item.account_name" :label="item.account_name" :value="item.chat_id" />
                </el-select>
              </el-form-item>
              <!--
                            <el-form-item :label="proxy.$t('lables.chat_id')" prop="chat_id">
                <el-input v-model="form.chat_id" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              -->
              <el-form-item :label="proxy.$t('server.stream.name')" prop="note">
                <el-input v-model="form.note" type="textarea" :rows="2" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.total_chance')" prop="total_chance">
                <el-input-number v-model="form.total_chance" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
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
import { getBotLocationRecords, saveBotLocation, updateBotLocation } from '@/api/telegram/telegram'
import { getChancePoint, saveChancePoint, updateChancePoint } from '@/api/lucky/chance.service'
import { getTelegramMember } from '@/api/lucky/lucky.service'
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
  getChancePoint(data.queryParams).then((res) => {
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
        updateChancePoint(form.value).then((response) => {
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

        saveChancePoint(form.value).then((response) => {
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
