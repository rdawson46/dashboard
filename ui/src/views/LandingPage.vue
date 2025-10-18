<template>
  <div class="landing-page">
    <header class="header">
      <div class="logo">Chat</div>
      <nav class="nav">
        <router-link to="/login" class="nav-link">Login</router-link>
        <router-link to="/register" class="nav-link primary">Register</router-link>
      </nav>
    </header>

    <section class="hero-section">
      <div class="hero-content">
        <h1 class="main-headline">Your Personal AI Teammate</h1>
        <h2 class="sub-headline">Automate, Create, and Accomplish More Than Ever Before.</h2>
        <p class="description">
          Chat is a new kind of intelligent assistant that can understand your goals, connect with your tools, and work alongside you to get things done.
        </p>
        <div class="cta-buttons">
          <router-link to="/register" class="cta-button">Get Started for Free</router-link>
          <a href="#features" class="cta-button secondary">Learn More</a>
        </div>
      </div>
      <div class="hero-image">
        <!-- Placeholder for a graphic -->
        <div class="glowing-orb"></div>
      </div>
    </section>

    <section id="features" class="features-section">
      <h2 class="section-title">A Tool for Every Task</h2>
      <p class="section-subtitle">Chat integrates with your workflow to provide seamless assistance.</p>
      <div class="features-grid">
        <div class="feature-card">
          <div class="feature-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
          </div>
          <h3 class="feature-title">Intelligent Chat</h3>
          <p class="feature-description">Converse naturally about complex tasks and ideas. Chat understands context and provides insightful responses.</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"></path></svg>
          </div>
          <h3 class="feature-title">Workspace Integration</h3>
          <p class="feature-description">Connects to your data, reads files, and understands your project structure to provide relevant help.</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 8V4H8"></path><path d="M8 4H4v4"></path><path d="M4 8v4"></path><path d="M4 12v4h4"></path><path d="M8 16h4"></path><path d="M12 16h4v-4"></path><path d="M16 12v-4"></path><path d="M16 8h-4"></path><path d="M12 12h4v4h-4z"></path></svg>
          </div>
          <h3 class="feature-title">Tool-Powered Actions</h3>
          <p class="feature-description">Go beyond text. Chat can execute programs, search the web, and interact with your data.</p>
        </div>
      </div>
    </section>

    <section class="final-cta-section">
        <h2 class="section-title">Ready to Build Faster?</h2>
        <p class="cta-description">Stop switching contexts. Start building with an AI that understands your code.</p>
        <router-link to="/register" class="cta-button large">Sign Up & Start Working</router-link>
    </section>

    <footer class="footer">
    </footer>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>

<style scoped>
:root {
  --primary-color: #4A90E2;
  --background-color: #121212;
  --text-color: #EAEAEA;
  --card-background: rgba(255, 255, 255, 0.05);
  --border-color: rgba(255, 255, 255, 0.1);
}

.landing-page {
  background-color: var(--background-color);
  color: var(--text-color);
  width: 100%;
  overflow-x: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  position: fixed;
  width: 100%;
  top: 0;
  left: 0;
  z-index: 1000;
  background: rgba(18, 18, 18, 0.8);
  backdrop-filter: blur(10px);
}

.logo {
  font-size: 1.5rem;
  font-weight: 600;
}

.nav {
  display: flex;
  gap: 1rem;
}

.nav-link {
  color: var(--text-color);
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  transition: background-color 0.3s ease;
}

.nav-link:hover {
  background-color: var(--card-background);
}

.nav-link.primary {
  background-color: var(--primary-color);
  color: white;
}

.hero-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 100vh;
  padding: 0 5%;
  padding-top: 80px; /* Header height */
  text-align: left;
}

.hero-content {
  max-width: 600px;
  animation: fade-in-up 1s ease-out forwards;
}

.hero-image {
  position: relative;
  width: 400px;
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.glowing-orb {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle, var(--primary-color) 0%, rgba(74,144,226,0) 60%);
  animation: pulse 4s infinite ease-in-out;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(0.95);
    opacity: 0.7;
  }
  50% {
    transform: scale(1.05);
    opacity: 1;
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.main-headline {
  font-size: 3.8rem;
  font-weight: 700;
  margin-bottom: 1rem;
  line-height: 1.2;
}

.sub-headline {
  font-size: 1.6rem;
  font-weight: 400;
  margin-bottom: 2rem;
  opacity: 0.8;
}

.description {
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 3rem;
  opacity: 0.7;
}

.cta-buttons {
  display: flex;
  gap: 1rem;
}

.cta-button {
  padding: 0.9rem 2.2rem;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  background-color: var(--primary-color);
  color: white;
  border: 1px solid var(--primary-color);
}

.cta-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 25px rgba(74, 144, 226, 0.4);
}

.cta-button.secondary {
  background-color: transparent;
  color: var(--primary-color);
}

.cta-button.secondary:hover {
  background-color: var(--primary-color);
  color: white;
}

.cta-button.large {
  padding: 1rem 3rem;
  font-size: 1.2rem;
}

/* Features Section */
.features-section {
  padding: 6rem 5%;
  text-align: center;
}

.section-title {
  font-size: 2.8rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.section-subtitle {
  font-size: 1.2rem;
  opacity: 0.7;
  margin-bottom: 4rem;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.feature-card {
  background: var(--card-background);
  padding: 2.5rem;
  border-radius: 12px;
  border: 1px solid var(--border-color);
  text-align: left;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.feature-icon {
  color: var(--primary-color);
  margin-bottom: 1.5rem;
}

.feature-icon svg {
  width: 32px;
  height: 32px;
}

.feature-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.feature-description {
  font-size: 1rem;
  line-height: 1.6;
  opacity: 0.7;
}

/* Final CTA Section */
.final-cta-section {
  padding: 6rem 5%;
  text-align: center;
  background: linear-gradient(to top, rgba(74,144,226,0.1), transparent);
}

.cta-description {
  font-size: 1.2rem;
  opacity: 0.7;
  margin-bottom: 2.5rem;
}

/* Footer */
.footer {
  text-align: center;
  padding: 2rem 0;
  border-top: 1px solid var(--border-color);
  font-size: 0.9rem;
  opacity: 0.6;
}

@media (max-width: 900px) {
  .hero-section {
    flex-direction: column;
    text-align: center;
    padding-top: 120px;
  }
  .hero-content {
    order: 2;
  }
  .hero-image {
    order: 1;
    width: 300px;
    height: 300px;
    margin-bottom: 2rem;
  }
  .main-headline {
    font-size: 3rem;
  }
  .cta-buttons {
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .header {
    padding: 1rem;
  }
  .nav {
    display: none; /* Simple responsive approach */
  }
}
</style>
