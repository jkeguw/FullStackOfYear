<template>
  <div
    ref="rulerContainer"
    class="simple-ruler"
    :style="{
      left: `${position.x}px`,
      top: `${position.y}px`
    }"
    @mousedown="startDrag"
  >
    <div class="ruler-header">
      <div class="handle" @mousedown.stop="startDrag"></div>
      <div class="length-control">
        <span class="length-label">长度:</span>
        <input type="range" min="10" max="50" v-model="rulerLength" class="length-slider" />
        <span>{{ rulerLength }}cm</span>
      </div>
    </div>
    <div class="ruler-body">
      <div class="ruler-scale" :style="{ width: `${rulerLength * 37.8}px` }">
        <div v-for="tick in ticks" :key="tick.value" class="ruler-tick" :style="tickStyle(tick)">
          <span v-if="tick.label" class="tick-label">{{ tick.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue';

interface Tick {
  value: number;
  label?: string;
}

export default defineComponent({
  name: 'SimpleRuler',
  props: {
    initialPosition: {
      type: Object,
      default: () => ({ x: 100, y: 100 })
    }
  },
  setup(props) {
    const rulerContainer = ref<HTMLElement | null>(null);
    const position = ref({ x: props.initialPosition.x, y: props.initialPosition.y });
    const rulerLength = ref(30); // 默认长度为30厘米
    const isDragging = ref(false);
    const dragOffset = ref({ x: 0, y: 0 });

    // 计算刻度
    const ticks = computed(() => {
      const result: Tick[] = [];
      // 1cm = 37.8px (这是一个近似值，实际应该根据显示器DPI校准)
      const pixelsPerCm = 37.8;
      const totalLength = rulerLength.value * pixelsPerCm;

      // 每厘米生成一个主刻度和9个小刻度
      for (let cm = 0; cm <= rulerLength.value; cm++) {
        // 添加厘米刻度（主刻度）
        result.push({
          value: cm * pixelsPerCm,
          label: cm.toString()
        });

        // 添加毫米刻度（小刻度）
        if (cm < rulerLength.value) {
          for (let mm = 1; mm < 10; mm++) {
            result.push({
              value: (cm * 10 + mm) * (pixelsPerCm / 10),
              label: undefined
            });
          }
        }
      }

      return result;
    });

    const tickStyle = (tick: Tick) => {
      // 厘米刻度更长，毫米刻度更短
      const height = tick.label ? '15px' : '8px';
      return {
        left: `${tick.value}px`,
        height
      };
    };

    // 拖动事件处理
    const startDrag = (event: MouseEvent) => {
      // 检查点击是否在刻度尺区域内
      const target = event.target as HTMLElement;
      const isScaleArea =
        target.classList.contains('ruler-scale') ||
        target.classList.contains('ruler-tick') ||
        target.classList.contains('tick-label') ||
        target.classList.contains('handle');

      if (!isScaleArea && !target.closest('.handle')) {
        return; // 如果不是刻度区域或拖动手柄，不启动拖动
      }

      isDragging.value = true;
      if (rulerContainer.value) {
        const rect = rulerContainer.value.getBoundingClientRect();
        dragOffset.value = {
          x: event.clientX - rect.left,
          y: event.clientY - rect.top
        };
      }
      document.addEventListener('mousemove', onDrag);
      document.addEventListener('mouseup', stopDrag);
    };

    const onDrag = (event: MouseEvent) => {
      if (isDragging.value) {
        // 计算新位置
        let newX = event.clientX - dragOffset.value.x;
        let newY = event.clientY - dragOffset.value.y;

        // 获取窗口边界
        const windowWidth = window.innerWidth;
        const windowHeight = window.innerHeight;

        // 获取尺子尺寸
        const rulerWidth = rulerContainer.value?.offsetWidth || 0;
        const rulerHeight = rulerContainer.value?.offsetHeight || 0;

        // 防止尺子移出视窗
        newX = Math.max(0, Math.min(newX, windowWidth - rulerWidth));
        newY = Math.max(0, Math.min(newY, windowHeight - rulerHeight));

        position.value = {
          x: newX,
          y: newY
        };
      }
    };

    const stopDrag = () => {
      isDragging.value = false;
      document.removeEventListener('mousemove', onDrag);
      document.removeEventListener('mouseup', stopDrag);
    };

    // 清理事件监听器
    onUnmounted(() => {
      document.removeEventListener('mousemove', onDrag);
      document.removeEventListener('mouseup', stopDrag);
    });

    return {
      rulerContainer,
      position,
      rulerLength,
      ticks,
      tickStyle,
      startDrag
    };
  }
});
</script>

<style scoped>
.simple-ruler {
  position: absolute;
  z-index: 1000;
  background-color: rgba(40, 40, 40, 0.8);
  border-radius: 4px;
  cursor: move;
  user-select: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.ruler-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px;
  border-bottom: 1px solid #555;
}

.handle {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #666;
  cursor: grab;
}

.handle:active {
  cursor: grabbing;
  background-color: #888;
}

.length-control {
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-size: 12px;
}

.length-slider {
  width: 80px;
  height: 4px;
}

.ruler-body {
  position: relative;
  padding: 4px 0;
}

.ruler-scale {
  position: relative;
  height: 40px;
  background-color: rgba(70, 70, 70, 0.5);
  border-radius: 2px;
}

.ruler-tick {
  position: absolute;
  bottom: 0;
  width: 1px;
  background-color: white;
}

.tick-label {
  position: absolute;
  top: -16px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: white;
}
</style>
