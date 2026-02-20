import client from './client'

export function getProfile() {
  return client.get('/user/profile')
}
