<template>
  <div class="review-detail-container">
    <div v-if="loading" class="loading-container">
      <el-spinner size="large" />
    </div>
    
    <div v-else-if="!review" class="empty-state">
      <el-empty description="找不到该评测" />
      <router-link to="/reviews">
        <el-button type="primary">返回评测列表</el-button>
      </router-link>
    </div>
    
    <template v-else>
      <div class="review-header">
        <div class="review-title-section">
          <h1 class="review-title">{{ review.title }}</h1>
          <div class="review-meta">
            <div class="author-info">
              <el-avatar :size="32" :src="review.author.avatar">
                {{ review.author.name.charAt(0) }}
              </el-avatar>
              <span class="author-name">{{ review.author.name }}</span>
            </div>
            <div class="published-date">
              发布于 {{ formatDate(review.publishedAt) }}
            </div>
            <div class="review-score" :class="getScoreClass(review.score)">
              {{ review.score }}/10
            </div>
          </div>
        </div>
        
        <div class="review-mouse-card">
          <el-card shadow="hover" class="mouse-info-card">
            <div class="mouse-card-content">
              <div class="mouse-image">
                <img 
                  v-if="review.mouse.imageUrl" 
                  :src="review.mouse.imageUrl" 
                  :alt="review.mouse.name" 
                  class="img-fluid"
                />
                <div v-else class="placeholder-image">
                  <i class="el-icon-mouse"></i>
                </div>
              </div>
              <div class="mouse-info">
                <h3 class="mouse-name">{{ review.mouse.brand }} {{ review.mouse.name }}</h3>
                <div class="mouse-specs">
                  <div class="spec-item">
                    <span class="spec-label">尺寸:</span>
                    <span class="spec-value">
                      {{ review.mouse.dimensions.length }}x{{ review.mouse.dimensions.width }}x{{ review.mouse.dimensions.height }}mm
                    </span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">重量:</span>
                    <span class="spec-value">{{ review.mouse.weight }}g</span>
                  </div>
                  <div class="spec-item">
                    <span class="spec-label">形状:</span>
                    <span class="spec-value">{{ review.mouse.shape.type }}</span>
                  </div>
                </div>
                <div class="mouse-actions">
                  <router-link :to="`/mice/${review.mouse.id}`">
                    <el-button size="small">查看详情</el-button>
                  </router-link>
                  <el-button 
                    size="small" 
                    type="primary" 
                    @click="addToComparison(review.mouse)"
                    :disabled="isInComparison(review.mouse.id)"
                  >
                    {{ isInComparison(review.mouse.id) ? '已添加比较' : '添加到比较' }}
                  </el-button>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      
      <div class="review-gallery">
        <review-gallery :images="review.images" />
      </div>
      
      <div class="review-content">
        <div class="content-section">
          <h2 class="section-title">概述</h2>
          <div class="section-content" v-html="review.summary"></div>
        </div>
        
        <div class="content-section">
          <h2 class="section-title">优点</h2>
          <ul class="pros-cons-list">
            <li v-for="(pro, index) in review.pros" :key="`pro-${index}`" class="pro-item">
              {{ pro }}
            </li>
          </ul>
        </div>
        
        <div class="content-section">
          <h2 class="section-title">缺点</h2>
          <ul class="pros-cons-list">
            <li v-for="(con, index) in review.cons" :key="`con-${index}`" class="con-item">
              {{ con }}
            </li>
          </ul>
        </div>
        
        <div class="content-section">
          <h2 class="section-title">详细评测</h2>
          <div class="section-content" v-html="review.content"></div>
        </div>
        
        <div class="content-section">
          <h2 class="section-title">规格参数</h2>
          <el-table :data="specsTableData" border stripe size="small">
            <el-table-column prop="property" label="参数" width="180" />
            <el-table-column prop="value" label="数值" />
          </el-table>
        </div>
        
        <div class="content-section">
          <h2 class="section-title">总结</h2>
          <div class="section-content" v-html="review.conclusion"></div>
          
          <div class="final-score-section">
            <div class="final-score" :class="getScoreClass(review.score)">
              {{ review.score }}/10
            </div>
            <div class="score-description">
              {{ getScoreDescription(review.score) }}
            </div>
          </div>
        </div>
      </div>
      
      <div class="review-actions">
        <el-button @click="$router.push('/reviews')">返回评测列表</el-button>
        <el-button type="primary" @click="addToComparison(review.mouse)">
          添加到比较
        </el-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useComparisonStore } from '@/stores';
import { getReview } from '@/api/review';
import { formatDate } from '@/utils/date';
import ReviewGallery from '@/components/review/ReviewGallery.vue';
import type { Review, MouseDevice } from '@/types/review';

// 路由
const route = useRoute();
const router = useRouter();
const comparisonStore = useComparisonStore();

// 状态
const review = ref<Review | null>(null);
const loading = ref(true);

// 计算属性
const specsTableData = computed(() => {
  if (!review.value) return [];
  
  const mouse = review.value.mouse;
  return [
    { property: '品牌', value: mouse.brand },
    { property: '型号', value: mouse.name },
    { property: '长度', value: `${mouse.dimensions.length}mm` },
    { property: '宽度', value: `${mouse.dimensions.width}mm` },
    { property: '高度', value: `${mouse.dimensions.height}mm` },
    { property: '重量', value: `${mouse.weight}g` },
    { property: '形状', value: mouse.shape.type },
    { property: '握法', value: mouse.shape.gripStyles.join(', ') },
    { property: '传感器', value: mouse.sensor },
    { property: 'DPI', value: `${mouse.dpi}` },
    { property: '按键开关', value: mouse.switches },
    { property: '连接方式', value: mouse.connectivity },
  ];
});

// 方法
function getScoreClass(score: number) {
  if (score >= 9) return 'score-excellent';
  if (score >= 7) return 'score-good';
  if (score >= 5) return 'score-average';
  return 'score-poor';
}

function getScoreDescription(score: number) {
  if (score >= 9) return '极力推荐';
  if (score >= 7) return '优秀';
  if (score >= 5) return '一般';
  if (score >= 3) return '较差';
  return '不推荐';
}

function addToComparison(mouse: MouseDevice) {
  if (comparisonStore.selectedMice.length >= 3) {
    window.ElMessage.warning('最多只能同时比较3个鼠标');
    return;
  }
  
  comparisonStore.addMouse(mouse);
  window.ElMessage.success(`已添加 ${mouse.brand} ${mouse.name} 到比较列表`);
}

function isInComparison(mouseId: string) {
  return comparisonStore.selectedMice.some(m => m.id === mouseId);
}

// 生命周期钩子
onMounted(async () => {
  const reviewId = route.params.id as string;
  
  if (!reviewId) {
    router.push('/reviews');
    return;
  }
  
  try {
    loading.value = true;
    const data = await getReview(reviewId);
    review.value = data.review as Review;
  } catch (error) {
    console.error('Error fetching review:', error);
    window.ElMessage.error('获取评测详情失败');
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.review-detail-container {
  padding: 1.5rem 1rem;
  max-width: 1200px;
  margin: 0 auto;
}

.loading-container, .empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  gap: 1.5rem;
}

.review-header {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

@media (min-width: 768px) {
  .review-header {
    grid-template-columns: 2fr 1fr;
  }
}

.review-title {
  font-size: 1.875rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.review-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1rem;
  color: var(--claude-gray-600);
}

.author-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.review-score {
  font-weight: 700;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
}

.mouse-card-content {
  display: flex;
  flex-direction: column;
}

@media (min-width: 640px) {
  .mouse-card-content {
    flex-direction: row;
    gap: 1rem;
  }
}

.mouse-image {
  width: 100%;
  height: 150px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f9f9f9;
  border-radius: 0.25rem;
  margin-bottom: 1rem;
}

@media (min-width: 640px) {
  .mouse-image {
    width: 150px;
    flex-shrink: 0;
    margin-bottom: 0;
  }
}

.placeholder-image {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #d0d0d0;
  font-size: 3rem;
}

.mouse-name {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.mouse-specs {
  margin-bottom: 1rem;
}

.spec-item {
  display: flex;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.spec-label {
  color: var(--claude-gray-600);
  width: 40px;
  margin-right: 0.5rem;
}

.mouse-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: auto;
}

.review-gallery {
  margin-bottom: 2rem;
}

.content-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--claude-gray-200);
}

.section-content {
  font-size: 1rem;
  line-height: 1.6;
}

.pros-cons-list {
  list-style-type: none;
  padding: 0;
}

.pros-cons-list li {
  position: relative;
  padding-left: 1.5rem;
  margin-bottom: 0.75rem;
}

.pros-cons-list li::before {
  position: absolute;
  left: 0;
  top: 0.2rem;
}

.pro-item::before {
  content: '✓';
  color: #67C23A;
}

.con-item::before {
  content: '✕';
  color: #F56C6C;
}

.final-score-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 2rem;
  padding: 1.5rem;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
}

.final-score {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.score-description {
  font-size: 1.25rem;
  font-weight: 600;
}

.score-excellent {
  color: #67C23A;
}

.score-good {
  color: #409EFF;
}

.score-average {
  color: #E6A23C;
}

.score-poor {
  color: #F56C6C;
}

.review-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--claude-gray-200);
}
</style>