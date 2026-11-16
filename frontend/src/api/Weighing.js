// 关于物流磅单接口
import http from "./http";
import path from "./path";
  
export default {
    all(query={}) {  
        return http.post(path.Weighing + 'listByPage',query)
    },
    list(query={}) {  
        return http.post(path.Weighing + 'list',query)
    },
    add(data={}) { 
        return http.post(path.Weighing + 'add',data)
    },
    one(data={}) { 
        return http.post(path.Weighing + 'getById',data)
    },
    edit(data={}) { 
        return http.post(path.Weighing + 'updateById',data)
    },
    del(id=[]) { 
        return http.post(path.Weighing + 'delBatch',{ids:id})
    },
}