import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore, useWalletStore } from '../utils/store'

const Profile: React.FC = () => {
  const user = useAuthStore((state) => state.user)
  const wallet = useWalletStore((state) => state.wallet)
  const logout = useAuthStore((state) => state.logout)
  const navigate = useNavigate()
  const [copiedField, setCopiedField] = useState<string | null>(null)

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text)
    setCopiedField(field)
    setTimeout(() => setCopiedField(null), 2000)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Profile</h1>

      {/* User Information */}
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">
          User Information
        </h2>

        {user ? (
          <div className="space-y-6">
            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Full Name</label>
              <p className="text-gray-800 dark:text-white font-semibold text-lg">{user.fullName}</p>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Email</label>
              <p className="text-gray-800 dark:text-white font-semibold">{user.email}</p>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">CNIC</label>
              <p className="text-gray-800 dark:text-white font-semibold">{user.cnic}</p>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Account Status</label>
              <p className="text-gray-800 dark:text-white">
                {user.isVerified ? (
                  <span className="text-green-600 font-semibold">✓ Verified</span>
                ) : (
                  <span className="text-orange-600 font-semibold">⏳ Pending Verification</span>
                )}
              </p>
            </div>
          </div>
        ) : (
          <p className="text-gray-600 dark:text-gray-400">No user data available</p>
        )}
      </div>

      {/* Wallet Information */}
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">
          Wallet Information
        </h2>

        {wallet ? (
          <div className="space-y-6">
            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Wallet Address</label>
              <div className="flex items-center justify-between p-3 bg-gray-100 dark:bg-gray-700 rounded font-mono text-sm break-all">
                <span className="text-gray-800 dark:text-white">{wallet.walletAddress}</span>
                <button
                  onClick={() => copyToClipboard(wallet.walletAddress, 'walletAddress')}
                  className="ml-2 px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-xs whitespace-nowrap"
                >
                  {copiedField === 'walletAddress' ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Balance</label>
              <p className="text-gray-800 dark:text-white font-semibold text-lg">
                {wallet?.balance || '0.00000000'} CRW
              </p>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Zakat Status</label>
              <p className="text-gray-800 dark:text-white">
                {wallet?.balance && wallet.balance > 0 ? (
                  <span className="text-green-600 font-semibold">✓ Deducted this month</span>
                ) : (
                  <span className="text-yellow-600 font-semibold">⏳ Pending deduction on 1st</span>
                )}
              </p>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Wallet ID</label>
              <div className="flex items-center justify-between p-3 bg-gray-100 dark:bg-gray-700 rounded font-mono text-sm break-all">
                <span className="text-gray-800 dark:text-white">{wallet.id}</span>
                <button
                  onClick={() => copyToClipboard(wallet.id, 'walletId')}
                  className="ml-2 px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 text-xs whitespace-nowrap"
                >
                  {copiedField === 'walletId' ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>

            <div>
              <label className="text-gray-600 dark:text-gray-400 block text-sm mb-2">Created</label>
              <p className="text-gray-800 dark:text-white">
                {wallet?.lastUpdated ? new Date(wallet.lastUpdated).toLocaleDateString() : 'N/A'} {wallet?.lastUpdated ? new Date(wallet.lastUpdated).toLocaleTimeString() : ''}
              </p>
            </div>
          </div>
        ) : (
          <p className="text-gray-600 dark:text-gray-400">No wallet data available</p>
        )}
      </div>

      {/* Security Section */}
      <div className="card">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">
          Security
        </h2>

        <div className="space-y-4">
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            Keep your private keys and seed phrases safe. Never share them with anyone.
          </p>

          <div className="p-4 bg-red-50 dark:bg-red-900 rounded border border-red-200 dark:border-red-700">
            <p className="text-sm text-red-900 dark:text-red-100">
              ⚠️ <strong>Warning:</strong> Never share your private key or seed phrase. Anyone with access to these can control your wallet.
            </p>
          </div>

          <button
            onClick={handleLogout}
            className="w-full px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 font-semibold"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  )
}

export default Profile
