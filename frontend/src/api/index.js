import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api', // This will be proxied to http://localhost:8080/api
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 允许携带cookies
})

// 添加请求拦截器，自动携带JWT
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 添加响应拦截器，处理JWT过期等
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // JWT过期或无效，清除本地token并跳转到登录页
      localStorage.removeItem('token')
      // 可以在这里触发登出事件或跳转到登录页
    }
    return Promise.reject(error)
  }
)

export default apiClient
