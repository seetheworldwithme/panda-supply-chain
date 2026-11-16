// 关于合同接口
import http from "./http";
import path from "./path";

export default {
    all(query = {}) {
        return http.post(path.listByPage + 'listByPage', query)
    },
    one(id = {}) {
        return http.post(path.listByPage + 'getById', id)
    },
    add(data = {}) {
        return http.post(path.listByPage + 'add', data)
    },
    edit(data = {}) {
        return http.post(path.listByPage + 'updateById', data)
    },
    del(id = []) {
        return http.post(path.listByPage + 'delBatch', { ids: id })
    },
    // 系统编码 
    code(data = {}) {
        return http.post(path.listByPage + 'getSystemCode', data)
    },
    // 合同编号
    list(data = {}) {
        return http.post(path.listByPage + 'list', data)
    },
    // 采购数据
    count(type = {}) {
        return http.post(path.listByPage + 'selectQuantity', type)
    },
    // 四流合一物流数据
    getLogistics(type = {}) {
        return http.post(path.listByPage + 'getLogistics', type)
    },
    // // 四流合一物流数据（新）
    // selAllOrders(type = {}) {
    //     return http.post(path.listByPage + 'selAllOrders', type)
    // },
    // 合计
    getSumPrice(type = {}) {
        return http.post(path.listByPage + 'getSumPrice', type)
    },
    // 四流合一图
    sear(id = {}) {
        return http.post(path.listByPage + 'searchFour', id)
    },
    // 删除一键导入
    delOnce(id = {}) {
        return http.post(path.listByPage + 'delOnce', id)
    },
    /*
    平台
    */
    // 交易产品份额
    countPro(type = {}) {
        return http.post(path.listByPage + 'countPro', type)
    },
    // 平台合同量
    countSumContract(type = {}) {
        return http.post(path.listByPage + 'countSumContract', type)
    },
}