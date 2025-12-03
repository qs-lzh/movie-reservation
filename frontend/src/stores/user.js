import { ref } from 'vue'
import { defineStore } from 'pinia'
import { authAPI } from '@/api/auth'
import apiClient from '@/api/index'

export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const isAuthenticated = ref(false)

  const login = async (username, password) => {
    try {
      const response = await authAPI.login(username, password)
      if (response.data.data) {
        // 成功登录后，从响应中获取用户信息
        const userData = response.data.data
        isAuthenticated.value = true
        user.value = {
          username: userData.username,
          role: userData.role
        }
        // 保存到 localStorage 以便页面刷新后恢复状态
        localStorage.setItem('user', JSON.stringify(user.value))
        return { success: true, message: 'Login successful' }
      } else {
        return { success: false, message: 'Login failed' }
      }
    } catch (error) {
      console.error('Login error:', error)
      let message = 'Login failed'
      if (error.response?.data?.message) {
        message = error.response.data.message
      } else if (error.response?.data) {
        message = error.response.data
      } else if (error.response?.statusText) {
        message = `${error.response.status}: ${error.response.statusText}`
      }
      return { success: false, message }
    }
  }

  const register = async (username, password, role = 'user') => {
    try {
      const response = await authAPI.register(username, password, role)

      // 注册成功
      if (response.data.data) {
        return { success: true, message: response.data.data || 'Registration successful' }
      } else {
        return { success: false, message: 'Registration failed' }
      }
    } catch (error) {
      console.error('Registration error:', error)
      let message = 'Registration failed'
      if (error.response?.data?.message) {
        message = error.response.data.message
      } else if (error.response?.data) {
        message = error.response.data
      } else if (error.response?.statusText) {
        message = `${error.response.status}: ${error.response.statusText}`
      }
      return { success: false, message }
    }
  }

  const logout = async () => {
    try {
      // 调用后端登出API
      await authAPI.logout()
      // 清除本地状态
      user.value = null
      isAuthenticated.value = false
      localStorage.removeItem('user')
      return { success: true, message: 'Logged out successfully' }
    } catch (error) {
      console.error('Logout error:', error)
      // 即使API调用失败，也清除本地状态
      user.value = null
      isAuthenticated.value = false
      localStorage.removeItem('user')
      return { success: false, message: error.response?.data?.message || 'Logout failed' }
    }
  }

  const checkAuth = () => {
    // 检查 localStorage 中是否有用户信息
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
        isAuthenticated.value = true
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        localStorage.removeItem('user')
      }
    }
  }

  return {
    user,
    isAuthenticated,
    login,
    register,
    logout,
    checkAuth
  }
})