<template>
  <div class="contract-create-container">
    <div class="page-header">
      <h2>创建合同</h2>
    </div>
    
    <el-card class="contract-form-card">
      <el-form 
        ref="contractFormRef" 
        :model="contractForm" 
        :rules="contractRules" 
        label-width="120px"
        v-loading="loading"
      >
        <!-- 合同基本信息 -->
        <el-divider content-position="center">基本信息</el-divider>
        <el-form-item label="合同类型" prop="contractType">
          <el-select v-model="contractForm.contractType" placeholder="请选择合同类型">
            <el-option label="购房合同" value="PURCHASE" />
            <el-option label="贷款合同" value="MORTGAGE" />
            <el-option label="租赁合同" value="LEASE" />
          </el-select>
        </el-form-item>
        <el-form-item label="合同标题" prop="title">
          <el-input v-model="contractForm.title" placeholder="请输入合同标题" />
        </el-form-item>
        
        <el-form-item label="合同模板">
          <el-select v-model="selectedTemplate" placeholder="请选择合同模板" @change="selectTemplate">
            <el-option 
              v-for="item in templateOptions" 
              :key="item.id" 
              :label="item.name" 
              :value="item.id" 
            />
          </el-select>
          <div class="form-tip">选择模板将自动填充合同标题和内容</div>
        </el-form-item>
        
        <el-form-item label="合同内容" prop="content">
          <el-input 
            v-model="contractForm.content" 
            type="textarea" 
            :rows="15" 
            placeholder="请输入合同内容"
          />
        </el-form-item>
        
        <!-- 提交按钮 -->
        <el-form-item>
          <el-button type="primary" @click="submitContract" :loading="submitting">创建合同</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import { createContract } from '@/api/contract'

const router = useRouter()
const userStore = useUserStore()
const contractFormRef = ref(null)

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 合同表单数据
const contractForm = reactive({
  contractType: 'PURCHASE',
  title: '',
  content: '',
  creatorCitizenID: userStore.user.citizenID || ''
})

// 选中的模板
const selectedTemplate = ref('')

// 表单验证规则
const contractRules = {
  contractType: [
    { required: true, message: '请选择合同类型', trigger: 'change' }
  ],
  title: [
    { required: true, message: '请输入合同标题', trigger: 'blur' },
    { min: 3, max: 100, message: '长度在 3 到 100 个字符', trigger: 'blur' }
  ],
  transactionId: [
    { required: true, message: '请选择关联交易', trigger: 'change' }
  ],
  buyerInfo: [
    { required: true, message: '请输入买方信息', trigger: 'blur' }
  ],
  buyerId: [
    { required: true, message: '请输入买方证件号', trigger: 'blur' }
  ],
  sellerInfo: [
    { required: true, message: '请输入卖方信息', trigger: 'blur' }
  ],
  sellerId: [
    { required: true, message: '请输入卖方证件号', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入合同金额', trigger: 'blur' },
    { type: 'number', min: 0, message: '金额必须大于等于0', trigger: 'blur' }
  ],
  templateId: [
    { required: true, message: '请选择合同模板', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入合同内容', trigger: 'blur' },
    { min: 100, message: '合同内容不能少于100字符', trigger: 'blur' }
  ]
}

// 交易选项
const transactionOptions = ref([])
// 模板选项
const templateOptions = ref([
  { id: 'T001', name: '标准购房合同模板' },
  { id: 'T002', name: '商业房屋购买合同模板' },
  { id: 'T003', name: '二手房买卖合同模板' },
  { id: 'T004', name: '住房按揭贷款合同模板' },
  { id: 'T005', name: '房屋租赁合同模板' }
])

// 获取交易列表
const fetchTransactions = async () => {
  try {
    loading.value = true
    
    // 调用API获取交易列表
    const { data } = await axios.get('/transactions', {
      params: {
        pageSize: 100,
        status: 'active' // 假设只获取活跃状态的交易
      }
    })
    
    if (data.code === 200) {
      transactionOptions.value = data.data.items || []
    } else {
      ElMessage.error(data.message || '获取交易列表失败')
    }
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
    ElMessage.error('获取交易列表失败')
  } finally {
    loading.value = false
  }
}

// 处理交易选择变化
const handleTransactionChange = (transactionId) => {
  // 查找选中的交易
  const selectedTransaction = transactionOptions.value.find(item => item.id === transactionId)
  
  if (selectedTransaction) {
    // 自动填充相关信息
    contractForm.title = `${selectedTransaction.title}购买合同`
    contractForm.amount = selectedTransaction.amount || 0
    
    // 如果有交易方信息，也可以自动填充
    if (selectedTransaction.buyer) {
      contractForm.buyerInfo = selectedTransaction.buyer.name || ''
      contractForm.buyerId = selectedTransaction.buyer.id || ''
    }
    
    if (selectedTransaction.seller) {
      contractForm.sellerInfo = selectedTransaction.seller.name || ''
      contractForm.sellerId = selectedTransaction.seller.id || ''
    }
  }
}

// 选择模板
const selectTemplate = (templateId) => {
  const template = templateOptions.value.find(t => t.id === templateId)
  if (template) {
    contractForm.title = `${template.name} - ${new Date().toLocaleDateString()}`
    // 这里可以加载模板内容
    contractForm.content = `这是${template.name}的默认内容，可以根据需要进行修改。\n\n合同条款：\n1. 甲方（买方）...\n2. 乙方（卖方）...\n3. 房产信息...\n4. 交易价格...\n5. 支付方式...\n6. 交付条件...\n7. 违约责任...\n8. 争议解决...`
  }
}

// 提交合同
const submitContract = async () => {
  if (!contractFormRef.value) return
  
  await contractFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      submitting.value = true
      
      // 调用API创建合同
      await createContract({
        title: contractForm.title,
        content: contractForm.content,
        contractType: contractForm.contractType,
        creatorCitizenID: contractForm.creatorCitizenID
      })
      
      ElMessage.success('合同创建成功')
      router.push('/contract')
    } catch (error) {
      console.error('Failed to create contract:', error)
      ElMessage.error('创建合同失败')
    } finally {
      submitting.value = false
    }
  })
}

// 重置表单
const resetForm = () => {
  if (contractFormRef.value) {
    contractFormRef.value.resetFields()
    contractForm.contractType = 'PURCHASE'
    contractForm.title = ''
    contractForm.content = ''
    selectedTemplate.value = ''
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 页面加载时获取相关数据
onMounted(() => {})
</script>

<style scoped>
.contract-create-container {
  padding: 20px;
}

.contract-form-card {
  margin-bottom: 20px;
}

.page-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
}
</style>
