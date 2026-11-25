// 接口地址chainCommonFile
//process.env.NODE_ENV   Vue打包的环境   192.168.0.215:8765   119.45.238.116:8766 

var env = "developmen_test";    // developmen_zsc

if(env == "developmen_test"){

    var baseUrl = "http://127.0.0.1:8765"          // 开发环境(本地:原件为 https://test.coalcloudtech.com/saas,备份见 path.js.orig-backup)
    var coludUrl = "http://127.0.0.1:8765"          // 掌上煤仓开发环境
}else if(env == "developmen_zsc"){

    var baseUrl = "https://zsc.coalcloudtech.com/zsc"          // 准生产开发环境 
    var coludUrl = "https://zsc.coalcloudtech.com/zsc"          // 掌上煤仓准生产开发环境  

}else{
    var baseUrl = "https://scm.coalcloudtech.com/saas"          // 生产环境
    var coludUrl = "https://scm.coalcloudtech.com/saas"          // 掌上煤仓生产环境
} 
  
var areaInfo = "/saas-common/chainAreaInfo/";                                   // 基础...行政区划表
var enterData = '/saas-common/chainBasicEnterprise/';                           // 平台...企业信息表
var glconfig = '/saas-common/chainGlobalConfiguration/';                        // 平台...全局配置表
var common = "/saas-common/chainCommonTempmx/";                                 // 基础维护...模板基础信息
var address = "/saas-common/chainCommonAddress/";                               // 合同...地址信息
var commonlist = '/saas-common/chainCommonDict/';                               // 基础维护...全部基础信息
var addfile = "/saas-common/chainCommonFile/";                                  // 基础...文件表
var file = "/saas-common/upload/";                                              // 基础...文件上传表
var adjust = '/saas-common/chainBasicBussiness/';                               // 平台...业务核算
var statistics = '/saas-common/chainBasicMainPage/';                            // 首页...基础统计
var commonlog = '/saas-common/chainCommonLog/';                                 // 首页...操作日志
var config = '/saas-common/chainScreenConfig/';                                 // 首页...欢迎辞
var eqbconfig = '/saas-common/chainEqbConfig/';                                 // 首页...e签宝配置
var eqbflows = '/saas-common/chainEqbFlow/';                                    // 首页...e签宝签署流程表
var eqbSingers = '/saas-common/chainEqbSingers/';                               // 合同管理...签署方信息表
var ller = '/saas-common/TenOcrController/';                                    // 平台企业管理...识别身份证等表


var supplier = "/saas-contract/chainContractSupplier/";                         // 合同...供应商信息
var listByPage = "/saas-contract/chainContract/";                               // 合同管理...合同
var cgContract = "/saas-contract/chainCgContract/";                             // 合同管理...采购合同
var xsContract = "/saas-contract/chainXsContract/";                             // 合同管理...销售合同

var logistics = '/saas-logistics/chainLogisticsBill/';                          // 物流...订单表
var waybill = '/saas-logistics/chainLogistics/';                                // 物流...运单表
var Weighing = '/saas-logistics/chainLogisticsWeighingList/';                   // 物流...磅单表
var logistihird = '/saas-logistics/chainLogisticsThirdWaybill/';                // 物流...数据表

var capital = '/saas-capital/chainCapitalSettlement/';                          // 资金流...结算管理
var cgcapital = '/saas-capital/chainCapitalCgSettlement/';                      // 资金流...采购结算
var xscapital = '/saas-capital/chainCapitalXsSettlement/';                      // 资金流...销售结算
var chainCap = '/saas-capital/chainCapitalDetailed/';                           // 资金流...收付款单管理
var cgchainCap = '/saas-capital/chainCapitalCgDetailed/';                       // 资金流...采购付款单管理
var xschainCap = '/saas-capital/chainCapitalXsDetailed/';                       // 资金流...销售收款单管理
var Account = '/saas-capital/chainCapitalAccounts/';                            // 资金流...账户信息
var count = '/saas-capital/chainCapitalCount/';                                 // 资金流...核销记录表
var settRec = '/saas-capital/chainBillSettRec/';                                // 资金流...结算关联记录表

var bill = "/saas-bill/chainBill/";                                             // 票据...发票
var jxbill = "/saas-bill/chainCgBill/";                                         // 票据...进项发票
var xxbill = "/saas-bill/chainXsBill/";                                         // 票据...销项发票
var billRate = "/saas-bill/chainBillRate/";                                     // 票据...缴税

var authAll = "/saas-auth/chainAuthUser/";                                      // 权限...用户表
var roleList = "/saas-auth/chainAuthRole/";                                     // 权限...角色权限
var roleUserList = "/saas-auth/chainAuthUserList/";                             // 权限...用户列表
var authMenu = "/saas-auth/chainAuthMenu/";                                     // 权限...菜单管理
var wxChat = "/saas-auth/chainAuthChat/";                                       // 小程序...账号登录

var busechlog = "/saas-warehouse/busScoreExchangeLog/";                         // 小程序...兑换列表

 
var confirmPay = "/saas-capital/chainPayment/";                                 // 招商银行接口
var Transinfo = "/saas-capital/chainCapitalTransinfo/";                         // 交易记录接口


var busrodect = '/saas-warehouse/busProduct/';                                  // 云仓煤种列表
var Prodect = '/saas-warehouse/busBigProdect/';                                 // 云仓列表
var ware = '/saas-warehouse/busWarehouse/';                                  	// 仓库表
var order = '/saas-warehouse/busOrder/';                                  		// 订单信息表
var buslog = '/saas-warehouse/busScoreLog/';                                    // 积分记录表
var busopin = '/saas-warehouse/busOpinion/'; 		                            // 意见建议表
var bustag = '/saas-warehouse/busTag/'; 		                                // 标签表
var busprolog = '/saas-warehouse/busProLog/'; 		                            // 用户访问表
var phonelog = '/saas-warehouse/busPhoneLog/'; 		                            // 电话访问表
var brelord = '/saas-warehouse/busRedEnvelopeRecord/'; 		                    // 红包记录表
var bcmlog = '/saas-warehouse/busCommonLog/'; 		                            // 煤种操作表
var brelrge = '/saas-warehouse/busRedEnvelopeRange/'; 		                    // 红包范围表
var verifyInfo = '/saas-warehouse/busVerifyInfo/'                               // 煤仓核实信息
var accessrecords = '/saas-warehouse/busVisitorRecord/'                         // 访问人数表
var getSignUrl = '/saas-common/chainEqbFlow/'                                   // 电子签章
var busEnterpriseLeasing = '/saas-warehouse/busEnterpriseLeasing/'             // 租赁售卖
var busconfig = '/saas-warehouse/busMemberConfig/'                             // 会员管理
var busbenfig = '/saas-warehouse/busMemberBenefits/'                           // 会员权限


 


export default {

    wxChat:baseUrl + wxChat,                                                    // 小程序...账号登录
    updatePwd: baseUrl + '/saas-auth/chainAuthUser/updatePwd',                  // 角色权限...修改密码
    code: baseUrl + '/saas-auth/code/',                                         // 登录获取验证码接口
    mobilelogin: baseUrl + '/saas-auth/authentication/mobile',                  // 手机号登录接口
    imglogin: baseUrl + '/saas-auth/code/image',                                // 验证码图形接口
    login: baseUrl + '/saas-auth/',                                             // 登录接口
    listByPage: baseUrl + listByPage,                                           // 合同接口
    cgContract: baseUrl + cgContract,                                           // 采购合同接口
    xsContract: baseUrl + xsContract,                                           // 销售合同接口
    basics: baseUrl + common,                                                   // 基础维护接口
    supplier: baseUrl + supplier,                                               // 合同...供应商信息
    address: baseUrl + address,                                                 // 合同...地址信息
    basicslist: baseUrl + commonlist,                                           // 基础维护...全部基础信息
    logistics: baseUrl + logistics,                                             // 物流...订单表
    waybill: baseUrl + waybill,                                                 // 物流...运单表
    Weighing: baseUrl + Weighing,                                               // 物流...磅单表
    logistihird: baseUrl + logistihird,                                         // 物流...数据表
    
    capital: baseUrl + capital,                                                 // 资金流...结算管理
    cgcapital: baseUrl + cgcapital,                                             // 资金流...采购结算
    xscapital: baseUrl + xscapital,                                             // 资金流...销售结算


    chainCap: baseUrl + chainCap,                                               // 资金流...收付款单管理
    cgchainCap: baseUrl + cgchainCap,                                           // 资金流...采购付款单管理
    xschainCap: baseUrl + xschainCap,                                           // 资金流...销售收款单管理
    addfile: baseUrl + addfile,                                                 // 基础...文件表
    file: baseUrl + file,                                                       // 基础...文件上传表
    Account: baseUrl + Account,                                                 // 资金流...账户信息
    count: baseUrl + count,                                                     // 资金流...核销记录表
    bill: baseUrl + bill,                                                       // 票据...发票
    jxbill: baseUrl + jxbill,                                                   // 票据...进项发票
    xxbill: baseUrl + xxbill,                                                   // 票据...销项发票
    settRec: baseUrl + settRec,                                                 // 资金流...结算关联记录表
    billRate: baseUrl + billRate,                                               // 票据...缴税
    areaInfo: baseUrl + areaInfo,                                               // 基础...行政区划表
    enterData: baseUrl + enterData,                                             // 平台...企业信息列表
    glconfig: baseUrl + glconfig,                                              // 平台...全局配置表
    authAll: baseUrl + authAll,                                                 // 权限...用户表
    roleList: baseUrl + roleList,                                               // 权限...角色列表
    roleUserList: baseUrl + roleUserList,                                       // 权限...用户列表
    authMenu: baseUrl + authMenu,                                               // 权限...路由表
    adjust: baseUrl + adjust,                                                   // 平台...业务核算
    statistics:baseUrl + statistics,                                            // 首页...基础统计 
    commonlog:baseUrl + commonlog,                                              // 首页...操作日志
    config:baseUrl + config,                                                    // 首页...欢迎词
    eqbconfig:baseUrl + eqbconfig,                                              // 首页...e签宝配置
    eqbflows:baseUrl + eqbflows,                                                // 首页...e签宝签署流程表
    eqbSingers:baseUrl + eqbSingers,                                            // 合同管理...签署方信息表
    confirmPay:baseUrl + confirmPay,                                            // 招商银行接口
    Transinfo:baseUrl + Transinfo,                                              // 交易记录接口
    ller:baseUrl + ller,                                                        // 平台企业管理...识别身份证等表
    busechlog:coludUrl  + busechlog,                                            // 小程序...兑换列表
    
    
    busrodect: coludUrl + busrodect,                                         	// 云仓煤种列表
    Prodect: coludUrl + Prodect,                                         		// 云仓列表接口
    ware: coludUrl + ware,                                         			    // 仓库表
    order: coludUrl + order,                                         			// 订单信息表
    buslog: coludUrl + buslog,                                         			// 积分记录表
    busopin: coludUrl + busopin, 	                                            // 意见接口
    bustag: coludUrl + bustag, 	                                                // 标签接口
    busprolog: coludUrl + busprolog, 	                                        // 用户访问接口
    phonelog: coludUrl + phonelog, 	                                            // 电话拨打接口
    brelord: coludUrl + brelord, 	                                            // 红包记录接口
    bcmlog: coludUrl + bcmlog, 	                                                // 煤种操作接口
    brelrge: coludUrl + brelrge, 	                                            // 红包范围接口
    verifyInfo: coludUrl + verifyInfo,                                          // 煤仓核实接口
    accessrecords: coludUrl + accessrecords,                                    // 访问人数表
    getSignUrl: coludUrl + getSignUrl,                                          // 电子签章
    busEnterpriseLeasing: coludUrl + busEnterpriseLeasing,                      // 租赁售卖
    busconfig: coludUrl + busconfig,                                            // 会员管理
    busbenfig: coludUrl + busbenfig,                                            // 会员权限
    baseUrl,
    coludUrl,
}