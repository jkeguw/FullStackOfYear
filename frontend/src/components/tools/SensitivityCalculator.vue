<template>
  <div class="sensitivity-calculator p-6 rounded-lg bg-gray-800 text-white">
    <div class="grid md:grid-cols-2 gap-6">
      <!-- 游戏灵敏度转换 -->
      <div class="card p-5 bg-gray-700 rounded-lg">
        <h3 class="text-lg font-medium mb-4">游戏灵敏度转换</h3>
        <el-form label-position="top">
          <el-form-item label="鼠标DPI">
            <el-input-number
              v-model="dpi"
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
                class="result-display py-2 px-3 bg-blue-900 rounded text-white text-lg font-bold text-center"
              >
                {{ convertedGameSensitivity }}
              </div>
            </el-form-item>
          </div>

          <div class="mt-3 text-center">
            <p class="text-sm text-gray-400">cm/360°: {{ cmPer360 }}</p>
          </div>
        </el-form>
      </div>

      <!-- 三阶校准法 -->
      <div class="card p-5 bg-gray-700 rounded-lg">
        <h3 class="text-lg font-medium mb-4">三阶校准法</h3>
        <div v-if="currentPhase === 0" class="phase-start">
          <p class="mb-4 text-gray-300">通过三个渐进阶段帮助您找到最适合的鼠标灵敏度</p>
          <el-form label-position="top">
            <el-form-item label="起始基准值">
              <el-input-number
                v-model="baseValue"
                :min="0.1"
                :max="10"
                :step="0.1"
                :precision="2"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="游戏">
              <el-select v-model="calibrationGame" class="w-full">
                <el-option
                  v-for="game in games"
                  :key="game.id"
                  :label="game.name"
                  :value="game.id"
                />
              </el-select>
            </el-form-item>
            <div class="text-center mt-4">
              <el-button type="primary" @click="startCalibration">开始校准</el-button>
            </div>
          </el-form>
        </div>

        <div v-else class="calibration-in-progress">
          <div class="phase-indicator mb-4">
            <div class="flex justify-between text-xs text-gray-400 mb-1">
              <span>阶段 {{ currentPhase }}/3</span>
              <span>基准值: {{ baseValue.toFixed(3) }}</span>
            </div>
            <el-progress :percentage="(currentPhase / 3) * 100" />
          </div>

          <div class="phase-description mb-6">
            <p v-if="currentPhase === 1" class="text-sm text-gray-300">
              <b>第一阶段（矫快敏）</b>: 选择哪个值感觉更适合您
            </p>
            <p v-if="currentPhase === 2" class="text-sm text-gray-300">
              <b>第二阶段（拔慢敏）</b>: 继续微调您的灵敏度
            </p>
            <p v-if="currentPhase === 3" class="text-sm text-gray-300">
              <b>第三阶段（细常敏）</b>: 最终精确调整
            </p>
          </div>

          <div class="choice-buttons grid grid-cols-2 gap-4 mb-6">
            <el-button @click="selectValue('left')" type="success" class="h-20">
              {{ leftValue.toFixed(3) }}
              <div class="text-xs mt-1">更快的灵敏度</div>
            </el-button>
            <el-button @click="selectValue('right')" type="warning" class="h-20">
              {{ rightValue.toFixed(3) }}
              <div class="text-xs mt-1">更慢的灵敏度</div>
            </el-button>
          </div>

          <div class="action-buttons flex justify-between">
            <el-button @click="resetCalibration" size="small">重置</el-button>
            <el-button v-if="currentPhase < 3" @click="nextPhase" size="small" type="info"
              >跳过此阶段</el-button
            >
          </div>
        </div>

        <div v-if="currentPhase === 4" class="calibration-result mt-4 p-4 bg-blue-900 rounded-lg">
          <h4 class="text-center font-bold mb-3">校准完成</h4>
          <div class="text-center">
            <p class="mb-2">
              您的理想灵敏度: <span class="text-xl font-bold">{{ baseValue.toFixed(3) }}</span>
            </p>
            <p class="text-sm text-gray-300">cm/360°: {{ calculateCmPer360(baseValue) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 灵敏快分法和极敏内推法 -->
    <div class="grid md:grid-cols-2 gap-6 mt-6">
      <div class="card p-5 bg-gray-700 rounded-lg">
        <h3 class="text-lg font-medium mb-4">灵敏快分法</h3>
        <div v-if="!fastBisectionStarted">
          <p class="mb-4 text-sm text-gray-300">
            通过9个步骤快速精确地确定您的理想灵敏度，无需大量测试
          </p>
          <el-form label-position="top">
            <el-form-item label="起始基准值">
              <el-input-number
                v-model="bisectionBaseValue"
                :min="0.1"
                :max="10"
                :step="0.1"
                :precision="2"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
            <div class="text-center mt-4">
              <el-button type="primary" @click="startFastBisection">开始快分法</el-button>
            </div>
          </el-form>
        </div>

        <div v-else>
          <div class="phase-indicator mb-4">
            <div class="flex justify-between text-xs text-gray-400 mb-1">
              <span>步骤 {{ bisectionStep }}/9</span>
              <span>基准值: {{ bisectionBaseValue.toFixed(3) }}</span>
            </div>
            <el-progress :percentage="(bisectionStep / 9) * 100" />
          </div>

          <div class="choice-buttons grid grid-cols-2 gap-4 mb-6">
            <el-button @click="selectBisectionValue('low')" type="success" class="h-16">
              {{ bisectionLowValue.toFixed(3) }}
              <div class="text-xs mt-1">更低的值</div>
            </el-button>
            <el-button @click="selectBisectionValue('high')" type="warning" class="h-16">
              {{ bisectionHighValue.toFixed(3) }}
              <div class="text-xs mt-1">更高的值</div>
            </el-button>
          </div>

          <div class="action-buttons flex justify-between">
            <el-button @click="resetFastBisection" size="small">重置</el-button>
          </div>

          <div v-if="bisectionStep > 9" class="calibration-result mt-4 p-4 bg-blue-900 rounded-lg">
            <h4 class="text-center font-bold mb-3">快分法完成</h4>
            <div class="text-center">
              <p class="mb-2">
                您的理想灵敏度:
                <span class="text-xl font-bold">{{ bisectionBaseValue.toFixed(3) }}</span>
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="card p-5 bg-gray-700 rounded-lg">
        <h3 class="text-lg font-medium mb-4">极敏内推法</h3>
        <p class="mb-4 text-sm text-gray-300">
          基于您的最快和最慢可承受灵敏度，计算出最适合的灵敏度值
        </p>

        <el-form label-position="top">
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
            <div class="result-card p-3 bg-blue-800 rounded-lg">
              <div class="text-sm text-center mb-1">瞄准场景推荐</div>
              <div class="text-xl font-bold text-center">{{ aimingScenarioSens.toFixed(3) }}</div>
              <div class="text-xs text-gray-400 text-center mt-1">(k + 3m) / 4</div>
            </div>

            <div class="result-card p-3 bg-blue-800 rounded-lg">
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

// 游戏数据
const games = [
  { id: 'csgo', name: 'CS:GO/CS2', coefficient: 0.022 },
  { id: 'valorant', name: 'Valorant', coefficient: 0.07 },
  { id: 'apex', name: 'Apex Legends', coefficient: 0.022 },
  { id: 'overwatch', name: 'Overwatch 2', coefficient: 0.0066 },
  { id: 'r6siege', name: 'Rainbow Six Siege', coefficient: 0.00223 }
];

// 灵敏度转换
const dpi = ref(800);
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

const cmPer360 = computed(() => {
  const sourceCoeff = getGameCoefficient(sourceGame.value);
  const cm = (360 * 2.54) / (sourceSensitivity.value * dpi.value * sourceCoeff);
  return cm.toFixed(2) + ' cm/360°';
});

// 三阶校准法
const currentPhase = ref(0);
const baseValue = ref(1.0);
const calibrationGame = ref('csgo');
const leftValue = ref(0);
const rightValue = ref(0);

const calculateCmPer360 = (sensitivity: number) => {
  const gameCoeff = getGameCoefficient(calibrationGame.value);
  const cm = (360 * 2.54) / (sensitivity * dpi.value * gameCoeff);
  return cm.toFixed(2) + ' cm/360°';
};

const calculatePhaseValues = () => {
  if (currentPhase.value === 1) {
    // 阶段1（矫快敏）
    leftValue.value = (baseValue.value * 360) / (360 + 180);
    rightValue.value = (baseValue.value * 360) / (360 - 56.25);
  } else if (currentPhase.value === 2) {
    // 阶段2（拔慢敏）
    leftValue.value = (baseValue.value * 360) / (360 + 22.5);
    rightValue.value = (baseValue.value * 360) / (360 - 45);
  } else if (currentPhase.value === 3) {
    // 阶段3（细常敏）
    leftValue.value = (baseValue.value * 360) / (360 + 11.25);
    rightValue.value = (baseValue.value * 360) / (360 - 22.5);
  }
};

const startCalibration = () => {
  currentPhase.value = 1;
  calculatePhaseValues();
};

const resetCalibration = () => {
  currentPhase.value = 0;
  baseValue.value = 1.0;
};

const selectValue = (side: string) => {
  if (side === 'left') {
    baseValue.value = leftValue.value;
  } else {
    baseValue.value = rightValue.value;
  }

  if (currentPhase.value < 3) {
    nextPhase();
  } else {
    currentPhase.value = 4; // 完成
  }
};

const nextPhase = () => {
  currentPhase.value++;
  calculatePhaseValues();
};

// 灵敏快分法
const bisectionStep = ref(0);
const bisectionBaseValue = ref(1.0);
const bisectionLowValue = ref(0);
const bisectionHighValue = ref(0);
const fastBisectionStarted = ref(false);

// 比例系数数组
const ratioCoefficients = [0.4375, 0.375, 0.3125, 0.1875, 0.125, 0.125, 0.125, 0.0625, 0.03125];

const calculateBisectionValues = () => {
  const range = bisectionBaseValue.value * ratioCoefficients[bisectionStep.value - 1];
  bisectionLowValue.value = bisectionBaseValue.value - range;
  bisectionHighValue.value = bisectionBaseValue.value + range;
};

const startFastBisection = () => {
  fastBisectionStarted.value = true;
  bisectionStep.value = 1;
  calculateBisectionValues();
};

const resetFastBisection = () => {
  fastBisectionStarted.value = false;
  bisectionStep.value = 0;
  bisectionBaseValue.value = 1.0;
};

const selectBisectionValue = (side: string) => {
  if (side === 'low') {
    bisectionBaseValue.value = bisectionLowValue.value;
  } else {
    bisectionBaseValue.value = bisectionHighValue.value;
  }

  if (bisectionStep.value < 9) {
    bisectionStep.value++;
    calculateBisectionValues();
  } else {
    bisectionStep.value = 10; // 完成
  }
};

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
.el-input-number :deep(.el-input__inner) {
  color: white !important;
  background-color: #4a4a4a !important;
}

.el-select :deep(.el-input__inner) {
  color: white !important;
  background-color: #4a4a4a !important;
}

.el-button {
  transition: transform 0.2s;
}

.el-button:hover {
  transform: translateY(-2px);
}

.result-display,
.result-card {
  transition: all 0.3s ease;
}

.result-display:hover,
.result-card:hover {
  transform: scale(1.03);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
