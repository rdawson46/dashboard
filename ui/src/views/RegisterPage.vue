<script setup>
import { ref, computed } from 'vue';
import { RouterLink } from 'vue-router';

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

const register = () => {
  if (passwordsMismatch.value) {
    alert("Passwords do not match!");
    return;
  }
  if (!rules.value.length || !rules.value.uppercase || !rules.value.special || !rules.value.number) {
    alert("Password does not meet the requirements.");
    return;
  }
  // Handle registration logic here
  alert(`Registering with username: ${username.value}`);
};
</script>

<template>
  <div class="register-container">
    <div class="register-box">
      <h1 class="title">Register</h1>
      <form @submit.prevent="register">
        <div class="input-group">
          <label for="username">Username</label>
          <input type="text" id="username" v-model="username" required>
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input :type="passwordFieldType" id="password" v-model="password" @input="validatePassword" required>
          <i :class="['fa-solid', passwordFieldIcon, 'password-toggle']" @click="togglePasswordVisibility"></i>
        </div>
        <div class="input-group">
          <label for="confirm-password">Confirm Password</label>
          <input :type="passwordFieldType" id="confirm-password" v-model="confirmPassword" @input="passwordTouched = true" required>
        </div>
        <div class="password-rules">
          <p :class="{ 'valid': rules.length && passwordTouched }">8 characters minimum</p>
          <p :class="{ 'valid': rules.number && passwordTouched }">One digit</p>
          <p :class="{ 'valid': rules.uppercase && passwordTouched }">One uppercase letter</p>
          <p :class="{ 'valid': rules.special && passwordTouched }">One special character</p>
        </div>
        <p v-if="passwordsMismatch && passwordTouched" class="error-message">Passwords do not match.</p>
        <button type="submit" class="register-button">Register</button>
      </form>
      <p class="login-link">Already have an account? <router-link to="/login">Login here</router-link></p>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100vw;
  background: linear-gradient(-45deg, #0f2027, #203a43, #2c5364, #1a2980);
  color: white;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.register-box {
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
  box-sizing: border-box; /* Added to fix width issue */
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

.register-button {
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

.register-button:hover {
  background-color: #0056b3;
}

.login-link {
  margin-top: 1.5rem;
  font-size: 0.9rem;
}

.login-link a {
  color: #007bff;
  text-decoration: none;
  font-weight: 600;
}
</style>
