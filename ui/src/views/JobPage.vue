<script setup>
import Sidebar from '@/components/sidebar.vue'
</script>

<template>
    <div class="job-page-container">
        <Sidebar />
        <h2>Hello</h2>
        <h3>No jobs found</h3>
        <button>create job</button>
    </div>
</template>


<style scoped>
.job-page-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  padding: 1.5rem;
  gap: 1.5rem;
}
</style>
