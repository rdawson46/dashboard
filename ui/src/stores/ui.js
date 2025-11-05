import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', () => {
  const isSidebarCollapsed = ref(false)

  function toggleSidebar() {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
  }

  const sidebarWidth = computed(() => (isSidebarCollapsed.value ? '100px' : '280px'))
  const sidebarWidthExpanded = '280px';
  const sidebarWidthCollapsed = '100px';


  return { isSidebarCollapsed, toggleSidebar, sidebarWidth, sidebarWidthCollapsed, sidebarWidthExpanded }
})
