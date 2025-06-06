import request from '@/utils/request'
//time-records
export function getTimeRecords(query) {
    return request({
        url: '/account/user-requests/time-records',
        method: 'get',
        params: query
    })
}
//bot-location
export function getBotLocationRecords(query) {
    return request({
        url: '/account/user-requests/bot-location',
        method: 'get',
        params: query
    })
}
export function saveBotLocation(data) {
    return request({
        url: '/account/user-requests/bot-location',
        method: 'post',
        data: data
    })
}
export function updateBotLocation(data) {
    return request({
        url: '/account/user-requests/bot-location',
        method: 'put',
        data: data
    })
}



export function getTelegramSetting(query) {
    return request({
        url: '/account/user-requests/setting',
        method: 'get',
        params: query
    })
}
export function saveTelegramSetting(data) {
    return request({
        url: '/account/user-requests/setting',
        method: 'post',
        data: data
    })
}
export function updateTelegramSetting(data) {
    return request({
        url: '/account/user-requests/setting',
        method: 'put',
        data: data
    })
}
export function getTelegram(query) {
    return request({
        url: '/account/user-requests/list',
        method: 'get',
        params: query
    })
}
export const createTelegram = (data) => {
    return request({
        url: '/account/user-requests',
        method: 'post',
        data: data,
    })
}
export function updateTelegram(data) {
    return request({
        url: '/account/user-requests',
        method: 'put',
        data: data
    })
}

export function getTelegramById(id) {
    return request({
        url: `/account/user-requests/list/${id}`,
        method: 'get',
    })
}

export function deleteTelegram(id) {
    return request({
        url: `/account/user-requests/delete/${id}`,
        method: 'delete',
    })
}
export function getSelectTelegram() {
    return request({
        url: '/account/user-requests/select',
        method: 'get',

    })
}



//bot-location
export function getPhoneRecords(query) {
    return request({
        url: '/account/user-requests/bot-phone',
        method: 'get',
        params: query
    })
}
export function savePhone(data) {
    return request({
        url: '/account/user-requests/bot-phone',
        method: 'post',
        data: data
    })
}
export function updatePhone(data) {
    return request({
        url: '/account/user-requests/bot-phone',
        method: 'put',
        data: data
    })
}
export function deletePhone(data) {
    return request({
        url: '/account/user-requests/bot-phone/delete',
        method: 'post',
        data: data
    })
}
