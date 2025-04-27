import axios from 'axios'
import { useUserStore } from '@/stores'
import { ElMessage } from 'element-plus'
import { getTestMouseData, mockApiResponse } from '@/data/testMice'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 30000, // 延长超时时间，避免操作超时
  xsrfCookieName: 'XSRF-TOKEN',
  xsrfHeaderName: 'X-XSRF-TOKEN',
  withCredentials: true // 支持跨域请求时携带cookie，用于CSRF保护
})

// 为了支持灵敏度工具不登录也能使用的本地模式
const localSensitivityAPI = {
  // 三阶校准法
  '/sensitivity/three-stage': (data: any) => {
    // 基于数学公式计算三阶校准法的值
    const { initialValue = 40, dpi = 800, stage = 1, direction, currentBase } = data

    // 计算游戏对应的灵敏度
    const calculateGameSens = (cm360: number) => {
      // 游戏灵敏度常量
      const gameFactor = {
        csgo: 0.022, // CS2/CSGO 灵敏度常量
        valorant: 0.07, // Valorant 灵敏度常量
        overwatch: 0.176, // Overwatch 灵敏度常量
        apex: 0.022, // Apex 灵敏度常量
        rainbow6: 0.0633 // R6 灵敏度常量
      }
      
      const result: Record<string, any> = {}
      
      // 计算各游戏的灵敏度值
      for (const [game, factor] of Object.entries(gameFactor)) {
        const value = (360 * factor) / (cm360 * 0.01 * dpi)
        result[game] = { game, value }
      }
      
      return result
    }

    // 根据阶段和当前基准值计算左右值的偏移量
    const stageAdjustment = [
      { stage: 1, leftValue: 0.6, rightValue: 1.5 },
      { stage: 2, leftValue: 0.85, rightValue: 1.2 },
      { stage: 3, leftValue: 0.95, rightValue: 1.05 }
    ][stage - 1] || { stage: 1, leftValue: 0.6, rightValue: 1.5 }

    // 计算左右值
    const leftValue = currentBase * stageAdjustment.leftValue
    const rightValue = currentBase * stageAdjustment.rightValue

    // 若提供了方向，则计算新的基准值
    let newBase = currentBase
    const history = []

    if (direction) {
      newBase = direction === 'left' ? leftValue : rightValue
      
      history.push({
        stage,
        direction,
        baseValue: currentBase,
        leftValue,
        rightValue,
        newBase,
        timestamp: new Date().toLocaleString()
      })
    }

    return {
      data: {
        currentBase: newBase,
        leftValue,
        rightValue,
        gameSens: calculateGameSens(newBase),
        history
      }
    }
  },
  
  // 二分法
  '/sensitivity/binary-method': (data: any) => {
    const { initialValue = 10, currentBase, currentStep = 0, choice } = data
    
    // 二分法比率 (从初始值分叉的比例)
    const binaryRatios = [0.7, 1.5, 0.8, 1.3, 0.85, 1.2, 0.9, 1.1, 0.95, 1.05]
    
    // 默认值
    const base = currentBase || initialValue
    
    // 计算当前步骤的低值和高值
    const lowValue = base * 0.9
    const highValue = base * 1.1
    
    // 若是最后一步，则返回完成状态
    const isComplete = currentStep >= 8
    
    // 如果没有选择，则只返回初始计算结果
    if (!choice) {
      return {
        data: {
          currentBase: base,
          lowValue,
          highValue,
          currentStep,
          isComplete,
          history: []
        }
      }
    }
    
    // 根据选择计算新的基准值
    const newBase = choice === 'low' ? lowValue : highValue
    const nextStep = currentStep + 1
    
    // 构建历史记录
    const history = [{
      step: currentStep + 1,
      baseValue: base,
      lowValue,
      highValue,
      choice,
      newBase,
      timestamp: new Date().toLocaleString()
    }]
    
    // 如果完成，则返回最终值
    const finalValue = isComplete ? newBase : undefined
    
    return {
      data: {
        currentBase: newBase,
        lowValue: newBase * 0.9,
        highValue: newBase * 1.1,
        currentStep: nextStep,
        isComplete: nextStep >= 8,
        finalValue,
        history
      }
    }
  },
  
  // 极敏内推法
  '/sensitivity/interpolation': (data: any) => {
    const { fastValue, slowValue, targetType } = data
    
    let result
    let formula
    let calculation
    
    if (targetType === 'aiming') {
      // 瞄准场景: sqrt(k × m²) = √(k × m²)
      result = Math.sqrt(fastValue * slowValue * slowValue)
      formula = '瞄准场景公式: √(k × m²)'
      calculation = `√(${fastValue} × ${slowValue}²) = ${result.toFixed(2)}`
    } else {
      // 游戏场景: sqrt(k² × m) = √(k² × m)
      result = Math.sqrt(fastValue * fastValue * slowValue)
      formula = '游戏场景公式: √(k² × m)'
      calculation = `√(${fastValue}² × ${slowValue}) = ${result.toFixed(2)}`
    }
    
    return {
      data: {
        result,
        formula,
        calculation
      }
    }
  },
  
  // 灵敏度转换
  '/sensitivity/convert': (data: any) => {
    const { dpi, sensitivity, sourceGame, targetGame } = data
    
    // 游戏灵敏度常量
    const gameFactor: Record<string, number> = {
      csgo: 0.022, // CS2/CSGO 灵敏度常量
      valorant: 0.07, // Valorant 灵敏度常量
      overwatch: 0.176, // Overwatch 灵敏度常量
      apex: 0.022, // Apex 灵敏度常量
      rainbow6: 0.0633 // R6 灵敏度常量
    }
    
    // 计算 cm/360°
    const calculateCM360 = (gameName: string, sens: number) => {
      if (gameName === 'cm360') return sens
      const factor = gameFactor[gameName] || 0.022 // 默认使用 CSGO 的常量
      return (360 * factor) / (sens * dpi * 0.01)
    }
    
    // 计算特定游戏的灵敏度
    const calculateGameSens = (gameName: string, cm360Value: number) => {
      if (gameName === 'cm360') return cm360Value
      const factor = gameFactor[gameName] || 0.022 // 默认使用 CSGO 的常量
      return (360 * factor) / (cm360Value * dpi * 0.01)
    }
    
    // 计算 cm/360°
    const cm360Value = calculateCM360(sourceGame, sensitivity)
    
    // 计算目标游戏的灵敏度
    const targetValue = calculateGameSens(targetGame, cm360Value)
    
    // 计算其他游戏的灵敏度
    const otherGames = Object.keys(gameFactor)
      .filter(game => game !== sourceGame && game !== targetGame)
      .map(game => ({
        game,
        value: calculateGameSens(game, cm360Value)
      }))
    
    return {
      data: {
        sourceGame,
        sourceValue: sensitivity,
        targetGame,
        targetValue,
        cm360Value,
        otherGames
      }
    }
  },
  
  // 灵敏度常量
  '/sensitivity/constants': () => {
    return {
      data: {
        csgoFactor: 0.022,
        valorantFactor: 0.07,
        overwatchFactor: 0.176,
        apexFactor: 0.022,
        r6Factor: 0.0633,
        threeStageAdjustments: [
          { stage: 1, leftValue: 0.6, rightValue: 1.5 },
          { stage: 2, leftValue: 0.85, rightValue: 1.2 },
          { stage: 3, leftValue: 0.95, rightValue: 1.05 }
        ],
        binaryRatios: [0.7, 1.5, 0.8, 1.3, 0.85, 1.2, 0.9, 1.1, 0.95, 1.05]
      }
    }
  }
}

// 拦截请求，提供本地处理
request.interceptors.request.use(
  config => {
    const userStore = useUserStore()
    
    // 对灵敏度工具的API请求进行本地处理
    if (config.url && config.url.startsWith('/sensitivity/')) {
      const handler = localSensitivityAPI[config.url as keyof typeof localSensitivityAPI]
      
      if (handler && config.method === 'post' && config.data) {
        // 使用拦截函数记录这个请求将被本地处理
        config.adapter = (config) => {
          return new Promise((resolve) => {
            const result = handler(config.data)
            resolve({
              data: result,
              status: 200,
              statusText: 'OK',
              headers: {},
              config
            })
          })
        }
      } else if (handler && config.method === 'get') {
        // 处理GET请求
        config.adapter = (config) => {
          return new Promise((resolve) => {
            const result = handler({})
            resolve({
              data: result,
              status: 200,
              statusText: 'OK',
              headers: {},
              config
            })
          })
        }
      }
    }
    
    // 处理设备相关API
    if (config.url && config.url === '/devices' && config.method === 'get') {
      config.adapter = () => {
        return mockApiResponse({
          data: getTestMouseData(),
          status: 200,
          statusText: 'OK',
          headers: {},
          config
        }, 800)
      }
    }
    
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器：处理响应和错误
request.interceptors.response.use(
  response => response.data,
  async error => {
    const originalRequest = error.config
    const userStore = useUserStore()
    
    // 处理401错误（未授权/token过期）
    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        // 尝试使用refreshToken获取新的accessToken
        const refreshToken = sessionStorage.getItem('refreshToken')
        
        if (refreshToken) {
          // 发送刷新token请求
          const response = await axios.post(
            `${import.meta.env.VITE_API_URL}/auth/refresh`,
            { refreshToken },
            { headers: { 'Content-Type': 'application/json' } }
          )
          
          const { accessToken } = response.data.data
          
          // 更新store中的token
          userStore.setToken(accessToken)
          
          // 更新原始请求的Authorization头
          originalRequest.headers.Authorization = `Bearer ${accessToken}`
          
          // 重试原始请求
          return axios(originalRequest)
        }
      } catch (refreshError) {
        // 如果刷新token失败，登出用户
        userStore.clearUser()
        sessionStorage.removeItem('refreshToken')
        ElMessage.error('您的会话已过期，请重新登录')
        
        // 重定向到登录页面
        window.location.href = '/login'
      }
    }
    
    // 显示错误信息
    ElMessage.error(error.response?.data?.message || 'request failed')
    return Promise.reject(error)
  }
)

export default request
export { request }