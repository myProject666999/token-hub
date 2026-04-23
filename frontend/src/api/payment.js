import request from '@/utils/request'

export function getPaymentMethods(includeDisabled = false) {
  return request({
    url: '/payment-methods',
    method: 'get',
    params: { include_disabled: includeDisabled }
  })
}

export function createPaymentMethod(data) {
  return request({
    url: '/admin/payment-methods',
    method: 'post',
    data
  })
}

export function updatePaymentMethod(id, data) {
  return request({
    url: `/admin/payment-methods/${id}`,
    method: 'put',
    data
  })
}

export function deletePaymentMethod(id) {
  return request({
    url: `/admin/payment-methods/${id}`,
    method: 'delete'
  })
}

export function createRechargeOrder(data) {
  return request({
    url: '/recharge/create',
    method: 'post',
    data
  })
}

export function simulatePayment(orderNo) {
  return request({
    url: `/recharge/simulate/${orderNo}`,
    method: 'post'
  })
}

export function getUserRechargeRecords(params) {
  return request({
    url: '/recharge/records',
    method: 'get',
    params
  })
}

export function getAllRechargeRecords(params) {
  return request({
    url: '/admin/recharge/records',
    method: 'get',
    params
  })
}

export function getPointsRate() {
  return request({
    url: '/points/rate',
    method: 'get'
  })
}

export function setPointsRate(rate) {
  return request({
    url: '/admin/points/rate',
    method: 'put',
    data: { rate }
  })
}

export function getPointsRecords(params) {
  return request({
    url: '/points/records',
    method: 'get',
    params
  })
}
