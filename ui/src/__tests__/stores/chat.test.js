import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useChatStore } from '@/stores/chat.js'

describe('useChatStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('state', () => {
    it('should have empty history initially', () => {
      const store = useChatStore()
      expect(store.history).toEqual([])
    })
  })

  describe('actions', () => {
    it('should add chat to history', () => {
      const store = useChatStore()
      const chat = { id: 1, description: 'Test Chat' }
      
      store.addChatToHistory(chat)
      
      expect(store.history).toHaveLength(1)
      expect(store.history[0]).toEqual(chat)
    })

    it('should add chat to beginning of history', () => {
      const store = useChatStore()
      store.addChatToHistory({ id: 1, description: 'First Chat' })
      store.addChatToHistory({ id: 2, description: 'Second Chat' })
      
      expect(store.history[0].id).toBe(2)
      expect(store.history[1].id).toBe(1)
    })

    it('should set entire history', () => {
      const store = useChatStore()
      const history = [
        { id: 1, description: 'Chat 1' },
        { id: 2, description: 'Chat 2' },
        { id: 3, description: 'Chat 3' }
      ]
      
      store.setHistory(history)
      
      expect(store.history).toEqual(history)
      expect(store.history).toHaveLength(3)
    })

    it('should replace existing history', () => {
      const store = useChatStore()
      store.addChatToHistory({ id: 1, description: 'Old Chat' })
      
      store.setHistory([{ id: 2, description: 'New Chat' }])
      
      expect(store.history).toHaveLength(1)
      expect(store.history[0].id).toBe(2)
    })
  })
})
