// 关于权限...角色列表接口
import http from "./http";
import path from "./path";

export default {
    // 新增企业角色
    addQy(query={}) {  
        return http.post(path.roleList + 'addEtcRole',query)
    },
    // 查询企业全部角色
    AllQy(query={}) {  
        return http.post(path.roleList + 'selEtcRoles',query)
    },
    // 删除企业单个角色
    delQy(id={}) {
        return http.post(path.roleList + 'delRoleById',id)
    },
    // 编辑企业单个角色
    editQy(query={}) {  
        return http.post(path.roleList + 'updateById',query)
    }, 
    // 查看企业单个角色
    oneQy(query={}) {  
        return http.post(path.roleList + 'getByRoleId',query)
    }, 
    // 查询角色菜单列表
    selMenus(query={}) {  
        return http.post(path.roleList + 'selMenus',query)
    }, 
    // 保存企业权限分配
    saveEtcAuths(query={}) {  
        return http.post(path.roleList + 'saveEtcAuths',query)
    }, 
    // 根据角色查询操作权限和菜单
    selAuthyByRole(query={}) {  
        return http.post(path.roleList + 'selAuthyByRole',query)
    }, 

    // 新增平台角色
    addPtRole(query={}) {  
        return http.post(path.roleList + 'addPtRole',query)
    },
    // 查询平台全部角色
    selPtRoles(query={}) {  
        return http.post(path.roleList + 'selPtRoles',query)
    },
    // 查看平台单个角色
    getRoleForPt(query={}) {  
        return http.post(path.roleList + 'getRoleForPt',query)
    }, 
    // 编辑平台单个角色
    updRoleForPt(query={}) {  
        return http.post(path.roleList + 'updRoleForPt',query)
    }, 
    
    
}