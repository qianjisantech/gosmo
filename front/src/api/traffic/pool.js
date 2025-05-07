import request from '@/utils/request'
export function TrafficPoolPage (body) {
    return request({
        url: '/traffic/pool/page',
        method: 'post',
        data: body
    })
}
