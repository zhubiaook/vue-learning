<template>
  <div class="iframe-wrapper">
    <iframe
      class="iframe-full"
      src="https://www.yuque.com/zhubiaook/operation"
      allowfullscreen
      loading="lazy"
      title="Operation Documentation"
    ></iframe>
  </div>
</template>

<script>
export default {
  name: "OperationFrame",
};
</script>

<style scoped>
.iframe-wrapper {
  width: 100%;
  height: 100vh;
  margin: 0;
  padding: 0;
  overflow: hidden;
  display: flex;
}

.iframe-full {
  flex: 1;
  border: none;
  height: 100%;
  width: 100%;
}
</style>
