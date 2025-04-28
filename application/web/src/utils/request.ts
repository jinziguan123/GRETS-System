import axios from 'axios'
import type {AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig, AxiosResponse} from 'axios'
import {ElMessage} from 'element-plus'
import type {ApiResponse} from '@/types/api'

// 创建axios实例
const service: AxiosInstance = axios.create({
    baseURL: 'http://localhost:8080/api/v1',
    timeout: 15000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
        // 在发送请求之前做些什么
        const token = localStorage.getItem('token')
        if (token && config.headers) {
            config.headers['Authorization'] = token
        }
        return config
    },
    (error: any) => {
        // 对请求错误做些什么
        console.error('Request error:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    (response: AxiosResponse): Promise<any> => {
        // 对响应数据做点什么
        const res = response.data

        // 如果是文件下载等二进制数据直接返回
        if (response.config.responseType === 'blob') {
            return Promise.resolve(response)
        }

        // 根据后端的响应格式调整这里的条件
        if (res.code && res.code !== 200) {
            ElMessage({
                message: res.message || '请求错误',
                type: 'error',
                duration: 5 * 1000
            })

            // 401: 未登录或token过期
            if (res.code === 401) {
                // 可以跳转到登录页或其他处理
                localStorage.removeItem('token')
                window.location.href = '/login'
            }
            return Promise.reject(new Error(res.message || '请求错误'))
        } else {
            return Promise.resolve(res)
        }
    },
    (error: any) => {
        // 对响应错误做点什么
        console.error('Response error:', error)
        ElMessage({
            message: error.response.data.message || '请求失败',
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
)

// 封装请求方法
const request = async <T = any>(config: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    try {
        const response = await service(config as InternalAxiosRequestConfig)
        return response.data as ApiResponse<T>
    } catch (error) {
        return Promise.reject(error)
    }
}

export default request 