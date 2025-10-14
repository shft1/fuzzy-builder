import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

export const api = axios.create({ baseURL: '/' })

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
  return Promise.reject(error)
})


