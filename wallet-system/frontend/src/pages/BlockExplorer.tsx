import { useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { api } from '../services/api'

interface Block {
  block_index: number
  hash: string
  previous_hash: string
  nonce: number
  merkle_root: string
  difficulty: number
  mined_by: string
  timestamp: number
  created_at: string
}

interface Transaction {
  transaction_hash: string
  sender_wallet: string
  receiver_wallet: string
  amount: number
  fee: number
  status: string
  created_at: string
}

export function BlockExplorer() {
  const { hash } = useParams()
  const navigate = useNavigate()
  const [blocks, setBlocks] = useState<Block[]>([])
  const [block, setBlock] = useState<Block | null>(null)
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)

  useEffect(() => {
    if (hash) {
      fetchBlock(hash)
    } else {
      fetchBlocks()
    }
  }, [hash, page])

  const fetchBlocks = async () => {
    try {
      setLoading(true)
      const response = await api.getBlocks(10, (page - 1) * 10)
      setBlocks(response.data?.data?.blocks || [])
    } catch (error) {
      console.error('Failed to fetch blocks:', error)
      setBlocks([])
    } finally {
      setLoading(false)
    }
  }

  const fetchBlock = async (blockHash: string) => {
    try {
      setLoading(true)
      const response = await api.getBlock(blockHash)
      setBlock(response.data?.data?.block || null)
      setTransactions(response.data?.data?.transactions || [])
    } catch (error) {
      console.error('Failed to fetch block:', error)
      setBlock(null)
      setTransactions([])
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (timestamp: number | string) => {
    const ts = typeof timestamp === 'string' ? parseInt(timestamp) : timestamp
    return new Date(ts * 1000).toLocaleString()
  }

  const truncateHash = (hash: string, length = 16) => {
    return hash.substring(0, length) + '...' + hash.substring(hash.length - 8)
  }

  const formatAddress = (address: string) => {
    return address.substring(0, 10) + '...' + address.substring(address.length - 8)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600 dark:text-gray-400">Loading...</div>
      </div>
    )
  }

  if (hash && block) {
    return (
      <div className="max-w-6xl mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/explorer')}
          className="mb-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          ← Back to Blocks
        </button>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
          <h1 className="text-2xl font-bold mb-4">Block #{block.block_index}</h1>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Hash</p>
                <p className="font-mono text-sm break-all">{block.hash}</p>
              </div>

              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Previous Hash</p>
                <p className="font-mono text-sm break-all">{block.previous_hash}</p>
              </div>

              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Merkle Root</p>
                <p className="font-mono text-sm break-all">{block.merkle_root}</p>
              </div>
            </div>

            <div>
              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Difficulty</p>
                <p className="text-lg font-semibold">{block.difficulty}</p>
              </div>

              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Nonce</p>
                <p className="text-lg font-semibold">{block.nonce}</p>
              </div>

              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Mined By</p>
                <p className="font-mono text-sm">{formatAddress(block.mined_by)}</p>
              </div>

              <div className="mb-4">
                <p className="text-gray-600 dark:text-gray-400 text-sm">Timestamp</p>
                <p className="text-sm">{formatDate(block.timestamp)}</p>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-bold mb-4">Transactions ({transactions.length})</h2>

          {transactions.length === 0 ? (
            <p className="text-gray-600 dark:text-gray-400">No transactions in this block</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b dark:border-gray-700">
                    <th className="text-left py-2 px-2">Hash</th>
                    <th className="text-left py-2 px-2">From</th>
                    <th className="text-left py-2 px-2">To</th>
                    <th className="text-right py-2 px-2">Amount</th>
                    <th className="text-right py-2 px-2">Fee</th>
                    <th className="text-left py-2 px-2">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {transactions.map((tx) => (
                    <tr key={tx.transaction_hash} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                      <td className="py-3 px-2">
                        <code className="text-xs font-mono">{truncateHash(tx.transaction_hash)}</code>
                      </td>
                      <td className="py-3 px-2">
                        <code className="text-xs font-mono">{formatAddress(tx.sender_wallet)}</code>
                      </td>
                      <td className="py-3 px-2">
                        <code className="text-xs font-mono">{formatAddress(tx.receiver_wallet)}</code>
                      </td>
                      <td className="py-3 px-2 text-right font-semibold">{tx.amount.toFixed(2)} CRW</td>
                      <td className="py-3 px-2 text-right">{tx.fee.toFixed(2)} CRW</td>
                      <td className="py-3 px-2">
                        <span className={`px-2 py-1 rounded text-xs font-semibold ${
                          tx.status === 'confirmed'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                            : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                        }`}>
                          {tx.status}
                        </span>
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

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Block Explorer</h1>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div className="mb-4 flex justify-between items-center">
          <h2 className="text-xl font-bold">Recent Blocks</h2>
          <div className="flex gap-2">
            <button
              onClick={() => setPage(Math.max(1, page - 1))}
              disabled={page === 1}
              className="px-4 py-2 bg-gray-300 dark:bg-gray-700 rounded disabled:opacity-50"
            >
              Previous
            </button>
            <span className="px-4 py-2">Page {page}</span>
            <button
              onClick={() => setPage(page + 1)}
              className="px-4 py-2 bg-gray-300 dark:bg-gray-700 rounded hover:bg-gray-400"
            >
              Next
            </button>
          </div>
        </div>

        {blocks.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-400">No blocks found</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b dark:border-gray-700">
                  <th className="text-left py-2 px-2">Height</th>
                  <th className="text-left py-2 px-2">Hash</th>
                  <th className="text-left py-2 px-2">Miner</th>
                  <th className="text-right py-2 px-2">Difficulty</th>
                  <th className="text-left py-2 px-2">Time</th>
                </tr>
              </thead>
              <tbody>
                {blocks.map((blk) => (
                  <tr
                    key={blk.hash}
                    onClick={() => navigate(`/explorer/${blk.hash}`)}
                    className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
                  >
                    <td className="py-3 px-2 font-semibold">{blk.block_index}</td>
                    <td className="py-3 px-2">
                      <code className="text-xs font-mono">{truncateHash(blk.hash)}</code>
                    </td>
                    <td className="py-3 px-2">
                      <code className="text-xs font-mono">{formatAddress(blk.mined_by)}</code>
                    </td>
                    <td className="py-3 px-2 text-right font-semibold">{blk.difficulty}</td>
                    <td className="py-3 px-2 text-sm">{formatDate(blk.timestamp)}</td>
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

export default BlockExplorer
