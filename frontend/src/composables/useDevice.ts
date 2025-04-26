import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getDevices,
  getMouseDevice,
  createMouseDevice,
  updateMouseDevice,
  deleteDevice,
  getDeviceReviews,
  getDeviceReview,
  createDeviceReview,
  updateDeviceReview,
  deleteDeviceReview,
  type Device,
  type MouseDevice,
  type DeviceReview,
  type DeviceListParams,
  type DeviceReviewListParams
} from '@/api/device'

/**
 * 外设管理钩子
 * 用于管理设备、评测和用户配置的操作
 */
export function useDevice() {
  // 设备相关状态
  const devices = ref<Device[]>([])
  const currentDevice = ref<MouseDevice | null>(null)
  const deviceLoading = ref(false)
  const deviceError = ref<string | null>(null)
  const devicePagination = reactive({
    total: 0,
    page: 1,
    pageSize: 20
  })

  // 评测相关状态
  const reviews = ref<DeviceReview[]>([])
  const currentReview = ref<DeviceReview | null>(null)
  const reviewLoading = ref(false)
  const reviewError = ref<string | null>(null)
  const reviewPagination = reactive({
    total: 0,
    page: 1,
    pageSize: 20
  })

  // 已移除用户设备配置相关状态

  // 设备相关方法
  
  // 获取设备列表
  const fetchDevices = async (params?: DeviceListParams) => {
    deviceLoading.value = true
    deviceError.value = null

    try {
      const defaultParams = {
        page: devicePagination.page,
        pageSize: devicePagination.pageSize
      }
      
      const mergedParams = { ...defaultParams, ...params }
      const res = await getDevices(mergedParams)
      
      devices.value = res.data.devices
      devicePagination.total = res.data.total
      devicePagination.page = res.data.page
      devicePagination.pageSize = res.data.pageSize
      
      return res.data
    } catch (err) {
      console.error('获取设备列表失败', err)
      deviceError.value = '获取设备列表失败，请重试'
      ElMessage.error(deviceError.value)
      return null
    } finally {
      deviceLoading.value = false
    }
  }

  // 获取鼠标设备详情
  const fetchMouseDevice = async (id: string) => {
    deviceLoading.value = true
    deviceError.value = null

    try {
      const res = await getMouseDevice(id)
      currentDevice.value = res.data
      return res.data
    } catch (err) {
      console.error('获取设备详情失败', err)
      deviceError.value = '获取设备详情失败，请重试'
      ElMessage.error(deviceError.value)
      return null
    } finally {
      deviceLoading.value = false
    }
  }

  // 创建鼠标设备
  const saveMouseDevice = async (data: Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>) => {
    deviceLoading.value = true
    deviceError.value = null

    try {
      const res = await createMouseDevice(data)
      ElMessage.success('设备创建成功')
      return res.data
    } catch (err) {
      console.error('创建设备失败', err)
      deviceError.value = '创建设备失败，请重试'
      ElMessage.error(deviceError.value)
      return null
    } finally {
      deviceLoading.value = false
    }
  }

  // 更新鼠标设备
  const updateDevice = async (id: string, data: Partial<Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>>) => {
    deviceLoading.value = true
    deviceError.value = null

    try {
      const res = await updateMouseDevice(id, data)
      ElMessage.success('设备更新成功')
      return res.data
    } catch (err) {
      console.error('更新设备失败', err)
      deviceError.value = '更新设备失败，请重试'
      ElMessage.error(deviceError.value)
      return null
    } finally {
      deviceLoading.value = false
    }
  }

  // 删除设备
  const removeDevice = async (id: string) => {
    deviceLoading.value = true
    deviceError.value = null

    try {
      await deleteDevice(id)
      ElMessage.success('设备删除成功')
      return true
    } catch (err) {
      console.error('删除设备失败', err)
      deviceError.value = '删除设备失败，请重试'
      ElMessage.error(deviceError.value)
      return false
    } finally {
      deviceLoading.value = false
    }
  }

  // 评测相关方法
  
  // 获取评测列表
  const fetchDeviceReviews = async (params?: DeviceReviewListParams) => {
    reviewLoading.value = true
    reviewError.value = null

    try {
      const defaultParams = {
        page: reviewPagination.page,
        pageSize: reviewPagination.pageSize
      }
      
      const mergedParams = { ...defaultParams, ...params }
      const res = await getDeviceReviews(mergedParams)
      
      reviews.value = res.data.reviews
      reviewPagination.total = res.data.total
      reviewPagination.page = res.data.page
      reviewPagination.pageSize = res.data.pageSize
      
      return res.data
    } catch (err) {
      console.error('获取评测列表失败', err)
      reviewError.value = '获取评测列表失败，请重试'
      ElMessage.error(reviewError.value)
      return null
    } finally {
      reviewLoading.value = false
    }
  }

  // 获取评测详情
  const fetchDeviceReview = async (id: string) => {
    reviewLoading.value = true
    reviewError.value = null

    try {
      const res = await getDeviceReview(id)
      currentReview.value = res.data
      return res.data
    } catch (err) {
      console.error('获取评测详情失败', err)
      reviewError.value = '获取评测详情失败，请重试'
      ElMessage.error(reviewError.value)
      return null
    } finally {
      reviewLoading.value = false
    }
  }

  // 创建评测
  const saveDeviceReview = async (data: Omit<DeviceReview, 'id' | 'userId' | 'status' | 'createdAt' | 'updatedAt'>) => {
    reviewLoading.value = true
    reviewError.value = null

    try {
      const res = await createDeviceReview(data)
      ElMessage.success('评测创建成功，等待审核')
      return res.data
    } catch (err) {
      console.error('创建评测失败', err)
      reviewError.value = '创建评测失败，请重试'
      ElMessage.error(reviewError.value)
      return null
    } finally {
      reviewLoading.value = false
    }
  }

  // 更新评测
  const updateReview = async (id: string, data: Partial<Omit<DeviceReview, 'id' | 'userId' | 'deviceId' | 'status' | 'createdAt' | 'updatedAt'>>) => {
    reviewLoading.value = true
    reviewError.value = null

    try {
      const res = await updateDeviceReview(id, data)
      ElMessage.success('评测更新成功，等待审核')
      return res.data
    } catch (err) {
      console.error('更新评测失败', err)
      reviewError.value = '更新评测失败，请重试'
      ElMessage.error(reviewError.value)
      return null
    } finally {
      reviewLoading.value = false
    }
  }

  // 删除评测
  const removeDeviceReview = async (id: string) => {
    reviewLoading.value = true
    reviewError.value = null

    try {
      await deleteDeviceReview(id)
      ElMessage.success('评测删除成功')
      return true
    } catch (err) {
      console.error('删除评测失败', err)
      reviewError.value = '删除评测失败，请重试'
      ElMessage.error(reviewError.value)
      return false
    } finally {
      reviewLoading.value = false
    }
  }

  // 用户设备配置相关方法已移除

  // 通用帮助方法
  
  // 获取设备类型本地化名称
  const getDeviceTypeName = (type: string): string => {
    switch (type) {
      case 'mouse': return '鼠标'
      case 'keyboard': return '键盘'
      case 'monitor': return '显示器'
      case 'mousepad': return '鼠标垫'
      case 'accessory': return '配件'
      default: return '未知设备'
    }
  }

  // 获取手型本地化名称
  const getHandSizeName = (handSize: string): string => {
    switch (handSize) {
      case 'small': return '小型手'
      case 'medium': return '中型手'
      case 'large': return '大型手'
      default: return '未知'
    }
  }

  // 获取握持方式本地化名称
  const getGripStyleName = (gripStyle: string): string => {
    switch (gripStyle) {
      case 'palm': return '手掌握持'
      case 'claw': return '爪式握持'
      case 'fingertip': return '指尖握持'
      default: return '未知握持'
    }
  }

  return {
    // 设备相关状态和方法
    devices,
    currentDevice,
    deviceLoading,
    deviceError,
    devicePagination,
    fetchDevices,
    fetchMouseDevice,
    saveMouseDevice,
    updateDevice,
    removeDevice,
    
    // 评测相关状态和方法
    reviews,
    currentReview,
    reviewLoading,
    reviewError,
    reviewPagination,
    fetchDeviceReviews,
    fetchDeviceReview,
    saveDeviceReview,
    updateReview,
    removeDeviceReview,
    
    // 用户设备配置相关状态和方法已移除
    
    // 帮助方法
    getDeviceTypeName,
    getHandSizeName,
    getGripStyleName
  }
}