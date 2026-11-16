// 关于首页...欢迎词接口
import http from "./http";
import path from "./path";

export default {
    // 查询
    getById(data={}) {   
        return http.post(path.config + 'getById',data) 
    },
    // 编辑
    edit(data={}) {   
        return http.post(path.config + 'updateById',data) 
    },
    
} 