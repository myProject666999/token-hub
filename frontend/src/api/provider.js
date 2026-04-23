import request from '@/utils/request'

export function getProviders(includeDisabled = false) {
  return request({
    url: '/providers',
    method: 'get',
    params: { include_disabled: includeDisabled }
  })
}

export function getProviderList(params) {
  return request({
    url: '/admin/providers',
    method: 'get',
    params
  })
}

export function getProvider(id) {
  return request({
    url: `/admin/providers/${id}`,
    method: 'get'
  })
}

export function createProvider(data) {
  return request({
    url: '/admin/providers',
    method: 'post',
    data
  })
}

export function updateProvider(id, data) {
  return request({
    url: `/admin/providers/${id}`,
    method: 'put',
    data
  })
}

export function deleteProvider(id) {
  return request({
    url: `/admin/providers/${id}`,
    method: 'delete'
  })
}

export function getModelList(params) {
  return request({
    url: '/admin/models',
    method: 'get',
    params
  })
}

export function getModel(id) {
  return request({
    url: `/admin/models/${id}`,
    method: 'get'
  })
}

export function createModel(data) {
  return request({
    url: '/admin/models',
    method: 'post',
    data
  })
}

export function updateModel(id, data) {
  return request({
    url: `/admin/models/${id}`,
    method: 'put',
    data
  })
}

export function deleteModel(id) {
  return request({
    url: `/admin/models/${id}`,
    method: 'delete'
  })
}
