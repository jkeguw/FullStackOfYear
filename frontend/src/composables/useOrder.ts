import { ref, reactive, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  createOrder,
  getOrder,
  getOrderList,
  updateOrderStatus,
  processPayment
} from '../api/order';
import type {
  CreateOrderRequest,
  OrderResponse,
  UpdateOrderStatusRequest,
  PaymentCompleteRequest
} from '../types/order';

export function useOrder() {
  return _useOrder();
}

export default function _useOrder() {
  const loading = ref(false);
  const orderDetail = ref<OrderResponse | null>(null);
  const orderList = ref<OrderResponse[]>([]);
  const pagination = reactive({
    total: 0,
    current: 1,
    pageSize: 10
  });

  // 创建订单
  const submitOrder = async (orderData: CreateOrderRequest) => {
    loading.value = true;
    try {
      const { data } = await createOrder(orderData);
      ElMessage.success('订单创建成功');
      return data;
    } catch (error: any) {
      ElMessage.error(error.message || '创建订单失败');
      throw error;
    } finally {
      loading.value = false;
    }
  };

  // 获取订单详情
  const fetchOrderDetail = async (orderId: string) => {
    if (!orderId) return;

    loading.value = true;
    try {
      const { data } = await getOrder(orderId);
      orderDetail.value = data;
      return data;
    } catch (error: any) {
      ElMessage.error(error.message || '获取订单详情失败');
      throw error;
    } finally {
      loading.value = false;
    }
  };

  // 获取订单列表
  const fetchOrderList = async (page = 1, pageSize = 10) => {
    loading.value = true;
    try {
      const { data } = await getOrderList(page, pageSize);
      orderList.value = data.orders;
      pagination.total = data.totalCount;
      pagination.current = data.currentPage;
      pagination.pageSize = data.pageSize;
      return data;
    } catch (error: any) {
      ElMessage.error(error.message || '获取订单列表失败');
      throw error;
    } finally {
      loading.value = false;
    }
  };

  // 更新订单状态
  const changeOrderStatus = async (orderId: string, statusData: UpdateOrderStatusRequest) => {
    loading.value = true;
    try {
      const { data } = await updateOrderStatus(orderId, statusData);

      // 如果正在查看的是当前订单，更新详情
      if (orderDetail.value && orderDetail.value.id === orderId) {
        orderDetail.value = data;
      }

      // 更新列表中的订单状态
      const index = orderList.value.findIndex((item) => item.id === orderId);
      if (index !== -1) {
        orderList.value[index] = data;
      }

      ElMessage.success('订单状态更新成功');
      return data;
    } catch (error: any) {
      ElMessage.error(error.message || '更新订单状态失败');
      throw error;
    } finally {
      loading.value = false;
    }
  };

  // 取消订单
  const cancelOrder = async (orderId: string, reason: string) => {
    try {
      const result = await ElMessageBox.confirm('确定要取消此订单吗？此操作不可逆。', '取消订单', {
        confirmButtonText: '确定取消',
        cancelButtonText: '保留订单',
        type: 'warning'
      });

      if (result === 'confirm') {
        return await changeOrderStatus(orderId, {
          status: 'cancelled',
          cancelReason: reason
        });
      }
    } catch {
      // 用户取消操作
      return null;
    }
  };

  // 支付订单
  const payOrder = async (orderId: string, paymentData: PaymentCompleteRequest) => {
    loading.value = true;
    try {
      const { data } = await processPayment(orderId, paymentData);

      // 如果正在查看的是当前订单，更新详情
      if (orderDetail.value && orderDetail.value.id === orderId) {
        orderDetail.value = data;
      }

      // 更新列表中的订单
      const index = orderList.value.findIndex((item) => item.id === orderId);
      if (index !== -1) {
        orderList.value[index] = data;
      }

      ElMessage.success('支付处理成功');
      return data;
    } catch (error: any) {
      ElMessage.error(error.message || '处理支付失败');
      throw error;
    } finally {
      loading.value = false;
    }
  };

  // 筛选订单状态
  const filterOrdersByStatus = (status: string) => {
    if (!status || status === 'all') {
      return orderList.value;
    }
    return orderList.value.filter((order) => order.status === status);
  };

  // 计算订单状态统计
  const orderStatusStats = computed(() => {
    const stats = {
      pending: 0,
      paid: 0,
      shipped: 0,
      delivered: 0,
      cancelled: 0,
      refunded: 0,
      total: orderList.value.length
    };

    orderList.value.forEach((order) => {
      if (Object.prototype.hasOwnProperty.call(stats, order.status)) {
        stats[order.status as keyof typeof stats]++;
      }
    });

    return stats;
  });

  return {
    loading,
    orderDetail,
    orderList,
    pagination,
    submitOrder,
    fetchOrderDetail,
    fetchOrderList,
    changeOrderStatus,
    cancelOrder,
    payOrder,
    filterOrdersByStatus,
    orderStatusStats
  };
}
