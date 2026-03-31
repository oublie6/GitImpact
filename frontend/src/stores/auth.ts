import { defineStore } from 'pinia'
import http from '../api/http'

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: localStorage.getItem('token') || '', user: null as any }),
  actions: {
    async login(username: string, password: string) {
      const res = await http.post('/api/auth/login', { username, password })
      this.token = res.data.data.token
      this.user = res.data.data.user
      localStorage.setItem('token', this.token)
    },
    async register(payload: { username: string; password: string; email: string }) { await http.post('/api/auth/register', payload) },
    logout() { this.token = ''; this.user = null; localStorage.removeItem('token') }
  }
})
