// import Cookies from 'js-cookie'

const TokenKey = 'Authorization'

export function getToken() {
  return localStorage.getItem(TokenKey)
}

export function setToken(token) {
  return localStorage.setItem(TokenKey, token)
}

export function removeToken() {
  return localStorage.removeItem(TokenKey)
}

// User init 

export function getUserInfo() {
  return localStorage.getItem('user_session')
}

export function setUserInfo(data) {
  return localStorage.setItem('user_session', JSON.stringify(data))
}
export function removeUserInfo() {
  return localStorage.removeItem('user_session')
}

export function setMenu(data) {
  return localStorage.setItem('menu', JSON.stringify(data))
}

export function getMenu() {
  return localStorage.getItem('menu')
}

export function removeMenu() {
  return localStorage.removeItem('menu')
}