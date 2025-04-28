<template>
  <div class="order-detail-page">
    <div class="page-container">
      <div class="page-header">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/orders' }">我的订单</el-breadcrumb-item>
          <el-breadcrumb-item>订单详情</el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="15" animated />
      </div>

      <template v-else-if="orderDetail">
        <div class="order-detail-card">
          <div class="order-status-bar">
            <div class="order-number">订单号：{{ orderDetail.orderNumber }}</div>
            <div class="order-status">
              <el-tag :type="getStatusTagType(orderDetail.status)">
                {{ getStatusText(orderDetail.status) }}
              </el-tag>
            </div>
          </div>

          <el-steps
            :active="getStatusStep(orderDetail.status)"
            finish-status="success"
            class="order-steps"
            align-center
          >
            <el-step title="提交订单" :description="formatDate(orderDetail.createdAt)" />
            <el-step
              title="付款成功"
              :description="orderDetail.paidAt ? formatDate(orderDetail.paidAt) : ''"
            />
            <el-step
              title="商品发货"
              :description="orderDetail.shippedAt ? formatDate(orderDetail.shippedAt) : ''"
            />
            <el-step
              title="交易完成"
              :description="orderDetail.deliveredAt ? formatDate(orderDetail.deliveredAt) : ''"
            />
          </el-steps>

          <div class="detail-section product-section">
            <h3 class="section-title">商品信息</h3>
            <div class="product-list">
              <div class="product-table-header">
                <div class="product-cell product-info">商品信息</div>
                <div class="product-cell price">单价</div>
                <div class="product-cell quantity">数量</div>
                <div class="product-cell subtotal">小计</div>
              </div>

              <div v-for="item in orderDetail.items" :key="item.productId" class="product-row">
                <div class="product-cell product-info">
                  <div class="product-image">
                    <el-image
                      :src="item.imageUrl || '/placeholder.png'"
                      :alt="item.name"
                      fit="cover"
                    />
                  </div>
                  <div class="product-name">{{ item.name }}</div>
                </div>
                <div class="product-cell price">¥{{ item.price.toFixed(2) }}</div>
                <div class="product-cell quantity">{{ item.quantity }}</div>
                <div class="product-cell subtotal">¥{{ item.subtotal.toFixed(2) }}</div>
              </div>
            </div>
          </div>

          <div class="order-sections">
            <div class="detail-section">
              <h3 class="section-title">订单信息</h3>
              <div class="info-list">
                <div class="info-item">
                  <span class="info-label">订单编号：</span>
                  <span class="info-value">{{ orderDetail.orderNumber }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">下单时间：</span>
                  <span class="info-value">{{ formatDate(orderDetail.createdAt) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">付款时间：</span>
                  <span class="info-value">
                    {{ orderDetail.paidAt ? formatDate(orderDetail.paidAt) : '未付款' }}
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">发货时间：</span>
                  <span class="info-value">
                    {{ orderDetail.shippedAt ? formatDate(orderDetail.shippedAt) : '未发货' }}
                  </span>
                </div>
                <div v-if="orderDetail.trackingNumber" class="info-item">
                  <span class="info-label">物流单号：</span>
                  <span class="info-value">{{ orderDetail.trackingNumber }}</span>
                </div>
                <div v-if="orderDetail.cancelledAt" class="info-item">
                  <span class="info-label">取消时间：</span>
                  <span class="info-value">{{ formatDate(orderDetail.cancelledAt) }}</span>
                </div>
                <div v-if="orderDetail.cancelReason" class="info-item">
                  <span class="info-label">取消原因：</span>
                  <span class="info-value">{{ orderDetail.cancelReason }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h3 class="section-title">收货信息</h3>
              <div class="info-list">
                <div class="info-item">
                  <span class="info-label">收货人：</span>
                  <span class="info-value">{{ orderDetail.shippingInfo.name }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">联系电话：</span>
                  <span class="info-value">{{ orderDetail.shippingInfo.phone }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">收货地址：</span>
                  <span class="info-value">
                    {{ formatAddress(orderDetail.shippingInfo) }}
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">配送方式：</span>
                  <span class="info-value">
                    {{ getShippingMethodText(orderDetail.shippingInfo.shippingMethod) }}
                  </span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h3 class="section-title">支付信息</h3>
              <div class="info-list">
                <div class="info-item">
                  <span class="info-label">支付方式：</span>
                  <span class="info-value">
                    {{ getPaymentMethodText(orderDetail.paymentInfo.method) }}
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">支付状态：</span>
                  <span class="info-value">
                    <el-tag
                      size="small"
                      :type="getPaymentStatusTagType(orderDetail.paymentInfo.paymentStatus)"
                    >
                      {{ getPaymentStatusText(orderDetail.paymentInfo.paymentStatus) }}
                    </el-tag>
                  </span>
                </div>
                <div v-if="orderDetail.paymentInfo.transactionId" class="info-item">
                  <span class="info-label">交易单号：</span>
                  <span class="info-value">{{ orderDetail.paymentInfo.transactionId }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="order-summary">
            <div class="summary-item">
              <span>商品小计：</span>
              <span>¥{{ orderDetail.subtotal.toFixed(2) }}</span>
            </div>
            <div class="summary-item">
              <span>配送费：</span>
              <span>¥{{ orderDetail.shippingFee.toFixed(2) }}</span>
            </div>
            <div class="summary-item">
              <span>税费：</span>
              <span>¥{{ orderDetail.tax.toFixed(2) }}</span>
            </div>
            <div v-if="orderDetail.discount > 0" class="summary-item discount">
              <span>优惠金额：</span>
              <span>-¥{{ orderDetail.discount.toFixed(2) }}</span>
            </div>
            <div class="summary-item total">
              <span>订单总价：</span>
              <span>¥{{ orderDetail.total.toFixed(2) }}</span>
            </div>
          </div>

          <div v-if="orderDetail.notes" class="order-notes">
            <h3 class="section-title">订单备注</h3>
            <p>{{ orderDetail.notes }}</p>
          </div>

          <div class="order-actions">
            <el-button @click="goBack">返回列表</el-button>

            <el-button
              v-if="orderDetail.status === 'pending'"
              type="primary"
              @click="handlePayment(orderDetail.id)"
            >
              立即支付
            </el-button>

            <el-button
              v-if="['pending', 'paid'].includes(orderDetail.status)"
              type="danger"
              @click="handleCancelOrder(orderDetail.id)"
            >
              取消订单
            </el-button>

            <el-button
              v-if="orderDetail.status === 'shipped'"
              type="success"
              @click="handleConfirmReceipt(orderDetail.id)"
            >
              确认收货
            </el-button>
          </div>
        </div>
      </template>

      <div v-else class="error-container">
        <el-empty description="订单不存在或已被删除" />
        <el-button type="primary" @click="goToOrderList">返回订单列表</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessageBox, ElMessage } from 'element-plus';
import { useOrder } from '@/composables/useOrder';
import { OrderStatus, OrderStatusText, PaymentMethod, PaymentMethodText } from '@/types/order';
import type { ShippingInfoResponse } from '@/types/order';
import { formatDateTime } from '@/utils/date';

const route = useRoute();
const router = useRouter();
const { loading, orderDetail, fetchOrderDetail, cancelOrder, changeOrderStatus } = useOrder();

// 从路由参数获取订单ID
const orderId = route.params.id as string;

// 初始化
onMounted(async () => {
  await fetchOrderDetail(orderId);
});

// 格式化日期
const formatDate = (date: string) => {
  return formatDateTime(new Date(date));
};

// 格式化地址
const formatAddress = (shippingInfo: ShippingInfoResponse) => {
  return `${shippingInfo.country} ${shippingInfo.state} ${shippingInfo.city} ${shippingInfo.address} ${shippingInfo.zipCode}`;
};

// 获取订单状态对应的步骤索引
const getStatusStep = (status: string) => {
  switch (status) {
    case OrderStatus.PENDING:
      return 0;
    case OrderStatus.PAID:
      return 1;
    case OrderStatus.SHIPPED:
      return 2;
    case OrderStatus.DELIVERED:
      return 3;
    case OrderStatus.CANCELLED:
    case OrderStatus.REFUNDED:
      return 0; // 取消或退款状态显示为第一步
    default:
      return 0;
  }
};

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

// 获取支付方式文本
const getPaymentMethodText = (method: string) => {
  return PaymentMethodText[method as PaymentMethod] || method;
};

// 获取配送方式文本
const getShippingMethodText = (method: string) => {
  const methodMap: Record<string, string> = {
    standard: '标准配送 (3-5个工作日)',
    express: '快速配送 (1-2个工作日)'
  };
  return methodMap[method] || method;
};

// 获取支付状态文本
const getPaymentStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待支付',
    success: '支付成功',
    failed: '支付失败',
    refunded: '已退款'
  };
  return statusMap[status] || status;
};

// 获取支付状态标签类型
const getPaymentStatusTagType = (status: string) => {
  switch (status) {
    case 'pending':
      return 'warning';
    case 'success':
      return 'success';
    case 'failed':
      return 'danger';
    case 'refunded':
      return 'info';
    default:
      return 'info';
  }
};

// 处理支付
const handlePayment = (orderId: string) => {
  router.push(`/checkout/payment/${orderId}`);
};

// 处理取消订单
const handleCancelOrder = async (orderId: string) => {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？取消后无法恢复。', '取消订单', {
      confirmButtonText: '确定取消',
      cancelButtonText: '继续保留',
      type: 'warning'
    });

    await cancelOrder(orderId, '用户主动取消');
    ElMessage.success('订单已取消');

    // 刷新订单详情
    await fetchOrderDetail(orderId);
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
    await ElMessageBox.confirm('确认已收到商品吗？', '确认收货', {
      confirmButtonText: '确认收货',
      cancelButtonText: '取消',
      type: 'info'
    });

    await changeOrderStatus(orderId, { status: OrderStatus.DELIVERED });
    ElMessage.success('确认收货成功');

    // 刷新订单详情
    await fetchOrderDetail(orderId);
  } catch (error) {
    // 用户取消操作或发生错误
    if (error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
};

// 返回上一页
const goBack = () => {
  router.back();
};

// 前往订单列表
const goToOrderList = () => {
  router.push('/orders');
};
</script>

<style scoped>
.order-detail-page {
  background-color: #f5f7fa;
  min-height: calc(100vh - 200px);
  padding: 30px 0;
}

.page-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  margin-bottom: 20px;
}

.loading-container,
.error-container {
  background-color: #fff;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 30px;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 50px 20px;
}

.order-detail-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
  margin-bottom: 30px;
}

.order-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.order-number {
  font-size: 16px;
  font-weight: 500;
}

.order-steps {
  margin-bottom: 30px;
  padding: 20px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-section {
  margin-bottom: 30px;
}

.product-section {
  margin-bottom: 30px;
}

.product-list {
  border: 1px solid #f0f0f0;
  border-radius: 4px;
}

.product-table-header {
  display: grid;
  grid-template-columns: 3fr 1fr 1fr 1fr;
  background-color: #f5f7fa;
  padding: 12px 20px;
  font-weight: 500;
}

.product-row {
  display: grid;
  grid-template-columns: 3fr 1fr 1fr 1fr;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
  align-items: center;
}

.product-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.product-image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid #f0f0f0;
}

.product-name {
  font-weight: 500;
}

.quantity {
  text-align: center;
}

.price,
.subtotal {
  text-align: right;
}

.subtotal {
  font-weight: 600;
  color: #ff6700;
}

.order-sections {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 30px;
  margin-bottom: 30px;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-item {
  display: flex;
  margin-bottom: 10px;
}

.info-label {
  width: 100px;
  color: #666;
}

.info-value {
  flex: 1;
}

.order-summary {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 30px;
}

.summary-item {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}

.summary-item span:first-child {
  margin-right: 30px;
  color: #666;
}

.summary-item.discount {
  color: #ff6700;
}

.summary-item.total {
  font-size: 18px;
  font-weight: 600;
  border-top: 1px solid #eee;
  padding-top: 10px;
  margin-top: 10px;
}

.order-notes {
  margin-bottom: 30px;
  padding: 0 20px;
}

.order-notes p {
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 4px;
  white-space: pre-line;
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}
</style>
