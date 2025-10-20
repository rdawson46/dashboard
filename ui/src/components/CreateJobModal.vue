<script setup>
import { ref } from 'vue';

const emit = defineEmits(['close', 'create-job']);

const jobName = ref('');
const jobType = ref('LLM');

// LLM job
const query = ref('');

// Tool Job
const toolType = ref('code')
const code = ref('')
const webQuery = ref('')

const freq = ref('15mins');

function createJob() {
    if (jobName.value && jobType.value) {
        emit('create-job', { name: jobName.value, type: jobType.value });
    }
}
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

                <Transition name="fade" mode="out-in">
                    <div v-if="jobType === 'LLM'">
                        <div class="form-group">
                            <label for="job-freq">Frequency</label>
                            <select id="job-freq" v-model="freq" required>
                                <option value="15mins">15 Minutes</option>
                                <option value="30mins">30 Minutes</option>
                                <option value="1hour">1 Hour</option>
                            </select>
                        </div>

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
                                    <label for="code">Code</label>
                                    <input type="text" id="code" v-model="code" required>
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

select {
  background-color: var(--bg-color-light);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  padding: 0.8rem 1rem;
  border-radius: 8px;
  font-size: 1rem;
  width: 100%;
}

select:focus {
  outline: none;
  border-color: var(--primary-color);
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
