<script setup>
import { onMounted, onUnmounted } from 'vue';

const props = defineProps({
    job: {
        type: Object,
        required: true
    }
});

const emit = defineEmits(['close']);

function handleKeydown(e) {
    if (e.key === 'Escape') {
        emit('close');
    }
}

onMounted(() => {
    window.addEventListener('keydown', handleKeydown);
});

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown);
});
</script>

<template>
    <div class="modal-backdrop">
        <div class="modal-content glass-card">
            <h2>{{ job.name }}</h2>
            <div class="job-details">
                <div class="detail-item">
                    <span class="label">Job Type:</span>
                    <span class="value">{{ job.task.task_type }}</span>
                </div>
                <div class="detail-item">
                    <span class="label">Frequency:</span>
                    <span class="value">{{ job.freq }}</span>
                </div>
                <div class="detail-item">
                    <span class="label">Status:</span>
                    <span class="value">{{ job.status }}</span>
                </div>
                <div v-if="job.task.task_type === 'LLM'" class="detail-item">
                    <span class="label">Query:</span>
                    <span class="value">{{ job.task.query }}</span>
                </div>
                 <div v-if="job.type === 'Tool'" class="detail-item">
                    <span class="label">Tool Type:</span>
                    <span class="value">{{ job.tool_type }}</span>
                </div>
                <div v-if="job.tool_type === 'code'" class="detail-item">
                    <span class="label">Code File:</span>
                    <span class="value">{{ job.code_file }}</span>
                </div>
                <div v-if="job.result">
                    <span class="label">Result:</span>
                    <span class="value">{{ job.result }}</span>
                </div>
            </div>
            <div class="form-actions">
                <button type="button" class="primary-btn" @click="$emit('close')">Close</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 100;
}

.modal-content {
    width: 90%;
    max-width: 500px;
    padding: 2rem;
}

h2 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 2rem;
}

.job-details {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
}

.detail-item {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
}

.detail-item:last-child {
    border-bottom: none;
}

.label {
    font-weight: 600;
    opacity: 0.8;
}

.value {
    font-weight: 400;
    max-width: 70%;
    overflow-wrap: break-word;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
}

.primary-btn {
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 0.8rem 1.5rem;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.3s ease;
}

.primary-btn:hover {
    background-color: var(--primary-color-light);
}

button {
    padding: 0.8rem 1.5rem;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.3s ease;
}
</style>
