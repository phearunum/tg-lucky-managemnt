<template>
  <video :width="width" :height="height" class="flv-player" ref="flvPlayer">
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
    width: {
      type: Number,
      default: '100%'
    },
    height: {
      type: Number,
      default: 'auto'
    },
    isLive: {
      type: Boolean,
      required: true
    },
    hasAudio: {
      type: Boolean,
      required: true
    }
  },

  data() {
    return {
      isSupported: false,
      flvPlayer: null
    }
  },

  watch: {
    source(newSource) {
      if (newSource) {
        this.initPlayer(newSource)
      }
    }
  },

  methods: {
    initPlayer(source) {
      if (flvjs.isSupported()) {
        this.flvPlayer = flvjs.createPlayer({
          url: source,
          type: 'flv',
          isLive: this.isLive,
          hasAudio: this.hasAudio
        })
        this.flvPlayer.attachMediaElement(this.$refs.flvPlayer)
        this.flvPlayer.load()
        this.flvPlayer.play()
      }
    },
    play() {
      if (this.flvPlayer) {
        this.flvPlayer.play()
      }
    },
    pause() {
      if (this.flvPlayer) {
        this.flvPlayer.pause()
      }
    },
    destroy() {
      if (this.flvPlayer) {
        this.flvPlayer.destroy()
        this.flvPlayer = null
      }
    }
  },

  mounted() {
    this.isSupported = flvjs.isSupported()
    if (this.source) {
      this.initPlayer(this.source)
    }
  },

  beforeDestroy() {
    this.destroy()
  }
}
</script>

<style lang="scss">
/* Add styles as needed */
</style>
