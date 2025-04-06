<template>
  <div class="transaction-create-container">
    <el-card class="box-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>创建房产交易</h3>
        </div>
      </template>
      
      <el-form 
        :model="transactionForm" 
        :rules="rules" 
        ref="transactionFormRef" 
        label-width="120px"
        label-position="right"
        status-icon
      >
        <el-form-item label="房产证号" prop="realtyCert">
          <el-input v-model="transactionForm.realtyCert" placeholder="请输入房产证号" :disabled="!!realtyId"></el-input>
        </el-form-item>
        
        <el-form-item label="买方身份证" prop="buyerCitizenID">
          <el-input v-model="transactionForm.buyerCitizenID" placeholder="请输入买方身份证号" disabled></el-input>
        </el-form-item>

        <el-form-item label="税费" prop="tax">
          <el-input-number v-model="transactionForm.tax" :min="0" :precision="2" :step="1000" style="width: 100%"></el-input-number>
        </el-form-item>
        
        <el-form-item label="交易价格" prop="price">
          <el-input-number v-model="transactionForm.price" :min="0" :precision="2" :step="10000" style="width: 100%"></el-input-number>
        </el-form-item>
        
        <div v-if="realtyInfo" class="realty-info-container">
          <h4>房产信息确认</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="房产证号">{{ realtyInfo.realtyCert }}</el-descriptions-item>
            <el-descriptions-item label="房产地址">{{ realtyInfo.address }}</el-descriptions-item>
            <el-descriptions-item label="房产类型">{{ realtyInfo.realtyType }}</el-descriptions-item>
            <el-descriptions-item label="建筑面积">{{ realtyInfo.area }} 平方米</el-descriptions-item>
            <el-descriptions-item label="参考价格">{{ formatPrice(realtyInfo.price) }}</el-descriptions-item>
            <el-descriptions-item label="当前所有者">{{ realtyInfo.currentOwnerCitizenID }}</el-descriptions-item>
            <el-descriptions-item label="房产状态">{{ getStatusText(realtyInfo.status) }}</el-descriptions-item>
          </el-descriptions>
        </div>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitLoading">提交交易申请</el-button>
          <el-button type="info" @click="searchRealty" :disabled="!transactionForm.realtyCert">查询房产信息</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import {useUserStore} from "@/stores/user.js";
import {createTransaction} from "@/api/transaction.js";

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()
const transactionFormRef = ref(null)
const loading = ref(false)
const submitLoading = ref(false)
const realtyInfo = ref(null)
const realtyId = ref(route.query.realtyId || '')

// 表单数据
const transactionForm = reactive({
  realtyCert: route.query.realtyCert,
  buyerCitizenID: userStore.user?.citizenID,
  paymentUUIDList: [],
  tax: 0,
  price: 0
})

// 表单验证规则
const rules = reactive({
  realtyCert: [
    { required: true, message: '请输入房产证号', trigger: 'blur' }
  ],
  buyerCitizenID: [
    { required: true, message: '请输入买方身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入交易价格', trigger: 'blur' },
    { type: 'number', min: 1, message: '价格必须大于0', trigger: 'blur' }
  ]
})

// 价格格式化
const formatPrice = (price) => {
  return `¥ ${parseFloat(price).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 查询房产信息
const searchRealty = async () => {
  if (!transactionForm.realtyCert) {
    ElMessage.warning('请输入房产证号')
    return
  }
  
  loading.value = true
  try {
    const response = await axios.get(`/api/realty/cert/${transactionForm.realtyCert}`)
    realtyInfo.value = response.data
    
    // 自动填充卖方信息
    if (realtyInfo.value && realtyInfo.value.currentOwnerCitizenID) {
      transactionForm.sellerCitizenID = realtyInfo.value.currentOwnerCitizenID
    }
    
    // 检查房产状态是否可交易
    if (realtyInfo.value.status !== 'NORMAL') {
      ElMessage.warning(`当前房产状态为"${getStatusText(realtyInfo.value.status)}"，可能无法进行交易`)
    }
  } catch (error) {
    console.error('获取房产信息失败:', error)
    ElMessage.error(error.response?.data?.message || '获取房产信息失败')
    realtyInfo.value = null
  } finally {
    loading.value = false
  }
}

// 提交表单
const submitForm = async () => {
  if (!transactionFormRef.value) return
  
  await transactionFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      submitLoading.value = true
      
      try {
        await createTransaction(transactionForm)
        ElMessage.success('交易申请提交成功')
        router.push(`/transaction`)
      } catch (error) {
        console.error('创建交易失败:', error)
        ElMessage.error(error.response?.data?.message || '创建交易失败')
      } finally {
        submitLoading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 如果从房产详情页进入，自动加载房产信息
onMounted(async () => {
  if (realtyId.value) {
    loading.value = true
    try {
      const response = await axios.get(`/api/realty/${realtyId.value}`)
      realtyInfo.value = response.data
      transactionForm.realtyCert = realtyInfo.value.realtyCert
      transactionForm.sellerCitizenID = realtyInfo.value.currentOwnerCitizenID
      
      // 获取当前登录用户信息，如果是投资者，自动填充为买方
      const userJson = localStorage.getItem('user')
      if (userJson) {
        try {
          const userInfo = JSON.parse(userJson)
          if (userInfo.role === 'INVESTOR' && userInfo.citizenId) {
            transactionForm.buyerCitizenID = userStore.user.value.citizenID
          }
        } catch (e) {
          console.error('解析用户信息失败', e)
        }
      }
      
      // 设置默认价格为房产参考价格
      if (realtyInfo.value.price) {
        transactionForm.price = realtyInfo.value.price
      }
    } catch (error) {
      console.error('获取房产信息失败:', error)
      ElMessage.error(error.response?.data?.message || '获取房产信息失败')
    } finally {
      loading.value = false
    }
  }
})
</script>

<style scoped>
.transaction-create-container {
  max-width: 800px;
  margin: 20px auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.box-card {
  margin-bottom: 20px;
}

.realty-info-container {
  margin: 20px 0;
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #f5f7fa;
}
</style>
