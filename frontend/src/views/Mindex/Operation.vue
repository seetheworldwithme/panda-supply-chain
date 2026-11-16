<template>
  <!-- 关联操作单页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
      <!-- <el-col :span="21" class="pur-right">
      </el-col> -->
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 结算信息列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
                <!-- tab选项卡 -->
                <div class="totalTab"
                    :class="{ active: item.id === leftsel }"
                    v-for="item in lefttabs"
                    :key="item.id"
                    @click="leftselect(item)" style="cursor: pointer;">
                    <span>{{ item.label }}</span>
                    <div class="stripnone" :style="item.width"
                    :class="{ strip: item.id === leftsel }">
                    </div>
                </div>
            </div>
            <!-- 采购结算 -->
            <div class="adBox" v-if="leftsel === 1">
                <span class="Totalsettlement" v-if="newwatertcountone != 0">结算总价合计: {{newwatertcountone|currency}} 元</span>
                <span class="Totalsettlement" v-else>结算总价合计: 0.00 元</span>
                <div class="Advanced">
                    <el-input
                    @keyup.enter.native="newwaterClickOppositeNameone"
                    size="small"
                    placeholder="搜索供应商"
                    prefix-icon="el-icon-user"
                    v-model="newwaterSearchoppositeNameone">
                    </el-input>
                    <el-date-picker
                    @change="DataTime"
                    size="small"
                    v-model="keyupTime"
                    type="month"
                    value-format="yyyy-MM"
                    placeholder="所属月份">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="jshxpurSearchs"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="jshxpurSearch">
                    </el-input>
                </div>
                <!-- 表格 -->
                <div class="pur-table">
                    <el-table
                        v-loading="settlementloading"
                        element-loading-text = "加载中..."
                        height="100%"
                        border
                        :data="settlementtableData"
                        @selection-change="purChangeone"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333'}"
                        style="width: 100%;cursor: pointer;">
                        <el-table-column type="selection"> </el-table-column>
                        <el-table-column prop="supplierName" label="供应商" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="month" label="所属月份" :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="settlementCount" label="结算数量" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="unitPrice" label="结算单价" :show-overflow-tooltip="true" >
                            <template slot-scope="scope">
                                <span>{{scope.row.unitPrice}} 元</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="sumPrice" label="结算总价" :show-overflow-tooltip="true" >
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice">{{scope.row.sumPrice|currency}} 元</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <span v-if="newwatertcontmoenyone != 0">选中合计: {{newwatertcontmoenyone|currency}} 元</span>
                    <span v-else>选中合计: 0.00 元</span>
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
            <!-- 销售结算 -->
            <div class="adBox" v-if="leftsel === 2">
                
                <span class="Totalsettlement" v-if="newwatertcount != 0">结算总价合计: {{newwatertcount|currency}} 元</span>
                <span class="Totalsettlement" v-else>结算总价合计: 0.00 元</span>
                <div class="Advanced">
                    <el-input
                    @keyup.enter.native="newwaterClickOppositeName"
                    size="small"
                    placeholder="搜索客户"
                    prefix-icon="el-icon-user"
                    v-model="newwaterSearchoppositeName">
                    </el-input>
                    <el-date-picker
                    @change="newwaterDataTime"
                    size="small"
                    v-model="newwaterkeyupTime"
                    type="month"
                    value-format="yyyy-MM"
                    placeholder="所属月份">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="newwatersearch"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="newwaterpurSearch"
                    >
                    </el-input>
                </div>
                <!-- 销售结算表格 -->
                <div class="tab-table">
                    <el-table
                    v-loading="newwaterloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="newwatertableData"
                    @selection-change="purChange"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                        <el-table-column type="selection"> </el-table-column>
                        <el-table-column prop="supplierName" label="客户" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="month" label="所属月份" :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="settlementCount" label="结算数量" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="unitPrice" label="结算单价" :show-overflow-tooltip="true" >
                            <template slot-scope="scope">
                                <span>{{scope.row.unitPrice}} 元</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="sumPrice" label="结算总价" :show-overflow-tooltip="true" >
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice">{{scope.row.sumPrice|currency}} 元</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <span v-if="newwatertcontmoeny != 0">选中合计: {{newwatertcontmoeny|currency}} 元</span>
                    <span v-else>选中合计: 0.00 元</span>
                    <div style="float:right;width:75%;">
                        <el-pagination
                        background
                        :page-sizes="[10, 30, 50, 100, 200]"
                        :page-size.sync="newwatertsize"
                        :pager-count="5"
                        :current-page.sync="newwatertcurrent"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="newwaterttotal"
                        @current-change='newwatertgetpurcon'
                        @size-change='newwatertgetpurcon'>
                        </el-pagination>
                    </div>
                </div> 
            </div>
        </div>
        <div class="rece-right"><span>《</span></div>
        <!-- 发票tab列表 -->
        <div class="rece">
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
            <!-- 进项发票 -->
            <div class="adBox" v-if="sel === 1">  
                <el-button class="adBoxbutton" type="primary" plain size="small" @click="openPur">关 联</el-button> 
                <span class="Totalsettlement" v-if="contmoeny != 0">含税金额合计: {{contmoeny|currency}} 元</span>
                <span class="Totalsettlement" v-else>含税金额合计: 0.00 元</span>
                <div class="Advanced">
                    <el-date-picker
                    size="small"
                    @change="Time"
                    v-model="TimeArray"
                    type="daterange"
                    unlink-panels
                    range-separator="~"
                    start-placeholder="开票开始日期"
                    end-placeholder="开票结束日期"
                    value-format="yyyy-MM-dd"
                    style="width:30%;float:left;margin-right:2%;">
                    </el-date-picker>
                    <el-date-picker 
                    @change="ClickOppositeName"
                    size="small"
                    v-model="SearchoppositeName"
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
                <!-- 进项发票表格 -->
                <div class="tab-table">
                    <el-table
                    v-loading="loading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    @selection-change="selectionpurChange"
                    :data="tableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                        <el-table-column type="selection"> </el-table-column>
                        <el-table-column prop="companyName" label="销方企业名称" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="month" label="所属月份" :show-overflow-tooltip="true" width="110">
                        </el-table-column>
                        <el-table-column prop="openDate" label="开票日期" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="baseMoney" label="金额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.baseMoney">{{scope.row.baseMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="taxAmount" label="税额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.taxAmount">{{scope.row.taxAmount|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="includMoney" label="含税金额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.includMoney">{{scope.row.includMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <span v-if="contmoenyone != 0">选中合计: {{contmoenyone|currency}} 元</span>
                    <span v-else>选中合计: 0.00 元</span>
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
            <!-- 销项发票 -->
            <div class="adBox" v-if="sel === 2">  
                <el-button class="adBoxbutton" type="primary" plain size="small" @click="openPurtwo">关 联</el-button>
                <span class="Totalsettlement" v-if="watertcontmoeny != 0">含税金额合计: {{watertcontmoeny|currency}} 元</span>
                <span class="Totalsettlement" v-else>含税金额合计: 0.00 元</span>
                <div class="Advanced">
                    <!-- <el-input
                    @keyup.enter.native="waterClickOppositeName"
                    size="small"
                    :placeholder="'搜索' + Seller + '方企业名称'"
                    prefix-icon="el-icon-user"
                    v-model="waterSearchoppositeName"
                    >
                    </el-input> -->
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
                <!-- 销项发票表格 -->
                <div class="tab-table">
                    <el-table
                    v-loading="waterloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    @selection-change="selectionpurChangeone"
                    :data="watertableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="overflow:visible !important;">
                    <el-table-column type="selection"> </el-table-column>
                        <el-table-column prop="companyName" label="购方企业名称" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="openDate" label="开票日期" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="baseMoney" label="金额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.baseMoney">{{scope.row.baseMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="taxAmount" label="税额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.taxAmount">{{scope.row.taxAmount|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="includMoney" label="含税金额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.includMoney">{{scope.row.includMoney|currency}} 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <span v-if="watertcount != 0">选中合计: {{watertcount|currency}} 元</span>
                    <span v-else>选中合计: 0.00 元</span>
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
            // 关联数据
            selData:{
                sments:[],
                bills:[],
            },                                 
            // 结算tab
            leftsel:1,                                      // 选项卡默认显示tab
            lefttabs: [
                { label: "采购结算", id: 1 ,width:"width:64px;"},
                { label: "销售结算", id: 2 ,width:"width:64px;"},
            ],
            // 发票tab
            sel:1,                                          // 选项卡默认显示tab
            tabs: [
                { label: "进项发票", id: 1 ,width:"width:64px;"},
                { label: "销项发票", id: 2 ,width:"width:64px;"},
            ],
            http:'',
            show:false,
            purTitle: "",                                   // 标题
            logo: "",                                       // 抽屉标题

            /* 
            采购结算数据列表
            */
            settlementshow:true,                            // 无数据隐藏分页
            settlementloading:false,
            settlementcurrent:1,                            // 当前页数
            settlementsize:50,                              // 显示条数 
            settlementtotal:0,                              // 总条数
            newwatertcontmoenyone:0,                        // 选中合计  
            newwatertcountone:0,                            // 结算总价合计金额  
            newwaterSearchoppositeNameone:'',               // 供应商搜索
            keyupTime:'',                                   // 时间搜索
            jshxpurSearch:"",                               // 全部搜索
            settlementtableData: [],                        // 表格数据
            multiplesettlement:[],                          // 选中单条数据

            /* 
            销售结算数据列表
            */
            newwaterpurSearch:'',                           // 全部搜索
            newwaterkeyupTime: '',                          // 时间搜索
            newwaterSearchoppositeName:'',                  // 客户搜索
            newwatertcontmoeny:0,                           // 选中合计  
            newwatertcount:0,                               // 结算总价合计金额  
            newwaterloading:false,                          // 加载状态
            newwatertcurrent:1,                             // 当前页数
            newwatertsize:50,                               // 显示条数
            newwaterttotal:0,                               // 总条数
            newwatertableData: [],                          // 表格数据

            /* 
            进项发票数据列表
            */
            banks:'',                           // 银行搜索    
            purSearch:'',                       // 全部搜索
            SearchoppositeName:'',              // 所属月份搜索
            TimeArray:[],                       // 进项时间段搜索
            contmoeny:0,                        // 含税金额合计 
            contmoenyone:0,                     // 选中合计   
            loading:false,                      // 加载状态
            current:1,                          // 当前页数
            size:50,                            // 显示条数
            total:0,                            // 总条数
            tableData: [],                      // 表格数据

            /* 
            销项发票数据列表
            */
            waterpurSearch:'',                       // 全部搜索
            waterkeyupTime: [],                      // 时间段搜索
            // waterSearchoppositeName:'',           // 对方户名搜索
            watertcontmoeny:0,                       // 含税金额合计
            watertcount:0,                           // 选中合计  
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
        this.settlementgetpurcon();                             // 采购结算信息
        this.amountcount();                                     // 采购结算总价合计
        this.getpurcon();                                       // 进项发票信息
        this.numcontmoeny();                                    // 进项含税金额合计
    },
    methods: {
        // 关联采购 进项
        openPur(){
            if(this.selData.sments.length > 0 && this.selData.bills.length > 0){
                if(Math.abs(this.newwatertcontmoenyone).toFixed(2) == Math.abs(this.contmoenyone).toFixed(2)){
                    this.api.capital.relSmentsBills(this.selData)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                message: '关联成功',
                                type: 'success'
                            });
                            this.settlementtableData = [];
                            this.settlementgetpurcon();                 // 采购结算信息
                            this.amountcount();                         // 采购结算总价合计
                            this.tableData = [];
                            this.getpurcon();                           // 进项发票信息
                            this.numcontmoeny();                        // 进项含税金额合计
                        }
                    })
                }else if(Math.abs(this.newwatertcontmoenyone).toFixed(2) != Math.abs(this.contmoenyone).toFixed(2)){
                    this.$message({
                        message: '采购结算总价和进项发票含税金额不一致,请检查!',
                        type: 'warning'
                    });
                }
            }else{
                this.$message({
                    message: '至少选中一条采购结算数据和进项发票数据,请检查!',
                    type: 'warning'
                });
            }
        },
        // 关联销售 销项
        openPurtwo(){
            if(this.selData.sments.length > 0 && this.selData.bills.length > 0){
                if(Math.abs(this.newwatertcontmoeny).toFixed(2) == Math.abs(this.watertcount).toFixed(2)){
                    this.api.capital.relSmentsBills(this.selData)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                message: '关联成功',
                                type: 'success'
                            });
                            this.newwatertableData = [];
                            this.newwatertgetpurcon()                   // 销售结算列表
                            this.amountcountone();                      // 销售结算总价
                            this.watertableData = [];
                            this.watertgetpurcon();                     // 销项发票列表
                            this.numcontmoenyone();                     // 销项含税金额合计
                        }
                    })
                }else if(Math.abs(this.newwatertcontmoeny).toFixed(2) != Math.abs(this.watertcount).toFixed(2)){
                    this.$message({
                        message: '销售结算总价和发票含税金额不一致,请检查!',
                        type: 'warning'
                    });
                }
            }else{
                this.$message({
                    message: '至少选中一条销售结算数据和销项发票数据,请检查!',
                    type: 'warning'
                });
            }
        },
        // 结算选项卡
        leftselect(item) {
            this.newwaterkeyupTime = "";
            this.newwaterSearchoppositeName = "";
            this.newwaterpurSearch = "";
            this.newwaterSearchoppositeNameone = "";
            this.jshxpurSearch = "";
            this.keyupTime = "";
            this.leftsel = item.id;
            if(this.leftsel == 1){
                this.sel = 1;
                this.settlementtableData = [];
                this.settlementgetpurcon();                                 // 采购结算列表
                this.amountcount();                                         // 采购结算总价
                this.tableData = [];
                this.getpurcon();                                           // 进项发票列表
                this.numcontmoeny();                                        // 进项含税金额合计
            }else if(this.leftsel == 2){
                this.sel = 2;
                this.newwatertableData = [];
                this.newwatertgetpurcon()                                   // 销售结算列表
                this.amountcountone();                                      // 销售结算总价
                this.watertableData = [];
                this.watertgetpurcon();                                     // 销项发票列表
                this.numcontmoenyone();                                     // 销项含税金额合计
            }
        },
        // 票据选项卡
        select(item) {
            this.waterkeyupTime = null;
            this.waterpurSearch = "";
            this.TimeArray = null;
            this.purSearch = "";
            this.SearchoppositeName = "";
            this.sel = item.id;
            if(this.sel == 1){
                this.leftsel = 1;
                this.settlementtableData = [];
                this.settlementgetpurcon();                                 // 采购结算列表
                this.tableData = [];
                this.getpurcon();                                           // 进项发票列表
                this.amountcount();                                         // 采购结算总价
                this.numcontmoeny();                                        // 进项含税金额合计
            }else if(this.sel == 2){
                this.leftsel = 2;
                this.newwatertableData = [];
                this.newwatertgetpurcon()                                   // 销售结算列表
                this.amountcountone();                                      // 销售结算总价
                this.watertableData = [];
                this.watertgetpurcon();                                     // 销项发票列表
                this.numcontmoenyone();                                     // 销项含税金额合计
            }
        },

        /*
        采购结算信息列表展示功能
        */
        // 采购结算表格数据
        settlementgetpurcon(){
            this.settlementloading = true;
            // 传参
            var listByPageData = {
                type:1,
                relBillFlag:false,
                supplierName:this.newwaterSearchoppositeNameone,
                searchStr:this.jshxpurSearch,
                month:this.keyupTime,
                page:{
                    current:this.settlementcurrent,
                    size:this.settlementsize
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageData)
            .then(res=>{
                this.settlementloading = false; 
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.settlementtableData = res.data.data.records;
                        this.settlementtotal = res.data.data.total;
                    }else{
                        this.settlementshow = false;
                    }
                }
            })
        },
       // 采购结算总价合计
        amountcount(){
            // this.newwatertcountone = 0;
            this.api.capital.amount({relBillFlag:false,type:1})
            .then(res=>{
                if(res.data.code == 200){
                    this.newwatertcountone = res.data.data.sumMoney;    
                } 
            })
        },
        // 采购结算选中数据
        purChangeone(val){
            let arr = [];
            this.newwatertcontmoenyone = 0;
            for(let i=0;i<val.length;i++){  
                arr.push({id:val[i].id,type:1});  
                this.newwatertcontmoenyone += val[i].sumPrice;
            }
            this.selData.sments = arr;
        },
        // 搜索供应商
        newwaterClickOppositeNameone() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
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
        销售结算列表展示功能
        */
        // 表格数据
        newwatertgetpurcon(){
            this.newwaterloading = true;
            // 传参
            var listByPageData = {
                type:2,
                relBillFlag:false,
                month:this.newwaterkeyupTime,
                supplierName:this.newwaterSearchoppositeName,
                searchStr:this.newwaterpurSearch,
                page:{
                current:this.newwatertcurrent,
                size:this.newwatertsize
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageData)
            .then(res=>{
                this.newwaterloading = false;
                if(res.data.code == 200){
                    this.newwatertableData = res.data.data.records;
                    this.newwaterttotal = res.data.data.total;
                }
            })
        },
       // 销售结算总价合计
        amountcountone(){
            this.api.capital.amount({relBillFlag:false,type:2})
            .then(res=>{
                if(res.data.code == 200){
                    this.newwatertcount = res.data.data.sumMoney;       // 销售结算总价
                }
            })
        },
        // 销售结算选中数据
        purChange(val){
            let arr = [];
            this.newwatertcontmoeny = 0;
            for(let i=0;i<val.length;i++){  
                arr.push({id:val[i].id,type:2});    
                this.newwatertcontmoeny += val[i].sumPrice;
            }
            this.selData.sments = arr;
        },
        // 日期搜索
        newwaterDataTime() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },
        // 客户搜索
        newwaterClickOppositeName() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },
        // 全部搜索
        newwatersearch() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },

        /*
        进项发票列表展示功能
        */
        // 表格数据
        getpurcon(){
            this.loading = true;
            if(this.TimeArray == null){
                // 传参
                var listByPageData = {
                    type:1,
                    startTime:"",
                    endTime:"",
                    relationFlag:false,
                    month:this.SearchoppositeName,
                    searchStr:this.purSearch,
                    page:{
                        current:this.current,
                        size:this.size
                    }
                }
            }else{
                // 传参
                var listByPageData = {
                    type:1,
                    startTime:this.TimeArray[0],
                    endTime:this.TimeArray[1],
                    relationFlag:false,
                    month:this.SearchoppositeName,
                    searchStr:this.purSearch,
                    page:{
                        current:this.current,
                        size:this.size
                    }
                }
            }
            // 渲染表格
            this.api.bill.all(listByPageData)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.tableData = res.data.data.records;
                    this.total = res.data.data.total;
                }
            })
        },
        // 进项含税金额合计
        numcontmoeny(){
            this.api.bill.amount({relationFlag:false,type:1})
            .then(res=>{
                if(res.data.code == 200){
                    this.contmoeny = res.data.data.sumMoney;       // 进项含税金额总价
                }
            })
        },
        // 进项选中数据列表
        selectionpurChange(val){
            let arr = [];
            this.contmoenyone = 0;
            for(let i=0;i<val.length;i++){   
                arr.push({id:val[i].id,type:1});  
                this.contmoenyone += val[i].includMoney;
            }
            this.selData.bills = arr;
        },
        // 供应商搜索
        ClickOppositeName() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 结算时间搜索
        Time() {
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
        销项发票列表展示功能
        */
        // 表格数据
        watertgetpurcon(){
            this.waterloading = true;
            if(this.waterkeyupTime == null){
                // 传参
                var listByPageData = {
                    type:2,
                    startTime:"",
                    endTime:"",
                    relationFlag:false,
                    searchStr:this.waterpurSearch,
                    page:{
                        current:this.watertcurrent,
                        size:this.watertsize
                    }
                }
            }else{
                // 传参
                var listByPageData = {
                    type:2,
                    startTime:this.waterkeyupTime[0],
                    endTime:this.waterkeyupTime[1],
                    relationFlag:false,
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
        },
        // 销项含税金额合计
        numcontmoenyone(){
            this.api.bill.amount({relationFlag:false,type:2})
            .then(res=>{
                if(res.data.code == 200){
                    this.watertcontmoeny = res.data.data.sumMoney;       // 销项含税金额合计
                }
            })
        },
        // 销项选中数据列表
        selectionpurChangeone(val){
            let arr = [];
            this.watertcount = 0;
            for(let i=0;i<val.length;i++){   
                arr.push({id:val[i].id,type:2});   
                this.watertcount += val[i].includMoney;
            }
            this.selData.bills = arr;
        },
        // 日期搜索
        waterDataTime() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // // 企业名称搜索
        // waterClickOppositeName() {
        //     this.watertableData = [];
        //     this.watertgetpurcon();
        //     this.watertcurrent = 1;
        // },
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
    .Totalsettlement{
        position: absolute;
        top: -47px;
        right: 32px;
        font-size: 14px;
        color: #333;
        // font-weight: 600;
    }
    .adBoxbutton{
        position: absolute;
        top: -125px;
        right: 32px;   
    }
    .tab-table{
        height: 76%;
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
            span{
                font-size: 14px;
                font-weight: 600;
                color: #333;
                line-height:33px;
            }
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