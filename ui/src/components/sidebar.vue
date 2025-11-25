<script setup>
import { ref, onMounted } from 'vue';
import 'vue3-toastify/dist/index.css';
import { useAuthStore } from '@/stores/auth';
import { RouterLink } from 'vue-router';
import { useRouter } from 'vue-router';
import { useNotify } from '@/composables/notify.js'
import { useUiStore } from '@/stores/ui';

const authStore = useAuthStore();
const router = useRouter();
const uiStore = useUiStore();

const history = ref([]);
const deletingChatId = ref(null);

async function getHistory() {
  try{
    const res = await fetch('/api/chatDescription', { credentials: 'include' });

    if (!res.ok) { 
        return
    }
    const data = await res.json();
    return data;
  } catch (e) {
    useNotify(`Can't get chat history`)
    throw e
  }
}

async function deleteChat(chatId) {
    if (!authStore.id) {
      useNotify('No user id')
    }

    deletingChatId.value = chatId;

    const body = {
      'userId': authStore.id.toString(),
      'ChatId': chatId.toString()
    }

    try {
      const response = await fetch('/api/deleteMessages', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
        credentials: 'include'
      })

      if (!response.ok) {
        useNotify(`Error: ${response.status} ${response.statusText}`);
        return;
      }

      const data = await response.json();

      if (!data['Status'] && data['Status'] != 'ok') {
        useNotify('Failed to delete chat')
        return
      }
      for (let i = 0; i < history.value.length; i++) {
        if (history.value[i].id === chatId) {
          history.value.splice(i, 1)
        }
      }
    } catch (e) {
      console.error('Error parsing JSON:', e);
      useNotify('Error processing server response.');
    } finally {
      deletingChatId.value = null;
    }
}

async function logout() {
  try {
    const response = await fetch('/api/logout', { credentials: 'include' })

    await router.push('/')
  } catch (e) {
    console.error(e)
  }
}

onMounted(async () => {
    try {
      const hist = await getHistory();
      history.value = hist;
    } catch (e) {
      console.error(e)
    }
});
</script>

<template>
    <nav class="sidebar glass-card" :class="{ collapsed: uiStore.isSidebarCollapsed }">
        <div class="logo">
            <i id="robot" class="fa-solid fa-robot"></i>
            <h1 v-show="!uiStore.isSidebarCollapsed">Chat</h1>
        </div>

        <ul class="nav-links">
            <li>
                <RouterLink :to="{ name: 'New Chat' }" :class="{ 'active': false }">
                    <i class="fa-solid fa-comments"></i>
                    <span class="link-text" v-show="!uiStore.isSidebarCollapsed">New Chat</span>
                </RouterLink>
            </li>
            <li>
                <a href="#">
                    <i class="fa-solid fa-user"></i>
                    <span class="link-text" v-show="!uiStore.isSidebarCollapsed">Profile</span>
                </a>
            </li>
            <li>
                <RouterLink :to="{ name: 'Models' }" :class="{ 'active': false }">
                    <i class="fa-solid fa-cogs"></i>
                    <span class="link-text" v-show="!uiStore.isSidebarCollapsed">Models</span>
                </RouterLink>
            </li>
            <li>
                <RouterLink :to="{ name: 'jobs' }" :class="{ 'active': false }">
                    <i class="fa-solid fa-clock"></i>
                    <span class="link-text" v-show="!uiStore.isSidebarCollapsed">Jobs</span>
                </RouterLink>
            </li>
        </ul>

        <div class="history-section" v-show="!uiStore.isSidebarCollapsed">
            <h3 class="history-title">History</h3>
            <div class="history">
                <ul>
                    <li class="item" v-for="h in history" :key="h.id">
                        <RouterLink :to="{ name: 'Existing Chat', params: { id: h.id } }" class="history-link">
                        {{ h.description }}
                        </RouterLink>
                        <i v-if="deletingChatId === h.id" class="fa-solid fa-spinner fa-spin space"></i>
                        <i v-else class="fa-solid fa-trash space" @click="deleteChat(h.id)"></i>
                    </li>
                </ul>
            </div>
        </div>

        <div class="sidebar-footer">
            <a @click="logout()" class="logout-link">
                <i class="fa-solid fa-right-from-bracket"></i>
                <span class="link-text" v-show="!uiStore.isSidebarCollapsed">Logout</span>
            </a>
            <button class="glass-card collapse-button" @click="uiStore.toggleSidebar">
                <i class="fa-solid" :class="uiStore.isSidebarCollapsed ? 'fa-chevron-right' : 'fa-chevron-left'"></i>
            </button>
        </div>
    </nav>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: 280px;
  min-width: 280px;
  padding: 2rem;
  transition: all 0.3s ease;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 100;
  border-radius: 0;
}

.sidebar.collapsed {
  width: 80px;
  min-width: 80px;
  padding: 2rem 1rem;
}

.sidebar.collapsed .logo {
    justify-content: center;
}

.logo {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
  transition: all 0.3s ease;
}

.logo #robot {
  font-size: 2rem;
  color: var(--primary-color-light);
  margin-right: 1rem;
  transition: margin 0.3s ease;
}

.sidebar.collapsed .logo #robot {
    margin-right: 0;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 700;
  transition: all 0.2s ease;
}

.nav-links {
  list-style: none;
  margin-bottom: 3rem;
}

.nav-links li a {
  display: flex;
  align-items: center;
  padding: 1rem;
  margin-bottom: 0.5rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--text-color);
  transition: background-color 0.3s ease, color 0.3s ease;
}

.sidebar.collapsed .nav-links li a {
    justify-content: center;
    padding: 1.2rem 1rem;
    margin-bottom: 0.2rem;
}

.nav-links li a:hover, .nav-links li a.active {
  background-color: var(--primary-color);
  color: white;
}

.nav-links li a i {
  margin-right: 1rem;
  font-size: 1.2rem;
  transition: margin 0.3s ease;
}

.sidebar.collapsed .nav-links li a i {
    margin-right: 0;
}

.link-text {
    transition: opacity 0.2s ease;
}

.history-section {
    flex-grow: 1;
    overflow-y: hidden;
    display: flex;
    flex-direction: column;
}

.history {
  flex-grow: 1;
  overflow-y: auto;
}

.history-title {
  font-size: 1rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.history ul {
  list-style: none;
}

.history ul li {
  padding: 0.5rem 0.25rem;
  cursor: pointer;
  opacity: 0.65;
  transition: opacity 0.3s ease;
  margin: 0.15rem;
  border-radius: 0.5rem;
}

.history-link {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 150px;
    display: inline-block;
}

.history ul li:hover {
  opacity: 1;
  background-color: rgba(255, 255, 255, 0.05);
}

.history::-webkit-scrollbar {
  width: 8px;
}

.history::-webkit-scrollbar-track {
  background: transparent;
}

.history::-webkit-scrollbar-thumb {
  background: #4f4f4f;
  border-radius: 4px;
}

.history::-webkit-scrollbar-thumb:hover {
  background: #6f6f6f;
}

.sidebar-footer {
    margin-top: auto;
    padding-top: 1rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.logout-link {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--text-color);
  transition: background-color 0.3s ease, color 0.3s ease;
}

.sidebar.collapsed .logout-link {
    justify-content: center;
}

.logout-link:hover {
    background-color: #ef4444;
    color: white;
}

.logout-link i {
    margin-right: 1rem;
    transition: margin 0.3s ease;
}

.sidebar.collapsed .logout-link i {
    margin-right: 0;
}

.collapse-button {
    margin-top: 1rem;
    padding: 0.5rem;
    font-size: 1.5rem;
    background-color: var(--bg-color-light);
    width: 100%;
    border: none;
    color: var(--text-color);
    cursor: pointer;
}

.item {
    display: flex;
    align-items: center;
}

.item i {
    visibility: hidden;
    margin-left: auto;
}

.item:hover i {
    visibility: visible;
}
</style>
