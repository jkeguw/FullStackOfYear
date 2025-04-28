<template>
  <div class="order-confirmation">
    <div class="confirmation-header">
      <el-result
        icon="success"
        title="订单提交成功"
        sub-title="感谢您的购买，我们将尽快处理您的订单"
      >
        <template #extra>
          <el-button type="primary" @click="goToOrderDetail">查看订单详情</el-button>
          <el-button @click="continueShopping">继续购物</el-button>
        </template>
      </el-result>
    </div>

    <div class="order-detail-card">
      <div class="order-detail-header">
        <h3>订单信息</h3>
        <div class="order-number">
          订单号: <span>{{ order.orderNumber }}</span>
        </div>
      </div>

      <div class="order-detail-content">
        <div class="detail-section">
          <h4>订单状态</h4>
          <el-tag :type="getStatusTagType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>

        <div class="detail-section">
          <h4>订单商品</h4>
          <div class="order-items">
            <div v-for="item in order.items" :key="item.productId" class="order-item">
              <div class="item-image">
                <img :src="item.imageUrl || '/placeholder.png'" :alt="item.name" />
              </div>
              <div class="item-info">
                <div class="item-name">{{ item.name }}</div>
                <div class="item-price">¥{{ item.price.toFixed(2) }} × {{ item.quantity }}</div>
              </div>
              <div class="item-subtotal">¥{{ item.subtotal.toFixed(2) }}</div>
            </div>
          </div>
        </div>

        <div class="detail-section price-section">
          <div class="price-row">
            <span>商品小计:</span>
            <span>¥{{ order.subtotal.toFixed(2) }}</span>
          </div>
          <div class="price-row">
            <span>配送费:</span>
            <span>¥{{ order.shippingFee.toFixed(2) }}</span>
          </div>
          <div class="price-row">
            <span>税费:</span>
            <span>¥{{ order.tax.toFixed(2) }}</span>
          </div>
          <div v-if="order.discount > 0" class="price-row discount">
            <span>折扣:</span>
            <span>-¥{{ order.discount.toFixed(2) }}</span>
          </div>
          <div class="price-row total">
            <span>订单总计:</span>
            <span>¥{{ order.total.toFixed(2) }}</span>
          </div>
        </div>

        <div class="detail-section">
          <h4>配送信息</h4>
          <div class="shipping-info">
            <div class="info-row">
              <span class="label">收货人:</span>
              <span>{{ order.shippingInfo.name }}</span>
            </div>
            <div class="info-row">
              <span class="label">联系电话:</span>
              <span>{{ order.shippingInfo.phone }}</span>
            </div>
            <div class="info-row">
              <span class="label">收货地址:</span>
              <span>{{ formatAddress(order.shippingInfo) }}</span>
            </div>
            <div class="info-row">
              <span class="label">配送方式:</span>
              <span>{{ getShippingMethodText(order.shippingInfo.shippingMethod) }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4>支付信息</h4>
          <div class="payment-info">
            <div class="info-row">
              <span class="label">支付方式:</span>
              <span>{{ getPaymentMethodText(order.paymentInfo.method) }}</span>
            </div>
            <div class="info-row">
              <span class="label">支付状态:</span>
              <el-tag :type="getPaymentStatusTagType(order.paymentInfo.paymentStatus)">
                {{ getPaymentStatusText(order.paymentInfo.paymentStatus) }}
              </el-tag>
            </div>
          </div>
        </div>

        <div v-if="order.notes" class="detail-section">
          <h4>订单备注</h4>
          <div class="order-notes">
            {{ order.notes }}
          </div>
        </div>
      </div>

      <div class="confirmation-footer">
        <div class="button-group">
          <el-button type="primary" @click="handlePay" v-if="canPay">前往支付</el-button>
          <el-button type="danger" @click="handleCancel" v-if="canCancel">取消订单</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { computed, defineProps, defineEmits } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';
import {
  OrderResponse,
  OrderStatus,
  OrderStatusText,
  PaymentMethod,
  PaymentMethodText
} from '@/types/order';
import type { ShippingInfoResponse } from '@/types/order';

const props = defineProps<{
  order: OrderResponse;
}>();

const emit = defineEmits<{
  (e: 'pay', orderId: string): void;
  (e: 'cancel', orderId: string): void;
}>();

const router = useRouter();

// 计算属性：是否可以支付
const canPay = computed(() => {
  return props.order.status === OrderStatus.PENDING;
});

// 计算属性：是否可以取消
const canCancel = computed(() => {
  return [OrderStatus.PENDING, OrderStatus.PAID].includes(props.order.status as OrderStatus);
});

// 获取订单状态文本
const getStatusText = (status: string) => {
  return OrderStatusText[status as OrderStatus] || status;
};

// 获取订单状态标签类型
const getStatusTagType = (status: string) => {
  switch (status) {
    case OrderStatus.PENDING:
      return 'warning';
    case OrderStatus.PAID:
      return 'success';
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

// 格式化地址
const formatAddress = (shippingInfo: ShippingInfoResponse) => {
  return `${shippingInfo.country} ${shippingInfo.state} ${shippingInfo.city} ${shippingInfo.address} ${shippingInfo.zipCode}`;
};

// 前往订单详情
const goToOrderDetail = () => {
  router.push(`/orders/${props.order.id}`);
};

// 继续购物
const continueShopping = () => {
  router.push('/mouse-database');
};

// 处理支付
const handlePay = () => {
  emit('pay', props.order.id);
};

// 处理取消订单
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消此订单吗？此操作不可逆。', '取消订单', {
      confirmButtonText: '确定取消',
      cancelButtonText: '继续保留',
      type: 'warning'
    });

    emit('cancel', props.order.id);
  } catch {
    // 用户取消操作
  }
};
</script>

<style scoped>
.order-confirmation {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px 0;
}

.confirmation-header {
  margin-bottom: 30px;
}

.order-detail-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.order-detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 15px;
  margin-bottom: 20px;
}

.order-detail-header h3 {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.order-number {
  font-size: 14px;
  color: #666;
}

.order-number span {
  font-weight: bold;
  color: #333;
}

.detail-section {
  margin-bottom: 30px;
}

.detail-section h4 {
  font-size: 16px;
  margin-bottom: 15px;
  color: #333;
  font-weight: 600;
  border-bottom: 1px dashed #eee;
  padding-bottom: 8px;
}

.order-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.order-item {
  display: flex;
  align-items: center;
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 10px;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
  overflow: hidden;
  background-color: #fff;
  margin-right: 15px;
  border: 1px solid #eee;
}

.item-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item-info {
  flex: 1;
}

.item-name {
  font-weight: 500;
  margin-bottom: 5px;
}

.item-price {
  color: #999;
  font-size: 13px;
}

.item-subtotal {
  font-weight: 600;
  color: #ff6700;
  font-size: 16px;
}

.price-section {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
}

.price-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.price-row.total {
  border-top: 1px solid #eee;
  padding-top: 8px;
  margin-top: 8px;
  font-weight: 600;
  font-size: 16px;
}

.price-row.discount {
  color: #ff6700;
}

.shipping-info,
.payment-info {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.info-row {
  display: flex;
  gap: 10px;
}

.label {
  color: #999;
  width: 80px;
}

.order-notes {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  white-space: pre-line;
}

.confirmation-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.button-group {
  display: flex;
  gap: 10px;
}
</style>
