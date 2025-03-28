import request from '@/utils/request'

// 用户登录
export function login(data) {
    return request({
        url: '/login',
        data: data,
        method: 'post',
    })
}

// 用户注册
export function register(data) {
    return request({
        url: '/register',
        data: data,
        method: 'post',
    })
}

// 获取用户信息
export function getUserInfo() {
    return request.get('/user/info')
}

// 更新用户信息
export function updateUserInfo(data) {
    return request.put('/user/info', data)
}