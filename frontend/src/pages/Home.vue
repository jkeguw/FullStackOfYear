<template>
  <div class="home-page">
    <div class="hero bg-gradient-to-r from-[#1E1E1E] to-[#333333] py-16 text-white">
      <div class="container mx-auto px-4">
        <h1 class="text-4xl font-bold mb-4">Mouse Hub</h1>
        <p class="text-xl mb-8">Everything you need to know about mice is here</p>

        <!-- Vertical button layout -->
        <div class="flex flex-col items-center gap-6 mt-12 max-w-md mx-auto space-y-6">
          <router-link to="/compare" class="feature-card w-full">
            <div class="rgb-border">
              <div class="card-content w-full">
                <div class="icon-container mb-4">
                  <i class="el-icon-sort text-4xl"></i>
                </div>
                <h3 class="text-xl font-bold">Mouse Comparison</h3>
                <div class="card-hover-info">
                  <p>Compare different mouse shapes and sizes to find the perfect fit for you</p>
                </div>
              </div>
            </div>
          </router-link>

          <router-link to="/similar" class="feature-card w-full">
            <div class="rgb-border">
              <div class="card-content w-full">
                <div class="icon-container mb-4">
                  <i class="el-icon-search text-4xl"></i>
                </div>
                <h3 class="text-xl font-bold">Find Similar</h3>
                <div class="card-hover-info">
                  <p>Discover mice with similar shapes based on your favorites</p>
                </div>
              </div>
            </div>
          </router-link>

          <router-link to="/database" class="feature-card w-full">
            <div class="rgb-border">
              <div class="card-content w-full">
                <div class="icon-container mb-4">
                  <i class="el-icon-mouse text-4xl"></i>
                </div>
                <h3 class="text-xl font-bold">Mouse Database</h3>
                <div class="card-hover-info">
                  <p>Browse the complete mouse database with detailed specs and comparison data</p>
                </div>
              </div>
            </div>
          </router-link>

          <router-link to="/reviews" class="feature-card w-full">
            <div class="rgb-border">
              <div class="card-content w-full">
                <div class="icon-container mb-4">
                  <i class="el-icon-star text-4xl"></i>
                </div>
                <h3 class="text-xl font-bold">User Reviews</h3>
                <div class="card-hover-info">
                  <p>Read detailed user reviews and share your own experience</p>
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 移除专业工具版块 -->

    <!-- 移除灵敏度计算工具介绍 -->
  </div>
</template>

<script setup lang="ts">
// 保留必要的导入
import { onMounted } from 'vue';

// 动态鼠标跟踪效果
onMounted(() => {
  const cards = document.querySelectorAll('.feature-card');

  cards.forEach((card) => {
    card.addEventListener('mousemove', handleMouseMove);
    card.addEventListener('mouseleave', handleMouseLeave);
  });

  function handleMouseMove(e: MouseEvent) {
    const card = e.currentTarget as HTMLElement;
    const rect = card.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    // 计算相对位置 (0-1)
    const relX = x / rect.width;
    const relY = y / rect.height;

    // 施加倾斜效果
    card.style.transform = `perspective(1000px) rotateX(${(relY - 0.5) * 10}deg) rotateY(${(relX - 0.5) * -10}deg)`;

    // 更新RGB边框动画
    const border = card.querySelector('.rgb-border') as HTMLElement;
    if (border) {
      border.style.setProperty('--mouse-x', `${relX * 100}%`);
      border.style.setProperty('--mouse-y', `${relY * 100}%`);
    }
  }

  function handleMouseLeave(e: MouseEvent) {
    const card = e.currentTarget as HTMLElement;
    card.style.transform = 'perspective(1000px) rotateX(0) rotateY(0)';
  }
});
</script>

<style scoped>
.hero {
  background: linear-gradient(to right, var(--claude-bg-dark), var(--claude-bg-medium));
  background-position: center;
  background-size: cover;
}

.feature-card {
  position: relative;
  border-radius: 0.75rem;
  transition: all 0.3s ease;
  transform-style: preserve-3d;
  cursor: pointer;
}

.rgb-border {
  position: relative;
  border-radius: 0.75rem;
  overflow: hidden;
  --mouse-x: 50%;
  --mouse-y: 50%;
}

.rgb-border::before {
  content: '';
  position: absolute;
  inset: 0;
  background: conic-gradient(
    from calc(var(--mouse-x) + var(--mouse-y)),
    #ff0000,
    #ffff00,
    #00ff00,
    #00ffff,
    #0000ff,
    #ff00ff,
    #ff0000
  );
  border-radius: 0.75rem;
  z-index: -1;
  animation: rgb-rotate 3s linear infinite;
  opacity: 0;
  transition: opacity 0.3s ease;
}

/* 确保按钮颜色覆盖整个按钮 */
.feature-card .card-content {
  width: 100%;
  height: 100%;
  background-color: rgba(42, 42, 42, 0.95);
}

.feature-card:hover .rgb-border::before {
  opacity: 1;
}

.card-content {
  background: rgba(42, 42, 42, 0.95);
  padding: 1.5rem;
  border-radius: 0.75rem;
  text-align: center;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
  transition:
    transform 0.3s ease,
    background-color 0.3s ease;
  width: 100%; /* 使用100%宽度以适应竖向布局 */
  max-width: 500px; /* 限制最大宽度 */
  /* 使按钮整体显示为长条形 */
  min-height: 80px;
}

.card-hover-info {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 1rem;
  background: rgba(18, 18, 18, 0.9);
  color: #fff;
  border-radius: 0 0 0.75rem 0.75rem;
  transform: translateY(100%);
  transition: transform 0.3s ease;
  z-index: 2;
}

.feature-card:hover .card-hover-info {
  transform: translateY(0);
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: linear-gradient(135deg, #7d5af3, #6a48e0);
  color: white;
}

@keyframes rgb-rotate {
  0% {
    filter: hue-rotate(0deg);
  }
  100% {
    filter: hue-rotate(360deg);
  }
}
</style>
