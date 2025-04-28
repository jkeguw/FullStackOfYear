<template>
  <div
    ref="rulerRef"
    class="draggable-ruler"
    :style="{
      left: `${position.x}px`,
      top: `${position.y}px`,
      display: visible ? 'block' : 'none'
    }"
  >
    <div class="ruler-header">
      <div class="handle"></div>
      <div class="length-control">
        <span class="length-label">长度:</span>
        <input type="range" min="10" max="50" v-model="rulerLength" class="length-slider" />
        <span>{{ rulerLength }}cm</span>
      </div>
      <div class="close-button" @click="hideRuler">×</div>
    </div>
    <div class="ruler-body">
      <div class="ruler-scale" :style="{ width: `${rulerLength * pixelsPerCm}px` }">
        <div v-for="tick in ticks" :key="tick.value" class="ruler-tick" :style="tickStyle(tick)">
          <span v-if="tick.label" class="tick-label">{{ tick.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted, onUnmounted, watch } from 'vue';

interface Tick {
  value: number;
  label?: string;
}

export default defineComponent({
  name: 'DraggableRuler',
  props: {
    initialPosition: {
      type: Object,
      default: () => ({ x: 100, y: 100 })
    },
    visible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:visible'],
  setup(props, { emit }) {
    const rulerRef = ref<HTMLElement | null>(null);
    const position = ref({ x: props.initialPosition.x, y: props.initialPosition.y });
    const rulerLength = ref(30); // 默认长度为30厘米
    const isDragging = ref(false);
    const dragOffset = ref({ x: 0, y: 0 });
    const pixelsPerCm = 37.8; // 1cm = 37.8px (校准值，可能需要根据显示器DPI调整)

    // 计算刻度
    const ticks = computed(() => {
      const result: Tick[] = [];
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
      // 检查点击是否在刻度尺区域内或是拖动手柄
      const target = event.target as HTMLElement;
      const isScaleArea =
        target.classList.contains('ruler-scale') ||
        target.classList.contains('ruler-tick') ||
        target.classList.contains('tick-label');
      const isHandle = target.classList.contains('handle');

      if (!isScaleArea && !isHandle) {
        return; // 如果不是刻度区域或拖动手柄，不启动拖动
      }

      isDragging.value = true;
      if (rulerRef.value) {
        const rect = rulerRef.value.getBoundingClientRect();
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
        const rulerWidth = rulerRef.value?.offsetWidth || 0;
        const rulerHeight = rulerRef.value?.offsetHeight || 0;

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

    const hideRuler = () => {
      emit('update:visible', false);
    };

    // 清理事件监听器
    onUnmounted(() => {
      document.removeEventListener('mousemove', onDrag);
      document.removeEventListener('mouseup', stopDrag);
    });

    return {
      rulerRef,
      position,
      rulerLength,
      ticks,
      tickStyle,
      startDrag,
      hideRuler,
      pixelsPerCm
    };
  }
});
</script>

<style scoped>
.draggable-ruler {
  position: fixed;
  z-index: 9999;
  background-color: rgba(40, 40, 40, 0.9);
  border-radius: 4px;
  cursor: move;
  user-select: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(80, 80, 80, 0.8);
}

.ruler-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  border-bottom: 1px solid #555;
  background-color: rgba(50, 50, 50, 0.95);
  border-radius: 4px 4px 0 0;
}

.handle {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #666;
  cursor: grab;
  display: flex;
  align-items: center;
  justify-content: center;
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

.close-button {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(255, 75, 75, 0.7);
  color: white;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.close-button:hover {
  background-color: rgba(255, 75, 75, 0.9);
}

.ruler-body {
  position: relative;
  padding: 8px 0;
  overflow: hidden;
}

.ruler-scale {
  position: relative;
  height: 40px;
  background-color: rgba(70, 70, 70, 0.7);
  border-radius: 0 0 2px 2px;
}

.ruler-tick {
  position: absolute;
  bottom: 0;
  width: 1px;
  background-color: white;
}

.tick-label {
  position: absolute;
  top: -18px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: white;
  white-space: nowrap;
}
</style>
