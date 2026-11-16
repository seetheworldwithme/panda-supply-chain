// 关于招商银行接口
import http from "./http";
import path from "./path";

export default {
    // 付款
    payment(query={}) {  
        return http.post(path.confirmPay + 'payment',query)
    },
    // 交易记录详情
    orderDetails(query={}) {  
        return http.post(path.confirmPay + 'orderDetails',query)
    },
    // 交易对账单查询（收款记录）
    billInfo(query={}) {  
        return http.post(path.confirmPay + 'billInfo',query)
    },
    // 解冻-确认付款
    GoConfirmPay(query={}) {  
        return http.post(path.confirmPay + 'confirmPay',query)
    },
    // 退款
    cancelPay(query={}) {  
        return http.post(path.confirmPay + 'cancelPay',query)
    },
    
}