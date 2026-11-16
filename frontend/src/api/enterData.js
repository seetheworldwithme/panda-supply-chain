// 关于平台...企业信息列表接口
import http from "./http";
import path from "./path";

export default {
    // 企业客户
    count(query={}) {  
        return http.post(path.enterData + 'selectCount',query)
    },
    // 企业信息列表
    enterData(query={}) {   
        return http.post(path.enterData + 'ptlist',query)
    },
    // 企业信息列表
    list(query={}) {  
        return http.post(path.enterData + 'list',query)
    },
    // 企业统计列表
    all(query={}) {  
        return http.post(path.enterData + 'listByPage',query)
    },
    // 单条平台企业信息
    byidData(query={}) {  
        return http.post(path.enterData + 'getById',query)
    },
    // 删除企业信息
    del(id=[]) {
        return http.post(path.enterData + 'delBatch',{ids:id})
    },
    onedel(data={}) {
        return http.post(path.enterData + 'deleteById',data)
    },
    // 根据id更新企业信息
    updateById(query={}) {  
        return http.post(path.enterData + 'updateById',query)
    },
    // 根据id批量更新签约企业
    updBatch(query=[]) {  
        return http.post(path.enterData + 'updBatch',query)
    },
    // 企查查查询企业
    seletcname(query=[]) {  
        return http.post(path.enterData + 'selEtcName',query)
    },
    // 企查查查询企业列表
    seletcinfo(query=[]) {  
        return http.post(path.enterData + 'selEtcInfo',query)
    },
    /*
    申请注册管理员
    */
    // 提交表单
    addNoAuth(query={}) {  
        return http.post(path.enterData + 'add',query)
    },
}


