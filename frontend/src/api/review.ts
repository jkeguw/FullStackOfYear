import request from '@/utils/request';
import type { Response } from './types';
import type { Review, ReviewListParams, ReviewListResponse } from '@/types/review';

// 重新导出Review类型供组件使用
export type { Review, ReviewListParams, ReviewListResponse };

export const getReviews = (params?: ReviewListParams) => {
  return request.get<Response<ReviewListResponse>>('/api/reviews', { params }).then((res) => res.data);
};

export const addReview = (data: Omit<Review, 'id'>) => {
  return request.post<Response<Review>>('/api/reviews', data);
};

export const getReview = (id: string) => {
  return request.get<Response<Review>>(`/api/reviews/${id}`).then((res) => res.data);
};
