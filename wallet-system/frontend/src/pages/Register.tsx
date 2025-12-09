import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { apiClient } from '../services/api'

interface RegisterFormData {
  email: string
  fullName: string
  cnic: string
  password: string
  confirmPassword: string
}

const Register: React.FC = () => {
  const navigate = useNavigate()
  const [error, setError] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit, watch, formState: { errors } } = useForm<RegisterFormData>()

  const password = watch('password')

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setLoading(true)
      setError('')

      if (data.password !== data.confirmPassword) {
        setError('Passwords do not match')
        return
      }

      const response = await apiClient.register(
        data.email,
        data.fullName,
        data.cnic,
        data.password
      )

      if (response.data.status === 'success') {
        navigate('/login', { state: { message: 'Registration successful! Please login.' } })
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl p-8 w-full max-w-md">
        <h1 className="text-3xl font-bold text-center text-gray-800 mb-8">
          Create Account
        </h1>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-gray-700 font-semibold mb-2">Full Name</label>
            <input
              {...register('fullName', { required: 'Full name is required' })}
              type="text"
              className="input-field"
              placeholder="John Doe"
            />
            {errors.fullName && <span className="error-text">{errors.fullName.message}</span>}
          </div>

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
            <label className="block text-gray-700 font-semibold mb-2">CNIC</label>
            <input
              {...register('cnic', { 
                required: 'CNIC is required',
                pattern: {
                  value: /^\d{5}-\d{7}-\d{1}$/,
                  message: 'Invalid CNIC format (12345-1234567-1)'
                }
              })}
              type="text"
              className="input-field"
              placeholder="12345-1234567-1"
            />
            {errors.cnic && <span className="error-text">{errors.cnic.message}</span>}
          </div>

          <div>
            <label className="block text-gray-700 font-semibold mb-2">Password</label>
            <input
              {...register('password', { 
                required: 'Password is required',
                minLength: {
                  value: 8,
                  message: 'Password must be at least 8 characters'
                }
              })}
              type="password"
              className="input-field"
              placeholder="Enter your password"
            />
            {errors.password && <span className="error-text">{errors.password.message}</span>}
          </div>

          <div>
            <label className="block text-gray-700 font-semibold mb-2">Confirm Password</label>
            <input
              {...register('confirmPassword', { 
                required: 'Please confirm your password',
                validate: (value) => value === password || 'Passwords do not match'
              })}
              type="password"
              className="input-field"
              placeholder="Confirm your password"
            />
            {errors.confirmPassword && <span className="error-text">{errors.confirmPassword.message}</span>}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="btn-primary w-full"
          >
            {loading ? 'Creating Account...' : 'Register'}
          </button>
        </form>

        <p className="text-center text-gray-600 mt-6">
          Already have an account?{' '}
          <a href="/login" className="text-blue-600 hover:underline font-semibold">
            Login here
          </a>
        </p>
      </div>
    </div>
  )
}

export default Register
