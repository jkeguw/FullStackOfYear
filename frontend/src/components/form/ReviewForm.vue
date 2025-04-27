<template>
  <div class="review-form">
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef" 
      label-width="100px" 
      class="max-w-2xl"
      v-loading="loading"
    >
      <el-form-item label="设备" prop="deviceId">
        <el-select 
          v-model="form.deviceId" 
          filterable 
          placeholder="选择要评测的设备"
          :loading="devicesLoading"
          @change="handleDeviceChange"
        >
          <el-option 
            v-for="device in devices" 
            :key="device.id" 
            :label="`${device.brand} ${device.name}`" 
            :value="device.id"
          >
            <div class="flex items-center">
              <div v-if="device.imageUrl" class="mr-2">
                <el-image 
                  :src="device.imageUrl" 
                  style="width: 30px; height: 30px"
                  fit="contain"
                  :preview-src-list="[device.imageUrl]"
                ></el-image>
              </div>
              <div>
                <div>{{ device.brand }} {{ device.name }}</div>
                <div class="text-xs text-gray-500">{{ getDeviceTypeName(device.type) }}</div>
              </div>
            </div>
          </el-option>
        </el-select>
      </el-form-item>
      
      <template v-if="selectedDevice">
        <div class="selected-device mb-6 p-4 bg-gray-50 rounded-md">
          <div class="flex">
            <div v-if="selectedDevice.imageUrl" class="mr-4">
              <el-image 
                :src="selectedDevice.imageUrl" 
                style="width: 80px; height: 80px"
                fit="contain"
                :preview-src-list="[selectedDevice.imageUrl]"
              ></el-image>
            </div>
            <div>
              <div class="text-lg font-medium">{{ selectedDevice.brand }} {{ selectedDevice.name }}</div>
              <div class="text-sm text-gray-500">{{ getDeviceTypeName(selectedDevice.type) }}</div>
              <div v-if="selectedDevice.description" class="text-sm mt-2">{{ selectedDevice.description }}</div>
            </div>
          </div>
        </div>
      </template>
      
      <el-form-item label="使用场景" prop="usage">
        <el-select 
          v-model="form.usage" 
          placeholder="选择您的主要使用场景"
        >
          <el-option label="FPS 游戏" value="fps_gaming"></el-option>
          <el-option label="MOBA 游戏" value="moba_gaming"></el-option>
          <el-option label="MMO 游戏" value="mmo_gaming"></el-option>
          <el-option label="综合游戏" value="general_gaming"></el-option>
          <el-option label="办公工作" value="office_work"></el-option>
          <el-option label="创意设计" value="creative_design"></el-option>
          <el-option label="日常使用" value="daily_use"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="总体评分" prop="score">
        <el-rate 
          v-model="form.score" 
          :max="5" 
          :texts="['极差', '差', '一般', '好', '极好']"
          :colors="rateColors"
          show-text
          :allow-half="true"
        ></el-rate>
      </el-form-item>
      
      <el-form-item label="优点" prop="pros">
        <el-tag
          v-for="tag in form.pros"
          :key="tag"
          closable
          @close="handleRemovePro(tag)"
          class="mr-1 mb-1"
          effect="light"
          type="success"
        >
          {{ tag }}
        </el-tag>
        
        <el-input
          v-if="proInputVisible"
          ref="proInputRef"
          v-model="proInput"
          class="w-80 mt-1"
          @keyup.enter="handleAddPro"
          @blur="handleAddPro"
        />
        
        <el-button v-else size="small" @click="showProInput" class="mt-1">
          + 添加优点
        </el-button>
      </el-form-item>
      
      <el-form-item label="缺点" prop="cons">
        <el-tag
          v-for="tag in form.cons"
          :key="tag"
          closable
          @close="handleRemoveCon(tag)"
          class="mr-1 mb-1"
          effect="light"
          type="danger"
        >
          {{ tag }}
        </el-tag>
        
        <el-input
          v-if="conInputVisible"
          ref="conInputRef"
          v-model="conInput"
          class="w-80 mt-1"
          @keyup.enter="handleAddCon"
          @blur="handleAddCon"
        />
        
        <el-button v-else size="small" @click="showConInput" class="mt-1">
          + 添加缺点
        </el-button>
      </el-form-item>
      
      <el-form-item label="详细评测" prop="content">
        <el-input 
          v-model="form.content" 
          type="textarea" 
          :rows="8" 
          placeholder="请详细描述您对该设备的使用体验、功能评价和推荐理由等"
        ></el-input>
        
        <div class="mt-2 text-sm text-gray-500 flex justify-between items-center">
          <div>
            <el-progress 
              :percentage="contentProgress" 
              :status="contentProgressStatus"
              :stroke-width="8"
              class="w-48"
            ></el-progress>
          </div>
          <div>
            字数: {{ form.content.length }} / 最小要求: 50
          </div>
        </div>
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" @click="submitForm" :loading="loading">提交评测</el-button>
        <el-button @click="resetForm">重置</el-button>
        <el-button v-if="isEdit" type="danger" @click="handleDelete">删除</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { useDevice } from '@/composables/useDevice'
import type { Device, DeviceReview } from '@/api/device'

const props = defineProps<{
  reviewId?: string
}>()

const emit = defineEmits<{
  (e: 'saved', data: DeviceReview): void
  (e: 'deleted'): void
}>()

// 表单引用
const formRef = ref<FormInstance>()
const proInputRef = ref<HTMLInputElement>()
const conInputRef = ref<HTMLInputElement>()

// 使用设备钩子
const { 
  getDeviceTypeName,
  devices,
  deviceLoading: devicesLoading,
  fetchDevices,
  reviewLoading: loading,
  saveDeviceReview,
  fetchDeviceReview,
  updateReview,
  removeDeviceReview
} = useDevice()

// 标签输入状态
const proInputVisible = ref(false)
const conInputVisible = ref(false)
const proInput = ref('')
const conInput = ref('')

// 编辑模式判断
const isEdit = computed(() => !!props.reviewId)

// 选中的设备
const selectedDevice = ref<Device | null>(null)

// 评分颜色
const rateColors = ['#F56C6C', '#E6A23C', '#909399', '#67C23A', '#409EFF']

// 内容进度状态
const contentProgress = computed(() => {
  const progress = (form.content.length / 50) * 100
  return Math.min(progress, 100)
})

const contentProgressStatus = computed(() => {
  if (form.content.length >= 150) return 'success'
  if (form.content.length >= 50) return 'warning'
  return 'exception'
})

// 表单数据初始化
const form = reactive({
  deviceId: '',
  score: 0,
  usage: '',
  pros: [] as string[],
  cons: [] as string[],
  content: ''
})

// 表单验证规则
const rules = reactive<FormRules>({
  deviceId: [
    { required: true, message: '请选择要评测的设备', trigger: 'change' }
  ],
  usage: [
    { required: true, message: '请选择使用场景', trigger: 'change' }
  ],
  score: [
    { required: true, message: '请进行评分', trigger: 'change' },
    { type: 'number', min: 1, message: '请至少给出1星评价', trigger: 'change' }
  ],
  pros: [
    { type: 'array', min: 1, message: '请至少添加一个优点', trigger: 'change' }
  ],
  cons: [
    { type: 'array', min: 1, message: '请至少添加一个缺点', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请填写详细评测内容', trigger: 'blur' },
    { min: 50, message: '评测内容不能少于50个字符', trigger: 'blur' }
  ]
})

// 生命周期钩子
onMounted(async () => {
  // 获取设备列表
  await fetchDevices()
  
  // 如果是编辑模式，获取评测详情
  if (props.reviewId) {
    await fetchReviewData()
  }
})

// 获取评测数据
const fetchReviewData = async () => {
  if (!props.reviewId) return
  
  const review = await fetchDeviceReview(props.reviewId)
  if (review) {
    // 更新表单数据
    form.deviceId = review.deviceId
    form.score = review.score
    form.usage = review.usage
    form.pros = [...review.pros]
    form.cons = [...review.cons]
    form.content = review.content
    
    // 获取设备详情
    handleDeviceChange(review.deviceId)
  }
}

// 设备选择变更
const handleDeviceChange = (deviceId: string) => {
  const device = devices.value.find(d => d.id === deviceId)
  selectedDevice.value = device || null
}

// 添加优点标签
const showProInput = () => {
  proInputVisible.value = true
  nextTick(() => {
    proInputRef.value?.focus()
  })
}

const handleAddPro = () => {
  if (proInput.value) {
    if (form.pros.indexOf(proInput.value) === -1) {
      form.pros.push(proInput.value)
    }
  }
  proInputVisible.value = false
  proInput.value = ''
}

const handleRemovePro = (tag: string) => {
  const index = form.pros.indexOf(tag)
  if (index !== -1) {
    form.pros.splice(index, 1)
  }
}

// 添加缺点标签
const showConInput = () => {
  conInputVisible.value = true
  nextTick(() => {
    conInputRef.value?.focus()
  })
}

const handleAddCon = () => {
  if (conInput.value) {
    if (form.cons.indexOf(conInput.value) === -1) {
      form.cons.push(conInput.value)
    }
  }
  conInputVisible.value = false
  conInput.value = ''
}

const handleRemoveCon = (tag: string) => {
  const index = form.cons.indexOf(tag)
  if (index !== -1) {
    form.cons.splice(index, 1)
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid): Promise<void> => {
    if (valid) {
      let result: DeviceReview | null = null
      
      if (isEdit.value && props.reviewId) {
        // 更新评测
        result = await updateReview(props.reviewId, {
          content: form.content,
          pros: form.pros,
          cons: form.cons,
          score: form.score,
          usage: form.usage
        })
      } else {
        // 创建新评测
        result = await saveDeviceReview({
          deviceId: form.deviceId,
          content: form.content,
          pros: form.pros,
          cons: form.cons,
          score: form.score,
          usage: form.usage
        })
      }
      
      if (result) {
        emit('saved', result)
      }
    } else {
      ElMessage.warning('请完成表单验证')
    }
  })
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  selectedDevice.value = null
}

// 删除评测
const handleDelete = () => {
  if (!props.reviewId) return
  
  ElMessageBox.confirm(
    '确定要删除此评测吗？删除后无法恢复！',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
  .then(async () => {
    const success = await removeDeviceReview(props.reviewId!)
    if (success) {
      emit('deleted')
    }
  })
  .catch(() => {
    // 取消删除，不做任何操作
  })
}
</script>

<style scoped>
.review-form {
  margin-bottom: 20px;
}
</style>