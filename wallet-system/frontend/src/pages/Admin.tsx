import React, { useEffect, useState } from 'react'
import { useAuthStore } from '../utils/store'
import { useNavigate } from 'react-router-dom'
import { apiClient } from '../services/api'

interface SystemLog {
  id: string
  log_type: string
  message: string
  wallet_address?: string
  ip_address?: string
  created_at: string
}

interface SystemStats {
  total_logs: number
  transaction_logs: number
  block_logs: number
  zakat_logs: number
  error_logs: number
  auth_logs: number
}

const Admin: React.FC = () => {
  const user = useAuthStore((state) => state.user)
  const navigate = useNavigate()
  const [logs, setLogs] = useState<SystemLog[]>([])
  const [stats, setStats] = useState<SystemStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [filterType, setFilterType] = useState<string>('ALL')

  useEffect(() => {
    // Simple admin check - in production, verify via backend
    if (!user || user.email !== 'admin@crypto-wallet.local') {
      navigate('/dashboard')
      return
    }
    
    loadData()
  }, [user, navigate])

  const loadData = async () => {
    try {
      setLoading(true)
      // Fetch logs
      const logsResponse = await apiClient.getSystemLogs(filterType, 50, 0)
      if (logsResponse.data.status === 'success') {
        setLogs(logsResponse.data.data?.logs || [])
      }

      // Fetch stats
      const statsResponse = await apiClient.getSystemLogStats()
      if (statsResponse.data.status === 'success') {
        setStats(statsResponse.data.data)
      }
    } catch (error) {
      console.error('Failed to load admin data:', error)
      // Fallback to empty state rather than showing error
      setLogs([])
    } finally {
      setLoading(false)
    }
  }

  const handleFilterChange = async (type: string) => {
    setFilterType(type)
    try {
      const response = await apiClient.getSystemLogs(type, 50, 0)
      if (response.data.status === 'success') {
        setLogs(response.data.data?.logs || [])
      }
    } catch (error) {
      console.error('Failed to filter logs:', error)
    }
  }

  const refreshData = () => {
    loadData()
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString()
  }

  const getLogTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      TRANSACTION_CREATED: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
      BLOCK_MINED: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
      ZAKAT_DEDUCTED: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
      ERROR: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
      AUTH: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
    }
    return colors[type] || 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600 dark:text-gray-400">Loading admin panel...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Admin Panel</h1>

      {/* System Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Total Logs
          </h3>
          <p className="text-3xl font-bold text-blue-600">{stats?.total_logs || 0}</p>
        </div>

        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Transactions
          </h3>
          <p className="text-3xl font-bold text-green-600">
            {stats?.transaction_logs || 0}
          </p>
        </div>

        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Blocks Mined
          </h3>
          <p className="text-3xl font-bold text-purple-600">
            {stats?.block_logs || 0}
          </p>
        </div>

        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Errors
          </h3>
          <p className="text-3xl font-bold text-red-600">
            {stats?.error_logs || 0}
          </p>
        </div>
      </div>

      {/* System Logs */}
      <div className="card">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-800 dark:text-white">System Logs</h2>
          <button
            onClick={refreshData}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
          >
            🔄 Refresh
          </button>
        </div>

        {/* Filter */}
        <div className="mb-6">
          <label className="text-gray-600 dark:text-gray-400 text-sm mb-2 block">
            Filter by Type
          </label>
          <select
            value={filterType}
            onChange={(e) => handleFilterChange(e.target.value)}
            className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-800 dark:text-white"
          >
            <option value="ALL">All Types</option>
            <option value="TRANSACTION_CREATED">Transaction Created</option>
            <option value="BLOCK_MINED">Block Mined</option>
            <option value="ZAKAT_DEDUCTED">Zakat Deducted</option>
            <option value="ERROR">Errors</option>
            <option value="AUTH">Auth Events</option>
          </select>
        </div>

        {/* Logs Table */}
        {logs.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-400">No logs found</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b dark:border-gray-700">
                  <th className="text-left py-3 px-4">Time</th>
                  <th className="text-left py-3 px-4">Type</th>
                  <th className="text-left py-3 px-4">Message</th>
                  <th className="text-left py-3 px-4">Wallet</th>
                  <th className="text-left py-3 px-4">IP</th>
                </tr>
              </thead>
              <tbody>
                {logs.map((log) => (
                  <tr key={log.id} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="py-3 px-4 text-xs">{formatDate(log.created_at)}</td>
                    <td className="py-3 px-4">
                      <span className={`px-3 py-1 rounded text-xs font-semibold ${getLogTypeColor(log.log_type)}`}>
                        {log.log_type}
                      </span>
                    </td>
                    <td className="py-3 px-4">{log.message}</td>
                    <td className="py-3 px-4 font-mono text-xs">
                      {log.wallet_address ? log.wallet_address.substring(0, 12) + '...' : 'N/A'}
                    </td>
                    <td className="py-3 px-4 text-xs">{log.ip_address || 'N/A'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* System Health */}
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">System Health</h2>

        <div className="space-y-4">
          <div className="flex items-center justify-between p-4 bg-green-50 dark:bg-green-900 rounded">
            <div>
              <p className="font-semibold text-gray-800 dark:text-white">Backend Server</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Running on port 8080</p>
            </div>
            <div className="w-3 h-3 rounded-full bg-green-600"></div>
          </div>

          <div className="flex items-center justify-between p-4 bg-green-50 dark:bg-green-900 rounded">
            <div>
              <p className="font-semibold text-gray-800 dark:text-white">Database</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Supabase PostgreSQL</p>
            </div>
            <div className="w-3 h-3 rounded-full bg-green-600"></div>
          </div>

          <div className="flex items-center justify-between p-4 bg-green-50 dark:bg-green-900 rounded">
            <div>
              <p className="font-semibold text-gray-800 dark:text-white">Blockchain</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">All nodes synchronized</p>
            </div>
            <div className="w-3 h-3 rounded-full bg-green-600"></div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Admin
