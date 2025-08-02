import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '../views/LandingPage.vue'
import LoginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import ChatPage from '../views/ChatPage.vue'
import NotFoundPage from '../views/NotFoundPage.vue'

const routes = [
    { path: '/', component: LandingPage },
    { path: '/login', component: LoginPage },
    { path: '/register', component: RegisterPage },
    { path: '/chat', component: ChatPage },
    { path: '/:pathMatch(.*)*', component: NotFoundPage }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router
