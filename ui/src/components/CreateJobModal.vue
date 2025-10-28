<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const emit = defineEmits(['close', 'create-job']);

const jobName = ref('');
const jobType = ref('LLM');

// LLM job
const query = ref('');

// Tool Job
const toolType = ref('code')
const codeFile = ref(null)
const webQuery = ref('')

const freq = ref('15mins');

function handleFileUpload(event) {
    const file = event.target.files[0];
    if (file) {
        if (file.size > 10 * 1024 * 1024) {
            alert("File size exceeds 10MB");
            event.target.value = ''; // Clear the file input
            codeFile.value = null;
            return;
        }
        if (!file.name.endsWith('.py')) {
            alert("Only .py files are allowed");
            event.target.value = ''; // Clear the file input
            codeFile.value = null;
            return;
        }
        codeFile.value = file;
    }
}

function createJob() {
    if (jobName.value && jobType.value) {
        const formData = new FormData();

        formData.append('name', jobName.value);
        formData.append('type', jobType.value);
        formData.append('freq', freq.value);

        if (jobType.value === 'LLM') {
            formData.append('query', query.value);

        } else if (jobType.value === 'Tool') {
            formData.append('toolType', toolType.value);

            if (toolType.value === 'code') {
                formData.append('codeFile', codeFile.value);
            } else if (toolType.value === 'search') {
                formData.append('webQuery', webQuery.value);
            }
        }

        emit('create-job', formData);
    }
}

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
            <h2>Create New Job</h2>
            <form @submit.prevent="createJob">
                <div class="form-group">
                    <label for="job-name">Job Name</label>
                    <input type="text" id="job-name" v-model="jobName" required>
                </div>

                <div class="form-group">
                    <label for="job-type">Job Type</label>
                    <select id="job-type" v-model="jobType" required>
                        <option value="LLM">LLM</option>
                        <option value="Tool">Tool</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="job-freq">Frequency</label>
                    <select id="job-freq" v-model="freq" required>
                        <option value="15mins">15 Minutes</option>
                        <option value="30mins">30 Minutes</option>
                        <option value="1hour">1 Hour</option>
                    </select>
                </div>

                <Transition name="fade" mode="out-in">
                    <div v-if="jobType === 'LLM'">
                        <div class="form-group">
                            <label for="job-name">Query</label>
                            <input type="text" id="query" v-model="query" required>
                        </div>
                    </div>
                    <div v-else-if="jobType === 'Tool'">
                        <div class="form-group">
                            <label for="tool-type">Tool</label>
                            <select id="tool-type" v-model="toolType" required>
                                <option value="code">Python Code</option>
                                <option value="search">Web Search</option>
                            </select>
                        </div>

                        <Transition name="fade" mode="out-in">
                            <div v-if="toolType === 'code'">
                                <div class="form-group">
                                    <label for="code-upload">Upload Code</label>
                                    <input type="file" id="code-upload" @change="handleFileUpload" accept=".py" required>
                                </div>
                            </div>

                            <div v-else-if="toolType === 'search'">
                                <div class="form-group">
                                    <label for="search">Search Query</label>
                                    <input type="text" id="search" v-model="webQuery" required>
                                </div>
                            </div>
                        </Transition>
                    </div>
                </Transition>
                <div class="form-actions">
                    <button type="button" class="cancel-btn" @click="$emit('close')">Cancel</button>
                    <button type="submit" class="create-btn">Create</button>
                </div>
            </form>
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

.form-group {
    margin-bottom: 1.5rem;
}

.form-group label {
    display: block;
    margin-bottom: 0.5rem;
    opacity: 0.8;
}

select, input[type="file"] {
  background-color: var(--bg-color-light);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  padding: 0.8rem 1rem;
  border-radius: 8px;
  font-size: 1rem;
  width: 100%;
}

select:focus, input[type="file"]:focus {
  outline: none;
  border-color: var(--primary-color);
}

input[type="file"]::file-selector-button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.3s ease;
    margin-right: 1rem;
}

input[type="file"]::file-selector-button:hover {
    background-color: var(--primary-color-light);
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
}

.cancel-btn {
    background-color: var(--bg-color-light);
}

.cancel-btn:hover {
    background-color: var(--border-color);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
