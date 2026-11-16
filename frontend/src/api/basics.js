// 关于基础维护接口
import http from "./http";
import path from "./path";

export default {
    all(query={}) {  
        return http.post(path.basics + 'formatlist',query)
    },
    list(query={}) {  
        return http.post(path.basics + 'list',query)
    },
    one(id={}) { 
        return http.post(path.basics,id)
    },
    add(data={}) {
        return http.post(path.basics + 'add',data)
    },
    edit(data={}) {
        return http.post(path.basics,data) 
    },
    del(id={}) {
        return http.post(path.basics,id)
    },
}