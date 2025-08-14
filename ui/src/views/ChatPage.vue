<script setup>
import Sidebar from '@/components/sidebar.vue'
import Chat from '@/components/Chat.vue'
</script>

<template>
  <div class="chat-page-container">
    <Sidebar />
    <Chat />
  </div>
</template>

<style scoped>
.chat-page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  padding: 1.5rem;
  gap: 1.5rem;
}
</style>
