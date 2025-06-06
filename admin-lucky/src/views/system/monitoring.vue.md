<template>
  <div class="app-container">
    <h1>Monitoring</h1>
    <button @click="myFunction(500, 800)">Open</button>
    <input v-model="iframeSrc" placeholder="Enter URL" />
    <button @click="loadIframe">Load URL</button>
    <button @click="openInNewTab">Open in New Tab</button>
    <iframe v-if="iframeLoaded" :src="iframeSrc" @load="onIframeLoad" @error="onIframeError"></iframe>
  </div>
</template>

<script setup name="user">
import { ref, getCurrentInstance } from 'vue'

const iframeSrc = ref('')
const iframeLoaded = ref(false)

const instance = getCurrentInstance()
const { proxy } = instance
function myFunction(width, height) {
  var top = parseInt(screen.availHeight - height - 100)
  var left = parseInt(screen.availWidth - width / 2)
  var features =
    'location=1, status=1,toolbar=yes,scrollbars=yes,resizable=0,menubar=yes, width=' +
    width +
    ', height=' +
    height +
    ', top=' +
    top +
    ', left=' +
    left
  window.open('http://192.168.26.24:3000/', 'App', features)
  window.close()
}
const loadIframe = () => {
  iframeLoaded.value = true
}

const openInNewTab = () => {
  if (iframeSrc.value) {
    window.open(iframeSrc.value, '_blank')
  }
}

const onIframeLoad = () => {
  console.log('Iframe loaded successfully')
}

const onIframeError = () => {
  console.log('Failed to load iframe')
}
</script>

<style scoped>
html,
body,
.app-container {
  height: 100%;
  margin: 0;
}

.app-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 20px;
}

input {
  width: 80%;
  padding: 10px;
  margin: 10px 0;
}

button {
  padding: 10px 20px;
  margin: 10px 0;
}

iframe {
  width: 100%;
  height: 100%;
  border: none;
  position: absolute;
  top: 0;
  left: 0;
}
</style>
