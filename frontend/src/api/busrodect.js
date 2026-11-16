// 关于云仓煤种列表接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.busrodect + 'listByPage',query)
    },
    edit(data={}) {
        return http.post(path.busrodect + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.busrodect + "delBatch",{ids:id})
    },
    // 云仓大屏
    warelist(data={}) {
        return http.post(path.busrodect + "warelist",data)
    },
    
    list(query={}) {  
        return http.post(path.busrodect + 'list',query)
    },
    one(id={}) { 
        return http.post(path.busrodect + 'getById',id)
    },
    add(data={}) {
        return http.post(path.busrodect + 'add',data)
    },
    // 批量上下架云仓
    updateBatch(data={}) {
        return http.post(path.busrodect + 'updateBatch',data)
    },
    
}