// 关于红包范围接口
import http from "./http";
import path from "./path";

export default {
    // 查询全部
    list(query={}) {  
        return http.post(path.brelrge + 'list',query)
    },
    // 更新红包配置
    updBatch(query={}) {  
        return http.post(path.brelrge + 'updBatch',query)
    },
 
}