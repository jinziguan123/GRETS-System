<template>
  <div class="payment-container">
    <div class="payment-header">
      <h2>支付记录列表</h2>
    </div>

    <!-- 搜索条件 -->
    <el-card class="search-card">
      <el-form :model="searchForm">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item label="支付UUID">
              <el-input v-model="searchForm.paymentUUID" placeholder="请输入支付UUID" clearable></el-input>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="关联交易UUID">
              <el-input v-model="searchForm.transactionUUID" placeholder="请输入交易UUID" clearable></el-input>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item label="支付类型">
              <el-select v-model="searchForm.paymentType" placeholder="请选择支付类型" clearable>
                <el-option label="房款" value="TRANSFER"></el-option>
                <el-option label="税费" value="TAX"></el-option>
                <el-option label="手续费" value="FEE"></el-option>
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :span="6">
            <el-form-item>
              <el-button type="primary" @click="handleSearch">搜索</el-button>
              <el-button @click="resetSearch">重置</el-button>
            </el-form-item>
          </el-col>
        </el-row>

      </el-form>

    </el-card>

    <!-- 支付列表 -->
    <el-card class="table-card">
      <div v-loading="loading">
        <el-table :data="paymentList" stripe border style="width: 100%">
          <el-table-column prop="paymentUUID" label="支付UUID" width="180" show-overflow-tooltip></el-table-column>
          <el-table-column prop="transactionUUID" label="交易UUID" width="180" show-overflow-tooltip></el-table-column>
          <el-table-column prop="payerCitizenIDHash" label="付款方" width="150" show-overflow-tooltip></el-table-column>
          <el-table-column prop="payerOrganization" label="付款方身份" width="150" show-overflow-tooltip>
            <template #default="scope">
              <span>{{ getCurrentOwnerOrganization(scope.row.payerOrganization) || '未知' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="receiverCitizenIDHash" label="收款方" width="150" show-overflow-tooltip></el-table-column>
          <el-table-column prop="receiverOrganization" label="收款方身份" width="150" show-overflow-tooltip>
            <template #default="scope">
              <span>{{ getCurrentOwnerOrganization(scope.row.receiverOrganization) || '未知' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="120">
            <template #default="scope">
              <span>{{ formatCurrency(scope.row.amount) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="paymentType" label="支付类型" width="100">
            <template #default="scope">
              <el-tag :type="getPaymentTypeTag(scope.row.paymentType)">
                {{ formatPaymentType(scope.row.paymentType) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="创建时间">
            <template #default="scope">
              {{ formatDateTime(scope.row.createTime) }}
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.pageNumber"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="pagination.total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/stores/user';
import {
  getPaymentList
} from '@/api/payment';

const router = useRouter();
const userStore = useUserStore();

// 分页参数
const pagination = reactive({
  pageNumber: 1,
  pageSize: 10,
  total: 0
});

// 搜索表单
const searchForm = reactive({
  paymentUUID: '',
  transactionUUID: '',
  paymentType: ''
});

// 数据和状态
const loading = ref(false);
const paymentList = ref<any[]>([]);
const detailsDialogVisible = ref(false);
const selectedPayment = ref<any>(null);

// 获取支付列表
const fetchPayments = async () => {
  loading.value = true;
  try {
    const response = await getPaymentList({
      pageNumber: pagination.pageNumber,
      pageSize: pagination.pageSize,
      ...searchForm
    });

    paymentList.value = response.paymentList || [];
    pagination.total = response.total || 0;
  } catch (error) {
    console.error('获取支付列表失败:', error);
    ElMessage.error('获取支付列表失败');
  } finally {
    loading.value = false;
  }
};

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

// 搜索处理
const handleSearch = () => {
  pagination.pageNumber = 1;
  fetchPayments();
};

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = '';
  });
  pagination.pageNumber = 1;
  fetchPayments();
};

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.pageSize = size;
  fetchPayments();
};

const handleCurrentChange = (current: number) => {
  pagination.pageNumber = current;
  fetchPayments();
};

// 创建支付
const handleCreatePayment = () => {
  router.push('/payment/create');
};

// 查看详情
const handleViewDetails = async (row: any) => {

};

// 格式化函数
const formatCurrency = (value: number | string): string => {
  const numValue = typeof value === 'string' ? parseFloat(value) : value;
  return numValue.toLocaleString('zh-CN', { style: 'currency', currency: 'CNY' });
};

const formatDateTime = (dateStr: string): string => {
  if (!dateStr) return '未知';
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

const formatDate = (dateStr: string): string => {
  if (!dateStr) return '未知';
  const date = new Date(dateStr);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

const formatPaymentType = (type: string): string => {
  const typeMap = {
    'TRANSFER': '房款',
    'TAX': '税费',
    'FEE': '手续费',
  };
  return typeMap[type] || type;
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

const getPaymentTypeTag = (type: string): string => {
  const tagMap = {
    'fullpayment': '',
    'loan': 'primary',
    'installment': 'warning',
    'cash': 'success',
    'transfer': 'info'
  };
  return tagMap[type] || '';
};

// 初始化数据
onMounted(() => {
  fetchPayments();
});
</script>

<style scoped>
.payment-container {
  padding: 20px;
}

.payment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.payment-details {
  margin-top: 20px;
}

.loan-info, .installment-info, .blockchain-info {
  margin-top: 20px;
}

h3 {
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 600;
  color: #606266;
}
</style>
