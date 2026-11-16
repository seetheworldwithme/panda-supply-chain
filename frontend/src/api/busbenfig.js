
// 会员权限表接口
import http from "./http";
import path from "./path";

export default {
    list(data={}) { 
        return http.post(path.busbenfig + 'list',data)
    },
    edit(data={}) { 
        return http.post(path.busbenfig + 'updateById',data)
    },
}