export const API_URL = 'http://localhost:8080/api'
export const SUPABASE_URL = ''
export const SUPABASE_ANON_KEY = ''

export const TRANSACTION_STATUS = {
  PENDING: 'pending',
  CONFIRMED: 'confirmed',
  FAILED: 'failed',
} as const

export const TRANSACTION_TYPES = {
  TRANSFER: 'transfer',
  ZAKAT: 'zakat',
} as const

export const ZAKAT_PERCENTAGE = 2.5
export const MIN_TRANSACTION_AMOUNT = 0.001
export const MAX_TRANSACTION_AMOUNT = 1000000000

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  DASHBOARD: '/dashboard',
  WALLET: '/wallet',
  SEND_MONEY: '/send-money',
  TRANSACTIONS: '/transactions',
  BLOCK_EXPLORER: '/block-explorer',
  REPORTS: '/reports',
  ADMIN: '/admin',
  PROFILE: '/profile',
} as const
