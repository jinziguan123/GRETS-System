import request from '@/utils/request'

// 获取交易列表
export function getTransactionList(params: {
  page?: number
  pageSize?: number
  status?: string
}) {
  return request({
    url: '/transactions',
    method: 'get',
    params
  })
}

// 获取交易详情
export function getTransactionDetail(transactionUUID: string) {
  return request({
    url: `/transactions/${transactionUUID}`,
    method: 'get'
  })
}

// 创建交易
export function createTransaction(data: {
  realtyCert: string
  sellerCitizenID: string
  buyerCitizenID: string
  contractID: string
  price: number
  tax: number
  expectedEndTime: string
}) {
  return request({
    url: '/transactions/createTransaction',
    method: 'post',
    data
  })
}

// 检查交易 (审核状态更改)
export function checkTransaction(transactionUUID: string, status: string) {
  return request({
    url: `/transactions/${transactionUUID}/check`,
    method: 'post',
    data: { status }
  })
}

// 完成交易
export function completeTransaction(data: {
  transactionUUID: string
}) {
  return request({
    url: `/transactions/completeTransaction`,
    method: 'post',
    data
  })
}

// 交易条件查询
export function queryTransactionList(data: {
  buyerCitizenID?: string
  sellerCitizenID?: string
  realtyCert?: string
  status?: string
  pageSize: number
  pageNumber: number
}){
  return request({
    url: '/transactions/queryTransactionList',
    method: 'post',
    data
  })
}

export function updateTransaction(data: {
  transactionUUID: string,
  status: string
}){
  return request({
    url: '/transactions/updateTransaction',
    method: 'post',
    data
  })
}