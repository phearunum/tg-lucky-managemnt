import request from '@/utils/request'
export function getEmployeeList(query) {
  return request({
    url: '/hr/employees/list',
    method: 'get',
    params: query
  })
}
export function getEmployeeByDept(query) {
  return request({
    url: `/hr/employees/dept`,
    method: 'get',
    params: query
  })
}
export const addEmployee = (data) => {
  return request({
    url: '/hr/employees/create',
    method: 'post',
    data: data,
  })
}
export function updateEmployee(data) {

  return request({
    url: '/hr/employees/update',
    method: 'put',
    data: data
  })
}

export function getEmployee(id) {
  return request({
    url: `/hr/employees/list/${id}`,
    method: 'get',
  })
}

export function deleteEmployee(id) {
  return request({
    url: `/hr/employees/delete/${id}`,
    method: 'delete',
  })
}
export function getSelectEmployee(query) {
  return request({
    url: '/employee/select',
    method: 'get',
    params: query
  })
}
