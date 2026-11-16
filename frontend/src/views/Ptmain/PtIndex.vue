<template>
    <div v-loading="loading" element-loading-text = "加载中...">
        <!-- 头部开始 -->
        <div class="pur-top">
            <span class="titleSpan">{{timedata}} </span>
            <el-button icon="el-icon-search" size="mini" @click="dialogVisibleclick">{{TimeFilter}}</el-button>
            <el-button size="mini" @click="echartsClick">数字大屏</el-button>
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
                    style="width:35%;margin-right:15px;"
                    size="small"
                    :clearable="false"
                    unlink-panels
                    v-model="screeninvoice"
                    :picker-options="pickerOptionsone"
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
                    style="width:180px;"
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
                        :picker-options="pickerOptionsone"
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
            <!-- 企业查询 -->
            <div class="screenBox">
                <el-select filterable v-model="enterpriseValue" placeholder="请选择企业名称" size="mini" style="min-width:100%;" @change="enterprises($event)" clearable @clear="enterprisesclear">
                    <el-option
                    v-for="item in enterpriseOptions"
                    :key="item.value"
                    :label="item.enterpriseName"
                    :value="item.id">
                    </el-option>
                </el-select>
            </div>
        </div>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
            <!-- 简报看板--平台 -->
            <!-- <div class="Purchasing" v-if="purShow == false">
                <img :src="qyimg" alt="">
                <span class="tity"> 简报看板</span>
                <ul class="Purchasing-conent">
                    <li style="width: 40%;">
                        <p>注册企业数</p>
                        <h5 v-if="chasingTonnage">{{chasingTonnage}}</h5>
                        <h5 v-else> 0 </h5>
                    </li>
                    <li>
                        <p>当月新增企业数</p>
                        <h5 v-if="arrsumNum">{{arrsumNum}}</h5>
                        <h5 v-else> 0 </h5>
                    </li> 
                    <li style="width:100%;">
                        <p>缴税金额累计</p>
                        <h5 v-if="rateBills" style="font-size:28px;">{{rateBills|newcurrency}} 亿</h5>
                        <h5 v-else style="font-size:28px;"> 0.00 亿</h5>
                        <span>(元)</span>
                    </li> 
                </ul>
            </div> -->
            <!-- 简报看板--企业 -->
            <!-- <div class="Purchasing" v-if="purShow == true">
                <img :src="qyimg" alt="">
                <span class="tity"> 简报看板</span>
                <ul class="Purchasing-conent">
                    <li>
                        <p>企业名称</p>
                        <h5 v-if="enterpriseName" style="font-size:18px;">{{enterpriseName}}</h5>
                        <h5 v-else>  </h5>
                    </li>
                    <li style="width: 35%;">
                        <p>上户时间</p>
                        <h5 v-if="regTime" style="font-size:18px;">{{regTime}}</h5>
                        <h5 v-else>  </h5>
                    </li>
                    <li>
                        <p>缴税金额累计</p>
                        <h5 v-if="rateBills" style="font-size:24px;">{{rateBills|newcurrency}} 万</h5>
                        <h5 v-else style="font-size:24px;"> 0.00 万</h5>
                        <span>(元)</span>
                    </li>  
                </ul>
            </div> -->
            <!-- 统计进度 -->
            <div class="quantity" style="width:32.6%;margin-right:0;">
                <img :src="tjimg" alt="">
                <span class="tity"> 业务量开票统计进度（亿元）</span>
                <ul class="quantity-conent">
                    <li style="width:auto;">
                        <p>录入完成</p>
                        <h5 v-if="TaxAmounts">{{TaxAmounts|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                    <li style="padding:15% 0 0 0;">
                        <el-progress :text-inside="true" :stroke-width="13" :percentage="ProportionbusCount" color="#6FCF97"></el-progress>
                    </li>
                    <li>
                        <p>业务完成</p>
                        <h5 v-if="newbusCount">{{newbusCount|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                </ul>
            </div>
            <div class="quantity" style="width:32.7%;margin-right:0;">
                <img :src="qyimg" alt="">
                <span class="tity"> 缴税量统计进度（亿元）</span>
                <ul class="quantity-conent">
                    <li style="width:auto;">
                        <p>录入完成</p>
                        <h5 v-if="rateBills">{{rateBills|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                    <li style="padding:15% 0 0 0;">
                        <el-progress :text-inside="true" :stroke-width="13" :percentage="newTax" color="#FF8058"></el-progress>
                    </li>
                    <li>
                        <p>业务完成</p>
                        <h5 v-if="rateBillsCount">{{rateBillsCount|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                </ul>
            </div>
            <div class="quantity" style="width:32.7%;margin-right:0;">
                <img :src="queryimg" alt="">
                <span class="tity"> 物流量统计进度（万吨）</span>
                <ul class="quantity-conent">
                    <li style="width:auto;">
                        <p>录入完成</p>
                        <h5 v-if="logisticsBills">{{logisticsBills|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                    <li style="padding:15% 0 0 0;">
                        <el-progress :text-inside="true" :stroke-width="13" :percentage="ProportlogisticsCount" color="#58AFFF"></el-progress>
                    </li>
                    <li>
                        <p>开票完成</p>
                        <h5 v-if="logisticsCount">{{logisticsCount|newcurrency}} </h5>
                        <h5 v-else> 0.00 </h5>
                    </li>
                </ul>
            </div>
            <div style="clear:both;"></div>
            <!-- 票据 -->
            <div class="billnote" v-show="purShow">
                <span>票据</span>
                <div id="billnote"></div>
            </div>
            <!-- 税负率 -->
            <div class="burden" :style="styleburden">
                <div id="burden"></div>
            </div>
            <div style="clear:both;"></div>
            <!-- 所得税率 -->
            <div class="allTaxrate" v-show="purShow == false">
                <div id="allTaxrate"></div>
            </div>
            <!-- 业务量 -->
            <div class="payment">
                <div id="payment"></div>
            </div>
            <!-- 客户统计工作 -->
            <div class="countMonData" v-show="purShow == false">
                <div id="countMonData"></div>
            </div>
            <!-- 网络货运占比 -->
            <!-- <div class="NetworkFreight" v-show="purShow == false">
                <div id="burdenLogins"></div>
            </div> -->
            <!-- 缴税量 -->
            <div class="payment">
                <div id="paymentTaxes"></div>
            </div>
            <!-- 物流 -->
            <div class="logistics">
                <div id="logistics"></div>
            </div>
            <!-- 产品交易份额 -->
            <div class="invoice">
                <span>产品交易份额</span>
                <div id="invoicess" key="invoicess"></div>
                <div class="invoiceBox" v-show="invoiceShow == 2">暂无数据</div>
            </div>
        </div>
        <!-- 内容结束 -->
    </div>
</template>
<script>
import qyimg from '../../assets/qiye.svg';
import tjimg from '../../assets/tongji.svg';
import queryimg from '../../assets/query2.svg';
var echarts = require('echarts');
export default {
    data() {
        return {
            /*
            筛选功能弹窗
            */
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
            screenlist:[1,2,3,4,5,6,7,8],                 // 类别
            newscreenlist:[1,2,3,4,5,6,7,8],              // 选中类别
            screeninvoice:[],                       // 时间
            newscreeninvoice:[],                    // 选中时间
            screenOptions: [
                {
                    value: 1,
                    label: '产品交易份额'
                },
                {
                    value: 2,
                    label: '业务量占比'
                },
                {
                    value: 3,
                    label: '物流量'
                },
                {
                    value: 4,
                    label: '税负率'
                },
                {
                    value: 5,
                    label: '缴税额'
                },
                {
                    value: 6,
                    label: '客户数占比'
                },
                {
                    value: 7,
                    label: '数据进度'
                },
                {
                    value: 8,
                    label: '所得税率'
                },
            ],
            qyimg:qyimg,
            tjimg:tjimg,
            queryimg:queryimg,
            purShow:false,
            username:'',
            pickerOptionsone:{
                disabledDate:(time)=>{
                    return time.getTime() > Date.now() || time.getTime() == Date.now();
                }
            },
            loading:false,
            timedata:'',
            // 地区筛选
            enterpriseOptions: [],          // 企业名称
            enterpriseValue: '',
            enterpriseName:'',              // 企业列表数据

            // 数据简报平台
            // chasingTonnage:0,               // 注册企业数量
            // arrsumNum:0,                    // 当月新增企业数量
            // regTime:'',

            // 统计进度
            TaxAmounts:0,                   // 系统录入业务量金额累计
            logisticsBills:0,               // 系统录入物流量累计
            rateBills:0,                    // 缴税录入金额累计
            newbusCount:0,                  // 实际业务量 
            logisticsCount:0,               // 实际物流量
            rateBillsCount:0,               // 实际缴税金额累计
            ProportionbusCount:0,           // 系统/实际业务量占比
            ProportlogisticsCount:0,        // 系统/实际物流量占比
            newTax:0,                       // 系统录入/实际缴税占比
            
            // 客户统计工作
            countMonment:[],                // 时间
            countMonData:[],                // 数据完成户数
            countMonDatas:[],                // 开票业务户数

            // 网络货运平台占比
            optionLogins: null,
            number: 0, // 播放所在下标
            timer: null,
            burdenLoginsTime:[],              // 公司名称
            burdenLoginsData:[],              // 网络货运占比
            allLoginsData:[],                 // 运输发票占比
            ReservedTon:[],                   // 网络货运未开票占比
            AllReservedTon:[],                // 运输发票未开票占比

            // 业务量
            datapayment:[],                 // 时间
            paymentData:[],                 // 系统业务量
            statisticsData:[],              // 实际业务量

            // 缴税量
            paymentdatapayment:[],          // 时间
            paymentserData:[],              // 录入缴税量
            taxPaymentData:[],              // 实际缴税量

            // 物流
            datalogistics:[],               // 时间
            logisticsData:[],               // 系统统计
            logstatisticsData:[],           // 实际统计

            // 产品交易份额
            invoiceShow:1,
            newColor:[],                // 颜色
            countPro:[],
            invoice:'',
            sumTaxAmount:'',            // 进项税额
            arrTaxAmount:'',            // 销项税额

            // 税负率
            burdenTime:[],              // 公司名称
            burdenData:[],              // 税负率
            styleburden:'width:99%;',

            // 所得税率
            allTaxrateTime:[],              // 公司名称
            allTaxrateData:[],              // 所得税率

            // 票据税额进项比
            invoice:'',
            sumTaxAmountbillnote:'',            // 进项税额
            arrTaxAmountbillnote:'',            // 销项税额
        }
    },
    created() {
        this.username = sessionStorage.getItem('username');
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
        this.username = sessionStorage.getItem('username');
        this.http  = this.$store.state.http;                 // IP地址
    },
    mounted(){
        // 行政区划
        this.api.statistics.enterData({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
        .then(res=>{
            if(res.data.code == 200){
                this.enterpriseOptions = res.data.data; // 企业名称
            }
        })
        // 获取当前时间
        let date = new Date();
        let year = date.getFullYear();
        let month = date.getMonth() + 1;
        let day = date.getDate();
        // let newmonth = month - 1;
        // let newmonths = newmonth< 10? '0'+ newmonth:newmonth
        // let invoiceDatestwo = {
        //     xzqhCode:sessionStorage.getItem('xzqhCode'),
        //     startTime:year + '-01-01',
        //     endTime:year + '-' + newmonths + '-01'
        // };
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
        // this.countCapital({xzqhCode:sessionStorage.getItem('xzqhCode')});
        // this.rateBill({xzqhCode:sessionStorage.getItem('xzqhCode')});      // 缴税金额累计
        this.logisticsBill(invoiceDate); // 统计进度数据
        // 客户统计工作
        this.countMonDataTaxes(invoiceDate);  
        // 网络货运平台占比
        this.burdenLoginsapi(invoiceDate);                  
        // 业务量
        this.paymentCapital(invoiceDate);
        // 缴税量
        this.paymentCapitalTaxes(invoiceDate);
        // 物流
        this.logis(invoiceDate);
        // 产品交易份额
        this.countBill(invoiceDate);
        // 税负率
        this.burdenapi(invoiceDate);
        // 所得税率
        this.allTaxrateapi(invoiceDate);        
    },
    methods: {
        echartsClick(){
            this.$router.push({name:"Dpmain", params: {show:true}})
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
            this.screenlist = [1,2,3,4,5,6,7,8];
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
            // this.enterpriseValue = '';
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
            let invoiceDatesOne = {
                xzqhCode:sessionStorage.getItem('xzqhCode'),
                startTime:this.newscreeninvoice[0],
                endTime:this.newscreeninvoice[1]
            } 
            for(let i=0;i<this.newscreenlist.length;i++){
                if(this.newscreenlist[i] == 1){
                    this.countBill({xzqhCode:sessionStorage.getItem('xzqhCode')});       // 产品交易份额
                }
                if(this.newscreenlist[i] == 2){
                    this.paymentCapital(invoiceDatesOne);   // 业务量
                }
                if(this.newscreenlist[i] == 3){
                    this.logis(invoiceDatesOne);            // 物流
                }
                if(this.newscreenlist[i] == 4){
                    this.burdenapi(invoiceDatesOne);        // 税负率
                }
                if(this.newscreenlist[i] == 5){
                    this.paymentCapitalTaxes(invoiceDatesOne);  // 缴税量
                }
                if(this.newscreenlist[i] == 6){
                    this.countMonDataTaxes(invoiceDatesOne);  // 客户统计工作
                }
                if(this.newscreenlist[i] == 7){
                    this.logisticsBill(invoiceDatesOne);      // 统计进度
                }
                if(this.newscreenlist[i] == 8){
                    this.allTaxrateapi(invoiceDatesOne);        // 所得税率
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
                    // 获取当前时间  b
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
                    this.topYeartwo = [this.yearValues + '-0' + newIndex +'-01',this.yearValues + '-0' + newIndex +'-28'];
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
        // 选中企业数据
        enterprises(id){
            if(id != ''){
                // 通过企业id显示不同图表
                let invoiceDatesAll = {
                    searchStr:'',
                    companyId:id,
                    startTime:this.newscreeninvoice[0],
                    endTime:this.newscreeninvoice[1]
                } 
                // this.rateBill({companyId:id,searchStr:''});      // 缴税金额累计
                // this.logisticsBill({companyId:id,searchStr:''}); // 统计进度数据
                for(let i=0;i<this.newscreenlist.length;i++){
                    if(this.newscreenlist[i] == 1){
                        this.countBill({companyId:id,searchStr:''});       // 产品交易份额
                    }
                    if(this.newscreenlist[i] == 2){
                        this.paymentCapital(invoiceDatesAll);       // 业务量
                    }
                    if(this.newscreenlist[i] == 3){
                        this.logis(invoiceDatesAll);                // 物流
                    }
                    if(this.newscreenlist[i] == 4){
                        this.burdenapi(invoiceDatesAll);            // 税负率
                    }
                    if(this.newscreenlist[i] == 5){
                        this.paymentCapitalTaxes(invoiceDatesAll);  // 缴税量
                    }
                    if(this.newscreenlist[i] == 6){
                        this.countMonDataTaxes(invoiceDatesAll);    // 客户统计工作
                    }
                    if(this.newscreenlist[i] == 7){
                        this.logisticsBill(invoiceDatesAll);      // 统计进度
                    }
                    if(this.newscreenlist[i] == 8){
                        this.allTaxrateapi(invoiceDatesAll);        // 所得税率
                    }
                }
                if(this.enterpriseValue != ''){
                    this.purShow = true;
                    this.styleburden = "width:69.2%;";
                    this.countBillbillnote(invoiceDatesAll);
                    
                    // 调用单个企业id，查询单条企业信息
                    this.api.enterData.byidData({id:id})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.enterpriseName = res.data.data.enterpriseName;
                            // this.regTime = res.data.data.regTime;
                        }
                    })
                }
            }
        },
        // 清空行政区划企业数据
        enterprisesclear(){
            this.enterpriseValue = '';      // 企业名称
            this.purShow = false;
            this.styleburden = "width:99%;";
            let value = sessionStorage.getItem('xzqhCode');
            let invoiceDatesClear = {
                xzqhCode:value,
                startTime:this.newscreeninvoice[0],
                endTime:this.newscreeninvoice[1]
            } 
            // this.countCapital({xzqhCode:value}); // 企业数量
            // this.rateBill({xzqhCode:value});      // 缴税金额累计
            // this.logisticsBill({xzqhCode:value}); // 统计进度数据
            for(let i=0;i<this.newscreenlist.length;i++){
                if(this.newscreenlist[i] == 1){
                    this.countBill({xzqhCode:value});              // 产品交易份额
                }
                if(this.newscreenlist[i] == 2){
                    this.paymentCapital(invoiceDatesClear);  // 业务量
                }
                if(this.newscreenlist[i] == 3){
                    this.logis(invoiceDatesClear);           // 物流
                }
                if(this.newscreenlist[i] == 4){
                    this.burdenapi(invoiceDatesClear);        // 税负率
                }
                if(this.newscreenlist[i] == 5){
                    this.paymentCapitalTaxes(invoiceDatesClear);  // 缴税量
                }
                if(this.newscreenlist[i] == 6){
                    this.countMonDataTaxes(invoiceDatesClear);  // 客户统计工作
                }
                if(this.newscreenlist[i] == 7){
                    this.logisticsBill(invoiceDatesClear);      // 统计进度
                }
                if(this.newscreenlist[i] == 8){
                    this.allTaxrateapi(invoiceDatesClear);      // 所得税率
                }
            }
        },
        /* 
        企业数量
        */
        // countCapital(type){
        //     this.loading = true;
        //     this.arrsumNum = 0;
        //     this.chasingTonnage = 0;
        //     this.api.statistics.count(type)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             this.loading = false;
        //             this.chasingTonnage = res.data.data.sumCount;           // 注册企业数量
        //             this.arrsumNum = res.data.data.monCount;                // 当月新增企业数量
        //         }
        //     })
        // },
        // 缴税金额累计
        // rateBill(type){
        //     this.loading = true;
        //     this.rateBills = 0;
        //     this.api.billRate.rateBill(type)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             this.loading = false;
        //             this.rateBills = res.data.data.sumMoney/100000000;          // 缴税金额累计
        //         }
        //     })
        // },   

        /* 
        统计进度数据
        */
        logisticsBill(type){
            this.newbusCount = 0 ;
            this.logisticsCount = 0;
            this.TaxAmounts = 0;
            this.logisticsBills = 0;
            this.rateBills = 0;
            this.rateBillsCount = 0;
            this.api.statistics.countActBus(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.newbusCount = res.data.data.busCount/100000000;                                    // 开票业务完成
                    this.logisticsCount = res.data.data.logisticsCount/10000;                                 // 开票物流完成
                    this.rateBillsCount = res.data.data.taxCount/100000000;                                 // 开票缴税完成
                    if(res.data.data.transCount != null){
                        this.logisticsBills = res.data.data.transCount/10000;          // 物流录入完成累计
                    }else{
                        this.logisticsBills = 0.00;                                     
                    };
                    if(this.logisticsBills - this.logisticsCount > 0){
                        this.ProportlogisticsCount = 100 *1;
                    }else if(this.logisticsBills - this.logisticsCount < 0){
                        this.ProportlogisticsCount = ((this.logisticsBills/this.logisticsCount)*100).toFixed(0)*1;
                    }else{
                        this.ProportlogisticsCount = 0;
                    }
                    this.api.bill.countBuses(type)
                    .then(res=>{
                        if(res.data.code == 200){
                            if(res.data.data != null){
                                this.TaxAmounts = res.data.data.baseMoney/100000000;   // 开票录入完成金额累计
                            }else{
                                this.TaxAmounts = 0;    
                            }
                            if(this.TaxAmounts - this.newbusCount > 0){
                                this.ProportionbusCount = 100 *1;
                            }else if(this.TaxAmounts - this.newbusCount < 0){
                                this.ProportionbusCount = ((this.TaxAmounts/this.newbusCount)*100).toFixed(0)*1;
                            }else{
                                this.ProportionbusCount = 0;
                            }
                        }
                    });
                    this.api.billRate.rateBill(type)
                    .then(res=>{
                        if(res.data.code == 200){
                            if(res.data.data != null){
                                this.rateBills = res.data.data.sumMoney/100000000;          // 录入缴税金额累计
                            }else{
                                this.rateBills = 0.00;
                            }
                            if(this.rateBills - this.rateBillsCount > 0){
                                this.newTax = 100 *1;
                            }else if(this.rateBills - this.rateBillsCount < 0){
                                this.newTax = ((this.rateBills/this.rateBillsCount)*100).toFixed(0)*1;
                            }else{
                                this.newTax = 0;
                            }
                        }
                    })
                }
            })
        },

        /* 
        物流简报数据
        */
        burdenLoginsapi(buildinvoiceDate){
            this.api.statistics.countNetLogist(buildinvoiceDate)
            .then(res=>{
                if(res.data.code == 200){
                    let allinit = [];
                    let obj = {};
                    let o;
                    for(let iev in res.data.data.netlist){
                        obj[res.data.data.netlist[iev].simpleName] = res.data.data.netlist[iev].lrate;
                    }                    
                    for(let key in res.data.data.translist){
                        if(obj[res.data.data.translist[key].simpleName]){
                            o = {
                                simpleName:res.data.data.translist[key].simpleName,
                                lrate:obj[ res.data.data.translist[key].simpleName ],           // 网络货运
                                yslrate:res.data.data.translist[key].lorate                     // 运输发票
                            }
                        }else{
                            o = {
                                simpleName:res.data.data.translist[key].simpleName,
                                lrate:0,                                                        // 网络货运
                                yslrate:res.data.data.translist[key].lorate                     // 运输发票
                            }                           
                        }
                        allinit.push(o)
                    }
                    this.burdenLoginsTime = [];   // 企业名称
                    this.burdenLoginsData = [];   // 网络货运占比
                    this.allLoginsData = [];      // 运输发票占比
                    this.ReservedTon = [];        // 网络货运未开票占比
                    this.AllReservedTon = [];     // 运输发票未开票占比
                    for(let key in allinit){
                        this.burdenLoginsTime.push(allinit[key].simpleName);
                        if((allinit[key].lrate*100).toFixed(2)*1 < 100){
                            this.burdenLoginsData.push((allinit[key].lrate*100).toFixed(2)*1);
                            this.ReservedTon.push(0);
                        }else{
                            this.burdenLoginsData.push(100);
                            let a = ((allinit[key].lrate*100)*1 - 100).toFixed(2)*1
                            this.ReservedTon.push(a);
                        } 
                        if((allinit[key].yslrate*100).toFixed(2)*1 < 100){
                            this.allLoginsData.push((allinit[key].yslrate*100).toFixed(2)*1);
                            this.AllReservedTon.push(0);
                        }else{
                            this.allLoginsData.push(100);
                            let b = ((allinit[key].yslrate*100)*1 - 100).toFixed(2)*1
                            this.AllReservedTon.push(b);
                        } 
                    }
                    this.burdenLogins();  // 网络货运平台占比
                }
            })
        },
        burdenLogins(){
            let a = document.getElementById("burdenLogins").removeAttribute("burdenLogins",'');
            let chartsmyChartLogins = echarts.init(document.getElementById('burdenLogins'));
            let optionLogins = {
                title: {
                    text: '物流量占比',
                    top: '0%',
                    left:'0%',
                    textAlign: 'left',
                    textStyle:{
                        color:'#333',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        crossStyle: {
                            color: '#333'
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        list.push(
                            a[0].axisValue
                        )
                        list.push(
                            '&nbsp&nbsp网络货运占比：' +
                            (a[0].value + a[1].value).toFixed(2) + `%<br/>` +
                            '&nbsp&nbsp网络货运未开票占比：' +
                            a[1].value + `%<br/>` + 
                            '&nbsp&nbsp运输发票占比：' +
                            (a[2].value + a[3].value).toFixed(2) + `%<br/>` + 
                            '&nbsp&nbsp运输发票未开票占比：' +
                            a[3].value + `%`
                        )
                        listItem = list.join('<br>')
                        return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                legend: {
                    selectedMode:false,
                    data:['网络货运占比','运输发票占比'],
                    textStyle: {
                        color: "#333"
                    }
                },
                grid:{
                    x:50,
                    y:75,
                    x2:70,
                    y2:50,
                },
                xAxis: [
                    {
                        type: 'category',
                        name:'企业',
                        data:  this.burdenLoginsTime,
                        axisPointer: {
                            type: 'line'
                        },
                        nameTextStyle:{
                            color:'#333'
                        },
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333',
                            },
                        },
                        nameTextStyle: {
                            color: '#333',
                        },
                        // axisPointer: {
                        //     label: {
                        //         show: true,
                        //         backgroundColor: 'rgba(247, 125, 88, .4)',
                        //     },
                        // },
                    }
                ],
                dataZoom: [
                    {
                        type: 'slider',
                        // backgroundColor:"#FCDEDE",
                        // handleColor: '#CF1724',                     // h滑动图标的颜色
                        // fillerColor:"#F5C4C4",                      // 选中范围的填充颜色。
                        // borderColor:"#CF1724",                      // 边框颜色。
                        show: true,                                 // flase直接隐藏图形
                        xAxisIndex: [0],
                        left: '5%',                                 // 滚动条靠左侧的百分比
                        bottom: -5,
                        start: 0,                                   // 滚动条的起始位置
                        end: 100,                                    // 滚动条的截止位置（按比例分割你的柱状图x轴长度）
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        min: 0,
                        max: 140,
                        interval:20,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333',
                            },
                        },
                        nameTextStyle: {
                            color: '#333',
                        },
                        axisPointer: {
                            label: {
                                show: true,
                                backgroundColor: '#333',
                            },
                        },
                        axisLabel: {
                            formatter: '{value} %'
                        },
                    }
                ],
                series:[
                    {
                        type: 'bar',
                        name:'',
                        // barWidth:'30%',
                        name:'网络货运占比',
                        stack: '网络货运占比',
                        color:'rgba(111, 207, 151, 1)',
                        data:this.burdenLoginsData,
                        // barWidth : 40,//柱图宽度
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                color:'#333',
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "70%"
                            },
                            data : [
                                {
                                    name: '',
                                    silent:false,             //鼠标悬停事件  true没有，false有
                                    lineStyle:{               //警戒线的样式  ，虚实  颜色
                                        type:"dashed",
                                        color:"#6FCF97"
                                    },
                                    yAxis:70,
                                }
                            ]
                        },
                        // barGap: '-100%', // 设置柱形重合的重要步骤
                    },
                    {
                        name: '网络货运未开票占比',
                        type: 'bar',
                        // barWidth:'30%',
                        stack: '网络货运占比',
                        data:this.ReservedTon,
                        color:'rgba(111, 207, 151, .4)',
                    },
                    {
                        type: 'bar',
                        name:'',
                        name:'运输发票占比',
                        stack: '运输发票占比',
                        color:'rgba(111, 178, 207, 1)',
                        data:this.allLoginsData,
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                color:'#333',
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "100%"
                            },
                            data : [
                                {
                                    name: '',
                                    silent:false,             //鼠标悬停事件  true没有，false有
                                    lineStyle:{               //警戒线的样式  ，虚实  颜色
                                        type:"dashed",
                                        color:"rgba(111, 178, 207, 1)"
                                    },
                                    yAxis:100,
                                }
                            ]
                        },
                        // label: {
                        //     normal: {
                        //         show: true, //是否显现，不显示的话设置成false
                        //         color:'#333',
                        //         position: "top", //显示的位置
                        //         distance: 0, //距离侄子的值 // label要显示的值比如： 20%
                        //         formatter: function(param) {
                        //             return param.value + '%';
                        //         },
                        //     }
                        // },
                        // barGap: '-100%', // 设置柱形重合的重要步骤
                    },
                    {
                        name: '运输发票未开票占比',
                        type: 'bar',
                        // barWidth:'30%',
                        stack: '运输发票占比',
                        data:this.AllReservedTon,
                        color:'rgba(111, 178, 207, .4)',
                    },
                ],
            };
            optionLogins && chartsmyChartLogins.setOption(optionLogins);
            // window.addEventListener('resize', () => {
            //     chartsmyChartLogins.resize();
            // });  

            // this.timer = setInterval(() => {
            //     chartsmyChartLogins.dispatchAction({
            //         type: 'showTip',
            //         seriesIndex: 0,
            //         dataIndex: this.number
            //     });
            //     this.number++;
            //     if (this.number > this.burdenLoginsTime.length) {
            //         this.number = 0;
            //     }
            // }, 2000);
        },

        /* 
        业务量占比数据
        */
        paymentCapital(type){
            this.loading = true;
            this.api.bill.business(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.datapayment = [];      // 图表时间
                    let newdate = [];           // 组装时间
                    this.paymentData = [];      // 系统业务量
                    this.statisticsData = [];   // 实际业务量
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
                        this.paymentData.push((data).toFixed(2))
                    }
                    this.api.statistics.countMonActBus(type)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.loading = false;
                            for(let itme in this.datapayment){
                                let data = res.data.data.busCount[this.datapayment[itme]]
                                this.statisticsData.push(data.toFixed(2))
                            }
                        }
                        this.payment();
                    }) 
                }
            })
        },
        payment(){
            // 业务量图表
            var chartDompayment = document.getElementById('payment');
            var myChartpayment = echarts.init(chartDompayment);
            var optionpayment;
            let barWidths;
            if(this.statisticsData.length <= 6){
                barWidths = 60;
            }else{
                barWidths = 25;
            }
            optionpayment = {
                title: {
                    text: '业务量占比'
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        crossStyle: {
                            color: '#333'
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        let b;
                        list.push(
                            a[0].axisValue
                        )
                        for (var i = 0; i <a.length; i++) {
                        list.push(
                            '<i style="display: inline-block;width: 10px;height: 10px;background: ' +
                                a[i].color +
                                ';margin-right: 5px;border-radius: 50%;}"></i><span style="width:120px; display:inline-block;">' +
                                a[i].seriesName + 
                                '</span>' + a[i].data + `（元）`
                            )
                        }
                        if(a[1].data/a[0].data == Infinity){
                            b = 0;
                        }else if(a[1].data == 0 && a[0].data == 0){
                            b = 0;
                        }else if(a[1].data/a[0].data >= 100){
                            b = 100;
                        }else{
                            b = ((a[1].data/a[0].data)*100).toFixed(2);
                        }
                        list.push(
                            '<br>&nbsp&nbsp进度占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                textStyle: {
                    color: '#333'
                },
                legend: {
                    selectedMode:false,
                    data: ['开票业务完成', '数据录入完成'],
                    textStyle: {
                        color: "#333"
                    }
                },
                grid:{
                    x:70,
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
                        name: '元',
                        // min: 0,
                        // max: 16000,
                        // interval: 4000,
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
                        name: '开票业务完成',
                        type: 'bar',
                        color: '#f5f5f5',
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        data: this.statisticsData,
                        z: 0,
                        silent: true,
                        animation: true, // 关闭动画效果
                        barWidth : barWidths,      //柱图宽度
                        // label: {
                        //     normal: {
                        //         show: true,
                        //         // formatter: '({c}-{c})%',
                        //         position: "top",
                        //         textStyle: {
                        //             color: "#333",
                        //             fontSize: 14
                        //         }
                        //     }
                        // }
                    },
                    {
                        name: '数据录入完成',
                        type: 'bar',
                        color:'#6FCF97',
                        data: this.paymentData,
                        barWidth : barWidths,      //柱图宽度
                        barGap: '-100%', // 设置柱形重合的重要步骤
                    }
                ]
            };
            optionpayment && myChartpayment.setOption(optionpayment);
        },

        /* 
        客户数占比数据
        */
        countMonDataTaxes(type){
            this.loading = true;
            this.api.statistics.countMonData(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.countMonment = [];     // 图表时间
                    let newdate = [];           // 组装时间
                    this.countMonData = [];     // 数据完成户数
                    this.countMonDatas = [];    // 开票业务户数
                    for(let key in res.data.data.openCount){
                        newdate.push({date:key.replace(/-/g, "")*1});
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        let length = newdate[i].date.toString().length
                        if(length > 6){
                            let daty = newdate[i].date.toString().substring(6,8);
                            this.countMonment.push(year + '-' + newdatas + '-' + daty)
                        }else{
                            this.countMonment.push(year + '-' + newdatas)
                        }
                    }
                    for(let itme in this.countMonment){
                        let datas = res.data.data.openCount[this.countMonment[itme]];
                        this.countMonDatas.push(datas);
                        let data = res.data.data.importCount[this.countMonment[itme]];
                        this.countMonData.push(data);
                    }
                    this.countMonDatasss();
                }
            })
        },
        countMonDatasss(){
            // 客户数占比数据图表
            document.getElementById("countMonData").removeAttribute("_echarts_instance_",'');
            var countMonchartDompayment = document.getElementById('countMonData');
            var countmyChartpayment = echarts.init(countMonchartDompayment);
            var countoptionpayment;
            let barWidths;
            if(this.countMonDatas.length <= 6){
                barWidths = 60;
            }else{
                barWidths = 25;
            }
            countoptionpayment = {
                title: {
                    text: '客户数占比'
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        crossStyle: {
                            color: '#333'
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        let b;
                        list.push(
                            a[0].axisValue
                        )
                        for (var i = 0; i <a.length; i++) {
                        list.push(
                            '<i style="display: inline-block;width: 10px;height: 10px;background: ' +
                                a[i].color +
                                ';margin-right: 5px;border-radius: 50%;}"></i><span style="width:120px; display:inline-block;">' +
                                a[i].seriesName + 
                                '</span>' + a[i].data + `户数`
                            )
                        }
                        if(a[1].data/a[0].data == Infinity){
                            b = 0;
                        }else if(a[1].data == 0 && a[0].data == 0){
                            b = 0;
                        }else if(a[1].data/a[0].data >= 100){
                            b = 100;
                        }else{
                            b = ((a[1].data/a[0].data)*100).toFixed(0);
                        }
                        list.push(
                            '<br>&nbsp&nbsp进度占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                textStyle: {
                    color: '#333'
                },
                legend: {
                    selectedMode:false,
                    data: ['开票业务户数', '数据完成户数'],
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
                        data:this.countMonment,
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
                        name: '客户数',
                        // min: 0,
                        // max: 50,
                        // interval: 10,
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
                        name: '开票业务户数',
                        type: 'bar',
                        color: '#f5f5f5',
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        data: this.countMonDatas,
                        z: 0,
                        silent: true,
                        animation: true, // 关闭动画效果
                        barWidth:barWidths,
                    },
                    {
                        name: '数据完成户数',
                        type: 'bar',
                        color:'#FF8058',
                        data: this.countMonData,
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        barWidth:barWidths,
                    }
                ]
            };
            countoptionpayment && countmyChartpayment.setOption(countoptionpayment);
        },

        /* 
        缴税额数据
        */
        paymentCapitalTaxes(type){
            this.loading = true;
            this.api.billRate.amount(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.paymentdatapayment = [];       // 图表时间
                    let newdate = [];                   // 组装时间
                    this.paymentserData = [];           // 录入缴税量
                    this.taxPaymentData = [];           // 完成缴税量
                    for(let key in res.data.data){
                        newdate.push({date:key.replace(/-/g, "")*1});
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        let length = newdate[i].date.toString().length
                        if(length > 6){
                            // let daty = newdate[i].date.toString().substring(6,8);
                            // this.paymentdatapayment.push(year + '-' + newdatas + '-' + daty)
                        }else{
                            this.paymentdatapayment.push(year + '-' + newdatas)
                        }
                    }
                    for(let itme in this.paymentdatapayment){
                        let data = res.data.data[this.paymentdatapayment[itme]]
                        this.paymentserData.push(data.toFixed(2))
                    }
                    this.api.statistics.countMonActBus(type).then(res=>{
                        if(res.data.code == 200){
                            this.loading = false;
                            for(let itme in this.paymentdatapayment){
                                let data = res.data.data.taxPayment[this.paymentdatapayment[itme]];
                                this.taxPaymentData.push(data.toFixed(2))
                            }
                        }
                        this.paymentTaxes();
                    }) 
                }
            })
        },
        paymentTaxes(){
            // 缴税额图表
            var chartDompaymentss = document.getElementById('paymentTaxes');
            var myChartpaymentss = echarts.init(chartDompaymentss);
            var optionpaymentss;
            let barWidths;
            if(this.taxPaymentData.length <= 6){
                barWidths = 60;
            }else{
                barWidths = 25;
            }
            optionpaymentss = {
                title: {
                    text: '缴税额'
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        crossStyle: {
                            color: '#333'
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        let b;
                        list.push(
                            a[0].axisValue
                        )
                        for (var i = 0; i <a.length; i++) {
                        list.push(
                            '<i style="display: inline-block;width: 10px;height: 10px;background: ' +
                                a[i].color +
                                ';margin-right: 5px;border-radius: 50%;}"></i><span style="width:120px; display:inline-block;">' +
                                a[i].seriesName + 
                                '</span>' + a[i].data + `（元）`
                            )
                        }
                        if(a[1].data/a[0].data == Infinity){
                            b = 0;
                        }else if(a[1].data == 0 && a[0].data == 0){
                            b = 0;
                        }else if(a[1].data/a[0].data >= 100){
                            b = 100;
                        }else{
                            b = ((a[1].data/a[0].data)*100).toFixed(2);
                        }
                        list.push(
                            '<br>&nbsp&nbsp进度占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                textStyle: {
                    color: '#333'
                },
                legend: {
                    selectedMode:false,
                    data: ['实际缴税完成', '数据录入完成'],
                    textStyle: {
                        color: "#333"
                    }
                },
                grid:{
                    x:86,
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
                        name: '元',
                        // min: 0,
                        // max: 16000,
                        // interval: 4000,
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
                        name: '实际缴税完成',
                        type: 'bar',
                        color: '#f5f5f5',
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        data: this.taxPaymentData,
                        z: 0,
                        silent: true,
                        animation: true, // 关闭动画效果
                        barWidth:barWidths,
                    },
                    {
                        name: '数据录入完成',
                        type: 'bar',
                        color:'#6F90CF',
                        data: this.paymentserData,
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        barWidth:barWidths,
                    }
                ]
            };
            optionpaymentss && myChartpaymentss.setOption(optionpaymentss);
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
                    this.datalogistics = [];        // 时间
                    this.logisticsData = [];        // 系统统计
                    this.logstatisticsData = [];    // 实际统计
                    let newdate = [];               // 组装时间数据
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
                        this.logisticsData.push(data.toFixed(2))      // 系统统计
                    }
                    this.api.statistics.countMonActBus(type)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.loading = false;
                            for(let itme in this.datalogistics){
                                let data = res.data.data.logisticsCount[this.datalogistics[itme]]
                                this.logstatisticsData.push(data.toFixed(2))
                            }
                        }
                        this.logissett();
                    }) 
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
                    text: '物流量'
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'cross',
                        crossStyle: {
                            color: '#333'
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        let b;
                        list.push(
                            a[0].axisValue
                        )
                        for (var i = 0; i <a.length; i++) {
                        list.push(
                            '<i style="display: inline-block;width: 10px;height: 10px;background: ' +
                                a[i].color +
                                ';margin-right: 5px;border-radius: 50%;}"></i><span style="width:120px; display:inline-block;">' +
                                a[i].seriesName + 
                                '</span>' + a[i].data + `（吨）`
                            )
                        }
                        if(a[0].data/a[1].data == Infinity || a[0].data/a[1].data == 0){
                            b = 0;
                        }else if(a[0].data == 0 && a[1].data == 0){
                            b = 0;
                        }else if(a[0].data/a[1].data >= 100){
                            b = 100;
                        }else{
                            b = ((a[0].data/a[1].data)*100).toFixed(2);
                        }
                        list.push(
                            '<br>&nbsp&nbsp进度占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                grid:{
                    x:74,
                    y:75,
                    x2:70,
                    y2:20,
                },
                legend: {
                    selectedMode:false,
                    data: ['开票物流完成', '数据录入完成'],
                    textStyle: {
                        color: "#333"
                    }
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
                        axisPointer: {
                            label: {
                                show: true,
                                backgroundColor: '#58afff'
                            },
                        },
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '吨',
                
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
                        name:'数据录入完成',
                        type: 'line',
                        color: '#58afff',
                        // smooth: true,
                        // symbol: 'none',
                        // itemStyle: {
                        //     normal: {
                        //         color: '#58afff'
                        //     }
                        // },
                        // areaStyle: {
                        //     color: '#58afff'
                        // },
                        data:this.logisticsData
                    },
                    {
                        name:'开票物流完成',
                        type: 'line',
                        color: '#BDBDBD',
                        // smooth: true,
                        // symbol: 'none',
                        // itemStyle: {
                        //     normal: {
                        //         color: '#EAF5FF'
                        //     }
                        // },
                        // areaStyle: {
                        //     color: '#f5f5f5'
                        // },
                        data:this.logstatisticsData
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
                    if(res.data.data.length == 0){
                        this.invoiceShow = 2;
                    }else{
                        this.invoiceShow = 1;
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
                        this.chartDom();  // 交易产品份额图表
                    }

                }
            })
        },
        chartDom(){
            var chartDom = document.getElementById('invoicess');
            var myChart = echarts.init(chartDom);
            var option;
            option = {
                // title: {
                //     text: '产品交易份额'
                // },
                tooltip: {
                    trigger: 'item'
                },
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                    left:0,
                    top: 20,
                    bottom: 10,
                },
                noDataLoadingOption: {
                    text: '暂无数据',
                    effect: 'bubble',
                    effectOption: {
                        effect: {
                            n: 0
                        }
                    }
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

        /*
        所得税率
        */
        allTaxrateapi(type){
            this.loading = true;
            this.api.statistics.countEtcTax(type)  
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.allTaxrateTime = [];
                    this.allTaxrateData = [];
                    for(let key in res.data.data){
                        this.allTaxrateTime.push(res.data.data[key].simpleName);
                        let color;
                        if((res.data.data[key].trate*100).toFixed(2) < 1){
                            color = "#EE6666";
                        }else{
                            color = "#FFF5F5";
                        }
                        let obj = {
                            value:(res.data.data[key].trate*100).toFixed(2),
                            itemStyle: {
                                color:color
                            }
                        }
                        this.allTaxrateData.push(obj);
                    }
                    this.allTaxrateburden();  // 所得税率
                }
            })
        },
        allTaxrateburden(){
            document.getElementById("allTaxrate").removeAttribute("_echarts_instance_");
            var allTaxrate = document.getElementById('allTaxrate');
            var ChartallTaxrate = echarts.init(allTaxrate);
            var optionChartallTaxrate;
            optionChartallTaxrate = {
                title: {
                    text: '平台企业所得税率',
                    left: '50%',
                    textAlign: 'center'
                },
                grid:{
                    x:50,
                    y:75,
                    x2:70,
                    y2:50,
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
                xAxis: [
                    {
                        type: 'category',
                        // axisLabel: { interval: 0, rotate: 30 },
                        name: '企业名称',
                        data:  this.allTaxrateTime,
                        axisPointer: {
                            type: 'line'
                        },
                        nameTextStyle:{
                            color:'#333'
                        }
                    }
                ],
                dataZoom: [
                    {
                        type: 'slider',
                        // backgroundColor:"#FCDEDE",
                        // handleColor: '#CF1724',                     // h滑动图标的颜色
                        // fillerColor:"#F5C4C4",                      // 选中范围的填充颜色。
                        // borderColor:"#CF1724",                      // 边框颜色。
                        show: true,                                 // flase直接隐藏图形
                        xAxisIndex: [0],
                        left: '5%',                                 // 滚动条靠左侧的百分比
                        bottom: -5,
                        start: 0,                                   // 滚动条的起始位置
                        end: 100,                                    // 滚动条的截止位置（按比例分割你的柱状图x轴长度）
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '税率',
                        // min: 0,
                        // max: 14,
                        // interval:2,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333'
                            }
                        },
                        axisLabel: {
                            formatter: '{value} %'
                        }
                    }
                ],
                series:[
                    {
                        type: 'bar',
                        barWidth:'30%',
                        itemStyle: {
                            emphasis:{
                                color:"#F5C4C4",    //移入后的颜色
                            }
                        },
                        data:this.allTaxrateData,
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "1%预警线"
                            },
                            data : [{
                                silent:true,             //鼠标悬停事件  true没有，false有
                                lineStyle:{               //警戒线的样式  ，虚实  颜色
                                    type:"dashed",
                                    color:"#CF1724"
                                },
                                name: '警戒线',
                                yAxis: 1
                            }]
                        },
                        // markPoint: {
                        //     data: [
                        //         {type: 'max', name: '最大值'},
                        //         {type: 'min', name: '最小值'}
                        //     ]
                        // },
                        // label: seriesLabel,
                        // axisLabel: {
                        //     formatter: '{value} %'
                        // },
                        label: {
                            normal: {
                                show: true, //是否显现，不显示的话设置成false
                                position: "top", //显示的位置
                                distance: 0, //距离侄子的值 // label要显示的值比如： 20%
                                formatter: function(param) {
                                    return param.value + '%';
                                },
                            }
                        },
                        // markLine: {
                        //     data: [
                        //         {type: 'average', name: '平均值'}
                        //     ]
                        // }
                    },
                ],
            };
            optionChartallTaxrate && ChartallTaxrate.setOption(optionChartallTaxrate);
        },

        /* 
        税负率数据
        */
        burdenapi(type){
            this.loading = true;
            this.api.bill.countPtTax(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.burdenTime = [];
                    this.burdenData = [];
                    for(let key in res.data.data){
                        this.burdenTime.push(key);
                        let color;
                        if((res.data.data[key]*100).toFixed(2) >= 9.5){
                            color = "#EE6666";
                        }else{
                            color = "#FFF5F5";
                        }
                        let obj = {
                            value:(res.data.data[key]*100).toFixed(2),
                            itemStyle: {
                                color:color
                            }
                        }
                        this.burdenData.push(obj);
                    }
                    this.burden();  // 税负率
                }
            })
        },
        burden(){
            document.getElementById("burden").removeAttribute("_echarts_instance_");
            var burden = document.getElementById('burden');
            var myChartburden = echarts.init(burden);
            var optionmyChartburden;
            optionmyChartburden = {
                title: {
                    text: '平台企业增值税税负率',
                    left: '50%',
                    textAlign: 'center'
                },
                grid:{
                    x:50,
                    y:75,
                    x2:70,
                    y2:50,
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
                xAxis: [
                    {
                        type: 'category',
                        // axisLabel: { interval: 0, rotate: 30 },
                        name: '企业名称',
                        data:  this.burdenTime,
                        axisPointer: {
                            type: 'line'
                        },
                        nameTextStyle:{
                            color:'#333'
                        }
                    }
                ],
                dataZoom: [
                    {
                        type: 'slider',
                        // backgroundColor:"#FCDEDE",
                        // handleColor: '#CF1724',                     // h滑动图标的颜色
                        // fillerColor:"#F5C4C4",                      // 选中范围的填充颜色。
                        // borderColor:"#CF1724",                      // 边框颜色。
                        show: true,                                 // flase直接隐藏图形
                        xAxisIndex: [0],
                        left: '5%',                                 // 滚动条靠左侧的百分比
                        bottom: -5,
                        start: 0,                                   // 滚动条的起始位置
                        end: 100,                                    // 滚动条的截止位置（按比例分割你的柱状图x轴长度）
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '税负率',
                        min: 0,
                        max: 14,
                        interval:2,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333'
                            }
                        },
                        axisLabel: {
                            formatter: '{value} %'
                        }
                    }
                ],
                series:[
                    {
                        type: 'bar',
                        barWidth:'30%',
                        itemStyle: {
                            emphasis:{
                                color:"#F5C4C4",    //移入后的颜色
                            }
                        },
                        data:this.burdenData,
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "9.5%预警线"
                            },
                            data : [{
                                silent:true,             //鼠标悬停事件  true没有，false有
                                lineStyle:{               //警戒线的样式  ，虚实  颜色
                                    type:"dashed",
                                    color:"#CF1724"
                                },
                                name: '警戒线',
                                yAxis: 9.5
                            }]
                        },
                        // markPoint: {
                        //     data: [
                        //         {type: 'max', name: '最大值'},
                        //         {type: 'min', name: '最小值'}
                        //     ]
                        // },
                        // label: seriesLabel,
                        // axisLabel: {
                        //     formatter: '{value} %'
                        // },
                        label: {
                            normal: {
                                show: true, //是否显现，不显示的话设置成false
                                position: "top", //显示的位置
                                distance: 0, //距离侄子的值 // label要显示的值比如： 20%
                                formatter: function(param) {
                                    return param.value + '%';
                                },
                            }
                        },
                        // markLine: {
                        //     data: [
                        //         {type: 'average', name: '平均值'}
                        //     ]
                        // }
                    },
                ],
            };
            optionmyChartburden && myChartburden.setOption(optionmyChartburden);
        },

        /* 
        票据税额进项比数据
        */
        countBillbillnote(type){
            this.loading = true;
            this.api.bill.count(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    if(res.data.data == '' || res.data.data.length <= 1){
                        this.sumTaxAmountbillnote = 0;
                        this.arrTaxAmountbillnote = 0;
                        this.chartDombillnote();  // 税额
                    }else{
                        for(let i=0;i<res.data.data.length;i++){
                            if(res.data.data[i].type == 1){
                                this.sumTaxAmountbillnote =  res.data.data[i].sumTaxAmount *1/10000; // 进项税额
                            }else if(res.data.data[i].type == 2){
                                this.arrTaxAmountbillnote = res.data.data[i].sumTaxAmount *1/10000; // 销项税额 
                            }
                        }
                        this.chartDombillnote();  // 税额
                    }
                }
            })
        },
        chartDombillnote(){
            document.getElementById("billnote").removeAttribute("_echarts_instance_");
            var chartDomss = document.getElementById('billnote');
            var myChartss = echarts.init(chartDomss);
            var optionss;
            var num;
            if(this.sumTaxAmountbillnote*1 == 0 || this.arrTaxAmountbillnote*1 == 0){
                num = 0;
            }else{
                let newnum = this.sumTaxAmountbillnote*1 / this.arrTaxAmountbillnote*1;
                num = (newnum*100).toFixed(2);
            }
            optionss = {
                series: [{
                    type: 'gauge',
                    startAngle: 180,
                    endAngle: 0,
                    min: 0,
                    max: 100,
                    splitNumber: 10,
                    axisLine: {
                        lineStyle: {
                            width: 10,
                            color: [
                                [0.3, '#EE6666'],
                                [1, '#FCDEDE'],
                        
                            ]
                        }
                    },
                    pointer: {
                        icon: 'path://M12.8,0.7l12,40.1H0.7L12.8,0.7z',
                        length: '25%',
                        width: 8,
                        offsetCenter: [0, '-60%'],
                        itemStyle: {
                            color: 'auto'
                        }
                    },
                    axisTick: {
                        length: 6,
                        lineStyle: {
                            color: 'auto',
                            width: 2
                        }
                    },
                    splitLine: {
                        length: 10,
                        lineStyle: {
                            color: 'auto',
                            width: 2
                        }
                    },
                    axisLabel: {
                        color: '#EE6666',
                        // fontSize: 14,
                        distance:-55,
                        formatter: function (value) {
                            if (value === 0) {
                                return '不足';
                            }
                            else if (value === 100) {
                                return '正常';
                            }
                        }
                    },
                    title: {
                        offsetCenter: [0, '-30%'],
                        // fontSize: 24
                    },
                    detail: {
                        // fontSize: 40,
                        offsetCenter: [0, '-7%'],
                        valueAnimation: true,
                        formatter: function (value) {
                            return value + '%';
                        },
                        color: 'auto'
                    },
                    data: [{
                        value:num,
                        name: '进项占比'
                    }]
                }]
            };
            optionss && myChartss.setOption(optionss);
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
            // width:45%;
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
        // .Purchasing{
        //     float: left;
        //     width:30%;
        //     height:321px;
        //     box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        //     background-color: #fff;
        //     margin:0 0 0.5% 0.5%;
        //     padding:24px 0 24px 48px;
        //     box-sizing: border-box;
        //     // overflow: hidden;
        //     .tity{
        //         font-size: 18px;
        //         font-weight: 600;
        //         color: #333;
        //         line-height: 32px;
        //         margin-left:8px;
        //     }
        //     span{
        //         line-height: 26px;
        //     }
        //     img{
        //         width: 32px;
        //         height: 32px;
        //         vertical-align: middle;
        //     }
        //     .Purchasing-conent{
        //         width: 100%;
        //         li{
        //             width: 60%;
        //             height: 100%;
        //             float: left;
        //             padding-top:20px;
        //             box-sizing: border-box;
        //             p{
        //                 font-size:16px;
        //                 line-height: 24px;
        //                 color: #999;
        //             }
        //             span{
        //                 font-size: 16px;
        //                 color: #333;
        //                 vertical-align: middle;
        //                 margin-left: 8px;
        //             }
        //             h5{
        //                 display: inline-block;
        //                 font-size: 32px;
        //                 color: #333;
        //                 line-height: 64px; 
        //                 font-weight: 600;
        //             }
        //         }
        //     }
        // }
        // 统计进度
        .quantity{
            float: left;
            // width:99%;
            // height:321px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;
            margin:0 0.5% 0.5% 0.5%;
            padding:24px 0 0 28px;
            box-sizing: border-box;
                            // overflow: hidden;
                // overflow-x: scroll;
            // overflow: hidden;
            .tity{
                font-size: 18px;
                font-weight: 600;
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
                display: flex;
                flex-wrap: nowrap;
                justify-content: first baseline;
                li{
                    width:33.3%;
                    height: 100%;
                    float: left;
                    padding:6% 20px;
                    box-sizing: border-box;
                    text-align: left;
                    p{
                        font-size:16px;
                        line-height: 24px;
                        color: #999;
                    }
                    span{
                        font-size: 16px;
                        color: #333;
                        vertical-align: middle;
                    }
                    h5{
                        display: inline-block;
                        font-size: 32px;
                        color: #333;
                        line-height: 64px; 
                        font-weight: 600;
                    }
                }
            }
        }
        // 客户统计工作
        .countMonData{
            float: left;
            width:99%;
            height:380px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;   
            margin:0 0.5% 0.5%;
            position: relative;
            #countMonData{
                float: left;
                width: 100%;
                height: 100%;
                padding: 14px;
                box-sizing: border-box;
            }
        }
        // 网络货运占比
        .NetworkFreight{
            float: left;
            width:99%;
            height:380px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;   
            margin:0 0.5% 0.5%;
            position: relative;
            #burdenLogins{
                float: left;
                width: 100%;
                height: 100%;
                padding: 14px;
                box-sizing: border-box;
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
            #paymentTaxes{
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
            margin:0 0.5% 0.5%;
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
            position: relative;
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
            #invoicess{
                width: 100%;
                height:90%;
            }
            .invoiceBox{
                width: 100px;
                height: 20px;
                color: #999;
                text-align: center;
                line-height: 20px;
                position: absolute;
                top: 0;
                left: 0;
                bottom: 0;
                right: 0;
                margin: auto;
            }
        }
        // 税负率
        .burden{
            float: left;
            width:99%;
            height:400px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;  
            margin:0 0.5% 0.5%;
            margin-right: 0;
            padding: 14px;
            box-sizing: border-box;
            #burden{
                width:100%;
                height:100%;
            }
        }
        // 所得税率
        .allTaxrate{
            float: left;
            width:99%;
            height:400px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;  
            margin:0 0.5% 0.5%;
            margin-right: 0;
            padding: 14px;
            box-sizing: border-box;
            #allTaxrate{
                width:100%;
                height:100%;
            }
        }
        // 票据进项比
        .billnote{
            float: left;
            width:29.2%;
            height:400px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
            background-color: #fff;  
            margin:0 0.5% 0.5%;
            margin-right: 0;
            padding: 14px;
            box-sizing: border-box;
            position: relative;
            overflow: hidden;
            span{
                font-size: 18px;
                color: #333;
                font-weight: 600;
                line-height: 32px;
            }
            #billnote{
                padding-top:10%;
                width: 100%;
                height:100%;
            }
        }
    }
    // 内容结束
</style>