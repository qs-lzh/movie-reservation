import apiClient from './index'

export const authAPI = {
  register: (username, password, role = 'user') => {
    return apiClient.post('/users/register', {
      username,
      password,
      user_role: role
    })
  },

  login: (username, password, key) => {
    return apiClient.post('/users/login', {
      username,
      password,
      key,
    })
  },

  logout: () => {
    return apiClient.post('/users/logout')
  }
}
