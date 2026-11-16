// 关于权限...用户表接口
import http from "./http";
import path from "./path";

export default {
    // 注销账号
    logOffUser(query={}) {  
        return http.post(path.authAll + 'logOffUser',query)
    },
    add(query={}) {  
        return http.post(path.authAll + 'add',query)
    },
    // 查询非企业账号列表
    ptlist(query={}) {  
        return http.post(path.authAll + 'ptlist',query)
    },
    // 生成邀请码
    createCode(query={}) {  
        return http.post(path.authAll + 'createCode',query)
    },
    // 企业设置...授权用户列表查询全部
    listByPage(query={}) {  
        return http.post(path.authAll + 'listByPage',query)
    }, 
    // 平台设置...授权用户列表查询全部
    supList(query={}) {  
        return http.post(path.authAll + 'supList',query)
    }, 
    // 设置...待授权用户列表修改单条
    updateById(query={}) {  
        return http.post(path.authAll + 'updateById',query)
    }, 
    // 设置...待授权用户列表删除单条
    deleteById(query={}) {  
        return http.post(path.authAll + 'deleteById',query)
    }, 
    // 取消经纪人
    cancelBuser(query={}) {  
        return http.post(path.authAll + 'cancelBuser',query)
    }, 
    // 经纪人列表（地区，煤种）
    list(query={}) {  
        return http.post(path.authAll + 'list',query)
    }, 
    /*
    忘记密码
    */
    // 忘记密码手机号验证码
    valCode(data={}){
        return http.post(path.authAll + 'valIdentity',data)
    },
    // 忘记密码提交验证
    valIdentity(data={}){
        return http.post(path.authAll + 'valCode',data)
    },
    // 修改密码
    updPwd(data={}){
        return http.post(path.authAll + 'updPwd',data)
    },
    /*
    注册申请
    */
    // 获取手机验证码  
    getPhoneCodeNoAuth(data={}){
        return http.post(path.authAll + 'getPhoneCode',data)
    },
    // 设置账号信息
    setAccountNoAuth(data={}){
        return http.post(path.authAll + 'setAccount',data)
    },
    // 加入企业
    joinCpyNoAuth(data={}){
        return http.post(path.authAll + 'joinCpy',data)
    },
    /*
    注册e签宝
    */
    // 提交注册
    regNoAuth(data={}){
        return http.post(path.authAll + 'reg',data)
    },
}