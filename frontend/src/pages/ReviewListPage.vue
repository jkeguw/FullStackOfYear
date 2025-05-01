<template>
  <div class="review-list-page">
    <div class="container mx-auto p-4">
      <div class="flex justify-between items-center mb-4">
        <h1 class="text-2xl font-bold text-white">Review List</h1>
        <router-link v-if="isAdmin" to="/reviews/create" class="el-button el-button--primary">
          <el-icon class="mr-1"><Plus /></el-icon> Create Review
        </router-link>
      </div>

      <el-card class="mb-4 bg-[#1A1A1A] border border-[#333333] text-white">
        <div class="flex flex-wrap gap-4 justify-between items-end">
          <div class="flex-1 min-w-[200px]">
            <el-input
              v-model="searchQuery"
              placeholder="Search reviews"
              clearable
              @clear="handleSearch"
              @keyup.enter="handleSearch"
              class="dark-theme-input"
            >
              <template #suffix>
                <el-icon @click="handleSearch"><Search /></el-icon>
              </template>
            </el-input>
          </div>

          <div class="flex flex-wrap gap-2">
            <el-select
              v-model="filters.type"
              placeholder="Review Type"
              clearable
              @change="handleFilterChange"
              class="dark-theme-select"
            >
              <el-option label="All Types" value="" />
              <el-option label="Mouse" value="mouse" />
              <el-option label="Keyboard" value="keyboard" />
              <el-option label="Monitor" value="monitor" />
              <el-option label="Mousepad" value="mousepad" />
              <el-option label="Accessory" value="accessory" />
            </el-select>

            <el-select
              v-model="filters.contentType"
              placeholder="Content Type"
              clearable
              @change="handleFilterChange"
              class="dark-theme-select"
            >
              <el-option label="All Content" value="" />
              <el-option label="Single Product Review" value="single" />
              <el-option label="Comparison Review" value="comparison" />
              <el-option label="Usage Experience" value="experience" />
              <el-option label="Gaming Experience" value="gaming" />
              <el-option label="Buying Guide" value="buying" />
            </el-select>

            <el-select v-model="sortBy" placeholder="Sort By" @change="handleSortChange" class="dark-theme-select">
              <el-option label="Newest First" value="createdAt:desc" />
              <el-option label="Oldest First" value="createdAt:asc" />
              <el-option label="Highest Rating" value="score:desc" />
              <el-option label="Lowest Rating" value="score:asc" />
              <el-option label="Most Viewed" value="viewCount:desc" />
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

      <div v-else-if="reviews.length === 0" class="text-center py-8 text-white">
        <el-empty description="No reviews available" />
        <router-link v-if="isAdmin" to="/reviews/create" class="el-button el-button--primary mt-4">
          Create First Review
        </router-link>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <el-card
          v-for="review in reviews"
          :key="review.id"
          class="review-card mb-4 hover:shadow-lg transition-shadow duration-300 bg-[#1A1A1A] border border-[#333333] text-white"
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
              <h3 class="text-lg font-bold truncate mr-2 text-white">{{ review.title || 'Unnamed Review' }}</h3>
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

            <p class="text-gray-400 mb-4 line-clamp-2">
              {{ review.content ? review.content.substring(0, 100) + '...' : 'No content available' }}
            </p>

            <div class="flex justify-between items-center text-sm text-gray-500">
              <span>
                <el-icon><View /></el-icon> {{ review.viewCount || 0 }}
              </span>
              <span>{{ formatDate(review.createdAt) }}</span>
            </div>

            <div class="mt-4 flex gap-2">
              <el-button type="primary" @click="viewReview(review.id)">View Details</el-button>
              <el-button v-if="isAdmin" plain @click="editReview(review.id)">Edit</el-button>
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

// State
const reviews = ref<Review[]>([]);
const totalCount = ref(0);
const totalPages = ref(0);
const currentPage = ref(1);
const pageSize = ref(12);
const searchQuery = ref('');
const sortBy = ref('createdAt:desc');
// Get admin privileges from user authentication state
const { isAdmin } = useAuth();
const filters = reactive({
  type: '',
  contentType: ''
});

// Lifecycle hooks
onMounted(() => {
  loadReviews();
});

// Load review data
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
    ElMessage.error('Failed to load review list');
    console.error(error);
  }
};

// Handle search
const handleSearch = () => {
  currentPage.value = 1;
  loadReviews();
};

// Handle filtering
const handleFilterChange = () => {
  currentPage.value = 1;
  loadReviews();
};

// Handle sorting
const handleSortChange = () => {
  loadReviews();
};

// Handle page change
const handleCurrentChange = (page: number) => {
  currentPage.value = page;
  loadReviews();
};

// Handle page size change
const handleSizeChange = (size: number) => {
  pageSize.value = size;
  currentPage.value = 1;
  loadReviews();
};

// View review details
const viewReview = (id: string) => {
  router.push(`/reviews/${id}`);
};

// Edit review
const editReview = (id: string) => {
  router.push(`/reviews/${id}/edit`);
};

// Get review type name
const getTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    mouse: 'Mouse',
    keyboard: 'Keyboard',
    monitor: 'Monitor',
    mousepad: 'Mousepad',
    accessory: 'Accessory'
  };
  return typeMap[type] || 'Unknown Type';
};

// Get content type name
const getContentTypeName = (contentType: string) => {
  const contentTypeMap: Record<string, string> = {
    single: 'Single Review',
    comparison: 'Comparison',
    experience: 'Experience',
    gaming: 'Gaming Review',
    buying: 'Buying Guide'
  };
  return contentTypeMap[contentType] || 'Unknown Type';
};

// Get tag type
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

// Get content tag type
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
  background-color: #121212;
  min-height: 100vh;
}

.review-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: transform 0.2s ease-in-out;
}

.review-card:hover {
  transform: translateY(-5px);
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

:deep(.el-pagination) {
  --el-pagination-bg-color: #1A1A1A;
  --el-pagination-text-color: white;
  --el-pagination-button-color: white;
  --el-pagination-hover-color: #409EFF;
}

:deep(.el-select) .el-input__inner {
  color: white !important;
  background-color: #242424 !important;
  border-color: #333333 !important;
}

:deep(.el-input__inner) {
  color: white !important;
  background-color: #242424 !important;
  border-color: #333333 !important;
}

:deep(.el-button--primary) {
  background-color: #409EFF;
}

:deep(.el-button) {
  color: white;
  border-color: #333333;
}

:deep(.el-button:not(.el-button--primary)) {
  background-color: #242424;
}

:deep(.el-select-dropdown) {
  background-color: #242424;
  border-color: #333333;
}

:deep(.el-select-dropdown__item) {
  color: white;
}

:deep(.el-select-dropdown__item.hover) {
  background-color: #333333;
}
</style>
