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
            <p>煤炭交易大市场</p> 
            <div class="TimesAll">
                <span>{{FirstScreenTime}}</span>
                <div class="buttonbox" @click="BackOff">
                    <i class="iconfont el-icon-alireturn"></i>
                </div>
            </div>
            <!-- <div class="TimesAllleft" @click="BackTime">
                <span v-if="backcolor">自动循环 </span>
                <span v-if="!backcolor" style="color:rgb(1,240,244);">自动循环 </span>
                <i v-if="backcolor" class="el-icon-turn-off"></i>
                <i v-if="!backcolor" class="el-icon-open" style="color:rgb(1,240,244);"></i>
            </div> -->
        </div>
        <div class="middle">
            <div class="echartsbutton" @click="ehartClick">云 仓 报 表</div>
            <div class="middleBox">
                <el-select @change="citys($event)" v-model="threeName" placeholder="三级小类名称" style="width:20%;margin:5px 0 0 12px;" :popper-append-to-body="false" popper-class="select-popper">
                    <el-option
                    v-for="item in options"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
                <el-select @change="entisetwos($event)" v-model="secName" placeholder="二级小类名称" style="width:20%;margin:5px 0 0 12px;" :popper-append-to-body="false" popper-class="select-popper">
                    <el-option
                    v-for="item in cityOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
                <!-- clearable @clear="clearent" -->
                <el-select @change="enterprises($event)" filterable v-model="fstName" placeholder="一级大类名称" style="width:20%;margin-top:5px;" :popper-append-to-body="false" popper-class="select-popper">
                    <el-option
                    v-for="item in provinceOptions"
                    :key="item.id"
                    :label="item.productName"
                    :value="item.id + ',' + item.productName">
                    </el-option>
                </el-select>
                <div class="elbutton" @click="elbuttonClick">重 置</div>
            </div>
        </div>
        <div class="nav">
           <div class="navleft">
               <div>
                    <ul id="con2" ref="con2" :class="{animone:animateone==true}">
                        <li v-for='(v,k) in itemsone' :key="k" @click="Checkli(v,k)">
                            <!-- 选中左侧列表 -->
                            <div class="navli" v-if="k==0" @mouseover="mouseovnavs" @mouseleave="mouseleavenavs">
                                <!-- 左侧视频 -->
                                <div class="flexnasleft">
                                    <img v-if="v.warehouseImg" :src="http + '/saas-common/upload/getVideo?fileName='+ v.warehouseImg + '&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f'" alt="">
                                    <!-- <div class="dlexicon" v-if="v.warehouseVideo" @click="videoClick(v)"><i class="el-icon-alishipinbofang"></i></div> -->
                                    <img v-else src="@/assets/defaultPic.png" alt="">
                                    <div class="flexnastitle">{{v.warehouseNb}}</div>
                                    <div class="flexbottomsize">{{v.productName}}</div>
                                </div>
                                <!-- 右侧内容 -->
                                <div class="flexnasright">
                                    <div class="flexnasrightBox">
                                        <!-- <span v-if="v.isShowAshContent">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                        <span v-if="v.isShowWaterContent">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                        <span v-if="v.isShowStd">硫(St,d)%≤{{v.maxStd}}，</span>
                                        <span v-if="v.isShowVolatiles">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                        <span v-if="v.isShowBondIndex">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                        <span v-if="v.isShowY">Ymm≥{{v.minY}}，</span>
                                        <span v-if="v.isShowCsr">热态CSR≥{{v.minCsr}}，</span>
                                        <span v-if="v.isShowHgi">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                        <span v-if="v.isShowLithofacies">岩相S≤{{v.lithofacies}}，</span>
                                        <span v-if="v.isShowRecoveryAnalysis">回收率分析 {{v.recoveryAnalysis}}，</span>
                                        <span v-if="v.isShowCalorificValue">热值≥{{v.calorificValue}}</span> -->
                                        <span v-if="v.maxAshContent!=null&&v.maxAshContent!=0">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                        <span v-if="v.maxWaterContent!=null&&v.maxWaterContent!=0">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                        <span v-if="v.maxStd!=null&&v.maxStd!=0">硫(St,d)%≤{{v.maxStd}}，</span>
                                        <span v-if="v.maxVolatiles!=null&&v.maxVolatiles!=0">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                        <span v-if="v.minBondIndex!=null&&v.minBondIndex!=0">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                        <span v-if="v.minY!=null&&v.minY!=0">Ymm≥{{v.minY}}，</span>
                                        <span v-if="v.minCsr!=null&&v.minCsr!=0">热态CSR≥{{v.minCsr}}，</span>
                                        <span v-if="v.minHgi!=null&&v.minHgi!=0">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                        <span v-if="v.lithofacies!=null&&v.lithofacies!=0">岩相S≤{{v.lithofacies}}，</span>
                                        <span v-if="v.recoveryAnalysis!=null&&v.recoveryAnalysis!=0">回收率分析 {{v.recoveryAnalysis}}，</span>
                                        <span v-if="v.calorificValue!=null&&v.calorificValue!=0">热值≥{{v.calorificValue}}</span>
                                    </div>
                                </div>
                            </div>
                            <!-- 未选中左侧列表 -->
                            <div class="navbloneli" v-else>
                                <!-- 左侧视频 -->
                                <div class="flexnaslefts">
                                    <img v-if="v.warehouseImg" :src="http + '/saas-common/upload/getVideo?fileName='+ v.warehouseImg + '&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f'" alt="">
                                    <img v-else src="@/assets/defaultPic.png"  alt="">
                                    <div class="flexnastitles">{{v.warehouseNb}}</div>
                                    <div class="flexbottomsizes">{{v.productName}}</div>
                                </div>
                                <!-- 右侧内容 -->
                                <div class="flexnasrights">
                                    <div class="flexnasrightBoxs">
                                        <!-- <span v-if="v.isShowAshContent">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                        <span v-if="v.isShowWaterContent">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                        <span v-if="v.isShowStd">硫(St,d)%≤{{v.maxStd}}，</span>
                                        <span v-if="v.isShowVolatiles">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                        <span v-if="v.isShowBondIndex">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                        <span v-if="v.isShowY">Ymm≥{{v.minY}}，</span>
                                        <span v-if="v.isShowCsr">热态CSR≥{{v.minCsr}}，</span>
                                        <span v-if="v.isShowHgi">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                        <span v-if="v.isShowLithofacies">岩相S≤{{v.lithofacies}}，</span>
                                        <span v-if="v.isShowRecoveryAnalysis">回收率分析 {{v.recoveryAnalysis}}，</span>
                                        <span v-if="v.isShowCalorificValue">热值≥{{v.calorificValue}}</span> -->
                                        <span v-if="v.maxAshContent!=null&&v.maxAshContent!=0">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                        <span v-if="v.maxWaterContent!=null&&v.maxWaterContent!=0">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                        <span v-if="v.maxStd!=null&&v.maxStd!=0">硫(St,d)%≤{{v.maxStd}}，</span>
                                        <span v-if="v.maxVolatiles!=null&&v.maxVolatiles!=0">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                        <span v-if="v.minBondIndex!=null&&v.minBondIndex!=0">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                        <span v-if="v.minY!=null&&v.minY!=0">Ymm≥{{v.minY}}，</span>
                                        <span v-if="v.minCsr!=null&&v.minCsr!=0">热态CSR≥{{v.minCsr}}，</span>
                                        <span v-if="v.minHgi!=null&&v.minHgi!=0">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                        <span v-if="v.lithofacies!=null&&v.lithofacies!=0">岩相S≤{{v.lithofacies}}，</span>
                                        <span v-if="v.recoveryAnalysis!=null&&v.recoveryAnalysis!=0">回收率分析 {{v.recoveryAnalysis}}，</span>
                                        <span v-if="v.calorificValue!=null&&v.calorificValue!=0">热值≥{{v.calorificValue}}</span>
                                    </div>
                                </div>
                            </div>
                        </li>
                    </ul>
               </div>
           </div>
           <div class="navright">
               <div class="navtitle">{{chname}}</div>
                <div id="sectionBox">
                    <ul id="con1" ref="con1" :class="{anim:animate==true}">
                        <li v-for='(v,key) in items' :key="key">
                            <!-- style="font-weight: 600;" -->
                            <!-- <div style="width:21%;height:48px;"></div> -->
                            <div class="sectionone">{{key + 1}}</div>
                            <div class="secttitle">{{v.productName}} <i v-if="v.ckhyFileName" @click="videoClick(v)" class="el-icon-alishipinbofang"></i>
                                <!-- v-if="v.lastTestDate" -->
                                <div class="sectbutton" v-if="v.lastTestDate" @click="TestEport(v)" style="margin-right:0;">检测报告</div> 
                                <!-- v-if="v.lastTransactionTime" -->
                                <div class="sectbutton" v-if="v.lastTransactionTime" @click="Testhistory(v)" >历史成交</div>
                            </div>
                            <div class="section">更新时间:<span>{{v.updateTime|Time}}</span></div>
                            <div class="section">出厂现汇<span>{{v.exFactoryPrice}}</span>元/吨</div>
                            <div class="section" style="margin-top:0;">现有库存 <span v-if="v.bleTonnage">{{v.bleTonnage}} 吨</span><span v-else>-- 吨</span></div>
                            <div class="section" style="margin-top:0;"> 平台可供应量 <span v-if="v.availableTonnage">{{v.availableTonnage}} 吨</span><span v-else>-- 吨</span></div>
                            <div class="sectconent">
                                <!-- <span v-if="v.isShowAshContent">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                <span v-if="v.isShowWaterContent">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                <span v-if="v.isShowStd">硫(St,d)%≤{{v.maxStd}}，</span>
                                <span v-if="v.isShowVolatiles">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                <span v-if="v.isShowBondIndex">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                <span v-if="v.isShowY">Ymm≥{{v.minY}}，</span>
                                <span v-if="v.isShowCsr">热态CSR≥{{v.minCsr}}，</span>
                                <span v-if="v.isShowHgi">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                <span v-if="v.isShowLithofacies">岩相S≤{{v.lithofacies}}，</span>
                                <span v-if="v.isShowRecoveryAnalysis">回收率分析 {{v.recoveryAnalysis}}，</span>
                                <span v-if="v.isShowCalorificValue">热值≥{{v.calorificValue}}</span> -->
                                <span v-if="v.maxAshContent!=null&&v.maxAshContent!=0">灰分(Ad)%≤{{v.maxAshContent}}，</span>
                                <span v-if="v.maxWaterContent!=null&&v.maxWaterContent!=0">水分(Mt)%≤{{v.maxWaterContent}}，</span>
                                <span v-if="v.maxStd!=null&&v.maxStd!=0">硫(St,d)%≤{{v.maxStd}}，</span>
                                <span v-if="v.maxVolatiles!=null&&v.maxVolatiles!=0">挥发分(Vdaf)%≤{{v.maxVolatiles}}，</span>
                                <span v-if="v.minBondIndex!=null&&v.minBondIndex!=0">粘结指数GR.I≥{{v.minBondIndex}}，</span>
                                <span v-if="v.minY!=null&&v.minY!=0">Ymm≥{{v.minY}}，</span>
                                <span v-if="v.minCsr!=null&&v.minCsr!=0">热态CSR≥{{v.minCsr}}，</span>
                                <span v-if="v.minHgi!=null&&v.minHgi!=0">哈氏可磨指数（HGI）≥{{v.minHgi}}，</span>
                                <span v-if="v.lithofacies!=null&&v.lithofacies!=0">岩相S≤{{v.lithofacies}}，</span>
                                <span v-if="v.recoveryAnalysis!=null&&v.recoveryAnalysis!=0">回收率分析 {{v.recoveryAnalysis}}，</span>
                                <span v-if="v.calorificValue!=null&&v.calorificValue!=0">热值≥{{v.calorificValue}}</span>
                            </div>
                        </li>
                    </ul>
                </div>
           </div>
        </div>
        <!-- 检测报告弹框 -->
        <div class="notification" v-if="notficShow" @mouseleave="mouselnots"> 
                <el-carousel trigger="click" height="685px">
                    <el-carousel-item v-for="(item,index) in imagelists" :key="index">
                        <img :src="http + '/saas-common/upload/getVideo?fileName='+ item.fileName + '&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f'" alt="">
                    </el-carousel-item>
                </el-carousel>
        </div>
        <!-- 历史成交弹窗 -->
        <div class="nothistory" v-if="historyShow" @mouseleave="mouseTesth">
            <el-timeline-item v-for="(v, index) in activities" :key="index" :timestamp="v.transactionTime" placement="top" color="#69FEF4">
            <div class="cards">
                <div class="cardflex">
                    (Ad)% {{v.fashContent}}，(Vdaf)% {{v.fvolatiles}}，(St,d)% {{v.fstd}}，GR.I {{v.fbondIndex}}，Ymm {{v.fy}}
                </div>
                <div class="cardone">{{v.tonnageTraded}}吨</div>
                <div class="cardtwo">{{v.transactionPrice}}元</div>
                <div class="cardtree">{{v.laboratory}}</div>
            </div>
            </el-timeline-item>
        </div>
        <!-- 视频弹窗 -->
        <div class="videoBox" v-if="videoShow" @mouseleave="videomous">
            <div class="videotitie">视频更新时间：{{videoatitie}}</div>
            <video :controls="true" preload="auto" :playsinline="true" width="950px" height="600px" :src="http + '/saas-common/upload/getVideo?fileName='+ videourl + '&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f'" x5-playsinline webkit-playsinline="true" playsinline="true" autoplay></video>
        </div>
        <!-- 云仓报表弹框 -->
        <div class="echartBox" v-if="drawer">
            <div class="offbutton" @click="offClick">退 出 全 屏</div>
            <div class="gdtBox">
                <el-table
                :data="tableData"
                :cell-style="cellStyle"
                :header-cell-style="{ background: '#0B2956', color: '#86ADE4',borderBottom:'1px solid #86ADE4', }"
                style="width:100%;">
                    <el-table-column prop="productName" label="品名" :show-overflow-tooltip="true" width="140">
                    </el-table-column>
                    <el-table-column label="出厂现汇价" :show-overflow-tooltip="true" width="140">
                    <template slot-scope="scope">
                        <span v-if="scope.row.exFactoryPrice">{{scope.row.exFactoryPrice}}元/吨</span>
                        <span v-if="scope.row.platePrice">{{scope.row.platePrice}}元/吨</span>
                        <span v-if="!scope.row.platePrice && !scope.row.exFactoryPrice">--</span>
                    </template>
                    </el-table-column>
                    <el-table-column label="更新时间" :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="scope.row.updateTime">{{scope.row.updateTime|Time}}</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="涨跌" :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="(scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1) > 0" style="color:rgb(255,122,122);"><img style="width:10px;" src="../../assets/sz.png" alt=""> {{scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1}} 元</span>
                            <span v-if="(scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1) < 0" style="color:#48fc07;"><img style="width:10px;" src="../../assets/xd1.png" alt=""> {{(scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1)*-1}} 元</span>
                            <span v-if="(scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1) == 0" style="color:#fff;">{{scope.row.exFactoryPrice*1 - scope.row.preExFactoryPrice*1}} 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="硫(St,d)%≤" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.maxStd">{{scope.row.maxStd}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="灰分(Ad)%≤" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.maxAshContent">{{scope.row.maxAshContent}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <!-- <el-table-column label="水分(Mt)%≤" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <span v-if="scope.row.maxWaterContent">{{scope.row.maxWaterContent}}</span>
                        <span v-else>--</span>
                    </template>
                    </el-table-column> -->
                    <el-table-column label="挥发分(Vdaf)%≤" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.maxVolatiles">{{scope.row.maxVolatiles}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="粘结指数GR.I≥" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.minBondIndex">{{scope.row.minBondIndex}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="胶质层Ymm≥" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.minY">{{scope.row.minY}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="强度CSR≥" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.minCsr">{{scope.row.minCsr}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="岩相标准差S≤" :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <span v-if="scope.row.lithofacies">{{scope.row.lithofacies}}</span>
                            <span v-else>--</span>
                        </template>
                    </el-table-column>
                    <!-- <el-table-column label="哈氏可磨指数(HGI)≥" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <span v-if="scope.row.minHgi">{{scope.row.minHgi}}</span>
                        <span v-else>--</span>
                    </template>
                    </el-table-column> -->
                    <!-- <el-table-column label="回收率%" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <span v-if="scope.row.recoveryAnalysis">{{scope.row.recoveryAnalysis}}</span>
                        <span v-else>--</span>
                    </template>
                    </el-table-column> -->
                    <!-- <el-table-column label="热值≥" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <span v-if="scope.row.calorificValue">{{scope.row.calorificValue}}</span>
                        <span v-else>--</span>
                    </template>
                    </el-table-column> -->
                <!-- 暂无数据 -->
                <!-- <div slot="empty" class="dataAvailable" style="background: #0B2956;">
                    <p>暂无数据，您可以新建试试!</p>
                </div> -->
                </el-table>
            </div>
        </div>
    </div>
</template>
<script>
export default{
    data(){
        return {
            http:"",
            backcolor:false,            // 自动循环
            FirstScreenTime:'',         // 时间屏幕
            routerboday:null,           // 跳转页面
            provinceOptions:[],             // 一级大类
            cityOptions: [],                // 二级小类
            options: [],                    // 三级小类
            fstName:'',
            secName:'',
            threeName:'',
            fstProId:'',
            bigProId:'',
            threeProId:'',

            // 详细信息滚动
            chname:'',      // 仓号名字
            animate:false,
            items:[],                   // 详细信息
            // 左侧煤仓信息
            animateone:false,
            itemsone:[],

            // 计时器
            inter:null,
            interone:null,
            // 检测报告弹框
            notficShow:false,
            imagelists:[],
            // 历史成交
            historyShow:false,
            activities:[], // 成交数据
            // 视频弹窗
            videoShow:false,
            videourl:'',
            videoatitie:'', // 更新时间
            // 云仓报表
            tableData: [],              // 表格数据
            drawer: false,              // 查看报表弹框
            loading:false,
        }
    },
    created(){
        this.http  = this.$store.state.http;                 // IP地址
        this.road();
        // 跳转
        // this.routerboday = setInterval(this.routername,60000);
        setInterval(this.formatDate, 1000);
        this.interone = setInterval(this.scrollone,16000) // 在钩子函数中调用我在method 里面写的scroll()方法
        this.getprop({fstProId:'',bigProId:'',threeProId:''});
    },
    mounted() {},
    methods:{
        // tableRowClassName({row, rowIndex}){//根据报警级别显示颜色
        //     console.log();
        //     // console.log(row);
        //     if(this.isInteger(((rowIndex + 1)/2)*1)) {
        //         return 'warning-row';
        //     }else{
        //         return 'success-row';
        //     }
        //     // return '';
        //     // if(row.row.id){
            
        //     // }
        // },
        isInteger(obj) {
            return Math.floor(obj) === obj
        },

        cellStyle(row){//根据报警级别显示颜色
            // console.log(row);
            // console.log(row.column);
            if(this.isInteger(((row.rowIndex + 1)/2)*1)) {
                return 'background:#1A3D71;color:#fff;border:0;font-size: 16px;'
            }else{
                return 'background:#0D2A55;color:#fff;border:0;font-size: 16px;'
            }
            // if(row.row.id){
            // return 'background:#0B2956;color:#fff;border:0;font-size: 16px;'
            // }
        },
        // 云仓报表
        ehartClick(){
            this.drawer = true;
            this.getpurcon();
        },
        offClick(){
            this.drawer = false;   
            this.tableData = []; 
        },
        // 云仓报表数据
        getpurcon(){
            this.loading = true;
            // 渲染表格
            this.api.busrodect.list()
            .then(res=>{
                if(res.data.code == 200){
                    this.loading = false; 
                    this.tableData = res.data.data;
                }
            })
        },
        // 重置
        elbuttonClick(){
            this.fstName = '';
            this.secName = '';
            this.threeName = '';
            this.fstProId = '';
            this.bigProId = '';
            this.threeProId = '';
            this.road();
            clearInterval(this.interone);
            this.interone = setInterval(this.scrollone,16000) // 在钩子函数中调用我在method 里面写的scroll()方法
            this.getprop({fstProId:'',bigProId:'',threeProId:''});
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
        // 一级清空
        // clearent(){
        //     console.log(111)
        //     this.fstProId = '';
        //     this.bigProId = '';
        //     this.threeProId = "";
        //     this.cityOptions = [];               // 二级小类
        //     this.options = [];                    // 三级小类
        //     this.threeName = null;
        //     this.secName = null;
        //     this.road();
        //     this.getprop({fstProId:'',bigProId:'',threeProId:''});
        // },
        // 一级选中
        enterprises(e){
            this.bigProId = '';
            this.threeProId = "";
            this.cityOptions = [];               // 二级小类
            this.options = [];                    // 三级小类
            this.threeName = null;
            this.secName = null;
            let arr = e.split(",");
            this.fstProId = arr[0]*1;
            this.fstName = arr[1];
            this.list(arr[0]*1);
            this.getprop({fstProId:this.fstProId,bigProId:'',threeProId:''});

        },
        // 二级选中
        entisetwos(e){
            this.options = [];                    // 三级小类
            this.threeName = null;
            this.threeProId = "";
            let arr = e.split(",");
            this.bigProId = arr[0]*1;
            this.secName = arr[1];
            this.twolists(arr[0]*1);
            this.getprop({fstProId:this.fstProId,bigProId:this.bigProId,threeProId:''});
        },
        // 三级选中
        citys(e){
            // console.log(e)
            let arr = e.split(",");
            this.threeProId = arr[0]*1;
            this.threeName = arr[1];
            this.getprop({fstProId:this.fstProId,bigProId:this.bigProId,threeProId:this.threeProId});
            // console.log(this.AddData)
        },
        // 云仓一级大类
        road(){
            this.provinceOptions = [];
            this.api.Prodect.list({parent_id:0})
            .then(res=>{
                if(res.data.code == 200){
                    this.provinceOptions = res.data.data;
                }
            })
        },
        // 仓库产品二级小类类
        list(id){
            this.cityOptions = [];
            this.api.Prodect.list({parent_id:id,fst_pro_id:this.fstProId})
            .then(res=>{
                if(res.data.code == 200){
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 仓库产品三级小类类
        twolists(id){
            this.options = [];
            this.api.Prodect.list({parent_id:id,fst_pro_id:this.fstProId})
            .then(res=>{
                if(res.data.code == 200){
                    this.options = res.data.data;
                }
            })
        },
        // 视频弹窗
        videoClick(data){
            if(data.ckhyFileName){
                this.videoShow = true;
                this.videourl = data.ckhyFileName;
                this.videoatitie = data.updateTime
                // console.log(this.videourl);
            }else{
                this.$message({
                    type: "warning",
                    message: "暂无视频内容!"
                });
            }
        },
        // 视频移出
        videomous(){
            this.videoShow = false;
        },
        // 左侧移入
        mouseovnavs(){
            clearInterval(this.interone);
        },
        // 左侧移出
        mouseleavenavs(){
            this.interone = setInterval(this.scrollone,16000) // 在钩子函数中调用我在method 里面写的scroll()方法
        },
        // 检测报告
        TestEport(id){
            clearInterval(this.interone);
            clearInterval(this.inter);
            this.notficShow = true;
            this.api.addfile.list({f_type:2,relation_id:id.id})
            .then(res=>{
                if(res.data.code == 200){
                    // console.log(res.data.data)
                    this.imagelists = res.data.data;
                }
            })
        },
        // 检测报告移出
        mouselnots(){
           this.notficShow = false; 
            this.inter = setInterval(this.scroll,12000) // 在钩子函数中调用我在method 里面写的scroll()方法
            // this.interone = setInterval(this.scrollone,16000) // 在钩子函数中调用我在method 里面写的scroll()方法
        },
        // 历史成交
        Testhistory(id){
            // console.log(this.interone,this.inter)
            clearInterval(this.interone);
            clearInterval(this.inter);
            this.historyShow = true;
            this.api.order.list({pro_id:id.id})
            .then(res=>{
                // console.log(res)
                if(res.data.code == 200){
                    this.activities = res.data.data;
                }
                
            })
        },
        // 历史成交移出
        mouseTesth(){
           this.historyShow = false; 
            this.inter = setInterval(this.scroll,12000) // 在钩子函数中调用我在method 里面写的scroll()方法
            // this.interone = setInterval(this.scrollone,8000) // 在钩子函数中调用我在method 里面写的scroll()方法
        },
        // 左侧数据列表
        getprop(data){
            // this.itemsone = [];
            let obj = {
                fstProId:data.fstProId,
                bigProId:data.bigProId,
                threeProId:data.threeProId, 
            }
            this.api.busrodect.warelist(obj)
            .then(res=>{
                if(res.data.code == 200){
                    this.itemsone = res.data.data;
                    this.items = [];
                    this.detailed(this.itemsone[0])
                }
            })
        },
        // 左侧信息滚动展示
        scrollone(){
            let con2 = this.$refs.con2;
            con2.style.marginBottom='-189px';
            this.animateone=!this.animateone;
            let that = this;
            setTimeout(function(){
                that.itemsone.push(that.itemsone[0]);
                this.items = [];
                that.itemsone.shift();
                that.detailed(that.itemsone[0])
                con2.style.marginBottom='0px';
                that.animateone=!that.animateone;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },2000)
        },
        // 左侧选中数据
        Checkli(data,index){
            // clearInterval(this.scrollone);
            this.itemsone.splice(index,1);
            this.itemsone.unshift(data);
            // console.log(data,index);
            this.detailed(this.itemsone[0])
            // this.interone = setInterval(this.scrollone,16000) // 在钩子函数中调用我在method 里面写的scroll()方法
            
        },
        // 右侧详细
        detailed(data){
            clearInterval(this.inter);
            this.inter = setInterval(this.scroll,12000) // 在钩子函数中调用我在method 里面写的scroll()方法
            // console.log(data)
            this.chname = data.productName;
            
            this.api.busrodect.list({warehouse_id:data.warehouseId})
            .then(res=>{
                // console.log(res);
                if(res.data.code == 200){
                    this.items = res.data.data;
                }
            })
        },
        // 详细信息滚动展示
        scroll(){
            let con1 = this.$refs.con1;
            con1.style.marginBottom = '-254px';
            this.animate=!this.animate;
            let that = this;
            setTimeout(function(){
                that.items.push(that.items[0]);
                that.items.shift();
                con1.style.marginBottom='0px';
                that.animate=!that.animate;  // 这个地方如果不把animate 取反会出现消息回滚的现象，此时把ul 元素的过渡属性取消掉就可以完美实现无缝滚动的效果了
            },1700)
        },

        // 向右页跳转
        rightGos(){
            this.$router.push({name:'Firstscreen'});
        },
        // // 向左页跳转
        leftGos(){
            this.$router.push({name:'Maps'});
        },
        // 自动循环
        BackTime(){
            // this.backcolor = !this.backcolor;
            // if(this.backcolor){
            //     clearInterval(this.routerboday);
            // }else{
            //     this.routerboday = setInterval(this.routername,60000);
            // }
        },
        // 跳转
        // routername(){
        //     this.$router.push({name:'Thirdscreen'});
        // },
        // 退出大屏6
        BackOff(){
            this.$router.push({name:'Dpmain'});
        },
    },
    // 销毁定时器
    beforeDestroy() {
        clearInterval(this.routerboday);
        if (this.formatDate) {
            clearInterval(this.formatDate); // 在Vue实例销毁前，清除时间屏幕的定时器
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
    position: relative;
    // 报表弹窗
    .echartBox{
        width: 100%;
        height: calc(100vh - 65px);
        position: absolute;
        top: 57px;
        left: 0;
        background: #0B2956;
        // border: 1px solid #69FEF4;
        padding: 18px 0 0 18px;
        box-sizing: border-box;
        overflow: hidden;
        z-index: 999;
        .offbutton{
            width: 120px;
            height: 36px;
            text-align: center;
            font-size: 18px;
            line-height: 36px;
            margin-bottom:12px;
            border-radius: 8px;
            border: 1px solid #69FEF4;
            color: #69FEF4;
            cursor: pointer;
        }
        .gdtBox{
            width: 100%;
            height: calc(100% - 64px);
            overflow: hidden;
            overflow-y: scroll;
            margin-left:8px;
        }

    }
    // 视频弹框样式开始
    .videoBox{
        width: 950px;
        height: 600px;
        position: absolute;
        top: 15%;
        left: 20%;
        background: #0B2956;
        border: 1px solid #69FEF4;
        .videotitie{
            width: auto;
            height: 40px;
            position: absolute;
            padding-left:10px;
            box-sizing: border-box;
            top: 0;
            left: 0;
            z-index: 99;
            font-size: 24px;
            line-height: 40px;
            color: #69FEF4;
        }
    }
    .notification{
        width: 500px;
        height: 700px;
        background: #0B2956;
        border: 1px solid #69FEF4;
        position: absolute;
        top: 12%;
        left: 20%;
        z-index: 99;
        // transition: all 2s;
        // transform: translate(500px,700px); // 向上浮动了300个像素
        img{
            width: 100%;
            // height: 100%;
        }
    }
    .nothistory{
        width: 500px;
        height: 700px;
        background: #0B2956;
        border: 1px solid #69FEF4;
        padding: 24px;
        box-sizing: border-box;
        position: absolute;
        top: 12%;
        right: 20%;
        z-index: 99;
        overflow: hidden;
        overflow-y: scroll;
        // transition: all 2s;
        // transform: translate(500px,700px); // 向上浮动了300个像素
        .cards{
            width: 100%;
            height: 160px;
            display: flex;
            flex-wrap: wrap;
            padding-left:14px;
            box-sizing: border-box;
            color: #FDFBFE;
            .cardflex{
                width: 100%;
                height: 124px;
                font-size: 22px;
                line-height: 50px;
                word-wrap:break-word;
                word-break: break-all;
            }
            .cardone{
                width: 25%;
                height: 36px;
                line-height: 36px;
                font-size: 16px;
            }
            .cardtwo{
                width: 30%;
                height: 36px;
                line-height: 36px;
                font-size: 16px;
            }
            .cardtree{
                width: 45%;
                height: 36px;
                line-height: 36px;
                font-size: 16px;
            }
        }
    }
    // 弹框样式结束
    .rightGo{
        width: 50px;
        height: 50px;
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
            font-size: 50px;
            line-height: 50px;
        }
    }
    .leftGo{
        width: 50px;
        height: 50px;
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
            font-size: 50px;
            line-height: 50px;
        }
    }
    // 头部
    .header{
        width: 100%;
        height: 51px;
        background-image: url(../../assets/yctitle.png);
        background-size:1600 51px;
        background-repeat: no-repeat;
        background-position: center;
        text-align: center;
        position: relative;
        @font-face {
            font-family: aa;
            src: url(../../assets/icon/FZMeiHei-M07.ttf);
        }
        p{
            font-family: aa;
            font-size: 38px;
            line-height: 51px;
            color: #fff;
            font-weight: 500;
        }
        .TimesAllleft{
            position: absolute;
            top: 20%;
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
    .middle{
        width: 100%;
        height: 50px;
        padding-top:10px;
        .echartsbutton{
            width: 120px;
            height: 40px;
            float: left;
            text-align: center;
            font-size: 18px;
            line-height: 40px;
            margin:5px 10px;
            border-radius: 8px;
            border: 1px solid #69FEF4;
            color: #69FEF4;
            cursor: pointer;
        }
        .middleBox{
            width:50%;
            height: 50px;
            float: right;
            .el-select{
                float: right;
            }
            .elbutton{
                width: 80px;
                height: 40px;
                float: right;
                text-align: center;
                font-size: 18px;
                line-height: 40px;
                margin:5px 10px;
                border-radius: 8px;
                border: 1px solid #69FEF4;
                color: #69FEF4;
                cursor: pointer;
            }
        }
    }
    // 内容
    .nav{
        width: 100%;
        height: calc(100% - 117px);
        display: flex;
        flex-wrap:wrap;
        padding: 18px;
        box-sizing: border-box;
        // background: chartreuse;
        .navleft{
            width: 55%;
            height: 100%;
            float: left;
            padding:0 34px 34px 0;
            box-sizing: border-box;
            overflow: hidden;
            .animone{
                transition: all 2s;
                transform: translateY(-189px); // 向上浮动了300个像素
            }
            #con2{
                width: 100%;
                display: flex;
                flex-wrap: wrap;
                li{
                    width:100%;
                    // height: 170px;
                    margin-bottom:24px;
                    box-sizing: border-box;
                    text-align: left;
                    cursor: pointer;
                    // 选中样式
                    .navli{
                        width: 100%;
                        height: 100%;
                        background-image: url(../../assets/li.png);
                        background-size: 100% 100%;
                        background-repeat: no-repeat;
                        opacity: 1;
                        // 左侧视频
                        .flexnasleft{
                            width: 220px;
                            height: 100%;
                            float: left;
                            position: relative;
                            img{
                                width: 100%;
                                height: 100%;
                            }
                            // .dlexicon{
                            //     width: 60px;
                            //     height: 60px;
                            //     position: absolute;
                            //     top: 0;
                            //     left: 0;
                            //     right: 0;
                            //     bottom: 0;
                            //     margin: auto;
                            //     i{
                            //         font-size: 60px;
                            //         color: #fff;
                            //     }
                            // }
                            .flexnastitle{
                                position: absolute;
                                top: -1px;
                                left: 0;
                                width: auto;
                                height: 42px;
                                padding: 0 14px;
                                box-sizing: border-box;
                                background-image: url(../../assets/Rectangle.png);
                                background-size: 100% 100%;
                                background-repeat: no-repeat;
                                font-size:20px ;
                                text-align: center;
                                line-height: 42px;
                                color: #69FEF4;
                                letter-spacing:0.2em;
                            }
                            .flexbottomsize{
                                width: 220px;
                                height: 37px;
                                position: absolute;
                                bottom: -1px;
                                left: 0;
                                background: rgba(11, 42, 86, 0.8);
                                font-size: 20px;
                                text-align: center;
                                line-height: 37px;
                                color: #69FEF4;
                                letter-spacing:0.2em;
                            }
                        }
                        // 右侧内容
                        .flexnasright{
                            width: calc(100% - 220px);
                            height: 100%;
                            float: right;
                            padding: 20px;
                            box-sizing: border-box;
                            .flexnasrightBox{
                                width: 100%;
                                height: 100%;
                                background: rgba(105, 254, 244, 0.2);
                                padding: 10px;
                                box-sizing: border-box;
                                border-radius: 8px;
                                word-wrap:break-word;
                                word-break: break-all;
                                span{
                                    font-size: 18px;
                                    line-height: 36px;
                                    color: #69FEF4;
                                }
                            }
                        }
                    }
                    // 未选中样式
                    .navbloneli{
                        width: 95%;
                        height: 98%;
                        margin-left:2.5%;
                        background-image: url(../../assets/liteo.png);
                        background-size: 100% 100%;
                        background-repeat: no-repeat;
                        // opacity: .9;
                        // 左侧视频
                        .flexnaslefts{
                            width: 220px;
                            height: 100%;
                            float: left;
                            position: relative;
                            img{
                                width: 100%;
                                height: 100%;
                            }
                            .flexnastitles{
                                position: absolute;
                                top: -1px;
                                left: 0;
                                width: auto;
                                height: 42px;
                                padding: 0 14px;
                                box-sizing: border-box;
                                background-image: url(../../assets/Rectangle.png);
                                background-size: 100% 100%;
                                background-repeat: no-repeat;
                                font-size:20px ;
                                text-align: center;
                                line-height: 42px;
                                color:  rgba(255, 255, 255, 0.5);;
                                letter-spacing:0.2em;
                            }
                            .flexbottomsizes{
                                width: 220px;
                                height: 37px;
                                position: absolute;
                                bottom: -1px;
                                left: 0;
                                background: rgba(11, 42, 86, 0.8);
                                font-size: 20px;
                                text-align: center;
                                line-height: 37px;
                                color:  rgba(255, 255, 255, 0.5);;
                                letter-spacing:0.2em;
                            }
                        }
                        // 右侧内容
                        .flexnasrights{
                            width: calc(100% - 220px);
                            height: 100%;
                            float: right;
                            padding: 20px;
                            box-sizing: border-box;
                            .flexnasrightBoxs{
                                width: 100%;
                                height: 100%;
                                background: #0B2956;
                                padding: 10px;
                                box-sizing: border-box;
                                border-radius: 8px;
                                word-wrap:break-word;
                                word-break: break-all;
                                span{
                                    font-size: 20px;
                                    line-height: 36px;
                                    color:  rgba(255, 255, 255, 0.5);;
                                }
                            }
                        }
                    }
                }
            } 
        }
        .navright{
            width: 45%;
            height: 100%;
            background-image: url(../../assets/liright.png);
            background-size: 100% 100%;
            background-repeat: no-repeat;
            float: left;
            padding: 20px;
            box-sizing: border-box;
            overflow: hidden;
            .navtitle{
                width:100%;
                height: 55px;
                padding-left:110px;
                background-image: url(../../assets/navtitle.png);
                background-size:100% 55px;
                background-repeat: no-repeat;
                font-size: 32px;
                line-height: 55px;
                color: #fff;
                letter-spacing:0.6em;
                margin-bottom:20px;
                box-sizing: border-box;
            }
            // 详细信息整个盒子大小
            #sectionBox{
                width: 100%;
                height: calc(100% - 75px);
                padding: 0 23px;
                box-sizing: border-box;
                overflow: hidden;
                overflow-y: scroll;
                border-radius: 8px;
                margin-left:28px;
            }
            .anim{
                transition: all 2s;
                transform: translateY(-254px); // 向上浮动了300个像素
                
            }
            #con1{
                border-radius: 4px;
                li{
                    width: 100%;
                    line-height: 32px;
                    // height: 254px;
                    display: flex;
                    flex-wrap: wrap;
                    background: #0B2956;
                    border-bottom:1px solid #1F4783;
                    padding:0 23px;
                    box-sizing: border-box;
                    .section{
                        width: 48%;
                        height: 48px;
                        background: rgba(175, 48, 48, .3);
                        border-radius: 8px;
                        margin:12px 2% 16px 0;
                        font-size: 18px;
                        line-height: 48px;
                        text-align: center;
                        color: #fff;
                    }
                    .sectionone{
                        width: 36px;
                        height: 36px;
                        margin:22px 0 0 0;
                        background: #3063AF;
                        color: #fff;
                        border-radius: 50%;
                        font-size: 24px;
                        line-height: 36px;
                        text-align: center;
                    }
                    .secttitle{
                        width: calc(100% - 36px);
                        height: 48px;
                        font-size: 24px;
                        line-height: 48px;
                        color: #FFF;
                        padding-left:10px;
                        box-sizing: border-box;
                        margin-top: 16px;
                        i{
                            font-size: 24px;
                            line-height: 48px;
                            cursor: pointer;
                        }
                        .sectbutton{
                            width:auto;
                            height: 48px;
                            margin-right:10px;
                            background-image: url(../../assets/selno.png);
                            background-size:100% 55px;
                            background-repeat: no-repeat;
                            line-height: 48px;
                            padding: 0 6px;
                            font-size: 18px;
                            color: #fff;
                            float: right;
                            text-align: center;
                            cursor: pointer;
                        }
                    }

                    .sectconent{
                        width: 100%;
                        height: calc(100% - 160px);
                        margin:0px 0 20px;
                        font-size: 22px;
                        line-height: 36px;
                        color: #fff;
                        word-wrap:break-word;
                        word-break: break-all;
                        span{
                            font-size: 24px;
                            line-height: 36px;
                            color: #fff;
                        }
                    }
                }
            } 
        }
     
    }
}
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
