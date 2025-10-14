import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Card, DatePicker, Form, Input, Select, Space } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

export function DefectNewPage() {
  const token = useAuthStore(s => s.token)
  const navigate = useNavigate()
  const [projects, setProjects] = useState<any[]>([])
  useEffect(() => {
    if (!token) return
    api.get('/api/projects', { headers: { Authorization: `Bearer ${token}` } }).then(r => setProjects(r.data))
  }, [token])
  const onFinish = async (v: any) => {
    const body = {
      title: v.title,
      description: v.description,
      project_id: Number(v.project_id),
      assigned_to: v.assigned_to ? Number(v.assigned_to) : undefined,
      priority: v.priority,
      due_date: v.due_date ? v.due_date.toISOString() : undefined,
    }
    await api.post('/api/defects', body, { headers: { Authorization: `Bearer ${token}` } })
    navigate('/defects')
  }
  return (
    <Card>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="title" label="Заголовок" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="description" label="Описание" rules={[{ required: true }]}>
          <Input.TextArea rows={4} />
        </Form.Item>
        <Space wrap>
          <Form.Item name="project_id" label="Проект" rules={[{ required: true }]}>
            <Select style={{ width: 240 }} options={projects.map(p => ({ value: p.id, label: p.name }))} />
          </Form.Item>
          <Form.Item name="assigned_to" label="Исполнитель">
            <Input type="number" min={1} />
          </Form.Item>
          <Form.Item name="priority" label="Приоритет" rules={[{ required: true }]}>
            <Select style={{ width: 180 }}
              options={[{value:'low',label:'Низкий'},{value:'medium',label:'Средний'},{value:'high',label:'Высокий'}]}
            />
          </Form.Item>
          <Form.Item name="due_date" label="Срок">
            <DatePicker showTime />
          </Form.Item>
        </Space>
        <Space>
          <Button onClick={() => navigate('/defects')}>Отмена</Button>
          <Button type="primary" htmlType="submit">Создать</Button>
        </Space>
      </Form>
    </Card>
  )
}


