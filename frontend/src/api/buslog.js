// 云仓积分兑换接口
import http from "./http";
import path from "./path";

export default {
    // 查询兑换记录
    exchangeRecords(query={}) {   
        return http.post(path.buslog + 'exchangeRecords',query)
    }, 
    one(id={}) { 
        return http.post(path.buslog + 'getById',id)
    },
    add(data={}) {
        return http.post(path.buslog + 'add',data)
    },
    edit(data={}) {
        return http.post(path.buslog + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.buslog + 'delBatch',{ids:id})
    },
}