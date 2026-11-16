// 关于核销记录管理接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.count + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.count + 'getById',id)
    },
    add(data={}) {
        return http.post(path.count + 'add',data)
    },
    edit(data={}) {
        return http.post(path.count + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.count + 'delBatch',{ids:id})
    },
    // 合计金额
    amount(data={}){
        return http.post(path.count + 'getSumPrice',data)
    }, 
}