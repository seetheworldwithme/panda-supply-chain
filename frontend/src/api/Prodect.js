// 关于云仓列表接口
import http from "./http";
import path from "./path";

export default {
    list(query={}) {  
        return http.post(path.Prodect + 'list',query)
    },
    alllist(query={}) {  
        return http.post(path.Prodect + 'allProductType',query)
    },
}