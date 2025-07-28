<template>
  <div class="auth-container">
    <div class="content-wrapper">
      <h1 class="welcome-header">Welcome to the Chat App</h1>
      <p class="welcome-subheader">Please login or register to continue</p>
      <div class="button-container">
        <button class="auth-button">
          Login
        </button>
        <button class="auth-button">
          Register
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  text-align: center;
}
.content-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
}
.welcome-header {
  font-size: 2.25rem; /* 36px */
  font-weight: 700;
  margin-bottom: 0.5rem;
}
.welcome-subheader {
  font-size: 1.125rem; /* 18px */
  margin-bottom: 1.5rem;
}
.button-container {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}
.auth-button {
  background-color: #3b82f6; /* approx blue-500 */
  color: white;
  font-weight: 700;
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s;
}
.auth-button:hover {
  background-color: #2563eb; /* approx blue-700 */
}
</style>
