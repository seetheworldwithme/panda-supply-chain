// 会员管理表接口
import http from "./http";
import path from "./path";

export default {
    list(data={}) { 
        return http.post(path.busconfig + 'list',data)
    },
    edit(data={}) { 
        return http.post(path.busconfig + 'updateById',data)
    },
}