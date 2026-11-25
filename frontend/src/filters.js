
import store from '../src/store/index';
var att;
var arr;


// console.log(store.state)
export default{
    install(Vue,options){
        // 保留两位小数
        Vue.filter('ToFiexd', function(value) {
            let time = null;
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            time = (value*1).toFixed(2);
            return time
        }),
        // 银行交易时间过滤
        Vue.filter('orderStrTime', function(value) {
            let time = null;
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            time = value.substr(0, 4) + '-' + value.substr(4, 2) + '-' + value.substr(6, 2) + ' ' + value.substr(8, 2) + ':' + value.substr(10, 2) + ':' + value.substr(12, 2);
            return time
        }),
        
        // 操作类型过滤
        Vue.filter('operationType', function(value) {
            let time = "";
            if(value == 1){
                time = "合同"
            }else if(value == 2){
                time = "发票"
            }else if(value == 3){
                time = "流水"
            }else if(value == 4){
                time = "结算单"
            }else if(value == 5){
                time = "运单"
            }else if(value == 6){
                time = "核销"
            }else if(value == 7){
                time = "税收"
            }else if(value == 8){
                time = "其他"
            }
            return time
        }),
        // 年月日过滤
        Vue.filter('Time', function(value) {
            let time = "";
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            time = value.substring(0,10);
            return time
        }),
        // 年月日时分秒过滤
        Vue.filter('Times', function(value) {
            let time = "";
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            time = value.substring(0,10)  + ' ' + value.substring(11,19)
            return time
        }),
        // 月份过滤
        Vue.filter('timeDate', function(value) {
            let time = "";
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            time = value.substring(0,7);
            return time
        }),
        // 银行卡号过滤
        Vue.filter('cardCode', function(value) {
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            let newStr = value.replace(/\d(?=(?:\d{4})+\b)/g, `$& `);
            return String(newStr)
        }),
        // 金额过滤
        Vue.filter('currency', function(value = '0', currencyType = '',limit=2) {
            let res;
            value = (value*1).toFixed(limit);
            let prev = value.toString().split('.')[0]; // 获取整数部分
            let next = value.toString().split('.')[1];
            res = prev.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') + '.' + next;
            return currencyType + res
        }),
        // 体表盘金额过滤
        Vue.filter('newcurrency', function(value = '0', currencyType = '',limit=2) {
            if (value == null || value === '') { return '' } // 空值不设防会在渲染期抛 TypeError,整表不刷新
            let res;
            value = (value*1).toFixed(limit);
            let prev = value.toString().split('.')[0]; // 获取整数部分
            let next = value.toString().split('.')[1];
            res = prev.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') + '.' + next;
            return currencyType + res;
        }),
        // 价格类型过滤
        Vue.filter('pricetype', function(value) {
            att = sessionStorage.getItem('data'); 
            arr = JSON.parse(att);
            let arrdata = arr.price_type
            return arrdata[value]
        })
        // 产品大类过滤
        Vue.filter('protype', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.pro_type
            return arrdata[value]
        })
        // 产品小类过滤
        Vue.filter('promiddletype', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.pro__middle_type
            return arrdata[value]
        })
        // 付款方式过滤
        Vue.filter('paymentmethod', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.payment_method
            return arrdata[value]
        })
        // 运输模式过滤
        Vue.filter('transportmode', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.transport_mode
            return arrdata[value]
        })
        // 承运工具过滤
        Vue.filter('carrier', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.carrier
            return arrdata[value]
        })
        // 支付方式过滤
        Vue.filter('paymenttype', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.payment_type
            return arrdata[value]
        })
        // 标的指标过滤
        Vue.filter('prosmalltype', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.pro__small_type
            return arrdata[value]
        })
        // 物流平台过滤
        Vue.filter('logisticsplatform', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.logistics_platform
                return arrdata[value]
        })
        // 文件关联类型
        Vue.filter('filerelation', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.file_relation
            return arrdata[value]
        })
        // 用途过滤
        Vue.filter('purpose', function(value) {
            att = sessionStorage.getItem('data'); 
            arr =JSON.parse(att);
            let arrdata = arr.purpose
            return arrdata[value]
        })
    }
}