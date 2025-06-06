<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-col :lg="24" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" plain icon="plus" @click="handleAdd">{{ $t('btn.add') }}</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-input v-model="queryParams.query" placeholder="Search" clearable prefix-icon="el-icon-search" style="margin-bottom: 20px" />
          </el-col>
          <el-col :span="1.5">
            <el-button type="primary" icon="search">{{ $t('btn.search') }}</el-button>
          </el-col>
        </el-row>
        <el-table
          style="height: 75vh"
          v-loading="loading"
          :data="userList"
          :header-cell-style="{
            background: 'var(--el-fill-color-light)',
            color: 'var(--el-text-color-primary)'
          }">
          <el-table-column :label="proxy.$t('layout.identity')" align="center" key="id" width="50" prop="id" />
          <el-table-column :label="proxy.$t('common.name')" align="center" key="bucketProject" prop="bucketProject" />
          <el-table-column :label="proxy.$t('menu.bucketName')" align="center" key="bucketName" prop="bucketName" />
          <el-table-column :label="proxy.$t('menu.bucketAddress')" align="center" key="bucketAddress" prop="bucketAddress" />
          <el-table-column :label="proxy.$t('menu.folder')" key="folders" prop="folders" />
          <el-table-column :label="proxy.$t('menu.bucketAccesskey')" align="center">
            <template #default="scope">
              <div>
                <el-tooltip class="box-item" effect="dark" :content="scope.row?.bucketAccesskey" placement="top-start">
                  <el-icon><DocumentCopy /></el-icon>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('menu.bucketScretekey')" align="center">
            <template #default="scope">
              <div>
                <el-tooltip class="box-item" effect="dark" :content="scope.row?.bucketScretekey" placement="top-start">
                  <el-icon><DocumentCopy /></el-icon>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('menu.bucketRegin')" align="center" key="bucketRegin" prop="bucketRegin" />
          <el-table-column :label="proxy.$t('common.status')" align="center">
            <template #default="scope">
              <el-switch
                :value="scope.row?.bucketStatus"
                inline-prompt
                :active-text="proxy.$t('common.yes')"
                :inactive-text="proxy.$t('common.no')" />
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('createon')" align="center" prop="createdAt" v-if="columns.showColumn('createdAt')" width="200">
            <template #default="scope">
              {{ formatDate(scope.row?.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('action')" align="center" width="150" class-name="small-padding fixed-width">
            <template #default="scope">
              <el-button
                size="small"
                text
                icon="edit"
                :title="$t('btn.edit')"
                @click.stop="handleUpdate(scope.row)"
                v-hasPermi="['upload:bucket:edit']">
              </el-button>
              <el-button
                size="small"
                text
                icon="delete"
                :title="$t('btn.delete')"
                @click.stop="handleDelete(scope.row)"
                v-hasPermi="['upload:bucket:delete']">
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.limit"
          :page-sizes="[100, 200, 300, 400]"
          :small="true"
          :disabled="disabled"
          background="var(--el-fill-color-light)"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleChangeSize"
          @current-change="handlePageChangeSize" />
      </el-col>

      <el-dialog :title="title" v-model="open" width="600px" append-to-body :close-on-click-modal="false">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="150px">
          <el-row>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('common.name')" prop="bucketProject">
                <el-input v-model="form.bucketProject" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketName')" prop="bucketName">
                <el-input v-model="form.bucketName" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketAddress')" prop="bucketAddress">
                <el-input v-model="form.bucketAddress" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.folder') + '( / )'" prop="folders">
                <el-input v-model="form.folders" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketPort')" prop="bucketPort">
                <el-input v-model="form.bucketPort" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketScretekey')" prop="bucketScretekey">
                <el-input v-model="form.bucketScretekey" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketAccesskey')" prop="bucketAccesskey">
                <el-input v-model="form.bucketAccesskey" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketRegin')" prop="bucketRegin">
                <el-input v-model="form.bucketRegin" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="8">
              <el-form-item :label="proxy.$t('common.status')" prop="bucketStatus" style="width: 300px">
                <el-select v-model="form.bucketStatus" placeholder="status" style="width: 200px">
                  <el-option v-for="item in activeOption" :key="item.statusCode" :label="item.statusName" :value="item.statusCode"> </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.staticFileName')" prop="isStaticName">
                <el-checkbox v-model="form.isStaticName" label="Use Satic Name " />
                <el-input v-model="form.staticFileName" :placeholder="proxy.$t('system.Enter')" />
                <el-text class="mx-1" type="warning" size="small">* Use upload with the static file name </el-text>
                <el-text class="mx-1" type="warning" size="small">* If not just leave it blank </el-text>
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
import { getBucket, createBucket, updateBucket, getBucketById } from '@/api/bucket/bucket'

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
const activeOption = ref()
activeOption.value = [
  {
    statusCode: false,
    statusName: proxy.$t('common.no')
  },
  {
    statusCode: true,
    statusName: proxy.$t('common.yes')
  }
]

const data = reactive({
  queryParams: {
    page: 0, //num_page
    limit: 15, //page
    query: undefined
  },
  form: {
    folders: '/'
  },
  rules: {
    bucketName: [{ required: true, min: 5, max: 20, message: `${proxy.$t('menu.min5Max20')}`, trigger: 'blur' }],
    bucketProject: [{ required: true, min: 5, max: 20, message: `${proxy.$t('menu.min5Max20')}`, trigger: 'blur' }],
    bucketAddress: [{ required: true, min: 5, max: 20, message: `${proxy.$t('menu.min5Max20')}`, trigger: 'blur' }],
    bucketAccesskey: [{ required: true, min: 5, max: 50, message: `${proxy.$t('menu.notEmpty')} [min:5 | max:50]`, trigger: 'blur' }],
    bucketScretekey: [{ required: true, min: 5, max: 50, message: `${proxy.$t('menu.notEmpty')} [min:5 | max:50]`, trigger: 'blur' }],
    bucketPort: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }],
    bucketPort: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }],
    folders: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }],
    bucketStatus: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }],
    bucketRegin: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }]
  }
})

const { queryParams, rules, form } = toRefs(data)

function getList(_page = 1) {
  loading.value = true
  getBucket(proxy.addDateRange(data.queryParams)).then((res) => {
    loading.value = false
    userList.value = res.data
    total.value = res.meta.total
  })
}
function handleChangeSize(val) {
  current.value = val
  proxy.$nextTick(() => {
    // getList(val)
  })
}
function handlePageChangeSize(val) {
  pageSize.value = val
  proxy.$nextTick(() => {
    //getList()
  })
}
function handleAdd() {
  open.value = true
  title.value = proxy.$t('system.Add New')
  form.value = {
    daylimit: 0,
    folders: '/'
  }
}
function handleUpdate(row) {
  const id = row.id || ids.value
  getBucketById(id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = 'Update Record'
  })
}
function submitForm() {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.id !== undefined && form.value.id > 0) {
        updateBucket(form.value).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            getList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
        })
      } else {
        createBucket(form.value).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            getList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
        })
      }
    }
  })
}

function handleDelete(row) {
  const ids = row.id || ids.value
  proxy
    .$confirm(proxy.$t('system.Do you want to delete') + ` = ${ids} ?`, proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('yes'),
      cancelButtonText: proxy.$t('no'),
      type: 'warning'
    })
    .then(() => {
      return deleteCardsetting(ids)
    })
    .then((response) => {
      if (response.status == 200) {
        proxy.$modal.msgSuccess(proxy.$t('system.Success'))
        getList()
      } else {
        proxy.$modal.msgError(proxy.$t('system.Failed'))
      }
    })
}
const formatDate = (dateString) => {
  const date = new Date(dateString)
  const options = {
    month: 'long',
    day: 'numeric',
    year: 'numeric'
  }
  const formattedDate = date.toLocaleDateString('en-US', options)

  const timeOptions = {
    hour: 'numeric',
    minute: 'numeric',
    hour12: true
  }
  const formattedTime = date.toLocaleTimeString('en-US', timeOptions)

  return `${formattedDate}, ${formattedTime}`
}

getList()
</script>
