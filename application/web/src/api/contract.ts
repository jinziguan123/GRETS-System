import request from '@/utils/request'

// 获取合同列表
export function getContractList(params: {
  page?: number
  pageSize?: number
  status?: string
  contractType?: string
}) {
  return request({
    url: '/api/v1/contracts',
    method: 'get',
    params
  })
}

// 获取合同详情
export function getContractDetail(contractID: string) {
  return request({
    url: `/api/v1/contracts/${contractID}`,
    method: 'get'
  })
}

// 创建合同
export function createContract(data: {
  contractType: string
  docContent: string
  title: string
  sellerCitizenID?: string
  buyerCitizenID?: string
}) {
  return request({
    url: '/api/v1/contracts',
    method: 'post',
    data
  })
}

// 更新合同状态
export function updateContractStatus(contractID: string, status: string) {
  return request({
    url: `/api/v1/contracts/${contractID}/status`,
    method: 'put',
    data: { status }
  })
}

// 上传合同文件
export function uploadContractFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/api/v1/contracts/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
} 