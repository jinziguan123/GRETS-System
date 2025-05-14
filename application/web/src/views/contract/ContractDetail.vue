<template>
  <div class="contract-detail-container">
    <div class="page-header">
      <div class="title-section">
        <h2>合同详情</h2>
        <div v-if="contract">
          <el-tag :type="getStatusTag(contract.status)">
            {{ getStatusName(contract.status) }}
          </el-tag>
        </div>
      </div>
      <div class="header-actions">
        <el-button 
          v-if="canAudit"
          type="warning" 
          @click="openUpdateStatusDialog"
        >
          修改状态
        </el-button>
        <el-button @click="goBack">返回</el-button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>
    
    <template v-else-if="contract">
      <el-tabs type="border-card" class="contract-tabs">
        <el-tab-pane label="基本信息">
          <el-card>
            <el-descriptions title="合同信息" :column="2" border>
              <el-descriptions-item label="合同编号">{{ contract.contractUUID }}</el-descriptions-item>
              <el-descriptions-item label="合同类型">{{ getContractTypeName(contract.contractType) }}</el-descriptions-item>
              <el-descriptions-item label="合同标题">{{ contract.title }}</el-descriptions-item>
              <el-descriptions-item label="创建日期">{{ formatDate(contract.createTime) }}</el-descriptions-item>
              <el-descriptions-item label="关联交易">{{ contract.transactionUUID }}</el-descriptions-item>
              <el-descriptions-item label="合同状态">
                <el-tag :type="getStatusTag(contract.status)">
                  {{ getStatusName(contract.status) }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
            
            <el-divider content-position="left">交易信息</el-divider>
            
            <el-descriptions :column="2" border>
              <el-descriptions-item label="买方">{{ transaction?.buyerCitizenIDHash }}</el-descriptions-item>
              <el-descriptions-item label="卖方">{{ transaction?.sellerCitizenIDHash }}</el-descriptions-item>
              <el-descriptions-item label="买方组织">{{ getCurrentOwnerOrganization(transaction?.buyerOrganization) }}</el-descriptions-item>
              <el-descriptions-item label="卖方组织">{{ getCurrentOwnerOrganization(transaction?.sellerOrganization) }}</el-descriptions-item>
              <el-descriptions-item label="交易金额">¥{{ formatCurrency(transaction?.amount) }}</el-descriptions-item>
              <el-descriptions-item label="生效日期" :span="1">{{ formatDate(transaction?.createTime) || '未生效' }}</el-descriptions-item>
            </el-descriptions>
            
            <el-divider content-position="left">相关资产</el-divider>
            
            <el-descriptions :column="1" border v-if="contract.property">
              <el-descriptions-item label="资产名称">{{ contract.property.title }}</el-descriptions-item>
              <el-descriptions-item label="交付日期">{{ formatDate(contract.property.deliveryDate) }}</el-descriptions-item>
            </el-descriptions>
            <el-empty v-else description="暂无关联资产信息" />
          </el-card>
        </el-tab-pane>
        
        <el-tab-pane label="合同内容">
          <el-card>
            <div v-if="contract.content" class="contract-content">
              <template v-if="contract.content.sections">
                <div v-for="(section, sIndex) in contract.content.sections" :key="`section-${sIndex}`" class="contract-section">
                  <h3>{{ section.title }}</h3>
                  <div v-for="(clause, cIndex) in section.clauses" :key="`clause-${sIndex}-${cIndex}`" class="contract-clause">
                    <h4>{{ clause.title }}</h4>
                    <p>{{ clause.content }}</p>
                  </div>
                </div>
              </template>
              <div v-else class="contract-text">
                {{ contract.content }}
              </div>
            </div>
            <el-empty v-else description="暂无合同内容" />
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </template>
    
    <el-empty v-else description="未找到合同信息" />
    
    <!-- 签署合同对话框 -->
    <el-dialog
      v-model="signDialogVisible"
      title="合同签署"
      width="600px"
      destroy-on-close
    >
      <div v-if="contract" class="sign-form">
        <p>您正在签署合同《{{ contract.title }}》。</p>
        <p>请确认以下信息无误：</p>
        <ul>
          <li>合同编号：{{ contract.id }}</li>
          <li>合同金额：¥{{ formatCurrency(contract.amount) }}</li>
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
    
    <!-- 修改合同状态对话框 -->
    <el-dialog
      v-model="updateStatusDialogVisible"
      title="修改合同状态"
      width="450px"
      destroy-on-close
    >
      <div class="update-status-form">
        <el-form :model="updateStatusForm" label-width="100px">
          <el-form-item label="当前状态">
            <el-tag :type="getStatusTag(contract?.status)">
              {{ getStatusName(contract?.status) }}
            </el-tag>
          </el-form-item>
          
          <el-form-item label="新状态" prop="newStatus">
            <el-select v-model="updateStatusForm.newStatus" placeholder="请选择新状态" style="width: 100%">
              <el-option label="正常" value="NORMAL"></el-option>
              <el-option label="冻结" value="FROZEN"></el-option>
              <el-option label="取消" value="CANCEL"></el-option>
              <el-option label="已完成" value="COMPLETED"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="修改原因" prop="reason">
            <el-input 
              v-model="updateStatusForm.reason" 
              type="textarea" 
              :rows="3" 
              placeholder="请输入修改原因"
            ></el-input>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="updateStatusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmUpdateStatus" :loading="updating">
            确认修改
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user'
import {ElMessage} from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'
import {getContractByUUID, updateContractStatus} from "@/api/contract.js";
import {getTransactionDetail} from "@/api/transaction.js";

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 加载状态
const loading = ref(true)
const signing = ref(false)
const updating = ref(false)

// 合同信息
const contract = ref(null)
const signDialogVisible = ref(false)
const transaction = ref(null)

// 合同ID
const contractUUID = computed(() => route.params.id)
const transactionUUID = ref(null)

// 修改状态相关
const updateStatusDialogVisible = ref(false)
const updateStatusForm = ref({
  newStatus: '',
  reason: ''
})

// 判断是否可以签署合同
const canSign = computed(() => {
  if (!contract.value || contract.value.status !== 'pending') return false
  
  // 根据组织类型判断签署权限
  if (userStore.hasOrganization('investor') && !contract.value.buyerSigned) {
    return true
  }
  
  if (userStore.hasOrganization('government') && !contract.value.sellerSigned) {
    return true
  }
  
  return false
})

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

// 判断是否可以审核合同
const canAudit = computed(() => {
  return userStore.hasOrganization('audit') && 
         contract.value && 
         ['NORMAL', 'FROZEN', 'PENDING', 'SIGNED'].includes(contract.value.status)
})

// 获取合同详情
const fetchContractDetail = async () => {
  try {
    loading.value = true
    
    // 调用API获取合同详情
    contract.value = await getContractByUUID(contractUUID.value)
    transactionUUID.value = contract.value.transactionUUID
    if (transactionUUID.value) {
      await getTransactionInfo()
    }
  } catch (error) {
    console.error('Failed to fetch contract details:', error)
    ElMessage.error('获取合同详情失败')
  } finally {
    loading.value = false
  }
}

// 获取交易信息
const getTransactionInfo = async () => {
  try {
    const response = await getTransactionDetail(transactionUUID.value)
    transaction.value = response.transaction
  } catch (error) {
    console.error('Failed to get transaction info:', error)
    ElMessage.error('获取交易信息失败')
  }
}

// 显示签署对话框
const showSignDialog = () => {
  signDialogVisible.value = true
}

// 确认签署
const confirmSign = async () => {
  signing.value = true
  
  try {
    const signerType = userStore.hasOrganization('investor') ? 'buyer' : 'seller'
    
    const { data } = await axios.post(`/contracts/${contractUUID.value}/sign`, {
      signerType
    })
    
    if (data.code === 200) {
      ElMessage.success('合同签署成功')
      signDialogVisible.value = false
      // 刷新合同详情
      fetchContractDetail()
    } else {
      ElMessage.error(data.message || '合同签署失败')
    }
  } catch (error) {
    console.error('Failed to sign contract:', error)
    ElMessage.error('合同签署失败')
  } finally {
    signing.value = false
  }
}

// 审核合同
const auditContract = () => {
  router.push({
    path: '/contract/audit',
    query: { id: contractUUID.value }
  })
}

// 下载文件
const downloadFile = (file) => {
  ElMessage.success(`开始下载文件：${file.name}`)
  // 实际实现可能涉及调用API下载文件
}

// 返回上一页
const goBack = () => {
  router.push('/contract')
}

// 获取合同类型标签类型
const getContractTypeTag = (type) => {
  const tagMap = {
    PURCHASE: 'success',
    MORTGAGE: 'primary',
    LEASE: 'info'
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
    NORMAL: 'primary',
    FROZEN: 'warning',
    COMPLETED: 'success',
    CANCEL: 'danger'
  }
  return tagMap[status] || ''
}

// 获取状态名称
const getStatusName = (status) => {
  const nameMap = {
    NORMAL: '正常',
    FROZEN: '冻结',
    COMPLETED: '已完成',
    CANCEL: '已取消'
  }
  return nameMap[status] || status
}

// 获取审计日志标题
const getAuditLogTitle = (action) => {
  const titleMap = {
    created: '创建合同',
    reviewed: '初审完成',
    approved: '审核通过',
    rejected: '审核拒绝',
    needRevision: '需要修改',
    revised: '已修改',
    signed: '签署合同',
    completed: '合同完成'
  }
  return titleMap[action] || action
}

// 获取审计日志类型
const getAuditLogType = (action) => {
  const typeMap = {
    created: 'info',
    reviewed: 'primary',
    approved: 'success',
    rejected: 'danger',
    needRevision: 'warning',
    revised: 'primary',
    signed: 'success',
    completed: 'success'
  }
  return typeMap[action] || 'info'
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

// 格式化日期时间
const formatDateTime = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 格式化货币
const formatCurrency = (value) => {
  if (!value && value !== 0) return '-'
  return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 打开修改状态对话框
const openUpdateStatusDialog = () => {
  updateStatusForm.value = {
    newStatus: '',
    reason: ''
  }
  updateStatusDialogVisible.value = true
}

// 确认修改状态
const confirmUpdateStatus = async () => {
  if (!updateStatusForm.value.newStatus) {
    ElMessage.warning('请选择新状态')
    return
  }
  
  if (!updateStatusForm.value.reason) {
    ElMessage.warning('请输入修改原因')
    return
  }
  
  updating.value = true
  
  try {
    await updateContractStatus({
      contractUUID: contractUUID.value,
      status: updateStatusForm.value.newStatus,
      reason: updateStatusForm.value.reason
    })
    
    ElMessage.success('合同状态修改成功')
    updateStatusDialogVisible.value = false
    
    // 刷新合同详情
    await fetchContractDetail()
  } catch (error) {
    console.error('Failed to update contract status:', error)
    ElMessage.error(error.response?.data?.message || '修改合同状态失败')
  } finally {
    updating.value = false
  }
}

// 页面加载时获取合同详情
onMounted(async () => {
  await fetchContractDetail()
})
</script>

<style scoped>
.contract-detail-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 15px;
}

.title-section h2 {
  margin: 0;
  font-size: 22px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.loading-container {
  padding: 20px 0;
}

.contract-tabs {
  margin-bottom: 20px;
}

.contract-section {
  margin-bottom: 20px;
}

.contract-section h3 {
  margin: 15px 0 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
  font-size: 18px;
}

.contract-clause {
  margin: 10px 0 10px 20px;
}

.contract-clause h4 {
  margin: 8px 0;
  font-size: 16px;
  font-weight: 500;
}

.log-details {
  margin: 5px 0;
  padding-left: 10px;
  border-left: 2px solid #EBEEF5;
}

.log-meta {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
}

.contract-text {
  white-space: pre-wrap;
  line-height: 1.6;
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

.update-status-form {
  padding: 10px 20px;
}
</style>
