// 认证状态仓库存放 JWT 与当前登录用户信息。
import { defineStore } from 'pinia'
import http from '../api/http'

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: localStorage.getItem('token') || '', user: null as any }),
  actions: {
    // login 调用后端登录接口，并把 JWT 持久化到 localStorage。
    async login(username: string, password: string) {
      const res = await http.post('/api/auth/login', { username, password })
      this.token = res.data.data.token
      this.user = res.data.data.user
      localStorage.setItem('token', this.token)
    },

    // register 只负责创建数据库用户，不自动登录。
    async register(payload: { username: string; password: string; email: string }) { await http.post('/api/auth/register', payload) },

    // logout 当前仅清理前端状态，服务端没有令牌吊销机制。
    logout() { this.token = ''; this.user = null; localStorage.removeItem('token') }
  }
})
