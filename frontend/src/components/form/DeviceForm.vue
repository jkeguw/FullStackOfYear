<template>
  <div class="device-form">
    <el-form 
      :model="form" 
      :rules="rules" 
      ref="formRef" 
      label-position="top" 
      :label-width="isMobile ? 'auto' : '120px'" 
      class="max-w-3xl"
      v-loading="loading"
    >
      <!-- 基础信息 -->
      <el-divider content-position="left">基础信息</el-divider>
      
      <el-form-item label="设备类型" prop="type">
        <el-select v-model="form.type" placeholder="选择设备类型" disabled>
          <el-option label="鼠标" value="mouse"></el-option>
          <el-option label="键盘" value="keyboard" disabled></el-option>
          <el-option label="显示器" value="monitor" disabled></el-option>
          <el-option label="鼠标垫" value="mousepad" disabled></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="设备名称" prop="name">
        <el-input v-model="form.name" placeholder="输入设备名称"></el-input>
      </el-form-item>
      
      <el-form-item label="品牌" prop="brand">
        <el-input v-model="form.brand" placeholder="输入品牌名称"></el-input>
      </el-form-item>
      
      <el-form-item label="设备描述" prop="description">
        <el-input 
          v-model="form.description" 
          type="textarea" 
          :rows="3" 
          placeholder="输入设备描述和概况"
        ></el-input>
      </el-form-item>
      
      <el-form-item label="图片URL" prop="imageUrl">
        <el-input v-model="form.imageUrl" placeholder="输入图片URL链接"></el-input>
      </el-form-item>
      
      <!-- 尺寸信息 -->
      <el-divider content-position="left">尺寸信息</el-divider>
      
      <el-row :gutter="20">
        <el-col :span="24" :sm="12">
          <el-form-item label="长度(mm)" prop="dimensions.length">
            <el-input-number 
              v-model="form.dimensions.length" 
              :min="0" 
              :precision="1" 
              :step="1" 
              class="w-full" 
              controls-position="right"
            />
          </el-form-item>
        </el-col>
        <el-col :span="24" :sm="12">
          <el-form-item label="宽度(mm)" prop="dimensions.width">
            <el-input-number 
              v-model="form.dimensions.width" 
              :min="0" 
              :precision="1" 
              :step="1" 
              class="w-full" 
              controls-position="right"
            />
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24" :sm="12">
          <el-form-item label="高度(mm)" prop="dimensions.height">
            <el-input-number 
              v-model="form.dimensions.height" 
              :min="0" 
              :precision="1" 
              :step="1" 
              class="w-full" 
              controls-position="right"
            />
          </el-form-item>
        </el-col>
        <el-col :span="24" :sm="12">
          <el-form-item label="重量(g)" prop="dimensions.weight">
            <el-input-number 
              v-model="form.dimensions.weight" 
              :min="0" 
              :precision="1" 
              :step="1" 
              class="w-full" 
              controls-position="right"
            />
          </el-form-item>
        </el-col>
      </el-row>
      
      <!-- 形状信息 -->
      <el-divider content-position="left">形状信息</el-divider>
      
      <el-form-item label="形状类型" prop="shape.type">
        <el-select v-model="form.shape.type" placeholder="选择形状类型">
          <el-option label="人体工学" value="ergonomic"></el-option>
          <el-option label="左右对称" value="ambidextrous"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="坑位位置" prop="shape.humpPlacement">
        <el-select v-model="form.shape.humpPlacement" placeholder="选择坑位位置">
          <el-option label="前段" value="front"></el-option>
          <el-option label="中段" value="center"></el-option>
          <el-option label="后段" value="back"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="前端开叉" prop="shape.frontFlare">
        <el-select v-model="form.shape.frontFlare" placeholder="选择前端开叉">
          <el-option label="窄" value="narrow"></el-option>
          <el-option label="中等" value="medium"></el-option>
          <el-option label="宽" value="wide"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="侧面曲线" prop="shape.sideCurvature">
        <el-select v-model="form.shape.sideCurvature" placeholder="选择侧面曲线">
          <el-option label="直线" value="straight"></el-option>
          <el-option label="曲线" value="curved"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="手型适配" prop="shape.handCompatibility">
        <el-select v-model="form.shape.handCompatibility" placeholder="选择手型适配">
          <el-option label="小型手" value="small"></el-option>
          <el-option label="中型手" value="medium"></el-option>
          <el-option label="大型手" value="large"></el-option>
          <el-option label="广泛适配" value="universal"></el-option>
        </el-select>
      </el-form-item>
      
      <!-- 技术参数 -->
      <el-divider content-position="left">技术参数</el-divider>
      
      <el-form-item label="连接方式" prop="technical.connectivity">
        <el-checkbox-group v-model="form.technical.connectivity">
          <el-checkbox label="wired">有线</el-checkbox>
          <el-checkbox label="wireless">2.4G无线</el-checkbox>
          <el-checkbox label="bluetooth">蓝牙</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
      
      <el-form-item label="传感器" prop="technical.sensor">
        <el-input v-model="form.technical.sensor" placeholder="输入传感器型号"></el-input>
      </el-form-item>
      
      <el-row :gutter="20">
        <el-col :span="24" :sm="12">
          <el-form-item label="最大DPI" prop="technical.maxDPI">
            <el-input-number 
              v-model="form.technical.maxDPI" 
              :min="100" 
              :step="100" 
              class="w-full" 
              controls-position="right"
            />
          </el-form-item>
        </el-col>
        <el-col :span="24" :sm="12">
          <el-form-item label="轮询率(Hz)" prop="technical.pollingRate">
            <el-select v-model="form.technical.pollingRate" placeholder="选择轮询率" class="w-full">
              <el-option label="125Hz" :value="125"></el-option>
              <el-option label="250Hz" :value="250"></el-option>
              <el-option label="500Hz" :value="500"></el-option>
              <el-option label="1000Hz" :value="1000"></el-option>
              <el-option label="2000Hz" :value="2000"></el-option>
              <el-option label="4000Hz" :value="4000"></el-option>
              <el-option label="8000Hz" :value="8000"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-form-item label="侧键数量" prop="technical.sideButtons">
        <el-input-number 
          v-model="form.technical.sideButtons" 
          :min="0" 
          :max="20" 
          class="w-full sm:w-60" 
          controls-position="right"
        />
      </el-form-item>
      
      <!-- 电池信息 -->
      <template v-if="form.technical.connectivity.includes('wireless') || form.technical.connectivity.includes('bluetooth')">
        <el-divider content-position="left">电池信息</el-divider>
        
        <el-form-item label="电池类型" prop="technical.battery.type">
          <el-select v-model="form.technical.battery.type" placeholder="选择电池类型">
            <el-option label="内置锂电池" value="lithium"></el-option>
            <el-option label="可更换电池" value="replaceable"></el-option>
          </el-select>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="24" :sm="12">
            <el-form-item label="电池容量(mAh)" prop="technical.battery.capacity">
              <el-input-number 
                v-model="form.technical.battery.capacity" 
                :min="0" 
                class="w-full" 
                controls-position="right"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24" :sm="12">
            <el-form-item label="电池寿命(小时)" prop="technical.battery.life">
              <el-input-number 
                v-model="form.technical.battery.life" 
                :min="0" 
                class="w-full" 
                controls-position="right"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </template>
      
      <!-- 推荐使用场景 -->
      <el-divider content-position="left">推荐使用场景</el-divider>
      
      <el-form-item label="游戏类型" prop="recommended.gameTypes">
        <el-select
          v-model="form.recommended.gameTypes"
          multiple
          collapse-tags
          placeholder="选择适合的游戏类型"
        >
          <el-option label="FPS射击游戏" value="fps"></el-option>
          <el-option label="MOBA游戏" value="moba"></el-option>
          <el-option label="RTS策略游戏" value="rts"></el-option>
          <el-option label="MMO角色扮演" value="mmo"></el-option>
          <el-option label="竞速游戏" value="racing"></el-option>
          <el-option label="格斗游戏" value="fighting"></el-option>
          <el-option label="休闲游戏" value="casual"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="握持方式" prop="recommended.gripStyles">
        <el-select
          v-model="form.recommended.gripStyles"
          multiple
          collapse-tags
          placeholder="选择适合的握持方式"
        >
          <el-option label="手掌握持" value="palm"></el-option>
          <el-option label="爪式握持" value="claw"></el-option>
          <el-option label="指尖握持" value="fingertip"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="手型大小" prop="recommended.handSizes">
        <el-select
          v-model="form.recommended.handSizes"
          multiple
          collapse-tags
          placeholder="选择适合的手型大小"
        >
          <el-option label="小型手" value="small"></el-option>
          <el-option label="中型手" value="medium"></el-option>
          <el-option label="大型手" value="large"></el-option>
        </el-select>
      </el-form-item>
      
      <el-form-item label="适合日常使用" prop="recommended.dailyUse">
        <el-switch v-model="form.recommended.dailyUse"></el-switch>
      </el-form-item>
      
      <el-form-item label="专业级设备" prop="recommended.professional">
        <el-switch v-model="form.recommended.professional"></el-switch>
      </el-form-item>
      
      <!-- 表单按钮 -->
      <el-form-item>
        <div class="flex flex-wrap gap-2">
          <el-button type="primary" @click="submitForm" :loading="loading">保存</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button v-if="isEdit" type="danger" @click="handleDelete">删除</el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus'
import { useDevice } from '@/composables/useDevice'
import type { MouseDevice } from '@/api/device'
import { useWindowSize } from '@vueuse/core'

const props = defineProps<{
  deviceId?: string
}>()

const emit = defineEmits<{
  (e: 'saved', data: MouseDevice): void
  (e: 'deleted'): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 响应式判断
const { width } = useWindowSize()
const isMobile = computed(() => width.value < 640)

// 使用设备钩子
const { 
  deviceLoading: loading, 
  saveMouseDevice, 
  fetchMouseDevice, 
  updateDevice, 
  removeDevice 
} = useDevice()

// 编辑模式判断
const isEdit = computed(() => !!props.deviceId)

// 表单数据初始化
const form = reactive({
  type: 'mouse',
  name: '',
  brand: '',
  description: '',
  imageUrl: '',
  dimensions: {
    length: 120,
    width: 60,
    height: 40,
    weight: 80
  },
  shape: {
    type: 'ambidextrous',
    humpPlacement: 'center',
    frontFlare: 'medium',
    sideCurvature: 'curved',
    handCompatibility: 'medium'
  },
  technical: {
    connectivity: ['wired'],
    sensor: '',
    maxDPI: 16000,
    pollingRate: 1000,
    sideButtons: 2,
    battery: {
      type: 'lithium',
      capacity: 500,
      life: 70
    }
  },
  recommended: {
    gameTypes: ['fps'],
    gripStyles: ['palm', 'claw'],
    handSizes: ['medium'],
    dailyUse: true,
    professional: false
  }
})

// 表单验证规则
const rules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入设备名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  brand: [
    { required: true, message: '请输入品牌名称', trigger: 'blur' },
    { min: 2, max: 30, message: '长度在 2 到 30 个字符', trigger: 'blur' }
  ],
  'dimensions.length': [
    { required: true, message: '请输入长度', trigger: 'blur' },
    { type: 'number', min: 0, message: '长度必须大于0', trigger: 'blur' }
  ],
  'dimensions.width': [
    { required: true, message: '请输入宽度', trigger: 'blur' },
    { type: 'number', min: 0, message: '宽度必须大于0', trigger: 'blur' }
  ],
  'dimensions.height': [
    { required: true, message: '请输入高度', trigger: 'blur' },
    { type: 'number', min: 0, message: '高度必须大于0', trigger: 'blur' }
  ],
  'dimensions.weight': [
    { required: true, message: '请输入重量', trigger: 'blur' },
    { type: 'number', min: 0, message: '重量必须大于0', trigger: 'blur' }
  ],
  'technical.connectivity': [
    { required: true, message: '请选择至少一种连接方式', trigger: 'change' },
    { type: 'array', min: 1, message: '请选择至少一种连接方式', trigger: 'change' }
  ],
  'technical.sensor': [
    { required: true, message: '请输入传感器型号', trigger: 'blur' }
  ],
  'technical.maxDPI': [
    { required: true, message: '请输入最大DPI', trigger: 'blur' },
    { type: 'number', min: 100, message: 'DPI必须大于100', trigger: 'blur' }
  ],
  'technical.pollingRate': [
    { required: true, message: '请选择轮询率', trigger: 'change' }
  ],
  'recommended.gameTypes': [
    { type: 'array', min: 1, message: '请选择至少一种游戏类型', trigger: 'change' }
  ],
  'recommended.gripStyles': [
    { type: 'array', min: 1, message: '请选择至少一种握持方式', trigger: 'change' }
  ],
  'recommended.handSizes': [
    { type: 'array', min: 1, message: '请选择至少一种手型大小', trigger: 'change' }
  ]
})

// 生命周期钩子
onMounted(async () => {
  if (props.deviceId) {
    await fetchDeviceData()
  }
})

// 获取设备数据
const fetchDeviceData = async () => {
  if (!props.deviceId) return
  
  const device = await fetchMouseDevice(props.deviceId)
  if (device) {
    // 更新表单数据
    form.name = device.name
    form.brand = device.brand
    form.description = device.description || ''
    form.imageUrl = device.imageUrl || ''
    form.dimensions = { ...device.dimensions }
    form.shape = { ...device.shape }
    // @ts-ignore - Optional vs required property issue
    form.technical = { ...device.technical }
    form.recommended = { ...device.recommended }
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid): Promise<void> => {
    if (valid) {
      let result: MouseDevice | null = null
      
      if (isEdit.value && props.deviceId) {
        // 更新设备
        result = await updateDevice(props.deviceId, {
          name: form.name,
          brand: form.brand,
          description: form.description,
          imageUrl: form.imageUrl,
          dimensions: form.dimensions,
          shape: form.shape,
          technical: form.technical,
          recommended: form.recommended
        })
      } else {
        // 创建新设备
        result = await saveMouseDevice({
          name: form.name,
          brand: form.brand,
          description: form.description,
          imageUrl: form.imageUrl,
          dimensions: form.dimensions,
          shape: form.shape,
          technical: form.technical,
          recommended: form.recommended
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
}

// 删除设备
const handleDelete = () => {
  if (!props.deviceId) return
  
  ElMessageBox.confirm(
    '确定要删除此设备吗？删除后无法恢复！',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  )
  .then(async () => {
    const success = await removeDevice(props.deviceId!)
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
.device-form {
  margin-bottom: 20px;
}

.device-form :deep(.el-divider__text) {
  font-weight: bold;
  color: #409EFF;
}
</style>