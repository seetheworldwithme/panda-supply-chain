// 关于基础...文件表接口
import http from "./http";
import path from "./path";
export default {
    all(query={}) {  
        return http.post(path.addfile + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.addfile + 'getById',id)
    },
    add(data={}) {
        return http.post(path.addfile + 'add',data)
    },
    edit(data={}) {
        return http.post(path.addfile + 'updateById',data) 
    },
    del(id={}) {
        return http.post(path.addfile + 'deleteById',id)
    },
    alldel(id={}){
        return http.post(path.addfile + 'delBatch',id)
    },
    updBthById(id=[]){
        return http.post(path.addfile + 'updBthById',id)
    },
    // 云仓检测报告
    list(data={}) {
        return http.post(path.addfile + 'list',data) 
    },
}