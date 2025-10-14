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
        {role === 'manager' && <Button type="primary" onClick={() => setOpen(true)}>Создать проект</Button>}
      </Space>
      <Table rowKey="id" dataSource={items} columns={[
        { title: 'Название', dataIndex: 'name' },
        { title: 'Описание', dataIndex: 'description' },
      ]} />
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
    </>
  )
}


