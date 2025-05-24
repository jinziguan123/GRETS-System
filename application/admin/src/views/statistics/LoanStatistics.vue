<template>
  <div class="loan-statistics-container">
    <div class="page-header">
      <h2>贷款业务统计</h2>
      <el-button-group>
        <el-button :type="timeRange === 'month' ? 'primary' : 'default'" @click="changeTimeRange('month')">
          本月
        </el-button>
        <el-button :type="timeRange === 'quarter' ? 'primary' : 'default'" @click="changeTimeRange('quarter')">
          本季度
        </el-button>
        <el-button :type="timeRange === 'year' ? 'primary' : 'default'" @click="changeTimeRange('year')">
          本年度
        </el-button>
        <el-button :type="timeRange === 'custom' ? 'primary' : 'default'" @click="changeTimeRange('custom')">
          自定义
        </el-button>
      </el-button-group>
    </div>
    
    <el-card v-if="timeRange === 'custom'" class="filter-card">
      <el-form :inline="true" size="small">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item label="贷款类型">
          <el-select v-model="loanType" placeholder="选择贷款类型" clearable @change="fetchData">
            <el-option 
              v-for="item in loanTypeOptions" 
              :key="item.value" 
              :label="item.label" 
              :value="item.value" 
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>
    
    <template v-else>
      <!-- 数据概览 -->
      <el-row :gutter="20" class="stats-overview">
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>贷款总笔数</span>
                <el-tooltip content="统计周期内的新增贷款总笔数" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">{{ formatNumber(overviewData.totalLoans) }}</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.totalLoansTrend > 0, 'trend-down': overviewData.totalLoansTrend < 0}">
                <el-icon><component :is="overviewData.totalLoansTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.totalLoansTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>贷款总金额</span>
                <el-tooltip content="统计周期内的新增贷款总金额" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">¥{{ formatNumber(overviewData.totalAmount / 10000) }}万</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.totalAmountTrend > 0, 'trend-down': overviewData.totalAmountTrend < 0}">
                <el-icon><component :is="overviewData.totalAmountTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.totalAmountTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>平均贷款额</span>
                <el-tooltip content="统计周期内的平均每笔贷款金额" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">¥{{ formatNumber(overviewData.averageLoanAmount / 10000) }}万</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.averageLoanAmountTrend > 0, 'trend-down': overviewData.averageLoanAmountTrend < 0}">
                <el-icon><component :is="overviewData.averageLoanAmountTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.averageLoanAmountTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>贷款利息收入</span>
                <el-tooltip content="统计周期内的贷款预计利息收入总额" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">¥{{ formatNumber(overviewData.totalInterest / 10000) }}万</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.totalInterestTrend > 0, 'trend-down': overviewData.totalInterestTrend < 0}">
                <el-icon><component :is="overviewData.totalInterestTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.totalInterestTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- 图表展示 -->
      <el-row :gutter="20" class="chart-row">
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>贷款趋势</span>
                <el-radio-group v-model="loanTrendType" size="small" @change="updateLoanTrendChart">
                  <el-radio-button label="count">笔数</el-radio-button>
                  <el-radio-button label="amount">金额</el-radio-button>
                </el-radio-group>
              </div>
            </template>
            <div class="chart-container" ref="loanTrendChart"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>贷款类型分布</span>
              </div>
            </template>
            <div class="chart-container" ref="loanTypeChart"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" class="chart-row">
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>贷款期限分布</span>
              </div>
            </template>
            <div class="chart-container" ref="loanTermChart"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>贷款金额区间分布</span>
              </div>
            </template>
            <div class="chart-container" ref="loanAmountRangeChart"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- 贷款风险指标 -->
      <el-card class="risk-card">
        <template #header>
          <div class="card-header">
            <span>贷款风险指标</span>
          </div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-progress 
              type="dashboard" 
              :percentage="riskData.approvalRate" 
              :color="getApprovalRateColor" 
              :stroke-width="10"
            >
              <template #default>
                <div class="progress-content">
                  <div class="progress-value">{{ riskData.approvalRate }}%</div>
                  <div class="progress-label">贷款审批通过率</div>
                </div>
              </template>
            </el-progress>
          </el-col>
          <el-col :span="8">
            <el-progress 
              type="dashboard" 
              :percentage="riskData.overdueRate" 
              :color="getOverdueRateColor" 
              :stroke-width="10"
            >
              <template #default>
                <div class="progress-content">
                  <div class="progress-value">{{ riskData.overdueRate }}%</div>
                  <div class="progress-label">贷款逾期率</div>
                </div>
              </template>
            </el-progress>
          </el-col>
          <el-col :span="8">
            <el-progress 
              type="dashboard" 
              :percentage="riskData.loanToValueRatio" 
              :color="getLTVColor" 
              :stroke-width="10"
            >
              <template #default>
                <div class="progress-content">
                  <div class="progress-value">{{ riskData.loanToValueRatio }}%</div>
                  <div class="progress-label">平均贷款成数</div>
                </div>
              </template>
            </el-progress>
          </el-col>
        </el-row>
      </el-card>
      
      <!-- 贷款排行榜 -->
      <el-card class="ranking-card">
        <template #header>
          <div class="card-header">
            <span>贷款排行榜</span>
            <el-tabs v-model="rankingType" @tab-change="handleRankingTypeChange">
              <el-tab-pane label="按贷款金额" name="amount"></el-tab-pane>
              <el-tab-pane label="按贷款期限" name="term"></el-tab-pane>
              <el-tab-pane label="按贷款利率" name="rate"></el-tab-pane>
            </el-tabs>
          </div>
        </template>
        <el-table :data="rankingData" style="width: 100%">
          <el-table-column type="index" width="50" />
          <el-table-column prop="borrowerName" label="借款人" width="120">
            <template #default="scope">
              {{ scope.row.borrowerName.replace(/(?<=.).(?=.)/g, '*') }}
            </template>
          </el-table-column>
          <el-table-column prop="propertyInfo" label="抵押房产信息" min-width="180">
            <template #default="scope">
              <div>{{ scope.row.district }} | {{ scope.row.propertyInfo }}</div>
              <div class="property-meta">{{ scope.row.roomType }} | {{ scope.row.area }}㎡</div>
            </template>
          </el-table-column>
          <el-table-column prop="loanAmount" label="贷款金额" sortable width="150">
            <template #default="scope">
              ¥{{ formatNumber(scope.row.loanAmount / 10000) }}万
            </template>
          </el-table-column>
          <el-table-column prop="loanTerm" label="贷款期限" sortable width="100">
            <template #default="scope">
              {{ scope.row.loanTerm }}年
            </template>
          </el-table-column>
          <el-table-column prop="interestRate" label="贷款利率" sortable width="100">
            <template #default="scope">
              {{ scope.row.interestRate }}%
            </template>
          </el-table-column>
          <el-table-column prop="loanDate" label="贷款日期" width="120">
            <template #default="scope">
              {{ formatDate(scope.row.loanDate) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button type="text" @click="viewLoanDetail(scope.row.id)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { QuestionFilled, CaretTop, CaretBottom } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const router = useRouter()

// 加载状态
const loading = ref(false)

// 筛选条件
const timeRange = ref('month')
const dateRange = ref([])
const loanType = ref('')

// 图表引用
const loanTrendChart = ref(null)
const loanTypeChart = ref(null)
const loanTermChart = ref(null)
const loanAmountRangeChart = ref(null)

// 图表实例
let trendChartInstance = null
let typeChartInstance = null
let termChartInstance = null
let amountRangeChartInstance = null

// 图表类型
const loanTrendType = ref('count')

// 排行榜类型
const rankingType = ref('amount')
const rankingData = ref([])

// 贷款类型选项
const loanTypeOptions = [
  { value: 'firstHouse', label: '首套房贷款' },
  { value: 'secondHouse', label: '二套房贷款' },
  { value: 'commercialProperty', label: '商业用房贷款' },
  { value: 'publicHousingFund', label: '公积金贷款' },
  { value: 'combinedLoan', label: '组合贷款' }
]

// 数据概览
const overviewData = reactive({
  totalLoans: 0,
  totalLoansTrend: 0,
  totalAmount: 0,
  totalAmountTrend: 0,
  averageLoanAmount: 0,
  averageLoanAmountTrend: 0,
  totalInterest: 0,
  totalInterestTrend: 0
})

// 风险数据
const riskData = reactive({
  approvalRate: 0,
  overdueRate: 0,
  loanToValueRatio: 0
})

// 计算属性：审批通过率颜色
const getApprovalRateColor = computed(() => {
  if (riskData.approvalRate >= 80) return '#67C23A'
  if (riskData.approvalRate >= 60) return '#E6A23C'
  return '#F56C6C'
})

// 计算属性：逾期率颜色
const getOverdueRateColor = computed(() => {
  if (riskData.overdueRate <= 3) return '#67C23A'
  if (riskData.overdueRate <= 7) return '#E6A23C'
  return '#F56C6C'
})

// 计算属性：贷款成数颜色
const getLTVColor = computed(() => {
  if (riskData.loanToValueRatio <= 65) return '#67C23A'
  if (riskData.loanToValueRatio <= 80) return '#E6A23C'
  return '#F56C6C'
})

// 修改时间范围
const changeTimeRange = (range) => {
  timeRange.value = range
  setDefaultDateRange(range)
  fetchData()
}

// 设置默认日期范围
const setDefaultDateRange = (range) => {
  const now = dayjs()
  
  switch(range) {
    case 'month':
      dateRange.value = [now.startOf('month').toDate(), now.endOf('month').toDate()]
      break
    case 'quarter':
      dateRange.value = [now.startOf('quarter').toDate(), now.endOf('quarter').toDate()]
      break
    case 'year':
      dateRange.value = [now.startOf('year').toDate(), now.endOf('year').toDate()]
      break
    case 'custom':
      // 保持当前选择
      if (!dateRange.value.length) {
        dateRange.value = [now.subtract(1, 'month').toDate(), now.toDate()]
      }
      break
  }
}

// 日期范围改变
const handleDateRangeChange = () => {
  if (dateRange.value && dateRange.value.length === 2) {
    fetchData()
  }
}

// 重置筛选条件
const resetFilters = () => {
  loanType.value = ''
  setDefaultDateRange('custom')
  fetchData()
}

// 获取统计数据
const fetchData = async () => {
  loading.value = true
  
  try {
    // 这里替换为实际的API调用
    setTimeout(() => {
      // 模拟数据
      overviewData.totalLoans = 342
      overviewData.totalLoansTrend = 5.7
      overviewData.totalAmount = 2158000000
      overviewData.totalAmountTrend = 8.2
      overviewData.averageLoanAmount = overviewData.totalAmount / overviewData.totalLoans
      overviewData.averageLoanAmountTrend = 2.3
      overviewData.totalInterest = 412000000
      overviewData.totalInterestTrend = 3.8
      
      // 风险数据
      riskData.approvalRate = 76
      riskData.overdueRate = 2.5
      riskData.loanToValueRatio = 72
      
      // 加载排行榜数据
      loadRankingData()
      
      // 绘制图表
      nextTick(() => {
        initCharts()
      })
      
      loading.value = false
    }, 1000)
  } catch (error) {
    console.error('Failed to fetch loan statistics data:', error)
    ElMessage.error('获取贷款统计数据失败')
    loading.value = false
  }
}

// 加载排行榜数据
const loadRankingData = () => {
  // 这里替换为实际的API调用
  const districtOptions = [
    '浦东新区', '黄浦区', '徐汇区', '长宁区', '静安区', '普陀区', 
    '虹口区', '杨浦区', '闵行区', '宝山区', '嘉定区'
  ]
  
  const mockData = []
  for (let i = 1; i <= 10; i++) {
    const district = districtOptions[Math.floor(Math.random() * districtOptions.length)]
    const area = Math.floor(80 + Math.random() * 150)
    const loanAmount = Math.floor(1000000 + Math.random() * 5000000)
    const loanTerm = [10, 15, 20, 25, 30][Math.floor(Math.random() * 5)]
    const interestRate = (3.8 + Math.random() * 1.5).toFixed(2)
    
    mockData.push({
      id: `L${String(i).padStart(3, '0')}`,
      borrowerName: ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十'][Math.floor(Math.random() * 8)],
      district: district,
      propertyInfo: `${district}某小区${i}号楼${Math.floor(Math.random() * 30) + 1}层`,
      roomType: `${Math.floor(Math.random() * 3) + 2}室${Math.floor(Math.random() * 2) + 1}厅`,
      area: area,
      loanAmount: loanAmount,
      loanTerm: loanTerm,
      interestRate: interestRate,
      loanDate: new Date(Date.now() - Math.floor(Math.random() * 90) * 24 * 60 * 60 * 1000)
    })
  }
  
  // 根据排行类型排序
  if (rankingType.value === 'amount') {
    mockData.sort((a, b) => b.loanAmount - a.loanAmount)
  } else if (rankingType.value === 'term') {
    mockData.sort((a, b) => b.loanTerm - a.loanTerm)
  } else if (rankingType.value === 'rate') {
    mockData.sort((a, b) => b.interestRate - a.interestRate)
  }
  
  rankingData.value = mockData
}

// 初始化图表
const initCharts = () => {
  initLoanTrendChart()
  initLoanTypeChart()
  initLoanTermChart()
  initLoanAmountRangeChart()
}

// 初始化贷款趋势图表
const initLoanTrendChart = () => {
  if (!loanTrendChart.value) return
  
  if (trendChartInstance) {
    trendChartInstance.dispose()
  }
  
  trendChartInstance = echarts.init(loanTrendChart.value)
  
  // 模拟数据
  const months = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
  const countData = [35, 42, 28, 38, 32, 40, 45, 36, 32, 48, 52, 56]
  const amountData = [150, 180, 120, 200, 140, 190, 210, 175, 160, 230, 260, 280].map(v => v * 1000000)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        const param = params[0]
        if (loanTrendType.value === 'count') {
          return `${param.axisValue}<br/>${param.marker}贷款笔数: ${param.value}`
        } else {
          return `${param.axisValue}<br/>${param.marker}贷款金额: ¥${formatNumber(param.value / 10000)}万`
        }
      }
    },
    xAxis: {
      type: 'category',
      data: months
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: function(value) {
          if (loanTrendType.value === 'count') {
            return value
          } else {
            return `${value / 10000}万`
          }
        }
      }
    },
    series: [
      {
        name: loanTrendType.value === 'count' ? '贷款笔数' : '贷款金额',
        type: 'line',
        smooth: true,
        data: loanTrendType.value === 'count' ? countData : amountData,
        markPoint: {
          data: [
            { type: 'max', name: '最大值' },
            { type: 'min', name: '最小值' }
          ]
        },
        markLine: {
          data: [{ type: 'average', name: '平均值' }]
        },
        areaStyle: {}
      }
    ],
    color: ['#409EFF']
  }
  
  trendChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    trendChartInstance && trendChartInstance.resize()
  })
}

// 更新贷款趋势图表
const updateLoanTrendChart = () => {
  initLoanTrendChart()
}

// 初始化贷款类型分布图表
const initLoanTypeChart = () => {
  if (!loanTypeChart.value) return
  
  if (typeChartInstance) {
    typeChartInstance.dispose()
  }
  
  typeChartInstance = echarts.init(loanTypeChart.value)
  
  // 模拟数据
  const loanTypeData = [
    { value: 45, name: '首套房贷款' },
    { value: 25, name: '二套房贷款' },
    { value: 10, name: '商业用房贷款' },
    { value: 15, name: '公积金贷款' },
    { value: 5, name: '组合贷款' }
  ]
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      data: loanTypeData.map(item => item.name)
    },
    series: [
      {
        name: '贷款类型',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
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
            fontSize: 16,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: loanTypeData
      }
    ]
  }
  
  typeChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    typeChartInstance && typeChartInstance.resize()
  })
}

// 初始化贷款期限分布图表
const initLoanTermChart = () => {
  if (!loanTermChart.value) return
  
  if (termChartInstance) {
    termChartInstance.dispose()
  }
  
  termChartInstance = echarts.init(loanTermChart.value)
  
  // 模拟数据
  const termData = [
    { value: 10, name: '10年' },
    { value: 15, name: '15年' },
    { value: 30, name: '20年' },
    { value: 35, name: '25年' },
    { value: 10, name: '30年' }
  ]
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c}% ({d}%)'
    },
    series: [
      {
        name: '贷款期限',
        type: 'pie',
        radius: '70%',
        center: ['50%', '50%'],
        roseType: 'radius',
        itemStyle: {
          borderRadius: 5
        },
        label: {
          formatter: '{b}: {c}%'
        },
        data: termData
      }
    ]
  }
  
  termChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    termChartInstance && termChartInstance.resize()
  })
}

// 初始化贷款金额区间分布图表
const initLoanAmountRangeChart = () => {
  if (!loanAmountRangeChart.value) return
  
  if (amountRangeChartInstance) {
    amountRangeChartInstance.dispose()
  }
  
  amountRangeChartInstance = echarts.init(loanAmountRangeChart.value)
  
  // 模拟数据
  const amountRanges = ['100万以下', '100-200万', '200-300万', '300-400万', '400-500万', '500万以上']
  const amountData = [5, 15, 30, 25, 15, 10]
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    xAxis: {
      type: 'category',
      data: amountRanges,
      axisLabel: {
        interval: 0,
        rotate: 0
      }
    },
    yAxis: {
      type: 'value',
      name: '占比(%)'
    },
    series: [
      {
        name: '贷款金额区间',
        type: 'bar',
        barWidth: '60%',
        data: amountData,
        itemStyle: {
          color: function(params) {
            const colorList = [
              '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452'
            ]
            return colorList[params.dataIndex]
          }
        },
        label: {
          show: true,
          position: 'top',
          formatter: '{c}%'
        }
      }
    ]
  }
  
  amountRangeChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    amountRangeChartInstance && amountRangeChartInstance.resize()
  })
}

// 排行榜类型改变
const handleRankingTypeChange = () => {
  loadRankingData()
}

// 查看贷款详情
const viewLoanDetail = (id) => {
  router.push(`/mortgage/${id}`)
}

// 数字格式化
const formatNumber = (num) => {
  if (num === undefined || num === null) return '-'
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 日期格式化
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 初始化
onMounted(() => {
  setDefaultDateRange('month')
  fetchData()
})
</script>

<style scoped>
.loan-statistics-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 22px;
}

.filter-card {
  margin-bottom: 20px;
}

.loading-container {
  padding: 20px 0;
}

.stats-overview {
  margin-bottom: 20px;
}

.stats-card {
  text-align: center;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  margin: 10px 0;
}

.stats-trend {
  font-size: 14px;
  color: #909399;
}

.trend-up {
  color: #f56c6c;
}

.trend-down {
  color: #67c23a;
}

.chart-row {
  margin-bottom: 20px;
}

.chart-card {
  height: 400px;
}

.chart-container {
  height: 320px;
}

.risk-card {
  margin-bottom: 20px;
  text-align: center;
}

.progress-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.progress-value {
  font-size: 22px;
  font-weight: bold;
}

.progress-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.ranking-card {
  margin-bottom: 20px;
}

.property-meta {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style> 