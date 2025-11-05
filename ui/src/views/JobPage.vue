<script setup>
import Sidebar from '@/components/sidebar.vue'
import CreateJobModal from '@/components/CreateJobModal.vue'
import EditJobModal from '@/components/EditJobModal.vue'
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useNotify } from '@/composables/notify.js'
import { useUiStore } from '@/stores/ui';

const uiStore = useUiStore();
const router = useRouter();

const totalItems = ref(0)
const currentPage = ref(1)

const loading = ref(true)
const jobs = ref([])
const limit = ref(10)
const offset = ref(0)


const hasMore = computed(() => maxIndex.value < totalItems.value)
const maxIndex = computed(() => (currentPage.value - 1) * limit.value + jobs.value.length)

const showCreateJobModal = ref(false)
const showEditJobModal = ref(false)
const selectedJob = ref(null)

function openEditJobModal(job) {
    selectedJob.value = job
    showEditJobModal.value = true
}

async function deleteJob(id) {
    const params = new URLSearchParams()
    params.append('jobId', id)

    const url = `/api/deleteJob?${params}`

    try {
        const response = await fetch(url, {
            method: 'POST',
            credentials: 'include',
        })

        if (!response.ok) {
          if (response.status === 401) {
            await router.replace('/login')
            return
          }
          useNotify(`Error: ${response.status} ${response.statusText}`);
          return;
        }

        const data = await response.json()

        if (!data.status && data.status !== "ok") {
            return
        }

        totalItems.value--
        if (jobs.value.length === 1 && currentPage.value > 1) {
            await prevPage()
        } else {
            await getJobList()
        }

        return 
    } catch (e) {
        useNotify('Could not delete job')
    }
}

async function handleCreateJob(jobForm) {
    try {
        const res = await fetch('/api/createJob', {
            method: 'POST',
            credentials: 'include',
            body: jobForm
        })

        if (!res.ok) {
            useNotify('Unable to create job')
            return
        }

        const data = await res.json()
        totalItems.value += 1
        if (jobs.value.length < limit.value) {
            jobs.value.push(data)
        }
    } catch (e) {
        useNotify('Unable to create job')
        console.error(e)
    }
    showCreateJobModal.value = false
}

async function handleEditJob(jobForm) {
    try {
        const res = await fetch('/api/updateJob', {
            method: 'POST',
            credentials: 'include',
            body: jobForm
        })

        if (!res.ok) {
            useNotify('Unable to edit job')
            return
        }

        const updatedJob = await res.json()
        const index = jobs.value.findIndex(j => j.id === updatedJob.id)
        if (index !== -1) {
            jobs.value[index] = updatedJob
        }
    } catch (e) {
        useNotify('Unable to edit job')
        console.error(e)
    }
    showEditJobModal.value = false
}

async function getJobList() {
    const params = new URLSearchParams()
    params.append('limit', limit.value)
    params.append('offset', offset.value)

    const url = `/api/jobList?${params}`

    try {
        const response = await fetch(url, { credentials: 'include' })

        if (!response.ok) {
          if (response.status === 401) {
            await router.replace('/login')
            return
          }
          useNotify(`Error: ${response.status} ${response.statusText}`);
          return;
        }

        const data = await response.json()

        if (!data.jobs || !Array.isArray(data.jobs)) {
            return
        }

        if (data.totalItems && typeof data.totalItems === 'number') {
            totalItems.value = data.totalItems
        }

        jobs.value = data.jobs
        return 
    } catch (e) {
        console.error(e)
    }
}

async function nextPage() {
    offset.value += limit.value
    currentPage.value += 1
    getJobList()
}

async function prevPage() {
    offset.value = Math.max(offset.value - limit.value, 0)
    currentPage.value -= 1
    getJobList()
}

onMounted(async () => {
    try {
        await getJobList()
        loading.value = false
    } catch (e) {
        console.error(e)
    }
})

</script>

<template>
    <div class="job-page-container">
        <Sidebar />
        <main class="jobs-main-content glass-card" :style="{ marginLeft: uiStore.sidebarWidth }">
            <header class="jobs-header">
                <h1>Jobs</h1>
                <button class="create-job-btn" @click="showCreateJobModal = true">
                    <i class="fa-solid fa-plus"></i>
                    <span>Create Job</span>
                </button>
            </header>
            <div class="jobs-table-container">
                <table class="jobs-table">
                    <thead>
                        <tr>
                            <th>Job Name</th>
                            <th>Type</th>
                            <th>Status</th>
                            <!--<th>Date</th>-->
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-for="job in jobs">
                            <tr>
                                <td>{{job.name}}</td>
                                <td>{{job.task.task_type}}</td>
                                <td><span class="status" :class="[`status-${job.status.toLowerCase()}`]">{{job.status}}</span></td>
                                <!--<td>{{job.date}}</td>-->
                                <td class="action-buttons">
                                    <button class="action-btn"><i class="fa-solid fa-play"></i></button>
                                    <button class="action-btn" @click="openEditJobModal(job)"><i class="fa-solid fa-pencil"></i></button>
                                    <button class="action-btn delete-btn" @click="deleteJob(job.id)"><i class="fa-solid fa-trash"></i></button>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            
                <div v-if="loading">
                    <i  class="fa-solid fa-spinner fa-spin"></i>
                </div>
            </div>

            <div class="pagination-controls">
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
        <CreateJobModal 
            v-if="showCreateJobModal" 
            @close="showCreateJobModal = false"
            @create-job="handleCreateJob"
        />
        <EditJobModal
            v-if="showEditJobModal"
            :job="selectedJob"
            @close="showEditJobModal = false"
            @edit-job="handleEditJob"
        />
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
    padding: 2rem;
    transition: margin-left 0.3s ease;
    border-radius: 0;
}

.jobs-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.jobs-header h1 {
    font-size: 2rem;
    font-weight: 700;
}

.create-job-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.jobs-table-container {
    width: 100%;
    overflow-x: auto;
}

.jobs-table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
}

.jobs-table th, .jobs-table td {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
}

.jobs-table th {
    font-size: 0.875rem;
    text-transform: uppercase;
    font-weight: 600;
    color: var(--text-color);
    opacity: 0.7;
}

.jobs-table tbody tr {
    transition: background-color 0.2s ease;
}

.jobs-table tbody tr:hover {
    background-color: rgba(255, 255, 255, 0.05);
}

.status {
    padding: 0.25rem 0.75rem;
    border-radius: 9999px;
    font-weight: 500;
    font-size: 0.875rem;
    display: inline-block;
}

.status-running {
    background-color: #3b82f6;
    color: white;
}

.status-pending {
    background-color: #f97316;
    color: white;
}

.status-completed {
    background-color: #22c55e;
    color: white;
}

.status-failed {
    background-color: #ef4444;
    color: white;
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
