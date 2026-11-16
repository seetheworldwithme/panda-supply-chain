// 关于全局配置表接口
import http from "./http";
import path from "./path";

export default {
    list(data={}) {  
        return http.post(path.glconfig + 'list',data)
    },
    updateById(data={}) {  
        return http.post(path.glconfig + 'updateById',data)
    },
}