<template>
  <el-dialog
    v-model="dialogVisible"
    title="支付交易"
    width="500px"
  >
    <div class="payment-dialog-content">
      <div class="payment-info-row">
        <div class="info-label">交易总金额:</div>
        <div class="info-value">{{ formatPrice(transaction?.price || 0) }}</div>
      </div>
      <div class="payment-info-row">
        <div class="info-label">已支付金额:</div>
        <div class="info-value">{{ formatPrice(totalPaid || 0) }}</div>
      </div>
      <div class="payment-info-row">
        <div class="info-label">剩余应付:</div>
        <div class="info-value highlight">{{ formatPrice(remainingAmount) }}</div>
      </div>
      
      <el-divider />
      
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="支付类型" prop="paymentType">
          <el-select v-model="form.paymentType" placeholder="请选择支付类型" style="width: 100%">
            <el-option label="房款支付" value="TRANSFER"></el-option>
            <el-option label="税费支付" value="TAX"></el-option>
            <el-option label="手续费支付" value="FEE"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="支付金额" prop="amount">
          <el-input-number 
            v-model="form.amount" 
            :min="1" 
            :max="remainingAmount" 
            :precision="2" 
            :step="1000" 
            style="width: 100%"
          ></el-input-number>
          <div class="form-tip">
            <span>最大金额: {{ formatPrice(remainingAmount) }}</span>
            <el-button type="text" @click="form.amount = remainingAmount">设为支付最大金额</el-button>
          </div>
        </el-form-item>
        
        <el-form-item label="备注" prop="remarks">
          <el-input 
            v-model="form.remarks"
            type="textarea" 
            :rows="3"
            placeholder="请输入备注信息"
          ></el-input>
        </el-form-item>
      </el-form>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="cancelPayment">取消</el-button>
        <el-button type="primary" @click="submitPayment" :loading="loading">确认支付</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import {payForTransaction} from "@/api/payment.js";
import {useUserStore} from "@/stores/user.js";

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  transaction: {
    type: Object,
    default: null
  },
  totalPaid: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['update:visible', 'payment-success'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const loading = ref(false)
const formRef = ref(null)

// 支付表单数据
const form = ref({
  paymentType: 'TRANSFER',
  amount: 0,
  remarks: ''
})

// 表单验证规则
const rules = {
  paymentType: [
    { required: true, message: '请选择支付类型', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入支付金额', trigger: 'blur' },
    { type: 'number', min: 1, message: '金额必须大于0', trigger: 'blur' }
  ],
  method: [
    { required: true, message: '请选择支付方式', trigger: 'change' }
  ]
}

// 剩余应付金额
const remainingAmount = computed(() => {
  if (!props.transaction) return 0
  const remaining = props.transaction.price - props.totalPaid
  return remaining > 0 ? remaining : 0
})

// 监听对话框可见性变化，重置表单
watch(dialogVisible, (val) => {
  if (val) {
    nextTick(() => {
      resetForm()
      // 默认设置为剩余金额的一半或全部（如果剩余金额很小）
      const defaultAmount = remainingAmount.value <= 10000 ? 
                           remainingAmount.value : 
                           Math.round(remainingAmount.value / 2)
      form.value.amount = defaultAmount
      
      // 根据剩余金额自动选择支付类型
      if (props.transaction && props.transaction.tax > 0 && props.totalPaid < props.transaction.tax) {
        form.value.paymentType = 'TAX'
      } else {
        form.value.paymentType = 'TRANSFER'
      }
    })
  }
})

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 取消支付
const cancelPayment = () => {
  dialogVisible.value = false
}

const userStore = useUserStore()
// 提交支付
const submitPayment = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid, fields) => {
    if (valid) {
      if (!props.transaction || !props.transaction.transactionUUID) {
        ElMessage.error('交易信息不完整')
        return
      }
      
      loading.value = true
      
      try {
        // 发送支付请求
        const response = await payForTransaction({
          transactionUUID: props.transaction.transactionUUID,
          amount: form.value.amount,
          payerCitizenID: userStore.user.citizenID,
          payerOrganization: userStore.user.organization,
          receiverCitizenIDHash: props.transaction.sellerCitizenIDHash,
          receiverOrganization: props.transaction.sellerOrganization,
          paymentType: form.value.paymentType,
          remarks: form.value.remarks,
        })
        
        ElMessage.success('支付成功')
        dialogVisible.value = false
        emit('payment-success')
      } catch (error) {
        console.error('支付失败:', error)
        ElMessage.error(error.response?.data?.message || '支付失败')
      } finally {
        loading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 创建支付API调用
const createPayment = async (paymentData) => {
  try {
    const response = await axios.post('/api/payment/create', paymentData)
    return response.data
  } catch (error) {
    // 使用模拟数据
    console.log('使用Mock数据')
    // 模拟1秒延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    return mockCreatePayment(paymentData)
  }
}

// Mock函数：创建支付
const mockCreatePayment = (paymentData) => {
  return {
    code: 200,
    data: {
      paymentID: 'PAY' + Date.now().toString().slice(-6),
      transactionUUID: paymentData.transactionUUID,
      paymentType: paymentData.paymentType,
      amount: paymentData.amount,
      payerCitizenID: paymentData.payerCitizenID,
      payeeCitizenID: paymentData.payeeCitizenID,
      description: paymentData.description,
      createTime: new Date().toISOString()
    },
    message: '支付创建成功'
  }
}

// 格式化价格
const formatPrice = (price) => {
  return `¥ ${parseFloat(price).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`
}
</script>

<style scoped>
.payment-dialog-content {
  padding: 10px;
}

.payment-info-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  font-size: 16px;
}

.info-label {
  color: #606266;
  font-weight: bold;
}

.info-value {
  color: #303133;
}

.info-value.highlight {
  color: #409EFF;
  font-weight: bold;
  font-size: 18px;
}

.form-tip {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #909399;
  font-size: 12px;
  margin-top: 5px;
}
</style> 