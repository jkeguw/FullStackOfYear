<template>
  <div class="review-gallery">
    <div v-if="!images.length" class="no-images">暂无图片</div>
    <div v-else class="gallery-container">
      <div class="gallery-thumbnails">
        <div
          v-for="(image, index) in images"
          :key="`thumb-${index}`"
          class="thumbnail"
          :class="{ active: index === activeIndex }"
          @click="setActiveImage(index)"
        >
          <img
            :src="image.thumbnailUrl || image.url"
            :alt="`图片 ${index + 1}`"
            class="img-fluid"
          />
        </div>
      </div>

      <div class="gallery-main-image">
        <div class="image-container">
          <img
            v-if="activeImage"
            :src="activeImage.url"
            :alt="activeImage.caption || `图片 ${activeIndex + 1}`"
            class="img-fluid"
          />
        </div>

        <div v-if="activeImage && activeImage.caption" class="image-caption">
          {{ activeImage.caption }}
        </div>

        <div class="gallery-controls">
          <el-button
            circle
            @click="prevImage"
            :disabled="activeIndex === 0"
            icon="el-icon-arrow-left"
          />
          <span class="image-counter">{{ activeIndex + 1 }} / {{ images.length }}</span>
          <el-button
            circle
            @click="nextImage"
            :disabled="activeIndex === images.length - 1"
            icon="el-icon-arrow-right"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineProps, watch } from 'vue';

interface GalleryImage {
  url: string;
  thumbnailUrl?: string;
  caption?: string;
}

const props = defineProps<{
  images: GalleryImage[];
  initialIndex?: number;
}>();

// 状态
const activeIndex = ref(props.initialIndex || 0);

// 计算属性
const activeImage = computed(() => {
  return props.images[activeIndex.value] || null;
});

// 方法
function setActiveImage(index: number) {
  if (index >= 0 && index < props.images.length) {
    activeIndex.value = index;
  }
}

function nextImage() {
  if (activeIndex.value < props.images.length - 1) {
    activeIndex.value++;
  }
}

function prevImage() {
  if (activeIndex.value > 0) {
    activeIndex.value--;
  }
}

// 监听图片数组变化
watch(
  () => props.images,
  () => {
    if (activeIndex.value >= props.images.length) {
      activeIndex.value = props.images.length > 0 ? 0 : -1;
    }
  }
);
</script>

<style scoped>
.review-gallery {
  margin-bottom: 2rem;
}

.no-images {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  color: var(--claude-gray-600);
}

.gallery-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

@media (min-width: 768px) {
  .gallery-container {
    grid-template-columns: 100px 1fr;
  }
}

.gallery-thumbnails {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
  padding-bottom: 0.5rem;
}

@media (min-width: 768px) {
  .gallery-thumbnails {
    flex-direction: column;
    overflow-x: hidden;
    max-height: 450px;
    overflow-y: auto;
  }
}

.thumbnail {
  width: 80px;
  height: 80px;
  flex-shrink: 0;
  cursor: pointer;
  border: 2px solid transparent;
  border-radius: 0.25rem;
  overflow: hidden;
  transition: border-color 0.2s ease;
}

.thumbnail.active {
  border-color: var(--claude-primary);
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery-main-image {
  position: relative;
}

.image-container {
  width: 100%;
  height: 0;
  padding-bottom: 60%;
  position: relative;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  overflow: hidden;
}

.image-container img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.image-caption {
  margin-top: 0.75rem;
  font-size: 0.875rem;
  color: var(--claude-gray-600);
  text-align: center;
}

.gallery-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.image-counter {
  font-size: 0.875rem;
  color: var(--claude-gray-600);
}
</style>
