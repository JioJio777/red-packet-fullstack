import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button, Card, List, Spin, Tag, Typography, message, Avatar, Divider } from 'antd'
import { ArrowLeftOutlined, GiftOutlined, UserOutlined } from '@ant-design/icons'
import { getRedPacketDetail, getRedPacketRecords, claimRedPacket } from '../api/redPacket'

const { Title, Text } = Typography

const TYPE_LABEL = { 1: '普通红包', 2: '拼手气红包' }
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

export default function RedPacketDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [detail, setDetail] = useState(null)
  const [records, setRecords] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [claiming, setClaiming] = useState(false)

  useEffect(() => {
    loadDetail()
  }, [id])

  useEffect(() => {
    if (detail) loadRecords()
  }, [id, page])

  async function loadDetail() {
    setLoading(true)
    try {
      const data = await getRedPacketDetail(id)
      setDetail(data)
    } catch (err) {
      message.error(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function loadRecords() {
    try {
      const data = await getRedPacketRecords(id, page)
      setRecords(data.list ?? [])
      setTotal(data.total)
    } catch {
      // 领取记录加载失败不阻断主流程
    }
  }

  async function handleClaim() {
    setClaiming(true)
    try {
      const data = await claimRedPacket(id)
      message.success(`领到红包 ${formatAmount(data.amount)}！`)
      // 重新拉取详情和记录
      await loadDetail()
      setPage(1)
      await loadRecords()
    } catch (err) {
      message.error(err.message)
    } finally {
      setClaimingt(false)
    }
  }

  if (loading) {
    return (
      <div style={styles.center}>
        <Spin size="large" />
      </div>
    )
  }

  if (!detail) return null

  const canClaim = detail.status === 1 && !detail.my_claim?.claimed

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        {/* 顶栏 */}
        <div style={styles.header}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>红包详情</Title>
          <div style={{ width: 72 }} />
        </div>

        {/* 红包信息卡片 */}
        <Card style={styles.heroCard}>
          <div style={styles.heroContent}>
            <GiftOutlined style={styles.heroIcon} />
            <Text style={styles.heroSender}>{detail.sender_name} 发出的红包</Text>
            <Text style={styles.heroType}>{TYPE_LABEL[detail.type]}</Text>
            <Title level={1} style={styles.heroAmount}>
              {formatAmount(detail.total_amount)}
            </Title>
            <div style={styles.heroStats}>
              <Text style={styles.heroStatText}>
                共 {detail.total_count} 个 · 已领 {detail.claimed_count} 个
              </Text>
            </div>
            <div style={{ marginTop: 8 }}>{STATUS_TAG[detail.status]}</div>
          </div>

          {/* 领取按钮 */}
          {canClaim && (
            <Button
              type="primary"
              danger
              size="large"
              block
              loading={claiming}
              onClick={handleClaim}
              style={styles.claimBtn}
            >
              领红包
            </Button>
          )}

          {/* 已领取提示 */}
          {detail.my_claim?.claimed && (
            <div style={styles.claimedTip}>
              <Text type="secondary">
                你已领取 {formatAmount(detail.my_claim.amount)}，
                领取时间 {formatTime(detail.my_claim.claimed_at)}
              </Text>
            </div>
          )}
        </Card>

        {/* 领取记录 */}
        <Card style={styles.recordsCard} title={`领取记录（${total} 人）`}>
          {records.length === 0 ? (
            <Text type="secondary">暂无人领取</Text>
          ) : (
            <List
              dataSource={records}
              pagination={
                total > 10
                  ? { current: page, pageSize: 10, total, onChange: setPage, size: 'small' }
                  : false
              }
              renderItem={(item) => (
                <List.Item
                  extra={
                    <Text strong style={{ color: '#f5222d' }}>
                      {formatAmount(item.amount)}
                    </Text>
                  }
                >
                  <List.Item.Meta
                    avatar={<Avatar icon={<UserOutlined />} />}
                    title={item.receiver_name}
                    description={formatTime(item.claimed_at)}
                  />
                </List.Item>
              )}
            />
          )}
        </Card>

        {/* 过期时间 */}
        <Divider />
        <Text type="secondary" style={{ fontSize: 12 }}>
          红包过期时间：{formatTime(detail.expired_at)}
        </Text>
      </div>
    </div>
  )
}

const styles = {
  page: { minHeight: '100vh', background: '#f5f5f5' },
  center: { minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' },
  container: { maxWidth: 480, margin: '0 auto', padding: '0 16px 32px' },
  header: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '16px 0',
  },
  heroCard: {
    borderRadius: 12,
    background: 'linear-gradient(135deg, #ff4d4f, #cf1322)',
    border: 'none',
    marginBottom: 16,
  },
  heroContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    paddingBottom: 16,
  },
  heroIcon: { fontSize: 48, color: 'rgba(255,255,255,0.8)', marginBottom: 8 },
  heroSender: { color: 'rgba(255,255,255,0.85)', fontSize: 14, marginBottom: 4 },
  heroType: { color: 'rgba(255,255,255,0.65)', fontSize: 12, marginBottom: 8 },
  heroAmount: { color: '#fff', margin: '0 0 4px' },
  heroStats: { marginBottom: 8 },
  heroStatText: { color: 'rgba(255,255,255,0.75)', fontSize: 13 },
  claimBtn: { marginTop: 8 },
  claimedTip: { marginTop: 12, textAlign: 'center' },
  recordsCard: { borderRadius: 12, marginBottom: 16 },
}
