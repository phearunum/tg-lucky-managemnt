<template>
  <div class="app-container">
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" @click="handleAdd" v-hasPermi="['system:menu:add']">{{ $t('btn.add') }}</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="info" plain icon="Sort" @click="toggleExpandAll">{{ $t('btn.expand') }}/{{ $t('btn.collapse') }}</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table
      :default-expand-all="isExpandAll"
      :data="menuList"
      node-key="id"
      row-key="id"
      lazy
      :load="loadMenu"
      :tree-props="{ children: 'children', hasChildren: 'alwaysShow' }"
      :header-cell-style="{
        background: 'var(--el-fill-color-light)',
        color: 'var(--el-text-color-primary)'
      }">
      <el-table-column prop="name" :label="$t('m.menuName')" :show-overflow-tooltip="true" width="200">
        <template #default="scope">
          <span v-if="scope.row.titleKey">
            {{ $t(scope.row.titleKey) }}
          </span>
          <span v-else>
            {{ scope.row.name }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="icon" :label="$t('m.icon')" align="center" width="60">
        <template #default="scope">
          <svg-icon :name="scope.row.icon" />
        </template>
      </el-table-column>
      <el-table-column prop="id" :label="$t('m.id')" :show-overflow-tooltip="true" width="80" align="center"></el-table-column>
      <el-table-column prop="menuType" :label="$t('m.menuType')" align="center" width="80">
        <template #default="scope">
          <el-tag type="danger" v-if="scope.row.menuType == 'M' && scope.row.isFrame == 1">{{ $t('m.link') }}</el-tag>
          <el-tag v-else-if="scope.row.menuType == 'C'">{{ $t('m.menu') }}</el-tag>
          <el-tag type="success" v-else-if="scope.row.menuType == 'M'">{{ $t('m.directory') }}</el-tag>
          <el-tag type="warning" v-else-if="scope.row.menuType == 'F'">{{ $t('m.button') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="orderNum" :label="$t('m.sort')" width="90" sortable align="center">
        <template #default="scope">
          <span v-show="editIndex != scope.$index" @click="editCurrRow(scope.$index)">{{ scope.row.orderNum }}</span>
          <el-input
            :ref="setColumnsRef"
            v-show="editIndex == scope.$index"
            v-model="scope.row.orderNum"
            @blur="handleChangeSort(scope.row)"></el-input>
        </template>
      </el-table-column>
      <el-table-column prop="perms" :label="$t('m.authorityID')" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="component" :label="$t('m.componentPath')" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="hidden" :label="$t('m.isShow')" width="70" align="center">
        <template #default="scope">
          <el-tag v-if="scope.row.hidden == false">true</el-tag>
          <el-tag type="danger" v-if="scope.row.hidden !== false">false</el-tag>
        </template></el-table-column
      >
      <el-table-column prop="menuStatus" :label="$t('m.menuState')" width="80" align="center">
        <template #default="scope">
          <el-tag v-if="scope.row.menuStatus == true">{{ scope.row.menuStatus }}</el-tag>
          <el-tag type="danger" v-if="scope.row.menuStatus !== true">false</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('common.API')" align="left" prop="apiURL" :show-overflow-tooltip="true">
        <template #default="scope">
          <span>{ {{ scope.row.apiURL }} }</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('btn.operate')" align="center" width="170" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button text size="small" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['system:menu:edit']"></el-button>
          <el-button text size="small" icon="Plus" @click="handleAdd(scope.row)" v-hasPermi="['system:menu:add']"></el-button>
          <el-button text size="small" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['system:menu:remove']"></el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog :title="title" v-model="open" width="680px" append-to-body>
      <el-form ref="menuRef" :model="form" :rules="rules" label-width="120px">
        <el-row>
          <el-col :lg="12">
            <el-form-item :label="$t('m.parentMenu')">
              <el-cascader
                class="w100"
                :options="menuOptions"
                :props="{ checkStrictly: true, value: 'id', label: 'name', emitPath: false }"
                placeholder="请选择上级菜单"
                clearable
                v-model="form.subof">
                <template #default="{ node, data }">
                  <span>{{ data.name }}</span>
                  <span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
                </template>
              </el-cascader>
              <!--
              <el-select v-model="form.subof" placeholder="Role" style="width: 100%">
                <el-option v-for="item in menuOptionsParent" :key="item.id" :label="item.name" :value="item.id"
                 >
                </el-option>
              </el-select>
               -->
            </el-form-item>
          </el-col>

          <el-col :lg="12">
            <el-form-item :label="$t('m.menuType')" prop="menuType">
              <el-radio-group v-model="form.menuType">
                <el-radio label="M">{{ $t('m.directory') }}</el-radio>
                <el-radio label="C">{{ $t('m.menu') }}</el-radio>
                <el-radio label="F">{{ $t('m.button') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :lg="24">
            <el-form-item :label="$t('m.apiURL')" prop="apiURL">
              <el-input v-model="form.apiURL" placeholder="Enter apiURL" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item :label="$t('m.menuName')" prop="name">
              <el-input v-model="form.name" placeholder="Enter menu name" />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item label="menu(key)" prop="menuNameKey">
              <template #label>
                <span>
                  <el-tooltip content="eg：menu.system" placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.menuNameKey') }}
                </span>
              </template>
              <el-input v-model="form.titleKey" placeholder="Enter menu key " />
            </el-form-item>
          </el-col>
          <el-col :lg="12">
            <el-form-item :label="$t('m.sort')" prop="orderNum">
              <el-input-number v-model="form.orderNum" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :lg="24" v-if="form.menuType != 'F'">
            <el-form-item :label="$t('m.icon')" prop="icon">
              <el-popover placement="bottom-start" :teleported="false" :width="540" :visible="showChooseIcon" trigger="click">
                <template #reference>
                  <el-input v-model="form.icon" placeholder="Icon" @click="showSelectIcon" readonly>
                    <template #prefix>
                      <svg-icon v-if="form.icon" :name="form.icon" />
                      <el-icon v-else>
                        <search />
                      </el-icon>
                    </template>
                  </el-input>
                </template>
                <icon-select ref="iconSelectRef" @selected="selected" />
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType != 'F'">
            <el-form-item>
              <template #label>
                <span>
                  <el-tooltip content="" placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.isFrame') }}
                </span>
              </template>
              <el-radio-group v-model="form.isFrame">
                <el-radio :label="1" value="1">{{ $t('common.yes') }}</el-radio>
                <el-radio :label="0" value="0">{{ $t('common.no') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType != 'F'">
            <el-form-item prop="path">
              <template #label>
                <span>
                  <el-tooltip content="" placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.routePath') }}
                </span>
              </template>
              <el-input v-model="form.path" placeholder="Rote path" />
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType == 'C'">
            <el-form-item prop="component">
              <template #label>
                <span>
                  <el-tooltip
                    content="The component path to access, such as: `system/user/index`, default in the `views` directory"
                    placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.componentPath') }}
                </span>
              </template>
              <el-input v-model="form.component" placeholder="Please enter component path" />
            </el-form-item>
          </el-col>

          <el-col :lg="12" v-if="form.menuType != 'M'">
            <el-form-item>
              <el-input v-model="form.perms" placeholder="Please enter the permission ID" maxlength="100" />

              <template #label>
                <span>
                  <el-tooltip
                    content="Permission characters defined in the controller, such as: [ActionPermissionFilter(Permission = 'system:user:delete')])"
                    placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.permissionStr') }}
                </span>
              </template>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType == 'C'">
            <el-form-item prop="isCache">
              <template #label>
                <span>
                  <el-tooltip
                    content="If you select Yes, it will be cached by `keep-alive`, and the `name` and address of the matching component need to be consistent"
                    placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.isCache') }}
                </span>
              </template>
              <el-radio-group v-model="form.noCache">
                <el-radio :label="false" value="false">{{ $t('common.yes') }}</el-radio>
                <el-radio :label="true" value="true">{{ $t('common.no') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType != 'F'">
            <el-form-item prop="visible">
              <template #label>
                <span>
                  <el-tooltip
                    content="If you choose to hide, the route will not appear in the sidebar, but it can still be accessed"
                    placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.isShow') }}
                </span>
              </template>
              <el-radio-group v-model="form.hidden">
                <el-radio :label="false" value="false">{{ $t('common.yes') }}</el-radio>
                <el-radio :label="true" value="true">{{ $t('common.no') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType != 'F'">
            <el-form-item>
              <template #label>
                <span>
                  <el-tooltip content="If you choose to disable, the route will not appear in the sidebar and cannot be accessed" placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.menuState') }}
                </span>
              </template>
              <el-radio-group v-model="form.menuStatus">
                <el-radio :label="true" value="true">{{ $t('common.yes') }}</el-radio>
                <el-radio :label="false" value="false">{{ $t('common.no') }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :lg="12" v-if="form.menuType != 'F'">
            <el-form-item>
              <template #label>
                <span>
                  <el-tooltip content="If you choose to disable, the route will not appear in the sidebar and cannot be accessed" placement="top">
                    <el-icon :size="15">
                      <questionFilled />
                    </el-icon>
                  </el-tooltip>
                  {{ $t('m.hasSub') }}
                </span>
              </template>
              <el-radio-group v-model="form.alwaysShow">
                <el-radio :label="true" value="true">{{ $t('common.yes') }}</el-radio>
                <el-radio :label="false" value="false">{{ $t('common.no') }}</el-radio>
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
  </div>
</template>

<script setup name="sysmenu">
import { addMenu, delMenu, getMenu, listMenu, updateMenu, changeMenuSort as changeSort, listMenuById } from '@/api/system/menu'
import SvgIcon from '@/components/SvgIcon'
import IconSelect from '@/components/IconSelect'

const { proxy } = getCurrentInstance()

const menuList = ref([])
const open = ref(false)
const loading = ref(true)
const showSearch = ref(true)
const title = ref('')
const menuOptions = ref([])
const menuOptionsParent = ref([])
const menuQueryOptions = ref([])
const isExpandAll = ref(false)
const refreshTable = ref(true)
const showChooseIcon = ref(false)
const iconSelectRef = ref(null)
const loadNodeMap = new Map()
const menuRef = ref(null)
const listRef = ref(null)
const state = reactive({
  form: {},
  queryParams: {
    num_page: 0, //num_page
    limit: 10,
    menuName: undefined,
    visible: undefined,
    menuTypeIds: 'M,C',
    parentId: undefined
  },
  rules: {
    menuName: [{ required: true, message: 'required', trigger: 'blur' }],
    menuNameKey: [{ pattern: /^[A-Za-z].+$/, message: 'required', trigger: 'blur' }],
    orderNum: [{ required: true, message: 'required', trigger: 'blur' }],
    path: [{ required: true, message: 'required', trigger: 'blur' }],
    visible: [{ required: true, message: 'required', trigger: 'blur' }]
  },
  sys_show_hide: [],
  sys_normal_disable: []
})

const { queryParams, form } = toRefs(state)

/** Get Menu List */
function getList() {
  loading.value = true
  if (queryParams.value.subof != undefined) {
    queryParams.value.menuTypeIds = ''
  }
  listMenu(queryParams.value).then((response) => {
    menuList.value = response.data
    menuOptionsParent.value = response.data
    loading.value = false
  })
}
/** Query menu drop-down tree structure */
function getTreeselect() {
  listMenu({
    num_page: 0, //num_page
    limit: 50,
    menuTypeIds: 'M,C,F'
  }).then((response) => {
    menuOptions.value = response.data
  })
}
function refreshMenu(pid) {
  loading.value = true

  getList()
}
/** cancel */
function cancel() {
  open.value = false
  reset()
}
/** reset from */
function reset() {
  form.value = {
    id: undefined,
    parentId: 0,
    menuName: undefined,
    icon: undefined,
    menuType: 'M',
    orderNum: 999,
    isFrame: 0,
    isCache: 0,
    visible: 0,
    status: 0
  }
  proxy.resetForm('menuRef')
}
/** Show dropdown icon */
function showSelectIcon() {
  iconSelectRef.value.reset()
  showChooseIcon.value = true
}
/** select icon */
function selected(name) {
  form.value.icon = name
  showChooseIcon.value = false
}
/** Search  */
function handleQuery() {
  getList()
}
/** Reset button action */
function resetQuery() {
  proxy.resetForm('queryRef')
  handleQuery()
}
/** add button  */
function handleAdds(row) {
  proxy.$refs.menuRef.handleAdd(row)
  open.value = true
}
function handleAdd(row) {
  reset()
  getTreeselect()
  console.log(row)
  if (row != null && row.id != undefined) {
    form.value.subof = row.id
  } else {
    form.value.parentId = 0
  }
  open.value = true
  title.value = proxy.$t('btn.add')
  proxy.$refs.menuRef.handleAdd(row)
}

async function handleUpdate(row) {
  reset()
  getTreeselect()
  getMenu(row.id).then((response) => {
    form.value = response.data
    open.value = true
    title.value = proxy.$t('btn.edit')
  })
}
function ensureLeadingSlash(path) {
  if (path.charAt(0) !== '/') {
    return '/' + path
  }
  return path
}
/** Submit form */
function submitForm() {
  proxy.$refs['menuRef'].validate((valid) => {
    if (valid) {
      if (form.value.menuType == 'M') {
        form.value.component = 'Layout'
        form.value.redirect = 'noRedirect'
        form.value.subof = 0
        form.value.isFrame = 0
        form.value.path = ensureLeadingSlash(form.value.path)
      }
      if (form.value.id != undefined) {
        updateMenu(form.value.id, form.value).then(() => {
          proxy.$modal.msgSuccess('Success')
          open.value = false
          refreshMenu(form.value.parentId)
        })
      } else {
        addMenu(form.value).then(() => {
          proxy.$modal.msgSuccess('Success')
          open.value = false
          refreshMenu(form.value.parentId)
        })
      }
    }
  })
}
/** Delete */
function handleDelete(row) {
  proxy.$modal
    .confirm('是否确认删除名称为"' + row.menuName + '"的数据项?')
    .then(function () {
      return delMenu(row.id)
    })
    .then(() => {
      getList()

      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}

// Dynamic ref setting value
const columnRefs = ref([])
const setColumnsRef = (el) => {
  if (el) {
    columnRefs.value.push(el)
  }
}
const editIndex = ref(-1)
// Show Edit Sort
function editCurrRow(rowId) {
  editIndex.value = rowId
  setTimeout(() => {
    columnRefs.value[rowId].focus()
  }, 100)
}
// 保存排序
function handleChangeSort(info) {
  editIndex.value = -1
  proxy
    .$confirm('Do you want to save ?')
    .then(function () {
      return changeSort({ value: info.orderNum, id: info.id })
    })
    .then(() => {
      handleQuery()
      proxy.$modal.msgSuccess('修改成功')
    })
    .catch(() => {
      handleQuery()
    })
}

const loadMenu = (row, treeNode, resolve) => {
  listMenuById(row.id).then((res) => {
    if (res.data[0].children == undefined) {
      resolve(null)
      return
    } else {
      loadNodeMap.set(row.id, { row, treeNode, resolve })
      resolve(res.data[0]?.children)
    }
  })
}

listMenu({ menuTypeIds: 'M,C' }).then((response) => {
  menuQueryOptions.value = response.data
})
getList()
</script>
