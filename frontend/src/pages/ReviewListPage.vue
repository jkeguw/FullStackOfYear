<template>
  <div class="review-list-page">
    <div class="container mx-auto p-4">
      <div class="flex justify-between items-center mb-4">
        <h1 class="text-2xl font-bold">评测列表</h1>
        <router-link v-if="isAdmin" to="/reviews/create" class="el-button el-button--primary">
          <el-icon class="mr-1"><Plus /></el-icon> 创建评测
        </router-link>
      </div>

      <el-card class="mb-4">
        <div class="flex flex-wrap gap-4 justify-between items-end">
          <div class="flex-1 min-w-[200px]">
            <el-input
              v-model="searchQuery"
              placeholder="搜索评测"
              clearable
              @clear="handleSearch"
              @keyup.enter="handleSearch"
            >
              <template #suffix>
                <el-icon @click="handleSearch"><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <div class="flex flex-wrap gap-2">
            <el-select
              v-model="filters.type"
              placeholder="评测类型"
              clearable
              @change="handleFilterChange"
            >
              <el-option label="全部类型" value="" />
              <el-option label="鼠标" value="mouse" />
              <el-option label="键盘" value="keyboard" />
              <el-option label="显示器" value="monitor" />
              <el-option label="鼠标垫" value="mousepad" />
              <el-option label="配件" value="accessory" />
            </el-select>

            <el-select
              v-model="filters.contentType"
              placeholder="内容类型"
              clearable
              @change="handleFilterChange"
            >
              <el-option label="全部内容" value="" />
              <el-option label="单品评测" value="single" />
              <el-option label="对比评测" value="comparison" />
              <el-option label="使用心得" value="experience" />
              <el-option label="游戏体验" value="gaming" />
              <el-option label="选购建议" value="buying" />
            </el-select>

            <el-select v-model="sortBy" placeholder="排序" @change="handleSortChange">
              <el-option label="最新发布" value="createdAt:desc" />
              <el-option label="最早发布" value="createdAt:asc" />
              <el-option label="评分最高" value="score:desc" />
              <el-option label="评分最低" value="score:asc" />
              <el-option label="最多浏览" value="viewCount:desc" />
            </el-select>
          </div>
        </div>
      </el-card>

      <el-skeleton v-if="loading" :loading="loading" animated :count="3">
        <template #template>
          <div class="mb-4">
            <el-skeleton-item variant="p" style="width: 100%; height: 200px" />
          </div>
        </template>
      </el-skeleton>

      <div v-else-if="reviews.length === 0" class="text-center py-8">
        <el-empty description="暂无评测数据" />
        <router-link v-if="isAdmin" to="/reviews/create" class="el-button el-button--primary mt-4">
          创建第一篇评测
        </router-link>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <el-card
          v-for="review in reviews"
          :key="review.id"
          class="review-card mb-4 hover:shadow-lg transition-shadow duration-300"
          shadow="hover"
          :body-style="{ padding: '0' }"
        >
          <div class="relative">
            <div class="review-image bg-gray-200" style="height: 200px">
              <img
                v-if="review.imageUrl"
                :src="review.imageUrl"
                alt="Review image"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                <el-icon size="64"><Picture /></el-icon>
              </div>
            </div>

            <div class="absolute top-2 right-2 flex gap-2">
              <el-tag size="small" :type="getTagType(review.type)">
                {{ getTypeName(review.type) }}
              </el-tag>
              <el-tag size="small" :type="getContentTagType(review.contentType)">
                {{ getContentTypeName(review.contentType) }}
              </el-tag>
            </div>
          </div>

          <div class="p-4">
            <div class="flex justify-between items-start mb-2">
              <h3 class="text-lg font-bold truncate mr-2">{{ review.title || '未命名评测' }}</h3>
              <div class="flex items-center">
                <el-rate
                  :model-value="review.score || 0"
                  disabled
                  text-color="#ff9900"
                  :show-score="false"
                  size="small"
                />
                <span class="text-sm ml-1 font-medium">{{ (review.score || 0).toFixed(1) }}</span>
              </div>
            </div>

            <p class="text-gray-600 mb-4 line-clamp-2">
              {{ review.content ? review.content.substring(0, 100) + '...' : '暂无内容' }}
            </p>

            <div class="flex justify-between items-center text-sm text-gray-500">
              <span>
                <el-icon><View /></el-icon> {{ review.viewCount || 0 }}
              </span>
              <span>{{ formatDate(review.createdAt) }}</span>
            </div>

            <div class="mt-4 flex gap-2">
              <el-button type="primary" @click="viewReview(review.id)"> 查看详情 </el-button>
              <el-button v-if="isAdmin" plain @click="editReview(review.id)"> 编辑 </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <div v-if="reviews.length > 0 && totalPages > 1" class="flex justify-center mt-6">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalCount"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { Plus, Search, View, Picture } from '@element-plus/icons-vue';
import { useReview } from '@/composables/useReview';
import { useAuth } from '@/composables/useAuth';
import { formatDate } from '@/utils/date';

interface Review {
  id: string;
  title?: string;
  content: string;
  score: number;
  type: string;
  contentType: string;
  imageUrl?: string;
  viewCount: number;
  createdAt: string;
  updatedAt: string;
}

const router = useRouter();
const { fetchReviews, loading } = useReview();

// 状态
const reviews = ref<Review[]>([]);
const totalCount = ref(0);
const totalPages = ref(0);
const currentPage = ref(1);
const pageSize = ref(12);
const searchQuery = ref('');
const sortBy = ref('createdAt:desc');
// 从用户认证状态获取管理员权限
const { isAdmin } = useAuth();
const filters = reactive({
  type: '',
  contentType: ''
});

// 生命周期钩子
onMounted(() => {
  loadReviews();
});

// 加载评测数据
const loadReviews = async () => {
  try {
    const [field, order] = sortBy.value.split(':');
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      sort: field,
      order: order,
      search: searchQuery.value,
      type: filters.type,
      contentType: filters.contentType
    };

    const response = await fetchReviews(params);
    reviews.value = response.data || [];
    totalCount.value = response.total || 0;
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } catch (error) {
    ElMessage.error('加载评测列表失败');
    console.error(error);
  }
};

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1;
  loadReviews();
};

// 处理筛选
const handleFilterChange = () => {
  currentPage.value = 1;
  loadReviews();
};

// 处理排序
const handleSortChange = () => {
  loadReviews();
};

// 处理页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page;
  loadReviews();
};

// 处理每页数量变化
const handleSizeChange = (size: number) => {
  pageSize.value = size;
  currentPage.value = 1;
  loadReviews();
};

// 查看评测详情
const viewReview = (id: string) => {
  router.push(`/reviews/${id}`);
};

// 编辑评测
const editReview = (id: string) => {
  router.push(`/reviews/${id}/edit`);
};

// 获取评测类型名称
const getTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    mouse: '鼠标',
    keyboard: '键盘',
    monitor: '显示器',
    mousepad: '鼠标垫',
    accessory: '配件'
  };
  return typeMap[type] || '未知类型';
};

// 获取内容类型名称
const getContentTypeName = (contentType: string) => {
  const contentTypeMap: Record<string, string> = {
    single: '单品评测',
    comparison: '对比评测',
    experience: '使用心得',
    gaming: '游戏体验',
    buying: '选购建议'
  };
  return contentTypeMap[contentType] || '未知类型';
};

// 获取标签类型
const getTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    mouse: 'primary',
    keyboard: 'success',
    monitor: 'warning',
    mousepad: 'info',
    accessory: 'danger'
  };
  return typeMap[type] || '';
};

// 获取内容标签类型
const getContentTagType = (contentType: string) => {
  const contentTypeMap: Record<string, string> = {
    single: '',
    comparison: 'success',
    experience: 'info',
    gaming: 'warning',
    buying: 'danger'
  };
  return contentTypeMap[contentType] || '';
};
</script>

<style scoped>
.review-list-page {
  padding-bottom: 40px;
}

.review-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.review-card .el-card__body {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
