<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!--部门数据-->
      <el-col :span="4" :xs="24">
        <div class="head-container">
          <el-input v-model="queryParams.q" placeholder="Search" clearable prefix-icon="el-icon-search" style="margin-bottom: 20px" />
        </div>
        <div class="head-container">
          <el-tree
            :data="deptOptions"
            :props="{ label: 'label', children: 'children' }"
            :expand-on-click-node="false"
            :filter-node-method="filterNode"
            ref="deptTreeRef"
            default-expand-all
            @node-click="handleNodeClick" />
        </div>
      </el-col>

      <el-col :lg="20" :xm="24">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" plain icon="Plus" @click="handleAdd" v-hasPermi="['system:user:add']">
              {{ $t('btn.add') }}
            </el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="success" plain icon="Edit" :disabled="single" @click="handleUpdate" v-hasPermi="['system:user:edit']">
              {{ $t('btn.edit') }}
            </el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="handleDelete" v-hasPermi="['system:user:delete']">
              {{ $t('btn.delete') }}
            </el-button>
          </el-col>
        </el-row>

        <el-table
          v-loading="loading"
          :data="userList"
          @selection-change="handleSelectionChange"
          :header-cell-style="{
            background: 'var(--el-fill-color-light)',
            color: 'var(--el-text-color-primary)'
          }">
          <el-table-column type="selection" width="50" align="center" />
          <el-table-column label="ID" align="center" key="id" prop="id" v-if="columns.showColumn('id')" />
          <el-table-column
            :label="proxy.$t('user.Frist Name')"
            align="center"
            key="fname"
            prop="fname"
            v-if="columns.showColumn('fname')"
            :show-overflow-tooltip="true" />
          <el-table-column
            :label="proxy.$t('user.Last Name')"
            align="center"
            key="lname"
            prop="lname"
            v-if="columns.showColumn('lname')"
            :show-overflow-tooltip="true" />
          <el-table-column
            :label="proxy.$t('user.username')"
            align="center"
            key="username"
            prop="username"
            v-if="columns.showColumn('username')"
            :show-overflow-tooltip="true" />

          <el-table-column
            label="TG Group"
            align="left"
            key="tg_group"
            prop="tg_group"
            width="350"
            v-if="columns.showColumn('tg_group')"
            :show-overflow-tooltip="true">
            <template #default="scope">
              <span v-if="scope.row.tg_group">[ {{ scope.row.tg_group }} ]</span>
              <span v-else>Unassigned</span>
            </template>
          </el-table-column>

          <el-table-column :label="proxy.$t('create on')" align="center" prop="createdAt" v-if="columns.showColumn('createdAt')" width="200" />
          <el-table-column prop="sex" label="Sex" align="center" v-if="columns.showColumn('sex')">
            <template #default="scope">
              <dict-tag :options="sexOptions" :value="scope.row.sex" />
            </template>
          </el-table-column>

          <el-table-column prop="email" label="Email" align="center" v-if="columns.showColumn('email')" />
          <el-table-column prop="loginDate" label="Last Login" align="center" v-if="columns.showColumn('loginDate')" />
          <el-table-column label="Action" align="center" width="200" class-name="small-padding fixed-width">
            <template #default="scope">
              <el-button v-if="scope.row.id !== 0" text icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['system:user:edit']"> </el-button>
              <el-button v-if="scope.row.id !== 0" text icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['system:user:delete']">
              </el-button>
              <el-button text icon="Key" @click="handleResetPwd(scope.row)" v-hasPermi="['system:user:edit']"> </el-button>
              <el-button @click="handleAssign(scope.row)" v-hasPermi="['system:user:edit']">
                <el-icon><Promotion /></el-icon> Assgin
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <pagination :total="total" v-model:page="queryParams.num_page" v-model:limit="queryParams.limit" @pagination="getList" />
      </el-col>
    </el-row>

    <!-- Add or modify user profile dialog -->
    <el-dialog :title="title" v-model="open" append-to-body>
      <el-form :model="form" :rules="rules" ref="userRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :lg="12">
            <el-form-item :label="$t('user.Frist Name')" prop="fname">
              <el-input v-model="form.fname" :placeholder="$t('user.Frist Name')" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item :label="$t('user.Last Name')" prop="lname">
              <el-input v-model="form.lname" :placeholder="$t('user.Last Name')" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item label="Gender">
              <el-radio-group v-model="form.sex">
                <el-radio v-for="sex in sexOptions" :key="sex.dictValue" :label="sex.dictValue">{{ sex.dictLabel }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>

          <el-col :lg="12">
            <el-form-item label="Role" prop="roleId">
              <el-select v-model="form.roleId" placeholder="Role" style="width: 100%" @change="selectRole($event)">
                <el-option v-for="item in roleOptions" :key="item.id" :label="item.role_name" :value="item.id"> </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item :label="$t('user.username')" prop="username">
              <el-input :disabled="form.id != undefined" v-model="form.username" :placeholder="$t('user.username')" />
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.id == undefined">
            <el-form-item :label="$t('user.password')" prop="password">
              <el-input v-model="form.password" :placeholder="$t('user.password')" type="password" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item label="TG Group" prop="tg_group">
              <el-input v-model="form.tg_group" type="textarea" placeholder="-3898595955,98598585,...." />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button text @click="cancel">{{ $t('btn.cancel') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('btn.submit') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="user">
import { getToken } from '@/utils/auth'
import { changeUserStatus, assignTelegramGroup, listUser, resetUserPwd, delUser, getUser, updateUser, addUser, exportUser } from '@/api/system/user'
import { getListSelect } from '@/api/system/role'
const { proxy } = getCurrentInstance()
import { Promotion } from '@element-plus/icons-vue'
const statusOptions = ref([])
const sexOptions = ref([])
/*
proxy.getDicts('sys_normal_disable').then((response) => {
  statusOptions.value = response.data
})
*/
sexOptions.value = [
  {
    dictCode: 'M',
    dictLabel: 'Male',
    dictValue: 'M'
  },
  {
    dictCode: 'F',
    dictLabel: 'Female',
    dictValue: 'F'
  }
]
const userList = ref([])
const open = ref(false)
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)
const title = ref('')
const dateRange = ref([])
const country_code = ref('')
const deptOptions = ref(undefined)
const initPassword = ref(undefined)
const postOptions = ref([])
const roleOptions = ref([])

const columns = ref([
  { key: 0, label: `用户编号`, visible: true, prop: 'id' },
  { key: 1, label: `用户名称`, visible: true, prop: 'username' },
  { key: 2, label: `用户昵称`, visible: true, prop: 'fname' },
  { key: 3, label: `部门`, visible: true, prop: 'country_code' },
  { key: 4, label: `手机号码`, visible: true, prop: 'phone' },
  { key: 5, label: `状态`, visible: true, prop: 'status' },
  { key: 6, label: `创建时间`, visible: true, prop: 'createdAt' },
  { key: 7, label: `性别`, visible: false, prop: 'sex' },
  { key: 8, label: `头像`, visible: false, prop: 'avatar' },
  { key: 9, label: `邮箱`, visible: false, prop: 'email' },
  { key: 10, label: `最后登录时间`, visible: false, prop: 'loginDate' }
])

const data = reactive({
  form: {},
  queryParams: {
    page: 0, //num_page
    limit: 10, //page
    username: undefined,
    phone: undefined,
    status: undefined,
    deptId: undefined,
    q: undefined
  },
  rules: {
    username: [
      { required: true, message: '用户名称不能为空', trigger: 'blur' },
      {
        min: 2,
        max: 20,
        message: '用户名称长度必须介于 2 和 20 之间',
        trigger: 'blur'
      }
    ],
    fname: [{ required: true, message: '用户昵称不能为空', trigger: 'blur' }],
    password: [
      { required: true, message: '用户密码不能为空', trigger: 'blur' },
      {
        min: 5,
        max: 20,
        message: '用户密码长度必须介于 5 和 20 之间',
        trigger: 'blur'
      }
    ],
    email: [
      {
        required: true,
        type: 'email',
        message: '请输入正确的邮箱地址',
        trigger: ['blur', 'change']
      }
    ]
  }
})

const { queryParams, form, rules } = toRefs(data)
/*
proxy.getConfigKey('sys.user.initPassword').then((response) => {
  initPassword.value = response.data
})
*/

/** 通过条件过滤节点  */
const filterNode = (value, data) => {
  if (!value) return true
  return data.label.indexOf(value) !== -1
}
/** 根据名称筛选部门树 */

watch(country_code, (val) => {
  proxy.$refs['deptTreeRef'].filter(val)
})

/** 查询部门下拉树结构 */
function getTreeselect() {
  getListSelect().then((response) => {
    deptOptions.value = response.data
  })
  getListSelect().then((response) => {
    console.log(response.data)
    roleOptions.value = response.data
  })
}

function getList() {
  loading.value = true
  listUser(proxy.addDateRange(queryParams.value, dateRange.value)).then((res) => {
    loading.value = false
    userList.value = res.data
    total.value = res.meta.total
  })
}

function handleNodeClick(data) {
  queryParams.value.deptId = data.id
  handleQuery()
}

function handleQuery() {
  queryParams.value.num_page = 1
  getList()
}
function resetQuery() {
  dateRange.value = []
  proxy.resetForm('queryRef')
  queryParams.value.deptId = undefined
  handleQuery()
}

function handleDelete(row) {
  const delete_id = row.id || ids.value
  proxy.$modal
    .confirm('Are you sure to delete the user ID as"' + delete_id + '"data item?', proxy.$t('btn.yes'), proxy.$t('btn.cancel'))
    .then(function () {
      return delUser(delete_id)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('Success')
    })
    .catch(() => {})
}
/** 导出按钮操作 */
function handleExport() {
  proxy.$modal
    .confirm('是否确认导出所有用户数据项?', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    .then(async () => {
      await exportUser(queryParams.value)
    })
}

function handleResetPwd(row) {
  proxy
    .$prompt(proxy.$t('user.changePwd') + ' ' + row.username, proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('common.ok'),
      cancelButtonText: proxy.$t('common.cancel'),
      closeOnClickModal: false,
      inputPattern: /^.{5,20}$/,
      inputErrorMessage: proxy.$t('system.Password Required')
    })
    .then(({ value }) => {
      console.log(row)
      resetUserPwd({
        id: row.id,
        password: value,
        username: row.username
      }).then((response) => {
        proxy.$modal.msgSuccess(`${proxy.$t('system.Success')}`)
      })
    })
    .catch(() => {})
}
function handleAssign(row) {
  proxy
    .$prompt('Assign Telegram Group ' + row.username + ' [9485850696,-30949494...]', proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('common.ok'),
      cancelButtonText: proxy.$t('common.cancel'),
      closeOnClickModal: false,
      inputPattern: /^.{5,100}$/,
      inputErrorMessage: proxy.$t('system.Password Required')
    })
    .then(({ value }) => {
      console.log(row)
      assignTelegramGroup({
        id: row.id,
        tg_group: value
      }).then((response) => {
        proxy.$modal.msgSuccess(`${proxy.$t('system.Success')}`)
      })
    })
    .catch(() => {})
}
/** 选择条数  */
function handleSelectionChange(selection) {
  ids.value = selection.map((item) => item.id)
  single.value = selection.length != 1
  multiple.value = !selection.length
}

function reset() {
  form.value = {
    id: undefined,
    deptId: undefined,
    username: undefined,
    fname: undefined,
    password: undefined,
    phone: undefined,
    email: undefined,
    sex: '1',
    role_id: [],
    tg_group: undefined
  }
  //proxy.resetForm('userRef')
}

function cancel() {
  open.value = false
  reset()
}
function handleAdd() {
  reset()
  title.value = proxy.$t('user.create_user')
  open.value = true
}
function handleUpdate(row) {
  reset()

  const id = row.id || ids.value
  getUser(id).then((response) => {
    var data = response.data
    form.value = {
      id: data.id,
      username: data.username,
      fname: data.fname,
      lname: data.lname,
      password: '',
      phone: data.phone,
      email: data.email,
      sex: data.sex,
      roleId: data.roleId,
      tg_group: data.tg_group
    }
    //roleOptions.value = response.data.roles
    //postOptions.value = response.data.posts
    open.value = true
    title.value = proxy.$t('user.update_user')
    form.password = ''
  })
}
/** 提交按钮 */
function submitForm() {
  proxy.$refs['userRef'].validate((valid) => {
    if (valid) {
      console.log('====' + form.value.id)

      if (form.value.id !== undefined) {
        updateUser(form.value).then((response) => {
          proxy.$modal.msgSuccess(proxy.$t('system.Success'))
          open.value = false
          getList()
        })
      } else {
        addUser(form.value).then((response) => {
          if (response.status == 200) {
            proxy.$modal.msgSuccess(proxy.$t('system.Success'))
          } else {
            proxy.$modal.msgError(proxy.$t('system.Failed'))
          }
          open.value = false
          getList()
        })
      }
    }
  })
}
/**
 * 解决编辑时角色选中不了问题
 */
function selectRole(e) {
  console.log(e, JSON.stringify(this.form))
  proxy.$forceUpdate()
}
getTreeselect()
getList()
</script>
