import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ChatInput from '@/components/ChatInput.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn()
}))

describe('ChatInput.vue', () => {
  it('should render input container', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.find('.chat-input-container').exists()).toBe(true)
  })

  it('should have centered class when isCentered is true', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: true
      }
    })
    expect(wrapper.find('.chat-input-container.middle').exists()).toBe(true)
  })

  it('should not have centered class when isCentered is false', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.find('.chat-input-container.middle').exists()).toBe(false)
  })

  it('should render message input', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.find('#message-input').exists()).toBe(true)
  })

  it('should have send button', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.find('#send-btn').exists()).toBe(true)
  })

  it('should have file options button', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.findAll('button')).toHaveLength(4)
  })

  it('should emit send-message when send button clicked', async () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    
    const input = wrapper.find('#message-input')
    input.element.innerText = 'Test message'
    await input.trigger('input')
    await wrapper.find('#send-btn').trigger('click')
    
    expect(wrapper.emitted('send-message')).toBeTruthy()
  })

  it('should not emit send-message when input is empty', async () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    
    await wrapper.find('#send-btn').trigger('click')
    
    expect(wrapper.emitted('send-message')).toBeFalsy()
  })

  it('should render web search button', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.text()).toContain('Web Search')
  })

  it('should render code button', () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    expect(wrapper.text()).toContain('Code')
  })

  it('should toggle searchActive on web search button click', async () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    
    const searchButton = wrapper.findAll('button').find(b => b.text().includes('Web Search'))
    await searchButton.trigger('click')
    
    expect(searchButton.classes()).toContain('active')
  })

  it('should toggle codeActive on code button click', async () => {
    const wrapper = mount(ChatInput, {
      props: {
        userMessageHistory: [],
        isCentered: false
      }
    })
    
    const codeButton = wrapper.findAll('button').find(b => b.text().includes('Code'))
    await codeButton.trigger('click')
    
    expect(codeButton.classes()).toContain('active')
  })
})
