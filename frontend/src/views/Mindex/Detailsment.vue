<template>
  <!-- 关联明细单页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 关联明细信息列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
                <b>关联记录</b>
            </div>
            <div class="Advanced">
                <el-date-picker
                @change="DataTime"
                size="small"
                v-model="keyupTime"
                type="date"
                value-format="yyyy-MM-dd"
                placeholder="操作日期">
                </el-date-picker>
                <el-input
                @keyup.enter.native="jshxpurSearchs"
                size="small"
                placeholder="搜索"
                prefix-icon="el-icon-search"
                v-model="jshxpurSearch">
                </el-input>
                <el-button size="small" @click="detile">解除关联</el-button>
            </div>
            <!-- 表格 -->
            <div class="pur-table">
                <el-table
                    v-loading="settlementloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="settlementtableData"
                    highlight-current-row
                    @current-change="settlementhandleCurrentChange"
                    @selection-change="purChange"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333'}"
                    style="width: 100%;cursor: pointer;">
                    <el-table-column type="selection"></el-table-column>
                    <el-table-column type="index" width="50" label="序号" :index="indexMethod">
                    </el-table-column>
                    <el-table-column prop="type" label="结算类型" :show-overflow-tooltip="true" width="100">
                        <template slot-scope="scope">
                            <span v-if="scope.row.type ==1">采购</span>
                            <span v-if="scope.row.type ==2">销售</span>
                        </template> 
                    </el-table-column>
                    <el-table-column prop="settNum" label="结算数量" :show-overflow-tooltip="true" width="100">
                    </el-table-column>
                    <el-table-column prop="billNum" label="发票数量" :show-overflow-tooltip="true" width="110">
                    </el-table-column>
                    <el-table-column prop="createTime" label="操作日期" sortable :show-overflow-tooltip="true">
                    </el-table-column>
                    <el-table-column prop="operater" label="操作人" :show-overflow-tooltip="true" >
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="pur-pages">
                <el-pagination
                background
                :page-sizes="[10, 30, 50, 100, 200]"
                :page-size.sync="settlementsize"
                :pager-count="5"
                :current-page.sync="settlementcurrent"
                layout="total, sizes, prev, pager, next, jumper"
                :total="settlementtotal"
                @current-change='settlementgetpurcon'
                @size-change='settlementgetpurcon'>
                </el-pagination>
            </div>
        </div>
        <div class="rece-right"><span>《</span></div>
        <!-- 结算 发票tab列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title" >
                <!-- tab选项卡 -->
                <div class="totalTab"
                    :class="{ active: item.id === sel }"
                    v-for="item in tabs"
                    :key="item.id"
                    @click="select(item)" style="cursor: pointer;">
                    <span>{{ item.label }}</span>
                    <div class="stripnone" :style="item.width"
                    :class="{ strip: item.id === sel }">
                    </div>
                </div>
            </div>
            <!-- 结算 -->
            <div class="adBox" v-show="sel === 1">
                <div class="sjhxjeMoeny">
                    <span v-if="contmoeny != 0">结算总价合计: {{contmoeny|currency}} 元</span>
                    <span v-else>结算总价合计: 0.00 元</span>
                </div>   
                <div class="Advanced">
                    <el-input
                    @keyup.enter.native="ClickOppositeName"
                    size="small"
                    :placeholder="'搜索'+ customer"
                    prefix-icon="el-icon-user"
                    v-model="SearchoppositeName"
                    >
                    </el-input>
                    <el-date-picker
                    @change="Time"
                    size="small"
                    v-model="TimeArray"
                    type="month"
                    value-format="yyyy-MM"
                    placeholder="所属月份">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="search"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="purSearch"
                    >
                    </el-input>
                </div>
                <!-- 结算表格 -->
                <div class="tab-table">
                    <el-table
                    v-loading="loading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="tableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                        <el-table-column prop="contractCode" label="合同编号" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="supplierName" :label="customer" sortable :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="settlementCount" label="结算数量" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span>{{scope.row.settlementCount}} 吨</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="month" label="所属月份" :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <!-- <el-table-column prop="goodsName" label="结算日期" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span>{{scope.row.setStartTime}} ~ {{scope.row.setEndTime}}</span>
                            </template>
                        </el-table-column> -->
                        <el-table-column prop="sumPrice" label="结算总价" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice !=0">{{scope.row.sumPrice|currency}} 元</span>
                                <span v-else>0.00 元</span>
                            </template> 
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <div style="float:right;width:75%;">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="size"
                        :pager-count="5"
                        :current-page.sync="current"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="total"
                        @current-change='getpurcon'
                        @size-change='getpurcon'>
                        </el-pagination>
                    </div>
                </div>   
            </div>  
            <!-- 发票 -->
            <div class="adBox" v-show="sel === 2"> 
                <div class="sjhxjeMoeny">
                    <span v-if="watertcontmoeny != 0">含税金额合计: {{watertcontmoeny|currency}} 元</span>
                    <span v-else>含税金额合计: 0.00 元</span>
                </div> 
                <div class="Advanced">
                    <el-input
                    @keyup.enter.native="waterClickOppositeName"
                    size="small"
                    :placeholder="'搜索' + Seller + '方企业名称'"
                    prefix-icon="el-icon-user"
                    v-model="waterSearchoppositeName"
                    >
                    </el-input>
                    <el-date-picker
                    size="small"
                    @change="waterDataTime"
                    v-model="waterkeyupTime"
                    type="daterange"
                    unlink-panels
                    range-separator="~"
                    start-placeholder="开票开始日期"
                    end-placeholder="开票结束日期"
                    value-format="yyyy-MM-dd"
                    style="width:30%;float:left;margin-right:2%;">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="watersearch"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="waterpurSearch"
                    >
                    </el-input>
                </div>
                <!-- 发票表格数据 -->
                <div class="tab-table">
                    <el-table
                    v-loading="waterloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="watertableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                        <el-table-column prop="companyName" :label="Seller + '方企业名称'" sortable :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="openDate" label="开票日期" sortable width="110" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="baseMoney" label="金额" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.baseMoney != 0">{{scope.row.baseMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template> 
                        </el-table-column>
                        <el-table-column prop="taxAmount" label="税额" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.taxAmount != 0">{{scope.row.taxAmount|currency}} 元</span>
                                <span v-else>---</span>
                            </template> 
                        </el-table-column>
                        <el-table-column prop="includMoney" label="含税金额" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.includMoney != 0">{{scope.row.includMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <div style="float:right;width:75%;">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="watertsize"
                        :pager-count="5"
                        :current-page.sync="watertcurrent"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="waterttotal"
                        @current-change='watertgetpurcon'
                        @size-change='watertgetpurcon'>
                        </el-pagination>
                    </div>
                </div>     
            </div>     
        </div>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            customer:'供应商 / 客户',
            Seller:'销 / 购',
            sel:1,                                          // 选项卡默认显示tab
            tabs: [
                { label: "结算", id: 1 ,width:"width:32px;"},
                { label: "发票", id: 2 ,width:"width:32px;"},
            ],
            http:'',
            show:false,
            purTitle: "",                                   // 标题
            logo: "",                                       // 抽屉标题

            /* 
            关联明细信息数据列表
            */
            multipleSelection:[],                           // 选中数据集合
            settlementloading:false,
            settlementcurrent:1,                            // 当前页数
            settlementsize:50,                              // 显示条数 
            settlementtotal:0,                              // 总条数
            keyupTime:'',                                   // 时间搜索
            jshxpurSearch:"",                               // 全部搜索
            settlementtableData: [],                        // 表格数据
            multiplesettlement:{},                          // 选中单条数据

            /* 
            结算数据列表
            */
            banks:'',                           // 银行搜索    
            purSearch:'',                       // 全部搜索
            SearchoppositeName:'',              // 供应商/客户搜索
            TimeArray:'',                       // 结算月份搜索
            contmoeny:0,                        // 结算合计金额  
            loading:false,                      // 加载状态
            current:1,                          // 当前页数
            size:50,                            // 显示条数
            total:0,                            // 总条数
            tableData: [],                      // 表格数据

            /* 
            发票数据列表
            */
            waterpurSearch:'',                       // 全部搜索
            waterkeyupTime: [],                      // 时间段搜索
            waterSearchoppositeName:'',              // 客户/供应商搜索
            watertcontmoeny:0,                       // 合计含税金额
            waterloading:false,                      // 加载状态
            watertcurrent:1,                         // 当前页数
            watertsize:50,                           // 显示条数
            waterttotal:0,                           // 总条数
            watertableData: [],                      // 表格数据
        };
    },
    created() {},
    mounted(){
        this.tblehei = "610";
        this.purTitle = this.$route.meta.title;                 // 标题
        this.logo = sessionStorage.getItem('companyName')       // 企业名称
        this.http  = this.$store.state.http;                    // IP地址
        this.settlementgetpurcon();                             // 关联记录信息
    },
    methods: {
        // 关联记录表格选中数据
        purChange(val) {
            this.multipleSelection = val.map(row => row.id); // 原先只追加不清:取消勾选的记录仍在 ids 里,会被一起解除关联
        },
        // 解除关联按钮
        detile(){ 
            if(this.multipleSelection.length > 0){
                this.$confirm("确定解除选中的 " + this.multipleSelection.length + " 条关联吗?", "提示", { type: 'warning' }) // 解除后不可恢复,原先没有确认
                .then(() => this.api.settRec.relieve({ids:this.multipleSelection}))
                .then(res=>{
                    if(res.data.code == 200){
                        this.multipleSelection = []; // 不清的话下次不勾选直接点会重发旧 id
                        this.$message({
                            message: '取消关联成功',
                            type: 'success'
                        });
                        this.jshxpurSearch = "";
                        this.keyupTime = "";
                        this.TimeArray = '';
                        this.SearchoppositeName = "";
                        this.purSearch = "";
                        this.waterSearchoppositeName = "";
                        this.waterpurSearch = "";
                        this.waterkeyupTime == null;
                        this.settlementtableData = [];
                        this.settlementgetpurcon();
                        this.tableData = [];
                        this.getpurcon();                       // 结算列表
                        this.watertableData = [];
                        this.watertgetpurcon();                 // 票据列表 
                    }
                }).catch(() => {}); // 取消确认
            }else{
                this.$message({
                    message: '至少选中一条关联记录数据!',
                    type: 'warning'
                }); 
            }
        },
        // 序号索引
        indexMethod(index) {
            return index + 1;
        },
        // 关联记录信息tab选中条数
        settlementhandleCurrentChange(currentRow) {
            if(!currentRow) return; // 列表刷新后原高亮行不在新数据里时,回调收到 null
            this.multiplesettlement = currentRow;
            if(this.multiplesettlement.type == 1){
                this.tabs = [
                    { label: "采购结算", id: 1 ,width:"width:64px;"},
                    { label: "进项发票", id: 2 ,width:"width:64px;"}
                ]
                this.customer = "供应商";
                this.Seller = "销";
            }else if(this.multiplesettlement.type == 2){
                this.tabs = [
                    { label: "销售结算", id: 1 ,width:"width:64px;"},
                    { label: "销项发票", id: 2 ,width:"width:64px;"},
                ]
                this.customer = "客户";
                this.Seller = "购"
            }
            this.sel = 1;
            this.getpurcon();                           // 结算列表
        },
        // 结算 票据选项卡
        select(item) {
            this.TimeArray = '';
            this.SearchoppositeName = "";
            this.purSearch = "";
            this.waterSearchoppositeName = "";
            this.waterpurSearch = "";
            this.waterkeyupTime == null;
            this.sel = item.id;
            if(this.sel == 1){
                this.tableData = [];
                this.getpurcon();                       // 结算列表
            }
            if(this.sel == 2){
                this.watertableData = [];
                this.watertgetpurcon();                 // 票据列表
            }
        },

        /*
        关联明细信息列表展示功能
        */
        // 关联明细信息表格数据
        settlementgetpurcon(){
            this.settlementloading = true;
            // 传参
            var listByPageData = {
                searchStr:this.jshxpurSearch,
                relDate:this.keyupTime,
                page:{
                    current:this.settlementcurrent,
                    size:this.settlementsize
                }
            }
            // 渲染表格
            this.api.settRec.all(listByPageData)
            .then(res=>{
                this.settlementloading = false; 
                if(res.data.code == 200){
                    this.settlementtableData = res.data.data.records;
                    this.settlementtotal = res.data.data.total;
                }
            })
        },
        // 搜索全部功能
        jshxpurSearchs() {
            this.settlementtableData = [];
            this.settlementcurrent = 1;
            this.settlementgetpurcon();
        },
        // 日期搜索
        DataTime() {
            this.settlementtableData = [];
            this.settlementcurrent = 1;
            this.settlementgetpurcon();
        },

        /*
        结算列表展示功能
        */
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                relationId:this.multiplesettlement.id,
                month:this.TimeArray,
                supplierName:this.SearchoppositeName,
                searchStr:this.purSearch,
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageData)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.tableData = res.data.data.records;
                    this.total = res.data.data.total;
                }
            })
            // 结算总价合计
            this.contmoeny = 0;
            this.api.capital.amount({relationId:this.multiplesettlement.id})
            .then(res=>{
                if(res.data.code == 200){
                   this.contmoeny = res.data.data.sumMoney;
                }
            })
        },
        // 供应商搜索
        ClickOppositeName() {
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
        },
        // 结算时间搜索
        Time() {
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
        },
        // 全部搜索
        search() {
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
        },

        /*
        发票列表展示功能
        */
        // 表格数据
        watertgetpurcon(){
            this.waterloading = true;
            if(this.waterkeyupTime == null){
                // 传参
                var listByPageData = {
                    relationId:this.multiplesettlement.id,
                    StartTime:'',
                    EndTime:'',
                    companyName:this.waterSearchoppositeName,
                    searchStr:this.waterpurSearch,
                    page:{
                        current:this.watertcurrent,
                        size:this.watertsize
                    }
                }
            }else{
                // 传参
                var listByPageData = {
                    relationId:this.multiplesettlement.id,
                    startTime:this.waterkeyupTime[0],
                    endTime:this.waterkeyupTime[1],
                    companyName:this.waterSearchoppositeName,
                    searchStr:this.waterpurSearch,
                    page:{
                        current:this.watertcurrent,
                        size:this.watertsize
                    }
                }
            }
            // 渲染表格
            this.api.bill.all(listByPageData)
            .then(res=>{
                this.waterloading = false;
                if(res.data.code == 200){
                    this.watertableData = res.data.data.records;
                    this.waterttotal = res.data.data.total;
                }
            })
            // 含税金额合计
            this.watertcontmoeny = 0;
            this.api.bill.amount({relationId:this.multiplesettlement.id})
            .then(res=>{
                if(res.data.code == 200){
                   this.watertcontmoeny = res.data.data.sumMoney;
                }
            })
        },
        // 日期搜索
        waterDataTime() {
            this.watertableData = [];
            this.watertcurrent = 1;
            this.watertgetpurcon();
        },
        // 企业名称搜索
        waterClickOppositeName() {
            this.watertableData = [];
            this.watertcurrent = 1;
            this.watertgetpurcon();
        },
        // 全部搜索
        watersearch() {
            this.watertableData = [];
            this.watertcurrent = 1;
            this.watertgetpurcon();
        },
    },
};
</script>
<style lang="less" scoped>

// 选项卡开始
.totalTab {
    width: auto;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    margin-right:3%;
    color: #333;
    text-align: center;
}
.totalTab.active {
    color: #ea222e;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:4px auto 0;
    border-radius:25px;
    display: none;
}
.strip {
    display: block;
}
.adBox{
    height: 100%;
    position: relative;
    .tab-table{
        height: 76%;
    }
    .sjhxjeMoeny{
        position: absolute;
        top: -65px;
        right: 32px;
        span{
            font-size: 14px;
            font-weight: 600;
            color: #333;
            line-height:50px;
        }  
    }
}
//选项卡结束
// 头部开始
.pur-top {
  width: 100%;
  height: 50px;
  padding: 9px 24px 9px 32px;
  background-color: #fff;
  box-sizing: border-box;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  border: 1px solid rgba(102, 102, 102, 0.15);
  .pur-size {
    font-size: 16px;
    line-height: 32px;
    color: #333;
    font-weight: 600;
  }
  .pur-right {
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
// 头部结束
// 内容开始
.pur-nav {
    width: calc(100% - 32px);
    height: calc(100vh - 142px);
    background-color: #fff;
    margin: 8px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding: 8px;
    .rece{
        width: 49%;
        height: 100%;
        border: 1px solid #E4E4E4;
        float: left;
        box-sizing: border-box;
        position: relative;
        .pur-title{
            width: 100%;
            height:67px;
            padding:16px 0 24px 32px;
            box-sizing: border-box;
            b{
                color: #333;
            }
        }
        .Advanced{
            width: 100%;
            height: 48px;
            padding:0 0 16px 32px;
            box-sizing: border-box;
            .el-select{
                float: left;
                width: 24%;
                margin-right: 1%;
            }
            .el-input {
                width: 24%;
                float: left;
                margin-right: 1%;
            }
            .el-date-picker{
                width: 24%;
                float: left;
                margin-right: 1%;
            }
        }
        .pur-table{
            width: 100%;
            height: 76%;
        }
        .pur-pages {
            width: 100%;
            height: 42px;
            padding:10px 18px 0 24px;
            box-sizing: border-box;
            .el-pagination {
                float: right;
            }
        }
    }
    .rece-right{
        width: 2%;
        height: 89%;
        float: left;
        transform: rotate(180deg);
        overflow: hidden;
        position: relative;
        box-sizing: border-box;
        span{
            position: absolute;
            right:20%;
            font-size: 24px;
            color: #999;
            font-weight: 300;
            line-height: 800px;
        }
    }
}
// 内容结束
</style>