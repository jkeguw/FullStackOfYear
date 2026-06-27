// 鼠标数据模型 - 与后端 models.MouseDevice 结构保持一致

// 鼠标形状类型
export type MouseShapeType = 'ergo' | 'ambi' | 'fingertip' | 'symmetrical' | 'asymmetrical';

// 握持方式
export type GripStyle = 'palm' | 'claw' | 'fingertip' | 'hybrid';

// 手部尺寸
export type HandSize = 'small' | 'medium' | 'large' | 'extra-large';

// 凸起位置
export type HumpPlacement = 'front' | 'center' | 'back' | 'none';

// 鼠标连接方式
export type Connectivity = 'wired' | 'wireless' | 'bluetooth' | 'dual';

// 鼠标形状特征
export interface MouseShape {
  type: MouseShapeType | string;
  humpPlacement: HumpPlacement | string;
  frontFlare: string;
  sideCurvature: string;
  handCompatibility: string;
  thumbRest?: boolean;
  ringFingerRest?: boolean;
}

// 鼠标技术规格
export interface MouseTechnical {
  connectivity: Connectivity[] | string[];
  sensor: string;
  maxDPI: number;
  pollingRate: number;
  sideButtons: number;
  weight?: number; // 冗余字段，部分数据源可能放在这里
  battery?: {
    type: string;
    capacity: number;
    life: number;
  };
}

// 鼠标推荐用途
export interface MouseRecommendation {
  gameTypes: string[];
  gripStyles: GripStyle[] | string[];
  handSizes: HandSize[] | string[];
  dailyUse: boolean;
  professional: boolean;
}

// 尺寸信息
export interface MouseDimensions {
  length: number;
  width: number;
  height: number;
  weight: number;
  gripWidth?: number;
}

// SVG数据
export interface MouseSvgData {
  topView: string;
  sideView: string;
}

// 鼠标设备 - 统一使用后端返回结构
export interface MouseDevice {
  id: string;
  name: string;
  brand: string;
  type: string;
  dimensions: MouseDimensions;
  shape: MouseShape;
  technical: MouseTechnical;
  recommended: MouseRecommendation;
  svgData?: MouseSvgData;
  imageUrl?: string;
  description?: string;
  createdAt?: string;
  updatedAt?: string;
  price?: number;
  weight?: number; // 冗余字段，兼容部分本地转换的代码
}

// 鼠标比较结果
export interface MouseComparisonResult {
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

// 鼠标比较模式
export type ComparisonMode = 'overlay' | 'sideBySide';

// 鼠标视图类型
export type ViewType = 'topView' | 'sideView';

// 比较状态
export interface ComparisonState {
  selectedMice: MouseDevice[];
  comparisonMode: ComparisonMode;
  viewType: ViewType;
  overlayOpacity: number;
  recentlyViewedMice: MouseDevice[];
}

// 完整的鼠标设备请求
export interface GetMouseDeviceResponse {
  code: number;
  message: string;
  data: MouseDevice;
}

// 鼠标列表请求
export interface GetMouseDevicesResponse {
  code: number;
  message: string;
  data: {
    total: number;
    page: number;
    pageSize: number;
    devices: MouseDevice[];
  };
}

// 比较请求
export interface CompareMouseResponse {
  code: number;
  message: string;
  data: MouseComparisonResult;
}

// 筛选参数
export interface MouseFilterParams {
  type?: MouseShapeType;
  brand?: string;
  weightMin?: number;
  weightMax?: number;
  connectivityType?: Connectivity;
  gripStyle?: GripStyle;
  sortBy?: 'price' | 'weight' | 'releaseDate' | 'popularity';
  sortOrder?: 'asc' | 'desc';
}
