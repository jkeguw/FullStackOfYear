<template>
  <div class="comparison-table">
    <h3 class="text-lg font-medium mb-3">Parameter Comparison</h3>

    <div v-if="similarityScore !== undefined" class="similarity-score mb-4">
      <div class="font-medium text-sm">Similarity Score</div>
      <el-progress
        :percentage="similarityScore"
        :color="getSimilarityColor(similarityScore)"
        :format="(percent) => `${percent}%`"
        :stroke-width="10"
      />
    </div>

    <el-table :data="tableData" border stripe>
      <el-table-column prop="property" label="Parameter" width="150">
        <template #default="scope">
          <el-tooltip
            v-if="getPropertyDescription(scope.row.property)"
            :content="getPropertyDescription(scope.row.property)"
            placement="top"
          >
            <div class="cursor-help underline-dotted">{{ scope.row.property }}</div>
          </el-tooltip>
          <span v-else>{{ scope.row.property }}</span>
        </template>
      </el-table-column>

      <el-table-column
        v-for="(label, index) in columnLabels"
        :key="index"
        :label="label"
        :min-width="120"
      >
        <template #default="scope">
          <div v-if="isHighlight(scope.row, index)" class="highlighted-value">
            {{ formatValue(scope.row.values[index], scope.row.property) }}
          </div>
          <div v-else>
            {{ formatValue(scope.row.values[index], scope.row.property) }}
          </div>
        </template>
      </el-table-column>

      <!-- Removed difference column -->
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue';
import type { MouseComparisonResult } from '@/models/MouseModel';

// Property descriptions
const propertyDescriptions: Record<string, string> = {
  Length: 'The front-to-back length of the mouse (mm)',
  Width: 'The side-to-side width of the mouse (mm)',
  Height: 'The height from the highest point to the desk surface (mm)',
  Weight: 'The weight of the mouse (g)',
  Grip_Width: 'The width where you grip the mouse, usually the middle width (mm)',
  Max_DPI: 'The maximum DPI value of the mouse sensor',
  Polling_Rate: 'The data refresh rate of the mouse (Hz)',
  Side_Buttons: 'The number of buttons on the side of the mouse',
  Shape_Type: 'The overall shape type of the mouse',
  Recommended_Grip: 'The grip posture suitable for the mouse',
  Hump_Position: 'The main protrusion position of the mouse',
  Connection_Type: 'The connection type of the mouse'
};

// Props definition
const props = defineProps<{
  data: Array<{
    property: string;
    values: any[];
    differencePercent: number;
  }>;
  columnLabels: string[];
  sortBy?: 'difference' | 'property';
  similarityScore?: number;
}>();

// Computed properties
const tableData = computed(() => {
  if (!props.data) return [];

  const data = [...props.data];

  // Sort
  if (props.sortBy === 'difference') {
    data.sort((a, b) => b.differencePercent - a.differencePercent);
  } else {
    // Sort by property name
    const propertyOrder = [
      'Length',
      'Width',
      'Height',
      'Weight',
      'Grip_Width',
      'Shape_Type',
      'Max_DPI',
      'Polling_Rate',
      'Side_Buttons',
      'Recommended_Grip'
    ];
    data.sort((a, b) => {
      const indexA = propertyOrder.indexOf(a.property);
      const indexB = propertyOrder.indexOf(b.property);
      if (indexA === -1 && indexB === -1) return a.property.localeCompare(b.property);
      if (indexA === -1) return 1;
      if (indexB === -1) return -1;
      return indexA - indexB;
    });
  }

  return data;
});

// Methods
function formatValue(value: any, property: string): string {
  if (value === undefined || value === null) return 'Unknown';

  // Format value based on property type
  if (property === 'Length' || property === 'Width' || property === 'Height' || property === 'Grip_Width') {
    return `${value}mm`;
  } else if (property === 'Weight') {
    return `${value}g`;
  } else if (property === 'Max_DPI') {
    return value.toLocaleString();
  } else if (property === 'Polling_Rate') {
    return `${value}Hz`;
  } else if (Array.isArray(value)) {
    return value.join(', ');
  }

  return String(value);
}

function formatDifference(value: number): string {
  if (value === 0) return 'Same';
  return `${value.toFixed(1)}%`;
}

function getDifferenceClass(value: number): string {
  if (value === 0) return 'text-green-500';
  if (value < 10) return 'text-blue-500';
  if (value < 25) return 'text-amber-500';
  return 'text-red-500';
}

function getSimilarityColor(score: number): string {
  if (score >= 90) return '#67C23A'; // Green
  if (score >= 75) return '#409EFF'; // Blue
  if (score >= 50) return '#E6A23C'; // Orange
  return '#F56C6C'; // Red
}

function getPropertyDescription(property: string): string {
  return propertyDescriptions[property] || '';
}

function isHighlight(row: any, index: number): boolean {
  if (row.differencePercent === 0) return false;

  // Find maximum/minimum values in numeric properties
  if (['Length', 'Width', 'Height', 'Weight', 'Grip_Width', 'Max_DPI', 'Polling_Rate'].includes(row.property)) {
    const numValues = row.values.map((v: any) => Number(v));
    const maxValue = Math.max(...numValues);
    const minValue = Math.min(...numValues);

    // For DPI and polling rate, higher is better
    if (['Max_DPI', 'Polling_Rate'].includes(row.property)) {
      return numValues[index] === maxValue;
    }

    // For weight, lower is better
    if (row.property === 'Weight') {
      return numValues[index] === minValue;
    }

    // For dimensions, highlight extreme values
    return numValues[index] === maxValue || numValues[index] === minValue;
  }

  return false;
}
</script>

<style scoped>
.underline-dotted {
  border-bottom: 1px dotted currentColor;
}

.highlighted-value {
  font-weight: 500;
  position: relative;
}

.highlighted-value::after {
  content: '•';
  position: absolute;
  top: -5px;
  right: -5px;
  color: var(--el-color-primary);
  font-size: 14px;
}
</style>
