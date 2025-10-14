import { api } from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import { Button, Card, Form, Input, Typography } from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'

export function LoginPage() {
  const setToken = useAuthStore(s => s.setToken)
  const setUser = useAuthStore(s => s.setUser)
  const navigate = useNavigate()
  const location = useLocation() as any

  const onFinish = async (values: any) => {
    const { data } = await api.post('/api/auth/login', values)
    setToken(data.token)
    // fetch user info
    const me = await api.get('/api/users/me', { headers: { Authorization: `Bearer ${data.token}` } })
    setUser(me.data)
    const to = location.state?.from?.pathname || '/dashboard'
    navigate(to, { replace: true })
  }

  return (
    <Card style={{ width: 360 }}>
      <Typography.Title level={4} style={{ textAlign: 'center' }}>Вход</Typography.Title>
      <Form layout="vertical" onFinish={onFinish}>
        <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="password" label="Пароль" rules={[{ required: true }]}>
          <Input.Password />
        </Form.Item>
        <Button type="primary" htmlType="submit" block>Войти</Button>
      </Form>
    </Card>
  )
}


