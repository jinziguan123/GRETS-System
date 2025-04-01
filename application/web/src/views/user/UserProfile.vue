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
          <el-input v-model="profileForm.name" placeholder="请输入真实姓名" />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user.ts'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { getUserInfo, updateUserInfo } from '@/api/user'
import type { UserInfo } from '@/types/api'

// 表单引用
const formRef = ref<FormInstance | null>(null)
const loading = ref<boolean>(false)
const userStore = useUserStore()

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
  name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
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
    const res = await getUserInfo()
    const userInfo = res.data
    
    // 填充表单数据
    Object.assign(profileForm, {
      name: userInfo.name,
      citizenID: userInfo.citizenID,
      phone: userInfo.phone,
      email: userInfo.email,
      organization: userInfo.organization,
      balance: userInfo.balance
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
        await updateUserInfo({
          name: profileForm.name,
          phone: profileForm.phone,
          email: profileForm.email
        })
        
        ElMessage.success('个人资料更新成功')
        
        // 更新Pinia中的用户信息
        userStore.updateUserInfo({
          name: profileForm.name,
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

// 组件挂载时获取用户信息
onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.user-profile-container {
  padding: 20px;
}

.profile-form {
  max-width: 500px;
}
</style>
