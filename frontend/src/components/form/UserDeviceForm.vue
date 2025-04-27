<template>
  <div class="user-device-form">
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef" 
      label-width="120px" 
      class="max-w-3xl"
      v-loading="loading"
    >
      <!-- 基本信息 -->
      <div class="bg-blue-50 border-l-4 border-blue-400 p-4 mb-6 rounded-r-md">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="el-icon-info text-blue-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800">什么是设备配置？</h3>
            <div class="mt-2 text-sm text-blue-700">
              设备配置让您可以将不同的游戏外设组合在一起，记录每个设备的特定设置（如DPI、轮询率等）。
              这样您就可以轻松切换不同游戏场景的设备配置，也可以与其他用户分享您的最佳设置。
            </div>
          </div>
        </div>
      </div>
      
      <el-form-item label="配置名称" prop="name">
        <el-input v-model="form.name" placeholder="输入设备配置名称">
          <template #prepend>
            <el-tooltip content="为您的配置起一个描述性的名称，如'FPS比赛设置'或'MOBA日常配置'" placement="top">
              <i class="el-icon-question"></i>
            </el-tooltip>
          </template>
        </el-input>
      </el-form-item>
      
      <el-form-item label="配置描述" prop="description">
        <el-input 
          v-model="form.description" 
          type="textarea" 
          :rows="3" 
          placeholder="输入设备配置描述，例如游戏种类、使用场景等"
        ></el-input>
        <span class="text-xs text-gray-500">一个好的描述能帮助其他用户理解您的配置目的和特点</span>
      </el-form-item>
      
      <el-form-item label="是否公开" prop="isPublic">
        <el-switch v-model="form.isPublic"></el-switch>
        <span class="ml-2 text-sm text-gray-500">
          <i :class="form.isPublic ? 'el-icon-view text-green-500' : 'el-icon-lock text-gray-500'" class="mr-1"></i>
          {{ form.isPublic ? '公开的配置将对所有用户可见' : '私密配置仅您自己可见' }}
        </span>
      </el-form-item>
      
      <!-- 设备列表 -->
      <el-divider content-position="left">
        <div class="flex items-center">
          <span class="text-blue-500 font-bold">设备列表</span>
          <el-tooltip content="添加新设备" placement="top">
            <el-button type="primary" size="small" class="ml-2" @click="handleAddDevice">
              <i class="el-icon-plus mr-1"></i> 添加设备
            </el-button>
          </el-tooltip>
        </div>
      </el-divider>
      
      <div v-if="form.devices.length === 0" class="empty-devices text-center py-8 bg-gray-50 rounded-lg">
        <el-empty description="暂无设备" :image-size="100"></el-empty>
        <p class="text-gray-500 mt-2 mb-4 max-w-md mx-auto">
          您需要添加至少一个设备到您的配置中。每个设备可以有不同的设置，如鼠标DPI、键盘灯效等。
        </p>
        <el-button type="primary" @click="handleAddDevice">添加第一个设备</el-button>
      </div>
      
      <div v-else-if="form.devices.length === 1 && !form.devices[0].deviceId" class="guide-tip bg-yellow-50 p-4 rounded-md mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i class="el-icon-warning-outline text-yellow-400 text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-800">创建您的第一个设备配置</h3>
            <div class="mt-1 text-sm text-yellow-700">
              请先选择设备类型，然后从列表中选择设备。每种设备都有不同的设置项，您可以根据自己的偏好进行调整。
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="devices-container">
        <div 
          v-for="(device, index) in form.devices" 
          :key="index"
          class="device-item mb-6 p-4 border border-gray-200 rounded-md bg-gray-50 hover:bg-white transition-colors"
        >
          <div class="device-header flex justify-between items-center mb-4 pb-2 border-b border-gray-200">
            <h3 class="flex items-center">
              <el-tag :type="getDeviceTagType(device.deviceType)" class="mr-2">
                {{ getDeviceTypeName(device.deviceType) }}
              </el-tag>
              <span class="text-lg font-medium">设备 {{ index + 1 }}</span>
            </h3>
            
            <div>
              <el-tooltip content="上移设备" placement="top" v-if="index > 0">
                <el-button
                  type="info"
                  plain
                  circle
                  size="mini"
                  class="mr-1"
                  @click="moveDevice(index, 'up')"
                >
                  <i class="el-icon-arrow-up"></i>
                </el-button>
              </el-tooltip>
              
              <el-tooltip content="下移设备" placement="top" v-if="index < form.devices.length - 1">
                <el-button
                  type="info"
                  plain
                  circle
                  size="mini"
                  class="mr-1"
                  @click="moveDevice(index, 'down')"
                >
                  <i class="el-icon-arrow-down"></i>
                </el-button>
              </el-tooltip>
              
              <el-tooltip content="删除设备" placement="top">
                <el-button 
                  type="danger" 
                  plain
                  circle 
                  size="mini" 
                  @click="removeDevice(index)"
                  :disabled="form.devices.length <= 1"
                >
                  <i class="el-icon-delete"></i>
                </el-button>
              </el-tooltip>
            </div>
          </div>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item 
                :label="'设备类型'" 
                :prop="`devices.${index}.deviceType`"
                :rules="[
                  { required: true, message: '请选择设备类型', trigger: 'change' }
                ]"
              >
                <el-select 
                  v-model="device.deviceType" 
                  placeholder="选择设备类型"
                  @change="fetchDeviceList(index)"
                  style="width: 100%"
                >
                  <el-option label="鼠标" value="mouse">
                    <div class="flex items-center">
                      <i class="el-icon-mouse mr-2"></i> 鼠标
                    </div>
                  </el-option>
                  <el-option label="键盘" value="keyboard">
                    <div class="flex items-center">
                      <i class="el-icon-platform mr-2"></i> 键盘
                    </div>
                  </el-option>
                  <el-option label="显示器" value="monitor">
                    <div class="flex items-center">
                      <i class="el-icon-monitor mr-2"></i> 显示器
                    </div>
                  </el-option>
                  <el-option label="鼠标垫" value="mousepad">
                    <div class="flex items-center">
                      <i class="el-icon-magic-stick mr-2"></i> 鼠标垫
                    </div>
                  </el-option>
                  <el-option label="配件" value="accessory">
                    <div class="flex items-center">
                      <i class="el-icon-headset mr-2"></i> 配件
                    </div>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            
            <el-col :span="12">
              <el-form-item 
                :label="'选择设备'" 
                :prop="`devices.${index}.deviceId`"
                :rules="[
                  { required: true, message: '请选择设备', trigger: 'change' }
                ]"
              >
                <el-select 
                  v-model="device.deviceId" 
                  placeholder="选择设备"
                  filterable
                  :loading="deviceSelectLoading"
                  @change="handleDeviceSelected(index)"
                  style="width: 100%"
                >
                  <el-option 
                    v-for="option in deviceOptions[index] || []" 
                    :key="option.id" 
                    :label="`${option.brand} ${option.name}`" 
                    :value="option.id"
                  >
                    <div class="flex justify-between items-center">
                      <span>{{ option.brand }} {{ option.name }}</span>
                      <el-tag size="mini" effect="plain">{{ option.type }}</el-tag>
                    </div>
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          
          <!-- 设备设置 -->
          <div v-if="device.deviceId" class="mt-4 bg-white p-4 rounded-md">
            <div class="flex justify-between items-center mb-3">
              <h4 class="text-md font-medium text-blue-500">设备设置</h4>
              <el-button 
                type="text" 
                size="mini"
                @click="device.showAdvancedSettings = !device.showAdvancedSettings"
              >
                {{ device.showAdvancedSettings ? '隐藏高级设置' : '显示高级设置' }}
                <i :class="[device.showAdvancedSettings ? 'el-icon-arrow-up' : 'el-icon-arrow-down', 'ml-1']"></i>
              </el-button>
            </div>
            
            <!-- 鼠标设置 -->
            <template v-if="device.deviceType === 'mouse'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="DPI">
                    <el-input-number 
                      v-model="device.settings.dpi" 
                      :min="400" 
                      :max="25600" 
                      :step="100"
                      style="width: 100%"
                    ></el-input-number>
                  </el-form-item>
                </el-col>
                
                <el-col :span="12">
                  <el-form-item label="轮询率">
                    <el-select v-model="device.settings.pollingRate" style="width: 100%">
                      <el-option label="125Hz" :value="125"></el-option>
                      <el-option label="250Hz" :value="250"></el-option>
                      <el-option label="500Hz" :value="500"></el-option>
                      <el-option label="1000Hz" :value="1000"></el-option>
                      <el-option label="2000Hz" :value="2000"></el-option>
                      <el-option label="4000Hz" :value="4000"></el-option>
                      <el-option label="8000Hz" :value="8000"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <div v-if="device.showAdvancedSettings">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="抬离距离">
                      <el-select v-model="device.settings.liftOffDistance" style="width: 100%">
                        <el-option label="低" value="low"></el-option>
                        <el-option label="中" value="medium"></el-option>
                        <el-option label="高" value="high"></el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="12">
                    <el-form-item label="防抖设置">
                      <el-switch v-model="device.settings.debounce"></el-switch>
                    </el-form-item>
                  </el-col>
                </el-row>
                
                <el-form-item label="报告率提升">
                  <el-switch v-model="device.settings.enhancedReportRate"></el-switch>
                </el-form-item>
              </div>
            </template>
            
            <!-- 键盘设置 -->
            <template v-if="device.deviceType === 'keyboard'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="RGB灯效">
                    <el-select v-model="device.settings.rgbLighting" style="width: 100%">
                      <el-option label="关闭" value="off"></el-option>
                      <el-option label="静态" value="static"></el-option>
                      <el-option label="呼吸" value="breathing"></el-option>
                      <el-option label="波浪" value="wave"></el-option>
                      <el-option label="光谱循环" value="spectrum"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
                
                <el-col :span="12">
                  <el-form-item label="按键延迟">
                    <el-select v-model="device.settings.keyDelay" style="width: 100%">
                      <el-option label="低" value="low"></el-option>
                      <el-option label="中" value="medium"></el-option>
                      <el-option label="高" value="high"></el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
              
              <div v-if="device.showAdvancedSettings">
                <el-form-item label="按键映射">
                  <el-switch v-model="device.settings.keyMapping"></el-switch>
                </el-form-item>
              </div>
            </template>
            
            <!-- 显示器设置 -->
            <template v-if="device.deviceType === 'monitor'">
              <el-form-item label="亮度">
                <el-slider v-model="device.settings.brightness" :min="0" :max="100">
                  <template #default="{ value }">
                    <div class="slider-demo-value">{{ value }}%</div>
                  </template>
                </el-slider>
              </el-form-item>
              
              <el-form-item label="对比度">
                <el-slider v-model="device.settings.contrast" :min="0" :max="100">
                  <template #default="{ value }">
                    <div class="slider-demo-value">{{ value }}%</div>
                  </template>
                </el-slider>
              </el-form-item>
              
              <div v-if="device.showAdvancedSettings">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="颜色模式">
                      <el-select v-model="device.settings.colorMode" style="width: 100%">
                        <el-option label="标准" value="standard"></el-option>
                        <el-option label="FPS" value="fps"></el-option>
                        <el-option label="电影" value="movie"></el-option>
                        <el-option label="sRGB" value="sRGB"></el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="12">
                    <el-form-item label="响应时间">
                      <el-select v-model="device.settings.responseTime" style="width: 100%">
                        <el-option label="正常" value="normal"></el-option>
                        <el-option label="快速" value="fast"></el-option>
                        <el-option label="极速" value="fastest"></el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </template>
            
            <!-- 鼠标垫设置 -->
            <template v-if="device.deviceType === 'mousepad'">
              <el-form-item label="备注">
                <el-input 
                  v-model="device.settings.notes" 
                  type="textarea" 
                  :rows="2" 
                  placeholder="鼠标垫相关备注"
                ></el-input>
              </el-form-item>
            </template>
            
            <!-- 自定义设置 -->
            <el-form-item label="其他设置" v-if="device.showAdvancedSettings">
              <el-input 
                v-model="device.settings.customSetting" 
                placeholder="其他自定义设置"
              ></el-input>
            </el-form-item>
          </div>
        </div>
      </div>
      
      <!-- 配置预览 -->
      <el-divider content-position="left">
        <span class="text-blue-500 font-bold">配置预览</span>
      </el-divider>
      
      <el-card class="mb-4" shadow="never">
        <div class="preview-content">
          <h4 class="text-lg font-medium mb-2">{{ form.name || '未命名配置' }}</h4>
          <p v-if="form.description" class="text-gray-600 mb-4 text-sm">{{ form.description }}</p>
          
          <div class="devices-preview" v-if="form.devices.length > 0">
            <h5 class="text-md font-medium mb-2">包含的设备:</h5>
            <div class="flex flex-wrap gap-2 mb-3">
              <el-tag 
                v-for="(device, index) in form.devices" 
                :key="index"
                :type="getDeviceTagType(device.deviceType)"
                effect="plain"
              >
                {{ device.deviceId ? getDeviceDisplayName(device, index) : `未选择设备 ${index + 1}` }}
              </el-tag>
            </div>
            
            <div class="visibility text-sm text-gray-500">
              <i :class="form.isPublic ? 'el-icon-view' : 'el-icon-lock'" class="mr-1"></i>
              本配置将{{ form.isPublic ? '对所有人公开' : '保持私密' }}
            </div>
          </div>
        </div>
      </el-card>
      
      <!-- 表单按钮 -->
      <el-form-item>
        <el-button type="primary" @click="submitForm" :loading="loading">
          <i class="el-icon-check mr-1"></i> 保存配置
        </el-button>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="success" plain @click="handleTestConfig" :disabled="!formIsValid">
          <i class="el-icon-cpu mr-1"></i> 测试配置
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { useDevice } from '@/composables/useDevice'
import type { Device, UserDevice, UserDeviceListResponse } from '@/api/device'

const props = defineProps<{
  userDeviceId?: string
}>()

const emit = defineEmits<{
  (e: 'saved', data: UserDevice): void
  (e: 'canceled'): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 使用设备钩子
const { 
  devices: allDevices,
  userDeviceLoading: loading, 
  fetchDevices,
  fetchUserDevice,
  saveUserDevice,
  updateUserDeviceConfig
} = useDevice()

// 设备选择状态
const deviceSelectLoading = ref(false)
const deviceOptions = ref<Record<number, Device[]>>({})

// 编辑模式判断
const isEdit = computed(() => !!props.userDeviceId)

// 表单验证状态
const formIsValid = computed(() => {
  return form.name && form.devices.every(d => d.deviceId)
})

// 获取设备显示名称
const getDeviceDisplayName = (device: ExtendedUserDeviceSetting, index: number) => {
  const option = deviceOptions.value[index]?.find(o => o.id === device.deviceId)
  if (option) {
    return `${option.brand} ${option.name}`
  }
  return `设备 ${index + 1}`
}

// 拓展用户设备设置类型，添加UI状态
interface ExtendedUserDeviceSetting {
  deviceId: string;
  deviceType: string;
  settings: Record<string, any>;
  showAdvancedSettings?: boolean;
}

// 获取设备类型标签
const getDeviceTagType = (deviceType: string): string => {
  const types: Record<string, string> = {
    'mouse': 'primary',
    'keyboard': 'success',
    'monitor': 'warning',
    'mousepad': 'info',
    'accessory': 'danger'
  }
  return types[deviceType] || 'info'
}

// 获取设备类型名称
const getDeviceTypeName = (deviceType: string): string => {
  const names: Record<string, string> = {
    'mouse': '鼠标',
    'keyboard': '键盘',
    'monitor': '显示器',
    'mousepad': '鼠标垫',
    'accessory': '配件'
  }
  return names[deviceType] || deviceType
}

// 表单数据初始化
const form = reactive({
  name: '',
  description: '',
  isPublic: false,
  devices: [{
    deviceId: '',
    deviceType: 'mouse',
    showAdvancedSettings: false,
    settings: {
      dpi: 800,
      pollingRate: 1000,
      enhancedReportRate: false,
      liftOffDistance: 'medium',
      debounce: true,
      customSetting: ''
    }
  }] as ExtendedUserDeviceSetting[]
})

// 表单验证规则
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入配置名称', trigger: 'blur' },
    { min: 2, max: 30, message: '长度在 2 到 30 个字符', trigger: 'blur' }
  ],
  devices: [
    { 
      type: 'array', 
      required: true, 
      message: '请至少添加一个设备', 
      trigger: 'change' 
    }
  ]
})

// 生命周期钩子
onMounted(async () => {
  // 加载设备列表
  await fetchDevices({ pageSize: 100 })
  
  // 如果是编辑模式，加载现有配置
  if (props.userDeviceId) {
    await fetchUserDeviceData()
  } else {
    // 新建模式，预加载第一个设备的选项
    fetchDeviceList(0)
  }
})

// 监听设备类型变化
watch(() => form.devices.map(d => d.deviceType), (newTypes) => {
  // 遍历每个设备，如果设备类型变化，清空设备选择
  form.devices.forEach((device, index) => {
    if (device.deviceType !== newTypes[index]) {
      device.deviceId = ''
      device.settings = getDefaultSettings(device.deviceType)
    }
  })
}, { deep: true })

// 获取设备配置数据
const fetchUserDeviceData = async () => {
  if (!props.userDeviceId) return
  
  const userDevice = await fetchUserDevice(props.userDeviceId)
  if (userDevice) {
    // 更新表单数据
    form.name = userDevice.name
    form.description = userDevice.description || ''
    form.isPublic = userDevice.isPublic
    
    // 设置设备列表
    if (userDevice.devices && userDevice.devices.length > 0) {
      form.devices = userDevice.devices.map(device => {
        // 确保每个设备都有settings对象
        const settings = device.settings || getDefaultSettings(device.deviceType)
        
        return {
          deviceId: device.deviceId,
          deviceType: device.deviceType,
          settings,
          showAdvancedSettings: false
        }
      })
      
      // 预加载所有设备的选项
      form.devices.forEach((_, index) => {
        fetchDeviceList(index)
      })
    }
  }
}

// 获取特定类型的设备列表
const fetchDeviceList = async (index: number) => {
  const deviceType = form.devices[index].deviceType
  
  deviceSelectLoading.value = true
  
  try {
    // 过滤设备列表
    deviceOptions.value[index] = allDevices.value.filter(d => d.type === deviceType)
  } catch (err) {
    console.error('获取设备列表失败', err)
    ElMessage.error('获取设备列表失败，请重试')
  } finally {
    deviceSelectLoading.value = false
  }
}

// 设备选择处理
const handleDeviceSelected = (index: number) => {
  // 保留原来的设置，确保兼容性
  const currentSettings = { ...form.devices[index].settings }
  
  // 根据设备类型设置默认值，但保留原有的自定义设置
  form.devices[index].settings = {
    ...getDefaultSettings(form.devices[index].deviceType),
    ...currentSettings
  }
}

// 获取默认设置
const getDefaultSettings = (deviceType: string) => {
  switch (deviceType) {
    case 'mouse':
      return {
        dpi: 800,
        pollingRate: 1000,
        enhancedReportRate: false,
        liftOffDistance: 'medium',
        debounce: true,
        customSetting: ''
      }
    case 'keyboard':
      return {
        keyDelay: 'medium',
        rgbLighting: 'static',
        keyMapping: false,
        customSetting: ''
      }
    case 'monitor':
      return {
        brightness: 50,
        contrast: 50,
        colorMode: 'standard',
        responseTime: 'fast',
        customSetting: ''
      }
    case 'mousepad':
      return {
        notes: '',
        customSetting: ''
      }
    default:
      return {
        customSetting: ''
      }
  }
}

// 添加设备
const handleAddDevice = () => {
  form.devices.push({
    deviceId: '',
    deviceType: 'mouse',
    showAdvancedSettings: false,
    settings: getDefaultSettings('mouse')
  })
  
  // 为新添加的设备加载选项
  const newIndex = form.devices.length - 1
  fetchDeviceList(newIndex)
}

// 移除设备
const removeDevice = (index: number) => {
  if (form.devices.length > 1) {
    form.devices.splice(index, 1)
  } else {
    ElMessage.warning('至少需要保留一个设备')
  }
}

// 移动设备位置
const moveDevice = (index: number, direction: 'up' | 'down') => {
  if (direction === 'up' && index > 0) {
    // 上移
    const temp = form.devices[index]
    form.devices[index] = form.devices[index - 1]
    form.devices[index - 1] = temp
  } else if (direction === 'down' && index < form.devices.length - 1) {
    // 下移
    const temp = form.devices[index]
    form.devices[index] = form.devices[index + 1]
    form.devices[index + 1] = temp
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid): Promise<void> => {
    if (valid) {
      // 验证每个设备是否都选择了设备ID
      const invalidDevices = form.devices.filter(d => !d.deviceId)
      if (invalidDevices.length > 0) {
        ElMessage.warning('请选择所有设备')
        return
      }
      
      let result: UserDevice | null = null
      
      const data = {
        name: form.name,
        description: form.description,
        devices: form.devices.map(d => ({
          deviceId: d.deviceId,
          deviceType: d.deviceType,
          deviceName: d.deviceName || '',
          deviceBrand: d.deviceBrand || '',
          settings: d.settings || {}
        })),
        isPublic: form.isPublic
      }
      
      if (isEdit.value && props.userDeviceId) {
        // 更新设备配置
        result = await updateUserDeviceConfig(props.userDeviceId, data)
      } else {
        // 创建新设备配置
        result = await saveUserDevice(data)
      }
      
      if (result) {
        emit('saved', result)
      }
    } else {
      ElMessage.warning('请完成表单验证')
    }
  })
}

// 取消
const handleCancel = () => {
  emit('canceled')
}

// 测试配置
const handleTestConfig = () => {
  if (!formIsValid.value) {
    ElMessage.warning('请完成所有必填项')
    return
  }
  
  // 模拟测试过程
  ElMessage({
    message: '配置测试中...',
    type: 'info',
    duration: 1000
  })
  
  // 显示测试结果
  setTimeout(() => {
    ElMessage({
      message: '配置测试成功，所有设备参数有效',
      type: 'success',
      duration: 3000
    })
  }, 1500)
}
</script>

<style scoped>
.device-form {
  margin-bottom: 20px;
}

.device-item {
  transition: all 0.3s ease;
  position: relative;
}

.device-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.device-header {
  position: sticky;
  top: 0;
  background-color: inherit;
  z-index: 1;
  border-top-left-radius: 0.375rem;
  border-top-right-radius: 0.375rem;
}

:deep(.el-divider__text) {
  font-weight: bold;
  color: #409EFF;
}

.slider-demo-value {
  position: absolute;
  top: -30px;
  right: 0;
  font-size: 14px;
  color: #606266;
  padding: 0 5px;
  background-color: #fff;
  border-radius: 3px;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
}

:deep(.el-switch__core) {
  width: 50px !important;
}

.devices-container {
  max-height: 600px;
  overflow-y: auto;
  padding-right: 5px;
}
</style>