<script setup>
import { ref, onMounted } from 'vue';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

const history = ref([]);

function notify(message) {
    toast(message, {
        autoClose: 1000,
        theme: 'dark',
    });
}

async function getHistory() {
  try{
    const res = await fetch('/api/chatDescription', { credentials: 'include' });
    const data = await res.json();
    return data;
  } catch (e) {
    notify(`Can't get chat history`)
    throw e
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
  <nav class="sidebar glass-card">
    <div class="logo">
      <i class="fa-solid fa-robot"></i>
      <h1>Chat</h1>
    </div>
    <ul class="nav-links">
      <li><a href="#" class="active"><i class="fa-solid fa-comments"></i><span>Chat</span></a></li>
      <li><a href="#"><i class="fa-solid fa-user"></i><span>Profile</span></a></li>
      <li><a href="/models"><i class="fa-solid fa-cogs"></i><span>Models</span></a></li>
    </ul>
    <div class="history">
      <h3>History</h3>
      <ul>
        <li v-for="h in history" :key="h.id">{{ h.description }}</li>
      </ul>
    </div>
    <div class="logout">
        <a href="/api/logout">
            <i class="fa-solid fa-right-from-bracket"></i>
            <span>Logout</span>
        </a>
    </div>
  </nav>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: 280px;
  padding: 2rem;
}

.logo {
  display: flex;
  align-items: center;
  margin-bottom: 3rem;
}

.logo i {
  font-size: 2rem;
  color: var(--primary-color);
  margin-right: 1rem;
}

.logo h1 {
  font-size: 1.5rem;
  font-weight: 700;
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
}

.history h3 {
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
  padding: 0.5rem 0;
  cursor: pointer;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.history ul li:hover {
  opacity: 1;
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

</style>
