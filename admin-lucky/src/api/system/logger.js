import request from '@/utils/request'
export function getLogger(query) {
    return request({
        url: '/log/list',
        method: 'get',
        params: query
    })
}