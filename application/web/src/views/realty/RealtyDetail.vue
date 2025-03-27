<template>
  <div class="realty-detail-container">
    <el-card class="box-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>房产详情</h3>
          <div class="header-actions">
            <el-button type="primary" @click="goEdit" v-if="hasEditPermission">编辑</el-button>
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>
      
      <el-descriptions :column="2" border>
        <el-descriptions-item label="不动产证号">{{ realty.realtyCert }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTagType(realty.status)">
            {{ getStatusText(realty.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ realty.address }}</el-descriptions-item>
        <el-descriptions-item label="房产类型">{{ realty.realtyType }}</el-descriptions-item>
        <el-descriptions-item label="面积">{{ realty.area }} 平方米</el-descriptions-item>
        <el-descriptions-item label="价格">{{ formatPrice(realty.price) }}</el-descriptions-item>
        <el-descriptions-item label="当前所有者">{{ realty.currentOwnerName || '未知' }}</el-descriptions-item>
        <el-descriptions-item label="登记日期">{{ formatDate(realty.registrationDate) }}</el-descriptions-item>
        <el-descriptions-item label="最后更新时间">{{ formatDate(realty.lastUpdateDate) }}</el-descriptions-item>
      </el-descriptions>

      <!-- 交易记录 -->
      <div class="transaction-history" v-if="transactions.length > 0">
        <h4>交易记录</h4>
        <el-table :data="transactions" style="width: 100%" stripe>
          <el-table-column prop="transactionUUID" label="交易ID" width="180" />
          <el-table-column prop="sellerName" label="卖方" width="120" />
          <el-table-column prop="buyerName" label="买方" width="120" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getTransactionStatusTagType(scope.row.status)">
                {{ getTransactionStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdTime" label="创建时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.createdTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="completedTime" label="完成时间" width="180">
            <template #default="scope">
              {{ scope.row.completedTime ? formatDate(scope.row.completedTime) : '-' }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template #default="scope">
              <el-button type="primary" link @click="viewTransaction(scope.row.transactionUUID)">查看</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 相关操作 -->
      <div class="actions" v-if="realty.status === 'NORMAL' && hasTransactionPermission">
        <el-divider />
        <el-button type="primary" @click="createTransaction">发起交易</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(true)
const realty = reactive({
  realtyCert: '',
  address: '',
  realtyType: '',
  price: 0,
  area: 0,
  status: '',
  registrationDate: '',
  lastUpdateDate: '',
  currentOwnerName: ''
})

// 交易记录
const transactions = ref([])

// 获取房产状态对应的Tag类型
const getStatusTagType = (status) => {
  const statusMap = {
    'NORMAL': 'success',
    'IN_TRANSACTION': 'warning',
    'MORTGAGED': 'info',
    'FROZEN': 'danger'
  }
  return statusMap[status] || ''
}

// 获取房产状态对应的文本
const getStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 获取交易状态对应的Tag类型
const getTransactionStatusTagType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'APPROVED': 'success',
    'REJECTED': 'danger',
    'COMPLETED': 'success'
  }
  return statusMap[status] || ''
}

// 获取交易状态对应的文本
const getTransactionStatusText = (status) => {
  const statusMap = {
    'PENDING': '待处理',
    'APPROVED': '已批准',
    'REJECTED': '已拒绝',
    'COMPLETED': '已完成'
  }
  return statusMap[status] || status
}

// 格式化价格
const formatPrice = (price) => {
  return `¥ ${price.toLocaleString()} 元`
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

// 编辑权限
const hasEditPermission = computed(() => {
  return userStore.role === 'GOVERNMENT'
})

// 交易权限
const hasTransactionPermission = computed(() => {
  return userStore.role === 'INVESTOR'
})

// 获取房产详情
const fetchRealtyDetail = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/api/realty/${route.params.id}`)
    Object.assign(realty, response.data)
    await fetchTransactionHistory()
  } catch (error) {
    console.error('获取房产详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取房产详情失败')
  } finally {
    loading.value = false
  }
}

// 获取交易历史
const fetchTransactionHistory = async () => {
  try {
    const response = await axios.get(`/api/transaction/list?realtyId=${route.params.id}`)
    transactions.value = response.data.records || []
  } catch (error) {
    console.error('获取交易历史失败:', error)
  }
}

// 前往编辑页面
const goEdit = () => {
  router.push(`/realty/edit/${route.params.id}`)
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 查看交易详情
const viewTransaction = (transactionId) => {
  router.push(`/transaction/detail/${transactionId}`)
}

// 创建交易
const createTransaction = () => {
  router.push(`/transaction/create?realtyId=${route.params.id}`)
}

onMounted(() => {
  fetchRealtyDetail()
})
</script>

<style scoped>
.realty-detail-container {
  max-width: 1000px;
  margin: 20px auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.box-card {
  margin-bottom: 20px;
}

.transaction-history {
  margin-top: 30px;
}

.actions {
  margin-top: 20px;
  text-align: right;
}
</style>
