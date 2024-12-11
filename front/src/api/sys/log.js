import request from '@/utils/axios.js'

export function getOpsLogQueryPage (body) {
  return request({
    url: '/log/ops/queryPage',
    method: 'post',
    data: body
  })
}