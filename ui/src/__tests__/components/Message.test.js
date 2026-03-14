import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import Message from '@/components/Message.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn()
}))

describe('Message.vue', () => {
  it('should render user message', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'user',
        content: '<p>Hello</p>'
      }
    })
    expect(wrapper.find('.message.user').exists()).toBe(true)
  })

  it('should render assistant message', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'assistant',
        content: '<p>Hi there</p>'
      }
    })
    expect(wrapper.find('.message.assistant').exists()).toBe(true)
  })

  it('should show loading spinner when loading is true', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'assistant',
        content: '',
        loading: true
      }
    })
    expect(wrapper.find('.spinner').exists()).toBe(true)
  })

  it('should not show spinner when content exists', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'assistant',
        content: '<p>Some content</p>',
        loading: true
      }
    })
    expect(wrapper.find('.spinner').exists()).toBe(false)
  })

  it('should render info message when role is info', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'info',
        content: '<p>Info content</p>',
        details: { total_duration: 1000000000, eval_count: 10 }
      }
    })
    expect(wrapper.find('.info-message').exists()).toBe(true)
  })

  it('should render tool calls when provided', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'assistant',
        content: '',
        tool_calls: [
          {
            id: 'call_1',
            function: { name: 'test_function', arguments: '{}' }
          }
        ]
      }
    })
    expect(wrapper.find('.tool-calls-list').exists()).toBe(true)
  })

  it('should toggle tool call expansion', async () => {
    const wrapper = mount(Message, {
      props: {
        role: 'assistant',
        content: '',
        tool_calls: [
          {
            id: 'call_1',
            function: { name: 'test_function', arguments: '{}' }
          }
        ]
      }
    })
    
    expect(wrapper.find('.tool-call-details').exists()).toBe(false)
    
    await wrapper.find('.tool-call-header').trigger('click')
    
    expect(wrapper.find('.tool-call-details').exists()).toBe(true)
  })

  it('should not render when role is info and no content', () => {
    const wrapper = mount(Message, {
      props: {
        role: 'info',
        content: '',
        details: { total_duration: 1000000000, eval_count: 10 }
      }
    })
    expect(wrapper.find('.message').exists()).toBe(false)
  })
})
