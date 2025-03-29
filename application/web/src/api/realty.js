import request from '@/utils/request.js'

export function createRealty(data){
    return request({
        url: '/realty/createRealty',
        data: data,
        method: 'post',
    })
}

export function queryRealtyList(data){
    return request({
        url: '/realty/queryRealtyList',
        method: 'POST',
        data: data
    })
}

export function getRealtyByID(id){
    return request({
        url: `/realty/${id}`,
        method: 'get'
    })
}

// 更新房产信息
export function updateRealty(id, data){
    return request({
        url: `/realty/${id}`,
        method: 'put',
        data: data
    })
}

// 查询交易列表
export function queryTransactionList(data){
    return request({
        url: '/transaction/list',
        method: 'post',
        data: data
    })
}