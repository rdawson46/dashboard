import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '@/views/LoginPage.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn(() => vi.fn())
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'Home', component: { template: 'Home' } },
    { path: '/chat', name: 'Chat', component: { template: 'Chat' } },
    { path: '/register', name: 'Register', component: { template: 'Register' } }
  ]
})

describe('LoginPage.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render login page', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.auth-page').exists()).toBe(true)
  })

  it('should render branding title', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.branding-title').text()).toBe('Welcome Back')
  })

  it('should have username input', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#username').exists()).toBe(true)
  })

  it('should have password input', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#password').exists()).toBe(true)
  })

  it('should have login button', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').text()).toBe('Login')
  })

  it('should have register link', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.sub-link').exists()).toBe(true)
    expect(wrapper.find('.sub-link').text()).toContain('Register here')
  })

  it('should have password toggle icon', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.password-toggle').exists()).toBe(true)
  })

  it('should toggle password visibility', async () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    
    const passwordInput = wrapper.find('#password')
    expect(passwordInput.attributes('type')).toBe('password')
    
    await wrapper.find('.password-toggle').trigger('click')
    
    expect(passwordInput.attributes('type')).toBe('text')
  })

  it('should have form', () => {
    const wrapper = mount(LoginPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('form').exists()).toBe(true)
  })
})
