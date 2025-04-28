<template>
  <div class="sensitivity-tool-page container mx-auto py-8 px-4">
    <h1 class="text-3xl font-bold mb-8 text-white">灵敏度工具</h1>

    <div class="mb-8 bg-[#1E1E1E] rounded-lg p-6">
      <h2 class="text-2xl font-bold mb-4 text-white">工具说明</h2>
      <p class="text-gray-300 mb-4">
        这个页面集成了多种灵敏度计算工具，帮助您找到最适合自己的游戏鼠标灵敏度设置。
        您可以使用DPI转换计算器在不同的鼠标DPI设置之间转换灵敏度，或者利用游戏灵敏度转换工具在不同游戏间保持相同的手感。
      </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <!-- DPI转换计算器 -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">DPI转换计算器</h2>
        <p class="text-gray-300 mb-4">更换鼠标DPI设置时，使用此工具保持相同的鼠标移动手感。</p>

        <el-form :model="dpiForm" label-position="top" class="text-white">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <el-form-item label="原始DPI" class="mb-3">
              <el-input-number
                v-model="dpiForm.sourceDPI"
                :min="100"
                :max="32000"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="目标DPI" class="mb-3">
              <el-input-number
                v-model="dpiForm.targetDPI"
                :min="100"
                :max="32000"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
          </div>
          <el-form-item label="原始灵敏度" class="mb-3">
            <el-input-number
              v-model="dpiForm.sensitivity"
              :min="0.1"
              :max="10"
              :step="0.1"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>
          <el-form-item class="mt-6 text-center">
            <div class="p-3 bg-[#2A2A2A] border border-[#444444] rounded-md">
              <span class="text-white">转换后灵敏度: </span>
              <span class="text-xl font-bold text-white">{{ convertedDpiSensitivity }}</span>
            </div>
          </el-form-item>
          <div class="mt-4 p-3 bg-[#2A2A2A] border border-[#444444] rounded-md">
            <p class="text-sm text-gray-300">
              cm/360° = <span class="font-semibold text-white">{{ cmPer360 }}</span>
            </p>
          </div>
        </el-form>
      </div>

      <!-- 游戏灵敏度转换 -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">游戏灵敏度转换</h2>
        <p class="text-gray-300 mb-4">在不同游戏之间转换灵敏度，保持一致的肌肉记忆和瞄准手感。</p>

        <el-form label-position="top" class="text-white">
          <el-form-item label="鼠标DPI">
            <el-input-number
              v-model="gameDpi"
              :min="100"
              :max="32000"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>

          <div class="grid grid-cols-2 gap-4">
            <el-form-item label="源游戏">
              <el-select v-model="sourceGame" class="w-full">
                <el-option
                  v-for="game in games"
                  :key="game.id"
                  :label="game.name"
                  :value="game.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="源游戏灵敏度">
              <el-input-number
                v-model="sourceSensitivity"
                :min="0.001"
                :max="100"
                :step="0.001"
                :precision="3"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <el-form-item label="目标游戏">
              <el-select v-model="targetGame" class="w-full">
                <el-option
                  v-for="game in games"
                  :key="game.id"
                  :label="game.name"
                  :value="game.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="计算结果">
              <div
                class="result-display py-2 px-3 bg-[#2A2A2A] border border-[#444444] rounded text-white text-lg font-bold text-center"
              >
                {{ convertedGameSensitivity }}
              </div>
            </el-form-item>
          </div>

          <div class="mt-3 text-center">
            <p class="text-sm text-gray-400">cm/360°: {{ gameCmPer360 }}</p>
          </div>
        </el-form>
      </div>
    </div>

    <!-- 三阶校准法和灵敏快分法 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 三阶校准法 -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">三阶校准法</h2>
        <p class="text-gray-300 mb-4">
          使用渐进式校准方法找到最适合您的灵敏度，提供精确的游戏体验。
        </p>

        <SensitivityCalculator />
      </div>

      <!-- 极敏内推法 -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">极敏内推法</h2>
        <p class="text-gray-300 mb-4">基于您确定的速度范围，使用数学公式推导最佳灵敏度值。</p>

        <el-form label-position="top" class="text-white">
          <el-form-item label="最快可接受灵敏度 (k)">
            <el-input-number
              v-model="fastestSens"
              :min="0.1"
              :max="10"
              :step="0.1"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>

          <el-form-item label="最慢可接受灵敏度 (m)">
            <el-input-number
              v-model="slowestSens"
              :min="0.1"
              :max="10"
              :step="0.1"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>

          <div class="grid grid-cols-2 gap-4 mt-6">
            <div class="result-card p-3 bg-[#2A2A2A] border border-[#444444] rounded-lg">
              <div class="text-sm text-center mb-1">瞄准场景推荐</div>
              <div class="text-xl font-bold text-center">{{ aimingScenarioSens.toFixed(3) }}</div>
              <div class="text-xs text-gray-400 text-center mt-1">(k + 3m) / 4</div>
            </div>

            <div class="result-card p-3 bg-[#2A2A2A] border border-[#444444] rounded-lg">
              <div class="text-sm text-center mb-1">游戏场景推荐</div>
              <div class="text-xl font-bold text-center">{{ gamingScenarioSens.toFixed(3) }}</div>
              <div class="text-xs text-gray-400 text-center mt-1">(3k + m) / 4</div>
            </div>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { convertDPI } from '@/utils/dpi';
import SensitivityCalculator from '@/components/tools/SensitivityCalculator.vue';

// DPI转换计算器
const dpiForm = ref({
  sourceDPI: 800,
  targetDPI: 1600,
  sensitivity: 1
});

const convertedDpiSensitivity = computed(() => {
  return convertDPI(
    dpiForm.value.sourceDPI,
    dpiForm.value.targetDPI,
    dpiForm.value.sensitivity
  ).toFixed(3);
});

const cmPer360 = computed(() => {
  // cm/360° = 360 × 2.54 / (灵敏度 × DPI × 游戏系数)
  // 使用CSGO的游戏系数0.022作为默认值
  const gameCoefficient = 0.022;
  const cm = (360 * 2.54) / (dpiForm.value.sensitivity * dpiForm.value.sourceDPI * gameCoefficient);
  return cm.toFixed(2) + ' cm';
});

// 游戏灵敏度转换
const games = [
  { id: 'csgo', name: 'CS:GO/CS2', coefficient: 0.022 },
  { id: 'valorant', name: 'Valorant', coefficient: 0.07 },
  { id: 'apex', name: 'Apex Legends', coefficient: 0.022 },
  { id: 'overwatch', name: 'Overwatch 2', coefficient: 0.0066 },
  { id: 'r6siege', name: 'Rainbow Six Siege', coefficient: 0.00223 }
];

const gameDpi = ref(800);
const sourceGame = ref('csgo');
const targetGame = ref('valorant');
const sourceSensitivity = ref(2.0);

const getGameCoefficient = (gameId: string) => {
  const game = games.find((g) => g.id === gameId);
  return game ? game.coefficient : 0.022;
};

const convertedGameSensitivity = computed(() => {
  const sourceCoeff = getGameCoefficient(sourceGame.value);
  const targetCoeff = getGameCoefficient(targetGame.value);

  const result = (sourceSensitivity.value * sourceCoeff) / targetCoeff;
  return result.toFixed(3);
});

const gameCmPer360 = computed(() => {
  const sourceCoeff = getGameCoefficient(sourceGame.value);
  const cm = (360 * 2.54) / (sourceSensitivity.value * gameDpi.value * sourceCoeff);
  return cm.toFixed(2) + ' cm/360°';
});

// 极敏内推法
const fastestSens = ref(2.0);
const slowestSens = ref(1.0);

const aimingScenarioSens = computed(() => {
  return (fastestSens.value + 3 * slowestSens.value) / 4;
});

const gamingScenarioSens = computed(() => {
  return (3 * fastestSens.value + slowestSens.value) / 4;
});
</script>

<style scoped>
.sensitivity-tool-page {
  color: var(--claude-text-light);
}

.el-input-number :deep(.el-input__inner) {
  color: white !important;
  background-color: var(--claude-bg-light) !important;
  border-color: var(--claude-border-light) !important;
}

.el-select :deep(.el-input__inner) {
  color: white !important;
  background-color: var(--claude-bg-light) !important;
  border-color: var(--claude-border-light) !important;
}

.result-display,
.result-card {
  transition: all 0.3s ease;
}

.result-display:hover,
.result-card:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
</style>
