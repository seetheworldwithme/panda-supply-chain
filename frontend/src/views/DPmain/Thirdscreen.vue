<template>
    <!-- 第二屏LED看板999 -->
    <div class="FirstScreen">
        <div class="rightGo" @click="rightGos">
            <i class="el-icon-arrow-right"></i>
        </div>
        <div class="leftGo" @click="leftGos">
            <i class="el-icon-arrow-left"></i>
        </div>
        <div class="header">
            <p>内部数据工作平台</p>
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
                <div class="sectionssbox">
                    <el-select v-model="sectionssboxvalue" size="mini" style="width:100%;" :popper-append-to-body="false" popper-class="select-popper" @change="sboxvalueprises($event)">
                        <el-option
                        v-for="item in sleectoptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                </div>
                <!-- 资料收集户数 -->
                <div class="sectionleft">
                    <div id="sectionleft"></div>
                </div>
                <!-- 资料数据工作 -->
                <div class="sectionright">
                    <div id="sectionright"></div>
                </div>
            </div>
            <!-- 网络货运平台占比 -->
            <div class="conentLogins">
                <div class="conentLoginsBox">
                    <span>{{chartTime}} </span>
                    <span style="font-size:14px;font-weight:200;color:#ea222e;" v-if="ifshowslc"> *</span>
                    <span style="font-size:14px;font-weight:200;" v-if="ifshowslc"> 超出100%为预留吨数</span>
                </div>
                <div class="tooleft" v-if="ifshowslc">
                    <div class="toolone" style="background-image: linear-gradient(rgba(111, 178, 207, 1), rgba(111, 178, 207, 0));"></div><span>运输发票占比</span>
                    <el-tooltip class="item" effect="dark" content="运输发票占开票吨数" placement="top" style="margin:0 2%;float:left;line-height: 20px;">
                        <i class="el-icon-question" style="color:#fff;font-size: 14px;"></i>
                    </el-tooltip>
                    <div class="toolone" style="margin-left:4%;background-image: linear-gradient(rgba(111, 207, 151, 1), rgba(111, 207, 151, 0));"></div><span>网络货运占比</span>
                    <el-tooltip class="item" effect="dark" content="网络货运占开票吨数" placement="top" style="margin:0 2%;float:left;line-height: 20px;">
                        <i class="el-icon-question" style="color:#fff;font-size: 14px;"></i>
                    </el-tooltip>
                </div>
                <!-- <div class="tooright" v-if="ifshowslc">

                </div> -->
                <!-- 右侧框 -->
                <div class="sleectright">
                    <el-select v-model="slctvaluesctp" filterable size="mini" style="width:65%;margin-right:3%;" :popper-append-to-body="false" popper-class="select-popper" @change="slctvalprises($event)">
                        <el-option
                        v-for="item in esctpoptions"
                        :key="item.id"
                        :label="item.enterpriseName"
                        :value="item.id">
                        </el-option>
                    </el-select>
                    <el-select v-model="sleectvalue" size="mini" style="width:30%;" :popper-append-to-body="false" popper-class="select-popper" @change="enterprises($event)">
                        <el-option
                        v-for="item in sleectoptions"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                </div>
                <div class="burdenLogins" v-if="ifshowslc">
                    <div id="burdenLogins"></div>
                </div>
                <div class="burdenLs" v-if="ifshowslc == false">
                    <div id="burdenLs"></div>
                </div>
                
            </div>
        </div>
    </div>
</template>
<script>
import ELoginVue from '../ELogin.vue';
var echarts = require('echarts');
export default{
    data(){
        return {
            backcolor:false,
            chartTime:'',
            routerboday:null,
            invoiceDate:{},             // 图表时间
            buildinvoiceDate:{},        // 网络货运平台占比时间
            FirstScreenTime:'',

            // 客户统计工作
            countoptionpayment: null,
            numbertwo: 0, // 播放所在下标
            timertwo: null,
            countMonment:[],                // 时间
            countMonData:[],                // 数据完成户数
            countMonDatas:[],                // 开票业务户数

            // 资料数据工作
            sectionssboxvalue:'',
            optionpayment: null,
            numberone: 0, // 播放所在下标
            timerone: null,
            datapayment:[],                 // 时间
            paymentData:[],                 // 数据录入完成
            statisticsData:[],              // 开票业务完成

            // 网络货运平台占比
            optionLogins: null,
            number: 0, // 播放所在下标
            timer: null,
            burdenLoginsTime:[],              // 公司名称--时间
            burdenLoginsData:[],              // 网络货运占比
            allLoginsData:[],                 // 运输发票占比

            sleectoptions:[],                   // 时间列表 
            sleectvalue:'',                    // 选中时间
            esctpoptions:[],                   // 企业列表
            slctvaluesctp:-1,                  // 选中企业
            ifshowslc:true,                    // 物流简报显示隐藏
            // 单个企业数据
            // tidTon:[],                      // 物流时间
            ReservedTon:[],                   // 开票数量吨数
            AllReservedTon:[],                // 运输发票吨数
            sssReservedTon:[],                // 网络货运吨数
            DataTime:null,
        }
    },
    created(){
        // 跳转
        this.routerboday = setInterval(this.routername,60000);
        this.api.statistics.enterData({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
        .then(res=>{
            if(res.data.code == 200){
                this.esctpoptions = res.data.data; // 企业名称
                this.esctpoptions.unshift({enterpriseName:'全部企业',id:-1})
            }
        })
        // 获取当前时间
        let date = new Date();
        let year = date.getFullYear();
        let month = date.getMonth() + 1;
        let day = date.getDate();
        let Years = (year - 2019) + 1;
        this.sleectvalue = year;
        this.sectionssboxvalue = year;
        for(let i=0;i<Years;i++){
            this.sleectoptions.push({
                value:(year - i),
                label:year - i + '年'
            })
        }
        let lastDay= new Date(year,month,0).getDate();   // 当月最后一天
        day = lastDay< 10? '0'+ lastDay:lastDay;
        month = month< 10? '0'+ month:month;
        let adate = year + "-" + month + "-" + day;
        // 客户数占比 业务量占比时间
        this.buildinvoiceDate = {
            xzqhCode:sessionStorage.getItem('xzqhCode'),
            startTime:year + '-01-01',
            endTime:adate
        };
        this.chartTime = '物流简报(' + year + '年01月-' + month + '月)';
        // 客户数占比 业务量占比时间
        // this.invoiceDate = {
        //     xzqhCode:sessionStorage.getItem('xzqhCode'),
        //     startTime:(year - 1) + '-01-01',
        //     endTime:adate
        // };
        setInterval(this.formatDate, 1000);
    },
    mounted() {
        // 客户统计工作
        this.countMonDataTaxes(this.buildinvoiceDate); 
        // 资料数据工作
        this.paymentCapital(this.buildinvoiceDate);
        this.burdenLoginsapi(this.buildinvoiceDate);                           // 网络货运平台占比
    },
    methods:{
        rightGos(){
            this.$router.push({name:'Maps'});
        },
        leftGos(){
            this.$router.push({name:'Secondscreen'});
        },
        // 跳转
        routername(){
            this.$router.push({name:'Firstscreen'})
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

        /* 
        客户统计工作量数据
        */
        // 年份下拉框 
        sboxvalueprises(id){
            let date = new Date();
            let year = date.getFullYear();
            let month = date.getMonth() + 1;
            // 全部企业
                if(id != year){
                    // 不是当年
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-12-31"
                    };
                }else{
                    // 当年情况
                    let lastDay= new Date(year,month,0).getDate();   // 当月最后一天
                    month = month<10? '0' + month:month;
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-" + month + "-" + lastDay
                    }; 
                }
                this.countMonDataTaxes(this.buildinvoiceDate);
                this.paymentCapital(this.buildinvoiceDate);
        },

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
            // 客户统计工作量数据图表
            let countmyChartpayment = echarts.init(document.getElementById('sectionleft'));
            this.countoptionpayment = {
                title: {
                    text: '客户数占比',
                    textStyle:{
                        color:'#fff',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(0, 0, 0, .6)',
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
                                color: 'rgba(111, 207, 207, .2)'
                            },{
                                offset: 0,
                                color: 'rgba(255, 255, 255, 0)'
                            }])
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
                            '<br>&nbsp&nbsp客户数占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                textStyle: {
                    color: '#fff'
                },
                legend: {
                    selectedMode:false,
                    data: ['开票业务户数', '数据完成户数'],
                    textStyle: {
                        color: "#fff"
                    }
                },
                grid:{
                    x:50,
                    y:70,
                    x2:60,
                    y2:50,
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
                                backgroundColor: 'rgba(111, 207, 207, .4)',
                            },
                        },
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
                                color: '#fff',
                            },
                        },
                        nameTextStyle: {
                            color: '#fff',
                        },
                        splitLine: { show: false },
                    }
                ],
                series: [
                    {
                        name: '开票业务户数',
                        type: 'bar',
                        color: 'rgba(255, 255, 255, 0.15)',
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        data: this.countMonDatas,
                        z: 0,
                        barWidth : 40,//柱图宽度
                        silent: true,
                        animation: true, // 关闭动画效果
                    },
                    {
                        name: '数据完成户数',
                        type: 'bar',
                        color: {
                            type: 'linear',
                            x: 0,
                            y: 0,
                            x2: 0,
                            y2: 1,
                            colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(111, 207, 207, 1)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(111, 207, 207, 0)',      // 100% 处的颜色
                                },
                            ],
                            globalCoord: false,                     // 缺省为 false
                        },
                        data: this.countMonData,
                        barWidth : 40,//柱图宽度
                        barGap: '-100%', // 设置柱形重合的重要步骤
                    }
                ]
            };
            this.countoptionpayment && countmyChartpayment.setOption(this.countoptionpayment);
            window.addEventListener('resize', () => {
                this.countoptionpayment.resize();
            });  
            this.timertwo = setInterval(() => {
                countmyChartpayment.dispatchAction({
                    type: 'showTip',
                    seriesIndex: 0,
                    dataIndex: this.numbertwo
                });
                this.numbertwo++;
                if (this.numbertwo > this.countMonment.length) {
                    this.numbertwo = 0;
                }
            }, 2500);
        },

        /* 
        资料数据工作
        */
        paymentCapital(type){
            this.loading = true;
            this.api.bill.business(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.datapayment = [];      // 图表时间
                    let newdate = [];           // 组装时间
                    this.paymentData = [];      // 数据录入完成
                    this.statisticsData = [];   // 开票业务完成
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
            // 资料数据工作图表
            let myChartpayment = echarts.init(document.getElementById('sectionright'));
            this.optionpayment = {
                title: {
                    text: '业务量占比',
                    textStyle:{
                        color:'#fff',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(0, 0, 0, .6)',
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
                                color: 'rgba(111, 138, 207,  .2)'
                            },{
                                offset: 0,
                                color: 'rgba(255, 255, 255, 0)'
                            }])
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
                            '<br>&nbsp&nbsp业务量占比：' +
                            b + `%`
                        )
                    listItem = list.join('<br>')
                    return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                textStyle: {
                    color: '#fff'
                },
                legend: {
                    selectedMode:false,
                    data: ['开票业务完成', '数据录入完成'],
                    textStyle: {
                        color: "#fff"
                    }
                },
                grid:{
                    x:70,
                    y:70,
                    x2:60,
                    y2:50,
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
                                backgroundColor: 'rgba(111, 138, 207,  .4)',
                            },
                        },
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
                                color: '#fff',
                            },
                        },
                        nameTextStyle: {
                            color: '#fff',
                        },
                        splitLine: { show: false },
                    }
                ],
                series: [
                    {
                        name: '开票业务完成',
                        type: 'bar',
                        color: {
                            type: 'linear',
                            x: 0,
                            y: 0,
                            x2: 0,
                            y2: 1,
                            colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(255, 255, 255, 0.15)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(255, 255, 255, 0.15)',      // 100% 处的颜色
                                },
                            ],
                            globalCoord: false,                     // 缺省为 false
                        },
                        barGap: '-100%', // 设置柱形重合的重要步骤
                        data: this.statisticsData,
                        barWidth : 40,//柱图宽度
                    },
                    {
                        name: '数据录入完成',
                        type: 'bar',
                        color: {
                            type: 'linear',
                            x: 0,
                            y: 0,
                            x2: 0,
                            y2: 1,
                            colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(111, 138, 207,  1)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(111, 138, 207, 0)',      // 100% 处的颜色
                                },
                            ],
                            globalCoord: false,                     // 缺省为 false
                        },
                        data: this.paymentData,
                        barWidth : 40,//柱图宽度
                        barGap: '-100%', // 设置柱形重合的重要步骤
                    }
                ]
            };
            this.optionpayment && myChartpayment.setOption(this.optionpayment);
            window.addEventListener('resize', () => {
                myChartpayment.resize();
            });  
            this.timerone = setInterval(() => {
                myChartpayment.dispatchAction({
                    type: 'showTip',
                    seriesIndex: 0,
                    dataIndex: this.numberone
                });
                this.numberone++;
                if (this.numberone > this.datapayment.length) {
                    this.numberone = 0;
                }
            }, 2000);
        },

        /* 
        网络货运平台占比数据
        */

        // 企业下拉框
        slctvalprises(DataTime){
            this.timer = null;
            this.optionLogins = null;
            this.number = 0;
            this.slctvaluesctp = DataTime;           // 选择企业的id
            let date = new Date();
            let year = date.getFullYear();
            // this.sleectvalue = year;               // 切换企业重置时间下拉框到当前年份
            // this.enterprises(year);                // 重置年份
            if(DataTime == -1){
            // 全部企业
                this.burdenLoginsTime = [];     
                this.burdenLoginsData = [];
                this.allLoginsData = [];   
                this.ifshowslc = true;
                this.buildinvoiceDate.companyId = null;  // 选择全部企业传值companyId为null
                console.log(this.buildinvoiceDate)
                this.burdenLoginsapi(this.buildinvoiceDate);
            }else{
            // 其他企业
                this.ReservedTon = [];                   // 开票数量吨数
                this.AllReservedTon = [];                // 运输发票吨数
                this.sssReservedTon = [];                // 网络货运吨数
                this.ifshowslc = false;
                this.buildinvoiceDate.companyId = DataTime;   // 选择其他企业传值companyId
                this.burdenLoginsapis(this.buildinvoiceDate); // 各家企业吨数图表
            }
        },
        // 年份下拉框 
        enterprises(id){
            let date = new Date();
            let year = date.getFullYear();
            let month = date.getMonth() + 1;
            if(this.slctvaluesctp == -1){
            // 全部企业
                if(id != year){
                    // 不是当年
                    this.chartTime = '物流简报('+ id + '年01月-12月)';
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-12-31"
                    };
                }else{
                    // 当年情况
                    let lastDay= new Date(year,month,0).getDate();   // 当月最后一天
                    month = month<10? '0' + month:month;
                    this.chartTime = '物流简报(' + id + '年01月-' + month +'月)';
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-" + month + "-" + lastDay
                    }; 
                }
                this.burdenLoginsapi(this.buildinvoiceDate);
                
            }else{
            // 其他企业
                if(id != year){
                    // 不是当年
                    this.chartTime = '物流简报('+ id + '年01月-12月)';
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-12-31"
                    };
                }else{
                    // 当年情况
                    let lastDay= new Date(year,month,0).getDate();   // 当月最后一天
                    month = month<10? '0' + month:month;
                    this.chartTime = '物流简报(' + id + '年01月-' + month +'月)';
                    this.buildinvoiceDate = {
                        companyId:this.DataTime,
                        xzqhCode:sessionStorage.getItem('xzqhCode'),
                        startTime:id + '-01-01',
                        endTime:id + "-" + month + "-" + lastDay
                    }; 
                }
                this.burdenLoginsapis(this.buildinvoiceDate);
            }
        },

        burdenLoginsapis(buildinvoiceDate){
            this.api.statistics.countNetLogist(buildinvoiceDate)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.buses.length != 0){
                        this.burdenLoginsTime = [];
                        for(let iev in res.data.data.buses){
                            this.burdenLoginsTime.push(res.data.data.buses[iev].month)                           // 时间
                            this.ReservedTon.push((res.data.data.buses[iev].logisticsCount*1).toFixed(2))        // 开票数量吨数
                            this.AllReservedTon.push((res.data.data.buses[iev].transCount*1).toFixed(2))         // 运输发票吨数
                        }
                        this.api.waybill.countMonBuses(buildinvoiceDate)
                        .then(res=>{
                            if(res.data.code == 200){
                                for(let key in this.burdenLoginsTime){
                                    if(res.data.data[this.burdenLoginsTime[key]] != 0){
                                        this.sssReservedTon.push((res.data.data[this.burdenLoginsTime[key]]).toFixed(2));
                                    }
                                }
                                if(this.ReservedTon.length == 0 && this.AllReservedTon.length == 0 && this.sssReservedTon.length == 0){
                                    this.$message({
                                        type: 'warning',
                                        message: '该企业物流简报暂无数据!'
                                    });
                                    this.ReservedTon = [];                   // 开票数量吨数
                                    this.AllReservedTon = [];                // 运输发票吨数
                                    this.sssReservedTon = [];                // 网络货运吨数
                                    this.burdenLoginss();  // 各家企业吨数图表
                                }else{
                                    this.burdenLoginss();  // 各家企业吨数图表
                                } 
                            }
                        })
                    }else{
                        this.$message({
                            type: 'warning',
                            message: '该企业物流简报暂无数据!'
                        });
                        this.ReservedTon = [];                   // 开票数量吨数
                        this.AllReservedTon = [];                // 运输发票吨数
                        this.sssReservedTon = [];                // 网络货运吨数
                        this.burdenLoginss();  // 各家企业吨数图表
                    }
                }
            })
        },
        burdenLoginss(){
            document.getElementById("burdenLs").removeAttribute("_echarts_instance_",'');
            var esA = document.getElementById('burdenLs');
            // esA.removeAttr("_echarts_instance_").empty();
            let chartsmyChartLoginsn = echarts.init(esA);
            this.optionLogins = {
                title: {
                    text: '单位(吨)',
                    top: '8%',
                    left:'3%',
                    textAlign: 'left',
                    textStyle:{
                        color:'#fff',
                    },
                },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(0, 0, 0, .6)',
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
                                color: 'rgba(255, 255, 255, .3)'
                            },{
                                offset: 0,
                                color: 'rgba(255, 255, 255, 0)'
                            }])
                        }
                    },
                },
                legend: {
                    selectedMode:false,
                    data: ['开票数量吨数','运输发票吨数','网络货运吨数'],
                    textStyle: {
                        color: "#fff"
                    }
                },
                xAxis: [
                    {
                        type: 'category',
                        data:  this.burdenLoginsTime,
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
                                backgroundColor: 'rgba(111, 207, 151, .4)',
                            },
                        },
                    }
                ],
                grid:{
                    x:75,
                    y:80,
                    x2:55,
                    y2:50,
                },
                yAxis: [
                    {
                        type: 'value',
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
                    }
                ],
                series: [
                    {
                    name: '开票数量吨数',
                    type: 'line',
                    smooth: true,
                    color:"rgba(189, 159, 102, 1)",
                    data:this.ReservedTon
                    },
                    {
                    name: '运输发票吨数',
                    type: 'line',
                    smooth: true,
                    color:"rgba(102, 173, 189, 1)",
                    data: this.AllReservedTon
                    },
                    {
                    name: '网络货运吨数',
                    type: 'line',
                    smooth: true,
                    color:"rgba(102, 189, 142, 1)",
                    data: this.sssReservedTon
                    }
                ]
            };
            this.optionLogins && chartsmyChartLoginsn.setOption(this.optionLogins);
            window.addEventListener('resize', () => {
                chartsmyChartLoginsn.resize();
            });  
            this.timer = setInterval(() => {
                chartsmyChartLoginsn.dispatchAction({
                    type: 'showTip',
                    seriesIndex: 0,
                    dataIndex: this.number
                });
                this.number++;
                if (this.number > this.burdenLoginsTime.length) {
                    this.number = 0;
                }
            }, 2000);
        },

        burdenLoginsapi(buildinvoiceDate){
            this.api.statistics.countNetLogist(buildinvoiceDate)
            .then(res=>{
                
                if(res.data.code == 200){
                    // console.log(res)
                    let allinit = [];
                    let obj = {};
                    let o;
                    for(let iev in res.data.data.netlist){
                        // console.log(res.data.data.netlist[iev])
                        obj[res.data.data.netlist[iev].simpleName] = res.data.data.netlist[iev].lrate;
                    } 
                    // console.log(obj)                   
                    for(let key in res.data.data.translist){
                        if(obj[res.data.data.translist[key].simpleName]){
                            o = {
                                simpleName:res.data.data.translist[key].simpleName,
                                lrate:obj[res.data.data.translist[key].simpleName],           // 网络货运
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
                    for(let key in allinit){
                        if(allinit[key].lrate*1 != 0 || allinit[key].yslrate*1 != 0){
                            this.burdenLoginsTime.push(allinit[key].simpleName);
                            this.burdenLoginsData.push((allinit[key].lrate*100).toFixed(2)*1);
                            this.allLoginsData.push((allinit[key].yslrate*100).toFixed(2)*1);
                        }
                    }
                    this.burdenLogins();  // 网络货运平台占比
                }
            })
        },
        burdenLogins(){
            document.getElementById("burdenLogins").removeAttribute("_echarts_instance_",'');
            var esAS = document.getElementById('burdenLogins');
            let chartsmyChartLogins = echarts.init(esAS);
            this.optionLogins = {
                // title: {
                //     text: this.chartTime,
                //     top: '0%',
                //     left:'2%',
                //     textAlign: 'left',
                //     textStyle:{
                //         color:'#fff',
                //     },
                // },
                tooltip: {
                    trigger: 'axis',
                    showContent: false,
                    showContent :true,
                    backgroundColor:'rgba(0, 0, 0, .6)',
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
                                color: 'rgba(255, 255, 255, .3)'
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
                            '&nbsp&nbsp网络货运占比：' +
                            a[0].value.toFixed(2) + `%<br/>` +
                            '&nbsp&nbsp运输发票占比：' +
                            a[1].value.toFixed(2) + `%<br/>`
                        )
                        listItem = list.join('<br>')
                        return '<div class="showBox">' + listItem + '</div>'
                    },
                },
                // legend: {
                //     selectedMode:false,
                //     data: ['网络货运占比','运输发票占比'],
                //     textStyle: {
                //         color: "#fff"
                //     },
                  
                // },
                grid:{
                    x:65,
                    y:80,
                    x2:60,
                    y2:50,
                },
                xAxis: [
                    {
                        type: 'category',
                        data:  this.burdenLoginsTime,
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
                                backgroundColor: 'rgba(111, 207, 151, .4)',
                            },
                        },
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
                        // barWidth:'30%',
                        name:'网络货运占比',
                        stack: '网络货运占比',
                        color: {
                            type: 'linear',
                            x: 0,
                            y: 0,
                            x2: 0,
                            y2: 1,
                            colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(111, 207, 151, 1)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(111, 207, 151, 0)',      // 100% 处的颜色
                                },
                            ],
                            globalCoord: false,                     // 缺省为 false
                        },
                        data:this.burdenLoginsData,
                        // barWidth : 40,//柱图宽度
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                color:'#fff',
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
                        type: 'bar',
                        // barWidth:'30%',
                        name:'运输发票占比',
                        stack: '运输发票占比',
                        color: {
                            type: 'linear',
                            x: 0,
                            y: 0,
                            x2: 0,
                            y2: 1,
                            colorStops: [
                                {
                                    offset: 0,
                                    color: 'rgba(111, 178, 207, 1)',     // 0% 处的颜色
                                },
                                {
                                    offset: 1,
                                    color: 'rgba(111, 178, 207, 0)',      // 100% 处的颜色
                                },
                            ],
                            globalCoord: false,                     // 缺省为 false
                        },
                        data:this.allLoginsData,
                        // barWidth : 40,//柱图宽度
                        markLine : {
                            symbol:"none",               //去掉警戒线最后面的箭头
                            label:{
                                color:'#fff',
                                position:"end",         //将警示值放在哪个位置，三个值“start”,"middle","end"  开始  中点 结束
                                formatter: "100%"
                            },
                            data : [
                                {
                                    name: '',
                                    silent:false,             //鼠标悬停事件  true没有，false有
                                    lineStyle:{               //警戒线的样式  ，虚实  颜色
                                        type:"dashed",
                                        color:"#6FB2CF"
                                    },
                                    yAxis:100,
                                },
                            ]
                        },
                        // label: {
                        //     normal: {
                        //         show: true, //是否显现，不显示的话设置成false
                        //         color:'#fff',
                        //         position: "bottom", //显示的位置
                        //         distance: 0, //距离侄子的值 // label要显示的值比如： 20%
                        //         formatter: function(param) {
                        //             return param.value + '%';
                        //         },
                        //     }
                        // },
                        // barGap: '-100%', // 设置柱形重合的重要步骤
                    },
                
                ],
            };
            this.optionLogins && chartsmyChartLogins.setOption(this.optionLogins);
            window.addEventListener('resize', () => {
                chartsmyChartLogins.resize();
            });  

            this.timer = setInterval(() => {
                chartsmyChartLogins.dispatchAction({
                    type: 'showTip',
                    seriesIndex: 0,
                    dataIndex: this.number
                });
                this.number++;
                if (this.number > this.burdenLoginsTime.length) {
                    this.number = 0;
                }
            }, 2000);
        },




    },
    // 销毁定时器
    beforeDestroy() {
        clearInterval(this.routerboday);
        clearInterval(this.timer);
        clearInterval(this.timerone);
        clearInterval(this.timertwo);
        if (this.formatDate) {
            clearInterval(this.formatDate); // 在Vue实例销毁前，清除时间定时器(大屏时间)
        }
    }
}
</script>
<style scoped lang="less">
// 下拉框样式开始
/deep/ .el-input__inner{
      background: rgb(2,28,61);
      color: #fff;
    }
/deep/ .el-select .el-input.is-focus .el-input__inner{
      border:1px solid #fff;
    }
 /deep/ .el-select .el-input__inner:focus{
     border:1px solid #fff;
 }
/deep/ .el-scrollbar__view{
     background: rgb(2,28,61);
     color: #fff;
     border:0;
} 
 .select-popper .el-select-dropdown__item{
    background: rgb(2,28,61);
    border:0;
}
.select-popper .el-select-dropdown__item.selected{
    color: cyan;
}
.select-popper .el-select-dropdown__item :hover{
    color: #fff;
}
/deep/ .select-popper .el-select-dropdown__empty{
    background: rgb(2,28,61);
     color: #fff;
}
// 下拉框样式结束

.FirstScreen{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    padding:0.5%;
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
        display: flex;
        flex-wrap:wrap;
        // 上半部分内容
        .section{
            width: 100%;
            height: 47%;
            margin: 2% 0 0 0;
            background-image: url(../../assets/k3.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            // padding: 0 1% 1%;
            box-sizing: border-box;
            display: flex;
            flex-wrap: nowrap;
            position: relative;
            .sectionssbox{
                width: 7%;
                height: 30px;
                position: absolute;
                top: 5.5%;
                right: 2%;
                z-index: 99;
            }
            .sectionleft{
                width: 50%;
                height: 100%;
                padding:1.5% 0 0 2%;
                box-sizing: border-box;
                #sectionleft{
                    width: 100%;
                    height: 100%;
                }  
            } 
            .sectionright{
                width: 50%;
                height: 100%;
                padding-top:1.5%;
                box-sizing: border-box;
                #sectionright{
                    width: 100%;
                    height: 100%;
                }
            }

        }
        // 下半部分
        .conentLogins{
            width: 100%;
            height: 50%;
            background-image: url(../../assets/k3.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding-top:1.3%;
            box-sizing: border-box;
            position: relative;
            .conentLoginsBox{
                width: 20%;
                height: 40px;
                position: absolute;
                top: 4%;
                left: 2.2%;
                font-size:18px;
                font-weight: 600;
                line-height: 40px;
                color: #fff;
                z-index: 10;
            }
            .sleectright{
                width: 22%;
                height: 40px;
                position: absolute;
                top: 5%;
                right: 2.2%;
                z-index: 10;
                // background: cornflowerblue;
            }
            .tooleft{
                width:26%;
                height: 30px;
                position: absolute;
                top: -80%;
                left: 10%;
                right: 0;
                bottom: 0;
                margin: auto;
                z-index: 10;
                padding: 5px;
                box-sizing: border-box;
                .toolone{
                    width: 30px;
                    height: 20px;
                    float: left;
                    z-index: 99;
                    border-radius:4px;
                    margin-right:2%;
                }
                span{
                    font-size: 14px;
                    color: #fff;
                    line-height: 20px;
                    float: left;
                }
            }
            // .tooright{
            //     width:40px;
            //     height: 40px;
            //     position: absolute;
            //     top: 5.5%;
            //     left: 43%;
            //     z-index: 10;
            // }
            .burdenLogins{
                width: 100%;
                height: 100%;
                #burdenLogins{
                    width: 100%;
                    height: 100%;
                }
            }
            .burdenLs{
                width: 100%;
                height: 100%;
                #burdenLs{
                    width: 100%;
                    height: 100%;
                }
            }
        }
  
    }
}
// 测试九宫格
.FirstScreen-box{
    width: 100%;
    height: 100vh;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1;
    li{
        width: 33.33%;
        height: 33.33%;
        border: 1px solid #ff6700;
        float: left;
        box-sizing: border-box;
    }
}
</style>
