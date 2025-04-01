import request from '@/utils/request'

// 获取支付列表
export function getPaymentList(params: {
  page?: number
  pageSize?: number
  transactionUUID?: string
}) {
  return request({
    url: '/api/v1/payments',
    method: 'get',
    params
  })
}

// 获取支付详情
export function getPaymentDetail(paymentUUID: string) {
  return request({
    url: `/api/v1/payments/${paymentUUID}`,
    method: 'get'
  })
}

// 创建支付
export function createPayment(data: {
  amount: number
  fromCitizenID: string
  toCitizenID: string
  paymentType: string
}) {
  return request({
    url: '/api/v1/payments',
    method: 'post',
    data
  })
}

// 为交易创建支付
export function payForTransaction(data: {
  transactionUUID: string
  amount: number
  fromCitizenID: string
  toCitizenID: string
  paymentType: string
}) {
  return request({
    url: '/api/v1/payments/transaction',
    method: 'post',
    data
  })
} 