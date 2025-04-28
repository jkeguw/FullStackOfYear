import { defineStore } from 'pinia';
import type { User, UserState } from '@/types/user';
import type { MouseDevice, ComparisonState, ComparisonMode, ViewType } from '@/models/MouseModel';

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null,
    profile: null,
    gameSettings: null,
    settings: null,
    privacy: null,
    devicePreferences: []
  }),
  actions: {
    setUser(user: User) {
      this.user = user;
    },
    setToken(token: string) {
      this.token = token;
    },
    clearUser() {
      this.user = null;
      this.token = null;
    }
  }
});

// 导出类型供其他文件使用
export type { MouseDevice, ComparisonState, ComparisonMode, ViewType };

export const useComparisonStore = defineStore('comparison', {
  state: (): ComparisonState => ({
    selectedMice: [],
    comparisonMode: 'overlay',
    viewType: 'topView',
    overlayOpacity: 0.5,
    recentlyViewedMice: []
  }),
  actions: {
    addMouse(mouse: MouseDevice) {
      // 限制最多3个鼠标
      if (this.selectedMice.length >= 3) {
        // 如果已经有3个鼠标，移除第一个
        this.selectedMice.shift();
      }

      if (this.selectedMice.findIndex((m) => m.id === mouse.id) === -1) {
        this.selectedMice.push(mouse);
      }

      // 同时添加到最近浏览
      this.addToRecentlyViewed(mouse);
    },
    removeMouse(mouseId: string) {
      this.selectedMice = this.selectedMice.filter((m) => m.id !== mouseId);
    },
    clearSelection() {
      this.selectedMice = [];
    },
    setComparisonMode(mode: ComparisonMode) {
      this.comparisonMode = mode;
    },
    setViewType(type: ViewType) {
      this.viewType = type;
    },
    setOverlayOpacity(opacity: number) {
      this.overlayOpacity = opacity;
    },
    isMouseSelected(mouseId: string) {
      return this.selectedMice.some((m) => m.id === mouseId);
    },
    addToRecentlyViewed(mouse: MouseDevice) {
      // 移除已存在的同一鼠标（如果有）
      this.recentlyViewedMice = this.recentlyViewedMice.filter((m) => m.id !== mouse.id);
      // 在开头添加鼠标
      this.recentlyViewedMice.unshift(mouse);
      // 保持最多10个记录
      if (this.recentlyViewedMice.length > 10) {
        this.recentlyViewedMice.pop();
      }
    }
  }
});

// 购物车状态仓库
export interface CartItem {
  productId: string;
  mouseId: string;
  name: string;
  price: number;
  quantity: number;
  image: string;
}

export interface CartState {
  items: CartItem[];
}

export const useCartStore = defineStore('cart', {
  state: (): CartState => ({
    items: []
  }),
  getters: {
    totalItems: (state) => state.items.reduce((total, item) => total + item.quantity, 0),
    totalPrice: (state) =>
      state.items.reduce((total, item) => total + item.price * item.quantity, 0)
  },
  actions: {
    addItem(item: CartItem) {
      const existingItem = this.items.find((i) => i.productId === item.productId);
      if (existingItem) {
        existingItem.quantity += item.quantity;
      } else {
        this.items.push(item);
      }
    },
    removeItem(productId: string) {
      this.items = this.items.filter((item) => item.productId !== productId);
    },
    updateQuantity(productId: string, quantity: number) {
      const item = this.items.find((i) => i.productId === productId);
      if (item) {
        item.quantity = quantity;
      }
    },
    clearCart() {
      this.items = [];
    }
  }
});
