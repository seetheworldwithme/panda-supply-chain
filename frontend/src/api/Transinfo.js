// 关于交易记录接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.Transinfo + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.Transinfo + 'getById',id)
    },
    add(data={}) {
        return http.post(path.Transinfo + 'add',data)
    },
    edit(data={}) {
        return http.post(path.Transinfo + 'updateById',data) 
    }, 
    del(id=[]) {
        return http.post(path.Transinfo + 'delBatch',{ids:id})
    },
    list(data={}) {
        return http.post(path.Transinfo + 'list',data)
    },
    

}