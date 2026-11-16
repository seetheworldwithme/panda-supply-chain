// 关于权限...用户列表接口
import http from "./http";
import path from "./path";

export default {
    // 用户关联角色接口
    assignRole(query={}) {  
        return http.post(path.roleUserList + 'assignRole',query)
    }, 
    // 查询用户拥有角色
    selUserRoles(query={}) {  
        return http.post(path.roleUserList + 'selUserRoles',query)
    }, 
}