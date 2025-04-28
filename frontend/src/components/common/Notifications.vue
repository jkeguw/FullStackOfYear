<template>
  <transition name="fade">
    <div
      v-if="show"
      class="fixed top-16 right-4 z-50 max-w-md shadow-lg rounded-lg overflow-hidden"
      :class="typeClasses"
    >
      <div class="p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <i :class="iconClass" class="text-xl"></i>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium" :class="textColorClass">
              {{ title }}
            </h3>
            <div v-if="message" class="mt-2 text-sm" :class="textClass">
              {{ message }}
            </div>
          </div>
        </div>
      </div>

      <div v-if="actions && actions.length" class="border-t border-gray-200 bg-gray-50 px-4 py-3">
        <div class="flex justify-end space-x-3">
          <button
            v-for="(action, index) in actions"
            :key="index"
            @click="action.onClick ? action.onClick() : close()"
            class="px-3 py-1 rounded-md text-sm font-medium"
            :class="action.primary ? primaryButtonClass : secondaryButtonClass"
          >
            {{ action.label }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';

interface Action {
  label: string;
  onClick?: () => void;
  primary?: boolean;
}

interface Props {
  type?: 'success' | 'info' | 'warning' | 'error';
  title: string;
  message?: string;
  duration?: number;
  actions?: Action[];
  autoClose?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'info',
  duration: 5000,
  autoClose: true,
  actions: () => []
});

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const show = ref(true);
let timer: number | null = null;

// 类型相关样式
const typeClasses = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-50 border-l-4 border-green-400';
    case 'warning':
      return 'bg-yellow-50 border-l-4 border-yellow-400';
    case 'error':
      return 'bg-red-50 border-l-4 border-red-400';
    case 'info':
    default:
      return 'bg-blue-50 border-l-4 border-blue-400';
  }
});

const textColorClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-800';
    case 'warning':
      return 'text-yellow-800';
    case 'error':
      return 'text-red-800';
    case 'info':
    default:
      return 'text-blue-800';
  }
});

const textClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'text-green-700';
    case 'warning':
      return 'text-yellow-700';
    case 'error':
      return 'text-red-700';
    case 'info':
    default:
      return 'text-blue-700';
  }
});

const iconClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'el-icon-circle-check text-green-400';
    case 'warning':
      return 'el-icon-warning text-yellow-400';
    case 'error':
      return 'el-icon-circle-close text-red-400';
    case 'info':
    default:
      return 'el-icon-info text-blue-400';
  }
});

const primaryButtonClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-green-600 hover:bg-green-700 text-white';
    case 'warning':
      return 'bg-yellow-600 hover:bg-yellow-700 text-white';
    case 'error':
      return 'bg-red-600 hover:bg-red-700 text-white';
    case 'info':
    default:
      return 'bg-blue-600 hover:bg-blue-700 text-white';
  }
});

const secondaryButtonClass = 'bg-white hover:bg-gray-50 text-gray-700';

// 自动关闭
onMounted(() => {
  if (props.autoClose) {
    timer = window.setTimeout(() => {
      close();
    }, props.duration);
  }
});

onBeforeUnmount(() => {
  if (timer) {
    clearTimeout(timer);
  }
});

// 关闭通知
const close = () => {
  show.value = false;
  setTimeout(() => {
    emit('close');
  }, 300); // 等待动画完成
};
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
