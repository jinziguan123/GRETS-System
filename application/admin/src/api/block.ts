import request from '@/utils/request'

/**
 * 查询区块列表
 * @param data 查询参数
 * @returns 区块列表数据
 */
export function queryBlockList(data: {
  pageSize: number
  pageNumber: number
  blockNumber?: string
  blockHash?: string
  provinceName?: string
  organization: string
}) {
  return request({
    url: '/blocks/queryBlockList',
    method: 'post',
    data
  })
}

/**
 * 获取区块详情
 * @param blockNumber 区块号
 * @returns 区块详情数据
 */
export function getBlockDetail(blockNumber: string) {
  return request({
    url: '/blocks/getBlockDetail',
    method: 'get',
    params: blockNumber
  })
}

/**
 * 查询区块交易列表
 * @param data 查询参数
 * @returns 区块交易列表数据
 */
export function queryBlockTransactionList(data: {
  blockNumber: string
  channelName: string
  organization: string
}) {
  return request({
    url: '/blocks/queryBlockTransactionList',
    method: 'post',
    data
  })
} 