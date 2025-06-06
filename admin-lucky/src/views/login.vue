<template>
  <div class="logins">
    <el-row :gutter="24">
      <el-col :xs="0" :sm="16" :md="16" :lg="16" :xl="16" class="bg-login-r" :style="b"> </el-col>
      <el-col :xs="24" :sm="8" :md="8" :lg="8" :xl="8" class="bg-login-l login" style="background-color: #169af0">
        <starBackground></starBackground>

        <el-form ref="loginRef" :model="loginForm" :rules="loginRules" class="login-form">
          <div style="font-size: 1.7rem; font-weight: 600; padding-bottom: 20px; padding-top: 20px; font-family: 'Kh Muol'">
            {{ proxy.$t('menu.systemTools') }}
          </div>
          <br />

          <LangSelect title="Multilingual settings" class="langSet" />
          <el-form-item prop="username">
            <el-input v-model="loginForm.username" type="text" size="default" auto-complete="off" :placeholder="$t('login.account')">
              <template #prefix>
                <svg-icon name="user" class="el-input__icon input-icon" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              size="default"
              auto-complete="off"
              :placeholder="$t('login.password')"
              @keyup.enter="handleLogin">
              <template #prefix>
                <svg-icon name="password" class="el-input__icon input-icon" />
              </template>
            </el-input>
          </el-form-item>
          <el-form-item style="width: 100%">
            <el-button :loading="loading" size="default" type="primary" style="width: 100%" @click.prevent="handleLogin">
              <span v-if="!loading">{{ $t('login.btnLogin') }}</span>
              <span v-else>logging in...</span>
            </el-button>
          </el-form-item>
          <el-form-item style="width: 100%; margin-top: 50px">
            <el-button :loading="loading" size="default" type="infor" style="width: 100%" @click.prevent="setBackground">
              Change Wallpaper
            </el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>

    <!--  bottom  -->
    <div class="el-login-footer">
      <div v-html="defaultSettings.copyright"></div>
    </div>
  </div>
</template>

<script setup name="login">
import Cookies from 'js-cookie'
import { encrypt, decrypt } from '@/utils/jsencrypt'
import defaultSettings from '@/settings'
import starBackground from '@/views/components/starBackground.vue'
import LangSelect from '@/components/LangSelect/index.vue'
import useUserStore from '@/store/modules/user'
import { removeToken } from '@/utils/auth'
const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const { proxy } = getCurrentInstance()

const loginForm = ref({
  username: '',
  password: '',
  rememberMe: false,
  code: '',
  uuid: ''
})
const host = window.location.host
const loginRules = {
  username: [{ required: true, trigger: 'blur', message: ' Enter usernmae' }],
  password: [{ required: true, trigger: 'blur', message: 'Enter password' }],
  code: [{ required: true, trigger: 'change', message: 'Enter code' }]
}

const codeUrl = ref('')
const loading = ref(false)
const captchaOnOff = ref('')
const register = ref(false)
const redirect = ref()
const b = {
  'background-image': `url(${Cookies.get('wallpaper') || '/static/images/background.jpeg'})`
}
redirect.value = route.query.redirect

function handleLogin() {
  proxy.$refs.loginRef.validate((valid) => {
    if (valid) {
      loading.value = true
      // cookie
      if (loginForm.value.rememberMe) {
        Cookies.set('username', loginForm.value.username, { expires: 30, path: host })
        Cookies.set('password', encrypt(loginForm.value.password), { expires: 30, path: host })
        Cookies.set('rememberMe', loginForm.value.rememberMe, { expires: 30, path: host })
      } else {
        // remove cookie
        Cookies.remove('username')
        Cookies.remove('password')
        Cookies.remove('rememberMe')
      }
      // action
      userStore
        .login(loginForm.value)
        .then((data) => {
          proxy.$modal.msgSuccess(proxy.$t('login.loginSuccess'))
          router.push({ path: redirect.value || '/' })
        })
        .catch((error) => {
          console.error(error)
          proxy.$modal.msgError(error.msg)
          loading.value = false
        })
    }
  })
}
function setBackground() {
  proxy
    .$prompt('System Wallpaper ', proxy.$t('system.System Message'), {
      confirmButtonText: proxy.$t('common.ok'),
      cancelButtonText: proxy.$t('common.cancel'),
      closeOnClickModal: false,
      inputPattern: /^.{5,100}$/,
      inputErrorMessage: proxy.$t('system. Required')
    })
    .then(({ value }) => {
      const backg = Cookies.set('wallpaper', value, { expires: 30, path: host })
      if (backg) {
        proxy.$modal.msgSuccess('Wallpaper changed successfully!')
      } else {
        proxy.$modal.msgError('Failed to change wallpaper.')
      }
      window.location.reload()
    })
    .catch(() => {})
}

function getCookie() {
  const username = Cookies.get('username')
  const password = Cookies.get('password')
  const rememberMe = Cookies.get('rememberMe')
  loginForm.value = {
    username: username === undefined ? loginForm.value.username : username,
    password: password === undefined ? loginForm.value.password : decrypt(password),
    rememberMe: rememberMe === undefined ? false : Boolean(rememberMe)
  }
}
function onAuth(type) {
  userStore.setAuthSource(type)

  switch (type) {
    default:
      window.location.href = import.meta.env.VITE_APP_BASE_API + '/auth/Authorization?authSource=' + type
      break
  }
}

getCookie()
</script>

<style lang="scss" scoped>
@import '@/assets/styles/login.scss';

.login-icon {
  width: 30px;
  margin-right: 10px;
  cursor: pointer;
}
.other-login {
  padding: 0px 10px 5px;
}
.title {
  color: white !important;
}
.box-login {
  box-shadow: rgba(99, 99, 99, 0.2) 0px 2px 8px 0px;
  height: 300px;
  background: green;
}
.bg-login-l {
  height: 100vh;

  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}
.bg-login-r {
  min-height: 100vh;
  justify-content: center;
  align-items: center;
  text-align: center;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  margin: auto;
  // background-image: url('@/assets/images/background.jpeg');
  background-repeat: no-repeat;
  background-position: center center;
  background-size: cover;
}
</style>
