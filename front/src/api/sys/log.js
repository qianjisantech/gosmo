import request from '@/utils/axios.js'

export function getOpsLogPage (body) {
  return request({
    url: '/log/ops/page',
    method: 'post',
    data: body
  })
}