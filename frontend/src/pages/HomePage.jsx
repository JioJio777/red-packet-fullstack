import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Card, Typography, Spin, message } from 'antd'
import { WalletOutlined, GiftOutlined, UnorderedListOutlined, LogoutOutlined } from '@ant-design/icons'
import { getProfile } from '../api/user'
import { useAuth } from '../hooks/useAuth'

const { Title, Text } = Typography

export default function HomePage() {
  const { logout } = useAuth()
  const navigate = useNavigate()
  const [profile, setProfile] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getProfile()
      .then(setProfile)
      .catch((err) => message.error(err.message))
      .finally(() => setLoading(false))
  }, [])

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  if (loading) {
    return (
      <div style={styles.center}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        {/* 顶栏 */}
        <div style={styles.header}>
          <Title level={4} style={{ margin: 0 }}>红包系统</Title>
          <Button
            type="text"
            icon={<LogoutOutlined />}
            onClick={handleLogout}
          >
            退出登录
          </Button>
        </div>

        {/* 余额卡片 */}
        <Card style={styles.balanceCard}>
          <div style={styles.balanceContent}>
            <WalletOutlined style={styles.walletIcon} />
            <div>
              <Text style={styles.welcomeText}>
                你好，{profile?.username}
              </Text>
              <div style={styles.balanceRow}>
                <Text style={styles.balanceLabel}>账户余额</Text>
                <Title level={2} style={styles.balanceAmount}>
                  ¥ {((profile?.balance ?? 0) / 100).toFixed(2)}
                </Title>
              </div>
            </div>
          </div>
        </Card>

        {/* 功能入口 */}
        <div style={styles.actions}>
          <Card
            hoverable
            style={styles.actionCard}
            onClick={() => navigate('/send')}
          >
            <div style={styles.actionContent}>
              <GiftOutlined style={{ ...styles.actionIcon, color: '#f5222d' }} />
              <Text strong>发红包</Text>
            </div>
          </Card>
          <Card
            hoverable
            style={styles.actionCard}
            onClick={() => navigate('/my-packets')}
          >
            <div style={styles.actionContent}>
              <UnorderedListOutlined style={{ ...styles.actionIcon, color: '#1890ff' }} />
              <Text strong>我的红包</Text>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}

const styles = {
  page: {
    minHeight: '100vh',
    background: '#f5f5f5',
  },
  center: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  container: {
    maxWidth: 480,
    margin: '0 auto',
    padding: '0 16px 32px',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '16px 0',
  },
  balanceCard: {
    borderRadius: 12,
    background: 'linear-gradient(135deg, #ff4d4f, #cf1322)',
    border: 'none',
    marginBottom: 24,
  },
  balanceContent: {
    display: 'flex',
    alignItems: 'center',
    gap: 16,
  },
  walletIcon: {
    fontSize: 48,
    color: 'rgba(255,255,255,0.8)',
  },
  welcomeText: {
    color: 'rgba(255,255,255,0.85)',
    fontSize: 14,
    display: 'block',
    marginBottom: 4,
  },
  balanceRow: {
    display: 'flex',
    flexDirection: 'column',
  },
  balanceLabel: {
    color: 'rgba(255,255,255,0.7)',
    fontSize: 12,
  },
  balanceAmount: {
    color: '#fff',
    margin: 0,
  },
  actions: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr',
    gap: 16,
  },
  actionCard: {
    borderRadius: 12,
    textAlign: 'center',
  },
  actionContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: 8,
    padding: '8px 0',
  },
  actionIcon: {
    fontSize: 32,
  },
}
