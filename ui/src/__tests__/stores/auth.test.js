import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth.js'

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('state', () => {
    it('should have null user initially', () => {
      const store = useAuthStore()
      expect(store.user).toBeNull()
    })

    it('should have isAuthenticated as false initially', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('getters', () => {
    it('should return null username when user is null', () => {
      const store = useAuthStore()
      expect(store.username).toBeNull()
    })

    it('should return null id when user is null', () => {
      const store = useAuthStore()
      expect(store.id).toBeNull()
    })

    it('should return username when user is set', () => {
      const store = useAuthStore()
      store.setUser({ id: 1, username: 'testuser' })
      expect(store.username).toBe('testuser')
    })

    it('should return id when user is set', () => {
      const store = useAuthStore()
      store.setUser({ id: 42, username: 'testuser' })
      expect(store.id).toBe(42)
    })
  })

  describe('actions', () => {
    it('should set user and mark as authenticated', () => {
      const store = useAuthStore()
      const user = { id: 1, username: 'testuser' }
      
      store.setUser(user)
      
      expect(store.user).toEqual(user)
      expect(store.isAuthenticated).toBe(true)
    })

    it('should logout and clear user', () => {
      const store = useAuthStore()
      store.setUser({ id: 1, username: 'testuser' })
      
      store.logout()
      
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })

    it('should fetchUser successfully when authenticated', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ id: 1, username: 'testuser' })
      })
      
      const store = useAuthStore()
      await store.fetchUser()
      
      expect(store.user).toEqual({ id: 1, username: 'testuser' })
      expect(store.isAuthenticated).toBe(true)
    })

    it('should handle fetchUser when not authenticated', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false
      })
      
      const store = useAuthStore()
      await store.fetchUser()
      
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })

    it('should handle fetchUser network error', async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error('Network error'))
      
      const store = useAuthStore()
      await store.fetchUser()
      
      expect(store.user).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    })
  })
})
