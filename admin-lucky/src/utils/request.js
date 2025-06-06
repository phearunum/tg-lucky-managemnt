import axios from 'axios'
import { ElMessageBox, ElMessage, ElLoading } from 'element-plus'
import { getToken, getUserInfo, setToken, removeToken } from '@/utils/auth'
import useUserStore from '@/store/modules/user'
import { blobValidate } from '@/utils/ruoyi'
import { saveAs } from 'file-saver'

let downloadLoadingInstance
// Solve the problem that the backend cannot obtain cookies across domains
// axios.defaults.withCredentials = true
axios.defaults.headers['Content-Type'] = 'application/json;charset=utf-8'
// 创建axios实例
const service = axios.create({
  // The request configuration in axios has the baseURL option, indicating the public part of the request URL
  baseURL: import.meta.env.VITE_APP_BASE_API,
  // time out
  timeout: 30000
})

// request interceptor
service.interceptors.request.use(
  (config) => {
    // Check if the token exists
    const token = getToken();
    if (token) {
      // Add query parameter using URL constructor
      //const url = new URL(config.url, config.baseURL);
      //url.searchParams.append('key', token);
      //config.url = url.toString();

      // Add headers to the request
      const { userId, roleId, companyId, username } = JSON.parse(getUserInfo());
      config.headers['Authorization'] = 'Bearer ' + token;
      config.headers['id'] = userId;
      config.headers['roleid'] = roleId;
      config.headers['companyId'] = companyId;
      config.headers['username'] = username;
    }
    return config;
  },
  (error) => {
    console.error(error);
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  (res) => {
    if (res.data.statusCode == 401) {
      ElMessage({
        message: res.data.message,
        type: 'error',
        duration: 3 * 1000,
        grouping: true
      })
      Promise.reject(res.data.message)
      return
    }

    if (res.status == 500) {
      ElMessage({
        message: 'Network error',
        type: 'error',
        duration: 3 * 1000,
        grouping: true
      })
      Promise.reject('Network error')
      return
    }
    if (res.status == 404) {
      ElMessage({
        message: 'Request Not Found',
        type: 'error',
        duration: 3 * 1000,
        grouping: true
      })
      Promise.reject('Network error')
      return
    }
    if (res.data.token !== undefined) {
      setToken(res.data.token)
    }
    if (res.headers.Authorization !== undefined) {
      setToken(res.headers.header.Authorization)
    }
    const { statusCode, message, code } = res
    if (res.request.responseType === 'blob' || res.request.responseType === 'arraybuffer') {
      return res
    }

    if (statusCode == 401 || code == 401) {
      ElMessageBox.confirm('System Unuathorized ', 'System Prompt', {
        confirmButtonText: 'Login again',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        useUserStore()
          // .logOut()
          .then(() => {
            location.href = import.meta.env.VITE_APP_ROUTER_PREFIX + 'index'
          })
      })
      return Promise.reject('Invalid session, or session has expired, please log in again')
    } else if (statusCode == 0 || statusCode == 1 || statusCode == 110 || statusCode == 101 || statusCode == 403 || statusCode == 500 || statusCode == 429) {
      ElMessage({
        message: message,
        type: 'error'
      })
      return Promise.reject(res.data)
    } else {
      return res.data
    }
  },
  (error) => {
    let { message } = error
    console.log('error stage')
    if (message == 'Network Error') {
      message = 'Backend interface connection exception'

    } else if (message.includes('timeout')) {
      message = 'System interface request timed out'
    } else if (message.includes('Request failed with status code 429')) {
      message = 'The request is too frequent, please try again later'
    } else if (message.includes('equest failed with status code 404')) {
      message = 'System interface request 404'
    } else if (message.includes('Unauthorized, Company License Disabled')) {
      message = 'Unauthorized, Company License Disabled'
    } else if (message.includes('Request failed with status code')) {
      console.log(message)
      message = 'system interface ' + message.substr(message.length - 3) + ' exception'

    } else if (message == 'equest failed with status code 401') {
      useUserStore()
        //  .logOut()
        .then(() => {
          location.href = import.meta.env.VITE_APP_ROUTER_PREFIX + 'index'
        })
    }
    if (message !== '') {
      ElMessage({
        message: message,
        type: 'error',
        duration: 3 * 1000,
        grouping: true
      })

    }

    return Promise.reject(error)
  }
)

/**
 * get方法，对应get请求
 * @param {String} url [请求的url地址]
 * @param {Object} params [请求时携带的参数]
 */
export function get(url, params) {
  return new Promise((resolve, reject) => {
    axios
      .get(url, {
        params: params
      })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export function post(url, params) {
  return new Promise((resolve, reject) => {
    axios
      .post(url, {
        params: params
      })
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

/**
 * 提交表单
 * @param {*} url
 * @param {*} data
 */
export function postForm(url, data, config) {
  return new Promise((resolve, reject) => {
    axios
      .post(url, data, config)
      .then((res) => {
        resolve(res.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

/**
 * 通用下载方法
 * @param {*} url 请求地址
 * @param {*} params 请求参数
 * @param {*} config 配置
 * @returns
 */
export async function downFile(url, params, config) {
  downloadLoadingInstance = ElLoading.service({ text: '正在下载数据，请稍候', background: 'rgba(0, 0, 0, 0.7)' })
  try {
    const resp = await service.get(url, {
      params,
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      responseType: 'blob',
      ...config
    })
    const { data } = resp

    var patt = new RegExp('filename=([^;]+\\.[^\\.;]+);*')
    var contentDisposition = decodeURI(resp.headers['content-disposition'])
    var result = patt.exec(contentDisposition)
    var fileName = result[1]
    fileName = fileName.replace(/\"/g, '')

    const isLogin = await blobValidate(data)
    if (isLogin) {
      const blob = new Blob([data])
      saveAs(blob, fileName)
    } else {
      const resText = await data.text()
      const rspObj = JSON.parse(resText)
      const errMsg = errorCode[rspObj.code] || rspObj.msg || errorCode['default']

      ElMessage({
        message: errMsg,
        type: 'error'
      })
    }
    downloadLoadingInstance.close()
  } catch (r) {
    console.error(r)
    ElMessage({
      message: '下载文件出现错误，请联系管理员！',
      type: 'error'
    })
    downloadLoadingInstance.close()
  }
}

export default service
