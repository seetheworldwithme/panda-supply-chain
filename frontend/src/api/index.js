//聚合文件
import updatePwd from './updatePwd'
import code from './code'
import mobilelogin from './mobilelogin'
import imglogin from './imglogin'
import login from './login'
import listByPage from './listByPage'
import cgContract from './cgContract'
import xsContract from './xsContract'
import basics from './basics'
import supplier from './supplier'
import address from './address'
import basicslist from './basicslist'
import logistics from './logistics'
import waybill from './waybill'
import Weighing from './Weighing'
import capital from './capital'
import cgcapital from './cgcapital'
import xscapital from './xscapital'
import chainCap from './chainCap'
import cgchainCap from './cgchainCap'
import xschainCap from './xschainCap'
import addfile from './addfile'
import file from './file'
import Account from './Account'
import count from './count'
import bill from './bill'
import jxbill from './jxbill'
import xxbill from './xxbill'
import settRec from './settRec'
import billRate from './billRate'
import areaInfo from './areaInfo'
import enterData from './enterData'
import authAll from './authAll'
import roleList from './roleList'
import adjust from './adjust'
import statistics from './statistics'
import roleUserList from './roleUserList'
import authMenu from './authMenu'
import commonlog from './commonlog'
import config from './config'
import eqbconfig from './eqbconfig'
import eqbflows from './eqbflows'
import eqbSingers from './eqbSingers'
import wxChat from './wxChat'
import logistihird from './logistihird'
import confirmPay from './confirmPay'
import Transinfo from './Transinfo'

import busrodect from './busrodect'
import Prodect from './Prodect'
import ware from './ware'
import order from './order'
import glconfig from './glconfig'
import buslog from './buslog'
import busopin from './busopin'
import bustag from './bustag'
import busprolog from './busprolog'
import phonelog from './phonelog'
import brelord from './brelord'
import bcmlog from './bcmlog'
import brelrge from './brelrge'
import verifyInfo from './verifyInfo'
import ller from './ller'
import accessrecords from './accessrecords'
import getSignUrl from './getSignUrl'
import busEnterpriseLeasing from './busEnterpriseLeasing'
import busconfig from './busconfig'
import busbenfig from './busbenfig'
import busechlog from './busechlog'










export default{
    updatePwd,      // 修改密码
    code,           // 获取手机号验证码
    mobilelogin,    // 手机号验证码登录
    imglogin,       // 获取图形验证码
    login,          // 账号登录
    listByPage,     // 合同
    cgContract,     // 采购合同
    xsContract,     // 销售合同
    basics,         // 基础维护
    supplier,       // 合同...供应商信息
    address,        // 合同...地址信息
    basicslist,     // 全部基础维护信息
    logistics,      // 物流...订单表
    waybill,        // 物流...运单表
    Weighing,       // 物流...磅单表
    logistihird,    // 物流...数据表
    capital,        // 资金流...结算管理
    cgcapital,      // 资金流...采购结算
    xscapital,      // 资金流...销售结算
    chainCap,       // 资金流...收付款单管理
    cgchainCap,     // 资金流...采购付款单管理
    xschainCap,     // 资金流...销售收款单管理
    addfile,        // 基础...文件表
    file,           // 基础...文件上传表
    Account,        // 资金流...账户信息
    count,          // 资金流...核销记录表
    bill,           // 票据...发票
    jxbill,         // 票据...进项发票
    xxbill,         // 票据...销项发票
    settRec,        // 资金流...结算关联记录表
    billRate,       // 票据...缴税
    areaInfo,       // 基础...行政区划表
    enterData,      // 平台...企业信息列表
    authAll,        // 权限...用户表
    adjust,         // 平台...业务核算
    statistics,     // 首页...基础统计
    roleList,       // 权限...角色列表
    roleUserList,   // 权限...用户列表
    authMenu,       // 权限...菜单管理
    commonlog,      // 首页...操作日志
    config,         // 首页...欢迎词
    eqbconfig,      // 首页...e签宝配置
    eqbflows,       // 首页...e签宝签署流程表
    eqbSingers,     // 合同管理...签署方信息表
    wxChat,         // 小程序...账号登录
    confirmPay,     // 招商银行接口
    Transinfo,      // 交易记录接口

    busrodect,      // 云仓煤种列表
    Prodect,        // 云仓列表
    ware,		    // 仓库表
    order,          // 云仓订单表
    glconfig,       // 平台...全局配置
    buslog,         // 云仓积分记录表
    busopin,        // 云仓意见建议表
    bustag,         // 标签表
    busprolog,      // 用户访问表
    phonelog,       // 电话访问表
    brelord,        // 红包记录表
    bcmlog,         // 煤种操作表
    brelrge,        // 红包范围表
    verifyInfo,     // 煤仓核实信息
    ller,           // 平台企业管理...识别身份证等表
    accessrecords,  // 访问人数表
    getSignUrl,     // 电子签章
    busEnterpriseLeasing, // 租赁售卖
    busconfig,      // 会员管理
    busbenfig,      // 会员权限
    busechlog,      // 兑换列表
}