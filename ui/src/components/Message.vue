<script setup>
import { ref, computed } from 'vue'
const { role, tool_calls, content, details, loading } = defineProps(['role', 'tool_calls', 'content', 'details', 'loading'])

const isToolCallsExpanded = ref(false)

async function copyToClip(text) {
  try {
    await navigator.clipboard.writeText(text)
    notify("Copied to clipboard!")
  } catch (error) {
    console.error(error)
  }
}

const formattedToolCalls = computed(() => {
  if (!tool_calls) {
    return []
  }
  return tool_calls.map(tool => {
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
    }
  })
})
</script>

<template>
    <div v-if="role !== 'info'" :class="['message', role, { 'loading': loading && role === 'assistant' }]">

        <div v-if="tool_calls && tool_calls.length > 0" class="tool-calls-container">
            <div @click="isToolCallsExpanded = !isToolCallsExpanded" class="tool-calls-header">
                <div class="tool-calls-title">
                    <i class="fa-solid fa-terminal"></i>
                    <span>Tool Calls ({{ tool_calls.length }})</span>
                </div>
                <i :class="['fa-solid', isToolCallsExpanded ? 'fa-chevron-down' : 'fa-chevron-right']" class="chevron-icon"></i>
            </div>
            <div v-if="isToolCallsExpanded" class="tool-calls-body">
                <div v-for="(tool_call, index) in formattedToolCalls" :key="index" class="tool-call-item">
                    <div class="tool-call-name">
                        <i class="fa-solid fa-bolt"></i>
                        <span>{{ tool_call.function.name }}</span>
                    </div>
                    <pre><code>{{ tool_call.function.arguments }}</code></pre>
                    <div v-if="tool_call.result" class="tool-result-in-dropdown">
                        <div class="tool-result-header">
                            <i class="fa-solid fa-gears"></i>
                            <strong>Tool Result</strong>
                        </div>
                        <pre class="tool-result-content">{{ tool_call.result }}</pre>
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

.tool-result-header {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border-color);
    font-size: 0.9rem;
    color: var(--text-color-secondary);
}

.tool-result-in-dropdown .tool-result-header {
    border-top: 1px solid var(--border-color);
    border-bottom: none;
}

.tool-result-header i {
    margin-right: 0.5rem;
}

.tool-result-content {
    padding: 1rem;
    white-space: pre-wrap;
    word-wrap: break-word;
    background-color: #1e293b !important;
    color: #f8fafc;
    margin: 0;
    border-bottom-left-radius: 6px;
    border-bottom-right-radius: 6px;
}

.message {
    padding: 5px 15px;
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

.tool-calls-container {
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 8px;
    margin-bottom: 10px;
    border: 1px solid var(--border-color);
}

.tool-calls-header {
    padding: 10px 15px;
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    align-items: center;
    transition: background-color 0.2s ease;
}

.tool-calls-header:hover {
    background-color: rgba(0, 0, 0, 0.1);
}

.tool-calls-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-weight: bold;
}

.chevron-icon {
    transition: transform 0.2s ease;
}

.tool-calls-body {
    padding: 0 15px 15px;
}

.tool-call-item {
    background-color: rgba(0, 0, 0, 0.15);
    border-radius: 6px;
    margin-top: 10px;
}

.tool-call-name {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 500;
    padding: 8px 12px;
    background-color: rgba(255, 255, 255, 0.05);
    border-top-left-radius: 6px;
    border-top-right-radius: 6px;
}

.tool-call-item pre {
    margin: 0;
    border-top-left-radius: 0;
    border-top-right-radius: 0;
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
