<template>
  <div class="comparison-table">
    <h3 class="text-lg font-medium mb-3">参数对比</h3>

    <div v-if="similarityScore !== undefined" class="similarity-score mb-4">
      <div class="font-medium text-sm">相似度评分</div>
      <el-progress
        :percentage="similarityScore"
        :color="getSimilarityColor(similarityScore)"
        :format="(percent) => `${percent}%`"
        :stroke-width="10"
      />
    </div>

    <el-table :data="tableData" border stripe>
      <el-table-column prop="property" label="参数" width="150">
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

      <!-- 移除了差异列 -->
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue';
import type { MouseComparisonResult } from '@/models/MouseModel';

// 属性说明
const propertyDescriptions: Record<string, string> = {
  长度: '鼠标的前后长度(mm)',
  宽度: '鼠标的左右宽度(mm)',
  高度: '鼠标的最高点到桌面的高度(mm)',
  重量: '鼠标的重量(g)',
  握宽: '鼠标的握持宽度，通常是中部宽度(mm)',
  最大DPI: '鼠标传感器的最大DPI值',
  轮询率: '鼠标数据刷新率(Hz)',
  侧键数量: '鼠标侧面按键的数量',
  形状类型: '鼠标的整体形状类型',
  推荐握持方式: '鼠标适合的握持姿势',
  凸起位置: '鼠标的主要凸起位置',
  连接方式: '鼠标的连接类型'
};

// Props定义
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

// 计算属性
const tableData = computed(() => {
  if (!props.data) return [];

  const data = [...props.data];

  // 排序
  if (props.sortBy === 'difference') {
    data.sort((a, b) => b.differencePercent - a.differencePercent);
  } else {
    // 按属性名称排序
    const propertyOrder = [
      '长度',
      '宽度',
      '高度',
      '重量',
      '握宽',
      '形状类型',
      '最大DPI',
      '轮询率',
      '侧键数量',
      '推荐握持方式'
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

// 方法
function formatValue(value: any, property: string): string {
  if (value === undefined || value === null) return '未知';

  // 根据属性类型格式化值
  if (property === '长度' || property === '宽度' || property === '高度' || property === '握宽') {
    return `${value}mm`;
  } else if (property === '重量') {
    return `${value}g`;
  } else if (property === '最大DPI') {
    return value.toLocaleString();
  } else if (property === '轮询率') {
    return `${value}Hz`;
  } else if (Array.isArray(value)) {
    return value.join(', ');
  }

  return String(value);
}

function formatDifference(value: number): string {
  if (value === 0) return '相同';
  return `${value.toFixed(1)}%`;
}

function getDifferenceClass(value: number): string {
  if (value === 0) return 'text-green-500';
  if (value < 10) return 'text-blue-500';
  if (value < 25) return 'text-amber-500';
  return 'text-red-500';
}

function getSimilarityColor(score: number): string {
  if (score >= 90) return '#67C23A'; // 绿色
  if (score >= 75) return '#409EFF'; // 蓝色
  if (score >= 50) return '#E6A23C'; // 橙色
  return '#F56C6C'; // 红色
}

function getPropertyDescription(property: string): string {
  return propertyDescriptions[property] || '';
}

function isHighlight(row: any, index: number): boolean {
  if (row.differencePercent === 0) return false;

  // 找出数值型属性中的最大/最小值
  if (['长度', '宽度', '高度', '重量', '握宽', '最大DPI', '轮询率'].includes(row.property)) {
    const numValues = row.values.map((v: any) => Number(v));
    const maxValue = Math.max(...numValues);
    const minValue = Math.min(...numValues);

    // 对于DPI和轮询率，最大值更好
    if (['最大DPI', '轮询率'].includes(row.property)) {
      return numValues[index] === maxValue;
    }

    // 对于重量，最小值更好
    if (row.property === '重量') {
      return numValues[index] === minValue;
    }

    // 对于尺寸，突出显示极值
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
