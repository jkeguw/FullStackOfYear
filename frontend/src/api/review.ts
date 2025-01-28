import request from '@/utils/request'
import type { Response } from './types'

export interface Review {
  id: string
  userId: string
  deviceId: string
  content: string
  rating: number
  images?: string[]
}

export const getReviews = (params?: any) => {
  return request.get<Response<Review[]>>('/reviews', { params })
    .then(res => res.data)
}

export const addReview = (data: Omit<Review, 'id'>) => {
  return request.post<Response<Review>>('/reviews', data)
}