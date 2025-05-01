<template>
  <div class="sensitivity-tool-page container mx-auto py-8 px-4">
    <h1 class="text-3xl font-bold mb-8 text-white">Sensitivity Tools</h1>

    <div class="mb-8 bg-[#1E1E1E] rounded-lg p-6">
      <h2 class="text-2xl font-bold mb-4 text-white">Tool Description</h2>
      <p class="text-gray-300 mb-4">
        This page integrates various sensitivity calculation tools to help you find the most suitable game mouse sensitivity settings for yourself.
        You can use the DPI Conversion Calculator to convert sensitivity between different mouse DPI settings, or use the Game Sensitivity Conversion tool to maintain the same feel across different games.
      </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <!-- DPI Conversion Calculator -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">DPI Conversion Calculator</h2>
        <p class="text-gray-300 mb-4">Use this tool to maintain the same mouse movement feel when changing your mouse DPI settings.</p>

        <el-form :model="dpiForm" label-position="top" class="text-white">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <el-form-item label="Original DPI" class="mb-3">
              <el-input-number
                v-model="dpiForm.sourceDPI"
                :min="100"
                :max="32000"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
            <el-form-item label="Target DPI" class="mb-3">
              <el-input-number
                v-model="dpiForm.targetDPI"
                :min="100"
                :max="32000"
                class="w-full"
                controls-position="right"
              />
            </el-form-item>
          </div>
          <el-form-item label="Original Sensitivity" class="mb-3">
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
              <span class="text-white">Converted Sensitivity: </span>
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

      <!-- Game Sensitivity Conversion -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">Game Sensitivity Conversion</h2>
        <p class="text-gray-300 mb-4">Convert sensitivity between different games to maintain consistent muscle memory and aiming feel.</p>

        <el-form label-position="top" class="text-white">
          <el-form-item label="Mouse DPI">
            <el-input-number
              v-model="gameDpi"
              :min="100"
              :max="32000"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>

          <div class="grid grid-cols-2 gap-4">
            <el-form-item label="Source Game">
              <el-select v-model="sourceGame" class="w-full">
                <el-option
                  v-for="game in games"
                  :key="game.id"
                  :label="game.name"
                  :value="game.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="Source Game Sensitivity">
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
            <el-form-item label="Target Game">
              <el-select v-model="targetGame" class="w-full">
                <el-option
                  v-for="game in games"
                  :key="game.id"
                  :label="game.name"
                  :value="game.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="Result">
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

    <!-- Three-Level Calibration and PSA Method -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Three-Level Calibration Method -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">Three-Level Calibration Method</h2>
        <p class="text-gray-300 mb-4">
          Use a progressive calibration method to find the sensitivity that's best for you, providing a precise gaming experience.
        </p>

        <SensitivityCalculator />
      </div>

      <!-- PSA Method -->
      <div class="bg-[#1E1E1E] rounded-lg p-6">
        <h2 class="text-2xl font-bold mb-4 text-white">PSA Method</h2>
        <p class="text-gray-300 mb-4">Derive the optimal sensitivity value using mathematical formulas based on your determined speed range.</p>

        <el-form label-position="top" class="text-white">
          <el-form-item label="Fastest Acceptable Sensitivity (k)">
            <el-input-number
              v-model="fastestSens"
              :min="0.1"
              :max="10"
              :step="0.1"
              class="w-full"
              controls-position="right"
            />
          </el-form-item>

          <el-form-item label="Slowest Acceptable Sensitivity (m)">
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
              <div class="text-sm text-center mb-1">Aiming Scenario Recommendation</div>
              <div class="text-xl font-bold text-center">{{ aimingScenarioSens.toFixed(3) }}</div>
              <div class="text-xs text-gray-400 text-center mt-1">(k + 3m) / 4</div>
            </div>

            <div class="result-card p-3 bg-[#2A2A2A] border border-[#444444] rounded-lg">
              <div class="text-sm text-center mb-1">Gaming Scenario Recommendation</div>
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

// DPI Conversion Calculator
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
  // cm/360° = 360 × 2.54 / (sensitivity × DPI × game coefficient)
  // Using CSGO's game coefficient 0.022 as the default value
  const gameCoefficient = 0.022;
  const cm = (360 * 2.54) / (dpiForm.value.sensitivity * dpiForm.value.sourceDPI * gameCoefficient);
  return cm.toFixed(2) + ' cm';
});

// Game Sensitivity Conversion
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

// PSA Method
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
