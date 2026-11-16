// 关于首页...基础统计接口
import http from "./http";
import path from "./path";

export default {
    // 查询累计实际业务量-物流量-纳税量
    countActBus(data={}) {   
        return http.post(path.statistics + 'countActBus',data) 
    },
    // 按月统计实际业务量-物流量-纳税量
    countMonActBus(data={}) { 
        return http.post(path.statistics + 'countMonActBus',data)
    },
    // 按月统计开票企业数量-资料完成企业数量(客户数占比)
    countMonData(data={}) { 
        return http.post(path.statistics + 'countMonData',data)
    },   
    // 按月统计网络货运占比
    countNetLogist(data={}) { 
        return http.post(path.statistics + 'countNetLogist',data)
    },  
    // 企业客户
    count(query={}) {  
        return http.post(path.statistics + 'selectCount',query)
    },
    // 企业信息列表
    enterData(query={}) {  
        return http.post(path.statistics + 'ptlist',query)
    },
    // 统计企业所得税占比
    countEtcTax(query={}) {  
        return http.post(path.statistics + 'countEtcTax',query)
    },
    // 物流收发货地轨迹
    countLogistArea(query={}){
        return http.post(path.statistics + 'countLogistArea',query)
    },
    // 物流平台轨迹
    countPtLogistArea(query={}){
        return http.post(path.statistics + 'countPtLogistArea',query)
    },
    // 查询价格走势图
    countPriceTrend(query={}){
        return http.post(path.statistics + 'countPriceTrend',query)
    },
    // 企业简介
    selEtcMsg(query={}){
        return http.post(path.statistics + 'selEtcMsg',query)
    }
}