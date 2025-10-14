import { ManagerOnly } from '@/components/ManagerOnly'
import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Form, Input, Modal, Space, Table } from 'antd'
import { useEffect, useState } from 'react'

type Project = { id: number; name: string; description: string }

export function ProjectsPage() {
  const token = useAuthStore(s => s.token)
  const role = useAuthStore(s => s.user?.role)
  const [items, setItems] = useState<Project[]>([])
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [editing, setEditing] = useState<Project | null>(null)
  const [search, setSearch] = useState('')
  useEffect(() => {
    if (!token) return
    api.get('/api/projects', { headers: { Authorization: `Bearer ${token}` } }).then(r => setItems(r.data))
  }, [token])
  const refresh = async () => {
    if (!token) return
    const r = await api.get('/api/projects', { headers: { Authorization: `Bearer ${token}` } })
    setItems(r.data)
  }
  const onCreate = async (values: any) => {
    setLoading(true)
    try {
      await api.post('/api/projects', values, { headers: { Authorization: `Bearer ${token}` } })
      setOpen(false)
      await refresh()
    } finally { setLoading(false) }
  }
  return (
    <>
      <Space style={{ marginBottom: 16 }}>
        <ManagerOnly><Button type="primary" onClick={() => setOpen(true)}>Создать проект</Button></ManagerOnly>
        <Input placeholder="Поиск" value={search} onChange={e=>setSearch(e.target.value)} allowClear style={{ width: 260 }} />
      </Space>
      <Table rowKey="id" dataSource={items.filter(p => p.name.toLowerCase().includes(search.toLowerCase()))}
        pagination={{ pageSize: 10 }}
        columns={[
          { title: 'Название', dataIndex: 'name', sorter: (a:Project,b:Project)=>a.name.localeCompare(b.name) },
          { title: 'Описание', dataIndex: 'description' },
          role === 'manager' ? { title: 'Действия', render: (_:any, r:Project) => <Space>
            <ManagerOnly><Button size="small" onClick={() => setEditing(r)}>Редактировать</Button></ManagerOnly>
            <ManagerOnly><Button size="small" danger onClick={() => Modal.confirm({ title: 'Удалить проект?', onOk: async () => { await api.delete(`/api/projects/${r.id}`, { headers: { Authorization: `Bearer ${token}` } }); refresh() } })}>Удалить</Button></ManagerOnly>
          </Space> } : {}
        ]}
      />
      <Modal title="Новый проект" open={open} onCancel={() => setOpen(false)} footer={null} destroyOnClose>
        <Form layout="vertical" onFinish={onCreate}>
          <Form.Item name="name" label="Название" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Описание" rules={[{ required: true }]}>
            <Input.TextArea rows={4} />
          </Form.Item>
          <Space>
            <Button onClick={() => setOpen(false)}>Отмена</Button>
            <Button type="primary" htmlType="submit" loading={loading}>Создать</Button>
          </Space>
        </Form>
      </Modal>

      <Modal title="Редактировать проект" open={!!editing} onCancel={() => setEditing(null)} footer={null} destroyOnClose>
        <Form layout="vertical" initialValues={editing ?? {}} onFinish={async (v)=>{
          if (!editing) return
          await api.put(`/api/projects/${editing.id}`, v, { headers: { Authorization: `Bearer ${token}` } })
          setEditing(null)
          refresh()
        }}>
          <Form.Item name="name" label="Название" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Описание" rules={[{ required: true }]}>
            <Input.TextArea rows={4} />
          </Form.Item>
          <Space>
            <Button onClick={() => setEditing(null)}>Отмена</Button>
            <Button type="primary" htmlType="submit">Сохранить</Button>
          </Space>
        </Form>
      </Modal>
    </>
  )
}


