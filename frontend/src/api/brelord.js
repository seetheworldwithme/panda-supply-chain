// 红包记录接口
import http from "./http";
import path from "./path";

export default {
    all(data={}) {
        return http.post(path.brelord + 'listByPage',data)
    },
    totoalMoney(data={}) {
        return http.post(path.brelord + 'totoalMoney',data)
    },
  
}