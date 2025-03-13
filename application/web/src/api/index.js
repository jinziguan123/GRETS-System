import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 从pinia获取token
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers['Authorization'] = `Bearer ${userStore.token}`
    }
    return config
  },
  error => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data

    // 如果返回的状态码不是200，说明接口请求有误
    if (res.code !== 200 && res.code !== 0) {
      ElMessage({
        message: res.message || '系统错误',
        type: 'error',
        duration: 5 * 1000
      })

      // 特定错误码处理
      if (res.code === 401) {
        // token过期或未登录
        const userStore = useUserStore()
        userStore.resetToken()
        location.reload()
      }
      return Promise.reject(new Error(res.message || '系统错误'))
    } else {
      return res
    }
  },
  error => {
    console.error('响应错误:', error)
    const { response } = error
    let message = '网络错误，请稍后再试'
    
    if (response && response.data) {
      message = response.data.message || `错误: ${response.status}`
    }
    
    ElMessage({
      message,
      type: 'error',
      duration: 5 * 1000
    })
    
    // 特定状态码处理
    if (response && response.status === 401) {
      // token过期或未登录
      const userStore = useUserStore()
      userStore.resetToken()
      location.reload()
    }
    
    return Promise.reject(error)
  }
)

// 封装GET请求
const get = (url, params) => {
  return service({
    url,
    method: 'get',
    params
  })
}

// 封装POST请求
const post = (url, data) => {
  return service({
    url,
    method: 'post',
    data
  })
}

// 封装PUT请求
const put = (url, data) => {
  return service({
    url,
    method: 'put',
    data
  })
}

// 封装DELETE请求
const del = (url, params) => {
  return service({
    url,
    method: 'delete',
    params
  })
}

// 用户相关API
export const userApi = {
  // 用户登录
  login: (data) => post('/user/login', data),
  
  // 用户注册
  register: (data) => post('/user/register', data),
  
  // 获取用户信息
  getUserInfo: () => get('/user/info'),
  
  // 获取用户详细资料
  getUserProfile: () => get('/user/profile'),
  
  // 更新用户资料
  updateProfile: (data) => put('/user/profile', data),
  
  // 修改密码
  changePassword: (data) => post('/user/change-password', data),
  
  // 退出登录
  logout: () => post('/user/logout')
}

// 房产相关API
export const realtyApi = {
  // 获取房产列表
  getRealtyList: (params) => get('/realty/list', params),
  
  // 获取房产详情
  getRealtyDetail: (id) => get(`/realty/${id}`),
  
  // 创建房产
  createRealty: (data) => post('/realty', data),
  
  // 更新房产信息
  updateRealty: (id, data) => put(`/realty/${id}`, data),
  
  // 删除房产
  deleteRealty: (id) => del(`/realty/${id}`),
  
  // 获取房产证书
  getCertificate: (id) => get(`/realty/${id}/certificate`),
  
  // 上传房产相关文件
  uploadDocument: (id, data) => post(`/realty/${id}/document`, data)
}

// 交易相关API
export const transactionApi = {
  // 获取交易列表
  getTransactionList: (params) => get('/transaction/list', params),
  
  // 获取交易详情
  getTransactionDetail: (id) => get(`/transaction/${id}`),
  
  // 创建交易
  createTransaction: (data) => post('/transaction', data),
  
  // 更新交易信息
  updateTransaction: (id, data) => put(`/transaction/${id}`, data),
  
  // 取消交易
  cancelTransaction: (id) => post(`/transaction/${id}/cancel`),
  
  // 确认交易
  confirmTransaction: (id) => post(`/transaction/${id}/confirm`),
  
  // 查询交易历史
  getTransactionHistory: (id) => get(`/transaction/${id}/history`)
}

// 合同相关API
export const contractApi = {
  // 获取合同列表
  getContractList: (params) => get('/contract/list', params),
  
  // 获取合同详情
  getContractDetail: (id) => get(`/contract/${id}`),
  
  // 创建合同
  createContract: (data) => post('/contract', data),
  
  // 更新合同信息
  updateContract: (id, data) => put(`/contract/${id}`, data),
  
  // 取消合同
  cancelContract: (id) => post(`/contract/${id}/cancel`),
  
  // 签署合同
  signContract: (id, data) => post(`/contract/${id}/sign`, data),
  
  // 获取合同签署历史
  getContractSignHistory: (id) => get(`/contract/${id}/sign-history`)
}

// 支付相关API
export const paymentApi = {
  // 获取支付列表
  getPaymentList: (params) => get('/payment/list', params),
  
  // 获取支付详情
  getPaymentDetail: (id) => get(`/payment/${id}`),
  
  // 创建支付
  createPayment: (data) => post('/payment', data),
  
  // 更新支付信息
  updatePayment: (id, data) => put(`/payment/${id}`, data),
  
  // 取消支付
  cancelPayment: (id) => post(`/payment/${id}/cancel`),
  
  // 确认支付
  confirmPayment: (id) => post(`/payment/${id}/confirm`),
  
  // 获取支付记录
  getPaymentRecords: (params) => get('/payment/records', params)
}

// 税费相关API
export const taxApi = {
  // 获取税费列表
  getTaxList: (params) => get('/tax/list', params),
  
  // 获取税费详情
  getTaxDetail: (id) => get(`/tax/${id}`),
  
  // 创建税费
  createTax: (data) => post('/tax', data),
  
  // 更新税费信息
  updateTax: (id, data) => put(`/tax/${id}`, data),
  
  // 计算交易税费
  calculateTax: (data) => post('/tax/calculate', data),
  
  // 支付税费
  payTax: (id) => post(`/tax/${id}/pay`),
  
  // 获取税费记录
  getTaxRecords: (params) => get('/tax/records', params)
}

// 抵押贷款相关API
export const mortgageApi = {
  // 获取贷款列表
  getMortgageList: (params) => get('/mortgage/list', params),
  
  // 获取贷款详情
  getMortgageDetail: (id) => get(`/mortgage/${id}`),
  
  // 创建贷款申请
  createMortgage: (data) => post('/mortgage', data),
  
  // 更新贷款信息
  updateMortgage: (id, data) => put(`/mortgage/${id}`, data),
  
  // 取消贷款申请
  cancelMortgage: (id) => post(`/mortgage/${id}/cancel`),
  
  // 审批贷款
  approveMortgage: (id, data) => post(`/mortgage/${id}/approve`, data),
  
  // 拒绝贷款
  rejectMortgage: (id, data) => post(`/mortgage/${id}/reject`, data),
  
  // 获取还款计划
  getRepaymentPlan: (id) => get(`/mortgage/${id}/plan`),
  
  // 记录还款
  recordRepayment: (id, data) => post(`/mortgage/${id}/repay`, data)
}

// 文件相关API
export const fileApi = {
  // 上传文件
  uploadFile: (data, config) => post('/file/upload', data, config),
  
  // 获取文件
  getFile: (id) => get(`/file/${id}`, { responseType: 'blob' }),
  
  // 删除文件
  deleteFile: (id) => del(`/file/${id}`),
}

// 仪表盘相关API
export const dashboardApi = {
  // 获取仪表盘数据
  getDashboardData: () => get('/dashboard'),
  
  // 获取交易统计
  getTransactionStats: (params) => get('/dashboard/transaction-stats', params),
  
  // 获取房产统计
  getRealtyStats: () => get('/dashboard/realty-stats'),
  
  // 获取用户统计
  getUserStats: () => get('/dashboard/user-stats')
}

// 管理员相关API
export const adminApi = {
  // 获取用户列表
  getUserList: (params) => get('/admin/users', params),
  
  // 创建用户
  createUser: (data) => post('/admin/users', data),
  
  // 更新用户信息
  updateUser: (id, data) => put(`/admin/users/${id}`, data),
  
  // 删除用户
  deleteUser: (id) => del(`/admin/users/${id}`),
  
  // 获取系统日志
  getSystemLogs: (params) => get('/admin/logs', params),
  
  // 获取系统设置
  getSettings: () => get('/admin/settings'),
  
  // 更新系统设置
  updateSettings: (data) => put('/admin/settings', data)
}

export default {
  userApi,
  realtyApi,
  transactionApi,
  contractApi,
  paymentApi,
  taxApi,
  mortgageApi,
  fileApi,
  dashboardApi,
  adminApi
} 