// 关于登录接口封装
import http from "./http";
import path from "./path";

var headers = {headers: {'Content-Type': 'application/x-www-form-urlencoded'}}
export default {
    login(data={}) {
        return http.post(path.login + "login",data,headers)
    },
    screenLogin(data={}) {
        return http.post(path.login + "screenLogin",data,headers)
    },
}