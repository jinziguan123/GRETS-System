<template>
  <div class="contract-audit-container">
    <div class="page-header">
      <h2>合同审计管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="refreshContracts">刷新数据</el-button>
      </div>
    </div>
    
    <!-- 筛选器 -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="contract-filter" size="small">
        <el-form-item label="合同状态">
          <el-select v-model="filterForm.status" placeholder="选择合同状态" clearable>
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
            <el-option label="需修改" value="needRevision" />
          </el-select>
        </el-form-item>
        <el-form-item label="合同类型">
          <el-select v-model="filterForm.contractType" placeholder="选择合同类型" clearable>
            <el-option label="购房合同" value="purchase" />
            <el-option label="贷款合同" value="mortgage" />
            <el-option label="租赁合同" value="lease" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="filterForm.keyword" placeholder="合同编号/内容关键词" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchContracts">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
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
          @row-click="handleRowClick" 
          row-key="id"
          v-loading="tableLoading"
        >
          <el-table-column prop="id" label="合同编号" width="160" />
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
          <el-table-column prop="title" label="合同标题" min-width="200" />
          <el-table-column prop="parties" label="合同方" width="260">
            <template #default="scope">
              <div>买方: {{ scope.row.parties.buyer.replace(/(?<=.).(?=.)/g, '*') }}</div>
              <div>卖方: {{ scope.row.parties.seller.replace(/(?<=.).(?=.)/g, '*') }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="合同金额" width="150">
            <template #default="scope">
              ¥{{ formatCurrency(scope.row.amount) }}
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建日期" width="120">
            <template #default="scope">
              {{ formatDate(scope.row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getStatusTag(scope.row.status)">
                {{ getStatusName(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="scope">
              <el-button 
                v-if="scope.row.status === 'pending'" 
                type="primary" 
                link 
                @click.stop="openAuditDialog(scope.row)"
              >
                审核
              </el-button>
              <el-button 
                type="info" 
                link 
                @click.stop="viewContractDetail(scope.row.id)"
              >
                查看
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
    
    <!-- 合同审核对话框 -->
    <el-dialog
      v-model="auditDialogVisible"
      title="合同审核"
      width="700px"
      destroy-on-close
    >
      <div v-if="currentContract" class="audit-dialog-content">
        <div class="contract-info">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="合同编号">{{ currentContract.id }}</el-descriptions-item>
            <el-descriptions-item label="合同类型">{{ getContractTypeName(currentContract.contractType) }}</el-descriptions-item>
            <el-descriptions-item label="合同标题">{{ currentContract.title }}</el-descriptions-item>
            <el-descriptions-item label="创建日期">{{ formatDate(currentContract.createTime) }}</el-descriptions-item>
            <el-descriptions-item label="买方">{{ currentContract.parties.buyer }}</el-descriptions-item>
            <el-descriptions-item label="卖方">{{ currentContract.parties.seller }}</el-descriptions-item>
            <el-descriptions-item label="合同金额">¥{{ formatCurrency(currentContract.amount) }}</el-descriptions-item>
            <el-descriptions-item label="关联资产">{{ currentContract.property.title }}</el-descriptions-item>
          </el-descriptions>
        </div>
        
        <div class="contract-highlights">
          <h4>合同要点</h4>
          <el-collapse>
            <el-collapse-item title="主要条款" name="main">
              <div class="highlight-content">
                <p v-for="(point, index) in currentContract.highlights.mainPoints" :key="`main-${index}`">
                  {{ index + 1 }}. {{ point }}
                </p>
              </div>
            </el-collapse-item>
            <el-collapse-item title="风险点" name="risks">
              <div class="highlight-content">
                <p v-for="(risk, index) in currentContract.highlights.risks" :key="`risk-${index}`" class="risk-point">
                  <el-tag type="danger" size="small">风险</el-tag> {{ risk }}
                </p>
              </div>
            </el-collapse-item>
            <el-collapse-item title="异常情况" name="anomalies">
              <div class="highlight-content">
                <p v-for="(anomaly, index) in currentContract.highlights.anomalies" :key="`anomaly-${index}`" class="anomaly-point">
                  <el-tag type="warning" size="small">异常</el-tag> {{ anomaly }}
                </p>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
        
        <el-divider />
        
        <div class="audit-form">
          <el-form ref="auditFormRef" :model="auditForm" :rules="auditRules" label-width="100px">
            <el-form-item label="审核结果" prop="result">
              <el-radio-group v-model="auditForm.result">
                <el-radio label="approved">通过</el-radio>
                <el-radio label="needRevision">需修改</el-radio>
                <el-radio label="rejected">拒绝</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="审核意见" prop="comments">
              <el-input
                v-model="auditForm.comments"
                type="textarea"
                :rows="4"
                placeholder="请输入审核意见..."
              />
            </el-form-item>
            
            <el-form-item 
              v-if="auditForm.result === 'needRevision'" 
              label="修改要求" 
              prop="revisionRequirements"
            >
              <el-input
                v-model="auditForm.revisionRequirements"
                type="textarea"
                :rows="4"
                placeholder="请详细说明需要修改的内容..."
              />
            </el-form-item>
            
            <el-form-item 
              v-if="auditForm.result === 'rejected'" 
              label="拒绝理由" 
              prop="rejectionReason"
            >
              <el-input
                v-model="auditForm.rejectionReason"
                type="textarea"
                :rows="4"
                placeholder="请详细说明拒绝的理由..."
              />
            </el-form-item>
          </el-form>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="auditDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitAudit" :loading="submitting">
            提交审核
          </el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 合同详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="合同详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentContract" class="contract-detail">
        <el-descriptions title="基本信息" :column="2" border>
          <el-descriptions-item label="合同编号">{{ currentContract.id }}</el-descriptions-item>
          <el-descriptions-item label="合同类型">{{ getContractTypeName(currentContract.contractType) }}</el-descriptions-item>
          <el-descriptions-item label="合同标题">{{ currentContract.title }}</el-descriptions-item>
          <el-descriptions-item label="签署日期">{{ formatDate(currentContract.signDate || currentContract.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="生效日期">{{ formatDate(currentContract.effectiveDate || currentContract.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTag(currentContract.status)">
              {{ getStatusName(currentContract.status) }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <el-divider content-position="left">交易信息</el-divider>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="买方" :span="1">{{ currentContract.parties.buyer }}</el-descriptions-item>
          <el-descriptions-item label="卖方" :span="1">{{ currentContract.parties.seller }}</el-descriptions-item>
          <el-descriptions-item label="合同金额" :span="1">¥{{ formatCurrency(currentContract.amount) }}</el-descriptions-item>
          <el-descriptions-item label="支付方式" :span="1">{{ currentContract.paymentMethod }}</el-descriptions-item>
          <el-descriptions-item label="交付日期" :span="1">{{ formatDate(currentContract.property.deliveryDate) }}</el-descriptions-item>
          <el-descriptions-item label="关联资产" :span="1">{{ currentContract.property.title }}</el-descriptions-item>
        </el-descriptions>
        
        <el-divider content-position="left">合同文本预览</el-divider>
        
        <div class="contract-document">
          <div class="document-preview">
            <el-tabs type="border-card">
              <el-tab-pane label="合同条款">
                <div class="contract-clauses">
                  <div v-for="(section, sIndex) in currentContract.content.sections" :key="`section-${sIndex}`" class="contract-section">
                    <h3>{{ section.title }}</h3>
                    <div v-for="(clause, cIndex) in section.clauses" :key="`clause-${sIndex}-${cIndex}`" class="contract-clause">
                      <h4>{{ clause.title }}</h4>
                      <p>{{ clause.content }}</p>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
              <el-tab-pane label="审计日志">
                <div class="audit-logs">
                  <el-timeline>
                    <el-timeline-item
                      v-for="(log, index) in currentContract.auditLogs"
                      :key="index"
                      :type="getAuditLogType(log.action)"
                      :timestamp="formatDateTime(log.timestamp)"
                    >
                      <h4>{{ getAuditLogTitle(log.action) }}</h4>
                      <p>{{ log.comments }}</p>
                      <p v-if="log.revisionRequirements" class="log-details">
                        <strong>修改要求：</strong> {{ log.revisionRequirements }}
                      </p>
                      <p v-if="log.rejectionReason" class="log-details">
                        <strong>拒绝理由：</strong> {{ log.rejectionReason }}
                      </p>
                      <div class="log-meta">
                        <span>操作人：{{ log.auditor }}</span>
                      </div>
                    </el-timeline-item>
                  </el-timeline>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import axios from 'axios'

// 加载状态
const loading = ref(false)
const tableLoading = ref(false)
const submitting = ref(false)

// 对话框状态
const auditDialogVisible = ref(false)
const detailDialogVisible = ref(false)

// 当前选中的合同
const currentContract = ref(null)

// 筛选表单
const filterForm = reactive({
  status: '',
  contractType: '',
  dateRange: [],
  keyword: ''
})

// 审核表单
const auditForm = reactive({
  result: '',
  comments: '',
  revisionRequirements: '',
  rejectionReason: ''
})

// 审核表单规则
const auditRules = {
  result: [
    { required: true, message: '请选择审核结果', trigger: 'change' }
  ],
  comments: [
    { required: true, message: '请输入审核意见', trigger: 'blur' },
    { min: 5, message: '审核意见不能少于5个字符', trigger: 'blur' }
  ],
  revisionRequirements: [
    { 
      required: true, 
      message: '请输入修改要求', 
      trigger: 'blur',
      validator: (rule, value, callback) => {
        if (auditForm.result === 'needRevision' && (!value || value.trim() === '')) {
          callback(new Error('请输入修改要求'))
        } else {
          callback()
        }
      } 
    }
  ],
  rejectionReason: [
    { 
      required: true, 
      message: '请输入拒绝理由', 
      trigger: 'blur',
      validator: (rule, value, callback) => {
        if (auditForm.result === 'rejected' && (!value || value.trim() === '')) {
          callback(new Error('请输入拒绝理由'))
        } else {
          callback()
        }
      } 
    }
  ]
}

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

// 合同数据
const contractsData = ref([])

// 重置筛选条件
const resetFilter = () => {
  filterForm.status = ''
  filterForm.contractType = ''
  filterForm.dateRange = []
  filterForm.keyword = ''
  searchContracts()
}

// 搜索合同
const searchContracts = () => {
  tableLoading.value = true
  
  // 这里实现实际的API调用
  setTimeout(() => {
    fetchContractsData()
    tableLoading.value = false
  }, 500)
}

// 刷新合同列表
const refreshContracts = () => {
  tableLoading.value = true
  
  // 这里实现实际的API调用
  setTimeout(() => {
    fetchContractsData()
    ElMessage.success('数据已刷新')
    tableLoading.value = false
  }, 500)
}

// 获取合同数据（连接实际API）
const fetchContractsData = async () => {
  try {
    // 构建查询参数
    const params = {
      status: filterForm.status,
      contractType: filterForm.contractType,
      pageSize: pageSize.value,
      pageNumber: currentPage.value,
    }
    
    if (filterForm.keyword) {
      params.keyword = filterForm.keyword
    }
    
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.startDate = filterForm.dateRange[0]
      params.endDate = filterForm.dateRange[1]
    }
    
    // 调用API
    const { data } = await axios.get('/contracts', { params })
    
    if (data.code === 200) {
      contractsData.value = data.data.items || []
      totalItems.value = data.data.total || 0
    } else {
      ElMessage.error(data.message || '获取合同列表失败')
    }
  } catch (error) {
    console.error('Failed to fetch contracts:', error)
    ElMessage.error('获取合同列表失败')
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
  viewContractDetail(row.id)
}

// 查看合同详情
const viewContractDetail = async (id) => {
  try {
    // 调用API获取合同详情
    const { data } = await axios.get(`/contracts/${id}`)
    
    if (data.code === 200) {
      currentContract.value = data.data
      detailDialogVisible.value = true
    } else {
      ElMessage.error(data.message || '获取合同详情失败')
    }
  } catch (error) {
    console.error('Failed to fetch contract details:', error)
    ElMessage.error('获取合同详情失败')
  }
}

// 打开审核对话框
const openAuditDialog = async (contract) => {
  try {
    // 调用API获取合同详情
    const { data } = await axios.get(`/contracts/${contract.id}`)
    
    if (data.code === 200) {
      currentContract.value = data.data
      
      // 重置审核表单
      auditForm.result = ''
      auditForm.comments = ''
      auditForm.revisionRequirements = ''
      auditForm.rejectionReason = ''
      
      auditDialogVisible.value = true
    } else {
      ElMessage.error(data.message || '获取合同详情失败')
    }
  } catch (error) {
    console.error('Failed to fetch contract details:', error)
    ElMessage.error('获取合同详情失败')
  }
}

// 提交审核
const submitAudit = async () => {
  if (!auditForm.result) {
    ElMessage.warning('请选择审核结果')
    return
  }
  
  try {
    submitting.value = true
    
    // 准备提交数据
    const auditData = {
      result: auditForm.result,
      comments: auditForm.comments,
      revisionRequirements: auditForm.revisionRequirements || '',
      rejectionReason: auditForm.rejectionReason || ''
    }
    
    // 调用审核API
    const { data } = await axios.post(`/contracts/${currentContract.value.id}/audit`, auditData)
    
    if (data.code === 200) {
      ElMessage.success('审核提交成功')
      auditDialogVisible.value = false
      
      // 更新表格数据
      const index = contractsData.value.findIndex(item => item.id === currentContract.value.id)
      if (index !== -1) {
        contractsData.value[index].status = auditForm.result
      }
      
      // 刷新数据
      setTimeout(() => {
        fetchContractsData()
      }, 500)
    } else {
      ElMessage.error(data.message || '审核提交失败')
    }
    
    submitting.value = false
  } catch (error) {
    console.error('审核提交失败:', error)
    ElMessage.error('审核提交失败，请重试')
    submitting.value = false
  }
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
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    needRevision: 'info'
  }
  return tagMap[status] || ''
}

// 获取状态名称
const getStatusName = (status) => {
  const nameMap = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    needRevision: '需修改'
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
    revised: '已修改'
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
    revised: 'primary'
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

// 页面加载时获取合同数据
onMounted(() => {
  loading.value = true
  // 模拟API调用
  setTimeout(() => {
    fetchContractsData()
    loading.value = false
  }, 1000)
})
</script>

<style scoped>
.contract-audit-container {
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

.audit-dialog-content {
  max-height: 60vh;
  overflow-y: auto;
}

.contract-info {
  margin-bottom: 20px;
}

.contract-highlights {
  margin-top: 20px;
}

.highlight-content {
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.risk-point {
  margin: 5px 0;
  color: #F56C6C;
}

.anomaly-point {
  margin: 5px 0;
  color: #E6A23C;
}

.contract-document {
  margin-top: 20px;
}

.document-preview {
  border: 1px solid #EBEEF5;
  border-radius: 4px;
}

.contract-section {
  margin-bottom: 20px;
}

.contract-clause {
  margin: 10px 0;
  padding-left: 15px;
}

.contract-clause h4 {
  margin: 8px 0;
  font-size: 14px;
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

.audit-logs {
  padding: 10px;
}

.property-meta {
  font-size: 12px;
  color: #909399;
}
</style> 