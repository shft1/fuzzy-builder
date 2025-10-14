import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button } from 'antd'

export function ReportsPage() {
  const token = useAuthStore(s => s.token)
  const onExport = async () => {
    const r = await api.get('/api/reports/defects', { responseType: 'blob', headers: { Authorization: `Bearer ${token}` } })
    const url = URL.createObjectURL(new Blob([r.data], { type: 'text/csv' }))
    const a = document.createElement('a')
    a.href = url
    a.download = 'defects.csv'
    a.click()
    URL.revokeObjectURL(url)
  }
  return <Button type="primary" onClick={onExport} disabled={!token}>Экспорт CSV</Button>
}


