import request from '@/utils/request'

// 用户登录
export function login(data) {
  return request.post('/login', data)
}

// 用户注册
export function register(data) {
  return request.post('/register', data)
}

// 获取用户信息
export function getUserInfo() {
  return request.get('/user/info')
}

// 更新用户信息
export function updateUserInfo(data) {
  return request.put('/user/info', data)
}