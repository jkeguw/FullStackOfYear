import { ref } from 'vue'
import { useUserStore } from '@/stores'
import { UserRole } from '@/constants'

export function useAuth() {
  const userStore = useUserStore()
  const isAdmin = ref(false)
  const isReviewer = ref(false)

  const checkRole = () => {
    isAdmin.value = userStore.user?.role === UserRole.ADMIN
    isReviewer.value = userStore.user?.role === UserRole.REVIEWER
  }

  return {
    isAdmin,
    isReviewer,
    checkRole
  }
}