// 关于修改密码
import http from "./http";
import path from "./path";

export default {
    updatePwd(data={}) {
        return http.post(path.updatePwd,data)
    }
}