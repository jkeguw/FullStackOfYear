import { createI18n } from 'vue-i18n'
import { nextTick } from 'vue'
import enUS from './locales/en-US'
import zhCN from './locales/zh-CN'
import enJson from './locales/en.json'
import zhJson from './locales/zh.json'

// 获取首选语言的优先级：
// 1. localStorage 中存储的语言设置
// 2. 浏览器语言设置
// 3. 默认为英文
const getPreferredLanguage = (): string => {
  // 检查localStorage
  const storedLang = localStorage.getItem('language')
  if (storedLang) {
    return storedLang
  }

  // 检查浏览器语言
  const browserLang = navigator.language
  if (browserLang.startsWith('zh')) {
    return 'zh-CN'
  }

  // 默认为英文
  return 'en-US'
}

// 支持的语言列表
export const SUPPORTED_LOCALES = [
  {
    code: 'en-US',
    name: 'English'
  },
  {
    code: 'zh-CN',
    name: '中文'
  }
]

export const DEFAULT_LOCALE = 'en-US'

// 创建i18n实例
const i18n = createI18n({
  legacy: false, // 使用组合式API
  globalInjection: true, // 全局注入 $t, $d 等方法
  locale: getPreferredLanguage(),
  fallbackLocale: DEFAULT_LOCALE,
  messages: {
    'en-US': enUS,
    'zh-CN': zhCN,
    'en': enJson,
    'zh': zhJson
  }
})

// 设置语言的方法
export async function setLocale(locale: string): Promise<void> {
  const targetLocale = locale === 'zh' ? 'zh-CN' : locale === 'en' ? 'en-US' : locale;
  
  if (i18n.global.locale.value !== targetLocale) {
    i18n.global.locale.value = targetLocale
    // 存储到localStorage中
    localStorage.setItem('language', targetLocale)
    // 设置文档语言
    document.querySelector('html')?.setAttribute('lang', targetLocale.split('-')[0])
    // 设置Cookie方便后端识别
    document.cookie = `locale=${targetLocale}; path=/; max-age=${60 * 60 * 24 * 30}`
  }
  return nextTick()
}

export default i18n