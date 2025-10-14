import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Card, Col, Row, Statistic } from 'antd'
import { useEffect, useState } from 'react'

export function DashboardPage() {
  const token = useAuthStore(s => s.token)
  const [data, setData] = useState<{ by_status?: Record<string, number>, by_project?: Record<string, number> }>({})
  useEffect(() => {
    if (!token) return
    api.get('/api/reports/analytics', { headers: { Authorization: `Bearer ${token}` } }).then(r => setData(r.data))
  }, [token])
  return (
    <Row gutter={16}>
      <Col span={8}><Card><Statistic title="По статусам" value={Object.values(data.by_status || {}).reduce((a, b) => a + b, 0)} /></Card></Col>
      <Col span={8}><Card><Statistic title="Проектов с дефектами" value={Object.keys(data.by_project || {}).length} /></Card></Col>
    </Row>
  )
}


