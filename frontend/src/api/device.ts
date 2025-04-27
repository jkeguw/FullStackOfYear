import request from '@/utils/request'
import type { Response } from './types'

// 基础设备类型
export interface Device {
  id: string
  name: string
  brand: string
  type: string
  imageUrl?: string
  description?: string
  createdAt: string
  updatedAt: string
}

// 鼠标设备类型
export interface MouseDevice extends Device {
  dimensions: {
    length: number
    width: number
    height: number
    weight: number
    gripWidth?: number
  }
  shape: {
    type: string
    humpPlacement: string
    frontFlare: string
    sideCurvature: string
    handCompatibility: string
  }
  technical: {
    connectivity: string[]
    sensor: string
    maxDPI: number
    pollingRate: number
    sideButtons: number
    weight?: number
    battery?: {
      type: string
      capacity: number
      life: number
    }
  }
  recommended: {
    gameTypes: string[]
    gripStyles: string[]
    handSizes: string[]
    dailyUse: boolean
    professional: boolean
  }
  svgData?: {
    topView: string
    sideView: string
  }
}

// 键盘设备类型
export interface KeyboardDevice extends Device {
  layout: string
  switches: string[]
  size: string
  technical: {
    connectivity: string[]
    keycaps: string
    hotSwap: boolean
    rgbLighting: boolean
    nKeyRollover: boolean
  }
  recommended: {
    gameTypes: string[]
    dailyUse: boolean
    portable: boolean
  }
}

// 显示器设备类型
export interface MonitorDevice extends Device {
  size: number
  resolution: {
    width: number
    height: number
  }
  technical: {
    refreshRate: number
    responseTime: number
    panelType: string
    aspectRatio: string
    curvature?: number
    hdrSupport: boolean
    adaptiveSync?: string
  }
  recommended: {
    gameTypes: string[]
    content: string[]
    proUse: boolean
  }
}

// 鼠标垫设备类型
export interface MousepadDevice extends Device {
  size: {
    length: number
    width: number
    height: number
  }
  material: string
  surface: string
  base: string
  recommended: string[]
}

// 设备评测类型
export interface DeviceReview {
  id: string
  deviceId: string
  userId: string
  content: string
  pros: string[]
  cons: string[]
  score: number
  usage: string
  status: string
  createdAt: string
  updatedAt: string
}

// 用户设备配置类型
export interface UserDevice {
  id: string
  userId: string
  name: string
  description?: string
  devices: UserDeviceSetting[]
  isPublic: boolean
  createdAt: string
  updatedAt: string
}

export interface UserDeviceSetting {
  deviceId: string
  deviceType: string
  deviceName: string
  deviceBrand: string
  settings?: Record<string, any>
}

// 响应类型
export interface DeviceListResponse {
  total: number
  page: number
  pageSize: number
  devices: Device[]
}

export interface DeviceReviewListResponse {
  total: number
  page: number
  pageSize: number
  reviews: DeviceReview[]
}

export interface UserDeviceListResponse {
  total: number
  page: number
  pageSize: number
  userDevices: UserDevice[]
}

// 请求参数类型
export interface DeviceListParams {
  page?: number
  pageSize?: number
  type?: string
  brand?: string
  sortBy?: string
  sortOrder?: string
}

export interface DeviceReviewListParams {
  deviceId?: string
  userId?: string
  page?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: string
}

export interface UserDeviceListParams {
  userId?: string
  page?: number
  pageSize?: number
  isPublic?: boolean
  sortBy?: string
  sortOrder?: string
}

// 设备相关API

// 获取设备列表
export const getDevices = (params?: DeviceListParams) => {
  return request.get<Response<DeviceListResponse>>('/devices', { params })
    .then(res => res.data)
}

// 获取鼠标设备
export const getMouseDevice = (id: string) => {
  return request.get<Response<MouseDevice>>(`/devices/mouse/${id}`)
    .then(res => res.data)
}

// 比较结果接口
export interface ComparisonResult {
  mice: MouseDevice[]
  differences: {
    [key: string]: {
      property: string
      values: any[]
      differencePercent: number
    }
  }
  similarityScore: number
}

// 比较鼠标
export const compareMice = (ids: string[]) => {
  return request.get<Response<ComparisonResult>>(`/mice/compare?ids=${ids.join(',')}`)
    .then(res => res.data)
}

// 创建鼠标设备
export const createMouseDevice = (data: Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Response<MouseDevice>>('/devices/mouse', data)
    .then(res => res.data)
}

// 更新鼠标设备
export const updateMouseDevice = (id: string, data: Partial<Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>>) => {
  return request.put<Response<MouseDevice>>(`/devices/mouse/${id}`, data)
    .then(res => res.data)
}

// 删除设备
export const deleteDevice = (id: string) => {
  return request.delete<Response<null>>(`/devices/${id}`)
    .then(res => res.data)
}

// 设备评测相关API

// 获取设备评测列表
export const getDeviceReviews = (params?: DeviceReviewListParams) => {
  return request.get<Response<DeviceReviewListResponse>>('/device-reviews', { params })
    .then(res => res.data)
}

// 获取单条设备评测
export const getDeviceReview = (id: string) => {
  return request.get<Response<DeviceReview>>(`/device-reviews/${id}`)
    .then(res => res.data)
}

// 创建设备评测
export const createDeviceReview = (data: Omit<DeviceReview, 'id' | 'userId' | 'status' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Response<DeviceReview>>('/device-reviews', data)
    .then(res => res.data)
}

// 更新设备评测
export const updateDeviceReview = (id: string, data: Partial<Omit<DeviceReview, 'id' | 'userId' | 'deviceId' | 'status' | 'createdAt' | 'updatedAt'>>) => {
  return request.put<Response<DeviceReview>>(`/device-reviews/${id}`, data)
    .then(res => res.data)
}

// 删除设备评测
export const deleteDeviceReview = (id: string) => {
  return request.delete<Response<null>>(`/device-reviews/${id}`)
    .then(res => res.data)
}

// 用户设备配置相关API

// 获取用户设备配置列表
export const getUserDevices = (params?: UserDeviceListParams) => {
  return request.get<Response<UserDeviceListResponse>>('/user-devices', { params })
    .then(res => res.data)
}

// 获取单个用户设备配置
export const getUserDevice = (id: string) => {
  return request.get<Response<UserDevice>>(`/user-devices/${id}`)
    .then(res => res.data)
}

// 创建用户设备配置
export const createUserDevice = (data: Omit<UserDevice, 'id' | 'userId' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Response<UserDevice>>('/user-devices', data)
    .then(res => res.data)
}

// 更新用户设备配置
export const updateUserDevice = (id: string, data: Omit<UserDevice, 'id' | 'userId' | 'createdAt' | 'updatedAt'>) => {
  return request.put<Response<UserDevice>>(`/user-devices/${id}`, data)
    .then(res => res.data)
}

// 删除用户设备配置
export const deleteUserDevice = (id: string) => {
  return request.delete<Response<null>>(`/user-devices/${id}`)
    .then(res => res.data)
}