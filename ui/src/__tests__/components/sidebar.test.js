import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import sidebar from '@/components/sidebar.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn()
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { name: 'New Chat', path: '/chat' },
    { name: 'Files', path: '/files' },
    { name: 'Models', path: '/models' },
    { name: 'jobs', path: '/jobs' }
  ]
})

describe('sidebar.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should render sidebar', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.sidebar').exists()).toBe(true)
  })

  it('should render logo with robot icon', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#robot').exists()).toBe(true)
  })

  it('should render navigation links', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.findAll('.nav-links a')).toHaveLength(4)
  })

  it('should have New Chat link', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.text()).toContain('New Chat')
  })

  it('should have Files link', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.text()).toContain('Files')
  })

  it('should have Models link', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.text()).toContain('Models')
  })

  it('should have Jobs link', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.text()).toContain('Jobs')
  })

  it('should have logout link', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.logout-link').exists()).toBe(true)
  })

  it('should have collapse button', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.collapse-button').exists()).toBe(true)
  })

  it('should render history section', () => {
    const wrapper = mount(sidebar, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.history-section').exists()).toBe(true)
  })
})
