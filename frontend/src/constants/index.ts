// status code
export const StatusCode = {
  SUCCESS: 200,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_ERROR: 500
} as const

// route name
export const RouteName = {
  HOME: 'Home',
  LOGIN: 'Login',
  REGISTER: 'Register',
  DEVICE_LIST: 'DeviceList',
  DEVICE_DETAIL: 'DeviceDetail',
  REVIEW_LIST: 'ReviewList'
} as const

// device type
export const DeviceType = {
  MOUSE: 'mouse',
  KEYBOARD: 'keyboard',
  MOUSEPAD: 'mousepad',
  HEADSET: 'headset'
} as const

// comment type
export const ReviewType = {
  NORMAL: 1,
  PROFESSIONAL: 2,
  VERIFIED: 3
} as const

// user character
export const UserRole = {
  USER: 'user',
  REVIEWER: 'reviewer',
  ADMIN: 'admin'
} as const