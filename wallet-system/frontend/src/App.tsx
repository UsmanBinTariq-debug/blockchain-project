import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './utils/store'
import Layout from './components/Shared/Layout'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import SendMoney from './pages/SendMoney'
import Wallet from './pages/Wallet'
import BlockExplorer from './pages/BlockExplorer'
import Transactions from './pages/Transactions'
import Reports from './pages/Reports'
import Profile from './pages/Profile'
import Admin from './pages/Admin'

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }

  return <>{children}</>
}

function App() {
  React.useEffect(() => {
    // On app load, if token present try to fetch wallet/profile
    const token = localStorage.getItem('auth_token')
    console.log('[App] useEffect - token:', token ? token.substring(0, 20) + '...' : null)
    // ignore legacy demo-token
    if (token && token !== 'demo-token') {
      console.log('[App] token valid, fetching wallet...')
      import('./services/api').then(({ apiClient }) => {
        apiClient.getWallet().then((resp) => {
          console.log('[App] wallet fetch response:', resp.data)
          if (resp.data?.status === 'success' && resp.data.data) {
            console.log('[App] setting wallet to store:', resp.data.data)
            // Lazy import store to avoid circular
            import('./utils/store').then(({ useWalletStore }) => {
              useWalletStore.setState({ wallet: resp.data.data, balance: resp.data.data.balance || 0 })
            })
          }
        }).catch((err) => {
          console.error('[App] wallet fetch error:', err)
        })
      })
    } else if (token === 'demo-token') {
      // cleanup legacy demo token
      console.log('[App] removing demo-token')
      localStorage.removeItem('auth_token')
    } else {
      console.log('[App] no token in localStorage')
    }
  }, [])
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route
          element={
            <ProtectedRoute>
              <Layout />
            </ProtectedRoute>
          }
        >
          <Route path="/" element={<Navigate to="/dashboard" />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/wallet" element={<Wallet />} />
          <Route path="/send-money" element={<SendMoney />} />
          <Route path="/transactions" element={<Transactions />} />
          <Route path="/explorer" element={<BlockExplorer />} />
          <Route path="/explorer/:hash" element={<BlockExplorer />} />
          <Route path="/reports" element={<Reports />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/admin" element={<Admin />} />
        </Route>
      </Routes>
    </Router>
  )
}

export default App
