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
      <el-form :model="searchForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="证书号">
              <el-input v-model="searchForm.realtyCert" placeholder="请输入不动产证号" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房产类型">
              <el-select v-model="searchForm.realtyType" placeholder="请选择房产类型" clearable style="width: 100%">
                <el-option label="住宅" value="HOUSE"></el-option>
                <el-option label="商铺" value="SHOP"></el-option>
                <el-option label="办公" value="OFFICE"></el-option>
                <el-option label="工业" value="INDUSTRIAL"></el-option>
                <el-option label="其他" value="OTHER"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="户型">
              <el-select v-model="searchForm.houseType" placeholder="请选择户型" clearable style="width: 100%">
                <el-option label="一室" value="single"></el-option>
                <el-option label="两室" value="double"></el-option>
                <el-option label="三室" value="triple"></el-option>
                <el-option label="四室及以上" value="multiple"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="省份">
              <el-input v-model="searchForm.province" placeholder="请输入省份" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="城市">
              <el-input v-model="searchForm.city" placeholder="请输入城市" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="区/县">
              <el-input v-model="searchForm.district" placeholder="请输入区/县" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="街道">
              <el-input v-model="searchForm.street" placeholder="请输入街道" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="小区">
              <el-input v-model="searchForm.community" placeholder="请输入小区" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="单元">
              <el-input v-model="searchForm.unit" placeholder="请输入单元" clearable />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="楼层">
              <el-input v-model="searchForm.floor" placeholder="请输入楼层" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房号">
              <el-input v-model="searchForm.room" placeholder="请输入房号" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="是否为新房">
              <el-switch v-model="searchForm.isNewHouse" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="面积范围">
              <el-input-number v-model="searchForm.minArea" :min="0.0" :step="10" placeholder="最小面积" style="width: 45%" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxArea" :min="0.0" :step="10" placeholder="最大面积" style="width: 45%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格范围">
              <el-input-number v-model="searchForm.minPrice" :min="0.0" :step="50000" :formatter="formatPrice" :parser="parsePrice" placeholder="最小价格" style="width: 45%" />
              <span class="separator">-</span>
              <el-input-number v-model="searchForm.maxPrice" :min="0.0" :step="50000" :formatter="formatPrice" :parser="parsePrice" placeholder="最大价格" style="width: 45%" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <div style="text-align: right; margin-top: 20px;">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
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
          :key="item.realtyCert" 
          class="realty-card"
          @click="viewDetails(item)"
        >
          <div class="realty-image">
            <img :src="getRandomImage(item.realtyCert)" alt="房产图片" />
            <div class="realty-status" :class="getStatusClass(item.status)">{{ getStatusText(item.status) }}</div>
            <div class="realty-type-tag">{{ item.isNewHouse ? '新房' : '二手房' }}</div>
          </div>
          <div class="realty-info">
            <h3 class="realty-title">{{ generateTitle(item) }}</h3>
            <div class="realty-address">{{ generateAddress(item) }}</div>
            <div class="realty-meta">
              <span class="house-type">{{ getHouseTypeText(item.houseType) }}</span>
              <span class="area">{{ item.area }}平米</span>
            </div>
            <div class="realty-price" v-if="userStore.hasOrganization(['investor', 'bank'])">
              <span>¥ {{ formatPriceText(item.price) }}</span>
            </div>
            <div class="realty-actions">
              <el-button type="primary" size="small" @click.stop="viewDetails(item)">详情</el-button>
              <el-button 
                v-if="userStore.hasOrganization('investor') && item.status === 'NORMAL'" 
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
          v-model:current-page="searchForm.pageNumber"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 30, 50]"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import {queryRealtyList} from "@/api/realty.ts";

const router = useRouter()
const userStore = useUserStore()

// 查询条件
const searchForm = reactive({
  realtyCert: '',
  realtyType: '',
  houseType: '',
  minPrice: null,
  maxPrice: null,
  minArea: null,
  maxArea: null,
  province: '',
  city: '',
  district: '',
  street: '',
  community: '',
  unit: '',
  floor: '',
  room: '',
  isNewHouse: false,
  pageSize: 10,
  pageNumber: 1
})

// 重置查询条件
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    if (key === 'pageSize') {
      searchForm[key] = 10
    } else if (key === 'pageNumber') {
      searchForm[key] = 1
    } else {
      searchForm[key] = key.startsWith('min') || key.startsWith('max') ? null : ''
    }
  })
  handleSearch()
}

// 格式化价格显示
const formatPrice = (value) => {
  if (value === null) return ''
  return `¥ ${value}`
}

// 解析价格
const parsePrice = (value) => {
  if (value === '') return null
  return value.replace(/[^\d]/g, '')
}

// 价格文本格式化
const formatPriceText = (price) => {
  if (!price) return '暂无报价'
  return (price / 10000).toFixed(2) + '万'
}

// 获取随机图片
const getRandomImage = (id) => {
  // 使用房产证号作为种子生成一个稳定的随机数，这样同一个房产始终显示相同的图片
  const seed = id.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  const imageId = (seed % 100) + 1 // 限制在1-100之间
  return `https://picsum.photos/id/${imageId}/300/200`
}

// 获取状态样式类
const getStatusClass = (status) => {
  const statusMap = {
    'NORMAL': 'status-normal',
    'IN_TRANSACTION': 'status-pending',
    'MORTGAGED': 'status-mortgaged',
    'FROZEN': 'status-frozen'
  }
  return statusMap[status] || ''
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'NORMAL': '可交易',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || '未知状态'
}

// 获取户型文本
const getHouseTypeText = (houseType) => {
  const houseTypeMap = {
    'single': '一室',
    'double': '两室',
    'triple': '三室',
    'multiple': '四室及以上'
  }
  return houseTypeMap[houseType] || houseType
}

// 生成房产标题
const generateTitle = (item) => {
  if (!item) return '未知房产'
  
  let title = ''
  
  // 房产类型
  const typeMap = {
    'HOUSE': '住宅',
    'SHOP': '商铺',
    'OFFICE': '办公',
    'INDUSTRIAL': '工业',
    'OTHER': '其他'
  }
  
  title += typeMap[item.realtyType] || '未知类型'
  
  // 户型信息
  if (item.realtyType === 'HOUSE') {
    title += ' - ' + getHouseTypeText(item.houseType)
  }
  
  // 小区/位置信息
  if (item.community) {
    title += ' - ' + item.community
  }
  
  return title
}

// 生成房产地址
const generateAddress = (item) => {
  if (!item) return '地址不详'
  
  const parts = []
  if (item.province) parts.push(item.province)
  if (item.city && item.city !== item.province) parts.push(item.city)
  if (item.district && item.district !== item.city) parts.push(item.district)
  if (item.street) parts.push(item.street)
  if (item.community) parts.push(item.community)
  if (item.unit) parts.push(item.unit + '单元')
  if (item.floor) parts.push(item.floor + '层')
  if (item.room) parts.push(item.room + '室')
  
  return parts.length > 0 ? parts.join(' ') : '地址不详'
}

// 房产列表数据
const realtyList = ref([])
const loading = ref(true)
const total = ref(0)

// 获取房产列表数据
const fetchData = async () => {
  loading.value = true
  try {
    const response = await queryRealtyList({
      realtyCert: searchForm.realtyCert,
      realtyType: searchForm.realtyType,
      houseType: searchForm.houseType,
      minPrice: searchForm.minPrice,
      maxPrice: searchForm.maxPrice,
      minArea: searchForm.minArea,
      maxArea: searchForm.maxArea,
      province: searchForm.province,
      city: searchForm.city,
      district: searchForm.district,
      street: searchForm.street,
      community: searchForm.community,
      unit: searchForm.unit,
      floor: searchForm.floor,
      room: searchForm.room,
      pageSize: searchForm.pageSize,
      pageNumber: searchForm.pageNumber
    })

    // 适配新的API返回格式
    if (response && response.realtyList) {
      realtyList.value = response.realtyList
      total.value = response.total || 0
    } else {
      realtyList.value = []
      total.value = 0
    }
  } catch (error) {
    ElMessage.error('获取房产列表失败，请稍后再试')
    realtyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  searchForm.pageNumber = 1
  fetchData()
}

// 改变每页显示数量
const handleSizeChange = (size) => {
  searchForm.pageSize = size
  fetchData()
}

// 改变页码
const handleCurrentChange = (page) => {
  searchForm.pageNumber = page
  fetchData()
}

// 查看详情
const viewDetails = (item) => {
  router.push(`/realty/${item.id}`)
}

// 开始交易
const startTransaction = (item) => {
  router.push({
    path: '/transaction/create',
    query: { realtyCert: item.realtyCert }
  })
}

// 编辑房产
const editRealty = (item) => {
  router.push(`/realty/${item.id}`)
}

// 初始加载
onMounted(() => {
  fetchData()
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
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.realty-type-tag {
  position: absolute;
  top: 10px;
  left: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-normal {
  background-color: rgba(25, 190, 107, 0.7);
}

.status-pending {
  background-color: #E6A23C;
}

.status-mortgaged {
  background-color: #409EFF;
}

.status-frozen {
  background-color: #F56C6C;
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
