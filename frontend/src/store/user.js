import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getProfile } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    userInfo.value = info
  }

  const login = async (username, password) => {
    const res = await loginApi(username, password)
    setToken(res.data.token)
    setUserInfo(res.data.user)
    return res
  }

  const fetchUserInfo = async () => {
    if (!token.value) return null
    try {
      const res = await getProfile()
      setUserInfo(res.data)
      return res.data
    } catch (e) {
      logout()
      return null
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    isAdmin,
    setToken,
    setUserInfo,
    login,
    fetchUserInfo,
    logout
  }
})
