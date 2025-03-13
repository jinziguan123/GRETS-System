<template>
  <div class="realty-list-container">
    <div class="page-header">
      <h1>房产列表</h1>
      <el-button type="primary" @click="$router.push('/realty/create')" v-if="hasPermission(['GovernmentMSP'])">
        添加房产
      </el-button>
    </div>
    
    <!-- 搜索条件 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent="handleSearch">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="输入房产名称/地址" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="房产类型" clearable>
            <el-option label="住宅" value="住宅" />
            <el-option label="商铺" value="商铺" />
            <el-option label="写字楼" value="写字楼" />
            <el-option label="厂房" value="厂房" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="房产状态" clearable>
            <el-option label="可用" value="available" />
            <el-option label="已售" value="sold" />
            <el-option label="抵押中" value="mortgaged" />
            <el-option label="冻结" value="frozen" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 房产列表 -->
    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <span>房产列表</span>
          <div>
            <el-switch
              v-model="viewMode"
              active-text="卡片视图"
              inactive-text="表格视图"
              inline-prompt
              style="margin-right: 10px;"
            />
            <el-button :icon="Refresh" circle @click="fetchList" />
          </div>
        </div>
      </template>
      
      <!-- 表格视图 -->
      <el-table
        v-if="viewMode === false"
        :data="realtyList"
        border
        stripe
        style="width: 100%"
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="名称" min-width="150" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="address" label="地址" min-width="200" />
        <el-table-column prop="area" label="面积(㎡)" width="100" />
        <el-table-column prop="price" label="价格(元)" width="120">
          <template #default="scope">
            {{ scope.row.price ? scope.row.price.toLocaleString() : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              @click="$router.push(`/realty/detail/${scope.row.id}`)"
            >
              查看
            </el-button>
            <el-button 
              v-if="hasPermission(['GovernmentMSP'])"
              type="warning" 
              link 
              @click="handleEdit(scope.row)"
            >
              编辑
            </el-button>
            <el-popconfirm
              v-if="hasPermission(['GovernmentMSP'])"
              title="确定删除该房产记录吗？"
              @confirm="handleDelete(scope.row.id)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 卡片视图 -->
      <div v-else class="card-view">
        <el-row :gutter="20">
          <el-col 
            v-for="item in realtyList" 
            :key="item.id" 
            :xs="24" 
            :sm="12" 
            :md="8" 
            :lg="6"
          >
            <el-card 
              class="realty-card" 
              shadow="hover" 
              @click="$router.push(`/realty/detail/${item.id}`)"
            >
              <img 
                src="https://via.placeholder.com/300x200" 
                class="realty-image"
              />
              <div class="realty-info">
                <h3 class="realty-title">{{ item.title }}</h3>
                <p class="realty-address">{{ item.address }}</p>
                <div class="realty-details">
                  <span class="realty-area">{{ item.area }}㎡</span>
                  <span class="realty-type">{{ item.type }}</span>
                  <el-tag 
                    :type="getStatusType(item.status)" 
                    size="small"
                  >
                    {{ getStatusText(item.status) }}
                  </el-tag>
                </div>
                <div class="realty-price">
                  ¥{{ item.price ? item.price.toLocaleString() : '-' }}
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { realtyApi } from '@/api'

const userStore = useUserStore()
const loading = ref(false)
const realtyList = ref([])
const viewMode = ref(false) // false:表格视图, true:卡片视图

// 权限检查
const hasPermission = (roles) => {
  if (!roles || roles.length === 0) return true
  return roles.includes(userStore.userRole)
}

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: '',
  status: ''
})

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 获取房产状态类型
const getStatusType = (status) => {
  const statusMap = {
    'available': 'success',
    'sold': 'info',
    'mortgaged': 'warning',
    'frozen': 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取房产状态文本
const getStatusText = (status) => {
  const statusMap = {
    'available': '可用',
    'sold': '已售',
    'mortgaged': '抵押中',
    'frozen': '冻结'
  }
  return statusMap[status] || status
}

// 获取房产列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    
    const res = await realtyApi.getRealtyList(params)
    realtyList.value = res.data || []
    pagination.total = res.total || 0
  } catch (error) {
    console.error('获取房产列表失败:', error)
    ElMessage.error('获取房产列表失败，请刷新页面重试')
    
    // 使用模拟数据（仅用于演示）
    realtyList.value = [
      { id: 'R001', title: '阳光花园 3室2厅', type: '住宅', address: '杭州市西湖区文三路138号阳光花园3幢1单元601', area: 120, price: 1200000, status: 'available' },
      { id: 'R002', title: '江南名府 4室2厅', type: '住宅', address: '杭州市滨江区滨盛路1509号江南名府12幢2单元1801', area: 180, price: 2500000, status: 'available' },
      { id: 'R003', title: '城市广场 商铺A12', type: '商铺', address: '杭州市江干区凯旋路166号城市广场A区12号', area: 85, price: 3500000, status: 'mortgaged' }
    ]
    pagination.total = realtyList.value.length
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchList()
}

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.currentPage = 1
  fetchList()
}

// 编辑房产
const handleEdit = (row) => {
  ElMessage.info('编辑房产功能正在开发中')
}

// 删除房产
const handleDelete = async (id) => {
  try {
    await realtyApi.deleteRealty(id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除房产失败:', error)
    ElMessage.error('删除失败，请稍后重试')
  }
}

// 分页大小改变
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchList()
}

// 页码改变
const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchList()
}

// 初始化
onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.realty-list-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
}

.search-card,
.list-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.card-view {
  min-height: 300px;
}

.realty-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.3s;
}

.realty-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
}

.realty-image {
  width: 100%;
  height: 180px;
  object-fit: cover;
  border-radius: 4px;
}

.realty-info {
  padding: 15px 0 5px;
}

.realty-title {
  margin: 0 0 10px;
  font-size: 16px;
  font-weight: bold;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.realty-address {
  margin: 0 0 10px;
  color: #606266;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.realty-details {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.realty-area,
.realty-type {
  margin-right: 10px;
  font-size: 14px;
  color: #909399;
}

.realty-price {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}
</style>
