import React, { useEffect, useState } from 'react'
import { useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'
import { exportTransactionsToCSV } from '../utils/csvExport'

const Transactions: React.FC = () => {
  const wallet = useWalletStore((state) => state.wallet)
  const [transactions, setTransactions] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const loadTransactions = async () => {
      if (!wallet) return
      try {
        setLoading(true)
        const walletAddress = (wallet as any).wallet_address || wallet.walletAddress
        const response = await apiClient.getTransactionHistory(walletAddress, 100, 0)
        if (response.data.status === 'success') {
          setTransactions(response.data.data?.transactions || [])
        }
      } catch (error) {
        console.error('Failed to load transactions:', error)
      } finally {
        setLoading(false)
      }
    }
    loadTransactions()
  }, [wallet])

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Transactions</h1>
        {transactions.length > 0 && (
          <button
            onClick={() => exportTransactionsToCSV(transactions.map((tx) => ({
              id: tx.id || '',
              transactionHash: tx.transaction_hash || tx.transactionHash || '',
              senderWallet: tx.sender_wallet || tx.senderWallet || '',
              receiverWallet: tx.receiver_wallet || tx.receiverWallet || '',
              amount: tx.amount || 0,
              fee: tx.fee || 0,
              signature: tx.signature || '',
              status: tx.status || 'pending',
              transactionType: tx.transaction_type || tx.transactionType || 'normal',
              createdAt: tx.created_at || tx.createdAt || new Date().toISOString(),
            })), `transactions_${new Date().toISOString().split('T')[0]}.csv`)}
            className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-semibold transition"
          >
            📥 Export CSV
          </button>
        )}
      </div>

      <div className="card">
        {loading ? (
          <p className="text-center text-gray-600 dark:text-gray-400">Loading...</p>
        ) : transactions.length === 0 ? (
          <p className="text-center text-gray-600 dark:text-gray-400">No transactions found</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-100 dark:bg-gray-700">
                <tr>
                  <th className="px-4 py-2 text-left">Hash</th>
                  <th className="px-4 py-2 text-left">From</th>
                  <th className="px-4 py-2 text-left">To</th>
                  <th className="px-4 py-2 text-right">Amount</th>
                  <th className="px-4 py-2 text-center">Status</th>
                  <th className="px-4 py-2 text-left">Date</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx: any) => (
                  <tr key={tx.transaction_hash} className="border-b hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="px-4 py-2 text-sm font-mono text-blue-600 dark:text-blue-400">
                      {(tx.transaction_hash || tx.transactionHash)?.slice(0, 8)}...
                    </td>
                    <td className="px-4 py-2 text-sm font-mono">
                      {((tx.sender_wallet || tx.senderWallet) || '').slice(0, 8)}...
                    </td>
                    <td className="px-4 py-2 text-sm font-mono">
                      {((tx.receiver_wallet || tx.receiverWallet) || '').slice(0, 8)}...
                    </td>
                    <td className="px-4 py-2 text-right text-sm font-semibold">
                      {(tx.amount || 0).toFixed(8)}
                    </td>
                    <td className="px-4 py-2 text-center">
                      <span className={`px-2 py-1 rounded text-xs font-semibold ${
                        (tx.status || tx.Status) === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                        (tx.status || tx.Status) === 'confirmed' ? 'bg-green-100 text-green-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {(tx.status || tx.Status) || 'pending'}
                      </span>
                    </td>
                    <td className="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
                      {new Date(tx.created_at || tx.createdAt).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default Transactions
