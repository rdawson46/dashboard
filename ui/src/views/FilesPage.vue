<script setup>
import { useUiStore } from '@/stores/ui';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from '@/components/sidebar.vue';
import { useNotify } from '@/composables/notify';

const uiStore = useUiStore();

const files = ref(null)

async function getFiles() {

}

onMounted(async () => {

})
</script>

<template>
    <div class="job-page-container">
        <Sidebar />
        <main class="jobs-main-content" :style="{ marginLeft: uiStore.sidebarWidth }">
            <header class="jobs-header">
                <h1>Files</h1>
                <button class="create-file-btn" @click="showCreateJobModal = true">
                    <i class="fa-solid fa-plus"></i>
                    <span>Upload File</span>
                </button>
            </header>

            <!-- loading -->
            <template v-if="files === null">
                <div>
                    <i  class="fa-solid fa-spinner fa-spin"></i>
                </div>
            </template>

            <!-- array -->
            <template v-else-if="Array.isArray(files) && files.length > 0"></template>

            <!-- empty array -->
            <template v-else-if="Array.isArray(files) && files.length === 0">
                <div>
                    <p>No Files Found</p>
                </div>
            </template>


        </main>
    </div>
</template>

<style scoped>
.job-page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  position: relative;
}

.jobs-main-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    padding: 2rem 3rem;
    transition: margin-left 0.3s ease;
    border-radius: 0;
}

.jobs-main-content h1 {
  margin-bottom: 2rem;
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-color);
}

.create-file-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}
</style>
