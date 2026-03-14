import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import NotFoundPage from '@/views/NotFoundPage.vue'

describe('NotFoundPage.vue', () => {
  it('should render 404 page', () => {
    const wrapper = mount(NotFoundPage)
    expect(wrapper.find('.main').exists()).toBe(true)
  })

  it('should display 404 heading', () => {
    const wrapper = mount(NotFoundPage)
    expect(wrapper.find('h1').text()).toContain('404')
  })

  it('should display page not found message', () => {
    const wrapper = mount(NotFoundPage)
    expect(wrapper.find('p').text()).toBe('The page you are looking for does not exist.')
  })

  it('should have sad face icon', () => {
    const wrapper = mount(NotFoundPage)
    expect(wrapper.find('.logo').exists()).toBe(true)
  })
})
