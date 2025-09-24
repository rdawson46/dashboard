<script setup>
    const { role, tool_calls, content, details } = defineProps(['role', 'tool_calls', 'content', 'details'])

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

    <div v-else-if="role !== 'info'" :class="['message', role]">

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

        <span v-if="content.length" v-html="content"></span>

        <i v-if="(!tool_calls || !tool_calls.length) && !content.length" class="fa-solid fa-spinner fa-spin-pulse"></i>
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
.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding-right: 1rem; /* for scrollbar */
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
</style>
