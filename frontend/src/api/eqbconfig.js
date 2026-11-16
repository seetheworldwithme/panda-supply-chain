// 关于收付款单管理接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.eqbconfig + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.eqbconfig + 'getById',id)
    },
    add(data={}) {
        return http.post(path.eqbconfig + 'add',data)
    },
    edit(data={}) {
        return http.post(path.eqbconfig + 'updateById',data) 
    }, 
    del(id=[]) {
        return http.post(path.eqbconfig + 'delBatch',{ids:id})
    },
}