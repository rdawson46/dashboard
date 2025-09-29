<script setup>
import { ref, onMounted, watch, watchEffect, nextTick, computed, reactive } from 'vue';
import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import { useAuthStore } from '@/stores/auth'
import { useStream } from '@/composables/stream.js'
import Message from '@/components/Message.vue'
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();

const authStore = useAuthStore();

const messageId = ref(route.params.id || null);

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
    inputValue.value = e.target.innerText;
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

watch(inputValue, (newValue) => {
    if (inputContainer.value && newValue !== inputContainer.value.innerText) {
        inputContainer.value.innerText = newValue;
        
        nextTick(() => {
            if (document.activeElement === inputContainer.value) {
                const selection = window.getSelection();

                if (selection) {
                    const range = document.createRange();
                    range.selectNodeContents(inputContainer.value);
                    range.collapse(false)
                    selection.removeAllRanges()
                    selection.addRange(range)
                }
            }
        });
    }
});

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

async function query() {
  if (!authStore.username || !authStore.id) {
    notify("Invalid username or id")
    chatMessages.value.pop();
    return
  }

  const chat = inputValue.value.trim();
  if (chat === '') return;

  apiMessages.push({ 'role': 'user', 'content': chat });
  chatMessages.value.push({ role: 'user', content: marked.parse(chat) });
  
  inputValue.value = '';
  historyIndex.value = -1;

  chatMessages.value.push({ role: 'assistant', content: '' });

  const body = {
    'model': modelSelector.value.value,
    'webSearch': searchActive.value,
    'code': codeActive.value,
    'username': authStore.username,
    'userId': authStore.id.toString(),
    'messageId': messageId.value ? messageId.value.toString() : null,
    'messages': apiMessages
  }

  await useStream(body, messageId, apiMessages, chatMessages);
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

async function getMessages(id) {
  if (!authStore.username || !authStore.id) {
    notify("Invalid username or id")
    chatMessages.value.pop();
    return
  }

  const body = {
    "chatId": messageId.value ? messageId.value.toString() : null,
    "userId": authStore.id.toString(),
  }

  try {
    const response = await fetch('/api/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body),
      credentials: 'include',
    })

    if (!response.ok) {
      notify(`Error: ${response.status} ${response.statusText}`);
      return;
    }

    return await response.json()
  } catch (error) {
    console.error(error)
  }
}

function parseLoadedMessages(messages) {
  let result = []
  if (!messages || !Array.isArray(messages) || !messages.length ) {
    return
  }

  for (let message of messages) {
    if (message.role == 'assistant') {
      let temp = {...message}
      temp.content = marked.parse(message.content)
      result.push(temp)
    } else {
      result.push(message)
    }
  }

  return result
}

watch(() => route.params.id, async (newId, oldId) => {
  if (newId && oldId != newId) {
    if (newId != messageId.value) {
        console.log(`setting new ID: ${newId}`)
        messageId.value = newId
        let mess = await getMessages(messageId.value)
        // apiMessages = mess
        apiMessages.splice(0, apiMessages.length, ...mess)
        chatMessages.value = parseLoadedMessages(mess)
    }
  }

  if (!newId) {
    messageId.value = newId
    apiMessages.splice(0, apiMessages.length)
    chatMessages.value = []
  }
}, { immediate: true })

onMounted(async () => {
  console.log(`Loaded with session: ${messageId.value}`)

  if (messageId.value) {
    let mess = await getMessages(messageId.value)
    // apiMessages = mess
    apiMessages.splice(0, apiMessages.length, ...mess)
    chatMessages.value = parseLoadedMessages(mess)
  }


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

    <div v-if="chatMessages.length" class="chat-messages" ref="chatContainer">
      <template v-for="(message, index) in chatMessages" :key="index">
          <Message v-bind="message"/>
      </template>
    </div>

    <div class="welcome-container" v-else>
        <h1>Hello <span class="username">{{authStore.username}}</span>, How can I help you?</h1>
    </div>


    <div class="input-area-main glass-card" :class="{ 'middle': !chatMessages || !chatMessages.length }">


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
  padding: 1rem 2rem 0 2rem;
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
    /*background-color: #2d2d2d;*/
    border-radius: 25px;
    padding: 5px 15px;
    width: 75%;
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
  display: flex;
  gap: 1rem;
  margin: 0.25rem 0;
}

.input-area-buttons button {
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

.input-area-buttons button.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.input-area-buttons button:last-child {
    margin-left: auto;
}
/*
   buttons from og
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
*/

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

.welcome-container {
    margin-top: auto;
    margin-bottom: auto;
    justify-content: center;
    display: flex;
}

.username {
    color: var(--primary-color-light);
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
