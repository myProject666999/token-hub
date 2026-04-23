import request from '@/utils/request'

export function login(username, password) {
  return request({
    url: '/login',
    method: 'post',
    data: { username, password }
  })
}

export function register(data) {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

export function getProfile() {
  return request({
    url: '/user/profile',
    method: 'get'
  })
}

export function updateProfile(data) {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

export function getUserPoints() {
  return request({
    url: '/user/points',
    method: 'get'
  })
}

export function getUserStatistics() {
  return request({
    url: '/user/statistics',
    method: 'get'
  })
}
