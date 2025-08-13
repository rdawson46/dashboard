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
    console.log(data)
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
  <nav class="navbar">
    <ul class="navbar-nav">
      <li class="logo">
        <div class="nav-link">
          <span class="link-text logo-text">Chat</span>
          <i class="fa-solid fa-angles-right fa-2xl"></i>
        </div>
      </li>
      <li class="nav-item">
        <div class="nav-link">
          <span class="link-text">Profile</span>
        </div>
      </li>
      <li class="nav-item">
        <div class="nav-link">
          <span class="link-text"><a href="/models">Model</a></span>
        </div>
      </li>

      <li class="nav-item">
        <div class="nav-link">
          <span class="link-text">Chats</span>
        </div>

      </li>
      <ul class="desc-list">

        <!-- TODO: make this not bad and add a delete button -->
        <!-- TODO: and paginate this api -->
        <template v-for="h in history">
          <li class="description">
            {{h.description}}
          </li>
        </template>

      </ul>
    </ul>
  </nav>
</template>

<style scoped>
.navbar {
  position: fixed;
  background-color: #2c3e50;
  transition: width 200ms ease;
  z-index: 1000;
  width: 16rem;
}

.navbar-nav {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100vh;
}

.nav-item {
  width: 100%;
}

.nav-link {
  display: flex;
  align-items: center;
  height: 5rem;
  color: #bdc3c7;
  text-decoration: none;
  filter: grayscale(100%) opacity(0.7);
  transition: all 200ms ease;
}

.nav-link:hover {
  filter: grayscale(0%) opacity(1);
  background: #34495e;
  color: #ecf0f1;
}

.link-text {
  margin-left: 1rem;
  white-space: nowrap;
  font-size: 1.2rem;
}

.nav-link i {
  min-width: 2rem;
  margin: 0 1.5rem;
  transition: 200ms ease;
}

.logo {
  font-weight: bold;
  text-transform: uppercase;
  margin-bottom: 1rem;
  text-align: center;
  color: #ecf0f1;
  background: #1a252f;
  font-size: 1.5rem;
  letter-spacing: 0.3ch;
  width: 100%;
}

.logo i {
  transform: rotate(-180deg);
  transition: transform 200ms ease;
}

.navbar:hover .logo i {
  transform: rotate(0);
}

@media only screen and (min-width: 600px) {
  .navbar {
    top: 0;
    height: 100vh;
  }

  .navbar:hover .link-text {
    display: inline;
  }
}

@media only screen and (max-width: 600px) {
  .navbar {
    bottom: 0;
    width: 100vw;
    height: 5rem;
  }

  .logo {
    display: none;
  }

  .navbar-nav {
    flex-direction: row;
  }

  .nav-link {
    justify-content: center;
  }

  main {
    margin: 0;
  }
}
</style>
