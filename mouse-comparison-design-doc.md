# 鼠标对比电商网站设计文档

## 1. 项目概述

本项目旨在开发一个类似eloshapes.com的鼠标对比网站，提供形状比较、相似度搜索、数据库查询和鼠标购买功能。网站允许用户直观比较不同鼠标的外形和参数，通过算法查找相似鼠标，并提供在线购买功能。

## 2. 技术栈

- **前端**：Vue3, TypeScript, Vite
- **后端**：Go
- **数据库**：MongoDB
- **状态管理**：Pinia
- **SVG处理**：D3.js / SVG.js
- **UI框架**：TailwindCSS

## 3. 系统架构

采用单体应用架构（非微服务），前后端分离：

```
前端应用
 ├── 公共模块（无需登录）
 │   ├── 鼠标比较
 │   ├── 相似度查找
 │   ├── 数据库浏览
 │   └── 评测查看
 └── 购物模块（需登录）
     ├── 用户账户
     ├── 购物车
     └── 订单管理
     
后端服务（统一API）
 ├── 控制器层
 │   ├── 鼠标数据控制器
 │   ├── 用户控制器
 │   ├── 订单控制器
 │   └── 支付控制器
 └── 服务层
     ├── 鼠标比较服务
     ├── 相似度计算服务
     ├── 购物车服务
     └── 订单服务
```

## 4. 最小功能实现单元拆分

### 前端组件

#### A. 页面布局与导航
1. **布局组件** ✅
   - 页面框架(Layout) ✅
   - 导航栏(Navbar) ✅
   - 页脚(Footer) ✅

2. **鼠标数据显示组件** ✅
   - 鼠标数据卡片组件(MouseCard) ✅
   - 参数表格组件(SpecsTable) ✅
   - 分页组件(Pagination) ✅

#### B. 鼠标比较功能
3. **SVG显示组件** ✅
   - SVG容器组件(SvgContainer) ✅
   - 鼠标俯视图组件(TopViewSvg) ✅
   - 鼠标侧视图组件(SideViewSvg) ✅

4. **比较工具组件** ✅
   - 鼠标选择器组件(MouseSelector) ✅
   - 可拖动尺子组件(DraggableRuler) ✅
   - 刻度尺组件(ScaleRuler) ✅
   - 重叠比较组件(OverlayComparison) ✅

#### C. 相似鼠标功能
5. **相似度搜索组件** ✅
   - 搜索表单组件(SimilaritySearchForm) ✅
   - 相似度结果组件(SimilarityResults) ✅

#### D. 数据库浏览功能
6. **数据浏览组件** ✅
   - 过滤面板组件(FilterPanel) ✅
   - 排序控件组件(SortControls) ✅
   - 列表/网格视图切换组件(ViewToggle) ✅

#### E. 评测功能
7. **评测组件** ✅
   - 评测列表组件(ReviewList) ✅
   - 评测详情组件(ReviewDetail) ✅
   - 评测图集组件(ReviewGallery) ✅

#### F. 购物功能
8. **产品购买组件** ⚠️ 待完成
   - 购买按钮组件(BuyButton) ⚠️ 框架完成
   - 购物车组件(ShoppingCart) ⚠️ 待完成
   - 加入购物车通知组件(AddToCartNotification) ⚠️ 待完成

9. **结账流程组件** ⚠️ 待完成
   - 结账表单组件(CheckoutForm) ⚠️ 待完成
   - 支付选择组件(PaymentOptions) ⚠️ 待完成
   - 订单确认组件(OrderConfirmation) ⚠️ 待完成

#### G. 用户功能
10. **用户认证组件** ✅
    - 登录组件(Login) ✅
    - 注册组件(Register) ✅
    - 用户菜单组件(UserMenu) ✅

11. **用户中心组件** ⚠️ 部分完成
    - 账户信息组件(AccountInfo) ✅
    - 订单历史组件(OrderHistory) ⚠️ 待完成
    - 订单详情组件(OrderDetail) ⚠️ 待完成

### 前端状态管理

12. **状态仓库** ✅
    - 鼠标数据仓库(mouseStore) ✅
    - 比较状态仓库(comparisonStore) ✅
    - 用户认证仓库(authStore) ✅
    - 购物车仓库(cartStore) ✅
    - 订单仓库(orderStore) ⚠️ 待完成

### 前端服务

13. **API服务** ✅
    - HTTP客户端服务(apiService) ✅
    - 鼠标数据服务(mouseDataService) ✅
    - 认证服务(authService) ✅
    - 购物车服务(cartService) ⚠️ 待完成
    - 订单服务(orderService) ⚠️ 待完成

14. **工具服务** ✅
    - SVG处理服务(svgService) ✅
    - 相似度算法服务(similarityService) ✅
    - 本地存储服务(storageService) ✅

### 后端控制器

15. **API控制器**
    - 鼠标数据控制器(MouseController)
    - 相似度控制器(SimilarityController)
    - 用户控制器(UserController)
    - 购物车控制器(CartController)
    - 订单控制器(OrderController)
    - 支付控制器(PaymentController)

### 后端服务

16. **业务逻辑服务**
    - 鼠标数据服务(MouseService)
    - 相似度计算服务(SimilarityService)
    - 用户服务(UserService)
    - 购物车服务(CartService)
    - 订单服务(OrderService)
    - 支付处理服务(PaymentService)

### 数据库模型

17. **数据模型**
    - 鼠标模型(Mouse)
    - 用户模型(User)
    - 购物车模型(Cart)
    - 订单模型(Order)
    - 产品模型(Product)
    - 评测模型(Review)

## 5. 关键功能设计

### 5.1 鼠标比较功能

**SVG比较方案**：
- 使用标准化的坐标系统确保所有鼠标SVG具有相同比例
- 支持透明度调整的重叠视图
- 支持并排比较视图
- 集成可拖动尺子工具和固定刻度尺

**参数比较方案**：
- 对比表格展示关键参数差异
- 参数差异可视化（如条形图对比）

### 5.2 相似鼠标查找功能

**相似度算法思路**：
- 基于SVG轮廓形状的特征提取
- 基于参数权重的多维相似度计算
- 可调整相似度权重（形状vs参数）

**相似度搜索流程**：
1. 用户选择基准鼠标
2. 后端计算相似度分数
3. 返回排序后的相似鼠标列表

### 5.3 购买流程整合

**无缝购买体验**：
- 在比较页面和相似度结果中直接添加购买按钮
- 未登录用户点击购买按钮时显示登录提示
- 登录后自动继续购买流程

**购物流程**：
1. 用户点击"加入购物车"按钮
2. 系统验证登录状态
3. 添加商品到购物车并显示通知
4. 用户进入购物车确认商品
5. 用户进入结账流程
6. 支付完成后显示订单确认

### 5.4 可拖动尺子组件

**设计方案**：
- 实现可拖拽的半透明尺子组件
- 支持任意方向放置
- 支持不同单位切换（mm/cm/inch）
- 与实际物理尺寸精确对应

## 6. 数据流设计

### 6.1 鼠标数据流

```
API请求 → 后端API → 数据库查询 → 后端处理 
→ 前端接收 → 状态管理 → UI渲染
```

### 6.2 购物流程数据流

```
加入购物车 → 购物车状态更新 → 服务器同步 
→ 结账流程 → 支付请求 → 订单创建 → 订单确认
```

## 7. 数据库设计

### 7.1 鼠标集合(Mouse)

```json
{
  "_id": "ObjectId",
  "name": "String",
  "brand": "String",
  "dimensions": {
    "length": "Number",
    "width": "Number", 
    "height": "Number"
  },
  "weight": "Number",
  "shape": "String",
  "humpPlacement": "String",
  "frontFlare": "String",
  "sideCurvature": "String",
  "handCompatibility": "String",
  "thumbRest": "String",
  "ringFingerRest": "String",
  "material": "String",
  "connectivity": "String",
  "sensor": "String",
  "sensorTechnology": "String",
  "sensorPosition": "String",
  "dpi": "Number",
  "pollingRate": "Number",
  "trackingSpeed": "Number", // IPS (inches per second)
  "acceleration": "Number", // G
  "sideButtons": "Number",
  "middleButtons": "Number",
  "svgData": {
    "topView": "String",
    "sideView": "String"
  },
  "relatedProductId": "ObjectId",
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

### 7.2 产品集合(Product)

```json
{
  "_id": "ObjectId",
  "mouseId": "ObjectId",
  "name": "String",
  "price": "Number",
  "stock": "Number",
  "description": "String",
  "images": ["String"],
  "specs": "Object",
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

### 7.3 用户集合(User)

```json
{
  "_id": "ObjectId",
  "username": "String",
  "email": "String",
  "passwordHash": "String",
  "addresses": ["Object"],
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

### 7.4 订单集合(Order)

```json
{
  "_id": "ObjectId",
  "userId": "ObjectId",
  "items": [{
    "productId": "ObjectId",
    "quantity": "Number",
    "price": "Number"
  }],
  "totalAmount": "Number",
  "status": "String",
  "shippingAddress": "Object",
  "paymentDetails": "Object",
  "createdAt": "Date",
  "updatedAt": "Date"
}
```

## 8. API接口设计

### 8.1 鼠标数据API

- `GET /api/mice` - 获取鼠标列表
- `GET /api/mice/:id` - 获取单个鼠标详情
- `GET /api/mice/compare?ids=id1,id2` - 获取比较数据
- `GET /api/mice/similar/:id` - 获取相似鼠标

### 8.2 用户API

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/users/me` - 获取当前用户信息
- `PUT /api/users/me` - 更新用户信息

### 8.3 购物API

- `GET /api/cart` - 获取购物车
- `POST /api/cart/items` - 添加商品到购物车
- `PUT /api/cart/items/:id` - 更新购物车商品
- `DELETE /api/cart/items/:id` - 删除购物车商品
- `POST /api/orders` - 创建订单
- `GET /api/orders` - 获取订单列表
- `GET /api/orders/:id` - 获取订单详情

## 9. 开发计划

### 9.1 阶段划分

**阶段1: 基础架构和鼠标比较功能** (2周)
- 前后端基础架构搭建
- 数据库设计和初始化
- 实现基本的鼠标数据API
- 开发SVG显示和比较组件
- 实现可拖动尺子和刻度尺

**阶段2: 相似度查找和数据库浏览功能** (1-2周)
- 实现相似度算法
- 开发相似度搜索界面
- 完善数据库浏览和过滤功能
- 开发鼠标详情页面

**阶段3: 用户认证和购物功能** (1-2周)
- 实现用户注册和登录
- 开发购物车功能
- 集成购买按钮到比较和相似度页面
- 实现结账流程

**阶段4: 集成测试和优化** (1周)
- 全站集成测试
- 用户界面优化
- 性能优化
- 部署准备

### 9.2 总体时间线

总计开发时间: **5-7周**

## 10. 技术挑战与解决方案

### 10.1 SVG比较精确性

**挑战**: 确保不同鼠标SVG以相同比例和参考点进行比较
**解决方案**: 
- 实现标准化处理流程，确保所有SVG使用统一坐标系统
- 根据实际物理尺寸缩放SVG

### 10.2 相似度算法有效性

**挑战**: 开发能精确反映形状相似度的算法
**解决方案**:
- 结合轮廓形状分析和参数权重
- 使用主成分分析(PCA)降低参数维度提高准确性

### 10.3 购物流程与鼠标比较的无缝集成

**挑战**: 在比较功能中自然融入购买流程
**解决方案**:
- 采用非侵入式设计，在适当位置添加购买按钮
- 使用状态管理确保购物车状态在全站保持一致

### 10.4 多语言支持

**挑战**: 提供全球用户友好的多语言体验
**解决方案**:
- 实现完整的i18n国际化解决方案
- 支持英文和中文两种主要语言
- 使用JSON文件存储翻译资源
- 前后端统一的语言切换机制
