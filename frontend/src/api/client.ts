import axios, { type AxiosError, type InternalAxiosRequestConfig, type AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: attach auth token
instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('assassin_auth')
    if (token) {
      try {
        const parsed = JSON.parse(token)
        if (parsed.token) {
          config.headers.Authorization = `Bearer ${parsed.token}`
        }
      } catch {
        // ignore parse error
      }
    }
    return config
  },
  (error) => Promise.reject(error),
)

// Response interceptor: unwrap data, handle common errors
instance.interceptors.response.use(
  (response) => response.data,
  (error: AxiosError<{ message?: string }>) => {
    const status = error.response?.status
    const msg = error.response?.data?.message || error.message

    switch (status) {
      case 401:
        localStorage.removeItem('assassin_auth')
        router.push('/login')
        ElMessage.error('登录已过期，请重新登录')
        break
      case 403:
        ElMessage.error('无访问权限')
        break
      case 404:
        ElMessage.warning('资源不存在')
        break
      case 422:
        ElMessage.warning(msg || '参数校验失败')
        break
      case 500:
        ElMessage.error('服务器错误，请稍后重试')
        break
      default:
        if (status) {
          ElMessage.error(msg || '请求失败')
        }
    }
    return Promise.reject(error)
  },
)

// Typed wrapper: the interceptor unwraps response.data, so all calls return T directly
const http = {
  get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.get(url, config) as unknown as Promise<T>
  },
  post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.post(url, data, config) as unknown as Promise<T>
  },
  put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return instance.put(url, data, config) as unknown as Promise<T>
  },
  delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.delete(url, config) as unknown as Promise<T>
  },
}

export default http
