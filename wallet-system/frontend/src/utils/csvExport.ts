/**
 * CSV Export Utilities
 * Handles exporting transaction history to CSV format
 */

interface Transaction {
  id: string
  transactionHash: string
  senderWallet: string
  receiverWallet: string
  amount: number
  fee: number
  signature: string
  status: string
  transactionType: string
  createdAt: string
}

/**
 * Export transactions to CSV file
 */
export const exportTransactionsToCSV = (transactions: Transaction[], fileName = 'transactions.csv') => {
  if (!transactions || transactions.length === 0) {
    alert('No transactions to export')
    return
  }

  // Prepare CSV headers
  const headers = [
    'Transaction Hash',
    'Date',
    'From',
    'To',
    'Amount (CRW)',
    'Fee (CRW)',
    'Type',
    'Status',
  ]

  // Prepare CSV rows
  const rows = transactions.map((tx) => [
    tx.transactionHash || '',
    new Date(tx.createdAt).toLocaleString(),
    tx.senderWallet || '',
    tx.receiverWallet || '',
    tx.amount.toFixed(2),
    tx.fee.toFixed(2),
    tx.transactionType || 'normal',
    tx.status || 'pending',
  ])

  // Create CSV content
  const csvContent = [
    headers.join(','),
    ...rows.map((row) => row.map((cell) => `"${cell}"`).join(',')),
  ].join('\n')

  // Download file
  downloadCSV(csvContent, fileName)
}

/**
 * Export monthly reports to CSV file
 */
export const exportMonthlyReportToCSV = (
  monthlyData: Array<{
    month: string
    incoming: number
    outgoing: number
    fee: number
    balance: number
  }>,
  fileName = 'monthly_report.csv'
) => {
  if (!monthlyData || monthlyData.length === 0) {
    alert('No report data to export')
    return
  }

  // Prepare CSV headers
  const headers = ['Month', 'Incoming (CRW)', 'Outgoing (CRW)', 'Fees (CRW)', 'Net (CRW)']

  // Prepare CSV rows
  const rows = monthlyData.map((data) => [
    data.month,
    data.incoming.toFixed(2),
    data.outgoing.toFixed(2),
    data.fee.toFixed(2),
    data.balance.toFixed(2),
  ])

  // Add summary row
  const totals = monthlyData.reduce(
    (acc, data) => ({
      month: 'TOTAL',
      incoming: acc.incoming + data.incoming,
      outgoing: acc.outgoing + data.outgoing,
      fee: acc.fee + data.fee,
      balance: acc.balance + data.balance,
    }),
    { month: 'TOTAL', incoming: 0, outgoing: 0, fee: 0, balance: 0 }
  )

  const csvContent = [
    headers.join(','),
    ...rows.map((row) => row.map((cell) => `"${cell}"`).join(',')),
    [
      totals.month,
      totals.incoming.toFixed(2),
      totals.outgoing.toFixed(2),
      totals.fee.toFixed(2),
      totals.balance.toFixed(2),
    ]
      .map((cell) => `"${cell}"`)
      .join(','),
  ].join('\n')

  // Download file
  downloadCSV(csvContent, fileName)
}

/**
 * Export zakat history to CSV file
 */
export const exportZakatReportToCSV = (
  zakatData: Array<{
    month: string
    zakatAmount: number
    percentage: number
    deductedDate: string
  }>,
  fileName = 'zakat_report.csv'
) => {
  if (!zakatData || zakatData.length === 0) {
    alert('No zakat data to export')
    return
  }

  // Prepare CSV headers
  const headers = ['Month', 'Amount (CRW)', 'Percentage', 'Date']

  // Prepare CSV rows
  const rows = zakatData.map((data) => [
    data.month,
    data.zakatAmount.toFixed(2),
    `${data.percentage}%`,
    new Date(data.deductedDate).toLocaleString(),
  ])

  // Calculate total zakat
  const totalZakat = zakatData.reduce((sum, data) => sum + data.zakatAmount, 0)

  const csvContent = [
    headers.join(','),
    ...rows.map((row) => row.map((cell) => `"${cell}"`).join(',')),
    ['TOTAL ZAKAT', totalZakat.toFixed(2), '', '']
      .map((cell) => `"${cell}"`)
      .join(','),
  ].join('\n')

  // Download file
  downloadCSV(csvContent, fileName)
}

/**
 * Helper function to trigger CSV download
 */
const downloadCSV = (csvContent: string, fileName: string) => {
  // Create blob
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })

  // Create object URL
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)

  // Set download attributes
  link.setAttribute('href', url)
  link.setAttribute('download', fileName)
  link.style.visibility = 'hidden'

  // Append to body and click
  document.body.appendChild(link)
  link.click()

  // Cleanup
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * Export wallet balance snapshot to CSV
 */
export const exportWalletSnapshot = (
  walletAddress: string,
  balance: number,
  zakatStatus: string,
  fileName = 'wallet_snapshot.csv'
) => {
  const timestamp = new Date().toISOString()

  const csvContent = [
    'Wallet Snapshot Report',
    '',
    ['Wallet Address', walletAddress].map((cell) => `"${cell}"`).join(','),
    ['Balance (CRW)', balance.toFixed(2)].map((cell) => `"${cell}"`).join(','),
    ['Zakat Status', zakatStatus].map((cell) => `"${cell}"`).join(','),
    ['Export Date', timestamp].map((cell) => `"${cell}"`).join(','),
  ].join('\n')

  downloadCSV(csvContent, fileName)
}
