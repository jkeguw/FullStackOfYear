<template>
  <div class="profile-tab">
    <el-form
      :model="form"
      :rules="rules"
      ref="formRef"
      label-width="100px"
      v-loading="loading || submitting"
    >
      <div class="mb-6 text-center">
        <el-avatar :size="100" :src="avatarUrl">
          {{ form.username?.substring(0, 2).toUpperCase() }}
        </el-avatar>
        <div class="mt-2">
          <el-upload
            class="avatar-uploader"
            action="#"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleAvatarChange"
          >
            <el-button size="small" type="primary">更换头像</el-button>
          </el-upload>
        </div>
        <div v-if="profile" class="text-sm text-gray-500 mt-2">
          资料完成度: {{ profile.completedRate }}%
        </div>
      </div>
      
      <el-divider />
      
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" disabled />
      </el-form-item>
      
      <el-form-item label="电子邮箱" prop="email">
        <el-input v-model="form.email" disabled />
      </el-form-item>
      
      <el-form-item label="注册时间">
        <el-input v-model="createdAtFormatted" disabled />
      </el-form-item>
      
      <el-divider>个人信息</el-divider>
      
      <el-form-item label="个人简介" prop="bio">
        <el-input
          v-model="form.bio"
          type="textarea"
          :rows="4"
          placeholder="介绍一下你自己"
          maxlength="1000"
          show-word-limit
        />
      </el-form-item>
      
      <el-form-item label="性别" prop="gender">
        <el-select v-model="form.gender" placeholder="选择性别">
          <el-option label="男" value="male" />
          <el-option label="女" value="female" />
          <el-option label="其他" value="other" />
          <el-option label="不愿透露" value="prefer_not_to_say" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="位置" prop="location">
        <el-input v-model="form.location" placeholder="您的所在地" />
      </el-form-item>
      
      <el-form-item label="个人网站" prop="website">
        <el-input v-model="form.website" placeholder="您的个人网站或社交媒体链接" />
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" @click="submitForm" :loading="submitting">保存</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import type { UserProfile } from '@/types/user'
import type { UpdateProfileRequest } from '@/api/user'
import { format } from 'date-fns'

// 组件属性
const props = defineProps<{
  profile: UserProfile | null
  loading: boolean
}>()

// 组件事件
const emit = defineEmits<{
  (e: 'update', data: UpdateProfileRequest): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 提交状态
const submitting = ref(false)

// 缓存用户名和邮箱
const username = ref('')
const email = ref('')
const createdAt = ref<Date | null>(null)

// 表单数据
const form = reactive({
  username: '',
  email: '',
  avatar: '',
  bio: '',
  gender: '',
  location: '',
  website: ''
})

// 表单验证规则
const rules = reactive<FormRules>({
  bio: [
    { max: 1000, message: '个人简介不能超过1000个字符', trigger: 'blur' }
  ],
  website: [
    { pattern: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/, message: '请输入有效的网址', trigger: 'blur' }
  ]
})

// 计算属性
const avatarUrl = computed(() => form.avatar || '')

const createdAtFormatted = computed(() => {
  if (!createdAt.value) return ''
  return format(createdAt.value, 'yyyy-MM-dd HH:mm:ss')
})

// 监听属性变化
watch(() => props.profile, (newProfile) => {
  if (newProfile) {
    updateFormData(newProfile)
  }
}, { immediate: true })

// 生命周期钩子
onMounted(() => {
  if (props.profile) {
    updateFormData(props.profile)
  }
})

// 更新表单数据
const updateFormData = (profile: UserProfile) => {
  form.avatar = profile.avatar || ''
  form.bio = profile.bio || ''
  form.gender = profile.gender || ''
  form.location = profile.location || ''
  form.website = profile.website || ''
  
  // 从全局状态获取用户名和邮箱
  // 实际应用中，这些数据应该从用户状态中获取
  username.value = username.value || 'user123'
  email.value = email.value || 'user@example.com'
  form.username = username.value
  form.email = email.value
  
  // 设置创建时间
  if (typeof profile.createdAt === 'string') {
    createdAt.value = new Date(profile.createdAt)
  }
}

// 头像更改处理
const handleAvatarChange = (file: any) => {
  // 在实际应用中，应该将文件上传到服务器并获取URL
  // 这里简单地使用文件的 URL
  const reader = new FileReader()
  reader.onload = (e) => {
    if (e.target && typeof e.target.result === 'string') {
      form.avatar = e.target.result
    }
  }
  reader.readAsDataURL(file.raw)
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid): Promise<void> => {
    if (valid) {
      submitting.value = true
      try {
        // 发送更新请求
        const updateData: UpdateProfileRequest = {}
        
        // 只包含已更改的字段
        if (form.avatar !== (props.profile?.avatar || '')) {
          updateData.avatar = form.avatar
        }
        if (form.bio !== (props.profile?.bio || '')) {
          updateData.bio = form.bio
        }
        if (form.gender !== (props.profile?.gender || '')) {
          updateData.gender = form.gender
        }
        if (form.location !== (props.profile?.location || '')) {
          updateData.location = form.location
        }
        if (form.website !== (props.profile?.website || '')) {
          updateData.website = form.website
        }
        
        // 如果有更改，则发送更新请求
        if (Object.keys(updateData).length > 0) {
          emit('update', updateData)
        } else {
          ElMessage.info('没有需要更新的信息')
        }
      } finally {
        submitting.value = false
      }
    } else {
      ElMessage.warning('请修正表单中的错误')
      return false
    }
  })
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  if (props.profile) {
    updateFormData(props.profile)
  }
}
</script>

<style scoped>
.profile-tab {
  max-width: 600px;
  margin: 0 auto;
}

.avatar-uploader {
  display: inline-block;
}
</style>