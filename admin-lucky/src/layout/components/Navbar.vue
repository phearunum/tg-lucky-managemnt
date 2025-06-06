<template>
  <div class="navbar" :data-theme="sideTheme" :class="appStore.device">
    <hamburger id="hamburger-container" :is-active="appStore.sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />
    <template v-if="appStore.device == 'desktop'">
      <breadcrumb id="breadcrumb-container" class="breadcrumb-container" v-if="!settingsStore.topNav" />
      <top-nav id="topmenu-container" class="topmenu-container" v-if="settingsStore.topNav" />
    </template>

    <div class="right-menu">
      <header-search id="header-search" class="right-menu-item" />
      <template v-if="appStore.device == 'desktop'">
        <!-- <zr-git title="source address" class="right-menu-item" />
         <zr-doc title="Document Address" class="right-menu-item" />-->
        <screenfull title="full screen" class="right-menu-item" />
      </template>
      <size-select title="Layout Size" class="right-menu-item" />
      <LangSelect title="Language Settings" class="right-menu-item" />
      <Notice title="Notice" class="right-menu-item" />

      <el-dropdown @command="handleCommand" class="right-menu-item avatar-container" trigger="hover">
        <span class="avatar-wrapper">
          <!-- <img :src="userStore.avatar" class="user-avatar" /> -->
          <img src="@/assets/images/profile.jpg" class="user-avatar" />
          <span class="name">{{ userStore.name }}</span>
          <el-icon><ArrowDown /></el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <router-link to="/user/profile">
              <el-dropdown-item>{{ $t('layout.personalCenter') }}</el-dropdown-item>
            </router-link>
            <el-dropdown-item command="setLayout">
              <span>{{ $t('layout.layoutSetting') }}</span>
            </el-dropdown-item>
            <el-dropdown-item command="copyToken">
              <span>Copy Token</span>
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <span>{{ $t('layout.logOut') }}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import Breadcrumb from '@/components/Breadcrumb'
import TopNav from '@/components/TopNav'
import Hamburger from '@/components/Hamburger'
import Screenfull from '@/components/Screenfull'
import SizeSelect from '@/components/SizeSelect'
import HeaderSearch from '@/components/HeaderSearch'
import ZrGit from '@/components/Zr/Git'
import ZrDoc from '@/components/Zr/Doc'
import Notice from '@/components/Notice/Index'
import LangSelect from '@/components/LangSelect/index'
import useAppStore from '@/store/modules/app'
import useUserStore from '@/store/modules/user'
import useSettingsStore from '@/store/modules/settings'
import { useClipboard } from '@vueuse/core'
const { proxy } = getCurrentInstance()
const appStore = useAppStore()
const userStore = useUserStore()
const settingsStore = useSettingsStore()

const sideTheme = computed(() => settingsStore.sideTheme)
function toggleSideBar() {
  appStore.toggleSideBar()
}

function handleCommand(command) {
  switch (command) {
    case 'setLayout':
      setLayout()
      break
    case 'logout':
      logout()
      break
    case 'copyToken':
      copyText(userStore.token)
      break
    default:
      break
  }
}

const { copy, isSupported } = useClipboard()
const copyText = async (val) => {
  if (isSupported) {
    copy(val)
    proxy.$modal.msgSuccess('复制成功！')
  } else {
    alert(val)
    proxy.$modal.msgError('Current browser does not support')
  }
}
function logout() {
  proxy
    .$confirm(proxy.$t('layout.logOutConfirm'), proxy.$t('common.tips'), {
      confirmButtonText: 'Sure',
      cancelButtonText: 'Cancel',
      type: 'warning'
    })
    .then(() => {
      userStore.logOut().then(() => {
        location.href = import.meta.env.VITE_APP_ROUTER_PREFIX + 'index'
      })
    })
    .catch(() => {})
}

const emits = defineEmits(['setLayout'])
function setLayout() {
  emits('setLayout')
}
</script>

<style lang="scss" scoped>
.el-menu {
  // display: inline-table;
  border-bottom: none;
  .el-menu-item {
    vertical-align: center;
    font-size: 1rem !important;
  }
}
.navbar {
  margin-top: -2px;
  height: var(--base-header-height);
  line-height: var(--base-header-height);
  overflow: hidden;
  position: relative;
  background: linear-gradient(50deg, var(--el-color-primary), #ffffff);

  background: var(--base-topBar-background);
  //background: var(--el-color-primary);
  box-shadow: 2px 1px 4px var(--el-color-primary);

  .hamburger-container {
    line-height: var(--base-header-height);
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }
  .topmenu-container {
    position: absolute;
    left: 50px;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    display: flex;
    justify-content: flex-end;
    align-items: center;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      padding: 0 8px;
      color: var(--base-topBar-color);
      vertical-align: text-bottom;
    }

    .avatar-container {
      .avatar-wrapper {
        display: flex;
        align-items: center;
        .user-avatar {
          cursor: pointer;
          width: 30px;
          height: 30px;
          border-radius: 50%;
          vertical-align: middle;
          margin-right: 5px;
        }
        .name {
          font-size: 12px;
        }
        i {
          cursor: pointer;
          margin-left: 10px;
        }
      }
    }
  }
}
</style>
