// 关于e签宝合同管理...签署方信息表
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.eqbSingers + 'list',query)
    },
}