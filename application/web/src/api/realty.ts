import request from '@/utils/request'

// 获取房产列表
export function queryRealtyList(data: {
  pageNumber?: number
  pageSize?: number
  status?: string
  realtyType?: string
  houseType?: string
  minPrice?: number
  maxPrice?: number
  minArea?: number
  maxArea?: number
  province?: string
  city?: string
  district?: string
  street?: string
  community?: string
  unit?: string
  floor?: string
  room?: string
}) {
  return request({
    url: '/realty/queryRealtyList',
    method: 'post',
    data
  })
}

// 获取房产详情
export function getRealtyDetail(id: string) {
  return request({
    url: `/realty/${id}`,
    method: 'get'
  })
}

export function QueryRealtyByOrganizationAndCitizenID(organization: string, citizenID: string) {
  return request({
    url: `/realty/QueryRealtyByOrganizationAndCitizenID`,
    method: 'get',
    params: {
      organization,
      citizenID
    }
  })
}

// 创建房产 (仅政府)
export function createRealty(data: {
  realtyCert: string
  realtyType: string
  currentOwnerCitizenID: string
  previousOwnersCitizenIDList?: string[]
  status: string
  contractID?: string
  description?: string
  price?: number
  area?: number
  houseType?: string
  images?: string[]
  isNewHouse?: boolean
  province?: string
  city?: string
  district?: string
  street?: string
  community?: string
  unit?: string
  floor?: string
  room?: string
  address?: string
}) {
  return request({
    url: '/realty/createRealty',
    method: 'post',
    data
  })
}

// 更新房产信息
export function updateRealty(realtyCert: string, data: {
  realtyType?: string
  status?: string
  description?: string
  price?: number
  houseType?: string
  images?: string[]
}) {
  return request({
    url: `/realty/${realtyCert}`,
    method: 'put',
    data
  })
}
