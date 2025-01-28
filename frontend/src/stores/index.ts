import { defineStore } from 'pinia'
import type { User, UserState } from '@/types/user'

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null
  }),
  actions: {
    setUser(user: User) {
      this.user = user
    },
    setToken(token: string) {
      this.token = token
    },
    clearUser() {
      this.user = null
      this.token = null
    }
  }
})