<template>
    <el-container>
        <el-header>
            {{companyName}}
            <div @click="signout" style="cursor: pointer;"><i class="el-icon-alituichu"></i> 退出</div>
        </el-header>
        <el-main v-loading="loading" element-loading-text = "加载中...">
            <!-- 头部开始 -->
            <div class="pur-top">
                <span class="titleSpan">{{timedata}} </span>
                <el-button icon="el-icon-search" size="mini" @click="dialogVisibleclick">{{TimeFilter}}</el-button>
                <!-- 筛选弹窗 -->
                <el-dialog
                :visible.sync="dialogVisibledata"
                :show-close="false"
                width="39.6%">
                    <div slot="title" class="header-title">
                        <span>{{timeSizes}} </span>
                        <el-date-picker
                        v-if="TimeShowTab == 1"
                        disabled
                        style="width:36%;margin-right:15px;"
                        size="small"
                        :clearable="false"
                        unlink-panels
                        v-model="screeninvoice"
                        :picker-options="pickerOptions0"
                        type="daterange"
                        value-format="yyyy-MM-dd"
                        range-separator="至"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期">
                        </el-date-picker>
                        <el-select v-model="yearValues" placeholder="请选择" size="small" style="width:25%;margin-right:15px;" v-if="TimeShowTab == 2">
                            <el-option
                            v-for="item in yearOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value">
                            </el-option>
                        </el-select>
                        <el-select
                        v-model="screenlist"
                        multiple
                        size="small"
                        collapse-tags
                        style="width:172px;"
                        placeholder="请选择字段">
                            <el-option
                            v-for="item in screenOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value">
                            </el-option>
                        </el-select>
                    </div>
                    <div class="elTabBox">
                        <el-tabs tab-position="left" style="height:330px;" @tab-click="elTabBoxClick" v-model="tabactiveName">
                            <el-tab-pane label="全年" name="0">
                                <ul>
                                    <li style="background-color: #CF1724;color:#fff;">全年</li>    
                                </ul> 
                            </el-tab-pane>
                            <el-tab-pane label="半年" name="1">
                                <ul>
                                    <li v-for="(item,index) in vipSet" :key="index" :class="index==clickIndex?'selectItem':''" @click="selectItem(item,index)">{{item}}</li>
                                </ul> 
                            </el-tab-pane>
                            <el-tab-pane label="季度" name="2">
                                <ul>
                                    <li v-for="(item,indexone) in vipSetone" :key="indexone" :class="indexone==clickIndexone?'selectItem':''" @click="selectItemone(item,indexone)">{{item}}</li>         
                                </ul> 
                            </el-tab-pane>
                            <el-tab-pane label="月份" name="3">
                                <ul>
                                    <li v-for="(item,indextwo) in vipSettwo" :key="indextwo" :class="indextwo==clickIndextwo?'selectItem':''" @click="selectItemtwo(item,indextwo)">{{item}}</li>
                                </ul>
                            </el-tab-pane>
                            <el-tab-pane label="自定义范围" name="4">
                                <div class="customBox" @click="customBoxclick">
                                    <span>开始时间 </span><span style="color:#CF1724;">{{screeninvoice[0]}} </span><span> ~ 结束时间 </span><span style="color:#CF1724;"> {{screeninvoice[1]}}</span>
                                </div>
                            </el-tab-pane>
                        </el-tabs>
                        <!-- 自定义时间范围 -->
                        <div class="zdyTime">
                            <el-date-picker
                            ref="tabTime"
                            :clearable="false"
                            unlink-panels
                            v-model="screeninvoice"
                            :picker-options="pickerOptions0"
                            type="daterange"
                            value-format="yyyy-MM-dd"
                            range-separator="至"
                            start-placeholder="开始日期"
                            end-placeholder="结束日期">
                            </el-date-picker>
                        </div>
                    </div>
                    <span slot="footer">
                        <el-button @click="dialogVisiblerice" size="small">重 置</el-button>
                        <el-button type="primary" @click="screenchange" size="small">确 定</el-button>
                    </span>
                </el-dialog>
                <!-- 地区查询 -->
                <div class="screenBox">
                    <el-select filterable v-model="enterpriseValue" placeholder="请选择企业名称" size="mini" style="width: 26%;" @change="enterprises($event)" clearable @clear="enterprisesclear">
                        <el-option
                        v-for="item in enterpriseOptions"
                        :key="item.value"
                        :label="item.enterpriseName"
                        :value="item.id">
                        </el-option>
                    </el-select>
                    <el-select v-model="villageValue" placeholder="请选择乡" size="mini" @change="villages($event)" v-if="citythree == true" clearable @clear="villagesclear">
                        <el-option
                        v-for="item in villageOptions"
                        :key="item.value"
                        :label="item.name"
                        :value="item.code">
                        </el-option>
                    </el-select>
                    <el-select v-model="countyValue" placeholder="请选择县 / 区" size="mini" @change="countys($event)" v-if="citytwo == true" clearable @clear="countysclear">
                        <el-option
                        v-for="item in countyOptions"
                        :key="item.value"
                        :label="item.name"
                        :value="item.code">
                        </el-option>
                    </el-select>
                    <el-select v-model="cityValue" placeholder="请选择市" size="mini" @change="citys($event)" v-if="cityone == true" clearable @clear="citysclear">
                        <el-option
                        v-for="item in cityOptions"
                        :key="item.value"
                        :label="item.name"
                        :value="item.code">
                        </el-option>
                    </el-select>
                </div>
            </div>
            <!-- 头部结束 -->
            <!-- 内容开始 -->
            <div class="pur-nav">
                <!-- 企业数量 -->
                <div class="Purchasing" v-if="purShow == true">
                    <img src="../assets/qiye.png" alt="">
                    <span class="tity"> 企业数量</span>
                    <ul class="Purchasing-conent">
                        <li>
                            <span>注册企业数量</span>
                        </li>
                        <li>
                            <span>当月新增企业数量</span>
                        </li>
                        <li style="">
                            <h5 v-if="chasingTonnage">{{chasingTonnage}}</h5>
                            <h5 v-else> 0 </h5>
                        </li>
                        <li>
                            <h5 v-if="arrsumNum">{{arrsumNum}}</h5>
                            <h5 v-else> 0 </h5>
                        </li>   
                    </ul>
                </div>
                <!-- 企业客户 -->
                <div class="Purchasing" v-if="purShow == false">
                    <img src="../assets/qiye.png" alt="">
                    <span class="tity"> 企业客户</span>
                    <ul class="Purchasing-conentss">
                        <li>
                            <span>企业名称</span>
                        </li>
                        <li style="width: 35%;">
                            <span>上户时间</span>
                        </li>
                        <li>
                            <h5 v-if="enterpriseName">{{enterpriseName}}</h5>
                            <h5 v-else>  </h5>
                        </li>
                        <li style="width: 35%;">
                            <h5 v-if="regTime">{{regTime}}</h5>
                            <h5 v-else>  </h5>
                        </li>   
                    </ul>
                </div>
                <!-- 总体业务量 -->
                <div class="quantity">
                    <img src="../assets/tongji.png" alt="">
                    <span class="tity"> 总体业务量</span>
                    <ul class="quantity-conent">
                        <li>
                            <span>业务量金额累计</span>
                        </li>
                        <li>
                            <span>纳税金额累计</span>
                        </li>
                        <li>
                            <span>物流量累计</span>
                        </li>
                        <li>
                            <h5 v-if="TaxAmounts">{{TaxAmounts|newcurrency}} 亿</h5>
                            <h5 v-else> 0.00 亿</h5>
                            <span>(元)</span>
                        </li>
                        <li>
                            <h5 v-if="rateBills">{{rateBills|newcurrency}} 万</h5>
                            <h5 v-else> 0.00 万</h5>
                            <span>(元)</span>
                        </li>
                        <li>
                            <h5 v-if="logisticsBills">{{logisticsBills|newcurrency}} 万</h5>
                            <h5 v-else> 0.00 万</h5>
                            <span>(吨)</span>
                        </li>
                    </ul>
                </div>
                <!-- 数量 -->
                <div class="payment">
                    <div id="payment"></div>
                    <!-- <div class="payment-time">
                        <el-dropdown trigger="click">
                            <span>{{paymenttime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item v-for="(v,k) in paymentdropdown" :key="k" @click.native="paymentChange(v.name,v.type)">{{v.name}}</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </div> -->
                </div>
                <!-- 物流 -->
                <div class="logistics">
                    <div id="logistics"></div>
                    <!-- <div class="logistics-time">
                        <el-dropdown trigger="click">
                            <span>{{logisticstime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item v-for="(v,k) in settdropdown" :key="k" @click.native="logisticsChange(v.name,v.type)">{{v.name}}</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </div> -->
                </div>
                <!-- 产品交易份额 -->
                <div class="invoice">
                    <span>产品交易份额</span>
                    <div id="invoice"></div>
                </div>
            </div>
            <!-- 内容结束 -->
        </el-main>
    </el-container>
</template>
<script>
    var echarts = require('echarts');
    export default {
        data() {
            return {
                TimeFilter:'', 
                timeSizes:'年份',                       // 筛选标题
                TimeShowTab:2,                
                tabactiveName:'0',                      // tab默认选中
                yearValues:'',                          // 选中全年年份
                yearOptions:[],                         // 年份下拉数据
                vipSet:["上半年","下半年"],              // 半年
                clickIndex:0,
                topYear:[],                             // 选中半年时间
                topYearSize:'上半年',                   // 选中半年名称
                vipSetone:["第一季度","第二季度","第三季度","第四季度"],      // 季度
                clickIndexone:0,
                topYearone:[],                          // 选中季度时间
                topYearSizeone:'第一季度',              // 选中季度名称
                vipSettwo:["一月","二月","三月","四月","五月","六月","七月","八月","九月","十月","十一月","十二月"],      // 月份
                clickIndextwo:0,
                topYeartwo:[],                          // 选中月份时间
                topYearSizetwo:'一月',                  // 选中月份名称
                dialogVisibledata:false,
                screenlist:[1,2,3],           // 类别
                newscreenlist:[1,2,3],        // 选中类别
                screeninvoice:[],               // 时间
                newscreeninvoice:[],            // 选中时间

                purShow:true,
                cityone:true,   // 市
                citytwo:true,   // 县
                citythree:true, // 乡
                companyName:'',
                pickerOptions0:{
                    disabledDate:(time)=>{
                        return time.getTime() > Date.now() || time.getTime() == Date.now();
                    }
                },
                screenOptions: [
                    {
                        value: 1,
                        label: '产品交易份额'
                    },
                    {
                        value: 2,
                        label: '数量'
                    },
                    {
                        value: 3,
                        label: '物流'
                    }
                ],
                loading:false,
                timedata:'',
                // 地区筛选
                cityOptions: [],                // 市
                cityValue: '',
                countyOptions: [],              // 县 区
                countyValue: '',
                villageOptions: [],             // 乡
                villageValue: '',
                enterpriseOptions: [],          // 企业名称
                enterpriseValue: '',


                // 企业数量
                chasingTonnage:0,               // 注册企业数量
                arrsumNum:0,                    // 当月新增企业数量

                // 企业客户
                enterpriseName:'',
                regTime:'',

                // 总体业务量
                TaxAmounts:0,                   // 业务量金额累计
                rateBills:0,                    // 纳税金额累计
                logisticsBills:0,               // 物流量累计

                // 数量
                datapayment:[],                 // 时间
                paymentData:[],                 // 业务量
                paymentserData:[],              // 纳税量

                // 物流
                datalogistics:[],               // 时间
                logisticsData:[],               // 物流自提

                // 产品交易份额
                newColor:[],                // 颜色
                countPro:[],
                invoice:'',
                sumTaxAmount:'',            // 进项税额
                arrTaxAmount:'',            // 销项税额
            }
        },
        created() {
            this.companyName = sessionStorage.getItem('companyName');
            // 根据系统不同时间来判断，所以需要用到日期内置对象
            // 2. 得到当前的小时数
            var h = new Date().getHours();
            // 3. 判断小时数改变文字信息
            if(h < 6){
                this.timedata = '深夜了 !';
            }else if (h < 9){
                this.timedata = '早上好 !';
            }else if (h < 12){
                this.timedata = '上午好 !';
            }else if (h < 14){
                this.timedata = '中午好 !';
            }else if (h < 18){
                this.timedata = '下午好 !';
            }else {
                this.timedata = '晚上好 !';
            };
            this.http  = this.$store.state.http;                 // IP地址
        },
        beforeDestroy(){
            // 各图表实例都是局部变量,离开页面时没人释放:按 echarts 在容器上打的标记逐个销毁
            this.$el.querySelectorAll('[_echarts_instance_]').forEach(el => echarts.getInstanceByDom(el).dispose());
        },
        mounted(){
            // 行政区划
            if(sessionStorage.getItem('level') == 3){
                this.api.areaInfo.city({parentCode:sessionStorage.getItem('xzqhCode')})
                .then(res=>{
                    if(res.data.code == 200){
                        this.cityOptions = res.data.data;   // 市
                    }
                })  
            } 
            if(sessionStorage.getItem('level') == 4){
                this.cityone = false;       // 市隐藏
                this.api.areaInfo.city({parentCode:sessionStorage.getItem('xzqhCode')})
                .then(res=>{
                    if(res.data.code == 200){
                        this.countyOptions = res.data.data;   // 县
                    }
                })
            }
            if(sessionStorage.getItem('level') == 5){
                this.cityone = false;       // 市隐藏
                this.citytwo = false;       // 县隐藏
                this.api.areaInfo.city({parentCode:sessionStorage.getItem('xzqhCode')})
                .then(res=>{
                    if(res.data.code == 200){
                        this.villageOptions = res.data.data;   // 乡
                    }
                })
            }
            if(sessionStorage.getItem('level') == 6){
                this.cityone = false;       // 市隐藏
                this.citytwo = false;       // 县隐藏
                this.citythree = false;     // 乡隐藏
                this.purShow = false;       // 显示企业客户
                this.api.statistics.enterData({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
                .then(res=>{
                    if(res.data.code == 200){
                        this.enterpriseOptions = res.data.data; // 企业名称
                    }
                })
            }

            // 获取当前时间
            let date = new Date();
            let year = date.getFullYear();
            let month = date.getMonth() + 1;
            let day = date.getDate();
            if (month < 10) {
                month = "0" + month;
            }
            if (day < 10) {
                day = "0" + day;
            }
            let nowDate = year + "-" + month + "-" + day;
            this.screeninvoice = ["2019-10-01",nowDate];
            this.TimeFilter = "2019-10-01 ~ " + nowDate;
            this.newscreeninvoice = ["2019-10-01",nowDate];
            let invoiceDate = {
                xzqhCode:sessionStorage.getItem('xzqhCode'),
                startTime:'2019-10-01',
                endTime:nowDate
            };
            this.yearValues = year;
            // 获取年份下拉框数据
            let arrYear = [];
            let Years = year - 2019;
            for (var i = 0; i <= Years; i++) {
                arrYear.push(year--)
            }
            arrYear.forEach(ele => {
                this.yearOptions.push({
                    value: ele,
                    label: ele + '年'
                })
            })
            // 企业数量
            this.countCapital({xzqhCode:sessionStorage.getItem('xzqhCode')});
            // 总体业务量
            this.countBillss({xzqhCode:sessionStorage.getItem('xzqhCode')});   // 业务量金额累计
            this.rateBill({xzqhCode:sessionStorage.getItem('xzqhCode')});      // 纳税金额累计
            this.logisticsBill({xzqhCode:sessionStorage.getItem('xzqhCode')}); // 物流量累计
            // 数量
            this.paymentCapital(invoiceDate);
            // 物流
            this.logis(invoiceDate);
            // 产品交易份额
            this.countBill({xzqhCode:sessionStorage.getItem('xzqhCode')});
        },
        methods: {
            // 退出
            signout(){
                sessionStorage.clear();
                this.$router.push({name:'Login'})
                window.location.reload(); // 整页刷新:路由表和按钮权限是按上一个账号裁剪过的,同标签页换账号会继承
            },
            // 处理时间排序函数
            compare(property){
                return function(a,b){
                    var value1 = a[property];
                    var value2 = b[property];
                    return value1 - value2;
                }
            },
            // 打开弹窗
            dialogVisibleclick(){
                this.screenlist = this.newscreenlist;
                this.screeninvoice = this.newscreeninvoice;
                this.dialogVisibledata = true;
                this.topYear = [this.yearValues + '-01-01',this.yearValues + '-06-30'];
            },
            // 重置筛选
            dialogVisiblerice(){
                this.screenlist = [1,2,3];
                // 获取当前时间
                let date = new Date();
                let year = date.getFullYear();
                let month = date.getMonth() + 1;
                let day = date.getDate();
                if (month < 10) {
                    month = "0" + month;
                }
                if (day < 10) {
                    day = "0" + day;
                }
                this.tabactiveName = "4";
                this.timeSizes = "时间";
                this.TimeShowTab = 1;
                this.$refs.tabTime.focus();
                let nowDate = year + "-" + month + "-" + day;
                this.screeninvoice = ["2019-10-01",nowDate];
            },
            // 确定筛选搜索
            screenchange(){
                // 清空行政划分
                this.cityValue = '';
                this.countyValue = '';
                this.villageValue = '';
                this.enterpriseValue = '';
                if(this.tabactiveName == 0){
                    this.screeninvoice = [this.yearValues + '-01-01',this.yearValues + '-12-31'];
                    this.TimeFilter = this.yearValues + '年';
                }
                if(this.tabactiveName == 1){
                    this.screeninvoice = this.topYear;
                    this.TimeFilter = this.yearValues + '年 ~ ' + this.topYearSize;
                }
                if(this.tabactiveName == 2){
                    this.screeninvoice = this.topYearone;
                    this.TimeFilter = this.yearValues + '年 ~ ' + this.topYearSizeone;
                }
                if(this.tabactiveName == 3){
                    this.screeninvoice = this.topYeartwo;
                    this.TimeFilter = this.yearValues + '年 ~ ' + this.topYearSizetwo;
                }
                if(this.tabactiveName == 4){
                    this.TimeFilter = this.screeninvoice[0] + ' ~ ' + this.screeninvoice[1];
                }
                // 选中筛选值赋值
                this.newscreenlist = this.screenlist;
                this.newscreeninvoice = this.screeninvoice;
                let invoiceDates = {
                    xzqhCode:sessionStorage.getItem('xzqhCode'),
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({xzqhCode:sessionStorage.getItem('xzqhCode')});       // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDates);  // 数量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDates);           // 物流
                    }
                }
                this.dialogVisibledata = false;
            },
            // 筛选tab栏
            elTabBoxClick(tab){
                if(tab.index == 4){
                    this.timeSizes = "时间";
                    this.TimeShowTab = 1;
                    this.$refs.tabTime.focus();
                }else{
                    this.timeSizes = "年份";
                    this.TimeShowTab = 2;
                    if(tab.index == 3){
                        // this.screeninvoice = [];
                        this.clickIndextwo = 0;
                        this.topYeartwo = [this.yearValues + '-01-01',this.yearValues + '-01-31'];
                        this.topYearSizetwo = "一月";
                    }
                    if(tab.index == 2){
                        // this.screeninvoice = [];
                        this.clickIndexone = 0;
                        this.topYearone = [this.yearValues + '-01-01',this.yearValues + '-03-31'];
                        this.topYearSizeone = "第一季度";
                    }
                    if(tab.index == 1){
                        // this.screeninvoice = [];
                        this.clickIndex = 0;
                        this.topYear = [this.yearValues + '-01-01',this.yearValues + '-06-30'];
                        this.topYearSize = "上半年";
                    }
                    if(tab.index == 0){
                        // this.screeninvoice = [];
                        // 获取当前时间
                        let date = new Date();
                        let year = date.getFullYear();
                        this.yearValues = year;
                    }
                }
            },
            // 半年
            selectItem(itme,index){
                this.clickIndex = index;
                if(index == 0){
                    this.topYear = [this.yearValues + '-01-01',this.yearValues + '-06-30'];
                    this.topYearSize = itme;
                }
                if(index == 1){
                    this.topYear = [this.yearValues + '-07-01',this.yearValues + '-12-31'];
                    this.topYearSize = itme;
                }
            },
            // 季度
            selectItemone(itme,index){
                this.clickIndexone = index;
                if(index == 0){
                    this.topYearone = [this.yearValues + '-01-01',this.yearValues + '-03-31'];
                    this.topYearSizeone = itme;
                }
                if(index == 1){
                    this.topYearone = [this.yearValues + '-04-01',this.yearValues + '-06-30'];
                    this.topYearSizeone = itme;
                }
                if(index == 2){
                    this.topYearone = [this.yearValues + '-07-01',this.yearValues + '-09-30'];
                    this.topYearSizeone = itme;
                }
                if(index == 3){
                    this.topYearone = [this.yearValues + '-10-01',this.yearValues + '-12-31'];
                    this.topYearSizeone = itme;
                }
            },
            // 月份
            selectItemtwo(itme,index){
                this.clickIndextwo = index;
                let newIndex = index + 1;
                if(newIndex < 10){
                    if(newIndex == 1 || newIndex == 3 || newIndex == 5 || newIndex == 7 || newIndex == 8){
                        this.topYeartwo = [this.yearValues + '-0' + newIndex +'-01',this.yearValues + '-0' + newIndex +'-31'];
                    }else if(newIndex == 2){
                        this.topYeartwo = [this.yearValues + '-0' + newIndex +'-01',this.yearValues + '-0' + newIndex + '-' + new Date(this.yearValues, 2, 0).getDate()]; // 二月最后一天按当年算:闰年是 29
                    }else{
                        this.topYeartwo = [this.yearValues + '-0' + newIndex +'-01',this.yearValues + '-0' + newIndex +'-30'];
                    }
                }else{
                    if(newIndex == 10 || newIndex == 12){
                        this.topYeartwo = [this.yearValues + '-' + newIndex +'-01',this.yearValues + '-' + newIndex +'-31'];
                    }else{
                        this.topYeartwo = [this.yearValues + '-' + newIndex +'-01',this.yearValues + '-' + newIndex +'-30'];
                    }
                    
                }
                this.topYearSizetwo = itme;
            },
            // 自定义时间范围
            customBoxclick(){
                this.$refs.tabTime.focus();
            },
            // 行政区划
            citys(value){
                if(value != ''){
                    this.enterpriseValue = '';
                    let invoiceDates = {
                        xzqhCode:value,
                        startTime:this.newscreeninvoice[0],
                        endTime:this.newscreeninvoice[1]
                    } 
                    this.countCapital({xzqhCode:value}); // 企业数量
                    // 总体业务量
                    this.countBillss({xzqhCode:value});   // 业务量金额累计
                    this.rateBill({xzqhCode:value});      // 纳税金额累计
                    this.logisticsBill({xzqhCode:value}); // 物流量累计
                    for(let i=0;i<this.newscreenlist.length;i++){
                        if(this.newscreenlist[i] == 1){
                            this.countBill({xzqhCode:value});              // 产品交易份额
                        }
                        if(this.newscreenlist[i] == 2){
                            this.paymentCapital(invoiceDates);  // 数量
                        }
                        if(this.newscreenlist[i] == 3){
                            this.logis(invoiceDates);           // 物流
                        }
                    }
                    this.api.areaInfo.county({parentCode:value})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.countyOptions = res.data.data; // 县 区
                        }
                    })
                    if(this.enterpriseValue == ''){
                        this.purShow = true;
                    }
                }
            },
            countys(value){
                if(value != ''){
                    this.enterpriseValue = '';
                    let invoiceDates = {
                        xzqhCode:value,
                        startTime:this.newscreeninvoice[0],
                        endTime:this.newscreeninvoice[1]
                    } 
                    this.countCapital({xzqhCode:value}); // 企业数量
                    // 总体业务量
                    this.countBillss({xzqhCode:value});   // 业务量金额累计
                    this.rateBill({xzqhCode:value});      // 纳税金额累计
                    this.logisticsBill({xzqhCode:value}); // 物流量累计
                    for(let i=0;i<this.newscreenlist.length;i++){
                        if(this.newscreenlist[i] == 1){
                            this.countBill({xzqhCode:value});       // 产品交易份额
                        }
                        if(this.newscreenlist[i] == 2){
                            this.paymentCapital(invoiceDates);  // 数量
                        }
                        if(this.newscreenlist[i] == 3){
                            this.logis(invoiceDates);           // 物流
                        }
                    }
                    this.api.areaInfo.village({parentCode:value})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.villageOptions = res.data.data; // 乡
                        }
                    })
                    if(this.enterpriseValue == ''){
                        this.purShow = true;
                    }
                }
            },
            villages(value){
                if(value != ''){
                    this.enterpriseValue = '';
                    let invoiceDates = {
                        xzqhCode:value,
                        startTime:this.newscreeninvoice[0],
                        endTime:this.newscreeninvoice[1]
                    }
                    this.countCapital({xzqhCode:value}); // 企业数量
                    // 总体业务量
                    this.countBillss({xzqhCode:value});   // 业务量金额累计
                    this.rateBill({xzqhCode:value});      // 纳税金额累计
                    this.logisticsBill({xzqhCode:value}); // 物流量累计
                    for(let i=0;i<this.newscreenlist.length;i++){
                        if(this.newscreenlist[i] == 1){
                            this.countBill({xzqhCode:value});       // 产品交易份额
                        }
                        if(this.newscreenlist[i] == 2){
                            this.paymentCapital(invoiceDates);  // 数量
                        }
                        if(this.newscreenlist[i] == 3){
                            this.logis(invoiceDates);           // 物流
                        }
                    }
                    this.api.statistics.enterData({countyCode:value,searchStr:''})
                    .then(res=>{
                        if(res.data.code == 200){
                            // this.enterpriseValue = '';
                            this.enterpriseOptions = res.data.data; // 企业名称
                        }
                    })
                    if(this.enterpriseValue == ''){
                        this.purShow = true;
                    }
                }
            },
            enterprises(id){
                if(id != ''){
                    // 通过企业id显示不同图表
                    let invoiceDates = {
                        searchStr:'',
                        companyId:id,
                        startTime:this.newscreeninvoice[0],
                        endTime:this.newscreeninvoice[1]
                    } 
                    // 总体业务量
                    this.countBillss({companyId:id,searchStr:''});   // 业务量金额累计
                    this.rateBill({companyId:id,searchStr:''});      // 纳税金额累计
                    this.logisticsBill({companyId:id,searchStr:''}); // 物流量累计
                    for(let i=0;i<this.newscreenlist.length;i++){
                        if(this.newscreenlist[i] == 1){
                            this.countBill({companyId:id,searchStr:''});       // 产品交易份额
                        }
                        if(this.newscreenlist[i] == 2){
                            this.paymentCapital(invoiceDates);  // 数量
                        }
                        if(this.newscreenlist[i] == 3){
                            this.logis(invoiceDates);           // 物流
                        }
                    }
                    if(this.enterpriseValue != ''){
                        this.purShow = false;
                        // 调用单个企业id，查询单条企业信息
                        this.api.enterData.byidData({id:id})
                        .then(res=>{
                            if(res.data.code == 200){
                                this.enterpriseName = res.data.data.enterpriseName;
                                this.regTime = res.data.data.regTime;
                            }
                        })
                    }
                }
            },
            // 清空行政区划
            citysclear(){
                this.cityValue = '';            // 市
                this.countyValue = '';          // 县 区
                this.countyOptions = [];
                this.villageValue = '';         // 乡
                this.villageOptions = [];
                this.enterpriseValue = '';      // 企业名称
                this.enterpriseOptions = [];
                let value = sessionStorage.getItem('xzqhCode')
                let invoiceDates = {
                    xzqhCode:value,
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                this.countCapital({xzqhCode:value}); // 企业数量
                // 总体业务量
                this.countBillss({xzqhCode:value});   // 业务量金额累计
                this.rateBill({xzqhCode:value});      // 纳税金额累计
                this.logisticsBill({xzqhCode:value}); // 物流量累计
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({xzqhCode:value});              // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDates);  // 数量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDates);           // 物流
                    }
                }
            },
            countysclear(){
                this.countyValue = '';          // 县 区
                this.villageValue = '';         // 乡
                this.villageOptions = [];
                this.enterpriseValue = '';      // 企业名称
                this.enterpriseOptions = [];
                let value
                if(this.cityValue == ''){
                    value = sessionStorage.getItem('xzqhCode')
                }else{
                    value = this.cityValue
                }
                let invoiceDates = {
                    xzqhCode:value,
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                this.countCapital({xzqhCode:value}); // 企业数量
                // 总体业务量
                this.countBillss({xzqhCode:value});   // 业务量金额累计
                this.rateBill({xzqhCode:value});      // 纳税金额累计
                this.logisticsBill({xzqhCode:value}); // 物流量累计
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({xzqhCode:value});              // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDates);  // 数量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDates);           // 物流
                    }
                }
            },
            villagesclear(){
                this.villageValue = '';         // 乡
                this.enterpriseValue = '';      // 企业名称
                this.enterpriseOptions = [];
                let value
                if(this.countyValue == ''){
                    value = sessionStorage.getItem('xzqhCode')
                }else{
                    value = this.countyValue
                }
                let invoiceDates = {
                    xzqhCode:value,
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                this.countCapital({xzqhCode:value}); // 企业数量
                // 总体业务量
                this.countBillss({xzqhCode:value});   // 业务量金额累计
                this.rateBill({xzqhCode:value});      // 纳税金额累计
                this.logisticsBill({xzqhCode:value}); // 物流量累计
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({xzqhCode:value});              // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDates);  // 数量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDates);           // 物流
                    }
                }
            },
            enterprisesclear(){
                this.enterpriseValue = '';      // 企业名称
                this.purShow = true;
                let value
                if(this.villageValue == ''){
                    value = sessionStorage.getItem('xzqhCode')
                }else{
                    value = this.villageValue
                }
                let invoiceDates = {
                    xzqhCode:value,
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                this.countCapital({xzqhCode:value}); // 企业数量
                // 总体业务量
                this.countBillss({xzqhCode:value});   // 业务量金额累计
                this.rateBill({xzqhCode:value});      // 纳税金额累计
                this.logisticsBill({xzqhCode:value}); // 物流量累计
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({xzqhCode:value});              // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDates);  // 数量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDates);           // 物流
                    }
                }
            },
            /* 
            企业数量
            */
            countCapital(type){
                this.loading = true;
                this.arrsumNum = 0;
                this.chasingTonnage = 0;
                this.api.enterData.count(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.loading = false;
                        this.chasingTonnage = res.data.data.sumCount;           // 注册企业数量
                        this.arrsumNum = res.data.data.monCount;                // 当月新增企业数量
                    }
                })
            },

            /* 
            总体业务量数据
            */
            countBillss(type){
                this.loading = true;
                this.TaxAmounts = 0;
                this.api.bill.countBuses(type)
                .then(res=>{
                    if(res.data.code == 200 && res.data.data != null){
                        this.loading = false;
                        this.TaxAmounts = res.data.data.baseMoney/100000000;    // 业务量金额累计
                    }
                })
            },
            rateBill(type){
                this.loading = true;
                this.rateBills = 0;
                this.api.billRate.rateBill(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.loading = false;
                        this.rateBills = res.data.data.sumMoney/10000;          // 纳税金额累计
                    }
                })
            },
            logisticsBill(type){
                this.loading = true;
                this.logisticsBills = 0;
                this.api.waybill.billAll(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.loading = false;
                        this.logisticsBills = res.data.data.tonnagePaid/10000;          // 物流量累计
                    }
                })
            },

            /* 
            数量数据
            */
            paymentCapital(type){
                this.loading = true;
                this.api.bill.business(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.datapayment = [];      // 图表时间
                        let newdate = [];           // 组装时间
                        this.paymentData = [];      // 业务量
                        this.paymentserData = [];   // 纳税量
                        for(let key in res.data.data){
                            newdate.push({date:key.replace(/-/g, "")*1});
                        }
                        newdate.sort(this.compare('date'));
                        for(let i=0;i<newdate.length;i++){
                            let year = newdate[i].date.toString().substring(0,4);
                            let newdatas =  newdate[i].date.toString().substring(4,6);
                            let length = newdate[i].date.toString().length
                            if(length > 6){
                                let daty = newdate[i].date.toString().substring(6,8);
                                this.datapayment.push(year + '-' + newdatas + '-' + daty)
                            }else{
                                this.datapayment.push(year + '-' + newdatas)
                            }
                        }
                        for(let itme in this.datapayment){
                            let data = res.data.data[this.datapayment[itme]]
                            this.paymentData.push((data/10000).toFixed(2))
                        }
                        this.api.billRate.amount(type)
                        .then(res=>{
                            if(res.data.code == 200){
                                this.loading = false;
                                for(let itme in this.datapayment){
                                    let data = res.data.data[this.datapayment[itme]]
                                    this.paymentserData.push((data/10000).toFixed(2))
                                }
                            }
                            this.payment();
                        }) 
                    }
                })
            },
            payment(){
                // 数量图表
                var chartDompayment = document.getElementById('payment');
                var myChartpayment = echarts.init(chartDompayment);
                var optionpayment;
                optionpayment = {
                    title: {
                        text: '数量'
                    },
                    tooltip: {
                        trigger: 'axis',
                        axisPointer: {
                            type: 'cross',
                            crossStyle: {
                                color: '#333'
                            }
                        }
                    },
                    legend: {
                        data: ['业务量', '纳税量'],
                        textStyle: {
                            color: "#333"
                        }
                    },
                    grid:{
                        x:50,
                        y:75,
                        x2:70,
                        y2:20,
                    },
                    xAxis: [
                        {
                            type: 'category',
                            name: '时间',
                            data:this.datapayment,
                            axisPointer: {
                                type: 'line'
                            },
                            nameTextStyle:{
                                color:'#333'
                            }
                        }
                    ],
                    yAxis: [
                        {
                            type: 'value',
                            name: '万元',
                            min: 0,
                            max: 10000,
                            interval: 2000,
                            axisLine: {
                                show: true,
                                lineStyle: {
                                    color: '#333'
                                }
                            },
                        }
                    ],
                    series: [
                        {
                            name: '业务量',
                            type: 'bar',
                            itemStyle: {
                                normal: {
                                    color: '#6FCF97'
                                }
                            },
                            data:this.paymentData
                        },
                        {
                            name: '纳税量',
                            type: 'bar',
                            itemStyle: {
                                normal: {
                                    color: '#2F80ED'
                                }
                            },
                            data:this.paymentserData
                        }
                    ]
                };
                optionpayment && myChartpayment.setOption(optionpayment);
            },

            /* 
            物流数据
            */
            logis(type){
                this.loading = true;
                this.api.waybill.countMonBuses(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.loading = false;
                        this.datalogistics = [];    // 时间
                        this.logisticsData = [];    // 物流量数据
                        let newdate = [];           // 组装时间数据
                        for(let key in res.data.data){
                            newdate.push({date:key.replace(/-/g, "")*1});
                        }
                        newdate.sort(this.compare('date'));
                        for(let i=0;i<newdate.length;i++){
                            let year = newdate[i].date.toString().substring(0,4);
                            let newdatas =  newdate[i].date.toString().substring(4,6);
                            let length = newdate[i].date.toString().length
                            if(length > 6){
                                let daty = newdate[i].date.toString().substring(6,8);
                                this.datalogistics.push(year + '-' + newdatas + '-' + daty)
                            }else{
                                this.datalogistics.push(year + '-' + newdatas)
                            }
                        }
                        for(let itme in this.datalogistics){
                            let data = res.data.data[this.datalogistics[itme]]
                            this.logisticsData.push(data.toFixed(2)/10000)
                        }
                        this.logissett(); 
                    }
                })
            },
            logissett(){
                // 物流
                var chartDomlogistics = document.getElementById('logistics');
                var myChartlogistics = echarts.init(chartDomlogistics);
                var optionlogistics;
                optionlogistics = {
                    title: {
                        text: '物流'
                    },
                    tooltip: {
                        trigger: 'axis',
                        axisPointer: {
                            type: 'cross',
                            crossStyle: {
                                color: '#333'
                            }
                        }
                    },
                    grid:{
                        x:50,
                        y:75,
                        x2:70,
                        y2:40,
                    },
                    xAxis: [
                        {
                            boundaryGap: false,
                            type: 'category',
                            name: '时间',
                            data: this.datalogistics,
                            axisPointer: {
                                type: 'line'
                            },
                            nameTextStyle:{
                                color:'#333'
                            },
                            // axisPointer: {
                            //     label: {
                            //         show: true,
                            //         backgroundColor: '#58afff'
                            //     },
                            // },
                        }
                    ],
                    yAxis: [
                        {
                            type: 'value',
                            name: '万吨',
                            min: 0,
                            max: 500,
                            interval: 100,
                            axisLine: {
                                show: true,
                                lineStyle: {
                                    color: '#333'
                                }
                            },
                        }
                    ],
                    series: [
                        {
                            type: 'line',
                            smooth: true,
                            symbol: 'none',
                            itemStyle: {
                                normal: {
                                    color: '#58afff'
                                }
                            },
                            areaStyle: {
                                color: '#EAF5FF'
                            },
                            data:this.logisticsData
                        }
                    ]
                };
                optionlogistics && myChartlogistics.setOption(optionlogistics);
            },

            /* 
            产品交易份额
            */
            countBill(type){
                
                this.loading = true;
                this.api.listByPage.countPro(type)
                .then(res=>{
                    if(res.data.code == 200){
                        this.loading = false;
                        this.countPro = [];
                        let color = ["#FB6B6C","#FDAD6D","#C16AF7","#686DFA","#FB8D6D","#6E58F9","#FED172","#6BE1FB","#88FA7A","#71FEE5"];
                        this.newColor = color.splice(0,res.data.data.length);
                        let obj;
                        let num = 0;
                        for(let k=0;k<res.data.data.length;k++){
                            num += res.data.data[k].count;
                        }
                        for(let i=0;i<res.data.data.length;i++){
                            let allnum = ((res.data.data[i].count)/num *100).toFixed(2) + '%';
                            obj = {
                                value:res.data.data[i].count,
                                name:res.data.data[i].d_key + allnum
                            }
                            this.countPro.push(obj);
                        }
                        this.chartDom();  // 税额
                    }
                })
            },
            chartDom(){
                var chartDom = document.getElementById('invoice');
                var myChart = echarts.init(chartDom);
                var option;
                option = {
                    tooltip: {
                        trigger: 'item'
                    },
                    legend: {
                        type: 'scroll',
                        orient: 'vertical',
                        left: 0,
                        top: 20,
                        bottom: 10,
                    },
                    color:this.newColor,
                    series: [
                        {
                            name: '基础维护',
                            type: 'pie',
                            radius: '60%',
                            center: ['77%', '50%'],
                            data:this.countPro
                        }
                    ]
                };
                option && myChart.setOption(option);
            },
        },
    }
</script>
<style  lang="less" scoped>
    .el-header{
        background-color: #CF1724;
        color: #fff;
        font-size: 14px;
        line-height: 60px;
        div{
            color: #fff;
            font-size: 14px;
            line-height: 60px;
            float: right;
            i{
                vertical-align: middle;
            }
        }
    }
    .el-main {
        background-color: #F5F5F5;
        padding: 0;
    }
    // 头部开始
    .pur-top {
        width: 100%;
        height: 50px;
        padding: 0 40px 0 32px;
        background-color: #fff;
        box-sizing: border-box;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        margin: 0 0 0.5% 0;
        .titleSpan {
            font-size: 14px;
            line-height: 50px;
            color: #333;
            font-weight: 600;
            margin-right:32px;
        }
        .screenBox{
            width:45%;
            height: 100%;
            float: right;
            padding:11px 0;
            box-sizing: border-box;
            .el-select{
                width: 23%;
                float: right;
                margin-left:1%;
            }
        }
        .Importfile{
            width: auto;
            height: 100%;
            float: right;
            cursor: pointer;
            line-height: 50px;
            font-weight: 600;
            i{
                font-size: 14px;
                color: #333;
                margin-right: 8px;
                font-weight: 600;
            }
            span{
                font-size: 14px;
                color: #333;
            }
        }
    }
    // 头部结束
    // 筛选样式开始
    ::v-deep .el-dialog__wrapper {
        .el-dialog {
            .el-dialog__body{
            padding: 0px;
            }
            .el-dialog__header{
            padding: 0px;
            }
            .el-dialog__footer{
                background-color: #F6F6F6;
                padding: 15px 20px;
                box-sizing: border-box;
            }
        }
    }
    .header-title{
        background-color: #F6F6F6;
        padding: 9px 16px;
        span{
            color:#333;
            line-height: 32px;
            font-size: 16px;
            margin-right:8px;
        }
    }
    .elTabBox{
        height: 330px;
        box-sizing: border-box;
        position: relative;
        overflow: hidden;
        .el-tabs{
            .el-tab-pane{
            padding:32px;
            ul{
                width: 100%;
                overflow: hidden;
                li{
                    width: 76px;
                    //    height: 28px;
                    //    border: 1px solid #333;
                    font-size: 14px;
                    text-align: center;
                    line-height: 16px;
                    color: #333;
                    border-radius: 4px;
                    float: left;
                    margin: 16px 16px;
                    padding: 6px;
                    box-sizing: border-box;
                    cursor: pointer;
                }
            }
                .selectItem{
                    background-color: #CF1724;
                    color:#fff;
                }
            }
        }
    }
    .customBox{
        span{
            font-size: 14px;
            color:#333;
            line-height: 32px;
            cursor: pointer;
        }
    }
    .zdyTime{
        width: 110px;
        height: 40px;
        position: absolute;
        left: 111px;
        top: -52px;
    }
    // 筛选样式结束
    // 内容开始
    .pur-nav {
        width:100%;
        height:calc(100vh - 123px);
        overflow-x: hidden;
        overflow-y: scroll;
        // 企业数量
        .Purchasing{
            float: left;
            width:33%;
            height:198px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;
            margin:0 0 0.5% 0.5%;
            padding:24px;
            box-sizing: border-box;
            .tity{
                font-size: 18px;
                font-weight: 500;
                color: #333;
                line-height: 32px;
                margin-left:8px;
            }
            span{
                line-height: 26px;
            }
            img{
                width: 32px;
                height: 32px;
                vertical-align: middle;
            }
            .Purchasing-conent{
                width: 100%;
                li{
                    width: 45%;
                    float: left;
                    padding-top:20px;
                    box-sizing: border-box;
                    text-align: center;
                    span{
                        font-size: 16px;
                        color: #333;
                        vertical-align: middle;
                        margin-left: 8px;
                    }
                    h5{
                        font-size: 32px;
                        color: #333;
                        line-height: 36px; 
                        font-weight: 600;
                    }
                }
            }
            .Purchasing-conentss{
                width: 100%;
                li{
                    width: 60%;
                    float: left;
                    padding-top:20px;
                    box-sizing: border-box;
                    text-align: center;
                    span{
                        font-size: 16px;
                        color: #333;
                        vertical-align: middle;
                        margin-left: 8px;
                    }
                    h5{
                        font-size: 18px;
                        color: #333;
                        line-height: 36px; 
                        font-weight: 600;
                    }
                } 
            }
        }
        // 总体业务量
        .quantity{
            float: left;
            width:65.5%;
            height:198px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;
            margin:0 0.5% 0.5% 0.5%;
            padding:24px;
            box-sizing: border-box;
            .tity{
                font-size: 18px;
                font-weight: 5 00;
                color: #333;
                line-height: 32px;
                margin-left:8px;
            }
            span{
                line-height: 26px;
            }
            img{
                width: 32px;
                height: 32px;
                vertical-align: middle;
            }
            .quantity-conent{
                width: 100%;
                li{
                    width: 30%;
                    height: 100%;
                    float: left;
                    padding-top:20px;
                    box-sizing: border-box;
                    text-align: center;
                    span{
                        font-size: 16px;
                        color: #333;
                        vertical-align: middle;
                        margin-left: 8px;
                    }
                    h5{
                        display: inline-block;
                        font-size: 32px;
                        color: #333;
                        line-height: 36px; 
                        font-weight: 600;
                    }
                }
            }
        }
        // 数量
        .payment{
            float: left;
            width:99%;
            height:280px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;   
            margin:0 0.5% 0.5%;
            position: relative;
            #payment{
                float: left;
                width: 100%;
                height: 100%;
                padding: 14px;
                box-sizing: border-box;
            }
            .payment-time{
                width:auto;
                height: 26px;
                background: #F0F0F0;
                border-radius: 60px;
                float: right;
                outline: none;
                padding: 2px 10px;
                box-sizing: border-box;
                cursor: pointer;
                color: #333;
                position: absolute;
                top:16px;
                right:16px;
                z-index: 1;
            }
        }
        // 物流
        .logistics{
            float: left;
            width:62.3%;
            height:340px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;  
            margin:0 0.5% 0.5% 0;
            margin-right: 0;
            position: relative;
            #logistics{
                float: left;
                width: 100%;
                height: 100%;
                padding: 14px;
                box-sizing: border-box;
            }
            .logistics-time{
                width:auto; 
                height: 26px;
                background: #F0F0F0;
                border-radius: 60px;
                float: right;
                outline: none;
                padding: 2px 10px;
                box-sizing: border-box;
                cursor: pointer;
                color: #333;
                position: absolute;
                top:16px;
                right:16px;
                z-index: 1;
            }
        }
        // 交易产品份额
        .invoice{
            float: left;
            width:36.2%;
            height:340px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;  
            margin:0 0.5% 0.5%;
            padding: 13px;
            box-sizing: border-box;
            overflow: hidden;
            span{
                font-size: 18px;
                color: #333;
                font-weight: 600;
                line-height: 32px;
            }
            .el-date-picker{
                width: 100%;
                float: right;
            }
            #invoice{
                width: 100%;
                height:90%;
                // padding-top:28px;
                // box-sizing: border-box;
                
            }
        }
    }
    // 内容结束
</style>