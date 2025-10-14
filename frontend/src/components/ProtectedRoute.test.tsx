import { useAuthStore } from '@/stores/auth'
import { render } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { ProtectedRoute } from './ProtectedRoute'

function App() {
  return (
    <Routes>
      <Route path="/login" element={<div>login</div>} />
      <Route path="/secure" element={<ProtectedRoute><div>secure</div></ProtectedRoute>} />
    </Routes>
  )
}

describe('ProtectedRoute', () => {
  it('redirects when no token', () => {
    useAuthStore.getState().logout()
    const { container } = render(<MemoryRouter initialEntries={["/secure"]}><App/></MemoryRouter>)
    expect(container.textContent).toContain('login')
  })
})


