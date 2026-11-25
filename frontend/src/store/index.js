import Vue from 'vue'
import Vuex from 'vuex'
import {staticRouterMap,activeRouterMap,governmentRouterMap,platformRouterMap,createRouter,afterRouteMap} from '../router' // 依次为静态路由表，企业路由表，政府路由表，平台路由表，路由表初始化
import router from '../router'
import store from '../store'
import api from '../api' 
import path from '../api/path'
var newhttps = ''
if(path.baseUrl=='saas'){
  newhttps = '/' + path.baseUrl;
} else {
  newhttps = path.baseUrl;
}
// var colhttps = '/' + path.coludUrl;
Vue.use(Vuex)   
export default new Vuex.Store({
  state:{ 
    // /saas 打包环境
    // http://119.45.238.116:8766
    // https://test.coalcloudtech.com/saas
    http:newhttps,
    // /bus 打包环境
    // http://119.45.238.116:8767
    // https://test.coalcloudtech.com/saas
    cloudhttp:newhttps,
    routes: staticRouterMap,          // 完整路由表 this.$router.options.routes
    addRoutes: [],                    // 系统最新路由表（不分企业平台）
    Basics:{},                        // 基础维护数据
    menus:[],                         // 企业菜单列表
    Ptmenus:[],                       // 平台菜单列表
    newstate:null,
  },
  mutations:{ 
    // 基础维护数据
    Basics(state,addBasics){
      state.Basics = addBasics;
    },

    // 设置企业路由菜单
    setRoutes(state,addRoutes){
      // if(addRoutes){
      //   state.newstate = true;
      // }else{
      //   state.newstate = false;
      // }
      state.addRoutes = addRoutes;
      state.routes = staticRouterMap.concat(addRoutes);
    },
    // 调取企业菜单列表
    setMenus(state,addMenusState){
      state.menus = [];
      state.menus = addMenusState;
    },

    // 调取平台菜单列表
    PtsetMenus(state,addMenusState){
      state.Ptmenus = addMenusState;
    },
    // 设置平台路由菜单
    setplatformRoutes(state,addRoutes){
      state.addRoutes = addRoutes;
      state.routes = staticRouterMap.concat(addRoutes);
    },

    // 政府路由菜单
    setgoverRoutes(state,addRoutes){
      state.addRoutes = addRoutes;
      state.routes = staticRouterMap.concat(addRoutes);
    },

    // 修改404页面
    // NewStates(state,data){
        
    // },
  },
  actions:{ 
    // 基础维护数据
    AllBasics:(store)=>{
      // 全部基础维护数据
      api.basicslist.list() 
      .then(res=>{
        if(res.data.code == 200){
          var newdata = res.data.data;
          var obj = {};                 // 全部基础维护信息
          var price_type = {};          // 价格类型
          var payment_method = {};      // 付款方式
          var transport_mode = {};      // 运输模式
          var carrier = {};             // 承运工具
          var payment_type = {};        // 支付方式
          var pro_type = {};            // 产品大类
          var pro__middle_type = {};    // 产品小类
          var pro__small_type = {};     // 标的指标
          var logistics_platform = {};  // 物流平台
          var file_relation = {};       // 文件关联类型
          var purpose = {};             // 用途
          // console.log(newdata)
          for(let i=0;i<newdata.length;i++){
            if(newdata[i].dtype == "price_type"){
                price_type[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "payment_method"){
                payment_method[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "transport_mode"){
                transport_mode[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "carrier"){
                carrier[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "payment_type"){
                payment_type[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "pro_type"){
                pro_type[newdata[i].dvalue] = newdata[i].dkey
            }else if(newdata[i].dtype == "pro__middle_type"){
                pro__middle_type[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "pro__small_type"){
                pro__small_type[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "logistics_platform"){
                logistics_platform[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "file_relation"){
                file_relation[newdata[i].dvalue] = newdata[i].dkey;
            }else if(newdata[i].dtype == "purpose"){
                purpose[newdata[i].dvalue] = newdata[i].dkey;
            }
          }
          obj.price_type = price_type
          obj.payment_method = payment_method
          obj.transport_mode = transport_mode
          obj.carrier = carrier
          obj.payment_type = payment_type
          obj.pro_type = pro_type
          obj.pro__middle_type = pro__middle_type
          obj.pro__small_type = pro__small_type
          obj.logistics_platform = logistics_platform
          obj.file_relation = file_relation
          obj.purpose = purpose
          store.commit('Basics',obj);
        }
      })
    },

    // 调取企业菜单列表
    asynSetMenus:()=>{
      api.roleList.selMenus({searchType:3}).then((res)=>{
        if(res.data.code == 200){
          let newbackMenus = res.data.data.menus;
          let addRouteMenus = [{}];
          function getMenusData(routes){
            let childMenus = [];
            for(let i=0;i<routes.length;i++){
              for(let j=0;j<newbackMenus.length;j++){
                if(routes[i].name == newbackMenus[j].routeName){
                  newbackMenus[j].import = -1;
                  newbackMenus[j].export = -1;
                  newbackMenus[j].add = -1;
                  newbackMenus[j].check = -1;
                  newbackMenus[j].edit = -1;
                  newbackMenus[j].delete = -1;
                  newbackMenus[j].upload = -1;
                  newbackMenus[j].fileDelete = -1;
                  newbackMenus[j].button = [];
                  let newObj = res.data.data.menuAuthy[newbackMenus[j].id];
                  if(newObj !== undefined){
                    if(newObj.hasOwnProperty('import')){
                      newbackMenus[j].import = 1;
                      // newbackMenus[j].button.push('导入')
                    }else{
                      newbackMenus[j].import = -1;
                    }
                    if(newObj.hasOwnProperty('export')){
                      newbackMenus[j].export = 1;
                      // newbackMenus[j].button.push('导出')
                    }else{
                      newbackMenus[j].export = -1;
                    }
                    if(newObj.hasOwnProperty('add')){
                      newbackMenus[j].add = 1;
                      // newbackMenus[j].button.push('新增')
                    }else{
                      newbackMenus[j].add = -1;
                    }
                    if(newObj.hasOwnProperty('check')){
                      newbackMenus[j].check = 1;
                      // newbackMenus[j].button.push('查看')
                    }else{
                      newbackMenus[j].check = -1;
                    }
                    if(newObj.hasOwnProperty('edit')){
                      newbackMenus[j].edit = 1;
                      // newbackMenus[j].button.push('编辑')
                    }else{
                      newbackMenus[j].edit = -1;
                    }  
                    if(newObj.hasOwnProperty('delete')){
                      newbackMenus[j].delete = 1;
                      // newbackMenus[j].button.push('删除')
                    }else{
                      newbackMenus[j].delete = -1;
                    }
                    if(newObj.hasOwnProperty('upload')){
                      newbackMenus[j].upload = 1;
                      // newbackMenus[j].button.push('上传')
                    }else{
                      newbackMenus[j].upload = -1;
                    }
                    if(newObj.hasOwnProperty('fileDelete')){
                      newbackMenus[j].fileDelete = 1;
                      // newbackMenus[j].button.push('文件删除')
                    }else{
                      newbackMenus[j].fileDelete = -1;
                    }
                  }
                  if(routes[i].children){
                    newbackMenus[j].children=getMenusData(routes[i].children);
                  }
                  // newbackMenus[j] = Object.assign(newbackMenus[j],res.data.data.menuAuthy[newbackMenus[j].id]);
                  childMenus.push(newbackMenus[j]);
                }
              }
            }
            return childMenus;
          }
          addRouteMenus = getMenusData(activeRouterMap[0].children);
          store.commit('setMenus', addRouteMenus)
        }
      })
    },
    // 调取平台端--企业菜单列表
    asynPtQySetMenus:()=>{
      api.roleList.selMenus().then((res)=>{
        if(res.data.code == 200){
          let newbackMenus = res.data.data.menus;
          let addRouteMenus = [{}];
          function getMenusData(routes){
            let childMenus = [];
            for(let i=0;i<routes.length;i++){
              for(let j=0;j<newbackMenus.length;j++){
                if(routes[i].name == newbackMenus[j].routeName){
                  newbackMenus[j].import = -1;
                  newbackMenus[j].export = -1;
                  newbackMenus[j].add = -1;
                  newbackMenus[j].check = -1;
                  newbackMenus[j].edit = -1;
                  newbackMenus[j].delete = -1;
                  newbackMenus[j].upload = -1;
                  newbackMenus[j].fileDelete = -1;
                  newbackMenus[j].button = [];
                  let newObj = res.data.data.menuAuthy[newbackMenus[j].id];
                  if(newObj !== undefined){
                    if(newObj.hasOwnProperty('import')){
                      newbackMenus[j].import = 1;
                      // newbackMenus[j].button.push('导入')
                    }else{
                      newbackMenus[j].import = -1;
                    }
                    if(newObj.hasOwnProperty('export')){
                      newbackMenus[j].export = 1;
                      // newbackMenus[j].button.push('导出')
                    }else{
                      newbackMenus[j].export = -1;
                    }
                    if(newObj.hasOwnProperty('add')){
                      newbackMenus[j].add = 1;
                      // newbackMenus[j].button.push('新增')
                    }else{
                      newbackMenus[j].add = -1;
                    }
                    if(newObj.hasOwnProperty('check')){
                      newbackMenus[j].check = 1;
                      // newbackMenus[j].button.push('查看')
                    }else{
                      newbackMenus[j].check = -1;
                    }
                    if(newObj.hasOwnProperty('edit')){
                      newbackMenus[j].edit = 1;
                      // newbackMenus[j].button.push('编辑')
                    }else{
                      newbackMenus[j].edit = -1;
                    }
                    if(newObj.hasOwnProperty('delete')){
                      newbackMenus[j].delete = 1;
                      // newbackMenus[j].button.push('删除')
                    }else{
                      newbackMenus[j].delete = -1;
                    }
                    if(newObj.hasOwnProperty('upload')){
                      newbackMenus[j].upload = 1;
                      // newbackMenus[j].button.push('上传')
                    }else{
                      newbackMenus[j].upload = -1;
                    }
                    if(newObj.hasOwnProperty('fileDelete')){
                      newbackMenus[j].fileDelete = 1;
                      // newbackMenus[j].button.push('文件删除')
                    }else{
                      newbackMenus[j].fileDelete = -1;
                    }
                  }
                  if(routes[i].children){
                    newbackMenus[j].children=getMenusData(routes[i].children);
                  }
                  // newbackMenus[j] = Object.assign(newbackMenus[j],res.data.data.menuAuthy[newbackMenus[j].id]);
                  childMenus.push(newbackMenus[j]);
                }
              }
            }
            return childMenus;
          }
          addRouteMenus = getMenusData(activeRouterMap[0].children);
          store.commit('setMenus', addRouteMenus)
        }
      })
    },
    // 登录设置企业路由菜单
    asynSetRoutes:(store)=>{
      // 把请求的 promise 交回去:登录页要等动态路由注册完再跳转
      return api.authMenu.selUserMenus()
      .then(res=>{
        if(res.data.code == 200){
          let backMenus = res.data.data.menus;
          let addRoutes;
          // console.log(backMenus)
          if(backMenus.length > 0){
            store.state.newstate = true;
            // console.log(store.state.newstate)
            addRoutes = activeRouterMap;
            // addRoutes[0].children[0].redirect = {name:activeRouterMap[0].children[0].children[0].name};
            // addRoutes[0].redirect = {name:activeRouterMap[0].children[0].name};
            function getRouteData(routes){
              let childs = [];
              for(var i=0;i<routes.length;i++){
                for(var j=0;j<backMenus.length;j++){
                  if(routes[i].name == backMenus[j].routeName){
                    if(backMenus[j].add){
                      routes[i].meta.add = 1;
                    }
                    if(backMenus[j].edit){
                      routes[i].meta.edit = 1;
                    }
                    if(backMenus[j].check){
                      routes[i].meta.check = 1;
                    }
                    if(backMenus[j].delete){
                      routes[i].meta.delete = 1;
                    }
                    if(backMenus[j].export){
                      routes[i].meta.export = 1;
                    }
                    if(backMenus[j].import){
                      routes[i].meta.import = 1;
                    }
                    if(backMenus[j].upload){
                      routes[i].meta.upload = 1;
                    }
                    if(backMenus[j].fileDelete){
                      routes[i].meta.fileDelete = 1;
                    }
                    if(routes[i].children){
                      // console.log(routes[i].children[0].name)
                      // if(routes[i].name == 'Mindex'){
                      //   routes[i].redirect.name = routes[i].children[0].name
                      // }else if(routes[i].name == 'Setup'){
                      //   routes[i].redirect.name = routes[i].children[1].name
                      // }
                      // console.log(routes[i].children)
                      routes[i].children = getRouteData(routes[i].children)
                    }
                    childs.push(routes[i])
                  }
                }
              }
              return childs; 
            }
            addRoutes[0].children = getRouteData(activeRouterMap[0].children);
            let arr = res.data.data.routelist;
            // 只有设置菜单  [{name:'',chinname:''},{name:'Setup',chinname:'UserList'}]
            // 首页菜单都有  [{name:'Mindex',chinname:'MainPage'},{name:'',chinname:'UserList'}]
            // 只有首页   [{name:'Mindex',chinname:'MainPage'},{name:'',chinname:''}]
            for(let i=0;i<arr.length;i++){
              if(arr[i].name){
                addRoutes[0].redirect.name = arr[i].name;
              } 
              if(arr[i].chinname && addRoutes[0].children.length >= 2){
                addRoutes[0].children[i].redirect = {name:arr[i].chinname};
              }else if(arr[i].chinname && addRoutes[0].children.length < 2){
                addRoutes[0].children[0].redirect = {name:arr[i].chinname};
              }
            }
            store.commit('setRoutes',addRoutes);
            router.matcher = createRouter().matcher;
            router.addRoutes(addRoutes);
            router.addRoutes(afterRouteMap);
            router.options.routes = store.state.routes;
          }else{
            store.state.newstate = false;
          }
        }
      })
    },

    // 调取平台菜单列表
    asynPtSetMenus:()=>{
      api.roleList.selMenus({searchType:2}).then((res)=>{
        // console.log(res)
        if(res.data.code == 200){
          // console.log(res.data.data)
          let PtnewbackMenus = res.data.data.menus;
          let PtaddRouteMenus = [{}];
          function PtgetMenusData(routes){
            let PtchildMenus = [];
            for(let i=0;i<routes.length;i++){
              for(let j=0;j<PtnewbackMenus.length;j++){
                if(routes[i].name == PtnewbackMenus[j].routeName){
                  PtnewbackMenus[j].import = -1;
                  PtnewbackMenus[j].export = -1;
                  PtnewbackMenus[j].add = -1;
                  PtnewbackMenus[j].check = -1;
                  PtnewbackMenus[j].edit = -1;
                  PtnewbackMenus[j].delete = -1;
                  PtnewbackMenus[j].upload = -1;
                  PtnewbackMenus[j].fileDelete = -1;
                  PtnewbackMenus[j].button = [];
                  let newObj = res.data.data.menuAuthy[PtnewbackMenus[j].id];
                  if(newObj !== undefined){
                    if(newObj.hasOwnProperty('import')){
                      PtnewbackMenus[j].import = 1;
                    }else{
                      PtnewbackMenus[j].import = -1;
                    }
                    if(newObj.hasOwnProperty('export')){
                      PtnewbackMenus[j].export = 1;
                    }else{
                      PtnewbackMenus[j].export = -1;
                    }
                    if(newObj.hasOwnProperty('add')){
                      PtnewbackMenus[j].add = 1;
                    }else{
                      PtnewbackMenus[j].add = -1;
                    }
                    if(newObj.hasOwnProperty('check')){
                      PtnewbackMenus[j].check = 1;
                    }else{
                      PtnewbackMenus[j].check = -1;
                    }
                    if(newObj.hasOwnProperty('edit')){
                      PtnewbackMenus[j].edit = 1;
                    }else{
                      PtnewbackMenus[j].edit = -1;
                    }
                    if(newObj.hasOwnProperty('delete')){
                      PtnewbackMenus[j].delete = 1;
                    }else{
                      PtnewbackMenus[j].delete = -1;
                    }
                    if(newObj.hasOwnProperty('upload')){
                      PtnewbackMenus[j].upload = 1;
                    }else{
                      PtnewbackMenus[j].upload = -1;
                    }
                    if(newObj.hasOwnProperty('fileDelete')){
                      PtnewbackMenus[j].fileDelete = 1;
                    }else{
                      PtnewbackMenus[j].fileDelete = -1;
                    }
                    
                  }
                  if(routes[i].children){
                    PtnewbackMenus[j].children = PtgetMenusData(routes[i].children);
                  }
                  // PtnewbackMenus[j] = Object.assign(PtnewbackMenus[j],res.data.data.menuAuthy[PtnewbackMenus[j].id]);
                  PtchildMenus.push(PtnewbackMenus[j])
                }
              }
            }
            return PtchildMenus;
          }
          PtaddRouteMenus = PtgetMenusData(platformRouterMap[0].children);
          // console.log(PtaddRouteMenus)
          store.commit('PtsetMenus', PtaddRouteMenus)
        }
      })
    },
    // 登录设置平台路由菜单
    asynSetplatformRoutes:(store)=>{
      // 把请求的 promise 交回去:登录页要等动态路由注册完再跳转
      return api.authMenu.selUserMenus()
      .then(res=>{
        if(res.data.code == 200){
          let PtbackMenus = res.data.data.menus;
          if(PtbackMenus.length != 0){
            store.state.newstate = true;
            let PtaddRoutes = platformRouterMap;
            // PtaddRoutes[0].children[0].redirect = {name:platformRouterMap[0].children[0].children[0].name};
            function PtgetRouteData(routes){
              let Ptchilds = [];
              for(var i=0;i<routes.length;i++){
                for(var j=0;j<PtbackMenus.length;j++){
                  if(routes[i].name == PtbackMenus[j].routeName){
                    if(PtbackMenus[j].add){
                      routes[i].meta.add = 1;
                    }
                    if(PtbackMenus[j].edit){
                      routes[i].meta.edit = 1;
                    }
                    if(PtbackMenus[j].check){
                      routes[i].meta.check = 1;
                    }
                    if(PtbackMenus[j].delete){
                      routes[i].meta.delete = 1;
                    }
                    if(PtbackMenus[j].export){
                      routes[i].meta.export = 1;
                    }
                    if(PtbackMenus[j].import){
                      routes[i].meta.import = 1;
                    }
                    if(PtbackMenus[j].upload){
                      routes[i].meta.upload = 1;
                    }
                    if(PtbackMenus[j].fileDelete){
                      routes[i].meta.fileDelete = 1;
                    }
                    if(routes[i].children){
                      // routes[i].redirect.name = routes[i].children[0].name
                      routes[i].children=PtgetRouteData(routes[i].children)
                    }
                    Ptchilds.push(routes[i])
                  }
                }
              }
              return Ptchilds;
            }
            PtaddRoutes[0].children = PtgetRouteData(platformRouterMap[0].children);
            let arr = res.data.data.routelist;
            // 只有设置菜单  [{name:'',chinname:''},{name:'Setup',chinname:'UserList'}]
            // 首页菜单都有  [{name:'Mindex',chinname:'MainPage'},{name:'',chinname:'UserList'}]
            // 只有首页   [{name:'Mindex',chinname:'MainPage'},{name:'',chinname:''}]
            for(let i=0;i<arr.length;i++){
              if(arr[i].name){
                PtaddRoutes[0].redirect.name = arr[i].name;
              } 
              if(arr[i].chinname && PtaddRoutes[0].children.length >= 2){
                PtaddRoutes[0].children[i].redirect = {name:arr[i].chinname};
              }else if(arr[i].chinname && PtaddRoutes[0].children.length < 2){
                PtaddRoutes[0].children[0].redirect = {name:arr[i].chinname};
              }
            }
            store.commit('setplatformRoutes',PtaddRoutes);
            router.matcher = createRouter().matcher;
            router.addRoutes(PtaddRoutes);
            router.addRoutes(afterRouteMap);
            router.options.routes = store.state.routes;
          }else{
            store.state.newstate = false;
          }
        }
      })
    },


    // 政府路由菜单
    asynSetgoverRoutes:(store)=>{
      var addRoutes = governmentRouterMap;
      store.commit('setgoverRoutes',addRoutes);
      router.matcher = createRouter().matcher;
      router.addRoutes(addRoutes);
      router.addRoutes(afterRouteMap);
      router.options.routes = store.state.routes;
    }
  },
  modules: {
  }
})
