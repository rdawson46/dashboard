import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import Chat from '@/components/Chat.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn(() => vi.fn())
}))

vi.mock('@/composables/stream.js', () => ({
  useStream: vi.fn()
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { name: 'New Chat', path: '/chat' },
    { name: 'Existing Chat', path: '/chat/:id' }
  ]
})

describe('Chat.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render chat container', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.chat-container').exists()).toBe(true)
  })

  it('should render chat header', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.chat-header').exists()).toBe(true)
  })

  it('should have Chat title', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.chat-header h2').text()).toBe('Chat')
  })

  it('should have model selector', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#modelSelector').exists()).toBe(true)
  })

  it('should render ChatInput component', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.findComponent({ name: 'ChatInput' }).exists()).toBe(true)
  })

  it('should render welcome message when no messages', () => {
    const wrapper = mount(Chat, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.welcome-container').exists()).toBe(true)
  })
})
