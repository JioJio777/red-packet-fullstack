import client from './client'

export function sendRedPacket(type, totalAmount, totalCount) {
  return client.post('/red-packets', { type, total_amount: totalAmount, total_count: totalCount })
}

export function claimRedPacket(id) {
  return client.post(`/red-packets/${id}/claim`)
}

export function getRedPacketDetail(id) {
  return client.get(`/red-packets/${id}`)
}

export function getRedPacketRecords(id, page = 1, pageSize = 10) {
  return client.get(`/red-packets/${id}/records`, { params: { page, page_size: pageSize } })
}
