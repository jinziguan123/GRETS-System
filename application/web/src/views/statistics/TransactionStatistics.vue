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
<!--        <el-form-item label="区域">-->
<!--          <el-select v-model="province" placeholder="选择区域" clearable @change="fetchData" style="width: 100px">-->
<!--            <el-option v-for="item in provinceOptions" :key="item.value" :label="item.label" :value="item.value" />-->
<!--          </el-select>-->
<!--        </el-form-item>-->
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
              <div>{{ scope.row.title }}</div>
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
// 导入API函数
import { queryTransactionStatistics } from '@/api/transaction'
import { queryRealtyList } from '@/api/realty'

const router = useRouter()

// 加载状态
const loading = ref(false)

// 筛选条件
const timeRange = ref('month')
const dateRange = ref([])
const province = ref('')

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

// 实际数据存储
const transactionData = ref([])
const realtyData = ref([])

// 地区选项
const provinceOptions = [
  { value: '上海市', label: '上海市' },
  { value: '北京市', label: '北京市' },
  { value: '天津市', label: '天津市' },
  { value: '重庆市', label: '重庆市' },
  { value: '江苏省', label: '江苏省' },
  { value: '浙江省', label: '浙江省' },
  { value: '安徽省', label: '安徽省' },
  { value: '福建省', label: '福建省' },
  { value: '江西省', label: '江西省' },
  { value: '山东省', label: '山东省' },
  { value: '河南省', label: '河南省' },
  { value: '湖北省', label: '湖北省' },
  { value: '湖南省', label: '湖南省' },
  { value: '广东省', label: '广东省' },
  { value: '广西壮族自治区', label: '广西壮族自治区' },
  { value: '海南省', label: '海南省' }
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
      // 这里由于startOf不兼容quarter，所以需要根据当前月份匹配对应季度
      // 计算季度开始日期
      const quarter = Math.floor(now.month() / 3) + 1
      const startMonth = (quarter - 1) * 3; // 季度开始月份（0-11）
      const startDate = dayjs(`${now.year()}-${startMonth + 1}-01`);

      // 计算季度结束日期
      const endMonth = startMonth + 2; // 季度结束月份（0-11）
      const lastDayOfEndMonth = startDate.endOf('month').date(); // 获取该月最后一天
      const endDate = dayjs(`${now.year()}-${endMonth + 1}-${lastDayOfEndMonth}`);
      dateRange.value = [startDate.toDate(), endDate.toDate()]
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
  province.value = ''
  setDefaultDateRange('custom')
  fetchData()
}

// 生成地址信息
const generateAddress = (item) => {
  // 根据省份类型决定地址格式
  if (item.province.includes('市')) {
    // 直辖市格式
    return `${item.province}${item.district}区${item.street}${item.community}${item.unit}单元${item.floor}楼${item.room}号`
  } else {
    // 普通省份格式
    return `${item.province}${item.city}${item.district}区${item.street}${item.community}${item.unit}单元${item.floor}楼${item.room}号`
  }
}

// 获取统计数据
const fetchData = async () => {
  loading.value = true
  
  try {
    // 构建日期查询参数
    const startDate = dateRange.value && dateRange.value[0] ? dayjs(dateRange.value[0]).format('YYYY-MM-DD') : undefined
    const endDate = dateRange.value && dateRange.value[1] ? dayjs(dateRange.value[1]).format('YYYY-MM-DD') : undefined
    
    // 使用单一API调用获取所有统计数据
    const result = await queryTransactionStatistics({
      startDate,
      endDate,
      province: province.value || undefined
    })

    if (!result) {
      throw new Error('获取数据失败')
    }
    
    // 从API返回结果中提取数据
    const { 
      totalTransactions, 
      totalAmount, 
      averagePrice, 
      totalTax,
      transactionDTOList = []
    } = result
    
    // 获取房产详情以补充图表所需信息
    const realtyResult = await queryRealtyList({
      pageSize: 1000,
      pageNumber: 1
    })
    
    const realtyItems = realtyResult.realtyList || []
    
    // 获取上一时间段数据进行比较（简化实现，假设上期数据）
    // 实际项目中，可以通过动态计算上一时间段，并重新请求API获取
    const prevTotalTransactions = totalTransactions * 0.9  // 模拟上期数据
    const prevTotalAmount = totalAmount * 0.93
    const prevAveragePrice = averagePrice * 0.97
    const prevTotalTax = totalTax * 0.95
    
    // 计算增长率
    const totalTransactionsTrend = totalTransactions > 0 && prevTotalTransactions > 0 
      ? ((totalTransactions - prevTotalTransactions) / prevTotalTransactions) * 100 
      : 0
    const totalAmountTrend = totalAmount > 0 && prevTotalAmount > 0 
      ? ((totalAmount - prevTotalAmount) / prevTotalAmount) * 100 
      : 0
    const averagePriceTrend = averagePrice > 0 && prevAveragePrice > 0 
      ? ((averagePrice - prevAveragePrice) / prevAveragePrice) * 100 
      : 0
    const totalTaxTrend = totalTax > 0 && prevTotalTax > 0 
      ? ((totalTax - prevTotalTax) / prevTotalTax) * 100 
      : 0
    
    // 更新数据概览
    overviewData.totalTransactions = totalTransactions || 0
    overviewData.totalTransactionsTrend = parseFloat(totalTransactionsTrend.toFixed(1))
    overviewData.totalAmount = totalAmount || 0
    overviewData.totalAmountTrend = parseFloat(totalAmountTrend.toFixed(1))
    overviewData.averagePrice = parseFloat((averagePrice || 0).toFixed(2))
    overviewData.averagePriceTrend = parseFloat(averagePriceTrend.toFixed(1))
    overviewData.totalTax = totalTax || 0
    overviewData.totalTaxTrend = parseFloat(totalTaxTrend.toFixed(1))
    
    // 存储交易数据用于图表
    transactionData.value = transactionDTOList
    realtyData.value = realtyItems
    
    // 加载排行榜数据
    loadRankingData()
    
    // 数据获取完成后，通过nextTick确保DOM已更新，然后渲染图表
    nextTick(() => {
      initCharts()
    })
    
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

// 获取户型文本
const getHouseTypeText = (houseType) => {
  const houseTypeMap = {
    'single': '一室',
    'double': '两室',
    'triple': '三室',
    'multiple': '四室及以上',
    'apartment': '公寓',
    'villa': '别墅',
    'office': '办公',
    'commercial': '商铺',
    'industrial': '工业用地',
    'other': '其他'
  }
  return houseTypeMap[houseType] || houseType
}

// 加载排行榜数据
const loadRankingData = () => {
  // 根据交易和房产数据构建排行榜数据
  const rankingItems = transactionData.value.map(transaction => {
    // 查找对应的房产信息
    const realty = realtyData.value.find(r => r.realtyCertHash === transaction.realtyCertHash) || {}

    // 计算单价，保留两位小数
    const unitPrice = realty.area && realty.area > 0 && transaction.price
      ? parseFloat((transaction.price / realty.area).toFixed(2))
      : 0

    return {
      id: transaction.transactionUUID,
      title: generateAddress(realty),
      roomType: getHouseTypeText(realty.houseType) || '未知户型',
      area: realty.area || 0,
      unitPrice: unitPrice,
      totalAmount: transaction.price || 0,
      transactionDate: transaction.createTime || new Date()
    }
  }).filter(item => item.totalAmount > 0) // 过滤无效数据

  // 根据排行类型排序
  if (rankingType.value === 'amount') {
    rankingItems.sort((a, b) => b.totalAmount - a.totalAmount)
  } else if (rankingType.value === 'area') {
    rankingItems.sort((a, b) => b.area - a.area)
  } else if (rankingType.value === 'price') {
    rankingItems.sort((a, b) => b.unitPrice - a.unitPrice)
  }
  
  // 获取前10个记录
  rankingData.value = rankingItems.slice(0, 10)
}

// 初始化图表
const initCharts = () => {
  // 确保DOM已更新
  nextTick(() => {
    initTransactionTrendChart()
    initDistrictChart()
    initPropertyTypeChart()
    initPriceHeatmapChart()
  })
}

// 初始化交易趋势图表
const initTransactionTrendChart = () => {
  if (!transactionTrendChart.value) return
  
  if (trendChartInstance) {
    trendChartInstance.dispose()
  }
  
  trendChartInstance = echarts.init(transactionTrendChart.value)
  
  // 如果没有数据，显示空图表
  if (transactionData.value.length === 0) {
    trendChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  // 根据timeRange生成不同的日期标签和数据
  const { labels, countData, amountData } = generateChartDataByTimeRange()
  
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
      data: labels
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
  
  // 确保图表能够正确渲染
  nextTick(() => {
    if (trendChartInstance) {
      trendChartInstance.resize()
    }
  })
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    trendChartInstance && trendChartInstance.resize()
  })
}

// 根据timeRange生成图表数据
const generateChartDataByTimeRange = () => {
  const today = dayjs()
  let labels = []
  const dataMap = {}
  
  // 根据不同的时间范围生成不同的标签和数据结构
  if (timeRange.value === 'month') {
    // 本月：显示每一天
    const daysInMonth = today.daysInMonth()
    const currentMonth = today.format('YYYY-MM')
    
    for (let day = 1; day <= daysInMonth; day++) {
      const dayStr = day < 10 ? `0${day}` : `${day}`
      const dateKey = `${currentMonth}-${dayStr}`
      const label = `${day}日`
      labels.push(label)
      dataMap[dateKey] = { count: 0, amount: 0 }
    }
    
    // 汇总交易数据（按天）
    transactionData.value.forEach(transaction => {
      const date = dayjs(transaction.createTime)
      // 只统计本月数据
      if (date.format('YYYY-MM') === currentMonth) {
        const dateKey = date.format('YYYY-MM-DD')
        if (dataMap[dateKey]) {
          dataMap[dateKey].count += 1
          dataMap[dateKey].amount += (transaction.price || 0)
        }
      }
    })
    
  } else if (timeRange.value === 'quarter') {
    // 本季度：显示本季度的各个月份
    const currentQuarter = Math.floor(today.month() / 3)
    const startMonth = currentQuarter * 3 // 季度开始月份（0-11）
    
    for (let i = 0; i < 3; i++) {
      const monthIndex = startMonth + i
      const month = dayjs().month(monthIndex).startOf('month')
      const monthKey = month.format('YYYY-MM')
      const label = month.format('M月')
      labels.push(label)
      dataMap[monthKey] = { count: 0, amount: 0 }
    }
    
    // 汇总交易数据（按月）
    transactionData.value.forEach(transaction => {
      const date = dayjs(transaction.createTime)
      const monthKey = date.format('YYYY-MM')
      if (dataMap[monthKey]) {
        dataMap[monthKey].count += 1
        dataMap[monthKey].amount += (transaction.price || 0)
      }
    })
    
  } else if (timeRange.value === 'year') {
    // 本年度：显示每个月
    const currentYear = today.format('YYYY')
    
    for (let month = 0; month < 12; month++) {
      const monthDate = dayjs().year(parseInt(currentYear)).month(month).startOf('month')
      const monthKey = monthDate.format('YYYY-MM')
      const label = monthDate.format('M月')
      labels.push(label)
      dataMap[monthKey] = { count: 0, amount: 0 }
    }
    
    // 汇总交易数据（按月）
    transactionData.value.forEach(transaction => {
      const date = dayjs(transaction.createTime)
      // 只统计本年数据
      if (date.format('YYYY') === currentYear) {
        const monthKey = date.format('YYYY-MM')
        if (dataMap[monthKey]) {
          dataMap[monthKey].count += 1
          dataMap[monthKey].amount += (transaction.price || 0)
        }
      }
    })
    
  } else {
    // 自定义时间范围或其他情况：显示所有月份
    // 获取日期范围的起止时间
    const startDate = dateRange.value && dateRange.value[0] ? dayjs(dateRange.value[0]) : dayjs().subtract(30, 'days')
    const endDate = dateRange.value && dateRange.value[1] ? dayjs(dateRange.value[1]) : dayjs()
    
    // 计算日期范围内的月份
    let currentDate = startDate.startOf('month')
    while (currentDate.isBefore(endDate) || currentDate.isSame(endDate, 'month')) {
      const monthKey = currentDate.format('YYYY-MM')
      const label = currentDate.format('YYYY年M月')
      labels.push(label)
      dataMap[monthKey] = { count: 0, amount: 0 }
      currentDate = currentDate.add(1, 'month')
    }
    
    // 汇总交易数据
    transactionData.value.forEach(transaction => {
      const date = dayjs(transaction.createTime)
      const monthKey = date.format('YYYY-MM')
      if (dataMap[monthKey]) {
        dataMap[monthKey].count += 1
        dataMap[monthKey].amount += (transaction.price || 0)
      }
    })
  }
  
  // 提取数据数组
  const countData = []
  const amountData = []
  
  // 根据标签顺序提取数据
  const sortedKeys = Object.keys(dataMap).sort()
  for (let i = 0; i < labels.length; i++) {
    const key = sortedKeys[i]
    countData.push(dataMap[key]?.count || 0)
    amountData.push(dataMap[key]?.amount || 0)
  }
  
  return { labels, countData, amountData }
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
  
  // 如果没有数据，显示空图表
  if (transactionData.value.length === 0 || realtyData.value.length === 0) {
    districtChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  // 按区域分组统计
  const provinceStats = {}
  
  // 初始化所有省份数据
  provinceOptions.forEach(opt => {
    provinceStats[opt.value] = { count: 0, amount: 0 }
  })
  
  // 统计各省份交易数据
  transactionData.value.forEach(transaction => {
    const realty = realtyData.value.find(r => r.realtyCertHash === transaction.realtyCertHash)
    
    if (realty && realty.province && provinceStats[realty.province]) {
      provinceStats[realty.province].count += 1
      provinceStats[realty.province].amount += (transaction.price || 0)
    }
  })  
  
  // 提取有交易的省份数据
  const pieData = []
  Object.keys(provinceStats).forEach(province => {
    if (provinceStats[province].count > 0) {
      pieData.push({
        name: province,
        value: districtChartType.value === 'count' 
          ? provinceStats[province].count 
          : provinceStats[province].amount
      })
    }
  })
  
  // 如果没有有效数据，显示空图表
  if (pieData.length === 0) {
    districtChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
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
        data: pieData
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
  
  // 如果没有数据，显示空图表
  if (transactionData.value.length === 0 || realtyData.value.length === 0) {
    typeChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  // 统计不同类型房产的交易数量
  const typeStats = {}
  
  // 房产类型中英文映射
  const typeMapping = {
    'single': '一室',
    'double': '两室',
    'triple': '三室',
    'multiple': '四室及以上',
    'apartment': '公寓',
    'villa': '别墅',
    'office': '办公',
    'commercial': '商铺',
    'industrial': '工业用地',
    'other': '其他'
  }
  
  // 获取已完成交易的房产证号
  const transactedRealtyCerts = transactionData.value.map(tx => tx.realtyCertHash)
  
  // 统计已交易房产的类型分布
  const transactedRealty = realtyData.value.filter(realty => 
    transactedRealtyCerts.includes(realty.realtyCertHash) && realty.houseType)
  
  transactedRealty.forEach(realty => {
    const type = realty.houseType || 'other'
    const displayType = typeMapping[type] || type
    
    if (!typeStats[displayType]) {
      typeStats[displayType] = 0
    }
    typeStats[displayType] += 1
  })
  
  // 提取数据
  const typeData = []
  Object.keys(typeStats).forEach(type => {
    if (typeStats[type] > 0) {
      typeData.push({
        name: type,
        value: typeStats[type]
      })
    }
  })
  
  // 如果没有有效数据，显示空图表
  if (typeData.length === 0) {
    typeChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      type: 'scroll',
      data: typeData.map(item => item.name)
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
  
  // 如果没有数据，显示空图表
  if (transactionData.value.length === 0 || realtyData.value.length === 0) {
    heatmapChartInstance.setOption({
      title: {
        text: '暂无数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  // 房产类型中英文映射
  const typeMapping = {
    'single': '一室',
    'double': '两室',
    'triple': '三室',
    'multiple': '四室及以上',
    'apartment': '公寓',
    'villa': '别墅',
    'office': '办公',
    'commercial': '商铺',
    'industrial': '工业用地',
    'other': '其他'
  }
  
  // 获取有效房产数据（有省份、户型和价格信息）
  const validRealty = realtyData.value.filter(realty => 
    realty.province && 
    realty.houseType && 
    realty.price && 
    realty.area && 
    realty.area > 0
  )
  
  // 提取所有出现的省份和户型
  const provinces = Array.from(new Set(validRealty.map(r => r.province)))
  const roomTypes = Array.from(new Set(validRealty.map(r => r.houseType)))
  
  // 转换户型为中文显示
  const roomTypeLabels = roomTypes.map(type => typeMapping[type] || type)
  
  // 如果数据不足，显示空图表
  if (provinces.length === 0 || roomTypes.length === 0) {
    heatmapChartInstance.setOption({
      title: {
        text: '暂无足够数据',
        left: 'center',
        top: 'center'
      }
    })
    return
  }
  
  // 计算每个省份、户型组合的平均单价
  const priceMap = {}
  
  validRealty.forEach(realty => {
    const key = `${realty.province}:${realty.houseType}`
    if (!priceMap[key]) {
      priceMap[key] = {
        totalPrice: 0,
        count: 0
      }
    }
    
    priceMap[key].totalPrice += (realty.price / realty.area)
    priceMap[key].count += 1
  })
  
  // 准备热力图数据
  const data = []
  for (let i = 0; i < provinces.length; i++) {
    for (let j = 0; j < roomTypes.length; j++) {
      const key = `${provinces[i]}:${roomTypes[j]}`
      if (priceMap[key] && priceMap[key].count > 0) {
        const avgPrice = priceMap[key].totalPrice / priceMap[key].count
        data.push([i, j, Math.round(avgPrice)])
      } else {
        // 对于没有数据的省份-户型组合，使用null表示
        data.push([i, j, null])
      }
    }
  }
  
  // 计算价格范围
  const prices = data.filter(item => item[2] !== null).map(item => item[2])
  const minPrice = Math.min(...prices) * 0.9
  const maxPrice = Math.max(...prices) * 1.1
  
  const option = {
    tooltip: {
      position: 'top',
      formatter: function(params) {
        const roomTypeLabel = roomTypeLabels[params.data[1]] || roomTypes[params.data[1]]
        if (params.data[2] === null) {
          return `${provinces[params.data[0]]}, ${roomTypeLabel}<br/>暂无数据`
        }
        return `${provinces[params.data[0]]}, ${roomTypeLabel}<br/>
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
      data: provinces,
      axisLabel: {
        interval: 0,
        rotate: 45
      }
    },
    yAxis: {
      type: 'category',
      data: roomTypeLabels
    },
    visualMap: {
      min: minPrice,
      max: maxPrice,
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
onMounted(async () => {
  // 设置默认日期范围
  setDefaultDateRange('month')
  // 获取数据
  await fetchData()
  
  // 确保图表在初始加载时能正确渲染
  nextTick(() => {
    if (transactionTrendChart.value && !trendChartInstance) {
      initTransactionTrendChart()
    }
    if (districtChart.value && !districtChartInstance) {
      initDistrictChart()
    }
    if (propertyTypeChart.value && !typeChartInstance) {
      initPropertyTypeChart()
    }
    if (priceHeatmapChart.value && !heatmapChartInstance) {
      initPriceHeatmapChart()
    }
  })
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