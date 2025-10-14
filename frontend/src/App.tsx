import { App as AntApp, ConfigProvider } from 'antd'
import ruRU from 'antd/locale/ru_RU'
import { Navigate, Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './components/ProtectedRoute'
import { AppLayout } from './layouts/AppLayout'
import { AuthLayout } from './layouts/AuthLayout'
import { DashboardPage } from './pages/Dashboard'
import { DefectNewPage } from './pages/DefectNew'
import { DefectsPage } from './pages/Defects'
import { LoginPage } from './pages/Login'
import { ProjectsPage } from './pages/Projects'
import { ReportsPage } from './pages/Reports'

export default function App() {
  return (
    <ConfigProvider locale={ruRU}>
      <AntApp>
        <Routes>
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
          </Route>
          <Route element={<ProtectedRoute><AppLayout /></ProtectedRoute>}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="/defects" element={<DefectsPage />} />
            <Route path="/defects/new" element={<DefectNewPage />} />
            <Route path="/reports" element={<ReportsPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </AntApp>
    </ConfigProvider>
  )
}


