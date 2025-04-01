import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import type { UserInfo } from '@/types/api'

export type UserRole = 'GOVERNMENT' | 'BANK' | 'INVESTOR' | 'AUDIT' | 'ADMINISTRATOR'
export type Organization = 'government' | 'bank' | 'investor' | 'audit' | 'administrator'

interface LoginResponse {
  code: number
  message: string
  data: {
    token: string
    user: UserInfo
  }
}

export const useUserStore = defineStore('user', () => {
  const router = useRouter()

  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<UserInfo | null>(JSON.parse(localStorage.getItem('userInfo') || 'null'))

  // 计算属性
  const isLoggedIn = computed<boolean>(() => !!token.value)
  const userRole = computed<string>(() => user.value?.role || '')
  const userId = computed<string>(() => user.value?.citizenID || '')
  const username = computed<string>(() => user.value?.name || '')
  const organization = computed<string>(() => user.value?.organization || '')

  // 判断是否有权限
  const hasRole = (roles: string | string[]): boolean => {
    if (!roles || (Array.isArray(roles) && roles.length === 0)) return true
    if (typeof roles === 'string') {
      return userRole.value === roles
    }
    return roles.includes(userRole.value)
  }

  // 判断是否属于某个组织
  const hasOrganization = (orgs: string | string[]): boolean => {
    if (!orgs) return true
    if (typeof orgs === 'string') {
      return organization.value === orgs
    }
    return orgs.includes(organization.value)
  }

  // 获取组织名称
  const getOrganizationName = (org?: string): string => {
    const orgMap: Record<string, string> = {
      'administrator': '系统管理员',
      'government': '政府监管部门',
      'investor': '投资者/买家',
      'bank': '银行机构',
      'audit': '审计监管部门'
    }
    return org ? orgMap[org] || '未知组织' : (orgMap[organization.value] || '未知组织')
  }

  // 登录
  const login = async (userInfo: { citizenID: string; password: string; organization: string }): Promise<LoginResponse> => {
    try {
      const { data } = await axios.post<LoginResponse>('/login', {
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
  const logout = (): void => {
    // 清除token和用户信息
    token.value = ''
    user.value = null

    // 清除localStorage
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')

    // 清除axios默认头部
    delete axios.defaults.headers.common['Authorization']

    // 跳转到登录页
    // router.push({ name: 'Login' })
    ElMessage.success('已退出登录')
  }

  // 更新用户信息
  const updateUserInfo = (newUserInfo: Partial<UserInfo>): void => {
    if (user.value) {
      user.value = { ...user.value, ...newUserInfo }
      localStorage.setItem('userInfo', JSON.stringify(user.value))
    } else {
      user.value = newUserInfo as UserInfo
      localStorage.setItem('userInfo', JSON.stringify(user.value))
    }
  }

  // 初始化axios头部
  const initAxiosHeaders = (): void => {
    if (token.value) {
      axios.defaults.headers.common['Authorization'] = `${token.value}`
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