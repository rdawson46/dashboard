import { describe, it, expect, vi } from 'vitest'
import { useNotify } from '@/composables/notify.js'

vi.mock('vue3-toastify', () => ({
  toast: vi.fn()
}))

describe('useNotify', () => {
  it('should call toast with message and options', async () => {
    const { toast } = await import('vue3-toastify')
    
    useNotify('Test message')
    
    expect(toast).toHaveBeenCalledWith('Test message', {
      autoClose: 1000,
      theme: 'dark'
    })
  })

  it('should call toast with different messages', async () => {
    const { toast } = await import('vue3-toastify')
    
    useNotify('Error occurred')
    expect(toast).toHaveBeenCalledWith('Error occurred', expect.any(Object))
    
    useNotify('Success!')
    expect(toast).toHaveBeenCalledWith('Success!', expect.any(Object))
  })
})
