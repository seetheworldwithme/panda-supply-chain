// 电话访问表接口
import http from "./http";
import path from "./path";

export default {
    list(data={}) { 
        return http.post(path.phonelog + 'list',data)
    },
    listByPage(data={}) { 
        return http.post(path.phonelog + 'listByPage',data)
    },
}