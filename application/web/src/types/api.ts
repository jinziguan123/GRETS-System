// 用户相关类型
export interface LoginData {
  citizenID: string
  password: string
  organization: string
}

export interface RegisterData {
  name: string
  citizenID: string
  password: string
  phone: string
  email: string
  publicKey: string
  role: string
  organization: string
  balance?: number
}

export interface UserInfo {
  id: number,
  citizenID: string
  name: string
  role: 'user' | 'admin'
  organization: string
  phone: string
  email: string
  balance?: number
  createdTime: string
  lastUpdateTime: string
  status: string
}

// 房产相关类型
export interface RealtyListParams {
  page?: number
  pageSize?: number
  status?: string
  realtyType?: string
}

export interface RealtyData {
  realtyCert: string
  realtyType: string
  currentOwnerCitizenID: string
  previousOwnersCitizenIDList?: string[]
  status: string
  contractID?: string
  description?: string
  price?: number
  houseType?: string
  images?: string[]
  isNewHouse?: boolean
}

// 交易相关类型
export interface TransactionListParams {
  page?: number
  pageSize?: number
  status?: string
}

export interface TransactionData {
  realtyCert: string
  sellerCitizenID: string
  buyerCitizenID: string
  contractID: string
  price: number
  tax: number
  expectedEndTime: string
}

// 支付相关类型
export interface PaymentListParams {
  page?: number
  pageSize?: number
  transactionUUID?: string
}

export interface PaymentData {
  amount: number
  fromCitizenID: string
  toCitizenID: string
  paymentType: string
}

// 合同相关类型
export interface ContractListParams {
  page?: number
  pageSize?: number
  status?: string
  contractType?: string
}

export interface ContractData {
  contractType: string
  docContent: string
  title: string
  sellerCitizenID?: string
  buyerCitizenID?: string
}

// 聊天相关类型
export interface ChatRoomData {
  realtyCert: string
  toUserID: string
}

export interface ChatMessageParams {
  page?: number
  pageSize?: number
}

export interface ChatMessageData {
  content: string
  messageType: 'text' | 'file'
  fileUrl?: string
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface Transaction {
  transactionUuid: string
  realtyCertHash: string
  sellerCitizenIDHash: string
  buyerCitizenIDHash: string
  contractUUID: string
  createTime: string
  updateTime: string
}

export interface QueryTransactionListResponse {
  total: number
  transactions: Transaction[]
}

export interface Challenge {
  did: string
  challenge: string
  nonce: string
  domain: string
}