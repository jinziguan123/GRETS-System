import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const router = useRouter()

  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => user.value.role || '')
  const userId = computed(() => user.value.id || '')
  const username = computed(() => user.value.name || '')
  const organization = computed(() => user.value.organization || '')

  // 判断是否有权限
  const hasRole = (roles) => {
    if (!roles || roles.length === 0) return true
    return roles.includes(userRole.value)
  }

  // 判断是否属于某个组织
  const hasOrganization = (orgs) => {
    if (!orgs) return true
    if (typeof orgs === 'string') {
      return organization.value === orgs
    }
    return orgs.includes(organization.value)
  }

  // 获取组织名称
  const getOrganizationName = (org) => {
    const orgMap = {
      'administrator': '系统管理员',
      'government': '政府监管部门',
      'investor': '投资者/买家',
      'bank': '银行机构',
      'audit': '审计监管部门'
    }
    return org ? orgMap[org] : (orgMap[organization.value] || '未知组织')
  }

  // 登录
  const login = async (userInfo) => {
    try {
      const { data } = await axios.post('/login', {
        // 注意：后端使用citizenID而不是username
        citizenID: userInfo.citizenID,
        password: userInfo.password,
        organization: userInfo.organization
      })

      if (data.code === 200) {
        // 保存token
        token.value = data.data.token
        localStorage.setItem('token', data.data.token)
        
        // 保存用户信息
        user.value = data.data.user
        localStorage.setItem('userInfo', JSON.stringify(data.data.user))
        
        // 设置axios认证头
        axios.defaults.headers.common['Authorization'] = `Bearer ${data.data.token}`
        
        return data
      } else {
        throw new Error(data.message || '登录失败')
      }
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  // 注销
  const logout = () => {
    // 清除token和用户信息
    token.value = ''
    user.value = {}

    // 清除localStorage
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')

    // 清除axios默认头部
    delete axios.defaults.headers.common['Authorization']

    // 跳转到登录页
    router.push({ name: 'Login' })
    ElMessage.success('已退出登录')
  }

  // 更新用户信息
  const updateUserInfo = (newUserInfo) => {
    user.value = { ...user.value, ...newUserInfo }
    localStorage.setItem('userInfo', JSON.stringify(user.value))
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
    organization,
    hasRole,
    hasOrganization,
    getOrganizationName,
    login,
    logout,
    updateUserInfo
  }
}) 