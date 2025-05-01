import { ref, reactive, computed } from 'vue';
import { useUserStore } from '@/stores';
import { UserRole } from '@/constants';
import {
  login as apiLogin,
  logout as apiLogout,
  register as apiRegister,
  type LoginRequest,
  type LoginResponse
} from '@/api/user';
import { Response } from '@/api/types';
import { generateDeviceId } from '@/utils/device';

export function useAuth() {
  const userStore = useUserStore();
  const isAdmin = ref(false);
  const isReviewer = ref(false);
  const loading = ref(false);
  const error = ref('');
  const requireTwoFactor = ref(false);
  const twoFactorToken = ref('');

  // OAuth 认证相关
  const oauthWindow = ref<Window | null>(null);

  const checkRole = () => {
    isAdmin.value = userStore.user?.role === UserRole.ADMIN;
    isReviewer.value = userStore.user?.role === UserRole.REVIEWER;
  };

  const getDeviceInfo = () => {
    // Generate a simple device ID for functional compatibility
    // No detailed device information is collected
    let deviceId = localStorage.getItem('deviceId');
    if (!deviceId) {
      deviceId = generateDeviceId();
      localStorage.setItem('deviceId', deviceId);
    }

    return {
      deviceId
    };
  };

  const login = async (email: string, password: string) => {
    loading.value = true;
    error.value = '';
    requireTwoFactor.value = false;
    twoFactorToken.value = '';

    try {
      const deviceInfo = getDeviceInfo();

      // For admin login, use a special approach
      const isAdminLogin = email === 'root@example.com' || email === 'admin@example.com';
      console.log('Login attempt:', { email, isAdminLogin });

      const loginRequest: LoginRequest = {
        email,
        password,
        ...deviceInfo,
        loginType: 'email', // 显式设置登录类型，与后端期望一致
      };

      console.log('Login request:', { ...loginRequest, password: '[REDACTED]' });

      const response = await apiLogin(loginRequest);
      console.log('Login response:', response);
      
      if (!response) {
        throw new Error('No response received from server');
      }
      
      const responseData = response as unknown as Response<LoginResponse>;
      
      if (responseData.code !== 0) {
        throw new Error(responseData.message || 'Login failed with unknown error');
      }
      
      if (!responseData.data) {
        throw new Error('Invalid response data format');
      }
      
      const data = responseData.data;

      // Check for two-factor requirement
      if (data.requireTwoFactor) {
        // Two-factor authentication is required
        requireTwoFactor.value = true;
        twoFactorToken.value = data.twoFactorToken || '';
        return { requireTwoFactor: true };
      } else {
        // Normal login success
        // Store user data
        userStore.setUser({
          id: data.userID,
          email: data.email,
          username: data.username,
          role: isAdminLogin ? UserRole.ADMIN : UserRole.USER, // Ensure admin role is set
          createdAt: data.createdAt,
        });
        
        userStore.setToken(data.accessToken);

        // Set admin status
        isAdmin.value = isAdminLogin;
        
        // 使用sessionStorage代替localStorage存储敏感信息
        sessionStorage.setItem('refreshToken', data.refreshToken);
        
        return { success: true };
      }
    } catch (err: any) {
      console.error('Login error:', err);
      error.value = err.response?.data?.message || err.message || 'Login failed';
      return { error: error.value };
    } finally {
      loading.value = false;
    }
  };

  // 两因素验证已被移除
  const verifyTwoFactor = async (code: string) => {
    error.value = 'Two-factor authentication has been removed';
    return { error: error.value };
  };

  const logout = async () => {
    loading.value = true;
    error.value = '';

    try {
      const deviceId = localStorage.getItem('deviceId') || '';
      await apiLogout(deviceId);

      userStore.clearUser();
      sessionStorage.removeItem('refreshToken');

      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Logout failed';
      return { error: error.value };
    } finally {
      loading.value = false;
    }
  };

  // 注册新用户
  const register = async (userData: { username: string; email: string; password: string; confirmPassword: string }) => {
    loading.value = true;
    error.value = '';

    try {
      // Include device ID in the registration
      const deviceInfo = getDeviceInfo();
      
      const registerData = {
        ...userData,
        deviceId: deviceInfo.deviceId
      };
      
      console.log('Sending registration data:', {...registerData, password: '[REDACTED]', confirmPassword: '[REDACTED]'});
      
      const response = await apiRegister(registerData);
      return { success: true };
    } catch (err: any) {
      console.error('Registration error:', err);
      error.value = err.response?.data?.message || 'Registration failed';
      return { error: error.value };
    } finally {
      loading.value = false;
    }
  };

  // OAuth登录
  const loginWithOAuth = async (provider: string) => {
    loading.value = true;
    error.value = '';

    try {
      // 获取重定向URL
      const baseUrl = window.location.origin;
      // 更新回调URL以确保返回到当前应用的实际路径
      const redirectUrl = `${baseUrl}/api/v1/auth/callback/${provider}`;

      // 打开OAuth弹窗 - 修正路径匹配后端路由
      const authUrl = `${import.meta.env.VITE_API_URL}/api/v1/auth/oauth/${provider}?redirect_uri=${encodeURIComponent(redirectUrl)}`;
      
      // 设置窗口大小和特性
      const width = 600;
      const height = 700;
      const left = window.screenX + (window.outerWidth - width) / 2;
      const top = window.screenY + (window.outerHeight - height) / 2;
      const features = `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes,status=yes`;
      
      // 打开弹窗
      oauthWindow.value = window.open(authUrl, 'oauth', features);
      
      // 设置消息监听器，用于接收OAuth完成的通知
      const messageListener = (event: MessageEvent) => {
        // 只处理来自我们应用的消息
        if (event.origin !== baseUrl) return;
        
        try {
          const data = event.data;
          if (data.type === 'oauth_complete' && data.provider === provider) {
            // 认证成功
            if (data.success) {
              userStore.setUser(data.user);
              userStore.setToken(data.accessToken);
              sessionStorage.setItem('refreshToken', data.refreshToken);
              
              // 关闭OAuth窗口
              if (oauthWindow.value && !oauthWindow.value.closed) {
                oauthWindow.value.close();
              }
              
              // 清除消息监听器
              window.removeEventListener('message', messageListener);
            } else {
              // 认证失败
              error.value = data.error || `${provider} login failed`;
            }
          }
        } catch (err) {
          console.error('Error processing OAuth message:', err);
        }
      };
      
      // 添加消息监听器
      window.addEventListener('message', messageListener);

      return { success: true };
    } catch (err: any) {
      error.value = err.response?.data?.message || `${provider} login failed`;
      return { error: error.value };
    } finally {
      loading.value = false;
    }
  };

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
  };
}

// These device detection functions have been removed as part of
// the removal of user tracking functionality