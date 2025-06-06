<template>
  <div>
    <el-button size="small" type="primary" plain @click.prevent="setBackgrounds"> Wallpaper </el-button>
  </div>
</template>

<script setup name="login">
import Cookies from 'js-cookie'
const { proxy } = getCurrentInstance()
const host = window.location.host
function setBackgrounds() {
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
      alert
      if (backg) {
        proxy.$modal.msgSuccess('Wallpaper changed successfully!')
      } else {
        proxy.$modal.msgError('Failed to change wallpaper.')
      }
      window.location.reload()
    })
    .catch((err) => {
      console.error(err)
    })
}
</script>
