import type { UploadRequestOptions } from 'element-plus'
import request from 'axios'

export const uploadFile = async ({ file }: UploadRequestOptions) => {
  const formData = new FormData()
  formData.append('file', file)
  return await request.post('/upload', formData)
}