<template>
  <div class="contract-list-container">
    <div class="page-header">
      <h2>合同管理</h2>
      <div class="header-actions">
        <el-button type="primary" v-if="hasPermissionToCreate" @click="createContract">创建合同</el-button>
        <el-button @click="refreshContracts">刷新数据</el-button>
      </div>
    </div>
    
    <!-- 筛选器 -->
    <el-card class="filter-card">
      <el-form :model="filterForm" class="contract-filter">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="合同编号">
              <el-input style="width: 200px" v-model="filterForm.contractUUID" placeholder="输入合同编号" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="合同状态">
              <el-select style="width: 200px" v-model="filterForm.status" placeholder="选择合同状态" clearable>
                <el-option label="正常" value="NORMAL" />
                <el-option label="冻结" value="FROZEN" />
                <el-option label="已完成" value="COMPLETED" />
                <el-option label="进行中" value="IN_PROGRESS" />
                <el-option label="已取消" value="CANCEL" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合同类型">
              <el-select style="width: 200px" v-model="filterForm.contractType" placeholder="选择合同类型" clearable>
                <el-option label="购房合同" value="purchase" />
                <el-option label="贷款合同" value="mortgage" />
                <el-option label="租赁合同" value="lease" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item>
            <el-button type="primary" @click="searchContracts">查询</el-button>
            <el-button @click="resetFilter">重置</el-button>
          </el-form-item>
        </el-row>
      </el-form>
    </el-card>
    
    <!-- 合同表格 -->
    <el-card class="contract-table-card">
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="10" animated />
      </div>
      
      <template v-else>
        <el-table 
          :data="contractsData" 
          style="width: 100%"
          row-key="id"
          v-loading="tableLoading"
        >
          <el-table-column prop="contractUUID" label="合同编号" width="180" />
          <el-table-column prop="contractType" label="合同类型" width="120">
            <template #default="scope">
              <el-tag 
                :type="getContractTypeTag(scope.row.contractType)"
                effect="plain"
              >
                {{ getContractTypeName(scope.row.contractType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="title" label="合同标题"/>
          <el-table-column prop="creatorCitizenIDHash" label="创建者ID" width="300px" />
          <el-table-column prop="createTime" label="创建日期" width="120">
            <template #default="scope">
              {{ formatDate(scope.row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="updateTime" label="更新日期" width="120">
            <template #default="scope">
              {{ formatDate(scope.row.updateTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.status)">
                {{ getStatusName(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button 
                type="primary" 
                link 
                @click.stop="viewContractDetail(scope.row.contractUUID)"
              >
                查看
              </el-button>
              <el-button 
                v-if="canSign(scope.row)"
                type="success" 
                link 
                @click.stop="signContract(scope.row)"
              >
                签署
              </el-button>
              <el-button 
                v-if="canAudit(scope.row)"
                type="warning" 
                link 
                @click.stop="auditContract(scope.row)"
              >
                审核
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="totalItems"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </template>
    </el-card>
    
    <!-- 签署合同对话框 -->
    <el-dialog
      v-model="signDialogVisible"
      title="合同签署"
      width="600px"
      destroy-on-close
    >
      <div v-if="currentContract" class="sign-form">
        <p>您正在签署合同《{{ currentContract.title }}》。</p>
        <p>请确认以下信息无误：</p>
        <ul>
          <li>合同编号：{{ currentContract.contractUUID }}</li>
          <li>合同类型：{{ getContractTypeName(currentContract.contractType) }}</li>
          <li>创建时间：{{ formatDate(currentContract.createTime) }}</li>
          <li>签署身份：{{ userStore.organization === 'investor' ? '买方' : '卖方' }}</li>
        </ul>
        <p class="sign-tip">点击"确认签署"，即表示您已阅读并同意本合同的全部条款。</p>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="signDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmSign" :loading="signing">
            确认签署
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { queryContractList, signContract } from '@/api/contract'

const router = useRouter()
const userStore = useUserStore()

// 加载状态
const loading = ref(false)
const tableLoading = ref(false)
const signing = ref(false)

// 对话框状态
const signDialogVisible = ref(false)

// 当前选中的合同
const currentContract = ref(null)

// 筛选表单
const filterForm = reactive({
  status: '',
  contractType: '',
  contractUUID: ''
})

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

// 合同数据
const contractsData = ref([])

// 判断是否有创建合同的权限
const hasPermissionToCreate = computed(() => {
  return userStore.hasOrganization(['investor'])
})

// 判断是否可以签署合同
const canSign = (contract) => {
  if (contract.status !== 'PENDING') return false
  
  // 根据组织类型判断签署权限
  if (userStore.hasOrganization('investor') && !contract.buyerSigned) {
    return true
  }
  
  if (userStore.hasOrganization('government') && !contract.sellerSigned) {
    return true
  }
  
  return false
}

// 判断是否可以审核合同
const canAudit = (contract) => {
  return userStore.hasOrganization('audit') && contract.status === 'SIGNED'
}

// 重置筛选条件
const resetFilter = () => {
  filterForm.status = ''
  filterForm.contractType = ''
  filterForm.contractUUID = ''
  searchContracts()
}

// 搜索合同
const searchContracts = () => {
  tableLoading.value = true
  fetchContractsData()
}

// 刷新合同列表
const refreshContracts = () => {
  tableLoading.value = true
  fetchContractsData()
}

// 获取合同数据
const fetchContractsData = async () => {
  try {
    // 构建查询参数
    const params = {
      status: filterForm.status,
      contractType: filterForm.contractType,
      contractUUID: filterForm.contractUUID,
      pageSize: pageSize.value,
      pageNumber: currentPage.value,
      creatorCitizenID: userStore.user.citizenID
    }
    
    tableLoading.value = true
    
    // 调用API
    const response = await queryContractList(params)
    
    if (response && response.contracts) {
      contractsData.value = response.contracts || []
      totalItems.value = response.total || 0
    } else {
      contractsData.value = []
      totalItems.value = 0
      ElMessage.warning('未获取到合同数据')
    }
  } catch (error) {
    console.error('Failed to fetch contracts:', error)
    ElMessage.error('获取合同列表失败')
    contractsData.value = []
    totalItems.value = 0
  } finally {
    tableLoading.value = false
    loading.value = false
  }
}

// 处理页面大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchContractsData()
}

// 处理页码变化
const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchContractsData()
}

// 处理行点击
const handleRowClick = (row) => {
  viewContractDetail(row.contractUUID)
}

// 创建合同
const createContract = () => {
  router.push('/contract/create')
}

// 查看合同详情
const viewContractDetail = (contractUUID) => {
  router.push(`/contract/${contractUUID}`)
}

// 确认签署
const confirmSign = async () => {
  if (!currentContract.value) return
  
  signing.value = true
  
  try {
    const signerType = userStore.hasOrganization('investor') ? 'buyer' : 'seller'
    
    const response = await signContract(currentContract.value.contractUUID, {
      signerType
    })
    
    ElMessage.success('合同签署成功')
    signDialogVisible.value = false
    refreshContracts()
  } catch (error) {
    console.error('Failed to sign contract:', error)
    ElMessage.error('合同签署失败')
  } finally {
    signing.value = false
  }
}

// 审核合同
const auditContract = (contract) => {
  router.push({
    path: '/contract/audit',
    query: { id: contract.contractUUID }
  })
}

// 获取合同类型标签类型
const getContractTypeTag = (type) => {
  const tagMap = {
    PURCHASE: 'success',
    MORTGAGE: 'primary',
    LEASE: 'primary'
  }
  return tagMap[type] || ''
}

// 获取合同类型名称
const getContractTypeName = (type) => {
  const nameMap = {
    PURCHASE: '购房合同',
    MORTGAGE: '贷款合同',
    LEASE: '租赁合同'
  }
  return nameMap[type] || type
}

// 获取状态标签类型
const getStatusTag = (status) => {
  const tagMap = {
    'NORMAL': 'primary',
    'FROZEN': 'danger',
    'COMPLETED': 'success',
    'IN_PROGRESS': 'warning'
  }
  return tagMap[status] || ''
}

// 获取状态名称
const getStatusName = (status) => {
  const nameMap = {
    'NORMAL': '正常',
    'FROZEN': '冻结',
    'COMPLETED': '已完成',
    'IN_PROGRESS': '进行中'
  }
  return nameMap[status] || status
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化货币
const formatCurrency = (value) => {
  if (!value && value !== 0) return '-'
  return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 页面加载时获取合同数据
onMounted(() => {
  loading.value = true
  fetchContractsData()
})
</script>

<style scoped>
.contract-list-container {
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

.header-actions {
  display: flex;
  gap: 10px;
}

.filter-card {
  margin-bottom: 20px;
}

.contract-filter {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.loading-container {
  padding: 20px 0;
}

.contract-table-card {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.sign-form {
  padding: 0 20px;
}

.sign-form ul {
  background-color: #f8f8f8;
  padding: 15px;
  border-radius: 4px;
  margin: 15px 0;
}

.sign-tip {
  margin-top: 20px;
  color: #E6A23C;
  font-size: 14px;
}
</style>
