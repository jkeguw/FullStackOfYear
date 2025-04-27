/**
 * 设备相关工具函数
 */

/**
 * 生成设备ID
 * 简单生成一个随机ID
 * @returns 设备ID
 */
export function generateDeviceId(): string {
  const randomPart = Math.random().toString(36).substring(2, 15)
  const timestamp = Date.now().toString(36)
  
  // 仅使用时间戳和随机部分
  return [timestamp, randomPart].join('-')
}