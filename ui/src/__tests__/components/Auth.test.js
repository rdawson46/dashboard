import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Auth from '@/components/Auth.vue'

describe('Auth.vue', () => {
  it('should render welcome header', () => {
    const wrapper = mount(Auth)
    expect(wrapper.text()).toContain('Welcome to the Chat App')
  })

  it('should render welcome subheader', () => {
    const wrapper = mount(Auth)
    expect(wrapper.text()).toContain('Please login or register to continue')
  })

  it('should render login button', () => {
    const wrapper = mount(Auth)
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(2)
    expect(buttons[0].text()).toBe('Login')
  })

  it('should render register button', () => {
    const wrapper = mount(Auth)
    const buttons = wrapper.findAll('button')
    expect(buttons[1].text()).toBe('Register')
  })

  it('should have auth-container class', () => {
    const wrapper = mount(Auth)
    expect(wrapper.find('.auth-container').exists()).toBe(true)
  })

  it('should have button-container with two buttons', () => {
    const wrapper = mount(Auth)
    expect(wrapper.find('.button-container').exists()).toBe(true)
    expect(wrapper.findAll('.auth-button')).toHaveLength(2)
  })
})
