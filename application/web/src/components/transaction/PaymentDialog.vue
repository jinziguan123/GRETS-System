<template>
  <el-dialog
    v-model="dialogVisible"
    title="支付交易"
    width="500px"
  >
    <div class="payment-dialog-content">
      <div class="payment-info-row">
        <div class="info-label">交易总金额:</div>
        <div class="info-value">{{ formatPrice(transaction?.price + transaction?.tax || 0) }}</div>
      </div>
      
      <!-- 分别显示房款和税费信息 -->
      <div class="payment-breakdown">
        <div class="breakdown-item">
          <span class="breakdown-label">房款:</span>
          <span class="breakdown-value">{{ formatPrice(transaction?.price || 0) }}</span>
        </div>
        <div class="breakdown-item">
          <span class="breakdown-label">税费:</span>
          <span class="breakdown-value">{{ formatPrice(transaction?.tax || 0) }}</span>
        </div>
      </div>
      
      <!-- 分别显示已支付的房款和税费 -->
      <div class="payment-info-row">
        <div class="info-label">已支付房款:</div>
        <div class="info-value">{{ formatPrice(totalPaidTransfer || 0) }}</div>
      </div>
      <div class="payment-info-row">
        <div class="info-label">已支付税费:</div>
        <div class="info-value">{{ formatPrice(totalPaidTax || 0) }}</div>
      </div>
      <div class="payment-info-row">
        <div class="info-label">已支付总额:</div>
        <div class="info-value">{{ formatPrice(totalPaid || 0) }}</div>
      </div>
      
      <!-- 根据选择的支付类型显示相应的剩余金额 -->
      <div class="payment-info-row">
        <div class="info-label">剩余应付:</div>
        <div class="info-value highlight">{{ formatPrice(currentRemainingAmount) }}</div>
      </div>
      
      <el-divider />
      
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="支付类型" prop="paymentType">
          <el-select v-model="form.paymentType" placeholder="请选择支付类型" style="width: 100%">
            <el-option label="房款支付" value="TRANSFER"></el-option>
            <el-option label="税费支付" value="TAX"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="支付金额" prop="amount">
          <el-input-number 
            v-model="form.amount" 
            :min="1" 
            :max="currentRemainingAmount" 
            :precision="2" 
            :step="1000" 
            style="width: 100%"
          ></el-input-number>
          <div class="form-tip">
            <span>当前可支付最大金额: {{ formatPrice(currentRemainingAmount) }}</span>
            <el-button type="text" @click="form.amount = currentRemainingAmount">设为最大金额</el-button>
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
  },
  // 添加房款和税费的已支付金额属性
  totalPaidTransfer: {
    type: Number,
    default: 0
  },
  totalPaidTax: {
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

// 剩余房款金额
const remainingTransferAmount = computed(() => {
  if (!props.transaction) return 0
  const remaining = props.transaction.price - (props.totalPaidTransfer || 0)
  return remaining > 0 ? remaining : 0
})

// 剩余税费金额
const remainingTaxAmount = computed(() => {
  if (!props.transaction) return 0
  const remaining = props.transaction.tax - (props.totalPaidTax || 0)
  return remaining > 0 ? remaining : 0
})

// 总剩余金额
const remainingAmount = computed(() => {
  return remainingTransferAmount.value + remainingTaxAmount.value
})

// 根据选择的支付类型返回当前应显示的剩余金额
const currentRemainingAmount = computed(() => {
  if (form.value.paymentType === 'TAX') {
    return remainingTaxAmount.value
  } else {
    return remainingTransferAmount.value
  }
})

// 监听对话框可见性变化，重置表单
watch(dialogVisible, (val) => {
  if (val) {
    nextTick(() => {
      resetForm()
      
      // 自动选择支付类型 - 优先选择还有剩余的税费
      if (remainingTaxAmount.value > 0) {
        form.value.paymentType = 'TAX'
        form.value.amount = Math.min(remainingTaxAmount.value, 10000) // 默认税费金额
      } else if (remainingTransferAmount.value > 0) {
        form.value.paymentType = 'TRANSFER'
        form.value.amount = Math.min(remainingTransferAmount.value, 50000) // 默认房款金额
      }
    })
  }
})

// 监听支付类型变化，自动调整金额
watch(() => form.value.paymentType, (newType) => {
  if (newType === 'TAX') {
    form.value.amount = Math.min(remainingTaxAmount.value, form.value.amount)
  } else {
    form.value.amount = Math.min(remainingTransferAmount.value, form.value.amount)
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
      
      // 检查金额是否超出剩余可支付金额
      if (form.value.paymentType === 'TAX' && form.value.amount > remainingTaxAmount.value) {
        ElMessage.error(`税费支付金额不能超过剩余应付税费：${formatPrice(remainingTaxAmount.value)}`)
        return
      } else if (form.value.paymentType === 'TRANSFER' && form.value.amount > remainingTransferAmount.value) {
        ElMessage.error(`房款支付金额不能超过剩余应付房款：${formatPrice(remainingTransferAmount.value)}`)
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
  padding: 8px 0;
  font-size: 16px;
}

.payment-breakdown {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  margin: 10px 0;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.breakdown-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.breakdown-label {
  color: #606266;
  font-size: 14px;
  margin-bottom: 5px;
}

.breakdown-value {
  font-weight: bold;
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