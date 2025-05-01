<template>
  <div class="mouse-comparison-container">
    <div class="comparison-header">
      <h2 class="text-xl font-bold">Mouse Comparison</h2>
      <div class="comparison-controls">
        <el-radio-group v-model="comparisonMode" size="small">
          <el-radio-button label="overlay">Overlay Comparison</el-radio-button>
          <el-radio-button label="sideBySide">Side-by-side Comparison</el-radio-button>
        </el-radio-group>

        <el-select
          v-if="comparisonMode === 'overlay'"
          v-model="viewType"
          size="small"
          placeholder="Select View"
        >
          <el-option label="Top View" value="topView" />
          <el-option label="Side View" value="sideView" />
        </el-select>

        <el-slider
          v-if="comparisonMode === 'overlay'"
          v-model="overlayOpacity"
          :min="0.2"
          :max="0.8"
          :step="0.05"
          :format-tooltip="(value) => `Transparency: ${Math.round(value * 100)}%`"
          class="w-32"
        />
      </div>
    </div>

    <div v-if="!selectedMice.length" class="empty-state">
      <el-empty description="Please select at least one mouse to compare" />
      <el-button type="primary" @click="openMouseSelector">Select Mouse</el-button>
    </div>

    <div v-else class="comparison-content">
      <div class="svg-comparison-area">
        <!-- Comparison view container -->
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div class="svg-container" v-html="comparisonSvg"></div>

        <!-- Ruler related functions and buttons have been removed -->

        <!-- Loading state -->
        <div v-if="loading" class="loading-overlay">
          <el-icon><Loading /></el-icon>
        </div>
      </div>

      <div class="comparison-details">
        <div>
          <el-button
            v-if="selectedMice.length < 3"
            size="small"
            type="primary"
            plain
            @click="openMouseSelector"
          >
            Add Mouse
          </el-button>
        </div>

        <div v-if="comparisonData" class="specs-comparison mt-6">
          <h3 class="text-lg font-medium mb-3">Parameter Comparison</h3>
          <div class="similarity-score mb-4">
            <div class="font-medium text-sm">Similarity Score</div>
            <el-progress
              :percentage="comparisonData.similarityScore"
              :color="getSimilarityColor(comparisonData.similarityScore)"
              :format="(percent) => `${percent}%`"
              :stroke-width="10"
            />
          </div>

          <el-table :data="comparisonTableData" border stripe>
            <el-table-column prop="property" label="Parameter" width="150" />
            <el-table-column
              v-for="(mouse, index) in selectedMice"
              :key="mouse.id"
              :label="`${mouse.brand} ${mouse.name}`"
              :prop="`values[${index}]`"
              :min-width="120"
            />
          </el-table>
        </div>
      </div>
    </div>

    <!-- Using MouseSelector to replace the original dialog -->
    <el-dialog v-model="mouseDialogVisible" title="Select Mouse" width="60%">
      <MouseSelector
        :initial-selected-mice="selectedMice"
        :max-selection="3"
        @select="handleMouseSelection"
        @cancel="mouseDialogVisible = false"
      />

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="mouseDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="handleDialogConfirm">Confirm</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { DeviceListResponse } from '@/api/device';
import { ref, computed, onMounted, watch } from 'vue';
import { useComparisonStore } from '@/stores';
import type {
  MouseDevice,
  MouseComparisonResult,
  ComparisonMode,
  ViewType
} from '@/models/MouseModel';
import svgService from '@/services/svgService';
import comparisonService from '@/services/comparisonService';
import { getDevices, getMouseSVG } from '@/api/device';
import MouseSelector from '@/components/comparison/MouseSelector.vue';

// State
const comparisonStore = useComparisonStore();
const comparisonMode = ref<ComparisonMode>(comparisonStore.comparisonMode);
const viewType = ref<ViewType>(comparisonStore.viewType);
const overlayOpacity = ref(comparisonStore.overlayOpacity);
const loading = ref(false);
const mouseDialogVisible = ref(false);
const availableMice = ref<MouseDevice[]>([]);
const comparisonData = ref<MouseComparisonResult | null>(null);

// Ruler tool related states have been removed

// Computed properties
const selectedMice = computed(() => comparisonStore.selectedMice);

// Mouse markers (ruler-dependent functionality has been removed)

// No longer need a separate function to load SVG data, now directly fetch comparison results via API
// This function is kept as a backup but is no longer used
const loadSvgData = async () => {
  console.warn('loadSvgData is deprecated, using API comparison instead');
  if (!selectedMice.value.length) return [];
  
  loading.value = true;
  
  try {
    const svgPromises = selectedMice.value.map(async (mouse) => {
      // Get SVG data from API
      try {
        const result = await getMouseSVG(mouse.id, viewType.value as 'top' | 'side');
        if (result && result.data) {
          return result.data.svgData;
        }
      } catch (err) {
        console.error(`Error fetching SVG for mouse ${mouse.id}:`, err);
      }
      
      // If retrieval fails, return placeholder SVG
      console.warn(`Mouse ${mouse.brand} ${mouse.name} is missing ${viewType.value} view SVG data`);
      return '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50" y="50" text-anchor="middle" fill="#999">No SVG data available</text></svg>';
    });
    
    return await Promise.all(svgPromises);
  } catch (error) {
    console.error('Error loading SVG data:', error);
    return [];
  } finally {
    loading.value = false;
  }
};

// Comparison SVG
const comparisonSvg = ref('');

// Update comparison SVG
const updateComparisonSvg = async () => {
  console.log('Starting updateComparisonSvg method', { 
    selectedMice: selectedMice.value.length, 
    viewType: viewType.value,
    comparisonMode: comparisonMode.value
  });
  
  // Set loading state at the start
  loading.value = true;
  
  if (!selectedMice.value.length) {
    console.log('No mice selected, clearing SVG');
    comparisonSvg.value = '';
    loading.value = false;
    return;
  }
  
  try {
    // Use new backend API for SVG comparison
    const deviceIds = selectedMice.value.map(mouse => mouse.id);
    console.log('Device IDs for comparison:', deviceIds);
    
    // Check if we have enough devices for comparison
    if (deviceIds.length < 2) {
      console.warn('Need at least 2 devices for comparison, only have', deviceIds.length);
      comparisonSvg.value = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50" y="50" text-anchor="middle" fill="#E6A23C">Please select at least two mice for comparison</text></svg>';
      return;
    }
    
    // Ensure view parameter is correctly passed as 'top' or 'side'
    // topView -> top, sideView -> side
    const view: 'top' | 'side' = viewType.value === 'topView' ? 'top' : 'side';
    console.log('Using view for comparison:', view);
    
    // Check for hardcoded SVG data in mouse objects first
    const hasSvgData = selectedMice.value.some(mouse => mouse.svgData && mouse.svgData[view]);
    console.log('Found hardcoded SVG data in mouse objects:', hasSvgData);
    
    if (comparisonMode.value === 'overlay') {
      const opacities = selectedMice.value.map((_, index) =>
        index === 0 ? 1.0 : overlayOpacity.value
      );
      const colors = ['#000000', '#FF0000', '#0000FF'];
      
      console.log('Overlay mode with opacities:', opacities, 'and colors:', colors.slice(0, selectedMice.value.length));
      
      try {
        // First check if mice have SVG data in their objects
        if (hasSvgData) {
          console.log('Using hardcoded SVG data from mouse objects');
          // Implement using hardcoded SVG data
          // This is a placeholder - you would need to implement how to use the hardcoded data
          throw new Error('Hardcoded SVG implementation not complete, falling back to API');
        }
        
        // Use API for overlay comparison if we have at least 2 devices
        if (deviceIds.length >= 2) {
          console.log('Using API for overlay SVG comparison');
          const result = await svgService.createOverlaySvg(
            deviceIds,
            view,
            opacities,
            colors.slice(0, selectedMice.value.length)
          );
          console.log('Successfully received overlay SVG from API');
          comparisonSvg.value = result;
        } else {
          throw new Error('At least two devices are needed for comparison');
        }
      } catch (apiError) {
        console.warn('Failed to get SVG from API, trying to use local SVG files:', apiError);
        // Fallback to local SVG files
        console.log('Falling back to local SVG files');
        const svgsPromises = selectedMice.value.map(mouse => {
          const fileName = mouse.id.replace(/-/g, ' ');
          const svgPath = `/svg/${fileName} ${view}.svg`;
          console.log(`Attempting to fetch SVG from: ${svgPath}`);
          
          return fetch(svgPath)
            .then(response => {
              if (!response.ok) {
                console.warn(`SVG fetch failed for ${mouse.id}: ${response.status}`);
                throw new Error(`Failed to load SVG: ${response.status}`);
              }
              return response.text();
            })
            .catch(err => {
              console.error(`Failed to load SVG: ${mouse.id}`, err);
              return '';
            });
        });
        
        const svgs = await Promise.all(svgsPromises);
        console.log(`Retrieved ${svgs.filter(svg => svg).length}/${svgs.length} SVGs successfully`);
        
        if (svgs.some(svg => svg)) {
          console.log('Creating overlay SVG from local files');
          const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
          svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
          svgContainer.setAttribute('viewBox', '0 0 1250 400');
          
          svgs.forEach((svgContent, index) => {
            if (!svgContent) return;
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = svgContent;
            const svgElement = tempDiv.querySelector('svg');
            const pathElement = svgElement?.querySelector('path');
            
            if (pathElement) {
              console.log(`Adding mouse ${index + 1} to overlay SVG`);
              const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
              g.setAttribute('opacity', opacities[index].toString());
              g.setAttribute('fill', 'none');
              g.setAttribute('stroke', colors[index]);
              g.setAttribute('stroke-width', '2');
              
              const clonedPath = pathElement.cloneNode(true);
              g.appendChild(clonedPath);
              svgContainer.appendChild(g);
            } else {
              console.warn(`No path element found in SVG for mouse ${index + 1}`);
            }
          });
          
          comparisonSvg.value = svgContainer.outerHTML;
          console.log('Successfully created overlay SVG from local files');
        } else {
          console.error('No SVG files could be loaded for any selected mice');
          throw new Error('Failed to load SVG files');
        }
      }
    } else {
      // Side by side comparison
      console.log('Using side-by-side comparison mode');
      
      try {
        // First check if mice have SVG data in their objects
        if (hasSvgData) {
          console.log('Using hardcoded SVG data from mouse objects for side-by-side view');
          // Implement using hardcoded SVG data
          // This is a placeholder - you would need to implement how to use the hardcoded data
          throw new Error('Hardcoded SVG implementation not complete, falling back to API');
        }
        
        // Use API for side-by-side comparison if we have at least 2 devices
        if (deviceIds.length >= 2) {
          console.log('Using API for side-by-side SVG comparison');
          const result = await svgService.createSideBySideSvg(
            deviceIds,
            view
          );
          console.log('Successfully received side-by-side SVG from API');
          comparisonSvg.value = result;
        } else {
          throw new Error('At least two devices are needed for comparison');
        }
      } catch (apiError) {
        console.warn('Failed to get side-by-side SVG from API, trying to use local SVG files:', apiError);
        console.log('Falling back to local SVG files for side-by-side view');
        
        // Fallback to local SVG files
        const svgsPromises = selectedMice.value.map(mouse => {
          const fileName = mouse.id.replace(/-/g, ' ');
          const svgPath = `/svg/${fileName} ${view}.svg`;
          console.log(`Attempting to fetch SVG from: ${svgPath}`);
          
          return fetch(svgPath)
            .then(response => {
              if (!response.ok) {
                console.warn(`SVG fetch failed for ${mouse.id}: ${response.status}`);
                throw new Error(`Failed to load SVG: ${response.status}`);
              }
              return response.text();
            })
            .catch(err => {
              console.error(`Failed to load SVG: ${mouse.id}`, err);
              return '';
            });
        });
        
        const svgs = await Promise.all(svgsPromises);
        console.log(`Retrieved ${svgs.filter(svg => svg).length}/${svgs.length} SVGs successfully for side-by-side view`);
        
        if (svgs.some(svg => svg)) {
          console.log('Creating side-by-side SVG from local files');
          const width = 1250;
          const height = 400;
          const padding = 20;
          const mouseWidth = (width - (padding * (selectedMice.value.length + 1))) / selectedMice.value.length;
          
          const svgContainer = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
          svgContainer.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
          svgContainer.setAttribute('viewBox', `0 0 ${width} ${height}`);
          
          svgs.forEach((svgContent, index) => {
            if (!svgContent) return;
            
            const tempDiv = document.createElement('div');
            tempDiv.innerHTML = svgContent;
            const svgElement = tempDiv.querySelector('svg');
            const pathElement = svgElement?.querySelector('path');
            
            if (pathElement) {
              console.log(`Adding mouse ${index + 1} to side-by-side SVG`);
              const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
              g.setAttribute('transform', `translate(${padding + (mouseWidth + padding) * index}, 0) scale(${mouseWidth/width})`);
              g.setAttribute('fill', 'none');
              g.setAttribute('stroke', '#000');
              g.setAttribute('stroke-width', '2');
              
              const clonedPath = pathElement.cloneNode(true);
              g.appendChild(clonedPath);
              
              // Add mouse name
              const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
              text.setAttribute('x', (mouseWidth / 2).toString());
              text.setAttribute('y', (height - 20).toString());
              text.setAttribute('text-anchor', 'middle');
              text.setAttribute('font-family', 'Arial, sans-serif');
              text.setAttribute('font-size', '14');
              text.textContent = `${selectedMice.value[index].brand} ${selectedMice.value[index].name}`;
              g.appendChild(text);
              
              svgContainer.appendChild(g);
            } else {
              console.warn(`No path element found in SVG for mouse ${index + 1}`);
            }
          });
          
          comparisonSvg.value = svgContainer.outerHTML;
          console.log('Successfully created side-by-side SVG from local files');
        } else {
          console.error('No SVG files could be loaded for any selected mice');
          throw new Error('Failed to load SVG files');
        }
      }
    }
  } catch (error) {
    console.error('Error creating comparison SVG:', error);
    comparisonSvg.value = '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text x="50" y="50" text-anchor="middle" fill="#f56c6c">Failed to render SVG comparison</text></svg>';
  } finally {
    // Make sure to reset loading state when done
    loading.value = false;
    console.log('Completed updateComparisonSvg method');
  }
};

// Comparison table data
const comparisonTableData = computed(() => {
  if (!comparisonData.value) return [];

  return Object.values(comparisonData.value.differences).sort((a, b) => {
    // Sort by difference percentage in descending order
    return b.differencePercent - a.differencePercent;
  });
});

// Methods
function openMouseSelector() {
  mouseDialogVisible.value = true;
}

function handleMouseSelection(_selectedMice) {
  // Handle selection through MouseSelector component
  updateComparisonData();
  updateComparisonSvg();
}

function handleDialogConfirm() {
  updateComparisonData();
  updateComparisonSvg();
  mouseDialogVisible.value = false;
}

function _removeMouse(mouseId: string) {
  comparisonStore.removeMouse(mouseId);
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
  } else {
    updateComparisonData();
  }
  updateComparisonSvg();
}

async function fetchAvailableMice() {
  loading.value = true;
  try {
    const response = await getDevices({ type: 'mouse' });
    if (!response || !response.data) {
      throw new Error('Failed to fetch mouse data: No data returned');
    }
    availableMice.value = response.data.devices as MouseDevice[];
  } catch (error) {
    console.error('Error fetching mice:', error);
    availableMice.value = []; // Set as empty array to avoid null reference
  } finally {
    loading.value = false;
  }
}

function updateComparisonData() {
  if (selectedMice.value.length < 2) {
    comparisonData.value = null;
    return;
  }

  // @ts-expect-error - Type inconsistency between MouseDevice definitions
  comparisonData.value = comparisonService.generateComparisonResult(selectedMice.value);
}

// Format difference percentage
function _formatDifference(value: number) {
  return value === 0 ? 'Identical' : `${value.toFixed(1)}%`;
}

// Get difference display style
function _getDifferenceClass(value: number) {
  if (value === 0) return 'text-green-500';
  if (value < 10) return 'text-blue-500';
  if (value < 25) return 'text-amber-500';
  return 'text-red-500';
}

// Get similarity score color
function getSimilarityColor(score: number) {
  if (score >= 90) return '#67C23A'; // Green
  if (score >= 75) return '#409EFF'; // Blue
  if (score >= 50) return '#E6A23C'; // Orange
  return '#F56C6C'; // Red
}

// Ruler related functionality has been removed

// Watch comparison mode and opacity changes
watch(comparisonMode, (newMode) => {
  comparisonStore.setComparisonMode(newMode);
  updateComparisonSvg();
});

watch(overlayOpacity, (newOpacity) => {
  comparisonStore.setOverlayOpacity(newOpacity);
  updateComparisonSvg();
});

watch(viewType, () => {
  updateComparisonSvg();
});

watch(selectedMice, () => {
  updateComparisonSvg();
}, { deep: true });

// Lifecycle hooks
onMounted(() => {
  fetchAvailableMice();
  updateComparisonData();
  updateComparisonSvg();
});
</script>

<style scoped>
.mouse-comparison-container {
  padding: 1rem;
}

.comparison-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.comparison-controls {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  gap: 1rem;
}

.comparison-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

@media (min-width: 1024px) {
  .comparison-content {
    grid-template-columns: 1fr 1fr;
  }
}

.svg-comparison-area {
  position: relative;
  background-color: #f9f9f9;
  border-radius: 0.5rem;
  padding: 1rem;
  min-height: 400px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.svg-container {
  max-width: 100%;
  height: auto;
  width: 100%;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 10px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
}

.mouse-card.selected {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.text-primary {
  color: #409eff;
}
</style>
