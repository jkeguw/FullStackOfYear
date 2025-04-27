<template>
  <div class="profile-page">
    <div class="container mx-auto p-4">
      <el-tabs v-model="activeTab" class="profile-tabs" @tab-click="handleTabClick">
        <el-tab-pane label="个人资料" name="profile">
          <div class="mx-auto max-w-3xl">
            <profile-tab 
              :loading="profileLoading" 
              :profile="profile" 
              @update="updateUserProfile"
            />
          </div>
        </el-tab-pane>
        <el-tab-pane label="系统设置" name="settings">
          <div class="mx-auto max-w-3xl">
            <settings-tab 
              :loading="settingsLoading" 
              :settings="settings"
              @update="updateUserSettings"
            />
          </div>
        </el-tab-pane>
        <el-tab-pane label="我的设备" name="devices">
          <div class="mx-auto max-w-4xl">
            <device-preferences-tab 
              :loading="devicePreferencesLoading" 
              :preferences="devicePreferences"
              @add="addUserDevicePreference"
              @update="updateUserDevicePreference"
              @delete="deleteUserDevicePreference"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

// 导入页面组件
import ProfileTab from '@/components/profile/ProfileTab.vue'
import SettingsTab from '@/components/profile/SettingsTab.vue'
import DevicePreferencesTab from '@/components/profile/DevicePreferencesTab.vue'

// 导入组合式API
import { useProfile } from '@/composables/useProfile'
import { 
  type UpdateProfileRequest, 
  type UpdateGameSettingsRequest,
  type AddSensitivityConfigRequest,
  type UpdateSensitivityConfigRequest,
  type UpdateSettingsRequest,
  type UpdatePrivacySettingsRequest,
  type AddDevicePreferenceRequest,
  type UpdateDevicePreferenceRequest
} from '@/api/user'

// 路由
const route = useRoute()
const router = useRouter()

// 获取用户资料相关钩子
const { 
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
  fetchProfile,
  updateProfile,
  fetchGameSettings,
  updateGameSettings,
  addSensitivityConfig,
  updateSensitivityConfig,
  deleteSensitivityConfig,
  fetchSettings,
  updateSettings,
  fetchPrivacySettings,
  updatePrivacySettings,
  fetchDevicePreferences,
  addDevicePreference,
  updateDevicePreference,
  deleteDevicePreference
} = useProfile()

// 安全设置状态
const securityLoading = ref(false)

// 当前活动的标签页
const activeTab = ref('profile')

// 生命周期钩子
onMounted(() => {
  // 从URL获取当前标签
  const tab = route.query.tab as string
  if (tab && ['profile', 'game-settings', 'settings', 'privacy', 'security', 'devices'].includes(tab)) {
    activeTab.value = tab
  }
  
  // 加载数据
  loadTabData(activeTab.value)
})

// 监听标签页变化
watch(activeTab, (newTab) => {
  loadTabData(newTab)
  
  // 更新URL
  router.replace({
    query: { ...route.query, tab: newTab }
  })
})

// 根据标签页加载相应数据
const loadTabData = async (tab: string) => {
  switch (tab) {
    case 'profile':
      if (!profile.value) await fetchProfile()
      break
    case 'game-settings':
      if (!gameSettings.value) await fetchGameSettings()
      break
    case 'settings':
      if (!settings.value) await fetchSettings()
      break
    case 'privacy':
      if (!privacySettings.value) await fetchPrivacySettings()
      break
    case 'security':
      // 加载安全设置（2FA状态、设备列表等）
      // 在实际应用中，应该从API获取
      break
    case 'devices':
      if (devicePreferences.value.length === 0) await fetchDevicePreferences()
      break
  }
}

// 处理标签页点击
const handleTabClick = () => {
  // 已在 watch 中处理数据加载
}

// 更新用户个人资料
const updateUserProfile = async (data: UpdateProfileRequest) => {
  await updateProfile(data)
}

// 更新游戏设置
const updateUserGameSettings = async (data: UpdateGameSettingsRequest) => {
  await updateGameSettings(data)
}

// 添加灵敏度配置
const addUserSensitivityConfig = async (data: AddSensitivityConfigRequest) => {
  await addSensitivityConfig(data)
}

// 更新灵敏度配置
const updateUserSensitivityConfig = async (game: string, data: UpdateSensitivityConfigRequest) => {
  await updateSensitivityConfig(game, data)
}

// 删除灵敏度配置
const deleteUserSensitivityConfig = async (game: string) => {
  await deleteSensitivityConfig(game)
}

// 更新用户设置
const updateUserSettings = async (data: UpdateSettingsRequest) => {
  await updateSettings(data)
}

// 更新隐私设置
const updateUserPrivacySettings = async (data: UpdatePrivacySettingsRequest) => {
  await updatePrivacySettings(data)
}

// 添加设备偏好
const addUserDevicePreference = async (data: AddDevicePreferenceRequest) => {
  await addDevicePreference(data)
}

// 更新设备偏好
const updateUserDevicePreference = async (id: string, data: UpdateDevicePreferenceRequest) => {
  await updateDevicePreference(id, data)
}

// 删除设备偏好
const deleteUserDevicePreference = async (id: string) => {
  await deleteDevicePreference(id)
}

// 安全设置相关方法
// 设置两因素认证
const setupTwoFactor = async () => {
  securityLoading.value = true
  try {
    // 调用API设置两因素认证
    // 在实际应用中应从后端API获取
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('两因素认证设置初始化成功')
  } catch (error) {
    ElMessage.error('两因素认证设置失败')
    console.error(error)
  } finally {
    securityLoading.value = false
  }
}

// 验证并启用两因素认证
const verifyTwoFactor = async (code: string) => {
  securityLoading.value = true
  try {
    // 调用API验证两因素认证
    // 在实际应用中应从后端API验证
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('两因素认证已成功启用')
  } catch (error) {
    ElMessage.error('验证失败，请检查验证码是否正确')
    console.error(error)
  } finally {
    securityLoading.value = false
  }
}

// 禁用两因素认证
const disableTwoFactor = async (password: string) => {
  securityLoading.value = true
  try {
    // 调用API禁用两因素认证
    // 在实际应用中应调用后端API
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('两因素认证已禁用')
  } catch (error) {
    ElMessage.error('禁用两因素认证失败')
    console.error(error)
  } finally {
    securityLoading.value = false
  }
}

// 修改密码
const changePassword = async (data: { currentPassword: string, newPassword: string }) => {
  securityLoading.value = true
  try {
    // 调用API修改密码
    // 在实际应用中应调用后端API
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('密码已成功修改')
  } catch (error) {
    ElMessage.error('修改密码失败')
    console.error(error)
  } finally {
    securityLoading.value = false
  }
}

// 移除设备
const removeDevice = async (deviceId: string) => {
  securityLoading.value = true
  try {
    // 调用API移除设备
    // 在实际应用中应调用后端API
    await new Promise(resolve => setTimeout(resolve, 500))
    ElMessage.success('设备已成功移除')
  } catch (error) {
    ElMessage.error('移除设备失败')
    console.error(error)
  } finally {
    securityLoading.value = false
  }
}
</script>

<style scoped>
.profile-tabs {
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}
</style>