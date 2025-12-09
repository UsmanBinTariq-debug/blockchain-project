import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { useAuthStore, useWalletStore } from '../utils/store'
import { apiClient } from '../services/api'

interface LoginFormData {
  email: string
  password: string
}

const Login: React.FC = () => {
  const navigate = useNavigate()
  const setToken = useAuthStore((state) => state.setToken)
  const [error, setError] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit, formState: { errors } } = useForm<LoginFormData>()

  const onSubmit = async (data: LoginFormData) => {
    try {
      setLoading(true)
      setError('')
      const response = await apiClient.login(data.email, data.password)
      
      if (response.data.status === 'success') {
        // Extract token from backend response
        const token = response.data.data?.token
        if (!token) {
          setError('No token received from server')
          return
        }

        // Save token to localStorage and store
        localStorage.setItem('auth_token', token)
        setToken(token)

        // Try to load wallet/profile immediately
        try {
          const walletResp = await apiClient.getWallet()
          if (walletResp.data.status === 'success' && walletResp.data.data) {
            useWalletStore.setState({ wallet: walletResp.data.data, balance: walletResp.data.data.balance || 0 })
          }
        } catch (e) {
          console.warn('Failed to load wallet after login', e)
        }

        navigate('/dashboard')
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center">
      <div className="bg-white rounded-lg shadow-xl p-8 w-full max-w-md">
        <h1 className="text-3xl font-bold text-center text-gray-800 mb-8">
          Crypto Wallet
        </h1>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-gray-700 font-semibold mb-2">Email</label>
            <input
              {...register('email', { 
                required: 'Email is required',
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: 'Invalid email address'
                }
              })}
              type="email"
              className="input-field"
              placeholder="your@email.com"
            />
            {errors.email && <span className="error-text">{errors.email.message}</span>}
          </div>

          <div>
            <label className="block text-gray-700 font-semibold mb-2">Password</label>
            <input
              {...register('password', { required: 'Password is required' })}
              type="password"
              className="input-field"
              placeholder="Enter your password"
            />
            {errors.password && <span className="error-text">{errors.password.message}</span>}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full"
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>

        <p className="text-center text-gray-600 mt-6">
          Don't have an account?{' '}
          <a href="/register" className="text-blue-600 hover:underline font-semibold">
            Register here
          </a>
        </p>
      </div>
    </div>
  )
}

export default Login
