import axios, { AxiosInstance } from 'axios'
import { API_URL } from '../utils/constants'
import { useAuthStore } from '../utils/store'

interface ApiResponse<T> {
  status: string
  message: string
  data?: T
  error?: string
}

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    // Add token to requests
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      return config
    })

    // Handle 401 responses globally: clear invalid token to avoid repeated 401s
    this.client.interceptors.response.use(
      (resp) => resp,
      (error) => {
        const status = error?.response?.status
        if (status === 401) {
          try {
            localStorage.removeItem('auth_token')
            // update zustand store
            const logout = useAuthStore.getState().logout
            if (logout) logout()
          } catch (e) {
            // swallow errors
          }
        }
        return Promise.reject(error)
      }
    )
  }

  async register(email: string, fullName: string, cnic: string, password: string) {
    return this.client.post<ApiResponse<any>>('/auth/register', {
      email,
      full_name: fullName,
      cnic,
      password,
    })
  }

  async login(email: string, password: string) {
    return this.client.post<ApiResponse<any>>('/auth/login', { email, password })
  }

  async getWallet() {
    return this.client.get<ApiResponse<any>>('/wallet/profile')
  }

  async getBalance(walletAddress: string) {
    return this.client.post<ApiResponse<any>>('/wallet/balance', {
      wallet_address: walletAddress,
    })
  }

  async getTransactionHistory(walletAddress: string, limit = 10, offset = 0) {
    return this.client.get<ApiResponse<any>>('/transaction/history', {
      params: { wallet_address: walletAddress, limit, offset },
    })
  }

  async sendTransaction(senderWallet: string, receiverWallet: string, amount: number, fee: number, note: string, signature: string) {
    return this.client.post<ApiResponse<any>>('/transaction/send', {
      sender_wallet: senderWallet,
      receiver_wallet: receiverWallet,
      amount,
      fee,
      note,
      signature,
    })
  }

  async getBlocks(limit = 10, offset = 0) {
    return this.client.get<ApiResponse<any>>('/blockchain/blocks', {
      params: { limit, offset },
    })
  }

  async getBlock(hash: string) {
    return this.client.get<ApiResponse<any>>(`/blockchain/blocks/${hash}`)
  }

  async mineBlock(minerAddress: string) {
    return this.client.post<ApiResponse<any>>('/blockchain/mine', {
      miner_address: minerAddress,
    })
  }

  async getPendingTransactions() {
    return this.client.get<ApiResponse<any>>('/transaction/pending')
  }

  async getMonthlyReport(walletAddress: string) {
    return this.client.get<ApiResponse<any>>('/reports/monthly', {
      params: { wallet_address: walletAddress },
    })
  }

  async getZakatReport(walletAddress: string) {
    return this.client.get<ApiResponse<any>>('/reports/zakat', {
      params: { wallet_address: walletAddress },
    })
  }

  async addBeneficiary(beneficiaryWalletId: string, nickname: string) {
    return this.client.post<ApiResponse<any>>('/beneficiary/add', {
      beneficiary_wallet_id: beneficiaryWalletId,
      nickname,
    })
  }

  async getBeneficiaries() {
    return this.client.get<ApiResponse<any>>('/beneficiary/list')
  }

  async getSystemLogs(type: string = 'ALL', limit: number = 50, offset: number = 0) {
    return this.client.get<ApiResponse<any>>('/system/logs', {
      params: { type, limit, offset },
    })
  }

  async getSystemLogStats() {
    return this.client.get<ApiResponse<any>>('/system/logs/stats')
  }

  async getSystemHealth() {
    return this.client.get<ApiResponse<any>>('/system/health')
  }
}

export const apiClient = new ApiClient()
export const api = apiClient

