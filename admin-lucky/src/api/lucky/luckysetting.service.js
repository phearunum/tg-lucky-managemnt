import request from '@/utils/request'

export function getLuckySetting(query) {
    return request({
        url: '/account/lucky-setting-winner/list',
        method: 'get',
        params: query
    })
}
export function saveLuckySetting(data) {
    return request({
        url: '/account/lucky-setting-winner/',
        method: 'post',
        data: data
    })
}
export function updateLuckySetting(data) {
    return request({
        url: '/account/lucky-setting-winner/update',
        method: 'put',
        data: data
    })
}

