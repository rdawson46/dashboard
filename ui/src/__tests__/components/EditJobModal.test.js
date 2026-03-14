import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EditJobModal from '@/components/EditJobModal.vue'

const mockJob = {
  id: 1,
  name: 'Test Job',
  task: { task_type: 'LLM', query: 'Test query' },
  freq: '15mins'
}

describe('EditJobModal.vue', () => {
  it('should render modal', () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.find('.modal-backdrop').exists()).toBe(true)
  })

  it('should pre-fill form with job data', async () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    await wrapper.vm.$nextTick()
    expect(wrapper.find('input#job-name').element.value).toBe('Test Job')
  })

  it('should emit close when cancel button clicked', async () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    
    await wrapper.find('.cancel-btn').trigger('click')
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should emit edit-job when form submitted', async () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    
    await wrapper.find('form').trigger('submit')
    
    expect(wrapper.emitted('edit-job')).toBeTruthy()
  })

  it('should have job name input', () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.find('input#job-name').exists()).toBe(true)
  })

  it('should have job type select', () => {
    const wrapper = mount(EditJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.find('select#job-type').exists()).toBe(true)
  })
})
