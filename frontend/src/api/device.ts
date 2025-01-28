import request from '@/utils/request'
import type { Response } from './types'

export interface Device {
  id: string
  name: string
  brand: string
  type: string
  specs: Record<string, any>
}

export const getDevices = (params?: any) => {
  return request.get<Response<Device[]>>('/devices', { params })
    .then(res => res.data)
}

export const addDevice = (data: Omit<Device, 'id'>) => {
  return request.post<Response<Device>>('/devices', data)
}