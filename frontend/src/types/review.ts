import { MouseDevice } from './mouse';

// 重新导出MouseDevice类型，以便其他组件可以从这里一起导入
export type { MouseDevice };

export interface ReviewAuthor {
  id: string;
  name: string;
  avatar?: string;
}

export interface ReviewImage {
  url: string;
  thumbnailUrl?: string;
  caption?: string;
}

export interface Review {
  id: string;
  userId: string;
  deviceId: string;
  title: string;
  content: string;
  summary: string;
  conclusion: string;
  pros: string[];
  cons: string[];
  score: number;
  rating: number;
  images: ReviewImage[];
  publishedAt: string;
  createdAt: string;
  updatedAt: string;
  type: string;
  contentType: string;
  viewCount: number;
  author: ReviewAuthor;
  mouse: MouseDevice;
}

export interface ReviewListParams {
  page: number;
  limit: number;
  sort: string;
  order: string;
  search: string;
  type: string;
  contentType: string;
}

export interface ReviewListResponse {
  reviews: Review[];
  total: number;
}
