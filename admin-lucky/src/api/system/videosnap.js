import request from '@/utils/request'

export function videoSettingList(query) {
    return request({
        url: '/account/videos-settings/list',
        method: 'get',
        params: query
    })
}
export function videoSettingByID(id) {
    return request({
        url: `/account/videos-settings/${id}`,
        method: 'get',
    })
}

export function addVideoSetting(data) {
    return request({
        url: '/account/videos-settings/',
        method: 'post',
        data: data
    })
}

export function updateVideoSetting(data) {
    return request({
        url: '/account/videos-settings/update',
        method: 'put',
        data: data
    })
}

export function deleteVideoSetting(id) {
    return request({
        url: `/account/videos-settings/delete/${id}`,
        method: 'delete',

    })
}