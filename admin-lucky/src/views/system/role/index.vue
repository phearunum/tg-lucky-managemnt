<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryForm" v-show="showSearch" :inline="true">
      <el-form-item label="Name" prop="role_name">
        <el-input v-model="queryParams.role_name" placeholder="Enter role name" clearable size="small" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="search" size="small" @click="handleQuery">{{ $t('btn.search') }}</el-button>
        <el-button icon="refresh" size="small" @click="resetQuery">{{ $t('btn.reset') }}</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="plus" size="small" @click="handleAdd" v-hasPermi="['system:role:add']">{{ $t('btn.add') }}</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table
      v-loading="loading"
      :data="roleList"
      highlight-current-row
      @selection-change="handleSelectionChange"
      :header-cell-style="{
        background: 'var(--el-fill-color-light)',
        color: 'var(--el-text-color-primary)'
      }">
      <el-table-column label="ID" prop="id" width="80" />
      <el-table-column label="Name" prop="role_name" />
      <el-table-column label="Role Key" prop="role_key" />
      <el-table-column label="Staus" align="center" width="90">
        <template #default="scope">
          <div v-if="scope.row.role_status == 1"><el-tag class="ml-2" type="success">Active</el-tag></div>
          <div v-if="scope.row.role_status == 0"><el-tag class="ml-2" type="warning">Inactive</el-tag></div>
        </template>
      </el-table-column>
      <el-table-column :label="proxy.$t('create on')" align="center" prop="createdAt" width="200" />
      <el-table-column label="Action" align="center" width="200">
        <template #default="scope">
          <div v-if="scope.row.role_name != 'admin'">
            <el-button size="small" text icon="edit" :title="$t('btn.edit')" @click.stop="handleUpdate(scope.row)"> </el-button>
            <el-button
              size="small"
              text
              icon="delete"
              :title="$t('btn.delete')"
              @click.stop="handleDelete(scope.row)"
              v-hasPermi="['system:role:remove']">
            </el-button>
            <!--  v-hasPermi="['system:role:edit', 'system:role:authorize', 'system:roleusers:list']" -->
            <el-dropdown size="small" @command="(command) => handleCommand(command, scope.row)">
              <span class="el-dropdown-link">
                {{ $t('btn.more') }}
                <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="handleDataScope" icon="circle-check">{{ $t('menu.menuPermi') }}</el-dropdown-item>
                  <!-- <el-dropdown-item command="handleAuthUser" icon="user">{{ $t('menu.assignUsers') }}</el-dropdown-item> -->
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <!-- 添加或修改角色配置对话框 -->
    <el-dialog :title="title" v-model="open" width="600px" append-to-body>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row>
          <el-col :lg="24">
            <el-form-item label="Role Name" prop="role_name">
              <el-input v-model="form.role_name" placeholder="Enter.." />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item label="Role Key" prop="role_key" key="role_key">
              <el-input v-model="form.role_key" placeholder="Role Key" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item label="Status">
              <el-radio-group v-model="form.role_status">
                <el-radio v-for="dict in statusOptions" :key="dict.status_value" :label="dict.status_value">{{ dict.status_name }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button text @click="cancel">{{ $t('btn.cancel') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('btn.submit') }}</el-button>
      </template>
    </el-dialog>
    <el-dialog title="Assignment of role permissions" v-model="showRoleScope" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="menu search">
          <el-input placeholder="Please enter keywords to filter" v-model="searchText"></el-input>
        </el-form-item>
        <el-form-item label="permission character">
          {{ form.roleKey }}
        </el-form-item>
        <el-form-item label="Menu permissions">
          <el-checkbox v-model="menuExpand" @change="handleCheckedTreeExpand($event, 'menu')">expand/collapse</el-checkbox>
          <el-checkbox v-model="menuNodeAll" @change="handleCheckedTreeNodeAll($event, 'menu')">Select all/Not select all</el-checkbox>
          <el-checkbox v-model="form.menuCheckStrictly" @change="handleCheckedTreeConnect($event, 'menu')">father-son linkage</el-checkbox>
          <el-tree
            class="tree-border"
            :data="menuOptions"
            show-checkbox
            ref="menuRef"
            node-key="id"
            :check-strictly="!form.menuCheckStrictly"
            empty-text="加载中，请稍后"
            :filter-node-method="menuFilterNode"
            :props="defaultProps"></el-tree>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button text @click="cancel">{{ $t('btn.cancel') }}</el-button>
        <el-button type="primary" @click="submitDataScope">{{ $t('btn.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup name="role">
import { listRole, getRole, delRole, addRole, updateRole, exportRole, dataScope, changeRoleStatus } from '@/api/system/role'
import { roleMenuTreeselect } from '@/api/system/menu'
const { proxy } = getCurrentInstance()

const loading = ref(true)
// 选中数组
const ids = ref([])
// 非单个禁用
const single = ref(true)
// 非多个禁用
const multiple = ref(true)
// 显示搜索条件
const showSearch = ref(true)
// 总条数
const total = ref(0)
// 角色表格数据
const roleList = ref([])
// 弹出层标题
const title = ref('')
// 是否显示弹出层
const open = ref(false)
const menuExpand = ref(true)
const menuNodeAll = ref(false)
const deptExpand = ref(true)
const deptNodeAll = ref(false)
// 日期范围
const dateRange = ref([])
// 状态数据字典
const statusOptions = ref([])
// 是否显示下拉菜单分配
const showRoleScope = ref(false)
// 数据范围选项
const dataScopeOptions = ref([
  {
    dictValue: '1',
    dictLabel: 'all'
  },
  {
    dictValue: '2',
    dictLabel: 'custom'
  },
  {
    dictValue: '3',
    dictLabel: 'This department'
  },
  {
    dictValue: '4',
    dictLabel: 'This department and below'
  },
  {
    dictValue: '5',
    dictLabel: 'Only me'
  }
])
// 菜单列表
const menuOptions = ref([])
// 部门列表
const deptOptions = ref([])

const queryParams = reactive({
  page: 0, //num_page
  limit: 10, //page
  q: '',
  role_name: undefined,
  roleKey: undefined,
  status: undefined
})
const searchText = ref('')

const state = reactive({
  form: {},
  rules: {
    role_name: [{ required: true, message: 'Role name cannot be empty', trigger: 'blur' }]
  },
  defaultProps: {
    children: 'children',
    label: 'label'
  }
})
const menuRef = ref()
const deptRef = ref()
const formRef = ref()
const { form, rules, defaultProps } = toRefs(state)

watch(searchText, (val) => {
  proxy.$refs.menuRef.filter(val)
})

/** 查询角色列表 */
function getList() {
  loading.value = true

  listRole(proxy.addDateRange(queryParams, dateRange.value)).then((response) => {
    roleList.value = response.data
    total.value = response.data.total
    loading.value = false
    console.log(roleList.value)
  })
}

function getMenuAllCheckedKeys() {
  // 目前被选中的菜单节点
  const checkedKeys = proxy.$refs.menuRef.getCheckedKeys()
  // 半选中的菜单节点
  const halfCheckedKeys = proxy.$refs.menuRef.getHalfCheckedKeys()
  checkedKeys.unshift.apply(checkedKeys, halfCheckedKeys)
  return checkedKeys
}
// 所有部门节点数据
function getDeptAllCheckedKeys() {
  // 目前被选中的部门节点
  const checkedKeys = proxy.$refs.deptRef.getCheckedKeys()
  // 半选中的部门节点
  const halfCheckedKeys = proxy.$refs.deptRef.getHalfCheckedKeys()
  checkedKeys.unshift.apply(checkedKeys, halfCheckedKeys)
  return checkedKeys
}

function handleStatusChange(row) {
  const text = row.status === '0' ? '启用' : '停用'

  proxy
    .$confirm('确认要"' + text + '""' + row.role_name + '"角色吗?', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    .then(function () {
      return changeRoleStatus(row.id, row.status)
    })
    .then(() => {
      proxy.$modal.msgSuccess(text + '成功')
    })
    .catch(function () {
      row.status = row.status === '0' ? '1' : '0'
    })
}
// 取消按钮
function cancel() {
  open.value = false
  showRoleScope.value = false
  reset()
}
// 表单重置
function reset() {
  form.value = {
    id: undefined,
    role_name: undefined,
    role_status: '0'
  }
  proxy.resetForm('form')
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  dateRange.value = []
  proxy.resetForm('queryForm')
  handleQuery()
}
// 多选框选中数据
function handleSelectionChange(selection) {
  ids.value = selection.map((item) => item.id)
  single.value = selection.length != 1
  multiple.value = !selection.length
}
// 更多操作触发
function handleCommand(command, row) {
  switch (command) {
    case 'handleDataScope':
      handleDataScope(row)
      break
    case 'handleAuthUser':
      handleAuthUser(row)
      break
    default:
      break
  }
}
// 树权限（展开/折叠）
function handleCheckedTreeExpand(value, type) {
  if (type == 'menu') {
    const treeList = menuOptions.value
    for (let i = 0; i < treeList.length; i++) {
      proxy.$refs.menuRef.store.nodesMap[treeList[i].id].expanded = value
    }
  } else if (type == 'dept') {
    const treeList = deptOptions.value
    for (let i = 0; i < treeList.length; i++) {
      proxy.$refs.deptRef.store.nodesMap[treeList[i].id].expanded = value
    }
  }
}

// 树权限（全选/全不选）
function handleCheckedTreeNodeAll(value, type) {
  if (type == 'menu') {
    proxy.$refs.menuRef.setCheckedNodes(value ? menuOptions.value : [])
  } else if (type == 'dept') {
    proxy.$refs.deptRef.setCheckedNodes(value ? deptOptions.value : [])
  }
}

// 树权限（父子联动）
function handleCheckedTreeConnect(value, type) {
  if (type == 'menu') {
    form.value.menuCheckStrictly = !!value
  } else if (type == 'dept') {
    form.value.deptCheckStrictly = !!value
  }
}
function getRoleMenuTreeselect(roleId) {
  return roleMenuTreeselect(roleId).then((response) => {
    menuOptions.value = response.data.menu
    return response
  })
}
// 菜单筛选
function menuFilterNode(value, data) {
  if (!value) return true
  return data.label.indexOf(value) !== -1
}

/** 新增按钮操作 */
function handleAdd() {
  reset()
  open.value = true
  title.value = proxy.$t('user.create_user')
  showRoleScope.value = false
}

function handleUpdate(row) {
  reset()
  const id = row.id || ids.value
  //const roleDeptTreeselect = 0
  getRole(id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = 'Update Record'
  })
}
/** 选择角色权限范围触发 */
function dataScopeSelectChange(value) {
  if (value !== '2') {
    proxy.$refs.deptRef.setCheckedKeys([])
  }
}
// 数据权限
function dataScopeFormat(row, column) {
  return proxy.selectDictLabel(dataScopeOptions.value, row.dataScope)
}
/** 分配角色权限按钮操作 */
// 新增 和上面代码基本相同
function handleDataScope(row) {
  if (row.id == 1) {
    showRoleScope.value = false
    return
  }
  reset()
  showRoleScope.value = true
  const id = row.id || ids.value
  const roleMenu = getRoleMenuTreeselect(id)

  roleMenu.then((res) => {
    const checkedKeys = res.data.checkedKeys
    checkedKeys.forEach((v) => {
      nextTick(() => {
        proxy.$refs.menuRef.setChecked(v, true, false)
      })
    })
  })
  form.value = {
    id: row.id,
    role_name: row.role_name,
    roleKey: row.roleKey,
    menuCheckStrictly: row.menuCheckStrictly
  }
}
const router = useRouter()
/** 分配用户操作 */
function handleAuthUser(row) {
  const id = row.id
  router.push({ path: '/system/roleusers', query: { id } })
}
/** 提交按钮 */
function submitForm() {
  proxy.$refs['formRef'].validate((valid) => {
    if (valid) {
      if (form.value.id !== undefined && form.value.id > 0) {
        form.value.type = 'edit'
        updateRole(form.value).then((response) => {
          proxy.$modal.msgSuccess('Success')
          open.value = false
          getList()
        })
      } else {
        form.value.type = 'add'
        addRole(form.value).then((response) => {
          open.value = false
          if (response.statusCode == 200) {
            proxy.$modal.msgSuccess(response.msg)
            getList()
          } else {
            proxy.$modal.msgError(response.msg)
          }
        })
      }
    }
  })
}

/** 提交按钮（菜单数据权限） */
function submitDataScope() {
  if (form.value.id != undefined) {
    form.value.menuIds = getMenuAllCheckedKeys()
    dataScope(form.value).then((response) => {
      proxy.$modal.msgSuccess('修改成功')
      getList()
      cancel()
    })
  } else {
    proxy.$modal.msgError('请选择角色')
  }
}

function handleDelete(row) {
  const ids = row.id || ids.value
  proxy
    .$confirm(proxy.$t('system.Do you want to delete') + ` = ${ids} ?`, proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('yes'),
      cancelButtonText: proxy.$t('no'),
      type: 'warning'
    })
    .then(function () {
      return delRole(ids)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
}

/** 导出按钮操作 */
function handleExport() {
  proxy
    .$confirm('是否确认导出所有角色数据项?', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    .then(function () {
      return exportRole(queryParams)
    })
    .then((response) => {
      proxy.download(response.data.path)
    })
}

getList()
statusOptions.value = [
  {
    status_value: 0,
    status_name: 'Disabled'
  },
  {
    status_value: 1,
    status_name: 'Active'
  }
]
</script>
<style scoped>
/* tree border */
.tree-border {
  margin-top: 5px;
  border: 1px solid #e5e6e7;
  background: var(--base-bg-main) none;
  border-radius: 4px;
  width: 100%;
}
.el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
.el-dropdown {
  vertical-align: middle;
}
</style>
