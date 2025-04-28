import { createI18n } from 'vue-i18n';
import { nextTick } from 'vue';
import enUS from './locales/en-US';
import enJson from './locales/en.json';

// Always use English
const getPreferredLanguage = (): string => {
  return 'en-US';
};

// Only support English
export const SUPPORTED_LOCALES = [
  {
    code: 'en-US',
    name: 'English'
  }
];

export const DEFAULT_LOCALE = 'en-US';

// Create i18n instance
const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en-US',
  fallbackLocale: DEFAULT_LOCALE,
  messages: {
    'en-US': enUS,
    en: enJson
  }
});

// Set locale method (maintained for compatibility)
export async function setLocale(locale: string): Promise<void> {
  const targetLocale = 'en-US';  // Always use English
  
  // Set document language
  document.querySelector('html')?.setAttribute('lang', 'en');
  // Set cookie for backend
  document.cookie = `locale=${targetLocale}; path=/; max-age=${60 * 60 * 24 * 30}`;
  
  return nextTick();
}

export default i18n;
