import axios from 'axios'

const client = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
})

// 请求拦截器：每次请求自动带上 token
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一处理后端返回的 { code, message, data } 格式
client.interceptors.response.use(
  (response) => {
    const { code, message, data } = response.data

    if (code !== 0) {
      // code 为 401 时清除 token，跳转登录页
      if (code === 401) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
      return Promise.reject(new Error(message || '请求失败'))
    }

    return data
  },
  (error) => {
    return Promise.reject(new Error(error.message || '网络错误'))
  },
)

export default client
