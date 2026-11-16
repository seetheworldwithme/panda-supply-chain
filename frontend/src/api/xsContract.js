// 关于销售合同接口
import http from "./http";
import path from "./path";
export default {
    all(query={}) {  
        return http.post(path.xsContract + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.xsContract + 'getById',id)
    },
    add(data={}) {
        return http.post(path.xsContract + 'add',data)
    },
    edit(data={}) {
        return http.post(path.xsContract + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.xsContract + 'delBatch',{ids:id})
    },
    // 系统编码 
    code(data={}){
        return http.post(path.xsContract + 'getSystemCode',data)
    },
    // 合同编号
    list(data={}){  
        return http.post(path.xsContract + 'list',data)
    },
    // 采购数据
    count(type={}){
        return http.post(path.xsContract + 'selectQuantity',type)
    },
    // 四流合一物流数据
    getLogistics(type={}){
        return http.post(path.xsContract + 'getLogistics',type)
    },
    // 合计
    getSumPrice(type={}){
        return http.post(path.xsContract + 'getSumPrice',type)
    },
    // 四流合一图
    sear(id={}) { 
        return http.post(path.xsContract + 'searchFour',id)
    },
}