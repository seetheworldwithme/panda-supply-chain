// 关于首页...操作日志接口
import http from "./http";
import path from "./path";

export default {
    // 查询操作日志
    all(data={}) {   
        return http.post(path.commonlog + 'listByPage',data) 
    },
} 