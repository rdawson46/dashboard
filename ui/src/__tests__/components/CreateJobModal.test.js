import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import CreateJobModal from '@/components/CreateJobModal.vue'

describe('CreateJobModal.vue', () => {
  it('should render modal', () => {
    const wrapper = mount(CreateJobModal)
    expect(wrapper.find('.modal-backdrop').exists()).toBe(true)
    expect(wrapper.find('.modal-content').exists()).toBe(true)
  })

  it('should render form with job name input', () => {
    const wrapper = mount(CreateJobModal)
    expect(wrapper.find('input#job-name').exists()).toBe(true)
  })

  it('should have job type select', () => {
    const wrapper = mount(CreateJobModal)
    expect(wrapper.find('select#job-type').exists()).toBe(true)
  })

  it('should have frequency select', () => {
    const wrapper = mount(CreateJobModal)
    expect(wrapper.find('select#job-freq').exists()).toBe(true)
  })

  it('should show LLM query input when job type is LLM', async () => {
    const wrapper = mount(CreateJobModal)
    
    await wrapper.find('select#job-type').setValue('LLM')
    
    expect(wrapper.find('input#query').exists()).toBe(true)
  })

  it('should show tool options when job type is Tool', async () => {
    const wrapper = mount(CreateJobModal)
    
    await wrapper.find('select#job-type').setValue('Tool')
    
    expect(wrapper.find('select#tool-type').exists()).toBe(true)
  })

  it('should emit close when cancel button clicked', async () => {
    const wrapper = mount(CreateJobModal)
    
    await wrapper.find('.cancel-btn').trigger('click')
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit create-job when form submitted', async () => {
    const wrapper = mount(CreateJobModal)
    
    await wrapper.find('input#job-name').setValue('Test Job')
    await wrapper.find('select#job-type').setValue('LLM')
    await wrapper.find('input#query').setValue('Test query')
    await wrapper.find('form').trigger('submit')
    
    expect(wrapper.emitted('create-job')).toBeTruthy()
  })

  it('should close on escape key', async () => {
    const wrapper = mount(CreateJobModal)
    
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    await wrapper.vm.$nextTick()
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })
})
