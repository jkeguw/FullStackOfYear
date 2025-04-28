// This file is maintained for development reference only
// All actual mouse data should be fetched from the backend API

import type { MouseDevice } from '@/models/MouseModel';

// Example mouse data structure for reference
export type ExampleMouseDataStructure = {
  id: string;
  name: string;
  brand: string;
  type: string;
  price: number;
  dimensions: {
    length: number;
    width: number;
    height: number;
    weight: number;
    gripWidth?: number;
  };
  weight: number;
  shape: {
    type: string;
    humpPlacement: string;
    frontFlare: string;
    sideCurvature: string;
    handCompatibility: string;
    thumbRest?: boolean;
  };
  technical: {
    connectivity: string[];
    sensor: string;
    maxDPI: number;
    pollingRate: number;
    sideButtons: number;
    weight: number;
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
  description: string;
};

// Empty array - data should be fetched from API
export const testMice: MouseDevice[] = [];

// Placeholder function - should not be used in production
export function getTestMouseData() {
  console.warn('getTestMouseData() is deprecated. Use API calls instead.');
  return {
    code: 0,
    message: 'success',
    data: {
      total: 0,
      page: 1,
      pageSize: 20,
      devices: []
    }
  };
}

// Helper function for tests only
export function mockApiResponse<T>(data: T, delay = 500): Promise<T> {
  console.warn('mockApiResponse() is deprecated. Use API calls instead.');
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, delay);
  });
}
