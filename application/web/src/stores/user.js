import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const router = useRouter()
  
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value.role || '')
  const userId = computed(() => user.value.id || '')
  const username = computed(() => user.value.username || '')
  
  // 判断是否有权限
  const hasRole = (roles) => {
    if (!roles || roles.length === 0) return true
    return roles.includes(userRole.value)
  }
  
  // 登录
  const login = async (credentials) => {
    try {
      const response = await axios.post('/api/v1/public/login', credentials)
      const data = response.data.data
      
      // 保存token和用户信息
      token.value = data.token
      user.value = data.user
      
      // 存储到localStorage
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
      localStorage.setItem('userRole', data.user.role)
      
      // 设置axios默认头部
      axios.defaults.headers.common['Authorization'] = `Bearer ${data.token}`
      
      ElMessage.success('登录成功')
      return true
    } catch (error) {
      console.error('登录失败:', error)
      ElMessage.error(error.response?.data?.message || '登录失败')
      return false
    }
  }
  
  // 注销
  const logout = () => {
    // 清除token和用户信息
    token.value = ''
    user.value = {}
    
    // 清除localStorage
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('userRole')
    
    // 清除axios默认头部
    delete axios.defaults.headers.common['Authorization']
    
    // 跳转到登录页
    router.push({ name: 'Login' })
    ElMessage.success('已退出登录')
  }
  
  // 更新用户信息
  const updateUserInfo = (newUserInfo) => {
    user.value = { ...user.value, ...newUserInfo }
    localStorage.setItem('user', JSON.stringify(user.value))
  }
  
  // 初始化axios头部
  const initAxiosHeaders = () => {
    if (token.value) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    }
  }
  
  // 初始化
  initAxiosHeaders()
  
  return {
    token,
    user,
    isLoggedIn,
    userRole,
    userId,
    username,
    hasRole,
    login,
    logout,
    updateUserInfo
  }
}) 