<template>
  <div class="sidebar-logo-container" :class="{ collapse: collapse }">
    <transition name="sidebarLogoFade">
      <router-link v-if="collapse" key="collapse" class="sidebar-logo-link" to="/">
        <img v-if="logo" :src="logo" class="sidebar-logo" />
        <SvgLogo color="#0d4794" width="500" height="300" style="margin-left: 20px; margin-top: -10px" />

        <!-- <h1 v-else class="sidebar-title">{{ title }}</h1> -->
      </router-link>
      <router-link v-else key="expand" class="sidebar-logo-link" to="/">
        <SvgLogo color="#0d4794" width="500" height="300" style="margin-left: 20px; margin-top: -10px" />
        <!--   <img v-if="logo" :src="SvgLogo" class="sidebar-logo" />     <h1 class="sidebar-title">SIT Admin ({{ companyCode }})</h1> -->
      </router-link>
    </transition>
  </div>
</template>

<script setup>
import logo from '@/assets/logo/logo.svg'
import SvgLogo from './svg.vue'
import useSettingsStore from '@/store/modules/settings'
import { getUserInfo } from '@/utils/auth'
const { companyCode } = JSON.parse(getUserInfo())

defineProps({
  collapse: {
    type: Boolean,
    required: true
  }
})

const title = ref(import.meta.env.VITE_APP_TITLE)
const settingsStore = useSettingsStore()
const sideTheme = computed(() => settingsStore.sideTheme)
const theme = computed(() => settingsStore.theme)
</script>

<style lang="scss" scoped>
.sidebarLogoFade-enter-active {
  transition: opacity 1.5s;
}

.sidebarLogoFade-enter,
.sidebarLogoFade-leave-to {
  opacity: 0;
}

.sidebar-logo-container {
  position: relative;
  width: 100%;
  height: 50px;
  line-height: 50px;
  background: var(--base-menu-background);
  background-color: var(--el-color-primary);
  box-shadow: 2px 1px 4px var(--el-color-primary);
  border-right: 1px solid #f7f7f7a2;
  text-align: center;
  overflow: hidden;
  color: white;

  & .sidebar-logo-link {
    height: 100%;
    width: 100%;

    & .sidebar-logo {
      width: 100px;
      height: 35px;
      vertical-align: middle;
      margin-right: 0px;
      right: 0;
    }

    & .sidebar-title {
      display: inline-block;
      margin: 0;
      color: var(--base-logo-title-color);
      font-weight: 600;
      line-height: 50px;
      font-size: 14px;
      font-family: Avenir, Helvetica Neue, Arial, Helvetica, sans-serif;
      vertical-align: middle;
    }
  }

  &.collapse {
    .sidebar-logo {
      margin-right: 0px;
    }
  }
}
</style>
