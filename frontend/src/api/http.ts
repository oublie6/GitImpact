// Axios 实例统一追加 JWT，供全部页面复用。
import axios from 'axios'

const apiBaseURL = import.meta.env.VITE_API_BASE_URL || '/api'

const http = axios.create({ baseURL: apiBaseURL })
http.interceptors.request.use((cfg) => {
  const t = localStorage.getItem('token')
  if (t) cfg.headers.Authorization = `Bearer ${t}`
  return cfg
})

export default http
