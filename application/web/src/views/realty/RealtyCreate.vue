<template>
  <div class="realty-create-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <h3>创建新房产</h3>
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
          <el-input v-model="realtyForm.realtyCert" placeholder="请输入不动产证号"></el-input>
        </el-form-item>
        
        <el-form-item label="房产地址" prop="address">
          <el-input v-model="realtyForm.address" placeholder="请输入房产地址"></el-input>
        </el-form-item>
        
        <el-form-item label="房产类型" prop="realtyType">
          <el-select v-model="realtyForm.realtyType" placeholder="请选择房产类型" style="width: 100%">
            <el-option label="住宅" value="HOUSE"></el-option>
            <el-option label="商铺" value="SHOP"></el-option>
            <el-option label="办公" value="OFFICE"></el-option>
            <el-option label="工业" value="INDUSTRIAL"></el-option>
            <el-option label="其他" value="OTHER"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="realtyForm.price" :min="0" :precision="2" :step="10000" style="width: 100%" placeholder="请输入价格（元）"></el-input-number>
        </el-form-item>
        
        <el-form-item label="面积" prop="area">
          <el-input-number v-model="realtyForm.area" :min="0" :precision="2" :step="1" style="width: 100%" placeholder="请输入面积（平方米）"></el-input-number>
        </el-form-item>
        
        <el-form-item label="当前所有者" prop="currentOwnerCitizenID">
          <el-input v-model="realtyForm.currentOwnerCitizenID" placeholder="请输入当前所有者身份证号"></el-input>
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
          <el-button type="primary" @click="submitForm" :loading="loading">创建</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import axios from 'axios'
import {createRealty} from "@/views/realty/api.js";

const router = useRouter()
const realtyFormRef = ref(null)
const loading = ref(false)

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
  realtyCert: [
    { required: true, message: '请输入不动产证号', trigger: 'blur' },
    { min: 5, max: 50, message: '长度应在 5 到 50 个字符之间', trigger: 'blur' }
  ],
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

// 提交表单
const submitForm = async () => {
  if (!realtyFormRef.value) return
  
  await realtyFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      loading.value = true
      try {
        await createRealty({
          realtyCert: realtyForm.realtyCert,
          address: realtyForm.address,
          realtyType: realtyForm.realtyType,
          price: realtyForm.price,
          area: realtyForm.area,
          currentOwnerCitizenID: realtyForm.currentOwnerCitizenID,
          status: realtyForm.status,
        }).then(res => {
            if (res.code === 200) {
                ElMessage.success('房产创建成功')
                router.push('/realty/list')
            } else {
                ElMessage.error(res.message || '创建房产失败')
            }
        })
      } catch (error) {
        console.error('创建房产失败:', error)
        ElMessage.error(error.response?.data?.message || '创建房产失败')
      } finally {
        loading.value = false
      }
    } else {
      console.log('表单验证失败:', fields)
      ElMessage.error('请完善表单信息')
    }
  })
}

// 重置表单
const resetForm = () => {
  if (realtyFormRef.value) {
    realtyFormRef.value.resetFields()
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}
</script>

<style scoped>
.realty-create-container {
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
