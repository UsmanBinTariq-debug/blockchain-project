import { create } from 'zustand'
import { User, Wallet } from '../types'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  token: string | null
  setUser: (user: User | null) => void
  setToken: (token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  // Treat legacy 'demo-token' as invalid and remove it
  token: (() => {
    const t = localStorage.getItem('auth_token')
    if (t === 'demo-token') {
      localStorage.removeItem('auth_token')
      return null
    }
    return t
  })(),
  isAuthenticated: (() => {
    const t = localStorage.getItem('auth_token')
    return !!t && t !== 'demo-token'
  })(),
  setUser: (user) => set({ user, isAuthenticated: !!user }),
  setToken: (token) => {
    localStorage.setItem('auth_token', token)
    set({ token, isAuthenticated: true })
  },
  logout: () => {
    localStorage.removeItem('auth_token')
    set({ user: null, isAuthenticated: false, token: null })
  },
}))

interface WalletState {
  wallet: Wallet | null
  balance: number
  setWallet: (wallet: Wallet | null) => void
  setBalance: (balance: number) => void
}

export const useWalletStore = create<WalletState>((set) => ({
  wallet: null,
  balance: 0,
  setWallet: (wallet) => set({ wallet }),
  setBalance: (balance) => set({ balance }),
}))

interface UIState {
  isDarkMode: boolean
  toggleDarkMode: () => void
  sidebarOpen: boolean
  toggleSidebar: () => void
}

export const useUIStore = create<UIState>((set) => ({
  isDarkMode: localStorage.getItem('dark_mode') === 'true',
  toggleDarkMode: () => set((state) => {
    const newValue = !state.isDarkMode
    localStorage.setItem('dark_mode', String(newValue))
    return { isDarkMode: newValue }
  }),
  sidebarOpen: true,
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
}))
