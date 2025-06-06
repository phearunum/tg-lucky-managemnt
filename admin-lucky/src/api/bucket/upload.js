import request from '@/utils/request'
export function getListUpload(query) {
    return request({
        url: '/upload/list',
        method: 'get',
        params: query
    })
}
export const createupload = (data) => {
    return request({
        url: '/upload',
        method: 'post',
        headers: {
            'Content-Type': 'multipart/form-data'
        },
        data: data,
    })
}
export function updateupload(data) {
    return request({
        url: '/upload',
        method: 'put',
        data: data
    })
}

export function getuploadById(id) {
    return request({
        url: `/upload/${id}`,
        method: 'get',
    })
}

export function deleteUpload(id) {
    return request({
        url: `/upload/delete/${id}`,
        method: 'delete',
    })
}
export function getSelectupload() {
    return request({
        url: '/upload/select',
        method: 'get',

    })
}
