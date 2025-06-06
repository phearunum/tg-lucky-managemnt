import request from '@/utils/request'

export function getTelegramMember(query) {
    return request({
        url: '/account/telegram-lucky/list',
        method: 'get',
        params: query
    })
}
export function saveTelegramMember(data) {
    return request({
        url: '/account/telegram-lucky/',
        method: 'post',
        data: data
    })
}
export function updateTelegramMember(data) {
    return request({
        url: '/account/telegram-lucky/update',
        method: 'put',
        data: data
    })
}

