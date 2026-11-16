  // 关于采购付款单管理接口
  import http from "./http";
  import path from "./path";
  
  export default {
      all(query={}) {  
          return http.post(path.cgchainCap + 'listByPage',query)
      },
      one(id={}) { 
          return http.post(path.cgchainCap + 'getById',id)
      },
      add(data={}) {
          return http.post(path.cgchainCap + 'add',data)
      },
      edit(data={}) {
          return http.post(path.cgchainCap + 'updateById',data) 
      }, 
      del(id=[]) {
          return http.post(path.cgchainCap + 'delBatch',{ids:id})
      },
      // 结算集合
      oneAll(id={}) { 
          return http.post(path.cgchainCap + 'getSments',id)
      },
      // 系统编码
      code(data={}){
          return http.post(path.cgchainCap + 'getSystemCode',data)
      },
      // 合计
      amount(data={}){
          return http.post(path.cgchainCap + 'getSumPrice',data)
      }, 
  }