import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import RegisterPage from '@/views/RegisterPage.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn(() => vi.fn())
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'Home', component: { template: 'Home' } },
    { path: '/chat', name: 'Chat', component: { template: 'Chat' } },
    { path: '/login', name: 'Login', component: { template: 'Login' } }
  ]
})

describe('RegisterPage.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render register page', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.auth-page').exists()).toBe(true)
  })

  it('should render branding title', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.branding-title').text()).toBe('Create Your Account')
  })

  it('should have username input', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#username').exists()).toBe(true)
  })

  it('should have password input', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#password').exists()).toBe(true)
  })

  it('should have confirm password input', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#confirm-password').exists()).toBe(true)
  })

  it('should have register button', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').text()).toBe('Register')
  })

  it('should have login link', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.sub-link').exists()).toBe(true)
    expect(wrapper.find('.sub-link').text()).toContain('Login here')
  })

  it('should show password rules when password is touched', async () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    
    expect(wrapper.find('.password-rules').exists()).toBe(false)
    
    await wrapper.find('#password').setValue('Test')
    
    expect(wrapper.find('.password-rules').exists()).toBe(true)
  })

  it('should have form', () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('form').exists()).toBe(true)
  })

  it('should show error when passwords mismatch', async () => {
    const wrapper = mount(RegisterPage, {
      global: {
        plugins: [router]
      }
    })
    
    await wrapper.find('#password').setValue('Test123!')
    await wrapper.find('#confirm-password').setValue('Different')
    
    expect(wrapper.find('.error-message').exists()).toBe(true)
  })
})
