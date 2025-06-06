import request from '@/utils/request'

// 查询菜单列表
export function listMenu(query) {
  return request({
    url: 'account/menus/list',
    method: 'get',
    params: query
  })
}
export function listMenuById(menuId) {
  return request({
    url: 'account/menus/list/' + menuId,
    method: 'get',

  })
}
// 查询菜单列表
export function MenuById(menuId) {
  return request({
    url: 'account/menus/' + menuId,
    method: 'get',
  })
}
// 查询菜单详细
export function getMenu(menuId) {
  return request({
    url: 'account/menus/id/' + menuId,
    method: 'get',
  })
}

// 查询菜单下拉树结构
export function treeselect() {
  return request({
    url: 'account/menus/treeSelect',
    method: 'get'
  })
}

// 根据角色ID查询菜单下拉树结构
export function roleMenuTreeselect(roleId) {
  return request({
    url: 'account/menus/roleMenuTreeselect/' + roleId,
    method: 'get',
  })
}

// 新增菜单
export const addMenu = (data) => {
  return request({
    url: 'account/menus/create',
    method: 'post',
    data: data,
  })
}

// 修改菜单
export function updateMenu(id = 0, data) {
  return request({
    url: 'account/menus/update',
    method: 'put',
    data: data
  })
}

// 删除菜单
export function delMenu(menuId) {
  return request({
    url: 'account/menus/delete/' + menuId,
    method: 'delete'
  })
}

//排序
export function changeMenuSort(data) {
  return request({
    url: 'account/menus/ChangeSort',
    method: 'GET',
    params: data
  })
}

// 获取路由
export const getRouters = (query) => {
  return request({
    url: 'account/menus/menu',
    method: 'get',
    params: query
  })
}
