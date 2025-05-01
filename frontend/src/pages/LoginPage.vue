<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold">Sign In</h2>
          <p class="text-gray-500">Welcome back! Please sign in to your account</p>
        </div>

        <!-- 登录表单 -->
        <div v-if="!requireTwoFactor">
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="rules"
            label-position="top"
            @submit.prevent="handleLogin"
          >
            <el-form-item label="Email" prop="email">
              <el-input v-model="loginForm.email" placeholder="Enter your email">
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="Password" prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="Enter your password"
                show-password
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>

            <div class="flex justify-between items-center mb-4">
              <el-checkbox v-model="rememberMe">Remember me</el-checkbox>
              <el-button text>Forgot password?</el-button>
            </div>

            <el-form-item>
              <el-button type="primary" native-type="submit" class="w-full" :loading="loading">
                Sign In
              </el-button>
            </el-form-item>
          </el-form>

          <div class="text-center mt-6">
            <p class="text-gray-600">
              Don't have an account?
              <router-link to="/register" class="text-blue-500 hover:text-blue-700"
                >Register now</router-link
              >
            </p>
          </div>
        </div>

        <!-- 两因素认证验证 -->
        <div v-if="requireTwoFactor">
          <div class="text-center">
            <h3 class="text-xl font-medium mb-4">Two-Factor Authentication</h3>
            <p class="mb-6">Please enter the verification code to continue</p>
            <el-input
              v-model="twoFactorCode"
              placeholder="Enter 6-digit code"
              class="mb-4 max-w-xs mx-auto"
              maxlength="6"
            ></el-input>
            <div class="flex justify-center gap-4">
              <el-button @click="cancelTwoFactor">Cancel</el-button>
              <el-button type="primary" @click="handleTwoFactorVerify(twoFactorCode)"
                >Verify</el-button
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage, FormInstance } from 'element-plus';
import { Lock, Message } from '@element-plus/icons-vue';
import { useAuth } from '@/composables/useAuth';

// 路由相关
const router = useRouter();
const route = useRoute();

// 登录表单
const loginFormRef = ref<FormInstance>();
const loginForm = reactive({
  email: '',
  password: ''
});

// 规则
const rules = {
  email: [
    { required: true, message: 'Please enter your email', trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Please enter your password', trigger: 'blur' },
    { min: 6, message: 'Password must be at least 6 characters', trigger: 'blur' }
  ]
};

// 状态
const rememberMe = ref(false);
const loading = ref(false);
const requireTwoFactor = ref(false);
const twoFactorCode = ref('');

// 使用认证相关的composable
const { login, verifyTwoFactor, requireTwoFactor: authRequireTwoFactor } = useAuth();

// 生命周期钩子
onMounted(() => {
  // 从localStorage中恢复邮箱
  const savedEmail = localStorage.getItem('rememberedEmail');
  if (savedEmail) {
    loginForm.email = savedEmail;
    rememberMe.value = true;
  }

  // 处理来自其他页面的重定向
  const redirect = route.query.redirect as string;
  if (redirect) {
    ElMessage.info('Please sign in to continue');
  }
});

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return;

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;

      try {
        const result = await login(loginForm.email, loginForm.password);

        if (result.requireTwoFactor) {
          // 需要两因素认证
          requireTwoFactor.value = true;
        } else if (result.success) {
          // 登录成功
          if (rememberMe.value) {
            localStorage.setItem('rememberedEmail', loginForm.email);
          } else {
            localStorage.removeItem('rememberedEmail');
          }

          ElMessage.success('Signed in successfully');

          // 重定向到首页或来源页面
          const redirect = (route.query.redirect as string) || '/';
          router.replace(redirect);
        } else if (result.error) {
          ElMessage.error(result.error);
        }
      } catch (error) {
        console.error('Login failed', error);
        ElMessage.error('Login failed, please check your credentials');
      } finally {
        loading.value = false;
      }
    }
  });
};

// 处理两因素认证验证
const handleTwoFactorVerify = async (code: string) => {
  loading.value = true;

  try {
    const result = await verifyTwoFactor(code);

    if (result.success) {
      // 验证成功
      ElMessage.success('Verification successful');

      // 保存邮箱（如果勾选了记住我）
      if (rememberMe.value) {
        localStorage.setItem('rememberedEmail', loginForm.email);
      }

      // 重定向到首页或来源页面
      const redirect = (route.query.redirect as string) || '/';
      router.replace(redirect);
    } else if (result.error) {
      ElMessage.error(result.error);
    }
  } catch (error) {
    console.error('Verification failed', error);
    ElMessage.error('Verification failed, please check your code');
  } finally {
    loading.value = false;
  }
};

// 处理恢复码验证
const handleRecoveryCodeVerify = async (code: string) => {
  loading.value = true;

  try {
    // 在实际应用中，应该有一个专用的恢复码验证API
    // 这里我们简单地使用相同的验证方法
    const result = await verifyTwoFactor(code);

    if (result.success) {
      // 验证成功
      ElMessage.success('Recovery code verified successfully');

      // 保存邮箱（如果勾选了记住我）
      if (rememberMe.value) {
        localStorage.setItem('rememberedEmail', loginForm.email);
      }

      // 重定向到首页或来源页面
      const redirect = (route.query.redirect as string) || '/';
      router.replace(redirect);
    } else if (result.error) {
      ElMessage.error(result.error);
    }
  } catch (error) {
    console.error('Recovery code verification failed', error);
    ElMessage.error('Recovery code verification failed, please check your code');
  } finally {
    loading.value = false;
  }
};

// 取消两因素认证
const cancelTwoFactor = () => {
  requireTwoFactor.value = false;
};

// OAuth登录功能已移除
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
