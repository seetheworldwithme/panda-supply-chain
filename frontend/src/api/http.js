// ajax相关封装与配置
import axios from 'axios'
import { Message } from 'element-ui';
// import VueRouter from 'vue-router'  //路由
// Vue.component(Message.name, Message);
// 任意请求失败(HTTP 错误,或业务失败 success=false)时通知页面:各页面普遍只在 code==200 分支里复位
// loading,失败后「加载中 / 保存中」会一直挂着,只能刷新。main.js 里的全局 mixin 负责解开。
const failHandlers = []
export const onApiFail = fn => { failHandlers.push(fn) }
export const offApiFail = fn => { const i = failHandlers.indexOf(fn); if (i >= 0) failHandlers.splice(i, 1) }
const emitFail = () => failHandlers.slice().forEach(fn => { try { fn() } catch (e) { console.error(e) } })

// 写请求防重复提交:同一个"写"接口、请求体完全相同的请求,在上一次还没返回时再次发出,直接丢弃后发的。
// 各页面的保存 / 核销 / 支付按钮普遍没有防重,双击会建出两条单据;这里在网络层统一兜底。
// 只针对写接口(按路径末段识别),查询接口并发重入是正常用法,不去重。
const WRITE_RE = /\/(?:add|updateById|updBatch|delBatch|deleteById|del|edit|save\w*|import\w*|writeoff|payment|launchFlow|confirmFlow|placeOrder|receiveAmount|generateAmount|setAccount|joinCpy|reg|updPwd|resetPwd|assignRole|saveEtcAuths|savePtAuths|updByOpenid|updByOpenId|updHeadImg|signIn|pay\w*|withdraw\w*|care\w*|like\w*|collect\w*|follow\w*|cancel\w*|relSmentsBills|relieve|clearData)$/i
const inflight = {} // 在途的写请求: key -> 发出时间
function writeKey(config) {
    const data = config.data, tag = Object.prototype.toString.call(data)
    if (config.method !== 'post' || !WRITE_RE.test(config.url || '') || data == null || typeof data !== 'object' || tag === '[object URLSearchParams]' || tag === '[object FormData]') {
        return ''
    }
    return config.url + '|' + JSON.stringify(data)
}
const settle = config => { if (config && config.dedupeKey) delete inflight[config.dedupeKey] }

const instance = axios.create({});
instance.interceptors.request.use(
    config => {
        if(sessionStorage.getItem("Authorization")){
            config.headers.Authorization = sessionStorage.getItem("Authorization")  // 将接口返回的token信息配置到接口请求中
            config.headers.routeName = sessionStorage.getItem('routeName');
        }else{
            config.headers.Signature = '328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f';
            config.headers.AuthorizationId = 'wxda0b1c468b83bc99';
            if(sessionStorage.getItem('routeName')){
                config.headers.routeName = sessionStorage.getItem('routeName');
            }
        }
        const key = writeKey(config)
        if (key) {
            if (inflight[key] && Date.now() - inflight[key] < 30000) return new Promise(() => {}) // 重复提交:页面只跟随第一次请求的结果
            inflight[key] = Date.now()
            config.dedupeKey = key
        }
        return config;
    },
    error => {
        return Promise.reject(error);
    }
);
instance.interceptors.response.use(    
    response => {   
        // 如果返回的状态码为200，说明接口请求成功，可以正常拿到数据     
        // 否则的话抛出错误
        settle(response.config)
        if (response.status === 200) {
           if(response.data && response.data.success === false){
               emitFail()
           }
           // 网关级未授权:HTTP 200 + body code 403 且无 data 键(token 失效/账号被停用)。
           // 原先只弹一条提示,用户被搁浅在当前页,刷新后白屏:这里清掉登录态并整页回到登录页。
           if(response.data && response.data.code === 403 && response.data.data === undefined){
               sessionStorage.removeItem('Authorization');
               window.location.reload();
               return Promise.reject(response);
           }
           // 二进制响应(图形验证码、文件流)没有 msg/code,不能当业务消息弹
           if(response.data && Object.prototype.toString.call(response.data) !== '[object Blob]' && response.data.msg != '成功' && response.data.code !=701){
                Message({
                    message:response.data.msg,
                    type:'warning', // 业务失败的消息原先用绿色「成功」样式弹出,有误导
                })
            }      
            return Promise.resolve(response);    
        } else{            
            return Promise.reject(response);        
        } 
    }, 
    err => {
        settle(err && err.config)
        emitFail()
        if (err && err.response) {
            switch (err.response.status) {
                case 602: err.message = err.response.status ; break;
                case 999: err.message = '请求错误(999)' ; break;
                case 401: err.message = '未授权，请退出重新登录(401)'; break;
                case 403: err.message = '请求未授权，访问失败！'; break;
                case 408: err.message = '请求超时(408)'; break;
                case 500: err.message = '服务器错误(500)'; break;
                case 501: err.message = '服务未实现(501)'; break;
                case 502: err.message = '网络错误(502)'; break;
                case 503: err.message = '服务不可用(503)'; break;
                case 504: err.message = '网络超时(504)'; break;
                case 505: err.message = 'HTTP版本不受支持(505)'; break;
                default: err.message = `连接出错(${err.response.status})!`;
            }
        }else{
            err.message = '连接服务器失败!'
        }
        Message({
            message:err.message,
            type:'error',
        })
        return Promise.reject(err);
    }   
)
export default instance