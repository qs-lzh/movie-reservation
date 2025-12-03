import apiClient from './index'

export const userAPI = {
  getUserInfo: () => {
    return apiClient.get('/users/me') // 这个端点可能需要在后端实现
  }
}