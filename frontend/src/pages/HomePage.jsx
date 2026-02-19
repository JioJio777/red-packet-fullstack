import { useAuth } from '../hooks/useAuth'

export default function HomePage() {
  const { logout } = useAuth()

  return (
    <div style={{ padding: 40 }}>
      <h2>首页（占位）</h2>
      <p>登录成功，欢迎来到红包系统。</p>
      <button onClick={logout} style={{ marginTop: 16 }}>退出登录</button>
    </div>
  )
}
