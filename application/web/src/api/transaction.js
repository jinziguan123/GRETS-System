import request from '@/utils/request.js'

export function queryTransactionList(data){
    return request({
        url: '/transaction/queryTransactionList',
        data: data,
        method: 'post',
    })
}

export function createTransaction(data){
    return request({
        url: '/transaction/createTransaction',
        method: 'post',
        data: data,
    })
}