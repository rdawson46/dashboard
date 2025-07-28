import { createRouter, createWebHistory } from 'vue-router'

// TODO: define the views
const routes = [
    { path: '/', component: LandingPage },
    { path: '/login', component: LoginPage },
    { path: '/register', component: RegisterPage },
    { 
        path: '/chat',
        component: ChatPage,
        meta: { requiresAuth: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    // TODO: swap this out later
    // will have to reach out to server to verify
    const token = localStorage.get('token')

    if ((to.path === '/login' || to.path === '/register') && token) {
        return next('/dashboard')
    }

    if (to.meta.requiresAuth && !token) {
        return next('/login')
    }

    next()
})

export default router
