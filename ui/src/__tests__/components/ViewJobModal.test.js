import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ViewJobModal from '@/components/ViewJobModal.vue'

const mockJob = {
  id: 1,
  name: 'Test Job',
  task: { task_type: 'LLM', query: 'Test query' },
  freq: '15mins',
  status: 'active'
}

describe('ViewJobModal.vue', () => {
  it('should render modal', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.find('.modal-backdrop').exists()).toBe(true)
  })

  it('should display job name', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.find('h2').text()).toBe('Test Job')
  })

  it('should display job type', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.text()).toContain('Job Type:')
  })

  it('should display frequency', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.text()).toContain('Frequency:')
  })

  it('should display status', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.text()).toContain('Status:')
  })

  it('should emit close when close button clicked', async () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    
    await wrapper.find('.primary-btn').trigger('click')
    
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('should display query for LLM jobs', () => {
    const wrapper = mount(ViewJobModal, {
      props: { job: mockJob }
    })
    expect(wrapper.text()).toContain('Query:')
  })
})
