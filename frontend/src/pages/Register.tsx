import { api } from '@/services/api'
import { Button, Card, Form, Input, Select, Typography, message } from 'antd'
import { Link, useNavigate } from 'react-router-dom'

export function RegisterPage() {
  const navigate = useNavigate()
  const onFinish = async (values: any) => {
    try {
      await api.post('/api/auth/register', values)
      message.success('Регистрация прошла успешно, войдите')
      navigate('/login')
    } catch {
      // уведомление покажет interceptor
    }
  }
  return (
    <Card style={{ width: 420 }}>
      <Typography.Title level={4} style={{ textAlign: 'center' }}>Регистрация</Typography.Title>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Пароль" rules={[{ required: true, min: 6 }]}>
          <Input.Password />
        </Form.Item>
        <Form.Item name="full_name" label="ФИО" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="role" label="Роль" rules={[{ required: true }]}>
          <Select options={[
            { value: 'engineer', label: 'Инженер' },
            { value: 'manager', label: 'Менеджер' },
            { value: 'observer', label: 'Наблюдатель' },
          ]} />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>Зарегистрироваться</Button>
      </Form>
      <div style={{ marginTop: 12, textAlign: 'center' }}>
        Уже есть аккаунт? <Link to="/login">Войти</Link>
      </div>
    </Card>
  )
}


