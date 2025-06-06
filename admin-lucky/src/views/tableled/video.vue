<template>
  <div class="app-container">
    <el-row :gutter="24">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item label="Table No:">
          <el-select v-model="form.tableNo" placeholder="table" @change="selectTable" style="width: 200px" filterable>
            <el-option v-for="item in videoList" :key="item.videoURL" :label="item.name" :value="item.videoURL" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="playVideo">Play PC</el-button>
          <el-button type="primary" @click="playVideoH5">Play H5</el-button>
          <el-button type="danger" @click="stopVideo">Stop Video</el-button>
        </el-form-item>
        <el-form-item>
          <div>Bogota(UTC-5) : {{ formattedDateTime }}</div>
        </el-form-item>
      </el-form>
    </el-row>
    <div>
      <vue-flv-player
        ref="flvPlayer"
        :hasAudio="true"
        controls
        :muted="false"
        :source="currentSource"
        :isLive="true"
        type="flv"
        style="margin-top: 0px; border-radius: 10px; background-color: white; box-shadow: 0 2px 3px 0 rgba(0, 0, 0, 0.1); transition: 0.3s" />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, toRefs, onMounted } from 'vue'
import { getSelectTables } from '@/api/tables/table-led'
import VueFlvPlayer from '@/components/vueFlv.vue'
const formattedDateTime = ref(null)
const videoList = ref([])
const data = reactive({
  form: {
    tableNo: '',
    signKey: 'GsRi8CMEYQeyBQNf',
    streamKey: ''
  }
})
const { form } = toRefs(data)
const currentSource = ref('') // Reactive variable for the current video source
const currentSourceH5 = ref('') // Reactive variable for the current video source
const currentSourcePc = ref('') // Reactive variable for the current video sour
const flvPlayer = ref(null) // Create a reference for the FLV player

// Fetch the list of videos
async function getList() {
  try {
    const res = await getSelectTables()
    const filteredData = res.data.filter((item) => item.videoURL !== null && item.videoURL !== '')

    console.log(filteredData)
    videoList.value = filteredData
  } catch (error) {
    console.error('Error fetching video list:', error)
  }
}
function updateDateTime() {
  const now = new Date()

  // Convert time to Bogota timezone (UTC-5)
  const bogotaTime = new Date(now.toLocaleString('en-US', { timeZone: 'America/Bogota' }))

  // Format the date as dd/mm/yyyy
  const formattedDate = bogotaTime.toLocaleDateString('en-GB')

  // Format the time as HH:MM:SS
  const formattedTime = bogotaTime.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false // 24-hour format
  })

  // Combine date and time into the final format
  formattedDateTime.value = `${formattedDate} ${formattedTime}`
}

// Play the selected video
function playVideo() {
  stopVideo()
  currentSource.value = currentSourcePc.value
  if (flvPlayer.value && currentSource.value) {
    flvPlayer.value.play() // Call the play method
  }
}
function playVideoH5() {
  stopVideo()
  currentSource.value = currentSourceH5.value
  if (flvPlayer.value && currentSource.value) {
    flvPlayer.value.play() // Call the play method
  }
  console.log(currentSource.value)
}
// Stop the video
function stopVideo() {
  if (flvPlayer.value) {
    flvPlayer.value.pause() // Call the pause method
  }
}

// Select table and set the video source
function selectTable(videoURL) {
  currentSourcePc.value = `wss://nms.devops.wine/live/${videoURL}.flv` //; // Set the current source to the selected video URL
  currentSourceH5.value = `wss://nms.devops.wine/live/${videoURL}-m.flv`
  playVideo() // Automatically play the selected video
}
// Fetch the video list when the component is mounted
onMounted(() => {
  updateDateTime()
  setInterval(updateDateTime, 1000)
  getList()
})
</script>

<style>
/* Add any additional styles here */
</style>
