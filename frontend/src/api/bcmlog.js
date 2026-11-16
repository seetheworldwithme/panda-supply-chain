// 煤种操作接口
import http from "./http";
import path from "./path";

export default {
    all(data={}) {
        return http.post(path.bcmlog + 'listByPage',data)
    },
  
}