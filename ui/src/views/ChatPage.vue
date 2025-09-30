<script setup>
import Sidebar from '@/components/sidebar.vue'
import Chat from '@/components/Chat.vue'
</script>

<template>
  <div class="chat-page-container">
    <Sidebar />
    <div class="chat-wrapper">
      <Chat />
    </div>
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

.chat-wrapper {
  flex: 1;
  min-width: 0; /* Prevents overflow in flex items */
  display: grid; /* Use grid to make child fill the space */
  transition: all 0.2s ease-out;
}
</style>