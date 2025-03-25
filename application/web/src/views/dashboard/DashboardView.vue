<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>系统仪表盘</h1>
      <p>欢迎回来，{{ username }}。您的组织类型：
        {{ organizationName }}
      </p>
    </div>
    
    <!-- 功能向导卡片 -->
    <el-card class="guide-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>快速开始</span>
        </div>
      </template>
      <div class="guide-content">
        <div v-if="userStore.hasOrganization('investor')" class="guide-section">
          <h3>投资者功能</h3>
          <el-button-group>
            <el-button type="primary" @click="$router.push('/realty')">浏览房产</el-button>
            <el-button type="primary" @click="$router.push('/transaction/create')">创建交易</el-button>
            <el-button type="primary" @click="$router.push('/contract/create')">创建合同</el-button>
            <el-button type="primary" @click="$router.push('/mortgage/create')">申请贷款</el-button>
          </el-button-group>
        </div>
        
        <div v-if="userStore.hasOrganization('government')" class="guide-section">
          <h3>政府监管功能</h3>
          <el-button-group>
            <el-button type="success" @click="$router.push('/realty/create')">添加房产</el-button>
            <el-button type="success" @click="$router.push('/tax/create')">添加税费</el-button>
            <el-button type="success" @click="$router.push('/statistics/transaction')">统计数据</el-button>
          </el-button-group>
        </div>
        
        <div v-if="userStore.hasOrganization('bank')" class="guide-section">
          <h3>银行功能</h3>
          <el-button-group>
            <el-button type="danger" @click="$router.push('/mortgage/approve')">贷款审批</el-button>
            <el-button type="danger" @click="$router.push('/payment')">支付管理</el-button>
            <el-button type="danger" @click="$router.push('/statistics/loan')">贷款统计</el-button>
          </el-button-group>
        </div>
        
        <div v-if="userStore.hasOrganization('audit')" class="guide-section">
          <h3>审计功能</h3>
          <el-button-group>
            <el-button type="warning" @click="$router.push('/contract/audit')">合同审核</el-button>
            <el-button type="warning" @click="$router.push('/realty')">房产查询</el-button>
          </el-button-group>
        </div>
        
        <div v-if="userStore.hasOrganization('administrator')" class="guide-section">
          <h3>管理员功能</h3>
          <el-button-group>
            <el-button type="info" @click="$router.push('/admin/users')">用户管理</el-button>
            <el-button type="info" @click="$router.push('/admin/system')">系统设置</el-button>
          </el-button-group>
        </div>
      </div>
    </el-card>
    
    <el-row :gutter="20">
      <!-- 统计卡片 -->
      <el-col :xs="24" :sm="12" :md="6" v-for="(card, index) in filteredStatisticCards" :key="index">
        <el-card class="statistic-card" shadow="hover">
          <div class="card-content">
            <div class="card-icon" :style="{ backgroundColor: card.color }">
              <el-icon><component :is="card.icon"></component></el-icon>
            </div>
            <div class="card-info">
              <div class="card-title">{{ card.title }}</div>
              <div class="card-value">{{ card.value }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="chart-row">
      <!-- 交易趋势图 - 对投资者、政府和银行显示 -->
      <el-col :xs="24" :lg="16" v-if="userStore.hasOrganization(['investor', 'government', 'bank'])">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>交易趋势</span>
              <el-radio-group v-model="transactionTimeRange" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">全年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <div ref="transactionTrendChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 分布饼图 - 所有组织都显示 -->
      <el-col :xs="24" :lg="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>房产状态分布</span>
            </div>
          </template>
          <div class="chart-container">
            <div ref="realtyDistributionChart" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20">
      <!-- 最近交易 - 对投资者、政府和银行显示 -->
      <el-col :xs="24" :lg="12" v-if="userStore.hasOrganization(['investor', 'government', 'bank'])">
        <el-card shadow="hover" class="latest-card">
          <template #header>
            <div class="card-header">
              <span>最近交易</span>
              <el-button type="primary" link @click="$router.push('/transaction')">
                查看全部
              </el-button>
            </div>
          </template>
          <el-table :data="latestTransactions" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="realtyName" label="房产名称" />
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="scope">¥{{ scope.row.amount.toLocaleString() }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <!-- 最近合同 - 对投资者、政府、银行和审计显示 -->
      <el-col :xs="24" :lg="12" v-if="userStore.hasOrganization(['investor', 'government', 'bank', 'audit'])">
        <el-card shadow="hover" class="latest-card">
          <template #header>
            <div class="card-header">
              <span>最近合同</span>
              <el-button type="primary" link @click="$router.push('/contract')">
                查看全部
              </el-button>
            </div>
          </template>
          <el-table :data="latestContracts" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="title" label="合同标题" />
            <el-table-column prop="date" label="签约日期" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import * as echarts from 'echarts/core'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import { LineChart, PieChart } from 'echarts/charts'
import { UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import {
  House,
  Document,
  Money,
  Wallet,
  CreditCard,
  DataAnalysis
} from '@element-plus/icons-vue'
import axios from 'axios'

// 注册必要的ECharts组件
echarts.use([
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  LineChart,
  PieChart,
  CanvasRenderer,
  UniversalTransition
])

const userStore = useUserStore()
const username = computed(() => userStore.username)

// 组织相关计算属性
const orgClass = computed(() => `org-${userStore.organization.value}`)
// 组织名称和样式
const userOrganization = computed(() => userStore.organization || '')
const organizationName = computed(() => {
  const orgMap = {
    'administrator': '系统管理员',
    'government': '政府监管部门',
    'investor': '投资者/买家',
    'bank': '银行机构',
    'audit': '审计监管部门'
  }
  return orgMap[userOrganization.value] || '未知组织'
})

// 统计卡片数据
const statisticCards = ref([
  { title: '房产总数', value: 0, icon: 'House', color: '#1976d2', orgs: ['investor', 'government', 'audit'] },
  { title: '交易总数', value: 0, icon: 'Money', color: '#ff9800', orgs: ['investor', 'government', 'bank'] },
  { title: '合同总数', value: 0, icon: 'Document', color: '#4caf50', orgs: ['investor', 'government', 'bank', 'audit'] },
  { title: '税费总额', value: '¥0', icon: 'Wallet', color: '#e91e63', orgs: ['government', 'investor'] },
  { title: '贷款总数', value: 0, icon: 'CreditCard', color: '#9c27b0', orgs: ['investor', 'bank'] },
  { title: '审计总数', value: 0, icon: 'DataAnalysis', color: '#795548', orgs: ['audit'] }
])

// 根据用户组织过滤统计卡片
const filteredStatisticCards = computed(() => {
  return statisticCards.value.filter(card => {
    return card.orgs.includes(userStore.organization.value)
  })
})

// 图表引用
const transactionTrendChart = ref(null)
const realtyDistributionChart = ref(null)

// 交易时间范围
const transactionTimeRange = ref('month')

// 最近交易和合同
const latestTransactions = ref([])
const latestContracts = ref([])

// 初始化
onMounted(async () => {
  // 获取仪表盘数据
  try {
    const response = await axios.get('/api/dashboard')
    if (response.data.code === 200) {
      updateDashboardData(response.data.data)
    } else {
      throw new Error(response.data.message || '获取数据失败')
    }
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
    initDefaultData()
  }
  
  // 初始化图表
  initTransactionTrendChart()
  initRealtyDistributionChart()
})

// 更新仪表盘数据
const updateDashboardData = (data) => {
  // 更新统计卡片
  if (data.statistics) {
    statisticCards.value[0].value = data.statistics.realtyCount || 0
    statisticCards.value[1].value = data.statistics.transactionCount || 0
    statisticCards.value[2].value = data.statistics.contractCount || 0
    statisticCards.value[3].value = `¥${(data.statistics.taxAmount || 0).toLocaleString()}`
    statisticCards.value[4].value = data.statistics.mortgageCount || 0
    statisticCards.value[5].value = data.statistics.auditCount || 0
  }
  
  // 更新最近交易和合同
  if (data.latestTransactions) {
    latestTransactions.value = data.latestTransactions
  }
  
  if (data.latestContracts) {
    latestContracts.value = data.latestContracts
  }
}

// 初始化默认数据
const initDefaultData = () => {
  // 根据用户组织设置不同的默认数据
  if (userStore.hasOrganization('investor')) {
    statisticCards.value[0].value = 15
    statisticCards.value[1].value = 8
    statisticCards.value[2].value = 6
    statisticCards.value[3].value = '¥120,000'
    statisticCards.value[4].value = 3
  } else if (userStore.hasOrganization('government')) {
    statisticCards.value[0].value = 127
    statisticCards.value[1].value = 85
    statisticCards.value[2].value = 94
    statisticCards.value[3].value = '¥1,358,000'
  } else if (userStore.hasOrganization('bank')) {
    statisticCards.value[1].value = 56
    statisticCards.value[2].value = 42
    statisticCards.value[4].value = 38
  } else if (userStore.hasOrganization('audit')) {
    statisticCards.value[0].value = 127
    statisticCards.value[2].value = 94
    statisticCards.value[5].value = 48
  }
  
  // 设置默认的最近交易
  latestTransactions.value = [
    { id: 'T-1001', realtyName: '翠湖豪庭 3号楼605', amount: 1580000, status: '完成' },
    { id: 'T-1002', realtyName: '阳光花园 12栋1203', amount: 2450000, status: '进行中' },
    { id: 'T-1003', realtyName: '蓝天公寓 B座502', amount: 1260000, status: '完成' },
    { id: 'T-1004', realtyName: '江南名府 6号楼1801', amount: 3280000, status: '待处理' }
  ]
  
  // 设置默认的最近合同
  latestContracts.value = [
    { id: 'C-2001', title: '翠湖豪庭605购房合同', date: '2024-05-12', status: '完成' },
    { id: 'C-2002', title: '阳光花园1203认购协议', date: '2024-05-10', status: '进行中' },
    { id: 'C-2003', title: '蓝天公寓502过户协议', date: '2024-05-08', status: '完成' },
    { id: 'C-2004', title: '江南名府1801租赁合同', date: '2024-05-05', status: '待处理' }
  ]
}

// 初始化交易趋势图
const initTransactionTrendChart = () => {
  if (!transactionTrendChart.value) return
  
  const chart = echarts.init(transactionTrendChart.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: function(params) {
        const data = params[0]
        return `${data.name}<br/>${data.seriesName}: ¥${data.value.toLocaleString()}`
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月'],
      axisLine: {
        lineStyle: {
          color: '#ddd'
        }
      },
      axisLabel: {
        color: '#666'
      }
    },
    yAxis: {
      type: 'value',
      axisLine: {
        show: false
      },
      axisLabel: {
        color: '#666',
        formatter: function(value) {
          if (value >= 10000) {
            return (value / 10000) + '万'
          }
          return value
        }
      },
      splitLine: {
        lineStyle: {
          color: '#eee'
        }
      }
    },
    series: [
      {
        name: '交易金额',
        type: 'line',
        data: [950000, 1250000, 1800000, 1400000, 2200000, 1800000, 2500000, 2100000, 1900000, 2300000, 2800000, 3100000],
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        itemStyle: {
          color: '#1976d2'
        },
        lineStyle: {
          width: 3,
          color: '#1976d2'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(25, 118, 210, 0.3)' },
              { offset: 1, color: 'rgba(25, 118, 210, 0.1)' }
            ]
          }
        }
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应窗口大小变化
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 初始化房产分布图
const initRealtyDistributionChart = () => {
  if (!realtyDistributionChart.value) return
  
  const chart = echarts.init(realtyDistributionChart.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      data: ['可出售', '已售出', '抵押中', '冻结中', '其他']
    },
    series: [
      {
        name: '房产状态',
        type: 'pie',
        radius: ['50%', '70%'],
        center: ['40%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 4,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 45, name: '可出售', itemStyle: { color: '#4CAF50' } },
          { value: 32, name: '已售出', itemStyle: { color: '#1976D2' } },
          { value: 15, name: '抵押中', itemStyle: { color: '#FF9800' } },
          { value: 5, name: '冻结中', itemStyle: { color: '#F44336' } },
          { value: 3, name: '其他', itemStyle: { color: '#9E9E9E' } }
        ]
      }
    ]
  }
  
  chart.setOption(option)
  
  // 响应窗口大小变化
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

// 获取状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    '完成': 'success',
    '进行中': 'primary',
    '待处理': 'info',
    '已取消': 'danger'
  }
  return statusMap[status] || 'info'
}
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.dashboard-header {
  margin-bottom: 20px;
}

.dashboard-header h1 {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 10px;
}

.dashboard-header p {
  margin: 0;
  color: #666;
}

.org-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  color: #fff;
  font-weight: bold;
  margin-left: 5px;
}

.org-administrator {
  background-color: #409EFF;
}

.org-government {
  background-color: #67C23A;
}

.org-investor {
  background-color: #E6A23C;
}

.org-bank {
  background-color: #F56C6C;
}

.org-audit {
  background-color: #909399;
}

.guide-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.guide-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.guide-section h3 {
  margin-top: 0;
  margin-bottom: 10px;
}

.statistic-card {
  margin-bottom: 20px;
  border: none;
  border-radius: 8px;
}

.card-content {
  display: flex;
  align-items: center;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 16px;
}

.card-icon :deep(.el-icon) {
  font-size: 24px;
  color: #fff;
}

.card-info {
  flex: 1;
}

.card-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.chart-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  width: 100%;
  height: 300px;
}

.chart {
  width: 100%;
  height: 100%;
}

.latest-card {
  margin-bottom: 20px;
}

.el-table {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 15px;
  }
  
  .chart-container {
    height: 250px;
  }
  
  .guide-content {
    flex-direction: column;
  }
  
  .guide-section {
    width: 100%;
    margin-bottom: 15px;
  }
  
  .el-button-group {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
  
  .el-button-group .el-button {
    margin-left: 0 !important;
    border-radius: 4px !important;
  }
}
</style> 