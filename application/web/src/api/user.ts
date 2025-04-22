import request from '@/utils/request'
import type { LoginData, RegisterData, UserInfo, ApiResponse } from '@/types/api'

// 用户登录
export function login(data: LoginData): Promise<ApiResponse<{ token: string, user?: UserInfo }>> {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 获取用户余额
export function getBalanceByCitizenIDHashAndOrganization(citizenID: string, organization: string): Promise<ApiResponse<{ balance: number }>> {
  return request({
    url: '/user/getBalance',
    method: 'get',
    params: { citizenID, organization }
  })
}
// 用户注册
export function register(data: RegisterData): Promise<ApiResponse<void>> {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 获取用户信息
export function getUserInfo(): Promise<ApiResponse<UserInfo>> {
  return request({
    url: '/users/info',
    method: 'get'
  })
}

// 更新用户信息
export function updateUserInfo(data: Partial<UserInfo>): Promise<ApiResponse<void>> {
  return request({
    url: '/users/info',
    method: 'put',
    data
  })
}

// 获取用户列表
export function getUserList(params: {
  organization?: string
  page?: number
  pageSize?: number
}): Promise<ApiResponse<{
  total: number
  list: UserInfo[]
}>> {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}