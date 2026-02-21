import client from './client'

export function getProfile() {
  return client.get('/user/profile')
}

export function getSentRedPackets(page = 1, pageSize = 10) {
  return client.get('/user/red-packets/sent', { params: { page, page_size: pageSize } })
}

export function getReceivedRedPackets(page = 1, pageSize = 10) {
  return client.get('/user/red-packets/received', { params: { page, page_size: pageSize } })
}
