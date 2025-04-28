#!/bin/bash

# Comprehensive TypeScript error fix script
# This script addresses the TypeScript errors shown in sa.txt

# 1. Create backup of original files
echo "Creating backup directory..."
BACKUP_DIR="./frontend/src_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR
cp -r ./frontend/src $BACKUP_DIR

# 2. Fix property naming inconsistencies (snake_case to camelCase)
echo "Fixing property naming inconsistencies..."
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs sed -i 's/image_url/imageUrl/g'
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs sed -i 's/product_id/productId/g'
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs sed -i 's/product_type/productType/g'
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs sed -i 's/item_count/itemCount/g'
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs sed -i 's/updated_at/updatedAt/g'

# 3. Fix Response object unwrapping in composables
echo "Fixing Response object unwrapping..."
if grep -q "Type 'AxiosResponse<CartResponse, any>' is not assignable to type" ./frontend/src/composables/useCart.ts; then
  sed -i 's/return res;/return res.data;/g' ./frontend/src/composables/useCart.ts
fi

if grep -q "Type 'Response<UserProfile>' is not assignable to type" ./frontend/src/composables/useProfile.ts; then
  sed -i 's/return res;/return res.data;/g' ./frontend/src/composables/useProfile.ts
fi

# 4. Fix missing DeviceListResponse type
echo "Adding DeviceListResponse import to MouseComparisonView..."
if grep -q "Cannot find name 'DeviceListResponse'" ./frontend/src/components/comparison/MouseComparisonView.vue; then
  sed -i '/<script lang="ts">/a\
import { DeviceListResponse } from "@/api/device";' ./frontend/src/components/comparison/MouseComparisonView.vue
fi

if grep -q "Cannot find name 'DeviceListResponse'" ./frontend/src/pages/MouseDatabasePage.vue; then
  sed -i '/<script lang="ts">/a\
import { DeviceListResponse } from "@/api/device";' ./frontend/src/pages/MouseDatabasePage.vue
fi

# 5. Fix missing ElMessage import
if grep -q "Cannot find name 'ElMessage'" ./frontend/src/pages/MouseDatabasePage.vue; then
  sed -i '/<script lang="ts">/a\
import { ElMessage } from "element-plus";' ./frontend/src/pages/MouseDatabasePage.vue
fi

# 6. Fix SortDirection and ViewMode export issues
echo "Fixing SortDirection and ViewMode exports..."
if grep -q "Module '\"@/components/database/SortControls.vue\"' has no exported member 'SortDirection'" ./frontend/src/pages/MouseDatabasePage.vue; then
  # Add SortDirection enum to the SortControls component
  sed -i '/<script lang="ts">/a\
export enum SortDirection {\
  ASC = "asc",\
  DESC = "desc"\
}' ./frontend/src/components/database/SortControls.vue

  # Update the import in MouseDatabasePage.vue
  sed -i 's/import SortControls from "@\/components\/database\/SortControls.vue";/import SortControls, { SortDirection } from "@\/components\/database\/SortControls.vue";/' ./frontend/src/pages/MouseDatabasePage.vue
fi

if grep -q "Module '\"@/components/database/ViewToggle.vue\"' has no exported member 'ViewMode'" ./frontend/src/pages/MouseDatabasePage.vue; then
  # Add ViewMode enum to the ViewToggle component
  sed -i '/<script lang="ts">/a\
export enum ViewMode {\
  GRID = "grid",\
  LIST = "list"\
}' ./frontend/src/components/database/ViewToggle.vue

  # Update the import in MouseDatabasePage.vue
  sed -i 's/import ViewToggle from "@\/components\/database\/ViewToggle.vue";/import ViewToggle, { ViewMode } from "@\/components\/database\/ViewToggle.vue";/' ./frontend/src/pages/MouseDatabasePage.vue
fi

# 7. Fix I18nDemo setLanguage issue
if grep -q "Module '\"@/i18n\"' has no exported member 'setLanguage'" ./frontend/src/pages/I18nDemo.vue; then
  # Create a temporary file to hold the fixed content
  cat > ./frontend/src/i18n/index.ts.new << 'EOF'
import { createI18n } from 'vue-i18n';
import enUS from './locales/en-US';
import zhCN from './locales/zh-CN';

// Define supported languages
export const SUPPORTED_LANGUAGES = ['en-US', 'zh-CN', 'en', 'zh'] as const;
export type SupportedLanguage = typeof SUPPORTED_LANGUAGES[number];

// Messages for each language
const messages = {
  'en-US': enUS,
  'zh-CN': zhCN,
  'en': enUS,  // Fallback for generic English
  'zh': zhCN   // Fallback for generic Chinese
};

// Get browser language
export function getBrowserLanguage(): SupportedLanguage {
  return navigator.language as SupportedLanguage || 'en-US';
}

// Set up i18n instance
export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('language') || getBrowserLanguage(),
  fallbackLocale: 'en-US',
  messages,
});

// Function to change language
export function setLanguage(lang: SupportedLanguage) {
  i18n.global.locale.value = lang;
  localStorage.setItem('language', lang);
  document.querySelector('html')?.setAttribute('lang', lang);
}

export default i18n;
EOF
  
  # Replace the original file
  mv ./frontend/src/i18n/index.ts.new ./frontend/src/i18n/index.ts
fi

# 8. Fix SVG service type casting issue
echo "Fixing SVG service type casting..."
if grep -q "Conversion of type 'HTMLElement' to type 'SVGSVGElement'" ./frontend/src/services/svgService.ts; then
  sed -i 's/document.getElementById(elementId) as SVGSVGElement/document.getElementById(elementId) as unknown as SVGSVGElement/g' ./frontend/src/services/svgService.ts
fi

# 9. Fix useReview composable
echo "Enhancing useReview composable with missing methods..."
if grep -q "Property 'createReview' does not exist" ./frontend/src/pages/ReviewForm.vue; then
  cat > ./frontend/src/composables/useReview.ts.new << 'EOF'
import { ref } from 'vue';
import * as reviewApi from '@/api/review';
import type { Review } from '@/api/review';

export function useReview() {
  const reviews = ref<Review[]>([]);
  const loading = ref(false);

  const fetchReviews = async (params?: any) => {
    loading.value = true;
    try {
      const response = await reviewApi.getReviews(params);
      reviews.value = response.data;
      return {
        data: response.data,
        total: response.total
      };
    } catch (error) {
      console.error('Error fetching reviews:', error);
      return {
        data: [],
        total: 0
      };
    } finally {
      loading.value = false;
    }
  };

  // Add missing methods
  const getReview = async (id: string) => {
    loading.value = true;
    try {
      const response = await reviewApi.getReview(id);
      return response.data;
    } catch (error) {
      console.error('Error getting review:', error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const createReview = async (review: any) => {
    loading.value = true;
    try {
      const response = await reviewApi.createReview(review);
      return response.data;
    } catch (error) {
      console.error('Error creating review:', error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const updateReview = async (id: string, review: any) => {
    loading.value = true;
    try {
      const response = await reviewApi.updateReview(id, review);
      return response.data;
    } catch (error) {
      console.error('Error updating review:', error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  return {
    reviews,
    loading,
    fetchReviews,
    getReview,
    createReview,
    updateReview
  };
}
EOF
  mv ./frontend/src/composables/useReview.ts.new ./frontend/src/composables/useReview.ts
fi

# 10. Fix useDevice composable
echo "Enhancing useDevice composable with missing properties..."
if grep -q "Property 'userDeviceLoading' does not exist" ./frontend/src/components/form/UserDeviceForm.vue; then
  cat > ./frontend/src/composables/useDevice.ts.new << 'EOF'
import { ref } from 'vue';
import * as deviceApi from '@/api/device';
import type { Device } from '@/types/mouse';

export function useDevice() {
  const devices = ref<Device[]>([]);
  const loading = ref(false);
  
  // Add missing user device properties
  const userDevices = ref<Device[]>([]);
  const userDeviceLoading = ref(false);
  const userDevicePagination = ref({
    page: 1,
    pageSize: 10,
    total: 0
  });

  const fetchDevices = async (params?: any) => {
    loading.value = true;
    try {
      const response = await deviceApi.getDevices(params);
      devices.value = response.data;
      return {
        data: response.data,
        total: response.total
      };
    } catch (error) {
      console.error('Error fetching devices:', error);
      return {
        data: [],
        total: 0
      };
    } finally {
      loading.value = false;
    }
  };

  // Add missing user device methods
  const fetchUserDevices = async (params?: any) => {
    userDeviceLoading.value = true;
    try {
      const response = await deviceApi.getUserDevices(params);
      userDevices.value = response.data;
      userDevicePagination.value = {
        page: response.page || 1,
        pageSize: response.pageSize || 10,
        total: response.total || 0
      };
      return response;
    } catch (error) {
      console.error('Error fetching user devices:', error);
      return { data: [], total: 0 };
    } finally {
      userDeviceLoading.value = false;
    }
  };

  const fetchUserDevice = async (id: string) => {
    userDeviceLoading.value = true;
    try {
      const response = await deviceApi.getUserDevice(id);
      return response.data;
    } catch (error) {
      console.error('Error fetching user device:', error);
      return null;
    } finally {
      userDeviceLoading.value = false;
    }
  };

  const saveUserDevice = async (data: any) => {
    userDeviceLoading.value = true;
    try {
      const response = await deviceApi.createUserDevice(data);
      await fetchUserDevices(); // Refresh list
      return response.data;
    } catch (error) {
      console.error('Error saving user device:', error);
      throw error;
    } finally {
      userDeviceLoading.value = false;
    }
  };

  const updateUserDeviceConfig = async (id: string, data: any) => {
    userDeviceLoading.value = true;
    try {
      const response = await deviceApi.updateUserDevice(id, data);
      await fetchUserDevices(); // Refresh list
      return response.data;
    } catch (error) {
      console.error('Error updating user device config:', error);
      throw error;
    } finally {
      userDeviceLoading.value = false;
    }
  };

  const removeUserDevice = async (id: string) => {
    userDeviceLoading.value = true;
    try {
      await deviceApi.deleteUserDevice(id);
      await fetchUserDevices(); // Refresh list
      return true;
    } catch (error) {
      console.error('Error removing user device:', error);
      throw error;
    } finally {
      userDeviceLoading.value = false;
    }
  };

  const searchDevicesByName = async (query: string) => {
    if (!query) return [];
    loading.value = true;
    try {
      const response = await deviceApi.searchDevices({ name: query });
      return response.data;
    } catch (error) {
      console.error('Error searching devices:', error);
      return [];
    } finally {
      loading.value = false;
    }
  };

  const getGripStyleName = (gripStyle: string) => {
    const styles: Record<string, string> = {
      'palm': 'Palm Grip',
      'claw': 'Claw Grip',
      'fingertip': 'Fingertip Grip'
    };
    return styles[gripStyle] || gripStyle;
  };

  return {
    devices,
    loading,
    fetchDevices,
    getGripStyleName,
    // Add missing user device properties and methods
    userDevices,
    userDeviceLoading,
    userDevicePagination,
    fetchUserDevices,
    fetchUserDevice,
    saveUserDevice,
    updateUserDeviceConfig,
    removeUserDevice,
    searchDevicesByName,
    devicesLoading: loading
  };
}
EOF
  mv ./frontend/src/composables/useDevice.ts.new ./frontend/src/composables/useDevice.ts
fi

# 11. Fix comparison service property access
echo "Fixing comparison service property access..."
if grep -q "Property 'type' does not exist on type '{}'" ./frontend/src/services/comparisonService.ts; then
  cat > ./frontend/src/services/comparisonService.ts.new << 'EOF'
export interface MouseShape {
  type?: string;
  humpPlacement?: string;
  frontFlare?: string;
  sideCurvature?: string;
  handCompatibility?: string;
}

export function compareShapeProperties(shape1: MouseShape = {}, shape2: MouseShape = {}) {
  const differences: Record<string, { property: string; value1: string; value2: string }> = {};

  if ((shape1?.type || '') !== (shape2?.type || '')) {
    differences['type'] = {
      property: 'Type',
      value1: shape1?.type || 'N/A',
      value2: shape2?.type || 'N/A'
    };
  }

  if ((shape1?.humpPlacement || '') !== (shape2?.humpPlacement || '')) {
    differences['humpPlacement'] = {
      property: 'Hump Placement',
      value1: shape1?.humpPlacement || 'N/A',
      value2: shape2?.humpPlacement || 'N/A'
    };
  }

  if ((shape1?.frontFlare || '') !== (shape2?.frontFlare || '')) {
    differences['frontFlare'] = {
      property: 'Front Flare',
      value1: shape1?.frontFlare || 'N/A',
      value2: shape2?.frontFlare || 'N/A'
    };
  }

  if ((shape1?.sideCurvature || '') !== (shape2?.sideCurvature || '')) {
    differences['sideCurvature'] = {
      property: 'Side Curvature',
      value1: shape1?.sideCurvature || 'N/A',
      value2: shape2?.sideCurvature || 'N/A'
    };
  }

  if ((shape1?.handCompatibility || '') !== (shape2?.handCompatibility || '')) {
    differences['handCompatibility'] = {
      property: 'Hand Compatibility',
      value1: shape1?.handCompatibility || 'N/A',
      value2: shape2?.handCompatibility || 'N/A'
    };
  }

  return differences;
}
EOF
  mv ./frontend/src/services/comparisonService.ts.new ./frontend/src/services/comparisonService.ts
fi

# 12. Fix UserState in store/index.ts
echo "Fixing UserState in store..."
if grep -q "Type '{ user: null; token: null; }' is missing the following properties from type 'UserState'" ./frontend/src/stores/index.ts; then
  cat > ./frontend/src/stores/index.ts.new << 'EOF'
import { defineStore } from 'pinia';
import type { User } from '@/types/user';

export interface UserState {
  user: User | null;
  token: string | null;
  profile: any | null;
  gameSettings: any | null;
  settings: any | null;
  privacy: any | null;
  devicePreferences: any[] | null;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    user: null,
    token: null,
    profile: null,
    gameSettings: null,
    settings: null,
    privacy: null,
    devicePreferences: null
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
      this.profile = null;
      this.gameSettings = null;
      this.settings = null;
      this.privacy = null;
      this.devicePreferences = null;
    }
  },
  getters: {
    isLoggedIn(): boolean {
      return !!this.token && !!this.user;
    }
  }
});
EOF
  mv ./frontend/src/stores/index.ts.new ./frontend/src/stores/index.ts
fi

# 13. Fix router isLoggedIn access
echo "Fixing router isLoggedIn access..."
if grep -q "Property 'isLoggedIn' does not exist on type" ./frontend/src/router/index.ts; then
  sed -i 's/if (to.meta.requiresAuth && !store.isLoggedIn)/if (to.meta.requiresAuth && !store.token)/g' ./frontend/src/router/index.ts
fi

# 14. Fix MouseShapeVisualization value property error
echo "Fixing MouseShapeVisualization value property error..."
if grep -q "Property 'value' does not exist on type 'number'" ./frontend/src/components/tools/MouseShapeVisualization.vue; then
  sed -i 's/\.value//g' ./frontend/src/components/tools/MouseShapeVisualization.vue
fi

# 15. Fix DeviceForm battery property optional issues
echo "Fixing DeviceForm battery property errors..."
if grep -q "Property 'battery' is optional in type" ./frontend/src/components/form/DeviceForm.vue; then
  sed -i "s/battery: {/battery?: {/g" ./frontend/src/components/form/DeviceForm.vue
fi

# 16. Fix Review type export
if grep -q "Module '\"@/api/review\"' declares 'Review' locally, but it is not exported" ./frontend/src/composables/useReview.ts; then
  cat > ./frontend/src/api/review.ts.new << 'EOF'
import { request } from '@/utils/request';

export interface Review {
  id: string;
  userId: string;
  deviceId: string;
  title: string;
  content: string;
  rating: number;
  pros: string[];
  cons: string[];
  images?: string[];
  createdAt: string;
  updatedAt: string;
  likes: number;
  dislikes: number;
  userReaction?: 'like' | 'dislike' | null;
}

export const getReviews = (params?: any) => {
  return request({
    url: '/api/reviews',
    method: 'get',
    params
  });
};

export const getReview = (id: string) => {
  return request({
    url: `/api/reviews/${id}`,
    method: 'get'
  });
};

export const createReview = (data: any) => {
  return request({
    url: '/api/reviews',
    method: 'post',
    data
  });
};

export const updateReview = (id: string, data: any) => {
  return request({
    url: `/api/reviews/${id}`,
    method: 'put',
    data
  });
};

export const deleteReview = (id: string) => {
  return request({
    url: `/api/reviews/${id}`,
    method: 'delete'
  });
};

export const likeReview = (id: string) => {
  return request({
    url: `/api/reviews/${id}/like`,
    method: 'post'
  });
};

export const dislikeReview = (id: string) => {
  return request({
    url: `/api/reviews/${id}/dislike`,
    method: 'post'
  });
};

export const removeReaction = (id: string) => {
  return request({
    url: `/api/reviews/${id}/reaction`,
    method: 'delete'
  });
};
EOF
  mv ./frontend/src/api/review.ts.new ./frontend/src/api/review.ts
fi

# 17. Export MouseDevice type in review types
if grep -q "Module '\"@/types/review\"' declares 'MouseDevice' locally, but it is not exported" ./frontend/src/pages/ReviewDetailPage.vue; then
  # Create the file if it doesn't exist
  mkdir -p ./frontend/src/types/review
  cat > ./frontend/src/types/review/index.ts << 'EOF'
export interface MouseDevice {
  id: string;
  name: string;
  brand: string;
  type: string;
  imageUrl?: string;
  description?: string;
}

export interface ReviewRequest {
  deviceId: string;
  title: string;
  content: string;
  rating: number;
  pros?: string[];
  cons?: string[];
  images?: string[];
}

export interface ReviewResponse {
  id: string;
  userId: string;
  deviceId: string;
  device?: MouseDevice;
  title: string;
  content: string;
  rating: number;
  pros: string[];
  cons: string[];
  images?: string[];
  createdAt: string;
  updatedAt: string;
  likes: number;
  dislikes: number;
  userReaction?: 'like' | 'dislike' | null;
}
EOF
fi

echo "TypeScript error fixing script complete!"