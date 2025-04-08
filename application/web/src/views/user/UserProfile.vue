<template>
  <div class="user-profile-container">
    <div class="card">
      <div class="card-title">个人资料</div>
      
      <el-form
        ref="formRef"
        :model="profileForm"
        :rules="profileRules"
        label-width="120px"
        class="profile-form"
        v-loading="loading"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="profileForm.name" disabled />
        </el-form-item>
        
        <el-form-item label="身份证号" prop="citizenID">
          <el-input v-model="profileForm.citizenID" disabled />
        </el-form-item>
        
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="profileForm.phone" placeholder="请输入手机号码" />
        </el-form-item>
        
        <el-form-item label="电子邮箱" prop="email">
          <el-input v-model="profileForm.email" placeholder="请输入电子邮箱" />
        </el-form-item>
        
        <el-form-item label="所属组织" prop="organization">
          <el-input v-model="organizationName" disabled />
        </el-form-item>
        
        <el-form-item v-if="showBalance" label="账户余额" prop="balance">
          <el-input v-model="profileForm.balance" disabled>
            <template #append>元</template>
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存修改</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 用户拥有的房产列表 -->
    <div class="card mt-20" v-loading="realtyLoading">
      <div class="card-title">我的房产</div>
      
      <div v-if="userRealties.length === 0" class="empty-realty">
        <el-empty description="暂无房产信息"></el-empty>
      </div>
      
      <div v-else>
        <el-table :data="userRealties" style="width: 100%">
          <el-table-column prop="realtyCert" label="不动产证号" width="180"></el-table-column>
          <el-table-column prop="realtyType" label="房产类型" width="120">
            <template #default="scope">
              {{ getRealtyTypeText(scope.row.realtyType) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusTagType(scope.row.status)">
                {{ getStatusText(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="area" label="面积" width="120">
            <template #default="scope">
              {{ scope.row.area }} 平方米
            </template>
          </el-table-column>
          <el-table-column prop="price" label="参考价格" width="150">
            <template #default="scope">
              {{ formatPrice(scope.row.price) }}
            </template>
          </el-table-column>
          <el-table-column prop="createTime" label="登记时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="120">
            <template #default="scope">
              <el-button type="primary" link @click="viewRealtyDetail(scope.row.id)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="realtyQuery.pageNumber"
            v-model:page-size="realtyQuery.pageSize"
            :page-sizes="[5, 10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            :total="realtyTotal"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user.ts'
import {ElMessage, type FormInstance, type FormRules} from 'element-plus'
import { updateUserInfo } from '@/api/user'
import { useRouter } from 'vue-router'
import { queryRealtyList } from '@/api/realty.js'
import CryptoJS from 'crypto-js'

// 表单引用
const formRef = ref<FormInstance | null>(null)
const loading = ref<boolean>(false)
const userStore = useUserStore()
const router = useRouter()

// 个人资料表单数据
interface ProfileForm {
  name: string
  citizenID: string
  phone: string
  email: string
  organization: string
  balance?: number
}

const profileForm = reactive<ProfileForm>({
  name: '',
  citizenID: '',
  phone: '',
  email: '',
  organization: '',
  balance: 0
})

// 组织名称
const organizationName = computed(() => {
  return userStore.getOrganizationName(profileForm.organization)
})

// 是否显示余额
const showBalance = computed(() => {
  return profileForm.organization === 'investor'
})

// 表单验证规则
const profileRules = reactive<FormRules>({
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入电子邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的电子邮箱地址', trigger: 'blur' }
  ]
})

// 获取用户信息
const fetchUserInfo = async () => {
  loading.value = true
  try {
    // 填充表单数据
    Object.assign(profileForm, {
      name: userStore.user?.name,
      citizenID: userStore.user?.citizenID,
      phone: userStore.user?.phone,
      email: userStore.user?.email,
      organization: userStore.user?.organization,
      balance: userStore.user?.balance
    })
  } catch (error) {
    console.error('Get user info error:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

// 提交表单
const submitForm = () => {
  if (!formRef.value) return
  
  formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        // 只提交电话和邮箱信息
        await updateUserInfo({
          phone: profileForm.phone,
          email: profileForm.email
        })
        
        ElMessage.success('个人资料更新成功')
        
        // 更新Pinia中的用户信息
        userStore.updateUserInfo({
          phone: profileForm.phone,
          email: profileForm.email
        })
      } catch (error) {
        console.error('Update profile error:', error)
        ElMessage.error('个人资料更新失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  fetchUserInfo()
}

// 房产列表相关
const userRealties = ref([])
const realtyLoading = ref(false)
const realtyTotal = ref(0)
const realtyQuery = reactive({
  pageSize: 5,
  pageNumber: 1
})

// 获取用户拥有的房产
const fetchUserRealties = async () => {
  realtyLoading.value = true
  try {
    // 使用用户身份证哈希作为筛选条件
    const user = userStore.user
    if (!user || !user.citizenID) {
      ElMessage.error('用户信息不完整，无法获取房产信息')
      return
    }
    
    const citizenIDHash = CryptoJS.SHA256(user.citizenID).toString()
    
    // 构造查询参数，使用适当的参数名称
    const params = {
      pageSize: realtyQuery.pageSize,
      pageNumber: realtyQuery.pageNumber,
    }
    
    // 添加额外的查询条件
    // 注意：实际参数名称可能需要根据API接口定义调整
    const queryParams = Object.assign({}, params, { owner: citizenIDHash })
    
    const response = await queryRealtyList(queryParams)
    
    // 确保响应数据结构正确
    userRealties.value = response.data?.realties || []
    realtyTotal.value = response.data?.total || 0
    
    if (userRealties.value.length === 0 && realtyTotal.value > 0 && realtyQuery.pageNumber > 1) {
      // 如果当前页没有数据，但总数大于0，可能是页码超出范围，回到第一页
      realtyQuery.pageNumber = 1
      fetchUserRealties()
    }
  } catch (error) {
    console.error('获取用户房产失败:', error)
    ElMessage.error('获取用户房产信息失败')
    userRealties.value = []
    realtyTotal.value = 0
  } finally {
    realtyLoading.value = false
  }
}

// 查看房产详情
const viewRealtyDetail = (id: number | string) => {
  router.push(`/realty/detail/${id}`)
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  realtyQuery.pageSize = size
  fetchUserRealties()
}

// 处理页码变化
const handleCurrentChange = (page: number) => {
  realtyQuery.pageNumber = page
  fetchUserRealties()
}

// 获取房产状态对应的Tag类型
const getStatusTagType = (status: string): 'success' | 'warning' | 'info' | 'primary' | 'danger' | undefined => {
  const statusMap: Record<string, 'success' | 'warning' | 'info' | 'primary' | 'danger'> = {
    'NORMAL': 'success',
    'IN_TRANSACTION': 'warning',
    'MORTGAGED': 'info',
    'FROZEN': 'danger'
  }
  return statusMap[status] || undefined
}

// 获取房产状态对应的文本
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'NORMAL': '正常',
    'IN_TRANSACTION': '交易中',
    'MORTGAGED': '已抵押',
    'FROZEN': '已冻结'
  }
  return statusMap[status] || status
}

// 获取房产类型文本
const getRealtyTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    'HOUSE': '住宅',
    'SHOP': '商铺',
    'OFFICE': '办公',
    'INDUSTRIAL': '工业',
    'OTHER': '其他'
  }
  return typeMap[type] || type
}

// 格式化价格
const formatPrice = (price: number) => {
  return `¥ ${parseFloat(price.toString()).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })} 元`
}

// 格式化日期
const formatDate = (date: string | Date | null | undefined) => {
  if (!date) return '-'
  const dateObj = new Date(date)
  return dateObj.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 组件挂载时获取用户信息和房产列表
onMounted(() => {
  fetchUserInfo()
  fetchUserRealties()
})
</script>

<style scoped>
.user-profile-container {
  padding: 20px;
}

.card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}

.card-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
  color: #303133;
}

.profile-form {
  max-width: 500px;
}

.mt-20 {
  margin-top: 20px;
}

.empty-realty {
  padding: 30px 0;
  text-align: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
