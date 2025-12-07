import apiClient from './index'

export const authAPI = {
  register: (username, password, role = 'user', key) => {
    return apiClient.post('/users/register', {
      username,
      password,
      user_role: role,
      key: key
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
