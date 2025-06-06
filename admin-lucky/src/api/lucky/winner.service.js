import request from '@/utils/request'

export function getWinner(query) {
    return request({
        url: '/account/lucky-winner/list',
        method: 'get',
        params: query
    })
}
export function saveWinner(data) {
    return request({
        url: '/account/lucky-winner/',
        method: 'post',
        data: data
    })
}
export function updateWinner(data) {
    return request({
        url: '/account/lucky-winner/update',
        method: 'put',
        data: data
    })
}

export function sendNotifWinner(data) {
    return request({
        url: '/account/lucky-winner/winner-send-notif',
        method: 'post',
        data: data
    })
}