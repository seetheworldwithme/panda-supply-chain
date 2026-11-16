// 关于登录获取手机验证码封装
import http from "./http";
import path from "./path";

// var headers = {headers: {'Content-Type': 'application/x-www-form-urlencoded'}}
export default {
    // 通过账号获取验证码
    code(data={}) {
        return http.post(path.code + 'phone',data)
    },
    // 贸易服务验证码
    refundMoney(data={}) {
        return http.post(path.code + 'refundMoney',data)
    }
}