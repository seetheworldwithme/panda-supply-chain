// 关于采购结算接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.cgcapital + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.cgcapital + 'getById',id)
    },
    add(data={}) {
        return http.post(path.cgcapital + 'add',data)
    },
    edit(data={}) {
        return http.post(path.cgcapital + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.cgcapital + 'delBatch',{ids:id})
    },
    // 系统编码
    code(data={}){
        return http.post(path.cgcapital + 'getSystemCode',data)
    },
    // 收付款单集合
    oneAll(id={}) { 
        return http.post(path.cgcapital + 'getDetails',id)
    },
    // 销进项发票集合
    twoAll(id={}) { 
        return http.post(path.cgcapital + 'getBills',id)
    },
    // 合计
    amount(data={}){
        return http.post(path.cgcapital + 'getSumPrice',data)
    },
    // 结算量查询
    count(type={}){
        return http.post(path.cgcapital + 'selectQuantity',type)
    },
    // 结算
    setall(type={}){
        return http.post(path.cgcapital + 'countSett',type)
    },
    // 结算单关联发票
    relSmentsBills(data={}){
        return http.post(path.cgcapital + 'relSmentsBills',data)
    },
    // 查询平台结算总金额
    getPtSumPrice(data={}){
        return http.post(path.cgcapital + 'getPtSumPrice',data)
    }

}