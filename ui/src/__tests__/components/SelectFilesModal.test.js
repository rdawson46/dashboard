import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SelectFilesModal from '@/components/SelectFilesModal.vue'

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn()
}))

global.fetch = vi.fn().mockResolvedValue({
  ok: true,
  json: async () => ({ files: [{ id: 1, file_name: 'test.txt' }] })
})

describe('SelectFilesModal.vue', () => {
  it('should not render when show is false', () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: false, selectedFiles: [] }
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(false)
  })

  it('should render when show is true', () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: true, selectedFiles: [] }
    })
    expect(wrapper.find('.modal-overlay').exists()).toBe(true)
  })

  it('should emit close when close button clicked', async () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: true, selectedFiles: [] }
    })
    
    await wrapper.findAll('button')[0].trigger('click')
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit close when clicking overlay', async () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: true, selectedFiles: [] }
    })
    
    await wrapper.find('.modal-overlay').trigger('click.self')
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should have title', () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: true, selectedFiles: [] }
    })
    expect(wrapper.find('h2').text()).toBe('Select Files')
  })

  it('should have cancel and add buttons', () => {
    const wrapper = mount(SelectFilesModal, {
      props: { show: true, selectedFiles: [] }
    })
    const buttons = wrapper.findAll('.modal-actions button')
    expect(buttons).toHaveLength(2)
  })
})
