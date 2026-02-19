import { Button, Card, Form, Input, message, Typography } from 'antd'
import { useNavigate, Link, Navigate } from 'react-router-dom'
import { login as loginApi } from '../api/auth'
import { useAuth } from '../hooks/useAuth'

const { Title } = Typography

export default function LoginPage() {
  const { isAuthenticated, login } = useAuth()
  const navigate = useNavigate()
  const [form] = Form.useForm()

  // 已登录则直接跳首页
  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  async function handleSubmit(values) {
    try {
      const data = await loginApi(values.username, values.password)
      login(data.token)
      navigate('/', { replace: true })
    } catch (err) {
      message.error(err.message)
    }
  }

  return (
    <div style={styles.page}>
      <Card style={styles.card}>
        <Title level={3} style={styles.title}>登录</Title>
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: '请输入用户名' },
              { max: 50, message: '用户名最多 50 个字符' },
            ]}
          >
            <Input placeholder="请输入用户名" autoComplete="username" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少 6 位' },
            ]}
          >
            <Input.Password placeholder="请输入密码" autoComplete="current-password" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" block>
              登录
            </Button>
          </Form.Item>
        </Form>
        <div style={styles.footer}>
          还没有账号？<Link to="/register">去注册</Link>
        </div>
      </Card>
    </div>
  )
}

const styles = {
  page: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: '#f5f5f5',
  },
  card: {
    width: 400,
    boxShadow: '0 2px 12px rgba(0,0,0,0.08)',
  },
  title: {
    textAlign: 'center',
    marginBottom: 24,
  },
  footer: {
    marginTop: 16,
    textAlign: 'center',
    color: '#888',
  },
}
