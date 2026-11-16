// 关于物流运单接口
import http from "./http";
import path from "./path";
  
export default {
    all(query={}) {  
        return http.post(path.waybill + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.waybill + 'getById',id)
    },
    edit(id={}) { 
        return http.post(path.waybill + 'updateById',id)
    },
    del(id=[]) { 
        return http.post(path.waybill + 'delBatch',{ids:id})
    },
    // 合计
    getSumPrice(type={}){
        return http.post(path.waybill + 'getSumPrice',type)
    },
    // 运单集合
    list(query={}) {  
        return http.post(path.waybill + 'list',query)
    },
    // 物流量金额累计
    billAll(data={}){ 
        return http.post(path.waybill + 'countBuses',data)
    },
    // 物流量
    countMonBuses(data={}){
        return http.post(path.waybill + 'countMonBuses',data)
    },
}