import request from '@/utils/request'

export function login(username, password) {
  const data = {
    username,
    password,
  }
  return request({
    url: 'account/auth/login',
    method: 'POST',
    data: data,
  })
}
export function getInfo() {
  return request({
    url: 'account/permiss/info',
    method: 'get',
  })
}

export function logout() {
  return request({
    url: 'account/auth/logout',
    method: 'get',

  })
}

/**
 * @returns
 */
export function register(data) {
  return request({
    url: 'account/auth/register',
    method: 'post',
    data: data
  })
}

/**
 * @param {*} data
 * @param {*} params
 * @returns
 */
export function oauthCallback(data, params) {
  return request({
    url: 'account/auth/callback',
    method: 'post',
    data: data,
    params: params
  })
}
