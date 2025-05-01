import request from '@/utils/request';
import type { Response } from './types';

// 基础设备类型
export interface Device {
  id: string;
  name: string;
  brand: string;
  type: string;
  imageUrl?: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

// 鼠标设备类型
export interface MouseDevice extends Device {
  dimensions: {
    length: number;
    width: number;
    height: number;
    weight: number;
    gripWidth?: number;
  };
  shape: {
    type: string;
    humpPlacement: string;
    frontFlare: string;
    sideCurvature: string;
    handCompatibility: string;
  };
  technical: {
    connectivity: string[];
    sensor: string;
    maxDPI: number;
    pollingRate: number;
    sideButtons: number;
    weight?: number;
    battery?: {
      type: string;
      capacity: number;
      life: number;
    };
  };
  recommended: {
    gameTypes: string[];
    gripStyles: string[];
    handSizes: string[];
    dailyUse: boolean;
    professional: boolean;
  };
  svgData?: {
    topView: string;
    sideView: string;
  };
}

// 键盘设备类型
export interface KeyboardDevice extends Device {
  layout: string;
  switches: string[];
  size: string;
  technical: {
    connectivity: string[];
    keycaps: string;
    hotSwap: boolean;
    rgbLighting: boolean;
    nKeyRollover: boolean;
  };
  recommended: {
    gameTypes: string[];
    dailyUse: boolean;
    portable: boolean;
  };
}

// 显示器设备类型
export interface MonitorDevice extends Device {
  size: number;
  resolution: {
    width: number;
    height: number;
  };
  technical: {
    refreshRate: number;
    responseTime: number;
    panelType: string;
    aspectRatio: string;
    curvature?: number;
    hdrSupport: boolean;
    adaptiveSync?: string;
  };
  recommended: {
    gameTypes: string[];
    content: string[];
    proUse: boolean;
  };
}

// 鼠标垫设备类型
export interface MousepadDevice extends Device {
  size: {
    length: number;
    width: number;
    height: number;
  };
  material: string;
  surface: string;
  base: string;
  recommended: string[];
}

// 设备评测类型
export interface DeviceReview {
  id: string;
  deviceId: string;
  userId: string;
  content: string;
  pros: string[];
  cons: string[];
  score: number;
  usage: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

// 用户设备配置类型
export interface UserDevice {
  id: string;
  userId: string;
  name: string;
  description?: string;
  devices: UserDeviceSetting[];
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UserDeviceSetting {
  deviceId: string;
  deviceType: string;
  deviceName: string;
  deviceBrand: string;
  settings: Record<string, any>;
  showAdvancedSettings?: boolean;
}

// SVG 数据接口
export interface SVGResponse {
  deviceId: string;
  deviceName: string;
  brand: string;
  view: 'top' | 'side';
  svgData: string;
  scale?: number;
}

export interface SVGCompareRequest {
  deviceIds: string[];
  view: 'top' | 'side';
}

export interface SVGCompareResponse {
  devices: SVGResponse[];
  scale: number;
}

// 响应类型
export interface DeviceListResponse {
  total: number;
  page: number;
  pageSize: number;
  devices: Device[];
}

export interface DeviceReviewListResponse {
  total: number;
  page: number;
  pageSize: number;
  reviews: DeviceReview[];
}

export interface UserDeviceListResponse {
  total: number;
  page: number;
  pageSize: number;
  userDevices: UserDevice[];
}

// 请求参数类型
export interface DeviceListParams {
  page?: number;
  pageSize?: number;
  type?: string;
  brand?: string;
  sortBy?: string;
  sortOrder?: string;
}

export interface DeviceReviewListParams {
  deviceId?: string;
  userId?: string;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: string;
}

export interface UserDeviceListParams {
  userId?: string;
  page?: number;
  pageSize?: number;
  isPublic?: boolean;
  sortBy?: string;
  sortOrder?: string;
}

// 设备相关API

// 获取设备列表
export const getDevices = (params?: DeviceListParams) => {
  return request.get<Response<DeviceListResponse>>('/api/devices', { params })
    .then((res) => {
      if (!res) {
        console.error('Failed to fetch device list: No response');
        // 返回空数据结构而不是null
        return {
          devices: [],
          total: 0,
          page: params?.page || 1,
          pageSize: params?.pageSize || 20
        };
      }
      
      // 确保data存在，否则返回默认值
      if (!res.data) {
        console.warn('Fetch device list: API returned no data field');
        return {
          devices: [],
          total: 0,
          page: params?.page || 1,
          pageSize: params?.pageSize || 20
        };
      }
      
      // 检查响应结构并处理
      const responseData = res.data;
      if (responseData.code === 0 && responseData.data) {
        // 直接返回data字段中的内容
        return responseData.data;
      }
      
      return responseData;
    });
};

// 获取鼠标设备
export const getMouseDevice = (id: string) => {
  return request.get<Response<MouseDevice>>(`/api/devices/${id}`)
    .then((res) => {
      if (!res) {
        console.error('Failed to fetch mouse device: No response');
        throw new Error('Failed to fetch mouse device: No response');
      }
      
      // 如果res.data为null，返回一个空对象
      if (res.data === null) {
        console.warn(`Mouse device fetch (ID: ${id}) returned null, converted to empty object`);
        return {} as any;
      }
      
      // 检查响应结构并处理
      const responseData = res.data;
      if (responseData.code === 0 && responseData.data) {
        // 直接返回data字段中的内容
        return responseData.data;
      }
      
      return responseData;
    });
};

// 比较结果接口
export interface ComparisonResult {
  mice: MouseDevice[];
  differences: {
    [key: string]: {
      property: string;
      values: any[];
      differencePercent: number;
    };
  };
  similarityScore: number;
}

// 比较鼠标
export const compareMice = (ids: string[]) => {
  return request
    .get<Response<ComparisonResult>>(`/api/devices/mice/compare?ids=${ids.join(',')}`)
    .then((res) => {
      if (!res) {
        throw new Error('Mouse comparison failed: No results returned');
      }
      
      // 确保res.data存在，如果不存在则提供默认值
      if (!res.data) {
        console.warn('compareMice: API response missing data field, using default empty data');
        return {
          mice: [],
          differences: {},
          similarityScore: 0
        };
      }
      
      // 检查响应结构并处理
      const responseData = res.data;
      if (responseData.code === 0 && responseData.data) {
        // 直接返回data字段中的内容
        return responseData.data;
      }
      
      return responseData;
    });
};

// 创建鼠标设备
export const createMouseDevice = (
  data: Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>
) => {
  return request.post<Response<MouseDevice>>('/api/devices/mouse', data).then((res) => res.data);
};

// 更新鼠标设备
export const updateMouseDevice = (
  id: string,
  data: Partial<Omit<MouseDevice, 'id' | 'type' | 'createdAt' | 'updatedAt'>>
) => {
  return request.put<Response<MouseDevice>>(`/api/devices/mouse/${id}`, data).then((res) => res.data);
};

// 删除设备
export const deleteDevice = (id: string) => {
  return request.delete<Response<null>>(`/api/devices/${id}`).then((res) => res.data);
};

// SVG 相关 API

// 获取鼠标SVG数据
export const getMouseSVG = (id: string, view: 'top' | 'side' = 'top') => {
  return request.get<Response<SVGResponse>>(`/api/devices/mice/${id}/svg?view=${view}`).then((res) => res.data);
};

// 比较多个鼠标的SVG
export const compareSVGs = (data: SVGCompareRequest) => {
  // 对请求进行额外的日志记录，方便调试
  console.log('SVG比较请求数据:', data);
  return request.post<Response<SVGCompareResponse>>('/api/devices/mice/svg/compare', data)
    .then((res) => {
      console.log('SVG比较响应:', res);
      return res.data;
    })
    .catch(error => {
      console.error('SVG比较请求失败:', error);
      // 重新抛出错误以便上层处理
      throw error;
    });
};

// 获取有SVG数据的鼠标列表
export const getSVGMouseList = (params?: { type?: string; brand?: string; views?: ('top' | 'side')[] }) => {
  let queryParams = '';
  if (params) {
    const parts = [];
    if (params.type) parts.push(`type=${params.type}`);
    if (params.brand) parts.push(`brand=${params.brand}`);
    if (params.views && params.views.length > 0) parts.push(`views=${params.views.join(',')}`);
    if (parts.length > 0) queryParams = `?${parts.join('&')}`;
  }
  return request.get<Response<{ devices: Device[]; total: number }>>(`/api/devices/mice/svg/list${queryParams}`)
    .then((res) => {
      if (!res) {
        throw new Error('获取SVG鼠标列表失败: 没有返回结果');
      }
      
      // 确保res.data存在，如果不存在则提供默认值
      if (!res.data) {
        console.warn('getSVGMouseList: API返回结果缺少data字段，已使用默认空数据');
        return { devices: [], total: 0 };
      }
      
      return res.data;
    });
};

// 设备评测相关API

// 获取设备评测列表
export const getDeviceReviews = (params?: DeviceReviewListParams) => {
  return request
    .get<Response<DeviceReviewListResponse>>('/api/device-reviews', { params })
    .then((res) => res.data);
};

// 获取单条设备评测
export const getDeviceReview = (id: string) => {
  return request.get<Response<DeviceReview>>(`/api/device-reviews/${id}`).then((res) => res.data);
};

// 创建设备评测
export const createDeviceReview = (
  data: Omit<DeviceReview, 'id' | 'userId' | 'status' | 'createdAt' | 'updatedAt'>
) => {
  return request.post<Response<DeviceReview>>('/api/device-reviews', data).then((res) => res.data);
};

// 更新设备评测
export const updateDeviceReview = (
  id: string,
  data: Partial<
    Omit<DeviceReview, 'id' | 'userId' | 'deviceId' | 'status' | 'createdAt' | 'updatedAt'>
  >
) => {
  return request.put<Response<DeviceReview>>(`/api/device-reviews/${id}`, data).then((res) => res.data);
};

// 删除设备评测
export const deleteDeviceReview = (id: string) => {
  return request.delete<Response<null>>(`/api/device-reviews/${id}`).then((res) => res.data);
};

// 注意：用户设备配置相关API已被后端移除