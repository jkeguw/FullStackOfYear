export interface MouseDimensions {
  length: number;
  width: number;
  height: number;
}

export interface MouseShape {
  type: string;
  gripStyles: string[];
}

export interface MouseBattery {
  type: string;
  capacity: number;
  life: number;
}

export interface MouseDevice {
  id: string;
  name: string;
  brand: string;
  type: string;
  dimensions: MouseDimensions;
  weight: number;
  shape: MouseShape;
  connectivity: string[];
  sensor: string;
  maxDPI: number;
  dpi: number;
  pollingRate: number;
  sideButtons: number;
  switches: string;
  battery?: MouseBattery;
  imageUrl?: string;
  description?: string;
  svgData?: {
    topView: string;
    sideView: string;
  };
  createdAt: string;
  updatedAt: string;
}

export interface MouseComparisonResult {
  similarityScore: number;
  differences: Record<string, {
    property: string;
    values: any[];
    differencePercent: number;
  }>;
}

export type ComparisonMode = 'overlay' | 'sideBySide';
export type ViewType = 'topView' | 'sideView';

export interface DeviceListResponse {
  devices: MouseDevice[];
  total: number;
}