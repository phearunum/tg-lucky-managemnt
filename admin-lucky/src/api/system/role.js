import request from '@/utils/request'

// 查询角色列表
export function listRole(query) {
  return request({
    url: 'account/roles/list',
    method: 'get',
    params: query
  })
}

// 查询角色详细
export function getRole(roleId) {
  return request({
    url: 'account/roles/list/' + roleId,
    method: 'get'
  })
}
export const getListSelect = () => {
  return request({
    url: 'account/roles/list',
    method: 'get'
  })
}

// 新增角色
export const addRole = (data) => {
  return request({
    url: 'account/roles/create',
    method: 'post',
    data: data,
  })
}

// 修改角色
export function updateRole(data) {
  return request({
    url: 'account/roles/update',
    method: 'put',
    data: data
  })
}

// 角色数据权限
export function dataScope(data) {
  return request({
    url: 'account/permiss/rolesacess',
    method: 'post',
    data: data
  })
}

// 角色状态修改
export function changeRoleStatus(roleId, status) {
  const data = {
    roleId,
    status
  }
  return request({
    url: '/systemaccount/roles/changeStatus',
    method: 'put',
    data: data
  })
}

// 删除角色
export function delRole(roleId) {
  return request({
    url: 'account/roles/delete/' + roleId,
    method: 'delete'
  })
}

// 导出角色
export function exportRole(query) {
  return request({
    url: '/systemaccount/roles/export',
    method: 'get',
    params: query
  })
}
