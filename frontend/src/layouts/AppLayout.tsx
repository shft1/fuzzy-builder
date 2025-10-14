import { useAuthStore } from '@/stores/auth'
import { Button, Dropdown, Layout, Menu } from 'antd'
import { Link, Outlet, useLocation } from 'react-router-dom'

const { Header, Sider, Content } = Layout

export function AppLayout() {
  const location = useLocation()
  const selected = [location.pathname]
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg">
        <div style={{ color: '#fff', padding: 16, fontWeight: 600 }}>Fuzzy</div>
        <Menu theme="dark" mode="inline" selectedKeys={selected}
          items={[
            { key: '/dashboard', label: <Link to="/dashboard">Дашборд</Link> },
            { key: '/projects', label: <Link to="/projects">Проекты</Link> },
            { key: '/defects', label: <Link to="/defects">Дефекты</Link> },
            { key: '/reports', label: <Link to="/reports">Отчеты</Link> }
          ]}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <UserMenu />
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

function UserMenu() {
  const user = useAuthStore(s => s.user)
  const logout = useAuthStore(s => s.logout)
  const items = [{ key: 'logout', label: <span onClick={logout}>Выйти</span> }]
  return (
    <Dropdown menu={{ items }} placement="bottomRight">
      <Button>{user?.role || 'Пользователь'}</Button>
    </Dropdown>
  )
}


