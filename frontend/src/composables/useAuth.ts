import { ref, reactive, computed } from 'vue'
import { useUserStore } from '@/stores'
import { UserRole } from '@/constants'
import { 
  login as apiLogin, 
  verifyTwoFactor as apiVerifyTwoFactor, 
  logout as apiLogout,
  register as apiRegister,
  type LoginRequest,
  type VerifyTwoFactorRequest,
  type LoginResponse
} from '@/api/user'
import { Response } from '@/api/types'
import { generateDeviceId } from '@/utils/device'

export function useAuth() {
  const userStore = useUserStore()
  const isAdmin = ref(false)
  const isReviewer = ref(false)
  const loading = ref(false)
  const error = ref('')
  const requireTwoFactor = ref(false)
  const twoFactorToken = ref('')
  
  // OAuth 认证相关
  const oauthWindow = ref<Window | null>(null)
  
  const checkRole = () => {
    isAdmin.value = userStore.user?.role === UserRole.ADMIN
    isReviewer.value = userStore.user?.role === UserRole.REVIEWER
  }
  
  const getDeviceInfo = () => {
    // Generate a simple device ID for functional compatibility
    // No detailed device information is collected
    let deviceId = localStorage.getItem('deviceId')
    if (!deviceId) {
      deviceId = generateDeviceId()
      localStorage.setItem('deviceId', deviceId)
    }
    
    return {
      deviceId
    }
  }
  
  const login = async (email: string, password: string) => {
    loading.value = true
    error.value = ''
    requireTwoFactor.value = false
    twoFactorToken.value = ''
    
    try {
      const deviceInfo = getDeviceInfo()
      
      const loginRequest: LoginRequest = {
        email,
        password,
        ...deviceInfo
      }
      
      const response = await apiLogin(loginRequest)
      const responseData = response as unknown as Response<LoginResponse>
      const data = responseData.data
      
      if (responseData.data.requireTwoFactor) {
        // Two-factor authentication is required
        requireTwoFactor.value = true
        twoFactorToken.value = responseData.data.twoFactorToken || ''
        return { requireTwoFactor: true }
      } else {
        // Normal login success
        userStore.setUser(responseData.data.user)
        userStore.setToken(responseData.data.accessToken)
        
        // 使用sessionStorage代替localStorage存储敏感信息
        sessionStorage.setItem('refreshToken', responseData.data.refreshToken)
        return { success: true }
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Login failed'
      return { error: error.value }
    } finally {
      loading.value = false
    }
  }
  
  const verifyTwoFactor = async (code: string) => {
    if (!twoFactorToken.value) {
      error.value = 'Invalid session. Please log in again.'
      return { error: error.value }
    }
    
    loading.value = true
    error.value = ''
    
    try {
      const deviceInfo = getDeviceInfo()
      
      const verifyRequest: VerifyTwoFactorRequest = {
        twoFactorToken: twoFactorToken.value,
        twoFactorCode: code,
        deviceId: deviceInfo.deviceId
      }
      
      const response = await apiVerifyTwoFactor(verifyRequest)
      const responseData = response as unknown as Response<LoginResponse>
      const data = responseData.data
      
      userStore.setUser(responseData.data.user)
      userStore.setToken(responseData.data.accessToken)
      sessionStorage.setItem('refreshToken', responseData.data.refreshToken)
      
      // Reset 2FA state
      requireTwoFactor.value = false
      twoFactorToken.value = ''
      
      return { success: true }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Verification failed'
      return { error: error.value }
    } finally {
      loading.value = false
    }
  }
  
  const logout = async () => {
    loading.value = true
    error.value = ''
    
    try {
      const deviceId = localStorage.getItem('deviceId') || ''
      await apiLogout(deviceId)
      
      userStore.clearUser()
      sessionStorage.removeItem('refreshToken')
      
      return { success: true }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Logout failed'
      return { error: error.value }
    } finally {
      loading.value = false
    }
  }
  
  // 注册新用户
  const register = async (userData: { username: string; email: string; password: string }) => {
    loading.value = true
    error.value = ''
    
    try {
      const response = await apiRegister(userData)
      return { success: true }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Registration failed'
      return { error: error.value }
    } finally {
      loading.value = false
    }
  }

  // OAuth登录
  const loginWithOAuth = async (provider: string) => {
    loading.value = true
    error.value = ''
    
    try {
      // 获取重定向URL
      const baseUrl = window.location.origin
      const redirectUrl = `${baseUrl}/auth/callback/${provider}`
      
      // 打开OAuth弹窗 - 修正路径匹配后端路由
      const authUrl = `${import.meta.env.VITE_API_URL}/api/v1/auth/oauth/${provider}/login?redirect_uri=${encodeURIComponent(redirectUrl)}`
      oauthWindow.value = window.open(authUrl, 'oauth', 'width=600,height=700')
      
      // 实际应用中应该监听OAuth回调
      // 本例中，我们简单地在新窗口中完成OAuth流程
      
      return { success: true }
    } catch (err: any) {
      error.value = err.response?.data?.message || `${provider} login failed`
      return { error: error.value }
    } finally {
      loading.value = false
    }
  }

  return {
    isAdmin,
    isReviewer,
    loading,
    error,
    requireTwoFactor,
    isAuthenticated: computed(() => !!userStore.token),
    user: computed(() => userStore.user),
    checkRole,
    login,
    register,
    loginWithOAuth,
    verifyTwoFactor,
    logout
  }
}

// These device detection functions have been removed as part of 
// the removal of user tracking functionality