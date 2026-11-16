// 关于物流订单接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.logistics + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.logistics + 'getById',id)
    },
    add(data={}) {
        return http.post(path.logistics + 'add',data)
    },
    edit(data={}) {
        return http.post(path.logistics + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.logistics + 'delBatch',{ids:id})
    },
    // 合同集合
    oneAll(id={}) { 
        return http.post(path.logistics + 'getContracts',id)
    },
    // 系统编码
    code(data={}){
        return http.post(path.logistics + 'getSystemCode',data)
    },
    // 物流
    setall(type={}){
        return http.post(path.logistics + 'countLogist',type)
    },
    // 合计
    getSumPrice(type={}){
        return http.post(path.logistics + 'getSumPrice',type)
    },
    //四流合一物流数据
    selAllOrders(query = {}) {
        return http.post(path.logistics + 'selAllOrders', query)
    },
    //四流合一物流条数
    selOrders(type = {}) {
        return http.post(path.logistics + 'selOrders', type)
    },
    /*
    企业
    */


}