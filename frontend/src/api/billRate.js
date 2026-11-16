// 关于缴税...发票接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {   
        return http.post(path.billRate + 'listByPage',query)
    }, 
    one(id={}) { 
        return http.post(path.billRate + 'getById',id)
    },
    add(data={}) {
        return http.post(path.billRate + 'add',data)
    },
    edit(data={}) {
        return http.post(path.billRate + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.billRate + 'delBatch',{ids:id})
    },
    /*
    平台
    */ 
    // 纳税金额累计数据
    rateBill(data={}) {
        return http.post(path.billRate + 'getSumPrice',data)
    },
    // 数量...纳税量
    amount(data={}) {
        return http.post(path.billRate + 'countMonBuses',data)
    },
}