// 关于采购合同接口
import http from "./http";
import path from "./path";
export default {
    all(query={}) {  
        return http.post(path.cgContract + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.cgContract + 'getById',id)
    },
    add(data={}) {
        return http.post(path.cgContract + 'add',data)
    },
    edit(data={}) {
        return http.post(path.cgContract + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.cgContract + 'delBatch',{ids:id})
    },
    // 系统编码 
    code(data={}){
        return http.post(path.cgContract + 'getSystemCode',data)
    },
    // 合同编号
    list(data={}){  
        return http.post(path.cgContract + 'list',data)
    },
    // 采购数据
    count(type={}){
        return http.post(path.cgContract + 'selectQuantity',type)
    },
    // 四流合一物流数据
    getLogistics(type={}){
        return http.post(path.cgContract + 'getLogistics',type)
    },
    // 合计
    getSumPrice(type={}){
        return http.post(path.cgContract + 'getSumPrice',type)
    },
    // 四流合一图
    sear(id={}) { 
        return http.post(path.cgContract + 'searchFour',id)
    },
}