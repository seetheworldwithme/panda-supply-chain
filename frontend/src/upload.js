import axios from "axios";
let headers = {Authorization:sessionStorage.getItem('Authorization')};
var identifier = `${Date.now()}${Math.random()}`;//文件唯一标识
var nifs = true;
let ind = 0;
// console.log(axios)
//正常上传
const upload = (url, data, headers = {}) => {
    return new Promise((resolve, reject) => {
        axios({
            url,
            method: "post",
            data,
            headers: {
                ...headers,
                'Content-Type': 'multipart/form-data'
            },
            name: 'file',
        }).then(res => {
            return resolve(res.data)
        }).catch(err => {
            return reject(err)
        })
    })
}
//分片上传
const uploadByPieces = async (url,{ fileName, file }) => {
    // 上传过程中用到的变量
    const chunkSize = 5 * 1024 * 1024; // 5MB一片
    const chunkCount = Math.ceil(file.size / chunkSize); // 总片数
    // 获取当前chunk数据
    const getChunkInfo = (file, index) => {
        let start = index * chunkSize;
        let end = Math.min(file.size, start + chunkSize);
        let chunk = file.slice(start, end);
        return { start, end, chunk };
    };
    // 分片上传接口
    const uploadChunk = (data) => {
        return new Promise((resolve, reject) => {
            axios({
                url,
                method: "post",
                data,
                headers: {
                    ...headers,
                },
                name: 'file',
            }).then(res => {
                // console.log(res)
                // if(res.data.code == 200){
                //     nifs = true;
                //     ind++
                //     funt()
                // }
                return resolve(res.data)
            }).catch(err => {
                return reject(err)
            })
        })
    }
    // 针对单个文件进行chunk上传
    const readChunk = (index) => {
        const { chunk } = getChunkInfo(file, index);
        let fetchForm = new FormData();
        // 生成文件标识，唯一，统一，后端通过同一唯一标识组装视频
        let array = identifier.split('.');//处理唯一标识
        let newentifier = array[0].toString() + array[1].toString() + '.mp4';
        fetchForm.append("file",chunk);
        fetchForm.append("fileName",newentifier);
        fetchForm.append("chunkNumber", (index*1 + 1));
        fetchForm.append("totalChunks", chunkCount);
        return uploadChunk(fetchForm)
    };
    // 针对每个文件进行chunk处理
    const promiseList = [];
    try {
        // promiseList.push(chunkCount[0]);
        for (let index = 0; index < chunkCount; ++index) {
            promiseList.push(readChunk(index));
        }
        const res = await Promise.all(promiseList)
        return res
    }catch (e) {
        return e
    }
    // const funt = () => {
    //     if(nifs){
    //         nifs = false;
    //         readChunk(ind)
    //     } 
    // }
}
 
export { upload, uploadByPieces }