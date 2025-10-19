<script setup>
import { ref, onMounted } from 'vue';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';
import { useAuthStore } from '@/stores/auth';
import { RouterLink } from 'vue-router';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const router = useRouter();

const history = ref([]);
const collapse = ref(false)
const deletingChatId = ref(null);

function notify(message) {
    toast(message, {
        autoClose: 1000,
        theme: 'dark',
    });
}

async function getHistory() {
  try{
    const res = await fetch('/api/chatDescription', { credentials: 'include' });

    if (!res.ok) { 
        return
    }
    const data = await res.json();
    return data;
  } catch (e) {
    notify(`Can't get chat history`)
    throw e
  }
}

async function deleteChat(chatId) {
    if (!authStore.id) {
      notify('No user id')
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
        notify(`Error: ${response.status} ${response.statusText}`);
        return;
      }

      const data = await response.json();

      if (!data['Status'] && data['Status'] != 'ok') {
        notify('Failed to delete chat')
        return
      }
      for (let i = 0; i < history.value.length; i++) {
        if (history.value[i].id === chatId) {
          history.value.splice(i, 1)
        }
      }
    } catch (e) {
      console.error('Error parsing JSON:', e);
      notify('Error processing server response.');
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
    <Transition name="slide-fade" mode="out-in" appear>
        <nav v-if="!collapse" class="sidebar glass-card">
            <div class="logo">
                <i id="robot" class="fa-solid fa-robot"></i>
                <h1>Chat</h1>
                <button class="glass-card collapse-button" @click="collapse = !collapse">
                    <i class="test fa-solid fa-bars"></i>
                </button>
            </div>
            <ul class="nav-links">
                <li>
                    <RouterLink :to="{ name: 'New Chat' }" :class="{ 'active': false }">
                    <i class="fa-solid fa-comments"></i><span>New Chat</span>
                    </RouterLink>
                </li>

                <li><a href="#"><i class="fa-solid fa-user"></i><span>Profile</span></a></li>
                <li><a href="/models"><i class="fa-solid fa-cogs"></i><span>Models</span></a></li>

                <li>
                    <RouterLink :to="{ name: 'jobs' }" :class="{ 'active': false }">
                        <i class="fa-solid fa-clock"></i>
                        <span>Jobs</span>
                    </RouterLink>
                </li>
            </ul>

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
            <div class="logout">
                <a @click="logout()">
                    <i class="fa-solid fa-right-from-bracket"></i>
                    <span>Logout</span>
                </a>
            </div>

        </nav>
        <div v-else>
            <button class="glass-card collapse-button" @click="collapse = !collapse">
                <i class="test fa-solid fa-bars"></i>
            </button>
        </div>
    </Transition>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: 280px;
  min-width: 280px;
  padding: 2rem;
}

.logo {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}

.logo #robot {
  font-size: 2rem;
  color: var(--primary-color-light);
  margin-right: 1rem;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 700;
}

.collapse-button {
    margin-left: auto;
    margin-right: 0;
    padding: 0.5rem;
    font-size: 1.5rem;
    background-color: var(--bg-color-light);
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

.nav-links li a:hover, .nav-links li a.active {
  background-color: var(--primary-color);
  color: white;
}

.nav-links li a i {
  margin-right: 1rem;
  font-size: 1.2rem;
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

.logout a {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  text-decoration: none;
  color: var(--text-color);
  transition: background-color 0.3s ease, color 0.3s ease;
}

.logout a:hover {
    background-color: #ef4444;
    color: white;
}

.logout a i {
    margin-right: 1rem;
}

.item {
    display: flex;
}

.item i {
    visibility: hidden;
    margin-left: auto;
}

.item:hover i {
    visibility: visible;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(-30px);
  opacity: 0;
}
</style>
