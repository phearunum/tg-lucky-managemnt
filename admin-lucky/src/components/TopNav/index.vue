<template>
  <el-menu
    :default-active="activeMenu"
    :active-text-color="theme"
    mode="horizontal"
    background-color="transparent"
    @select="handleSelect"
    :ellipsis="false">
    <template v-for="(item, index) in topMenus">
      <el-menu-item :style="{ '--theme': theme }" :index="item.path" :key="index" v-if="index < visibleNumber">
        <svg-icon :name="item.meta.icon" />
        <!-- {{ item.meta.title }} -->
        <template v-if="item.meta.titleKey" #title> {{ $t(item.meta.titleKey) }} </template>
        <template v-else-if="item.meta.title" #title>
          {{ item.meta.title }}
        </template>
      </el-menu-item>
    </template>

    <!-- Top menu collapsed beyond quantity -->
    <el-sub-menu :style="{ '--theme': theme }" index="more" v-if="topMenus.length > visibleNumber">
      <template #title>{{ $t('btn.more') }}</template>
      <template v-for="(item, index) in topMenus">
        <el-menu-item :index="item.path" :key="index" v-if="index >= visibleNumber">
          <svg-icon :name="item.meta.icon" />
          <span style="margin-left: 5px">{{ item.meta.title }}</span>
        </el-menu-item>
      </template>
    </el-sub-menu>
  </el-menu>
</template>

<script setup>
import { constantRoutes } from '@/router'
import { isHttp } from '@/utils/validate'
import { useRouter } from 'vue-router'
import { getNormalPath } from '@/utils/ruoyi'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'
import usePermissionStore from '@/store/modules/permission'

// Initial number of top bar
const visibleNumber = ref(5)
// Whether it is the first load
const isFrist = ref(false)
// index of the currently active menu
const currentIndex = ref(undefined)

const appStore = useAppStore()
const settingsStore = useSettingsStore()
const permissionStore = usePermissionStore()
const route = useRoute()
const router = useRouter()

// theme color
const theme = computed(() => settingsStore.theme)
// all routing information
const routers = computed(() => permissionStore.topbarRouters)

// Display the menu at the top
const topMenus = computed(() => {
  let topMenus = []
  routers.value.map((menu) => {
    if (menu.hidden !== true) {
      // Compatible with the internal jump of the first-level menu on the top bar
      if (menu.path === '/') {
        topMenus.push(menu.children[0])
      } else {
        topMenus.push(menu)
      }
    }
  })
  return topMenus
})

// set sub-routes
const childrenMenus = computed(() => {
  let childrenMenus = []
  routers.value.map((router) => {
    for (let item in router.children) {
      if (router.children[item].parentPath === undefined) {
        if (router.path === '/') {
          router.children[item].path = getNormalPath('/redirect/' + router.children[item].path)
        } else {
          if (!isHttp(router.children[item].path)) {
            router.children[item].path = getNormalPath(router.path + '/' + router.children[item].path)
          }
        }
        router.children[item].parentPath = router.path
      }
      childrenMenus.push(router.children[item])
    }
  })
  return constantRoutes.concat(childrenMenus)
})

const activeMenu = computed(() => {
  const path = route.path
  let activePath = defaultRouter.value
  if (path.lastIndexOf('/') > 0) {
    const tmpPath = path.substring(1, path.length)
    activePath = '/' + tmpPath.substring(0, tmpPath.indexOf('/'))
  } else if ('/index' == path || '' == path) {
    if (!isFrist.value) {
      isFrist.value = true
    } else {
      activePath = 'index'
    }
  }
  let routes = activeRoutes(activePath)
  if (routes.length === 0) {
    activePath = currentIndex.value || defaultRouter.value

    activeRoutes(activePath)
  }
  return activePath
})

const defaultRouter = computed(() => {
  let router
  Object.keys(routers.value).some((key) => {
    if (!routers.value[key].hidden) {
      router = routers.value[key].path
      return true
    }
  })

  return router
})
function setVisibleNumber() {
  const width = document.body.getBoundingClientRect().width / 3
  visibleNumber.value = parseInt(width / 85)
}
function handleSelect(key, keyPath) {
  currentIndex.value = key
  if (isHttp(key)) {
    window.open(key, '_blank')
  } else if (key.indexOf('/redirect') !== -1) {
    router.push({ path: key.replace('/redirect', '') }).catch((err) => {})
  } else {
    activeRoutes(key)
  }
}

function activeRoutes(key) {
  var routes = []
  if (childrenMenus.value && childrenMenus.value.length > 0) {
    childrenMenus.value.map((item) => {
      if (key == item.parentPath || (key == 'index' && '' == item.path)) {
        routes.push(item)
      }
    })
  }
  if (routes.length > 0) {
    permissionStore.setSidebarRouters(routes)
  }
  return routes
}

onMounted(() => {
  window.addEventListener('resize', setVisibleNumber)
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', setVisibleNumber)
})

onMounted(() => {
  setVisibleNumber()
})
</script>

<style lang="scss">
// 修改默认样式
.topmenu-container.el-menu--horizontal > .el-menu-item {
  height: 40px !important;
  line-height: 50px !important;
  color: #999093 !important;
  padding: 0 5px !important;
  margin: 0 10px !important;
}
.el-menu--horizontal > .el-menu-item .svg-icon {
  margin-right: 5px;
}
/* sub-menu item */
.topmenu-container.el-menu--horizontal > .el-sub-menu .el-sub-menu__title {
  height: 40px !important;
  line-height: 40px !important;
  color: #999093 !important;
  padding: 0 5px !important;
  margin: 0 10px !important;

}
</style>
