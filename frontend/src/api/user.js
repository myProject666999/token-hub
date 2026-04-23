import request from '@/utils/request'

export function getUserList(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function updateUserStatus(id, status) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data: { status }
  })
}

export function getAPIKeys() {
  return request({
    url: '/api-keys',
    method: 'get'
  })
}

export function createAPIKey(name) {
  return request({
    url: '/api-keys',
    method: 'post',
    data: { name }
  })
}

export function deleteAPIKey(id) {
  return request({
    url: `/api-keys/${id}`,
    method: 'delete'
  })
}

export function updateAPIKeyStatus(id, status) {
  return request({
    url: `/api-keys/${id}/status`,
    method: 'put',
    data: { status }
  })
}

export function getUserCallLogs(params) {
  return request({
    url: '/call-logs',
    method: 'get',
    params
  })
}

export function getUserCallStatistics() {
  return request({
    url: '/call-logs/statistics',
    method: 'get'
  })
}

export function getAllCallLogs(params) {
  return request({
    url: '/admin/call-logs',
    method: 'get',
    params
  })
}

export function getDailyStatistics(params) {
  return request({
    url: '/admin/call-logs/daily',
    method: 'get',
    params
  })
}

export function getOverviewStatistics() {
  return request({
    url: '/admin/call-logs/overview',
    method: 'get'
  })
}
