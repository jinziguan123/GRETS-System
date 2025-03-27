import request from '@/utils/request'

export function createRealty(data){
    return request.post('/realty/createRealty', data)
}