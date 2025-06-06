import router from './router'
import { ElMessage } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { getToken } from '@/utils/auth'
import { isHttp } from '@/utils/validate'
import useUserStore from '@/store/modules/user'
import useSettingsStore from '@/store/modules/settings'
import usePermissionStore from '@/store/modules/permission'
NProgress.configure({ showSpinner: false });

const whiteList = ['/login', '/auth-redirect', '/bind', '/register', '/socialLogin'];

router. beforeEach((to, from, next) => {
  NProgress. start()
  if (getToken()) {
    to.meta.title && useSettingsStore().setTitle(to.meta.title)
    /* has token*/
    if (to.path === '/login') {
      next({ path: '/' })
      NProgress. done()
    } else {
      if (useUserStore(). roles. length === 0) {
        // Determine whether the current user has pulled user_info information

        useUserStore().getInfo().then(() => {

          usePermissionStore().generateRoutes().then(accessRoutes => {
            // Generate an accessible routing table based on roles permissions
            accessRoutes. forEach(route => {
              if (!isHttp(route. path)) {
                router.addRoute(route) // Dynamically add an accessible routing table
              }
            })
            next({ ...to, replace: true }) // hack method to ensure that addRoutes has completed
          })
        }).catch(err => {
          useUserStore(). logOut(). then(() => {
            ElMessage.error(err != undefined ? err : 'Login failed')
            next({ path: '/' })
          })
        })
      } else {
        next()
      }
    }
  } else {
    // no token
    if (whiteList. indexOf(to. path) !== -1) {
      // In the login-free whitelist, enter directly
      next()
    } else {
      next(`/login?redirect=${to.fullPath}`) // otherwise all redirect to the login page
      NProgress. done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
