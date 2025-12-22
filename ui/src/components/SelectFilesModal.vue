<script setup>
import { ref, onMounted, defineEmits, defineProps } from 'vue';
import { useNotify } from '@/composables/notify.js';

const props = defineProps({
  show: Boolean,
  selectedFiles: Array
});

const emit = defineEmits(['close', 'add-files']);

const files = ref([]);
const localSelectedFiles = ref([]);

async function fetchFiles() {
  try {
    const response = await fetch('/api/getFileList', { credentials: 'include' });
    if (!response.ok) {
      throw new Error('Failed to fetch files.');
    }
    const data = await response.json();
    files.value = data.files;
  } catch (error) {
    useNotify(`Error: ${error.message}`, 'error');
  }
}

function addFiles() {
  emit('add-files', localSelectedFiles.value);
  emit('close');
}

function isSelected(file) {
  return localSelectedFiles.value.some(selected => selected.id === file.id);
}

function toggleSelection(file) {
  if (isSelected(file)) {
    localSelectedFiles.value = localSelectedFiles.value.filter(selected => selected.id !== file.id);
  } else {
    localSelectedFiles.value.push(file);
  }
}

onMounted(() => {
    fetchFiles(); 
    localSelectedFiles.value = [...props.selectedFiles];
});
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content glass-card">
      <h2>Select Files</h2>

      <div class="files-list">

        <div v-for="file in files" :key="file.id" class="file-item" :class="{ selected: isSelected(file) }" @click="toggleSelection(file)">
          <i class="fa-solid fa-file"></i>
          <span>{{ file.file_name }}</span>
        </div>

      </div>

      <div class="modal-actions">
        <button @click="$emit('close')">Cancel</button>
        <button @click="addFiles" class="primary">Add Files</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  width: 500px;
  max-width: 90%;
  padding: 2rem;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.files-list {
  max-height: 400px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.file-item:hover {
  background-color: var(--bg-color-light);
}

.file-item.selected {
  background-color: var(--primary-color);
  color: white;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1rem;
}

.modal-actions button {
    background: var(--bg-color-light);
    color: var(--text-color);
    border: 1px solid var(--border-color);
    padding: 0.5rem 1rem;
    border-radius: 8px;
    cursor: pointer;
}

.modal-actions button.primary {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
}
</style>
