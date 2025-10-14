import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Table, Tag } from 'antd'
import { useEffect, useState } from 'react'

type Defect = { id: number; title: string; status: string; priority: string; project_id: number }

export function DefectsPage() {
  const token = useAuthStore(s => s.token)
  const [items, setItems] = useState<Defect[]>([])
  useEffect(() => {
    if (!token) return
    api.get('/api/defects', { headers: { Authorization: `Bearer ${token}` } }).then(r => setItems(r.data))
  }, [token])
  return (
    <Table rowKey="id" dataSource={items} columns={[
      { title: 'Заголовок', dataIndex: 'title' },
      { title: 'Проект', dataIndex: 'project_id' },
      { title: 'Статус', dataIndex: 'status', render: (v: string) => <Tag color={v==='closed'?'green':v==='in_progress'?'blue':v==='on_review'?'orange':'default'}>{v}</Tag> },
      { title: 'Приоритет', dataIndex: 'priority' },
    ]} />
  )
}


