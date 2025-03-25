<template>
  <div class="statistics-container">
    <div class="page-header">
      <h2>交易统计分析</h2>
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
        <el-form-item label="区域">
          <el-select v-model="district" placeholder="选择区域" clearable @change="fetchData">
            <el-option v-for="item in districtOptions" :key="item.value" :label="item.label" :value="item.value" />
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
                <span>总交易量</span>
                <el-tooltip content="统计周期内的房产交易总数量" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">{{ formatNumber(overviewData.totalTransactions) }}</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.totalTransactionsTrend > 0, 'trend-down': overviewData.totalTransactionsTrend < 0}">
                <el-icon><component :is="overviewData.totalTransactionsTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.totalTransactionsTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>总交易额</span>
                <el-tooltip content="统计周期内的房产交易总金额" placement="top">
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
                <span>平均单价</span>
                <el-tooltip content="统计周期内的房产交易平均单价" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">¥{{ formatNumber(overviewData.averagePrice) }}/㎡</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.averagePriceTrend > 0, 'trend-down': overviewData.averagePriceTrend < 0}">
                <el-icon><component :is="overviewData.averagePriceTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.averagePriceTrend) }}%
              </span>
              较上期
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stats-card">
            <template #header>
              <div class="card-header">
                <span>税收总额</span>
                <el-tooltip content="统计周期内的房产交易产生的税收总额" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <div class="stats-value">¥{{ formatNumber(overviewData.totalTax / 10000) }}万</div>
            <div class="stats-trend">
              <span :class="{'trend-up': overviewData.totalTaxTrend > 0, 'trend-down': overviewData.totalTaxTrend < 0}">
                <el-icon><component :is="overviewData.totalTaxTrend > 0 ? 'CaretTop' : 'CaretBottom'" /></el-icon>
                {{ Math.abs(overviewData.totalTaxTrend) }}%
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
                <span>交易量趋势</span>
                <el-radio-group v-model="transactionTrendType" size="small" @change="updateTransactionTrendChart">
                  <el-radio-button label="count">数量</el-radio-button>
                  <el-radio-button label="amount">金额</el-radio-button>
                </el-radio-group>
              </div>
            </template>
            <div class="chart-container" ref="transactionTrendChart"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>区域分布</span>
                <el-radio-group v-model="districtChartType" size="small" @change="updateDistrictChart">
                  <el-radio-button label="count">数量</el-radio-button>
                  <el-radio-button label="amount">金额</el-radio-button>
                </el-radio-group>
              </div>
            </template>
            <div class="chart-container" ref="districtChart"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" class="chart-row">
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>房产类型分布</span>
              </div>
            </template>
            <div class="chart-container" ref="propertyTypeChart"></div>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="chart-card">
            <template #header>
              <div class="card-header">
                <span>单价热力图</span>
              </div>
            </template>
            <div class="chart-container" ref="priceHeatmapChart"></div>
          </el-card>
        </el-col>
      </el-row>
      
      <!-- 交易排行 -->
      <el-card class="ranking-card">
        <template #header>
          <div class="card-header">
            <span>交易排行榜</span>
            <el-tabs v-model="rankingType" @tab-change="handleRankingTypeChange">
              <el-tab-pane label="按交易额" name="amount"></el-tab-pane>
              <el-tab-pane label="按面积" name="area"></el-tab-pane>
              <el-tab-pane label="按单价" name="price"></el-tab-pane>
            </el-tabs>
          </div>
        </template>
        <el-table :data="rankingData" style="width: 100%">
          <el-table-column type="index" width="50" />
          <el-table-column prop="title" label="房产信息" min-width="200">
            <template #default="scope">
              <div>{{ scope.row.district }} | {{ scope.row.title }}</div>
              <div class="property-meta">{{ scope.row.roomType }} | {{ scope.row.area }}㎡</div>
            </template>
          </el-table-column>
          <el-table-column prop="totalAmount" label="交易金额" sortable width="150">
            <template #default="scope">
              ¥{{ formatNumber(scope.row.totalAmount / 10000) }}万
            </template>
          </el-table-column>
          <el-table-column prop="unitPrice" label="单价" sortable width="150">
            <template #default="scope">
              ¥{{ formatNumber(scope.row.unitPrice) }}/㎡
            </template>
          </el-table-column>
          <el-table-column prop="transactionDate" label="交易日期" width="120">
            <template #default="scope">
              {{ formatDate(scope.row.transactionDate) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="scope">
              <el-button type="text" @click="viewTransactionDetail(scope.row.id)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
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
const district = ref('')

// 图表引用
const transactionTrendChart = ref(null)
const districtChart = ref(null)
const propertyTypeChart = ref(null)
const priceHeatmapChart = ref(null)

// 图表实例
let trendChartInstance = null
let districtChartInstance = null
let typeChartInstance = null
let heatmapChartInstance = null

// 图表类型
const transactionTrendType = ref('count')
const districtChartType = ref('count')

// 排行榜类型
const rankingType = ref('amount')
const rankingData = ref([])

// 地区选项
const districtOptions = [
  { value: '浦东新区', label: '浦东新区' },
  { value: '黄浦区', label: '黄浦区' },
  { value: '徐汇区', label: '徐汇区' },
  { value: '长宁区', label: '长宁区' },
  { value: '静安区', label: '静安区' },
  { value: '普陀区', label: '普陀区' },
  { value: '虹口区', label: '虹口区' },
  { value: '杨浦区', label: '杨浦区' },
  { value: '闵行区', label: '闵行区' },
  { value: '宝山区', label: '宝山区' },
  { value: '嘉定区', label: '嘉定区' },
  { value: '金山区', label: '金山区' },
  { value: '松江区', label: '松江区' },
  { value: '青浦区', label: '青浦区' },
  { value: '奉贤区', label: '奉贤区' },
  { value: '崇明区', label: '崇明区' }
]

// 数据概览
const overviewData = reactive({
  totalTransactions: 0,
  totalTransactionsTrend: 0,
  totalAmount: 0,
  totalAmountTrend: 0,
  averagePrice: 0,
  averagePriceTrend: 0,
  totalTax: 0,
  totalTaxTrend: 0
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
  district.value = ''
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
      overviewData.totalTransactions = 856
      overviewData.totalTransactionsTrend = 12.5
      overviewData.totalAmount = 4562000000
      overviewData.totalAmountTrend = 8.3
      overviewData.averagePrice = 52680
      overviewData.averagePriceTrend = 3.2
      overviewData.totalTax = 152480000
      overviewData.totalTaxTrend = 5.7
      
      // 加载排行榜数据
      loadRankingData()
      
      // 绘制图表
      nextTick(() => {
        initCharts()
      })
      
      loading.value = false
    }, 1000)
  } catch (error) {
    console.error('Failed to fetch statistics data:', error)
    ElMessage.error('获取统计数据失败')
    loading.value = false
  }
}

// 加载排行榜数据
const loadRankingData = () => {
  // 这里替换为实际的API调用
  const mockData = []
  for (let i = 1; i <= 10; i++) {
    const district = districtOptions[Math.floor(Math.random() * districtOptions.length)].value
    const area = Math.floor(80 + Math.random() * 150)
    const unitPrice = Math.floor(45000 + Math.random() * 20000)
    const totalAmount = area * unitPrice
    
    mockData.push({
      id: `T${String(i).padStart(3, '0')}`,
      title: `${district}某小区${i}号楼${Math.floor(Math.random() * 30) + 1}层`,
      district: district,
      roomType: `${Math.floor(Math.random() * 3) + 2}室${Math.floor(Math.random() * 2) + 1}厅`,
      area: area,
      unitPrice: unitPrice,
      totalAmount: totalAmount,
      transactionDate: new Date(Date.now() - Math.floor(Math.random() * 90) * 24 * 60 * 60 * 1000)
    })
  }
  
  // 根据排行类型排序
  if (rankingType.value === 'amount') {
    mockData.sort((a, b) => b.totalAmount - a.totalAmount)
  } else if (rankingType.value === 'area') {
    mockData.sort((a, b) => b.area - a.area)
  } else if (rankingType.value === 'price') {
    mockData.sort((a, b) => b.unitPrice - a.unitPrice)
  }
  
  rankingData.value = mockData
}

// 初始化图表
const initCharts = () => {
  initTransactionTrendChart()
  initDistrictChart()
  initPropertyTypeChart()
  initPriceHeatmapChart()
}

// 初始化交易趋势图表
const initTransactionTrendChart = () => {
  if (!transactionTrendChart.value) return
  
  if (trendChartInstance) {
    trendChartInstance.dispose()
  }
  
  trendChartInstance = echarts.init(transactionTrendChart.value)
  
  // 模拟数据
  const months = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
  const countData = [120, 132, 101, 134, 90, 230, 210, 182, 191, 234, 290, 330]
  const amountData = [2200, 2400, 1900, 2600, 1800, 4600, 4200, 3800, 3900, 4800, 5200, 6000].map(v => v * 10000)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        const param = params[0]
        if (transactionTrendType.value === 'count') {
          return `${param.axisValue}<br/>${param.marker}交易量: ${param.value}`
        } else {
          return `${param.axisValue}<br/>${param.marker}交易额: ¥${formatNumber(param.value / 10000)}万`
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
          if (transactionTrendType.value === 'count') {
            return value
          } else {
            return `${value / 10000}万`
          }
        }
      }
    },
    series: [
      {
        name: transactionTrendType.value === 'count' ? '交易量' : '交易额',
        type: 'line',
        smooth: true,
        data: transactionTrendType.value === 'count' ? countData : amountData,
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

// 更新交易趋势图表
const updateTransactionTrendChart = () => {
  initTransactionTrendChart()
}

// 初始化区域分布图表
const initDistrictChart = () => {
  if (!districtChart.value) return
  
  if (districtChartInstance) {
    districtChartInstance.dispose()
  }
  
  districtChartInstance = echarts.init(districtChart.value)
  
  // 模拟数据
  const districtNames = districtOptions.map(item => item.value)
  const countData = districtNames.map(() => Math.floor(Math.random() * 100) + 20)
  const amountData = districtNames.map(() => (Math.floor(Math.random() * 1000) + 200) * 1000000)
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: function(params) {
        if (districtChartType.value === 'count') {
          return `${params.name}<br/>${params.marker}交易量: ${params.value} (${params.percent}%)`
        } else {
          return `${params.name}<br/>${params.marker}交易额: ¥${formatNumber(params.value / 10000)}万 (${params.percent}%)`
        }
      }
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      type: 'scroll'
    },
    series: [
      {
        name: districtChartType.value === 'count' ? '交易量' : '交易额',
        type: 'pie',
        radius: ['40%', '70%'],
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
        data: districtNames.map((name, index) => {
          return {
            name: name,
            value: districtChartType.value === 'count' ? countData[index] : amountData[index]
          }
        })
      }
    ]
  }
  
  districtChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    districtChartInstance && districtChartInstance.resize()
  })
}

// 更新区域分布图表
const updateDistrictChart = () => {
  initDistrictChart()
}

// 初始化房产类型分布图表
const initPropertyTypeChart = () => {
  if (!propertyTypeChart.value) return
  
  if (typeChartInstance) {
    typeChartInstance.dispose()
  }
  
  typeChartInstance = echarts.init(propertyTypeChart.value)
  
  // 模拟数据
  const typeNames = ['一室一厅', '两室一厅', '两室两厅', '三室一厅', '三室两厅', '四室两厅', '复式', '别墅']
  const typeData = [
    { value: 25, name: '一室一厅' },
    { value: 85, name: '两室一厅' },
    { value: 120, name: '两室两厅' },
    { value: 150, name: '三室一厅' },
    { value: 180, name: '三室两厅' },
    { value: 90, name: '四室两厅' },
    { value: 40, name: '复式' },
    { value: 20, name: '别墅' }
  ]
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      type: 'scroll',
      data: typeNames
    },
    series: [
      {
        name: '房产类型',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
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
        data: typeData
      }
    ]
  }
  
  typeChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    typeChartInstance && typeChartInstance.resize()
  })
}

// 初始化价格热力图
const initPriceHeatmapChart = () => {
  if (!priceHeatmapChart.value) return
  
  if (heatmapChartInstance) {
    heatmapChartInstance.dispose()
  }
  
  heatmapChartInstance = echarts.init(priceHeatmapChart.value)
  
  // 模拟数据
  const districts = districtOptions.map(item => item.value)
  const roomTypes = ['一室一厅', '两室一厅', '两室两厅', '三室一厅', '三室两厅', '四室两厅', '复式', '别墅']
  
  const data = []
  for (let i = 0; i < districts.length; i++) {
    for (let j = 0; j < roomTypes.length; j++) {
      // 生成不同户型在不同区域的平均单价
      let basePrice = 45000
      
      // 不同区域价格差异
      if (['静安区', '黄浦区', '徐汇区'].includes(districts[i])) {
        basePrice += 20000
      } else if (['长宁区', '浦东新区', '虹口区'].includes(districts[i])) {
        basePrice += 10000
      } else if (['松江区', '青浦区', '嘉定区'].includes(districts[i])) {
        basePrice -= 15000
      }
      
      // 不同户型价格差异
      if (roomTypes[j] === '别墅') {
        basePrice += 30000
      } else if (roomTypes[j] === '复式') {
        basePrice += 15000
      } else if (roomTypes[j].startsWith('四室')) {
        basePrice += 5000
      } else if (roomTypes[j].startsWith('一室')) {
        basePrice -= 5000
      }
      
      // 添加随机波动
      basePrice += Math.random() * 10000 - 5000
      
      data.push([i, j, Math.round(basePrice)])
    }
  }
  
  const option = {
    tooltip: {
      position: 'top',
      formatter: function(params) {
        return `${districts[params.data[0]]}, ${roomTypes[params.data[1]]}<br/>
                平均单价: ¥${formatNumber(params.data[2])}/㎡`
      }
    },
    grid: {
      left: '3%',
      right: '7%',
      bottom: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: districts,
      axisLabel: {
        interval: 0,
        rotate: 45
      }
    },
    yAxis: {
      type: 'category',
      data: roomTypes
    },
    visualMap: {
      min: 30000,
      max: 90000,
      calculable: true,
      orient: 'horizontal',
      right: 'center',
      bottom: '0%',
      color: ['#ff0000', '#ffff00', '#0000ff']
    },
    series: [
      {
        name: '单价热力图',
        type: 'heatmap',
        data: data,
        label: {
          show: false
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  }
  
  heatmapChartInstance.setOption(option)
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    heatmapChartInstance && heatmapChartInstance.resize()
  })
}

// 排行榜类型改变
const handleRankingTypeChange = () => {
  loadRankingData()
}

// 查看交易详情
const viewTransactionDetail = (id) => {
  router.push(`/transaction/${id}`)
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
.statistics-container {
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

.ranking-card {
  margin-bottom: 20px;
}

.property-meta {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style> 