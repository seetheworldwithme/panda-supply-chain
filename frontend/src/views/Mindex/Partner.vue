<template>
    <div>
        <!-- tab选项卡 -->
        <div class="menu-list">
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
        <div class="parNavs">
            <div v-show="sel === 1">
                <!-- 供应商档案 -->
                <div class="parNavstrip">
                    <el-col :span="16" class="par-right">
                        <el-button size="small" @click="deletPart" v-if="$route.meta.delete">删除</el-button>
                        <el-button size="small" style="margin-left: 24px" @click="addPart" v-if="$route.meta.add">新建</el-button>
                        <el-input
                        @keyup.enter.native="search"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="purSearch"
                        >
                        </el-input>
                    </el-col>
                    <!-- 供应商档案表格 -->
                    <div class="pur-nav"
                    v-loading="loading"
                    element-loading-text = "加载中...">
                        <div class="pur-table">
                            <el-table
                            height="100%"
                            border
                            :data="tableDatass"
                            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                            style="width: 100%;"
                            @selection-change="purChange">
                            <el-table-column type="selection" fixed="left"></el-table-column>
                            <el-table-column prop="id" v-if="show">
                            </el-table-column>
                            <el-table-column v-if="$route.meta.edit"
                            label="操作"
                            width="100" fixed="left">
                            <template slot-scope="scope">
                                <el-button @click="EditPart(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                            </template>
                            </el-table-column>
                            <el-table-column label="企业名称" :show-overflow-tooltip="true" width="260">
                                <template slot-scope="scope">
                                    <span v-if="$route.meta.check"
                                        style="color: #cd1724;cursor: pointer;"
                                        @click="purSee(scope.$index, scope.row)"
                                        >{{ scope.row.enterpriseName }}</span
                                    >
                                    <span v-else>{{scope.row.enterpriseName}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="企业性质" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.enterpriseNature == 1">国企</span>
                                    <span v-else-if="scope.row.enterpriseNature == 2">民营企业</span>
                                    <span v-else-if="scope.row.enterpriseNature == 3">央企</span>
                                    <span v-else-if="scope.row.enterpriseNature == 4">外企</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="统一信用代码" :show-overflow-tooltip="true" width="260">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.creditCode">{{scope.row.creditCode}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="经营地址" :show-overflow-tooltip="true" width="180">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.businessAddress">{{scope.row.businessAddress}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="煤矿" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.coalMine">{{scope.row.coalMine}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="年产量" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.annualOutput">{{scope.row.annualOutput}} 吨</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="煤场" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.coalYard">{{scope.row.coalYard}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品大类" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.pro1">{{scope.row.pro1|protype}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品小类" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.pro2">{{scope.row.pro2|promiddletype}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品名称" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                <span v-if="scope.row.pro2 == 20">{{scope.row.proRemark}}</span>
                                <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="标的指标" width="260" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <div v-if="!scope.row.sulfur && !scope.row.volatilization && !scope.row.water">---</div>
                                    <div v-else>
                                        <div style="display: inline-block;"><span v-if="scope.row.sulfur">硫分St.d% : {{scope.row.sulfur}}</span><span v-else></span></div>
                                        <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.volatilization">挥发分Vdaf% : {{scope.row.volatilization}}</span><span v-else></span></div>
                                        <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.water">水分Mt% : {{scope.row.water}}</span><span v-else></span></div>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column label="创建时间" sortable :show-overflow-tooltip="true" width="160">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.createTime">{{scope.row.createTime}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <!-- 暂无数据 -->
                            <div slot="empty" class="dataAvailable" v-if="dataAva">
                                <div class="dataAva">
                                    <img src="../../assets/zanwushuju.png" alt="">
                                    <p>暂无数据，您可以新建合同试试</p>
                                    <el-button size="small" style="margin-left: 24px" @click="addPart" v-if="$route.meta.add">新建</el-button>
                                </div>
                            </div>
                            </el-table>
                        </div>
                        <!-- 分页 -->
                        <div class="pur-pages">
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
            </div>
            <div v-show="sel === 2">
                <!-- 客户档案 -->
                <div class="parNavstrip">
                    <el-col :span="16" class="par-right">
                        <el-button size="small" @click="deletcustomer" v-if="$route.meta.delete">删除</el-button>
                        <el-button size="small" style="margin-left: 24px" @click="addcustomer" v-if="$route.meta.add">新建</el-button>
                        <el-input
                        @keyup.enter.native="searchs"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="purSearchtwo"
                        >
                        </el-input>
                    </el-col>
                    <!-- 客户档案表格 -->
                    <div class="pur-nav"
                    v-loading="loadingone"
                    element-loading-text = "加载中...">
                        <div class="pur-table">
                            <el-table
                            height="100%"
                            border
                            :data="tableData"
                            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                            style="width:100%;"
                            @selection-change="customerChange">
                            <el-table-column type="selection"></el-table-column>
                            <el-table-column prop="id" v-if="show">
                            </el-table-column>
                            <el-table-column v-if="$route.meta.edit"
                            label="操作"
                            width="100">
                            <template slot-scope="scope">
                                <el-button @click="Editcustomer(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
                            </template>
                            </el-table-column>
                            <el-table-column label="企业名称" :show-overflow-tooltip="true" width="260">
                                <template slot-scope="scope">
                                    <span v-if="$route.meta.check"
                                        style="color: #cd1724;cursor: pointer;"
                                        @click="customerpurSee(scope.$index, scope.row)"
                                        >{{ scope.row.enterpriseName }}</span
                                    >
                                    <span v-else>{{scope.row.enterpriseName}}</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="企业性质" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.enterpriseNature == 1">国企</span>
                                    <span v-else-if="scope.row.enterpriseNature == 2">民营企业</span>
                                    <span v-else-if="scope.row.enterpriseNature == 3">央企</span>
                                    <span v-else-if="scope.row.enterpriseNature == 4">外企</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="统一信用代码" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.creditCode">{{scope.row.creditCode}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="经营地址" :show-overflow-tooltip="true" width="160">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.businessAddress">{{scope.row.businessAddress}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品大类" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.pro1">{{scope.row.pro1|protype}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品小类" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.pro2">{{scope.row.pro2|promiddletype}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="产品名称" width="120" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                <span v-if="scope.row.pro2 == 20">{{scope.row.proRemark}}</span>
                                <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="标的指标" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <div v-if="!scope.row.sulfur && !scope.row.volatilization && !scope.row.water">---</div>
                                    <div v-else>
                                        <div style="display: inline-block;"><span v-if="scope.row.sulfur">硫分St.d% : {{scope.row.sulfur}}</span><span v-else></span></div>
                                        <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.volatilization">挥发分Vdaf% : {{scope.row.volatilization}}</span><span v-else></span></div>
                                        <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.water">水分Mt% : {{scope.row.water}}</span><span v-else></span></div>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column prop="createtime" label="创建时间" sortable :show-overflow-tooltip="true" width="160">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.createTime">{{scope.row.createTime}}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <!-- 暂无数据 -->
                            <div slot="empty" class="dataAvailable" v-if="dataAvaone">
                                <div class="dataAva">
                                    <img src="../../assets/zanwushuju.png" alt="">
                                    <p>暂无数据，您可以新建合同试试</p>
                                    <el-button size="small" style="margin-left: 24px" v-if="$route.meta.add">新建</el-button>
                                    <!-- <el-button size="small" v-if="$route.meta.import">导入</el-button> -->
                                </div>
                            </div>
                            </el-table>
                        </div>
                        <!-- 分页 -->
                        <div class="pur-pages">
                            <el-pagination
                            background
                            :page-sizes="[10, 30, 50, 100, 200]"
                            :page-size.sync="sizetype"
                            :pager-count="5" 
                            :current-page.sync="currenttype"
                            layout="total, sizes, prev, pager, next, jumper"
                            :total="totaltype"
                            @current-change='getpurcontype'
                            @size-change='getpurcontype'>
                            </el-pagination>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 供应商档案抽屉 查看功能开始 -->
        <el-drawer :title="logo" :visible.sync="drawerss" size="700px">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
            <el-tabs v-model="activeName" @tab-click="handleClick">
                <el-tab-pane label="基本信息" name="first">
                <el-row>
                    <el-col :span="24">
                    <div class="see-top">基础信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>企业名称：</label>
                            <span>{{seeDatass.enterpriseName}}</span>
                        </li>
                        <li>
                            <label>企业性质：</label>
                            <span v-if="seeDatass.enterpriseNature">{{seeDatass.enterpriseNature}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>统一信用代码：</label>
                            <span v-if="seeDatass.creditCode">{{seeDatass.creditCode}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>经营地址：</label>
                            <span v-if="seeDatass.businessAddress">{{seeDatass.businessAddress}}</span>
                            <span v-else>---</span>
                        </li>
                    </ul>
                    </el-col>
                    <el-col :span="24">
                        <div class="see-top">生产经营信息</div>
                        <ul class="see-lable">
                            <li>
                                <label>煤矿：</label>
                                <span v-if="seeDatass.coalMine">{{seeDatass.coalMine}}</span>
                                <span v-else>---</span>
                            </li>
                            <li>
                                <label>年产量：</label>
                                <span v-if="seeDatass.annualOutput">{{seeDatass.annualOutput}}</span>
                                <span v-else>---</span>
                            </li>
                            <li>
                                <label>煤场：</label>
                                <span v-if="seeDatass.coalYard">{{seeDatass.coalYard}}</span>
                                <span v-else>---</span>
                            </li>
                            <li>
                                <label>产品大类：</label>
                                <span v-if="seeDatass.pro1">{{seeDatass.pro1}}</span>
                                <span v-else>---</span>
                            </li>
                            <li>
                                <label>产品小类：</label>
                                <span v-if="seeDatass.pro2">{{seeDatass.pro2}}</span>
                                <span v-else>---</span>
                            </li>
                            <li>
                                <label>标的指标：</label>
                                <div class="see-libox">
                                    <span>硫分St.d% : </span>
                                    <div>
                                        <p v-if="seeDatass.sulfur">{{seeDatass.sulfur}} </p><p v-else>--- </p>
                                    </div>
                                    <span>挥发分Vdaf% : </span>
                                    <div>
                                        <span v-if="seeDatass.volatilization">{{seeDatass.volatilization}} </span><span v-else>--- </span>
                                    </div>
                                    <span>水分Mt% : </span>
                                    <div>
                                        <span v-if="seeDatass.water">{{seeDatass.water}} </span><span v-else>--- </span>
                                    </div>
                              
                                </div>
                            </li>
                        </ul>
                    </el-col>
                </el-row>
                </el-tab-pane>
                <el-tab-pane label="动态" name="second">
                    <el-row>
                        <el-col :span="24">
                            <div class="see-top" style="margin-bottom:24px;">动态信息</div>
                            <el-timeline>
                                <el-timeline-item v-for="v,i in timeData" :key="i" :timestamp="v.createtime" placement="top">
                                    <el-card>
                                        <h4>{{v.operation}}</h4>
                                    </el-card>
                                </el-timeline-item>
                            </el-timeline>
                        </el-col>
                </el-row>
                </el-tab-pane>
                <!-- 团队模块 V.0.0.2 上线  -->
                <!-- <el-tab-pane label="团队" name="third">团队</el-tab-pane> -->
            </el-tabs>
            <div class="xia">《</div>
            </div>
        </el-drawer>
        <!-- 供应商档案抽屉 查看功能结束 -->

        <!-- 客户档案抽屉 查看功能开始 -->
        <el-drawer :title="logo" :visible.sync="customerdrawer" size="700px">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
            <el-tabs v-model="customeractiveName" @tab-click="customerhandle">
                <el-tab-pane label="基本信息" name="first">
                <el-row>
                    <el-col :span="24">
                    <div class="see-top">基础信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>企业名称：</label>
                            <span>{{customerData.enterpriseName}}</span>
                        </li>
                        <li>
                            <label>企业性质：</label>
                            <span v-if="customerData.enterpriseNature">{{customerData.enterpriseNature}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>统一信用代码：</label>
                            <span v-if="customerData.creditCode">{{customerData.creditCode}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>经营地址：</label>
                            <span v-if="customerData.businessAddress">{{customerData.businessAddress}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>产品大类：</label>
                            <span v-if="customerData.pro1">{{customerData.pro1}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>产品小类：</label>
                            <span v-if="customerData.pro2">{{customerData.pro2}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>标的指标：</label>
                            <div class="see-libox">
                                <span>硫分St.d% : </span>
                                <div>
                                    <p v-if="customerData.sulfur">{{customerData.sulfur}} </p><p v-else>--- </p>
                                </div>
                                <span>挥发分Vdaf% : </span>
                                <div>
                                    <span v-if="customerData.volatilization">{{customerData.volatilization}} </span><span v-else>--- </span>
                                </div>
                                <span>水分Mt% : </span>
                                <div>
                                    <span v-if="customerData.water">{{customerData.water}} </span><span v-else>--- </span>
                                </div>
                            
                            </div>
                        </li>
                    </ul>
                    </el-col>
                </el-row>
                </el-tab-pane>
                <el-tab-pane label="动态" name="second">
                    <el-row>
                        <el-col :span="24">
                            <div class="see-top" style="margin-bottom:24px;">动态信息</div>
                            <el-timeline>
                                <el-timeline-item v-for="v,i in customertimeData" :key="i" :timestamp="v.createtime" placement="top">
                                    <el-card>
                                        <h4>{{v.operation}}</h4>
                                    </el-card>
                                </el-timeline-item>
                            </el-timeline>
                        </el-col>
                </el-row>
                </el-tab-pane>
                <!-- 团队模块 V.0.0.2 上线  -->
                <!-- <el-tab-pane label="团队" name="third">团队</el-tab-pane> -->
            </el-tabs>
            <div class="xia">《</div>
            </div>
        </el-drawer>
        <!-- 客户档案抽屉 查看功能结束 -->
    </div>
</template>
<script>
export default {
    data(){
        return {
            sel: '',   // 选项卡默认显示tab
            tabs: [
                { label: "供应商档案", id: 1 ,width:"width:80px;"},
                { label: "客户档案", id: 2 ,width:"width:64px;"},
            ],
            show:false,  //隐藏表格id
            logo:'平遥星皓洗煤厂', //标题

            /*
            供应商档案数据
            */
            dataAva:false,
            loading:false,
            current:1,                          // 当前页数
            size:10,                            // 显示条数
            total:0,                            // 总条数
            delmul:[],                          // 批量删除
            activeName: "first", // 供应商档案查看tab栏
            purSearch:'', // 供应商档案搜索
            drawerss: false, // 供应商档案查看抽屉
            multiplePart:[],  //供应商档案选中表格数据
            tableDatass: [],        // 供应商档案表格数据
            seeDatass:{
                enterpriseName:'',
                enterpriseNature:'',
                creditCode:'',
                businessAddress:'',
                coalMine:'',
                annualOutput:'',
                coalYard:'',
                pro1: "",
                pro2:"",
                sulfur:"",
                volatilization:"",
                water:"",
            },  // 供应商档案查看数据
            timeData:[
                // {
                //     createtime:'2021-01-02',
                //     operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ],  // 供应商档案 查看 动态

            /*
            客户档案数据
            */
            dataAvaone:false,
            loadingone:false,
            currenttype:1,                          // 当前页数
            sizetype:10,                            // 显示条数
            totaltype:0,                            // 总条数
            delmultype:[],                          // 批量删除
            customeractiveName: "first", // 客户档案查看tab栏
            purSearchtwo:'', // 客户档案搜索
            customerdrawer: false, // 客户档案查看抽屉
            customerPart:[],  //客户档案选中表格数据
            tableData: [],    // 客户档案表格数据
            customerData:{
                enterpriseName:'',
                enterpriseNature:'',
                creditCode:'',
                businessAddress:'',
                pro1: "",
                pro2:"",
                sulfur:"",
                volatilization:"",
                water:"",
            },  // 客户档案查看数据
            customertimeData:[
                // {
                //     createtime:'2021-01-02',
                //     operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ],  // 客户档案 查看 动态
        }
    },
    mounted(){
        
        if(this.$route.query.sel == 1){
            this.sel = 1;
            this.getpurcon(); // 渲染供应商表格
        }else if(this.$route.query.sel == 2){
            this.sel = 2;
            this.getpurcontype();
        }else{
            this.sel = 1;
            this.getpurcon(); // 渲染供应商表格
        }
    },
    methods: {
        //选项卡
        select(item) {
            this.purSearch = "";
            this.purSearchtwo = "";
            this.sel = item.id;
            if(this.sel == 1){
                this.getpurcon(); // 渲染供应商表格
            }else if(this.sel == 2){
                this.getpurcontype(); // 渲染客户表格
            }
        },

        /*
        供应商档案功能方法
        */
      // 供应商表格数据
      getpurcon(){
        this.loading = true;
        // 传参
        var listByPageData = {
          partnerType:1,
          searchStr:this.purSearch,
          page:{
            current:this.current,
            size:this.size
          }
        }
        // 渲染表格
        this.api.supplier.all(listByPageData)
        .then(res=>{
            this.loading = false;
            if(res.data.code == 200){
                if(res.data.data.records != ""){
                    this.tableDatass = res.data.data.records;
                    this.total = res.data.data.total;
                }else{
                    this.dataAva = true;
                }
            }
        })
      },
        //供应商档案搜索功能
        search() {
            this.tableDatass = [];
            this.current = 1;
            this.getpurcon();
        },
        purChange(val) {
            // 供应商档案表格选中数据
            this.multiplePart = val;
            this.delmul = [];
            for(let i=0;i<this.multiplePart.length;i++){
                this.delmul.push(this.multiplePart[i].id) 
            }
        },
        // 供应商档案查看
        purSee(index, row) {
            this.drawerss = true;
            this.api.supplier.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                    this.seeDatass = res.data.data;
                }
            })
        },
        // 供应商档案抽屉...动态
        handleClick(tab) {
            // console.log(tab.label);
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "动态") {
                // console.log(2);
            } else if (tab.label == "团队") {
                // console.log(3);
            }
        },
        //新建供应商档案
        addPart(){
            this.$router.push({name:"Addsupplier"})
        },
        // 编辑供应商档案
        EditPart(row){
            this.api.supplier.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                let purdata = res.data.data;
                purdata.id = row.id;
                sessionStorage.setItem('queryAll',JSON.stringify(purdata));
                this.$router.push({name:"Editsupplier"})
                }
            })
        },
        //供应商档案 删除功能
        deletPart() {
            if(this.multiplePart.length == ""){
                this.$message({
                    type: 'warning',
                    message: '请选择要删除的数据!'
                });
            }else{
                this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.api.supplier.del(this.delmul)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                            type: 'success',
                            message: '删除成功!'
                            });
                            this.tableDatass = [];
                            this.current = 1;
                            this.getpurcon();
                            if(this.tableDatass == []){
                                this.dataAva = true;
                            }
                        }
                    })
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });          
                });
            }
        },

        /*
        客户档案功能方法
        */
        // 客户表格数据
        getpurcontype(){
            this.loadingone = true;
            // 传参
            var listByPageData = {
                partnerType:2,
                searchStr:this.purSearchtwo,
                page:{
                    current:this.currenttype,
                    size:this.sizetype
                }
            }
            // 渲染表格
            this.api.supplier.all(listByPageData)
            .then(res=>{
                this.loadingone = false;
            if(res.data.code == 200){
                if(res.data.data.records != ""){
                    this.tableData = res.data.data.records
                    this.totaltype = res.data.data.total
                }else{
                    this.dataAvaone = true;
                }

            }
            })
        },
        // 客户档案搜索功能
        searchs(){
            this.tableData = [];
            this.getpurcontype();
            this.currenttype = 1
        },
        customerChange(val) {
            // 供应商档案表格选中数据
            this.customerPart = val;
            this.delmultype = [];
            for(let i=0;i<this.customerPart.length;i++){
                this.delmultype.push(this.customerPart[i].id) 
            }
        },
        // 客户档案查看
        customerpurSee(index, row) {
            this.customerdrawer = true;
            this.api.supplier.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                    this.customerData = res.data.data;
                }
            })
        },
        // 客户档案抽屉
        customerhandle(tab) {
            // console.log(tab.label);
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "动态") {
                // console.log(2);
            } else if (tab.label == "团队") {
                // console.log(3);
            }
        },
        //新建客户档案
        addcustomer(){
            this.$router.push({name:"Addcustomer"})
        },
        // 编辑客户档案
        Editcustomer(row){
            this.api.supplier.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                let purdata = res.data.data;
                purdata.id = row.id;
                sessionStorage.setItem('queryAll',JSON.stringify(purdata));
                this.$router.push({name:"Editcustomer"})
                }
            })
        },
        //客户档案 删除功能
        deletcustomer() {
            if(this.customerPart.length == []){
                this.$message({
                    type: 'warning',
                    message: '请选择要删除的数据!'
                });
            }else{
                this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    this.api.supplier.del(this.delmultype)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: 'success',
                                message: '删除成功!'
                            });
                            this.tableData = [];
                            this.getpurcontype();
                            this.currenttype = 1;
                            if(this.tableData == []){
                            this.dataAvaone = true;
                            }
                        }
                    })
                }).catch(() => {
                    this.$message({
                        type: 'info',
                        message: '已取消删除'
                    });          
                });
            }
        },  
    }
}
</script>
<style lang="less" scoped>

// 选项卡开始
.menu-list{
    width: 100%;
    height: 48px;
    background-color: #fff;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    border: 1px solid rgba(102, 102, 102, 0.15);
    overflow: hidden;
}
.totalTab {
    width: 10%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin:9px 0 0 35px;
}
.totalTab.active {
    color: #ea222e;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:10px auto 0;
    display: none;
}
.strip {
    display: block;
}
.pur-nav {
  width: calc(100% - 32px);
  height: calc(100vh - 142px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px 8px 8px 8px;
  .pur-table{
    width:100%;
    height: 94%;
  } 
}
  .pur-pages {
    width: 100%;
    height: 6%;
    padding: 6px 0 0 12px;
    box-sizing:border-box;
    span{
      font-size: 14px;
      line-height:33px;
      color: #333;
    }
    .el-pagination {
      float: right;
    }
  }

//选项卡结束
// 暂无数据样式开始
.dataAvailable{
width: 100%;
height: 245px;
position: relative;
margin-top:185px;
margin-bottom:279px;
.dataAva{
    position: absolute;
    top:75%;
    left: 25%;
    transform: translate(0, -50%);
    text-align: center;
    img{
    width: 371px;
    height: 200px;
    margin: auto;
    }
    p{
    font-size: 16px;
    line-height: 24px;
    color: #333;
    }
    .el-button{
    margin-bottom:100px;
    }
}
}
// 暂无数据样式结束
// 查看抽屉样式开始
.pur-tiao {
  width: 100%;
  height: 8px;
  background: #f5f5f5;
}
.pur-drawerNav {
  width: 100%;
  padding: 24px 0 24px 24px;
  height: calc(100% - 8px);
  .el-row {
    height: calc(100vh - 200px);
    overflow-y: scroll;
  }
  .xia{
    width: 30px;
    margin: auto;
    font-size: 16px;
    color: #828282;
    line-height: 18px;
    transform: rotate(270deg);
  }
}
.see-top {
  height: 40px;
  font-weight: 500;
  font-size: 14px;
  line-height: 45px;
  border-bottom: 1px dashed #999999;
}
.see-lable {
  padding: 15px 0;
  label {
      display: inline-block;
    width: 18%;
    font-size: 14px;
    color: #666;
  }
  li {
    height: 40px;
    font-size: 14px;
    line-height: 40px;
    color: #333;
    .see-libox{
      width: 76%;
      display: inline-block;
      span{
        color: #666;
      }
      div{
        display: inline-block;
        margin-right: 15px;
      }
    }
  }
}
// 查看抽屉样式结束
.parNavstrip{
    position: relative;
    .par-right{
        position: absolute;
        top: -50px;
        right: 24px;
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