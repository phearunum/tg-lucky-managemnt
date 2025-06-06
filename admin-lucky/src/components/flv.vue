<template>
  <div>
    <video :width="width" :height="height" class="flv-player" ref="flvPlayer">
      <p>Your browser does not support the video tag.</p>
    </video>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import flvjs from 'flv.js'

const props = defineProps({
  source: {
    type: String,
    required: true
  },
  type: {
    type: String,
    default: 'flv',
    required: true
  },
  width: {
    type: Number,
    default: 800
  },
  height: {
    type: Number,
    default: 600
  },
  isLive: {
    type: Boolean,
    default: false
  }
})

const flvPlayer = ref(null)
const videoElement = ref(null)

onMounted(() => {
  if (flvjs.isSupported()) {
    flvPlayer.value = flvjs.createPlayer({
      url: props.source,
      type: props.type,
      isLive: props.isLive
    })
    flvPlayer.value.attachMediaElement(videoElement.value)
    flvPlayer.value.load()
    flvPlayer.value.play()
  }
})

onBeforeUnmount(() => {
  if (flvPlayer.value) {
    flvPlayer.value.unload()
    flvPlayer.value.destroy()
    flvPlayer.value = null
  }
})
</script>

<style scoped>
.flv-player {
  width: 100%;
  height: auto;
}
</style>
