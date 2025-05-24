<template>
  <div class="transaction-list-container">
    <el-card >
      <template #header>
        <div class="card-header">
          <h3>交易列表</h3>
          <div class="header-buttons">
            <el-button type="primary" @click="refreshData">刷新数据</el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索框 -->
      <el-form :model="searchForm" class="search-form">
        <el-row>
          <el-col>
            <el-form-item label="交易编号">
              <el-input v-model="searchForm.transactionUUID" placeholder="输入交易编号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col>
            <el-form-item label="房产证号">
              <el-input v-model="searchForm.realtyCert" placeholder="输入房产证号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col>
            <el-form-item label="交易状态">
              <el-select v-model="searchForm.status" placeholder="选择交易状态" clearable style="width: 150px">
                <el-option label="待处理" value="PENDING"></el-option>
                <el-option label="已批准" value="APPROVED"></el-option>
                <el-option label="已拒绝" value="REJECTED"></el-option>
                <el-option label="已完成" value="COMPLETED"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col>
            <el-form-item label="卖方身份证">
              <el-input v-model="searchForm.sellerCitizenID" placeholder="输入卖方身份证号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col>
            <el-form-item label="买方身份证">
              <el-input v-model="searchForm.buyerCitizenID" placeholder="输入买方身份证号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item>
              <el-button type="primary" @click="handleSearch">搜索</el-button>
              <el-button @click="resetSearch">重置</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      
      <!-- 交易列表表格 -->
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        border
        stripe
        highlight-current-row
      >
        <el-table-column prop="transactionUUID" label="交易编号" width="auto"></el-table-column>
        <el-table-column prop="realtyCertHash" label="房产证号" width="auto"></el-table-column>
        <el-table-column prop="sellerCitizenIDHash" label="卖方身份证" width="auto"></el-table-column>
        <el-table-column prop="buyerCitizenIDHash" label="买方身份证" width="auto"></el-table-column>
        <el-table-column prop="status" label="交易状态" width="auto">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="160">
          <template #default="scope">
            {{ formatDate(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="auto">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewDetail(scope.row)" :disabled="!isTransactionParty(scope.row)">查看详情</el-button>
            <el-button 
              v-if="canCheckTransaction(scope.row)"
              type="success" 
              size="small" 
              @click="checkTransaction(scope.row)"
            >审核</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import {queryTransactionList} from "@/api/transaction.js";
import CryptoJS from 'crypto-js'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const searchForm = reactive({
  status: '',
  realtyCert: '',
  sellerCitizenID: '',
  buyerCitizenID: '',
  transactionUUID: '',
  pageSize: pageSize.value,
  pageNumber: currentPage.value
})

// 获取当前用户信息
const userInfo = computed(() => {
  const userJson = localStorage.getItem('userInfo')
  if (!userJson) return null
  try {
    return JSON.parse(userJson)
  } catch (e) {
    return null
  }
})

// 交易查询条件
const transactionQuery = reactive({
  realtyCert: '',
  status: '',
  pageSize: pageSize.value,
  pageNumber: currentPage.value
})

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchTransactionList()
}

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  currentPage.value = 1
  fetchTransactionList()
}

// 处理每页显示数量变化
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchTransactionList()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchTransactionList()
}

// 获取交易列表
const fetchTransactionList = async () => {
  loading.value = true
  try {
    // 只传递非空参数
    const params = { ...searchForm }
    Object.keys(params).forEach(key => {
      if (params[key] === '') {
        delete params[key]
      }
    })

    const response = await queryTransactionList(params)
    tableData.value = response.transactions || []
    total.value = response.total || 0
  } catch (error) {
    console.error('获取交易历史失败:', error)
    ElMessage.error('获取交易历史失败')
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  fetchTransactionList()
}

// 查看交易详情
const viewDetail = (row) => {
  router.push(`/transaction/${row.transactionUUID}`)
}

// 审核交易
const checkTransaction = (row) => {
  router.push(`/transaction/audit/${row.transactionUUID}`)
}

// 价格格式化
const formatPrice = (price) => {
  return `¥ ${parseFloat(price).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`
}

// 日期格式化
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'REJECTED': 'danger',
    'IN_PROGRESS': 'warning',
    'COMPLETED': 'success',
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'PENDING': '待处理',
    'REJECTED': '已拒绝',
    'IN_PROGRESS': '进行中',
    'COMPLETED': '已完成',
  }
  return statusMap[status] || status
}

// 判断当前用户是否可以审核交易
const canCheckTransaction = (transaction) => {
  if (!userInfo.value) return false
  
  // 政府组织且交易状态为待审核
  if (userInfo.value.role === 'GOVERNMENT' && transaction.status === 'PENDING') {
    return true
  }
  
  return false
}

// 判断当前用户是否为交易的买方或卖方
const isTransactionParty = (transaction) => {
  if (!userInfo.value || !userInfo.value.citizenID) return false
  
  // 获取用户身份证号的SHA256哈希值
  const userCitizenIDHash = CryptoJS.SHA256(userInfo.value.citizenID).toString(CryptoJS.enc.Hex)
  
  // 检查用户是否为买方或卖方
  const isBuyer = transaction.buyerCitizenIDHash === userCitizenIDHash
  const isSeller = transaction.sellerCitizenIDHash === userCitizenIDHash
  
  // 政府和审计可以查看所有交易
  const isAdmin = userInfo.value.organization === 'government' || userInfo.value.organization === 'audit'
  
  return isBuyer || isSeller || isAdmin
}

onMounted(() => {
  fetchTransactionList()
})
</script>

<style scoped>
.transaction-list-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
