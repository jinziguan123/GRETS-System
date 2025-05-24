<template>
  <div class="realty-edit-container">
    <el-card class="box-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>编辑房产信息</h3>
        </div>
      </template>
      
      <el-form 
        :model="realtyForm" 
        :rules="rules" 
        ref="realtyFormRef" 
        label-width="120px"
        label-position="right"
        status-icon
      >
        <el-form-item label="不动产证号" prop="realtyCert">
          <el-input v-model="realtyForm.realtyCert" disabled></el-input>
        </el-form-item>
        
        <el-form-item label="房产地址" prop="address">
          <el-input v-model="realtyForm.address"></el-input>
        </el-form-item>
        
        <el-form-item label="房产类型" prop="realtyType">
          <el-select v-model="realtyForm.realtyType" placeholder="请选择房产类型" style="width: 100%">
            <el-option label="住宅" value="住宅"></el-option>
            <el-option label="商铺" value="商铺"></el-option>
            <el-option label="办公" value="办公"></el-option>
            <el-option label="工业" value="工业"></el-option>
            <el-option label="其他" value="其他"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="realtyForm.price" :min="0" :precision="2" :step="10000" style="width: 100%"></el-input-number>
        </el-form-item>
        
        <el-form-item label="面积" prop="area">
          <el-input-number v-model="realtyForm.area" :min="0" :precision="2" :step="1" style="width: 100%"></el-input-number>
        </el-form-item>
        
        <el-form-item label="当前所有者" prop="currentOwnerCitizenID">
          <el-input v-model="realtyForm.currentOwnerCitizenID"></el-input>
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-select v-model="realtyForm.status" placeholder="请选择房产状态" style="width: 100%">
            <el-option label="正常" value="NORMAL"></el-option>
            <el-option label="交易中" value="IN_TRANSACTION"></el-option>
            <el-option label="已抵押" value="MORTGAGED"></el-option>
            <el-option label="已冻结" value="FROZEN"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="submitLoading">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const realtyFormRef = ref(null)
const loading = ref(true)
const submitLoading = ref(false)

// 表单数据
const realtyForm = reactive({
  realtyCert: '',
  address: '',
  realtyType: '',
  price: 0,
  area: 0,
  currentOwnerCitizenID: '',
  status: 'NORMAL',
  previousOwnersCitizenIDList: []
})

// 表单验证规则
const rules = reactive({
  address: [
    { required: true, message: '请输入房产地址', trigger: 'blur' }
  ],
  realtyType: [
    { required: true, message: '请选择房产类型', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格必须大于等于0', trigger: 'blur' }
  ],
  area: [
    { required: true, message: '请输入面积', trigger: 'blur' },
    { type: 'number', min: 0, message: '面积必须大于等于0', trigger: 'blur' }
  ],
  currentOwnerCitizenID: [
    { required: true, message: '请输入当前所有者身份证号', trigger: 'blur' },
    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择房产状态', trigger: 'change' }
  ]
})

// 获取房产详情
const fetchRealtyDetail = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/api/realty/${route.params.id}`)
    Object.assign(realtyForm, response.data)
  } catch (error) {
    console.error('获取房产详情失败:', error)
    ElMessage.error(error.response?.data?.message || '获取房产详情失败')
    router.push('/realty/list')
  } finally {
    loading.value = false
  }
}

// 提交表单
const submitForm = async () => {
  if (!realtyFormRef.value) return
  
  await realtyFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      submitLoading.value = true
      
      try {
        const response = await axios.put(`/api/realty/${route.params.id}`, realtyForm)
        ElMessage.success('房产信息更新成功')
        router.push(`/realty/detail/${route.params.id}`)
      } catch (error) {
        console.error('更新房产信息失败:', error)
        ElMessage.error(error.response?.data?.message || '更新房产信息失败')
      } finally {
        submitLoading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 返回上一页
const goBack = () => {
  ElMessageBox.confirm('确定要取消编辑吗？未保存的内容将丢失', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    router.back()
  }).catch(() => {})
}

onMounted(() => {
  fetchRealtyDetail()
})
</script>

<style scoped>
.realty-edit-container {
  max-width: 800px;
  margin: 20px auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.box-card {
  margin-bottom: 20px;
}
</style>
