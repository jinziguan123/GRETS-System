<template>
  <div class="dashboard">
    <h1 class="page-title">系统仪表盘</h1>
    
    <!-- 欢迎卡片 -->
    <el-card class="welcome-card">
      <template #header>
        <div class="welcome-header">
          <h2>欢迎使用房地产交易系统</h2>
          <el-tag type="success">{{ userRole }}</el-tag>
        </div>
      </template>
      <div class="welcome-content">
        <p>您好，{{ username }}！今天是 {{ currentDate }}，祝您工作愉快。</p>
        <p>当前系统版本：v1.0.0</p>
      </div>
    </el-card>
    
    <!-- 数据统计卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #409EFF;">
            <el-icon><House /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-title">房产总数</div>
            <div class="stat-value">{{ stats.realtyCount }}</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #67C23A;">
            <el-icon><Sell /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-title">交易总数</div>
            <div class="stat-value">{{ stats.transactionCount }}</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #E6A23C;">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-title">合同总数</div>
            <div class="stat-value">{{ stats.contractCount }}</div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background-color: #F56C6C;">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-title">用户总数</div>
            <div class="stat-value">{{ stats.userCount }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 最近交易 -->
    <el-card class="recent-card">
      <template #header>
        <div class="card-header">
          <h3>最近交易</h3>
          <router-link to="/transaction">
            <el-button type="primary" text>查看全部</el-button>
          </router-link>
        </div>
      </template>
      <el-table :data="recentTransactions" style="width: 100%">
        <el-table-column prop="id" label="交易ID" width="180" />
        <el-table-column prop="realtyId" label="房产ID" width="180" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="scope">
            {{ formatCurrency(scope.row.amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" />
      </el-table>
    </el-card>
    
    <!-- 最近房产 -->
    <el-card class="recent-card">
      <template #header>
        <div class="card-header">
          <h3>最近房产</h3>
          <router-link to="/realty">
            <el-button type="primary" text>查看全部</el-button>
          </router-link>
        </div>
      </template>
      <el-table :data="recentRealties" style="width: 100%">
        <el-table-column prop="id" label="房产ID" width="180" />
        <el-table-column prop="location" label="位置" width="220" />
        <el-table-column prop="area" label="面积" width="120">
          <template #default="scope">
            {{ scope.row.area }} m²
          </template>
        </el-table-column>
        <el-table-column prop="totalPrice" label="总价" width="120">
          <template #default="scope">
            {{ formatCurrency(scope.row.totalPrice) }}
          </template>
        </el-table-column>
        <el-table-column prop="realtyStatus" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getRealtyStatusType(scope.row.realtyStatus)">
              {{ scope.row.realtyStatus }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { House, Sell, Document, User } from '@element-plus/icons-vue'
import api from '@/api'

const userStore = useUserStore()

// 用户信息
const username = computed(() => userStore.username)
const userRole = computed(() => {
  const roleMap = {
    'AdminMSP': '系统管理员',
    'GovernmentMSP': '政府机构',
    'AgencyMSP': '中介机构',
    'AuditMSP': '审计机构',
    'BankMSP': '银行',
    'InvestorMSP': '投资者'
  }
  return roleMap[userStore.userRole] || userStore.userRole
})

// 当前日期
const currentDate = computed(() => {
  const date = new Date()
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

// 统计数据
const stats = ref({
  realtyCount: 0,
  transactionCount: 0,
  contractCount: 0,
  userCount: 0
})

// 最近交易
const recentTransactions = ref([])

// 最近房产
const recentRealties = ref([])

// 格式化货币
const formatCurrency = (value) => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency: 'CNY'
  }).format(value)
}

// 获取交易状态类型
const getStatusType = (status) => {
  const statusMap = {
    'PENDING': 'warning',
    'APPROVED': 'success',
    'REJECTED': 'danger',
    'COMPLETED': 'primary'
  }
  return statusMap[status] || 'info'
}

// 获取房产状态类型
const getRealtyStatusType = (status) => {
  const statusMap = {
    'AVAILABLE': 'success',
    'SOLD': 'info',
    'PENDING': 'warning',
    'UNAVAILABLE': 'danger'
  }
  return statusMap[status] || 'info'
}

// 加载仪表盘数据
const loadDashboardData = async () => {
  try {
    // 模拟数据，实际项目中应该从API获取
    stats.value = {
      realtyCount: 128,
      transactionCount: 56,
      contractCount: 42,
      userCount: 35
    }
    
    // 获取最近交易
    const transactionResponse = await api.transaction.getList({
      pageSize: 5,
      pageNumber: 1
    })
    recentTransactions.value = transactionResponse.data.list || []
    
    // 获取最近房产
    const realtyResponse = await api.realty.getList({
      pageSize: 5,
      pageNumber: 1
    })
    recentRealties.value = realtyResponse.data.list || []
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    // 使用模拟数据
    recentTransactions.value = [
      {
        id: 'TX20230501001',
        realtyId: 'R20230501001',
        amount: 1500000,
        status: 'COMPLETED',
        createTime: '2023-05-01 10:30:00'
      },
      {
        id: 'TX20230502002',
        realtyId: 'R20230502002',
        amount: 2300000,
        status: 'PENDING',
        createTime: '2023-05-02 14:20:00'
      },
      {
        id: 'TX20230503003',
        realtyId: 'R20230503003',
        amount: 1800000,
        status: 'APPROVED',
        createTime: '2023-05-03 09:15:00'
      }
    ]
    
    recentRealties.value = [
      {
        id: 'R20230501001',
        location: '北京市朝阳区建国路89号',
        area: 120,
        totalPrice: 1500000,
        realtyStatus: 'SOLD'
      },
      {
        id: 'R20230502002',
        location: '上海市浦东新区陆家嘴环路1000号',
        area: 150,
        totalPrice: 2300000,
        realtyStatus: 'AVAILABLE'
      },
      {
        id: 'R20230503003',
        location: '广州市天河区天河路385号',
        area: 100,
        totalPrice: 1800000,
        realtyStatus: 'PENDING'
      }
    ]
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.page-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #303133;
}

.welcome-card {
  margin-bottom: 20px;
}

.welcome-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-header h2 {
  margin: 0;
  font-size: 18px;
}

.welcome-content {
  line-height: 1.8;
}

.stat-row {
  margin-bottom: 20px;
}

.stat-card {
  height: 100px;
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 15px;
  color: white;
  font-size: 24px;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.recent-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
}
</style> 