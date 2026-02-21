import { Button, Card, Form, InputNumber, Radio, Typography, message } from 'antd'
import { ArrowLeftOutlined, GiftOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { sendRedPacket } from '../api/redPacket'

const { Title, Text } = Typography

export default function SendRedPacketPage() {
  const navigate = useNavigate()
  const [form] = Form.useForm()

  async function handleSubmit(values) {
    // 页面输入单位是"元"，发送给后端前转成"分"
    const totalAmountFen = Math.round(values.totalAmount * 100)
    try {
      const data = await sendRedPacket(values.type, totalAmountFen, values.totalCount)
      message.success('红包发送成功！')
      // 发送成功后跳转到该红包的详情页
      navigate(`/red-packets/${data.id}`)
    } catch (err) {
      message.error(err.message)
    }
  }

  return (
    <div style={styles.page}>
      <div style={styles.container}>
        {/* 顶栏 */}
        <div style={styles.header}>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/')}
          >
            返回
          </Button>
          <Title level={4} style={{ margin: 0 }}>发红包</Title>
          <div style={{ width: 72 }} />
        </div>

        <Card style={styles.card}>
          <div style={styles.iconRow}>
            <GiftOutlined style={styles.giftIcon} />
          </div>

          <Form
            form={form}
            layout="vertical"
            onFinish={handleSubmit}
            initialValues={{ type: 1, totalCount: 1 }}
          >
            <Form.Item name="type" label="红包类型">
              <Radio.Group>
                <Radio.Button value={1}>普通红包（每人等额）</Radio.Button>
                <Radio.Button value={2}>拼手气红包（随机）</Radio.Button>
              </Radio.Group>
            </Form.Item>

            <Form.Item
              name="totalAmount"
              label="总金额（元）"
              rules={[
                { required: true, message: '请输入总金额' },
                { type: 'number', min: 0.01, message: '金额至少 0.01 元' },
              ]}
            >
              <InputNumber
                style={{ width: '100%' }}
                min={0.01}
                precision={2}
                prefix="¥"
                placeholder="请输入总金额"
              />
            </Form.Item>

            <Form.Item
              name="totalCount"
              label="红包个数"
              rules={[
                { required: true, message: '请输入红包个数' },
                { type: 'number', min: 1, max: 100, message: '个数范围 1 ~ 100' },
              ]}
            >
              <InputNumber
                style={{ width: '100%' }}
                min={1}
                max={100}
                precision={0}
                placeholder="请输入红包个数"
              />
            </Form.Item>

            {/* 每份金额提示（仅普通红包显示） */}
            <Form.Item noStyle shouldUpdate>
              {() => {
                const type = form.getFieldValue('type')
                const amount = form.getFieldValue('totalAmount')
                const count = form.getFieldValue('totalCount')
                if (type === 1 && amount > 0 && count > 0) {
                  return (
                    <div style={styles.hint}>
                      <Text type="secondary">
                        每人可领 ¥ {(amount / count).toFixed(2)}
                      </Text>
                    </div>
                  )
                }
                return null
              }}
            </Form.Item>

            <Form.Item style={{ marginTop: 24, marginBottom: 0 }}>
              <Button type="primary" htmlType="submit" block size="large" danger>
                发红包
              </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>
    </div>
  )
}

const styles = {
  page: {
    minHeight: '100vh',
    background: '#f5f5f5',
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
  card: {
    borderRadius: 12,
  },
  iconRow: {
    textAlign: 'center',
    marginBottom: 24,
  },
  giftIcon: {
    fontSize: 56,
    color: '#f5222d',
  },
  hint: {
    marginTop: -16,
    marginBottom: 16,
    textAlign: 'right',
  },
}
