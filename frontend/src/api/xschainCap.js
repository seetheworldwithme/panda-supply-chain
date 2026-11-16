  // 关于销售收款单管理接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.xschainCap + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.xschainCap + 'getById',id)
    },
    add(data={}) {
        return http.post(path.xschainCap + 'add',data)
    },
    edit(data={}) {
        return http.post(path.xschainCap + 'updateById',data) 
    }, 
    del(id=[]) {
        return http.post(path.xschainCap + 'delBatch',{ids:id})
    },
    // 结算集合
    oneAll(id={}) { 
        return http.post(path.xschainCap + 'getSments',id)
    },
    // 系统编码
    code(data={}){
        return http.post(path.xschainCap + 'getSystemCode',data)
    },
    // 合计
    amount(data={}){
        return http.post(path.xschainCap + 'getSumPrice',data)
    }, 
}