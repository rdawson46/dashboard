import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '../views/LandingPage.vue'
import LoginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import ChatPage from '../views/ChatPage.vue'
import NotFoundPage from '../views/NotFoundPage.vue'

import { useAuthStore } from '@/stores/auth';

const routes = [
    { path: '/', component: LandingPage },
    {
        path: '/login',
        component: LoginPage,
        name: "Login"
    },
    { path: '/register', component: RegisterPage },
    {
        path: '/chat',
        component: ChatPage,
        name: "New Chat",
        meta: { requiresAuth: true }
    },
    {
        path: '/chat/:id',
        component: ChatPage,
        name: "Existing Chat",
        props: true,
        meta: { requiresAuth: true }
    },
    { path: '/:pathMatch(.*)*', component: NotFoundPage }
]

export const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(async (to) => {
    const authStore = useAuthStore()

    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        await authStore.fetchUser()

        if (!authStore.isAuthenticated) {
            return { path: '/login' }
        }
    }
})
