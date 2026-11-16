<template>
    <!-- 南京LED看板 -->
    <div class="FirstScreen">
        <div class="header">
            <p>煤炭大宗数字交易平台</p>
            <div class="TimesAll">
                <span>{{FirstScreenTime}}</span>
                <div class="buttonbox" @click="BackOff">
                    <i class="iconfont el-icon-alituichu"></i>
                </div>
            </div>
        </div>
        <div class="nav">
            <!-- 淘煤网阳光竞价统计图 -->
            <div class="section-left">
                <div class="section-title">
                    <img :src="zmlog" alt="">
                    淘煤网阳光竞价统计图
                </div>
                <div id="paymentTaxes"></div>
            </div>
            <!-- 淘煤网竞价业务流程图 -->
            <div class="section">
                <div class="section-title">
                    <img :src="zmlog" alt="">
                    淘煤网竞价业务流程图
                </div>
                <div class="section-img">
                    <!-- <div id="sectionImg"></div> -->
                    <img :src="jjywlct" alt="">
                </div>
            </div>
            <!-- 供应链业务数据 -->
            <div class="conent">
                <div class="conent-title">供应链业务数据</div>
                <div class="conent-left">
                    <img :src="gylywt" alt="">    
                </div>    
                <div class="conent-right">
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
            </div>
            <!-- 煤矿企业 -->
            <div class="footer-right">
                <div class="section-title">
                    <img :src="zmlog" alt="">
                    煤矿企业
                </div>
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
            <!-- 企业列表信息 -->
            <div class="footer">
                <div class="section-title">洗煤厂企业</div>
                <div id="sectionBoxtwo">
                    <ul id="con3" ref="con3" :class="{animtwo:animate==true}">
                        <li v-for='item in itemstwo'>
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
import gylywt from '../../assets/gylywt.svg';
import zmlog from '../../assets/zmlog.svg';
import jjywlct from '../../assets/jjywlct.svg';

var echarts = require('echarts');
export default{
    components: { countTo },
    data(){
        return {
            jjywlct:jjywlct,
            zmlog:zmlog,
            gylywt:gylywt,
            FirstScreenTime:'',
            animate:false,
            animateone:false,
            // 煤矿信息
            itemsone:[
                {
                    img:require('../../assets/mk1.jpg'),
                    name:'保利裕丰煤业有限公司',
                    nameone:'所属地: 临汾市乡宁县',
                    nametwo:'主营: 低硫焦煤',
                },
                {
                    img:require('../../assets/mk2.jpg'),
                    name:'中阳荣欣焦化有限公司',
                    nameone:'所属地: 吕梁市中阳县',
                    nametwo:'主营: 低硫焦煤',
                },
                {
                    img:require('../../assets/mk3.jpg'),
                    name:'沁新能源集团股份有限公司',
                    nameone:'所属地: 长治市沁源县',
                    nametwo:'主营: 优质主焦煤',
                },
                {
                    img:require('../../assets/mk4.jpg'),
                    name:'保利平山煤业股份有限公司',
                    nameone:'所属地: 晋城市沁水县',
                    nametwo:'主营: 低硫无烟煤',
                },
            ],
            // 洗煤厂信息
            itemstwo:[
                {
                    img:require('../../assets/xm1.jpg'),
                    name:'介休市佳晟煤化有限公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 原煤',
                },
                {
                    img:require('../../assets/xm2.jpg'),
                    name:'介休晋兴达洗煤有限公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 原煤',
                },
                {
                    img:require('../../assets/xm3.jpg'),
                    name:'介休市东海煤化有限公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 主焦 / 肥煤',
                },
                {
                    img:require('../../assets/xm4.jpg'),
                    name:'孝义市金汇达洗煤有限公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 原煤',
                },
                {
                    img:require('../../assets/xm5.jpg'),
                    name:'介休市昌鑫洗煤有限公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 精煤',
                },
                {
                    img:require('../../assets/xm6.jpg'),
                    name:'灵石锦源煤化有限公司',
                    nameone:'所属地: 山西晋中市',
                    nametwo:'主营: 精煤',
                },
                {
                    img:require('../../assets/xm7.jpg'),
                    name:'山西平遥兴浩洗煤有限公司',
                    nameone:'所属地: 山西平遥县',
                    nametwo:'主营: 主焦 / 瘦焦',
                },
                {
                    img:require('../../assets/xm8.jpg'),
                    name:'山西晋燃煤炭有限责任公司',
                    nameone:'所属地: 山西介休市',
                    nametwo:'主营: 主焦 / 瘦焦 / 肥煤',
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
        }
    },
    created(){
        setInterval(this.formatDate, 1000);
        setInterval(this.scroll,3500) // 在钩子函数中调用我在method 里面写的scroll()方法
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
    },
    mounted() {
        this.paymentTaxes();                 // 淘煤网阳光竞价图
        // this.sectionImg();
    },
    methods:{
        // 退出大屏
        BackOff(){
            this.$router.push({name:'DPLogin'}) 
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
            this.api.capital.getPtSumPrice()
            .then(res=>{
                if(res.data.code == 200){
                    this.getPtSumPricesData =  (res.data.data.sumPrice/100000000).toFixed(2)*1;
                }
            })
        },

        // 平台合同量
        countSumContracts(){
            this.countSumContractData = 0;
            this.api.listByPage.countSumContract()
            .then(res=>{
                if(res.data.code == 200){
                    this.countSumContractData =  res.data.data.count*1;
                }
            })
        },

        // 资金流水累计
        countYearBusess(){
            this.Capitalflow = 0;
            this.api.chainCap.getPtSumPrice()
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
                    this.logisticsCount = res.data.data.logisticsCount;         // 物流量累计
                    this.newbusCount = res.data.data.busCount/10000;            // 业务量累计
                    this.rateBills = res.data.data.taxCount/10000;              // 缴税金额累计
                }
            })
        },

        // 洗煤厂信息滚动展示
        scroll(){
            let con3 = this.$refs.con3;
            con3.style.marginLeft = '-15.33%';
            this.animate=!this.animate;
            let that = this;
            setTimeout(function(){
                that.itemstwo.push(that.itemstwo[0]);
                that.itemstwo.shift();
                con3.style.marginLeft='0px';
                that.animate=!that.animate;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },500)
        },

        // 煤矿企业信息滚动展示
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

        sectionImg(){
            var chartsectionImg = document.getElementById('sectionImg');
            var myChartsectionImg = echarts.init(chartsectionImg);
            var optionsectionImg;

            // 节点图标
            let appSymbolone = "image://data:image/png;base64,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"
            let masterSymbol2 = "image://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAANEAAADQCAQAAAAzb0zBAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAJcEhZcwAAFiUAABYlAUlSJPAAAAAHdElNRQflCgkSEAzhM17eAAAWf0lEQVR42u2dMXPbxrbHf7mToRsywTVdRE4huHhyY43hURqnMdzEzfWImle9StQnkPQJKH0CyZ9AVPVKKZM0168g3CS3iCJ65MZKYbi40i3MBAnVkI1fgSUELBYgSAEEZOGvGZuEKAjYn87Zswe7Zz/7yDVVjSo1atyiQpVbVIAKt3yfGDAE+sCQCwb06XNBP+8Ln1SfXStEdercoUqd2hXO0qNPT/x7DXQdEFXQucMc1YCNpKEBPX7njA9c5H2T0SoyIhfN/JUsJql69LA5Y5j3TYdVTEQ15tG5m8NvPuOUs2LZVNEQ1VhAp57zVfR4i10UUMVBVGEhJ8uJUkEsqhiI5lhkLvVgIB2d8pbzPC8gb0TT286AC3oM+YshAy4YAkMGvk+4Y6UqiPFTldqUUWGPE07zaqI8EVVY5MFETTbggnP+4pyLAIxJdIvb1PmCOnMT/VyfN7zLw+3lhWgyPH3O+cB56kPNOercnWAg3MfmZNaY8kCUHM+AHjbvM0/a1JnjXmKrOuLtLDHNGlFSPAN+4x2/T+3OptEtvuIe8wmurs8pR7O6rNkiesDS2AZw4eQZQ81zj4Wxn+pzNJsQYnaI5lgaG7mdYfPbTC0nSreYZ3HsEPoDL7N3ebNBVGGJxdhP5G87KtVZHGtPJ1kHELNApPMk1r0NeMObQtiOSjUWWIiN+TJ2eVkjqmLGurc+R7wvLJ5LLbAUi+mUX7KypWwRxdvPzDrclBSPKbO7yQ5RfP8z4M3swtYUtRTr9E74Of1fmRWi2zyLvJVi9z3jVGMpJoTo80PaDi8bRHHjnzNeXb8pHpJqmJG5iCE/pevw0kcU5+D6WIULrKdVXM90lKYTTxtRlWeRA74Tfr227k2lOJd3hpWWw0sXUXQP9CnZj1/RtpRar5QmougQ+1OzH7+ibanPyzQen6SH6AHfKo8P+PlajX6m0SKPI75jXf3e00K0xJLyeI+X1z5+S6Ia/4hweFcOHdJBFAXoU3Zwsip8G+HwrggpDURRgH7mJPuWKZSiWuJKkK6OSH1ZA15+khHcOM1jKkOmK0C6KiI1oD4/3ogeSKWoXmlqSFdDpAZ0U0KEKKUM6SqIogD9eGNChCjV+E6ZY5kK0vSI1IDe8q8bDwigwnMlpCnGSX+b8hL0CECvSkAADPlBmVkwJ1/1MR2iKk8UR9/yKueGKZKiID2nOtmJpkFU5bkisCwByVJDqvCcyiSnmQbRc0W80uNfebdIAaWGVOO7SU4yOaLHCkD9MoqL0FA5ALkbmXZVaFJEDxRPVEtAcVIP48dPofQ0GaKq4oHDTc4kJFOfl4o/4W+Thg2TIKrwXHHUKgGNVU8xeavCs2RhwySIvlX0Qkc3Mlk6uVSLXeoReXFJyREtKLznybWcrpiPjrBDxxaTLDtLiqiqIN7PYu7lJyxVl2COd3ZJEX0TcnIDfsz7nq+ZhoqwoTbe2SVDpHJyP5dhwsTq8Wvo2FhnlwSRysmdfvKzerLRiaJHGuPskiAKO7myF5pe1qTObjyiusLJWWU2YWoNsULHYlfVjkcUTvmdlGOhK+m9wtnF5OzGIVpQOLlfKXU1hb3QXfSoD8cjUgUKR6WTu7KGij/zx1FBQzyi+yEbKiO5dBTuLGpRq7LiEKltqFQ6+il0ZFFtR3GIvgkdOSqHq6mpF5pOXVHbUTSiaijY7vMm7/sCNDoc0PIdMdhjAy3mZwyaiu64ww6NHO8k3Ksr7SgakcqGZh8oNENHDEwavPcdMWmygxNzlnX2eCedy8Rkg/WZ39GlhqE/eaUdRSFS2dCsAwWDDnvs+d6bmKyKdyYmBgA60I05jy7gWIGjDQD2Z3xPQZ0ksaPPI35YZUOzloMONLHZBmAH0/ueC87iKfAQAlYly4XalgaMy+IMeWrIGykkq7Aot7QaUSWUfZ29DYHNCsfAFn+y6x2zxCvT610M4qxoZEPbgaOGsD175ncV1EmogGIIkdrR6aHxUD7BdpdNAFoeDos18TWaWGmiAS0+el97gXOsohO2IbcPOkRXfM1Sqv5oPnhAbUXyeCgPG3K1yzImGns8BcCkI75jSP+rpbMFnqvUPVAmAFtshX7C5t5M71BlRwG3rUI0VxAbcrXGMY7XrYf/ypcBm31AYx045Hvfd93g/AU20KTFNm1QhuAjOTO+v3B/dJc5f+5Bhei+9L6fa2bbZoWu13Bt1rxXq4COCaLhTdaBfQ69n23SBGzRk7XQ2cPGEuAOJRwaDWaPSGVHejyicLh9nnNOwQpcfFO8cj12w/cJDcDX5+gChesiWyJAsFhHx0UflJ7TQHbIe6nN73N0uU1PGFG4BmOR8nKmL/AGt9v/XmAxgMvYzuBAQFlHE9YGK+hsAHKEl6/eSogq3L9MD4URyePbs9xsyPS97uJ4KEw04cxMEa2ZWLjjo673E5robwxfQLGJzYawoXZOd6XSOefSMGc+GlE99Ig2v4cPHd/rp1iiR9lDB/aFc9uixR4aj+hi4B/CWjhogIODhgYcsgvs4tDiKQaOFIhrADmNlN5JiO5SGbk6GVH+aZ8oNYQ1/AnAjojbLNp00DhgBR18oQKsADY2BseALcZY0OaQZXZxWAkMebUc7+40VGTRc3UyIl16n2cs5/YWywJNw8vOgd99bbFGB91LCV3KEnd0AMCKzz4cHqGh0WEl5xTQSEN+j3J1QUThPUje5njZWwDoAsYhNm6uoOsb+ViAhSVSqrbCTblBw6aUItrAoYVGh6cSpLhsX5b6RVp34rm6IKJwZq44c30OOcRNib4K5QTWOEaDwKDV1R4GsM0ubgBh0BWwtnAHtgc8FUf0XO+vx0Dt6oI5Ojn1URxAI+nAupeNGwUUo+SqI316jybgMM8x7/iDY/Z8Y58tkZPo5AzH1TBkvyJX50dUDVlRnm5OJTPi+GjQue7r8jU6Ypir0RSZbVlNuoAmeiv3E3Zud3cmvb/jPjvyOzp50DosnBWZADzCAQ4wRHSH95hco+XFbU4gwerQ5TW25+ZGWqGDzuu8bwwI/3FUqHMej+iMomkVsHyN7N7UOk2gjY7JBvve9/dZ5VCAsSPOaLOG6QUmqoaanYahAawuI5LdXH4Xq9aOGLQCaCKCg1V2cR83uEmeHZGTgy2R6omX5UV0X+Z9g6EB7Bz4+6JqKOAulpvbYANo06bBnug9uqzTBhxWsLHYB29GQ1TOWo/8DZpopvwkt/gdKn4rkt1cvxBz5tbFHIM9VrBxRE/TBKBLVwxYR6OeNstsKx6Su+H2QzHg/XvE7zKA/MZFoAq853gfjagINrTnPXrQOabNJg5gsYImOv6nHLDpuSqLe57tjLDoEdGcLB2N/J37f6TH4nf9iOT0af7Bwo54ILeGTgudJg0OeRXo/m0eiZy2Jv7XmWcNMIUz9MuN6zSlCzTF+fLVvyVEdb+jkxGlvc3tpDJEZ/8UG7DYYhVNPEeNlwX4H0u4YLpYIQAGpofLDdy/J1/JrX7nEpEczQ1zR9TlBetsi2a1adJmPdFz0e/FTxzymm7sRCybViC/7QQy5XlIbvUKtVHhTXk7gbMCLNnXOPAC6MtjBgYaWiBE/lPYgg2E5yTEqePLWDheti5P/bfkz6wRosfS09ZMthwroC4nTNq+SSz5XlHwmd3RyNHdkT6Yf7AwG1l5X0BIHyREtdHQVS68NUx0ulLpSx6N3h4hkjMLeQcLN1e/S++/cBHJAfewtKLc1JcWtFRcRPKalg95X+eNlmQeUVZUKj9Jrs5FJFfdLnuiPCWt3HMRycFCaUV5Sorp1Ij+yvsqb7SUiGSVVlQgqSO6ElGeUlpRiahIUgbdMqKyVlaeUiK6Ff+hUnlq2l3ASmUn5bioVJGkdHSlCqwSUeFVIiq8SkTFkzQEKhEVT9IQyEUkD1Un2tSyVLZyEclD1VuTn6hUaiodXeGlRHQhfWjCLZZLpSolItnR1ZKdq1QmShQulH1RnpIMxEUkz4AsHV2eSoSotKI8pUQkO7rbeV/ljZYyXJDnEc/WilpeaQudHel7muLznTEr9UyOE6xu1emIBchazJnGnyd9SRNP1VZUyzi/0AyUA2x6q7SNUOMfiAJMbqHa0dcT75WmOHvXK30WrQbHYs1rkz8Cm4RcyvAgzlKhmcGjoFvujbINu7sYohE/8hGdPVFyyQytkltDEzaxwZ74gob32m3Ctq+Q+kf+QMMMHPkYWOSrscMB+2IVbZtNtninsBedPJYny8bx12gJ2O8SlHqmk4a7bLLHa3ZB1IPb4gnwEE2geyVql9o8pcMuDVEaA+AjL0Llzrq8EK8MbGm1nRkoNrhBCw0Lhw0BTsNGpyPVeHTXkB8rr1/PsDqDbEUXI0R/jflg2mpjsMwucIADaHTRMDkEDDRvuwG3WLOD23eMtC4a/XsPmyOQmuyJ6nN+rUqvbXQ0rzbda17jsMoxK4HlyA994IP6k+z0tfS+N0Ik20zWiPAad40usMFDGqJ29i4PA/V/uwA4XhH1Dt97Nob4/w/ALfwMT3wwbTZxsHyN+lTYWAsjUIyzTVNaVmnyIoc6xPKY9EMUojsJTpaO1rEBA4eWaKKHoviY4S3h3/LqbLv6UvQculiGvyXOtAtseDgM1mnjiHqQo2r2IzXpSvXtHBqBnV3yWQkb6ejkcKFCLbMaQLpXw+eQNaDFPq8wMHDQsTGEe9FZxa1KsoXmc1aW753tVexuYfKCZZrCNhqseju3uNJCcV447mt7WAzyQKSofzFa1K9Y759Z+WfXHem02UcHWji8ACx20FjjHfd8DqjJHp+NOaPGMTbbWKLH2uYhG+yzEbFMv8UWa54T66CL6M6vg5gyHJuh3i4tKepffO57E0R0NzNEXe7hFhUzeYKBg8Mq0GWNd+wp6wRDsIyFq6deTdRHXnGMFTrsAYeRA9wdNrwdWFrss8KxIp7bjyg2o9PKMBQPBQv+GkBybzQ39nRX1xZwzCH7uFXpYY2DiCgK/BFcMMIb1Z4zWGYDjTYmDd7xIlAD0v2pPUy2sWixASJ8f0SHY7YC+0gcRlyDSSvDEhpfSe/P/IjCu/Vm1xsFG83d6Mnyisk8jPzslz470nzHNQyWMTFweMEuDm5ftINbleQ1FhYa62wAa6zTwmKbNg4Gq7zgEVts0WRzbBUgg+zKCtZDqbdA4c0+fWn4qoc2h01fJhp/x6HNR9xx/yENGhENpfPEe635Gu0Y1wptLOa9PJ+DhYGNTkO4pgaHbGNjCZt1gwy3ZtAWbW9AHIVgk0N0sisrKHuuDwyDtVHPQxmG7LVKGweY5xWwA2zymh0syZnooljgmu9IQ7zqskIXm7aoAOmXxjaWyMY5PBJHHUye0MS1XQs3Ydr2ei8r4mod3AFtVlLWSfcjOpOqz+iZXYpbQ1GjgclTGixjsk8Hg0fY7LIu8gwjrbLLNv6d9IKObvRZx4cQLnc7cgJHNd4BXeHo3NLRGxz6KnAFzxKUkVlBtIq6gPDn8gHpR7Ipv2nQoYuBwRo2oLGPhS7eOTzyXYsGtHnBLsu8922Vq02dhXbjPwcAk2WaJA+jzUByKl2FA7QeBBFF1IlOXTodMWjdpUMXm9cie9DELcOveVvZ6KzjsCnGMO98CZnwsyXQpPTqPGo56Jg8pIFG1wswomWIuqkaO2Q3oJXd3JlqW48zCdH9TKrSOXRFeLvBLk2esOy5LUf0Gq+893Zo1xTYEtkFWzqu+cIJ931YTbGrnnsV3UQNrvkc7GZmoyL5D0qMSy+zCwBz0gYt8EMhKg27anqFMV2bkzcXaKBL7kqnGdp0V6MZW/9eLUPgzq6s4DzPpCP/6w56gojgf6So7qZUd8xfzyQr8kqfyhOG5aTP/XIK/kxUjXJzYURywc3KTBJBpcJbtXodzN9C35D7nkVKZS95V3ffVq3hlRGyHd2dSZbhZiu8q7uvwwkjCufl9Lzv4JOXvKv70B9vhhENFa6uDBmyVHhXd9u/VkW1BOwX6X2l7I8yVbh1A55Mheg89HivRJSldOn9WbD91Qsp5Z0oKyFTLJWWFuJCBYhCdBqa5b1EqWwkt2w/GaIhb6QjtdKOMlHYho7kj0StGD8p7WgmCttQKG0dhUhlR2XQkLaWFP1QaEpPdN0FlR2V46M0FR4P9VVzF6MRhe2oHB+lq2+S2FB89ZKwHS2Wa8lTU0IbikeksiMz7zv7ZPQsdOSNemppfA2gsB3dLZ8fpaKF0PODftTE0nhEQ34NHTPLoOHKqiqGMEdRHx5XSesklK+rlSOkK+sbxcZ4ketQxhc7+yl0ZLF0dlfSgiJT8zL64+MRnSt8ZOnsppfKyZ3ErUFJUjLwSFE6o3R20yrs5PqKHt+nJIhUQcNimVadSg8U7XYUv41KssKbJ4o5qd+Ww9iJVZVWsgKcjluwmrQ2qqWoQvys7JEmUiU0HRv60cH2SEkR9RXT0+tljzSRTEVlpZ/HL1ZNXmH4vSKyWywTq4m1pJjsdpRk8v8kRaCPFMQfl2OkRNIVHieBk4PJEA15qYg9npVhw1hVpXVPAIPRyodxmqyUek8Rfld4XkKKVZXnijqZCXohV5NWuz9R9Ei1MraLUZXnijDhKHltmMk3JDhSFBOs87yEpFSFZwpAZ8l6IVeTIxryUmGidR7n3RoFVIXninUl/cnWnE+zrUdfGTbcV3SJN1tRgH6crHDPdDuv9JQrYEtIfqkBofRBsZp2c5zTSEhlnwTRgKzJywLLK8Yn0ZIyAdTjhxu/u3KVZ0pAR5OECSNdBVEUpD4/hHZEuklSh9lTAroqohJSWLeVYfbUgK6OKArSkH8WqO7J7LTAY+WOG1MDSgNRFKQrXdY1VSYtkQai6Es74ejGhA6VyH1arvinmg6iaEg3pVeK6oFS8CVpIYqGNOSnzMpJF0UPWIrY8ymFaufpIYJ5zIgLPeWXT9aWoh3cgB/T2L8mTURQ4x8R5t7H+iQjvDnljAT3jifMxUUpXURQ47vImkEnnHxStlRhKXLuxhn/Fz87LrnSRgTwOPLC+5M8yiq4okZAkHKhxSwQwWLM06NPoV+6zbeKCnKuBvyc7p9hNojieiWAI95eW0xx7g16kz9sGKesEI27levp8ios8iBmS9UTfk2rB7pUdogAFhSVBS51vTCNw5NZzJotIneZS9waiuuBaRwesHmVvv24yhoRjLMlF9NZYfum8XgyHvPNAtF4W4IB7wsY6c1xj/8as513Jv2PX7NBBONiPFdnnAaLTuamCgvokYH15fUeZZ8zmR0iGO/ywLWntzkmiyrUuc/8GNuZYS86W0RJPPuoAc5nDiopHBjwhjfZurdLzRoRJOmZLpviHJtepvudu1c0j67YyU59TTPEA/kgchslKSaAPr1MUNWZ404iu3E1czyQHyJwi3kujO2b/BrQ43f+okdv6qCiRpU6d6hTTYzG/d054IF8EblKEkKoNKDHkAsG9OkDF8AgAK7CLaDKLSpU+IIqtQmxjHTGKe/zwANFQAQwx/2C1nEY8Bvv8n0YWQxEADXmWCzUFiJn2PyWl+1cqjiIXNWY537uoM44z6ffUaloiFzVmOdeLmvRz7B5P5ONuxOrmIhc3eIr7in29slCfd7zb/5TFMvxq8iIRqpxm6+pZ2BVbhBfUDQjXQdEl6pT5WvqVK7UW/XpccEHzovl0KJ0vRD5VadCjRq3xFinymgkNJI7ShoypM9QjJ8urgcWv/4fzVKiOxGSYkAAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjEtMTAtMDlUMTA6MTY6MTIrMDg6MDBstkADAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIxLTEwLTA5VDEwOjE2OjEyKzA4OjAwHev4vwAAACB0RVh0c29mdHdhcmUAaHR0cHM6Ly9pbWFnZW1hZ2ljay5vcme8zx2dAAAAGHRFWHRUaHVtYjo6RG9jdW1lbnQ6OlBhZ2VzADGn/7svAAAAGHRFWHRUaHVtYjo6SW1hZ2U6OkhlaWdodAAyMDhzDJ1bAAAAF3RFWHRUaHVtYjo6SW1hZ2U6OldpZHRoADIwOZf6/ZAAAAAZdEVYdFRodW1iOjpNaW1ldHlwZQBpbWFnZS9wbmc/slZOAAAAF3RFWHRUaHVtYjo6TVRpbWUAMTYzMzc0NTc3MsYz5TkAAAASdEVYdFRodW1iOjpTaXplADc1OTdCQrfU/C8AAABGdEVYdFRodW1iOjpVUkkAZmlsZTovLy9hcHAvdG1wL2ltYWdlbGMvaW1ndmlldzJfOV8xNjMxNzU3NTI1MzMyMTQ0OV8zOV9bMF3b+16MAAAAAElFTkSuQmCC"
            let slaveSymbol3 = "image://data:image/png;base64,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fTNX0AAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMS0xMC0wOVQxMDoyNzoyMSswODowMKXxUEIAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjEtMTAtMDlUMTA6Mjc6MjErMDg6MDDUrOj+AAAAIHRFWHRzb2Z0d2FyZQBodHRwczovL2ltYWdlbWFnaWNrLm9yZ7zPHZ0AAAAYdEVYdFRodW1iOjpEb2N1bWVudDo6UGFnZXMAMaf/uy8AAAAYdEVYdFRodW1iOjpJbWFnZTo6SGVpZ2h0ADIwOHMMnVsAAAAXdEVYdFRodW1iOjpJbWFnZTo6V2lkdGgAMjA44P3NBgAAABl0RVh0VGh1bWI6Ok1pbWV0eXBlAGltYWdlL3BuZz+yVk4AAAAXdEVYdFRodW1iOjpNVGltZQAxNjMzNzQ2NDQxZOT29wAAABJ0RVh0VGh1bWI6OlNpemUAOTM5NUJCARq8kQAAAEV0RVh0VGh1bWI6OlVSSQBmaWxlOi8vL2FwcC90bXAvaW1hZ2VsYy9pbWd2aWV3Ml85XzE2MzE3NTc1MjY0NDQyMzU3XzVfWzBd4g/dPQAAAABJRU5ErkJggg=="
            let proxySymbol4 = "image://data:image/png;base64,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Pl2BRwF5iBO/qJMTFEzRdXEYgixPKnokDYKCxkUwGNe5vEK4z0rAvn+deFNy4UiHIFjgh4DbPOUfX6wcFQITz/q8TKeYW1oEQYk/Bgh6j6M03qW5wwJG4BMnjrkFIT2CoMQPgZVOvzQ5KrASmBZ6OU7MLQhpEgSqzLEhviQ5ChrxAJykW0o0bYJgmT8E3rzJ36+8HBUxFAVcbPR5k7ZTlD5BYdYI7Anfq0pSgeWAIQWkanlGyIKgcA19VZVdgWUpb9eLE94opjInRjYEASyxEiJHV42k8Kfp8yqr6GN2BEFJuaJyhKtBUoEllkM30DZ5lYXs2MiSIIh682ySLubWJkWptUxlx0bWBEXLEcznKGmBu/x7KDmZ2R0vsicIovw6GxecYs5F/LvAEnqgKz2636OMd7cApkUQRCs7gD7nvJ/GYwdigXssRsgN9Hk9Les5PYLATsOIIskuCxl356C0UKAcQ6UB9HnHu6wV2wjTJSieRRrCxOQy0x1d7TtaRKccg5qpkwPTJ8jukpUYamSIHpdc8CF1ososcDvBfcyAHJgNQWCX2Yyj7ry44BO/csnl2K5EgTJlvmaBYmxiYGbkwOwIshHHcVChz2d6fKZPjx7wGej7aCtwCyhyiwIFvqZIKSEpQ1xg8vNsyIFZEwS235RV0drJ0OcDJzP1KZkHggBKLEiFm2aLGcvNCPNBkI0Si9yb+e6TF3yIqp8zTcwTQTZKLHJ3JtJ0gcn5/FBjY/4IsnGLf+Mu5anIU49z/o9/zoNCkzGvBA1R4lt+TzkDiepzySc+cj6fxAwx7wSNUKbI7ylPuEdyj0s+85EP86bKgnB1CPKiTIESJW45Y5siw5HPEPaoaMCAHgNnvPT5qpDixf8Dm7jxxVBYOGMAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjEtMTAtMDlUMTA6Mjk6MjArMDg6MDAdT2tFAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIxLTEwLTA5VDEwOjI5OjIwKzA4OjAwbBLT+QAAACB0RVh0c29mdHdhcmUAaHR0cHM6Ly9pbWFnZW1hZ2ljay5vcme8zx2dAAAAGHRFWHRUaHVtYjo6RG9jdW1lbnQ6OlBhZ2VzADGn/7svAAAAGHRFWHRUaHVtYjo6SW1hZ2U6OkhlaWdodAAyMDhzDJ1bAAAAF3RFWHRUaHVtYjo6SW1hZ2U6OldpZHRoADIwOOD9zQYAAAAZdEVYdFRodW1iOjpNaW1ldHlwZQBpbWFnZS9wbmc/slZOAAAAF3RFWHRUaHVtYjo6TVRpbWUAMTYzMzc0NjU2MCAXztQAAAASdEVYdFRodW1iOjpTaXplADkxMTRCQr+srSkAAABGdEVYdFRodW1iOjpVUkkAZmlsZTovLy9hcHAvdG1wL2ltYWdlbGMvaW1ndmlldzJfOV8xNjMxNTkyMzM2ODY1MzYyM183OV9bMF0Oejq5AAAAAElFTkSuQmCC"
            let proxySymbol5 = "image://data:image/png;base64,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+z2UFPmhpmfID8nL9U6zf9BAOQYuKXKhpodp6Jl2g6e+0uUuKYJwejpJfLcxpt2B61Rp7EqKwuh5H54WEo34oR4Zq8oFkiP+lmlBzSJjkceKcQ8cT58NNQtBYRRd1xjdfUW+G8xEz6wEhVM0U6cuJR6GLNKZ8UnMSlA4RT06KaXbFx9VDClbx8bMumR2glSlgGwMeDud53LJsIgRsqlPZ1rXYIw0CAo3jtdB1YVZniFv03CW0iEIavwp5C26yqouXLWlpj3SIiiKogHHs4t6AXGf1RC90U8y4xON9AiCCt+Hlr045acrJUfhsgO/8H/p7U6bJkEQ7tNdLTkKl53UrW7aBMES34duVDiYMpuuSFhkNVR2UnIM/EifIKjxOCJT4TIruyglDj3epb/7XxYEQZSqgxEnKS0EmycqrEhr3/04SalQgYCsCII6jyP2ZL1cFmkSOQM6WQWIsyPILn8RVURmwDG9wkvSJHLghJ/T89pEZEkQRIVBbAwwUyuRkT4mk5Oh7NjImiCItkc2TvmlcHNIi9xjKZIcOOZTdrJjYx4EQY3vJy7ZKI7Cq7CMHupKu8jEZ5MxH4IAlpUrzEScYqay2H06VKizSn2C3MxBsY0xP4IgLklDzjH95fHngAp30CeqNLt3n+YZoZ8vQVBjmeUYJAH0MLnIfEapziJ6DKkBm5zMrU4Q8yYIkpEEA/r06KeuUuoscjuWzNjIgRzIhyCAmyzFUnd+9PjCv+nTn1r51ahS5zZ1qrGJgdzIgfwIsrEsFWyKhyFfGfCVIQMGwFdgGKCtwk2gyk0qVPiWKrWEpLjo8SlHtyVnggDqrCRQNPPEkL9zlvf4LH+CwFZ408lSVuhh8vd8lFoQxSDIRo1FVnJfUtnjIi97o0KRCLJRYzEnaSqM1PhRPIJs3OQ/uEt9LvI04JwzvhSNGhtFJchFjVv8njq3UncjhvT5wj/4ZzGJcVF0gsaoU+U2dWozSdWAPl/5zMU8Ap1p4P8BVxqpKF8/s5cAAAAldEVYdGRhdGU6Y3JlYXRlADIwMjEtMTAtMDlUMTA6MzI6NTIrMDg6MDC5e1dCAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDIxLTEwLTA5VDEwOjMyOjUyKzA4OjAwyCbv/gAAACB0RVh0c29mdHdhcmUAaHR0cHM6Ly9pbWFnZW1hZ2ljay5vcme8zx2dAAAAGHRFWHRUaHVtYjo6RG9jdW1lbnQ6OlBhZ2VzADGn/7svAAAAGHRFWHRUaHVtYjo6SW1hZ2U6OkhlaWdodAAyMDkEC63NAAAAF3RFWHRUaHVtYjo6SW1hZ2U6OldpZHRoADIwOOD9zQYAAAAZdEVYdFRodW1iOjpNaW1ldHlwZQBpbWFnZS9wbmc/slZOAAAAF3RFWHRUaHVtYjo6TVRpbWUAMTYzMzc0Njc3MtSGStcAAAASdEVYdFRodW1iOjpTaXplADc4NzdCQq+bT60AAABGdEVYdFRodW1iOjpVUkkAZmlsZTovLy9hcHAvdG1wL2ltYWdlbGMvaW1ndmlldzJfOV8xNjMxNzU3NTI5NTY4OTc4NF83Ml9bMF1RUCgoAAAAAElFTkSuQmCC"
            let proxySymbol6 = "image://data:image/png;base64,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LNKHrkFcUFT0N90Fns+7ZRYCBaWKgxeTj7RNFhcFsaNcZGljffTVlaMIOptmktssFEWR1OMvH3Uwq3azpyRo1ipgZ4ooCVY/6mA2Y4JmlSKIMop9nuWcL38xscbXyoobUxOUBUVRJMGLTHvBbwPuKSYwwkwEZUNRlE36uNLTFHgQ8arORFBWFEWT1Oe7j8J1iLJAMxOUHUVxJJ1lti5pUXGPnYiaTxko++wogq3IYk3nnH2wsVK0ghvwLIspoVlS5PZ4q8X9Q1V4a5iRv3jCvrgoZEsRFPl9ZFXYmbXygqHATmTSkN7kfXFRyJqi+Af/kOySzoPIlzFTTzZ7iiDadQC44Mdbr/JWMCOzmA94mW08mA9F8Qrvdqu8AtuhhXFj5NCrkhdF7tT86Jwnt1Xl7cTQA+e8TDPtajLkRxHAFjsxsnTbaIr/NQOe51PQIF+KoIipyOI+xm2hKZ4euOJ/spcfF3lTlObnLTZN8bYHcpQfFzdBUZJdAuhzvYCe3ip3+W1COe9c7I8fN0MRwIZyZUAQV8mrbW4IBXS2EssDXXGWfzWAm6MIklUeuPL0Zq5lENbYTJSdHOKfKNwsRVDkXqpMq/MhKh05MOA1r/NVb2PcNEWQxjKNm+Iam6vcbVSRDdZZS0HODdMD86HIbZS0NAH06HFFL/PZEEXWWGctUfmOcOP0wPwoApem9M0DMKDHL1wxoDf1+FOBIqt8weeK+o/x954DPTBfisB9j5NdCBVcmnoMeM+AIUMGIdoK3AFWKFLgMwqsUpjqXj3e8PM86IH5U+Rig7sLkB5XjRtxrOOwGBTBLPKUF66w5yc7YywORS5WucvG3Inqc8HVopSoWzSKXKxyl7W5FLe64hp7sebRLiZFLiZ1iWdBn2uueDt/tRbGIlM0wiqrrCvKMGWBHte8U2XHXhzcBopGuMPnrPIZq3w+UUQjo0+P97yjx/tFlBoZt4kiP+6wQoFV7lCkCKwwioL8cCMlN27q02fALwxvAy1+/D8utguCF3WayQAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMS0xMC0wOVQxMDozNjoxNSswODowMPF9x0wAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjEtMTAtMDlUMTA6MzY6MTUrMDg6MDCAIH/wAAAAIHRFWHRzb2Z0d2FyZQBodHRwczovL2ltYWdlbWFnaWNrLm9yZ7zPHZ0AAAAYdEVYdFRodW1iOjpEb2N1bWVudDo6UGFnZXMAMaf/uy8AAAAYdEVYdFRodW1iOjpJbWFnZTo6SGVpZ2h0ADIwOQQLrc0AAAAXdEVYdFRodW1iOjpJbWFnZTo6V2lkdGgAMjA5l/r9kAAAABl0RVh0VGh1bWI6Ok1pbWV0eXBlAGltYWdlL3BuZz+yVk4AAAAXdEVYdFRodW1iOjpNVGltZQAxNjMzNzQ2OTc1QHzyfgAAABJ0RVh0VGh1bWI6OlNpemUAOTE4N0JCwOJc+gAAAEZ0RVh0VGh1bWI6OlVSSQBmaWxlOi8vL2FwcC90bXAvaW1hZ2VsYy9pbWd2aWV3Ml85XzE2MzE3NTc1MjUwMDMyNjEyXzI4X1swXf0jTsYAAAAASUVORK5CYII="
            let proxySymbol7 = "image://data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAANAAAADRCAQAAAAX8fRaAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAJcEhZcwAAFiUAABYlAUlSJPAAAAAHdElNRQflCgkSJDECceE4AAAZFUlEQVR42u2dv3PbSJbHP+Ny0Yk4xpoTnJQIDlZKRjWc0iZ2IiiZSdYlui66SNRfIOkvIPUXSP4LSCW3ye1JUzORLyCdjJPVmLNyMrrAUCJtYHqxRydkMhegAQLdDRAkAYGS+HWVRTaBBtBfvNevf7z3vvidOWYZD/O+gYlQpECJAo8oUqBA0S930WcADOgDn+nR5zN9ugzyvvHx8cWtkaAiJUoUKbHAownr6POZHl269Ojm/UDJMOsEFVhiiSeUJiYlCn26XHM163I1qwQVKPGURUo3cK0u13zgOu9H1mP2CCqwgpmBxIzGFRdc8TnvBghjlghyqVnK+S5mjKbZIGg2qAniigvsWeid8idokaf8MQeFNhp9LjnP29rLl6Bl1iaSmz4Dugzo0QN6DMS73hO/u+OhBaBAkUcUxchpElxxwUV+TZQXQQXW+HpMuenS5SM9PvlEjIsSBUosURyTrB5nefVLeRA0Hjl9MVr5RD/Vu1ikxBKLie+jxwW/3TxJN03QOORcYXOdeR9QYpGnLCY6NgeSbpKgpOT0ueSKy5QlJh6PWGaJ5QR31+PsJvukmyPI5Jk/mRmFPl3OUldm42CZVcyRR90gSTdD0BOej7TWrrD53xypGaLIIusjX6aPvL4JZZc9QQXWWYs9os8lv83cXNgiq6yMOOacs6wHs1kTtIgV+y72ec/7mZAbHYp8jRl7/5kruywJGiU7s03OECsjFN4Ff8tO2WVHULzs3BZyPMST1OMtdjYXzoageNm5beR4WGclhqQzzrK4aBYELfB9zFTKOb/cQnJcFFmPMRx6/Ji+qkufoBWeRQ73rjibOWttXBT5LvL1G/Bz2iZD2gQ9i1RtfX7hPN2L5Ya4HillVZcmQQW+ixyO2ry5tYpNhzhld0U7PVWXHkELvIh4q/q8ycrGyRXRcpRib5QWQU94EdHzXPE/d0p2goiWox6v05mHT4egFSxt+V3qd6KwxrOIX9ppGAxpELTOurY8tbdoxlHkzxGqLgWDYXqCoui5a2ZBHApYEYsUU1M0LUFR9GQ0rp5hZNQS0xGkv6k+b/PcB5MbooboU1E0DUF6eu5Lz6NDie+0vdEUFE1OUBQ9P028KeouIMpgmJiiSQnS09Pl9b2mB1Kn6MFEN7ESQc/9lh4XPf6qVfHrIxfQtZiEoCfaYWmXn+6NWR2PAT9qKbIm2Xw8PkELfK8pndMTRBRFL1gYt6pxCdJPic7pkaGnqMALCuNVNC5Bun0GvTk9GugpKvLdeNWMR9C6Zr1nTk8UBlqbdilyclWLcQjS2W73fdwTD33rrI1jzyUnaEHDfH9Ozwj0eK3RL8+TGwvJCdItyL2Z0zMSXd4qZQW+T2osJCVI55lwdicXstPHhWYOoRQx960gGUGmZq+Ofe8WFCaH7lVeS+Y0loQgXe/T403eT32r0NZ0BlYSNZeEoD9p1NvctB4PA42xUExicI8maEVjFL6dGwdjo8svStnqaDU3iqAFTWd2ced36mSDc01PNFLNjSJoTVFvvblxMDHUnqg4ypqLJ2hBY72dzdXbxBjQVspGWHPxBKkLC7mGRbkDuNZ0D7EyFEfQirLANFdv0+NMseaW4ubm4ghSmZ2rt+kx0Ez9rEebCtEEqTv3u3P1lgouFCe2YrTDaBRBCxqxe533k90Z/KyUrEXJUBRBq4r8XMzVW2roKqZCIUqG9ASp8jM3D9KFaipEyJCeoLskPyZVjefBCYdUcryrAe+lkggZ0qUGmBX52WNrrOMdXmpKa1SBHZqBMosKUOY0h6fycK6EZlvjXI38oyNoSZGffNRbOcJvLwq2pswUctIOlbr1/pDLU3kY8F4ayBRYU1taR5A8/unlal63E6zbGpHqahsDaEp1bIma80UiGVIJWpkR+fFwHFJOepgRBJlUATiQSsuATSfX59LLkCkLg0qQ3FXlKz/wOEEExKgjtjFR5acKQFt7ls1NQpWhFbm1ZfeTRV5IleTnzNhke6zjbZ6Gvpu8w8DhW2zAwBHlHyIJlWvIHqobz4/heQbZzF5Vqri90zs1DOAVNlDhnZCcagKJvDmcK+Mh6e7CKq6gGNj5j3/2hTFsBsziClu8GtGHWFQBmzoAh5g0cDilBsCpL08eqqCUZY8BF1KnshoOsxkmyFQqyH9x2xEK6pAyNZ7iALuYmLHqyKQBwCYANUzA5lTIj60ZMVXJgyCwJYIKrAZbPUyQrOA+zog7sEGLMvAKB6+5DT4ox7XZAaBMQxCxi4Epxj2bmEJ+DpQzzdye7ZpraU11OYqgBWXx9T2zgDInmMABdaAilJaBoRxpir8GZfF9z//tAJs9QVsz70cK4YPU7ksUhkouaCSYyql23vcOGLwL0FOmgWttfaH5tynOaQtl5WCLT23qwBE72GxizpShcKGZOPURlKBV5cR8Exw9BqDDAbsccEQN2BNy09Ac/0bIFsBLwMbG5ANgC9UHTU6xaACbISPDrdXO5TkH/IPlUElAooYELSg7ED6QLwzxt04Tm13qgI2BK+vt0LEVjFDzur+atADYCfzmYGEA76hyrFwrH/xdImiJomc9PwwUhjHgMteb9uAANmWOAJtNKmxh+eaAB0vbxF7f1Q6V7uFQA5oQoChPdOlLMwqmZygM+yBZwdl537Vo8n8CZVqAwyY2RwBU+T30z9Sc36CM13cZlKkK0wHqwo5r+iXutfJ6JVVh8CXKk6CCYsHlreA8o+WSXUHKfuCleSNZYoeKBNXEyGaZdxiirgO/36ljsAu0/ImgfPGbNEXgW3IeQeruxvzDJxuAw6FvKtuB3/4ZK+EGJ2LsY4gJnmGNHvb4BguDE771f3Fye1ZVyS26UuURJI/Kr3JPUVkGwI5YSKjELlg7vupyv9l0+JWONDn0kneY/IqBkztBAz4po6EQQbIFl7+BYIj72OeEY3YDv7zEoMwJ0AzNCTiBz8ds0+YNNp1IWXPYwfKHvXINNw15uGq6GxxdglQT+yrHW3VRBsCmzbfS0rcjhqv4f10El/bqgTmEaLR9+25Z1JwfLnke+l50eyGXIHUPdv5zcBuAO56xpV+GvVI5VG4GCHK0dZqR0mTEnHUz6Ol7IZcgeQyUPz1DCQrDEtOg0GZfNKghxjsvtfUYmJT5hjJlDL6IuJqpvdrN4lKx5HyCvpIOzd+Cs8S0ZidU6hncB2xjccIBTZ+yncCxLikmy5HReMNwzXAnVwmCjxJBJfAIeqIcmjfcxe62VPqKP7DBDjZNWpg02BUbQHZCx1qcKDU6dITFpqICkPsmElksvgKXoJISQ2QWJAh0O9fqgDsvcMoeniJ0zeqO3/gd/3hHGNjq5q0yln+8u0qUd2ABeSxUoEjvIShxY/K34Kr+CqiMPTYo+2rL4ZgtTMo+UU32xZmXtGNMbLDFngX8upp5PzafFTOh9xC1B8rfRGhzzLZ2IrPDIQAOp7zhFIc9LCpsUAYcXonjXo68hkMnYL47vJyB2ceuZE+X3G1X30uT3W9nYCcCQoZclKng7W/bw4nYbWphj9XIZd9MtwPqMU98LY2FbF5/8Tv8u8TbjzPQB91PyLsSe/zlAeowNX8Vd18hJ4UqwgNlJ/Yg92nS+4ueGs/ngWLD/V/ed3mvIcvQgipBN5DhfY5IKErugSYWzxz5QdZfxQeK6+q8B8oTcusXHvClVDSXoDwht/4jVYLmBOUJ2YpbmCw9zRxZQaPiZDN7bsXlCbn1H80laMah9kHzaL55Qm79Rw+Uxbq5mZ0nNH3QHDONOUEzjjlBM445QTOOOUEzjgeKYTdmlsI5UoUydf1AMeweJa1rjgwgt35/ruJmHKoEzVVcnpBnRnsqQWMnM54jSzxQl4jyvqV7DWWPlWrFzQnKE8oOkbmKmy1ozOy5ipslKDtEVBX3JGldc2QAZYfIAz5JRV8mrWsETP+TEeOGaCaO8WEmPC4Jkl4z+f0ZKcUqkcVDo+IKE4+ELH73HVlMPvjeN2U+BBxcTKqBh2mEnBWrkt/2ELWYWL3DK1bEfXyIrMdFhX8mytxg8E74I40+rpHguFEoKB1M9yHQk2yHYqb+DRYNWhHeOI1ANJ0g3FBkGyO8f4Zv+6j3fhsjwiPVlGIydKhyKQU6U+/C4Qd2I/xfx4HsZ/LR9VH9JBG0NCVBpgjE1AqU2VPEpN7miGNsmvxBeHlPBxNLCXbuwdJIQi30bUfrKnnKNpWpnShlgkQgC3k/8HSGdgWDHQwOORCNUGaXvVBwmXDj/B5bX406xyIk0iGPNQFhk+FIRHB0HY9NhYh9//03xb1WMGmGpMIMxQCrS9Q1/DonfRnVpHQ8RHXYSpSlXUKZExFibwvYweTQD7RiscupT5AXim+DXzUqoRP6ZtCgwisRWaQO1LFC8RMNKrHpN6r+W73lqzWHNmGjw6BMXbqbMieccuorS30UhR1N2fbEBo3sLXylJ2gSO24YkNLyb/pE+0huzF2DBvsjVEKVQ2CfI9xQ5jvUadPgHUcBOTrkVNtMbg0NbD+CQvxxLjoc8C+8SAvhmFqbdJQ+UvcEyYJn6CDbcJ89gmT//NLYvVCZfwFuWDOvQbyHKQcsoaqIwWMAW3yjqelXQciuCHxpizqq7ABtNqlTp8oBbWwcXrEbUE9hbAfoSYYOHQwRDajuvwZlWkIf1Cds+iRQ41V0vUgjsn/++GbCht8QQ+f5jqZxvIDIwf/DcABXTei7Y5sqdRrURO1H1NjTNpwZkOYhhn2Re7V96fcqNQz2MajjhkDfpU4nUvrSgzbi2EPxsRR76Gh4+j1Iik7F1cXfFrYf5brGvxTrzI1zYEZofjf2tQ24PcqulqCaNoT5FkMjxQSFIBPbl9s6VWwsvxdUoRv9bIzdei5kw6ILHkHTmwmbuBFvvOZ02MHNuNDkDQgFOGwE9912IxtUfVLlkUQjEPE6KG+OyOjg4oAWlmIomEIpqjj26axrEuDU/U9NvqGCKYwKFQbpzm/IPVBAguSHKw7jNo8BM9C8bmAVk0Ml/Cu473Ybkxab2P5DWrTYjOkzhiZu25c+cGP01BQyDqdIAWBQYRsLhwNsapxg0+YHP6K9i2U6ofvw0JxIhkqKkf0JPIIGymyCOUG0ERMvZF5wtHAoRgtDlWZS5QBbWFhlvPewpu21hhg2hhMqdzhWXjGD8kQjJjeojAW0eSkIaWJRpUIV92XY9K/haOuoTnBdXVDFHgxjlspJUpYnIGgDgxM2cUTTuLlGDJFiafgwLWyavgqycMOSV7BGDO/akb/sISsbJ6aub/wmVK3IXcq0aVLFUnJQHtNhIxATy/IjA6UB+X5FtBePILkX+iqYgSMRTMrsUKMmxjcGLRFNd4vNAD0NTKHYmjSBXTqY1HjjzzxkjeDoRr7iSxzc5FA70m8tbI4CxoxFunkk/036LrTQUILCKFAaM2LPLnCKTYsfaGNwiMFLWhywTYMdHLx41jsiEukObmasb6lQo5PpGCOIg5grOf6n1ojIx9ukSdBiVMw+j6CBkmbIHIsgkwpNHKG5TU4whdzYvKTFOzaxcTgOZFIwaWBxQIcOGzQxNX3JKAxH7QaegbuMO4/mQh6qjkqtlvR5LWmebjrIiRn8uOUPA0WL0ilvx7jAIV5eq1MsTnDYpCMaqcMmJ7T41g+aN7SSvMHoJkfUqdLmmDbe7EPZz19iEp4d95LRVEOGctXvXYZxttsSQUehb+7xG8gDAUZYgLuYqfZA8sDGl94hQRdK0tXFMWTo18Db32ZfSuLXYTNkhJdp0OaAo0DZnqDIVRwObwgGqVTDVXbEWfWpGmZbmAKqNMhTuVboW4VmijFOFxUT22vLUB7V/5AOOx9LhvSImguIXtyaftlrnDsxRIYIRznaHnF+mvdpSfF+P/Lf3scgQXLS1QH/OfdYvRFsSyZCILlwcPO8nNK2oEl7O0f6WFEsuAATQYJ6Sp+zzBzZQ02OGphmC7ufyEmdliZaXZ1jHIxIjhom6GJUZuk5UseflJKYZOsD/iEdvDr3F8oUC0o3ImVPlz3s/i59L0iTqHOkiyXFQPgt/FUm6FoxFNbmMpQh1qXviqGm+qjKhkJh3g9lhhVlBuFMPkQlSDUU1pkjG6jyI49FNQQNwmYeUJQmIrJAWt4BepiB65gxR2V5DypU+blQD3qoOfGcr6Wua113aqo4YlksJVeVJuxIaWpaYi3piMdiJ0KFQ6zQ+o1Fy9/Ca/LB3+tQDpS7iwbDaV03BbuHqpLWM22omknTyro4CToZujk1t802y4F/u2yFfjexxNy2EVgLMse6hgeLRmiXXBCNRA4qk0MnP5qNOg+1J6sytMb5jU2chhOpm9KvtbF266TtaZEWFjSv/JnuQD1BA94rq0Nr+gpSgAk8FlJg46qnIcqhif/o3W46pO1pkR5Wk8lPFEF6Gfotk4jAQ5eOD8AXwGVow1TYw62FIzI/uoaFe66Btz/gKVl5WqQJVX76Ua9/FEE6GbL4KYObdfeZ7WL6m3DDidSdwOeaWEhz9xlY/p5PEyuQQysLT4t0oc6/vY/aKPowshJVhpbGWgRPCreRq4ENtnIi9ab4W/aXt10rrsmyIKCKFVBcaXtapI0VZdjSi7aSowka8AvPpDKLv2ZuKmxGlJucYAc2CkcjbU+LdKE3DyI3Wj+MqeocU1qrKLKewj6FKBiYmFr/bNvf+1NPQFDanhbpYk2TEChmlPkwtrK/SUnvYA07AzVnsIzFBxyO/Y3nZYbd9GPA5tuEb3Q2nhbpYFGzOvA67oR4gq45VypMW8252xfB5hVNHAzKtBnOFwwxbDB3tuEb4WTl7rF2vey8jVxZe1pMhoKy3xsu4p3l4gmCM2VLQ9pqzgH22cAQqqVKjT+MOMfy5czbtmiL+YaGT1C2nhaT4blGvY0YXY4iaMBbhfU1uinOzblmdtmfqNzwFZu3T7StGA7VkbVm7WkxCVY0k85no/ywRscsvdD0Oc8zDFpW9j+9YZNNNhUnRQ81f8ahQis0E+16WhjUgCZNTtkSpGxxStOffG1gCi+GoaeFQ004bqWLBcUmhovRL3qSoLJtTejm7zNaZ7UwMYRH6AfatGlHjOjL1P0NwR1/D7cL19Nihz0s8D0t4ACDhqDSoEVV8bTY4RV7PM7A0+KFsrg9Ur1BMoJ6Gm1cymh+u0GbfeEGHweTk4DDis1LTH8IKntatKiIVOo2L7F4h4nrl7fpGw0mLfY4oEOdNk1qKa8jP1N6nwTqDUb3QS4uNdbcGp9TSspewRB7W2rCueuAQwh5DDWlKdMW4Yz3HQ445JIj0ve0mB7rGuM6gXqDpATBGabyDjzjYypjojI1HPYxhEUFdZpU2WDL71nqoTMsYFPqJY786Zq0PS2mxaJG2/SSWsJfJJ5fL/FnTTKo/0p1hnv04NAITIWOe3Y+nhYLvNCot78k9aJPThCsaeyQHj/O0xLGQE/P2+SdwzipAc411RYzs+fuAgpaes7H6bvHy91wppmWKPFiTpEWenoS9z4uxiNowGuN7ixpVN8c8EwT86g37qLnuNlPerzWJPJcnTiA0N3Fhsb9rc9P44bYGT89TVcronOKwtjQeie+GT8CUtJxUBAXPNIotVXg7dynFSjwTEvP2STze+OY2UGsa6d6uvx47ykq8EIbb+9ssm1rkxIURdF9Hxct8H2a9ExD0JwiFfph6RT0TEdQNEXtDPYtzD4W+U6bO2YKeqYlKIqiKW/qVuJrnmfREtMSFE1RGoFkbgsK2gUFSOFFnZ6gaIruS28UZRikokfSICiaogE/Z+76lTdWeBaRtWyMOetopEMQLGNF3OY5Z3d2bBSt2vq8TefVTIsgKPJnrYl5d626RazIJ36dVg6m9AiKo+juyVG07EBXO+c/IdIkCOBZ5G33OLsz/VG07MA5v2jm+ydG2gTpF8Y9XPC3W2/XLWCxFPFbn19S2unkI32C4lUdnN2gO3LaKLCmuLUNkWLPM0QWBEGB5zHBL26rsvua9Zg0wCmrNg/ZEASwwnqMHN02kuKfJjWjWkV2BLmOKnFBZG4LSfHkwNUkK6VJkSVBSR6uxxn2zPZJ8X0OZGIWhJE1QaPlyI2RNnvWXQlzBDlwwdss+p0gsifIfdjvYuUI4CrpdvLMUWCRtUhT2kOXn29ifuRmCILRyg5cWfot12mhRZ7yxxFycwOKbYibI8iNO5fEqygfmpJRA33e8z5rxTbETRIESXqkYUNcY3OVed9UZJEllhNQc+PkwM0T5DZJUpIAunS5opv6GN0lZnGk2vWQAzmQD0HgkpS8cQD6dPnEFX26E5vlBYqU+Ionmqy/8dfOhRzIjyCARywnMBx0cEnq0uczfQYM6CukFXgELFCkwJcUKFGY6Fo9LvIiB/IlyMUiqzcQtHYyXHGW91Jj/gSB2x+sTZCgOjtccZ2n3AwxGwS5KLHK8kRqKE30uOAqb7kZYpYIclFikdVcpOmK61mixsXsEeRiXCN4GvS45orLWVBoKmaVIA8uUaVMJKrLNR+5zm6pIA3MOkFDLFLiS0o8GWsEI6NHl898pMvn2ZQYGbeHoCBKFCjxiCJFoEBBjHqCcEdG7jipR48+nxjcDlKC+H+2n9a+PKxA2AAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMS0xMC0wOVQxMDozNjo0OSswODowMH49o1wAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjEtMTAtMDlUMTA6MzY6NDkrMDg6MDAPYBvgAAAAIHRFWHRzb2Z0d2FyZQBodHRwczovL2ltYWdlbWFnaWNrLm9yZ7zPHZ0AAAAYdEVYdFRodW1iOjpEb2N1bWVudDo6UGFnZXMAMaf/uy8AAAAYdEVYdFRodW1iOjpJbWFnZTo6SGVpZ2h0ADIwOQQLrc0AAAAXdEVYdFRodW1iOjpJbWFnZTo6V2lkdGgAMjA44P3NBgAAABl0RVh0VGh1bWI6Ok1pbWV0eXBlAGltYWdlL3BuZz+yVk4AAAAXdEVYdFRodW1iOjpNVGltZQAxNjMzNzQ3MDA5seZ0eAAAABJ0RVh0VGh1bWI6OlNpemUAODUxM0JChD/OyQAAAEZ0RVh0VGh1bWI6OlVSSQBmaWxlOi8vL2FwcC90bXAvaW1hZ2VsYy9pbWd2aWV3Ml85XzE2MzE3NTc1MjkyNDAyNzIwXzEzX1swXYuxiNAAAAAASUVORK5CYII="
            let proxySymbol8 = "image://data:image/png;base64,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eQJ77tNHSUnP5Th993h7N5wrhiUqPC0oUkJNT38863c8goa8UAhnReGCFYDHipXy/XEnPcfd/aTPC4X18SAw3lsA4Ili+duAn8ZNsTP+9jQ9pYgWFAXxRLk68eX4GZCSOapBXHFPodQeAK+KXViBEo+V9JynmVpM7gcFsa6M2O7x42dPUYmnyiwt5+nC1tISFEZRnx8/a9d1gR+ypGcSggqKZKjd0gnomYygcIo6OeX5mW8s8r1y75gJ6JmUoDCKJnypjxLf8l0eNTEpQeEUfU6JZErKCQXIoKFOTlA4RZ9LbxRmGGSiR7IgKJyiIT9nuPRrPrHK45Bdy15lkfc1G4JgGSPkNS85/2R9o3DVNuBVNk0zK4KgzF+UJuana9UtYoR+8YuscvdnR1AURZ+eHIXLDvSUY/4pkSVBoI5esdHn/JPpj8JlBy75NVm0QTJkTRCsRcwOXfHLR2/XLWCwFHJuwK9ZJ4TPnqBoVQfnmS1Hnj5KrEnL2kbIsOcZIQ+CoMR3EfHbH6uy+5b1iG2AM1ZtLvIhCGCV9Qg5+thIiv6azIxqGfkRZC9UiVoH8bGQFE0O3KSZKU2KPAlK8nF9zjHntk+K7nMgF7MgiLwJipcjO0fa/Fl3FfQYcuCKV3n0O37kT5D9sd9HyhHATdJw8txRYpG1UFPaRY+fpzE+Mh2CIF7ZgS1Lv890WGiRFf4UIzdTUGwjTI8gO+9ckuSAs6EpGTUw4DWv81ZsI0yTIEjSI40q4haTm9z7pjKLLLGcgJqpkwPTJ8iukqQkAfTocUMvcx/dJmYxVu26mAE5MBuCwCYpeeUADOjxjhsG9FKb5SXKVPiG+4pdf6OfPRNyYHYEAdxjOYHhoIJNUo8BHxgwZMhAIq3EPWCBMiW+pESFUqpn9bmaFTkwW4JsLPJgRklr43HD+aynGmdPENj9wVruG1SPgxtuZyk3I8wHQTYqPGA5lRrKEn2uuJm13IwwTwTZqLDIg5lI0w2380SNjfkjyMa4RvAk6HPLDdfzoNBkzCtBLmyiKrlIVI9b3nKb47byGWDeCRphkQpfUuH+WB6MiD49PvCWHh/mU2JEfDwE+VGhRIV7lCkDJUqO1+OH7RnZflKfPgPeMfw4SPHj/wHHDFkk2cDSuwAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMS0xMC0wOVQxMDozNzozMiswODowMJk9lYEAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjEtMTAtMDlUMTA6Mzc6MzIrMDg6MDDoYC09AAAAIHRFWHRzb2Z0d2FyZQBodHRwczovL2ltYWdlbWFnaWNrLm9yZ7zPHZ0AAAAYdEVYdFRodW1iOjpEb2N1bWVudDo6UGFnZXMAMaf/uy8AAAAYdEVYdFRodW1iOjpJbWFnZTo6SGVpZ2h0ADIwOQQLrc0AAAAXdEVYdFRodW1iOjpJbWFnZTo6V2lkdGgAMjA44P3NBgAAABl0RVh0VGh1bWI6Ok1pbWV0eXBlAGltYWdlL3BuZz+yVk4AAAAXdEVYdFRodW1iOjpNVGltZQAxNjMzNzQ3MDUyW0NZtQAAABJ0RVh0VGh1bWI6OlNpemUAOTI4MUJCg89amAAAAEZ0RVh0VGh1bWI6OlVSSQBmaWxlOi8vL2FwcC90bXAvaW1hZ2VsYy9pbWd2aWV3Ml85XzE2MzM3NDI3ODg2ODY0MTcyXzY5X1swXROe2j4AAAAASUVORK5CYII="

            // 节点数据
            let nodeDataList = [
            {
                    name: "第一步",
                    linkTargetName: '第二步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '0',
                    },
                    value: [10, 390.01],
                    draggable: false,
                    fixed: true,
                    symbol: appSymbolone,
                    symbolSize: 80
                }, {
                    name: "第五步",
                    linkTargetName: '第六步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '0',
                    },
                    value: [10, 140],
                    draggable: false,
                    symbol: proxySymbol5,
                    symbolSize: 80
                },
                //节点B,C,D ....n
                {
                    name: "第二步",
                    linkTargetName: '第三步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '1',
                    },
                    symbol: masterSymbol2,
                    symbolSize: 80,
                    draggable: false,
                    value: [100, 390]
                }, {
                    name: "第六步",
                    linkTargetName: '第七步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '1',
                    },
                    symbol: proxySymbol6,
                    symbolSize: 80,
                    draggable: false,
                    value: [100, 140]
                },    {
                    name: "第三步",
                    linkTargetName: '第四步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '1',
                    },
                    symbol: slaveSymbol3,
                    symbolSize: 80,
                    draggable: false,
                    value: [200, 390]
                }, {
                    name: "第七步",
                    linkTargetName: '第八步',
                    linkValue: ' ',
                    coordConfig: {
                        level: '1',
                    },
                    symbol: proxySymbol7,
                    symbolSize: 80,
                    draggable: false,
                    value: [200, 140]
                },
                {
                    name: "第四步",
                    linkTargetName: '第五步',
                    linkValue: ' ',
                    // coordConfig: {level: 2},
                    symbolSize: 80,
                    symbol: proxySymbol4,
                    draggable: false,
                    value: [320, 390]
                }, {
                    name: "第八步",
                    linkTargetName: '第八步',
                    linkValue: ' ',
                    // coordConfig: {level: 2},
                    symbolSize: 80,
                    symbol: proxySymbol8,
                    draggable: false,
                    value: [320, 140]
                }
            ]

            let getCoordDataList = function() {
                let coorDataDict = {}
                let defaultConfig = {
                    type: 'lines', //块1,2...n到节点A,B...N
                    coordinateSystem: 'cartesian2d',
                    // animationDelay: 1000,
                    z: -1,
                    effect: {
                        show: true,
                        smooth: true,
                        trailLength: 0,
                        symbol: "arrow",
                        color: '#6FCF97',
                        symbolSize: 10,
                        period: 1,
                        delay: 3500,
                        loop: true,
                    },
                    lineStyle: {
                        normal: {
                            curveness: 0,
                            color: '#6FCF97',
                            width: 4,
                        }
                    },
                    data: []
                }
                nodeDataList.map(item => {
                    if (item.coordConfig !== undefined) {
                        if (!(item.coordConfig.level in coorDataDict)) {
                            let coorConfig = JSON.parse(JSON.stringify(defaultConfig))
                            if (item.coordConfig.lineStyle !== undefined) {
                                coorConfig.lineStyle = item.coordConfig.lineStyle
                            }
                            if (item.coordConfig.effect !== undefined) {
                                coorConfig.effect = item.coordConfig.effect
                            }
                            coorDataDict[item.coordConfig.level] = coorConfig
                        }
                        
                        let coordData = [
                            item.value,
                            nodeDataList.filter(i => i.name === item.linkTargetName)[0].value
                        ]
                        coorDataDict[item.coordConfig.level].data.push({
                            coords: coordData
                        })
                        if (item.coordConfig.bilateral) {
                            coorDataDict[item.coordConfig.level].data.push({
                                coords: coordData.reverse()
                            })
                        }
                    }

                })
                return Object.values(coorDataDict)
            }

            // console.log(getCoordDataList())

            optionsectionImg = {
                // backgroundColor: '#ff6700',
                // title: {
                //     text: "MySQL应用架构拓扑图",
                //     top: "top",
                //     left: "center"
                // },
                itemStyle: {
                    normal: {
                        color: "#67C23A",
                    },
                    shadowBlur: 0
                },
                textStyle: {
                    color: '#fff',
                    // fontSize: 16,
                    fontWeight: 600,
                },
                legend: [{
                    formatter: function(name) {
                        return echarts.format.truncateText(name, 200, '12px', '…');
                    },
                    tooltip: {
                        show: true
                    },
                    selectedMode: 'false',
                    bottom: 20,
                }],
                animationDuration: 500,
                animationEasingUpdate: 'quinticInOut',
                xAxis: {
                    show: false,
                    type: 'value'
                },
                yAxis: {
                    show: false,
                    type: 'value'
                },
                series: [{
                    type: 'graph',
                    coordinateSystem: 'cartesian2d',
                    legendHoverLink: false,
                    hoverAnimation: false,
                    nodeScaleRatio: false,
                    //建头
                    edgeSymbol: ['circle', 'none'],
                    edgeSymbolSize: [2, 15],
                    edgeLabel: {
                        show: false,
                        normal: {
                            show: true,
                            position: 'middle',
                            textStyle: {
                                fontSize: 18
                            },
                            formatter: "{c}"
                        }
                    },
                    focusNodeAdjacency: true,
                    roam: false,
                    //圆形上面的文字	
                    label: {
                        normal: {
                            position: "bottom",
                            show: true,
                            textStyle: {
                                fontSize: 12
                            },
                        }
                    },
                    // itemStyle: {
                    //     normal: {
                    //         color: "#409eff",
                    //     },
                    //     shadowBlur: 0
                    // },
                    // lineStyle: {
                    //     normal: {
                    //         width: 0,
                    //         shadowColor: 'none',
                    //         color: '#ff0000'
                    //     }
                    // },
                    data: nodeDataList,

                    links: nodeDataList.map(item => {
                        return {
                            source: item.name,
                            value: item.linkValue,
                            target: item.targetName
                        }
                    })
                }].concat(getCoordDataList())
            }

            optionsectionImg && myChartsectionImg.setOption(optionsectionImg);
        },

        /* 
        淘煤网阳光竞价图
        */
        paymentTaxes(){
            document.getElementById("paymentTaxes").removeAttribute("_echarts_instance_");
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
                    x:80,
                    y:30,
                    x2:40,
                    y2:50,
                },
                xAxis: [
                    {
                        type: 'category',
                        data:['2018','2019','2020'],
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
                        name: "竞价场次",
                        type: "bar",
                        barWidth: 16,
                        itemStyle: {
                            normal: {
                                color: "#CD5D73"
                            }
                        },
                        data: [508,1451,928],
                        z: 10,
                        zlevel: 0,
                        "label": {
                            show: true,
                            position: "top",
                            distance: 8,
                            color: "#fff"
                        },
                        barGap:1.5,
                    },
                    {
                        name: "有效成交吨数",
                        type: "bar",
                        barWidth: 16,
                        itemStyle: {
                            normal: {
                                color: "#1570BB"
                            }
                        },
                        data: [220.16,447,441.48],
                        z: 7,
                        zlevel: 0,
                        "label": {
                            show: true,
                            position: "top",
                            distance: 8,
                            color: "#fff"
                        },
                        barGap:1.5,
                    },
                    {
                        name: "实现增收",
                        type: "bar",
                        barWidth: 16,
                        itemStyle: {
                            normal: {
                                color: "#16CDC2"
                            }
                        },
                        data: [2282.89,3627,6119.82],
                        z: 10,
                        zlevel: 0,
                        "label": {
                            show: true,
                            position: "top",
                            distance: 8,
                            color: "#fff"
                        },
                        barGap:1.5,
                    }
                ]
            };
            optionpaymentss && myChartpaymentss.setOption(optionpaymentss);
        },
    },
    // 销毁定时器
    beforeDestroy() {
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
    background-image: url(../../assets/backimg.jpg);
    background-size: cover;
    background-repeat: no-repeat;
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
            width: 59.9%;
            height: 32.5%;
            margin: 0.5% 0.1% 0.2% 0;
            background-image: url(../../assets/k1.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 0 1% 1%;
            box-sizing: border-box;
            .section-title{
                width: 100%;
                font-size: 18px;
                color: #fff;
                padding:2% 0 1% 4%;
                box-sizing: border-box;
                img{
                    height: 100%;
                    vertical-align: middle;
                    margin-right:1%;
                }
            }
            .section-img{
                width: 95%;
                height: 78%;
                // padding:0.5% 0 1.5% 8%;
                padding-left: 5%;
                padding-bottom:1%;
                box-sizing: border-box;
                // background: chocolate;
                margin: auto;
                // #sectionImg{
                //     width: 100%;
                //     height: 100%;
                // }
                img{
                    height: 100%;
                }
            }
        }
        .section-left{
            width: 39.9%;
            height: 32.5%;
            margin: 0.5% 0.1% 0.2% 0;
            background-image: url(../../assets/k2.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            // background: bisque;
            .section-title{
                width: 100%;
                font-size: 18px;
                color: #fff;
                padding:3% 0 0 8%;
                box-sizing: border-box;
                img{
                    height: 100%;
                    vertical-align: middle;
                    margin-right:1%;
                }
            }
            #paymentTaxes{
                width:100%;
                height: 83%;
                padding-top:2%;
                box-sizing: border-box;
            }
            // background: cyan;
        }
        // 中部
        .conent{
            width: 99.9%;
            height: 36%;
            margin:0 0.1% 0.2% 0;
            background-image: url(../../assets/k3.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 1.5%;
            box-sizing: border-box;
            .conent-title{
                width: 100%;
                font-size: 18px;
                color: #fff;
                line-height: 28px;
                padding-left:2%;
            }
            .conent-left{
                width: 25%;
                height: calc(100% - 28px);
                float: left;
                padding: 1.5%;
                box-sizing: border-box;
                img{
                    width: 100%;
                    height: 100%;
                    margin: auto;
                }
            }
            .conent-right{
                width: 75%;
                height: calc(100% - 28px);
                float: left;
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
                            font-size: 16px;
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
            // #burden{
            //     width: 100%;
            //     height: 100%;
            // }
            
        }
        // 下部
        .footer{
            width: 69.9%;
            height:30%;
            margin:0.2% 0.1% 0.2% 0;
            background-image: url(../../assets/k4.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 2% 0 0 2%;
            box-sizing: border-box;
            .section-title{
                width: 100%;
                font-size: 18px;
                color: #fff;
                line-height: 28px;
                padding:0 0 1% 1%;
                box-sizing: border-box;
            }
            #sectionBoxtwo{
                width:97%;
                height: calc(100% - 58px);
                overflow: hidden;
                .animtwo{
                    transition: all 1s;
                }
                #con3{
                    width: 100%;
                    display: flex;
                    flex-wrap: nowrap;
                    li{
                        width:20%;
                        margin: 0.5%;
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
        .footer-right{
            width: 29.9%;
            height: 30%;
            margin:0.2% 0.1% 0.2% 0;
            background-image: url(../../assets/k5.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            padding: 2% 2% 0 2%;
            box-sizing: border-box;
            overflow: hidden;
            .section-title{
                width: 100%;
                font-size: 18px;
                color: #fff;
                padding:0 0 1% 4%;
                box-sizing: border-box;
                img{
                    height: 100%;
                    vertical-align: middle;
                    margin-right:1%;
                }
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
                        margin: 2%;
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
</style>
