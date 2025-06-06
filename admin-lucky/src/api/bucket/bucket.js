import request from '@/utils/request'
export function getBucket(query) {
    return request({
        url: '/bucket/list',
        method: 'get',
        params: query
    })
}
export const createBucket = (data) => {
    return request({
        url: '/bucket',
        method: 'post',
        data: data,
    })
}
export function updateBucket(data) {
    return request({
        url: '/bucket',
        method: 'put',
        data: data
    })
}

export function getBucketById(id) {
    return request({
        url: `/bucket/list/${id}`,
        method: 'get',
    })
}

export function deleteBucket(id) {
    return request({
        url: `/bucket/delete/${id}`,
        method: 'delete',
    })
}
export function getSelectBucket() {
    return request({
        url: '/bucket/select',
        method: 'get',

    })
}
