import request from '@/utils/request'
import type { Response } from './types'
import type { User } from '@/types/user'

export const login = (data: { username: string; password: string }) => {
  return request.post<Response<{ token: string; user: User }>>('/users/login', data)
}

export const register = (data: { username: string; password: string }) => {
  return request.post<Response<User>>('/users/register', data)
}