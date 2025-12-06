import apiClient from './index'

export const authAPI = {
  register: (username, password, role = 'user') => {
    return apiClient.post('/users/register', {
      username,
      password,
      user_role: role
    })
  },

  login: (username, password) => {
    console.log("dafssjdsjsaaaaaaaaaaaaa")
    return apiClient.post('/users/login', {
      username,
      password
    })
  },

  logout: () => {
    return apiClient.post('/users/logout')
  }
}
