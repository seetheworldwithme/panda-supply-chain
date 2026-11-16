// 关于全部基础维护接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.basicslist + 'formatlist',query)
    },
    list(query={}) {  
        return http.post(path.basicslist + 'list',query)
    },
    one(id={}) { 
        return http.post(path.basicslist,id)
    },
    add(data={}) {
        return http.post(path.basicslist + 'add',data)
    },
    edit(data={}) {
        return http.post(path.basicslist,data) 
    },
    del(id={}) {
        return http.post(path.basicslist,id)
    },
}