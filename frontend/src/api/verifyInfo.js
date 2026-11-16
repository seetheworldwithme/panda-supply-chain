// 关于云仓订单列表接口
import http from "./http";
import path from "./path";

export default {
    list(query={}) {  
        return http.post(path.verifyInfo + 'list',query)
    },
    all(query={}) {  
        return http.post(path.verifyInfo + 'listByPage',query)
    },
    one(query={}) {  
        return http.post(path.verifyInfo + 'getById',query)
    },
    
}