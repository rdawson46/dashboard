<script setup>
import { ref, computed } from 'vue';
import { RouterLink } from 'vue-router';
import { toast } from 'vue3-toastify';
import 'vue3-toastify/dist/index.css';

import { router } from '@/router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const passwordFieldType = ref('password');
const passwordTouched = ref(false);


const rules = ref({
  length: false,
  uppercase: false,
  special: false,
  number: false
});


const validatePassword = () => {
  passwordTouched.value = true;
  rules.value.length = password.value.length >= 8;
  rules.value.uppercase = /[A-Z]/.test(password.value);
  rules.value.number = /[0-9]/.test(password.value);
  rules.value.special = /[!@#$%^&*(),.?":{}|<>]/.test(password.value);
};


const passwordsMismatch = computed(() => {
  return password.value && confirmPassword.value && password.value !== confirmPassword.value;
});


const passwordFieldIcon = computed(() => {
  return passwordFieldType.value === 'password' ? 'fa-eye' : 'fa-eye-slash';
});


const togglePasswordVisibility = () => {
  passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
};


function notify(message) {
    toast(message, {
        autoClose: 1000,
        theme: 'dark',
    });
}


const register = async () => {
  if (passwordsMismatch.value) {
    notify("Passwords do not match!");
    return;
  }
  if (!rules.value.length || !rules.value.uppercase || !rules.value.special || !rules.value.number) {
    notify("Password does not meet the requirements.");
    return;
  }

  try {
    const formData = new FormData();
    formData.append('username', username.value)
    formData.append('password', password.value)
    formData.append('confirm', confirmPassword.value)

    const res = await fetch('/api/register', {
        method: 'POST',
        credentials: 'include',
        body: formData
      }
    )

    if (!res.ok) {
      notify(`Error occured`)
      return
    }

    const data = await res.json()
    authStore.setUser(data.user)
    router.push('/chat')

  } catch (e) {
    notify(`Error occured`)
    console.error(e)
  }
};
</script>

<template>
  <div class="register-page">
    <div class="glass-card register-box">
      <h1 class="title">Create Account</h1>
      <form @submit.prevent="register">
        <div class="input-group">
          <input type="text" id="username" v-model="username" placeholder="Username" required>
        </div>
        <div class="input-group">
          <input :type="passwordFieldType" id="password" v-model="password" @input="validatePassword" placeholder="Password" required>
          <i :class="['fa-solid', passwordFieldIcon, 'password-toggle']" @click="togglePasswordVisibility"></i>
        </div>
        <div class="input-group">
          <input :type="passwordFieldType" id="confirm-password" v-model="confirmPassword" placeholder="Confirm Password" required>
        </div>
        <div class="password-rules" v-if="passwordTouched">
          <p :class="{ 'valid': rules.length }">8 characters minimum</p>
          <p :class="{ 'valid': rules.number }">One digit</p>
          <p :class="{ 'valid': rules.uppercase }">One uppercase letter</p>
          <p :class="{ 'valid': rules.special }">One special character</p>
        </div>
        <p v-if="passwordsMismatch && passwordTouched" class="error-message">Passwords do not match.</p>
        <button type="submit" class="register-button">Register</button>
      </form>
      <p class="login-link">Already have an account? <router-link to="/login">Login here</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100%;
}

.register-box {
  width: 100%;
  max-width: 450px;
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

.password-rules {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  text-align: left;
  font-size: 0.8rem;
  margin-bottom: 1.5rem;
  color: var(--text-color);
  opacity: 0.8;
}

.password-rules p {
  transition: color 0.3s ease;
}

.password-rules p.valid {
  color: var(--secondary-color);
  text-shadow: 0 0 5px var(--secondary-color);
}

.error-message {
  color: #ef4444;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.register-button {
  width: 100%;
  margin-top: 1rem;
}

.login-link {
  margin-top: 2rem;
  font-size: 0.9rem;
  color: var(--text-color);
  opacity: 0.8;
}

.login-link a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.login-link a:hover {
  text-decoration: underline;
}
</style>
