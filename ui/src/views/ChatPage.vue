<script setup>
import Sidebar from '@/components/sidebar.vue'
import Chat from '@/components/Chat.vue'
import { useUiStore } from '@/stores/ui';

const uiStore = useUiStore();
</script>

<template>
  <div class="chat-page-container">
    <Sidebar />
    <div class="chat-content" :style="{ marginLeft: uiStore.sidebarWidth }">
      <Chat />
    </div>
  </div>
</template>

<style scoped>
.chat-page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
}

.chat-content {
  width: 100%;
  height: 100%;
  padding: 1.5rem;
  transition: margin-left 0.3s ease;
}
</style>
