<template>
  <!-- 核销明细单页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 核销明细列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
                <b>核销明细</b>
            </div>
            <div class="Advanced">
                <el-date-picker
                @change="DataTime"
                size="small"
                v-model="keyupTime"
                type="date"
                value-format="yyyy-MM-dd HH:mm:ss"
                placeholder="选择核销日期">
                </el-date-picker>
                <el-input
                @keyup.enter.native="jshxpurSearchs"
                size="small"
                placeholder="请输入内容"
                prefix-icon="el-icon-search"
                v-model="jshxpurSearch">
                </el-input>
                <el-button size="small" @click="writeOff">取消核销</el-button>
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
                    :header-cell-style="{ background: '#F6F6F6', color: '#333'}"
                    style="width: 100%;cursor: pointer;"
                    @selection-change="settlementChange"
                    ref="contractTable"
                    :row-key="settlementkeyId">
                    <el-table-column type="selection" width="40" :reserve-selection="true"> </el-table-column>
                    <el-table-column prop="smentNum" label="结算单数量" :show-overflow-tooltip="true" width="100">
                    </el-table-column>
                    <el-table-column prop="detailNum" label="银行流水" :show-overflow-tooltip="true" >
                    </el-table-column>
                    <el-table-column prop="type" label="类型" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.type == 1">应收单</span>
                            <span v-if="scope.row.type == 2">应付单</span>
                        </template>  
                    </el-table-column>
                    <el-table-column prop="money" label="核销金额" sortable :show-overflow-tooltip="true" width="110">
                        <template slot-scope="scope">
                            <span v-if="scope.row.money !=0">{{scope.row.money|currency}} 元</span>
                            <span v-else>0.00 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="核销时间" :show-overflow-tooltip="true" >
                    </el-table-column>
                    <el-table-column prop="userName" label="操作人" :show-overflow-tooltip="true" >
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="pur-pages">
                <div style="float:right;width:75%;">
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
        </div>
        <div class="rece-right"><span>《</span></div>
        <!-- 结算 流水tab列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
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
                    <span v-if="contmoeny != 0">总价合计: {{contmoeny|currency}} 元</span>
                    <span v-else>总价合计: 0.00 元</span>
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
                        <el-table-column prop="systemCode" label="合同编号" sortable :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="goodsName" label="货物名称" sortable :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="supplierName" :label="customer" sortable :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="khxje" label="剩余可核销金额" sortable :show-overflow-tooltip="true" width="160">
                            <template slot-scope="scope">
                                <span v-if="scope.row.khxje !=0">{{scope.row.khxje|currency}} 元</span>
                                <span v-else>0.00 元</span>
                            </template>  
                        </el-table-column>
                        <el-table-column prop="sumPrice" label="结算总价" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice !=0">{{scope.row.sumPrice|currency}} 元</span>
                                <span v-else>0.00 元</span>
                            </template> 
                        </el-table-column>
                        <el-table-column prop="hxStatus" label="核销状态" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.hxStatus == 1">未核销</span>
                                <span v-if="scope.row.hxStatus == 2">核销中</span>
                                <span v-if="scope.row.hxStatus == 3">已核销</span>
                            </template> 
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages" v-if="paginationshow">
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
            <!-- 流水 -->
            <div class="adBox" v-show="sel === 2"> 
                <div class="sjhxjeMoeny">
                    <span v-if="watertcontmoeny != 0">金额合计: {{watertcontmoeny|currency}} 元</span>
                    <span v-else>金额合计: 0.00 元</span>
                </div> 
                <div class="Advanced">
                    <el-select @change="waterkeyupBanks" v-model="waterBanks" filterable placeholder="请选择银行名称" size="small" :clearable='true'>
                        <el-option
                        v-for="item in optionsBanks"
                        :key="item.bank_name"
                        :label="item.bank_name"
                        :value="item.bank_name">
                        </el-option>
                    </el-select>
                    <el-date-picker
                    @change="waterDataTime"
                    size="small"
                    v-model="waterkeyupTime"
                    type="date"
                    value-format="yyyy-MM-dd"
                    placeholder="选择交易日期">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="waterClickOppositeName"
                    size="small"
                    placeholder="搜索对方户名"
                    prefix-icon="el-icon-user"
                    v-model="waterSearchoppositeName"
                    >
                    </el-input>
                    <el-input
                    @keyup.enter.native="watersearch"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="waterpurSearch"
                    >
                    </el-input>
                </div>
                <!-- 流水表格 -->
                <div class="tab-table">
                    <el-table
                    v-loading="waterloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="watertableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                        <el-table-column prop="accountName" label="账户名称" sortable :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="transTime" label="交易日期" sortable width="110" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="oppositeName" label="对方户名" sortable :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="khxje" label="剩余可核销金额" sortable :show-overflow-tooltip="true" width="160">
                            <template slot-scope="scope">
                                <span v-if="scope.row.khxje !=0">{{scope.row.khxje|currency}} 元</span>
                                <span v-else>0.00 元</span>
                            </template>  
                        </el-table-column>
                        <el-table-column prop="money" label="金额" sortable :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.money !=0">{{scope.row.money|currency}} 元</span>
                                <span v-else>0.00 元</span>
                            </template> 
                        </el-table-column>
                        <el-table-column prop="hxStatus" label="核销状态" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.hxStatus == 1">未核销</span>
                                <span v-if="scope.row.hxStatus == 2">核销中</span>
                                <span v-if="scope.row.hxStatus == 3">已核销</span>
                            </template> 
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages" v-if="watertpaginationshow">
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
            // tblehei:'610',
            sel:1,   // 选项卡默认显示tab
            tabs: [
                { label: "结算", id: 1 ,width:"width:32px;"},
                { label: "流水", id: 2 ,width:"width:32px;"},
            ],
            http:'',
            show:false,
            purTitle: "",                                   // 标题
            logo: "",                                       // 抽屉标题
            writeOffVal:"",                                 // 表格选中数据

            /* 
            核销明细数据列表
            */
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
            SearchoppositeName:'',              // 供应商搜索
            contmoeny:0,                        // 结算合计金额  
            paginationshow:true,                // 无数据隐藏分页
            loading:false,                      // 加载状态
            current:1,                          // 当前页数
            size:50,                            // 显示条数
            total:0,                            // 总条数
            tableData: [],                      // 表格数据

            /* 
            流水数据列表
            */
            waterBanks:'',                           // 银行搜索    
            optionsBanks:[],                         // 银行名称列表数据回显
            waterpurSearch:'',                       // 全部搜索
            waterkeyupTime: '',                      // 时间搜索
            waterSearchoppositeName:'',              // 对方户名搜索
            watertcontmoeny:0,                       // 合计金额  
            watertpaginationshow:true,               // 无数据隐藏分页
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
        this.settlementgetpurcon();                             // 核销明细
        // 获取银行搜索列表
        this.api.chainCap.banks()
        .then(res=>{
            if(res.data.code == 200){
                this.optionsBanks = res.data.data;
            }
        })
    },
    methods: {
        // 表格选中数据
        settlementChange(val){
            this.writeOffVal = val;
        },
        settlementkeyId(row){
            return row.id;
        },
        // 取消核销
        writeOff(){
            console.log(this.writeOffVal)
            if(this.writeOffVal.length > 0){
                this.$confirm('正在恢复核销金额，请谨慎操作！', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning',
                showClose:false,
                }).then(() => {
                    let writarr = [];
                    for(let i=0;i<this.writeOffVal.length;i++){
                        writarr.push({hx_code:this.writeOffVal[i].systemCode})
                    }
                    this.api.chainCap.cancelBalancing(writarr)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this;
                            that.$message({
                                type: "success",
                                message: "取消核销成功!",
                                duration:1000,  
                                onClose(){
                                    that.$refs.contractTable.clearSelection();
                                    that.settlementtableData = [];
                                    that.settlementgetpurcon();
                                }
                            });
                        }
                    })
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已撤回取消核销'
                    });          
                });
            }else{
                this.$message({
                    type: 'warning',
                    message: '请选择要取消的核销金额'
                });  
            }
        },
        // 结算 流水选项卡
        select(item) {
            this.purSearch = "";
            this.SearchoppositeName = "";
            this.waterBanks = "";
            this.waterpurSearch = "";
            this.waterkeyupTime = "";
            this.waterkeyupTime = "";
            this.waterSearchoppositeName = "";
            this.sel = item.id;
            if(this.sel == 1 && this.multiplesettlement.systemCode){
                this.getpurcon();                       // 结算列表
            }else if(this.sel == 2 && this.multiplesettlement.systemCode){
                this.watertgetpurcon();                 // 流水列表
            }
        },
        // 核销明细选中条数
        settlementhandleCurrentChange(currentRow) {
            this.multiplesettlement = currentRow;
            if(this.multiplesettlement.type == 1){
                this.tabs[0].width = "width:64px;"
                this.tabs[1].width = "width:80px;"    
                this.tabs[0].label = "销售结算"
                this.tabs[1].label = "销售收款单"
                this.customer = "客户";
            }else if(this.multiplesettlement.type == 2){
                this.tabs[0].width = "width:64px;"
                this.tabs[1].width = "width:80px;"  
                this.tabs[0].label = "采购结算"
                this.tabs[1].label = "采购付款单"
                this.customer = "供应商";
            }
            this.sel = 1;
            this.getpurcon();                           // 结算列表
        },

        /*
        核销明细列表展示功能
        */
        // 核销明细表格数据
        settlementgetpurcon(){
            this.settlementloading = true;
            // 传参
            var listByPageData = {
                searchStr:this.jshxpurSearch,
                createTime:this.keyupTime,
                page:{
                    current:this.settlementcurrent,
                    size:this.settlementsize
                }
            }
            // 渲染表格
            this.api.count.all(listByPageData)
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
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },
        // 日期搜索
        DataTime() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },

        /*
        结算列表展示功能
        */
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                hxCode:this.multiplesettlement.systemCode,
                oppositeName:this.SearchoppositeName,
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
                    if(res.data.data.records != ""){
                        this.contmoeny = 0;
                        for(let i=0;i<res.data.data.records.length;i++){
                            this.contmoeny += res.data.data.records[i].sumPrice;
                        }
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                    }else{
                        this.paginationshow = false;
                    }
                }
            })
        },
        // 供应商搜索
        ClickOppositeName() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 全部搜索
        search() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },

        /*
        流水列表展示功能
        */
        // 表格数据
        watertgetpurcon(){
            this.waterloading = true;
            // 传参
            var listByPageData = {
                hxCode:this.multiplesettlement.systemCode,
                bankName:this.waterBanks,
                transTime:this.waterkeyupTime,
                oppositeName:this.waterSearchoppositeName,
                searchStr:this.waterpurSearch,
                page:{
                current:this.watertcurrent,
                size:this.watertsize
                }
            }
            // 渲染表格
            this.api.chainCap.all(listByPageData)
            .then(res=>{
                this.waterloading = false;
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.watertcontmoeny = 0;
                        for(let i=0;i<res.data.data.records.length;i++){
                            this.watertcontmoeny += res.data.data.records[i].money;
                        }
                        this.watertableData = res.data.data.records;
                        this.waterttotal = res.data.data.total;
                    }else{
                        this.watertpaginationshow = false;
                    }
                }
            })
        },
        // 银行名称搜索
        waterkeyupBanks() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // 日期搜索
        waterDataTime() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // 对方户名
        waterClickOppositeName() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // 全部搜索
        watersearch() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
    },
};
</script>
<style lang="less" scoped>

// 选项卡开始
.totalTab {
    width:auto;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    margin-right: 3%;
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