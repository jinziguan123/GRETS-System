import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    // 只返回数据部分
    return response.data
  },
  error => {
    // 处理错误响应
    const response = error.response
    if (response) {
      switch (response.status) {
        case 400:
          ElMessage.error(response.data.message || '请求参数错误')
          break
        case 401:
          ElMessage.error(response.data.message || '未授权，请重新登录')
          // 清除token并跳转到登录页
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          localStorage.removeItem('userRole')
          window.location.href = '/auth/login'
          break
        case 403:
          ElMessage.error(response.data.message || '禁止访问')
          break
        case 404:
          ElMessage.error(response.data.message || '请求的资源不存在')
          break
        case 500:
          ElMessage.error(response.data.message || '服务器内部错误')
          break
        default:
          ElMessage.error(response.data.message || `未知错误: ${response.status}`)
      }
    } else {
      ElMessage.error('网络错误，请检查您的网络连接')
    }
    return Promise.reject(error)
  }
)

// API模块
export default {
  // 用户相关API
  user: {
    login: (data) => api.post('/public/login', data),
    getProfile: () => api.get('/user/profile'),
    updateProfile: (data) => api.put('/user/profile', data),
    changePassword: (data) => api.put('/user/password', data)
  },
  
  // 房产相关API
  realty: {
    getList: (params) => api.get('/realty', { params }),
    getById: (id) => api.get(`/realty/${id}`),
    create: (data) => api.post('/realty', data),
    update: (id, data) => api.put(`/realty/${id}`, data)
  },
  
  // 交易相关API
  transaction: {
    getList: (params) => api.get('/realty/transaction', { params }),
    getById: (id) => api.get(`/realty/transaction/${id}`),
    create: (data) => api.post('/realty/transaction', data),
    audit: (id, data) => api.put(`/realty/transaction/${id}/audit`, data),
    complete: (id, data) => api.put(`/realty/transaction/${id}/complete`, data)
  },
  
  // 合同相关API
  contract: {
    getList: (params) => api.get('/realty/contract', { params }),
    getById: (id) => api.get(`/realty/contract/${id}`),
    create: (data) => api.post('/realty/contract', data),
    sign: (id, data) => api.put(`/realty/contract/${id}/sign`, data)
  },
  
  // 支付相关API
  payment: {
    getList: (params) => api.get('/realty/payment', { params }),
    getById: (id) => api.get(`/realty/payment/${id}`),
    create: (data) => api.post('/realty/payment', data),
    confirm: (id, data) => api.put(`/realty/payment/${id}/confirm`, data)
  },
  
  // 税费相关API
  tax: {
    getList: (params) => api.get('/realty/tax', { params }),
    getById: (id) => api.get(`/realty/tax/${id}`),
    create: (data) => api.post('/realty/tax', data),
    pay: (id, data) => api.put(`/realty/tax/${id}/pay`, data)
  },
  
  // 抵押贷款相关API
  mortgage: {
    getList: (params) => api.get('/realty/mortgage', { params }),
    getById: (id) => api.get(`/realty/mortgage/${id}`),
    create: (data) => api.post('/realty/mortgage', data),
    approve: (id, data) => api.put(`/realty/mortgage/${id}/approve`, data)
  },
  
  // 管理员相关API
  admin: {
    getUserList: (params) => api.get('/admin/user', { params }),
    getUserById: (id) => api.get(`/admin/user/${id}`),
    createUser: (data) => api.post('/admin/user', data),
    updateUser: (id, data) => api.put(`/admin/user/${id}`, data),
    deleteUser: (id) => api.delete(`/admin/user/${id}`)
  },
  
  // 文件相关API
  file: {
    upload: (formData) => api.post('/file/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }),
    getFile: (id) => api.get(`/file/${id}`)
  },
  
  // 系统信息API
  system: {
    getInfo: () => api.get('/public/system-info')
  }
} 