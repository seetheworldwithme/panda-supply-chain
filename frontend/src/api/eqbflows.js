// 关于e签宝签署流程表接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.eqbflows + 'listByPage',query)
    }, 
    one(id={}) { 
        return http.post(path.eqbflows + 'getById',id)
    },
    add(data={}) {
        return http.post(path.eqbflows + 'add',data)
    },
    edit(data={}) {
        return http.post(path.eqbflows + 'updateById',data) 
    }, 
    del(id=[]) {
        return http.post(path.eqbflows + 'delBatch',{ids:id})
    },
    // 查询发起方信息
    selLanucnInfo(query={}) {  
        return http.post(path.eqbflows + 'selLanucnInfo',query)
    },
    // 发起签署流程
    launchFlow(query={}) {  
        return http.post(path.eqbflows + 'launchFlow',query)
    },
    // 确认签署流程
    confirmFlow(query={}) {  
        return http.post(path.eqbflows + 'v3/confirmFlow',query)
    },
    // 查询签署地址
    selSignUrl(query={}) {  
        return http.post(path.eqbflows + 'v3/selSignUrl',query)
    },
    // 查询签署人状态
    selSignPeople(query={}) {  
        return http.post(path.eqbflows + 'selSignPeople',query)
    },
    // 签署文件下载
    vDownFile(query={}) {  
        return http.post(path.eqbflows + 'v3DownFile',query)
    },
}