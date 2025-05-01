<template>
  <div class="register-page">
    <div class="flex min-h-screen bg-gray-100">
      <div class="m-auto w-full max-w-md p-6 bg-white rounded-lg shadow-md">
        <div class="text-center mb-8">
          <h2 class="text-2xl font-bold">Create Account</h2>
          <p class="text-gray-600 mt-2">Join our community and enjoy more features</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          @submit.prevent="handleRegister"
        >
          <el-form-item label="Username" prop="username">
            <el-input v-model="form.username" placeholder="Enter your username" :prefix-icon="User" />
          </el-form-item>

          <el-form-item label="Email" prop="email">
            <el-input v-model="form.email" placeholder="Enter your email" :prefix-icon="Message" />
          </el-form-item>

          <el-form-item label="Password" prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="Enter your password"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item label="Confirm Password" prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="Confirm your password"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>

          <el-form-item>
            <el-checkbox v-model="form.agreeTerms"
              >I have read and agree to the <a href="#" class="text-blue-500">Terms of Service</a> and <a
                href="#"
                class="text-blue-500"
                >Privacy Policy</a
              ></el-checkbox
            >
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              class="w-full"
              :loading="loading"
              :disabled="!form.agreeTerms"
              @click="handleRegister"
            >
              Register
            </el-button>
          </el-form-item>
        </el-form>

        <div class="text-center mt-6">
          <p class="text-gray-600">
            Already have an account?
            <router-link to="/login" class="text-blue-500 hover:text-blue-700"
              >Sign in now</router-link
            >
          </p>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { User, Message, Lock } from '@element-plus/icons-vue';
import { useAuth } from '@/composables/useAuth';

const router = useRouter();
const route = useRoute();
const { register, loginWithOAuth } = useAuth();

const formRef = ref();
const loading = ref(false);

// 表单数据
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeTerms: false
});

// 验证规则
const validatePass = (rule: any, value: string, callback: (error?: Error) => void) => {
  if (value === '') {
    callback(new Error('Please enter a password'));
  } else if (value.length < 8) {
    callback(new Error('Password must be at least 8 characters'));
  } else {
    if (form.confirmPassword !== '') {
      if (formRef.value) formRef.value.validateField('confirmPassword');
    }
    callback();
  }
};

const validateConfirmPass = (rule: any, value: string, callback: (error?: Error) => void) => {
  if (value === '') {
    callback(new Error('Please confirm your password'));
  } else if (value !== form.password) {
    callback(new Error('Passwords do not match'));
  } else {
    callback();
  }
};

const rules = {
  username: [
    { required: true, message: 'Please enter a username', trigger: 'blur' },
    { min: 3, max: 20, message: 'Length should be 3 to 20 characters', trigger: 'blur' }
  ],
  email: [
    { required: true, message: 'Please enter an email address', trigger: 'blur' },
    { type: 'email', message: 'Please enter a valid email address', trigger: 'blur' }
  ],
  password: [{ validator: validatePass, trigger: 'blur' }],
  confirmPassword: [{ validator: validateConfirmPass, trigger: 'blur' }]
};

// 处理注册
const handleRegister = async () => {
  if (!form.agreeTerms) {
    ElMessage.warning('Please read and agree to the Terms of Service and Privacy Policy');
    return;
  }

  try {
    if (!formRef.value) return;
    await formRef.value.validate();

    loading.value = true;
    await register({
      username: form.username,
      email: form.email,
      password: form.password,
      confirmPassword: form.confirmPassword
    });

    ElMessage.success('Registration successful! Please sign in.');

    // 注册成功后跳转到登录页面
    router.push('/login');
  } catch (error: any) {
    const message = error.response?.data?.message || 'Registration failed, please try again later';
    ElMessage.error(message);
  } finally {
    loading.value = false;
  }
};

// 第三方登录功能已移除
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
