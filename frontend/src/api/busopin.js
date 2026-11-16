// 云仓意见建议接口
import http from "./http";
import path from "./path";

export default {
    one(id={}) { 
        return http.post(path.busopin + 'getById',id)
    },
    add(data={}) {
        return http.post(path.busopin + 'add',data)
    },
    edit(data={}) {
        return http.post(path.busopin + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.busopin + 'delBatch',{ids:id})
    },
    all(data={}) {
        return http.post(path.busopin + 'listByPage',data)
    },
    
}