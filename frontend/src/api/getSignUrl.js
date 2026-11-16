// 关于e签宝合同管理...签署方信息表
import http from "./http";
import path from "./path";

export default {
    getSignUrl(query={}) {  
        return http.post(path.getSignUrl + 'getSignUrl',query)
    },
    v3DownFile(query={}) {  
        return http.post(path.getSignUrl + 'v3DownFile',query)
    },
}