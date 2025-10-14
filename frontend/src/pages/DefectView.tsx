import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Card, Descriptions, Select, Space, Table, Tag } from 'antd'
import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'

export function DefectViewPage() {
  const { id } = useParams()
  const token = useAuthStore(s => s.token)
  const [d, setD] = useState<any>(null)
  const [attachments, setAttachments] = useState<any[]>([])
  const [status, setStatus] = useState<string>('in_progress')
  useEffect(() => {
    if (!id || !token) return
    api.get(`/api/defects/${id}`, { headers: { Authorization: `Bearer ${token}` } }).then(r => setD(r.data))
    api.get(`/api/defects/${id}/attachments`, { headers: { Authorization: `Bearer ${token}` } }).then(r => setAttachments(r.data))
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
      <Space style={{ marginTop: 16 }}>
        <Select style={{ width: 240 }} value={status} onChange={setStatus}
          options={[{value:'in_progress',label:'В работе'},{value:'on_review',label:'На проверке'},{value:'closed',label:'Закрыта'}]} />
        <Button type="primary" onClick={async ()=>{
          await api.put(`/api/defects/${id}/status`, { status }, { headers: { Authorization: `Bearer ${token}` } })
          const r = await api.get(`/api/defects/${id}`, { headers: { Authorization: `Bearer ${token}` } })
          setD(r.data)
        }}>Сменить статус</Button>
      </Space>
      <Card title="Вложения" style={{ marginTop: 16 }}>
        <Table rowKey="id" dataSource={attachments} columns={[
          { title: 'Файл', dataIndex: 'filename' },
          { title: 'Действия', render: (_:any, a:any) => <a href={`/api/attachments/${a.id}/download`} target="_blank">Скачать</a> }
        ]} pagination={false} />
      </Card>
    </Card>
  )
}


