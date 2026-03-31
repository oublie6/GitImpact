import axios from 'axios'
const http = axios.create({ baseURL: 'http://127.0.0.1:8080' })
http.interceptors.request.use((cfg) => { const t = localStorage.getItem('token'); if (t) cfg.headers.Authorization = `Bearer ${t}`; return cfg })
export default http
