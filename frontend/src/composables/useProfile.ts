import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getUserProfile,
  updateUserProfile,
  getGameSettings,
  updateGameSettings,
  addSensitivityConfig,
  updateSensitivityConfig,
  deleteSensitivityConfig,
  getUserSettings,
  updateUserSettings,
  getPrivacySettings,
  updatePrivacySettings,
  getDevicePreferences,
  addDevicePreference,
  updateDevicePreference,
  deleteDevicePreference,
  type UpdateProfileRequest,
  type UpdateGameSettingsRequest,
  type AddSensitivityConfigRequest,
  type UpdateSensitivityConfigRequest,
  type UpdateSettingsRequest,
  type UpdatePrivacySettingsRequest,
  type AddDevicePreferenceRequest,
  type UpdateDevicePreferenceRequest
} from '@/api/user'
import type {
  UserProfile,
  UserGameSettings,
  UserSettings,
  UserPrivacySettings,
  SensitivityConfig,
  UserDevicePreference
} from '@/types/user'

/**
 * 用户资料管理钩子
 */
export function useProfile() {
  // 状态
  const profile = ref<UserProfile | null>(null)
  const gameSettings = ref<UserGameSettings | null>(null)
  const settings = ref<UserSettings | null>(null)
  const privacySettings = ref<UserPrivacySettings | null>(null)
  const devicePreferences = ref<UserDevicePreference[]>([])
  
  // 加载状态
  const profileLoading = ref(false)
  const gameSettingsLoading = ref(false)
  const settingsLoading = ref(false)
  const privacyLoading = ref(false)
  const devicePreferencesLoading = ref(false)
  
  // 错误状态
  const error = ref<string | null>(null)
  
  // 加载用户个人资料
  const fetchProfile = async (): Promise<UserProfile | null> => {
    profileLoading.value = true
    error.value = null
    
    try {
      const res = await getUserProfile()
      profile.value = res.data
      return profile.value
    } catch (err) {
      console.error('获取用户资料失败', err)
      error.value = '获取用户资料失败，请刷新重试'
      ElMessage.error(error.value)
      return null
    } finally {
      profileLoading.value = false
    }
  }
  
  // 更新用户个人资料
  const updateProfile = async (data: UpdateProfileRequest): Promise<UserProfile | null> => {
    profileLoading.value = true
    error.value = null
    
    try {
      const res = await updateUserProfile(data)
      profile.value = res.data
      ElMessage.success('个人资料已更新')
      return profile.value
    } catch (err) {
      console.error('更新用户资料失败', err)
      error.value = '更新用户资料失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      profileLoading.value = false
    }
  }
  
  // 加载游戏设置
  const fetchGameSettings = async (): Promise<UserGameSettings | null> => {
    gameSettingsLoading.value = true
    error.value = null
    
    try {
      const res = await getGameSettings()
      gameSettings.value = res.data
      return gameSettings.value
    } catch (err) {
      console.error('获取游戏设置失败', err)
      error.value = '获取游戏设置失败，请刷新重试'
      ElMessage.error(error.value)
      return null
    } finally {
      gameSettingsLoading.value = false
    }
  }
  
  // 更新游戏设置
  const updateGameSettingsData = async (data: UpdateGameSettingsRequest) => {
    gameSettingsLoading.value = true
    error.value = null
    
    try {
      const res = await updateGameSettings(data)
      gameSettings.value = res.data
      ElMessage.success('游戏设置已更新')
      return res.data
    } catch (err) {
      console.error('更新游戏设置失败', err)
      error.value = '更新游戏设置失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      gameSettingsLoading.value = false
    }
  }
  
  // 添加灵敏度配置
  const addSensitivityConfigData = async (data: AddSensitivityConfigRequest) => {
    gameSettingsLoading.value = true
    error.value = null
    
    try {
      const res = await addSensitivityConfig(data)
      // 更新本地游戏设置
      if (gameSettings.value) {
        if (!gameSettings.value.sensitivityConfigs) {
          gameSettings.value.sensitivityConfigs = []
        }
        gameSettings.value.sensitivityConfigs.push(res.data)
      }
      ElMessage.success('灵敏度配置已添加')
      return res.data
    } catch (err) {
      console.error('添加灵敏度配置失败', err)
      error.value = '添加灵敏度配置失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      gameSettingsLoading.value = false
    }
  }
  
  // 更新灵敏度配置
  const updateSensitivityConfigData = async (game: string, data: UpdateSensitivityConfigRequest) => {
    gameSettingsLoading.value = true
    error.value = null
    
    try {
      const res = await updateSensitivityConfig(game, data)
      // 更新本地游戏设置
      if (gameSettings.value && gameSettings.value.sensitivityConfigs) {
        const index = gameSettings.value.sensitivityConfigs.findIndex(config => config.game === game)
        if (index !== -1) {
          gameSettings.value.sensitivityConfigs[index] = res.data
        }
      }
      ElMessage.success('灵敏度配置已更新')
      return res.data
    } catch (err) {
      console.error('更新灵敏度配置失败', err)
      error.value = '更新灵敏度配置失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      gameSettingsLoading.value = false
    }
  }
  
  // 删除灵敏度配置
  const deleteSensitivityConfigData = async (game: string) => {
    gameSettingsLoading.value = true
    error.value = null
    
    try {
      await deleteSensitivityConfig(game)
      // 更新本地游戏设置
      if (gameSettings.value && gameSettings.value.sensitivityConfigs) {
        gameSettings.value.sensitivityConfigs = gameSettings.value.sensitivityConfigs.filter(
          config => config.game !== game
        )
      }
      ElMessage.success('灵敏度配置已删除')
      return true
    } catch (err) {
      console.error('删除灵敏度配置失败', err)
      error.value = '删除灵敏度配置失败，请重试'
      ElMessage.error(error.value)
      return false
    } finally {
      gameSettingsLoading.value = false
    }
  }
  
  // 加载用户设置
  const fetchSettings = async () => {
    settingsLoading.value = true
    error.value = null
    
    try {
      const res = await getUserSettings()
      settings.value = res.data
      return res.data
    } catch (err) {
      console.error('获取用户设置失败', err)
      error.value = '获取用户设置失败，请刷新重试'
      ElMessage.error(error.value)
      return null
    } finally {
      settingsLoading.value = false
    }
  }
  
  // 更新用户设置
  const updateSettingsData = async (data: UpdateSettingsRequest) => {
    settingsLoading.value = true
    error.value = null
    
    try {
      const res = await updateUserSettings(data)
      settings.value = res.data
      ElMessage.success('用户设置已更新')
      return res.data
    } catch (err) {
      console.error('更新用户设置失败', err)
      error.value = '更新用户设置失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      settingsLoading.value = false
    }
  }
  
  // 加载隐私设置
  const fetchPrivacySettings = async () => {
    privacyLoading.value = true
    error.value = null
    
    try {
      const res = await getPrivacySettings()
      privacySettings.value = res.data
      return res.data
    } catch (err) {
      console.error('获取隐私设置失败', err)
      error.value = '获取隐私设置失败，请刷新重试'
      ElMessage.error(error.value)
      return null
    } finally {
      privacyLoading.value = false
    }
  }
  
  // 更新隐私设置
  const updatePrivacySettingsData = async (data: UpdatePrivacySettingsRequest) => {
    privacyLoading.value = true
    error.value = null
    
    try {
      const res = await updatePrivacySettings(data)
      privacySettings.value = res.data
      ElMessage.success('隐私设置已更新')
      return res.data
    } catch (err) {
      console.error('更新隐私设置失败', err)
      error.value = '更新隐私设置失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      privacyLoading.value = false
    }
  }
  
  // 加载设备偏好
  const fetchDevicePreferences = async () => {
    devicePreferencesLoading.value = true
    error.value = null
    
    try {
      const res = await getDevicePreferences()
      devicePreferences.value = res.data
      return res.data
    } catch (err) {
      console.error('获取设备偏好失败', err)
      error.value = '获取设备偏好失败，请刷新重试'
      ElMessage.error(error.value)
      return null
    } finally {
      devicePreferencesLoading.value = false
    }
  }
  
  // 添加设备偏好
  const addDevicePreferenceData = async (data: AddDevicePreferenceRequest) => {
    devicePreferencesLoading.value = true
    error.value = null
    
    try {
      const res = await addDevicePreference(data)
      devicePreferences.value.unshift(res.data)
      ElMessage.success('设备偏好已添加')
      return res.data
    } catch (err) {
      console.error('添加设备偏好失败', err)
      error.value = '添加设备偏好失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      devicePreferencesLoading.value = false
    }
  }
  
  // 更新设备偏好
  const updateDevicePreferenceData = async (id: string, data: UpdateDevicePreferenceRequest) => {
    devicePreferencesLoading.value = true
    error.value = null
    
    try {
      const res = await updateDevicePreference(id, data)
      // 更新本地设备偏好
      const index = devicePreferences.value.findIndex(pref => pref.id === id)
      if (index !== -1) {
        devicePreferences.value[index] = res.data
      }
      ElMessage.success('设备偏好已更新')
      return res.data
    } catch (err) {
      console.error('更新设备偏好失败', err)
      error.value = '更新设备偏好失败，请重试'
      ElMessage.error(error.value)
      return null
    } finally {
      devicePreferencesLoading.value = false
    }
  }
  
  // 删除设备偏好
  const deleteDevicePreferenceData = async (id: string) => {
    devicePreferencesLoading.value = true
    error.value = null
    
    try {
      await deleteDevicePreference(id)
      // 更新本地设备偏好
      devicePreferences.value = devicePreferences.value.filter(pref => pref.id !== id)
      ElMessage.success('设备偏好已删除')
      return true
    } catch (err) {
      console.error('删除设备偏好失败', err)
      error.value = '删除设备偏好失败，请重试'
      ElMessage.error(error.value)
      return false
    } finally {
      devicePreferencesLoading.value = false
    }
  }
  
  // 加载用户所有数据
  const fetchAllUserData = async () => {
    error.value = null
    
    await Promise.all([
      fetchProfile(),
      fetchGameSettings(),
      fetchSettings(),
      fetchPrivacySettings(),
      fetchDevicePreferences()
    ])
  }
  
  // 判断手型大小
  const getHandSizeLabel = (handSize: string): string => {
    switch (handSize) {
      case 'small': return '小型手'
      case 'medium': return '中型手'
      case 'large': return '大型手'
      default: return '未知'
    }
  }
  
  // 判断握持方式
  const getGripStyleLabel = (gripStyle: string): string => {
    switch (gripStyle) {
      case 'palm': return '手掌握持'
      case 'claw': return '爪式握持'
      case 'fingertip': return '指尖握持'
      default: return '未知'
    }
  }
  
  // 格式化单位
  const formatWithUnit = (value: number, unit: string): string => {
    const formattedValue = value.toFixed(1)
    return `${formattedValue} ${unit}`
  }
  
  // 计算DPI转换后的灵敏度
  const convertDPI = (sourceDPI: number, targetDPI: number, sensitivity: number): number => {
    return (sensitivity * sourceDPI) / targetDPI
  }

  return {
    // 状态
    profile,
    gameSettings,
    settings,
    privacySettings,
    devicePreferences,
    profileLoading,
    gameSettingsLoading,
    settingsLoading,
    privacyLoading,
    devicePreferencesLoading,
    error,
    
    // 方法
    fetchProfile,
    updateProfile,
    fetchGameSettings,
    updateGameSettings: updateGameSettingsData,
    addSensitivityConfig: addSensitivityConfigData,
    updateSensitivityConfig: updateSensitivityConfigData,
    deleteSensitivityConfig: deleteSensitivityConfigData,
    fetchSettings,
    updateSettings: updateSettingsData,
    fetchPrivacySettings,
    updatePrivacySettings: updatePrivacySettingsData,
    fetchDevicePreferences,
    addDevicePreference: addDevicePreferenceData,
    updateDevicePreference: updateDevicePreferenceData,
    deleteDevicePreference: deleteDevicePreferenceData,
    fetchAllUserData,
    
    // 辅助方法
    getHandSizeLabel,
    getGripStyleLabel,
    formatWithUnit,
    convertDPI
  }
}