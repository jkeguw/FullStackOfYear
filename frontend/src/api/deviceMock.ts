// Mock API implementation for testing
import { DeviceListParams, DeviceListResponse, MouseSVGResponse, SVGCompareRequest, SVGCompareResponse, SVGMouseListResponse } from './types';
import { getMockMouseData, getMockMouseSvg, getMockCompareSvgs, mockMice } from '@/data/mockMice';

// Export all original API types
export * from './device';

// API functions using mock data
export const getSVGMouseList = async (params?: { type?: string; brand?: string; views?: ('top' | 'side')[] }): Promise<SVGMouseListResponse> => {
  console.log('Using mock data: getSVGMouseList', params);
  return {
    code: 0,
    message: 'success',
    data: {
      devices: mockMice,
      total: mockMice.length
    }
  };
};

export const getMouseSVG = async (id: string, view: 'top' | 'side' = 'top'): Promise<MouseSVGResponse> => {
  console.log(`Using mock data: getMouseSVG, id=${id}, view=${view}`);
  const result = getMockMouseSvg(id, view);
  if (!result) {
    return {
      code: 404,
      message: 'Mouse not found',
      data: null
    };
  }
  return result;
};

export const compareSVGs = async (data: SVGCompareRequest): Promise<SVGCompareResponse> => {
  console.log(`Using mock data: compareSVGs, deviceIds=${data.deviceIds.join(',')}, view=${data.view}`);
  return getMockCompareSvgs(data.deviceIds, data.view);
};

export const getDevices = async (params?: DeviceListParams): Promise<DeviceListResponse> => {
  console.log('Using mock data: getDevices', params);
  const mockData = getMockMouseData();
  return {
    code: 0,
    message: 'success',
    data: mockData.data
  };
};

export const getMouseDevice = async (id: string) => {
  console.log(`Using mock data: getMouseDevice, id=${id}`);
  const mouse = mockMice.find(m => m.id === id);
  if (!mouse) {
    return {
      code: 404,
      message: 'Mouse not found',
      data: null
    };
  }
  
  return {
    code: 0,
    message: 'success',
    data: mouse
  };
};