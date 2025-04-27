# 鼠标对比网站测试计划

本文档列出了鼠标对比网站项目需要编写的测试，以实现较高的测试覆盖率。

## 后端测试

### 服务层测试

1. **购物车服务测试 (`services/cart/service_test.go`)**
   - `TestCartService_AddToCart` - 测试添加商品到购物车
   - `TestCartService_GetUserCart` - 测试获取用户购物车
   - `TestCartService_UpdateCartItem` - 测试更新购物车商品数量
   - `TestCartService_RemoveCartItem` - 测试从购物车移除商品
   - `TestCartService_ClearCart` - 测试清空购物车
   - `TestCartService_GetCartByID` - 测试通过ID获取购物车

2. **订单服务测试 (`services/order/service_test.go`)**
   - `TestOrderService_CreateOrder` - 测试创建订单
   - `TestOrderService_GetOrder` - 测试获取订单详情
   - `TestOrderService_GetOrderByNumber` - 测试通过订单号获取订单
   - `TestOrderService_ListUserOrders` - 测试获取用户订单列表
   - `TestOrderService_UpdateOrderStatus` - 测试更新订单状态
   - `TestOrderService_ProcessPayment` - 测试处理支付

3. **相似度服务测试 (`services/similarity/service_test.go`)**
   - `TestSimilarityService_CalculateDimensionSimilarity` - 测试尺寸相似度计算
   - `TestSimilarityService_CalculateShapeSimilarity` - 测试形状相似度计算
   - `TestSimilarityService_CalculateTechnicalSimilarity` - 测试技术参数相似度计算
   - `TestSimilarityService_CalculateOverallSimilarity` - 测试综合相似度计算
   - `TestSimilarityService_FindSimilarMice` - 测试查找相似鼠标功能

### 处理器层测试

1. **购物车处理器测试 (`handlers/cart/handlers_test.go`)**
   - `TestCartHandler_GetCart` - 测试获取购物车API
   - `TestCartHandler_AddItem` - 测试添加商品API
   - `TestCartHandler_UpdateItem` - 测试更新商品API
   - `TestCartHandler_RemoveItem` - 测试删除商品API
   - `TestCartHandler_ClearCart` - 测试清空购物车API

2. **订单处理器测试 (`handlers/order/handlers_test.go`)**
   - `TestOrderHandler_CreateOrder` - 测试创建订单API
   - `TestOrderHandler_GetOrder` - 测试获取订单API
   - `TestOrderHandler_GetOrderByNumber` - 测试通过订单号获取API
   - `TestOrderHandler_ListUserOrders` - 测试订单列表API
   - `TestOrderHandler_UpdateOrderStatus` - 测试更新订单状态API
   - `TestOrderHandler_ProcessPayment` - 测试处理支付API

3. **设备相似度处理器测试**
   - `TestDeviceHandler_CompareMice` - 测试比较鼠标API
   - `TestDeviceHandler_FindSimilarMice` - 测试查找相似鼠标API

### 集成测试

1. **购物流程集成测试 (`tests/integration/cart_test.go`)**
   - `TestCartWorkflow` - 测试完整购物车流程

2. **订单流程集成测试 (`tests/integration/order_test.go`)**
   - `TestOrderWorkflow` - 测试完整订单流程

3. **相似度功能集成测试 (`tests/integration/similarity_test.go`)**
   - `TestSimilarityWorkflow` - 测试相似度查找流程

## 前端测试 (使用Vitest)

### 组件测试

1. **购物车组件测试**
   - `CartDrawer.spec.ts` - 测试购物车抽屉组件
   - `CartIcon.spec.ts` - 测试购物车图标组件
   - `AddToCartButton.spec.ts` - 测试添加到购物车按钮

2. **结账组件测试**
   - `CheckoutForm.spec.ts` - 测试结账表单
   - `OrderSummary.spec.ts` - 测试订单摘要
   - `PaymentMethods.spec.ts` - 测试支付方式选择
   - `ShippingForm.spec.ts` - 测试配送信息表单

3. **相似度组件测试**
   - `SimilarityResults.spec.ts` - 测试相似度结果显示
   - `ComparisonView.spec.ts` - 测试比较视图

### 服务/Composables测试

1. **购物车服务测试**
   - `useCart.spec.ts` - 测试购物车composable
   - `cartService.spec.ts` - 测试购物车API服务

2. **订单服务测试**
   - `useOrder.spec.ts` - 测试订单composable
   - `orderService.spec.ts` - 测试订单API服务

3. **相似度服务测试**
   - `useSimilarity.spec.ts` - 测试相似度composable
   - `similarityService.spec.ts` - 测试相似度计算服务

### 端到端测试 (使用Cypress或Playwright)

1. **购物流程E2E测试**
   - `cartWorkflow.spec.ts` - 测试浏览商品->添加购物车->更新数量->移除商品

2. **结账流程E2E测试**
   - `checkoutWorkflow.spec.ts` - 测试购物车->结账->支付->订单确认

3. **鼠标比较E2E测试**
   - `mouseComparison.spec.ts` - 测试选择鼠标->比较->查找相似

## 前端测试配置说明

对于前端测试，建议使用以下配置：

1. **安装Vitest及相关依赖**:
   ```bash
   npm install -D vitest @vue/test-utils @testing-library/vue happy-dom
   ```

2. **配置package.json**:
   ```json
   "scripts": {
     "dev": "vite",
     "build": "vite build",
     "test": "vitest run",
     "test:watch": "vitest",
     "test:coverage": "vitest run --coverage"
   }
   ```

3. **Vitest配置**:
   创建`vitest.config.ts`文件:
   ```typescript
   import { defineConfig } from 'vitest/config'
   import vue from '@vitejs/plugin-vue'

   export default defineConfig({
     plugins: [vue()],
     test: {
       globals: true,
       environment: 'happy-dom',
       coverage: {
         provider: 'istanbul',
         reporter: ['text', 'json', 'html'],
       },
     },
   })
   ```

## 后端测试约定

1. 使用标准Go测试框架
2. 使用testify提供断言和模拟功能
3. 使用gomock生成mock对象
4. 数据库测试使用实际MongoDB或内存版本
5. 测试函数命名规则: Test{Service/Handler}_{Method}_{Scenario}