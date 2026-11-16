// 关于仓库表接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.ware + 'listByPage',query)
    },
    edit(data={}) {
        return http.post(path.ware + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.ware + "delBatch",{ids:id})
    },
    list(query={}) {  
        return http.post(path.ware + 'list',query)
    },
    one(id={}) { 
        return http.post(path.ware,id)
    },
    add(data={}) {
        return http.post(path.ware + 'add',data)
    },
}