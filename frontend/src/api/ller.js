// 平台企业管理...识别身份证等接口
import http from "./http";
import path from "./path";

export default {
    getOcrResult(query={}) {  
        return http.post(path.ller + 'getOcrResult',query)
    },
}