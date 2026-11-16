// 标签表接口
import http from "./http";
import path from "./path";

export default {
    one(id={}) { 
        return http.post(path.bustag + 'getById',id)
    },
    add(data={}) {
        return http.post(path.bustag + 'add',data)
    },
    edit(data={}) {
        return http.post(path.bustag + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.bustag + 'delBatch',{ids:id})
    },
    all(data={}) {
        return http.post(path.bustag + 'listByPage',data)
    },
    list(data={}) {
        return http.post(path.bustag + 'list',data)
    },
}