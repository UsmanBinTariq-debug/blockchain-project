import React from 'react'
import { Outlet, useNavigate, NavLink } from 'react-router-dom'
import { useAuthStore, useUIStore } from '../../utils/store'

const Layout: React.FC = () => {
  const navigate = useNavigate()
  const logout = useAuthStore((state) => state.logout)
  const isDarkMode = useUIStore((state) => state.isDarkMode)
  const toggleDarkMode = useUIStore((state) => state.toggleDarkMode)
  const sidebarOpen = useUIStore((state) => state.sidebarOpen)
  const toggleSidebar = useUIStore((state) => state.toggleSidebar)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className={isDarkMode ? 'dark' : ''}>
      <div className="flex h-screen bg-gray-50 dark:bg-gray-900">
        {/* Sidebar */}
        <aside className={`${sidebarOpen ? 'w-64' : 'w-20'} relative bg-gray-800 text-white transition-all duration-300`}>
          <div className="p-6 flex items-center justify-between">
            {sidebarOpen && <h1 className="text-2xl font-bold">CWallet</h1>}
            <button onClick={toggleSidebar} className="p-2 hover:bg-gray-700 rounded">
              ☰
            </button>
          </div>

          <nav className="mt-8 space-y-2">
            {[
              { label: 'Dashboard', path: '/dashboard', icon: '📊' },
              { label: 'Wallet', path: '/wallet', icon: '💰' },
              { label: 'Send Money', path: '/send-money', icon: '📤' },
              { label: 'Transactions', path: '/transactions', icon: '📋' },
              { label: 'Explorer', path: '/explorer', icon: '🔍' },
              { label: 'Reports', path: '/reports', icon: '📈' },
              { label: 'Profile', path: '/profile', icon: '👤' },
              { label: 'Admin', path: '/admin', icon: '⚙️' },
            ].map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                className={({ isActive }) =>
                  `flex items-center px-4 py-3 rounded transition ${isActive ? 'bg-gray-700' : 'hover:bg-gray-700'}`
                }
              >
                <span className="text-xl">{item.icon}</span>
                {sidebarOpen && <span className="ml-3">{item.label}</span>}
              </NavLink>
            ))}
          </nav>

          <div className="absolute bottom-0 left-0 right-0 p-4 space-y-2 border-t border-gray-700">
            <button
              onClick={toggleDarkMode}
              className="w-full flex items-center justify-center px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded transition"
            >
              {isDarkMode ? '☀️' : '🌙'}
            </button>
            <button
              onClick={handleLogout}
              className="w-full flex items-center justify-center px-4 py-2 bg-red-600 hover:bg-red-700 rounded transition"
            >
              {sidebarOpen ? 'Logout' : '🚪'}
            </button>
          </div>
        </aside>

        {/* Main Content */}
        <main className="flex-1 overflow-auto">
          <div className="p-8">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}

export default Layout
