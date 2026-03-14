import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FileOptionsMenu from '@/components/FileOptionsMenu.vue'

describe('FileOptionsMenu.vue', () => {
  it('should render menu when open is true', () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: true,
        rag: false
      }
    })
    expect(wrapper.find('.options-menu').exists()).toBe(true)
  })

  it('should not render menu when open is false', () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: false,
        rag: false
      }
    })
    expect(wrapper.find('.options-menu').exists()).toBe(false)
  })

  it('should emit trigger-upload when upload button clicked', async () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: true,
        rag: false
      }
    })
    
    await wrapper.findAll('button')[0].trigger('click')
    
    expect(wrapper.emitted('trigger-upload')).toBeTruthy()
  })

  it('should emit trigger-select-files when select files button clicked', async () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: true,
        rag: false
      }
    })
    
    await wrapper.findAll('button')[1].trigger('click')
    
    expect(wrapper.emitted('trigger-select-files')).toBeTruthy()
  })

  it('should emit update:rag when checkbox is toggled', async () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: true,
        rag: false
      }
    })
    
    await wrapper.find('input[type="checkbox"]').trigger('change')
    
    expect(wrapper.emitted('update:rag')).toBeTruthy()
  })

  it('should have RAG toggle', () => {
    const wrapper = mount(FileOptionsMenu, {
      props: {
        open: true,
        rag: false
      }
    })
    expect(wrapper.find('.rag-toggle').exists()).toBe(true)
  })
})
