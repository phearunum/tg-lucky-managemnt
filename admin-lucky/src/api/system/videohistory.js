import request from '@/utils/request'

export function videoRecordsList(query) {
    return request({
        url: '/account/video/list',
        method: 'get',
        params: query
    })
}
export function videoRecordsByID(id) {
    return request({
        url: `/account/video/${id}`,
        method: 'get',
    })
}

export function addvideoRecords(data) {
    return request({
        url: '/account/video/',
        method: 'post',
        data: data
    })
}

export function updatevideoRecords(data) {
    return request({
        url: '/account/video/update',
        method: 'put',
        data: data
    })
}

export function deletevideoRecords() {
    return request({
        url: '/account/video/delete',
        method: 'delete',
        data: data
    })
}