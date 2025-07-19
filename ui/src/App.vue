<script setup>
import { ref, onMounted, watch, nextTick } from 'vue';
import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/rose-pine-moon.css';
import Sidebar from './components/sidebar.vue';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

const searchActive = ref(false);
const codeActive = ref(false);

const chatContainer = ref(null);
const inputContainer = ref(null);

const apiMessages = [];
const chatMessages = ref([]); // For rendering: { role, content }

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

async function stream(url, body) {
  body['webSearch'] = searchActive.value
  body['code'] = codeActive.value
  try {
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
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
    assistantMessage.html = '';

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
                apiMessages.push({ 'role': 'assistant', 'content': fullResponse });
                
                // TODO: add a pop up for the data on hover
                chatMessages.value.push({ role: 'info', content: '<i class="fa-solid fa-circle-info"></i>', details: data });
                return;
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
  const chat = inputContainer.value.innerText.trim();
  if (chat === '') return;

  apiMessages.push({ 'role': 'user', 'content': chat });
  chatMessages.value.push({ role: 'user', content: marked.parse(chat) });
  
  inputContainer.value.innerText = '';

  chatMessages.value.push({ role: 'assistant', content: '<i class="fa-solid fa-spinner fa-spin-pulse"></i>' });
  await stream('/api/stream', { "messages": apiMessages });
}

onMounted(() => {
  inputContainer.value.focus();
});
</script>

<template>
  <div id="app-container">
    <Sidebar />
    <main id="main-content">
      <div id="chat-container" ref="chatContainer">
        <template v-for="(message, index) in chatMessages" :key="index">
          <div v-if="message.role !== 'info'" :class="[message.role, 'message']" v-html="message.content"></div>

          <div v-else class="info">
              <div class="hover-info" v-html="message.content"></div>
              <div class='hidden-info'>
                  Time: {{message.details.total_duration / 1_000_000_000}}(s)
                  <br>
                  Tokens: {{message.details.eval_count}}
              </div>
          </div>
        </template>
      </div>

      <div class="input-area-main">
        <div class="input-area-sub">
            <div ref="inputContainer" id="message-input" contenteditable="true" @keydown.enter.exact.prevent="query"></div>
            <button id="send-btn" @click='query'><i class="fa-solid fa-paper-plane"></i></button>
        </div>

        <div class="input-area-buttons">
            <button @click='searchActive = !searchActive' :class="{ active: searchActive }">
                <i class="fa-solid fa-globe"></i>
                Web Search
            </button>
            <button @click='codeActive = !codeActive' :class="{ active: codeActive }">
                <i class="fa-solid fa-code"></i>
                Code
            </button>
        </div>

      </div>
    </main>
  </div>
</template>

<style scoped>
#app-container {
  display: flex;
  height: 100vh;
  background-color: #1e1e1e;
  width: 100vw;
}

#main-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100% - 40px;
  padding: 20px;
  margin-left: 6rem; /* Adjust for sidebar width */
  transition: margin-left 200ms ease;
  align-items: center;
}

#chat-container {
  flex-grow: 1;
  overflow-y: auto;
  padding: 10px;
  border-radius: 8px;
  margin-bottom: 20px;
  width: 90%;
}

.input-area-main {
  align-items: left;
  background-color: #2d2d2d;
  border-radius: 25px;
  padding: 5px 15px;
  width: 50%;
}

.input-area-sub {
  display: flex;
  align-items: center;
}

.input-area-buttons {
    width: fit-content;
}

.active {
    background-color: #4a90e2;
}

.input-area-buttons > button {
    border: 1px solid #24b4fb;
    border-radius: 0.9em;
    cursor: pointer;
    padding: 0.8em 1.2em 0.8em 1em;
    font-size: 600;
    margin: 0 5px;
}

.input-area-buttons > button:hover {
    background-color: #0071e2;
}

#message-input {
  flex-grow: 1;
  padding: 10px;
  border: none;
  resize: none;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  font-size: 1rem;
  line-height: 1.5;
  min-height: 40px;
  max-height: 200px;
  overflow-y: auto;
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
  color: #4a90e2;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 10px;
  transition: color 0.2s ease;
}

#send-btn:hover {
  color: #81b2f3;
}

@media only screen and (min-width: 600px) {
  .navbar:hover ~ #main-content {
    margin-left: 16rem;
  }
}
</style>

<style>
.message {
  padding: 3px 20px;
  margin-bottom: 12px;
  border-radius: 18px;
  word-wrap: break-word;
  max-width: 85%;
  line-height: 1.6;
  animation: fadeIn 0.5s ease-in-out;
  width: fit-content;
}

.user {
  background-color: #cf6f6a;
  color: white;
  margin-left: auto;
  align-self: flex-end;
  text-align: right;
}

.assistant {
  background-color: #3a3a3a;
  color: #f0f0f0;
  margin-right: auto;
  align-self: flex-start;
  text-align: left;
  margin-bottom: 0.5em;
}

.info {
  color: #f0f0f0;
  margin-right: auto;
  align-self: flex-start;
  text-align: left;
  padding: 12px 20px;
  margin-bottom: 12px;
  width: fit-content;
  padding-top: 0;
  position: relative;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

pre {
  background-color: #1a1a1a;
  padding: 10px;
  border-radius: 15px;
  overflow-x: auto;
  font-family: 'Fira Code', monospace;
}

code {
  border-radius: 10px;
  font-family: 'Fira Code', monospace;
}

.hidden-info {
    display: none;
}

.hover-info:hover + .hidden-info {
    display: block;
    position: absolute;
    z-index: 100;
    min-width: 25%;
    width: fit-content;
    background-color: #232323;
    padding: 5px;
    border-radius: 5px;
    top: -60px;
    white-space: nowrap;
}
</style>
