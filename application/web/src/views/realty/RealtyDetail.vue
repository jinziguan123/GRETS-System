<template>
  <div class="realty-detail-container">
    <el-row :gutter="20">
      <!-- 左侧房产详情 -->
      <el-col :span="18">
        <el-card class="box-card" v-loading="loading">
          <template #header>
            <div class="card-header">
              <h3>房产详情</h3>
              <div class="header-actions">
                <el-button type="success" @click="openEditDialog" v-if="canEdit">编辑</el-button>
                <el-button @click="goBack">返回</el-button>
                <el-button type="primary" @click="openDrawer">查看交易记录</el-button>
              </div>
            </div>
          </template>
          
          <!-- 房产图片展示 -->
          <div class="realty-images" v-if="realty.images && realty.images.length > 0">
            <el-carousel :interval="4000" type="card" height="300px">
              <el-carousel-item v-for="(image, index) in realty.images" :key="index">
                <img :src="image" class="carousel-image" />
              </el-carousel-item>
            </el-carousel>
          </div>
          <div class="no-images" v-else>
            <el-empty description="暂无房产图片"></el-empty>
          </div>
          
          <!-- 房产基本信息 -->
          <el-descriptions class="margin-top" title="基本信息" :column="2" border>
            <el-descriptions-item label="不动产证号">{{ realty.realtyCert }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusTagType(realty.status)">
                {{ getStatusText(realty.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="房产类型">{{ getRealtyTypeText(realty.realtyType) }}</el-descriptions-item>
            <el-descriptions-item label="户型">{{ getHouseTypeText(realty.houseType) }}</el-descriptions-item>
            <el-descriptions-item label="面积">{{ realty.area }} 平方米</el-descriptions-item>
            <el-descriptions-item label="参考价格">{{ formatPrice(realty.price) }}</el-descriptions-item>
            <el-descriptions-item label="地址" :span="2">{{ generateAddress(realty) }}</el-descriptions-item>
          </el-descriptions>
          
          <!-- 房产其他信息 -->
          <el-descriptions class="margin-top" title="其他信息" :column="2" border>
            <el-descriptions-item label="当前所有者ID">{{ realty.currentOwnerCitizenIDHash || '未知' }}</el-descriptions-item>
            <el-descriptions-item label="登记日期">{{ formatDate(realty.registrationDate) }}</el-descriptions-item>
            <el-descriptions-item label="最后更新时间">{{ formatDate(realty.lastUpdateDate) }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ realty.description || '暂无描述' }}</el-descriptions-item>
          </el-descriptions>
          
          <!-- 相关操作 -->
          <div class="actions" v-if="realty.status === 'NORMAL' && hasTransactionPermission">
            <el-divider />
            <el-button type="primary" @click="createTransaction">发起交易</el-button>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧交易记录抽屉 -->
      <el-drawer
        v-model="drawer"
        title="房产交易记录"
        direction="rtl"
        size="50%"
      >
        <!-- 交易记录查询条件 -->
        <el-form :model="transactionQuery" label-width="100px" class="demo-form-inline">
          <el-form-item label="交易状态">
            <el-select v-model="transactionQuery.status" placeholder="选择交易状态" clearable style="width: 100%">
              <el-option label="待处理" value="PENDING"></el-option>
              <el-option label="已批准" value="APPROVED"></el-option>
              <el-option label="已拒绝" value="REJECTED"></el-option>
              <el-option label="已完成" value="COMPLETED"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="fetchTransactionHistory">查询</el-button>
            <el-button @click="resetTransactionQuery">重置</el-button>
          </el-form-item>
        </el-form>
        
        <!-- 交易记录列表 -->
        <div v-if="transactionLoading" class="loading-container">
          <el-skeleton :rows="5" animated />
        </div>
        <div v-else-if="transactions.length === 0" class="empty-data">
          <el-empty description="暂无交易记录" />
        </div>
        <div v-else>
          <el-table :data="transactions" style="width: 100%" stripe>
            <el-table-column prop="transactionUUID" label="交易ID" width="180" />
            <el-table-column prop="sellerCitizenIDHash" label="卖方ID" width="120">
              <template #default="scope">
                {{ scope.row.sellerCitizenIDHash || '未知' }}
              </template>
            </el-table-column>
            <el-table-column prop="buyerCitizenIDHash" label="买方ID" width="120">
              <template #default="scope">
                {{ scope.row.buyerCitizenIDHash || '未知' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getTransactionStatusTagType(scope.row.status)">
                  {{ getTransactionStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.createTime) }}
              </template>
            </el-table-column>
            <el-table-column prop="completedTime" label="完成时间" width="180">
              <template #default="scope">
                {{ scope.row.completedTime ? formatDate(scope.row.completedTime) : '-' }}
              </template>
            </el-table-column>
            <el-table-column fixed="right" label="操作" width="120">
              <template #default="scope">
                <el-button type="primary" link @click="viewTransaction(scope.row.transactionUUID)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页 -->
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="transactionQuery.pageNumber"
              v-model:page-size="transactionQuery.pageSize"
              :page-sizes="[10, 20, 30, 50]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="transactionTotal"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </el-drawer>
    </el-row>
    
    <!-- 编辑房产对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑房产信息"
      width="60%"
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="120px">
        <el-form-item label="不动产证号" prop="realtyCert">
          <el-input v-model="editForm.realtyCert" disabled></el-input>
        </el-form-item>
        
        <el-form-item label="房产类型" prop="realtyType" v-if="isGovernment">
          <el-select v-model="editForm.realtyType" placeholder="请选择房产类型" style="width: 100%">
            <el-option label="住宅" value="HOUSE"></el-option>
            <el-option label="商铺" value="SHOP"></el-option>
            <el-option label="办公" value="OFFICE"></el-option>
            <el-option label="工业" value="INDUSTRIAL"></el-option>
            <el-option label="其他" value="OTHER"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态" prop="status" v-if="canEditStatus">
          <el-select v-model="editForm.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="正常" value="NORMAL"></el-option>
            <el-option label="交易中" value="IN_TRANSACTION"></el-option>
            <el-option label="已抵押" value="MORTGAGED"></el-option>
            <el-option label="已冻结" value="FROZEN"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="价格" prop="price" v-if="isOwner && realty.status === 'NORMAL'">
          <el-input-number v-model="editForm.price" :min="0" :step="10000" style="width: 100%"></el-input-number>
        </el-form-item>
        
        <el-form-item label="户型" prop="houseType" v-if="isOwner && realty.status === 'NORMAL'">
          <el-select v-model="editForm.houseType" placeholder="请选择户型" style="width: 100%">
            <el-option label="一室" value="single"></el-option>
            <el-option label="两室" value="double"></el-option>
            <el-option label="三室" value="triple"></el-option>
            <el-option label="四室及以上" value="multiple"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description" v-if="isOwner && realty.status === 'NORMAL'">
          <el-input v-model="editForm.description" type="textarea" :rows="3"></el-input>
        </el-form-item>
        
        <el-form-item label="图片" prop="images" v-if="isOwner && realty.status === 'NORMAL'">
          <el-upload
            action="/api/upload"
            list-type="picture-card"
            :file-list="fileList"
            :on-preview="handlePictureCardPreview"
            :on-remove="handleRemove"
            :on-success="handleUploadSuccess"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <el-dialog v-model="previewVisible" append-to-body>
            <img w-full :src="previewUrl" alt="Preview Image" />
          </el-dialog>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEditForm" :loading="submitLoading">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getRealtyByID, updateRealty, queryTransactionList } from "@/views/realty/api.js"

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(true)
const drawer = ref(false)
const editDialogVisible = ref(false)
const submitLoading = ref(false)
const previewVisible = ref(false)
const previewUrl = ref('')
const editFormRef = ref(null)
const transactionLoading = ref(false)
const transactionTotal = ref(0)

// 房产信息
const realty = reactive({
  id: 0,
  realtyCert: '',
  realtyType: '',
  price: 0,
  area: 0,
  province: '',
  city: '',
  district: '',
  street: '',
  community: '',
  unit: '',
  floor: '',
  room: '',
  houseType: '',
  status: '',
  description: '',
  images: [],
  currentOwnerCitizenIDHash: '',
  registrationDate: '',
  lastUpdateDate: ''
})

// 交易记录
const transactions = ref([])

// 交易查询条件
const transactionQuery = reactive({
  realtyCert: '',
  status: '',
  pageSize: 10,
  pageNumber: 1
})

// 编辑表单
const editForm = reactive({
  realtyCert: '',
  realtyType: '',
  status: '',
  price: 0,
  houseType: '',
  description: '',
  images: []
})

// 编辑表单校验规则
const editRules = {
  realtyType: [
    { required: true, message: '请选择房产类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格必须大于0', trigger: 'blur' }
  ]
}

// 上传文件列表
const fileList = ref([])

// 计算属性：是否为政府部门
const isGovernment = computed(() => {
  return userStore.hasOrganization('government')
})

// 计算属性：是否为房产拥有者
const isOwner = computed(() => {
  // 这里需要根据实际情况判断当前用户是否为房产拥有者
  // 假设通过比较当前用户ID和房产所有者ID
  return userStore.citizenID === realty.currentOwnerCitizenIDHash
})

// 计算属性：是否有编辑权限
const canEdit = computed(() => {
  return isGovernment.value || isOwner.value
})

// 计算属性：是否可以编辑状态
const canEditStatus = computed(() => {
  // 政府部门可以在任何情况下修改状态
  // 房产拥有者只能在房产状态为正常时修改状态
  return isGovernment.value || (isOwner.value && realty.status === 'NORMAL')
})

// 计算属性：是否有交易权限
const hasTransactionPermission = computed(() => {
  return userStore.hasOrganization('investor')
})

// 获取房产状态对应的Tag类型
const getStatusTagType = (status) => {
  const statusMap = {
    'NORMAL': 'success',
    'IN_TRANSACTION': 'warning',
    'MORTGAGED': 'info',
    'FROZEN': 'danger'
  }
  return statusMap[status] || ''
}

// 获取房产状态对应的文本
const getStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 获取房产类型文本
const getRealtyTypeText = (type) => {
  const typeMap = {
    'HOUSE': '住宅',
    'SHOP': '商铺',
    'OFFICE': '办公',
    'INDUSTRIAL': '工业',
    'OTHER': '其他'
  }
  return typeMap[type] || type
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

// 获取交易状态对应的Tag类型
const getTransactionStatusTagType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'APPROVED': 'success',
    'REJECTED': 'danger',
    'COMPLETED': 'success'
  }
  return statusMap[status] || ''
}

// 获取交易状态对应的文本
const getTransactionStatusText = (status) => {
  const statusMap = {
    'PENDING': '待处理',
    'APPROVED': '已批准',
    'REJECTED': '已拒绝',
    'COMPLETED': '已完成'
  }
  return statusMap[status] || status
}

// 格式化价格
const formatPrice = (price) => {
  return `¥ ${price.toLocaleString()} 元`
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString()
}

// 生成地址信息
const generateAddress = (item) => {
  return `${item.province || ''}${item.province ? '省' : ''}${item.city || ''}${item.city ? '市' : ''}${item.district || ''}${item.district ? '区' : ''}${item.street || ''}${item.community || ''}${item.unit || ''}${item.unit ? '单元' : ''}${item.floor || ''}${item.floor ? '楼' : ''}${item.room || ''}`
}

// 获取房产详情
const fetchRealtyDetail = async () => {
  loading.value = true
  try {
    const response = await getRealtyByID(route.params.id)
    if (response.code === 200) {
      Object.assign(realty, response.data)
      
      // 更新交易查询条件中的房产证号
      transactionQuery.realtyCert = realty.realtyCert
      
      // 初始化文件列表
      if (realty.images && realty.images.length > 0) {
        fileList.value = realty.images.map((url, index) => ({
          name: `image-${index + 1}`,
          url
        }))
      }
    } else {
      ElMessage.error(response.message || '获取房产详情失败')
    }
  } catch (error) {
    console.error('获取房产详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取房产详情失败')
  } finally {
    loading.value = false
  }
}

// 打开抽屉并获取交易历史
const openDrawer = () => {
  drawer.value = true
  fetchTransactionHistory()
}

// 获取交易历史
const fetchTransactionHistory = async () => {
  transactionLoading.value = true
  try {
    // 只传递非空参数
    const params = { ...transactionQuery }
    Object.keys(params).forEach(key => {
      if (params[key] === '') {
        delete params[key]
      }
    })
    
    const response = await queryTransactionList(params)
    if (response.code === 200) {
      transactions.value = response.data.list || []
      transactionTotal.value = response.data.total || 0
    } else {
      ElMessage.error(response.message || '获取交易历史失败')
      transactions.value = []
      transactionTotal.value = 0
    }
  } catch (error) {
    console.error('获取交易历史失败:', error)
    ElMessage.error('获取交易历史失败')
    transactions.value = []
    transactionTotal.value = 0
  } finally {
    transactionLoading.value = false
  }
}

// 重置交易查询条件
const resetTransactionQuery = () => {
  Object.keys(transactionQuery).forEach(key => {
    if (key !== 'realtyCert' && key !== 'pageSize' && key !== 'pageNumber') {
      transactionQuery[key] = ''
    }
  })
  transactionQuery.pageNumber = 1
  fetchTransactionHistory()
}

// 改变每页显示数量
const handleSizeChange = (size) => {
  transactionQuery.pageSize = size
  fetchTransactionHistory()
}

// 改变页码
const handleCurrentChange = (page) => {
  transactionQuery.pageNumber = page
  fetchTransactionHistory()
}

// 打开编辑对话框
const openEditDialog = () => {
  // 初始化编辑表单数据
  editForm.realtyCert = realty.realtyCert
  editForm.realtyType = realty.realtyType
  editForm.status = realty.status
  editForm.price = realty.price
  editForm.houseType = realty.houseType
  editForm.description = realty.description
  editForm.images = realty.images || []
  
  editDialogVisible.value = true
  
  // 等待DOM更新后设置表单引用
  nextTick(() => {
    if (editFormRef.value) {
      editFormRef.value.clearValidate()
    }
  })
}

// 提交编辑表单
const submitEditForm = async () => {
  if (!editFormRef.value) return
  
  await editFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        // 构建更新请求数据
        const updateData = {
          realtyCert: editForm.realtyCert
        }
        
        // 根据权限添加不同字段
        if (isGovernment.value) {
          updateData.realtyType = editForm.realtyType
          updateData.status = editForm.status
        } else if (isOwner.value && realty.status === 'NORMAL') {
          updateData.price = editForm.price
          updateData.houseType = editForm.houseType
          updateData.description = editForm.description
          updateData.images = editForm.images
          
          // 如果拥有者可以修改状态
          if (canEditStatus.value) {
            updateData.status = editForm.status
          }
        }
        
        const response = await updateRealty(realty.id, updateData)
        if (response.code === 200) {
          ElMessage.success('更新房产信息成功')
          editDialogVisible.value = false
          // 刷新房产详情
          fetchRealtyDetail()
        } else {
          ElMessage.error(response.message || '更新房产信息失败')
        }
      } catch (error) {
        console.error('更新房产信息失败:', error)
        ElMessage.error(error.response?.data?.message || '更新房产信息失败')
      } finally {
        submitLoading.value = false
      }
    } else {
      ElMessage.error('请检查表单填写是否正确')
    }
  })
}

// 图片预览
const handlePictureCardPreview = (file) => {
  previewUrl.value = file.url
  previewVisible.value = true
}

// 移除图片
const handleRemove = (file, fileList) => {
  // 更新编辑表单中的图片列表
  editForm.images = fileList.map(file => file.url)
}

// 上传成功回调
const handleUploadSuccess = (response, file, fileList) => {
  if (response.code === 200) {
    // 假设上传成功后返回图片URL
    const imageUrl = response.data.url
    // 更新编辑表单中的图片列表
    editForm.images.push(imageUrl)
  } else {
    ElMessage.error(response.message || '上传图片失败')
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 查看交易详情
const viewTransaction = (transactionId) => {
  router.push(`/transaction/detail/${transactionId}`)
}

// 创建交易
const createTransaction = () => {
  router.push(`/transaction/create?realtyCert=${realty.realtyCert}`)
}

onMounted(() => {
  fetchRealtyDetail()
})
</script>

<style scoped>
.realty-detail-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.box-card {
  margin-bottom: 20px;
}

.margin-top {
  margin-top: 20px;
}

.realty-images {
  margin-bottom: 20px;
}

.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-images {
  height: 200px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 20px;
}

.transaction-history {
  margin-top: 30px;
}

.loading-container {
  padding: 20px 0;
}

.empty-data {
  padding: 40px 0;
  text-align: center;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.actions {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
