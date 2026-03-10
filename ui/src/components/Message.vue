<script setup>
import { ref, computed } from 'vue'
import { useNotify } from '@/composables/notify.js'
const { role, tool_calls, content, details, loading } = defineProps(['role', 'tool_calls', 'content', 'details', 'loading'])

async function copyToClip(text) {
  try {
    await navigator.clipboard.writeText(text)
    useNotify("Copied to clipboard!")
  } catch (error) {
    console.error(error)
  }
}

const expandedToolCalls = ref({})

function toggleToolCall(index) {
  expandedToolCalls.value[index] = !expandedToolCalls.value[index]
}

const formattedToolCalls = computed(() => {
  if (!tool_calls) {
    return []
  }
  return tool_calls.map((tool, index) => {
    let args = tool.function.arguments
    try {
      const parsedArgs = JSON.parse(args)
      args = JSON.stringify(parsedArgs, null, 2)
    } catch (e) { /* Do nothing, leave as is */ }

    let result = tool.result
    if (result) {
        try {
            const parsedResult = JSON.parse(result)
            result = JSON.stringify(parsedResult, null, 2)
        } catch (e) { /* Do nothing, leave as is */ }
    }

    return {
      ...tool,
      function: {
        ...tool.function,
        arguments: args
      },
      result: result,
      isExpanded: !!expandedToolCalls.value[index]
    }
  })
})
</script>

<template>
    <div v-if="role !== 'info'" :class="['message', role, { 'loading': loading && role === 'assistant' }]">

        <div v-if="tool_calls && tool_calls.length > 0" class="tool-calls-list">
            <div v-for="(tool_call, index) in formattedToolCalls" :key="index" class="tool-call-item">
                <div @click="toggleToolCall(index)" class="tool-call-header">
                    <div class="tool-call-info">
                        <div class="tool-icon">
                            <i class="fa-solid fa-screwdriver-wrench"></i>
                        </div>
                        <span class="tool-text">Used <strong>{{ tool_call.function.name }}</strong></span>
                    </div>
                    <div class="tool-action">
                        <i :class="['fa-solid', tool_call.isExpanded ? 'fa-chevron-up' : 'fa-chevron-down']" class="chevron-icon"></i>
                    </div>
                </div>
                <div v-if="tool_call.isExpanded" class="tool-call-details">
                    <div class="details-section">
                        <div class="details-label">Input</div>
                        <pre><code>{{ tool_call.function.arguments }}</code></pre>
                    </div>
                    <div v-if="tool_call.result" class="details-section">
                        <div class="details-label">Output</div>
                        <pre class="tool-result-content"><code>{{ tool_call.result }}</code></pre>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="content && content.length" class="content-wrapper">
            <span v-html="content"></span>
        </div>

        <div v-if="loading && (!content || !content.length)" class="spinner">
            <div class="dot"></div>
            <div class="dot"></div>
            <div class="dot"></div>
        </div>
    </div>

    <div v-else class="info-message">
        <div class="info-icon">
            <i class="fa-solid fa-circle-info"></i>
            <div class="tooltip">
                Time: {{ (details.total_duration / 1_000_000_000).toFixed(2) }}(s)<br>
                Tokens: {{ details.eval_count }}<br>
                Tokens/Sec: {{ (details.eval_count / (details.total_duration / 1_000_000_000)).toFixed(2) }}
            </div>
        </div>
        <i class="fa-solid fa-clipboard" @click="copyToClip(content)"></i>
    </div>
</template>

<style scoped>
@keyframes message-fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes pulse {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1.0);
  }
}

.spinner {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
}

.spinner .dot {
  width: 8px;
  height: 8px;
  margin: 0 4px;
  background-color: var(--primary-color);
  border-radius: 50%;
  animation: pulse 1.4s infinite ease-in-out both;
}

.spinner .dot:nth-child(1) {
  animation-delay: -0.32s;
}

.spinner .dot:nth-child(2) {
  animation-delay: -0.16s;
}


@keyframes shimmer-wave {
    0% {
        transform: translateX(-100%);
    }
    100% {
        transform: translateX(100%);
    }
}

.message {
    padding: 12px 16px;
    margin-bottom: 12px;
    border-radius: 12px;
    word-wrap: break-word;
    max-width: 85%;
    line-height: 1.6;
    width: fit-content;
}

.message.user {
  justify-content: flex-end;
  background-color: var(--primary-color);
  margin-left: auto;
  align-items: flex-end;
  text-align: left;
  color: white;
}

.message.assistant {
  font-size: 1.05rem;
  justify-content: flex-start;
  border-bottom-left-radius: 4px;
  margin-right: auto;
  align-self: flex-start;
  text-align: left;
  animation: message-fade-in 0.5s ease-out;
}

.message.assistant.loading {
  position: relative;
  overflow: hidden;
}

.message.assistant.loading::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.15), transparent);
    animation: shimmer-wave 2.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

.info-message {
    display: flex;
    margin-right: auto;
    text-align: left;
    padding: 12px 20px;
    margin-bottom: 12px;
    padding-top: 0;
    position: relative;
}

.info-message i {
  margin-right: 20px;
  cursor: pointer;
  transition: color 0.3s ease;
}

.info-message i:hover {
  color: var(--primary-color);
}

.info-icon {
  position: relative;
}

.tooltip {
  position: absolute;
  bottom: 100%;
  background-color: var(--bg-color-light);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.3s ease, visibility 0.3s ease;
  z-index: 100;
}

.info-icon:hover .tooltip {
  opacity: 1;
  visibility: visible;
}

.tool-calls-list {
    margin-bottom: 12px;
    width: 100%;
}

.tool-call-item {
    background-color: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--border-color);
    border-radius: 10px;
    margin-bottom: 8px;
    overflow: hidden;
    width: 100%;
}

.tool-call-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 14px;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.tool-call-header:hover {
    background-color: rgba(255, 255, 255, 0.05);
}

.tool-call-info {
    display: flex;
    align-items: center;
    gap: 12px;
}

.tool-icon {
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(109, 40, 217, 0.15);
    color: var(--primary-color-light);
    border-radius: 8px;
    font-size: 0.9rem;
}

.tool-text {
    font-size: 0.95rem;
    color: var(--text-color);
}

.tool-action {
    display: flex;
    align-items: center;
    gap: 10px;
    opacity: 0.7;
}

.chevron-icon {
    font-size: 0.8rem;
    transition: transform 0.2s ease;
}

.tool-call-details {
    padding: 0 14px 14px;
    background-color: rgba(0, 0, 0, 0.1);
    border-top: 1px solid var(--border-color);
}

.details-section {
    margin-top: 12px;
}

.details-label {
    font-size: 0.7rem;
    font-weight: 700;
    text-transform: uppercase;
    color: var(--text-color);
    opacity: 0.5;
    margin-bottom: 6px;
    letter-spacing: 0.05em;
}

.tool-call-item pre {
    margin: 0;
    padding: 10px;
    font-size: 0.85rem;
    background-color: #0f172a !important;
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.tool-result-content {
    margin-top: 4px;
}
</style>

<style>
pre {
  background-color: #1e293b !important;
  padding: 8px;
  border-radius: 14px;
  overflow-x: auto;
  margin: 0.5rem;
}

code {
  font-family: 'Fira Code', monospace;
  border-radius: 10px;
}

.message ul {
    padding-left: 1rem;
}

.message ol {
    padding-left: 1rem;
}
</style>
