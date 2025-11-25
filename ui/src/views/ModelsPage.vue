<script setup>
import { useUiStore } from '@/stores/ui';
import { ref, onMounted } from 'vue';
import Sidebar from '@/components/sidebar.vue';

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

  } catch (error) {
    console.error(error)
  }
}

async function getModelInfo() {
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
        console.log(response_data)

        let cap = response_data['Capabilities']
        let modelInfo = response_data['ModelInfo']

        let arch = modelInfo['general.architecture']

        architecture.value = arch
        // parameters.value = modelInfo['general.parameter_count']
        parameters.value = modelInfo['general.size_label']
        contextLength.value = modelInfo[`${arch}.context_length`]
        embeddingLength.value = modelInfo[`${arch}.embedding_length`] 
        // quantization.value // have to format this?

        capabilities.value = cap

        license.value = modelInfo['general.license']
    } catch (e) {

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

            <select name="model" id="modelSelector" v-model="modelSelector">
                <option v-for="model in modelList" :value="model.model" :key="model.model">{{ model.name }}</option>
            </select>

            <button @click="getModelInfo()">Test</button>

            <div>Architecture: {{ architecture }}</div>
            <div>Parameters: {{ parameters }}</div>
            <div>Context Length: {{ contextLength }}</div>
            <div>Embedding Length: {{ embeddingLength }}</div>
            <div>Quantization: {{ quantization }}</div>


            <div>Capabilities:</div>
            <div v-for="c in capabilities">
                {{ c }}
            </div>

            <div>License: {{ license }}</div>
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
    padding: 2rem;
    transition: margin-left 0.3s ease;
    border-radius: 0;
}

select {
    background-color: var(--bg-color-light);
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.5rem 1rem;
    border-radius: 8px;
    font-size: 1rem;
}
</style>
