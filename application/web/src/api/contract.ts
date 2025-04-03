import request from '@/utils/request'

// 获取合同列表
export function queryContractList(data: {
  contractUUID?: string
  docHash?: string
  contractType?: string
  creatorCitizenID?: string
  status?: string
  pageSize?: number
  pageNumber?: number
}) {
  return request({
    url: '/contracts/queryContractList',
    method: 'post',
    data
  })
}

// 获取合同详情
export function getContractDetail(id: string) {
  return request({
    url: `/contracts/${id}`,
    method: 'get'
  })
}

// 创建合同
export function createContract(data: {
  title: string
  content: string
  contractType: string
  creatorCitizenID: string
}) {
  return request({
    url: '/contracts/createContract',
    method: 'post',
    data
  })
}

// 签署合同
export function signContract(id: string, data: {
  signerType: string
}) {
  return request({
    url: `/contracts/${id}/sign`,
    method: 'post',
    data
  })
}

// 审核合同
export function auditContract(id: string, data: {
  result: string
  comments?: string
  revisionRequirements?: string
  rejectionReason?: string
}) {
  return request({
    url: `/contracts/${id}/audit`,
    method: 'post',
    data
  })
}

// 更新合同状态
export function updateContractStatus(contractID: string, status: string) {
  return request({
    url: `/contracts/${contractID}/status`,
    method: 'put',
    data: { status }
  })
}

// 上传合同文件
export function uploadContractFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/contracts/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getContractByID(id: number) {
  return request({
    url: `/contracts/${id}`,
    method: 'get'
  })
}

export function getContractByUUID(uuid: string) {
  return request({
    url: `/contracts/getContractByUUID/${uuid}`,
    method: 'get'
  })
}