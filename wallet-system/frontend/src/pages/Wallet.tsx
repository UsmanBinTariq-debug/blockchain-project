import { useState, useEffect } from 'react'
import { QRCodeSVG as QRCode } from 'qrcode.react'
import { useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'

const Wallet = () => {
  const wallet = useWalletStore((state) => state.wallet)
  const setWallet = useWalletStore((state) => state.setWallet)
  const [copied, setCopied] = useState(false)
  const [showQR, setShowQR] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadWalletData()
  }, [])

  const loadWalletData = async () => {
    try {
      setLoading(true)
      const token = localStorage.getItem('auth_token')
      console.log('[Wallet] Loading wallet data - token:', token ? token.substring(0, 20) + '...' : 'NO TOKEN')
      
      if (!token) {
        console.error('[Wallet] No auth token found')
        setLoading(false)
        return
      }

      const response = await apiClient.getWallet()
      console.log('[Wallet] API response:', response.data)
      console.log('[Wallet] Response data keys:', Object.keys(response.data?.data || {}))
      
      if (response.data?.status === 'success' && response.data.data) {
        // Map snake_case from backend to camelCase for frontend
        const walletData = {
          walletAddress: response.data.data.wallet_address,
          userId: response.data.data.user_id,
          walletId: response.data.data.wallet_id,
          balance: response.data.data.balance,
          lastUpdated: response.data.data.last_updated,
        }
        console.log('[Wallet] Mapped wallet data:', walletData)
        setWallet(walletData)
      } else {
        console.error('[Wallet] Invalid response:', response.data)
      }
    } catch (error) {
      console.error('[Wallet] Failed to load wallet:', error)
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Wallet</h1>

      {loading ? (
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p className="text-gray-600 dark:text-gray-400">Loading wallet data...</p>
          </div>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Wallet Address */}
            <div className="card">
              <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-4">
                Wallet Address
              </h2>
              <div className="space-y-4">
                <div className="p-4 bg-gray-100 dark:bg-gray-700 rounded-lg break-all font-mono text-sm">
                  {wallet?.walletAddress ? wallet.walletAddress : 'Unable to load wallet address'}
                </div>
                <button
                  onClick={() => copyToClipboard(wallet?.walletAddress || '')}
                  className="btn-primary w-full"
                  disabled={!wallet?.walletAddress}
                >
                  {copied ? 'Copied!' : 'Copy Address'}
                </button>
                <button
                  onClick={() => setShowQR(!showQR)}
                  className="btn-secondary w-full"
                  disabled={!wallet?.walletAddress}
                >
                  {showQR ? 'Hide QR Code' : 'Show QR Code'}
                </button>
              </div>
              
              {/* QR Code Display */}
              {showQR && wallet?.walletAddress && (
                <div className="mt-6 flex flex-col items-center">
                  <div className="p-4 bg-white rounded-lg">
                    <QRCode
                      value={wallet.walletAddress}
                      size={200}
                      level="H"
                      includeMargin={true}
                    />
                  </div>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-2 text-center">
                    Scan to share your wallet address
                  </p>
                </div>
              )}
            </div>

            {/* Quick Actions */}
            <div className="card">
              <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-4">
                Quick Actions
              </h2>
              <div className="space-y-2">
                <button 
                  onClick={() => copyToClipboard(wallet?.walletAddress || '')} 
                  className="btn-secondary w-full"
                  disabled={!wallet?.walletAddress}
                >
                  Copy Wallet ID
                </button>
                <button className="btn-secondary w-full">View Transactions</button>
              </div>
            </div>
          </div>

          {/* Wallet Information */}
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-800 dark:text-white mb-4">
              Wallet Information
            </h2>
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-gray-600 dark:text-gray-400">User ID</label>
                  <p className="text-gray-800 dark:text-white font-mono text-sm">
                    {wallet?.userId?.substring(0, 8)}...
                  </p>
                </div>
                <div>
                  <label className="text-gray-600 dark:text-gray-400">Last Updated</label>
                  <p className="text-gray-800 dark:text-white">
                    {wallet?.lastUpdated ? new Date(wallet.lastUpdated).toLocaleDateString() : 'N/A'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  )
}

export default Wallet
