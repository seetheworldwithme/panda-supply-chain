// 关于销售结算接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.xscapital + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.xscapital + 'getById',id)
    },
    add(data={}) {
        return http.post(path.xscapital + 'add',data)
    },
    edit(data={}) {
        return http.post(path.xscapital + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.xscapital + 'delBatch',{ids:id})
    },
    // 系统编码
    code(data={}){
        return http.post(path.xscapital + 'getSystemCode',data)
    },
    // 收付款单集合
    oneAll(id={}) { 
        return http.post(path.xscapital + 'getDetails',id)
    },
    // 销进项发票集合
    twoAll(id={}) { 
        return http.post(path.xscapital + 'getBills',id)
    },
    // 合计
    amount(data={}){
        return http.post(path.xscapital + 'getSumPrice',data)
    },
    // 结算量查询
    count(type={}){
        return http.post(path.xscapital + 'selectQuantity',type)
    },
    // 结算
    setall(type={}){
        return http.post(path.xscapital + 'countSett',type)
    },
    // 结算单关联发票
    relSmentsBills(data={}){
        return http.post(path.xscapital + 'relSmentsBills',data)
    },
    // 查询平台结算总金额
    getPtSumPrice(data={}){
        return http.post(path.xscapital + 'getPtSumPrice',data)
    }

}