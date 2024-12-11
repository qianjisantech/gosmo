import request from '@/utils/axios.js'

export function getCurrentUserInfo () {
  return request({
    url: '/sys/user/info',
    method: 'get'
  })
}

export function updateCurrentUserInfo (body) {
  return request({
    url: '/api/system/user/current/info/update',
    method: 'post',
    data: body
  })
}

export function changeCurrentUserPassword (body) {
  return request({
    url: '/api/system/user/current/password/change',
    method: 'post',
    data: body
  })
}

export function getUserQueryPage (body) {
  return request({
    url: '/sys/user/queryPage',
    method: 'post',
    data: body
  })
}

export function createUser (body) {
  return request({
    url: '/api/system/user/create',
    method: 'post',
    data: body
  })
}

export function updateUser (body) {
  return request({
    url: '/api/system/user/update',
    method: 'post',
    data: body
  })
}

export function deleteUser (parameter) {
  return request({
    url: '/api/system/user/delete',
    method: 'post',
    params: parameter
  })
}

export function batchEnableUser (body) {
  return request({
    url: '/api/system/user/batch/enable',
    method: 'post',
    data: body
  })
}

export function batchDisableUser (body) {
  return request({
    url: '/api/system/user/batch/disable',
    method: 'post',
    data: body
  })
}