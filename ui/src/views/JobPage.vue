<script setup>
import Sidebar from '@/components/sidebar.vue'
import { ref } from 'vue';

const jobs = ref([
    {
        Name: "News Update",
        Type: "LLM",
        Status: "Running",
        Date: "2025-10-17",
    },
    {
        Name: "Health Check",
        Type: "Tool",
        Status: "Pending",
        Date: "2025-10-15",
    },
    {
        Name: "Daily Summary",
        Type: "LLM",
        Status: "Completed",
        Date: "2025-10-10",
    },
    {
        Name: "Internet Search",
        Type: "Tool",
        Status: "Failed",
        Date: "2025-10-05",
    }
])

</script>

<template>
    <div class="job-page-container">
        <Sidebar />
        <main class="jobs-main-content glass-card">
            <header class="jobs-header">
                <h1>Jobs</h1>
                <button class="create-job-btn">
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
                            <th>Date</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        <template v-for="job in jobs">
                            <tr>
                                <td>{{job.Name}}</td>
                                <td>{{job.Type}}</td>
                                <td><span class="status" :class="[`status-${job.Status.toLowerCase()}`]">{{job.Status}}</span></td>
                                <td>{{job.Date}}</td>
                                <td class="action-buttons">
                                    <button class="action-btn"><i class="fa-solid fa-pencil"></i></button>
                                    <button class="action-btn delete-btn"><i class="fa-solid fa-trash"></i></button>
                                </td>
                            </tr>
                        </template>
                    </tbody>
                </table>
            </div>
        </main>
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
