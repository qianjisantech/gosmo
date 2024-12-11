import request from '@/utils/axios'

export function getCommonOptionProject (body) {
    return request({
        url: '/basicdata/common/option/project',
        method: 'post',
        data:body
    })
}
export function getCommonOptionEs (body) {
    return request({
        url: '/basicdata/common/option/es',
        method: 'post',
        data:body
    })
}