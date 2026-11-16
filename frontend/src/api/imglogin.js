// 关于登录获取手机图形验证码
import http from "./http";
import path from "./path";

export default { 
    imglogin(data={}) {
        return http.post(path.imglogin,data,{responseType: 'blob'})
    }
}