// 结算单关联记录接口
import http from "./http";
import path from "./path";

export default {
    // 全部数据
    all(data={}){
        return http.post(path.settRec + 'listByPage',data)
    },
    // 解除关联
    relieve(data={}){
        return http.post(path.settRec + 'relieve',data)
    },
}