import request from '@/utils/request'

export function deskSettingList(query) {
    return request({
        url: '/account/desk-settings/list',
        method: 'get',
        params: query
    })
}
export function deskSettingByID(id) {
    return request({
        url: `/account/desk-settings/${id}`,
        method: 'get',
    })
}

export function adddeskSetting(data) {
    return request({
        url: '/account/desk-settings/',
        method: 'post',
        data: data
    })
}

export function updatedeskSetting(data) {
    return request({
        url: '/account/desk-settings/update',
        method: 'put',
        data: data
    })
}

export function deletedeskSetting(id) {
    return request({
        url: `/account/desk-settings/delete/${id}`,
        method: 'delete',

    })
}