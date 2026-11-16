<template>
    <!-- 第二屏LED看板 -->
    <div class="FirstScreen">
        <div class="rightGo" @click="rightGos">
            <i class="el-icon-arrow-right"></i>
        </div>
        <div class="leftGo" @click="leftGos">
            <i class="el-icon-arrow-left"></i>
        </div>
        <div class="header">
            <p>煤炭大宗数字交易平台</p> 
            <div class="TimesAll">
                <span>{{FirstScreenTime}}</span>
                <div class="buttonbox" @click="BackOff">
                    <i class="iconfont el-icon-alireturn"></i>
                </div>
            </div>
            <div class="TimesAllleft" @click="BackTime">
                <span v-if="backcolor">自动循环 </span>
                <span v-if="!backcolor" style="color:rgb(1,240,244);">自动循环 </span>
                <i v-if="backcolor" class="el-icon-turn-off"></i>
                <i v-if="!backcolor" class="el-icon-open" style="color:rgb(1,240,244);"></i>
            </div>
        </div>
        <div class="nav">
            <!-- 数据看板 -->
            <div class="section">
                <ul>
                    <li>
                        <p>合同份数累计(份)</p>
                        <span><countTo :startVal='0' :endVal='countSumContractData' :duration='2000' :useEasing="true" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>物流量累计(万吨)</p>
                        <span><countTo :startVal='0' :endVal='logisticsCount' :duration='4000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>结算金额累计(亿元)</p>
                        <span><countTo :startVal='0' :endVal='getPtSumPricesData' :duration='2000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>注册企业数量</p>
                        <span><countTo :startVal='0' :endVal='chasingTonnage' :duration='2000' :useEasing="true" :autoplay="true"></countTo></span>
                    </li> 
                </ul>
                <ul>
                    <li>
                        <p>业务量金额累计(亿元)</p>
                        <span><countTo :startVal='0' :endVal='newbusCount' :duration='3000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>缴税金额累计(亿元)</p>
                        <span><countTo :startVal='0' :endVal='rateBills' :duration='1000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>资金流水累计(亿元)</p>
                        <span><countTo :startVal='0' :endVal='Capitalflow' :duration='1500' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                    <li>
                        <p>本年度业务量金额累计(亿元)</p>
                        <span><countTo :startVal='0' :endVal='countYearBusesData' :duration='1500' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                    </li>
                </ul>
            </div>
            <!-- 物流 -->
            <div class="section-left">
                <div id="paymentTaxes"></div>
            </div>
            <!-- 税负率 -->
            <div class="conent">
                <div id="burden"></div>
                <div class="burdenTimes">
                    <el-tooltip class="item" effect="dark" content="(增值税税额 / 不含税金额)*100% （数据来源：缴税明细）" placement="top" style="margin:0 2%;">
                        <i class="el-icon-question" style="color:#fff;font-size: 16px;"></i>
                    </el-tooltip>
                    <span>累计税负率(年)：</span>
                    <el-dropdown trigger="click" @command="burdenTimesClick">
                    <span class="el-dropdown-link">
                         {{burValue}}<i class="el-icon-arrow-down el-icon--right"></i>
                    </span>
                    <el-dropdown-menu slot="dropdown" style="background-color:rgba(0,0,0,.3);border:1px solid rgba(0,0,0,.3);">
                        <el-dropdown-item v-for="(item,key) in burdenTimesoptions" :key="key" :command="item.value" style="color:rgb(6,216,206);" class="burdenTimessize">{{item.label}}</el-dropdown-item>
                    </el-dropdown-menu>
                    </el-dropdown>
                </div>
            </div>
            <!-- 企业列表信息 -->
            <div class="footer">
                <dl id="sectionBox-top">
                    <dd style="width:25%;">企业名称</dd>
                    <dd style="width:10%;">企业简称</dd>
                    <dd style="width:10%;">注册时间</dd>
                    <dd style="width:20%;">税号</dd>
                    <dd style="width:10%;">上户时间</dd>
                    <dd style="width:15%;">注册资本</dd>
                </dl>
                <div id="sectionBox">
                    <ul id="con1" ref="con1" :class="{anim:animate==true}">
                        <li v-for='item in items'>
                            <div class="sectionsize" style="width:25%;">{{item.enterpriseName}}</div>
                            <dd class="sectionsize" style="width:10%;">{{item.simpleName}}</dd>
                            <div class="sectionsize" style="width:10%;">{{item.regTime}}</div>
                            <div class="sectionsize" style="width:20%;">{{item.enginenocopy}}</div>
                            <div class="sectionsize" style="width:10%;">{{item.shTime}}</div>
                            <div class="sectionsize" style="width:15%;">{{item.regMoney}}</div>
                        </li>
                    </ul>
                </div>
            </div>
            <div class="footer-right">
                <div class="sectionBox-top"> 洗煤厂企业</div>
                <div id="sectionBoxone">
                    <ul id="con2" ref="con2" :class="{animone:animateone==true}">
                        <li v-for='item in itemsone'>
                            <img :src="item.img" alt="">
                            <p>{{item.name}}</p>
                            <p>{{item.nameone}}</p>
                            <p>{{item.nametwo}}</p>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import countTo from 'vue-count-to';
var echarts = require('echarts');
export default{
    components: { countTo },
    data(){
        return {
            routerboday:null,
            invoiceDate:{},             // 图表时间
            buildinvoiceDate:{},        // 税负率时间
            FirstScreenTime:'',
            animate:false,
            items:[],                   // 企业信息
            animateone:false,
            backcolor:false,
            // 洗煤厂信息
            itemsone:[
                {
                    img:require('../../assets/xm8.jpg'),
                    name:'晋燃煤炭有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 主焦 / 瘦焦 / 肥煤',
                },
                {
                    img:require('../../assets/xm7.jpg'),
                    name:'兴浩洗煤有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 主焦 / 瘦焦',
                },
                {
                    img:require('../../assets/xm3.jpg'),
                    name:'东海煤化有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 主焦 / 肥煤',
                },
                {
                    img:require('../../assets/xm5.jpg'),
                    name:'昌鑫洗煤有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 精煤',
                },
                {
                    img:require('../../assets/xm1.jpg'),
                    name:'佳晟煤化有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 原煤',
                },
                {
                    img:require('../../assets/xm2.jpg'),
                    name:'晋兴达洗煤有限责任公司',
                    nameone:'所属地: 山西晋中',
                    nametwo:'主营: 原煤',
                },
            ],

            countSumContractData:0,     // 合同份数累计
            logisticsCount:0,           // 物流量累计 
            getPtSumPricesData:0,       // 结算金额累计
            newbusCount:0,              // 业务量金额累计
            rateBills:0,                // 缴税金额累计
            Capitalflow:0,              // 资金流水累计
            chasingTonnage:0,           // 注册企业数量
            countYearBusesData:0,       // 统计本年度业务量金额累计

            // 物流量
            paymentdatapayment:[],          // 时间
            paymentserData:[],              // 物流量

            // 税负率
            option: null,
            number: 0, // 播放所在下标
            timer: null,
            burdenTime:[],              // 公司名称
            burdenData:[],              // 税负率
            burdenTimesoptions:[
                {
                    value:1,
                    label:'2019 - 至今'
                }
            ],      // 税负率时间筛选
            burValue:'',
        }
    },
    created(){
        // 跳转
        this.routerboday = setInterval(this.routername,60000);
        // 获取当前时间
        let date = new Date();
        let year = date.getFullYear();
        let month = date.getMonth() + 1;
        let day = date.getDate();
        let Years = (year - 2019) + 1;
        this.burValue = year + '年'
        for(let i=0;i<Years;i++){
            this.burdenTimesoptions.push({
                value:(year - i),
                label:year - i + '年'
            })
        }
        if (month < 10) {
            month = "0" + month;
        }
        if (day < 10) {
            day = "0" + day;
        }
        if(month == 0){
            year = year - 1
            month = 12;
            var lastDay= new Date(year,month,0).getDate();   // 当月最后一天
            day = lastDay;
        }
        let nowDate = year + "-" + month + "-" + day;
        // 将半年的时间单位换算成毫秒
        var halfYear = 365 / 2 * 24 * 3600 * 1000;
        var pastResult = date - halfYear;  // 半年前的时间（毫秒单位）
        // 日期函数，定义起点为半年前
        var pastDate = new Date(pastResult),
            pastYear = pastDate.getFullYear(),
            pastMonth = pastDate.getMonth() + 1,
            pastDay = 1;
        if (pastMonth < 10) {
            pastMonth = "0" + pastMonth;
        }
        if (pastDay < 10) {
            pastDay = "0" + pastDay;
        }
        let bnDate = pastYear + '-' + pastMonth + '-' + pastDay;
        this.invoiceDate = {
            xzqhCode:sessionStorage.getItem('xzqhCode'),
            startTime:bnDate,
            endTime:nowDate
        };
        this.buildinvoiceDate = {
            xzqhCode:sessionStorage.getItem('xzqhCode'),
            startTime:'2019-10-01',
            endTime:nowDate
        };
        setInterval(this.formatDate, 1000);
        // 企业信息列表
        var listByPageData = {
            selEtcType:2,
            isAuthentication:true,
            searchStr:"",
            page:{
                current:1,
                size:100
            }
        }
        // 渲染表格
          this.api.enterData.all(listByPageData)
          .then(res=>{
            if(res.data.code == 200){
                if(res.data.data.records != ""){
                    for(let i=0;i<res.data.data.records.length;i++){
                        let str = res.data.data.records[i].rateNumber;
                        let enginenocopy = str.replace(/(^(?:(?![IOZSV])[\dA-Z]){2}\d{6}(?:(?![IOZSV])[\dA-Z]){10}$)|(^\d{15}$)/, "\$1****\$2");
                        res.data.data.records[i].enginenocopy = enginenocopy;
                    }
                  this.items = res.data.data.records
                  console.log(this.items,'items');
                }
            }
          })
        // this.api.statistics.enterData({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
        // .then(res=>{
        //     if(res.data.code == 200){
        //         for(let i=0;i<res.data.data.length;i++){
        //             let str = res.data.data[i].rateNumber;
        //             let enginenocopy = str.replace(/(^(?:(?![IOZSV])[\dA-Z]){2}\d{6}(?:(?![IOZSV])[\dA-Z]){10}$)|(^\d{15}$)/, "\$1****\$2");
        //             res.data.data[i].enginenocopy = enginenocopy;
        //         }
        //         this.items = res.data.data; // 企业名称
        //         console.log(this.items,'items');
        //     }
        // })
        setInterval(this.scroll,2000) // 在钩子函数中调用我在method 里面写的scroll()方法
        setInterval(this.scrollone,3000) // 在钩子函数中调用我在method 里面写的scroll()方法
        this.getPtSumPrices();                      // 结算金额累计
        setInterval(this.getPtSumPrices, 10000);
        this.countSumContracts();                   // 合同份数
        setInterval(this.countSumContracts, 15000);
        this.countYearBusess();                     // 资金流水累计
        setInterval(this.countYearBusess, 200000);
        this.logisticsBill()                        // 物流量累计
        setInterval(this.logisticsBill, 25000);
        this.taxartbarbox();                        // 注册企业数
        setInterval(this.taxartbarbox, 30000);
        this.oneYearBusess();             // 统计本年度业务量金额累计
        setInterval(this.oneYearBusess, 200000);
        let aobj = {
            xzqhCode:sessionStorage.getItem('xzqhCode'),
            startTime:year + '-01-01',
            endTime:year + '-12-31',
        }
        this.burdenapi(aobj);                           // 税负率
    },
    mounted() {
        this.paymentCapitalTaxes();                 // 物流量
    },
    methods:{
        burdenTimesClick(value){
            if(value == 1){
                this.burValue = "2019 - 至今";
                this.burdenapi(this.buildinvoiceDate);                           // 税负率
            }else{
                this.burValue = value + " 年";
                let obj = {
                    xzqhCode:sessionStorage.getItem('xzqhCode'),
                    startTime:value + '-01-01',
                    endTime:value + '-12-31',
                }
                this.burdenapi(obj);                           // 税负率
            }
        },
        rightGos(){
            this.$router.push({name:'Thirdscreen'});
        },
        leftGos(){
            this.$router.push({name:'Firstscreen'});
        },
        // 跳转
        routername(){
            this.$router.push({name:'Thirdscreen'})
        },
        BackTime(){
            this.backcolor = !this.backcolor;
            if(this.backcolor){
                clearInterval(this.routerboday);
            }else{
                this.routerboday = setInterval(this.routername,60000);
            }
            
        },
        // 退出大屏3

        BackOff(){
            this.$router.push({name:'Dpmain'})
            
        },
        // 处理时间排序函数
        compare(property){
            return function(a,b){
                var value1 = a[property];
                var value2 = b[property];
                return value1 - value2;
            }
        },
        // 时间屏幕
        formatDate() {
            let date = new Date();
            let year = date.getFullYear(); // 年
            let month = date.getMonth() + 1; // 月
            month = month < 10 ? "0" + month : month; // 如果只有一位，则前面补零
            let day = date.getDate(); // 日
            let hour = date.getHours(); // 时
            hour = hour < 10 ? "0" + hour : hour; // 如果只有一位，则前面补零
            let minute = date.getMinutes(); // 分
            minute = minute < 10 ? "0" + minute : minute; // 如果只有一位，则前面补零
            let second = date.getSeconds(); // 秒
            second = second < 10 ? "0" + second : second; // 如果只有一位，则前面补零
            let week = new Date().getDay();
            let weeks = ["日","一","二","三","四","五","六"];
            let getWeek = "星期" + weeks[week];
            this.FirstScreenTime = `${year}/${month}/${day} ${hour}:${minute}:${second} ${getWeek}`;
        },

        /*
        注册企业数量
        */
        taxartbarbox(){
            this.chasingTonnage = 0;
            this.api.statistics.count({xzqhCode:sessionStorage.getItem('xzqhCode')})
            .then(res=>{
                if(res.data.code == 200){
                    this.chasingTonnage = res.data.data.sumCount*1;           // 注册企业数量
                }
            })
        },

        // 统计本年度业务量金额累计
        oneYearBusess(){
            this.countYearBusesData = 0;
            this.api.bill.countYearBuses()
            .then(res=>{
                if(res.data.code == 200){
                    this.countYearBusesData =  (res.data.data.baseMoney/100000000).toFixed(2)*1;
                }
            })
        },

        // 结算金额累计
        getPtSumPrices(){
            this.getPtSumPricesData = 0;
            let xzqhCodes = sessionStorage.getItem('xzqhCode');
            this.api.capital.getPtSumPrice({xzqhCode:xzqhCodes})
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data){
                        this.getPtSumPricesData =  (res.data.data.sumPrice/100000000).toFixed(2)*1;
                    }else{
                        this.getPtSumPricesData = 0;
                    }
                }
            })
        },

        // 平台合同量
        countSumContracts(){
            this.countSumContractData = 0;
            let xzqhCodes = sessionStorage.getItem('xzqhCode');
            this.api.listByPage.countSumContract({xzqhCode:xzqhCodes})
            .then(res=>{
                if(res.data.code == 200){
                    this.countSumContractData =  res.data.data.count*1;
                }
            })
        },

        // 资金流水累计
        countYearBusess(){
            this.Capitalflow = 0;
            let xzqhCodes = sessionStorage.getItem('xzqhCode');
            this.api.chainCap.getPtSumPrice({xzqhCode:xzqhCodes})
            .then(res=>{
                if(res.data.code == 200){
                    this.Capitalflow =  (res.data.data.sumMoney/100000000).toFixed(2)*1;
                }
            })
        },

        // 物流量 业务量金额累计
        logisticsBill(){
            this.logisticsCount = 0;
            this.newbusCount = 0;
            this.rateBills = 0;
            this.api.statistics.countActBus({xzqhCode:sessionStorage.getItem('xzqhCode')})
            .then(res=>{
                if(res.data.code == 200){
                    this.logisticsCount = res.data.data.logisticsCount/10000;         // 物流量累计
                    this.newbusCount = res.data.data.busCount/100000000;            // 业务量累计
                    this.rateBills = res.data.data.taxCount/100000000;              // 缴税金额累计
                }
            })
        },

        // 企业信息滚动展示
        scroll(){
            let con1 = this.$refs.con1;
            con1.style.marginTop = '-32px';
            this.animate=!this.animate;
            let that = this;
            setTimeout(function(){
                that.items.push(that.items[0]);
                that.items.shift();
                con1.style.marginTop='0px';
                that.animate=!that.animate;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },500)
        },

        // 洗煤厂信息滚动展示
        scrollone(){
            let con2 = this.$refs.con2;
            con2.style.marginLeft='-33%';
            this.animateone=!this.animateone;
            let that = this;
            setTimeout(function(){
                that.itemsone.push(that.itemsone[0]);
                that.itemsone.shift();
                con2.style.marginLeft='0px';
                that.animateone=!that.animateone;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },1000)
        },

        /* 
        物流量数据
        */
        paymentCapitalTaxes(){
            this.api.statistics.countMonActBus(this.invoiceDate)
            .then(res=>{
                if(res.data.code == 200){
                    this.datalogistics = [];        // 时间
                    this.logstatisticsData = [];    // 实际统计
                    let newdate = [];               // 组装时间数据
                    for(let key in res.data.data.logisticsCount){
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
                        let data = res.data.data.logisticsCount[this.datalogistics[itme]]/10000
                        this.logstatisticsData.push(data.toFixed(2))
                    }
                }
                this.paymentTaxes();
                
            }) 
        },
        paymentTaxes(){
            // 物流量图表
            var chartDompaymentss = document.getElementById('paymentTaxes');
            var myChartpaymentss = echarts.init(chartDompaymentss);
            var optionpaymentss;
            optionpaymentss = {
                textStyle: {
                    color: '#fff'
                },
                legend: {
                    textStyle: {
                        color: "#fff"
                    }
                },
                grid:{
                    x:60,
                    y:60,
                    x2:60,
                    y2:50,
                },
                xAxis: [
                    {
                        type: 'category',
                        name: '时间',
                        data:this.datalogistics,
                        axisPointer: {
                            type: 'line'
                        },
                        nameTextStyle:{
                            color:'#fff'
                        }
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '物流（万吨）',
                        // min: 0,
                        // max: 2000,
                        // interval: 500,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#fff'
                            }
                        },
                        nameTextStyle:{
                            color:'#fff'
                        },
                        axisPointer: {
                            label: {
                                show: true,
                                backgroundColor: '#fff'
                            },
                        },
                        splitLine: {show: false},
                    }
                ],
                series: [
                    {
                        name: "",
                        type: "bar",
                        barWidth: 20,
                        itemStyle: {
                            normal: {
                                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                                        offset: 0,
                                        color: "#07ecd9"
                                    },
                                    {
                                        offset: 1,
                                        color: "#034881"
                                    }
                                ])
                            }
                        },
                        data: this.logstatisticsData,
                        z: 10,
                        zlevel: 0,
                        "label": {
                            "show": true,
                            "position": "top",
                            "distance": 8,
                            fontSize:16,
                            "color": "#01fff4"
                        }
                    },
                    {
                        // 分隔
                        type: "pictorialBar",
                        itemStyle: {
                            normal:{
                                color:"#0F375F"
                            }
                        },
                        symbolRepeat: "fixed",
                        symbolMargin: 6,
                        symbol: "rect",
                        symbolClip: true,
                        symbolSize: [20, 2],
                        symbolPosition: "start",
                        symbolOffset: [0, -1],
                        data: this.logstatisticsData,
                        width: 25,
                        z: 0,
                        zlevel: 1,
                    },
                ]
            };
            optionpaymentss && myChartpaymentss.setOption(optionpaymentss);
        },

        /* 
        税负率数据
        */
        burdenapi(Time){
            this.api.bill.countPtTax(Time)
            .then(res=>{
                if(res.data.code == 200){
                    this.burdenTime = [];
                    this.burdenData = [];
                    for(let key in res.data.data){
                        this.burdenTime.push(key);
                        let color = {};
                        if((res.data.data[key]*100).toFixed(2)*1 < 9.5){
                            color = {
                                type: 'linear',
                                x: 0,
                                y: 0,
                                x2: 0,
                                y2: 1,
                                colorStops: [
                                    {
                                        offset: 0,
                                        color: 'rgba(255, 228, 228, 1)',     // 0% 处的颜色
                                    },
                                    {
                                        offset: 1,
                                        color: 'rgba(255, 228, 228, 0.5)',      // 100% 处的颜色
                                    },
                                ],
                                globalCoord: false,                     // 缺省为 false
                            }                                          //渐变颜色;
                        }else if(9.5 <= (res.data.data[key]*100).toFixed(2)*1 && (res.data.data[key]*100).toFixed(2)*1 < 10){
                            color = {
                                type: 'linear',
                                x: 0,
                                y: 0,
                                x2: 0,
                                y2: 1,
                                colorStops: [
                                    {
                                        offset: 0,
                                        color: 'rgba(238, 151, 102, 1)',     // 0% 处的颜色
                                    },
                                    {
                                        offset: 1,
                                        color: 'rgba(238, 151, 102, .5)',      // 100% 处的颜色
                                    },
                                ],
                                globalCoord: false,                     // 缺省为 false
                            } 
                        }else if((res.data.data[key]*100).toFixed(2)*1 >= 10){
                            color = {
                                type: 'linear',
                                x: 0,
                                y: 0,
                                x2: 0,
                                y2: 1,
                                colorStops: [
                                    {
                                        offset: 0,
                                        color: 'rgba(238, 102, 102, 1)',     // 0% 处的颜色
                                    },
                                    {
                                        offset: 1,
                                        color: 'rgba(238, 102, 102, 0.5)',      // 100% 处的颜色
                                    },
                                ],
                                globalCoord: false,                     // 缺省为 false
                            }                                          //渐变颜色;
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
            let chartsmyChart = echarts.init(document.getElementById('burden'));
            this.option = {
                title: {
                    text: '企业增值税税负率 *',
                    top: '6%',
                    left:'50%',
                    textAlign: 'center',
                    textStyle:{
                        color:'rgba(255, 228, 228, 0.9)',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(238, 102, 102, .4)',
                    borderWidth:0,
                    textStyle:{
                        color:'#fff',
                    },
                    confine:true,
                    show: true,
                    extraCssText:{width:"10px"},
                    axisPointer: { 
                        type: 'shadow',
                        shadowStyle: {
                            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                                offset: 1,
                                color: 'rgba(238, 102, 102, .6)'
                            },{
                                offset: 0,
                                color: 'rgba(255, 255, 255, 0)'
                            }])
                        }
                    },
                    formatter: function(a) {
                        let list = []
                        let listItem = '';
                        list.push(
                            a[0].axisValue
                        )
                        list.push(
                            '&nbsp&nbsp税负率：' +
                            a[0].value + `%`
                        )
                        listItem = list.join('<br>')
                        return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                grid:{
                    x:65,
                    y:50,
                    x2:80,
                    y2:50,
                },
                xAxis: [
                    {
                        type: 'category',
                        name: '企业',
                        data:  this.burdenTime,
                        axisPointer: {
                            type: 'line'
                        },
                        nameTextStyle:{
                            color:'#fff'
                        },
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#fff',
                            },
                        },
                        nameTextStyle: {
                            color: '#fff',
                        },
                        axisPointer: {
                            label: {
                                show: true,
                                backgroundColor: 'rgba(238, 102, 102, .6)',
                            },
                        },
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '',
                        min: 0,
                        max: 16,
                        interval:4,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#fff',
                            },
                        },
                        nameTextStyle: {
                            color: '#fff',
                        },
                        axisPointer: {
                            label: {
                                show: true,
                                backgroundColor: '#fff',
                            },
                        },
                        axisLabel: {
                            formatter: '{value} %'
                        },
                        splitLine: { show: false },
                    }
                ],
                series:[
                    {
                        type: 'bar',
                        barWidth:'30%',
                        data:this.burdenData,
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                color:'#fff',
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "9.5%"
                            },
                            data : [{
                                silent:false,             //鼠标悬停事件  true没有，false有
                                lineStyle:{               //警戒线的样式  ，虚实  颜色
                                    type:"dashed",
                                    color:"rgba(252, 218, 38, 1)"
                                },
                                name: '警戒线',
                                yAxis: 9.5
                            }]
                        },
                        label: {
                            normal: {
                                show: true, //是否显现，不显示的话设置成false
                                color:'#fff',
                                position: "top", //显示的位置
                                distance: 0, //距离侄子的值 // label要显示的值比如： 20%
                                formatter: function(param) {
                                    return param.value + '%';
                                },
                            }
                        },
                    },
                ],
            };
            this.option && chartsmyChart.setOption(this.option);
            window.addEventListener('resize', () => {
                chartsmyChart.resize();
            });
            this.timer = setInterval(() => {
                chartsmyChart.dispatchAction({
                    type: 'showTip',
                    seriesIndex: 0,
                    dataIndex: this.number
                });
                this.number++;
                if (this.number > this.burdenTime.length) {
                    this.number = 0;
                }
            }, 2000);
        },
    },
    // 销毁定时器
    beforeDestroy() {
        clearInterval(this.routerboday);
        clearInterval(this.timer);
        if (this.formatDate || this.getPtSumPrices || this.countSumContracts || this.countYearBusess || this.logisticsBill || this.oneYearBusess || this.taxartbarbox) {
            clearInterval(this.formatDate); // 在Vue实例销毁前，清除时间定时器
            clearInterval(this.getPtSumPrices); // 在Vue实例销毁前，清除时间定时器
            clearInterval(this.countSumContracts); 
            clearInterval(this.countYearBusess); 
            clearInterval(this.oneYearBusess); 
            clearInterval(this.taxartbarbox); 
        }
    }
}
</script>
<style scoped lang="less">
.FirstScreen{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    padding:0.5%;
    box-sizing: border-box;
    background-image: url(../../assets/lan.jpg);
    background-size: cover;
    background-repeat: no-repeat;
    position: relative;
    .rightGo{
        width: 50px;
        height: 60px;
        box-sizing: border-box;
        background: #333;
        opacity: .4;
        position: absolute;
        right: 2%;
        top: 50%;
        z-index: 999;
        text-align: center;
        cursor: pointer;
        i{
            margin-right: 4px;
            color: #fff;
            font-size: 60px;
            line-height: 60px;
        }
    }
    .leftGo{
        width: 50px;
        height: 60px;
        box-sizing: border-box;
        background: #333;
        opacity: .4;
        position: absolute;
        left: 2%;
        top: 50%;
        z-index: 999;
        text-align: center;
        cursor: pointer;
        i{
            margin-right: 4px;
            color: #fff;
            font-size: 60px;
            line-height: 60px;
        }
    }
    // 头部
    .header{
        width: 100%;
        height: 66px;
        background-image: url(../../assets/titles.png);
        background-size:100% 132px;
        background-repeat: no-repeat;
        text-align: center;
        position: relative;
        @font-face {
            font-family: aa;
            src: url(../../assets/icon/FZMeiHei-M07.ttf);
        }
        p{
            font-family: aa;
            font-size: 30px;
            line-height: 58px;
            color: #fff;
            font-weight: 500;
        }
        .TimesAllleft{
            position: absolute;
            top: 40%;
            left: 1%;
            text-align: left;
            color: #fff;
            cursor: pointer;
            i{
                font-size: 32px;
                line-height: 40px;
                vertical-align: middle;
            }
            span{
                font-size: 18px;
                font-weight: 600;
                line-height: 40px;
                vertical-align: middle;
            }
        }
        .TimesAll{
            position: absolute;
            top: 40%;
            right: 0;
            text-align: left;
            span{
                font-size: 20px;
                color: #fff;
            }
            .buttonbox{
                width: 60px;
                padding-right:5%;
                box-sizing: border-box;
                float: right;
                text-align: right;
                font-size: 22px;
                color: #fff;
                cursor: pointer;
            }
        }
    }
    // 内容
    .nav{
        width: 100%;
        height: calc(100% - 66px);
        display: flex;
        flex-wrap:wrap;
        // 上半部分内容
        .section{
            width: 66.9%;
            height: 27.5%;
            margin: 0.5% 0.1% 0.2% 0;
            background-image: url(../../assets/k1.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 0 1% 1%;
            box-sizing: border-box;
            ul{
                width: 100%;
                height: 50%;
                display: flex;
                flex-wrap: nowrap;
                li{
                    width: 25%;
                    height: 100%;
                    padding: 2% 0 0 2%;
                    box-sizing: border-box;
                    p{
                        font-size: 14px;
                        line-height: 28px;
                        color: #fff;
                    }
                    span{
                        font-size: 38px;
                        line-height: 58px;
                        color: #fff;
                    }
                }
            }
        }
        .section-left{
            width: 32.9%;
            height: 27.5%;
            margin: 0.5% 0.1% 0.2% 0;
            background-image: url(../../assets/k2.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            #paymentTaxes{
                width: 100%;
                height: 100%;
            }
            // background: cyan;
        }
        .conent{
            width: 99.9%;
            height: 36%;
            margin:0 0.1% 0.2% 0;
            background-image: url(../../assets/k3.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            position: relative;
            // background: darkgoldenrod;
            #burden{
                width: 100%;
                height: 100%;
            }
            .burdenTimes{
                width: 16%;
                height: 40px;
                // background: burlywood;
                position: absolute;
                top: 4%;
                right:3%;
                cursor:pointer;
                span{
                    font-size: 16px;
                    line-height:40px;
                    color: #fff;
                }
            }
        }
        .footer{
            width: 66.9%;
            height:35%;
            margin:0.2% 0.1% 0.2% 0;
            background-image: url(../../assets/k4.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 2% 0 0 2%;
            box-sizing: border-box;
            #sectionBox-top{
                width: 100%;
                display: flex;
                flex-wrap: nowrap;
                margin-bottom:1%;
                // 标题
                dd{
                    
                    font-size: 16px;
                    line-height: 24px;
                    color: #fff;
                }
            }
            // 整个盒子大小
            #sectionBox{
                width:100%;
                height: calc(100% - 78px);
                overflow: hidden;
            }
            .anim{
                transition: all 0.5s;
                // font-size: 48px;
                // background: aqua;
            }
            #con1{
                li{
                    width: 100%;
                    line-height: 32px;
                    height: 32px;
                    display: flex;
                    flex-wrap: nowrap;
                    .sectionsize{
                        font-size: 14px;
                        line-height: 32px;
                        color: #fff;
                    }
                }
            } 
        }
        .footer-right{
            width: 32.9%;
            height: 35%;
            margin:0.2% 0.1% 0.2% 0;
            background-image: url(../../assets/k5.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 2% 2% 0 2%;
            box-sizing: border-box;
            overflow: hidden;
            // background: darksalmon;
            .sectionBox-top{
                width: 100%;
                display: flex;
                flex-wrap: nowrap;
                margin-bottom:4%;
                font-size: 16px;
                line-height: 24px;
                color: #fff;
            }
            #sectionBoxone{
                width:100%;
                height: calc(100% - 58px);
                overflow: hidden;
                .animone{
                    transition: all 1s;
                }
                #con2{
                    width: 100%;
                    display: flex;
                    flex-wrap: nowrap;
                    li{
                        width:33%;
                        margin: 3%;
                        box-sizing: border-box;
                        text-align: left;
                        .img{
                            width:100%;
                        }
                        p{
                            font-size: 12px;
                            color: #fff;
                            line-height: 20px;
                            margin-top:5%;
                        }
                    }
                } 
            }

        }
    }
}
// 测试九宫格
// .FirstScreen-box{
//     width: 100%;
//     height: 100vh;
//     position: fixed;
//     top: 0;
//     left: 0;
//     z-index: 1;
//     li{
//         width: 33.33%;
//         height: 33.33%;
//         border: 1px solid #ff6700;
//         float: left;
//         box-sizing: border-box;
//     }
// }
</style>
