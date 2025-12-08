import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api', // This will be proxied to http://localhost:8080/api
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 允许携带cookies
})

// 请求拦截器不再需要设置Authorization header since we're using cookies
apiClient.interceptors.request.use(
  (config) => {
    // We don't need to manually set headers since withCredentials: true will automatically include cookies
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 添加响应拦截器，处理JWT过期等
apiClient.interceptors.response.use(
  (response) => {
    // Check if the response contains a new JWT token in cookies
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // JWT过期或无效，跳转到登录页
      // Note: We don't need to manually remove token since it's in cookies
      // The server will handle the cookie invalidation
    }
    return Promise.reject(error)
  }
)

export default apiClient
