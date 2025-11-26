<script setup>
import { useUiStore } from '@/stores/ui';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Sidebar from '@/components/sidebar.vue';
import { useNotify } from '@/composables/notify';

const modelList = ref([])
const modelSelector = ref(null)

const architecture = ref("")
const parameters = ref("")
const contextLength = ref("")
const embeddingLength = ref("")
const quantization = ref("")

const capabilities = ref([])
const license = ref("")

const uiStore = useUiStore();
const router = useRouter();

async function getModelList() {
  try {
    const response = await fetch('/api/modelList', { credentials: 'include' })

    if (!response.ok) {
      if (response.status === 401) {
        await router.replace('/login')
        return
      }
      useNotify(`Error: ${response.status} ${response.statusText}`);
      return;
    }

    const data = await response.json();

    if (!Array.isArray(data.models) || data.models.length === 0) {
      throw new Error('No model list')
    }

    modelList.value = data.models
    const names = data.models.map(i => i.name)
    modelSelector.value = data.preference || names[0]

    if (modelSelector.value) {
      await getModelInfo();
    }

  } catch (error) {
    console.error(error)
  }
}

async function getModelInfo() {
    if (!modelSelector.value) return;

    const data = {
        'model': modelSelector.value,
    }

    try {
        const response = await fetch('/api/modelInfo', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
            credentials: 'include'
        })

        const response_data = await response.json()
        
        let cap = response_data['Capabilities'] || []
        let modelInfo = response_data['ModelInfo'] || {}

        let arch = modelInfo['general.architecture'] || 'N/A'

        architecture.value = arch
        parameters.value = modelInfo['general.size_label'] || 'N/A'
        contextLength.value = modelInfo[`${arch}.context_length`] || 'N/A'
        embeddingLength.value = modelInfo[`${arch}.embedding_length`] || 'N/A'
        quantization.value = 'Coming Soon'

        capabilities.value = cap

        license.value = modelInfo['general.license'] || 'N/A'
    } catch (e) {
        console.error("Failed to get model info:", e)
        useNotify("Error fetching model details.");
    }
}

onMounted(async () => {
  await getModelList()
});
</script>

<template>
    <div class="job-page-container">
        <Sidebar />
        <main class="jobs-main-content" :style="{ marginLeft: uiStore.sidebarWidth }">
            <h1>Models</h1>

            <div class="model-selection-container">
              <label for="modelSelector">Select a Model</label>
              <select name="model" id="modelSelector" v-model="modelSelector" @change="getModelInfo">
                  <option v-for="model in modelList" :value="model.model" :key="model.model">{{ model.name }}</option>
              </select>
            </div>

            <div class="model-info-container">
              <div class="info-card">
                <h2>Model Details</h2>
                <div class="details-grid">
                  <div>
                    <strong>Architecture</strong>
                    <span>{{ architecture }}</span>
                  </div>
                  <div>
                    <strong>Parameters</strong>
                    <span>{{ parameters }}</span>
                  </div>
                  <div>
                    <strong>Context Length</strong>
                    <span>{{ contextLength }}</span>
                  </div>
                  <div>
                    <strong>Embedding Length</strong>
                    <span>{{ embeddingLength }}</span>
                  </div>
                  <div>
                    <strong>Quantization</strong>
                    <span>{{ quantization }}</span>
                  </div>
                </div>
              </div>

              <div class="info-card">
                <h2>Capabilities</h2>
                <div class="capabilities-list">
                  <span v-if="!capabilities.length">No capabilities listed.</span>
                  <span v-for="c in capabilities" :key="c" class="capability-badge">{{ c }}</span>
                </div>
              </div>
              
              <div class="info-card">
                <h2>License</h2>
                <p>{{ license }}</p>
              </div>
            </div>
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

.model-selection-container {
  margin-bottom: 2.5rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.model-selection-container label {
  font-weight: 500;
  font-size: 0.9rem;
  color: var(--text-color);
}

select#modelSelector {
    background-color: var(--bg-color-light);
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.75rem 1rem;
    border-radius: 8px;
    font-size: 1rem;
    min-width: 350px;
    cursor: pointer;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

select#modelSelector:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-color-light);
}

.model-info-container {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
}

.info-card {
  background-color: var(--bg-color-light);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.info-card:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.info-card h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
}

.details-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.details-grid div {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.95rem;
}

.details-grid strong {
    color: #9ca3af;
    font-weight: 500;
}

.details-grid span {
    color: var(--text-color);
    font-weight: 500;
    text-align: right;
}

.capabilities-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.capability-badge {
  background-color: var(--secondary-color);
  color: var(--bg-color);
  padding: 0.35rem 0.85rem;
  border-radius: 9999px;
  font-size: 0.85rem;
  font-weight: 600;
}

.info-card p {
    color: #9ca3af;
    line-height: 1.6;
}
</style>
