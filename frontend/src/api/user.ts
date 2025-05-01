import request from '@/utils/request';
import type { Response } from './types';
import type {
  User,
  UserProfile,
  UserSettings
} from '@/types/user';

// 认证相关
export interface LoginRequest {
  email: string;
  password: string;
  deviceId: string;
  loginType?: string; // 添加可选的loginType字段与后端匹配
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
  requireTwoFactor?: boolean;
  twoFactorToken?: string;
}

export const login = (data: LoginRequest) => {
  return request.post<Response<LoginResponse>>('/api/auth/login', data);
};

// 注意: 两因素认证功能已被移除
// export interface VerifyTwoFactorRequest {
//   twoFactorToken: string;
//   twoFactorCode: string;
//   deviceId: string;
// }
// 
// export const verifyTwoFactor = (data: VerifyTwoFactorRequest) => {
//   return request.post<Response<LoginResponse>>('/api/auth/verify-2fa', data);
// };

export interface RegisterRequest {
  username: string; 
  password: string; 
  email: string;
  confirmPassword: string;
  deviceId?: string;
}

export const register = (data: RegisterRequest) => {
  // Ensure deviceId is included
  if (!data.deviceId) {
    const deviceId = localStorage.getItem('deviceId') || `device_${new Date().getTime()}`;
    localStorage.setItem('deviceId', deviceId);
    data.deviceId = deviceId;
  }
  return request.post<Response<User>>('/api/auth/register', data);
};

export const logout = (deviceId: string) => {
  return request.post<Response<null>>('/api/auth/logout', { deviceId });
};

export const refreshToken = (refreshToken: string) => {
  return request.post<Response<{ token: string; refreshToken: string }>>('/api/auth/refresh', {
    refreshToken
  });
};

// 用户个人资料
export const getUserProfile = () => {
  return request.get<Response<UserProfile>>('/api/users/me');
};

export interface UpdateProfileRequest {
  avatar?: string;
  bio?: string;
  gender?: string;
  website?: string;
  customFields?: Record<string, any>;
}

export const updateUserProfile = (data: UpdateProfileRequest) => {
  return request.put<Response<UserProfile>>('/api/users/me', data);
};

// 注意：游戏设置、用户设置、隐私设置和设备偏好功能已被移除

