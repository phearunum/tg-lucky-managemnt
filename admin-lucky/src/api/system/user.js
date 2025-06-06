import request from '@/utils/request'
import { praseStrZero } from '@/utils/ruoyi'
import { downFile } from '@/utils/request'

// 查询用户列表
export function listUser(query) {
  return request({
    url: 'account/users/list',
    method: 'get',
    params: query
  })
}

// 查询用户详细
export function getUser(userId) {
  return request({
    url: 'account/users/list/' + praseStrZero(userId),
    method: 'get'
  })
}

// 新增用户
export function addUser(data) {
  return request({
    url: 'account/users/create',
    method: 'post',
    data: data
  })
}

// 修改用户
export function updateUser(data) {
  return request({
    url: 'account/users/update',
    method: 'put',
    data: data
  })
}

export function delUser(userId) {
  return request({
    url: 'account/users/delete/' + userId,
    method: 'delete'
  })
}

// 导出用户
export async function exportUser(query) {
  // return request({
  //   url: '/system/User/export',
  //   method: 'get',
  //   params: query
  // })
  await downFile('/system/user/export', { ...query })
}

// 用户密码重置
export function resetUserPwd(data) {

  return request({
    url: 'account/users/change-password',
    method: 'put',
    data: data
  })
}

// 用户状态修改
export function changeUserStatus(userId, status) {
  const data = {
    userId,
    status
  }
  return request({
    url: '/system/user/changeStatus',
    method: 'put',
    data: data
  })
}

// 查询用户个人信息
export function getUserProfile() {
  return request({
    url: '/system/user/Profile',
    method: 'get'
  })
}

// 修改用户个人信息
export function updateUserProfile(data) {
  return request({
    url: '/system/user/profile',
    method: 'put',
    data: data
  })
}

// 用户密码重置
export function updateUserPwd(oldPassword, newPassword) {
  const data = {
    oldPassword,
    newPassword
  }
  return request({
    url: '/system/user/profile/updatePwd',
    method: 'put',
    params: data
  })
}

// 用户头像上传
export function uploadAvatar(data) {
  return request({
    url: '/system/user/profile/avatar',
    method: 'post',
    data: data
  })
}

// 下载用户导入模板
export function importTemplate() {
  return request({
    url: '/system/user/importTemplate',
    method: 'get',
    responseType: 'blob' //1.首先设置responseType对象格式为 blob:
  })
}

export function assignTelegramGroup(data) {

  return request({
    url: 'account/users/assign-group',
    method: 'put',
    data: data
  })
}