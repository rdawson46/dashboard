import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '@/views/LandingPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'Home', component: LandingPage },
    { path: '/login', name: 'Login', component: { template: 'Login' } },
    { path: '/register', name: 'Register', component: { template: 'Register' } }
  ]
})

describe('LandingPage.vue', () => {
  it('should render landing page', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.landing-page').exists()).toBe(true)
  })

  it('should render header', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.header').exists()).toBe(true)
  })

  it('should render logo', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.logo').exists()).toBe(true)
  })

  it('should have login link', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.findAll('.nav-link')).toHaveLength(2)
  })

  it('should render main headline', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.main-headline').text()).toContain('Your Personal AI Teammate')
  })

  it('should render features section', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('#features').exists()).toBe(true)
  })

  it('should have feature cards', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.findAll('.feature-card')).toHaveLength(3)
  })

  it('should have CTA buttons', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.cta-buttons').exists()).toBe(true)
  })

  it('should have Get Started button', () => {
    const wrapper = mount(LandingPage, {
      global: {
        plugins: [router]
      }
    })
    expect(wrapper.find('.cta-button').text()).toBe('Get Started for Free')
  })
})
