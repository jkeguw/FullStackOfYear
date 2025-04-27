export interface User {
  id: string
  username: string
  email: string
  role: string
  status?: UserStatus
  profile?: UserProfile
  gameSettings?: UserGameSettings
  settings?: UserSettings
  privacy?: UserPrivacySettings
  createdAt?: string
  updatedAt?: string
}

export interface UserStatus {
  emailVerified: boolean
}

export interface UserProfile {
  avatar?: string
  bio?: string
  gender?: string
  website?: string
  customFields?: Record<string, any>
  completedRate: number
}

export interface UserGameSettings {
  preferredGames?: string[]
  defaultDPI: number
  preferredGripStyle?: 'palm' | 'claw' | 'fingertip'
  mouseAcceleration: boolean
  pollRate: number
  sensitivityConfigs?: SensitivityConfig[]
}

export interface SensitivityConfig {
  game: string
  sensitivity: number
  dpi: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface UserSettings {
  language: string
  theme: string
  measurementUnit: string
  notificationSettings?: NotificationSettings
}

export interface NotificationSettings {
  emailNotifications: boolean
  pushNotifications: boolean
  newReviews: boolean
  replies: boolean
  systemUpdates: boolean
}

export interface UserPrivacySettings {
  profileVisibility: 'public' | 'friends' | 'private'
  deviceListVisibility: 'public' | 'friends' | 'private'
  reviewHistoryVisibility: 'public' | 'friends' | 'private'
  showOnlineStatus: boolean
  showActivity: boolean
}

export interface UserDevicePreference {
  id: string
  deviceId: string
  deviceType: string
  deviceName?: string
  deviceBrand?: string
  isFavorite: boolean
  isWishlist: boolean
  rating?: number
  notes?: string
  createdAt: string
  updatedAt: string
}

export interface UserState {
  user: User | null
  token: string | null
  profile: UserProfile | null
  gameSettings: UserGameSettings | null
  settings: UserSettings | null
  privacy: UserPrivacySettings | null
  devicePreferences: UserDevicePreference[]
}