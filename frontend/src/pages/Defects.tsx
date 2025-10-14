import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Select, Space, Table, Tag } from 'antd'
import { useEffect, useMemo, useState } from 'react'

type Defect = { id: number; title: string; status: string; priority: string; project_id: number }

export function DefectsPage() {
  const token = useAuthStore(s => s.token)
  const [items, setItems] = useState<Defect[]>([])
  const [status, setStatus] = useState<string | undefined>()
  const [priority, setPriority] = useState<string | undefined>()
  const query = useMemo(() => {
    const p = new URLSearchParams()
    if (status) p.set('status', status)
    if (priority) p.set('priority', priority)
    return p.toString()
  }, [status, priority])
  useEffect(() => {
    if (!token) return
    const url = '/api/defects' + (query ? `?${query}` : '')
    api.get(url, { headers: { Authorization: `Bearer ${token}` } }).then(r => setItems(r.data))
  }, [token, query])
  return (
    <>
      <Space style={{ marginBottom: 16 }} wrap>
        <Select placeholder="Статус" allowClear style={{ width: 200 }} value={status} onChange={setStatus}
          options={[
            { value: 'new', label: 'Новая' },
            { value: 'in_progress', label: 'В работе' },
            { value: 'on_review', label: 'На проверке' },
            { value: 'closed', label: 'Закрыта' },
          ]}
        />
        <Select placeholder="Приоритет" allowClear style={{ width: 200 }} value={priority} onChange={setPriority}
          options={[
            { value: 'low', label: 'Низкий' },
            { value: 'medium', label: 'Средний' },
            { value: 'high', label: 'Высокий' },
          ]}
        />
        <Button onClick={() => { setStatus(undefined); setPriority(undefined) }}>Сбросить</Button>
      </Space>
      <Table rowKey="id" dataSource={items} columns={[
        { title: 'Заголовок', dataIndex: 'title' },
        { title: 'Проект', dataIndex: 'project_id' },
        { title: 'Статус', dataIndex: 'status', render: (v: string) => <Tag color={v==='closed'?'green':v==='in_progress'?'blue':v==='on_review'?'orange':'default'}>{v}</Tag> },
        { title: 'Приоритет', dataIndex: 'priority' },
      ]} />
    </>
  )
}


