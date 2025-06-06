<template>
  <div>
    <el-dropdown trigger="hover" @command="handleLanguageChange" style="vertical-align: middle">
      <svg-icon class-name="size-icon" name="language" />
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-for="item of langOptions" :key="item.value" :disabled="lang === item.value" :command="item.value">
            {{ item.label }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import useAppStore from '@/store/modules/app'
const appStore = useAppStore()
const lang = computed(() => useAppStore().lang)

const { proxy } = getCurrentInstance()
const langOptions = ref([
  { label: 'Khmer', value: 'km' },
  { label: 'English', value: 'en' },
  { label: '简体中文', value: 'zh-cn' },
 // { label: '繁體中文', value: 'zh-tw' },
])

function handleLanguageChange(lang) {
  proxy.$modal.loading('Please wait, language setting...')
  appStore.setLang(lang)
  setTimeout('window.location.reload()', 1000)
}
</script>

<style lang="scss" scoped>
.size-icon--style {
  font-size: 18px;
  line-height: 50px;
  padding-right: 7px;
}
</style>
