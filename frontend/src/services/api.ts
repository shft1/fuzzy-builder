import { useAuthStore } from '@/stores/auth'
import { notification } from 'antd'
import axios from 'axios'

const BASE = (import.meta as any).env?.VITE_API_BASE || '/'
export const api = axios.create({ baseURL: BASE })

api.interceptors.request.use(config => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers = config.headers || {}
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(undefined, (error) => {
  if (error?.response?.status === 401) {
    useAuthStore.getState().logout()
  }
  // basic notification surface
  const msg = error?.response?.data?.error || error?.message || 'Ошибка запроса'
  notification.error({ message: 'Ошибка', description: msg })
  return Promise.reject(error)
})


