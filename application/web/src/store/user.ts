import { defineStore } from 'pinia'
import { login, getUserInfo } from '@/api/user'
import type { UserInfo as UserInfoType, LoginData } from '@/types/api'

// 定义用户角色类型
export type UserRole = 'GOVERNMENT' | 'BANK' | 'INVESTOR' | 'AUDITOR' | 'THIRDPARTY'

// 定义用户信息类型
interface UserInfo {
  citizenID: string
  name: string
  role: UserRole
  organization: string
  phone: string
  email: string
  balance?: number
  createdTime: string
  lastUpdateTime: string
  status: string
}

// 定义store的state类型
interface UserState {
  token: string
  userInfo: UserInfoType | null
  roles: string[]
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    roles: []
  }),
  
  getters: {
    isLoggedIn: (state: UserState): boolean => !!state.token,
    isGovernment: (state: UserState): boolean => state.userInfo?.role === 'GOVERNMENT',
    isInvestor: (state: UserState): boolean => state.userInfo?.role === 'INVESTOR',
    isBank: (state: UserState): boolean => state.userInfo?.role === 'BANK',
    isAuditor: (state: UserState): boolean => state.userInfo?.role === 'AUDITOR',
    canCreateRealty: (state: UserState): boolean => state.userInfo?.role === 'GOVERNMENT',
    canUpdateRealty: (state: UserState): boolean => ['GOVERNMENT', 'INVESTOR'].includes(state.userInfo?.role || ''),
    canCreateTransaction: (state: UserState): boolean => ['GOVERNMENT', 'INVESTOR'].includes(state.userInfo?.role || ''),
    canCheckTransaction: (state: UserState): boolean => ['GOVERNMENT', 'INVESTOR', 'AUDITOR'].includes(state.userInfo?.role || ''),
    canCreatePayment: (state: UserState): boolean => ['INVESTOR', 'BANK'].includes(state.userInfo?.role || '')
  },
  
  actions: {
    // 设置token
    setToken(token: string): void {
      this.token = token
      localStorage.setItem('token', token)
    },
    
    // 清除token
    clearToken(): void {
      this.token = ''
      localStorage.removeItem('token')
    },
    
    // 登录
    async loginAction(loginData: LoginData): Promise<void> {
      try {
        const response = await login(loginData)
        this.setToken(response.data.token)
      } catch (error) {
        console.error('Login error:', error)
        throw error
      }
    },
    
    // 登出
    logout(): void {
      this.clearToken()
      this.userInfo = null
      this.roles = []
    },
    
    // 获取用户信息
    async getUserInfoAction(): Promise<UserInfoType> {
      try {
        const response = await getUserInfo()
        this.userInfo = response.data
        this.roles = [response.data.role]
        return response.data
      } catch (error) {
        console.error('Get user info error:', error)
        throw error
      }
    },
    
    // 页面刷新时重新获取用户信息
    async restoreUserInfo(): Promise<void> {
      if (this.token && !this.userInfo) {
        await this.getUserInfoAction()
      }
    }
  }
}) 