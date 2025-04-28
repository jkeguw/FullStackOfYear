<template>
  <div class="checkout-form">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
      class="checkout-inner-form"
    >
      <!-- 配送信息部分 -->
      <h3 class="form-section-title">收货信息</h3>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="收货人" prop="shippingInfo.name">
            <el-input v-model="formData.shippingInfo.name" placeholder="请输入收货人姓名" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="联系电话" prop="shippingInfo.phone">
            <el-input v-model="formData.shippingInfo.phone" placeholder="请输入联系电话" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="电子邮箱" prop="shippingInfo.email">
        <el-input v-model="formData.shippingInfo.email" placeholder="请输入电子邮箱（选填）" />
      </el-form-item>

      <el-form-item label="详细地址" prop="shippingInfo.address">
        <el-input v-model="formData.shippingInfo.address" placeholder="请输入详细地址" />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="国家/地区" prop="shippingInfo.country">
            <el-select v-model="formData.shippingInfo.country" placeholder="请选择国家/地区">
              <el-option label="中国" value="China" />
              <el-option label="美国" value="USA" />
              <el-option label="日本" value="Japan" />
              <el-option label="韩国" value="Korea" />
              <el-option label="其他" value="Other" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="省/州" prop="shippingInfo.state">
            <el-input v-model="formData.shippingInfo.state" placeholder="请输入省/州" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="城市" prop="shippingInfo.city">
            <el-input v-model="formData.shippingInfo.city" placeholder="请输入城市" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="邮政编码" prop="shippingInfo.zipCode">
        <el-input v-model="formData.shippingInfo.zipCode" placeholder="请输入邮政编码" />
      </el-form-item>

      <el-form-item label="配送方式" prop="shippingInfo.shippingMethod">
        <el-radio-group v-model="formData.shippingInfo.shippingMethod">
          <el-radio label="standard">标准配送 (3-5个工作日)</el-radio>
          <el-radio label="express">快速配送 (1-2个工作日)</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 支付方式部分 -->
      <h3 class="form-section-title">支付方式</h3>
      <el-form-item label="支付方式" prop="paymentMethod">
        <el-radio-group v-model="formData.paymentMethod">
          <el-radio label="alipay">
            <i class="payment-icon alipay-icon"></i>
            支付宝
          </el-radio>
          <el-radio label="wechat">
            <i class="payment-icon wechat-icon"></i>
            微信支付
          </el-radio>
          <el-radio label="credit_card">
            <i class="payment-icon credit-card-icon"></i>
            信用卡
          </el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 信用卡信息（仅在选择信用卡时显示） -->
      <div v-if="formData.paymentMethod === 'credit_card'" class="credit-card-form">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="卡号" prop="cardNumber">
              <el-input
                v-model="formData.cardNumber"
                placeholder="请输入信用卡卡号"
                maxlength="19"
                @input="formatCreditCard"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="有效期" prop="cardExpiry">
              <el-input
                v-model="formData.cardExpiry"
                placeholder="MM/YY"
                maxlength="5"
                @input="formatExpiry"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="CVV码" prop="cardCvv">
              <el-input v-model="formData.cardCvv" placeholder="CVV" maxlength="4" show-password />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- 订单备注 -->
      <el-form-item label="订单备注">
        <el-input
          v-model="formData.notes"
          type="textarea"
          :rows="3"
          placeholder="如有特殊要求，请在此说明（选填）"
        />
      </el-form-item>

      <!-- 提交按钮 -->
      <div class="form-actions">
        <el-button @click="$emit('cancel')">返回购物车</el-button>
        <el-button type="primary" :loading="loading" @click="submitForm"> 提交订单 </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, defineEmits, defineProps } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import type { CartItem } from '@/types/cart';
import type { CreateOrderRequest, OrderItemRequest, ShippingInfoRequest } from '@/types/order';

const props = defineProps<{
  cartItems: CartItem[];
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'submit', formData: CreateOrderRequest): void;
  (e: 'cancel'): void;
}>();

const formRef = ref<FormInstance>();

// 表单数据
const formData = reactive({
  shippingInfo: {
    name: '',
    phone: '',
    email: '',
    address: '',
    city: '',
    state: '',
    zipCode: '',
    country: 'China',
    shippingMethod: 'standard'
  } as ShippingInfoRequest,
  paymentMethod: 'alipay',
  notes: '',
  // 信用卡字段
  cardNumber: '',
  cardExpiry: '',
  cardCvv: ''
});

// 表单验证规则
const rules = reactive<FormRules>({
  'shippingInfo.name': [
    { required: true, message: '请输入收货人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度应在2到20个字符之间', trigger: 'blur' }
  ],
  'shippingInfo.phone': [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  'shippingInfo.email': [{ type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }],
  'shippingInfo.address': [
    { required: true, message: '请输入详细地址', trigger: 'blur' },
    { min: 5, max: 100, message: '地址长度应在5到100个字符之间', trigger: 'blur' }
  ],
  'shippingInfo.city': [{ required: true, message: '请输入城市', trigger: 'blur' }],
  'shippingInfo.state': [{ required: true, message: '请输入省/州', trigger: 'blur' }],
  'shippingInfo.zipCode': [
    { required: true, message: '请输入邮政编码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '请输入正确的邮政编码', trigger: 'blur' }
  ],
  'shippingInfo.country': [{ required: true, message: '请选择国家/地区', trigger: 'change' }],
  'shippingInfo.shippingMethod': [{ required: true, message: '请选择配送方式', trigger: 'change' }],
  paymentMethod: [{ required: true, message: '请选择支付方式', trigger: 'change' }],
  // 信用卡表单验证（仅在选择信用卡支付时验证）
  cardNumber: [
    { required: true, message: '请输入卡号', trigger: 'blur' },
    { min: 16, message: '请输入正确的信用卡号', trigger: 'blur' }
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

// 格式化信用卡号
const formatCreditCard = (value: string) => {
  // 移除所有非数字字符
  const cleaned = value.replace(/\D/g, '');
  // 每4位添加一个空格
  const formatted = cleaned.replace(/(\d{4})(?=\d)/g, '$1 ');
  formData.cardNumber = formatted;
};

// 格式化有效期
const formatExpiry = (value: string) => {
  // 移除所有非数字字符
  const cleaned = value.replace(/\D/g, '');
  // 添加斜杠
  if (cleaned.length >= 2) {
    formData.cardExpiry = `${cleaned.substring(0, 2)}/${cleaned.substring(2, 4)}`;
  } else {
    formData.cardExpiry = cleaned;
  }
};

// 构建订单数据
const buildOrderData = (): CreateOrderRequest => {
  // 把购物车商品转换为订单商品
  const orderItems: OrderItemRequest[] = props.cartItems.map((item) => ({
    productId: item.productId,
    productType: item.productType,
    quantity: item.quantity
  }));

  return {
    items: orderItems,
    shippingInfo: formData.shippingInfo,
    paymentMethod: formData.paymentMethod,
    notes: formData.notes
  };
};

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid): Promise<void> => {
    if (valid) {
      if (props.cartItems.length === 0) {
        ElMessage.warning('购物车为空，无法提交订单');
        return;
      }

      // 构建订单数据
      const orderData = buildOrderData();

      // 触发提交事件
      emit('submit', orderData);
    } else {
      ElMessage.error('请正确填写表单信息');
    }
  });
};
</script>

<style scoped>
.checkout-form {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.checkout-inner-form {
  max-width: 800px;
  margin: 0 auto;
}

.form-section-title {
  margin-top: 20px;
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 10px;
}

.credit-card-form {
  background-color: #f9f9f9;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.payment-icon {
  display: inline-block;
  width: 24px;
  height: 24px;
  margin-right: 6px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  vertical-align: middle;
}

.alipay-icon {
  background-image: url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzMiAzMiI+PHBhdGggZD0iTTMyIDEyLjUxMmMwIC4xNDYtMTIuMjc0IDkuOTA1LTEyLjI3NCA5LjkwNWExLjQxNyAxLjQxNyAwIDAgMS0uODc3LjMxIDEuMzg4IDEuMzg4IDAgMCAxLS44NDktLjI5MUwwIDEyLjUxMlYxOC43MmMwICAxLjQyOS43OCAgMi43NDcgMi4wMzcgMy40MzJsMTIuNTIzIDYuODIxLjAwMy0uMDA2LjAwNy4wMDQgMTIuNTE4LTYuODJhNC4wMTYgNC4wMTYgMCAwIDAgMi4wMzctMy40MzFWMTIuNTEyaC0uMDAxeiIgZmlsbD0iIzAwYWVlZiIvPjxwYXRoIGQ9Ik0wIDE4LjcyVjguNzU2YzAtMS40MjkuNzgtMi43NDcgMi4wMzctMy40MzJsMTIuNTIzLTYuODIxLjAwMy4wMDYuMDA3LS4wMDQgMTIuNTE4IDYuODJhNC4wMTYgNC4wMTYgMCAwIDEgMi4wMzcgMy40MzF2OS45NjNsLTEyLjI3NCAxMC4yMTVhMS4zODggMS4zODggMCAwIDEtLjg0OS4yOTEgMS40MTcgMS40MTcgMCAwIDEtLjg3Ny0uMzFMMCAyMi45MjdWMTguNzJ6bTE5LjAyNC04Ljg4M2MuNTA2IDAgLjkxNi40MjIuOTE2Ljk0MnMtLjQxLjk0Mi0uOTE2Ljk0MmMtLjUwNiAwLS45MTYtLjQyMi0uOTE2LS45NDJzLjQxLS45NDIuOTE2LS45NDJ6bS0xLjgzMyAwYy41MDYgMCAuOTE2LjQyMi45MTYuOTQycy0uNDEuOTQyLS45MTYuOTQyYy0uNTA2IDAtLjkxNi0uNDIyLS45MTYtLjk0MnMuNDEtLjk0Mi45MTYtLjk0MnptLTEuODMzIDBjLjUwNiAwIC45MTYuNDIyLjkxNi45NDJzLS40MS45NDItLjkxNi45NDJjLS41MDYgMC0uOTE2LS40MjItLjkxNi0uOTQycy40MS0uOTQyLjkxNi0uOTQyeiIgZmlsbD0iIzAwYWVlZiIvPjwvc3ZnPg==');
}

.wechat-icon {
  background-image: url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzMiAzMiI+PHBhdGggZD0iTTIxLjc1IDExLjc1Yy0uMDYyIDAtLjEyNS4wMDUtLjE4NS4wMTVhNi4zNDMgNi4zNDMgMCAwIDAtMy43MDYtMS4wMTVjLTMuNTEzIDAtNi4zNTkgMi40MS02LjM1OSA1LjM4NWExMS42MTEgMTEuNjExIDAgMCAwIDEuMDk0IDMuMzMzLjQ5LjQ5IDAgMCAxIC4wMDIuNDUybC0uMjU4LjYwNS43NS0uNDI5YS43MTEuNzExIDAgMCAxIC42MTQtLjA1IDkuMDk4IDkuMDk4IDAgMCAwIDMgLjc1NyA0LjY1MiA0LjY1MiAwIDAgMS0uMzMzLTEuNzI3YzAtMi44MjggMi43MjQtNS4yMiA2LjE0NS01LjMyYS4yNS4yNSAwIDAgMCAuMjQyLS4yNDUuMjUuMjUgMCAwIDAtLjI1LS4yNTV6bS01LjE0MyAyLjUxMmEuNzUuNzUgMCAxIDEgMC0xLjVhLjc1Ljc1IDAgMCAxIDAgMS41em00IDEuNWEuNzUuNzUgMCAwIDEtLjc1LS43NS43NS43NSAwIDAgMSAuNzUtLjc1LjYwMy42MDMgMCAwIDEgLjEzLjAzLjc1Ljc1IDAgMCAxIC42Mi43Mi43NS43NSAwIDAgMS0uNzUuNzV6IiBmaWxsPSIjMDcwIi8+PHBhdGggZD0iTTE1Ljc5IDEwLjg2NkM5LjI5NCAxMS4xOTEgMy43NSAxNS45MDQgMy43NSAyMS43NWMwIDMuODk3IDIuOTcgNy4yNDIgNy4yMTIgOS4xM2wuMDM4LS4yNy0xLjAyOC0zLjM2Yy0xLjkzMy0uMDU0LTMuNDcyLTEuNjI1LTMuNDcyLTMuNmEzLjYgMy42IDAgMCAxIDMuNi0zLjZjMS43MSAwIDMuMTUgMS4xOTggMy41MiAyLjhaIiBmaWxsPSIjMDcwIi8+PHBhdGggZD0iTTI3LjY2NyAyMS43NWMwIDEuNzEtMS4xOTggMy4xNDktMi44IDMuNTJhNC44IDQuOCAwIDAgMS00Ljc4IDBsLS4wMzgtLjAzYy0uMDM4IDAtLjA1My4wMTYtLjA5MS4wM2E0LjggNC44IDAgMCAxLTQuNzkgMGMtLjc1Ny0uMzItMS40MjgtLjgzNy0xLjkxLTEuNS0uNDgtLjY0LS43OTktMS4zNzItLjc5MS0yLjE3MyAwLS4xNTQuMDEtLjMwOC4wMy0uNDYybC4wNDUtLjEzYy4zNzQtMS43MjMgMS45MTItMy4wMDggMy43MjUtMy4wMDggMS44MDcgMCAzLjM3NSAxLjMyMSAzLjcyNSAzLjEzOC4wNi4xNS4wOS4zMDguMDkuNDYybC0uMDQ1LjI3Yy4wNDUuMDc3LjEwNi4xNTMuMTY1LjIzLjAwOS4wMDQuMDE5LjAwNi4wMjkuMDA2cy4wMi0uMDAxLjAyOS0uMDA2Yy43MS0uNzUgMS43LTEuMjA5IDIuNzkyLTEuMjA5IDIuMTMuMDAyIDMuNjE1IDEuNjc0IDMuNjE1IDMuNjAyem0tOS42LjI2N2EuOC44IDAgMSAxIDAtMS42LjguOCAwIDAgMSAwIDEuNnptNC44IDEuNmEuOC44IDAgMSAxIDAtMS42LjguOCAwIDAgMSAwIDEuNnoiIGZpbGw9IiMwNzAiLz48L3N2Zz4=');
}

.credit-card-icon {
  background-image: url('data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iaXNvLTg4NTktMSI/Pg0KPCEtLSBHZW5lcmF0b3I6IEFkb2JlIElsbHVzdHJhdG9yIDE5LjAuMCwgU1ZHIEV4cG9ydCBQbHVnLUluIC4gU1ZHIFZlcnNpb246IDYuMDAgQnVpbGQgMCkgIC0tPg0KPHN2ZyB2ZXJzaW9uPSIxLjEiIGlkPSJDYXBhXzEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4Ig0KCSB2aWV3Qm94PSIwIDAgNTEyLjE2IDUxMi4xNiIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNTEyLjE2IDUxMi4xNjsiIHhtbDpzcGFjZT0icHJlc2VydmUiPg0KPHBhdGggc3R5bGU9ImZpbGw6IzAwRjJBOTsiIGQ9Ik00NzguNjA5LDEyOC4xNkg0Mi41NUMxOS4wOTcsMTI4LjE2LDAsMTQ3LjI1OSwwLDE3MC43MTJ2MTcwLjczOQ0KCWMwLDIzLjQ1MiwxOS4wOTcsNDIuNTQ5LDQyLjU1LDQyLjU0OWg0MzYuMDU5YzIzLjQ1MiwwLDQyLjU1MS0xOS4wOTcsNDIuNTUxLTQyLjU0OVYxNzAuNzEyDQoJQzUyMS4xNiwxNDcuMjU5LDUwMi4wNjIsMTI4LjE2LDQ3OC42MDksMTI4LjE2eiIvPg0KPGc+DQoJPHBhdGggc3R5bGU9ImZpbGw6IzA2RDZBMDsiIGQ9Ik00My42MjIsMzQxLjQ1Yy0wLjU5NC0wLjAwMS0xLjA3MS0wLjQ4MS0xLjA3MS0xLjA3NVYxNzEuNzg2DQoJCWMwLTAuNTkyLDAuNDc4LTEuMDcxLDEuMDcxLTEuMDcxaDQyNC45MTVjMC41OTQsMCwxLjA3MiwwLjQ3OSwxLjA3MiwxLjA3MXYxNjguNTg4YzAsMC41OTQtMC40NzksMS4wNzUtMS4wNzIsMS4wNzVINDMuNjIyeiIvPg0KCTxwYXRoIHN0eWxlPSJmaWxsOiMwNkQ2QTA7IiBkPSJNNDkwLjA4NSwxNzBjMC0xLjY1NC0xLjMyNi00LTIuOTYtNGgtNDYyLjA5Yy0xLjYzNCwwLTIuOTYsMi4zNDYtMi45Niw0djQwLjE2aDQ2OC4wMVYxNzB6Ii8+DQo8L2c+DQo8cmVjdCB4PSIyMi4wNzUiIHk9IjI0MC4xNiIgc3R5bGU9ImZpbGw6IzA2QzZBMTsiIHdpZHRoPSI0NjguMDEiIGhlaWdodD0iNDAiLz4NCjxwYXRoIGQ9Ik0wLDIzNS4xNmg1MTIuMTZ2MTZjMCwyMy40NTItMTkuMDk3LDQyLjU1LTQyLjU1MSw0Mi41NUg0Mi41NUMxOS4wOTcsMjkzLjcxLDAsMjc0LjYxMywwLDI1MS4xNlYyMzUuMTZ6Ii8+DQo8Y2lyY2xlIHN0eWxlPSJmaWxsOiNGRkREM0E7IiBjeD0iNDAwLjA4IiBjeT0iMjI0LjE2IiByPSI0OCIvPg0KPGNpcmNsZSBzdHlsZT0iZmlsbDojRkY5RjQwOyIgY3g9IjQzMi4wOCIgY3k9IjIyNC4xNiIgcj0iNDgiLz4NCjxwYXRoIHN0eWxlPSJmaWxsOiNGRjdCMDA7IiBkPSJNNDMyLjA4LDE3Ni4xNmMtMjYuNTEsMC00OCwyMS40OS00OCw0OGMwLDI2LjUxLDIxLjQ5LDQ4LDQ4LDQ4VjE3Ni4xNnoiLz4NCjxwYXRoIHN0eWxlPSJmaWxsOiNGRkNEMDA7IiBkPSJNNDAwLjA4LDE3Ni4xNmMtMjYuNTEsMC00OCwyMS40OS00OCw0OGMwLDI2LjUxLDIxLjQ5LDQ4LDQ4LDQ4VjE3Ni4xNnoiLz4NCjxwYXRoIHN0eWxlPSJmaWxsOiNGRjlBMDA7IiBkPSJNNDE2LjA4LDIyNC4xNmMwLDI0LjQzNi0xNi45MTgsNDQuODczLTM5LjY0Niw1MC4yOTZjNy42MjEsMi45NTgsMTYuMTAzLDQuMzY1LDI0Ljk5OSw0LjAzNw0KCWMxLjA0NS0wLjAzOSwyLjA4OS0wLjEwOCwzLjEzLTAuMjEzYzIyLjc5Mi01LjQxLDM5Ljc1LTI1Ljg1NywzOS43NS01MC4zMzZjMC0yNC4zMzMtMTYuNzY1LTQ0LjY4Ny0zOS4zNzYtNTAuMjE2DQoJYy0xLjE0OS0wLjEwNi0yLjMwMS0wLjE3NS0zLjQ1Ni0wLjIxMmMtOC44NzMtMC4zMzYtMTcuMzMzLDEuMDY5LTI0LjkzOSw0LjAxOUM0MDguMzQ0LDE3OS41NDcsNDE2LjA4LDE5OS44MTQsNDE2LjA4LDIyNC4xNnoiLz4NCjxnPg0KPC9nPg0KPGc+DQo8L2c+DQo8Zz4NCjwvZz4NCjxnPg0KPC9nPg0KPGc+DQo8L2c+DQo8Zz4NCjwvZz4NCjxnPg0KPC9nPg0KPGc+DQo8L2c+DQo8Zz4NCjwvZz4NCjxnPg0KPC9nPg0KPGc+DQo8L2c+DQo8Zz4NCjwvZz4NCjxnPg0KPC9nPg0KPGc+DQo8L2c+DQo8Zz4NCjwvZz4NCjwvc3ZnPg0K');
}
</style>
