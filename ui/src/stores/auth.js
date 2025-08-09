import { defineStore } from "pinia";

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        isAuthenticated: false
    }),

    actions: {
        async fetchUser() {
            // TODO: implement this on the go end
            const res = await fetch('/api/me', { credentials: 'include' })

            if (!res.ok) {
                this.user = null
                this.isAuthenticated = false
                return
            }

            const data = await res.json()

            this.user = data.user
            this.isAuthenticated = true
        },

        setUser(user) {
            this.user = user
            this.isAuthenticated = true
        },

        logout() {
            this.user = null
            this.isAuthenticated = false
        }
    }
})
