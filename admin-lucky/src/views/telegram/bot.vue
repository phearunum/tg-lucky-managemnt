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
              <el-table-column :label="proxy.$t('Bot Name')" align="left" width="180" sortable key="name" prop="name" />
              <el-table-column :label="proxy.$t('Chat ID')" align="left" sortable key="token_bot" width="250" prop="token_bot" />
              <el-table-column :label="proxy.$t('lables.allow')" align="center" key="allow" prop="allow">
                <template #default="scope">
                  {{ scope.row?.allow }}
                </template>
              </el-table-column>

              <el-table-column :label="proxy.$t('menu.location')" align="right" key="lat" prop="lat">
                <template #default="scope">
                  <a :href="'https://www.google.com/maps?q=' + scope.row.lat + ',' + scope.row.long" target="_blank">
                    {{ scope.row.lat }}, {{ scope.row.long }} <el-icon><Position /></el-icon>
                  </a>
                </template>
              </el-table-column>

              <el-table-column :label="proxy.$t('created_at')" align="center" key="created_at" prop="created_at">
                <template #default="scope">
                  {{ FormatDate24(scope?.row?.created_at) }}
                </template>
              </el-table-column>
              <el-table-column :label="proxy.$t('updated_at')" align="center" key="updated_at" prop="updated_at">
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
        <el-form ref="formRef" :model="form" :rules="rules" label-width="150px">
          <el-row>
            <el-col :lg="24">
              <el-form-item :label="proxy.$t('server.stream.name')" prop="name">
                <el-input v-model="form.name" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.token')" prop="token_bot">
                <el-input v-model="form.token_bot" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.lat')" prop="lat">
                <el-input v-model="form.lat" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.long')" prop="long">
                <el-input v-model="form.long" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
              <el-form-item :label="proxy.$t('lables.allow')" prop="allow">
                <el-input-number v-model="form.allow" :placeholder="proxy.$t('system.Enter')" controls-position="right" />
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
  getBotLocationRecords(data.queryParams).then((res) => {
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
const submitForm = () => {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.id !== undefined && form.value.id > 0) {
        //    form.value.type = 'edit'
        updateBotLocation(form.value).then((response) => {
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

        saveBotLocation(form.value).then((response) => {
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
