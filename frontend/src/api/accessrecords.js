// 关于熊猫煤仓...访问人数表接口
import http from "./http";
import path from "./path";
export default {
    all(query={}) {  
        return http.post(path.accessrecords + 'listByPage',query)
    },

}