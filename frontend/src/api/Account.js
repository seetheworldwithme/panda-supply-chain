// 关于账户信息接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.Account + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.Account + 'getById',id)
    },
    add(data={}) {
        return http.post(path.Account + 'add',data)
    },
    edit(data={}) {
        return http.post(path.Account + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.Account + 'delBatch',{ids:id})
    },
}