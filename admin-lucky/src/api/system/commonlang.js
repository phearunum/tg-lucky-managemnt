import request from '@/utils/request'


export function listCommonLang(query) {
  return request({
    url: 'system/CommonLang/list',
    method: 'get',
    params: query,
  })
}

export function listLangByLocale(locale) {
  return request({
    url: 'system/CommonLang/list/' + locale,
    method: 'get',
  })
}


export function addCommonLang(data) {
  return request({
    url: 'system/CommonLang',
    method: 'post',
    data: data,
  })
}


export function updateCommonLang(data) {
  return request({
    url: 'system/CommonLang',
    method: 'PUT',
    data: data,
  })
}

/**
 * 获取多语言配置详情
 * @param {Id}
 */
export function getCommonLang(id) {
  return request({
    url: 'system/CommonLang/' + id,
    method: 'get'
  })
}
/**
 * @param {key}
 */
export function getCommonLangByKey(key) {
  return request({
    url: 'system/CommonLang/key/' + key,
    method: 'get'
  })
}


/**
 * @param {primary key} pid
 */
export function delCommonLang(pid) {
  return request({
    url: 'system/CommonLang/' + pid,
    method: 'delete'
  })
}

// Export multilingual configuration
export function exportCommonLang(query) {
  return request({
    url: 'system/CommonLang/export',
    method: 'get',
    params: query
  })
}