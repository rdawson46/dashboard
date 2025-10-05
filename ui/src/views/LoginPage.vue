<script setup>
import { ref, computed } from 'vue';
import { RouterLink } from 'vue-router';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

import { router } from '@/router'
import { useAuthStore } from '@/stores/auth'

const username = ref('');
const password = ref('');
const passwordFieldType = ref('password');

const authStore = useAuthStore()

const passwordFieldIcon = computed(() => {
  return passwordFieldType.value === 'password' ? 'fa-eye' : 'fa-eye-slash';
});


function notify(message) {
    toast(message, {
        autoClose: 1000,
        theme: 'dark',
    });
}

const togglePasswordVisibility = () => {
  passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
};


const login = async () => {
  if (!username.value.length && !password.value.length) return

  try {
    const formData = new FormData();
    formData.append('username', username.value)
    formData.append('password', password.value)

    const res = await fetch('/api/login', {
        method: 'POST',
        credentials: 'include',
        body: formData
      }
    )

    if (!res.ok) {
      const data = await res.json()
      notify(data.error)
      console.log(res)
      return
    }

    const data = await res.json()
    authStore.setUser(data)
    router.push('/chat')
  } catch (e) {
    notify("An error has occured while logging in")
    console.log(e)
    return
  }
};
</script>

<template>
  <div class="login-page">
    <div class="glass-card login-box">
      <h1 class="title">Login</h1>
      <form @submit.prevent="login">
        <div class="input-group">
          <input type="text" id="username" v-model="username" placeholder="Username" required>
        </div>
        <div class="input-group">
          <input :type="passwordFieldType" id="password" v-model="password" placeholder="Password" required>
          <i :class="['fa-solid', passwordFieldIcon, 'password-toggle']" @click="togglePasswordVisibility"></i>
        </div>
        <button type="submit" class="login-button">Login</button>
      </form>
      <p class="register-link">Don't have an account? <router-link to="/register">Register here</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100%;
}

.login-box {
  width: 100%;
  max-width: 400px;
  text-align: center;
  padding: 3rem;
}

.title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 2rem;
  color: var(--text-color);
}

.input-group {
  position: relative;
  margin-bottom: 1.5rem;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 1rem;
  transform: translateY(-50%);
  cursor: pointer;
  color: var(--text-color);
  opacity: 0.7;
}

.login-button {
  width: 100%;
  margin-top: 1rem;
}

.register-link {
  margin-top: 2rem;
  font-size: 0.9rem;
  color: var(--text-color);
  opacity: 0.8;
}

.register-link a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
