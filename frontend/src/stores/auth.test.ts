import { describe, expect, it } from 'vitest'
import { useAuthStore } from './auth'

describe('auth store', () => {
  it('sets and clears token', () => {
    useAuthStore.getState().setToken('abc')
    expect(useAuthStore.getState().token).toBe('abc')
    useAuthStore.getState().logout()
    expect(useAuthStore.getState().token).toBeNull()
  })
})


