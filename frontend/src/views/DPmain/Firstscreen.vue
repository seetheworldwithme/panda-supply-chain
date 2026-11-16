<template>
    <!-- 第一屏LED看板 -->
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
            <div class="section">
                <!-- 代码特效 -->
                <div class="codeEffects">
                    <div class="codeBox1"></div>
                    <div class="codeBox2"></div>
                    <div class="codeBox3"></div>
                    <div class="codeBox4"></div>
                    <div class="codeBox5"></div>
                </div>
                <!-- 中间左边内容 -->
                <div class="section-left">
                    <dl id="sectionBox-top">
                        <dd style="width:50%;">企业名称</dd>
                        <dd style="width:30%;">简称</dd>
                        <dd style="width:20%;">注册时间</dd>
                    </dl>
                    <div id="sectionBox">
                        <ul id="con1" ref="con1" :class="{anim:animate==true}">
                            <li v-for='item in items'>
                                <div class="sectionsize" style="width:50%;">{{item.enterpriseName}}</div>
                                <div class="sectionsize" style="width:30%;">{{item.simpleName}}</div>
                                <div class="sectionsize" style="width:20%;">{{item.regTime}}</div>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="central-section">
                    <div class="boxImgone"></div>
                    <div class="boxImgtwo"></div>
                    <!-- 中间盒子 -->
                    <div class="boxImg">
                        <p>结算金额累计(亿元)</p>
                        <b><countTo :startVal='0' :endVal='getPtSumPricesData' :duration='3000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></b>
                        <p style="font-size:20px;line-height:28px;margin-top:6%;">本年度业务量金额累计(元)</p>
                        <h4><countTo :startVal='0' :endVal='countYearBusesData' :duration='3000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></h4>
                    </div>
                    <!-- 数字展板左 -->
                    <div class="Digit-leftone">
                        <div class="leftoneSize">
                            <p>缴税金额累计(亿元)</p>
                            <span><countTo :startVal='0' :endVal='rateBills' :duration='1000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                        </div>
                    </div>
                    <div class="Digit-lefttwo">
                        <div class="leftoneSize">
                            <p>业务量金额累计(亿元)</p>
                            <span><countTo :startVal='0' :endVal='newbusCount' :duration='2000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                        </div>
                    </div>
                    <!-- 数字展板右 -->
                    <div class="Digit-rightone">
                        <div class="leftoneSize">
                            <p>合同份数累计(份)</p>
                            <span><countTo :startVal='0' :endVal='countSumContractData' :duration='2000' :useEasing="true" :autoplay="true"></countTo></span>
                        </div>
                    </div>
                    <div class="Digit-righttwo">
                        <div class="leftoneSize">
                            <p>物流量累计(万吨)</p>
                            <span><countTo :startVal='0' :endVal='logisticsCount' :duration='4000' :useEasing="true" :decimals="2" :autoplay="true"></countTo></span>
                        </div>
                    </div>
                    <!-- 粒子盒子 -->
                    <div class="particleBox">
                        <vue-particles
                        color="#00FFFF"
                        :particleOpacity="0.3"
                        :particlesNumber="28"
                        shapeType="circle"
                        :particleSize="10"
                        linesColor="#fff"
                        :linesWidth="0"
                        :lineLinked="false"
                        :lineOpacity="0"
                        :linesDistance="500"
                        :moveSpeed="2"
                        :hoverEffect="false"
                        :clickEffect="false">
                        </vue-particles>
                    </div>
                </div>
                <!-- 中间右边内容 -->
                <div class="section-right">
                    <!-- 业务量 -->
                    <div id="payment"></div>
                    <!-- <div id="invoicess" key="invoicess"></div>
                    <div class="invoiceBox" v-show="invoiceShow == 2">暂无数据</div> -->
                </div>
            </div>
            <!-- 底部内容 -->
            <div class="footer">
                <div class="footer-left">
                    <!-- 物流量 -->
                    <div id="paymentTaxes"></div>
                </div>
                <div class="footerBox">
                    <!-- 企业数 -->
                    <div id="barecharts"></div>
                </div>
                <div class="footer-right">
                    <!-- 缴税量 -->
                    <div id="chart_arc"></div>
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
            backcolor:false,
            routerboday:null,
            invoiceDate:{},             // 图表时间
            FirstScreenTime:'',
            animate:false,
            items:[],
            rateBills:0,                // 缴税金额累计
            logisticsCount:0,           // 物流量累计 
            getPtSumPricesData:0,       // 结算金额累计
            countSumContractData:0,     // 平台合同量
            newbusCount:0,              // 业务量累计
            countYearBusesData:0,       // 统计本年度业务量金额累计
            // 产品交易份额
            // invoiceShow:1,
            // newColor:[],                // 颜色
            // countPro:[],
            // invoice:'',
            // sumTaxAmount:'',            // 进项税额
            // arrTaxAmount:'',            // 销项税额

            // 业务量
            datapayment:[],                 // 时间
            statisticsData:[],              // 开票业务量

            // 物流量
            paymentdatapayment:[],          // 时间
            paymentserData:[],              // 物流量

            // 缴税
            option: null,
            number: 0, // 播放所在下标
            timer: null,
            datalogistics: [],
            logstatisticsData: [],

            // 企业新增数量
            // timerone:null,
            chasingTonnage:0,   // 注册企业数量
            arrsumNum:0,        // 当月新增数量
            yeararrsumNum:0,    // 全年新增数量
        }
    },
    created(){
        // 跳转
        this.routerboday = setInterval(this.routername,60000)
        setInterval(this.formatDate, 1000);
        // 企业信息列表
        this.api.statistics.enterData({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
        .then(res=>{
            if(res.data.code == 200){
                this.items = res.data.data; // 企业名称
            }
        })
        setInterval(this.scroll,2000) // 在钩子函数中调用我在method 里面写的scroll()方法
        this.getPtSumPrices();              // 结算金额累计
        setInterval(this.getPtSumPrices, 10000);
        this.countSumContracts();           // 平台合同量
        setInterval(this.countSumContracts, 15000);
        this.countYearBusess();             // 统计本年度业务量金额累计
        setInterval(this.countYearBusess, 200000);
        this.logisticsBill()   // 物流量累计
        setInterval(this.logisticsBill, 25000);
        this.taxartbarbox();   // 注册新增企业数
        setInterval(this.taxartbarbox, 30000);
    },
    mounted() {
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
        if(month == 0){
            year = year - 1
            month = 12;
            var lastDay= new Date(year,month,0).getDate();   // 当月最后一天
            day = lastDay
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
        // 产品交易份额
        // this.countBill({xzqhCode:sessionStorage.getItem('xzqhCode')});
        this.paymentCapital();             // 业务量
        this.paymentCapitalTaxes();        // 物流量
        this.logis();                      // 缴税图表
    },
    methods:{
        rightGos(){
            this.$router.push({name:'Secondscreen'});
        },
        leftGos(){
            this.$router.push({name:'warehouse'});
        },
        // 跳转
        routername(){
            this.$router.push({name:'Secondscreen'});
        },
        BackTime(){
            this.backcolor = !this.backcolor;
            if(this.backcolor){
                clearInterval(this.routerboday);
            }else{
                this.routerboday = setInterval(this.routername,60000);
            }
            
        },
        // 退出大屏
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

        // 统计本年度业务量金额累计
        countYearBusess(){
            this.countYearBusesData = 0;
            this.api.bill.countYearBuses()
            .then(res=>{
                if(res.data.code == 200){
                    this.countYearBusesData =  res.data.data.baseMoney*1;
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
            con1.style.marginTop = -32 + 'px';
            this.animate=!this.animate;
            let that = this;
            setTimeout(function(){
                that.items.push(that.items[0]);
                that.items.shift();
                con1.style.marginTop='0px';
                that.animate=!that.animate;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },500)
        },

        /* 
        缴税数据
        */
        logis(){
            this.api.statistics.countMonActBus(this.invoiceDate)
            .then(res=>{
                if(res.data.code == 200){
                    this.paymentdatapayment = [];       // 图表时间
                    let newdate = [];                   // 组装时间
                    this.paymentserData = [];           // 缴税量
                    for(let key in res.data.data.taxPayment){
                        newdate.push({date:key.replace(/-/g, "")*1});
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        let length = newdate[i].date.toString().length
                        if(length > 6){
                            let daty = newdate[i].date.toString().substring(6,8);
                            this.paymentdatapayment.push(year + '-' + newdatas + '-' + daty)
                        }else{
                            this.paymentdatapayment.push(year + '-' + newdatas)
                        }
                    }
                    for(let itme in this.paymentdatapayment){
                        let data = res.data.data.taxPayment[this.paymentdatapayment[itme]]/10000
                        this.paymentserData.push(data.toFixed(2))
                    }
                    this.getEchart();
                }
            })
        },
        getEchart() {
            let chartsmyChart = echarts.init(document.getElementById('chart_arc'));
            this.option = {
                title: {
                    text: '缴税额',
                    textStyle: {
                    color: '#fff',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(0,0,0,0.4)',
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
                                color: 'rgba(111, 144, 207, .5)'
                            },{
                                offset: 0,
                                color: 'rgba(255, 255, 255, 0)'
                            }])
                        }
                    },
                },
                color: ['#6F90CF'],
                grid: {
                    x: 25,
                    y: 75,
                    x2: 40,
                    y2: 0,
                    containLabel: true //轴上的数值
                },
                xAxis: {
                    type: 'category',
                    name: '时间',
                    data: this.paymentdatapayment,
                    boundaryGap: false,
                    axisTick: {
                        show: false
                    }, 
                    axisLine: {
                        lineStyle: {
                        color: '#fff',
                        }
                    },
                    splitLine: { show: false },
                },
                yAxis: {
                type: 'value',
                name: '万元',
                axisLine: {
                    show: true,
                    lineStyle: {
                    color: '#fff',
                    },
                },
                splitLine: { show: false },
                },
                series: [
                    {
                        // name: '人数',
                        type: 'line',
                        data: this.paymentserData,
                        symbolSize: 10,
                        itemStyle: {
                            opacity: 1,
                        },
                        emphasis: {
                            // label: {
                            //     show: true,
                            //     color: '#fff',
                            //     position: "top", 
                            //     distance: 12, 
                            //     // fontSize: 14
                            // },
                            itemStyle: {
                                color: '#6F90CF',
                                borderColor: '#fff',
                                borderWidth: 1,
                                opacity: .5
                            },
                        },
                        areaStyle: {
                            normal: {
                                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
                                    offset: 0,
                                    color: '#6F90CF'
                                },{
                                    offset: 1,
                                    color: 'rgba(0, 0, 0, 0)'
                                }]),
                            }
                        },
                        smooth: true,
                    }
                ]
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
                if (this.number > this.logstatisticsData.length) {
                    this.number = 0;
                }
            }, 1000);
        },

        /* 
        产品交易份额
        */
        // countBill(type){
        //     this.loading = true;
        //     this.api.listByPage.countPro(type)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             this.loading = false;
        //             if(res.data.data.length == 0){
        //                 this.invoiceShow = 2;
        //             }else{
        //                 this.invoiceShow = 1;
        //                 this.countPro = [];
        //                 let color = ["#FB6B6C","#FDAD6D","#C16AF7","#686DFA","#FB8D6D","#6E58F9","#FED172","#6BE1FB","#88FA7A","#71FEE5"];
        //                 this.newColor = color.splice(0,res.data.data.length);
        //                 let obj;
        //                 let num = 0;
        //                 for(let k=0;k<res.data.data.length;k++){
        //                     num += res.data.data[k].count;
        //                 }
        //                 for(let i=0;i<res.data.data.length;i++){
        //                     let allnum = ((res.data.data[i].count)/num *100).toFixed(2) + '%';
        //                     obj = {
        //                         value:res.data.data[i].count,
        //                         name:res.data.data[i].d_key + allnum
        //                     }
        //                     this.countPro.push(obj);
        //                 }
        //                 this.chartDom();  // 交易产品份额图表
        //             }

        //         }
        //     })
        // },
        // chartDom(){
        //     var chartDom = document.getElementById('invoicess');
        //     var myChart = echarts.init(chartDom);
        //     var optionChart;
        //     optionChart = {
        //         tooltip: {
        //             trigger: 'item'
        //         },
        //         legend: {
        //             type: 'scroll',
        //             orient: 'vertical',
        //             right:0,
        //             top: 20,
        //             bottom: 10,
        //             textStyle: {
        //                 color: '#fff'
        //             },
        //         },
        //         grid:{
        //             x:35,
        //             y:40,
        //             x2:50,
        //             y2:100,
        //         },
        //         noDataLoadingOption: {
        //             text: '暂无数据',
        //             effect: 'bubble',
        //             effectOption: {
        //                 effect: {
        //                     n: 0
        //                 }
        //             }
        //         },
        //         color:this.newColor,
        //         series: [
        //             {
        //                 name: '基础维护',
        //                 type: 'pie',
        //                 // roseType: 'area',
        //                 radius: [20, 100],
        //                 radius: '60%',
        //                 center: ['35%', '50%'],
        //                 data:this.countPro,
        //                 itemStyle: {
        //                     normal: {
        //                         label: {
        //                             textStyle: {
        //                                 color:'#fff'
        //                             }
        //                         }
        //                     }
        //                 },
        //             }
        //         ]
        //     };
        //     optionChart && myChart.setOption(optionChart);
        // },

        /* 
        业务量数据
        */
        paymentCapital(){
            this.api.statistics.countMonActBus(this.invoiceDate)
            .then(res=>{
                if(res.data.code == 200){
                    this.datapayment = [];      // 图表时间
                    let newdate = [];           // 组装时间
                    this.statisticsData = [];   // 实际业务量
                    for(let key in res.data.data.busCount){
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
                        let data = res.data.data.busCount[this.datapayment[itme]]/10000;
                        this.statisticsData.push(data.toFixed(2))
                    }
                }
                this.payment();
            }) 
        },
        payment(){
            // 业务量图表
            var chartDompayment = document.getElementById('payment');
            var myChartpayment = echarts.init(chartDompayment);
            var optionpayment;
            optionpayment = {
                title: {
                    text: '业务量',
                    textStyle: {
                    color: '#fff',
                    },
                },
                textStyle: {
                    color: '#fff',
                },
                grid: {
                    x: 50,
                    y: 75,
                    x2: 40,
                    y2: 30,
                },
                xAxis: {
                    type: 'category',
                    name: '时间',
                    data: this.datapayment,
                    axisPointer: {
                        type: 'line',
                    },
                    nameTextStyle: {
                        color: '#fff',
                    },
                },

                yAxis: [
                    {
                        type: 'value',
                        name: '万元',
                        // min: 0,
                        // max: 16000,
                        // interval: 4000,
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
                        splitLine: { show: false },
                    },
                ],
                series: [
                    {
                        type: 'pictorialBar',
                        barWidth: '60%',
                        label: {
                            normal: {
                            show: true,             // 是否显现，不显示的话设置成false
                            position: 'top',        // 显示的位置
                            distance: 0,            // 距离侄子的值 // label要显示的值比如： 20%
                            color: '#fff',
                            },
                        },
                        itemStyle: {
                            normal: {
                            color: {
                                type: 'linear',
                                x: 0,
                                y: 0,
                                x2: 0,
                                y2: 1,
                                colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(48, 236, 166, 1)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(0, 34, 66, 0.2)',      // 100% 处的颜色
                                },
                                ],
                                globalCoord: false,                     // 缺省为 false
                            },                                          //渐变颜色
                            },
                        },
                        symbol:'path://M12.000,-0.000 C12.000,-0.000 16.074,60.121 22.731,60.121 C26.173,60.121 -3.234,60.121 0.511,60.121 C7.072,60.121 12.000,-0.000 12.000,-0.000 Z',
                        data:this.statisticsData,
                        animationDuration: 1600,
                        animationEasing: "cubicInOut",
                    },
                ],
            };
            optionpayment && myChartpayment.setOption(optionpayment);
        },

        /*
        企业新增数
        */
        taxartbarbox(){
            this.arrsumNum = 0;
            this.chasingTonnage = 0;
            this.yeararrsumNum = 0;
            this.api.statistics.count({xzqhCode:sessionStorage.getItem('xzqhCode')})
            .then(res=>{
                if(res.data.code == 200){
                    this.chasingTonnage = res.data.data.sumCount*1;           // 注册企业数量
                    this.arrsumNum = res.data.data.monCount*1;                // 当月新增数量
                    this.yeararrsumNum = res.data.data.yearCount*1;           // 全年新增数量
                    this.taxartbar(this.chasingTonnage,this.arrsumNum,this.yeararrsumNum);
                }
            })
        },
        taxartbar(one,two,tree){
            // 企业图表
            var chartDomoptionsss = document.getElementById('barecharts');
            var myChartoptionsss = echarts.init(chartDomoptionsss);
            var optionsss;
            let oneData = one;
            let twoData = two;
            let treeData = tree;

            // let oneData = 29;
            // let twoData = 15;
            // let treeData = 30;
            optionsss = {
                grid: {
                    top: 60,
                    bottom: 80,
                },
                tooltip: {},
                xAxis: {
                    data: ['注册企业数量', '当月新增数量', '全年新增数量'],
                    axisTick: {
                        show: false,
                    },
                    axisLine: {
                        show: false,
                    },
                    axisLabel: {
                        interval: 0,
                        textStyle: {
                            color: '#fff',
                            // fontSize: 20,
                        },
                        margin: 60, //刻度标签与轴线之间的距离。
                    },
                },
                yAxis: {
                    splitLine: {
                        show: false,
                    },
                    axisTick: {
                        show: false,
                    },
                    axisLine: {
                        show: false,
                    },
                    axisLabel: {
                        show: false,
                    },
                },
                series: [
                    {
                        name: '',
                        type: 'pictorialBar',
                        symbolSize: [45, 20],
                        symbolOffset: [0, -10],
                        z: 12,
                        data: [
                            {
                                name: '注册企业数量',
                                value:oneData,
                                trueVal: '98',
                                symbolPosition: 'end',
                                itemStyle: {
                                    color: '#00fff5', //圆柱顶部颜色
                                },
                            },       
                            {
                                name: '当月新增数量',
                                value: twoData,
                                trueVal: '100',
                                symbolPosition: 'end',
                                itemStyle: {
                                    color: '#ffcc00', //圆柱顶部颜色
                                },
                            },
                            {
                                name: '全年新增数量',
                                value:treeData,
                                trueVal: '90',
                                symbolPosition: 'end',
                                itemStyle: {
                                    color: '#b9b7ff', //圆柱顶部颜色
                                },
                            },
                        ],
                    },
                    {
                        name: '',
                        type: 'pictorialBar',
                        symbolSize: [45, 20],
                        symbolOffset: [0, 10],
                        z: 12,
                        data: [
                            {
                                name: '注册企业数量',
                                value: 10,
                                trueVal: '98',
                                itemStyle: {
                                    color: '#43bafe', //圆柱底部颜色
                                },          
                            },  
                            {
                                name: '当月新增数量',
                                value: 10,
                                trueVal: '100',
                                itemStyle: {
                                    color: '#ff7800', //圆柱底部颜色
                                },
                            },
                            {
                                name: '全年新增数量',
                                value: 10,
                                trueVal: '90',
                                itemStyle: {
                                    color: '#e9a5ff', //圆柱底部颜色
                                },
                            },
                        ],
                    },
                    {
                        name: '',
                        type: 'pictorialBar',
                        symbolSize: [80, 35],
                        symbolOffset: [0, 28],
                        z: 11,
                        data: [
                            {
                                name: '注册企业数量',
                                value: 10,
                                trueVal: '98',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#43bafe', //底部内圆圈颜色
                                    borderWidth: 5,
                                },    
                            },    
                            {
                                name: '当月新增数量',
                                value: 10,
                                trueVal: '100',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#ff7800', //底部内圆圈颜色
                                    borderWidth: 5,
                                },
                            },
                            {
                                name: '全年新增数量',
                                value:10,
                                trueVal: '90',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#e9a5ff', //底部内圆圈颜色
                                    borderWidth: 5,
                                },
                            },
                        ],
                    },
                    {
                        name: '',
                        type: 'pictorialBar',
                        symbolSize: [100, 50],
                        symbolOffset: [0, 42],
                        z: 10,
                        data: [
                            {
                                name: '注册企业数量',
                                value: 10,
                                trueVal: '98',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#43bafe', //底部外圆圈颜色
                                    borderType: 'dashed',
                                    borderWidth: 5,
                                },  
                            },     
                            {
                                name: '当月新增数量',
                                value: 10,
                                trueVal: '100',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#ff7800', //底部外圆圈颜色
                                    borderType: 'dashed',
                                    borderWidth: 5,
                                },
                            },
                            {
                                name: '全年新增数量',
                                value: 10,
                                trueVal: '90',
                                itemStyle: {
                                    color: 'transparent',
                                    borderColor: '#e9a5ff', //底部外圆圈颜色
                                    borderType: 'dashed',
                                    borderWidth: 5,
                                },
                            },
                        ],
                    },
                    {
                        type: 'bar',
                        silent: true,
                        barWidth: 45,
                        // barGap: '-10%',
                        data: [
                            {
                                name: '注册企业数量',
                                value: oneData,
                                trueVal: '98',
                                label: {
                                    show: true,
                                    position: 'top',
                                    distance: 18,
                                    textStyle: {
                                        color: '#00fff5', //柱子对应数值颜色
                                        // fontSize: 40,
                                    },
                                }, 
                                itemStyle: {
                                    normal: {
                                        color: {
                                            x: 0,
                                            y: 0,
                                            x2: 0,
                                            y2: 1,
                                            type: 'linear',
                                            global: false,
                                            colorStops: [
                                                {
                                                    offset: 0,
                                                    color: 'rgba(0,255,245,0.5)',
                                                },
                                                {
                                                    offset: 1,
                                                    color: '#43bafe', //底部渐变颜色
                                                },
                                            ],
                                        },
                                    },
                                },
                            },
                            {
                                name: '当月新增数量',
                                value:twoData,
                                trueVal: '100',
                                label: {
                                    show: true,
                                    position: 'top',
                                    distance: 18,
                                    textStyle: {
                                        color: '#ffcc00', //柱子对应数值颜色
                                        // fontSize: 40,
                                    },
                                },
                                itemStyle: {
                                    color: {
                                        x: 0,
                                        y: 0,
                                        x2: 0,
                                        y2: 1,
                                        type: 'linear',
                                        global: false,
                                        colorStops: [
                                            {
                                                offset: 0,
                                                color: 'rgba(255,204,0,0.5)',
                                            },
                                            {
                                                offset: 1,
                                                color: '#ff7800', //底部渐变颜色
                                            },
                                        ],
                                    },
                                },
                            },
                            {
                                name: '全年新增数量',
                                value: treeData,
                                trueVal: '90',
                                label: {
                                    show: true,
                                    position: 'top',
                                    distance: 18,
                                    textStyle: {
                                        color: '#b9b7ff', //柱子对应数值颜色
                                        // fontSize: 40,
                                    },
                                },
                                itemStyle: {
                                    color: {
                                        x: 0,
                                        y: 0,
                                        x2: 0,
                                        y2: 1,
                                        type: 'linear',
                                        global: false,
                                        colorStops: [
                                            {
                                                offset: 0,
                                                color: 'rgba(185,183,255,0.5)',
                                            },
                                            {
                                                offset: 1,
                                                color: '#e9a5ff', //底部渐变颜色
                                            },
                                        ],
                                    },
                                },
                            },
                        ],
                    },
                ],
            };
            optionsss && myChartoptionsss.setOption(optionsss);

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
                title: {
                    text: '物流量',
                    textStyle: { 
                        color:'#fff'
                    }
                },
                textStyle: {
                    color: '#fff'
                },
                legend: {
                    textStyle: {
                        color: "#fff"
                    }
                },
                grid:{
                    x:50,
                    y:75,
                    x2:40,
                    y2:20,
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
                        name: '万吨',
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
                        type: 'line',
                        symbolSize: 150,
                        symbol:'image://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAFAAAABRCAYAAABFTSEIAAAACXBIWXMAAAsSAAALEgHS3X78AAAEp0lEQVR42u3cz4sjRRTA8W9Vd3Vn8mMmjj9WQWSRZQ+CsH+B7MnDIgiCd0E8CYJ/gOAIelo8ehUP/gF6WLw5/gMueFP2sIcF0dHd2Z1kknR11fOQZJJJMtlZd03H7HtQpNOTnpn+8Lrm1etmjIig8e/DKoECKqACKqCGAiqgAiqghgIqoAIqoIYCKqACKqCGAiqgAiqghgIqoAJudKTr+osZMNPvBUQBHwHsPF9fB9R0DeHMOQ6T6WOrhEzXBM4swDOL0M6CrArRVoq3t2dGUIb9fTvatg8ZZup1PDBgzPmy98mey6qfzjLz2WaWjEUZKEvGyi9nWyneMOvGIyFQo2Sbg4MUSChpU9IeTTUpJdsEajPZOJeJG5uBZj7rLLduWS5dGm6XNLEELOFUFj54ACJCaychkpDSASK3bwsXL0YgVpWJKwM0iy9Zy8HdGru7jvt3Pbu7w0wES7drTwAbjTHMGCsQcIAnYTC1/wRx0wEnl27JNgZI8HQ6Kc1mQq83RNzaMjPzXqDbjTQaJRFLxIyyMSxAXEkWrhrQzAAmo5HOjCQf7jflILxOkohL+aUPgV4vEGNJo+E5PAy02+UIMEwBxo0CPDP7Dg5SnEtpt1PA0e87XO25FOoh8IYIH2Y5b45RzGAQBiIltZoHxqMcjbksXAVgdc2EQMYzzzdotyeZWKuleULXJtwT4SODfC2QCWR+IF9KnjuX1Xbo99Op7LVE8iXlz0YBTk5SyLEEjo5OLuccEoFUvHfO+reuUPx4zftXAIcx1hdcF+/TvFab4A0Bs0VwqyhpVnkJT89/Q4DDQ0e77YCMwIUsFMeFZD856699URRvX4nxE4A/jbnxXp7v4Zw3ReGNSDHI8wFQjIafuoyn58L/fB6sth/Ybg9fez2TRC6QZcZYvgHsazF+MP7YCyLXcM7gvSXLDGBqYDg+NhwdmSpPoTrAkub0W+f4FSB1fDucIunMHSLpO8WAH0rSy8u+19MBCHB4OHzd2pI+CEUhpigEiN+l6WcdY252jLn5s7Wf472ImPcN8pUl/tEHoV4XWq1Ke4KrLmPsTA3oODpytFoOyJKSyzHyMSIxteWngMW5cSEdDJQUhTdZVgxOz3/+jFJm4+bA2e5JpNU6WZ4Fw99JwnWMKccwpeddP+B7GZTNUPKqybJy0O+Hs1YfMz9swwvpB8fbGDG0GuGkkK7V0hxSmZQpABI8l2z0v3sJf50qpAMJCd2qCulql3LD1lRGQjm7lEsDz0rkxTQOfiPPxUBcuJTbbhss/Y1eyi3NwsmKInmkZsKk5gtPUzNhvp11507CSy/X6XYStpvFudpZw1ZWIOF4Cq6SdtbKbioJyAhRTu3u9yMJXerN+ugvaQQsjcZ8Q3VnZwxlSDhe1lB9GjrSw5b+1avT8+Jw+979nNaOI6U3KpTrWAosxVQmygK4ld8X0ZtK/7eViExD7O1NQPb3T7fsl4/4sBpwYzPwjFbTo95Yl9l9Vd1YN1X/147HebSjary1AHyc5qc+XLQEQx9ve8Kg6xr6hKoCKqACKqCGAiqgAiqghgIqoAIqoIYCKqACKqCGAiqgAiqghgIq4JrHP8fEWV8FMTmOAAAAAElFTkSuQmCC',
                        itemStyle: {
                            normal: {
                                color: '#00f8ff'
                            }
                        },
                        data:this.logstatisticsData,        // this.logstatisticsData
                        // animationDuration: 3000,
                        // animationEasing: "cubicInOut",
                        label: {
                            normal: {
                                show: true,             // 是否显现，不显示的话设置成false
                                position: "top",        // 显示的位置
                                distance: -60,            // 距离侄子的值 // label要显示的值比如： 20%
                                color:'#fff'
                            }
                        },
                    }
                ]
            };
            optionpaymentss && myChartpaymentss.setOption(optionpaymentss);
        },
    },
    // 销毁定时器
    beforeDestroy() {
        clearInterval(this.routerboday);
        clearInterval(this.timer);
        // clearInterval(this.timerone);
        if (this.formatDate || this.getPtSumPrices || this.countSumContracts || this.countYearBusess || this.logisticsBill || this.taxartbarbox) {
            clearInterval(this.formatDate); // 在Vue实例销毁前，清除时间定时器
            clearInterval(this.getPtSumPrices); // 在Vue实例销毁前，清除时间定时器
            clearInterval(this.countSumContracts); 
            clearInterval(this.countYearBusess); 
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
    padding:1%;
    box-sizing: border-box;
    background-image: url(../../assets/backimg.jpg);
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
        // 上半部分内容
        .section{
            width: 100%;
            height: 65%;
            display: flex;
            flex-wrap:wrap;
            flex-direction: column-reverse;
            // 代码特效
            .codeEffects{
                width: 80%;
                height: 50%;
                position: fixed;
                top: 13%;
                left: 10%;
                z-index: 1;
                overflow: hidden;
                opacity: .3;
                .codeBox1{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 0;
                    background: url(../../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
                @keyframes mymove
                {
                    0%  {top:-20%;}
                    10% {top:20%;}
                    20% {top:40%;}
                    30% {top:60%;}
                    40% {top:80%;}
                    50% {top:100%;}
                    60% {top:80%;}
                    70% {top:60%;}
                    80% {top:40%;}
                    90% {top:20%;}
                    100% {top:-20%;}
                }
                .codeBox2{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 20%;
                    background: url(../../assets/code1.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove1 14s infinite alternate linear;
                }
                @keyframes mymove1
                {
                    0%  {bottom:-20%;}
                    10% {bottom:20%;}
                    20% {bottom:40%;}
                    30% {bottom:60%;}
                    40% {bottom:80%;}
                    50% {bottom:100%;}
                    60% {bottom:80%;}
                    70% {bottom:60%;}
                    80% {bottom:40%;}
                    90% {bottom:20%;}
                    100% {bottom:-20%;}
                }
                .codeBox3{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 40%;
                    background: url(../../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
                .codeBox4{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 60%;
                    background: url(../../assets/code1.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove1 14s infinite alternate linear;
                }
                .codeBox5{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 80%;
                    background: url(../../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
            }
            // 左侧内容
            .section-left{
                width: 33%;
                height: 58%;
                padding:3.5% 0 0 2%;
                box-sizing: border-box;
                #sectionBox-top{
                    width: 100%;
                    display: flex;
                    flex-wrap: nowrap;
                    margin-bottom:5%;
                    dd{
                        font-size: 16px;
                        line-height: 24px;
                        color: #fff;
                        background: '#FABC61';
                    }
                }
                #sectionBox{
                    width:100%;
                    height: calc(100% - 100px);
                    overflow: hidden;
                }
                .anim{
                    transition: all 0.5s;
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
            // 中间内容
            .central-section{
                width:34%;
                height: 100%;
                position: relative;
                // 粒子盒子
                .particleBox{
                    width: 140%;
                    height: 100%;
                    position: absolute;
                    top: 60%;
                    left: -20%;
                    z-index: 2;
                    opacity: 0.4;
                }
                // 数字展板左
                .Digit-leftone{
                    width: 70%;
                    height:116px;
                    background-image: url(../../assets/leftone.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    position: absolute;
                    top: 3%;
                    left: -48%;
                    .leftoneSize{
                        padding: 3% 6%;
                        box-sizing: border-box;
                        p{
                          font-family: aa;  
                          font-size: 16px;
                          line-height: 28px;
                          color: #fff;
                        }
                        span{
                          font-family: bb;  
                          font-size: 38px;
                          color: #fff;
                        }
                    }
                }
                .Digit-lefttwo{
                    width: 92%;
                    height:116px;
                    background-image: url(../../assets/leftsize.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    position: absolute;
                    top: 26%;
                    left: -86.5%;
                    .leftoneSize{
                        padding: 3% 6%;
                        box-sizing: border-box;
                        p{
                          font-family: aa;  
                          font-size: 16px;
                          line-height: 28px;
                          color: #fff;
                        }
                        span{
                          font-family: bb;  
                          font-size: 38px;
                          color: #fff;
                        }
                    }
                }
                // 数字展板右
                .Digit-rightone{
                    width: 67%;
                    height:116px;
                    background-image: url(../../assets/rightone.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    position: absolute;
                    top: 3%;
                    right: -45%;
                    .leftoneSize{
                        width: 58%;
                        padding: 3% 4%;
                        box-sizing: border-box;
                        float: right;
                        p{
                          font-family: aa;  
                          font-size: 16px;
                          color: #fff;
                          line-height: 28px;
                        }
                        span{
                          font-family: bb;  
                          font-size: 38px;
                          color: #fff;
                        }
                    }
                }
                .Digit-righttwo{
                    width: 95%;
                    height: 116px;
                    background-image: url(../../assets/sizekaung.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    position: absolute;
                    top: 20%;
                    right: -87%;
                    .leftoneSize{
                        width: 43%;
                        padding: 4% 4%;
                        box-sizing: border-box;
                        float: right;
                        p{
                          font-family: aa;  
                          font-size: 16px;
                          color: #fff;
                          line-height: 28px;
                        }
                        span{
                          font-family: bb;  
                          font-size: 38px;
                          color: #fff;
                        }
                    }
                }
                .boxImgone{
                    // background: burlywood;
                    width:100%;
                    height:100%;
                    position: absolute;
                    left: 0;
                    right: 0;
                    top: 25%;
                    bottom: 0;
                    margin: auto;
                    background-image: url(../../assets/yuan2.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    animation: turn 10s infinite linear;
                }
                @keyframes turn
                {
                    0% {  transform: rotate(0deg); }
                    100% {  transform: rotate(360deg); }
                }
                .boxImgtwo{
                    width:100%;
                    height:100%;
                    position: absolute;
                    left: 0;
                    right: 0;
                    top: 25%;
                    bottom: 0;
                    margin: auto;
                    background-image: url(../../assets/yuan1.png);
                    background-size: 100% 100%;
                    background-repeat: no-repeat;
                    animation: turns 10s infinite linear;
                }
                @keyframes turns
                {
                    0% {  transform: rotate(360deg); }
                    100% {  transform: rotate(0deg); }
                }
                .boxImg{
                    width: 80%;
                    height: 22%;
                    position: absolute;
                    left: 0;
                    right: 0;
                    top: 0;
                    bottom: 0;
                    margin: auto;
                    text-align: center;
                    @font-face {
                        font-family: bb;
                        src: url(../../assets/icon/QuartzRegular.ttf);
                    }
                    b{
                        font-family: bb;
                        font-size:120px;
                        color: #fff;
                        font-weight: 500;
                        line-height: 140px;
                    }
                    h4{
                        font-family: bb;
                        font-size: 38px;
                        // line-height: 48px;
                        color: #FCC026;
                    }
                    p{
                        font-family: aa;
                        font-size: 20px;
                        color: #fff;
                        line-height: 32px;
                    }
                }
            }
            // 右侧内容
            .section-right{
                width: 33%;
                height: 58%;
                padding:1%;
                box-sizing: border-box;
                z-index: 10;
                // .el-date-picker{
                //     width: 100%;
                //     float: right;
                // }
                #payment{
                    width: 100%;
                    height:100%;
                }
                // .invoiceBox{
                //     width: 100px;
                //     height: 20px;
                //     color: #999;
                //     text-align: center;
                //     line-height: 20px;
                //     position: absolute;
                //     top: 0;
                //     left: 0;
                //     bottom: 0;
                //     right: 0;
                //     margin: auto;
                // }
            }
        }
        // 下半部分内容
        .footer{
            width: 100%;
            height: 35%;
            display: flex;
            flex-wrap:wrap;
            flex-direction: column-reverse;
            z-index: 10;
            .footer-left{
                width: 33.5%;
                height: 100%;
                padding: 1% 0 0 0;
                box-sizing: border-box;
                #paymentTaxes{
                    width: 100%;
                    height: 100%;
                }
            }
            .footerBox{
                width: 33%;
                height: 80%;
                z-index: 10;
                #barecharts{
                    width: 100%;
                    height: 100%;
                }
            }
            .footer-right{
                width: 33.5%;
                height: 100%;
                padding: 1% 0 0 1%;
                box-sizing: border-box;
                z-index: 10;
                #chart_arc{
                    width: 100%;
                    height: 100%;
                }
            }
        }
    }
}
// .sn-container {
//   left: 50px;
//   top: 690px;
//   width: 895px;
//   height: 400px;
//   .chartsdom {
//     width: 100%;
//     height: 95%;
//   }
// }
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
