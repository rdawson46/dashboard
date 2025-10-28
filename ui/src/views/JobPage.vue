<script setup>
import Sidebar from '@/components/sidebar.vue'
import CreateJobModal from '@/components/CreateJobModal.vue'
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useNotify } from '@/composables/notify.js'

const router = useRouter();

const jobs = ref([])
const limit = ref(10)
const offset = ref(0)
const loading = ref(true)
const showCreateJobModal = ref(false)

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

        // remove item from list by id
        const modifiedJobs = jobs.value.filter(j => j.id != id);
        jobs.value = modifiedJobs

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
        jobs.value.push(data)
    } catch (e) {
        useNotify('Unable to create job')
        console.error(e)
    }
    showCreateJobModal.value = false
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

        jobs.value = data.jobs
        return 
    } catch (e) {
        console.error(e)
    }
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
        <main class="jobs-main-content glass-card">
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
                                    <button class="action-btn"><i class="fa-solid fa-pencil"></i></button>
                                    <button class="action-btn delete-btn" @click="deleteJob(job.id)"><i class="fa-solid fa-trash"></i></button>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            
                <div v-if="loading">
                    <i  class="fa-solid fa-spinner fa-spin-pulse"></i>
                </div>
            </div>
        </main>
        <CreateJobModal 
            v-if="showCreateJobModal" 
            @close="showCreateJobModal = false"
            @create-job="handleCreateJob"
        />
    </div>
</template>

<style scoped>
.job-page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  padding: 1.5rem;
  gap: 1.5rem;
  position: relative;
}

.jobs-main-content {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    padding: 2rem;
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
</style>
