// 关于平台...业务核算接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.adjust + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.adjust + 'getById',id)
    },
    add(data={}) {
        return http.post(path.adjust + 'add',data)
    },
    edit(data={}) {
        return http.post(path.adjust + 'updateById',data) 
    },
    del(id={}) {
        return http.post(path.adjust + 'deleteById',id)
    },
    // 查询企业列表
    selEcts(data={}) {
        return http.post(path.adjust + 'selEcts',data)
    },
}