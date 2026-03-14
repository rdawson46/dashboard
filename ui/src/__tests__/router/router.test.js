import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

const routerModule = await import('@/router/index.js')
const router = routerModule.router

describe('router/index.js', () => {
  let pinia

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)
    vi.clearAllMocks()
  })

  describe('route definitions', () => {
    it('should have a root route', () => {
      const routes = router.getRoutes()
      const rootRoute = routes.find(r => r.path === '/')
      expect(rootRoute).toBeDefined()
    })

    it('should have a login route', () => {
      const routes = router.getRoutes()
      const loginRoute = routes.find(r => r.path === '/login')
      expect(loginRoute).toBeDefined()
      expect(loginRoute.name).toBe('Login')
    })

    it('should have a register route', () => {
      const routes = router.getRoutes()
      const registerRoute = routes.find(r => r.path === '/register')
      expect(registerRoute).toBeDefined()
    })

    it('should have a new chat route', () => {
      const routes = router.getRoutes()
      const chatRoute = routes.find(r => r.path === '/chat')
      expect(chatRoute).toBeDefined()
      expect(chatRoute.name).toBe('New Chat')
      expect(chatRoute.meta.requiresAuth).toBe(true)
    })

    it('should have an existing chat route with param', () => {
      const routes = router.getRoutes()
      const chatWithIdRoute = routes.find(r => r.path === '/chat/:id')
      expect(chatWithIdRoute).toBeDefined()
      expect(chatWithIdRoute.name).toBe('Existing Chat')
      expect(chatWithIdRoute.meta.requiresAuth).toBe(true)
    })

    it('should have a jobs route', () => {
      const routes = router.getRoutes()
      const jobsRoute = routes.find(r => r.path === '/jobs')
      expect(jobsRoute).toBeDefined()
      expect(jobsRoute.name).toBe('jobs')
      expect(jobsRoute.meta.requiresAuth).toBe(true)
    })

    it('should have a models route', () => {
      const routes = router.getRoutes()
      const modelsRoute = routes.find(r => r.path === '/models')
      expect(modelsRoute).toBeDefined()
      expect(modelsRoute.name).toBe('Models')
      expect(modelsRoute.meta.requiresAuth).toBe(true)
    })

    it('should have a files route', () => {
      const routes = router.getRoutes()
      const filesRoute = routes.find(r => r.path === '/files')
      expect(filesRoute).toBeDefined()
      expect(filesRoute.name).toBe('Files')
      expect(filesRoute.meta.requiresAuth).toBe(true)
    })

    it('should have a catch-all route for 404', () => {
      const routes = router.getRoutes()
      const notFoundRoute = routes.find(r => r.path.match(/^.*pathMatch/))
      expect(notFoundRoute).toBeDefined()
    })

    it('should have routes array', () => {
      const routes = router.getRoutes()
      expect(routes.length).toBeGreaterThan(0)
    })

    it('should have 9 total routes', () => {
      const routes = router.getRoutes()
      expect(routes.length).toBe(9)
    })
  })

  describe('route metadata', () => {
    it('should mark chat routes as requiring auth', () => {
      const routes = router.getRoutes()
      const newChatRoute = routes.find(r => r.name === 'New Chat')
      const existingChatRoute = routes.find(r => r.name === 'Existing Chat')
      
      expect(newChatRoute.meta.requiresAuth).toBe(true)
      expect(existingChatRoute.meta.requiresAuth).toBe(true)
    })

    it('should mark jobs route as requiring auth', () => {
      const routes = router.getRoutes()
      const jobsRoute = routes.find(r => r.name === 'jobs')
      
      expect(jobsRoute.meta.requiresAuth).toBe(true)
    })

    it('should mark models route as requiring auth', () => {
      const routes = router.getRoutes()
      const modelsRoute = routes.find(r => r.name === 'Models')
      
      expect(modelsRoute.meta.requiresAuth).toBe(true)
    })

    it('should mark files route as requiring auth', () => {
      const routes = router.getRoutes()
      const filesRoute = routes.find(r => r.name === 'Files')
      
      expect(filesRoute.meta.requiresAuth).toBe(true)
    })
  })
})
