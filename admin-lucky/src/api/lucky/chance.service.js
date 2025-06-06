import request from '@/utils/request'

export function getChancePoint(query) {
    return request({
        url: '/account/chance-points/list',
        method: 'get',
        params: query
    })
}
export function saveChancePoint(data) {
    return request({
        url: '/account/chance-points/',
        method: 'post',
        data: data
    })
}
export function updateChancePoint(data) {
    return request({
        url: '/account/chance-points/update',
        method: 'put',
        data: data
    })
}
