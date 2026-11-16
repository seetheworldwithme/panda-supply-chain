// 关于微信登录接口
import http from "./http";
import path from "./path";

export default {
    // 微信登录接口
    wxChat(query={}) {  
        return http.post(path.wxChat + 'loginByAccount',query)
    },
    
    getWxConfig(query={}) {  
        return http.post(path.wxChat + 'getWxConfig',query)
    },
}