import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Table } from 'antd'
import { useEffect, useState } from 'react'

type Project = { id: number; name: string; description: string }

export function ProjectsPage() {
  const token = useAuthStore(s => s.token)
  const [items, setItems] = useState<Project[]>([])
  useEffect(() => {
    if (!token) return
    api.get('/api/projects', { headers: { Authorization: `Bearer ${token}` } }).then(r => setItems(r.data))
  }, [token])
  return (
    <Table rowKey="id" dataSource={items} columns={[
      { title: 'Название', dataIndex: 'name' },
      { title: 'Описание', dataIndex: 'description' },
    ]} />
  )
}


