// 关于登录接口封装
import http from "./http";
import path from "./path";

var headers = {headers: {'Content-Type': 'application/x-www-form-urlencoded'}} 
export default {
    mobilelogin(data={}) {
        return http.post(path.mobilelogin,data,headers)
    }
}