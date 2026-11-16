// 关于收付款单管理接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.chainCap + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.chainCap + 'getById',id)
    },
    add(data={}) {
        return http.post(path.chainCap + 'add',data)
    },
    edit(data={}) {
        return http.post(path.chainCap + 'updateById',data) 
    }, 
    del(id=[]) {
        return http.post(path.chainCap + 'delBatch',{ids:id})
    },
    // 结算集合
    oneAll(id={}) { 
        return http.post(path.chainCap + 'getSments',id)
    },
    // 系统编码
    code(data={}){
        return http.post(path.chainCap + 'getSystemCode',data)
    },
    // 合计
    amount(data={}){
        return http.post(path.chainCap + 'getSumPrice',data)
    },
    
    // 收付
    setall(type={}){
        return http.post(path.chainCap + 'countDetail',type)
    },
    // 获取银行名称列表
    banks(data={}){
        return http.post(path.chainCap + 'getBanks',data)
    },   
    // 金额 核销
    writeoff(data={}){
        return http.post(path.chainCap + 'balancingAccounts',data)
    }, 
    // 金额 取消核销
    cancelBalancing(data=[]){
        return http.post(path.chainCap + 'cancelBalancing',data)
    }, 
    // 根据合同查询流水
    getDetails(data={}){
        return http.post(path.chainCap + 'getDetails',data)
    }, 
    // 查询平台流水合计金额
    getPtSumPrice(data={}){
        return http.post(path.chainCap + 'getPtSumPrice',data)
    }, 
}