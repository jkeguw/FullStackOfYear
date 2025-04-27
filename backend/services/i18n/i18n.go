package i18n

import (
	"fmt"
	"strings"
)

// Service 服务接口
type Service interface {
	GetSupportedLocales() []string
	GetDefaultLocale() string
	Translate(locale, key string, params map[string]interface{}) string
	T(locale, key string, params map[string]interface{}) string
}

type service struct {
	translations map[string]map[string]string
	locales      []string
	defaultLocale string
}

// DefaultService 默认i18n服务实现
type DefaultService struct {
	translations map[string]map[string]string
	locales      []string
	defaultLocale string
}

// NewService 创建新的i18n服务
func NewService() Service {
	svc := &service{
		translations: make(map[string]map[string]string),
		locales:      []string{"en-US", "zh-CN"},
		defaultLocale: "en-US",
	}
	
	// 简化：预设一些基本翻译
	svc.translations["en-US"] = map[string]string{
		"common.welcome": "Welcome",
		"common.hello": "Hello",
		"common.success": "Success",
		"common.error": "Error",
		"common.bad_request": "Bad Request",
	}
	
	svc.translations["zh-CN"] = map[string]string{
		"common.welcome": "欢迎",
		"common.hello": "你好",
		"common.success": "成功",
		"common.error": "错误",
		"common.bad_request": "请求错误",
	}
	
	return svc
}

// NewDefaultService 创建默认i18n服务
func NewDefaultService() *DefaultService {
	svc := &DefaultService{
		translations: make(map[string]map[string]string),
		locales:      []string{"en-US", "zh-CN"},
		defaultLocale: "en-US",
	}
	
	// 简化：预设一些基本翻译
	svc.translations["en-US"] = map[string]string{
		"common.welcome": "Welcome",
		"common.hello": "Hello",
		"common.success": "Success",
		"common.error": "Error",
		"common.bad_request": "Bad Request",
	}
	
	svc.translations["zh-CN"] = map[string]string{
		"common.welcome": "欢迎",
		"common.hello": "你好",
		"common.success": "成功",
		"common.error": "错误",
		"common.bad_request": "请求错误",
	}
	
	return svc
}

// GetSupportedLocales 返回支持的语言列表
func (s *service) GetSupportedLocales() []string {
	return s.locales
}

// GetDefaultLocale 返回默认语言
func (s *service) GetDefaultLocale() string {
	return s.defaultLocale
}

// Translate 翻译指定的键
func (s *service) Translate(locale, key string, params map[string]interface{}) string {
	// 如果语言不支持，使用默认语言
	if _, ok := s.translations[locale]; !ok {
		locale = s.defaultLocale
	}
	
	// 查找翻译文本
	translation, ok := s.translations[locale][key]
	if !ok {
		// 如果找不到翻译，返回键名
		return key
	}
	
	// 应用参数
	if params != nil {
		for k, v := range params {
			placeholder := fmt.Sprintf("{{%s}}", k)
			translation = strings.ReplaceAll(translation, placeholder, fmt.Sprintf("%v", v))
		}
	}
	
	return translation
}

// T 是 Translate 的别名，用于与已有代码兼容
func (s *service) T(locale, key string, params map[string]interface{}) string {
	return s.Translate(locale, key, params)
}

// Translate DefaultService实现
func (s *DefaultService) Translate(locale, key string, params map[string]interface{}) string {
	// 如果语言不支持，使用默认语言
	if _, ok := s.translations[locale]; !ok {
		locale = s.defaultLocale
	}
	
	// 查找翻译文本
	translation, ok := s.translations[locale][key]
	if !ok {
		// 如果找不到翻译，返回键名
		return key
	}
	
	// 应用参数
	if params != nil {
		for k, v := range params {
			placeholder := fmt.Sprintf("{{%s}}", k)
			translation = strings.ReplaceAll(translation, placeholder, fmt.Sprintf("%v", v))
		}
	}
	
	return translation
}

// T DefaultService实现
func (s *DefaultService) T(locale, key string, params map[string]interface{}) string {
	return s.Translate(locale, key, params)
}

// GetSupportedLocales DefaultService实现
func (s *DefaultService) GetSupportedLocales() []string {
	return s.locales
}

// GetDefaultLocale DefaultService实现
func (s *DefaultService) GetDefaultLocale() string {
	return s.defaultLocale
}