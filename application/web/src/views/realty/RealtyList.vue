<template>
  <div class="realty-list">
    <div class="page-header">
      <h2>房产列表</h2>
      <el-button 
        type="primary" 
        v-if="userStore.hasOrganization('government')"
        @click="router.push('/realty/create')"
      >
        添加房产
      </el-button>
    </div>

    <!-- 搜索条件 -->
    <el-card class="filter-container">
      <el-form :model="searchForm" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="地区">
              <el-select v-model="searchForm.district" placeholder="选择地区" clearable>
                <el-option v-for="item in districtOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="户型">
              <el-select v-model="searchForm.roomType" placeholder="选择户型" clearable>
                <el-option v-for="item in roomTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col :span="12">
            <el-form-item label="面积">
              <el-input-number v-model="searchForm.minArea" :min="0" :step="10" placeholder="最小" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxArea" :min="0" :step="10" placeholder="最大" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格">
              <el-input-number v-model="searchForm.minPrice" :min="0" :step="50000" :formatter="formatPrice" placeholder="最小" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxPrice" :min="0" :step="50000" :formatter="formatPrice" placeholder="最大" />
            </el-form-item>
          </el-col>
        </el-row>
        <div style="float: right">
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-form-item>
        </div>
      </el-form>
    </el-card>

    <!-- 房产列表 -->
    <el-card class="realty-cards-container">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else-if="realtyList.length === 0" class="empty-data">
        <el-empty description="暂无房产数据" />
      </div>
      <div v-else class="realty-grid">
        <el-card 
          v-for="item in realtyList" 
          :key="item.id" 
          class="realty-card"
          @click="viewDetails(item)"
        >
          <div class="realty-image">
            <img :src="item.imageUrl || 'https://via.placeholder.com/300x200?text=暂无图片'" alt="房产图片" />
            <div class="realty-status" :class="getStatusClass(item.status)">{{ getStatusText(item.status) }}</div>
          </div>
          <div class="realty-info">
            <h3 class="realty-title">{{ item.title }}</h3>
            <div class="realty-address">{{ item.district }} {{ item.address }}</div>
            <div class="realty-meta">
              <span class="room-type">{{ item.roomType }}</span>
              <span class="area">{{ item.area }}平米</span>
              <span class="floor">{{ item.floor }}层</span>
            </div>
            <div class="realty-price" v-if="userStore.hasOrganization(['investor', 'bank'])">
              <span>¥ {{ formatPriceText(item.expectedPrice) }}</span>
            </div>
            <div class="realty-actions">
              <el-button type="primary" size="small" @click.stop="viewDetails(item)">详情</el-button>
              <el-button 
                v-if="userStore.hasOrganization('investor') && item.status === 'available'" 
                type="success" 
                size="small" 
                @click.stop="startTransaction(item)"
              >
                交易
              </el-button>
              <el-button 
                v-if="userStore.hasOrganization('government')" 
                type="warning" 
                size="small" 
                @click.stop="editRealty(item)"
              >
                编辑
              </el-button>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

// 查询条件
const searchForm = reactive({
  district: '',
  roomType: '',
  minArea: null,
  maxArea: null,
  minPrice: null,
  maxPrice: null
})

// 重置查询条件
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = key.startsWith('min') || key.startsWith('max') ? null : ''
  })
  handleSearch()
}

// 格式化价格显示
const formatPrice = (value) => {
  if (value === null) return ''
  return `¥ ${value}`
}

const formatPriceText = (price) => {
  if (!price) return '暂无报价'
  return (price / 10000).toFixed(2) + '万'
}

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

// 户型选项
const roomTypeOptions = [
  { value: '一室一厅', label: '一室一厅' },
  { value: '两室一厅', label: '两室一厅' },
  { value: '两室两厅', label: '两室两厅' },
  { value: '三室一厅', label: '三室一厅' },
  { value: '三室两厅', label: '三室两厅' },
  { value: '四室两厅', label: '四室两厅' },
  { value: '复式', label: '复式' },
  { value: '别墅', label: '别墅' }
]

// 分页参数
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 房产列表
const realtyList = ref([])
const loading = ref(false)

// 获取房产列表
const fetchRealtyList = async () => {
  loading.value = true
  try {
    // 模拟API请求
    setTimeout(() => {
      // 这里替换为实际API调用
      const mockData = []
      for (let i = 1; i <= 20; i++) {
        mockData.push({
          id: i,
          title: `优质${i % 8 == 0 ? '别墅' : (i % 4 == 0 ? '复式' : '普通住宅')}`,
          district: districtOptions[i % districtOptions.length].value,
          address: `某某路${i}号`,
          roomType: roomTypeOptions[i % roomTypeOptions.length].value,
          area: Math.floor(70 + Math.random() * 100),
          floor: Math.floor(1 + Math.random() * 30),
          expectedPrice: Math.floor(300 + Math.random() * 700) * 10000,
          status: i % 5 === 0 ? 'sold' : (i % 3 === 0 ? 'pending' : 'available'),
          imageUrl: `https://picsum.photos/id/${i + 10}/300/200`
        })
      }
      
      // 根据查询条件过滤
      let filteredData = mockData.filter(item => {
        let match = true
        if (searchForm.district && item.district !== searchForm.district) match = false
        if (searchForm.roomType && item.roomType !== searchForm.roomType) match = false
        if (searchForm.minArea && item.area < searchForm.minArea) match = false
        if (searchForm.maxArea && item.area > searchForm.maxArea) match = false
        if (searchForm.minPrice && item.expectedPrice < searchForm.minPrice) match = false
        if (searchForm.maxPrice && item.expectedPrice > searchForm.maxPrice) match = false
        return match
      })
      
      total.value = filteredData.length
      
      // 分页处理
      const startIndex = (currentPage.value - 1) * pageSize.value
      const endIndex = startIndex + pageSize.value
      realtyList.value = filteredData.slice(startIndex, endIndex)
      
      loading.value = false
    }, 800)
  } catch (error) {
    console.error('Failed to fetch realty list:', error)
    loading.value = false
    ElMessage.error('获取房产列表失败')
  }
}

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchRealtyList()
}

// 改变每页显示数量
const handleSizeChange = (size) => {
  pageSize.value = size
  fetchRealtyList()
}

// 改变页码
const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchRealtyList()
}

// 查看详情
const viewDetails = (item) => {
  router.push(`/realty/${item.id}`)
}

// 开始交易
const startTransaction = (item) => {
  router.push({
    path: '/transaction/create',
    query: { realtyId: item.id }
  })
}

// 编辑房产
const editRealty = (item) => {
  router.push(`/realty/${item.id}/edit`)
}

// 获取状态样式类
const getStatusClass = (status) => {
  const statusMap = {
    available: 'status-available',
    pending: 'status-pending',
    sold: 'status-sold'
  }
  return statusMap[status] || ''
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    available: '可交易',
    pending: '交易中',
    sold: '已售出'
  }
  return statusMap[status] || '未知状态'
}

// 初始加载
onMounted(() => {
  fetchRealtyList()
})
</script>

<style scoped>
.realty-list {
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

.filter-container {
  margin-bottom: 20px;
}

.separator {
  margin: 0 5px;
}

.realty-cards-container {
  margin-top: 20px;
}

.loading-container {
  padding: 20px 0;
}

.empty-data {
  padding: 40px 0;
  text-align: center;
}

.realty-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.realty-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
}

.realty-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.realty-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.realty-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.realty-status {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: white;
}

.status-available {
  background-color: #67C23A;
}

.status-pending {
  background-color: #E6A23C;
}

.status-sold {
  background-color: #909399;
}

.realty-info {
  padding: 15px;
}

.realty-title {
  margin: 0 0 10px;
  font-size: 18px;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.realty-address {
  color: #606266;
  font-size: 14px;
  margin-bottom: 10px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.realty-meta {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  color: #909399;
  font-size: 13px;
}

.realty-price {
  font-size: 18px;
  color: #F56C6C;
  font-weight: bold;
  margin: 10px 0;
}

.realty-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 15px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
