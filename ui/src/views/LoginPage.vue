<script setup>
import { ref, computed } from 'vue';
import { RouterLink } from 'vue-router';

const username = ref('');
const password = ref('');
const passwordFieldType = ref('password');
const passwordTouched = ref(false);

const passwordFieldIcon = computed(() => {
  return passwordFieldType.value === 'password' ? 'fa-eye' : 'fa-eye-slash';
});

const togglePasswordVisibility = () => {
  passwordFieldType.value = passwordFieldType.value === 'password' ? 'text' : 'password';
};

const login = () => {
  alert(`Logging in with username: ${username.value}`);
};
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <h1 class="title">Login</h1>
      <form @submit.prevent="login">
        <div class="input-group">
          <label for="username">Username</label>
          <input type="text" id="username" v-model="username" required>
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input :type="passwordFieldType" id="password" v-model="password" @input="validatePassword" required>
          <i :class="['fa-solid', passwordFieldIcon, 'password-toggle']" @click="togglePasswordVisibility"></i>
        </div>
        <button type="submit" class="login-button">Login</button>
      </form>
      <p class="register-link">Don't have an account? <router-link to="/register">Register here</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100vw; 
  background: linear-gradient(-45deg, #0f2027, #203a43, #2c5364, #1a2980);
  color: white;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.login-box {
  background: rgba(0, 0, 0, 0.3);
  padding: 2.5rem;
  border-radius: 15px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  width: 100%;
  max-width: 400px;
  text-align: center;
  animation: fade-in 0.5s ease-out forwards;
}

@keyframes fade-in {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

.title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
}

.input-group {
  position: relative;
  margin-bottom: 1.5rem;
  text-align: left;
}

.input-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  opacity: 0.8;
}

.input-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: white;
  font-size: 1rem;
  transition: border-color 0.2s, box-shadow 0.2s;
  box-sizing: border-box; 
}

.input-group input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
}

.password-toggle {
  position: absolute;
  top: 70%;
  right: 1rem;
  transform: translateY(-50%);
  cursor: pointer;
  opacity: 0.6;
}

.password-rules {
  text-align: left;
  font-size: 0.8rem;
  margin-bottom: 1rem;
  color: #ff6b6b;
}

.password-rules p.valid {
  color: #63e6be;
}

.error-message {
  color: #ff6b6b;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.login-button {
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  background-color: #007bff;
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.login-button:hover {
  background-color: #0056b3;
}

.register-link {
  margin-top: 1.5rem;
  font-size: 0.9rem;
}

.register-link a {
  color: #007bff;
  text-decoration: none;
  font-weight: 600;
}
</style>
