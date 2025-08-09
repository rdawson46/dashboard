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
  max-height: 100vh;
  max-width: 100vw;
  background: linear-gradient(
      135deg,
      rgba(24, 42, 60, 0.6),
      rgba(24, 42, 60, 0.3),
      rgba(24, 42, 60, 0),
      rgba(24, 42, 60, 0.3),
      rgba(24, 42, 60, 0.6)
  );
}
</style>
