import { useAuthStore } from '@/stores/auth'
import { Button } from 'antd'

export function ReportsPage() {
  const token = useAuthStore(s => s.token)
  const onExport = () => {
    const url = `/api/reports/defects`
    const a = document.createElement('a')
    a.href = url
    a.download = 'defects.csv'
    a.target = '_blank'
    a.rel = 'noopener'
    a.click()
  }
  return <Button type="primary" onClick={onExport} disabled={!token}>Экспорт CSV</Button>
}


