// 关于小程序...兑换列表接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.busechlog + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.busechlog + 'getById',id)
    },
    edit(data={}) {
        return http.post(path.busechlog + 'updateById',data) 
    },
}