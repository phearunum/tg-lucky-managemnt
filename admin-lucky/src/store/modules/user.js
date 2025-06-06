import { login, logout, getInfo, oauthCallback } from '@/api/system/login'
import { getToken, setToken, removeToken, getUserInfo, setUserInfo, removeUserInfo, removeMenu } from '@/utils/auth'
import defAva from '@/assets/images/profile.jpg'
import Cookies from 'js-cookie'
import { encrypt } from '@/utils/jsencrypt'
const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: '',
    token: getToken(),
    name: '',
    avatar: defAva,
    roles: [],
    permissions: [],
    menu: [],
    userId: 0,
    roleId: 0,
    authSource: '',
    userSession: getUserInfo()
  }),
  actions: {
    setAuthSource(source) {
      this.authSource = source
    },

    login(userInfo) {
      const username = userInfo.username.trim()
      const password = userInfo.password
      return new Promise((resolve, reject) => {
        login(username, password)
          .then((res) => {
            if (res.statusCode == 200) {
              setToken(res.data.token)
              console.log(res.data.user)
              const user = { username: res.data.user.username, userId: res.data.user.id, roleId: res.data.user.role_id, companyId: res.data.user.company_id, companyName: 'App', companyCode: res.data.user.username };

              setUserInfo(user)
              this.token = res.data.token
              this.userId = res.data.user.id
              this.avatar = "@/assets/images/profile.jpg"
              resolve() //then
            } else {
              console.log('login error ', res)
              reject(res) //catch
            }
          })
          .catch((error) => {
            reject(error)
          })
      })
    },

    oauthLogin(data, param) {
      return new Promise((resolve, reject) => {
        oauthCallback(data, param).then((res) => {
          const { code, data } = res
          if (code == 200) {
            setToken(data.token)
            this.token = data.token
            this.userId = data.data.id
            Cookies.set('username', data.userName, { expires: 30 })
            Cookies.set('password', encrypt(data.password), { expires: 30 })
            Cookies.set('rememberMe', true, { expires: 30 })
            resolve(res) //then
          } else {
            console.log('login error ', res)
            reject(res) //catch processing
          }
        })
      })
    },

    getInfo() {

      if (!getUserInfo()) {

      }
      return new Promise((resolve, reject) => {
        getInfo()
          .then((res) => {
            const data = res.data
            // const avatar = data.user.avatar == '' ? defAva : data.user.avatar

            if (data.roles && data.roles.length > 0) {

              this.roles = data.roles
              this.permissions = data.permissions
            } else {
              this.roles = ['ROLE_DEFAULT']
            }

            // this.name = data.user.nickName
            this.avatar = 'avatar'
            //   this.userInfo = data.user
            //  this.userId = data.user.id
            // this.roleId = data.user.role_id
            resolve(res)
          })
          .catch((error) => {
            console.error(error)
            reject('Failed to get user information')
          })
      })
    },

    logOut() {
      return new Promise((resolve, reject) => {

        logout(this.token)
          .then((res) => {
            this.token = ''
            this.roles = []
            this.permissions = []
            removeToken()
            removeUserInfo()
            removeMenu()
            resolve(res)
          })
          .catch((error) => {
            reject(error)
          })
      })
    },
    // front end log out
    fedLogOut() {
      return new Promise((resolve) => {
        this.token = ''
        removeToken()
        resolve()
        removeMenu()
      })
    }
  }
})
export default useUserStore
