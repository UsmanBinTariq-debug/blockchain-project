import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'

interface Transaction {
  transaction_hash: string
  sender_wallet: string
  receiver_wallet: string
  amount: number
  fee: number
  status: string
  created_at: string
}

const Dashboard: React.FC = () => {
  const wallet = useWalletStore((state) => state.wallet)
  const balance = useWalletStore((state) => state.balance)
  const [recentTransactions, setRecentTransactions] = useState<Transaction[]>([])
  const [loadingTxs, setLoadingTxs] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    loadRecentTransactions()
  }, [wallet])

  const loadRecentTransactions = async () => {
    if (!wallet?.walletAddress) return
    
    try {
      setLoadingTxs(true)
      const response = await apiClient.getTransactionHistory(wallet.walletAddress, 5, 0)
      if (response.data?.status === 'success' && response.data.data?.transactions) {
        setRecentTransactions(response.data.data.transactions)
      }
    } catch (error) {
      console.error('Failed to load recent transactions:', error)
    } finally {
      setLoadingTxs(false)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  const formatAddress = (address: string) => {
    return address?.substring(0, 10) + '...' + address?.substring(address.length - 8) || 'Unknown'
  }

  const isOutgoing = (tx: Transaction) => tx.sender_wallet === wallet?.walletAddress

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Dashboard</h1>

      {/* Balance Card */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card bg-gradient-to-br from-blue-500 to-blue-600 text-white">
          <h2 className="text-lg font-semibold mb-2">Total Balance</h2>
          <p className="text-4xl font-bold">{balance.toFixed(8)} CRW</p>
          <p className="text-sm mt-2 opacity-75">Crypto Wallet Token</p>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold mb-2 text-gray-800 dark:text-white">Monthly Zakat</h2>
          <p className="text-2xl font-bold text-green-600">{(balance * 0.025).toFixed(8)}</p>
          <p className="text-sm mt-2 text-gray-600 dark:text-gray-400">2.5% Deducted Monthly</p>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold mb-2 text-gray-800 dark:text-white">Active Wallet</h2>
          <p className="text-sm font-mono break-all text-gray-700 dark:text-gray-300">
            {wallet?.walletAddress?.substring(0, 16)}...
          </p>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-4">Quick Actions</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <a href="/send-money" className="btn-primary text-center">Send Money</a>
          <a href="/wallet" className="btn-secondary text-center">View Wallet</a>
          <a href="/transactions" className="btn-secondary text-center">Transactions</a>
          <a href="/explorer" className="btn-secondary text-center">Explorer</a>
        </div>
      </div>

      {/* Recent Transactions */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold text-gray-800 dark:text-white">Recent Transactions</h2>
          <button
            onClick={() => navigate('/transactions')}
            className="text-blue-600 dark:text-blue-400 hover:underline text-sm"
          >
            View All
          </button>
        </div>
        
        {loadingTxs ? (
          <p className="text-gray-600 dark:text-gray-400">Loading...</p>
        ) : recentTransactions.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-400">No transactions yet</p>
        ) : (
          <div className="space-y-2">
            {recentTransactions.map((tx) => (
              <div
                key={tx.transaction_hash}
                className="flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded hover:bg-gray-50 dark:hover:bg-gray-700"
              >
                <div className="flex-1">
                  <div className="flex items-center gap-3">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-semibold ${
                      isOutgoing(tx) ? 'bg-red-500' : 'bg-green-500'
                    }`}>
                      {isOutgoing(tx) ? '↑' : '↓'}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-800 dark:text-white">
                        {isOutgoing(tx) ? 'Sent to' : 'Received from'} {formatAddress(isOutgoing(tx) ? tx.receiver_wallet : tx.sender_wallet)}
                      </p>
                      <p className="text-xs text-gray-600 dark:text-gray-400">{formatDate(tx.created_at)}</p>
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <p className={`font-semibold ${isOutgoing(tx) ? 'text-red-600' : 'text-green-600'}`}>
                    {isOutgoing(tx) ? '-' : '+'}{tx.amount.toFixed(2)} CRW
                  </p>
                  <span className={`text-xs px-2 py-1 rounded ${
                    tx.status === 'confirmed' 
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' 
                      : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                  }`}>
                    {tx.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default Dashboard
