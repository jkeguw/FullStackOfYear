<template>
  <div class="review-form-page">
    <div class="container mx-auto p-4">
      <el-card class="w-full max-w-4xl mx-auto">
        <template #header>
          <div class="flex justify-between items-center">
            <h1 class="text-xl font-bold">
              {{ isEditMode ? '编辑评测' : '创建新评测' }}
            </h1>
            <div>
              <el-button @click="$router.go(-1)"> 取消 </el-button>
              <el-button type="primary" :loading="loading" @click="submitReview">
                {{ isEditMode ? '保存修改' : '提交评测' }}
              </el-button>
            </div>
          </div>
        </template>

        <el-form
          ref="formRef"
          v-loading="loading"
          :model="form"
          :rules="rules"
          label-position="top"
        >
          <el-form-item label="评测标题" prop="title">
            <el-input
              v-model="form.title"
              placeholder="请输入评测标题"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <el-form-item label="评测类型" prop="type">
              <el-select v-model="form.type" placeholder="请选择评测类型" class="w-full">
                <el-option label="鼠标" value="mouse" />
                <el-option label="键盘" value="keyboard" />
                <el-option label="显示器" value="monitor" />
                <el-option label="鼠标垫" value="mousepad" />
                <el-option label="配件" value="accessory" />
              </el-select>
            </el-form-item>

            <el-form-item label="内容类型" prop="contentType">
              <el-select v-model="form.contentType" placeholder="请选择内容类型" class="w-full">
                <el-option label="单品评测" value="single" />
                <el-option label="对比评测" value="comparison" />
                <el-option label="使用心得" value="experience" />
                <el-option label="游戏体验" value="gaming" />
                <el-option label="选购建议" value="buying" />
              </el-select>
            </el-form-item>
          </div>

          <el-form-item label="选择设备" prop="deviceId">
            <el-select
              v-model="form.deviceId"
              placeholder="请选择要评测的设备"
              class="w-full"
              filterable
              remote
              :remote-method="searchDevices"
              :loading="devicesLoading"
            >
              <el-option
                v-for="device in deviceOptions"
                :key="device.id"
                :label="device.name"
                :value="device.id"
              >
                <div class="flex items-center">
                  <el-tag size="small" class="mr-2">
                    {{ getDeviceTypeLabel(device.type) }}
                  </el-tag>
                  <span>{{ device.name }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="总体评分" prop="score">
            <div class="flex items-center">
              <el-rate
                v-model="form.score"
                :colors="['#99A9BF', '#F7BA2A', '#FF9900']"
                :show-text="true"
                :texts="['很差', '较差', '一般', '不错', '很好']"
              />
              <span class="ml-4 text-lg">{{ form.score.toFixed(1) }}</span>
            </div>
          </el-form-item>

          <el-form-item label="主要使用场景" prop="usage">
            <el-input
              v-model="form.usage"
              placeholder="请描述主要使用场景，如：FPS游戏、办公、设计等"
              maxlength="50"
              show-word-limit
            />
          </el-form-item>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- 优点 -->
            <el-form-item label="优点" prop="pros">
              <div class="flex items-center mb-2">
                <el-input v-model="newPro" placeholder="请输入优点" @keyup.enter="addPro" />
                <el-button class="ml-2" @click="addPro"> 添加 </el-button>
              </div>

              <el-tag
                v-for="(pro, index) in form.pros"
                :key="index"
                closable
                type="success"
                class="mr-2 mb-2"
                @close="form.pros.splice(index, 1)"
              >
                {{ pro }}
              </el-tag>

              <div v-if="form.pros.length === 0" class="text-gray-400 text-sm">
                暂无优点，请添加
              </div>
            </el-form-item>

            <!-- 缺点 -->
            <el-form-item label="缺点" prop="cons">
              <div class="flex items-center mb-2">
                <el-input v-model="newCon" placeholder="请输入缺点" @keyup.enter="addCon" />
                <el-button class="ml-2" @click="addCon"> 添加 </el-button>
              </div>

              <el-tag
                v-for="(con, index) in form.cons"
                :key="index"
                closable
                type="danger"
                class="mr-2 mb-2"
                @close="form.cons.splice(index, 1)"
              >
                {{ con }}
              </el-tag>

              <div v-if="form.cons.length === 0" class="text-gray-400 text-sm">
                暂无缺点，请添加
              </div>
            </el-form-item>
          </div>

          <el-form-item label="评测内容" prop="content">
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="10"
              placeholder="请详细描述您对该设备的评测内容..."
              maxlength="5000"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="上传图片（可选）">
            <el-upload
              action="/api/v1/uploads"
              list-type="picture-card"
              :headers="uploadHeaders"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
              :limit="5"
            >
              <el-icon><Plus /></el-icon>
              <template #file="{ file }">
                <div class="relative">
                  <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
                  <div class="absolute top-0 right-0">
                    <el-icon @click.stop="handleRemove(file)"><Close /></el-icon>
                  </div>
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Plus, Close } from '@element-plus/icons-vue';
import { useReview } from '@/composables/useReview';
import { useDevice } from '@/composables/useDevice';

interface Props {
  deviceId?: string;
  reviewId?: string;
}

const props = defineProps<Props>();
const router = useRouter();
const route = useRoute();

// 获取评测相关操作
const { createReview, updateReview, getReview, loading } = useReview();
// 获取设备相关操作
const { searchDevicesByName, devicesLoading } = useDevice();

// 表单引用
const formRef = ref();

// 设备选项
const deviceOptions = ref<Array<{ id: string; name: string; type: string }>>([]);

// 新的优点和缺点
const newPro = ref('');
const newCon = ref('');

// 表单数据
const form = reactive({
  title: '',
  type: 'mouse',
  contentType: 'single',
  deviceId: '',
  score: 3,
  usage: '',
  pros: [] as string[],
  cons: [] as string[],
  content: '',
  images: [] as string[]
});

// 验证规则
const rules = {
  title: [
    { required: true, message: '请输入评测标题', trigger: 'blur' },
    { min: 3, max: 100, message: '标题长度应在3到100个字符之间', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择评测类型', trigger: 'change' }],
  contentType: [{ required: true, message: '请选择内容类型', trigger: 'change' }],
  deviceId: [{ required: true, message: '请选择要评测的设备', trigger: 'change' }],
  score: [{ required: true, type: 'number', message: '请设置评分', trigger: 'change' }],
  usage: [{ required: true, message: '请描述主要使用场景', trigger: 'blur' }],
  content: [
    { required: true, message: '请填写评测内容', trigger: 'blur' },
    { min: 50, message: '评测内容至少50个字符', trigger: 'blur' }
  ]
};

// 是否为编辑模式
const isEditMode = computed(() => !!props.reviewId);

// 上传头部信息
const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${localStorage.getItem('token')}`
  };
});

// 生命周期钩子
onMounted(async () => {
  // 如果URL中有deviceId，则设置到表单中
  if (props.deviceId) {
    form.deviceId = props.deviceId;
    // 获取设备信息
    await loadDeviceDetails(props.deviceId);
  }

  // 如果是编辑模式，则加载评测详情
  if (isEditMode.value) {
    await loadReviewDetails();
  }
});

// 加载设备详情
const loadDeviceDetails = async (deviceId: string) => {
  try {
    // 此处应调用API获取设备详情，用于展示设备名称
    // 在实际应用中，应该从API获取
    // 模拟数据
    deviceOptions.value = [
      {
        id: deviceId,
        name: '加载中...',
        type: 'unknown'
      }
    ];
  } catch (error) {
    console.error('加载设备详情失败', error);
  }
};

// 加载评测详情
const loadReviewDetails = async () => {
  if (!props.reviewId) return;

  try {
    const review = await getReview(props.reviewId);
    if (!review) {
      ElMessage.error('找不到指定的评测');
      router.push('/reviews');
      return;
    }

    // 填充表单数据
    form.title = review.title || '';
    form.type = review.type;
    form.contentType = review.contentType;
    form.deviceId = review.deviceId;
    form.score = review.score;
    form.usage = review.usage;
    form.pros = review.pros || [];
    form.cons = review.cons || [];
    form.content = review.content;
    form.images = review.images || [];

    // 加载设备详情
    await loadDeviceDetails(review.deviceId);
  } catch (error) {
    ElMessage.error('加载评测详情失败');
    console.error(error);
  }
};

// 搜索设备
const searchDevices = async (query: string) => {
  if (query) {
    try {
      const results = await searchDevicesByName(query);
      deviceOptions.value = results;
    } catch (error) {
      console.error('搜索设备失败', error);
    }
  } else {
    deviceOptions.value = [];
  }
};

// 添加优点
const addPro = () => {
  if (newPro.value.trim()) {
    form.pros.push(newPro.value.trim());
    newPro.value = '';
  }
};

// 添加缺点
const addCon = () => {
  if (newCon.value.trim()) {
    form.cons.push(newCon.value.trim());
    newCon.value = '';
  }
};

// 提交评测
const submitReview = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    const reviewData = {
      title: form.title,
      type: form.type,
      contentType: form.contentType,
      deviceId: form.deviceId,
      score: form.score,
      usage: form.usage,
      pros: form.pros,
      cons: form.cons,
      content: form.content,
      images: form.images
    };

    if (isEditMode.value && props.reviewId) {
      await updateReview(props.reviewId, reviewData);
      ElMessage.success('评测已更新');
    } else {
      await createReview(reviewData);
      ElMessage.success('评测已提交');
    }

    // 返回评测列表页
    router.push('/reviews');
  } catch (error: any) {
    if (error.message) {
      ElMessage.error(error.message);
    } else {
      ElMessage.error('提交失败，请检查表单内容');
    }
    console.error(error);
  }
};

// 上传前检查
const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/');
  const isLt5M = file.size / 1024 / 1024 < 5;

  if (!isImage) {
    ElMessage.error('只能上传图片文件!');
    return false;
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB!');
    return false;
  }
  return true;
};

// 上传成功回调
const handleUploadSuccess = (response: any, file: any) => {
  if (response && response.data && response.data.url) {
    form.images.push(response.data.url);
    ElMessage.success('图片上传成功');
  }
};

// 上传失败回调
const handleUploadError = (error: any) => {
  ElMessage.error('图片上传失败，请重试');
  console.error('上传失败', error);
};

// 移除图片
const handleRemove = (file: any) => {
  const index = form.images.findIndex((url) => url === file.url);
  if (index !== -1) {
    form.images.splice(index, 1);
  }
};

// 获取设备类型标签
const getDeviceTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    mouse: '鼠标',
    keyboard: '键盘',
    monitor: '显示器',
    mousepad: '鼠标垫',
    accessory: '配件',
    unknown: '未知'
  };
  return typeMap[type] || '未知';
};
</script>

<style scoped>
.review-form-page {
  padding-bottom: 40px;
}
</style>
