import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const service = axios.create({
  headers: {
    'Content-Type': 'application/json',
  },
  baseURL: 'http://localhost:8080/api/v1',
  timeout: 60000,
})

// 请求拦截器
service.interceptors.request.use(
    config => {
      // 从pinia获取token
      const userStore = useUserStore()
      if (userStore.token) {
        config.headers['Authorization'] = `${userStore.token}`
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

export default service
