// DID工具函数
import CryptoJS from 'crypto-js'

// 密钥对接口
export interface KeyPair {
  privateKey: string
  publicKey: string
}

// DID认证挑战接口
export interface DIDAuthChallenge {
  challenge: string
  nonce: string
  domain: string
  expiresAt: string
}

// DID认证响应接口
export interface DIDAuthResponse {
  did: string
  challenge: string
  signature: string
  publicKey: string
}

/**
 * 生成ECDSA密钥对（简化版本，实际应使用更安全的库）
 */
export function generateKeyPair(): KeyPair {
  // 这里使用简化的密钥生成，实际项目中应使用专业的加密库
  const privateKey = CryptoJS.lib.WordArray.random(32).toString()
  const publicKey = CryptoJS.SHA256(privateKey).toString()
  
  return {
    privateKey,
    publicKey
  }
}

/**
 * 使用私钥签名消息
 */
export function signMessage(message: string, privateKey: string): string {
  // 简化的签名实现，实际应使用ECDSA
  const hash = CryptoJS.SHA256(message + privateKey).toString()
  return hash
}

/**
 * 验证签名
 */
export function verifySignature(message: string, signature: string, publicKey: string): boolean {
  // 简化的验证实现
  const expectedSignature = CryptoJS.SHA256(message + publicKey).toString()
  return signature === expectedSignature
}

/**
 * 生成DID标识符
 */
export function generateDID(organization: string, citizenID: string): string {
  const identifier = CryptoJS.SHA256(citizenID + organization).toString().substring(0, 16)
  return `did:grets:${organization}:${identifier}`
}

/**
 * 创建DID文档
 */
export function createDIDDocument(did: string, organization: string, role: string, publicKey: string) {
  const now = new Date().toISOString()
  
  return {
    '@context': ['https://www.w3.org/ns/did/v1'],
    id: did,
    publicKey: [
      {
        id: `${did}#keys-1`,
        type: 'EcdsaSecp256k1VerificationKey2019',
        controller: did,
        publicKeyHex: publicKey
      }
    ],
    verificationMethod: [
      {
        id: `${did}#vm-1`,
        type: 'EcdsaSecp256k1VerificationKey2019',
        controller: did,
        publicKeyMultibase: publicKey
      }
    ],
    authentication: [`${did}#vm-1`],
    service: [
      {
        id: `${did}#grets-service`,
        type: 'GretsService',
        serviceEndpoint: 'https://grets.example.com/api/v1'
      }
    ],
    organization,
    role,
    created: now,
    updated: now
  }
}

/**
 * 生成哈希
 */
export function generateHash(data: string): string {
  return CryptoJS.SHA256(data).toString()
}

/**
 * 保存密钥对到本地存储（加密）
 */
export function saveKeyPair(keyPair: KeyPair, password: string): void {
  const encrypted = CryptoJS.AES.encrypt(JSON.stringify(keyPair), password).toString()
  localStorage.setItem('did_keypair', encrypted)
}

/**
 * 从本地存储加载密钥对（解密）
 */
export function loadKeyPair(password: string): KeyPair | null {
  try {
    const encrypted = localStorage.getItem('did_keypair')
    if (!encrypted) return null
    
    const decrypted = CryptoJS.AES.decrypt(encrypted, password).toString(CryptoJS.enc.Utf8)
    return JSON.parse(decrypted)
  } catch (error) {
    console.error('加载密钥对失败:', error)
    return null
  }
}

/**
 * 检查是否存在密钥对
 */
export function hasKeyPair(): boolean {
  return !!localStorage.getItem('did_keypair')
}

/**
 * 删除密钥对
 */
export function removeKeyPair(): void {
  localStorage.removeItem('did_keypair')
}

/**
 * 创建认证挑战响应
 */
export function createAuthResponse(
  did: string,
  challenge: DIDAuthChallenge,
  privateKey: string,
  publicKey: string
): DIDAuthResponse {
  const message = `${did}:${challenge.challenge}:${challenge.nonce}`
  const signature = signMessage(message, privateKey)
  
  return {
    did,
    challenge: challenge.challenge,
    signature,
    publicKey
  }
}

/**
 * 验证DID格式
 */
export function validateDID(did: string): boolean {
  const didRegex = /^did:grets:[a-zA-Z]+:[a-f0-9]{16}$/
  return didRegex.test(did)
}

/**
 * 解析DID
 */
export function parseDID(did: string): { method: string; organization: string; identifier: string } | null {
  if (!validateDID(did)) return null
  
  const parts = did.split(':')
  return {
    method: parts[1],
    organization: parts[2],
    identifier: parts[3]
  }
} 