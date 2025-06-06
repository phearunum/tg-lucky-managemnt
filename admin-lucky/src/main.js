import { createApp } from 'vue'
import Cookies from 'js-cookie'
import { uuid } from 'vue-uuid';

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
// import 'dayjs/locale/zh-cn'
import 'dayjs/locale/en'
import '@/assets/styles/index.scss' // global css

import App from './App'
import router from './router'
import directive from './directive' // directive
// register command
import plugins from './plugins' // plugins
import { downFile } from '@/utils/request'
//import signalR from '@/signalr/signalr'
import i18n from './i18n/index'
import pinia from '@/store/index'
// svg icon
import 'virtual:svg-icons-register'
import SvgIcon from '@/components/SvgIcon/index.vue'
import elementIcons from '@/components/SvgIcon/svgicon'

import './permission' // permission control

import { getConfigKey } from '@/api/system/config'
import { getDicts } from '@/api/system/dict/data'
import { parseTime, resetForm, addDateRange, handleTree, selectDictLabel, download } from '@/utils/ruoyi'
import md5 from 'js-md5';
// pagination component
import Pagination from '@/components/Pagination'
// Customize table tool component
import RightToolbar from '@/components/RightToolbar'
// file upload component
import FileUpload from '@/components/FileUpload'
// image upload component
import ImageUpload from '@/components/ImageUpload'
// image preview component
import ImagePreview from '@/components/ImagePreview'
// dictionary tag component
import DictTag from '@/components/DictTag'
import VueFlvPlayer from '@/components/vueFlv.vue';
// el-date-picker shortcut option
import dateOptions from '@/utils/dateOptions'
// Ifram
import { io } from 'socket.io-client'
// @ts-ignore


const app = createApp(App)
//signalR.init(import.meta.env.VITE_APP_SOCKET_API)
//app.config.globalProperties.signalr = signalR
// Global method mount
app.config.globalProperties.getConfigKey = getConfigKey
app.config.globalProperties.getDicts = getDicts
app.config.globalProperties.download = download
app.config.globalProperties.downFile = downFile
app.config.globalProperties.parseTime = parseTime
app.config.globalProperties.resetForm = resetForm
app.config.globalProperties.handleTree = handleTree
app.config.globalProperties.addDateRange = addDateRange
app.config.globalProperties.selectDictLabel = selectDictLabel
app.config.globalProperties.dateOptions = dateOptions
app.config.globalProperties.$uuid = uuid;
app.config.globalProperties.$md5 = md5;

// Websocket 
/*
app.config.globalProperties.$soketio = io(`${import.meta.env.VITE_APP_SOCKET_API}`, {
    path: '/ws/socket.io',
    transports: ['websocket'], // Ensure it only uses WebSocket transport (optional)
    reconnectionAttempts: 3,    // Retry 3 times if the connection fails
    timeout: 20000,             // Set timeout to 20 seconds
    autoConnect: true,          // Auto-connect on initialization

});
*/
// Global component mount
app.component('DictTag', DictTag)
app.component('Pagination', Pagination)
app.component('UploadFile', FileUpload)
app.component('UploadImage', ImageUpload)
app.component('ImagePreview', ImagePreview)
app.component('RightToolbar', RightToolbar)
app.component('svg-icon', SvgIcon)
app.component('flvPlayer', VueFlvPlayer)




directive(app)
//console.log(md5("1234566"))
app.use(pinia).use(router).use(plugins).use(ElementPlus, {}).use(elementIcons).use(i18n).mount('#app')

