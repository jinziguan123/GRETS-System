<template>
  <div class="payment-create-container">

    <el-card class="form-card">
      <el-form 
        ref="formRef" 
        :model="paymentForm" 
        :rules="rules" 
        label-width="120px" 
        label-position="right"
        v-loading="loading"
      >
        <!-- 基本信息 -->
        <h3>基本信息</h3>
        <el-divider></el-divider>
        
        <el-form-item label="关联交易" prop="transactionID">
          <el-select 
            v-model="paymentForm.transactionID" 
            filterable 
            remote 
            placeholder="请选择关联交易" 
            :remote-method="searchTransactions"
            :loading="transactionsLoading"
            @change="handleTransactionChange"
            style="width: 100%"
          >
            <el-option 
              v-for="item in transactionOptions" 
              :key="item.transactionUUID" 
              :label="item.transactionUUID" 
              :value="item.transactionUUID"
            >
              <div class="transaction-option">
                <div>交易ID: {{ item.transactionUUID }}</div>
                <div>房产ID: {{ item.realtyCertHash }}</div>
                <div>总价: {{ formatCurrency(item.price) }}</div>
                <div>
                  状态: 
                  <el-tag size="small" :type="getStatusType(item.status)">
                    {{ formatStatus(item.status) }}
                  </el-tag>
                </div>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="支付类型" prop="paymentType">
          <el-radio-group v-model="paymentForm.paymentType" @change="handlePaymentTypeChange">
            <el-radio-button label="fullpayment">全款</el-radio-button>
            <el-radio-button label="loan">贷款</el-radio-button>
            <el-radio-button label="installment">分期</el-radio-button>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="支付金额" prop="amount">
          <el-input-number 
            v-model="paymentForm.amount" 
            :precision="2" 
            :min="0.01" 
            :max="999999999" 
            style="width: 200px"
          ></el-input-number>
        </el-form-item>
        
        <el-form-item v-if="transactionSelected">
          <el-alert
            :title="`目前交易已支付金额: ${formatCurrency(currentTransactionPaidAmount)} / ${formatCurrency(selectedTransaction?.price || 0)}`"
            type="info"
            :closable="false"
          />
        </el-form-item>
        
        <!-- 贷款信息 -->
        <template v-if="paymentForm.paymentType === 'loan'">
          <h3>贷款信息</h3>
          <el-divider></el-divider>
          
          <el-form-item label="贷款银行" prop="loanInfo.bankID">
            <el-select v-model="paymentForm.loanInfo.bankID" placeholder="请选择贷款银行">
              <el-option label="中国银行" value="bank1"></el-option>
              <el-option label="工商银行" value="bank2"></el-option>
              <el-option label="建设银行" value="bank3"></el-option>
              <el-option label="农业银行" value="bank4"></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="贷款金额" prop="loanInfo.loanAmount">
            <el-input-number 
              v-model="paymentForm.loanInfo.loanAmount" 
              :precision="2" 
              :min="0.01" 
              :max="999999999" 
              style="width: 200px"
              @change="updateAmount"
            ></el-input-number>
          </el-form-item>
          
          <el-form-item label="贷款利率(%)" prop="loanInfo.interestRate">
            <el-input-number 
              v-model="paymentForm.loanInfo.interestRate" 
              :precision="2" 
              :min="0.01" 
              :max="20" 
              style="width: 200px"
            ></el-input-number>
          </el-form-item>
          
          <el-form-item label="贷款期限(月)" prop="loanInfo.loanTerm">
            <el-input-number 
              v-model="paymentForm.loanInfo.loanTerm" 
              :precision="0" 
              :min="1" 
              :max="360" 
              style="width: 200px"
              @change="calculateMonthlyPayment"
            ></el-input-number>
          </el-form-item>
          
          <el-form-item label="开始日期" prop="loanInfo.startDate">
            <el-date-picker 
              v-model="paymentForm.loanInfo.startDate" 
              type="date" 
              placeholder="选择开始日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              @change="updateEndDate"
            ></el-date-picker>
          </el-form-item>
          
          <el-form-item label="结束日期" prop="loanInfo.endDate">
            <el-date-picker 
              v-model="paymentForm.loanInfo.endDate" 
              type="date" 
              placeholder="选择结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              :disabled="true"
            ></el-date-picker>
          </el-form-item>
          
          <el-form-item label="月供金额" prop="loanInfo.monthlyAmount">
            <el-input-number 
              v-model="paymentForm.loanInfo.monthlyAmount" 
              :precision="2" 
              :min="0" 
              :disabled="true"
              style="width: 200px"
            ></el-input-number>
          </el-form-item>
        </template>
        
        <!-- 分期付款信息 -->
        <template v-if="paymentForm.paymentType === 'installment'">
          <h3>分期付款信息</h3>
          <el-divider></el-divider>
          
          <el-form-item label="分期期数" prop="installmentInfo.totalInstallments">
            <el-input-number 
              v-model="paymentForm.installmentInfo.totalInstallments" 
              :precision="0" 
              :min="2" 
              :max="60" 
              style="width: 200px"
              @change="generateInstallmentPlan"
            ></el-input-number>
          </el-form-item>
          
          <el-form-item label="首付比例(%)" prop="installmentInfo.downPaymentPercent">
            <el-input-number 
              v-model="paymentForm.installmentInfo.downPaymentPercent" 
              :precision="0" 
              :min="10" 
              :max="90" 
              style="width: 200px"
              @change="generateInstallmentPlan"
            ></el-input-number>
          </el-form-item>
          
          <el-form-item label="首次付款日期" prop="installmentInfo.firstPaymentDate">
            <el-date-picker 
              v-model="paymentForm.installmentInfo.firstPaymentDate" 
              type="date" 
              placeholder="选择首次付款日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              @change="generateInstallmentPlan"
            ></el-date-picker>
          </el-form-item>
          
          <h4>分期计划预览</h4>
          <el-table :data="installmentPlan" style="width: 100%" border>
            <el-table-column prop="installmentNumber" label="期数" width="80"></el-table-column>
            <el-table-column prop="amount" label="金额" width="150">
              <template #default="scope">
                {{ formatCurrency(scope.row.amount) }}
              </template>
            </el-table-column>
            <el-table-column prop="dueDate" label="付款日期" width="180"></el-table-column>
            <el-table-column prop="description" label="说明"></el-table-column>
          </el-table>
        </template>
        
        <div class="form-actions">
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" @click="submitForm">提交</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import { useUserStore } from '@/stores/user';
import { createPayment } from '@/api/payment';
import { queryTransactionList, getTransactionById } from '@/api/transaction';

const router = useRouter();
const userStore = useUserStore();
const formRef = ref<FormInstance>();

// 表单数据
const paymentForm = reactive({
  transactionID: '',
  paymentType: 'fullpayment',
  amount: 0,
  loanInfo: {
    bankID: '',
    loanAmount: 0,
    interestRate: 4.9,
    loanTerm: 12,
    startDate: '',
    endDate: '',
    monthlyAmount: 0
  },
  installmentInfo: {
    totalInstallments: 6,
    downPaymentPercent: 30,
    firstPaymentDate: new Date().toISOString().split('T')[0],
    installments: []
  }
});

// 安装计划预览
const installmentPlan = ref<any[]>([]);

// 表单验证规则
const rules = reactive<FormRules>({
  transactionID: [
    { required: true, message: '请选择关联交易', trigger: 'change' }
  ],
  paymentType: [
    { required: true, message: '请选择支付类型', trigger: 'change' }
  ],
  amount: [
    { required: true, message: '请输入支付金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '金额必须大于0', trigger: 'blur' }
  ],
  'loanInfo.bankID': [
    { required: true, message: '请选择贷款银行', trigger: 'change' }
  ],
  'loanInfo.loanAmount': [
    { required: true, message: '请输入贷款金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '贷款金额必须大于0', trigger: 'blur' }
  ],
  'loanInfo.interestRate': [
    { required: true, message: '请输入贷款利率', trigger: 'blur' },
    { type: 'number', min: 0.01, max: 20, message: '利率必须在0.01%到20%之间', trigger: 'blur' }
  ],
  'loanInfo.loanTerm': [
    { required: true, message: '请输入贷款期限', trigger: 'blur' },
    { type: 'number', min: 1, max: 360, message: '期限必须在1到360个月之间', trigger: 'blur' }
  ],
  'loanInfo.startDate': [
    { required: true, message: '请选择开始日期', trigger: 'change' }
  ],
  'installmentInfo.totalInstallments': [
    { required: true, message: '请输入分期期数', trigger: 'blur' },
    { type: 'number', min: 2, max: 60, message: '分期期数必须在2到60之间', trigger: 'blur' }
  ],
  'installmentInfo.downPaymentPercent': [
    { required: true, message: '请输入首付比例', trigger: 'blur' },
    { type: 'number', min: 10, max: 90, message: '首付比例必须在10%到90%之间', trigger: 'blur' }
  ],
  'installmentInfo.firstPaymentDate': [
    { required: true, message: '请选择首次付款日期', trigger: 'change' }
  ]
});

// 数据和状态
const loading = ref(false);
const transactionsLoading = ref(false);
const transactionOptions = ref<any[]>([]);
const selectedTransaction = ref<any>(null);
const transactionSelected = ref(false);
const currentTransactionPaidAmount = ref(0);

// 搜索交易
const searchTransactions = async (query: string) => {
  if (query.length < 3) return;
  
  transactionsLoading.value = true;
  try {
    const response = await queryTransactionList({
      pageNumber: 1,
      pageSize: 20,
      keyword: query
    });
    
    if (response && response.data && response.data.list) {
      // 只显示进行中的交易
      transactionOptions.value = response.data.list.filter(
        (tx: any) => tx.status === 'processing' || tx.status === '进行中'
      );
    }
  } catch (error) {
    console.error('搜索交易失败:', error);
    ElMessage.error('搜索交易失败');
  } finally {
    transactionsLoading.value = false;
  }
};

// 处理交易选择
const handleTransactionChange = async (transactionID: string) => {
  if (!transactionID) {
    transactionSelected.value = false;
    selectedTransaction.value = null;
    return;
  }
  
  loading.value = true;
  try {
    const response = await getTransactionById(transactionID);
    if (response && response.data) {
      selectedTransaction.value = response.data;
      transactionSelected.value = true;
      
      // 设置最大支付金额为交易总价减去已支付金额
      const paidAmount = response.data.paidAmount || 0;
      currentTransactionPaidAmount.value = paidAmount;
      const remainingAmount = response.data.price - paidAmount;
      
      // 默认设置为剩余金额
      paymentForm.amount = remainingAmount > 0 ? remainingAmount : 0;
      
      if (paymentForm.paymentType === 'loan') {
        paymentForm.loanInfo.loanAmount = paymentForm.amount;
        calculateMonthlyPayment();
      } else if (paymentForm.paymentType === 'installment') {
        generateInstallmentPlan();
      }
    }
  } catch (error) {
    console.error('获取交易详情失败:', error);
    ElMessage.error('获取交易详情失败');
    transactionSelected.value = false;
    selectedTransaction.value = null;
  } finally {
    loading.value = false;
  }
};

// 处理支付类型变化
const handlePaymentTypeChange = (type: string) => {
  if (type === 'loan') {
    paymentForm.loanInfo.loanAmount = paymentForm.amount;
    if (paymentForm.loanInfo.startDate) {
      updateEndDate();
    }
    calculateMonthlyPayment();
  } else if (type === 'installment') {
    generateInstallmentPlan();
  }
};

// 更新贷款金额时同步支付金额
const updateAmount = () => {
  if (paymentForm.paymentType === 'loan') {
    paymentForm.amount = paymentForm.loanInfo.loanAmount;
    calculateMonthlyPayment();
  }
};

// 计算贷款月供
const calculateMonthlyPayment = () => {
  const principal = paymentForm.loanInfo.loanAmount;
  const monthlyRate = paymentForm.loanInfo.interestRate / 100 / 12;
  const term = paymentForm.loanInfo.loanTerm;
  
  if (principal && term && monthlyRate) {
    // 等额本息计算公式
    const monthlyPayment = principal * monthlyRate * Math.pow(1 + monthlyRate, term) / 
                           (Math.pow(1 + monthlyRate, term) - 1);
    paymentForm.loanInfo.monthlyAmount = parseFloat(monthlyPayment.toFixed(2));
  } else {
    paymentForm.loanInfo.monthlyAmount = 0;
  }
};

// 更新贷款终止日期
const updateEndDate = () => {
  if (!paymentForm.loanInfo.startDate) return;
  
  const startDate = new Date(paymentForm.loanInfo.startDate);
  const endDate = new Date(startDate);
  endDate.setMonth(endDate.getMonth() + paymentForm.loanInfo.loanTerm);
  
  paymentForm.loanInfo.endDate = endDate.toISOString().split('T')[0];
};

// 生成分期
const generateInstallmentPlan = () => {
  // 实现生成分期计划的逻辑
};

// 其他辅助函数
const formatCurrency = (value: number | string): string => {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return numValue.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' });
};

const formatStatus = (status: string): string => {
  const statusMap = {
    'pending': '待处理',
    'completed': '已完成',
    'failed': '失败',
    'cancelled': '已取消',
    'approved': '已批准'
  };
  return statusMap[status] || status;
};

const getStatusType = (status: string): string => {
  const typeMap = {
    'pending': 'warning',
    'completed': 'success',
    'failed': 'danger',
    'cancelled': 'info',
    'approved': 'success'
  };
  return typeMap[status] || '';
};

const goBack = () => {
  router.back();
};

const submitForm = () => {
  // 实现提交表单的逻辑
};
</script>

<style scoped>
.payment-create-container {
  padding: 20px;
}

.create-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.form-card {
  margin-bottom: 20px;
}

.transaction-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-actions {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
