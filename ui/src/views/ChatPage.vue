<script setup>
import Sidebar from '@/components/sidebar.vue'
import Chat from '@/components/Chat.vue'
</script>

<template>
  <div id="app-container">
    <Sidebar/>
    <Chat/>
  </div>
</template>

<style scoped>
#app-container {
  display: flex;
  height: 100vh;
  background-color: #1e1e1e;
  width: 100vw;
}
</style>
