<template>
  <div class="scale-ruler" :style="{ width: `${customWidth}px` }">
    <div class="ruler-container">
      <div class="ruler-background">
        <div v-for="tick in ticks" :key="tick.value" class="ruler-tick" :style="tickStyle(tick)">
          <span v-if="tick.label" class="tick-label">{{ tick.label }}</span>
        </div>
      </div>
      <div v-if="showMarkers" class="ruler-markers">
        <div
          v-for="(marker, index) in markers"
          :key="index"
          class="ruler-marker"
          :style="{ left: `${marker.position}px` }"
          :title="marker.label"
        >
          <div class="marker-line" :style="{ backgroundColor: marker.color }"></div>
          <div class="marker-label" :style="{ color: marker.color }">{{ marker.label }}</div>
        </div>
      </div>
    </div>
    <div class="ruler-controls mb-3">
      <div class="unit-controls">
        <select v-model="selectedUnit" class="unit-select">
          <option value="mm">mm</option>
          <option value="cm">cm</option>
          <option value="inch">inch</option>
        </select>
      </div>
      <div class="zoom-controls">
        <button @click="zoomIn" class="zoom-button">+</button>
        <button @click="zoomOut" class="zoom-button">-</button>
      </div>
    </div>

    <div class="ruler-settings">
      <div class="length-setting">
        <label for="ruler-length" class="text-sm text-white">尺子长度 (像素)</label>
        <div class="flex gap-2">
          <input
            id="ruler-length"
            type="number"
            v-model.number="customWidth"
            min="200"
            max="1200"
            step="50"
            class="length-input"
          />
          <button @click="resetSize" class="reset-button">重置</button>
        </div>
      </div>

      <div class="real-size-display">
        <div class="text-xs text-gray-400 mt-1">实际长度: {{ realLength }}</div>
      </div>
    </div>

    <div class="custom-marker-section mt-4" v-if="markers.length === 0">
      <div class="flex gap-2 mb-2">
        <input
          type="number"
          v-model.number="newMarkerPosition"
          placeholder="标记位置(px)"
          class="marker-input"
        />
        <input type="text" v-model="newMarkerLabel" placeholder="标记名称" class="marker-input" />
        <button @click="addMarker" class="add-marker-button">添加</button>
      </div>
      <div v-if="customMarkers.length > 0" class="custom-markers-list">
        <div
          v-for="(marker, index) in customMarkers"
          :key="index"
          class="flex justify-between items-center p-1 border-b border-gray-600"
        >
          <span class="text-sm">{{ marker.label }}: {{ marker.position }}px</span>
          <button @click="removeMarker(index)" class="text-red-400 text-xs">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, PropType, watch } from 'vue';

interface Tick {
  value: number;
  label?: string;
}

interface Marker {
  position: number;
  label: string;
  color: string;
}

export default defineComponent({
  name: 'ScaleRuler',
  props: {
    width: {
      type: Number,
      default: 400
    },
    initialUnit: {
      type: String,
      default: 'mm'
    },
    markers: {
      type: Array as PropType<Marker[]>,
      default: () => []
    },
    showMarkers: {
      type: Boolean,
      default: true
    }
  },
  setup(props, { emit }) {
    const selectedUnit = ref(props.initialUnit);
    const zoomLevel = ref(1);
    const customWidth = ref(props.width);
    const customMarkers = ref<Marker[]>([]);
    const newMarkerPosition = ref<number>(0);
    const newMarkerLabel = ref('');

    // Calculate the ticks based on selected unit and zoom level
    const ticks = computed(() => {
      const result: Tick[] = [];

      let step = 0;
      let majorStep = 0;
      let unitConversion = 0;

      switch (selectedUnit.value) {
        case 'mm':
          step = 5 * zoomLevel.value;
          majorStep = 10;
          unitConversion = 1;
          break;
        case 'cm':
          step = 10 * zoomLevel.value;
          majorStep = 1;
          unitConversion = 0.1;
          break;
        case 'inch':
          step = 25.4 * zoomLevel.value;
          majorStep = 0.5;
          unitConversion = 1 / 25.4;
          break;
      }

      const numTicks = Math.floor(customWidth.value / step);

      for (let i = 0; i <= numTicks; i++) {
        const value = i * step;
        const unitValue = ((i * step) / zoomLevel.value) * unitConversion;
        const isMajor = Math.abs(Math.round(unitValue / majorStep) * majorStep - unitValue) < 0.001;

        result.push({
          value,
          label: isMajor ? unitValue.toFixed(1) : undefined
        });
      }

      return result;
    });

    const realLength = computed(() => {
      let unit = '';
      let length = 0;

      switch (selectedUnit.value) {
        case 'mm':
          unit = '毫米';
          length = (customWidth.value / (5 * zoomLevel.value)) * 10;
          break;
        case 'cm':
          unit = '厘米';
          length = (customWidth.value / (10 * zoomLevel.value)) * 1;
          break;
        case 'inch':
          unit = '英寸';
          length = (customWidth.value / (25.4 * zoomLevel.value)) * 1;
          break;
      }

      return `${length.toFixed(1)} ${unit}`;
    });

    const tickStyle = (tick: Tick) => {
      const height = tick.label ? '15px' : '8px';
      return {
        left: `${tick.value}px`,
        height
      };
    };

    const zoomIn = () => {
      if (zoomLevel.value < 2) {
        zoomLevel.value *= 1.25;
      }
    };

    const zoomOut = () => {
      if (zoomLevel.value > 0.5) {
        zoomLevel.value /= 1.25;
      }
    };

    const resetSize = () => {
      customWidth.value = props.width;
    };

    const addMarker = () => {
      if (newMarkerPosition.value <= 0 || newMarkerPosition.value >= customWidth.value) {
        return;
      }

      if (!newMarkerLabel.value) {
        newMarkerLabel.value = `标记 ${customMarkers.value.length + 1}`;
      }

      customMarkers.value.push({
        position: newMarkerPosition.value,
        label: newMarkerLabel.value,
        color: '#ff5252'
      });

      newMarkerPosition.value = 0;
      newMarkerLabel.value = '';
    };

    const removeMarker = (index: number) => {
      customMarkers.value.splice(index, 1);
    };

    return {
      selectedUnit,
      ticks,
      tickStyle,
      zoomIn,
      zoomOut,
      showMarkers: props.showMarkers,
      markers: computed(() => {
        // 如果有外部传入的marker，优先使用外部的
        return props.markers.length > 0 ? props.markers : customMarkers.value;
      }),
      customWidth,
      resetSize,
      realLength,
      customMarkers,
      newMarkerPosition,
      newMarkerLabel,
      addMarker,
      removeMarker
    };
  }
});
</script>

<style scoped>
.scale-ruler {
  border: 1px solid #555;
  border-radius: 4px;
  padding: 8px;
  background-color: rgba(40, 40, 40, 0.8);
  margin: 10px 0;
}

.ruler-container {
  position: relative;
  height: 40px;
  margin-bottom: 8px;
  overflow: hidden;
}

.ruler-background {
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
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: white;
}

.ruler-markers {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  pointer-events: none;
}

.ruler-marker {
  position: absolute;
  top: 0;
  height: 100%;
}

.marker-line {
  position: absolute;
  width: 2px;
  height: 100%;
  left: 0;
  background-color: #ff5252;
}

.marker-label {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  font-size: 10px;
  color: #ff5252;
  background-color: rgba(40, 40, 40, 0.7);
  padding: 2px 4px;
  border-radius: 2px;
}

.ruler-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.unit-select {
  background-color: #444;
  color: white;
  border: 1px solid #666;
  border-radius: 3px;
  padding: 2px 4px;
}

.zoom-controls {
  display: flex;
  gap: 8px;
}

.zoom-button {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #555;
  color: white;
  border: none;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.zoom-button:hover {
  background-color: #666;
}

.ruler-settings {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.length-setting {
  display: flex;
  flex-direction: column;
}

.length-input {
  background-color: #444;
  color: white;
  border: 1px solid #666;
  border-radius: 3px;
  padding: 2px 4px;
  width: 100px;
}

.reset-button {
  background-color: #555;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 0 8px;
  font-size: 12px;
  cursor: pointer;
}

.reset-button:hover {
  background-color: #666;
}

.marker-input {
  background-color: #444;
  color: white;
  border: 1px solid #666;
  border-radius: 3px;
  padding: 2px 4px;
  font-size: 12px;
}

.add-marker-button {
  background-color: #5252ff;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 0 8px;
  font-size: 12px;
  cursor: pointer;
}

.add-marker-button:hover {
  background-color: #6565ff;
}

.custom-markers-list {
  max-height: 100px;
  overflow-y: auto;
  background-color: rgba(60, 60, 60, 0.5);
  border-radius: 3px;
  margin-top: 4px;
  font-size: 12px;
}
</style>
