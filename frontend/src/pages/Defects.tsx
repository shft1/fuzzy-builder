import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Modal, Select, Space, Table, Tag, Upload, message } from 'antd'
import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

type Defect = { id: number; title: string; status: string; priority: string; project_id: number }

export function DefectsPage() {
  const token = useAuthStore(s => s.token)
  const [items, setItems] = useState<Defect[]>([])
  const [status, setStatus] = useState<string | undefined>()
  const [priority, setPriority] = useState<string | undefined>()
  const [changingId, setChangingId] = useState<number | null>(null)
  const [newStatus, setNewStatus] = useState<string>('in_progress')
  const [attachDefect, setAttachDefect] = useState<number | null>(null)
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
        <Button type="primary"><Link to="/defects/new">Создать дефект</Link></Button>
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
        { title: 'Заголовок', dataIndex: 'title', render: (_:any, r:Defect) => <Link to={`/defects/${r.id}`}>{r.title}</Link> },
        { title: 'Проект', dataIndex: 'project_id' },
        { title: 'Статус', dataIndex: 'status', render: (v: string) => <Tag color={v==='closed'?'green':v==='in_progress'?'blue':v==='on_review'?'orange':'default'}>{v}</Tag> },
        { title: 'Приоритет', dataIndex: 'priority' },
        { title: 'Действия', render: (_: any, r: Defect) => (
          <Space>
            <Button size="small" onClick={() => { setChangingId(r.id); setNewStatus('in_progress') }}>Статус</Button>
            <Button size="small" onClick={() => setAttachDefect(r.id)}>Вложение</Button>
          </Space>
        )}
      ]} />

      <Modal title="Сменить статус" open={changingId!=null} onCancel={() => setChangingId(null)} onOk={async () => {
        try {
          await api.put(`/api/defects/${changingId}/status`, { status: newStatus }, { headers: { Authorization: `Bearer ${token}` } })
          setChangingId(null)
          const url = '/api/defects' + (query ? `?${query}` : '')
          const r = await api.get(url, { headers: { Authorization: `Bearer ${token}` } })
          setItems(r.data)
        } catch(e:any){ message.error('Не удалось сменить статус') }
      }}>
        <Select style={{ width: 240 }} value={newStatus} onChange={setNewStatus}
          options={[
            { value: 'in_progress', label: 'В работе' },
            { value: 'on_review', label: 'На проверке' },
            { value: 'closed', label: 'Закрыта' },
          ]}
        />
      </Modal>

      <Modal title="Загрузить вложение" open={attachDefect!=null} onCancel={() => setAttachDefect(null)} footer={null}>
        <Upload name="file" multiple={false} action={`/api/defects/${attachDefect}/attachments`} 
          headers={{ Authorization: `Bearer ${token||''}` }}
          onChange={(info) => { if (info.file.status==='done') { message.success('Загружено'); setAttachDefect(null) } }}> 
          <Button type="primary">Выбрать файл</Button>
        </Upload>
      </Modal>
    </>
  )
}


