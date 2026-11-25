<template>
    <div>
        <!-- tab选项卡 -->
        <div class="menu-list">
            <div
            class="totalTab"
            :class="{ active: item.id === sel }"
            v-for="item in tabs"
            :key="item.id"
            @click="select(item)"
            style="cursor: pointer"
            >
            <span>{{ item.label }}</span>
            <div
                class="stripnone"
                :style="item.width"
                :class="{ strip: item.id === sel }"
            ></div>
            </div>
        </div>
        <div class="parNavs">
            <!-- 会员详情 -->
            <div v-show="sel === 1">
                <div class="parNavstrip">
                    <el-col :span="13" class="par-right">
                        <el-date-picker
                        clearable
                        size="small"
                        @change="usrTie"
                        v-model="usrTimArray"
                        type="daterange"
                        unlink-panels
                        range-separator="~"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd"
                        style="width:40%;float:right;margin-right: 10px;">
                        </el-date-picker>
                        <div class="ljnum">筛选查询:</div>
                    </el-col>
                    <!-- 会员详情表格 -->
                    <div class="pur-nav"
                    v-loading="loading"
                    element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                            height="100%"
                            border
                            :data="hairtableData"
                            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                            style="width:100%;">
                            <el-table-column prop="id" v-if="show"></el-table-column>
                            <el-table-column label="访问记录详情">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.visitCount>0">{{scope.row.user_phone}} 访问会员页面<span v-if="scope.row.visitCount>5" style="color:#CF1724;">{{scope.row.visitCount}}</span><span v-else>{{scope.row.visitCount}}</span>次;</span>
                                    <span v-if="scope.row.clickCount>0">点击会员支付按钮<span v-if="scope.row.clickCount>5" style="color:#CF1724;">{{scope.row.clickCount}}</span><span v-else>{{scope.row.clickCount}}</span>次</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="最后一次访问时间" sortable>
                                <template slot-scope="scope">
                                  <span v-if="scope.row.visitTime">{{scope.row.visitTime|Times}}</span>
                                </template>
                            </el-table-column>
                            <!-- 暂无数据 -->
                            <div slot="empty" class="dataAvailable" v-if="dataAva">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="" />
                                <p>暂无数据!</p>
                            </div>
                            </div>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="par-pages">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="hairsize"
                        :current-page.sync="haircurrent"
                        :pager-count="5"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="hairtotal"
                        @current-change='hairgetpurcon'
                        @size-change='hairgetpurcon'>
                        </el-pagination>
                    </div>
                    </div>
                </div>
            </div>
            <!-- 煤种详情 -->
            <div v-show="sel === 2">
                <div class="parNavstrip">
                    <el-col :span="13" class="par-right">
                        <el-date-picker
                        clearable
                        size="small"
                        @change="usrTie"
                        v-model="usrTimArray"
                        type="daterange"
                        unlink-panels
                        range-separator="~"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd"
                        style="width:40%;float:right;margin-right: 10px;">
                        </el-date-picker>
                        <div class="ljnum">筛选查询:</div>
                    </el-col>
                    <!-- 煤种详情表格 -->
                    <div class="pur-nav"
                    v-loading="loadinghair"
                    element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                        height="100%"
                        border
                        :data="collecttableData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                        <el-table-column prop="id" v-if="show"> </el-table-column>
                        <el-table-column label="访问记录详情" width="130">
                            <template slot-scope="scope">
                                <span v-if="scope.row.callCount != 0" @click="openphone(scope.row)" style="color:#ea222e;cursor: pointer;">拨打电话记录({{scope.row.callCount}})</span>
                                <span v-else>--</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="访问记录详情">
                            <template slot-scope="scope">
                                <span v-if="scope.row.visitCount">{{scope.row.user_phone}}对{{scope.row.enterprise_name}}{{scope.row.warehouse_nb}}访问了<span v-if="scope.row.visitCount>5" style="color:#CF1724;">{{scope.row.visitCount}}</span><span v-else>{{scope.row.visitCount}}</span>次;</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="最后一次访问时间" sortable>
                                <template slot-scope="scope">
                                  <span v-if="scope.row.visitTime">{{scope.row.visitTime|Times}}</span>
                                </template>
                            </el-table-column>
                        <!-- 暂无数据 -->
                        <div slot="empty" class="dataAvailable" v-if="dataAvahair">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="" />
                                <p>暂无数据!</p>
                            </div>
                        </div>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="par-pages">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="collectsize"
                        :pager-count="5"
                        :current-page.sync="collectcurrent"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="collecttotal"
                        @current-change='collectgetpurcon'
                        @size-change='collectgetpurcon'>
                        </el-pagination>
                    </div>
                    </div>
                </div>
            </div>
            <!-- 售卖租赁 -->
            <div v-show="sel === 3">
                <div class="parNavstrip">
                    <el-col :span="13" class="par-right">
                        <el-date-picker
                        clearable
                        size="small"
                        @change="usrTie"
                        v-model="usrTimArray"
                        type="daterange"
                        unlink-panels
                        range-separator="~"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd"
                        style="width:40%;float:right;margin-right: 10px;">
                        </el-date-picker>
                        <div class="ljnum">筛选查询:</div>
                    </el-col>
                    <!-- 售卖租赁点表格 -->
                    <div class="pur-nav"
                        v-loading="loadingpretend"
                        element-loading-text = "加载中...">
                    <div class="pur-table">
                        <el-table
                        height="100%"
                        border
                        :data="pretendtableData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                        <el-table-column prop="id" v-if="show"> </el-table-column>
                        <el-table-column label="访问记录详情">
                            <template slot-scope="scope">
                                <span>{{scope.row.createTime}}</span>
                                <span v-if="scope.row.visitCount>0">
                                    {{scope.row.user_phone}}访问了租赁售类页面
                                    <span v-if="scope.row.visitCount>5" style="color:#CF1724;">{{scope.row.visitCount}}</span><span v-else>{{scope.row.visitCount}}</span>
                                    次;
                                </span>
                                <span v-if="scope.row.callCount>0">拨打了客服电话15536587788电话
                                    <span v-if="scope.row.callCount>5" style="color:#CF1724;">{{scope.row.callCount}}</span><span v-else>{{scope.row.callCount}}</span>次</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="最后一次访问时间" sortable>
                            <template slot-scope="scope">
                                <span v-if="scope.row.visitTime">{{scope.row.visitTime|Times}}</span>
                            </template>
                        </el-table-column>
                        <!-- 暂无数据 -->
                        <div slot="empty" class="dataAvailable" v-if="dataAvapretend">
                            <div class="dataAva">
                                <img src="../../assets/zanwushuju.png" alt="" />
                                <p>暂无数据!</p>
                            </div>
                        </div>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="par-pages">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="pretendsize"
                        :pager-count="5"
                        :current-page.sync="pretendcurrent"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="pretendtotal"
                        @current-change='pretendgetpurcon'
                        @size-change='pretendgetpurcon'>
                        </el-pagination>
                    </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 煤种详情弹窗 -->
        <el-dialog
            :closeOnClickModal="false"
            title="煤种详情电话访问记录"
            :visible.sync="paymentShow"
            width="40%"
            :before-close="paymentShowClose">
            <el-table style="cursor: pointer;" :data="DraftsDastas" border v-loading="draftsLoading" element-loading-text = "加载中...">
            <el-table-column label="序号" type="index" width="50"></el-table-column>
            <el-table-column label="访问记录详情">
                <template slot-scope="scope">
                    <span v-if="scope.row.callCount>0 && scope.row.isService == 0 && scope.row.isPtjjr == 0">拨打了联系人{{scope.row.target_user_name}}{{scope.row.target_user_phone}}电话
                        <span v-if="scope.row.callCount>5" style="color:#CF1724;">{{scope.row.callCount}}</span><span v-else>{{scope.row.callCount}}</span>
                        次;
                    </span>
                    <span v-if="scope.row.callCount>0 && scope.row.isService == 1">拨打了客服电话{{scope.row.target_user_phone}}电话
                        <span v-if="scope.row.callCount>5" style="color:#CF1724;">{{scope.row.callCount}}</span><span v-else>{{scope.row.callCount}}</span>
                        次;
                    </span>
                    <span v-if="scope.row.callCount>0 && scope.row.isPtjjr == 1">拨打了经纪人{{scope.row.target_user_name}}电话{{scope.row.target_user_phone}}电话
                        <span v-if="scope.row.callCount>5" style="color:#CF1724;">{{scope.row.callCount}}</span><span v-else>{{scope.row.callCount}}</span>
                        次;
                    </span>
                </template>
            </el-table-column>
            <div slot="empty">
                <div class="dataAva">
                    <p>暂无电话访问记录数据！</p>
                </div>
            </div>
            </el-table>  
        </el-dialog>
    </div>
  </template>
  <script>
  export default {
    data() {
      return {
        show: false,
        sel: '', // 选项卡默认显示tab
        tabs: [
          { label: "会员详情", id: 1, width: "width:48px;" },
          { label: "煤种详情", id: 2, width: "width:48px;" },
          { label: "售卖租赁", id: 3, width: "width:64px;" },
        ],
        usrTimArray:[],             // 时间搜索
        /*
        会员详情数据
        */
        dataAva:false,
        loading:false,
        haircurrent:1,                         // 当前页数
        hairsize:50,                           // 显示条数
        hairtotal:0,                           // 总条数
        hairtableData: [],                     // 会员详情表格数据
  
        /*
        煤种详情数据
        */
        dataAvahair:false,
        loadinghair:false,
        collectcurrent:1,                      // 当前页数
        collectsize:50,                        // 显示条数
        collecttotal:0,                        // 总条数
        collecttableData: [],                  // 煤种详情表格数据

        // 电话访问弹窗
        paymentShow: false,
        draftsLoading:false,            // 加载动画
        DraftsDastas:[],                // 草稿箱数据
  
        /*
        售卖租赁点数据
        */
        dataAvapretend:false,
        loadingpretend:false,
        pretendcurrent:1,                       // 当前页数
        pretendsize:50,                         // 显示条数
        pretendtotal:0,                         // 总条数
        pretendtableData: [],                   // 售卖租赁点表格数据
      };
    },
    mounted(){
      if(this.$route.query.sel == 1){
        this.sel = 1;
        this.hairgetpurcon();
      }else if(this.$route.query.sel == 2){
        this.sel = 2;
        this.collectgetpurcon();
      }else if(this.$route.query.sel == 3){
        this.sel = 3;
        this.pretendgetpurcon();
      }else{
        this.sel = 1;
        this.hairgetpurcon();
      }
    },
    methods: {
        //选项卡
        select(item) {
            this.sel = item.id;
            this.usrTimArray = [];
            if(this.sel == 1){
                this.hairgetpurcon();
            }else if(this.sel == 2){
                this.collectgetpurcon();
            }else if(this.sel == 3){
                this.pretendgetpurcon();
            }
        },
        // 日期搜搜
        usrTie() {
            if(this.sel == 1){
                this.hairtableData = [];
                this.haircurrent = 1;
                this.hairgetpurcon();
            }else if(this.sel == 2){
                this.collecttableData = [];
                this.collectcurrent = 1;
                this.collectgetpurcon();
            }else if(this.sel == 3){
                this.pretendtableData = [];
                this.pretendcurrent = 1;
                this.pretendgetpurcon();
            }
        },
        openphone(data){
            this.paymentShow = true;
            this.draftsLoading = true;
            this.api.busprolog.selCallData({
                userPhone:data.user_phone, 
                proId:data.proId
            })
            .then(res=>{
                this.draftsLoading = false;
                if(res.data.code == 200){
                    this.DraftsDastas = res.data.data;
                }
            })
        },
        paymentShowClose(){
            this.paymentShow = false;
            this.DraftsDastas = [];
        },
        /*
        会员详情方法
        */
        hairgetpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                selType:1,
                searchStr:'',
                startTime:this.usrTimArray != null?this.usrTimArray[0]:"",
                endTime:this.usrTimArray != null?this.usrTimArray[1]:"",
                page:{
                    current:this.haircurrent,
                    size:this.hairsize
                }
            }
            // 渲染表格
            this.api.busprolog.selViata(listByPageData)
            .then(res=>{
            this.loading = false;
            if(res.data.code == 200){
                if(res.data.data.records != ""){
                this.hairtableData = res.data.data.records;
                this.hairtotal = res.data.data.total;
                }else{
                    this.dataAva = true;
                }
            }
            })
        },
    
        /*
        煤种详情方法
        */
        collectgetpurcon(){
            this.loadinghair = true;
            // 传参
            var listByPageData = {
                selType:2,
                searchStr:'',
                startTime:this.usrTimArray != null?this.usrTimArray[0]:"",
                endTime:this.usrTimArray != null?this.usrTimArray[1]:"",
                page:{
                    current:this.collectcurrent,
                    size:this.collectsize
                }
            }
            // 渲染表格
            this.api.busprolog.selViata(listByPageData)
            .then(res=>{
                this.loadinghair = false;
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                    this.collecttableData = res.data.data.records;
                    this.collecttotal = res.data.data.total;
                    }else{
                    this.dataAvahair = true;
                    }
                }
            })
        },

    
        /*
        售卖租赁点方法
        */
        pretendgetpurcon(){
            this.loadingpretend = true;
            // 传参
            var listByPageData = {
                selType:3,
                searchStr:'',
                startTime:this.usrTimArray != null?this.usrTimArray[0]:"",
                endTime:this.usrTimArray != null?this.usrTimArray[1]:"",
                page:{
                    current:this.pretendcurrent,
                    size:this.pretendsize
                }
            }
            // 渲染表格
            this.api.busprolog.selViata(listByPageData)
            .then(res=>{
            this.loadingpretend = false;
            if(res.data.code == 200){
                if(res.data.data.records != ""){
                this.pretendtableData = res.data.data.records;
                this.pretendtotal = res.data.data.total;
                }else{
                this.dataAvapretend = true;
                }
    
            }
            })
        },
    },
  };
  </script>
  <style lang="less" scoped>
  // 草稿箱样式开始
    .Draftspages{
        width: 100%;
        height: 40px;
        padding-top: 10px;
        box-sizing: border-box;
        .el-pagination{
            float: right;
        }
    }
  // 选项卡开始
  .menu-list {
    width: 100%;
    height: 48px;
    background-color: #fff;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    border: 1px solid rgba(102, 102, 102, 0.15);
  }
  .totalTab {
    width: 7.6%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin: 9px 0 0 35px;
  }
  .totalTab.active {
    color: #ea222e;
  }
  .stripnone {
    height: 3px;
    background-color: #ea222e;
    margin: 11px auto 0;
    display: none;
  }
  .strip {
    display: block;
  }
  .pur-nav {
    width: calc(100% - 32px);
    height: calc(100vh - 144px);
    background-color: #fff;
    margin: 8px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding: 8px 8px 8px 8px;
    .pur-table{
      width:100%;
      height: 94%;
    } 
    .par-pages {
      width: 100%;
      height: 6%;
      padding: 6px 0 0 12px;
      box-sizing:border-box;
      .el-pagination {
        float: right;
      }
    }
  }
  //选项卡结束
  // 暂无数据样式开始
  .dataAvailable {
    width: 100%;
    height: 245px;
    position: relative;
    margin-top:185px;
    margin-bottom:279px;
    .dataAva {
      position: absolute;
      top: 75%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;
      img {
        width: 371px;
        height: 200px;
        margin: auto;
      }
      p {
        font-size: 16px;
        line-height: 24px;
        color: #333;
      }
      .el-button {
        margin-bottom: 100px;
      }
    }
  }
  // 暂无数据样式结束
  .parNavstrip {
    position: relative;
    .par-right {
      position: absolute;
      top: -50px;
      right: 24px;
    .ljnum{
        float: right;
        padding-right:16px;
        font-size: 16px;
        line-height: 32px;
    }
      .el-input {
        width: 200px;
        float: right;
      }
      .el-button {
        margin-left: 16px;
        float: right;
      }
    }
  }
  </style>