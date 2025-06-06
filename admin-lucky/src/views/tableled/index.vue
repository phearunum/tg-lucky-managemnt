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
            <el-table-column :label="proxy.$t('server.desk.id')" align="center" sortable key="id" width="100" prop="id" />
            <el-table-column :label="proxy.$t('server.desk.desk_no')" align="center" sortable key="desk_no" prop="desk_no" />
            <el-table-column :label="proxy.$t('server.desk.desk_streamkey')" align="center" sortable key="desk_streamkey" prop="desk_streamkey" />
            <el-table-column
              :label="proxy.$t('server.desk.desk_stream_server')"
              align="center"
              sortable
              key="desk_stream_server"
              prop="desk_stream_server"
              width="250">
              <template #default="scope">
                <span>{{ scope.row?.stream_server?.server_name }} : {{ scope.row?.stream_server?.rtmp }}</span>
              </template>
            </el-table-column>

            <el-table-column
              :label="proxy.$t('create on')"
              align="center"
              sortable
              prop="created_at"
              v-if="columns.showColumn('create_at')"
              width="200" />
            <el-table-column
              :label="proxy.$t('update on')"
              align="center"
              sortable
              prop="updated_at"
              v-if="columns.showColumn('updated_at')"
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
              <el-form-item :label="proxy.$t('server.desk.desk_no')" prop="desk_no">
                <el-input v-model="form.desk_no" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('server.desk.desk_streamkey')" prop="desk_streamkey">
                <el-input v-model="form.desk_streamkey" :placeholder="proxy.$t('system.Enter')" style="width: 250px" />
                <el-button type="warning" @click="GenerateStreamingKey"
                  >Genertae<el-icon><Connection /></el-icon
                ></el-button>
              </el-form-item>

              <el-form-item :label="proxy.$t('server.desk.desk_stream_server')" prop="desk_stream_server">
                <el-select v-model="form.desk_stream_server" placeholder="Select Server" style="width: 100%" @change="selectRole($event)">
                  <el-option v-for="item in ServerOption" :key="item.id" :label="item.server_name + ' : ' + item.rtmp" :value="item.id">
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
import { videoSettingList } from '@/api/system/videosnap'
import { deskSettingList, adddeskSetting, updatedeskSetting, deskSettingByID, deletedeskSetting } from '@/api/tables/desks'
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
const ServerOption = ref([])
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
const loadServerOption = () => {
  videoSettingList().then((data) => {
    ServerOption.value = data.data
  })
}
function GenerateStreamingKey() {
  var keySign = 'GsRi8CMEYQeyBQNf'
  const _streamKey = proxy.$md5(`${keySign + form.value.desk_no}`)
  form.value.desk_streamkey = _streamKey
}
const loadList = (_page = 1) => {
  loading.value = true
  deskSettingList(data.queryParams).then((res) => {
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
  deskSettingByID(id).then((response) => {
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
        updatedeskSetting(form.value).then((response) => {
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

        adddeskSetting(form.value).then((response) => {
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
      return deletedeskSetting(ids)
    })
    .then((response) => {
      if (response.statusCode == '200') {
        proxy.$modal.msgSuccess(proxy.$t('system.Success'))
        loadList()
      } else {
        proxy.$modal.msgError(proxy.$t('system.Failed'))
      }
    })
}

onMounted(() => {
  loadList(), loadServerOption()
})
</script>
