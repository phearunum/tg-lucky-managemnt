<template>
  <div>
    <svg-icon name="github" @click="goto" />
  </div>
</template>
<script>
export default {
  setup() {
    const url = ref("/#");

    function goto() {
      window.open(url.value);
    }

    return {
      goto,
    };
  },
};
</script>
