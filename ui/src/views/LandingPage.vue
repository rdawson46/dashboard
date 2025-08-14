<template>
  <div class="landing-page">
    <div class="glass-card hero-content">
      <h1 class="main-headline">Unleash Your Potential</h1>
      <h2 class="sub-headline">The Future of Assistance is Here.</h2>
      <p class="description">
        Go Beyond Chat. Leverage the Power of an LLM with Real-World Tools. 
        Automate complex tasks, supercharge your workflow, and unlock new possibilities.
      </p>
      <div class="cta-buttons">
        <router-link to="/login" class="cta-button">Login</router-link>
        <router-link to="/register" class="cta-button secondary">Register</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.landing-page {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100%;
  text-align: center;
}

.hero-content {
  max-width: 600px;
  padding: 3rem;
  animation: fade-in 1s ease-out forwards;
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

.main-headline {
  font-size: 3.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: var(--text-color);
}

.sub-headline {
  font-size: 1.5rem;
  font-weight: 400;
  margin-bottom: 2rem;
  color: var(--text-color);
  opacity: 0.8;
}

.description {
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 3rem;
  color: var(--text-color);
  opacity: 0.7;
}

.cta-buttons {
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.cta-button {
  padding: 0.8rem 2rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  background-color: var(--primary-color);
  color: white;
}

.cta-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.cta-button.secondary {
  background-color: transparent;
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
}

.cta-button.secondary:hover {
  background-color: var(--primary-color);
  color: white;
}
</style>