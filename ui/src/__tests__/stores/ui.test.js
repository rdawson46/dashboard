import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUiStore } from '@/stores/ui.js'

describe('useUiStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('state', () => {
    it('should have sidebar collapsed initially', () => {
      const store = useUiStore()
      expect(store.isSidebarCollapsed).toBe(true)
    })
  })

  describe('computed', () => {
    it('should return collapsed width when collapsed', () => {
      const store = useUiStore()
      expect(store.sidebarWidth).toBe('80px')
    })

    it('should return expanded width when not collapsed', () => {
      const store = useUiStore()
      store.toggleSidebar()
      expect(store.sidebarWidth).toBe('280px')
    })

    it('should have correct collapsed width constant', () => {
      const store = useUiStore()
      expect(store.sidebarWidthCollapsed).toBe('80px')
    })

    it('should have correct expanded width constant', () => {
      const store = useUiStore()
      expect(store.sidebarWidthExpanded).toBe('280px')
    })
  })

  describe('actions', () => {
    it('should toggle sidebar from collapsed to expanded', () => {
      const store = useUiStore()
      
      store.toggleSidebar()
      
      expect(store.isSidebarCollapsed).toBe(false)
    })

    it('should toggle sidebar from expanded to collapsed', () => {
      const store = useUiStore()
      store.toggleSidebar() // expanded
      
      store.toggleSidebar() // collapsed again
      
      expect(store.isSidebarCollapsed).toBe(true)
    })

    it('should toggle multiple times', () => {
      const store = useUiStore()
      
      store.toggleSidebar()
      store.toggleSidebar()
      store.toggleSidebar()
      
      expect(store.isSidebarCollapsed).toBe(false)
    })
  })
})
