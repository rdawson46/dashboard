<script setup>
import { ref, onMounted, watch, nextTick, computed, reactive } from 'vue';
import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

const searchActive = ref(false);
const codeActive = ref(false);

const chatContainer = ref(null);
const inputContainer = ref(null);

const apiMessages = reactive([]);
const chatMessages = ref([]);

const inputValue = ref('');
const userMessageHistory = computed(() => apiMessages.filter(msg => msg.role === 'user'));
const historyIndex = ref(-1);
let pristineInput = '';

const modelList = ref([])
const modelSelector = ref(null)

const handleInput = (e) => {
    inputValue.value = e.target.value;
    if (historyIndex.value !== -1) {
        historyIndex.value = -1;
        pristineInput = '';
    }
};

function navigateHistoryUp() {
    if (!userMessageHistory.value.length) return;

    if (historyIndex.value === -1) {
        pristineInput = inputValue.value;
        historyIndex.value = userMessageHistory.value.length - 1;
    } else if (historyIndex.value > 0) {
        historyIndex.value--;
    } else {
        return;
    }
    inputValue.value = userMessageHistory.value[historyIndex.value].content;
}

function navigateHistoryDown() {
    if (historyIndex.value === -1) return;

    if (historyIndex.value < userMessageHistory.value.length - 1) {
        historyIndex.value++;
        inputValue.value = userMessageHistory.value[historyIndex.value].content;
    } else {
        historyIndex.value = -1;
        inputValue.value = pristineInput;
    }
}

const marked = new Marked(
  markedHighlight({
    langPrefix: 'hljs language-',
    highlight(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    }
  })
);

watch(chatMessages, async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
}, { deep: true });


function notify(message) {
  toast(message, {
    autoClose: 1000,
    theme: 'dark',
  });
}

async function copyToClip(text) {
  try {
    await navigator.clipboard.writeText(text)
    notify("Copied to clipboard!")
  } catch (error) {
    console.error(error)
  }
}

async function stream(url, body) {
  body['model'] = modelSelector.value.value
  body['webSearch'] = searchActive.value
  body['code'] = codeActive.value

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
      credentials: 'include'
    });

    if (!response.ok) {
      notify(`Error: ${response.status} ${response.statusText}`);
      chatMessages.value.pop();
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    let fullResponse = '';
    
    const assistantMessage = chatMessages.value[chatMessages.value.length - 1];
    assistantMessage.content = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';

      for (const line of lines) {
        if (line.startsWith("data: ")) {
          const jsonStr = line.substring(6);
          if (jsonStr.trim()) {
            try {
              const data = JSON.parse(jsonStr);
              if (data.done) {
                if (data.message.content.length) {
                  fullResponse += data.message.content
                  assistantMessage.content = marked.parse(fullResponse);
                }
                apiMessages.push({ 'role': 'assistant', 'content': fullResponse });
                
                chatMessages.value.push({ role: 'info', content: fullResponse, details: data });
                return;
              }

              if (data.message.tool_calls) {
                for (const toolCall of data.message.tool_calls) {
                  let { name } = toolCall;
                  console.log(toolCall)
                }
                continue
              }

              let token = data.message.content;
              fullResponse += token;
              assistantMessage.content = marked.parse(fullResponse);
            } catch (e) {
              console.error('Error parsing JSON:', e);
              notify('Error processing server response.');
            }
          }
        }
      }
    }
  } catch (error) {
    console.error('Error during fetch:', error);
    notify('Error when querying agent.');
    chatMessages.value.pop();
  } 
}

async function query() {
  const chat = inputValue.value.trim();
  if (chat === '') return;

  apiMessages.push({ 'role': 'user', 'content': chat });
  chatMessages.value.push({ role: 'user', content: marked.parse(chat) });
  
  inputValue.value = '';
  historyIndex.value = -1;

  chatMessages.value.push({ role: 'assistant', content: '<i class="fa-solid fa-spinner fa-spin-pulse"></i>' });
  await stream('/api/stream', { "messages": apiMessages });
}

async function getModelList() {
  try {
    const response = await fetch('/api/modelList', { credentials: 'include' })

    if (!response.ok) {
      notify(`Error: ${response.status} ${response.statusText}`);
      return;
    }

    const data = await response.json();
    modelList.value = data
  } catch (error) {
    console.error(error)
  }
}

onMounted(async () => {
  inputContainer.value.focus();
  await getModelList()
});
</script>

<template>
  <div class="chat-container">
    <div class="chat-header">
      <h2>Chat</h2>
      <select name="model" id="modelSelector" ref="modelSelector">
        <option v-for="model in modelList" :value="model.model" :key="model.model">{{ model.name }}</option>
      </select>
    </div>

    <div class="chat-messages" ref="chatContainer">
      <template v-for="(message, index) in chatMessages" :key="index">

        <div v-if="message.role !== 'info'" :class="['message', message.role]" v-html="message.content"></div>

        <div v-else class="info-message">
          <div class="info-icon">
            <i class="fa-solid fa-circle-info"></i>
            <div class="tooltip">
              Time: {{ (message.details.total_duration / 1_000_000_000).toFixed(2) }}(s)<br>
              Tokens: {{ message.details.eval_count }}<br>
              Tokens/Sec: {{ (message.details.eval_count / (message.details.total_duration / 1_000_000_000)).toFixed(2) }}
            </div>
          </div>
          <i class="fa-solid fa-clipboard" @click="copyToClip(message.content)"></i>
        </div>

      </template>
    </div>

    <!--
    <div class="chat-input-area">
      <div class="input-wrapper glass-card">
        <textarea
          ref="inputContainer"
          id="message-input"
          placeholder="Type your message..."
          v-model="inputValue"
          @input="handleInput"
          @keydown.down.exact.prevent="navigateHistoryDown()"
          @keydown.up.exact.prevent="navigateHistoryUp()"
          @keydown.enter.exact.prevent="query"
        ></textarea>
        <button id="send-btn" @click="query"><i class="fa-solid fa-paper-plane"></i></button>
      </div>
      <div class="input-options">
        <button @click='searchActive = !searchActive' :class="{ active: searchActive }">
          <i class="fa-solid fa-globe"></i>
          <span>Web Search</span>
        </button>
        <button @click='codeActive = !codeActive' :class="{ active: codeActive }">
          <i class="fa-solid fa-code"></i>
          <span>Code</span>
        </button>
      </div>
    </div>
    -->
    <div class="input-area-main" :class="{ 'middle': !chatMessages.length }">


      <div class="input-area-sub">
        <div
          ref="inputContainer"
          id="message-input"
          placeholder="Type your message..."
          contenteditable="true"
          @input="handleInput"
          @keydown.down.exact.prevent="navigateHistoryDown()"
          @keydown.up.exact.prevent="navigateHistoryUp()"
          @keydown.enter.exact.prevent="query"
        ></div>
      </div>


      <div class="input-area-buttons">
        <button @click='searchActive = !searchActive' :class="{ active: searchActive }">
          <i class="fa-solid fa-globe"></i>
          <span>Web Search</span>
        </button>
        <button @click='codeActive = !codeActive' :class="{ active: codeActive }">
          <i class="fa-solid fa-code"></i>
          <span>Code</span>
        </button>
        <button id="send-btn" @click="query"><i class="fa-solid fa-paper-plane"></i></button>
      </div>


    </div>


  </div>
</template>

<style scoped>
.chat-container {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 80%;
  padding: 2rem;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.chat-header h2 {
  font-size: 1.8rem;
  font-weight: 700;
}

.chat-header select {
  background-color: var(--bg-color-light);
  border: 1px solid var(--border-color);
  color: var(--text-color);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-size: 1rem;
}

.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding-right: 1rem; /* for scrollbar */
}

.message {
    padding: 3px 20px;
    margin-bottom: 12px;
    border-radius: 18px;
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

.chat-input-area {
  margin-top: 1.5rem;
}

.input-wrapper {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  border-radius: 16px;
}

/*
#message-input {
  flex-grow: 1;
  border: none;
  color: var(--text-color);
  font-size: 1rem;
  padding: 1rem;
  resize: none;
}

#message-input:focus {
  outline: none;
}

#send-btn {
  background: var(--primary-color);
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 1.2rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

#send-btn:hover {
  background: var(--primary-color-light);
}

.input-options {
  display: flex;
  gap: 1rem;
}

.input-options button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--bg-color-light);
  color: var(--text-color);
  border: 1px solid var(--border-color);
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.input-options button.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}
*/

.input-area-main {
    align-items: left;
    background-color: #2d2d2d;
    border-radius: 25px;
    padding: 5px 15px;
    width: 60%;
    transition: 500ms ease-out;
    margin: auto;
}

input-area-sub {
    display: flex;
    align-items: center;
}

.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.input-area-buttons {
    width: 100%;
    display: flex;
}

.input-area-buttons button:last-child {
    margin-left: auto;
}

.input-area-buttons > button:hover {
    background-color: var(--bg-color-light);
}

#message-input {
    flex-grow: 1;
    padding: 5px 5px 0 5px;
    border: none;
    resize: none;
    font-family: 'Sqgoe UI', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: 1rem;
    line-height: 1.5;
    max-height: 200px;
    min-height: 40px;
    color: #f0f0f0;
    background-color: transparent;
    text-align: left;
}

#message-input:focus {
    outline: none;
}

#send-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    padding: 10px;
    transition: color 0.2s ease;
}

#send-btn:hover {
    color: #81b2f3;
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
