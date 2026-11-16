// 关于行政区划接口
import http from "./http";
import path from "./path";

export default {
    // 省  
    all(query={}) {  
        return http.post(path.areaInfo + 'selByPet',query)
    },
    // 市
    city(query={}) {  
        return http.post(path.areaInfo + 'selByParent',query)
    },
    // 县 区
    county(query={}) {  
        return http.post(path.areaInfo + 'selByParent',query)
    },
    // 乡
    village(query={}) {  
        return http.post(path.areaInfo + 'selByParent',query)
    },
    // 注册页面省
    selByParentNoAuth(query={}) {  
        return http.post(path.areaInfo + 'selByParent',query)
    },
    // 注册页面市县区乡
    selByPetNoAuth(query={}) {  
        return http.post(path.areaInfo + 'selByPet',query)
    },
}