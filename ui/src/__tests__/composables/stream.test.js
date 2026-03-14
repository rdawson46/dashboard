import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('@/stores/chat.js', () => ({
  useChatStore: vi.fn(() => ({
    addChatToHistory: vi.fn()
  }))
}))

vi.mock('@/composables/notify.js', () => ({
  useNotify: vi.fn(() => vi.fn())
}))

describe('stream.js', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.resetModules()
  })

  it('should export useStream function', async () => {
    const streamModule = await import('@/composables/stream.js')
    expect(streamModule.useStream).toBeDefined()
  })

  it('useStream should be an async function', async () => {
    const streamModule = await import('@/composables/stream.js')
    const useStream = streamModule.useStream
    
    const isAsync = useStream.constructor.name === 'AsyncFunction'
    expect(isAsync).toBe(true)
  })
})
