import request from "@/utils/axios.js";

export function getDictList (body) {
    return request({
        url: '/sys/dict/list',
        method: 'post',
        data:body
    })
}