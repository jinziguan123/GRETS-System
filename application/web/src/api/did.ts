import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// DID注册请求接口
export interface DIDRegisterRequest {
  citizenID: string
  name: string
  phone: string
  email?: string
  password: string
  organization: string
  role?: string
  balance?: number
  publicKey: string
}

// DID注册响应接口
export interface DIDRegisterResponse {
  did: string
  didDocument: any
  credentials: any[]
  message: string
}

// DID登录请求接口
export interface DIDLoginRequest {
  did: string
  challenge: string
  signature: string
  publicKey: string
}

// DID登录响应接口
export interface DIDLoginResponse {
  token: string
  did: string
  didDocument: any
  user: {
    did: string
    name: string
    organization: string
    role: string
    citizenID?: string
  }
}

// 获取挑战请求接口
export interface GetChallengeRequest {
  did: string
  domain?: string
}

// 获取挑战响应接口
export interface GetChallengeResponse {
  challenge: string
  nonce: string
  domain: string
  expiresAt: string
}

// 创建DID请求接口
export interface CreateDIDRequest {
  citizenID: string
  organization: string
  role: string
  publicKey: string
  name: string
  phone: string
  email?: string
}

// 创建DID响应接口
export interface CreateDIDResponse {
  did: string
  didDocument: any
  credentials: any[]
}

/**
 * DID注册（兼容传统注册）
 */
export function didRegister(data: DIDRegisterRequest): Promise<ApiResponse<DIDRegisterResponse>> {
  return request({
    url: '/did/register',
    method: 'post',
    data
  })
}

/**
 * 获取认证挑战
 */
export function getChallenge(data: GetChallengeRequest): Promise<ApiResponse<GetChallengeResponse>> {
  return request({
    url: '/did/challenge',
    method: 'post',
    data
  })
}

/**
 * DID登录
 */
export function didLogin(data: DIDLoginRequest): Promise<ApiResponse<DIDLoginResponse>> {
  return request({
    url: '/did/login',
    method: 'post',
    data
  })
}

/**
 * 创建DID
 */
export function createDID(data: CreateDIDRequest): Promise<ApiResponse<CreateDIDResponse>> {
  return request({
    url: '/did/create',
    method: 'post',
    data
  })
}

/**
 * 解析DID
 */
export function resolveDID(did: string): Promise<ApiResponse<{ didDocument: any; metadata: any }>> {
  return request({
    url: `/did/resolve/${encodeURIComponent(did)}`,
    method: 'get'
  })
}

/**
 * 根据用户信息获取DID
 */
export function getDIDByUser(citizenID: string, organization: string): Promise<ApiResponse<{ did: string }>> {
  return request({
    url: '/did/user',
    method: 'get',
    params: { citizenID, organization }
  })
}

/**
 * 获取凭证
 */
export function getCredentials(data: { did: string; credentialType?: string }): Promise<ApiResponse<{ credentials: any[] }>> {
  return request({
    url: '/credentials/get',
    method: 'post',
    data
  })
}

/**
 * 签发凭证
 */
export function issueCredential(data: {
  issuerDid: string
  subjectDid: string
  credentialType: string
  claims: Record<string, any>
  expirationDate?: string
}): Promise<ApiResponse<{ credential: any }>> {
  return request({
    url: '/credentials/issue',
    method: 'post',
    data
  })
}

/**
 * 撤销凭证
 */
export function revokeCredential(data: { credentialId: string; reason?: string }): Promise<ApiResponse<void>> {
  return request({
    url: '/credentials/revoke',
    method: 'post',
    data
  })
}

/**
 * 验证展示
 */
export function verifyPresentation(data: { presentation: any }): Promise<ApiResponse<{ valid: boolean; reason?: string }>> {
  return request({
    url: '/credentials/verify',
    method: 'post',
    data
  })
} 