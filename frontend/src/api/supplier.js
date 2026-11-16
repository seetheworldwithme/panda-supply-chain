// 关于合同...供应商信息接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.supplier + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.supplier + 'getById',id)
    },
    add(data={}) {
        return http.post(path.supplier + 'add',data)
    },
    edit(data={}) {
        return http.post(path.supplier + 'updateById',data) 
    },
    del(id={}) {
        return http.post(path.supplier + 'delBatch',{ids:id})
    },
}