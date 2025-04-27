<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold">账户登录</h2>
          <p class="text-gray-500">欢迎回来！请登录您的账户</p>
        </div>
        
        <!-- 登录表单 -->
        <div v-if="!requireTwoFactor">
          <el-form
            :model="loginForm"
            :rules="rules"
            ref="loginFormRef"
            label-position="top"
            @submit.prevent="handleLogin"
          >
            <el-form-item label="邮箱" prop="email">
              <el-input 
                v-model="loginForm.email" 
                placeholder="请输入邮箱"
              >
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <el-form-item label="密码" prop="password">
              <el-input 
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                show-password
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            
            <div class="flex justify-between items-center mb-4">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <el-button text>忘记密码？</el-button>
            </div>
            
            <el-form-item>
              <el-button 
                type="primary" 
                native-type="submit" 
                class="w-full" 
                :loading="loading"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
          
          <div class="divider my-6 text-center">
            <span class="px-2 bg-[#1E1E1E] text-gray-500">或</span>
          </div>
          
          <div class="social-login">
            <el-button 
              class="w-full mb-3"
              @click="oauthLogin('google')"
            >
              <template #icon>
                <img src="/google-icon.svg" class="w-5 h-5 mr-2" alt="Google" />
              </template>
              使用Google账号登录
            </el-button>
            
            <div class="text-center mt-6">
              <p class="text-gray-600">
                还没有账户？
                <router-link to="/register" class="text-blue-500 hover:text-blue-700">立即注册</router-link>
              </p>
            </div>
          </div>
        </div>
        
        <!-- 两因素认证验证 -->
        <div v-if="requireTwoFactor">
          <div class="text-center">
            <h3 class="text-xl font-medium mb-4">两步验证</h3>
            <p class="mb-6">请输入验证码以继续登录</p>
            <el-input 
              v-model="twoFactorCode" 
              placeholder="请输入6位验证码"
              class="mb-4 max-w-xs mx-auto"
              maxlength="6"
            ></el-input>
            <div class="flex justify-center gap-4">
              <el-button @click="cancelTwoFactor">取消</el-button>
              <el-button type="primary" @click="handleTwoFactorVerify(twoFactorCode)">验证</el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, FormInstance } from 'element-plus'
import { Lock, Message } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth'

// 路由相关
const router = useRouter()
const route = useRoute()

// 登录表单
const loginFormRef = ref<FormInstance>()
const loginForm = reactive({
  email: '',
  password: ''
})

// 规则
const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6个字符', trigger: 'blur' }
  ]
}

// 状态
const rememberMe = ref(false)
const loading = ref(false)
const requireTwoFactor = ref(false)
const twoFactorCode = ref('')

// 使用认证相关的composable
const { login, verifyTwoFactor, requireTwoFactor: authRequireTwoFactor } = useAuth()

// 生命周期钩子
onMounted(() => {
  // 从localStorage中恢复邮箱
  const savedEmail = localStorage.getItem('rememberedEmail')
  if (savedEmail) {
    loginForm.email = savedEmail
    rememberMe.value = true
  }
  
  // 处理来自其他页面的重定向
  const redirect = route.query.redirect as string
  if (redirect) {
    ElMessage.info('请先登录以继续')
  }
})

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      
      try {
        const result = await login(loginForm.email, loginForm.password)
        
        if (result.requireTwoFactor) {
          // 需要两因素认证
          requireTwoFactor.value = true
        } else if (result.success) {
          // 登录成功
          if (rememberMe.value) {
            localStorage.setItem('rememberedEmail', loginForm.email)
          } else {
            localStorage.removeItem('rememberedEmail')
          }
          
          ElMessage.success('登录成功')
          
          // 重定向到首页或来源页面
          const redirect = route.query.redirect as string || '/'
          router.replace(redirect)
        } else if (result.error) {
          ElMessage.error(result.error)
        }
      } catch (error) {
        console.error('登录失败', error)
        ElMessage.error('登录失败，请检查您的凭据')
      } finally {
        loading.value = false
      }
    }
  })
}

// 处理两因素认证验证
const handleTwoFactorVerify = async (code: string) => {
  loading.value = true
  
  try {
    const result = await verifyTwoFactor(code)
    
    if (result.success) {
      // 验证成功
      ElMessage.success('验证成功')
      
      // 保存邮箱（如果勾选了记住我）
      if (rememberMe.value) {
        localStorage.setItem('rememberedEmail', loginForm.email)
      }
      
      // 重定向到首页或来源页面
      const redirect = route.query.redirect as string || '/'
      router.replace(redirect)
    } else if (result.error) {
      ElMessage.error(result.error)
    }
  } catch (error) {
    console.error('验证失败', error)
    ElMessage.error('验证失败，请检查您的验证码')
  } finally {
    loading.value = false
  }
}

// 处理恢复码验证
const handleRecoveryCodeVerify = async (code: string) => {
  loading.value = true
  
  try {
    // 在实际应用中，应该有一个专用的恢复码验证API
    // 这里我们简单地使用相同的验证方法
    const result = await verifyTwoFactor(code)
    
    if (result.success) {
      // 验证成功
      ElMessage.success('恢复码验证成功')
      
      // 保存邮箱（如果勾选了记住我）
      if (rememberMe.value) {
        localStorage.setItem('rememberedEmail', loginForm.email)
      }
      
      // 重定向到首页或来源页面
      const redirect = route.query.redirect as string || '/'
      router.replace(redirect)
    } else if (result.error) {
      ElMessage.error(result.error)
    }
  } catch (error) {
    console.error('恢复码验证失败', error)
    ElMessage.error('恢复码验证失败，请检查您的恢复码')
  } finally {
    loading.value = false
  }
}

// 取消两因素认证
const cancelTwoFactor = () => {
  requireTwoFactor.value = false
}

// OAuth登录
const oauthLogin = async (provider: string) => {
  try {
    loading.value = true
    const { loginWithOAuth } = useAuth()
    await loginWithOAuth(provider)
    // OAuth登录通常会在新窗口完成，这里不需要额外处理
  } catch (error) {
    console.error('OAuth登录失败', error)
    ElMessage.error(`${provider}登录失败，请稍后重试`)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--claude-bg-dark);
  padding: 1rem;
}

.login-container {
  width: 100%;
  max-width: 480px;
}

.login-card {
  background-color: var(--claude-bg-medium);
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
  border: 1px solid var(--claude-border-dark);
  color: var(--claude-text-light);
}

.divider {
  position: relative;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: var(--claude-border-light);
  z-index: -1;
}
</style>