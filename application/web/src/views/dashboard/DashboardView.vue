<template>
  <div class="dashboard-container">
    <div class="welcome-section">
      <h1>欢迎使用GRETS系统</h1>
      <p>基于区块链的房地产交易系统</p>
    </div>
    
    <div class="stats-section">
      <el-row :gutter="20">
        <!-- 统计卡片 -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><House /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.realtyCount }}</div>
              <div class="stat-label">房产总数</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><Sell /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.transactionCount }}</div>
              <div class="stat-label">交易总数</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.contractCount }}</div>
              <div class="stat-label">合同总数</div>
            </div>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <div class="stat-card">
            <div class="stat-icon">
              <el-icon><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ formatCurrency(stats.totalAmount) }}</div>
              <div class="stat-label">成交总额</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    
    <div class="quick-access">
      <h2>快速访问</h2>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="6" v-for="menu in quickAccessMenus" :key="menu.path">
          <div class="quick-card" @click="router.push(menu.path)">
            <el-icon>
              <component :is="menu.icon"></component>
            </el-icon>
            <span>{{ menu.title }}</span>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user.ts'
import { House, Sell, Document, Money, CreditCard, DataAnalysis } from '@element-plus/icons-vue'
import {queryContractList} from "@/api/contract.ts";
import {queryRealtyList} from "@/api/realty.ts";
import {queryTransactionList} from "@/api/transaction.ts";
import {getTotalPaymentAmount} from "@/api/payment.ts";

const router = useRouter()
const userStore = useUserStore()

// 统计数据
interface Stats {
  realtyCount: number
  transactionCount: number
  contractCount: number
  totalAmount: number
}

const stats = reactive<Stats>({
  realtyCount: 0,
  transactionCount: 0,
  contractCount: 0,
  totalAmount: 0
})

// 快速访问菜单
interface QuickMenu {
  title: string
  path: string
  icon: string
  organizations: string[]
}

const menus: QuickMenu[] = [
  { title: '房产列表', path: '/realty', icon: 'House', organizations: ['government', 'investor', 'bank', 'audit'] },
  { title: '交易列表', path: '/transaction', icon: 'Sell', organizations: ['government', 'investor', 'bank', 'audit'] },
  { title: '合同列表', path: '/contract', icon: 'Document', organizations: ['government', 'investor', 'bank', 'audit'] },
  { title: '支付列表', path: '/payment', icon: 'Money', organizations: ['investor', 'bank'] },
  { title: '统计分析', path: '/statistics/transaction', icon: 'DataAnalysis', organizations: ['government'] }
]

// 根据用户组织过滤菜单
const quickAccessMenus = computed(() => {
  return menus.filter(menu => {
    // 检查用户组织是否有权限访问该菜单
    return menu.organizations.includes(userStore.organization)
  })
})

// 格式化货币
const formatCurrency = (value: number): string => {
  return value.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' })
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const contractResponse = await queryContractList({
      pageNumber: 1,
      pageSize: 1000,
    })
    stats.contractCount = contractResponse.total

    const realtyResponse = await queryRealtyList({
      pageNumber: 1,
      pageSize: 1000,
    })
    stats.realtyCount = realtyResponse.total

    const transactionResponse = await queryTransactionList({
      pageNumber: 1,
      pageSize: 1000,
    })
    stats.transactionCount = transactionResponse.total

    const totalAmount = await getTotalPaymentAmount()
    stats.totalAmount = totalAmount.totalAmount
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.welcome-section {
  text-align: center;
  margin-bottom: 40px;
}

.welcome-section h1 {
  font-size: 28px;
  margin-bottom: 10px;
  color: #409EFF;
}

.welcome-section p {
  font-size: 16px;
  color: #606266;
}

.stats-section {
  margin-bottom: 30px;
}

.stat-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  height: 100px;
}

.stat-icon {
  font-size: 40px;
  margin-right: 20px;
  color: #409EFF;
}

.stat-info {
  flex-grow: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.quick-access {
  margin-top: 20px;
}

.quick-access h2 {
  font-size: 18px;
  margin-bottom: 20px;
  color: #303133;
}

.quick-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 120px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.quick-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.2);
}

.quick-card .el-icon {
  font-size: 32px;
  margin-bottom: 10px;
  color: #409EFF;
}

.quick-card span {
  font-size: 16px;
}

@media (max-width: 768px) {
  .stat-card {
    height: auto;
  }
}
</style> 