#!/bin/bash

# 修复所有类型错误的脚本

echo "开始修复前端类型问题..."

# 1. 修复MouseComparisonView.vue中的 DeviceListResponse 导入问题
if grep -q "Cannot find name 'DeviceListResponse'" frontend/src/components/comparison/MouseComparisonView.vue; then
  sed -i '3i import { DeviceListResponse } from "@/api/device";' frontend/src/components/comparison/MouseComparisonView.vue
  echo "已修复 MouseComparisonView.vue 中的 DeviceListResponse 导入"
fi

# 2. 修复MouseSelector.vue中的 image_url 属性问题
if grep -q "Property 'image_url' does not exist" frontend/src/components/comparison/MouseSelector.vue; then
  sed -i 's/image_url/imageUrl/g' frontend/src/components/comparison/MouseSelector.vue
  echo "已修复 MouseSelector.vue 中的 image_url 属性问题"
fi

# 3. 修复 DeviceForm.vue 中的 battery 必须属性问题
if grep -q "Property 'battery' is optional" frontend/src/components/form/DeviceForm.vue; then
  # 这里需要根据实际代码结构调整修复方法
  # 可能需要将 battery 设为可选或确保其必填
  echo "需要手动检查并修复 DeviceForm.vue 中的 battery 属性问题"
fi

# 4. 修复 UserDeviceForm.vue 中缺少属性问题
if grep -q "Property 'userDeviceLoading' does not exist" frontend/src/components/form/UserDeviceForm.vue; then
  echo "需要手动检查并修复 UserDeviceForm.vue 中的属性问题"
fi

# 5. 修复 MouseShapeVisualization.vue 中的 value 属性问题
if grep -q "Property 'value' does not exist on type 'number'" frontend/src/components/tools/MouseShapeVisualization.vue; then
  sed -i 's/\.value//g' frontend/src/components/tools/MouseShapeVisualization.vue
  echo "已修复 MouseShapeVisualization.vue 中的 value 属性问题"
fi

# 6. 修复 i18n/index.ts 中的类型问题
if grep -q "Type 'string' is not assignable to type" frontend/src/i18n/index.ts; then
  echo "需要手动检查并修复 i18n/index.ts 中的语言类型问题"
fi

# 7. 修复 store 中的 isLoggedIn 问题
# 在router/index.ts中使用token替代isLoggedIn

# 8. 修复 services/comparisonService.ts 中的类型问题
if grep -q "Property 'type' does not exist on type '{}'" frontend/src/services/comparisonService.ts; then
  echo "需要手动检查并修复 services/comparisonService.ts 中的类型问题"
fi

# 9. 修复 services/svgService.ts 中的类型转换问题
if grep -q "Conversion of type 'HTMLElement' to type 'SVGSVGElement'" frontend/src/services/svgService.ts; then
  echo "需要手动检查并修复 services/svgService.ts 中的类型转换问题"
fi

echo "脚本执行完成。请根据错误提示手动检查并修复剩余问题。"