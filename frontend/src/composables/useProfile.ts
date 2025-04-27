import { ref, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import {
  getUserProfile,
  updateUserProfile,
  type UpdateProfileRequest
} from '@/api/user';
import type {
  UserProfile,
  UserGameSettings,
  UserSettings,
  UserPrivacySettings,
  SensitivityConfig,
  UserDevicePreference
} from '@/types/user';

/**
 * 用户资料管理钩子
 */
export function useProfile() {
  // 状态
  const profile = ref<UserProfile | null>(null);
  const gameSettings = ref<UserGameSettings | null>(null);
  const settings = ref<UserSettings | null>(null);
  const privacySettings = ref<UserPrivacySettings | null>(null);
  const devicePreferences = ref<UserDevicePreference[]>([]);

  // 加载状态
  const profileLoading = ref(false);
  const gameSettingsLoading = ref(false);
  const settingsLoading = ref(false);
  const privacyLoading = ref(false);
  const devicePreferencesLoading = ref(false);

  // 错误状态
  const error = ref<string | null>(null);

  // 加载用户个人资料
  const fetchProfile = async (): Promise<UserProfile | null> => {
    profileLoading.value = true;
    error.value = null;

    try {
      const res = await getUserProfile();
      profile.value = res.data;
      return profile.value;
    } catch (err) {
      console.error('获取用户资料失败', err);
      error.value = '获取用户资料失败，请刷新重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      profileLoading.value = false;
    }
  };

  // 更新用户个人资料
  const updateProfile = async (data: UpdateProfileRequest): Promise<UserProfile | null> => {
    profileLoading.value = true;
    error.value = null;

    try {
      const res = await updateUserProfile(data);
      profile.value = res.data;
      ElMessage.success('个人资料已更新');
      return profile.value;
    } catch (err) {
      console.error('更新用户资料失败', err);
      error.value = '更新用户资料失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      profileLoading.value = false;
    }
  };

  // 加载游戏设置 - 注意：功能已被移除
  const fetchGameSettings = async (): Promise<UserGameSettings | null> => {
    gameSettingsLoading.value = true;
    error.value = null;

    try {
      console.log('游戏设置功能已被移除');
      gameSettings.value = null;
      return null;
    } catch (err) {
      console.error('获取游戏设置失败', err);
      error.value = '获取游戏设置失败，请刷新重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      gameSettingsLoading.value = false;
    }
  };

  // 更新游戏设置 - 注意：功能已被移除
  const updateGameSettingsData = async (data: any) => {
    gameSettingsLoading.value = true;
    error.value = null;

    try {
      console.log('游戏设置功能已被移除');
      ElMessage.warning('游戏设置功能已被移除');
      return null;
    } catch (err) {
      console.error('更新游戏设置失败', err);
      error.value = '更新游戏设置失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      gameSettingsLoading.value = false;
    }
  };

  // 添加灵敏度配置 - 注意：功能已被移除
  const addSensitivityConfigData = async (data: any) => {
    gameSettingsLoading.value = true;
    error.value = null;

    try {
      console.log('灵敏度配置功能已被移除');
      ElMessage.warning('灵敏度配置功能已被移除');
      return null;
    } catch (err) {
      console.error('添加灵敏度配置失败', err);
      error.value = '添加灵敏度配置失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      gameSettingsLoading.value = false;
    }
  };

  // 更新灵敏度配置 - 注意：功能已被移除
  const updateSensitivityConfigData = async (
    game: string,
    data: any
  ) => {
    gameSettingsLoading.value = true;
    error.value = null;

    try {
      console.log('灵敏度配置功能已被移除');
      ElMessage.warning('灵敏度配置功能已被移除');
      return null;
    } catch (err) {
      console.error('更新灵敏度配置失败', err);
      error.value = '更新灵敏度配置失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      gameSettingsLoading.value = false;
    }
  };

  // 删除灵敏度配置 - 注意：功能已被移除
  const deleteSensitivityConfigData = async (game: string) => {
    gameSettingsLoading.value = true;
    error.value = null;

    try {
      console.log('灵敏度配置功能已被移除');
      ElMessage.warning('灵敏度配置功能已被移除');
      return false;
    } catch (err) {
      console.error('删除灵敏度配置失败', err);
      error.value = '删除灵敏度配置失败，请重试';
      ElMessage.error(error.value);
      return false;
    } finally {
      gameSettingsLoading.value = false;
    }
  };

  // 加载用户设置 - 注意：功能已被移除
  const fetchSettings = async () => {
    settingsLoading.value = true;
    error.value = null;

    try {
      console.log('用户设置功能已被移除');
      settings.value = null;
      return null;
    } catch (err) {
      console.error('获取用户设置失败', err);
      error.value = '获取用户设置失败，请刷新重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      settingsLoading.value = false;
    }
  };

  // 更新用户设置 - 注意：功能已被移除
  const updateSettingsData = async (data: any) => {
    settingsLoading.value = true;
    error.value = null;

    try {
      console.log('用户设置功能已被移除');
      ElMessage.warning('用户设置功能已被移除');
      return null;
    } catch (err) {
      console.error('更新用户设置失败', err);
      error.value = '更新用户设置失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      settingsLoading.value = false;
    }
  };

  // 加载隐私设置 - 注意：功能已被移除
  const fetchPrivacySettings = async () => {
    privacyLoading.value = true;
    error.value = null;

    try {
      console.log('隐私设置功能已被移除');
      privacySettings.value = null;
      return null;
    } catch (err) {
      console.error('获取隐私设置失败', err);
      error.value = '获取隐私设置失败，请刷新重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      privacyLoading.value = false;
    }
  };

  // 更新隐私设置 - 注意：功能已被移除
  const updatePrivacySettingsData = async (data: any) => {
    privacyLoading.value = true;
    error.value = null;

    try {
      console.log('隐私设置功能已被移除');
      ElMessage.warning('隐私设置功能已被移除');
      return null;
    } catch (err) {
      console.error('更新隐私设置失败', err);
      error.value = '更新隐私设置失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      privacyLoading.value = false;
    }
  };

  // 加载设备偏好 - 注意：功能已被移除
  const fetchDevicePreferences = async () => {
    devicePreferencesLoading.value = true;
    error.value = null;

    try {
      console.log('设备偏好功能已被移除');
      devicePreferences.value = [];
      return [];
    } catch (err) {
      console.error('获取设备偏好失败', err);
      error.value = '获取设备偏好失败，请刷新重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      devicePreferencesLoading.value = false;
    }
  };

  // 添加设备偏好 - 注意：功能已被移除
  const addDevicePreferenceData = async (data: any) => {
    devicePreferencesLoading.value = true;
    error.value = null;

    try {
      console.log('设备偏好功能已被移除');
      ElMessage.warning('设备偏好功能已被移除');
      return null;
    } catch (err) {
      console.error('添加设备偏好失败', err);
      error.value = '添加设备偏好失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      devicePreferencesLoading.value = false;
    }
  };

  // 更新设备偏好 - 注意：功能已被移除
  const updateDevicePreferenceData = async (id: string, data: any) => {
    devicePreferencesLoading.value = true;
    error.value = null;

    try {
      console.log('设备偏好功能已被移除');
      ElMessage.warning('设备偏好功能已被移除');
      return null;
    } catch (err) {
      console.error('更新设备偏好失败', err);
      error.value = '更新设备偏好失败，请重试';
      ElMessage.error(error.value);
      return null;
    } finally {
      devicePreferencesLoading.value = false;
    }
  };

  // 删除设备偏好 - 注意：功能已被移除
  const deleteDevicePreferenceData = async (id: string) => {
    devicePreferencesLoading.value = true;
    error.value = null;

    try {
      console.log('设备偏好功能已被移除');
      ElMessage.warning('设备偏好功能已被移除');
      return false;
    } catch (err) {
      console.error('删除设备偏好失败', err);
      error.value = '删除设备偏好失败，请重试';
      ElMessage.error(error.value);
      return false;
    } finally {
      devicePreferencesLoading.value = false;
    }
  };

  // 加载用户所有数据
  const fetchAllUserData = async () => {
    error.value = null;

    // 只加载仍然存在的功能
    await fetchProfile();
    
    // 记录已被移除的功能
    console.log('游戏设置、用户设置、隐私设置和设备偏好功能已被移除');
  };

  // 判断手型大小
  const getHandSizeLabel = (handSize: string): string => {
    switch (handSize) {
      case 'small':
        return '小型手';
      case 'medium':
        return '中型手';
      case 'large':
        return '大型手';
      default:
        return '未知';
    }
  };

  // 判断握持方式
  const getGripStyleLabel = (gripStyle: string): string => {
    switch (gripStyle) {
      case 'palm':
        return '手掌握持';
      case 'claw':
        return '爪式握持';
      case 'fingertip':
        return '指尖握持';
      default:
        return '未知';
    }
  };

  // 格式化单位
  const formatWithUnit = (value: number, unit: string): string => {
    const formattedValue = value.toFixed(1);
    return `${formattedValue} ${unit}`;
  };

  // 计算DPI转换后的灵敏度
  const convertDPI = (sourceDPI: number, targetDPI: number, sensitivity: number): number => {
    return (sensitivity * sourceDPI) / targetDPI;
  };

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
  };
}
