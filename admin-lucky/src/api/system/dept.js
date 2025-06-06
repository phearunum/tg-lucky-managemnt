import request from '@/utils/request'

export function listDept(query) {
  return request({
    url: '/hr/department/list',
    method: 'get',
    params: query
  })
}

export function getDeptById(id) {
  return request({
    url: '/hr/department/list/' + id,
    method: 'get',

  })
}

export function treeselect() {
  return request({
    url: '/hr/department/list',
    method: 'get'
  })
}
export function getDepartmentSelect(data) {
  return request({
    url: '/hr/department/list',
    method: 'get',
    params: data
  })
}
export function dept_create(data) {
  return request({
    url: '/hr/department/create',
    method: 'post',
    data: data
  })
}
export function update(data) {
  return request({
    url: '/hr/department/update',
    method: 'put',
    data: data
  })
}


export function deleted(userId) {
  return request({
    url: '/hr/department/delete/' + userId,
    method: 'delete'
  })
}
