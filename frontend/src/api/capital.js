// 关于结算管理接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.capital + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.capital + 'getById',id)
    },
    add(data={}) {
        return http.post(path.capital + 'add',data)
    },
    edit(data={}) {
        return http.post(path.capital + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.capital + 'delBatch',{ids:id})
    },
    // 系统编码
    code(data={}){
        return http.post(path.capital + 'getSystemCode',data)
    },
    // 收付款单集合
    oneAll(id={}) { 
        return http.post(path.capital + 'getDetails',id)
    },
    // 销进项发票集合
    twoAll(id={}) { 
        return http.post(path.capital + 'getBills',id)
    },
    // 合计
    amount(data={}){
        return http.post(path.capital + 'getSumPrice',data)
    },
    // 结算量查询
    count(type={}){
        return http.post(path.capital + 'selectQuantity',type)
    },
    // 结算
    setall(type={}){
        return http.post(path.capital + 'countSett',type)
    },
    // 结算单关联发票
    relSmentsBills(data={}){
        return http.post(path.capital + 'relSmentsBills',data)
    },
    // 查询平台结算总金额
    getPtSumPrice(data={}){
        return http.post(path.capital + 'getPtSumPrice',data)
    }

}