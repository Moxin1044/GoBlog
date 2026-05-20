import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const username = ref(localStorage.getItem('username') || '')
  const role = ref(localStorage.getItem('role') || '')
  const userId = ref(Number(localStorage.getItem('userId')) || 0)

  function setLogin(data: { token: string; username: string; role: string; userId: number }) {
    token.value = data.token
    username.value = data.username
    role.value = data.role
    userId.value = data.userId
    localStorage.setItem('token', data.token)
    localStorage.setItem('username', data.username)
    localStorage.setItem('role', data.role)
    localStorage.setItem('userId', String(data.userId))
  }

  function logout() {
    token.value = ''
    username.value = ''
    role.value = ''
    userId.value = 0
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('role')
    localStorage.removeItem('userId')
  }

  const isLoggedIn = () => !!token.value
  const isAdmin = () => role.value === 'admin' || role.value === 'superadmin'
  const isSuperAdmin = () => role.value === 'superadmin'

  return { token, username, role, userId, setLogin, logout, isLoggedIn, isAdmin, isSuperAdmin }
})
