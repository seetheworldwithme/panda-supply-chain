import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../views/Login.vue'
import DPLogin from '../views/DPLogin.vue'
import ELogin from '../views/ELogin.vue'
import wxindex from '../views/wxindex.vue'

import NotFont from '../views/NotFont.vue'
import { Message } from 'element-ui'
import privacys from '../views/privacys.vue'
import sergrent from '../views/sergrent.vue'



Vue.use(VueRouter)
// 路由菜单状态：0为显示，1为隐藏，2为显示禁用状态
// 静态路由表
export const staticRouterMap=[
  {
    // 登录首页 
    path: '/Login',
    name: 'Login',
    component: Login 
  },
  {
    // 大屏登录首页 
    path: '/DPLogin',
    name: 'DPLogin',
    component: DPLogin
  },
  { 
    // E签宝
    path: '/ELogin',
    name: 'ELogin',
    component: ELogin
  }, 
  { 
    // 微信app
    path: '/wxindex',
    name: 'wxindex',
    component: wxindex
  },
  { 
    // 微信app隐私协议
    path: '/privacys',
    name: 'privacys',
    component: privacys
  },
  { 
    // 微信app服务协议
    path: '/sergrent',
    name: 'sergrent',
    component: sergrent
  },
  
  // {
  //   // 找不到页面跳转404
  //   path: '/NotFont',
  //   name: 'NotFont', 
  //   component: NotFont
  // },
  {
    // E签宝首页
    path:'/EsignIndex',
    name:'EsignIndex',
    component: ()=> import('../views/Emain/EsignIndex.vue'),
  },
  {
    // E签宝电子签章列表
    path:'/Sponsor',
    name:'Sponsor',
    component: ()=> import('../views/Emain/Sponsor.vue'),
  },
  {
    // E签宝发起签署
    path:'/InitiateStart ',
    name:'InitiateStart',
    component: ()=> import('../views/Emain/InitiateStart.vue'),
  },
  {
    // E签宝暂无数据页面
    path:'/Unknown',
    name:'Unknown',
    component: ()=> import('../views/Emain/Unknown.vue'),
  },
  {
    // 登录......平台服务协议
    path: '/Service',
    name: 'Service',
    component:resolve => require(['@/views/Service.vue'],resolve)
  },
  {
    // 登录......隐私政策
    path: '/Policy',
    name: 'Policy',
    component:resolve => require(['@/views/Policy.vue'],resolve)
  },
  {
    // 登录......立即注册
    path: '/Register',
    name: 'Register',
    component:resolve => require(['@/views/Register.vue'],resolve)
  },
  {
    // 登录......忘记密码
    path: '/Forget',
    name: 'Forget',
    component:resolve => require(['@/views/Forget.vue'],resolve)
  },
  {
    // 支付......支付成功
    path: '/Paymentyes',
    name: 'Paymentyes',
    component:resolve => require(['@/views/Paymentyes.vue'],resolve)
  },
  {
    // 大屏Main主页面
    path:'/Dpmain',
    name:'Dpmain',
    // meta:{isShow:0,title:'主页'},
    // redirect:{name:'Mindex'},
    component: ()=> import('../views/DPmain/Dpmain.vue')
  },
  {
    // 南京大屏Main主页面
    path:'/NjMain',
    name:'NjMain',
    // meta:{isShow:0,title:'主页'},
    // redirect:{name:'Mindex'},
    component: ()=> import('../views/DPmain/NjMain.vue')
  },
  {
    // 第一屏看板页面
    path:'/Firstscreen',
    name:'Firstscreen',
    component: ()=> import('../views/DPmain/Firstscreen.vue'),
  },
  {
    // 第二屏看板页面
    path:'/Secondscreen',
    name:'Secondscreen',
    component: ()=> import('../views/DPmain/Secondscreen.vue'),
  },
  {
    // 第三屏看板页面
    path:'/Thirdscreen',
    name:'Thirdscreen',
    component: ()=> import('../views/DPmain/Thirdscreen.vue'),
  },
  {
    // 第四屏欢迎页面
    path:'/Welcome',
    name:'Welcome',
    component: ()=> import('../views/DPmain/Welcome.vue'),
  },
  {
    // 第五屏地图页面
    path:'/Maps',
    name:'Maps',
    component: ()=> import('../views/DPmain/Maps.vue'),
  },
  {
    // 第六屏云仓页面
    path:'/warehouse',
    name:'warehouse',
    component: ()=> import('../views/DPmain/warehouse.vue'),
  }
]
// 企业动态路由表
export const activeRouterMap = [
  {
    // Main主页面
    path:'/Main',
    name:'Main',
    meta:{isShow:0,title:'主页'},
    redirect:{name:''},
    component: ()=> import('../views/Main.vue'),
    children:[
      {
        // 首页
        path:'Mindex',
        name:'Mindex',
        meta:{isShow:0,title:'首 页',icon:'el-icon-alishouye'},
        redirect:{name:''},
        component:()=>import('../views/AsidsMain.vue'),
        children:[
          {
            // 首页......仪表盘
            path:'MainPage',
            name:'MainPage',
            meta:{isShow:0,title:'仪表盘',icon:'el-icon-alishouye',routeName:'MainPage'},
            component:resolve => require(['@/views/Mindex/MainPage.vue'],resolve)
          },
          {
            // 首页......项目管理
            path:'Menter',
            name:'Menter',
            meta:{isShow:2,title:'项目管理'},
          },
          {
            // 首页......四流合一
            path:'Relation',
            name:'Relation',
            meta:{isShow:0,title:'四流合一',icon:'el-icon-aliweibiaoti-4',routeName:'Relation'},
            component:resolve => require(['@/views/Mindex/Relation.vue'],resolve)
          },
          // {
          //   // 首页......供应链业务
          //   path:'ChainBus',
          //   name:'ChainBus',
          //   meta:{isShow:0,title:'供应链业务',icon:'el-icon-aliweibiaoti-4',routeName:'ChainBus'},
          //   redirect:{name:'Relation'},
          //   component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
          //   children:[
          //     {
          //       // 首页......供应链业务.....四流合一
          //       path:"Relation",
          //       name:"Relation",
          //       meta:{isShow:0,title:'四流合一',routeName:'Relation'},
          //       component:resolve => require(['@/views/Mindex/Relation.vue'],resolve)
          //     },
          //     // {
          //     //   // 首页......供应链业务.....供应链业务内容
          //     //   path:"ChainBus",
          //     //   name:"ChainBus",
          //     //   meta:{isShow:0,title:'供应链业务',routeName:'ChainBus'},
          //     //   component:()=> import('../views/Mindex/ChainBus.vue')
          //     // }
          //   ]
          // },
          // {
          //   // 首页......供应链金融
          //   path:'ChainFinance',
          //   name:'ChainFinance',
          //   meta:{isShow:0,title:'供应链金融',icon:'el-icon-aliVector',routeName:'ChainFinance'},
          //   component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
          //   children:[
          //     // {
          //     //   //首页......供应链金融.....供应链金融内容
          //     //   path:"ChainFinance",
          //     //   name:"ChainFinance",
          //     //   meta:{isShow:0,title:'供应链金融',routeName:'ChainFinance'},
          //     //   component:()=> import('../views/Mindex/ChainFinance.vue')
          //     // },
          //   ]
          // },
          // {
          //   // 首页......业务管理
          //   path:'Xservice',
          //   name:'Xservice',
          //   meta:{isShow:2,title:'业务管理'},
          // },
          // {
          //   //首页......基础维护
          //   path:'Base',
          //   name:'Base',
          //   meta:{isShow:0,title:'基础维护',icon:'el-icon-caidan3',routeName:'Base'},
          //   component:()=>import('../views/Mindex/Base.vue')
          // },
          {
            // 首页......合同管理.....采购合同......新建采购合同
            path:"AddPurchaseContract",
            name:"AddPurchaseContract",
            meta:{isShow:1,title:'新建采购合同',routeName:'PurchaseContract'},
            component:resolve => require(['@/views/Mindex/AddPurchaseContract.vue'],resolve),
          },
          {
            // 首页......合同管理.....采购合同......编辑采购合同
            path:"EditPurchaseContract",
            name:"EditPurchaseContract",
            meta:{isShow:1,title:'编辑采购合同',routeName:'PurchaseContract'},
            component:resolve => require(['@/views/Mindex/EditPurchaseContract.vue'],resolve)
          },
          {
            // 首页......合同管理
            path:'Contract',
            name:'Contract',
            meta:{isShow:0,title:'合同管理',icon:'el-icon-alihetong'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            // component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                // 首页......合同管理.....签署合同
                path:"SignContract",
                name:"SignContract",
                meta:{isShow:0,title:'电子签章',routeName:'SignContract'},
                component:resolve => require(['@/views/Mindex/SignContract.vue'],resolve)
              },
              {
                // 首页......合同管理.....采购合同
                path:"PurchaseContract",
                name:"PurchaseContract",
                meta:{isShow:0,title:'采购合同',routeName:'PurchaseContract'},
                component:resolve => require(['@/views/Mindex/PurchaseContract.vue'],resolve)
              },
              {
                //首页......合同管理.....销售合同
                path:"SalesContract",
                name:"SalesContract",
                meta:{isShow:0,title:'销售合同',routeName:'SalesContract'},
                component:resolve => require(['@/views/Mindex/SalesContract.vue'],resolve)
              },
              // {
              //   // 首页......合同管理.....上传文件
              //   path:"FileUpload",
              //   name:"FileUpload",
              //   meta:{isShow:0,title:'上传文件',routeName:'FileUpload'},
              //   component:resolve => require(['@/views/Mindex/FileUpload.vue'],resolve)
              // },
            ]
          },
          {
            //首页......合同管理.....签署合同......发起签署合同
            path:"AddSignContract",
            name:"AddSignContract",
            meta:{isShow:1,title:'发起电子签署',routeName:'SignContract'},
            component:resolve => require(['@/views/Mindex/AddSignContract.vue'],resolve)
          },
          {
            //首页......合同管理.....销售合同......新建销售合同
            path:"AddSalesContract",
            name:"AddSalesContract",
            meta:{isShow:1,title:'新建销售合同',routeName:'SalesContract'},
            component:resolve => require(['@/views/Mindex/AddSalesContract.vue'],resolve),
            // component:()=> import('../views/Mindex/AddSalesContract.vue')
          },
          {
            //首页......合同管理.....销售合同......编辑销售合同
            path:"EditSalesContract",
            name:"EditSalesContract",
            meta:{isShow:1,title:'编辑销售合同',routeName:'SalesContract'},
            component:resolve => require(['@/views/Mindex/EditSalesContract.vue'],resolve),
            // component:()=> import('../views/Mindex/EditSalesContract.vue')
          },
          // {
          //   // 首页......资金监管
          //   path:'BillManagement',
          //   name:'BillManagement',
          //   meta:{isShow:0,title:'资金监管',icon:'el-icon-aliVector'},
          //   component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
          //   children:[
          //     {
          //       // 首页......资金监管.....付款记录
          //       path:"PaymentRecord",
          //       name:"PaymentRecord",
          //       meta:{isShow:0,title:'付款记录',routeName:'PaymentRecord'},
          //       component:resolve => require(['@/views/Mindex/PaymentRecord.vue'],resolve)
          //     },
          //     {
          //       // 首页......资金监管.....收款记录
          //       path:"MakeCollections",
          //       name:"MakeCollections",
          //       meta:{isShow:0,title:'收款记录',routeName:'MakeCollections'},
          //       component:resolve => require(['@/views/Mindex/MakeCollections.vue'],resolve)
          //     },
          //     {
          //       // 首页......资金监管.....对账记录
          //       path:"Reconciliation",
          //       name:"Reconciliation",
          //       meta:{isShow:0,title:'对账记录',routeName:'Reconciliation'},
          //       component:resolve => require(['@/views/Mindex/Reconciliation.vue'],resolve)
          //     },
          //   ]
          // },
          {
            //首页......物流管理
            path:'Logistics',
            name:'Logistics',
            meta:{isShow:0,title:'物流管理',icon:'el-icon-aliwuliu'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            // component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......物流管理.....订单管理
                path:"Waybill",
                name:"Waybill",
                meta:{isShow:0,title:'订单管理',routeName:'Waybill'},
                component:resolve => require(['@/views/Mindex/Waybill.vue'],resolve),
              },
              {
                //首页......物流管理.....运单管理
                path:"WaybillManag",
                name:"WaybillManag",
                meta:{isShow:0,title:'运单管理',routeName:'WaybillManag'},
                component:resolve => require(['@/views/Mindex/WaybillManag.vue'],resolve),
              },
              {
                //首页......物流管理.....磅单管理
                path:"Poundrder",
                name:"Poundrder",
                meta:{isShow:0,title:'磅单管理',routeName:'Poundrder'},
                component:resolve => require(['@/views/Mindex/Poundrder.vue'],resolve),
              },
              
              {
                //首页......物流管理.....手动刷新
                path:"liveData",
                name:"liveData",
                meta:{isShow:0,title:'手动刷新',routeName:'liveData'},
                component:resolve => require(['@/views/Mindex/liveData.vue'],resolve),
              },
            ]
          },

          {
            //首页......物流管理.....磅单管理......新建磅单管理
            path:"AddPoundrder",
            name:"AddPoundrder",
            meta:{isShow:1,title:'新建磅单管理',routeName:'Poundrder'},
            component:resolve => require(['@/views/Mindex/AddPoundrder.vue'],resolve),
          },
          {
            //首页......物流管理.....磅单管理......编辑磅单管理
            path:"EditPoundrder",
            name:"EditPoundrder",
            meta:{isShow:1,title:'编辑磅单管理',routeName:'Poundrder'},
            component:resolve => require(['@/views/Mindex/EditPoundrder.vue'],resolve),
          },

          {
            //首页......物流管理.....订单管理......新建订单信息...关联运单信息
            path:"RelationBill",
            name:"RelationBill",
            meta:{isShow:1,title:'新建关联运单',routeName:'Waybill'},
            component:resolve => require(['@/views/Mindex/RelationBill.vue'],resolve),
          },
          {
            //首页......物流管理.....订单管理......新建订单信息
            path:"AddWaybill",
            name:"AddWaybill",
            meta:{isShow:1,title:'新建订单信息',routeName:'Waybill'},
            component:resolve => require(['@/views/Mindex/AddWaybill.vue'],resolve),
          },
          {
            //首页......物流管理.....订单管理......运单详情
            path:"WaybillDetails",
            name:"WaybillDetails",
            meta:{isShow:1,title:'运单详情',routeName:'Waybill'},
            component:resolve => require(['@/views/Mindex/WaybillDetails.vue'],resolve),
          },
          {
            //首页......物流管理.....订单管理......编辑订单信息
            path:"EditWaybill",
            name:"EditWaybill",
            meta:{isShow:1,title:'编辑订单信息',routeName:'Waybill'},
            component:resolve => require(['@/views/Mindex/EditWaybill.vue'],resolve)
          },
          {
            //首页......物流管理.....订单管理......编辑订单信息...关联运单信息
            path:"EditRelationBill",
            name:"EditRelationBill",
            meta:{isShow:1,title:'编辑关联运单',routeName:'Waybill'},
            component:resolve => require(['@/views/Mindex/EditRelationBill.vue'],resolve)
          },
          {
            //首页......结算管理.....采购结算......新建采购结算
            path:"AddPurchaseBill",
            name:"AddPurchaseBill",
            meta:{isShow:1,title:'新建采购结算',routeName:'PurchaseBill'},
            component:resolve => require(['@/views/Mindex/AddPurchaseBill.vue'],resolve),
            // component:()=> import('../views/Mindex/AddPurchaseBill.vue')
          },
          {
            //首页......结算管理.....采购结算......编辑采购结算
            path:"EditPurchaseBill",
            name:"EditPurchaseBill",
            meta:{isShow:1,title:'编辑采购结算',routeName:'PurchaseBill'},
            component:resolve => require(['@/views/Mindex/EditPurchaseBill.vue'],resolve)
          },
          {
            //首页......结算管理
            path:'Settlement',
            name:'Settlement',
            meta:{isShow:0,title:'结算管理',icon:'el-icon-aliicon-settlement'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            children:[
              {
                //首页......结算管理.....采购结算
                path:"PurchaseBill",
                name:"PurchaseBill",
                meta:{isShow:0,title:'采购结算',routeName:'PurchaseBill'},
                component:resolve => require(['@/views/Mindex/PurchaseBill.vue'],resolve)
              },
              {
                //首页......结算管理.....销售结算
                path:"SalesBill",
                name:"SalesBill",
                meta:{isShow:0,title:'销售结算',routeName:'SalesBill'},
                component:resolve => require(['@/views/Mindex/SalesBill.vue'],resolve)
              },
              {
                //首页......结算管理.....关联操作
                path:"Operation",
                name:"Operation",
                meta:{isShow:0,title:'关联操作',routeName:'Operation'},
                component:resolve => require(['@/views/Mindex/Operation.vue'],resolve)
              },
              {
                //首页......结算管理.....关联明细
                path:"Detailsment",
                name:"Detailsment",
                meta:{isShow:0,title:'关联明细',routeName:'Detailsment'},
                component:resolve => require(['@/views/Mindex/Detailsment.vue'],resolve)
              }
            ]
          },
          {
            //首页......结算管理.....销售结算......新建销售结算
            path:"AddSalesBill",
            name:"AddSalesBill",
            meta:{isShow:1,title:'新建销售结算',routeName:'SalesBill'},
            component:resolve => require(['@/views/Mindex/AddSalesBill.vue'],resolve),
            // component:()=> import('../views/Mindex/AddSalesBill.vue')
          },
          {
            //首页......结算管理.....销售结算......编辑销售结算
            path:"EditSalesBill",
            name:"EditSalesBill",
            meta:{isShow:1,title:'编辑销售结算',routeName:'SalesBill'},
            component:resolve => require(['@/views/Mindex/EditSalesBill.vue'],resolve),
            // component:()=> import('../views/Mindex/EditSalesBill.vue')
          },
          {
            //首页......票据管理.....进项发票......新建进项发票
            path:"AddInBillStatistics",
            name:"AddInBillStatistics",
            meta:{isShow:1,title:'新建进项发票',routeName:'InBillStatistics'},
            component:resolve => require(['@/views/Mindex/AddInBillStatistics.vue'],resolve),
            // component:()=> import('../views/Mindex/AddInBillStatistics.vue')
          },
          {
            //首页......票据管理.....进项发票......编辑进项发票
            path:"EditInBillStatistics",
            name:"EditInBillStatistics",
            meta:{isShow:1,title:'编辑进项发票',routeName:'InBillStatistics'},
            component:resolve => require(['@/views/Mindex/EditInBillStatistics.vue'],resolve),
            // component:()=> import('../views/Mindex/EditInBillStatistics.vue')
          },
          {
            //首页......票据管理
            path:'Bill',
            name:'Bill',
            meta:{isShow:0,title:'票据管理',icon:'el-icon-alipiaoju'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            // component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......票据管理.....进项发票
                path:"InBillStatistics",
                name:"InBillStatistics",
                meta:{isShow:0,title:'进项发票',routeName:'InBillStatistics'},
                component:resolve => require(['@/views/Mindex/InBillStatistics.vue'],resolve),
                // component:()=> import('../views/Mindex/InBillStatistics.vue')
              },
              {
                //首页......票据管理.....销项发票
                path:"SalesStatistics",
                name:"SalesStatistics",
                meta:{isShow:0,title:'销项发票',routeName:'SalesStatistics'},
                component:resolve => require(['@/views/Mindex/SalesStatistics.vue'],resolve),
                // component:()=> import('../views/Mindex/SalesStatistics.vue')
              }
            ]
          },
          {
            //首页......票据管理.....销项发票......新建销项发票
            path:"AddSalesStatistics",
            name:"AddSalesStatistics",
            meta:{isShow:1,title:'新建销项发票',routeName:'SalesStatistics'},
            component:resolve => require(['@/views/Mindex/AddSalesStatistics.vue'],resolve),
            // component:()=> import('../views/Mindex/AddSalesStatistics.vue')
          },
          {
            //首页......票据管理.....销项发票......编辑销项发票
            path:"EditSalesStatistics",
            name:"EditSalesStatistics",
            meta:{isShow:1,title:'编辑销项发票',routeName:'SalesStatistics'},
            component:resolve => require(['@/views/Mindex/EditSalesStatistics.vue'],resolve),
            // component:()=> import('../views/Mindex/EditSalesStatistics.vue')
          }, 
          {
            //首页......收支管理.....销售收款单......新建销售收款单
            path:"AddReceipt",
            name:"AddReceipt",
            meta:{isShow:1,title:'新建收款单',routeName:'Receipt'},
            component:resolve => require(['@/views/Mindex/AddReceipt.vue'],resolve),
            // component:()=> import('../views/Mindex/AddReceipt.vue')
          },
          {
            //首页......收付管理.....销售收款单......编辑销售收款单
            path:"EditReceipt",
            name:"EditReceipt",
            meta:{isShow:1,title:'编辑收款单',routeName:'Receipt'},
            component:resolve => require(['@/views/Mindex/EditReceipt.vue'],resolve),
            // component:()=> import('../views/Mindex/EditReceipt.vue')
          },
          {
            //首页......收付管理
            path:'Revenue',
            name:'Revenue',
            meta:{isShow:0,title:'收付管理',icon:'el-icon-aliicon-Receiptandpayment'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            // component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......收付管理.....销售收款单
                path:"Receipt",
                name:"Receipt",
                meta:{isShow:0,title:'收款单',routeName:'Receipt'},
                component:resolve => require(['@/views/Mindex/Receipt.vue'],resolve),
                // component:()=> import('../views/Mindex/Receipt.vue')
              },
              {
                //首页......收付管理.....采购付款单
                path:"PaymentSlip",
                name:"PaymentSlip",
                meta:{isShow:0,title:'付款单',routeName:'PaymentSlip'},
                component:resolve => require(['@/views/Mindex/PaymentSlip.vue'],resolve),
                // component:()=> import('../views/Mindex/PaymentSlip.vue')
              }
            ]
          },
          {
            //首页......收付管理.....采购付款单......新建采购付款单
            path:"AddPaymentSlip",
            name:"AddPaymentSlip",
            meta:{isShow:1,title:'新建付款单',routeName:'PaymentSlip'},
            component:resolve => require(['@/views/Mindex/AddPaymentSlip.vue'],resolve),
            // component:()=> import('../views/Mindex/AddPaymentSlip.vue')
          },
          {
            //首页......收付管理.....采购付款单......编辑采购付款单
            path:"EditPaymentSlip",
            name:"EditPaymentSlip",
            meta:{isShow:1,title:'编辑付款单',routeName:'PaymentSlip'},
            component:resolve => require(['@/views/Mindex/EditPaymentSlip.vue'],resolve),
            // component:()=> import('../views/Mindex/EditPaymentSlip.vue')
          },
          {
            //首页......往来管理
            path:'Contacts',
            name:'Contacts',
            meta:{isShow:0,title:'往来管理',icon:'el-icon-alihexiao'},
            component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......往来管理.....应收单
                path:"Receivable",
                name:"Receivable",
                meta:{isShow:0,title:'应收单',routeName:'Receivable'},
                component:()=> import('../views/Mindex/Receivable.vue')
              },
              {
                //首页......往来管理.....应付单
                path:"Copewith",
                name:"Copewith",
                meta:{isShow:0,title:'应付单',routeName:'Copewith'},
                component:()=> import('../views/Mindex/Copewith.vue')
              },
              {
                //首页......往来管理.....明细
                path:"Contactsdetailed",
                name:"Contactsdetailed",
                meta:{isShow:0,title:'明细',routeName:'Contactsdetailed'},
                component:()=> import('../views/Mindex/Contactsdetailed.vue')
              },
            ]
          },
          {
            //首页......合作方管理
            path:'XTner',
            name:'XTner',
            meta:{isShow:0,title:'合作方管理',icon:'el-icon-alihezuofang'},
            component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......合作方管理.....合作方维护
                path:"Partner",
                name:"Partner",
                meta:{isShow:0,title:'合作方维护',routeName:'Partner'},
                component:resolve => require(['@/views/Mindex/Partner.vue'],resolve),
              },
              {
                //首页......合作方管理.....地址维护
                path:"Address",
                name:"Address",
                meta:{isShow:0,title:'地址维护',routeName:'Address'},
                component:resolve => require(['@/views/Mindex/Address.vue'],resolve),
              },
              {
                //首页......合作方管理.....账户维护
                path:"Account",
                name:"Account",
                meta:{isShow:0,title:'账户维护',routeName:'Account'}, 
                component:resolve => require(['@/views/Mindex/Account.vue'],resolve),
              },
            ]
          },
          {
            //首页......合作方维护.....供应商档案......新建供应商档案
            path:"Addsupplier",
            name:"Addsupplier",
            meta:{isShow:1,title:'新建供应商档案',routeName:'Partner'},
            component:resolve => require(['@/views/Mindex/Addsupplier.vue'],resolve),
            // component:()=> import('../views/Mindex/Addsupplier.vue')
          },
          {
            //首页......合作方维护.....供应商档案......编辑供应商档案
            path:"Editsupplier",
            name:"Editsupplier",
            meta:{isShow:1,title:'编辑供应商档案',routeName:'Partner'},
            component:resolve => require(['@/views/Mindex/Editsupplier.vue'],resolve),
            // component:()=> import('../views/Mindex/Editsupplier.vue')
          },
          {
            //首页......合作方维护.....客户档案......新建客户档案
            path:"Addcustomer",
            name:"Addcustomer",
            meta:{isShow:1,title:'新建客户档案',routeName:'Partner'},
            component:resolve => require(['@/views/Mindex/Addcustomer.vue'],resolve),
            // component:()=> import('../views/Mindex/Addcustomer.vue')
          },
          {
            //首页......合作方维护.....客户档案......编辑客户档案
            path:"Editcustomer",
            name:"Editcustomer",
            meta:{isShow:1,title:'编辑客户档案',routeName:'Partner'},
            component:resolve => require(['@/views/Mindex/Editcustomer.vue'],resolve),
            // component:()=> import('../views/Mindex/Editcustomer.vue')
          },
          {
            //首页......地址维护.....新建地址维护
            path:"Addhair",
            name:"Addhair",
            meta:{isShow:1,title:'新建地址维护',routeName:'Address'},
            component:resolve => require(['@/views/Mindex/Addhair.vue'],resolve),
            // component:()=> import('../views/Mindex/Addhair.vue')
          },
          {
            //首页......地址维护.....发货地......编辑发货地
            path:"Edithair",
            name:"Edithair",
            meta:{isShow:1,title:'编辑发货地',routeName:'Address'},
            component:resolve => require(['@/views/Mindex/Edithair.vue'],resolve),
            // component:()=> import('../views/Mindex/Edithair.vue')
          },
          // {
          //   //首页......地址维护.....收货地......新建收货地
          //   path:"Addcollect",
          //   name:"Addcollect",
          //   meta:{isShow:1,title:'新建收货地',routeName:'Address'},
          //   component:()=> import('../views/Mindex/Addcollect.vue')
          // },
          {
            //首页......地址维护.....收货地......编辑收货地
            path:"Editcollect",
            name:"Editcollect",
            meta:{isShow:1,title:'编辑收货地',routeName:'Address'},
            component:resolve => require(['@/views/Mindex/Editcollect.vue'],resolve),
            // component:()=> import('../views/Mindex/Editcollect.vue')
          },
          // {
          //   //首页......地址维护.....装货地点......新建装货地点
          //   path:"Addpretend",
          //   name:"Addpretend",
          //   meta:{isShow:1,title:'新建装货地点',routeName:'Address'},
          //   component:()=> import('../views/Mindex/Addpretend.vue')
          // },
          {
            //首页......地址维护.....装货地点......编辑装货地点
            path:"Editpretend",
            name:"Editpretend",
            meta:{isShow:1,title:'编辑装货地点',routeName:'Address'},
            component:resolve => require(['@/views/Mindex/Editpretend.vue'],resolve),
            // component:()=> import('../views/Mindex/Editpretend.vue')
          },
          // {
          //   //首页......地址维护.....卸货地点......新建卸货地点
          //   path:"Adddischarge",
          //   name:"Adddischarge",
          //   meta:{isShow:1,title:'新建卸货地点',routeName:'Address'},
          //   component:()=> import('../views/Mindex/Adddischarge.vue')
          // },
          {
            //首页......地址维护.....卸货地点......编辑卸货地点
            path:"Editdischarge",
            name:"Editdischarge",
            meta:{isShow:1,title:'编辑卸货地点',routeName:'Address'},
            component:resolve => require(['@/views/Mindex/Editdischarge.vue'],resolve),
            // component:()=> import('../views/Mindex/Editdischarge.vue')
          },
          {
            //首页......账户维护.....新建账户维护
            path:"AddAccount",
            name:"AddAccount",
            meta:{isShow:1,title:'新建账户维护',routeName:'Account'},
            component:resolve => require(['@/views/Mindex/AddAccount.vue'],resolve),
            // component:()=> import('../views/Mindex/AddAccount.vue')
          },
          {
            //首页......账户维护.....编辑账户维护
            path:"EditAccount",
            name:"EditAccount",
            meta:{isShow:1,title:'编辑账户维护',routeName:'Account'},
            component:resolve => require(['@/views/Mindex/EditAccount.vue'],resolve),
            // component:()=> import('../views/Mindex/EditAccount.vue')
          },
          // {
          //   //首页......统计管理
          //   path:'XStatistics',
          //   name:'XStatistics',
          //   meta:{isShow:2,title:'统计管理'},
          // },
          {
            //首页......业务统计
            path:'Accounting',
            name:'Accounting',
            meta:{isShow:0,title:'业务统计',icon:'el-icon-alishuju1'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            // component:()=>import('../views/Mindex/MindexView.vue'),
            children:[
              {
                //首页......业务统计.....缴税明细
                path:"Coststatistics",
                name:"Coststatistics",
                meta:{isShow:0,title:'缴税明细',routeName:'Coststatistics'},
                component:resolve => require(['@/views/Mindex/Coststatistics.vue'],resolve)
              },
              // {
              //   //首页......业务核算.....收入统计
              //   path:"IncomeStatistics",
              //   name:"IncomeStatistics",
              //   meta:{isShow:0,title:'收入统计',routeName:'IncomeStatistics'},
              //   component:resolve => require(['@/views/Mindex/IncomeStatistics.vue'],resolve),
              //   // component:()=> import('../views/Mindex/IncomeStatistics.vue')
              // },
              // {
              //   //首页......业务核算.....收益核算
              //   path:"IncomeAccounting",
              //   name:"IncomeAccounting",
              //   meta:{isShow:0,title:'收益核算',routeName:'IncomeAccounting'},
              //   component:resolve => require(['@/views/Mindex/IncomeAccounting.vue'],resolve),
              //   // component:()=> import('../views/Mindex/IncomeAccounting.vue')
              // },
            ]
          },
          {
            //首页......业务统计.....缴税明细......新建缴税明细
            path:"AddCoststatistics",
            name:"AddCoststatistics",
            meta:{isShow:1,title:'新建缴税明细',routeName:'Coststatistics'},
            component:resolve => require(['@/views/Mindex/AddCoststatistics.vue'],resolve),
          },
          {
            //首页......业务统计.....缴税明细......编辑缴税明细
            path:"EditCoststatistics",
            name:"EditCoststatistics",
            meta:{isShow:1,title:'编辑缴税明细',routeName:'Coststatistics'},
            component:resolve => require(['@/views/Mindex/EditCoststatistics.vue'],resolve),
          },
        ]
      },
      {
        //设置
        path:'Setup',
        name:'Setup',
        meta:{isShow:0,title:'设 置',icon:'el-icon-alishezhi1'},
        redirect:{name:''},
        component:resolve => require(['@/views/CommonView.vue'],resolve),
        children:[
          // {
          //   //设置......应用权限
          //   path:'Application',
          //   name:'Application',
          //   meta:{isShow:0,title:'应用权限',icon:'el-icon-aliapplication',routeName:'Application'},
          //   component:resolve => require(['@/views/setup/Application.vue'],resolve)
          // },
          {
            //设置......用户列表
            path:'UserList',
            name:'UserList',
            meta:{isShow:0,title:'用户列表',icon:'el-icon-aliyonghuliebiao',routeName:'UserList'},
            component:resolve => require(['@/views/setup/UserList.vue'],resolve)
          },
          {
            //设置......角色权限
            path:'Roleperm',
            name:'Roleperm',
            meta:{isShow:0,title:'角色权限',icon:'el-icon-alifenghuangxiangmutubiao_quanxian',routeName:'Roleperm'},
            component:resolve => require(['@/views/setup/Roleperm.vue'],resolve)
          },
          {
            //设置......企业档案
            path:'EnterpriseArchives',
            name:'EnterpriseArchives',
            meta:{isShow:0,title:'企业档案',icon:'el-icon-alidangan',routeName:'EnterpriseArchives'},
            component:resolve => require(['@/views/setup/EnterpriseArchives.vue'],resolve)
          },
          // {
          //   //设置......合同模板
          //   path:'ContractTemplate',
          //   name:'ContractTemplate',
          //   meta:{isShow:0,title:'合同模板',icon:'el-icon-alifabuweimoban',routeName:'ContractTemplate'},
          //   component:resolve => require(['@/views/setup/ContractTemplate.vue'],resolve)
          // },
          {
            //设置......邀请码生成
            path:'InvitationCode',
            name:'InvitationCode',
            meta:{isShow:0,title:'邀请码生成',icon:'el-icon-aliyaoqing',routeName:'InvitationCode'},
            component:resolve => require(['@/views/setup/InvitationCode.vue'],resolve)
          },
          {
            //首页......日志管理
            path:'JournalManagement',
            name:'JournalManagement',
            meta:{isShow:0,title:'日志管理',icon:'el-icon-alifabuweimoban'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            children:[
              {
                //首页......日志管理.....日志明细
                path:"Journal",
                name:"Journal",
                meta:{isShow:0,title:'日志明细',routeName:'Journal'},
                component:resolve => require(['@/views/Mindex/Journal.vue'],resolve)
              }
            ]
          },
        ]
      },
    ]
  }
]
// 平台动态路由表
export const platformRouterMap = [
  {
    // Main主页面
    path:'/Main',
    name:'Main',
    meta:{isShow:0,title:'主页'},
    redirect:{name:''},
    component: ()=> import('../views/Main.vue'),
    children:[
      {
        //首页
        path:'Mindex',
        name:'Mindex',
        meta:{isShow:0,title:'首 页',icon:'el-icon-alishouye'},
        redirect:{name:''},
        component:()=>import('../views/AsidsMain.vue'),
        children:[
          {
            //首页......仪表盘
            path:'PtIndex',
            name:'PtIndex',
            meta:{isShow:0,title:'仪表盘',icon:'el-icon-alishouye',routeName:'PtIndex'},
            component:resolve => require(['@/views/Ptmain/PtIndex.vue'],resolve)
          }, 
          // {
          //   //首页......测试
          //   path:'Demo',
          //   name:'Demo',
          //   meta:{isShow:0,title:'测试',icon:'el-icon-alishouye'},
          //   component:resolve => require(['../views/Demo.vue'],resolve)
          // }, 
          // {
          //   //首页......项目管理
          //   path:'Menter',
          //   name:'Menter',
          //   meta:{isShow:2,title:'项目管理'},
          // },
          {
            //首页......供应链业务
            path:'ChainBus',
            name:'ChainBus',
            meta:{isShow:0,title:'供应链业务',icon:'el-icon-aliweibiaoti-4',routeName:'ChainBus'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            children:[
              {
                // 首页......供应链业务.....企业信息
                path:"Qycensus",
                name:"Qycensus",
                meta:{isShow:0,title:'企业信息',routeName:'Qycensus'},
                component:resolve => require(['@/views/Ptmain/Qycensus.vue'],resolve)
              },
              {
                // 首页......供应链业务.....业务核算
                path:"BusinessAdjust",
                name:"BusinessAdjust",
                meta:{isShow:0,title:'业务核算',routeName:'BusinessAdjust'},
                component:resolve => require(['@/views/Ptmain/BusinessAdjust.vue'],resolve)
              },
              // {
              //   // 首页......供应链业务.....物流统计
              //   path:"LogStat",
              //   name:"LogStat",
              //   meta:{isShow:0,title:'物流统计',routeName:'LogStat'},
              //   component:resolve => require(['@/views/Ptmain/LogStat.vue'],resolve)
              // },
            ]
          },
          {
            // 首页......供应链业务.....企业统计......编辑企业统计
            path:"EditQycensus",
            name:"EditQycensus",
            meta:{isShow:1,title:'编辑企业统计',routeName:'Qycensus'},
            component:resolve => require(['@/views/Ptmain/EditQycensus.vue'],resolve)
          },
          // {
          //   // 首页......贸易服务TradeServices
          //   path:'TradeServices',
          //   name:'TradeServices',
          //   meta:{isShow:0,title:'贸易服务',icon:'el-icon-aliVector',routeName:'TradeServices'},
          //   component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
          //   children:[
          //     {
          //       // 首页......贸易服务.....交易记录
          //       path:"TradeContract",
          //       name:"TradeContract",
          //       meta:{isShow:0,title:'交易记录',routeName:'TradeContract'},
          //       component:resolve => require(['@/views/Ptmain/TradeContract.vue'],resolve)
          //     },
          //     // {
          //     //   // 首页......贸易服务.....解冻信息
          //     //   path:"Thaw",
          //     //   name:"Thaw",
          //     //   meta:{isShow:0,title:'解冻信息',routeName:'Thaw'},
          //     //   component:resolve => require(['@/views/Ptmain/Thaw.vue'],resolve)
          //     // },
          //   ]
          // },
          {   
            // 首页......熊猫煤仓
            path:'Palmcoal',
            name:'Palmcoal',
            meta:{isShow:0,title:'熊猫煤仓',icon:'el-icon-alimeikuangqiye',routeName:'Palmcoal'},
            component:resolve => require(['@/views/Mindex/MindexView.vue'],resolve),
            children:[
              {
                // 首页......熊猫煤仓.....供货信息
                path:"varietyList",
                name:"varietyList",
                meta:{isShow:0,title:'供货信息',routeName:'varietyList'},
                component:resolve => require(['@/views/Ptmain/varietyList.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....采购信息
                path:"procureList",
                name:"procureList",
                meta:{isShow:0,title:'采购信息',routeName:'procureList'},
                component:resolve => require(['@/views/Ptmain/procureList.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....云仓管理
                path:"Cloudcloud",
                name:"Cloudcloud",
                meta:{isShow:0,title:'云仓管理',routeName:'Cloudcloud'},
                component:resolve => require(['@/views/Ptmain/Cloudcloud.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....用户列表
                path:"userasp",
                name:"userasp",
                meta:{isShow:0,title:'用户列表',routeName:'userasp'},
                component:resolve => require(['@/views/Ptmain/userasp.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....意见建议
                path:"advocate",
                name:"advocate",
                meta:{isShow:0,title:'意见建议',routeName:'advocate'},
                component:resolve => require(['@/views/Ptmain/advocate.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....标签列表
                path:"labeltag",
                name:"labeltag",
                meta:{isShow:0,title:'标签列表',routeName:'labeltag'},
                component:resolve => require(['@/views/Ptmain/labeltag.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....煤种操作
                path:"coaltoper",
                name:"coaltoper",
                meta:{isShow:0,title:'煤种操作',routeName:'coaltoper'},
                component:resolve => require(['@/views/Ptmain/coaltoper.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....红包记录
                path:"redenvecord",
                name:"redenvecord",
                meta:{isShow:0,title:'红包记录',routeName:'redenvecord'},
                component:resolve => require(['@/views/Ptmain/redenvecord.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....订单记录
                path:"ordergoods",
                name:"ordergoods",
                meta:{isShow:0,title:'订单记录',routeName:'ordergoods'},
                component:resolve => require(['@/views/Ptmain/ordergoods.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....租赁售卖
                path:"rentalSales",
                name:"rentalSales",
                meta:{isShow:0,title:'租赁售卖',routeName:'rentalSales'},
                component:resolve => require(['@/views/Ptmain/rentalSales.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....会员管理
                path:"vipment",
                name:"vipment",
                meta:{isShow:0,title:'会员管理',routeName:'vipment'},
                component:resolve => require(['@/views/Ptmain/vipment.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....访问数据
                path:"visit",
                name:"visit",
                meta:{isShow:0,title:'访问数据',routeName:'visit'},
                component:resolve => require(['@/views/Ptmain/visit.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....兑换明细
                path:"chdetails",
                name:"chdetails",
                meta:{isShow:0,title:'兑换明细',routeName:'chdetails'},
                component:resolve => require(['@/views/Ptmain/chdetails.vue'],resolve)
              },
              {
                // 首页......熊猫煤仓.....商品列表
                path:"Produlist",
                name:"Produlist",
                meta:{isShow:0,title:'商品列表',routeName:'Produlist'},
                component:resolve => require(['@/views/Ptmain/Produlist.vue'],resolve)
              }
            ]
          },
        ]
      },
      {
        // 设置
        path:'Setup',
        name:'Setup',
        meta:{isShow:0,title:'设 置',icon:'el-icon-alishezhi1'},
        redirect:{name:''},
        component:resolve => require(['@/views/CommonView.vue'],resolve),
        children:[ 
          {
            //设置......平台用户列表
            path:'PtUserList',
            name:'PtUserList',
            meta:{isShow:0,title:'用户列表',icon:'el-icon-aliyonghuliebiao',routeName:'PtUserList'},
            component:resolve => require(['@/views/setup/PtUserList.vue'],resolve)
          },
          {
            //设置......菜单管理
            path:'menuPower',
            name:'menuPower',
            meta:{isShow:0,title:'菜单管理',icon:'el-icon-alifenghuangxiangmutubiao_quanxian',routeName:'menuPower'},
            component:resolve => require(['@/views/setup/menuPower.vue'],resolve)
          },
          {
            //设置......签章配置
            path:'ContractDeploy',
            name:'ContractDeploy',
            meta:{isShow:0,title:'签章配置',icon:'el-icon-aliicon-supplychain',routeName:'ContractDeploy'},
            component:resolve => require(['@/views/setup/ContractDeploy.vue'],resolve)
          },
          {
            //首页......合同配置.....新建合同配置
            path:"AddContractDeploy",
            name:"AddContractDeploy",
            meta:{isShow:1,title:'新建合同配置',routeName:'ContractDeploy'},
            component:resolve => require(['@/views/setup/AddContractDeploy.vue'],resolve)
          },
          {
            //首页......合同配置.....编辑合同配置
            path:"EditContractDeploy",
            name:"EditContractDeploy",
            meta:{isShow:1,title:'编辑合同配置',routeName:'ContractDeploy'},
            component:resolve => require(['@/views/setup/EditContractDeploy.vue'],resolve)
          },
          {
            //设置......平台角色权限
            path:'PtRoleperm',
            name:'PtRoleperm',
            meta:{isShow:0,title:'平台角色权限',icon:'el-icon-alifenghuangxiangmutubiao_quanxian',routeName:'PtRoleperm'},
            component:resolve => require(['@/views/setup/PtRoleperm.vue'],resolve)
          },
          {
            //设置......企业角色权限
            path:'QyRoleperm',
            name:'QyRoleperm',
            meta:{isShow:0,title:'企业角色权限',icon:'el-icon-alifenghuangxiangmutubiao_quanxian',routeName:'QyRoleperm'},
            component:resolve => require(['@/views/setup/QyRoleperm.vue'],resolve)
          },
          {
            //设置......全局配置
            path:'globalConfig',
            name:'globalConfig',
            meta:{isShow:0,title:'全局配置',icon:'el-icon-alifenghuangxiangmutubiao_quanxian',routeName:'globalConfig'},
            component:resolve => require(['@/views/setup/globalConfig.vue'],resolve)
          }
        ]
      },
    ]
  }
]
// 政府动态路由表
export const governmentRouterMap = [
  {
    // 政府看板页面
    path:'/Xzmain',
    name:'Xzmain',
    component: ()=> import('../views/Xzmain.vue'),
  }
]
//错误页路由
export const afterRouteMap=[
  {
    path:'*',
    name:'NotFont',
    component:()=>import('../views/NotFont.vue')
  }
]
//初始化路由表
export const createRouter = () => new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes:staticRouterMap
});

// 原始路由表
const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes:staticRouterMap
})
router.beforeEach((to, from, next) => {
  var Authorization = sessionStorage.getItem('Authorization');
  // 注册申请   || to.name =='Register'
      if(Authorization || to.name =='Login'|| to.name == 'NotFont' || to.name =='Policy' || to.name =='Service' || to.name =='Forget' || to.name =='DPLogin' || to.name =='ELogin' || to.name == 'Sponsor' || to.name == 'Unknown' || to.name == 'Paymentyes' || to.name =='Register' || to.name =='wxindex' || to.name =='privacys' || to.name =='sergrent'){
        next();
      }else{
        Message({
          type:'warning',
          duration:1,
        })
        next({name:'Login'})
      }
})
export default router
