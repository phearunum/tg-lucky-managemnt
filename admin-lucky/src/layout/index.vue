<template>
  <el-container :class="classObj" class="app-layout" :style="{ '--current-color': theme }">
    <!-- 移动端打开菜单遮罩 -->
    <el-drawer v-if="device === 'mobile'" :size="220" v-model="menuDrawer" :with-header="false" modal-class="sidebar-mobile" direction="ltr">
      <sidebar />
    </el-drawer>
    <sidebar v-else-if="!sidebar.hide" />

    <el-container class="main-container flex-center" :class="{ hasTagsView: needTagsView, sidebarHide: sidebar.hide }">
      <el-header :class="{ 'fixed-header': fixedHeader }">
        <navbar @setLayout="setLayout" />
        <tags-view v-if="needTagsView" />
      </el-header>
      <el-main class="app-main" :style="wallpaper">
        <router-view v-slot="{ Component, route }">
          <transition name="fade-transform" mode="out-in">
            <keep-alive :include="cachedViews">
              <component v-if="!route.meta.link" :is="Component" :key="route.path" />
            </keep-alive>
          </transition>
        </router-view>
        <iframe-toggle />
      </el-main>
      <el-footer v-if="showFooter">
        <div v-html="defaultSettings.copyright"></div>
      </el-footer>
      <settings ref="settingRef" />
    </el-container>
  </el-container>
</template>

<script setup>
import { useWindowSize } from '@vueuse/core'
import Sidebar from './components/Sidebar/index.vue'
import { Navbar, Settings, TagsView } from './components'
import defaultSettings from '@/settings'
import iframeToggle from './components/IframeToggle/index'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'
import useTagsViewStore from '@/store/modules/tagsView'
import Cookies from 'js-cookie'
const wallpaper = {
  'background-image': `url(${Cookies.get('wallpaper') || '/static/images/background.jpeg'})`
}
const menuDrawer = computed({
  get: () => useAppStore().sidebar.opened,
  set: (val) => {
    useAppStore().toggleSideBar(val)
  }
})
const settingsStore = useSettingsStore()
const theme = computed(() => settingsStore.theme)
const sidebar = computed(() => useAppStore().sidebar)
const device = computed(() => useAppStore().device)
const needTagsView = computed(() => settingsStore.tagsView)
const fixedHeader = computed(() => settingsStore.fixedHeader)
const showFooter = computed(() => settingsStore.showFooter)

// appMain 模块 start
const route = useRoute()
useTagsViewStore().addCachedView(route)
const cachedViews = computed(() => {
  return useTagsViewStore().cachedViews
})
//appMain 模块结束

const classObj = computed(() => ({
  hideSidebar: !sidebar.value.opened,
  openSidebar: sidebar.value.opened,
  mobile: device.value === 'mobile'
}))

const { width, height } = useWindowSize()
const WIDTH = 792 // refer to Bootstrap's responsive design

watchEffect(() => {
  if (device.value === 'mobile' && sidebar.value.opened) {
    // useAppStore().closeSideBar()
  }
  if (width.value - 1 < WIDTH) {
    useAppStore().toggleDevice('mobile')
    // useAppStore().closeSideBar()
  } else {
    useAppStore().toggleDevice('desktop')
  }
})

const settingRef = ref(null)
function setLayout() {
  settingRef.value.openSetting()
}
</script>

<style lang="scss">
@import '@/assets/styles/mixin.scss';

.main-container {
  min-height: 100%;
  width: 100%;
  flex-direction: column;
  position: relative;
}

.app-layout {
  @include clearfix;
  // position: relative;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;
  flex: 1;

  &.mobile.openSidebar {
    position: fixed;
    top: 0;
  }
}

// 固定header
.fixed-header {
  position: sticky;
  position: -webkit-sticky;
  z-index: 9;
}

.mobile .fixed-header {
  width: 100%;
}
.app-main {
  /* 50= navbar  50  */
  // min-height: calc(100vh - 50px);
  width: 100%;
  position: relative;
  height: 100%;
  overflow-x: hidden;
  //background-image: url('https://e0.pxfuel.com/wallpapers/196/899/desktop-wallpaper-for-computers-white-lotus-flower.jpg');
  min-height: calc(100vh - 100px);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-attachment: fixed;
}
.sidebar-mobile {
  .el-drawer__body {
    padding: 0;
  }
  @media screen and (max-width: 700px) {
    .el-drawer {
      width: var(--base-sidebar-width) !important;
    }
  }
}

.el-header {
  --el-header-padding: 0 0px !important;
  // --el-header-height: 50px !important;
}
.el-footer {
  --el-footer-height: var(--base-footer-height);
  line-height: var(--base-footer-height);
  text-align: center;
  color: #ccc;
  font-size: 14px;
  border-top: 1px solid #e7eaec;
  letter-spacing: 0.1rem;
}
.hasTagsView {
  // .app-main {
  //   min-height: calc(100vh - 84px - var(--base-footer-height));
  // }
  .el-header {
    --el-header-height: var(--el-header-height) + var(--el-tags-height) !important;
  }
}
</style>
