// 关于票据...进项发票接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.xxbill + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.xxbill + 'getById',id)
    },
    add(data={}) {
        return http.post(path.xxbill + 'add',data)
    },
    edit(data={}) {
        return http.post(path.xxbill + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.xxbill + 'delBatch',{ids:id})
    },
    // 系统编码
    code(data={}){
        return http.post(path.xxbill + 'getSystemCode',data)
    },
    // 结算集合
    oneAll(id={}) { 
        return http.post(path.xxbill + 'getSments',id)
    },
    // 合计
    amount(data={}){
        return http.post(path.xxbill + 'getSumPrice',data)
    },
}