// 关于物流管理...数据表接口
import http from "./http";
import path from "./path";

export default {
    // 列表页
    listByPage(query={}) {  
        return http.post(path.logistihird + 'listByPage',query)
    },
    // 查询企业全部角色
    clearData(query={}) {  
        return http.post(path.logistihird + 'clearData',query)
    },
}