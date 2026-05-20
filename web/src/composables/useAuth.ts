import { useUserStore } from '@/stores/user'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

export function useAuth() {
  const userStore = useUserStore()
  const router = useRouter()
  const isLoggedIn = computed(() => userStore.isLoggedIn())
  const isAdmin = computed(() => userStore.isAdmin())
  const isSuperAdmin = computed(() => userStore.isSuperAdmin())

  const logout = () => {
    userStore.logout()
    router.push('/')
  }

  return { isLoggedIn, isAdmin, isSuperAdmin, logout, userStore }
}
