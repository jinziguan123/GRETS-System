<template>
  <div class="realty-detail-container">
    <!-- 页面标题和操作按钮 -->
    <div class="page-header">
      <h3>房产详情</h3>
      <div class="header-actions">
        <el-button type="warning" @click="openEditDialog" v-if="canEdit">编辑</el-button>
        <el-button @click="goBack">返回</el-button>
      </div>
    </div>

    <el-row :gutter="20" v-loading="loading">
      <!-- 区域1: 房产图片展示 -->
      <el-col :span="10">
        <el-card class="image-card">
          <div v-if="realty.images && realty.images.length > 0" class="realty-images">
            <el-carousel :interval="4000" arrow="always" indicator-position="outside">
              <el-carousel-item v-for="(image, index) in realty.images" :key="index">
                <img :src="image" class="carousel-image" />
              </el-carousel-item>
            </el-carousel>
          </div>
          <div v-else class="no-images">
            <el-empty description="暂无房产图片">
              <template #image>
                <div class="empty-image">
                  <i class="el-icon-picture" style="font-size: 48px; color: #909399;"></i>
                </div>
              </template>
            </el-empty>
          </div>
          
          <!-- 房产证号和状态 -->
          <div class="property-status-bar">
            <div class="cert-number">
              <span class="label">不动产证号:</span>
              <span class="value">{{ realty.realtyCert }}</span>
            </div>
            <div class="status">
              <el-tag :type="getStatusTagType(realty.status)">
                {{ getStatusText(realty.status) }}
              </el-tag>
              <el-tag type="info" color="green" effect="dark" style="margin-left: 8px;" v-if="realty.isNewHouse">
                {{ '新房' }}
              </el-tag>
              <el-tag type="warning" effect="dark" style="margin-left: 8px;" v-if="!realty.isNewHouse">
                {{ '二手房' }}
              </el-tag>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 区域2: 房产基本信息 -->
      <el-col :span="9">
        <el-card class="info-card">
          <div class="card-title">基本信息</div>
          
          <div class="info-row">
            <div class="info-label">房产类型</div>
            <div class="info-value">{{ getRealtyTypeText(realty.realtyType) }}</div>
          </div>
          
          <div class="info-row">
            <div class="info-label">户型</div>
            <div class="info-value">{{ getHouseTypeText(realty.houseType) }}</div>
          </div>
          
          <div class="info-row">
            <div class="info-label">面积</div>
            <div class="info-value">{{ realty.area }} 平方米</div>
          </div>
          
          <div class="info-row">
            <div class="info-label">参考价格</div>
            <div class="info-value">{{ formatPrice(realty.price) }}</div>
          </div>
          
          <div class="info-row">
            <div class="info-label">地址</div>
            <div class="info-value address-value">{{ generateAddress(realty) }}</div>
          </div>
          
          <div class="info-row" v-if="realty.relContractUUID">
            <div class="info-label">关联合同</div>
            <div class="info-value">
              <el-button type="primary" link @click="openContractDialog">查看合同详情</el-button>
            </div>
          </div>
          
          <div class="info-row description-row">
            <el-row>
              <el-col :span="24">
                <div class="info-label">描述</div>
              </el-col>
              <el-col>
                <div class="info-value description-value">{{ realty.description || '暂无描述' }}</div>
              </el-col>
            </el-row>
          </div>
          
          <!-- 相关操作 -->
          <div class="actions" v-if="realty.status === 'PENDING_SALE' && hasTransactionPermission">
            <el-divider />
            <el-button type="primary" @click="startChatRoom" :disabled="isOwner">
              {{ isOwner ? '不能与自己聊天' : '咨询房主' }}
            </el-button>
          </div>
        </el-card>
      </el-col>
      
      <!-- 区域3: 房产其他信息（竖向排列） -->
      <el-col :span="5">
        <el-card class="other-info-card">
          <div class="card-title">其他信息</div>
          
          <div class="other-info-item">
            <div class="other-info-label">当前所有者ID哈希</div>
            <div class="other-info-value">{{ realty.currentOwnerCitizenIDHash || '未知' }}</div>
          </div>

          <div class="other-info-item">
            <div class="other-info-label">当前所有者组织</div>
            <div class="other-info-value">{{ getCurrentOwnerOrganization(realty.currentOwnerOrganization) || '未知' }}</div>
          </div>
          
          <div class="other-info-item">
            <div class="other-info-label">登记日期</div>
            <div class="other-info-value">{{ formatDate(realty.createTime) }}</div>
          </div>
          
          <div class="other-info-item">
            <div class="other-info-label">最后更新时间</div>
            <div class="other-info-value">{{ formatDate(realty.lastUpdateTime) }}</div>
          </div>
        </el-card>
        
        <!-- 查看交易记录按钮（最右侧居中） -->
        <div class="transaction-record-button" @click="openDrawer">
          <el-tooltip content="查看交易记录" placement="left">
            <el-button type="primary" circle>
              <el-icon><Document /></el-icon>
            </el-button>
          </el-tooltip>
          <div class="button-label">交易记录</div>
        </div>
      </el-col>
    </el-row>
    
    <!-- 交易记录抽屉 -->
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
            <el-option label="进行中" value="IN_PROGRESS"></el-option>
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
          <el-table-column prop="updateTime" label="更新时间" width="180">
            <template #default="scope">
              {{ scope.row.updateTime ? formatDate(scope.row.updateTime) : '-' }}
            </template>
          </el-table-column>
<!--          <el-table-column fixed="right" label="操作" width="120">-->
<!--            <template #default="scope">-->
<!--              <el-button type="primary" link @click="viewTransaction(scope.row.transactionUUID)">查看</el-button>-->
<!--            </template>-->
<!--          </el-table-column>-->
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
            <el-option label="正常" value="NORMAL" v-if="isOwner || isGovernment || isAudit"></el-option>
            <el-option label="挂牌" value="PENDING_SALE" v-if="isOwner"></el-option>
            <el-option label="已抵押" value="IN_MORTGAGE" v-if="isGovernment"></el-option>
            <el-option label="已冻结" value="FROZEN" v-if="isGovernment || isAudit"></el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="关联合同" prop="relContractUUID" v-if="isOwner && editForm.status === 'PENDING_SALE'" required>
          <el-select 
            v-model="editForm.relContractUUID" 
            placeholder="请选择关联合同" 
            style="width: 100%"
            @change="handleContractChange"
            filterable
            clearable
            :loading="contractLoading"
          >
            <el-option 
              v-for="contract in contractList" 
              :key="contract.contractUUID" 
              :label="contract.title" 
              :value="contract.contractUUID"
              :disabled="contract.status !== 'NORMAL'"
            ></el-option>
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
            action="http://localhost:8080/api/v1/picture/upload"
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

    <!-- 合同详情对话框 -->
    <el-dialog
      v-model="contractDialogVisible"
      title="合同详情"
      width="60%"
    >
      <div v-loading="contractLoading">
        <template v-if="currentContract">
          <div class="contract-header">
            <h2 class="contract-title">{{ currentContract.title }}</h2>
            <el-tag :type="getContractStatusType(currentContract.status)">
              {{ getContractStatusText(currentContract.status) }}
            </el-tag>
          </div>
          
          <el-descriptions :column="2" border>
            <el-descriptions-item label="合同编号">{{ currentContract.contractUUID }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(currentContract.createTime) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(currentContract.updateTime) }}</el-descriptions-item>
            <el-descriptions-item label="合同类型">{{ getContractTypeText(currentContract.contractType) }}</el-descriptions-item>
            <el-descriptions-item label="创建者" :span="2">{{ currentContract.creatorCitizenIDHash }}</el-descriptions-item>
          </el-descriptions>
          
          <div class="contract-content">
            <h3>合同内容</h3>
            <div class="content-box" v-html="currentContract.content"></div>
          </div>
        </template>
        <el-empty v-else description="暂无合同信息" />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="contractDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="downloadContract" :disabled="!currentContract">
            下载合同
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 添加验资对话框 -->
    <el-dialog
      v-model="verificationDialogVisible"
      title="添加验资"
      width="60%"
    >
      <el-form :model="verificationForm" :rules="verificationRules" ref="verificationFormRef" label-width="120px">
        <el-form-item label="验资金额" prop="verificationAmount">
          <el-input v-model="verificationForm.verificationAmount" type="number"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="verificationDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitVerification" :loading="verificationLoading">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Document } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getRealtyDetail, updateRealty } from "@/api/realty.js"
import CryptoJS from 'crypto-js'
import {queryTransactionList} from "@/api/transaction.js";
import {queryContractList, getContractDetail, getContractByUUID} from '@/api/contract'
import { verifyCapital, createChatRoom } from '@/api/chat'

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
const contractList = ref([])
const contractDialogVisible = ref(false)
const contractLoading = ref(false)
const currentContract = ref(null)

// 添加验资相关的响应式数据
const verificationDialogVisible = ref(false)
const verificationLoading = ref(false)
const verificationForm = reactive({
  verificationAmount: 100000 // 默认验资金额10万
})

// 验资表单规则
const verificationRules = {
  verificationAmount: [
    { required: true, message: '请输入验资金额', trigger: 'blur' },
    { type: 'number', min: 1, message: '验资金额必须大于0', trigger: 'blur' }
  ]
}

const verificationFormRef = ref(null)

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
  lastUpdateDate: '',
  isNewHouse: false
})

watch(realty.status, (newStatus) => {
  if (newStatus !== 'PENDING_SALE') {
    editForm.relContractUUID = ''
  }
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
  images: [],
  relContractUUID: '',
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

// 获取合同列表
const fetchContractList = async (query) => {
  const response = await queryContractList({
    status: 'NORMAL',
    pageSize: 100,
    pageNumber: 1,
    creatorCitizenID: userStore.user.citizenID,
    contractUUID: query,
    excludeAlreadyUsedFlag: true
  })
  contractList.value = response.contracts || []
}

// 上传文件列表
const fileList = ref([])

// 计算属性：是否为政府部门
const isGovernment = computed(() => {
  return userStore.hasOrganization('government')
})

// 计算属性：是否为审计
const isAudit = computed(() => {
  return userStore.hasOrganization('audit')
})

// 计算属性：是否为房产拥有者
const isOwner = computed(() => {
  // 这里需要根据实际情况判断当前用户是否为房产拥有者
  // 假设通过比较当前用户ID和房产所有者ID
  return CryptoJS.SHA256(userStore.user.citizenID).toString() === realty.currentOwnerCitizenIDHash
      && userStore.user.organization === realty.currentOwnerOrganization
})

// 计算属性：是否有编辑权限
const canEdit = computed(() => {
  return isGovernment.value || isOwner.value || isAudit.value
})

// 计算属性：是否可以编辑状态
const canEditStatus = computed(() => {
  // 政府部门可以在任何情况下修改状态
  // 房产拥有者只能在房产状态为正常时修改状态
  return isGovernment.value || (isOwner.value && (realty.status === 'NORMAL' || realty.status === 'PENDING_SALE')) || isAudit.value
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
    'IN_SALE': '在售',
    'IN_MORTGAGE': '抵押中',
    'FROZEN': '已冻结',
    'PENDING_SALE': '挂牌'
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
    'IN_PROGRESS': 'primary',
    'REJECTED': 'danger',
    'COMPLETED': 'success'
  }
  return statusMap[status] || ''
}

// 获取交易状态对应的文本
const getTransactionStatusText = (status) => {
  const statusMap = {
    'PENDING': '待处理',
    'IN_PROGRESS': '进行中',
    'REJECTED': '已拒绝',
    'COMPLETED': '已完成'
  }
  return statusMap[status] || status
}

// 获取房产所有者组织对应的文本
const getCurrentOwnerOrganization = (organization) => {
  const organizationMap = {
    'government': '政府监管部门',
    'investor': '投资者/买家',
    'bank': '银行机构',
    'thirdparty': '第三方机构',
    'audit': '审计机构'
  }
  return organizationMap[organization] || organization
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
  // 根据省份类型决定地址格式
  if (item.province.includes('市')) {
    // 直辖市格式
    return `${item.province}${item.district}区${item.street}${item.community}${item.unit}单元${item.floor}楼${item.room}号`
  } else {
    // 普通省份格式
    return `${item.province}${item.city}${item.district}区${item.street}${item.community}${item.unit}单元${item.floor}楼${item.room}号`
  }
}

// 获取房产详情
const fetchRealtyDetail = async () => {
  loading.value = true
  try {
    const response = await getRealtyDetail(route.params.id)
    Object.assign(realty, response)

    // 更新交易查询条件中的房产证号
    transactionQuery.realtyCert = realty.realtyCert

    // 初始化文件列表
    if (realty.images && realty.images.length > 0) {
      fileList.value = realty.images.map((url, index) => ({
        name: `image-${index + 1}`,
        url
      }))
    } else {
      // 如果没有图片，使用默认图片
      realty.images = ['http://localhost:8089/i/2025/04/28/680f3098ac95a.png']
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
    
    const response = await queryTransactionList(transactionQuery)
    transactions.value = response.transactions || []
    transactionTotal.value = response.total || 0
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
  editForm.relContractUUID = realty.relContractUUID || ''
  
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
        await updateRealty(realty.id, editForm)
        ElMessage.success('更新房产信息成功')
        await fetchRealtyDetail()
        editDialogVisible.value = false
        refresh()
      } catch (error) {
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

function refresh() {
  location.reload()
}

// 上传成功回调
const handleUploadSuccess = (response, file, fileList) => {
  console.log('上传成功:', response)
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

// 获取合同状态对应的Tag类型
const getContractStatusType = (status) => {
  const statusMap = {
    'NORMAL': 'success',
    'FROZEN': 'danger',
    'COMPLETED': 'info'
  }
  return statusMap[status] || ''
}

// 获取合同状态对应的文本
const getContractStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'FROZEN': '已冻结',
    'COMPLETED': '已完成'
  }
  return statusMap[status] || status
}

// 获取合同类型对应的文本
const getContractTypeText = (contractType) => {
  const typeMap = {
    'PURCHASE': '房屋买卖',
    'MORTGAGE': '抵押合同',
    'LEASE': '租赁合同',
  }
  return typeMap[contractType] || contractType
}

// 打开合同详情对话框
const openContractDialog = async () => {
  if (!realty.relContractUUID) {
    ElMessage.warning('该房产没有关联合同')
    return
  }
  
  contractDialogVisible.value = true
  contractLoading.value = true
  
  try {
    const response = await getContractByUUID(realty.relContractUUID)
    currentContract.value = response
  } catch (error) {
    console.error('获取合同详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取合同详情失败')
    currentContract.value = null
  } finally {
    contractLoading.value = false
  }
}

// 下载合同
const downloadContract = () => {
  if (!currentContract.value) return
  
  // 创建一个链接用于下载
  const element = document.createElement('a')
  const file = new Blob([currentContract.value.content], {type: 'text/plain'})
  element.href = URL.createObjectURL(file)
  element.download = `合同_${currentContract.value.contractUUID}.txt`
  document.body.appendChild(element)
  element.click()
  document.body.removeChild(element)
}

// 处理合同选择变更
const handleContractChange = (contractUUID) => {
  editForm.relContractUUID = contractUUID
}

// 打开聊天室
const startChatRoom = () => {
  if (isOwner.value) {
    ElMessage.warning('不能与自己聊天')
    return
  }
  
  // 打开验资对话框
  verificationForm.verificationAmount = Math.max(realty.price * 0.1, 100000) // 默认为房价的10%，最低10万
  verificationDialogVisible.value = true
}

// 提交验资
const submitVerification = async () => {
  if (!verificationFormRef.value) return
  
  await verificationFormRef.value.validate(async (valid) => {
    if (valid) {
      verificationLoading.value = true
      try {
        // 1. 先进行验资
        const verifyResult = await verifyCapital({
          userCitizenID: userStore.user.citizenID,
          userOrganization: userStore.user.organization,
          realtyCert: realty.realtyCert,
          verificationAmount: verificationForm.verificationAmount
        })
        
        if (!verifyResult.success) {
          ElMessage.error(verifyResult.message)
          return
        }
        
        // 2. 验资成功，创建聊天室
        const chatRoomResult = await createChatRoom({
          userCitizenID: userStore.user.citizenID,
          userOrganization: userStore.user.organization,
          realtyCert: realty.realtyCert,
          verificationAmount: verificationForm.verificationAmount
        })
        
        ElMessage.success('验资成功，聊天室已创建！')
        verificationDialogVisible.value = false
        
        // 3. 跳转到聊天室页面
        router.push(`/chat/room/${chatRoomResult.roomUUID}`)
        
      } catch (error) {
        console.error('验资或创建聊天室失败:', error)
        ElMessage.error(error.response?.data?.message || '操作失败')
      } finally {
        verificationLoading.value = false
      }
    } else {
      ElMessage.error('请检查表单填写是否正确')
    }
  })
}

onMounted(() => {
  fetchRealtyDetail()
  fetchContractList()
})
</script>

<style scoped>
.realty-detail-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h3 {
  margin: 0;
  font-size: 22px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

/* 区域1: 房产图片展示 */
.image-card {
  height: 100%;
}

.realty-images {
  margin-bottom: 15px;
}

.carousel-image {
  width: 100%;
  height: 350px;
  object-fit: cover;
}

.no-images {
  height: 350px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 15px;
}

.empty-image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 200px;
}

.property-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
  padding: 10px 0;
  border-top: 1px solid #ebeef5;
}

.cert-number .label {
  color: #909399;
  margin-right: 5px;
}

.cert-number .value {
  font-weight: bold;
}

/* 区域2: 房产基本信息 */
.info-card {
  height: 100%;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 20px;
  color: #303133;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 10px;
}

.info-row {
  display: flex;
  margin-bottom: 15px;
}

.info-label {
  flex: 0 0 100px;
  color: #909399;
}

.info-value {
  flex: 1;
  color: #303133;
}

.address-value {
  word-break: break-all;
}

.description-row {
  flex-direction: column;
}

.description-value {
  margin-top: 10px;
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  min-height: 80px;
}

/* 区域3: 房产其他信息 */
.other-info-card {
  height: calc(100% - 150px);
  margin-bottom: 20px;
}

.other-info-item {
  margin-bottom: 20px;
}

.other-info-label {
  color: #909399;
  margin-bottom: 5px;
}

.other-info-value {
  color: #303133;
  font-weight: bold;
  word-break: break-all;
}

/* 交易记录按钮 */
.transaction-record-button {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 130px;
  background-color: #f5f7fa;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.transaction-record-button:hover {
  background-color: #ecf5ff;
}

.button-label {
  margin-top: 10px;
  font-size: 14px;
  color: #409EFF;
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

/* 合同详情样式 */
.contract-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.contract-title {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.contract-content {
  margin-top: 20px;
}

.contract-content h3 {
  margin-bottom: 10px;
  font-size: 16px;
  color: #606266;
}

.content-box {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
  min-height: 200px;
  white-space: pre-wrap;
  font-family: 'Courier New', Courier, monospace;
  line-height: 1.5;
}
</style>
