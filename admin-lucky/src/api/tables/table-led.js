import request from '@/utils/request'
export function getTables(query) {
    return request({
        url: '/tables/lists',
        method: 'get',
        params: query
    })
}
export const createTables = (data) => {
    return request({
        url: '/tables',
        method: 'post',
        data: data,
    })
}
export function updateTables(data) {
    return request({
        url: '/tables',
        method: 'put',
        data: data
    })
}

export function getTablesById(id) {
    return request({
        url: `/tables/list/${id}`,
        method: 'get',
    })
}

export function deleteTables(id) {
    return request({
        url: `/tables/delete/${id}`,
        method: 'DELETE',
    })
}
export function getSelectTables() {
    return request({
        url: '/tables/select',
        method: 'get',

    })
}
