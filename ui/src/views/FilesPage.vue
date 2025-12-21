<script setup>
import { useUiStore } from '@/stores/ui';
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from '@/components/sidebar.vue';
import { useNotify } from '@/composables/notify';

const uiStore = useUiStore();
const router = useRouter();

const loading = ref(true);
const isUploading = ref(false);
const fileInput = ref(null);
const files = ref([]);
const totalItems = ref(0);
const currentPage = ref(1);
const limit = ref(10);
const offset = ref(0);

const hasMore = computed(() => maxIndex.value < totalItems.value);
const maxIndex = computed(() => (currentPage.value - 1) * limit.value + files.value.length);

function triggerFileUpload() {
    fileInput.value.click();
}

async function uploadFile(event) {
    const selectedFiles = event.target.files;
    if (!selectedFiles.length) return;

    isUploading.value = true;

    const uploadPromises = Array.from(selectedFiles).map(file => {
        const formData = new FormData();
        formData.append('file', file);
        return fetch('/api/uploadFile', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        }).then(async response => {
            if (!response.ok) {
                const error = await response.text();
                throw new Error(`Failed to upload ${file.name}: ${error}`);
            }
        });
    });

    try {
        await Promise.all(uploadPromises);
        useNotify(`${selectedFiles.length} file(s) uploaded successfully!`);
        
        // Go to first page and refresh
        offset.value = 0;
        currentPage.value = 1;
        await getFiles();

    } catch (error) {
        useNotify(`Error uploading file(s): ${error.message}`);
    } finally {
        isUploading.value = false;
        if (fileInput.value) {
            fileInput.value.value = '';
        }
    }
};

async function getFiles() {
    loading.value = true;
    const params = new URLSearchParams();
    params.append('limit', limit.value);
    params.append('offset', offset.value);

    const url = `/api/getFileList?${params}`;

    try {
        const response = await fetch(url, { credentials: 'include' });

        if (!response.ok) {
            if (response.status === 401) {
                await router.replace('/login');
            } else {
                useNotify(`Error: ${response.status} ${response.statusText}`);
            }
            return;
        }

        const data = await response.json();

        if (data.files && Array.isArray(data.files)) {
            files.value = data.files;
        }
        if (data.totalItems && typeof data.totalItems === 'number') {
            totalItems.value = data.totalItems;
        }
    } catch (e) {
        useNotify('Could not fetch files.');
        console.error(e);
    } finally {
        loading.value = false;
    }
}

async function deleteFile(fileId) {
    if (!confirm('Are you sure you want to delete this file?')) {
        return;
    }

    const params = new URLSearchParams();
    params.append('fileId', fileId);

    const url = `/api/deleteFile?${params}`;

    try {
        const response = await fetch(url, {
            method: 'POST',
            credentials: 'include',
        });

        if (!response.ok) {
            if (response.status === 401) {
                await router.replace('/login');
            } else {
                useNotify(`Error: ${response.status} ${response.statusText}`);
            }
            return;
        }

        useNotify('File deleted successfully.');
        
        totalItems.value--;
        if (files.value.length === 1 && currentPage.value > 1) {
            await prevPage();
        } else if (totalItems.value === 0) {
            files.value = [];
        } else {
            await getFiles();
        }

    } catch (e) {
        useNotify('Could not delete file.');
        console.error(e);
    }
}

function viewFile(fileId) {
    const url = `/api/getFile?fileId=${fileId}`;
    window.open(url, '_blank');
}

async function nextPage() {
    if (!hasMore.value) return;
    offset.value += limit.value;
    currentPage.value += 1;
    await getFiles();
}

async function prevPage() {
    if (offset.value === 0) return;
    offset.value = Math.max(offset.value - limit.value, 0);
    currentPage.value -= 1;
    await getFiles();
}

const formatTimestamp = (timestamp) => {
    if (!timestamp) return 'N/A';
    return new Date(timestamp).toLocaleString();
}


onMounted(async () => {
    await getFiles();
});
</script>

<template>
    <div class="page-container">
        <Sidebar />
        <main class="main-content" :style="{ marginLeft: uiStore.sidebarWidth }">
            <input type="file" ref="fileInput" @change="uploadFile" style="display: none" multiple />
            <header class="page-header">
                <h1>Files</h1>
                <button class="action-button" @click="triggerFileUpload" :disabled="isUploading">
                    <i v-if="isUploading" class="fa-solid fa-spinner fa-spin"></i>
                    <i v-else class="fa-solid fa-plus"></i>
                    <span>Upload File</span>
                </button>
            </header>

            <div v-if="loading" class="loading-spinner">
                <i class="fa-solid fa-spinner fa-spin"></i>
            </div>

            <div v-else-if="files.length > 0" class="table-container">
                <table class="data-table">
                    <thead>
                        <tr>
                            <th>File Name</th>
                            <th>Type</th>
                            <th>Created</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="file in files" :key="file.id">
                            <td>{{ file.file_name }}</td>
                            <td>{{ file.content_type }}</td>
                            <td>{{ formatTimestamp(file.created_at) }}</td>
                            <td class="action-buttons">
                                <button class="action-btn" @click="viewFile(file.id)">
                                    <i class="fa-solid fa-eye"></i>
                                </button>
                                <button class="action-btn delete-btn" @click="deleteFile(file.id)">
                                    <i class="fa-solid fa-trash"></i>
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div v-else class="no-data">
                <p>No Files Found</p>
            </div>

            <div v-if="!loading && files.length > 0" class="pagination-controls">
                <button class="pagination-btn" @click="prevPage" :disabled="offset === 0">
                    <i class="fa-solid fa-arrow-left"></i>
                    <span>Prev</span>
                </button>
                <button class="pagination-btn" @click="nextPage" :disabled="!hasMore">
                    <span>Next</span>
                    <i class="fa-solid fa-arrow-right"></i>
                </button>
            </div>
        </main>
    </div>
</template>

<style scoped>
.page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  position: relative;
}

.main-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    padding: 2rem 3rem;
    transition: margin-left 0.3s ease;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.page-header h1 {
    font-size: 2rem;
    font-weight: 700;
    color: var(--text-color);
}

.action-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.action-button:disabled {
    cursor: not-allowed;
}

.loading-spinner {
    text-align: center;
    margin-top: 5rem;
    font-size: 2rem;
}

.table-container {
    width: 100%;
    overflow-x: auto;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
}

.data-table th, .data-table td {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
}

.data-table th {
    font-size: 0.875rem;
    text-transform: uppercase;
    font-weight: 600;
    color: var(--text-color);
    opacity: 0.7;
}

.data-table tbody tr:hover {
    background-color: rgba(255, 255, 255, 0.05);
}

.action-buttons {
    display: flex;
    gap: 0.5rem;
}

.action-btn {
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.5rem;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.action-btn:hover {
    background-color: var(--primary-color);
    border-color: var(--primary-color);
    color: white;
}

.delete-btn:hover {
    background-color: #ef4444;
    border-color: #ef4444;
}

.no-data {
    text-align: center;
    margin-top: 5rem;
    color: var(--text-color-secondary);
}

.pagination-controls {
    display: flex;
    justify-content: flex-end;
    margin-top: 1rem;
    gap: 0.5rem;
}

.pagination-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: none;
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.5rem 1rem;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.pagination-btn:not(:disabled):hover {
    background-color: var(--primary-color);
    border-color: var(--primary-color);
    color: white;
}

.pagination-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>
