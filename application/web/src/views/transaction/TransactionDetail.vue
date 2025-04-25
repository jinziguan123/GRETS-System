<template>
  <div class="transaction-detail-container">
    <!-- 加载中状态 -->
    <el-card v-if="loading" class="box-card">
      <div class="loading-container">
        <el-skeleton :rows="10" animated/>
      </div>
    </el-card>

    <!-- 只有交易参与方、政府和审计可见的内容 -->
    <el-card v-else-if="isAuthorized" class="box-card">
      <template #header>
        <div class="card-header">
          <h3>交易详情</h3>
          <div class="header-buttons">
            <el-button
                v-if="canAcceptTransaction"
                type="success"
                @click="handleAccept"
            >同意交易
            </el-button>
            <el-button
                v-if="canCompleteTransaction"
                type="primary"
                @click="handleComplete"
            >完成交易
            </el-button>
            <el-button
                v-if="canCancelTransaction"
                type="warning"
                @click="handleCancel"
            >取消交易
            </el-button>
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>

      <div v-if="transactionInfo">
        <!-- 支付进度 -->
        <div class="progress-section">
          <h4>支付进度</h4>
          
          <!-- 房款支付进度 -->
          <div class="progress-item">
            <div class="progress-title">房款支付进度</div>
            <div class="progress-info">
              <div class="progress-text">
                <span>已支付: {{ formatPrice(totalTransferAmount) }}</span>
                <span>总金额: {{ formatPrice(transactionInfo.price) }}</span>
              </div>
              <el-progress
                  :percentage="transferPercentage"
                  :format="percentFormat"
                  :status="transferPercentage >= 100 ? 'success' : ''"
                  :stroke-width="20"
              ></el-progress>
            </div>
          </div>
          
          <!-- 税费支付进度 -->
          <div class="progress-item">
            <div class="progress-title">税费支付进度</div>
            <div class="progress-info">
              <div class="progress-text">
                <span>已支付: {{ formatPrice(totalTaxAmount) }}</span>
                <span>总金额: {{ formatPrice(transactionInfo.tax) }}</span>
              </div>
              <el-progress
                  :percentage="taxPercentage"
                  :format="percentFormat"
                  :status="taxPercentage >= 100 ? 'success' : ''"
                  :stroke-width="20"
              ></el-progress>
            </div>
          </div>
          
          <!-- 总支付进度 -->
          <div class="progress-item total-progress">
            <div class="progress-title">总支付进度</div>
            <div class="progress-info">
              <div class="progress-text">
                <span>已支付: {{ formatPrice(totalPaidAmount) }}</span>
                <span>总金额: {{ formatPrice(transactionInfo.price + transactionInfo.tax) }}</span>
              </div>
              <el-progress
                  :percentage="paymentPercentage"
                  :format="percentFormat"
                  :status="paymentPercentage >= 100 ? 'success' : ''"
                  :stroke-width="20"
              ></el-progress>
            </div>
          </div>
        </div>

        <!-- 交易基本信息 -->
        <el-descriptions title="交易基本信息" :column="2" border>
          <el-descriptions-item label-width="auto" label="交易编号">{{ transactionInfo.transactionUUID }}</el-descriptions-item>
          <el-descriptions-item label-width="auto" label="交易状态">
            <el-tag :type="getStatusType(transactionInfo.status)">
              {{ getStatusText(transactionInfo.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="房产证号">{{ transactionInfo.realtyCertHash }}</el-descriptions-item>
          <el-descriptions-item label="交易价格">{{ formatPrice(transactionInfo.price) }}</el-descriptions-item>
          <el-descriptions-item label="卖方身份证">{{ transactionInfo.sellerCitizenIDHash }}</el-descriptions-item>
          <el-descriptions-item label="买方身份证">{{ transactionInfo.buyerCitizenIDHash }}</el-descriptions-item>
          <el-descriptions-item label="卖方组织">{{ getCurrentOwnerOrganization(transactionInfo.sellerOrganization) }}</el-descriptions-item>
          <el-descriptions-item label="买方组织">{{ getCurrentOwnerOrganization(transactionInfo.buyerOrganization) }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(transactionInfo.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(transactionInfo.updateTime) }}</el-descriptions-item>
          <el-descriptions-item label="预计完成时间">
            {{ formatDate(transactionInfo.completedTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="税费">{{ formatPrice(transactionInfo.tax) }}</el-descriptions-item>
        </el-descriptions>

        <!-- 房产信息 -->
        <el-descriptions v-if="realtyInfo" title="房产信息" :column="3" border class="section-mt">
          <el-descriptions-item label="房产地址">{{ generateAddress(realtyInfo) }}</el-descriptions-item>
          <el-descriptions-item label="房产类型">{{ getRealtyTypeText(realtyInfo.realtyType) }}</el-descriptions-item>
          <el-descriptions-item label="建筑面积">{{ realtyInfo.area }} 平方米</el-descriptions-item>
          <el-descriptions-item label="房产状态">{{ getRealtyStatusText(realtyInfo.status) }}</el-descriptions-item>
          <el-descriptions-item label="参考价格">{{ formatPrice(realtyInfo.price) }}</el-descriptions-item>
          <el-descriptions-item>
            <template #label>
              <el-button type="primary" size="small" @click="viewRealtyDetail">查看房产详情</el-button>
            </template>
          </el-descriptions-item>
        </el-descriptions>

        <!-- 合同信息 -->
        <div v-if="contractInfo" class="section-mt">
          <h4>合同信息</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="合同编号">{{ contractInfo.contractUUID }}</el-descriptions-item>
            <el-descriptions-item label="合同类型">{{
                getContractTypeText(contractInfo.contractType)
              }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(contractInfo.createTime) }}</el-descriptions-item>
            <el-descriptions-item label="创建人">{{ contractInfo.creatorCitizenID }}</el-descriptions-item>
            <el-descriptions-item label="合同标题">{{ contractInfo.title }}</el-descriptions-item>
            <el-descriptions-item label="合同状态">
              <el-tag :type="getContractStatusType(contractInfo.status)">
                {{ getContractStatusText(contractInfo.status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
          <div class="contract-content section-mt">
            <h5>合同内容</h5>
            <div class="content-box">{{ contractInfo.content }}</div>
          </div>
        </div>

        <!-- 支付记录 -->
        <div class="section-mt">
          <h4>支付记录</h4>
          <div v-if="payments.length > 0">
            <el-table :data="payments" border style="width: 100%">
              <el-table-column prop="paymentUUID" label="支付编号" width="auto"></el-table-column>
              <el-table-column prop="paymentType" label="支付类型" width="auto">
                <template #default="scope">
                  {{ getPaymentTypeText(scope.row.paymentType) }}
                </template>
              </el-table-column>
              <el-table-column prop="amount" label="支付金额" width="auto">
                <template #default="scope">
                  {{ formatPrice(scope.row.amount) }}
                </template>
              </el-table-column>
              <el-table-column prop="payerCitizenIDHash" label="付款方" width="auto"></el-table-column>
              <el-table-column prop="payerOrganization" label="付款方组织" width="auto">
                <template #default="scope">
                  {{ getCurrentOwnerOrganization(scope.row.payerOrganization) }}
                </template>
              </el-table-column>
              <el-table-column prop="receiverCitizenIDHash" label="收款方" width="auto"></el-table-column>
              <el-table-column prop="receiverOrganization" label="收款方组织" width="auto">
                <template #default="scope">
                  {{ getCurrentOwnerOrganization(scope.row.receiverOrganization) }}
                </template>
              </el-table-column>
              <el-table-column prop="createTime" label="支付时间" width="auto">
                <template #default="scope">
                  {{ formatDate(scope.row.createTime) }}
                </template>
              </el-table-column>
              <el-table-column prop="remarks" label="备注" width="auto">
                <template #default="scope">
                  {{ scope.row.remarks || '无' }}
                </template>
              </el-table-column>
            </el-table>
          </div>
          <el-empty v-else description="暂无支付记录"></el-empty>
        </div>

        <!-- 买方发起支付按钮或完成交易按钮 -->
        <div v-if="isBuyer && transactionInfo.status === 'IN_PROGRESS'" class="payment-action section-mt">
          <el-button
              v-if="paymentPercentage >= 100"
              type="success"
              size="large"
              @click="handleBuyerComplete"
          >完成交易
          </el-button>
          <el-button
              v-else
              type="primary"
              size="large"
              @click="showPaymentDialog"
          >发起支付
          </el-button>
        </div>
      </div>

      <el-empty v-else description="未找到交易信息"></el-empty>
    </el-card>

    <!-- 非授权用户提示 -->
    <el-card v-else class="unauthorized-card">
      <el-result
          icon="error"
          title="无权访问"
          sub-title="您不是此交易的参与方，无权查看交易详情"
      >
        <template #extra>
          <el-button type="primary" @click="goBack">返回上一页</el-button>
        </template>
      </el-result>
    </el-card>

    <!-- 支付组件对话框 -->
    <payment-dialog
        v-model:visible="paymentDialogVisible"
        :transaction="transactionInfo"
        :total-paid="totalPaidAmount"
        :total-paid-transfer="totalTransferAmount"
        :total-paid-tax="totalTaxAmount"
        @payment-success="handlePaymentSuccess"
    />
  </div>
</template>

<script setup>
import {ref, reactive, computed, onMounted} from 'vue'
import {useRouter, useRoute} from 'vue-router'
import {ElMessage, ElMessageBox} from 'element-plus'
import axios from 'axios'
import {completeTransaction, getTransactionDetail, updateTransaction} from "@/api/transaction.js";
import PaymentDialog from '@/components/transaction/PaymentDialog.vue'
import CryptoJS from 'crypto-js'
import {getPaymentList} from "@/api/payment.js";
import {getRealtyDetail} from "@/api/realty.js";

const router = useRouter()
const route = useRoute()
const transactionUUID = ref(route.params.id)
const loading = ref(true)
const transactionInfo = ref(null)
const realtyInfo = ref(null)
const contractInfo = ref(null)
const audits = ref([])
const payments = ref([])
const submitLoading = ref(false)
const paymentDialogVisible = ref(false)

// 支付相关计算属性
const totalPaidAmount = computed(() => {
  if (!payments.value || payments.value.length === 0) return 0
  return payments.value.reduce((total, payment) => total + payment.amount, 0)
})

// 分别计算房款和税费的已支付金额
const totalTransferAmount = computed(() => {
  if (!payments.value || payments.value.length === 0) return 0
  return payments.value
    .filter(payment => payment.paymentType === 'TRANSFER')
    .reduce((total, payment) => total + payment.amount, 0)
})

const totalTaxAmount = computed(() => {
  if (!payments.value || payments.value.length === 0) return 0
  return payments.value
    .filter(payment => payment.paymentType === 'TAX')
    .reduce((total, payment) => total + payment.amount, 0)
})

const paymentPercentage = computed(() => {
  if (!transactionInfo.value || !transactionInfo.value.price || transactionInfo.value.price === 0) return 0
  const percentage = (totalPaidAmount.value / (transactionInfo.value.price + transactionInfo.value.tax)) * 100
  return Math.min(parseFloat(percentage.toFixed(2)), 100) // 最大不超过100%，保留两位小数
})

// 计算房款支付进度百分比
const transferPercentage = computed(() => {
  if (!transactionInfo.value || !transactionInfo.value.price || transactionInfo.value.price === 0) return 0
  const percentage = (totalTransferAmount.value / transactionInfo.value.price) * 100
  return Math.min(parseFloat(percentage.toFixed(2)), 100) // 最大不超过100%，保留两位小数
})

// 计算税费支付进度百分比
const taxPercentage = computed(() => {
  if (!transactionInfo.value || !transactionInfo.value.tax || transactionInfo.value.tax === 0) return 0
  const percentage = (totalTaxAmount.value / transactionInfo.value.tax) * 100
  return Math.min(parseFloat(percentage.toFixed(2)), 100) // 最大不超过100%，保留两位小数
})

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

// 获取用户信息
const userInfo = computed(() => {
  const userJson = localStorage.getItem('userInfo')
  if (!userJson) return null
  try {
    return JSON.parse(userJson)
  } catch (e) {
    return null
  }
})

// 判断当前用户是否为交易的买方
const isBuyer = computed(() => {
  if (!userInfo.value || !userInfo.value.citizenID || !transactionInfo.value) return false

  // 获取用户身份证号的SHA256哈希值
  const userCitizenIDHash = CryptoJS.SHA256(userInfo.value.citizenID).toString(CryptoJS.enc.Hex)

  return userCitizenIDHash === transactionInfo.value.buyerCitizenIDHash
})

// 判断当前用户是否为交易的卖方
const isSeller = computed(() => {
  if (!userInfo.value || !userInfo.value.citizenID || !transactionInfo.value) return false

  // 获取用户身份证号的SHA256哈希值
  const userCitizenIDHash = CryptoJS.SHA256(userInfo.value.citizenID).toString(CryptoJS.enc.Hex)

  return userCitizenIDHash === transactionInfo.value.sellerCitizenIDHash
})

// 判断当前用户是否为政府或审计
const isGovOrAudit = computed(() => {
  if (!userInfo.value || !userInfo.value.organization) return false
  return ['government', 'audit'].includes(userInfo.value.organization)
})

// 判断用户是否有权限查看此交易
const isAuthorized = computed(() => {
  return isBuyer.value || isSeller.value || isGovOrAudit.value;
})

// 判断是否可以审核交易
const canCheckTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false

  // 政府组织且交易状态为待审核
  return userInfo.value.organization === 'government' && transactionInfo.value.status === 'PENDING'
})

// 判断是否可以完成交易 - 修改为需要两种支付都完成
const canCompleteTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false

  // 政府组织且交易状态为已审核且房款和税费都已支付足够金额
  return userInfo.value.organization === 'government' &&
      transactionInfo.value.status === 'IN_PROGRESS' &&
      transferPercentage.value >= 100 &&
      taxPercentage.value >= 100
})

// 判断是否可以取消交易
const canCancelTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false

  // 任何角色在交易未完成前均可取消
  return ['PENDING', 'IN_PROGRESS'].includes(transactionInfo.value.status) &&
      (isBuyer.value || isSeller.value || userInfo.value.organization === 'government')
})

// 判断是否可以同意交易
const canAcceptTransaction = computed(() => {
  if (!userInfo.value || !transactionInfo.value) return false

  // 卖方且交易状态为待审核或进行中
  return isSeller.value && ['PENDING'].includes(transactionInfo.value.status)
})

// 获取交易详情
const fetchTransactionDetail = async () => {
  loading.value = true
  try {
    const response = await getTransactionDetail(transactionUUID.value)
    transactionInfo.value = response.transaction

    // 获取房产信息
    if (transactionInfo.value.realtyCertHash) {
      await fetchRealtyDetail(transactionInfo.value.realtyCertHash)
    }

    // 获取合同信息
    if (transactionInfo.value.contractUUID) {
      await fetchContractDetail(transactionInfo.value.contractUUID)
    }

    // 获取审核记录
    // await fetchAudits()

    // 获取支付记录
    await fetchPayments()

    // 检查登录状态和权限
    if (!isAuthorized.value && userInfo.value) {
      console.log("权限验证失败，但用户已登录")
    }

  } catch (error) {
    console.error('获取交易详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取交易详情失败')
  } finally {
    loading.value = false
  }
}

// 获取房产详情
const fetchRealtyDetail = async (realtyCertHash) => {
  try {
    const response = await getRealtyDetail(realtyCertHash)
    realtyInfo.value = response
  } catch (error) {
    console.error('获取房产详情失败:', error)
  }
}

// 生成地址信息
const generateAddress = (item) => {
  return `${item.province || ''}${item.province ? '省' : ''}${item.city || ''}${item.city ? '市' : ''}${item.district || ''}${item.district ? '区' : ''}${item.street || ''}${item.community || ''}${item.unit || ''}${item.unit ? '单元' : ''}${item.floor || ''}${item.floor ? '楼' : ''}${item.room || ''}${item.room ? '号' : ''}`
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

// 获取合同详情
const fetchContractDetail = async (contractUUID) => {
  try {
    const response = await axios.get(`/api/contract/getContractByUUID/${contractUUID}`)
    contractInfo.value = response.data
  } catch (error) {
    console.error('获取合同详情失败:', error)
    // 使用mock数据
    contractInfo.value = mockContractDetail(contractUUID)
  }
}

// 获取审核记录
const fetchAudits = async () => {
  try {
    const response = await axios.get(`/api/audit/transaction/${transactionUUID.value}`)
    audits.value = response.data || []
  } catch (error) {
    console.error('获取审核记录失败:', error)
    // 使用mock数据
    audits.value = mockAudits(transactionUUID.value)
  }
}

// 获取支付记录
const fetchPayments = async () => {
  try {
    const response = await getPaymentList({
      transactionUUID: transactionUUID.value,
      pageNumber: 1,
      pageSize: 1000
    })
    payments.value = response.paymentList || []
  } catch (error) {
    console.error('获取支付记录失败:', error)
    // 使用mock数据
    payments.value = mockPayments(transactionUUID.value)
  }
}

// 处理支付成功
const handlePaymentSuccess = async () => {
  ElMessage.success('支付成功')
  await fetchPayments() // 刷新支付记录
}

// 同意交易
const handleAccept = async () => {
  ElMessageBox.confirm('确定要同意此次交易吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'success'
  }).then(async () => {
    loading.value = true

    try {
      await updateTransaction({
        transactionUUID: transactionUUID.value,
        status: 'IN_PROGRESS',
      })
      ElMessage.success('已同意交易')
      await fetchTransactionDetail()
    } catch (error) {
      console.error('同意交易失败:', error)
      ElMessage.error(error.response?.data?.message || '同意交易失败')
    } finally {
      loading.value = false
    }
  }).catch(() => {
  })
}

function refresh(){
  location.reload()
}

// 完成交易
const handleComplete = async () => {
  ElMessageBox.confirm('确定要完成此次交易吗？完成后将不可撤销', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value = true

    try {
      await completeTransaction({transactionUUID: transactionUUID.value})
      ElMessage.success('交易已完成')
      await fetchTransactionDetail()
    } catch (error) {
      console.error('完成交易失败:', error)
      ElMessage.error(error.response?.data?.message || '完成交易失败')
    } finally {
      loading.value = false
    }
  }).catch(() => {
  })
}

// 取消交易
const handleCancel = async () => {
  ElMessageBox.confirm('确定要取消此次交易吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    loading.value = true

    try {
      await updateTransaction({
        transactionUUID: transactionUUID.value,
        status: 'REJECTED',
      })
      ElMessage.success('交易已取消')
      await fetchTransactionDetail()
    } catch (error) {
      console.error('取消交易失败:', error)
      ElMessage.error(error.response?.data?.message || '取消交易失败')
    } finally {
      loading.value = false
    }
  }).catch(() => {
  })
}

// 查看房产详情
const viewRealtyDetail = () => {
  if (realtyInfo.value) {
    router.push(`/realty/${realtyInfo.value.realtyCertHash}`)
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 显示支付对话框
const showPaymentDialog = () => {
  paymentDialogVisible.value = true
}

// Mock函数：获取合同详情
const mockContractDetail = (contractUUID) => {
  return {
    contractUUID: contractUUID,
    title: '房产买卖合同',
    content: '这是一份房产买卖合同的内容，包含买卖双方的权利和义务...',
    contractType: 'PURCHASE',
    status: 'SIGNED',
    creatorCitizenID: '310110200312345678',
    createTime: new Date().toISOString(),
    updateTime: new Date().toISOString()
  }
}

// Mock函数：获取审核记录
const mockAudits = (transactionUUID) => {
  return [
    {
      id: 1,
      transactionUUID: transactionUUID,
      auditorName: '张三',
      auditorOrg: 'GOVERNMENT',
      result: 'APPROVED',
      opinion: '审核通过，符合交易要求',
      createTime: new Date(Date.now() - 3600000).toISOString()
    }
  ]
}

// Mock函数：获取支付记录
const mockPayments = (transactionUUID) => {
  return [
    {
      paymentID: 'PAY001',
      transactionUUID: transactionUUID,
      paymentType: 'PAYMENT',
      amount: 100000,
      payerCitizenID: '310110200306241518',
      payeeCitizenID: '310110200306241519',
      createTime: new Date(Date.now() - 7200000).toISOString()
    },
    {
      paymentID: 'PAY002',
      transactionUUID: transactionUUID,
      paymentType: 'TAX',
      amount: 5000,
      payerCitizenID: '310110200306241518',
      payeeCitizenID: 'GOVERNMENT',
      createTime: new Date(Date.now() - 3600000).toISOString()
    }
  ]
}

// 价格格式化
const formatPrice = (price) => {
  return `¥ ${parseFloat(price).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })}`
}

// 日期格式化
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取交易状态标签类型
const getStatusType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'REJECTED': 'danger',
    'IN_PROGRESS': 'warning',
    'COMPLETED': 'success',
  }
  return statusMap[status] || 'info'
}

// 获取交易状态文本
const getStatusText = (status) => {
  const statusMap = {
    'PENDING': '待审核',
    'REJECTED': '已拒绝',
    'IN_PROGRESS': '进行中',
    'COMPLETED': '已完成',
  }
  return statusMap[status] || status
}

// 获取房产状态文本
const getRealtyStatusText = (status) => {
  const statusMap = {
    'NORMAL': '正常',
    'IN_SALE': '交易中',
    'IN_MORTGAGE': '已抵押',
    'FROZEN': '已冻结',
    'PENDING_SALE': '挂牌'
  }
  return statusMap[status] || status
}

// 获取合同状态标签类型
const getContractStatusType = (status) => {
  const statusMap = {
    'PENDING': 'info',
    'SIGNED': 'success',
    'REJECTED': 'danger',
    'EXPIRED': 'warning'
  }
  return statusMap[status] || 'info'
}

// 获取合同状态文本
const getContractStatusText = (status) => {
  const statusMap = {
    'PENDING': '待签署',
    'SIGNED': '已签署',
    'REJECTED': '已拒绝',
    'EXPIRED': '已过期'
  }
  return statusMap[status] || status
}

// 获取支付类型文本
const getPaymentTypeText = (type) => {
  const typeMap = {
    'TRANSFER': '房款',
    'TAX': '税费',
  }
  return typeMap[type] || type
}

// 获取合同类型文本
const getContractTypeText = (type) => {
  const typeMap = {
    'PURCHASE': '购房合同',
    'MORTGAGE': '抵押合同',
    'LEASE': '租赁合同'
  }
  return typeMap[type] || type
}

// 买方完成交易按钮判断逻辑也需要修改
const handleBuyerComplete = async () => {
  // 检查两种支付是否都已完成
  if (transferPercentage.value < 100 || taxPercentage.value < 100) {
    ElMessage.warning('请先完成所有房款和税费的支付')
    return
  }
  
  ElMessageBox.confirm('交易金额已支付完成，确定要完成此次交易吗？完成后将不可撤销', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'success'
  }).then(async () => {
    loading.value = true

    try {
      await completeTransaction({transactionUUID: transactionUUID.value});
      ElMessage.success('交易已完成');
      await fetchTransactionDetail();
    } catch (error) {
      console.error('完成交易失败:', error);
      ElMessage.error(error.response?.data?.message || '完成交易失败');
    } finally {
      loading.value = false;
    }
  }).catch(() => {
  });
}

onMounted(() => {
  fetchTransactionDetail()
})
</script>

<style scoped>
.transaction-detail-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-mt {
  margin-top: 30px;
}

.progress-section {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.progress-item {
  margin-bottom: 20px;
}

.progress-item:last-child {
  margin-bottom: 0;
}

.progress-title {
  font-weight: bold;
  margin-bottom: 10px;
  color: #606266;
}

.total-progress {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px dashed #d9d9d9;
}

.progress-info {
  margin-top: 10px;
}

.progress-text {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-weight: bold;
}

.contract-content {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
}

.content-box {
  background-color: white;
  padding: 15px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  min-height: 150px;
  white-space: pre-wrap;
}

.payment-action {
  display: flex;
  justify-content: center;
  margin: 30px 0;
}

.unauthorized-card {
  max-width: 800px;
  margin: 100px auto;
}

.header-buttons {
  display: flex;
  gap: 10px;
}
</style>
