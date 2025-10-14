import { useAuthStore } from '@/stores/auth'
import { Navigate, useLocation } from 'react-router-dom'

export function ProtectedRoute({ children }: { children: JSX.Element }) {
  const token = useAuthStore(s => s.token)
  const location = useLocation()
  if (!token) return <Navigate to="/login" replace state={{ from: location }} />
  return children
}


