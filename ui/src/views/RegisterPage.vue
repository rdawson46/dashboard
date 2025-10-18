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
      const data = await res.json()
      notify(data.error)
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
  <div class="auth-page">
    <div class="auth-container">
      <div class="branding-panel">
        <h1 class="branding-title">Create Your Account</h1>
        <p class="branding-description">Join our platform and start your journey with a next-generation AI assistant.</p>
      </div>
      <div class="form-panel">
        <div class="form-box glass-card">
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
            <button type="submit" class="auth-button">Register</button>
          </form>
          <p class="sub-link">Already have an account? <router-link to="/login">Login here</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Using styles consistent with LandingPage.vue and LoginPage.vue */
.auth-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100%;
  background-color: var(--background-color);
  color: var(--text-color);
  padding: 2rem 0;
}

.auth-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  width: 100%;
  max-width: 1200px;
  align-items: center;
  gap: 4rem;
  padding: 2rem;
}

.branding-panel {
  text-align: left;
  padding-right: 2rem;
  animation: fade-in-up 0.8s ease-out;
}

.logo {
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary-color);
  text-decoration: none;
  margin-bottom: 2rem;
  display: inline-block;
}

.branding-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  line-height: 1.2;
}

.branding-description {
  font-size: 1.1rem;
  opacity: 0.7;
  line-height: 1.6;
}

.form-panel {
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-box {
  width: 100%;
  max-width: 450px;
  text-align: center;
  padding: 3rem;
  animation: fade-in 1s ease-out forwards;
}

.title {
  font-size: 2.5rem;
  font-weight: 600;
  margin-bottom: 2rem;
}

.input-group {
  position: relative;
  margin-bottom: 1.5rem;
}

input {
  width: 100%;
  padding: 0.9rem 1rem;
  background-color: rgba(255, 255, 255, 0.07);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color);
  font-size: 1rem;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.3);
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 1rem;
  transform: translateY(-50%);
  cursor: pointer;
  opacity: 0.7;
}

.password-rules {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
  text-align: left;
  font-size: 0.8rem;
  margin-bottom: 1.5rem;
  color: rgba(255, 255, 255, 0.6);
}

.password-rules p {
  transition: color 0.3s ease;
  position: relative;
  padding-left: 18px;
}

.password-rules p::before {
  content: '○';
  position: absolute;
  left: 0;
  top: 1px;
  font-size: 12px;
  transition: color 0.3s ease;
}

.password-rules p.valid {
  color: #4ade80; /* A nice green color */
}

.password-rules p.valid::before {
  content: '●';
  color: #4ade80;
}

.error-message {
  color: #ef4444;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  text-align: left;
}

.auth-button {
  width: 100%;
  padding: 0.9rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  background-color: var(--primary-color);
  color: white;
  border: none;
  cursor: pointer;
  margin-top: 1rem;
}

.auth-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(74, 144, 226, 0.4);
}

.sub-link {
  margin-top: 2rem;
  font-size: 0.9rem;
  opacity: 0.8;
}

.sub-link a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.sub-link a:hover {
  text-decoration: underline;
}

@keyframes fade-in {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 900px) {
  .auth-container {
    grid-template-columns: 1fr;
    gap: 2rem;
    text-align: center;
  }
  .branding-panel {
    padding-right: 0;
    max-width: 500px;
    margin: 0 auto;
  }
  .branding-title {
    font-size: 2.5rem;
  }
}
</style>
