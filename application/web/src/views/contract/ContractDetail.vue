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
        <el-button @click="goBack">返回</el-button>
        <el-button 
          v-if="canSign"
          type="success" 
          @click="showSignDialog"
        >
          签署合同
        </el-button>
        <el-button 
          v-if="canAudit"
          type="warning" 
          @click="auditContract"
        >
          审核合同
        </el-button>
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
              <el-descriptions-item label="关联交易">{{ contract.transactionId }}</el-descriptions-item>
              <el-descriptions-item label="合同状态">
                <el-tag :type="getStatusTag(contract.status)">
                  {{ getStatusName(contract.status) }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
            
            <el-divider content-position="left">交易信息</el-divider>
            
            <el-descriptions :column="2" border>
              <el-descriptions-item label="买方">{{ contract.parties?.buyer }}</el-descriptions-item>
              <el-descriptions-item label="卖方">{{ contract.parties?.seller }}</el-descriptions-item>
              <el-descriptions-item label="合同金额">¥{{ formatCurrency(contract.amount) }}</el-descriptions-item>
              <el-descriptions-item label="支付方式">{{ contract.paymentMethod || '未指定' }}</el-descriptions-item>
              <el-descriptions-item label="签署日期" :span="1">{{ formatDate(contract.signDate) || '未签署' }}</el-descriptions-item>
              <el-descriptions-item label="生效日期" :span="1">{{ formatDate(contract.effectiveDate) || '未生效' }}</el-descriptions-item>
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
        
<!--        <el-tab-pane label="签署记录">-->
<!--          <el-card>-->
<!--            <div v-if="contract.auditLogs && contract.auditLogs.length > 0">-->
<!--              <el-timeline>-->
<!--                <el-timeline-item-->
<!--                  v-for="(log, index) in contract.auditLogs"-->
<!--                  :key="index"-->
<!--                  :type="getAuditLogType(log.action)"-->
<!--                  :timestamp="formatDateTime(log.timestamp)"-->
<!--                >-->
<!--                  <h4>{{ getAuditLogTitle(log.action) }}</h4>-->
<!--                  <p>{{ log.comments }}</p>-->
<!--                  <p v-if="log.revisionRequirements" class="log-details">-->
<!--                    <strong>修改要求：</strong> {{ log.revisionRequirements }}-->
<!--                  </p>-->
<!--                  <p v-if="log.rejectionReason" class="log-details">-->
<!--                    <strong>拒绝理由：</strong> {{ log.rejectionReason }}-->
<!--                  </p>-->
<!--                  <div class="log-meta">-->
<!--                    <span>操作人：{{ log.auditor }}</span>-->
<!--                  </div>-->
<!--                </el-timeline-item>-->
<!--              </el-timeline>-->
<!--            </div>-->
<!--            <el-empty v-else description="暂无签署/审核记录" />-->
<!--          </el-card>-->
<!--        </el-tab-pane>-->
        
        <el-tab-pane label="附件">
          <el-card>
            <div v-if="contract.attachments && contract.attachments.length > 0">
              <el-table :data="contract.attachments" style="width: 100%">
                <el-table-column prop="name" label="文件名" min-width="200" />
                <el-table-column prop="type" label="类型" width="120" />
                <el-table-column prop="size" label="大小" width="120">
                  <template #default="scope">
                    {{ formatFileSize(scope.row.size) }}
                  </template>
                </el-table-column>
                <el-table-column prop="uploadTime" label="上传时间" width="180">
                  <template #default="scope">
                    {{ formatDateTime(scope.row.uploadTime) }}
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="120" fixed="right">
                  <template #default="scope">
                    <el-button type="primary" link @click="downloadFile(scope.row)">下载</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
            <el-empty v-else description="暂无合同附件" />
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
  </div>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {useUserStore} from '@/stores/user'
import {ElMessage} from 'element-plus'
import axios from 'axios'
import dayjs from 'dayjs'
import {getContractByUUID} from "@/api/contract.js";

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 加载状态
const loading = ref(true)
const signing = ref(false)

// 合同信息
const contract = ref(null)
const signDialogVisible = ref(false)

// 合同ID
const contractUUID = computed(() => route.params.id)

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

// 判断是否可以审核合同
const canAudit = computed(() => {
  return userStore.hasOrganization('audit') && 
         contract.value && 
         contract.value.status === 'signed'
})

// 获取合同详情
const fetchContractDetail = async () => {
  try {
    loading.value = true
    
    // 调用API获取合同详情
    contract.value = await getContractByUUID(contractUUID.value)
  } catch (error) {
    console.error('Failed to fetch contract details:', error)
    ElMessage.error('获取合同详情失败')
  } finally {
    loading.value = false
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
    purchase: 'success',
    mortgage: 'primary',
    lease: 'info'
  }
  return tagMap[type] || ''
}

// 获取合同类型名称
const getContractTypeName = (type) => {
  const nameMap = {
    purchase: '购房合同',
    mortgage: '贷款合同',
    lease: '租赁合同'
  }
  return nameMap[type] || type
}

// 获取状态标签类型
const getStatusTag = (status) => {
  const tagMap = {
    NORMAL: 'primary',
    FROZEN: 'warning',
    COMPLETED: 'success'
  }
  return tagMap[status] || ''
}

// 获取状态名称
const getStatusName = (status) => {
  const nameMap = {
    NORMAL: '正常',
    FROZEN: '冻结',
    COMPLETED: '已完成',
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

// 页面加载时获取合同详情
onMounted(() => {
  if (contractUUID.value) {
    fetchContractDetail()
  } else {
    loading.value = false
    ElMessage.error('合同ID不能为空')
  }
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
</style>
