// 用户访问表接口
import http from "./http";
import path from "./path";

export default {
    selVisisted(data={}) { 
        return http.post(path.busprolog + 'selVisisted',data)
    },
    add(data={}) {
        return http.post(path.busprolog + 'add',data)
    },
    selCountByUser(data={}) {
        return http.post(path.busprolog + 'selCountByUser',data)
    },
    selMemberData(data={}) {
        return http.post(path.busprolog + 'selMemberData',data)
    },
    selVisitedPop(data={}) {
        return http.post(path.busprolog + 'selVisitedPop',data)
    },
    selViata(data={}) {
        return http.post(path.busprolog + 'selVisitedData',data)
    },
    selCallData(data={}) {
        return http.post(path.busprolog + 'selCallData',data)
    },
}