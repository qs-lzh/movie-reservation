import apiClient from './index'

export const authAPI = {
  register: (username, password, role = 'user', key, adminRolePassword) => {
    return apiClient.post('/users/register', {
      username,
      password,
      user_role: role,
      key: key,
      admin_role_password: adminRolePassword,
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
