<template>
  <div class="register-page">
    <div class="flex min-h-screen bg-gray-100">
      <div class="m-auto w-full max-w-md p-6 bg-white rounded-lg shadow-md">
        <div class="text-center mb-8">
          <h2 class="text-2xl font-bold">创建账号</h2>
          <p class="text-gray-600 mt-2">加入我们的社区，享受更多功能</p>
        </div>
        
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          @submit.prevent="handleRegister"
        >
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="form.username" 
              placeholder="请输入用户名"
              :prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input 
              v-model="form.email" 
              placeholder="请输入邮箱"
              :prefix-icon="Message"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="form.password" 
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="form.confirmPassword" 
              type="password"
              placeholder="请再次输入密码"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item>
            <el-checkbox v-model="form.agreeTerms">我已阅读并同意<a href="#" class="text-blue-500">服务条款</a>和<a href="#" class="text-blue-500">隐私政策</a></el-checkbox>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              class="w-full"
              :loading="loading"
              :disabled="!form.agreeTerms"
              @click="handleRegister"
            >
              注册
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="text-center mt-6">
          <p class="text-gray-600">
            已有账号？ 
            <router-link to="/login" class="text-blue-500 hover:text-blue-700">立即登录</router-link>
          </p>
        </div>
        
        <div class="divider my-6 flex items-center">
          <div class="flex-1 h-px bg-gray-300"></div>
          <span class="px-4 text-gray-500 text-sm">或通过第三方账号登录</span>
          <div class="flex-1 h-px bg-gray-300"></div>
        </div>
        
        <div class="flex justify-center space-x-4">
          <el-button class="social-button" @click="loginWithGoogle">
            <img src="https://upload.wikimedia.org/wikipedia/commons/5/53/Google_%22G%22_Logo.svg" alt="Google" class="w-5 h-5 mr-2" />
            Google
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Message, Lock } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const route = useRoute()
const { register, loginWithOAuth } = useAuth()

const formRef = ref()
const loading = ref(false)

// 表单数据
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeTerms: false
})

// 验证规则
const validatePass = (rule: any, value: string, callback: Function) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else if (value.length < 8) {
    callback(new Error('密码长度不能少于8个字符'))
  } else {
    if (form.confirmPassword !== '') {
      if (formRef.value) formRef.value.validateField('confirmPassword')
    }
    callback()
  }
}

const validateConfirmPass = (rule: any, value: string, callback: Function) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { validator: validatePass, trigger: 'blur' }
  ],
  confirmPassword: [
    { validator: validateConfirmPass, trigger: 'blur' }
  ]
}

// 处理注册
const handleRegister = async () => {
  if (!form.agreeTerms) {
    ElMessage.warning('请阅读并同意服务条款和隐私政策')
    return
  }
  
  try {
    if (!formRef.value) return
    await formRef.value.validate()
    
    loading.value = true
    await register({
      username: form.username,
      email: form.email,
      password: form.password
    })
    
    ElMessage.success('注册成功！请查收验证邮件')
    
    // 注册成功后，可以选择直接登录用户或者跳转到登录页面
    // 这里选择跳转到登录页面，让用户验证邮箱后登录
    router.push('/login')
    
  } catch (error: any) {
    const message = error.response?.data?.message || '注册失败，请稍后重试'
    ElMessage.error(message)
  } finally {
    loading.value = false
  }
}

// 第三方登录
const loginWithGoogle = async () => {
  try {
    loading.value = true
    await loginWithOAuth('google')
    
    // OAuth登录通常会直接将用户重定向到Google的认证页面
    // 所以这里的代码可能不会执行
  } catch (error: any) {
    const message = error.response?.data?.message || 'Google登录失败，请稍后重试'
    ElMessage.error(message)
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--background-color);
}

.social-button {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}
</style>