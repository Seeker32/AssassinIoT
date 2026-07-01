import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi, getUserInfo } from '@/api/auth'
import type { LoginRequest, UserInfo } from '@/api/types/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<UserInfo | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  async function login(payload: LoginRequest) {
    const res = await loginApi(payload)
    token.value = res.token
    user.value = res.user
    const redirect = (router.currentRoute.value.query.redirect as string) || '/dashboard'
    router.push(redirect)
  }

  async function fetchUserInfo() {
    const res = await getUserInfo()
    user.value = res
  }

  function logout() {
    token.value = null
    user.value = null
    router.push('/login')
  }

  return { token, user, isAuthenticated, login, fetchUserInfo, logout }
}, {
  persist: {
    key: 'assassin_auth',
    storage: localStorage,
    pick: ['token'],
  },
})
