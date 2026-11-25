// ajax相关封装与配置
import axios from 'axios'
import { Message } from 'element-ui';
// import VueRouter from 'vue-router'  //路由
// Vue.component(Message.name, Message);
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
        if (response.status === 200) {
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
                    type:'success',
                })
            }      
            return Promise.resolve(response);    
        } else{            
            return Promise.reject(response);        
        } 
    }, 
    err => {
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