// 关于基础...文件上传表接口
import { Header } from "element-ui";
import http from "./http";
import path from "./path";


export default {
    // 下载
    // download(id) { 
    //     return http.get(path.file + 'downByIdNoAuth?fileName=' + id ,{responseType: 'blob'})
    // },
    all(query={}) {  
        return http.post(path.file + 'listByPage',query)
    },
    one(id={}) { 
        return http.post(path.file + 'getByRelId',id)
    },
    add(data={}) {
        return http.post(path.file + 'add',data)
    },
    edit(data={}) {
        return http.post(path.file + 'updateById',data) 
    },
    del(id=[]) {
        return http.post(path.file + 'delBatch',{ids:id})
    },
    // 图片回显
    imageShow(data={}) {
        return http.post(path.file + 'getImgByRelId',data) 
    },
    // 文件预览
    preFIle(data={}) {
        return http.post(path.file + 'preFIle',data) 
    },

    // pdf文件预览
    pree(data={}) {
        return http.get(path.file + 'previewStream',data) 
    },
}