<script setup>
import { ref, onMounted, watch, nextTick, computed, reactive } from 'vue';
import { Marked } from 'marked';
import { markedHighlight } from "marked-highlight";
import hljs from 'highlight.js';
import 'highlight.js/styles/atom-one-dark.css';
import 'vue3-toastify/dist/index.css';
import { useAuthStore } from '@/stores/auth'
import { useStream } from '@/composables/stream.js'
import Message from '@/components/Message.vue'
import ChatInput from '@/components/ChatInput.vue'
import { useRoute, useRouter } from 'vue-router';
import { useNotify } from '@/composables/notify.js'

const route = useRoute();
const router = useRouter();

const authStore = useAuthStore();

const messageId = ref(route.params.id || null);

const chatContainer = ref(null);
const chatInput = ref(null);

const apiMessages = reactive([]);
const chatMessages = ref([]);

const userMessageHistory = computed(() => apiMessages.filter(msg => msg.role === 'user'));

const modelList = ref([])
const modelSelector = ref(null)

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

async function handleSendMessage({ content, webSearch, code }) {
  if (!authStore.username || !authStore.id) {
    useNotify("Invalid username or id")
    return
  }

  const chat = content;
  if (chat === '') return;

  apiMessages.push({ 'role': 'user', 'content': chat });
  chatMessages.value.push({ role: 'user', content: marked.parse(chat) });
  
  chatMessages.value.push({ role: 'assistant', content: '', loading: true });

  const body = {
    'model': modelSelector.value,
    'webSearch': webSearch,
    'code': code,
    'username': authStore.username,
    'userId': authStore.id.toString(),
    'messageId': messageId.value ? messageId.value.toString() : null,
    'messages': apiMessages
  }

  await useStream(body, messageId, apiMessages, chatMessages, router);
}

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

async function getMessages(id) {
  if (!authStore.username || !authStore.id) {
    useNotify("Invalid username or id")
    return
  }

  const body = {
    "chatId": id ? id : null,
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
      if (response.status === 401) {
        await router.replace('/login')
        return [];
      }
      useNotify(`Error: ${response.status} ${response.statusText}`);
      router.push({ name: 'New Chat' })
      return [];
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

        if (mess.length) {
            // apiMessages = mess
            apiMessages.splice(0, apiMessages.length, ...mess)
            chatMessages.value = parseLoadedMessages(mess)
        } else {

        }
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


  chatInput.value?.focus();
  await getModelList()
});
</script>

<template>
  <div class="chat-container">
    <div class="chat-header">
      <h2>Chat</h2>
      <select name="model" id="modelSelector" v-model="modelSelector">
        <option v-for="model in modelList" :value="model.model" :key="model.model">{{ model.name }}</option>
      </select>
    </div>

    <Transition name="fade" mode="out-in">
      <div v-if="chatMessages && chatMessages.length" class="chat-messages" ref="chatContainer">
        <template v-for="(message, index) in chatMessages" :key="index">
            <Message v-bind="message"/>
        </template>
      </div>

      <div class="welcome-container" v-else>
          <h1>Hello <span class="username">{{authStore.username}}</span>, How can I help you?</h1>
      </div>
    </Transition>


    <ChatInput
        ref="chatInput"
        :user-message-history="userMessageHistory"
        :is-centered="!chatMessages || !chatMessages.length"
        @send-message="handleSendMessage"
    />


  </div>
</template>

<style scoped>
.chat-container {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 1rem 2rem 120px 2rem;
  position: relative;
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

.chat-input-area {
  margin-top: 1.5rem;
}

.input-wrapper {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  border-radius: 16px;
}

.welcome-container {
    padding-top: 20vh;
    justify-content: center;
    display: flex;
}

.username {
    color: var(--primary-color-light);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
