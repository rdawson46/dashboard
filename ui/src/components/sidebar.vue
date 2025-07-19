<template>
  <nav class="navbar">
    <ul class="navbar-nav">
      <li class="logo">
        <a href="#" class="nav-link">
          <span class="link-text logo-text">Gemini Pro</span>
          <i class="fa-solid fa-angles-right fa-2xl"></i>
        </a>
      </li>
      <li class="nav-item">
        <a href="#" class="nav-link">
          <i class="fa-solid fa-house fa-2xl"></i>
          <span class="link-text">Dashboard</span>
        </a>
      </li>
      <li class="nav-item">
        <a href="#" class="nav-link">
          <i class="fa-solid fa-comment fa-2xl"></i>
          <span class="link-text">Chat</span>
        </a>
      </li>
      <li class="nav-item">
        <a href="#" class="nav-link">
          <i class="fa-solid fa-gear fa-2xl"></i>
          <span class="link-text">Settings</span>
        </a>
      </li>
      <li class="nav-item" id="themeButton">
        <a href="#" class="nav-link">
          <i class="fa-solid fa-moon fa-2xl"></i>
          <span class="link-text">Theme</span>
        </a>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'Sidebar'
};
</script>

<style scoped>
.navbar {
  position: fixed;
  background-color: #2c3e50;
  transition: width 200ms ease;
  z-index: 1000;
}

.navbar-nav {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.nav-item {
  width: 100%;
}

.nav-item:last-child {
  margin-top: auto;
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
  display: none;
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
  transform: rotate(0deg);
  transition: transform 200ms ease;
}

.navbar:hover .logo i {
  transform: rotate(-180deg);
}

@media only screen and (min-width: 600px) {
  .navbar {
    top: 0;
    width: 6rem;
    height: 100vh;
  }

  .navbar:hover {
    width: 16rem;
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