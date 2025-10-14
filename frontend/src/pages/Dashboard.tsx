import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Card, Col, DatePicker, Row, Select, Space, Statistic } from 'antd'
import { BarElement, CategoryScale, Chart, Legend, LinearScale, Tooltip } from 'chart.js'
import { useEffect, useState } from 'react'
import { Bar } from 'react-chartjs-2'
Chart.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend)

export function DashboardPage() {
  const token = useAuthStore(s => s.token)
  const [data, setData] = useState<{ by_status?: Record<string, number>, by_project?: Record<string, number> }>({})
  const [project, setProject] = useState<string | undefined>()
  const [range, setRange] = useState<any>(null)
  useEffect(() => {
    if (!token) return
    // Note: backend фильтры можно добавить позже; пока фильтры локальные (демо)
    api.get('/api/reports/analytics', { headers: { Authorization: `Bearer ${token}` } }).then(r => setData(r.data))
  }, [token, project, range])
  return (
    <>
      <Space style={{ marginBottom: 16 }}>
        <Select placeholder="Проект" allowClear style={{ width: 220 }} value={project} onChange={setProject}
          options={Object.keys(data.by_project || {}).map(k => ({ value: k, label: k }))}
        />
        <DatePicker.RangePicker onChange={setRange} />
      </Space>
      <Row gutter={16}>
        <Col span={8}><Card><Statistic title="Всего дефектов" value={Object.values(data.by_status || {}).reduce((a, b) => a + b, 0)} /></Card></Col>
        <Col span={8}><Card><Statistic title="Проектов с дефектами" value={Object.keys(data.by_project || {}).length} /></Card></Col>
      </Row>
      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card title="Дефекты по статусам">
            <Bar data={{
              labels: Object.keys(data.by_status || {}),
              datasets: [{ label: 'Статусы', data: Object.values(data.by_status || {}), backgroundColor: '#1677ff' }]
            }} options={{ plugins: { legend: { display: false } } }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Дефекты по проектам">
            <Bar data={{
              labels: Object.keys(data.by_project || {}),
              datasets: [{ label: 'Проекты', data: Object.values(data.by_project || {}), backgroundColor: '#52c41a' }]
            }} options={{ plugins: { legend: { display: false } }, scales: { x: { ticks: { autoSkip: true, maxTicksLimit: 10 } } } }} />
          </Card>
        </Col>
      </Row>
    </>
  )
}


