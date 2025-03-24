<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>系统仪表盘</h1>
      <p>欢迎回来，{{ username }}。您当前的角色是：{{ userRoleDisplay }}</p>
    </div>
    
    <el-row :gutter="20">
      <!-- 统计卡片 -->
      <el-col :xs="24" :sm="12" :md="6" v-for="(card, index) in statisticCards" :key="index">
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
      <!-- 交易趋势图 -->
      <el-col :xs="24" :lg="16">
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
      
      <!-- 分布饼图 -->
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
      <!-- 最近交易 -->
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover" class="latest-card">
          <template #header>
            <div class="card-header">
              <span>最近交易</span>
              <el-button type="primary" link @click="$router.push('/transaction/list')">
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
      
      <!-- 最近合同 -->
      <el-col :xs="24" :lg="12">
        <el-card shadow="hover" class="latest-card">
          <template #header>
            <div class="card-header">
              <span>最近合同</span>
              <el-button type="primary" link @click="$router.push('/contract/list')">
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
  Wallet
} from '@element-plus/icons-vue'

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
const userRoleDisplay = computed(() => {
  const roleMap = {
    'AdminMSP': '系统管理员',
    'GovernmentMSP': '政府监管',
    'AgencyMSP': '房产中介',
    'InvestorMSP': '投资者',
    'BankMSP': '银行机构'
  }
  return roleMap[userStore.userRole] || userStore.userRole
})

// 统计卡片数据
const statisticCards = ref([
  { title: '房产总数', value: 0, icon: 'House', color: '#1976d2' },
  { title: '交易总数', value: 0, icon: 'Money', color: '#ff9800' },
  { title: '合同总数', value: 0, icon: 'Document', color: '#4caf50' },
  { title: '税费总额', value: '¥0', icon: 'Wallet', color: '#e91e63' }
])

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
    const dashboardData = await dashboardApi.getDashboardData()
    updateDashboardData(dashboardData)
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
    // 使用模拟数据作为后备
    updateDashboardData(getMockData())
  }
  
  // 初始化图表
  initTransactionTrendChart()
  initRealtyDistributionChart()
})

// 更新仪表盘数据
const updateDashboardData = (data) => {
  // 更新统计卡片
  statisticCards.value[0].value = data.statistics.realtyCount
  statisticCards.value[1].value = data.statistics.transactionCount
  statisticCards.value[2].value = data.statistics.contractCount
  statisticCards.value[3].value = `¥${data.statistics.taxAmount.toLocaleString()}`
  
  // 更新最近交易和合同
  latestTransactions.value = data.latestTransactions
  latestContracts.value = data.latestContracts
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

// 模拟数据（作为后备）
const getMockData = () => {
  return {
    statistics: {
      realtyCount: 127,
      transactionCount: 85,
      contractCount: 94,
      taxAmount: 1358000
    },
    latestTransactions: [
      { id: 'T-1001', realtyName: '翠湖豪庭 3号楼605', amount: 1580000, status: '完成' },
      { id: 'T-1002', realtyName: '阳光花园 12栋1203', amount: 2450000, status: '进行中' },
      { id: 'T-1003', realtyName: '蓝天公寓 B座502', amount: 1260000, status: '完成' },
      { id: 'T-1004', realtyName: '江南名府 6号楼1801', amount: 3280000, status: '待处理' }
    ],
    latestContracts: [
      { id: 'C-2001', title: '翠湖豪庭605购房合同', date: '2024-05-12', status: '完成' },
      { id: 'C-2002', title: '阳光花园1203认购协议', date: '2024-05-10', status: '进行中' },
      { id: 'C-2003', title: '蓝天公寓502过户协议', date: '2024-05-08', status: '完成' },
      { id: 'C-2004', title: '江南名府1801租赁合同', date: '2024-05-05', status: '待处理' }
    ]
  }
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
}
</style> 