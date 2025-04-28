<template>
  <div class="device-detail-page">
    <div class="container mx-auto p-4">
      <el-page-header
        :title="'返回设备列表'"
        :content="device?.name || '设备详情'"
        @back="router.back()"
      />

      <div v-if="loading" class="mt-6">
        <el-skeleton :rows="10" animated />
      </div>

      <div v-else-if="!device" class="mt-6 text-center">
        <el-empty description="未找到设备信息" />
        <el-button class="mt-4" @click="router.push('/devices')">返回设备列表</el-button>
      </div>

      <template v-else>
        <!-- 设备详情 -->
        <div class="mt-6 grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- 设备图片 -->
          <div class="md:col-span-1">
            <el-card class="h-full">
              <div class="device-image h-80 flex items-center justify-center">
                <el-image
                  v-if="device.imageUrl"
                  :src="device.imageUrl"
                  fit="contain"
                  class="max-h-full"
                  :preview-src-list="[device.imageUrl]"
                ></el-image>
                <div v-else class="text-center text-gray-400">
                  <i class="el-icon-picture-outline text-5xl"></i>
                  <div>无图片</div>
                </div>
              </div>

              <div v-if="hasAdminPermission" class="mt-4 flex justify-center">
                <el-button type="primary" @click="handleEditDevice">编辑设备</el-button>
                <el-button type="danger" @click="handleDeleteDevice">删除设备</el-button>
              </div>
            </el-card>
          </div>

          <!-- 设备基本信息 -->
          <div class="md:col-span-2">
            <el-card>
              <div class="flex justify-between items-center">
                <div>
                  <h1 class="text-2xl font-bold">{{ device.brand }} {{ device.name }}</h1>
                  <div class="mt-1 flex items-center">
                    <el-tag>{{ getDeviceTypeName(device.type) }}</el-tag>
                    <span class="text-gray-500 ml-2"
                      >添加时间: {{ formatDate(device.createdAt) }}</span
                    >
                  </div>
                </div>
              </div>

              <div v-if="device.description" class="mt-4 text-gray-600">
                {{ device.description }}
              </div>

              <el-divider />

              <!-- 尺寸信息 -->
              <div class="mb-6">
                <h3 class="text-lg font-medium mb-2">尺寸信息</h3>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
                  <div class="stat-box">
                    <div class="text-sm text-gray-500">长度</div>
                    <div class="text-xl font-bold">{{ device.dimensions.length }} mm</div>
                  </div>
                  <div class="stat-box">
                    <div class="text-sm text-gray-500">宽度</div>
                    <div class="text-xl font-bold">{{ device.dimensions.width }} mm</div>
                  </div>
                  <div class="stat-box">
                    <div class="text-sm text-gray-500">高度</div>
                    <div class="text-xl font-bold">{{ device.dimensions.height }} mm</div>
                  </div>
                  <div class="stat-box">
                    <div class="text-sm text-gray-500">重量</div>
                    <div class="text-xl font-bold">{{ device.dimensions.weight }} g</div>
                  </div>
                </div>
              </div>

              <!-- 形状信息 -->
              <div class="mb-6">
                <h3 class="text-lg font-medium mb-2">形状信息</h3>
                <div class="grid grid-cols-2 gap-4">
                  <div class="shape-info-item">
                    <div class="font-medium">形状类型:</div>
                    <div>{{ getShapeTypeName(device.shape.type) }}</div>
                  </div>
                  <div class="shape-info-item">
                    <div class="font-medium">坑位位置:</div>
                    <div>{{ getHumpPlacementName(device.shape.humpPlacement) }}</div>
                  </div>
                  <div class="shape-info-item">
                    <div class="font-medium">前端开叉:</div>
                    <div>{{ getFrontFlareName(device.shape.frontFlare) }}</div>
                  </div>
                  <div class="shape-info-item">
                    <div class="font-medium">侧面曲线:</div>
                    <div>{{ getSideCurvatureName(device.shape.sideCurvature) }}</div>
                  </div>
                  <div class="shape-info-item">
                    <div class="font-medium">手型适配:</div>
                    <div>{{ getHandCompatibilityName(device.shape.handCompatibility) }}</div>
                  </div>
                </div>
              </div>

              <!-- 技术参数 -->
              <div class="mb-6">
                <h3 class="text-lg font-medium mb-2">技术参数</h3>
                <div class="grid grid-cols-2 gap-4">
                  <div class="tech-info-item">
                    <div class="font-medium">连接方式:</div>
                    <div class="flex flex-wrap gap-1">
                      <el-tag
                        v-for="conn in device.technical.connectivity"
                        :key="conn"
                        size="small"
                      >
                        {{ getConnectivityName(conn) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="tech-info-item">
                    <div class="font-medium">传感器:</div>
                    <div>{{ device.technical.sensor }}</div>
                  </div>
                  <div class="tech-info-item">
                    <div class="font-medium">最大DPI:</div>
                    <div>{{ device.technical.maxDPI.toLocaleString() }}</div>
                  </div>
                  <div class="tech-info-item">
                    <div class="font-medium">轮询率:</div>
                    <div>{{ device.technical.pollingRate }} Hz</div>
                  </div>
                  <div class="tech-info-item">
                    <div class="font-medium">侧键数量:</div>
                    <div>{{ device.technical.sideButtons }}</div>
                  </div>
                </div>

                <!-- 电池信息 (如果有) -->
                <template v-if="device.technical.battery">
                  <h4 class="text-md font-medium mt-4 mb-2">电池信息</h4>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="tech-info-item">
                      <div class="font-medium">电池类型:</div>
                      <div>{{ getBatteryTypeName(device.technical.battery.type) }}</div>
                    </div>
                    <div class="tech-info-item">
                      <div class="font-medium">电池容量:</div>
                      <div>{{ device.technical.battery.capacity }} mAh</div>
                    </div>
                    <div class="tech-info-item">
                      <div class="font-medium">电池寿命:</div>
                      <div>{{ device.technical.battery.life }} 小时</div>
                    </div>
                  </div>
                </template>
              </div>

              <!-- 推荐信息 -->
              <div>
                <h3 class="text-lg font-medium mb-2">推荐信息</h3>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="recommend-info-item">
                    <div class="font-medium">适合游戏类型:</div>
                    <div class="flex flex-wrap gap-1 mt-1">
                      <el-tag
                        v-for="game in device.recommended.gameTypes"
                        :key="game"
                        size="small"
                        type="success"
                      >
                        {{ getGameTypeName(game) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="recommend-info-item">
                    <div class="font-medium">适合握持方式:</div>
                    <div class="flex flex-wrap gap-1 mt-1">
                      <el-tag
                        v-for="grip in device.recommended.gripStyles"
                        :key="grip"
                        size="small"
                        type="warning"
                      >
                        {{ getGripStyleName(grip) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="recommend-info-item">
                    <div class="font-medium">适合手型大小:</div>
                    <div class="flex flex-wrap gap-1 mt-1">
                      <el-tag
                        v-for="hand in device.recommended.handSizes"
                        :key="hand"
                        size="small"
                        type="info"
                      >
                        {{ getHandSizeName(hand) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="recommend-info-item">
                    <div class="font-medium">其他特性:</div>
                    <div class="mt-1">
                      <el-tag
                        v-if="device.recommended.dailyUse"
                        size="small"
                        effect="plain"
                        class="mr-1"
                      >
                        适合日常使用
                      </el-tag>
                      <el-tag
                        v-if="device.recommended.professional"
                        size="small"
                        type="danger"
                        effect="plain"
                      >
                        专业级设备
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </div>

        <!-- 评测部分 -->
        <div class="mt-6">
          <el-card>
            <template #header>
              <div class="flex justify-between items-center">
                <h2 class="text-xl font-bold">用户评测</h2>
                <el-button type="primary" @click="handleAddReview">撰写评测</el-button>
              </div>
            </template>

            <!-- 评测列表 -->
            <div v-if="reviews.length === 0 && !reviewsLoading" class="text-center py-8">
              <el-empty description="暂无评测" />
              <el-button class="mt-4" type="primary" @click="handleAddReview"
                >成为第一个评测者</el-button
              >
            </div>

            <div v-else>
              <div v-loading="reviewsLoading">
                <div
                  v-for="review in reviews"
                  :key="review.id"
                  class="review-item mb-6 pb-6 border-b border-gray-200 last:border-0"
                >
                  <div class="flex justify-between">
                    <div class="review-header flex items-center">
                      <div class="reviewer-avatar mr-2">
                        <el-avatar :size="40">{{ review.userId.substring(0, 2) }}</el-avatar>
                      </div>
                      <div>
                        <div class="font-medium">用户 {{ review.userId.substring(0, 8) }}</div>
                        <div class="text-xs text-gray-500">{{ formatDate(review.createdAt) }}</div>
                      </div>
                    </div>

                    <div class="review-score">
                      <el-rate
                        v-model="review.score"
                        disabled
                        :colors="rateColors"
                        show-score
                      ></el-rate>
                    </div>
                  </div>

                  <div class="review-usage mt-2">
                    <el-tag size="small" effect="plain">{{ getUsageName(review.usage) }}</el-tag>
                  </div>

                  <div class="review-content mt-3">
                    {{ review.content }}
                  </div>

                  <div class="review-tags mt-4 flex flex-wrap gap-2">
                    <div>
                      <div class="text-sm text-gray-500 mb-1">优点:</div>
                      <div class="flex flex-wrap gap-1">
                        <el-tag
                          v-for="pro in review.pros"
                          :key="pro"
                          size="small"
                          type="success"
                          effect="light"
                        >
                          {{ pro }}
                        </el-tag>
                      </div>
                    </div>

                    <div class="ml-6">
                      <div class="text-sm text-gray-500 mb-1">缺点:</div>
                      <div class="flex flex-wrap gap-1">
                        <el-tag
                          v-for="con in review.cons"
                          :key="con"
                          size="small"
                          type="danger"
                          effect="light"
                        >
                          {{ con }}
                        </el-tag>
                      </div>
                    </div>
                  </div>

                  <!-- 评测操作 -->
                  <div v-if="canEditReview(review)" class="review-actions mt-4 flex justify-end">
                    <el-button size="small" @click="handleEditReview(review.id)">编辑</el-button>
                    <el-button size="small" type="danger" @click="handleDeleteReview(review.id)"
                      >删除</el-button
                    >
                  </div>
                </div>

                <!-- 评测分页 -->
                <div v-if="reviews.length > 0" class="review-pagination mt-4 flex justify-center">
                  <el-pagination
                    v-model:current-page="reviewPagination.page"
                    v-model:page-size="reviewPagination.pageSize"
                    :page-sizes="[5, 10, 20]"
                    layout="total, sizes, prev, pager, next"
                    :total="reviewPagination.total"
                    @size-change="handleReviewSizeChange"
                    @current-change="handleReviewPageChange"
                  ></el-pagination>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </template>
    </div>

    <!-- 创建/编辑设备对话框 -->
    <el-dialog
      v-model="deviceDialogVisible"
      title="编辑设备"
      width="80%"
      :before-close="closeDeviceDialog"
    >
      <DeviceForm :device-id="deviceId" @saved="handleDeviceSaved" @deleted="handleDeviceDeleted" />
    </el-dialog>

    <!-- 创建/编辑评测对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      :title="currentReviewId ? '编辑评测' : '撰写评测'"
      width="80%"
      :before-close="closeReviewDialog"
    >
      <ReviewForm
        :review-id="currentReviewId"
        @saved="handleReviewSaved"
        @deleted="handleReviewDeleted"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { format } from 'date-fns';
import DeviceForm from '@/components/form/DeviceForm.vue';
import ReviewForm from '@/components/form/ReviewForm.vue';
import { useDevice } from '@/composables/useDevice';
import type { DeviceReview, MouseDevice } from '@/api/device';

// 路由
const route = useRoute();
const router = useRouter();

// 设备ID
const deviceId = computed(() => route.params.id as string);

// 设备钩子
const {
  currentDevice: device,
  deviceLoading: loading,
  reviews,
  reviewLoading: reviewsLoading,
  reviewPagination,
  getDeviceTypeName,
  getHandSizeName,
  getGripStyleName,
  fetchMouseDevice,
  fetchDeviceReviews,
  removeDevice,
  removeDeviceReview
} = useDevice();

// 对话框控制
const deviceDialogVisible = ref(false);
const reviewDialogVisible = ref(false);
const currentReviewId = ref<string | undefined>(undefined);

// 权限检查（实际应该从用户状态获取）
const hasAdminPermission = ref(true);
const currentUserId = ref('user123'); // 模拟当前用户ID

// 评分颜色
const rateColors = ['#F56C6C', '#E6A23C', '#909399', '#67C23A', '#409EFF'];

// 生命周期钩子
onMounted(() => {
  loadData();
});

// 监听路由变化
watch(
  () => route.params.id,
  (newId) => {
    if (newId) {
      loadData();
    }
  }
);

// 加载数据
const loadData = async () => {
  if (!deviceId.value) return;

  await fetchMouseDevice(deviceId.value);
  loadReviews();
};

// 加载评测数据
const loadReviews = async () => {
  if (!deviceId.value) return;

  await fetchDeviceReviews({
    deviceId: deviceId.value,
    page: reviewPagination.page,
    pageSize: reviewPagination.pageSize,
    sortBy: 'createdAt',
    sortOrder: 'desc'
  });
};

// 评测分页处理
const handleReviewSizeChange = (size: number) => {
  reviewPagination.pageSize = size;
  loadReviews();
};

const handleReviewPageChange = (page: number) => {
  reviewPagination.page = page;
  loadReviews();
};

// 编辑设备
const handleEditDevice = () => {
  deviceDialogVisible.value = true;
};

// 关闭设备对话框
const closeDeviceDialog = () => {
  deviceDialogVisible.value = false;
};

// 设备保存成功处理
const handleDeviceSaved = (device: MouseDevice) => {
  deviceDialogVisible.value = false;
  ElMessage.success(`设备 ${device.name} 已更新`);
  loadData(); // 刷新数据
};

// 设备删除成功处理
const handleDeviceDeleted = () => {
  deviceDialogVisible.value = false;
  ElMessage.success('设备已删除');
  router.push('/devices'); // 返回设备列表
};

// 删除设备
const handleDeleteDevice = () => {
  ElMessageBox.confirm('确定要删除此设备吗？删除后无法恢复！', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      if (!deviceId.value) return;

      const success = await removeDevice(deviceId.value);
      if (success) {
        ElMessage.success('设备已删除');
        router.push('/devices'); // 返回设备列表
      }
    })
    .catch(() => {
      // 取消删除，不做任何操作
    });
};

// 添加评测
const handleAddReview = () => {
  currentReviewId.value = undefined;
  reviewDialogVisible.value = true;
};

// 编辑评测
const handleEditReview = (id: string) => {
  currentReviewId.value = id;
  reviewDialogVisible.value = true;
};

// 关闭评测对话框
const closeReviewDialog = () => {
  reviewDialogVisible.value = false;
  currentReviewId.value = undefined;
};

// 评测保存成功处理
const handleReviewSaved = (review: DeviceReview) => {
  reviewDialogVisible.value = false;
  ElMessage.success(currentReviewId.value ? '评测已更新' : '评测已提交，等待审核');
  loadReviews(); // 刷新评测列表
};

// 评测删除成功处理
const handleReviewDeleted = () => {
  reviewDialogVisible.value = false;
  ElMessage.success('评测已删除');
  loadReviews(); // 刷新评测列表
};

// 删除评测
const handleDeleteReview = (id: string) => {
  ElMessageBox.confirm('确定要删除此评测吗？删除后无法恢复！', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      const success = await removeDeviceReview(id);
      if (success) {
        ElMessage.success('评测已删除');
        loadReviews(); // 刷新评测列表
      }
    })
    .catch(() => {
      // 取消删除，不做任何操作
    });
};

// 检查是否可以编辑评测
const canEditReview = (review: DeviceReview) => {
  return review.userId === currentUserId.value || hasAdminPermission.value;
};

// 格式化日期
const formatDate = (date: string | Date) => {
  return format(new Date(date), 'yyyy-MM-dd');
};

// 辅助方法
const getShapeTypeName = (type: string): string => {
  const types: Record<string, string> = {
    ergonomic: '人体工学',
    ambidextrous: '左右对称'
  };
  return types[type] || type;
};

const getHumpPlacementName = (placement: string): string => {
  const placements: Record<string, string> = {
    front: '前段',
    center: '中段',
    back: '后段'
  };
  return placements[placement] || placement;
};

const getFrontFlareName = (flare: string): string => {
  const flares: Record<string, string> = {
    narrow: '窄',
    medium: '中等',
    wide: '宽'
  };
  return flares[flare] || flare;
};

const getSideCurvatureName = (curvature: string): string => {
  const curvatures: Record<string, string> = {
    straight: '直线',
    curved: '曲线'
  };
  return curvatures[curvature] || curvature;
};

const getHandCompatibilityName = (compatibility: string): string => {
  const compatibilities: Record<string, string> = {
    small: '小型手',
    medium: '中型手',
    large: '大型手',
    universal: '广泛适配'
  };
  return compatibilities[compatibility] || compatibility;
};

const getConnectivityName = (connectivity: string): string => {
  const connectivities: Record<string, string> = {
    wired: '有线',
    wireless: '2.4G无线',
    bluetooth: '蓝牙'
  };
  return connectivities[connectivity] || connectivity;
};

const getBatteryTypeName = (type: string): string => {
  const types: Record<string, string> = {
    lithium: '内置锂电池',
    replaceable: '可更换电池'
  };
  return types[type] || type;
};

const getGameTypeName = (type: string): string => {
  const types: Record<string, string> = {
    fps: 'FPS射击游戏',
    moba: 'MOBA游戏',
    rts: 'RTS策略游戏',
    mmo: 'MMO角色扮演',
    racing: '竞速游戏',
    fighting: '格斗游戏',
    casual: '休闲游戏'
  };
  return types[type] || type;
};

const getUsageName = (usage: string): string => {
  const usages: Record<string, string> = {
    fps_gaming: 'FPS 游戏',
    moba_gaming: 'MOBA 游戏',
    mmo_gaming: 'MMO 游戏',
    general_gaming: '综合游戏',
    office_work: '办公工作',
    creative_design: '创意设计',
    daily_use: '日常使用'
  };
  return usages[usage] || usage;
};
</script>

<style scoped>
.stat-box {
  padding: 10px;
  border-radius: 4px;
  background-color: #f5f7fa;
  transition: all 0.3s ease;
}

.stat-box:hover {
  background-color: #ecf5ff;
  transform: translateY(-2px);
}

.shape-info-item,
.tech-info-item,
.recommend-info-item {
  margin-bottom: 8px;
}

.device-image {
  overflow: hidden;
  background-color: #f5f5f5;
  border-radius: 4px;
}

/* 评测样式 */
.review-item {
  transition: all 0.3s ease;
}

.review-item:hover {
  background-color: #f9f9f9;
  padding-left: 8px;
  padding-right: 8px;
  margin-left: -8px;
  margin-right: -8px;
}
</style>
