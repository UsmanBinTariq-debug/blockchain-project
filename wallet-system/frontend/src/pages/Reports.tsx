import React, { useEffect, useState } from 'react'
import { useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'
import { exportMonthlyReportToCSV, exportZakatReportToCSV } from '../utils/csvExport'

interface MonthlyData {
  month: string
  incoming: number
  outgoing: number
  balance: number
  fee: number
}

interface ZakatData {
  month: string
  zakatAmount: number
  percentage: number
  deductedDate: string
}

const Reports: React.FC = () => {
  const wallet = useWalletStore((state) => state.wallet)
  const balance = useWalletStore((state) => state.balance)
  const [monthlyData, setMonthlyData] = useState<MonthlyData[]>([])
  const [zakatData, setZakatData] = useState<ZakatData[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (wallet?.walletAddress) {
      loadReportData()
    }
  }, [wallet])

  const loadReportData = async () => {
    try {
      setLoading(true)
      
      // Get transaction history for monthly summary
      if (wallet?.walletAddress) {
        const response = await apiClient.getTransactionHistory(wallet.walletAddress, 1000, 0)
        
        if (response.data?.data?.transactions) {
          const txns = response.data.data.transactions
          const monthMap: Record<string, MonthlyData> = {}
          
          txns.forEach((tx: any) => {
            const date = new Date(tx.created_at)
            const monthKey = date.toISOString().slice(0, 7)
            
            if (!monthMap[monthKey]) {
              monthMap[monthKey] = {
                month: monthKey,
                incoming: 0,
                outgoing: 0,
                balance: 0,
                fee: 0,
              }
            }
            
            const isOutgoing = tx.sender_wallet === wallet.walletAddress
            if (isOutgoing) {
              monthMap[monthKey].outgoing += tx.amount
              monthMap[monthKey].fee += tx.fee
            } else {
              monthMap[monthKey].incoming += tx.amount
            }
            monthMap[monthKey].balance = monthMap[monthKey].incoming - monthMap[monthKey].outgoing
          })
          
          setMonthlyData(Object.values(monthMap).sort((a, b) => b.month.localeCompare(a.month)))
        }
      }
      
      // Mock zakat data - in production this would come from the API
      const mockZakatData: ZakatData[] = [
        { month: new Date().toISOString().slice(0, 7), zakatAmount: balance * 0.025, percentage: 2.5, deductedDate: new Date().toLocaleDateString() },
      ]
      setZakatData(mockZakatData)
    } catch (error) {
      console.error('Failed to load reports:', error)
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (amount: number) => {
    return amount.toFixed(2)
  }

  const getMonthName = (monthStr: string) => {
    const [year, month] = monthStr.split('-')
    const date = new Date(parseInt(year), parseInt(month) - 1)
    return date.toLocaleString('default', { month: 'long', year: 'numeric' })
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-gray-600 dark:text-gray-400">Loading reports...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Reports</h1>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Total Received
          </h3>
          <p className="text-2xl font-bold text-green-600">
            {formatCurrency(monthlyData.reduce((sum, m) => sum + m.incoming, 0))} CRW
          </p>
        </div>

        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Total Sent
          </h3>
          <p className="text-2xl font-bold text-red-600">
            {formatCurrency(monthlyData.reduce((sum, m) => sum + m.outgoing, 0))} CRW
          </p>
        </div>

        <div className="card">
          <h3 className="text-sm font-semibold text-gray-600 dark:text-gray-400 mb-2">
            Net Balance
          </h3>
          <p className="text-2xl font-bold text-blue-600">
            {formatCurrency(monthlyData.reduce((sum, m) => sum + m.balance, 0))} CRW
          </p>
        </div>
      </div>

      {/* Monthly Summary */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-gray-800 dark:text-white">
            Monthly Summary
          </h2>
          {monthlyData.length > 0 && (
            <button
              onClick={() => exportMonthlyReportToCSV(monthlyData, `monthly_report_${new Date().toISOString().split('T')[0]}.csv`)}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-semibold transition"
            >
              📥 Export CSV
            </button>
          )}
        </div>

        {monthlyData.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-400">No transaction data available</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b dark:border-gray-700">
                  <th className="text-left py-2 px-2">Month</th>
                  <th className="text-right py-2 px-2">Incoming</th>
                  <th className="text-right py-2 px-2">Outgoing</th>
                  <th className="text-right py-2 px-2">Fees</th>
                  <th className="text-right py-2 px-2">Net</th>
                </tr>
              </thead>
              <tbody>
                {monthlyData.map((month) => (
                  <tr key={month.month} className="border-b dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700">
                    <td className="py-3 px-2 font-semibold">{getMonthName(month.month)}</td>
                    <td className="py-3 px-2 text-right text-green-600">
                      +{formatCurrency(month.incoming)}
                    </td>
                    <td className="py-3 px-2 text-right text-red-600">
                      -{formatCurrency(month.outgoing)}
                    </td>
                    <td className="py-3 px-2 text-right text-gray-600 dark:text-gray-400">
                      -{formatCurrency(month.fee)}
                    </td>
                    <td className="py-3 px-2 text-right font-semibold">
                      {month.balance >= 0 ? '+' : ''}{formatCurrency(month.balance)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Zakat Report */}
      <div className="card">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-gray-800 dark:text-white">
            Zakat Report
          </h2>
          {zakatData.length > 0 && (
            <button
              onClick={() => exportZakatReportToCSV(zakatData, `zakat_report_${new Date().toISOString().split('T')[0]}.csv`)}
              className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg text-sm font-semibold transition"
            >
              📥 Export CSV
            </button>
          )}
        </div>

        {zakatData.length === 0 ? (
          <p className="text-gray-600 dark:text-gray-400">No zakat data available</p>
        ) : (
          <div className="space-y-4">
            {zakatData.map((zakat, idx) => (
              <div key={idx} className="p-4 border border-gray-200 dark:border-gray-700 rounded">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <p className="text-gray-600 dark:text-gray-400 text-sm">Month</p>
                    <p className="font-semibold text-gray-800 dark:text-white">
                      {getMonthName(zakat.month)}
                    </p>
                  </div>
                  <div>
                    <p className="text-gray-600 dark:text-gray-400 text-sm">Zakat Amount (2.5%)</p>
                    <p className="font-semibold text-green-600">
                      {formatCurrency(zakat.zakatAmount)} CRW
                    </p>
                  </div>
                  <div>
                    <p className="text-gray-600 dark:text-gray-400 text-sm">Deduction Date</p>
                    <p className="font-semibold text-gray-800 dark:text-white">
                      {zakat.deductedDate}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Notes */}
      <div className="card bg-blue-50 dark:bg-blue-900">
        <p className="text-sm text-blue-900 dark:text-blue-100">
          <strong>Note:</strong> Zakat is calculated at 2.5% of your wallet balance and is deducted on the 1st of each month.
          Monthly summaries are based on your transaction history.
        </p>
      </div>
    </div>
  )
}

export default Reports
