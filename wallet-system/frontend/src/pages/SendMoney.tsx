import React, { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'

interface SendMoneyFormData {
  recipientAddress: string
  amount: number
  fee: number
  note?: string
}

const SendMoney: React.FC = () => {
  const wallet = useWalletStore((state) => state.wallet)
  const balance = useWalletStore((state) => state.balance)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>('')
  const [success, setSuccess] = useState<string>('')
  const { register, handleSubmit, watch, formState: { errors } } = useForm<SendMoneyFormData>()

  const rawAmount = watch('amount') || 0
  const rawFee = watch('fee') || 0

  const toNumber = (v: any) => {
    const n = typeof v === 'number' ? v : parseFloat(String(v || '0'))
    return Number.isFinite(n) ? n : 0
  }

  const amount = toNumber(rawAmount)
  const fee = toNumber(rawFee)

  const onSubmit = async (data: SendMoneyFormData) => {
    try {
      setLoading(true)
      setError('')
      setSuccess('')

      if (!wallet) {
        setError('Wallet not loaded')
        return
      }

      if (amount + fee > balance) {
        setError('Insufficient balance')
        return
      }

      // In production, this would sign the transaction with the private key
      const signature = 'demo-signature'

      const response = await apiClient.sendTransaction(
        (wallet as any).wallet_address || wallet.walletAddress,
        data.recipientAddress,
        Number(data.amount),
        Number(data.fee),
        data.note || '',
        signature
      )

      if (response.data.status === 'success') {
        setSuccess('Transaction sent successfully!')
        
        // Refresh wallet balance
        const walletResp = await apiClient.getWallet()
        if (walletResp.data?.status === 'success') {
          const { useWalletStore } = await import('../utils/store')
          useWalletStore.setState({ 
            wallet: walletResp.data.data, 
            balance: walletResp.data.data.balance || 0 
          })
        }
        
        // Clear success message after 2 seconds
        setTimeout(() => {
          setSuccess('')
        }, 2000)
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to send transaction')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 dark:text-white">Send Money</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Balance Info */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-2">
            Your Balance
          </h2>
          <p className="text-3xl font-bold text-blue-600">{balance.toFixed(8)}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">CRW Available</p>
        </div>

        {/* Transaction Total */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-2">
            Total Cost
          </h2>
          <p className="text-3xl font-bold text-green-600">{(amount + fee).toFixed(8)}</p>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">Amount + Fee</p>
        </div>

        {/* Remaining Balance */}
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-800 dark:text-white mb-2">
            Remaining
          </h2>
          <p className="text-3xl font-bold text-orange-600">
            {(balance - amount - fee).toFixed(8)}
          </p>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">After Transaction</p>
        </div>
      </div>

      {/* Send Form */}
      <div className="card max-w-2xl">
        <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-6">
          Send Transaction
        </h2>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
            {success}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-gray-700 dark:text-gray-300 font-semibold mb-2">
              Recipient Wallet Address
            </label>
            <input
              {...register('recipientAddress', { 
                required: 'Recipient address is required',
                pattern: {
                  value: /^[0-9a-fA-F]{64}$/,
                  message: 'Invalid wallet address'
                }
              })}
              type="text"
              className="input-field"
              placeholder="64-character wallet address"
            />
            {errors.recipientAddress && <span className="error-text">{errors.recipientAddress.message}</span>}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-gray-700 dark:text-gray-300 font-semibold mb-2">
                Amount
              </label>
              <input
                {...register('amount', { 
                  required: 'Amount is required',
                  min: {
                    value: 0.001,
                    message: 'Minimum amount is 0.001'
                  }
                })}
                type="number"
                step="0.00000001"
                className="input-field"
                placeholder="0.00000000"
              />
              {errors.amount && <span className="error-text">{errors.amount.message}</span>}
            </div>

            <div>
              <label className="block text-gray-700 dark:text-gray-300 font-semibold mb-2">
                Network Fee
              </label>
              <input
                {...register('fee', { 
                  required: 'Fee is required',
                  min: {
                    value: 0,
                    message: 'Fee cannot be negative'
                  }
                })}
                type="number"
                step="0.00000001"
                className="input-field"
                placeholder="0.00000001"
              />
              {errors.fee && <span className="error-text">{errors.fee.message}</span>}
            </div>
          </div>

          <div>
            <label className="block text-gray-700 dark:text-gray-300 font-semibold mb-2">
              Note (Optional)
            </label>
            <textarea
              {...register('note')}
              rows={3}
              className="input-field"
              placeholder="Add a note for this transaction"
            />
          </div>

          <button
            type="submit"
            disabled={loading || !wallet}
            className="btn-primary w-full"
          >
            {loading ? 'Sending...' : 'Send Transaction'}
          </button>
        </form>
      </div>
    </div>
  )
}

export default SendMoney
