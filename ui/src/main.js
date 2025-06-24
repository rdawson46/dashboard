import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

createApp(App).mount("#app")

/*
const app = createApp(App)
import { createPinia } from 'pinia'
import router from './router'
app.use(createPinia())
// TODO: make router
app.use(router)
app.mount("#app")
*/
