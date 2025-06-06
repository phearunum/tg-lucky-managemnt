<template>
  <video :width="width" :height="height" class="flv-player" v-bind="$attrs" ref="flvPlayer">
    <p>Your browser does not support the video tag.</p>
  </video>
</template>

<script>
import flvjs from 'flv.js'

export default {
  name: 'vue-flv-player',

  props: {
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
      type: [String, Number],
      default: '100%'
    },
    height: {
      type: [String, Number],
      default: 'auto'
    },
    mediaDataSource: {
      type: Object,
      default: () => ({})
    },
    config: {
      type: Object,
      default: () => ({})
    },
    isLive: {
      type: Boolean,
      default: true
    },
    enableStashBuffer: {
      type: Boolean,
      default: true
    },
    hasAudio: {
      type: Boolean,
      default: true
    }
  },

  data() {
    return {
      isSupported: false,
      flvPlayer: null
    }
  },

  watch: {
    source: 'init'
  },

  methods: {
    init() {
      if (this.isSupported && this.source) {
        const videoElement = this.$refs.flvPlayer

        this.flvPlayer = flvjs.createPlayer(
          {
            url: this.source,
            type: this.type,
            isLive: this.isLive,
            enableStashBuffer: this.enableStashBuffer,
            hasAudio: this.hasAudio,
            ...this.mediaDataSource
          },
          this.config
        )

        this.flvPlayer.attachMediaElement(videoElement)
        this.load()

        // Add event listeners for error handling
        this.flvPlayer.on(flvjs.Events.ERROR, (errorType, errorDetails) => {
          console.error(`FLV Error: ${errorType} - ${errorDetails}`)
          // Implement reconnection logic if necessary
        })
      }
    },
    load() {
      this.flvPlayer.load()
    },
    play() {
      this.flvPlayer.play()
    },
    pause() {
      this.flvPlayer.pause()
    },
    destroy() {
      if (this.flvPlayer) {
        this.flvPlayer.destroy()
        this.flvPlayer = null
      }
    }
  },

  created() {
    this.isSupported = flvjs.isSupported()
  },

  mounted() {
    this.init()
  },

  beforeDestroy() {
    this.destroy()
  }
}
</script>

<style lang="scss"></style>
