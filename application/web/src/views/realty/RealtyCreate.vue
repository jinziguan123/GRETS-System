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
        <!-- 基础信息 -->
        <h4>基础信息</h4>
        <el-form-item label="不动产证号" prop="realtyCert">
          <el-input v-model="realtyForm.realtyCert" placeholder="请输入不动产证号"></el-input>
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
        
        <el-form-item label="户型" prop="houseType">
          <el-select v-model="realtyForm.houseType" placeholder="请选择户型" style="width: 100%">
            <el-option label="一室" value="single"></el-option>
            <el-option label="两室" value="double"></el-option>
            <el-option label="三室" value="triple"></el-option>
            <el-option label="四室及以上" value="multiple"></el-option>
          </el-select>
        </el-form-item>

        <el-row>
          <el-col :span="12">
            <el-form-item label="面积" prop="area">
              <el-input-number v-model="realtyForm.area" :min="0" :precision="2" :step="1" style="width: 100%" placeholder="请输入面积（平方米）"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格" prop="price">
              <el-input-number v-model="realtyForm.price" :min="0" :precision="2" :step="10000" style="width: 100%" placeholder="请输入价格（元）"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>

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
        
        <!-- 地址信息 -->
        <h4>地址信息</h4>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="省份" prop="province">
              <el-input v-model="realtyForm.province" placeholder="请输入省份"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="城市" prop="city">
              <el-input v-model="realtyForm.city" placeholder="请输入城市"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="区/县" prop="district">
              <el-input v-model="realtyForm.district" placeholder="请输入区/县"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-form-item label="街道" prop="street">
          <el-input v-model="realtyForm.street" placeholder="请输入街道"></el-input>
        </el-form-item>
        
        <el-form-item label="小区" prop="community">
          <el-input v-model="realtyForm.community" placeholder="请输入小区名称"></el-input>
        </el-form-item>
        
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="单元" prop="unit">
              <el-input v-model="realtyForm.unit" placeholder="请输入单元"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="楼层" prop="floor">
              <el-input v-model="realtyForm.floor" placeholder="请输入楼层"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="房号" prop="room">
              <el-input v-model="realtyForm.room" placeholder="请输入房号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 附加信息 -->
        <h4>附加信息</h4>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="realtyForm.description" 
            type="textarea" 
            rows="3" 
            placeholder="请输入房产描述"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="房产图片" prop="images">
          <el-upload
            action="#"
            list-type="picture-card"
            :auto-upload="false"
            :limit="5"
            :file-list="fileList"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
          >
            <el-icon><Plus /></el-icon>
            <template #tip>
              <div class="el-upload__tip">
                只能上传 jpg/png 文件，且不超过 5MB
              </div>
            </template>
          </el-upload>
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
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { createRealty } from "@/views/realty/api.js"

const router = useRouter()
const realtyFormRef = ref(null)
const loading = ref(false)
const fileList = ref([])

// 表单数据
const realtyForm = reactive({
  realtyCert: '',
  realtyType: '',
  houseType: '',
  price: 0,
  area: 0,
  currentOwnerCitizenID: '',
  status: 'NORMAL',
  province: '',
  city: '',
  district: '',
  street: '',
  community: '',
  unit: '',
  floor: '',
  room: '',
  description: '',
  // 这些将作为计算属性或默认值
  address: '',
  images: [],
  previousOwnersCitizenIDList: []
})

// 计算完整地址
const updateAddress = () => {
  realtyForm.address = `${realtyForm.province}省${realtyForm.city}市${realtyForm.district}区${realtyForm.street}${realtyForm.community}${realtyForm.unit}单元${realtyForm.floor}楼${realtyForm.room}`
}

// 监听地址相关字段变化
const addressFields = ['province', 'city', 'district', 'street', 'community', 'unit', 'floor', 'room']
addressFields.forEach(field => {
  watch(() => realtyForm[field], () => {
    updateAddress()
  })
})

// 文件上传相关方法
const handleFileChange = (file) => {
  console.log('文件变化:', file)
  // 实际环境中，这里会上传文件并获取URL
  // 这里只是模拟
}

const handleFileRemove = (file) => {
  console.log('移除文件:', file)
}

// 表单验证规则
const rules = reactive({
  realtyCert: [
    { required: true, message: '请输入不动产证号', trigger: 'blur' },
    { min: 5, max: 50, message: '长度应在 5 到 50 个字符之间', trigger: 'blur' }
  ],
  realtyType: [
    { required: true, message: '请选择房产类型', trigger: 'change' }
  ],
  houseType: [
    { required: true, message: '请选择户型', trigger: 'change' }
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
  ],
  province: [
    { required: true, message: '请输入省份', trigger: 'blur' }
  ],
  city: [
    { required: true, message: '请输入城市', trigger: 'blur' }
  ],
  district: [
    { required: true, message: '请输入区/县', trigger: 'blur' }
  ],
  street: [
    { required: true, message: '请输入街道', trigger: 'blur' }
  ],
  community: [
    { required: true, message: '请输入小区', trigger: 'blur' }
  ],
  unit: [
    { required: true, message: '请输入单元', trigger: 'blur' }
  ],
  floor: [
    { required: true, message: '请输入楼层', trigger: 'blur' }
  ],
  room: [
    { required: true, message: '请输入房号', trigger: 'blur' }
  ]
})

// 提交表单
const submitForm = async () => {
  if (!realtyFormRef.value) return

  // 更新地址，确保是最新的
  updateAddress()
  
  await realtyFormRef.value.validate(async (valid, fields) => {
    if (valid) {
      loading.value = true
      try {
        // 准备请求参数，确保包含所有必要字段
        const requestData = {
          realtyCert: realtyForm.realtyCert,
          address: realtyForm.address,
          realtyType: realtyForm.realtyType,
          price: realtyForm.price,
          area: realtyForm.area,
          currentOwnerCitizenID: realtyForm.currentOwnerCitizenID,
          status: realtyForm.status,
          houseType: realtyForm.houseType,
          province: realtyForm.province,
          city: realtyForm.city,
          district: realtyForm.district,
          street: realtyForm.street, 
          community: realtyForm.community,
          unit: realtyForm.unit,
          floor: realtyForm.floor,
          room: realtyForm.room,
          description: realtyForm.description,
          // 默认传递空数组
          images: [],
          previousOwnersCitizenIDList: []
        }

        await createRealty(requestData).then(res => {
          if (res.code === 200) {
            ElMessage.success('房产创建成功')
            router.push('/realty')
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
    fileList.value = []
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}
</script>

<style scoped>
.realty-create-container {
  max-width: 900px;
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

h4 {
  margin-top: 20px;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
  color: #409EFF;
}
</style>
