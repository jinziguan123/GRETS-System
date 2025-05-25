/**
 * DID密钥管理工具
 * 与后端crypto.go保持一致的ECDSA P-256密钥生成和处理
 */

// 密钥对接口
export interface KeyPair {
  publicKey: string    // 65字节十六进制字符串 (04前缀 + 32字节X + 32字节Y)
  privateKey: string   // 32字节十六进制字符串
  cryptoKeyPair?: CryptoKeyPair  // Web Crypto API密钥对（用于签名）
}

// 签名结果接口
export interface SignatureResult {
  signature: string   // 64字节十六进制字符串 (32字节r + 32字节s)
  message: string     // 原始消息
  publicKey: string   // 签名者公钥
}

/**
 * 生成ECDSA P-256密钥对
 * 与后端crypto.go的GenerateKeyPair()保持一致
 */
export async function generateKeyPair(): Promise<KeyPair> {
  try {
    // 使用Web Crypto API生成ECDSA P-256密钥对
    const cryptoKeyPair = await window.crypto.subtle.generateKey(
      {
        name: 'ECDSA',
        namedCurve: 'P-256'
      },
      true, // 可导出
      ['sign', 'verify']
    )

    // 导出公钥
    const publicKeyBuffer = await window.crypto.subtle.exportKey('raw', cryptoKeyPair.publicKey)
    const publicKeyArray = new Uint8Array(publicKeyBuffer)
    
    // 确保公钥是65字节（04前缀 + 32字节X + 32字节Y）
    if (publicKeyArray.length !== 65 || publicKeyArray[0] !== 0x04) {
      throw new Error('生成的公钥格式不正确')
    }
    
    const publicKeyHex = arrayBufferToHex(publicKeyBuffer)

    // 导出私钥
    const privateKeyBuffer = await window.crypto.subtle.exportKey('pkcs8', cryptoKeyPair.privateKey)
    const privateKeyHex = extractPrivateKeyFromPKCS8(privateKeyBuffer)

    return {
      publicKey: publicKeyHex,
      privateKey: privateKeyHex,
      cryptoKeyPair
    }
  } catch (error: any) {
    throw new Error(`密钥生成失败: ${error.message}`)
  }
}

/**
 * 从十六进制字符串恢复密钥对
 * 与后端crypto.go的HexToPrivateKey()和HexToPublicKey()保持一致
 */
export async function restoreKeyPair(privateKeyHex: string, publicKeyHex: string): Promise<KeyPair> {
  try {
    // 验证密钥格式
    if (!isValidPrivateKey(privateKeyHex)) {
      throw new Error('无效的私钥格式')
    }
    
    if (!isValidPublicKey(publicKeyHex)) {
      throw new Error('无效的公钥格式')
    }

    // 重新导入密钥到Web Crypto API
    const cryptoKeyPair = await importKeyPair(privateKeyHex, publicKeyHex)

    return {
      publicKey: publicKeyHex,
      privateKey: privateKeyHex,
      cryptoKeyPair
    }
  } catch (error: any) {
    throw new Error(`密钥恢复失败: ${error.message}`)
  }
}

/**
 * 使用私钥签名消息
 * 与后端crypto.go的SignMessage()保持一致
 */
export async function signMessage(keyPair: KeyPair, message: string): Promise<SignatureResult> {
  try {
    let cryptoKeyPair = keyPair.cryptoKeyPair
    
    if (!cryptoKeyPair) {
      // 如果没有cryptoKeyPair，尝试恢复
      const restoredKeyPair = await restoreKeyPair(keyPair.privateKey, keyPair.publicKey)
      cryptoKeyPair = restoredKeyPair.cryptoKeyPair!
    }

    // 计算消息的SHA-256哈希
    const messageBuffer = new TextEncoder().encode(message)
    const hashBuffer = await window.crypto.subtle.digest('SHA-256', messageBuffer)

    // 使用私钥签名哈希
    const signatureBuffer = await window.crypto.subtle.sign(
      {
        name: 'ECDSA',
        hash: 'SHA-256'
      },
      cryptoKeyPair.privateKey,
      hashBuffer
    )

    // 转换签名格式为r+s格式（64字节）
    const signature = formatSignature(signatureBuffer)

    return {
      signature,
      message,
      publicKey: keyPair.publicKey
    }
  } catch (error: any) {
    throw new Error(`签名失败: ${error.message}`)
  }
}

/**
 * 验证签名
 * 与后端crypto.go的VerifySignature()保持一致
 */
export async function verifySignature(publicKeyHex: string, message: string, signatureHex: string): Promise<boolean> {
  try {
    // 验证输入格式
    if (!isValidPublicKey(publicKeyHex)) {
      throw new Error('无效的公钥格式')
    }
    
    if (!isValidSignature(signatureHex)) {
      throw new Error('无效的签名格式')
    }

    // 导入公钥
    const publicKey = await importPublicKey(publicKeyHex)

    // 计算消息的SHA-256哈希
    const messageBuffer = new TextEncoder().encode(message)
    const hashBuffer = await window.crypto.subtle.digest('SHA-256', messageBuffer)

    // 转换签名格式
    const signatureBuffer = parseSignature(signatureHex)

    // 验证签名
    const isValid = await window.crypto.subtle.verify(
      {
        name: 'ECDSA',
        hash: 'SHA-256'
      },
      publicKey,
      signatureBuffer,
      hashBuffer
    )

    return isValid
  } catch (error: any) {
    console.error('签名验证失败:', error)
    return false
  }
}

/**
 * 生成SHA-256哈希
 * 与后端crypto.go的GenerateHash()保持一致
 */
export async function generateHash(data: string): Promise<string> {
  const buffer = new TextEncoder().encode(data)
  const hashBuffer = await window.crypto.subtle.digest('SHA-256', buffer)
  return arrayBufferToHex(hashBuffer)
}

/**
 * 生成DID标识符
 */
export function generateDID(organization: string, identifier?: string): string {
  const id = identifier || generateRandomIdentifier()
  return `did:grets:${organization}:${id}`
}

/**
 * 解析DID标识符
 */
export function parseDID(did: string): { method: string; organization: string; identifier: string } | null {
  const parts = did.split(':')
  if (parts.length !== 4 || parts[0] !== 'did' || parts[1] !== 'grets') {
    return null
  }
  
  return {
    method: parts[1],
    organization: parts[2],
    identifier: parts[3]
  }
}

// ==================== 内部工具函数 ====================

/**
 * ArrayBuffer转十六进制字符串
 */
function arrayBufferToHex(buffer: ArrayBuffer): string {
  const byteArray = new Uint8Array(buffer)
  return Array.from(byteArray, byte => byte.toString(16).padStart(2, '0')).join('')
}

/**
 * 十六进制字符串转ArrayBuffer
 */
function hexToArrayBuffer(hex: string): ArrayBuffer {
  const bytes = new Uint8Array(hex.length / 2)
  for (let i = 0; i < hex.length; i += 2) {
    bytes[i / 2] = parseInt(hex.substr(i, 2), 16)
  }
  return bytes.buffer
}

/**
 * 从PKCS8格式中提取32字节私钥
 */
function extractPrivateKeyFromPKCS8(pkcs8Buffer: ArrayBuffer): string {
  const pkcs8Array = new Uint8Array(pkcs8Buffer)
  
  // PKCS8格式中，私钥通常在最后32字节
  // 这是一个简化的提取方法，实际应该解析ASN.1结构
  const privateKeyBytes = pkcs8Array.slice(-32)
  
  return Array.from(privateKeyBytes, byte => byte.toString(16).padStart(2, '0')).join('')
}

/**
 * 将私钥转换为PKCS8格式
 */
function privateKeyToPKCS8(privateKeyHex: string): ArrayBuffer {
  // 这是一个简化的PKCS8包装
  // 实际应该构建完整的ASN.1结构
  const privateKeyBytes = hexToArrayBuffer(privateKeyHex)
  
  // P-256私钥的PKCS8前缀（简化版）
  const prefix = new Uint8Array([
    0x30, 0x81, 0x87, 0x02, 0x01, 0x00, 0x30, 0x13, 0x06, 0x07, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x02,
    0x01, 0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07, 0x04, 0x6d, 0x30, 0x6b, 0x02,
    0x01, 0x01, 0x04, 0x20
  ])
  
  const suffix = new Uint8Array([
    0xa1, 0x44, 0x03, 0x42, 0x00
  ])
  
  const result = new Uint8Array(prefix.length + 32 + suffix.length + 64)
  result.set(prefix, 0)
  result.set(new Uint8Array(privateKeyBytes), prefix.length)
  result.set(suffix, prefix.length + 32)
  
  return result.buffer
}

/**
 * 导入密钥对到Web Crypto API
 */
async function importKeyPair(privateKeyHex: string, publicKeyHex: string): Promise<CryptoKeyPair> {
  // 导入私钥
  const privateKeyPKCS8 = privateKeyToPKCS8(privateKeyHex)
  const privateKey = await window.crypto.subtle.importKey(
    'pkcs8',
    privateKeyPKCS8,
    {
      name: 'ECDSA',
      namedCurve: 'P-256'
    },
    true,
    ['sign']
  )

  // 导入公钥
  const publicKey = await importPublicKey(publicKeyHex)

  return { privateKey, publicKey }
}

/**
 * 导入公钥到Web Crypto API
 */
async function importPublicKey(publicKeyHex: string): Promise<CryptoKey> {
  const publicKeyBuffer = hexToArrayBuffer(publicKeyHex)
  
  return await window.crypto.subtle.importKey(
    'raw',
    publicKeyBuffer,
    {
      name: 'ECDSA',
      namedCurve: 'P-256'
    },
    true,
    ['verify']
  )
}

/**
 * 格式化签名为r+s格式（64字节）
 */
function formatSignature(signatureBuffer: ArrayBuffer): string {
  // Web Crypto API返回的是DER格式，需要转换为r+s格式
  const signature = new Uint8Array(signatureBuffer)
  
  // 简化的DER解析（实际应该完整解析ASN.1）
  // DER格式: 0x30 [length] 0x02 [r-length] [r] 0x02 [s-length] [s]
  let offset = 2 // 跳过0x30和总长度
  
  // 解析r
  offset++ // 跳过0x02
  const rLength = signature[offset++]
  let r = signature.slice(offset, offset + rLength)
  offset += rLength
  
  // 解析s
  offset++ // 跳过0x02
  const sLength = signature[offset++]
  let s = signature.slice(offset, offset + sLength)
  
  // 确保r和s都是32字节
  if (r.length > 32) r = r.slice(-32)
  if (s.length > 32) s = s.slice(-32)
  
  const rPadded = new Uint8Array(32)
  const sPadded = new Uint8Array(32)
  rPadded.set(r, 32 - r.length)
  sPadded.set(s, 32 - s.length)
  
  const result = new Uint8Array(64)
  result.set(rPadded, 0)
  result.set(sPadded, 32)
  
  return arrayBufferToHex(result.buffer)
}

/**
 * 解析r+s格式签名为DER格式
 */
function parseSignature(signatureHex: string): ArrayBuffer {
  const signatureBytes = hexToArrayBuffer(signatureHex)
  const r = new Uint8Array(signatureBytes, 0, 32)
  const s = new Uint8Array(signatureBytes, 32, 32)
  
  // 构建DER格式
  // 移除前导零
  const rTrimmed = trimLeadingZeros(r)
  const sTrimmed = trimLeadingZeros(s)
  
  // 如果最高位是1，需要添加0x00前缀
  const rFinal = (rTrimmed[0] & 0x80) ? new Uint8Array([0x00, ...rTrimmed]) : rTrimmed
  const sFinal = (sTrimmed[0] & 0x80) ? new Uint8Array([0x00, ...sTrimmed]) : sTrimmed
  
  const totalLength = 2 + rFinal.length + 2 + sFinal.length
  const result = new Uint8Array(2 + totalLength)
  
  let offset = 0
  result[offset++] = 0x30 // SEQUENCE
  result[offset++] = totalLength
  result[offset++] = 0x02 // INTEGER
  result[offset++] = rFinal.length
  result.set(rFinal, offset)
  offset += rFinal.length
  result[offset++] = 0x02 // INTEGER
  result[offset++] = sFinal.length
  result.set(sFinal, offset)
  
  return result.buffer
}

/**
 * 移除字节数组的前导零
 */
function trimLeadingZeros(bytes: Uint8Array): Uint8Array {
  let start = 0
  while (start < bytes.length && bytes[start] === 0) {
    start++
  }
  return bytes.slice(start)
}

/**
 * 验证私钥格式（32字节十六进制）
 */
function isValidPrivateKey(privateKeyHex: string): boolean {
  return /^[0-9a-fA-F]{64}$/.test(privateKeyHex)
}

/**
 * 验证公钥格式（65字节十六进制，04前缀）
 */
function isValidPublicKey(publicKeyHex: string): boolean {
  return /^04[0-9a-fA-F]{128}$/.test(publicKeyHex)
}

/**
 * 验证签名格式（64字节十六进制）
 */
function isValidSignature(signatureHex: string): boolean {
  return /^[0-9a-fA-F]{128}$/.test(signatureHex)
}

/**
 * 生成随机标识符
 */
function generateRandomIdentifier(): string {
  const array = new Uint8Array(16)
  window.crypto.getRandomValues(array)
  return arrayBufferToHex(array.buffer)
} 