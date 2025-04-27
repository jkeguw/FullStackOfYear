<template>
  <div class="static-ruler-container">
    <div class="ruler" :style="rulerStyles">
      <div class="ruler-markings">
        <div 
          v-for="tick in ticks" 
          :key="tick.value" 
          class="ruler-tick"
          :class="{'major-tick': tick.isMajor}"
          :style="tickStyle(tick)"
        >
          <span v-if="tick.isMajor" class="tick-label">{{ tick.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

interface RulerTick {
  value: number;
  label: string;
  isMajor: boolean;
}

const props = defineProps({
  rulerLength: {
    type: Number,
    default: 30 // 默认标尺长度，单位厘米
  },
  pixelsPerCm: {
    type: Number,
    default: 37.795 // 默认一厘米等于37.795像素
  },
  height: {
    type: Number,
    default: 80 // 标尺高度，像素
  },
  mode: {
    type: String,
    default: 'horizontal', // horizontal or vertical
    validator: (value: string) => ['horizontal', 'vertical'].includes(value)
  }
});

// 计算标尺整体样式
const rulerStyles = computed(() => {
  const width = props.mode === 'horizontal' ? `${props.rulerLength * props.pixelsPerCm}px` : '80px';
  const height = props.mode === 'vertical' ? `${props.rulerLength * props.pixelsPerCm}px` : `${props.height}px`;
  
  return {
    width,
    height,
  };
});

// 生成刻度
const ticks = computed(() => {
  const result: RulerTick[] = [];
  // 计算像素
  const pxPerCm = props.pixelsPerCm;
  
  // 每厘米一个主刻度，每毫米一个小刻度
  for (let cm = 0; cm <= props.rulerLength; cm++) {
    // 添加厘米刻度（主刻度）
    result.push({
      value: cm * pxPerCm,
      label: cm.toString(),
      isMajor: true
    });
    
    // 添加毫米刻度（小刻度），但最后一个厘米标记后不再添加毫米刻度
    if (cm < props.rulerLength) {
      for (let mm = 1; mm < 10; mm++) {
        result.push({
          value: (cm * 10 + mm) * (pxPerCm / 10),
          label: '', // 毫米标记通常不需要标签
          isMajor: false
        });
      }
    }
  }
  
  return result;
});

// 计算每个刻度的样式
const tickStyle = (tick: RulerTick) => {
  if (props.mode === 'horizontal') {
    const height = tick.isMajor ? '18px' : '10px';
    return {
      left: `${tick.value}px`,
      height,
      top: '0px',
    };
  } else {
    const width = tick.isMajor ? '18px' : '10px';
    return {
      top: `${tick.value}px`,
      width,
      left: '0px',
    };
  }
};
</script>

<style scoped>
.static-ruler-container {
  display: flex;
  justify-content: center;
  margin: 20px 0;
}

.ruler {
  position: relative;
  background-color: #f5f5f5;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  overflow: hidden;
}

.ruler-markings {
  position: relative;
  width: 100%;
  height: 100%;
}

.ruler-tick {
  position: absolute;
  background-color: #333;
  width: 1px;
}

.major-tick {
  width: 2px;
  background-color: #000;
}

.tick-label {
  position: absolute;
  font-size: 10px;
  font-family: Arial, sans-serif;
  color: #333;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}

/* Vertical ruler adjustments */
.static-ruler-container:has(.ruler[style*="vertical"]) .ruler-tick {
  height: 1px;
  width: auto;
}

.static-ruler-container:has(.ruler[style*="vertical"]) .major-tick {
  height: 2px;
}

.static-ruler-container:has(.ruler[style*="vertical"]) .tick-label {
  top: 50%;
  left: 20px;
  transform: translateY(-50%);
}
</style>