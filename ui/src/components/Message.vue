<script setup>
    const { role, tool_calls, content, details, loading } = defineProps(['role', 'tool_calls', 'content', 'details', 'loading'])

    async function copyToClip(text) {
      try {
        await navigator.clipboard.writeText(text)
        notify("Copied to clipboard!")
      } catch (error) {
        console.error(error)
      }
    }
</script>

<template>
    <div v-if="role == 'tool'" class="glass-card tool-result">
        🔨 <b>Tool Call Response</b> 🔧
        <br>
        Response: {{content}}
    </div>

    <div v-else-if="role !== 'info'" :class="['message', role, { 'loading': loading && role === 'assistant' }]">

        <div v-if="tool_calls" class="glass-card tool-call">
            🔨 <b>Tool Calls</b> 🔧

            <template v-for="tool_call in tool_calls">
                <div>
                    {{tool_call.function.name}}
                    <br>
                    {{tool_call.function.arguments}}
                </div>
            </template>
        </div>

        <div v-if="content.length" class="content-wrapper">
            <span v-html="content"></span>
        </div>

        <div v-if="(!tool_calls || !tool_calls.length) && !content.length" class="spinner">
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

.tool-call, .tool-result {
    padding: 1rem;
}

.tool-result {
    width: 85%;
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
