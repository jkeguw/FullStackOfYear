<template>
  <div class="order-list-page">
    <div class="page-container">
      <h2 class="page-title">我的订单</h2>
      
      <div class="page-tabs">
        <el-tabs v-model="activeTab" @tab-click="handleTabClick">
          <el-tab-pane name="all" label="全部订单"></el-tab-pane>
          <el-tab-pane name="pending" label="待付款"></el-tab-pane>
          <el-tab-pane name="paid" label="待发货"></el-tab-pane>
          <el-tab-pane name="shipped" label="待收货"></el-tab-pane>
          <el-tab-pane name="completed" label="已完成"></el-tab-pane>
        </el-tabs>
      </div>
      
      <div v-if="loading" class="loading-orders">
        <el-skeleton :rows="10" animated />
      </div>
      
      <template v-else>
        <div v-if="filteredOrders.length === 0" class="empty-orders">
          <el-empty description="暂无订单" />
          <el-button type="primary" @click="goToShop">去购物</el-button>
        </div>
        
        <div v-else class="order-list">
          <div v-for="order in filteredOrders" :key="order.id" class="order-card">
            <div class="order-header">
              <div class="order-header-left">
                <div class="order-date">{{ formatDate(order.createdAt) }}</div>
                <div class="order-number">订单号：{{ order.orderNumber }}</div>
              </div>
              <div class="order-status">
                <el-tag :type="getStatusTagType(order.status)">
                  {{ getStatusText(order.status) }}
                </el-tag>
              </div>
            </div>
            
            <div class="order-items">
              <div v-for="item in order.items" :key="item.productId" class="order-item">
                <el-image 
                  :src="item.imageUrl || '/placeholder.png'" 
                  :alt="item.name"
                  class="item-image"
                />
                <div class="item-info">
                  <div class="item-name">{{ item.name }}</div>
                  <div class="item-price">¥{{ item.price.toFixed(2) }} × {{ item.quantity }}</div>
                </div>
              </div>
            </div>
            
            <div class="order-footer">
              <div class="order-total">
                <span>共{{ getTotalQuantity(order.items) }}件商品，总计：</span>
                <span class="price">¥{{ order.total.toFixed(2) }}</span>
              </div>
              
              <div class="order-actions">
                <el-button 
                  v-if="order.status === 'pending'" 
                  type="primary" 
                  size="small"
                  @click="handlePayment(order.id)"
                >
                  立即支付
                </el-button>
                
                <el-button 
                  v-if="['pending', 'paid'].includes(order.status)" 
                  type="danger" 
                  size="small"
                  @click="handleCancelOrder(order.id)"
                >
                  取消订单
                </el-button>
                
                <el-button 
                  v-if="order.status === 'shipped'" 
                  type="success" 
                  size="small"
                  @click="handleConfirmReceipt(order.id)"
                >
                  确认收货
                </el-button>
                
                <el-button 
                  size="small"
                  @click="viewOrderDetail(order.id)"
                >
                  查看详情
                </el-button>
              </div>
            </div>
          </div>
          
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="pagination.current"
              v-model:page-size="pagination.pageSize"
              :page-sizes="[5, 10, 20, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="pagination.total"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { useOrder } from '@/composables/useOrder';
import { OrderStatus, OrderStatusText } from '@/types/order';
import type { OrderItemResponse } from '@/types/order';
import { formatDateTime } from '@/utils/date';

const router = useRouter();
const {
  loading,
  orderList,
  pagination,
  fetchOrderList,
  cancelOrder,
  changeOrderStatus,
} = useOrder();

// 活动标签页
const activeTab = ref('all');

// 初始化
onMounted(async () => {
  await fetchOrderList();
});

// 根据标签过滤订单
const filteredOrders = computed(() => {
  if (activeTab.value === 'all') {
    return orderList.value;
  }
  
  if (activeTab.value === 'completed') {
    return orderList.value.filter(order => order.status === OrderStatus.DELIVERED);
  }
  
  return orderList.value.filter(order => order.status === activeTab.value);
});

// 获取订单状态文字
const getStatusText = (status: string) => {
  return OrderStatusText[status as OrderStatus] || status;
};

// 获取订单状态标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case OrderStatus.PENDING:
      return 'warning';
    case OrderStatus.PAID:
      return 'primary';
    case OrderStatus.SHIPPED:
      return 'info';
    case OrderStatus.DELIVERED:
      return 'success';
    case OrderStatus.CANCELLED:
      return 'danger';
    case OrderStatus.REFUNDED:
      return 'danger';
    default:
      return 'info';
  }
};

// 计算订单商品总数量
const getTotalQuantity = (items: OrderItemResponse[]) => {
  return items.reduce((sum, item) => sum + item.quantity, 0);
};

// 格式化日期
const formatDate = (date: string) => {
  return formatDateTime(new Date(date));
};

// 处理标签页点击
const handleTabClick = () => {
  // 重新加载当前页
  fetchOrderList(pagination.current, pagination.pageSize);
};

// 处理页码变化
const handleCurrentChange = (page: number) => {
  fetchOrderList(page, pagination.pageSize);
};

// 处理每页数量变化
const handleSizeChange = (size: number) => {
  fetchOrderList(1, size);
};

// 处理支付
const handlePayment = (orderId: string) => {
  router.push(`/checkout/payment/${orderId}`);
};

// 处理取消订单
const handleCancelOrder = async (orderId: string) => {
  try {
    await ElMessageBox.confirm(
      '确定要取消此订单吗？取消后无法恢复。',
      '取消订单',
      {
        confirmButtonText: '确定取消',
        cancelButtonText: '继续保留',
        type: 'warning'
      }
    );
    
    await cancelOrder(orderId, '用户主动取消');
    ElMessage.success('订单已取消');
    
    // 刷新订单列表
    await fetchOrderList(pagination.current, pagination.pageSize);
    
  } catch (error) {
    // 用户取消操作或发生错误
    if (error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
};

// 处理确认收货
const handleConfirmReceipt = async (orderId: string) => {
  try {
    await ElMessageBox.confirm(
      '确认已收到商品吗？',
      '确认收货',
      {
        confirmButtonText: '确认收货',
        cancelButtonText: '取消',
        type: 'info'
      }
    );
    
    await changeOrderStatus(orderId, { status: OrderStatus.DELIVERED });
    ElMessage.success('确认收货成功');
    
    // 刷新订单列表
    await fetchOrderList(pagination.current, pagination.pageSize);
    
  } catch (error) {
    // 用户取消操作或发生错误
    if (error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
};

// 查看订单详情
const viewOrderDetail = (orderId: string) => {
  router.push(`/orders/${orderId}`);
};

// 前往商店
const goToShop = () => {
  router.push('/mouse-database');
};
</script>

<style scoped>
.order-list-page {
  background-color: #f5f7fa;
  min-height: calc(100vh - 200px);
  padding: 30px 0;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #333;
}

.loading-orders {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.empty-orders {
  background-color: #fff;
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.order-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  overflow: hidden;
}

.order-header {
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fafafa;
}

.order-header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.order-date {
  color: #666;
  font-size: 14px;
}

.order-number {
  font-weight: 500;
}

.order-items {
  padding: 15px 20px;
}

.order-item {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px dashed #f0f0f0;
}

.order-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.item-image {
  width: 70px;
  height: 70px;
  border-radius: 4px;
  border: 1px solid #eee;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  margin-bottom: 5px;
}

.item-price {
  color: #666;
  font-size: 14px;
}

.order-footer {
  padding: 15px 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fafafa;
}

.order-total {
  font-size: 14px;
}

.price {
  color: #ff6700;
  font-weight: 600;
  font-size: 16px;
}

.order-actions {
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
}
</style>