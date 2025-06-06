<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-col :lg="24" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" plain size="small" @click="handleAdd">
              <el-icon style="font-size: 15px"><UploadFilled /></el-icon> {{ $t('btn.upload') }}
            </el-button>
          </el-col>
          <el-col :span="1.5">
            <el-input
              v-model="queryParams.query"
              :placeholder="proxy.$t('btn.search')"
              clearable
              prefix-icon="el-icon-search"
              style="margin-bottom: 20px" />
          </el-col>
          <el-col :span="1.5">
            <el-button type="primary" icon="search" size="small" @click="getList(1)">{{ $t('btn.search') }}</el-button>
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
          <el-table-column :label="proxy.$t('company.card.id')" align="center" key="id" width="50" prop="id" />

          <el-table-column :label="proxy.$t('menu.bucketName')" align="center" key="fileName" prop="fileName">
            <template #default="scope">{{ scope.row?.bucket.bucketProject }}</template>
          </el-table-column>
          <el-table-column :label="proxy.$t('menu.fileName')" key="fileName" prop="fileName" />
          <el-table-column :label="proxy.$t('menu.path')" key="filePath" prop="filePath" />
          <el-table-column :label="proxy.$t('menu.originalFile')" key="originalFile" prop="originalFile" />

          <el-table-column :label="proxy.$t('menu.type')" key="mimeType" prop="mimeType" />
          <el-table-column :label="proxy.$t('menu.size')" key="fileSize" prop="fileSize" />

          <el-table-column :label="proxy.$t('createon')" align="center" prop="createdAt" v-if="columns.showColumn('createdAt')" width="200">
            <template #default="scope">
              {{ scope.row?.createdAt }}
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('menu.createBy')" align="center" key="username" prop="username" />
          <el-table-column :label="proxy.$t('common.status')" align="center">
            <template #default="scope">
              <el-tag v-if="scope.row?.status == 'success'" type="success">
                {{ proxy.$t('menu.success') }}
              </el-tag>
              <el-tag v-else type="danger">
                {{ proxy.$t('menu.failed') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="proxy.$t('action')" align="center" width="150" class-name="small-padding fixed-width">
            <template #default="scope">
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
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
          <el-row>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('common.name')" prop="fileName">
                <el-input v-model="form.fileName" :placeholder="proxy.$t('system.Enter')" />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('common.name')" prop="file">
                <input type="file" :placeholder="proxy.$t('system.Enter')" @change="handleFileChange" required />
              </el-form-item>
            </el-col>
            <el-col :lg="20">
              <el-form-item :label="proxy.$t('menu.bucketName')" prop="bucketId">
                <el-select v-model="form.bucketId" placeholder="status" style="width: 100%">
                  <el-option
                    v-for="item in bucketSelect"
                    :key="item.id"
                    :label="item.bucketProject + '-(' + item.bucketName + ')'"
                    :value="item.id">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <template #footer>
          <el-button text @click="open = false">{{ $t('btn.cancel') }}</el-button>
          <el-button type="primary" :loading="uploading" @click="submitForm">{{ $t('btn.submit') }}</el-button>
        </template>
      </el-dialog>
    </el-row>
  </div>
</template>

<script setup name="user">
import { getBucketById, getSelectBucket } from '@/api/bucket/bucket'
import { createupload, getListUpload, deleteUpload } from '@/api/bucket/upload'
import { getUserInfo } from '@/utils/auth'
const { companyCode } = JSON.parse(getUserInfo())

const { proxy } = getCurrentInstance()
const userList = ref([])
const bucketSelect = ref([])
const ids = ref([])
const uploading = ref(false)
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
    statusCode: 0,
    statusName: proxy.$t('common.no')
  },
  {
    statusCode: 1,
    statusName: proxy.$t('common.yes')
  }
]
function httpRequest(option) {
  this.fileList.push(option)
}
const handleFileChange = (event) => {
  form.value.file = event.target.files[0]
}
const data = reactive({
  queryParams: {
    page: 0, //num_page
    limit: 15, //page
    query: ''
  },
  form: {
    fileName: '',
    file: null,
    bucketId: undefined,
    username: companyCode
  },
  rules: {
    fileName: [{ required: true, min: 5, max: 20, message: `${proxy.$t('menu.min5Max20')}`, trigger: 'blur' }],
    bucketId: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }],
    file: [{ required: true, message: `${proxy.$t('menu.notEmpty')}`, trigger: 'blur' }]
  }
})

const { queryParams, rules, form } = toRefs(data)
const router = useRouter()
const getBucketSelect = () => {
  getSelectBucket().then((bucket) => {
    bucketSelect.value = bucket.data
  })
}
const getList = (_page = 1) => {
  loading.value = true

  getListUpload(data.queryParams).then((res) => {
    loading.value = false
    userList.value = res.data
    total.value = res.meta.total
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
    hour24: true
  }
  const formattedTime = date.toLocaleTimeString('en-US', timeOptions)

  return `${formattedDate}, ${formattedTime}`
}
const handleChangeSize = (val) => {
  current.value = val
  proxy.$nextTick(() => {
    getList(val)
  })
}
const handlePageChangeSize = (val) => {
  pageSize.value = val
  proxy.$nextTick(() => {
    getList()
  })
}
const handleAdd = () => {
  open.value = true
  title.value = proxy.$t('system.Add New')
}
function handleUpdate(row) {
  const id = row.id || ids.value
  getBucketById(id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = 'Update Record'
  })
}
const submitForm = () => {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      uploading.value = true
      if (form.value.id !== undefined && form.value.id > 0) {
        createupload(form.value).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            uploading.value = false
            getList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
            uploading.value = false
          }
        })
      } else {
        const formData = new FormData()
        formData.append('fileName', form.value.fileName)
        formData.append('file', form.value.file)
        formData.append('bucketId', form.value.bucketId)
        formData.append('username', companyCode)

        createupload(formData).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
            uploading.value = false
            getList()
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
            uploading.value = false
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
      return deleteUpload(ids)
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

getList()
getBucketSelect()
</script>
