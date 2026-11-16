// 租赁售卖接口
import http from "./http";
import path from "./path";
 
export default {
    listByPage(data={}) { 
        return http.post(path.busEnterpriseLeasing + 'listByPage',data)
    },
    updateById(data={}) {
        return http.post(path.busEnterpriseLeasing + 'updateById',data) 
    },
    getById(id={}) { 
        return http.post(path.busEnterpriseLeasing + 'getById',id)
    },
    del(id={}) {
        return http.post(path.busEnterpriseLeasing + 'deleteById',id)
    },
}