<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="landing-container">
    <div class="hero-content">
      <h1 class="main-headline">Unleash Your Potential. The Future of Assistance is Here.</h1>
      <h2 class="sub-headline">Go Beyond Chat. Leverage the Power of an LLM with Real-World Tools.</h2>
      <p class="description">
        This is not just another chatbot. This is your new digital partner, an advanced Large Language Model with the unique ability to use tools on your behalf. Automate complex tasks, supercharge your workflow, and unlock new possibilities. Experience a revolutionary way to interact with technology and transform your daily life.
      </p>
      <div class="cta-buttons">
        <router-link to="/login" class="cta-button login">Login</router-link>
        <router-link to="/register" class="cta-button register">Register</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes gradient-animation {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes float-animation {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.landing-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100%;
  background: linear-gradient(-45deg, #0f2027, #203a43, #2c5364, #1a2980);
  background-size: 400% 400%;
  animation: gradient-animation 15s ease infinite;
  color: white;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  text-align: center;
  overflow: hidden;
}

.landing-container::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    backdrop-filter: blur(5px);
}

.hero-content {
  max-width: 800px;
  padding: 2rem;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 15px;
  backdrop-filter: blur(5px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: float-animation 6s ease-in-out infinite, fade-in 1s ease-out forwards;
  z-index: 10;
}

.main-headline {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
}

.sub-headline {
  font-size: 1.5rem;
  font-weight: 300;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.description {
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 2.5rem;
  opacity: 0.8;
}

.cta-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.cta-button {
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 25px;
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  color: white;
  cursor: pointer;
  transition: transform 0.2s ease, background-color 0.2s ease;
}

.cta-button:hover {
  transform: translateY(-3px);
}

.login {
  background-color: #007bff;
}

.login:hover {
  background-color: #0056b3;
}

.register {
  background-color: transparent;
  border: 2px solid #007bff;
}

.register:hover {
  background-color: #007bff;
}
</style>

