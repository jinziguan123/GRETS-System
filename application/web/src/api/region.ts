import request from '@/utils/request'

// 获取所有省份
export function getProvinces() {
  return request({
    url: '/regions/provinces',
    method: 'get'
  })
}

// 获取指定省份的城市
export function getCitiesByProvince(provinceCode: string) {
  return request({
    url: `/regions/cities/${provinceCode}`,
    method: 'get'
  })
} 