<template>
  <div id="app-container">
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
#app-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
}

.main-content {
  width: 100%;
  height: 100%;
  overflow-y: auto;
}
</style>
