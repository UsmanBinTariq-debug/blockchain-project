export interface User {
  id: string
  email: string
  fullName: string
  cnic: string
  walletId: string
  publicKey: string
  isVerified: boolean
}

export interface Wallet {
  id: string
  userId: string
  walletAddress: string
  balance: number
  lastUpdated: string
}

export interface Transaction {
  id: string
  transactionHash: string
  senderWallet: string
  receiverWallet: string
  amount: number
  fee: number
  note?: string
  signature: string
  status: 'pending' | 'confirmed' | 'failed'
  createdAt: string
  transactionType: 'transfer' | 'zakat'
}

export interface Block {
  id: string
  blockIndex: number
  timestamp: number
  previousHash: string
  hash: string
  nonce: number
  merkleRoot: string
  difficulty: number
  minedBy?: string
  createdAt: string
}

export interface UTXO {
  id: string
  transactionHash: string
  outputIndex: number
  walletAddress: string
  amount: number
  isSpent: boolean
  spentInTransaction?: string
}

export interface ZakatTransaction {
  id: string
  walletAddress: string
  amount: number
  zakatPercentage: number
  monthYear: string
  createdAt: string
}

export interface Beneficiary {
  id: string
  userId: string
  beneficiaryWalletId: string
  nickname: string
  createdAt: string
}
