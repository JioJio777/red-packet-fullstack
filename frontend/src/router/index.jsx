import { createBrowserRouter, Navigate } from 'react-router-dom'
import ProtectedRoute from '../components/ProtectedRoute'
import HomePage from '../pages/HomePage'

// 登录和注册页面后续步骤再补充
const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <HomePage />
      </ProtectedRoute>
    ),
  },
  {
    path: '/login',
    element: <div>登录页（待实现）</div>,
  },
  {
    path: '/register',
    element: <div>注册页（待实现）</div>,
  },
  {
    path: '*',
    element: <Navigate to="/" replace />,
  },
])

export default router
