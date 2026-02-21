import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Card, List, Tabs, Tag, Typography, Spin, message, Avatar } from 'antd'
import { ArrowLeftOutlined, GiftOutlined } from '@ant-design/icons'
import { getSentRedPackets, getReceivedRedPackets } from '../api/user'

const { Text, Title } = Typography

const STATUS_TAG = {
  1: <Tag color="green">可领取</Tag>,
  2: <Tag color="default">已抢完</Tag>,
  3: <Tag color="warning">已过期</Tag>,
}

function formatAmount(fen) {
  return '¥ ' + (fen / 100).toFixed(2)
}

function formatTime(iso) {
  return new Date(iso).toLocaleString('zh-CN', { hour12: false })
}

function SentList({ navigate }) {
  const [list, setList] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    getSentRedPackets(page)
      .then((data) => {
        setList(data.list ?? [])
        setTotal(data.total)
      })
      .catch((err) => message.error(err.message))
      .finally(() => setLoading(false))
  }, [page])

  if (loading) return <div style={styles.spinWrap}><Spin /></div>

  return (
    <List
      dataSource={list}
      locale={{ emptyText: '还没有发过红包' }}
      pagination={
        total > 10
          ? { current: page, pageSize: 10, total, onChange: setPage, size: 'small' }
          : false
      }
      renderItem={(item) => (
        <List.Item
          style={styles.listItem}
          onClick={() => navigate(`/red-packets/${item.id}`)}
        >
          <List.Item.Meta
            avatar={<Avatar icon={<GiftOutlined />} style={{ background: '#f5222d' }} />}
            title={
              <span>
                {item.type === 1 ? '普通红包' : '拼手气红包'}
                <span style={styles.tagWrap}>{STATUS_TAG[item.status]}</span>
              </span>
            }
            description={formatTime(item.created_at)}
          />
          <div style={styles.amountCol}>
            <Text strong style={{ color: '#f5222d' }}>{formatAmount(item.total_amount)}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>
              {item.total_count - item.remaining_count}/{item.total_count} 个已领
            </Text>
          </div>
        </List.Item>
      )}
    />
  )
}

function ReceivedList({ navigate }) {
  const [list, setList] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    getReceivedRedPackets(page)
      .then((data) => {
        setList(data.list ?? [])
        setTotal(data.total)
      })
      .catch((err) => message.error(err.message))
      .finally(() => setLoading(false))
  }, [page])

  if (loading) return <div style={styles.spinWrap}><Spin /></div>

  return (
    <List
      dataSource={list}
      locale={{ emptyText: '还没有领过红包' }}
      pagination={
        total > 10
          ? { current: page, pageSize: 10, total, onChange: setPage, size: 'small' }
          : false
      }
      renderItem={(item) => (
        <List.Item
          style={styles.listItem}
          onClick={() => navigate(`/red-packets/${item.red_packet_id}`)}
        >
          <List.Item.Meta
            avatar={<Avatar icon={<GiftOutlined />} style={{ background: '#fa8c16' }} />}
            title={`来自 ${item.sender_name} 的红包`}
            description={`领取于 ${formatTime(item.claimed_at)}`}
          />
          <Text strong style={{ color: '#f5222d' }}>{formatAmount(item.amount)}</Text>
        </List.Item>
      )}
    />
  )
}

export default function MyPacketsPage() {
  const navigate = useNavigate()

  const tabs = [
    { key: 'sent', label: '我发出的', children: <SentList navigate={navigate} /> },
    { key: 'received', label: '我收到的', children: <ReceivedList navigate={navigate} /> },
  ]

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        {/* 顶栏 */}
        <div style={styles.header}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={() => navigate('/')}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>我的红包</Title>
          <div style={{ width: 72 }} />
        </div>

        <Card style={styles.card}>
          <Tabs items={tabs} />
        </Card>
      </div>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f5f5f5' },
  container: { maxWidth: 480, margin: '0 auto', padding: '0 16px 32px' },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '16px 0',
  },
  card: { borderRadius: 12 },
  spinWrap: { padding: 40, textAlign: 'center' },
  listItem: { cursor: 'pointer' },
  tagWrap: { marginLeft: 8 },
  amountCol: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-end',
    gap: 2,
  },
}
