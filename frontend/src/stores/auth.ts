import { create } from 'zustand';

type User = { user_id: number; role: string } | null

type AuthState = {
  token: string | null
  user: User
  setToken: (t: string | null) => void
  setUser: (u: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: null,
  setToken: (t) => { if (t) localStorage.setItem('token', t); else localStorage.removeItem('token'); set({ token: t }) },
  setUser: (u) => set({ user: u }),
  logout: () => { localStorage.removeItem('token'); set({ token: null, user: null }) },
}))


