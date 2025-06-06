<template>
  <div>
    <svg-icon name="question" @click="goto"/>
  </div>
</template>
<script setup>
const url = ref('#');

function goto() {
  window.open(url.value)
}
</script>
