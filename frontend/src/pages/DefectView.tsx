import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Card, Descriptions, Tag } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'

export function DefectViewPage() {
  const { id } = useParams()
  const token = useAuthStore(s => s.token)
  const [d, setD] = useState<any>(null)
  useEffect(() => {
    if (!id || !token) return
    api.get(`/api/defects/${id}`, { headers: { Authorization: `Bearer ${token}` } }).then(r => setD(r.data))
  }, [id, token])
  if (!d) return null
  return (
    <Card title={`Дефект #${d.id}`}>
      <Descriptions bordered column={1}>
        <Descriptions.Item label="Заголовок">{d.title}</Descriptions.Item>
        <Descriptions.Item label="Описание">{d.description}</Descriptions.Item>
        <Descriptions.Item label="Проект">{d.project_id}</Descriptions.Item>
        <Descriptions.Item label="Статус"><Tag color={d.status==='closed'?'green':d.status==='in_progress'?'blue':d.status==='on_review'?'orange':'default'}>{d.status}</Tag></Descriptions.Item>
        <Descriptions.Item label="Приоритет">{d.priority}</Descriptions.Item>
      </Descriptions>
    </Card>
  )
}


