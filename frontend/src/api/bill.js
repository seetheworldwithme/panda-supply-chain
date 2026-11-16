// 关于票据...发票接口
import http from "./http";
import path from "./path";

export default {
    all(query = {}) {
        return http.post(path.bill + 'listByPage', query)
    },
    one(id = {}) {
        return http.post(path.bill + 'getById', id)
    },
    add(data = {}) {
        return http.post(path.bill + 'add', data)
    },
    edit(data = {}) {
        return http.post(path.bill + 'updateById', data)
    },
    del(id = []) {
        return http.post(path.bill + 'delBatch', { ids: id })
    },
    // 系统编码
    code(data = {}) {
        return http.post(path.bill + 'getSystemCode', data)
    },
    // 结算集合
    oneAll(id = {}) {
        return http.post(path.bill + 'getSments', id)
    },
    // 合计
    amount(data = {}) {
        return http.post(path.bill + 'getSumPrice', data)
    },
    // 仪表盘
    count(type = {}) {
        return http.post(path.bill + 'selectBusCount', type)
    },
    // 进项占比提示看板
    selJxzb(type = {}) {
        return http.post(path.bill + 'selJxzb', type)
    },
    // 税负率提示看板
    selZeroTax(type = {}) {
        return http.post(path.bill + 'selZeroTax', type)
    },
    // 税负率统计 
    countTaxRate(type = {}) {
        return http.post(path.bill + 'countTaxRate', type)
    },
    // 合同关联信息列表
    bill(query = {}) {
        return http.post(path.bill + 'searchRelBill', query)
    },
    // selAllOrders(type = {}) {
    //     return http.post(path.listByPage + 'selAllOrders', type)
    // },
    //根据合同id查询全部订单列表
    // selAllOrders(query = {}) {
    //     return http.post(path.bill + 'selAllOrders', query)
    // },
    // 合同列表
    // list(query={}) {  
    //     return http.post(path.bill + 'list',query)
    // },
    //获取票据数据
    // list(query={}){
    //     return http.post(path.bill+'list',query)
    // },
    /**
     * 平台
     */
    // 总体业务量...金额累计数据接口
    countBuses(query = {}) {
        return http.post(path.bill + 'countBuses', query)
    },
    // 数量...业务量
    business(query = {}) {
        return http.post(path.bill + 'countMonBuses', query)
    },
    // 税负率
    countPtTax(query = {}) {
        return http.post(path.bill + 'countPtTax', query)
    },
    // 统计本年度业务量金额累计
    countYearBuses(query = {}) {
        return http.post(path.bill + 'countYearBuses', query)
    },
}