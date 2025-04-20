import request from '@/utils/request'

// 获取支付列表
export function getPaymentList(data: {
  pageNumber?: number
  pageSize?: number
  transactionUUID?: string
  paymentType?: string
  payerCitizenID?: string
  receiverCitizenID?: string
}) {
  return request({
    url: '/payments/queryPaymentList',
    method: 'post',
    data
  })
}

// 获取支付详情
export function getPaymentDetail(paymentUUID: string) {
  return request({
    url: `/payments/${paymentUUID}`,
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
    url: '/payments',
    method: 'post',
    data
  })
}

// 为交易创建支付
export function payForTransaction(data: {
  transactionUUID: string
  amount: number
  payerCitizenID: string
  payerOrganization: string
  receiverCitizenIDHash: string
  receiverOrganization: string
  paymentType: string
  remarks?: string
}) {
  return request({
    url: '/payments/payForTransaction',
    method: 'post',
    data
  })
}

export function getTotalPaymentAmount(){
  return request({
    url: '/payments/getTotalPaymentAmount',
    method: 'get'
  })
}