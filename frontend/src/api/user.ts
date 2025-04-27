import request from '@/utils/request'
import type { Response } from './types'
import type { 
  User, 
  UserProfile, 
  UserGameSettings, 
  UserSettings, 
  UserPrivacySettings, 
  SensitivityConfig,
  UserDevicePreference
} from '@/types/user'

// 认证相关
export interface LoginRequest {
  email: string;
  password: string;
  deviceId: string;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  user: User;
  requireTwoFactor?: boolean;
  twoFactorToken?: string;
}

export const login = (data: LoginRequest) => {
  return request.post<Response<LoginResponse>>('/auth/login', data)
}

export interface VerifyTwoFactorRequest {
  twoFactorToken: string;
  twoFactorCode: string;
  deviceId: string;
}

export const verifyTwoFactor = (data: VerifyTwoFactorRequest) => {
  return request.post<Response<LoginResponse>>('/auth/verify-2fa', data)
}

export const register = (data: { username: string; password: string; email: string }) => {
  return request.post<Response<User>>('/auth/signup', data)
}

export const logout = (deviceId: string) => {
  return request.post<Response<null>>('/auth/logout', { deviceId })
}

export const refreshToken = (refreshToken: string) => {
  return request.post<Response<{ token: string; refreshToken: string }>>('/auth/refresh', { refreshToken })
}

// 用户个人资料
export const getUserProfile = () => {
  return request.get<Response<UserProfile>>('/users/me')
}

export interface UpdateProfileRequest {
  avatar?: string
  bio?: string
  gender?: string
  website?: string
  customFields?: Record<string, any>
}

export const updateUserProfile = (data: UpdateProfileRequest) => {
  return request.put<Response<UserProfile>>('/users/me', data)
}

// 游戏设置
export const getGameSettings = () => {
  return request.get<Response<UserGameSettings>>('/users/game-settings')
}

export interface UpdateGameSettingsRequest {
  preferredGames?: string[]
  defaultDPI?: number
  preferredGripStyle?: 'palm' | 'claw' | 'fingertip'
  mouseAcceleration?: boolean
  pollRate?: number
}

export const updateGameSettings = (data: UpdateGameSettingsRequest) => {
  return request.put<Response<UserGameSettings>>('/users/game-settings', data)
}

// 灵敏度配置
export interface AddSensitivityConfigRequest {
  game: string
  sensitivity: number
  dpi?: number
  isActive: boolean
}

export const addSensitivityConfig = (data: AddSensitivityConfigRequest) => {
  return request.post<Response<SensitivityConfig>>('/users/game-settings/sensitivity', data)
}

export interface UpdateSensitivityConfigRequest {
  sensitivity?: number
  dpi?: number
  isActive?: boolean
}

export const updateSensitivityConfig = (game: string, data: UpdateSensitivityConfigRequest) => {
  return request.put<Response<SensitivityConfig>>(`/users/game-settings/sensitivity/${game}`, data)
}

export const deleteSensitivityConfig = (game: string) => {
  return request.delete<Response<null>>(`/users/game-settings/sensitivity/${game}`)
}

// 用户设置
export const getUserSettings = () => {
  return request.get<Response<UserSettings>>('/users/settings')
}

export interface UpdateSettingsRequest {
  language?: string
  theme?: string
  measurementUnit?: string
  notifications?: {
    emailNotifications?: boolean
    pushNotifications?: boolean
    newReviews?: boolean
    replies?: boolean
    systemUpdates?: boolean
  }
}

export const updateUserSettings = (data: UpdateSettingsRequest) => {
  return request.put<Response<UserSettings>>('/users/settings', data)
}

// 隐私设置
export const getPrivacySettings = () => {
  return request.get<Response<UserPrivacySettings>>('/users/privacy')
}

export interface UpdatePrivacySettingsRequest {
  profileVisibility?: 'public' | 'friends' | 'private'
  deviceListVisibility?: 'public' | 'friends' | 'private'
  reviewHistoryVisibility?: 'public' | 'friends' | 'private'
  showOnlineStatus?: boolean
  showActivity?: boolean
}

export const updatePrivacySettings = (data: UpdatePrivacySettingsRequest) => {
  return request.put<Response<UserPrivacySettings>>('/users/privacy', data)
}

// 设备偏好
export const getDevicePreferences = () => {
  return request.get<Response<UserDevicePreference[]>>('/users/device-preferences')
}

export interface AddDevicePreferenceRequest {
  deviceId: string
  deviceType: string
  isFavorite: boolean
  isWishlist: boolean
  rating?: number
  notes?: string
}

export const addDevicePreference = (data: AddDevicePreferenceRequest) => {
  return request.post<Response<UserDevicePreference>>('/users/device-preferences', data)
}

export interface UpdateDevicePreferenceRequest {
  isFavorite?: boolean
  isWishlist?: boolean
  rating?: number
  notes?: string
}

export const updateDevicePreference = (id: string, data: UpdateDevicePreferenceRequest) => {
  return request.put<Response<UserDevicePreference>>(`/users/device-preferences/${id}`, data)
}

export const deleteDevicePreference = (id: string) => {
  return request.delete<Response<null>>(`/users/device-preferences/${id}`)
}