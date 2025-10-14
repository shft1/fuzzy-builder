import { useAuthStore } from '@/stores/auth'
import { ReactNode } from 'react'

export function ManagerOnly({ children }: { children: ReactNode }) {
  const role = useAuthStore(s => s.user?.role)
  if (role !== 'manager') return null
  return <>{children}</>
}


