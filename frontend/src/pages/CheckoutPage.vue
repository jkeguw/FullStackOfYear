<template>
  <div class="checkout-page">
    <div class="checkout-content">
      <h2 class="page-title">结算</h2>
      
      <div v-if="loading" class="loading-section">
        <el-skeleton :rows="15" animated />
      </div>
      
      <div v-else>
        <!-- 订单确认步骤 -->
        <div v-if="currentStep === 'confirmation' && orderDetail" class="confirmation-section">
          <OrderConfirmation 
            :order="orderDetail" 
            @pay="handlePayment"
            @cancel="handleCancelOrder"
          />
        </div>
        
        <!-- 模拟支付页面 -->
        <div v-else-if="currentStep === 'payment'" class="payment-section">
          <div class="payment-container">
            <h3>支付订单</h3>
            <div class="payment-header">
              <div class="logo-container">
                <img 
                  :src="getPaymentLogo(paymentMethod)" 
                  :alt="paymentMethod" 
                  class="payment-logo"
                />
              </div>
              <div class="payment-info">
                <div class="order-number">订单编号: {{ orderDetail?.orderNumber }}</div>
                <div class="payment-amount">支付金额: <span>¥{{ orderDetail?.total.toFixed(2) }}</span></div>
              </div>
            </div>
            
            <div class="payment-options">
              <div class="payment-form">
                <el-form 
                  v-if="paymentMethod === 'credit_card'"
                  ref="creditCardFormRef"
                  :model="creditCardForm"
                  :rules="creditCardRules"
                  label-width="120px"
                >
                  <el-form-item label="卡号" prop="cardNumber">
                    <el-input
                      v-model="creditCardForm.cardNumber"
                      placeholder="1234 5678 9012 3456"
                      maxlength="19"
                      @input="formatCreditCard"
                    />
                  </el-form-item>
                  <el-form-item label="持卡人姓名" prop="cardName">
                    <el-input
                      v-model="creditCardForm.cardName"
                      placeholder="持卡人姓名"
                    />
                  </el-form-item>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-form-item label="有效期" prop="cardExpiry">
                        <el-input
                          v-model="creditCardForm.cardExpiry"
                          placeholder="MM/YY"
                          maxlength="5"
                          @input="formatExpiry"
                        />
                      </el-form-item>
                    </el-col>
                    <el-col :span="12">
                      <el-form-item label="安全码" prop="cardCvv">
                        <el-input
                          v-model="creditCardForm.cardCvv"
                          placeholder="CVV"
                          maxlength="4"
                          show-password
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </el-form>
                
                <div v-else-if="paymentMethod === 'alipay'" class="qrcode-container">
                  <div class="qrcode">
                    <el-image 
                      src="https://t.alipayobjects.com/images/rmsweb/T1BbJfXndiXXXXXXXX.png" 
                      fit="cover"
                    />
                  </div>
                  <div class="scan-text">请使用支付宝扫码支付</div>
                </div>
                
                <div v-else-if="paymentMethod === 'wechat'" class="qrcode-container">
                  <div class="qrcode">
                    <el-image 
                      src="https://res.wx.qq.com/wxdoc/dist/assets/img/demo.ef5c5bef.jpg" 
                      fit="cover"
                    />
                  </div>
                  <div class="scan-text">请使用微信扫码支付</div>
                </div>
              </div>
            </div>
            
            <div class="payment-actions">
              <el-button @click="goBackToConfirmation">返回</el-button>
              <el-button 
                type="primary"
                @click="completePayment"
                :loading="processingPayment"
              >
                确认支付
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 结账表单步骤 -->
        <div v-else class="form-section">
          <el-steps :active="stepIndex" finish-status="success" class="checkout-steps">
            <el-step title="确认商品" />
            <el-step title="填写信息" />
            <el-step title="提交订单" />
          </el-steps>
          
          <!-- 购物车商品确认步骤 -->
          <div v-if="currentStep === 'cart'" class="cart-review">
            <h3>确认订单商品</h3>
            
            <div v-if="cartItems.length === 0" class="empty-cart">
              <el-empty description="您的购物车为空" />
              <el-button type="primary" @click="goToShop">去选购商品</el-button>
            </div>
            
            <template v-else>
              <div class="cart-items">
                <div class="cart-header">
                  <span class="item-title">商品信息</span>
                  <span class="item-price">单价</span>
                  <span class="item-quantity">数量</span>
                  <span class="item-subtotal">小计</span>
                </div>
                
                <div v-for="item in cartItems" :key="item.product_id" class="cart-item">
                  <div class="item-info">
                    <el-image 
                      :src="item.image_url || '/placeholder.png'" 
                      :alt="item.name"
                      class="item-image"
                    />
                    <div class="item-name">{{ item.name }}</div>
                  </div>
                  <div class="item-price">¥{{ item.price.toFixed(2) }}</div>
                  <div class="item-quantity">{{ item.quantity }}</div>
                  <div class="item-subtotal">¥{{ (item.price * item.quantity).toFixed(2) }}</div>
                </div>
              </div>
              
              <div class="cart-summary">
                <div class="summary-item">
                  <span>商品数量:</span>
                  <span>{{ totalItems }} 件</span>
                </div>
                <div class="summary-item">
                  <span>商品总计:</span>
                  <span>¥{{ cartTotal.toFixed(2) }}</span>
                </div>
              </div>
              
              <div class="cart-actions">
                <el-button @click="goToCart">返回购物车</el-button>
                <el-button type="primary" @click="nextStep">下一步</el-button>
              </div>
            </template>
          </div>
          
          <!-- 结账表单步骤 -->
          <div v-else-if="currentStep === 'form'" class="checkout-form-container">
            <CheckoutForm 
              :cart-items="cartItems"
              :loading="submitting"
              @submit="submitOrder"
              @cancel="previousStep"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useCart } from '@/composables/useCart';
import { useOrder } from '@/composables/useOrder';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import CheckoutForm from '@/components/order/CheckoutForm.vue';
import OrderConfirmation from '@/components/order/OrderConfirmation.vue';
import type { CreateOrderRequest, OrderResponse, PaymentCompleteRequest } from '@/types/order';

const router = useRouter();
const { cart, cartItems, total: cartTotal, clearCart } = useCart();
const { loading: orderLoading, orderDetail, submitOrder: createOrder, payOrder, cancelOrder } = useOrder();

// 步骤控制
type CheckoutStep = 'cart' | 'form' | 'confirmation' | 'payment';
const currentStep = ref<CheckoutStep>('cart');
const loading = ref(false);
const submitting = ref(false);
const processingPayment = ref(false);

// 信用卡表单
const creditCardFormRef = ref<FormInstance>();
const creditCardForm = reactive({
  cardNumber: '',
  cardName: '',
  cardExpiry: '',
  cardCvv: ''
});

// 信用卡表单验证规则
const creditCardRules = reactive<FormRules>({
  cardNumber: [
    { required: true, message: '请输入卡号', trigger: 'blur' },
    { min: 16, message: '请输入正确的信用卡号', trigger: 'blur' }
  ],
  cardName: [
    { required: true, message: '请输入持卡人姓名', trigger: 'blur' }
  ],
  cardExpiry: [
    { required: true, message: '请输入有效期', trigger: 'blur' },
    { pattern: /^\d{2}\/\d{2}$/, message: '请输入正确的有效期格式 (MM/YY)', trigger: 'blur' }
  ],
  cardCvv: [
    { required: true, message: '请输入CVV码', trigger: 'blur' },
    { pattern: /^\d{3,4}$/, message: '请输入正确的CVV码', trigger: 'blur' }
  ]
});

// 当前支付方式
const paymentMethod = computed(() => {
  return orderDetail.value?.paymentInfo.method || 'alipay';
});

// 计算商品总数
const totalItems = computed(() => {
  return cartItems.value.reduce((sum, item) => sum + item.quantity, 0);
});

// 当前步骤索引
const stepIndex = computed(() => {
  switch (currentStep.value) {
    case 'cart': return 0;
    case 'form': return 1;
    case 'confirmation': return 2;
    case 'payment': return 2;
    default: return 0;
  }
});

// 初始化
onMounted(async () => {
  // 如果直接从URL访问，检查是否有购物车商品
  if (cartItems.value.length === 0) {
    ElMessage.warning('购物车为空，无法结算');
    router.push('/mouse-database');
  }
});

// 下一步
const nextStep = () => {
  if (currentStep.value === 'cart') {
    if (cartItems.value.length === 0) {
      ElMessage.warning('购物车为空，无法结算');
      return;
    }
    currentStep.value = 'form';
  }
};

// 上一步
const previousStep = () => {
  if (currentStep.value === 'form') {
    currentStep.value = 'cart';
  } else if (currentStep.value === 'confirmation') {
    currentStep.value = 'form';
  } else if (currentStep.value === 'payment') {
    currentStep.value = 'confirmation';
  }
};

// 提交订单
const submitOrder = async (orderData: CreateOrderRequest) => {
  try {
    submitting.value = true;
    const order = await createOrder(orderData);
    
    orderDetail.value = order;
    currentStep.value = 'confirmation';
    
    // 清空购物车
    await clearCart();
    
  } catch (error) {
    console.error('创建订单失败:', error);
    ElMessage.error('创建订单失败，请重试');
  } finally {
    submitting.value = false;
  }
};

// 处理支付
const handlePayment = (orderId: string) => {
  currentStep.value = 'payment';
};

// 格式化信用卡号
const formatCreditCard = (value: string) => {
  // 移除所有非数字字符
  const cleaned = value.replace(/\D/g, '');
  // 每4位添加一个空格
  const formatted = cleaned.replace(/(\d{4})(?=\d)/g, '$1 ');
  creditCardForm.cardNumber = formatted;
};

// 格式化有效期
const formatExpiry = (value: string) => {
  // 移除所有非数字字符
  const cleaned = value.replace(/\D/g, '');
  // 添加斜杠
  if (cleaned.length >= 2) {
    creditCardForm.cardExpiry = `${cleaned.substring(0, 2)}/${cleaned.substring(2, 4)}`;
  } else {
    creditCardForm.cardExpiry = cleaned;
  }
};

// 获取支付方式图标
const getPaymentLogo = (method: string) => {
  switch (method) {
    case 'alipay':
      return 'https://img.alicdn.com/imgextra/i1/O1CN01Oo9Jgi1k3dumTHdUm_!!6000000004623-2-tps-200-66.png';
    case 'wechat':
    case 'wechat_pay':
      return 'https://res.wx.qq.com/a/wx_fed/assets/res/NTI4MWU5.png';
    case 'credit_card':
      return 'https://www.visa.com.cn/dam/VCOM/regional/ap/china/global-elements/images/visa-logo-blue-1x-desktopt.png';
    default:
      return '';
  }
};

// 返回确认页面
const goBackToConfirmation = () => {
  currentStep.value = 'confirmation';
};

// 完成支付
const completePayment = async () => {
  // 如果是信用卡支付，验证表单
  if (paymentMethod.value === 'credit_card' && creditCardFormRef.value) {
    const valid = await creditCardFormRef.value.validate().catch(() => false);
    if (!valid) {
      ElMessage.error('请正确填写信用卡信息');
      return;
    }
  }
  
  if (!orderDetail.value) {
    ElMessage.error('订单信息不存在');
    return;
  }
  
  try {
    processingPayment.value = true;
    
    // 模拟支付过程
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    // 生成一个模拟的交易ID
    const transactionId = `TX${Date.now().toString().substring(5)}`;
    
    // 更新订单支付状态
    const paymentData: PaymentCompleteRequest = {
      orderId: orderDetail.value.id,
      transactionId,
      paymentStatus: 'success'
    };
    
    const updatedOrder = await payOrder(orderDetail.value.id, paymentData);
    orderDetail.value = updatedOrder;
    currentStep.value = 'confirmation';
    
    ElMessage.success('支付成功');
    
  } catch (error) {
    console.error('支付失败:', error);
    ElMessage.error('支付处理失败，请重试');
  } finally {
    processingPayment.value = false;
  }
};

// 处理取消订单
const handleCancelOrder = async (orderId: string) => {
  try {
    loading.value = true;
    await cancelOrder(orderId, '用户取消');
    ElMessage.success('订单已取消');
    router.push('/orders');
  } catch (error) {
    console.error('取消订单失败:', error);
    ElMessage.error('取消订单失败');
  } finally {
    loading.value = false;
  }
};

// 前往购物车
const goToCart = () => {
  router.push('/cart');
};

// 前往商店
const goToShop = () => {
  router.push('/mouse-database');
};
</script>

<style scoped>
.checkout-page {
  min-height: calc(100vh - 200px);
  background-color: #f5f7fa;
  padding: 30px 0;
}

.checkout-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-title {
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
  font-weight: 600;
}

.checkout-steps {
  margin-bottom: 30px;
}

.loading-section {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.cart-review, .checkout-form-container {
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.cart-review h3 {
  font-size: 18px;
  margin-bottom: 20px;
  color: #333;
  font-weight: 600;
}

.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 0;
}

.cart-items {
  margin-bottom: 20px;
}

.cart-header {
  display: grid;
  grid-template-columns: 4fr 1fr 1fr 1fr;
  padding: 10px 15px;
  background-color: #f5f7fa;
  border-radius: 8px 8px 0 0;
  font-weight: 500;
  color: #666;
}

.cart-item {
  display: grid;
  grid-template-columns: 4fr 1fr 1fr 1fr;
  padding: 15px;
  border-bottom: 1px solid #f0f0f0;
  align-items: center;
}

.item-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
  border: 1px solid #eee;
}

.item-name {
  font-weight: 500;
}

.item-subtotal {
  font-weight: 600;
  color: #ff6700;
}

.cart-summary {
  margin: 20px 0;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.summary-item:last-child {
  margin-bottom: 0;
  font-weight: 600;
  font-size: 16px;
}

.cart-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.payment-section {
  max-width: 800px;
  margin: 0 auto;
}

.payment-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
}

.payment-container h3 {
  text-align: center;
  font-size: 22px;
  margin-bottom: 30px;
  color: #333;
}

.payment-header {
  display: flex;
  margin-bottom: 30px;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 20px;
  align-items: center;
}

.logo-container {
  width: 100px;
  margin-right: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.payment-logo {
  max-width: 100%;
  max-height: 60px;
  object-fit: contain;
}

.payment-info {
  flex: 1;
}

.order-number {
  color: #666;
  margin-bottom: 10px;
}

.payment-amount {
  font-size: 18px;
}

.payment-amount span {
  font-weight: 600;
  color: #ff6700;
  font-size: 24px;
}

.payment-form {
  max-width: 500px;
  margin: 0 auto;
}

.qrcode-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px 0;
}

.qrcode {
  width: 200px;
  height: 200px;
  border: 1px solid #eee;
  margin-bottom: 20px;
}

.scan-text {
  font-size: 16px;
  color: #666;
}

.payment-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 30px;
}
</style>