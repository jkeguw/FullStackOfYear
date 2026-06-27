import request from '@/utils/request';
import type { Response } from './types';
import type {
  MouseDevice,
  MouseComparisonResult
} from '@/models/MouseModel';

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

// 鼠标设备类型从统一模型导入
export type { MouseDevice };

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

// 鼠标列表响应（设备列表中的鼠标设备）
export interface MouseListResponse {
  total: number;
  page: number;
  pageSize: number;
  devices: MouseDevice[];
}

// 对比/相似度相关类型
export interface PropertyDiff {
  property: string;
  values: any[];
  differencePercent: number;
}

export interface ComparisonResponse {
  mice: MouseDevice[];
  differences: Record<string, PropertyDiff>;
  similarityScore: number;
}

export interface SimilarMouse {
  mouse: MouseDevice;
  similarityScore: number;
  keyDifferences: PropertyDiff[];
}

export interface SimilarityResponse {
  reference: MouseDevice;
  similarMice: SimilarMouse[];
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

// 获取设备列表
export const getDevices = (params: DeviceListParams = {}) => {
  return request.get<Response<DeviceListResponse>>('/api/devices', { params });
};

// 获取鼠标设备详情
export const getMouseDevice = (id: string) => {
  return request.get<Response<MouseDevice>>(`/api/devices/mice/${id}`);
};

// 比较多个鼠标
export const compareMice = (ids: string[]) => {
  return request.get<Response<ComparisonResponse>>('/api/devices/mice/compare', {
    params: { ids: ids.join(',') }
  });
};

// 查找相似鼠标
export const findSimilarMice = (id: string, limit = 5) => {
  return request.get<Response<SimilarityResponse>>(`/api/devices/mice/${id}/similar`, {
    params: { limit }
  });
};

// 获取鼠标SVG
export const getMouseSVG = (id: string, view: 'top' | 'side') => {
  return request.get<Response<SVGResponse>>(`/api/devices/mice/${id}/svg`, {
    params: { view }
  });
};

// SVG比较
export const compareSVG = (data: SVGCompareRequest) => {
  return request.post<Response<SVGCompareResponse>>('/api/devices/svg/compare', data);
};

// 获取设备评测列表
export const getDeviceReviews = (deviceId: string, page = 1, pageSize = 10) => {
  return request.get<Response<DeviceReviewListResponse>>(`/api/device-reviews`, {
    params: { deviceId, page, pageSize }
  });
};
