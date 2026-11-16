// 关于权限...菜单管理接口
import http from "./http";
import path from "./path";

export default {
    // 获取当前账户对应菜单
    selUserMenus(query={}) {  
        return http.post(path.authMenu + 'selUserMenus',query)
    },
    // 查询全部菜单列表
    listByPage(query={}) {  
        return http.post(path.authMenu + 'listByPage',query)
    },
    // 增加菜单
    add(query={}) {  
        return http.post(path.authMenu + 'add',query)
    },
    // 查询单个菜单信息
    getById(query={}) {  
        return http.post(path.authMenu + 'getById',query)
    },
    // 编辑单个菜单信息
    edit(query={}) {  
        return http.post(path.authMenu + 'updateById',query)
    },
    // 删除菜单
    del(query={}) {  
        return http.post(path.authMenu + 'delBatch',query)
    },
    // 查询全部父级菜单
    list(query={}) {  
        return http.post(path.authMenu + 'list',query)
    },
    // 查询操作权限列表
    selAuthys(query={}) {  
        return http.post(path.authMenu + 'selAuthys',query)
    },

}