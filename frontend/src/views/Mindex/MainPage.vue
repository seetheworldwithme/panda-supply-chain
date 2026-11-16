<template>
  <!-- 账户维护页面 -->
  <div
    v-loading="loading"
    element-loading-text = "加载中...">
    <!-- 头部开始 -->
    <div class="pur-top">
        <span class="titleSpan">{{timedata}} {{username}}</span>
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
                style="width:40%;margin-right:15px;"
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
                style="width:148px;"
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
        <div class="Importfile" @click="Importdailog">
            <i class="el-icon-ali-daoru"></i>
            <span>一键导入</span>
        </div>
      <!-- 导入弹窗功能 -->
      <el-dialog
        title="一键导入"
        :visible.sync="ImportVisible"
        width="60%"
        :before-close="ImporthandleClose">
        <div slot="title" class="Timedialog">
            一键导入
        </div>
        <div class="Timedialog-box">
            <el-steps :active="active" align-center process-status="wait" finish-status="finish">
            <el-step title="上传文件"></el-step>
            <el-step title="导入数据"></el-step>
            <el-step title="完成导入"></el-step>
            </el-steps>
            <!-- 步骤一 -->
            <div class="importone" v-show="oneShow">
            <el-form :model="fileForm" :rules="rules" ref="fileForm">
                <el-col :span="3" style="margin-left:35px;">
                    <span style="color:red;">* </span>
                    <span style="line-height: 32px;">选择导入文件</span>
                </el-col>
                <el-col :span="19">
                    <el-form-item prop="importnav" style="margin-bottom:0px;">
                        <el-input placeholder="仅限xlsx文件" v-model="fileForm.importnav" size="small">
                        <!-- importOnce    importExcelCg -->
                            <el-upload slot="append"
                            ref='upload'
                            :auto-upload="resfile"
                            :action="http + '/saas-contract/chainContract/importOnce'"
                            :headers="header"
                            name='file'
                            :show-file-list="false"
                            :on-success="addFileLet"
                            :on-change="supbeforedemoUpload"
                            :limit="1"><el-button size="mini" style="color: #333;">选择文件</el-button>
                            </el-upload>
                        </el-input>
                    </el-form-item>
                </el-col>
                <el-col :span="24">
                    <div class="downBox">
                        <a :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=YJDR.xlsx'">下载导入模板</a>
                    </div>
                    <!-- 一键删除导入文件 -->
                    <div class="importclear">
                        <el-col :span="3" style="margin-left:35px;">
                            <span style="line-height:32px;padding-left:12px;">删除导入文件</span>
                        </el-col>
                        <el-col :span="19">
                            <el-input placeholder="请输入已完成导入文件的批次号" v-model="importclearnumber" size="small">
                                <el-button size="mini" slot="append" @click="numberClearFile" style="color: #333;">确认删除</el-button>
                            </el-input>
                        </el-col>
                    </div>
                    <div class="importone-nav">
                    <span>注意事项： </span>
                    <p style="color:#CF1724;">1.请勿擅自修改导入模板的表头字段，否则会导致导入失败!</p>
                    <p>2.尽量格式书写规范，字体要统一字体，模版不允许有空白行，若有请选择删除 </p>
                    <p>3.字段类型为多选项的，各 个选项用逗号隔开，且请确保各 个选项与系统自定义的值匹配</p>
                    <p>4.日期或月份格式为XXXX/XX/XX，XXXX-XX-XX</p>
                    <p>5.必填字段若无数据时请输入“---”</p>
                    </div>
                </el-col>
                <div style="clear:both;"></div>
                <div class="importone-bottom">
                <el-button @click="Importble" size="small">取 消</el-button>
                <el-button type="primary" @click="ImportNext('fileForm')" size="small">导 入</el-button>
                </div>
            </el-form>
            </div>
            <!-- 步骤二 -->
            <div class="importtwo" v-show="twoShow">
            <div class="importtwo-Box">
                <el-progress type="circle" :percentage="perstep" color="#CF1724"></el-progress>
                <span>{{pernav}}</span>
            </div>
            </div>
            <!-- 步骤三...失败 -->
            <div class="importstree" v-show="streeShow">
                <i class="el-icon-cuowu"></i>
                <span> 导入格式错误,错误原因：</span>
                <div class="error-size">
                    <p v-for="(v,k) in cation" :key="k"> {{ v}}</p>
                </div>
                <div class="stree-ok">
                    <el-button size="mini" @click="ImportNextstree">关闭</el-button>
                </div>
            </div>
            <!-- 步骤三...成功 -->
            <div class="importfour" v-show="fourShow">
                <div class="importfour-Box">
                    <div v-if="sucessShow">
                        <el-progress type="circle" :percentage="100" status="success"></el-progress>
                        <span>数据完成导入，请在列表中点击查看</span>
                    </div>
                    <div v-if="errorShow" class="err-box">
                        <div v-for="(vol,ki) in sucessData" :key="ki">
                            <el-alert :title="vol + '条'" type="success" :closable="false"></el-alert>
                        </div>
                        <div style="clear:both"></div>
                        <div v-for="(val,kia) in errorData" :key="'info-'+ kia">
                            <el-alert
                                :title="val + '条'"
                                type="error"
                                :closable="false">
                            </el-alert>
                        </div>
                        <div style="clear:both"></div>
                        <div class="error-size">
                            <p v-for="(v,k) in cation" :key="k"> {{ v}}</p>
                        </div>
                    </div>
                </div>
                <div class="four-ok">
                    <el-button size="mini" @click="ImportNextstree">完成</el-button>
                </div>
            </div>
        </div>
      </el-dialog>
    </div>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <!-- 物流地图 -->
        <div class="cityMap">
            <div id="maps"></div>
        </div>  
        <!-- 实体企业简介 -->
        <div class="citynews">
            <div class="citynewsTitle">实体企业简介</div>
            <div class="citynewsImage">
                <div class="img">
                    <el-carousel trigger="click" height="94px" indicator-position="none" arrow="never" :interval="2000">
                        <el-carousel-item v-for="(item,index) in images" :key="index">
                            <img :src="item.url" alt="" style="width:100%;height:100%;">
                        </el-carousel-item>
                    </el-carousel>
                </div>
                <ul class="citynewsBox">
                    <li>
                        <b>名称：</b>
                        <span>{{creditData.name}} </span>
                    </li>
                    <li>
                        <b>缴税信用：</b>
                        <span>{{credit}}</span>
                    </li>
                    <li>
                        <b>煤种：</b>
                        <span>{{creditData.pro}}</span>
                    </li>
                    <li>
                        <b>网络货运平台：</b>
                        <span>{{creditData.carriagePlatform}}</span>
                    </li>
                    <li style="width: 100%;">
                        <b>地址：</b>
                        <span>{{creditData.address}}</span>
                    </li>
                </ul>
            </div>
            <div class="citynewsEchartBox">
                <div id="citynewsEchart"></div>
            </div>
        </div>
        <!-- 采购数据 -->
        <div class="Purchasing">
            <span class="tity">采购数据</span>
            <div class="Purchasing-time">
                <el-dropdown trigger="click">
                    <span>{{settlementtime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                    <el-dropdown-menu slot="dropdown">
                        <div>
                            <el-dropdown-item v-for="(v,k) in dropdown" :key="'1' + k" @click.native="settlemenchange(v.name,v.type)">{{v.name}}</el-dropdown-item>
                        </div>
                    </el-dropdown-menu>
                </el-dropdown>
            </div>
            <ul class="Purchasing-conent">
                <li>
                    <img src="@/assets/query1.svg" alt="">
                    <span>采购合同份数</span>
                    <div class="styletop">
                        <h5 v-if="chasingTonnage">{{chasingTonnage|currency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <p>采购合同吨数</p>
                    </div>
                </li>
                <li>
                    <h4 v-if="chasingnum"> {{chasingnum}} </h4>
                    <h4 v-else> 0 </h4>
                    <div class="styletop">
                        <h5 v-if="arrsumNum">{{arrsumNum|currency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <p>采购结算吨数</p>
                    </div>
                </li>  
            </ul>
        </div>
        <!-- 销售数据 -->
        <div class="Purchasing">
            <span class="tity">销售数据</span>
            <div class="Purchasing-time">
                <el-dropdown trigger="click">
                    <span>{{saletime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                    <el-dropdown-menu slot="dropdown">
                        <div>
                            <el-dropdown-item v-for="(vol,kol) in dropdown" :key="'2' +kol" @click.native="salechange(vol.name,vol.type)">{{vol.name}}</el-dropdown-item>
                        </div>
                    </el-dropdown-menu>
                </el-dropdown>
            </div>
            <ul class="Purchasing-conent">
                <li>
                    <img src="@/assets/query2.svg" alt="">
                    <span>销售合同份数</span>
                    <div class="styletop">
                        <h5 v-if="saleTonnage">{{saleTonnage|currency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <p>销售合同吨数</p>
                    </div>
                </li>
                <li>
                    <h4 v-if="salenum"> {{salenum}} </h4>
                    <h4 v-else> 0 </h4>
                    <div class="styletop">
                        <h5 v-if="sumNum">{{sumNum|currency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <p>销售结算吨数</p>
                    </div>
                </li>
            </ul>
        </div>
        <!-- 提示看板 -->
        <div class="Tips">
            <span class="tity">提示看板</span><img src="../../assets/tishi.png" alt="">
            <div class="tips-box" v-if="tipslist.length == 0">
                <p>暂无重要提醒数据<br/>实时更新</p>
            </div>
            <div v-else>
                <ul class="tips-list" style="overflow:auto">
                    <li v-if="tipslist !=[]" v-for="(i,ik) in tipslist" :key="'3' +ik" >
                        <div>
                            <div class="dian"></div>
                            <span>{{i.time}}<span style="color:#CF1724;">{{ i.att }}</span>{{i.last}}</span>
                            
                        </div>
                    </li>
                    <li v-if="newtipslist !=[]" v-for="(val,kal) in newtipslist" :key="'4' +kal" >
                        <div>
                            <div class="dian"></div>
                            <span>{{val.time}}<span style="color:#CF1724;">{{ val.att }}</span>{{val.last}}</span>
                            
                        </div>
                    </li>
                </ul>
            </div>
        </div>
        <!-- 筛选器 -->
        <div class="screenBox">
            <b>筛选器</b>
            <el-button icon="el-icon-search" size="mini" @click="dialogVisibleclick">{{TimeFilter}}</el-button>
        </div>
        <!-- 税额统计 -->
        <div class="quantity">
            <span class="tity">税额统计</span>
            <!-- <div class="quantity-time">
                <el-dropdown trigger="click">
                    <span>{{quantime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                    <el-dropdown-menu slot="dropdown">
                        <div>
                            <el-dropdown-item v-for="(v,k) in dropdown" :key="'5' +k" @click.native="dropdownchange(v.name,v.type)">{{v.name}}</el-dropdown-item>
                        </div>
                    </el-dropdown-menu>
                </el-dropdown>
            </div> -->
            <ul class="quantity-conent">
                <li>  
                    <div class="dian" style="background-color:#FB7F5B;"></div>
                    <div style="float:left;">
                        <h5 v-if="sumTaxAmount">{{sumTaxAmount|newcurrency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <span>进项税额 (万元)</span>
                    </div>
                </li>
                <li>
                    <div class="dian" style="background-color:#FAD337;;"></div>
                    <div style="float:left;">
                        <h5 v-if="arrTaxAmount">{{arrTaxAmount|newcurrency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <span>销项税额 (万元)</span>
                    </div>
                </li>
                <li>
                    <div class="dian" style="background-color:#FB7F5B;"></div>
                    <div style="float:left;">
                        <h5 v-if="sumIncludMoney">{{sumIncludMoney|newcurrency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <span>进项含税金额 (万元)</span>
                    </div>
                </li>
                <li>
                    <div class="dian" style="background-color:#FAD337;"></div>
                    <div style="float:left;">
                        <h5 v-if="arrIncludMoney">{{arrIncludMoney|newcurrency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <span>销项含税金额 (万元)</span>
                    </div>
                </li>
                <li>
                    <div class="dian" style="background-color:#FAD337;"></div>
                    <div style="float:left;">
                        <h5 v-if="rateBill">{{rateBill|newcurrency}}</h5>
                        <h5 v-else> 0.00 </h5>
                        <span>缴税金额累计 (万元)</span>
                    </div>
                </li>
            </ul>
        </div>
        <!-- 税负率 -->
        <div class="burden">
            <div id="burden"></div>
            <div class="burden-time">
                <span>{{burdenTime[0]}} ~ {{burdenTime[this.burdenTime.length - 1] }} 累计税负率为 <span v-if="burdenCumulative >= 10" style="color:#CF1724;">{{burdenCumulative}} %</span><span v-else>{{burdenCumulative}} %</span></span>
                <!-- <el-date-picker
                style="width:50%;float: right;"
                size="small"
                :clearable="false"
                unlink-panels
                v-model="burdeninvoice"
                :picker-options="pickerOptions0"
                type="daterange"
                @change="burdenchange"
                value-format="yyyy-MM-dd"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期">
                </el-date-picker> -->
            </div>
        </div>
        <!-- 票据 -->
        <div class="invoice">
            <span>票据</span>
            <!-- <el-date-picker
            style="width:50%;float: right;"
            size="small"
            v-model="invoice"
            :clearable="false"
            type="daterange"
            unlink-panels
            @change="invoicechange"
            :picker-options="pickerOptions0"
            value-format="yyyy-MM-dd"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期">
            </el-date-picker> -->
            <div id="invoice"></div>
        </div>
        <!-- 结算额 -->
        <div class="settlement">
            <div id="settlement"></div>
            <div class="settlement-right">
                <!-- <div class="quantity-time">
                    <el-dropdown trigger="click">
                        <span>{{setttime}}<i class="el-icon-arrow-down el-icon--right"></i></span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item v-for="(v,k) in settdropdown" :key="k" @click.native="settChange(v.name,v.type)">{{v.name}}</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </div> -->
                <div style="clear:both;"></div>
                <div class="dian" style="background-color: #72CE99;"></div><span class="dian-size">采购结算额累计</span>
                <div style="clear:both;"></div>
                <h5 v-if="serienum"> ￥{{serienum|newcurrency}}</h5>
                <h5 v-else> ￥ 0.00</h5>
                <div class="dian" style="background-color:#58AFFF;"></div><span class="dian-size"> 销售结算额累计</span>
                <div style="clear:both;"></div>
                <h5 v-if="riesnum"> ￥{{riesnum|newcurrency}}</h5>
                <h5 v-else> ￥ 0.00</h5>
            </div>
        </div>
        <!-- 收付 -->
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
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
var echarts = require('echarts');
import china from '../../assets/map/china.json'
// import jingjinji from '../../assets/map/jingjinji.json'
export default {
    data() {
        return {
            credit:'优秀',
            creditData:{},
            images:[],
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
            screenlist:[1,2,3,4,5,6],               // 类别
            newscreenlist:[1,2,3,4,5,6],            // 选中类别
            screeninvoice:[],                       // 时间
            newscreeninvoice:[],                    // 选中时间
            // 类别
            screenOptions: [
                {
                    value: 1,
                    label: '票据'
                },
                {
                    value: 2,
                    label: '收付'
                }, 
                {
                    value: 3,
                    label: '税负率'
                }, 
                {
                    value: 4,
                    label: '物流'
                }, 
                {
                    value: 5,
                    label: '结算额'
                }, 
                {
                    value: 6,
                    label: '税额统计'
                },
            ],

            // 时间筛选
            loading:false,
            username:'',
            timedata:'',
            dropdown:[
                {name:'本月',type:1},
                {name:'上月',type:2},
                {name:'近半年',type:3},
                {name:'全部',type:5},
            ],

            // 采购煤价趋势
            citynewdataArr:[],                 // 时间
            citynewseriesData:[],              // 采购数据
            citynewserienum :0,
            citynewriesData:[],                // 销售数据

            // 采购数据
            settlementtime:'全部',
            chasingnum:'',              // 合同份数
            chasingTonnage:'',          // 合同吨数
            arrsumNum:'',               // 结算吨数

            // 销售数据
            saletime:'全部',
            salenum:'',                 // 合同份数
            saleTonnage:'',             // 合同吨数
            sumNum:'',                  // 结算吨数

            // 提示看板
            count: 0,
            tipslist:[],                // 看板数据
            newtipslist:[],

            // 票据含税金额
            quantime:'上月',
            sumIncludMoney:'',          // 进项含税金额
            arrIncludMoney:'',          // 销项含税金额

            // 税负率
            burdeninvoice:'',
            pickerOptions0:{
                disabledDate:(time)=>{
                    return time.getTime() > Date.now() || time.getTime() == Date.now();
                }
            },
            burdenTime:[],              // 时间
            burdenData:[],              // 税负率
            burdenCumulative:0,         // 累计税负率

            // 票据税额
            invoice:'',
            sumTaxAmount:'',            // 进项税额
            arrTaxAmount:'',            // 销项税额
            rateBill:'',                // 缴税累计金额

            // 结算
            setttime:'近半年',
            dataArr:[],                 // 时间
            seriesData:[],              // 采购数据
            serienum :0,
            riesData:[],                // 销售数据
            riesnum :0,
            // settdropdown:[
            //     {name:'本年',type:4},
            //     {name:'近半年',type:3},
            //     {name:'全部',type:5},
            // ],

            // 收付
            paymenttime:'本月',
            datapayment:[],             // 时间
            paymentData:[],             // 收款数据
            paymentserData:[],          // 付款数据
            paymentdropdown:[
                {name:'本年',type:4},
                {name:'本月',type:1},
                {name:'全部',type:5},
            ],

            // 物流
            logisticstime:'近半年',
            datalogistics:[],           // 时间
            logisticsData:[],           // 物流自提
            Selfnum:0,                  // 自提总吨数
            mountsData:[],              // 物流自提
            Tonum:0,                    // 送到总吨数

            // 一键导入
            ImportVisible:false,        // 导入弹窗
            ImportVisible:false,        // 取消导入
            header:{},                  // 上传token
            resfile:false,              // 是否立即上传
            active:1,                   // 步骤条
            oneShow:true,               // 第一步
            // 导入上传验证
            fileForm:{
                importnav:'',
            },
            rules: {
                importnav: [
                    { required: true, message: '请上传xlsx 或 xls文件', trigger: 'blur' },
                ],
            },
            twoShow:false,              // 第二步
            perstep:0,                  // 进度%
            pernav:'正在导入...',
            streeShow:false,            // 第三步...失败
            cation:'',                  // 导入错误提示
            fourShow:false,             // 第三步...成功

            sucessShow:false,           // 第三步...全部成功
            errorShow:false,            // 第三步...非全部成功
            sucessData:[],              // 正确数据
            errorData:[],               // 错误数据 

            // 删除导入文件批次号
            importclearnumber:'',
        };
    },
    created() {
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
        // 提示看板
        this.tskbAll();
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
        // 获取初始年月日
        let nowDate = year + "-" + month + "-" + day;
        // this.invoice = ["2019-10-01",nowDate];
        // this.burdeninvoice = ["2019-10-01",nowDate];
        this.screeninvoice = ["2019-10-01",nowDate];
        this.TimeFilter = "2019-10-01 ~ " + nowDate;
        this.newscreeninvoice = ["2019-10-01",nowDate];
        let invoiceDate = {
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
        this.api.statistics.selEtcMsg()
        .then(res=>{
            if(res.data.code == 200){
                this.creditData = res.data.data;
            }
        })
        this.api.file.imageShow({relation_id:sessionStorage.getItem('companyId')})
        .then(res=>{
            if(res.data.code == 200){
                for(let i in res.data.data.lbts){
                    this.images.push(res.data.data.lbts[i])
                }
            }
        })
        this.getMapsEchart();                      // 地图图表 
        // 采购煤价趋势
        this.settcitynews(invoiceDate);
        // 票据税额接口
        this.countBill(invoiceDate);
        // 税负率
        this.burdenapi(invoiceDate);
        // 收付
        this.paymentCapital(invoiceDate);
        // 物流
        this.logis(invoiceDate);
        // 结算额
        this.settCapital(invoiceDate);
        // 税额统计接口
        this.countBillss(invoiceDate);
        // 采购合同数据
        this.countCapital(5);
        // 采购结算量数据
        this.countlistByPage(5);
        
    },
    methods: {
        /* 
        地图图表
        */
        getMapsEchart() {
            var chartsmyChart = echarts.init(document.getElementById('maps'));
            echarts.registerMap('china',china)
            // echarts.registerMap('jingjinji',jingjinji);
            this.api.statistics.countLogistArea()
            .then(res=>{
            var a = {};
            var b = [];
            var c = [];
            var color = [];//航线的颜色
            var d = [];
            for(let i=0;i<res.data.data.length;i++){
                a[res.data.data[i].loadingArea] = res.data.data[i].loadingLnglog.split(",").map(Number);
                a[res.data.data[i].unloadingArea] = res.data.data[i].unloadingLnglog.split(",").map(Number);
                if(d.length>0 && d[d.length-1][0].name != res.data.data[i].loadingArea){
                    c.push([d[0][0].name,d])
                    // let r = Math.floor(Math.random() * 256); //随机生成256以内r值
                    // let g = Math.floor(Math.random() * 256); //随机生成256以内g值
                    // let l = Math.floor(Math.random() * 256); //随机生成256以内b值
                    // // let alpha = Math.random(); //随机生成1以内a值
                    // color.push(`rgb(${r},${g},${l})`) //返回rgba(r,g,b,a)格式颜色
                    color.push('rgb(166,200,76)')
                    d = []
                }
                b.push([{name:res.data.data[i].loadingArea}, {name:res.data.data[i].unloadingArea,value:res.data.data[i].logCount}])
                d.push(b[i]);
            }
            for(let o=0;o<c.length;o++){
                c[o][1].push([{name:''}, {name:c[o][0],value:100}])
            }
            var planePath = 'path://M422.50606886 606.35552995c3.56517245 5.18570537 11.91091704 4.49697889 13.97709653 4.2133856 0.20256662 35.40864456 0.40513324 59.55458523 0.40513324 59.55458524s-2.51182603 22.80900101 21.02641479 28.15675969 80.13535347 10.69551735 112.2219055 0c0 0 16.3673826-4.98313877 17.46124234-26.37417346 0.28359327-5.30724536 0.40513323-28.84548619 0.48615987-61.25614481 2.75490597 0.28359327 10.0878175 0.5266732 13.36939668-4.29441227 3.03849923-4.41595223-5.87443188-2.95747261-5.87443188-2.9574726s-4.33492559 0.81026648-7.4949648 1.13437305c0.24307993-99.82482858-0.1620533-273.58647217-0.16205328-273.58647217s0.24307994-11.06013726-18.87920866-11.06013724c-8.18369131 0-100.47304176-0.3646199-100.47304176-0.36461992s-29.21010609-1.05334641-31.72193212 10.33089745c-1.78258622 8.14317798-1.01283308 178.1775958-0.44564658 274.7208452-3.20055254-0.32410658-8.06215133-1.2153997-8.06215131-1.21539968s-8.8724178-1.41796632-5.83391855 2.99798592zM450.41974859 684.91086378c-10.41192409-7.73804476-9.35857768-13.49093666-9.35857767-13.49093664s0.81026648-9.03447109 13.36939669-6.23905179S473.3908029 691.87915539 473.3908029 691.87915539s-12.5186169 0.76975314-22.97105431-6.9682916z m25.40185372-26.05006687c-2.14720613-3.4841458-0.5266732-7.90009805-0.52667321-7.90009802l35.81377779 3.20055252 0 7.90009805c-0.04051332 0-33.13989845 0.28359327-35.28710459-3.20055255zM517.388272 695.84946108s-1.17488637 1.01283308-4.53749223 1.49899295c-4.49697889 0.64821317-5.9149452-1.49899296-5.91494517-1.49899295s-1.78258622 0.81026648-15.88122274 0c-18.55510208-1.05334641-12.43759026-6.27956512-12.43759024-6.27956512 4.29441227 3.60568578 28.60240623 3.4841458 28.60240623 3.48414582 6.68469833-3.72722575 10.00679084 0 10.00679085 0 9.48011766 2.79541931 30.87115236-3.4841458 30.87115235-3.48414581-3.76773908 8.22420464-30.70909905 6.27956512-30.70909906 6.27956511z m29.89883257-37.27225743c-2.79541931 3.07901257-34.19324486 3.4841458-34.19324485 3.48414581l0-8.02163802c2.67387934 0.64821317 35.00351133-3.07901257 35.00351133-3.07901255l-0.81026648 7.61650476zM550.89279036 691.83864207c2.39028607-9.88525088 14.70633635-23.65978081 17.90688889-25.76647362 3.20055254-2.14720613 13.08580344-0.64821317 13.89606988 6.15802514 1.90412619 16.32686928-31.80295877 19.60844847-31.80295877 19.60844848zM585.89630167 604.77551034s-3.20055254 7.73804476-5.87443187 0.81026645-5.2262187-30.18242584-5.22621869-30.18242584l10.16884414 3.72722574 0.93180641 25.64493365z m-11.10065056-34.35529817c-1.98515284-9.23703771 0-44.60516894 0-44.60516893 5.06416541-1.86361286 8.021638 6.80623831 8.021638 6.80623831l1.98515284 42.74155607-10.00679084-4.94262545zM568.51608598 354.40317241c0-4.13235898 3.36260583-7.49496481 7.49496481-7.49496481s7.49496481 3.36260583 7.4949648 7.49496481l0 38.20406387-14.94941628 0L568.55659931 354.40317241z m0 40.67537659l14.94941629 0L583.46550228 414.76802411l-14.9494163 0-1e-8-19.68947511z m-1e-8 22.20130115l14.94941631 0-1e-8 19.60844847-14.9494163 0 0-19.60844847z m0 22.07976121l14.94941631 0-1e-8 18.6361287-14.9494163 0 0-18.6361287z m0 21.10744143l14.94941631 0-1e-8 19.16280192-14.9494163 0 0-19.16280192z m0 21.63411463l14.94941631 0-1e-8 19.28434188-14.9494163 0 0-19.28434188z m0 21.7556546l14.94941631 0-1e-8 4.98313877c0 4.13235898-3.36260583 7.49496481-7.4949648 7.4949648s-7.49496481-3.36260583-7.49496481-7.4949648l0-4.98313877z m-74.54451483-134.74731325c0-1.78258622 1.45847963-3.20055254 3.20055254-3.20055254l33.95016491 0c1.78258622 0 3.20055254 1.45847963 3.20055254 3.20055254l-1e-8 33.13989844c0 1.78258622-1.45847963 3.20055254-3.20055253 3.20055254l-33.9501649 0c-1.78258622 0-3.20055254-1.45847963-3.20055255-3.20055254l1e-8-33.13989844z m0 129.03493466c0-1.78258622 1.45847963-3.20055254 3.20055255-3.20055253l33.9501649 0c1.78258622 0 3.20055254 1.45847963 3.20055254 3.20055254l-1e-8 33.13989846c0 1.78258622-1.45847963 3.20055254-3.20055253 3.20055253l-33.95016491-1e-8c-1.78258622 0-3.20055254-1.45847963-3.20055254-3.20055252l0-33.13989846zM570.37969887 577.7531237s7.73804476 30.99269231 8.83190446 43.83541579S566.1257999 639.21183513 538.33366014 641.64263451c-27.79213977 2.39028607-55.30068628 2.95747261-74.82810813-4.01081901s-18.9602353-10.41192409-17.0966224-18.71715533 9.35857768-41.16153644 9.35857768-41.16153647 76.69172098 13.65298994 114.61219159 0zM440.29141778 354.40317241c0-4.13235898 3.36260583-7.49496481 7.49496481-7.49496481S455.2813474 350.27081343 455.28134741 354.40317241l0 38.20406387-14.98992962 0L440.29141778 354.40317241z m1e-8 40.67537659l14.94941629 0L455.24083408 414.76802411 440.29141778 414.76802412l1e-8-19.68947513z m0 22.20130115l14.94941629-1e-8 0 19.60844848L440.29141778 436.88829862l1e-8-19.60844847z m0 22.07976121l14.94941629 0 0 18.6361287L440.29141778 457.99574006l1e-8-18.6361287z m0 21.10744143l14.94941629 0 0 19.16280192L440.29141778 479.6298547l1e-8-19.16280192z m0 21.63411463l14.94941629 0 0 19.28434188L440.29141778 501.3855093l1e-8-19.28434188z m0 21.7556546l14.94941629 0 0 4.98313877c0 4.13235898-3.36260583 7.49496481-7.49496481 7.4949648s-7.49496481-3.36260583-7.4949648-7.4949648l0-4.98313877z m0.81026646 28.76445953s2.95747261-8.66985118 8.02163802-6.8062383c0 0 1.98515284 35.40864456 0 44.60516892l-10.00679086 4.94262545 1.98515284-42.74155607z m-2.10669281 46.46878181l10.16884414-3.72722574s-2.55233937 23.25464757-5.22621869 30.18242585c-2.67387933 6.92777829-5.87443188-0.81026648-5.87443189-0.81026646l0.93180644-25.64493365z';
            var convertData = function (data) {	
                var res = [];
                for (var i = 0; i < data.length; i++) {
                    var dataItem = data[i];
                    var fromCoord = a[dataItem[0].name];
                    var toCoord = a[dataItem[1].name];
                    if (fromCoord && toCoord) {
                        res.push({
                            fromName: dataItem[0].name,
                            toName: dataItem[1].name,
                            coords: [fromCoord, toCoord],
                            value: dataItem[1].value
                        });
                    }
                }
                return res;	
            };
    	    var series = [];
            c.forEach(function (item, i) {  
                series.push({
                    name: item[0] + ' 货运',
                    type: 'lines',
                    zlevel: 1,
                    // effect: {
                    //     show: true,
                    //     period: 6,
                    //     trailLength: 0.7,
                    //     color: 'red',   //arrow箭头的颜色
                    //     symbolSize: 3
                    // },
                    lineStyle: {
                        normal: {
                            color: color[i],
                            width: 0,
                            curveness: 0.2
                        }
                    },
                    data: convertData(item[1])
                },

                // {
                //      name: '接通率',
                //      type: "map",
                //      map: "china",
                //      zoom: 1.152,
                //      // top: '30',
                //      zlevel: 1,
                //      left: 'center',
                //      align: 'right',
                //      label: {
                //          normal: {
                //              show: false
                //          },
                //          emphasis: {
                //              show: false
                //          },
                //      },
                //      itemStyle: {
                //          normal: {
                //              borderColor: "#112b3b", //省市边界线
                //              borderColor: "#a7e4e6", //省市边界线
                //              shadowColor: 'rgba(166, 230, 236, 0.6)',
                //              shadowOffsetX: 0,
                //              shadowOffsetY: 0,
                //              shadowBlur: 120
                //          },
                //          emphasis: {
                //              areaColor: "#4aafde",
                //              borderColor: "red"
                //          },
                //      },
                //      data:convertData(item[1])
                //      },

                {
                    name: item[0] + ' 货运',
                    type: 'lines',
                    zlevel: 2,
                    symbol: ['none', 'arrow'],
                    symbolSize: 10,
                    effect: {
                        show: true,
                        period: 6,
                        trailLength: 0,
                        symbol: planePath,
                        symbolSize: 15
                    },
                    lineStyle: {
                        normal: {
                            color: color[i],
                            width: 2,
                            opacity: 0.6,
                            curveness: 0.2
                        }
                    },
                    data: convertData(item[1])
                },
                {
                    name: item[0] + ' 货运',
                    type: 'effectScatter',
                    coordinateSystem: 'geo',
                    // zoom: 1.6,
                    zlevel: 2,
                    rippleEffect: {
                        brushType: 'stroke'
                    },
                    label: {
                        normal: {
                            show: true,
                            position: 'right',
                            formatter: '{b}'
                        }
                    },
                    symbolSize: function (val) {
                        return val[2] / 8;
                    },
                    itemStyle: {
                        normal: {
                            color: color[i],
                        },
                        emphasis: {
                            areaColor: '#2B91B7'
                        }
                    },
                    data: item[1].map(function (dataItem) {
                        return {
                            name: dataItem[1].name,
                            value: a[dataItem[1].name].concat(100)
                            // value: a[dataItem[1].name].concat([dataItem[1].value])
                        };
                    })
                });
            });
            var option = {
                backgroundColor: 'rgb(64,74,89)',
                title : {
                    text: '业务路线地图',
                    // subtext: '数据纯属虚构',
                    left: 'left',
                    top:'3%',
                    left:'2%',
                    textStyle : {
                        color: '#fff'
                    }
                },
                tooltip : {
                    trigger: 'item', 
                    formatter:function(params, ticket, callback){
                        if(params.seriesType=="effectScatter") {
                            return " 收发地："+params.data.name;
                        }else if(params.seriesType=="lines"){
                            return params.data.fromName+" -> "+params.data.toName+"<br/>" + '总趟数：' + params.data.value;
                        }else{
                            return params.name;
                        }
                    } 
                },
                // grid:{
                //     x:10,
                //     y:15,
                //     x2:10,
                //     y2:10,
                // },
                // legend: {
                //     orient: 'vertical',
                //     top: 'bottom',
                //     left: 'right',
                // // data:['大同 Top3', '西宁 Top3', '银川 Top3'],
                //     textStyle: {
                //         color: '#333'
                //     },
                //     selectedMode: 'multiple'
                // },

    //  visualMap: {
    //      //图例值控制
    //      show: true,
    //      type: 'piecewise',
    //      seriesIndex: 0,
    //      pieces: [{
    //              min: 401
    //          }, // 不指定 max，表示 max 为无限大（Infinity）。
    //          {
    //              min: 301,
    //              max: 400
    //          },
    //          {
    //              min: 201,
    //              max: 300
    //          },
    //          {
    //              min: 101,
    //              max: 200
    //          },
    //          {
    //              max: 100
    //          }
    //      ],

    //      showLabel: false,
    //      calculable: true,
    //  },

                geo: {
                    map: 'china',
                    label: {
                        emphasis: {
                            show: true,
                            color:'#fff'
                        }
                    },
                    roam: true,
        //                     label: {
        //     emphasis: {
        //         show: false
        //     }
        // },
        // layoutSize: "100%",
        //     normal: {
        //         borderColor: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{
        //             offset: 0,
        //             color: '#00F6FF'
        //         }, {
        //             offset: 1,
        //             color: '#53D9FF'
        //         }], false),
        //         borderWidth: 3,
        //         shadowColor: 'rgba(10,76,139,1)',
        //         shadowOffsetY: 0,
        //         shadowBlur: 60
        //     }

                    itemStyle: {
                        normal: {
                            areaColor: 'rgb(50,60,72)',
                            borderColor: '#fff',
                            borderWidth: 0.1,
                        },
                        emphasis: {
                            areaColor: 'rgb(166,200,76)'
                        }
                    }
                },
                series: series
            };
            option && chartsmyChart.setOption(option);
            });

        },

        /* 
        采购煤价趋势数据
        */
        settcitynews(type){
            this.loading = true;
            this.api.statistics.countPriceTrend(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.citynewdataArr = [];      // 时间
                    this.citynewseriesData = [];   // 采购
                    this.citynewriesData = [];     // 销售
                    let newdate = [];           // 组装采购数据
                    let newtwodate = [];        // 新采购数据
                    for(let key in res.data.data.cgprice){
                        let obj = {
                            date:key.replace(/-/g, "")*1,
                            sumPrice:res.data.data.cgprice[key]
                        }
                        newdate.push(obj);
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumPrice:newdate[i].sumPrice
                        }
                        newtwodate.push(obj);
                    }
                    for(let i=0;i<newtwodate.length;i++){
                        this.citynewdataArr.push(newtwodate[i].date);            // 时间
                        let num = newtwodate[i].sumPrice*1;
                        this.citynewseriesData.push(num.toFixed(2))                    // 采购数据
                    }
                    let xslbnewdate = [];                                       // 组装销售数据
                    let xslbnewtwodate = [];                                    // 新销售数据
                    for(let key in res.data.data.xxprice){
                        let obj = {
                            date:key.replace(/-/g, "")*1,
                            sumPrice:res.data.data.xxprice[key]
                        }
                        xslbnewdate.push(obj);
                    }
                    xslbnewdate.sort(this.compare('date'));
                    for(let i=0;i<xslbnewdate.length;i++){
                        let year = xslbnewdate[i].date.toString().substring(0,4);
                        let newdatas =  xslbnewdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumPrice:xslbnewdate[i].sumPrice
                        }
                        xslbnewtwodate.push(obj);
                    }
                    for(let i=0;i<xslbnewtwodate.length;i++){
                        let nums = xslbnewtwodate[i].sumPrice*1;
                        this.citynewriesData.push(nums.toFixed(2))                   // 销售数据
                    }
                    this.setcitynews();
                }
            })
        },
        // 采购煤价趋势图表
        setcitynews(){
            var settlementcity = document.getElementById('citynewsEchart');
            var myChartcitynewsEchart = echarts.init(settlementcity);
            var optioncitynewsEchart;
            optioncitynewsEchart = {
                title: {
                    text: '采购煤价趋势'
                },
                grid:{
                    x:50,
                    y:75,
                    x2:70,
                    y2:20,
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
                // toolbox: {
                //     iconStyle:{
                //         borderColor:'#CF1724'
                //     },
                //     feature: {
                //         magicType: {
                //             show: true, 
                //             type: ['line', 'bar'],
                //             emphasis: {iconStyle:{textFill:'#CF1724'}}
                //         },
                //     }
                // },
                legend: {
                    data: ['采购单价', '销售单价'],
                    textStyle: {
                        color: "#333"
                    }
                },
                xAxis: [
                    {
                        type: 'category',
                        name: '时间',
                        data: this.citynewdataArr,
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
                        // max: 6000,
                        // interval:1000,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333'
                            }
                        },
                    }
                ],
                series:[
                {
                    name: '采购单价',
                    type: 'line',
                    itemStyle: {
                        normal: {
                            color: '#72CE99'
                        }
                    },
                    data:this.citynewseriesData
                },
                {
                    name: '销售单价',
                    type: 'line',
                    itemStyle: {
                        normal: {
                            color:'#58AFFF'
                        }
                    },
                    data: this.citynewriesData,
                },
            ],
            };
            optioncitynewsEchart && myChartcitynewsEchart.setOption(optioncitynewsEchart);
        },

        /*
        筛选功能
        */
        // 打开弹窗
        dialogVisibleclick(){
            this.screenlist = this.newscreenlist;
            this.screeninvoice = this.newscreeninvoice;
            this.dialogVisibledata = true;
            this.topYear = [this.yearValues + '-01-01',this.yearValues + '-06-30'];
        },
        // 重置筛选
        dialogVisiblerice(){
            this.screenlist = [1,2,3,4,5,6];
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
            this.newscreenlist = this.screenlist;
            this.newscreeninvoice = this.screeninvoice;
            let invoiceDates = {
                startTime:this.newscreeninvoice[0],
                endTime:this.newscreeninvoice[1]
            } 
            for(let i=0;i<this.newscreenlist.length;i++){
                if(this.newscreenlist[i] == 1){
                    this.countBill(invoiceDates);       // 票据
                }else if(this.newscreenlist[i] == 2){
                    this.paymentCapital(invoiceDates);  // 收付
                }else if(this.newscreenlist[i] == 3){
                    this.burdenapi(invoiceDates);       // 税负率
                }else if(this.newscreenlist[i] == 4){
                    this.logis(invoiceDates);           // 物流
                }else if(this.newscreenlist[i] == 5){
                    this.settCapital(invoiceDates);     // 结算额
                }else if(this.newscreenlist[i] == 6){
                    this.countBillss(invoiceDates);     // 税额统计
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

        /*
        提示看板数据
        */
        tskbAll(){
            // 进项占比报警
            this.api.bill.selJxzb()
            .then(res=>{
                if(res.data.code == 200 && res.data.data != []){
                    for(let i=0;i<res.data.data.length;i++){
                        let str = res.data.data[i];
                        let att = str.substring(15,27);
                        let time = str.substring(0,15);
                        let last = str.substring(27,str.length);
                        let obj = {
                            time:time,
                            att:att,
                            last:last
                        }
                        this.newtipslist.push(obj)
                    }
                }
            });
            // 税负率报警
            this.api.bill.selZeroTax()
            .then(res=>{
                if(res.data.code == 200 && res.data.data != []){
                    for(let i=0;i<res.data.data.length;i++){
                    let str = res.data.data[i];
                    let att = str.substring(12,18);
                    let time = str.substring(0,12);
                    let last = str.substring(18,str.length);
                    let obj = {
                        time:time,
                        att:att,
                        last:last
                    }
                        this.tipslist.push(obj)
                    }
                }
            })
        },
        /* 
        采购 销售结算数据
        */
        countCapital(type){
            this.loading = true;
            this.api.capital.count({timeType:type})
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.arrsumNum = '';
                    this.sumNum = '';
                    for(let i=0;i<res.data.data.length;i++){
                        if(res.data.data[i].type == 1){
                            this.arrsumNum = res.data.data[i].sumNum;   // 采购结算吨数
                        }
                        if(res.data.data[i].type == 2){
                            this.sumNum = res.data.data[i].sumNum;      // 销售结算吨数
                        }
                    }
                }
            })
        },
        // 采购合同数据
        countlistByPage(type){
            this.loading = true;
            this.api.listByPage.count({timeType:type})
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.chasingnum = '';
                    this.chasingTonnage = '';
                    this.salenum = '';
                    this.saleTonnage = '';
                    for(let i=0;i<res.data.data.fslist.length;i++){
                        if(res.data.data.fslist[i].contract_type == 1){
                            this.chasingnum = res.data.data.fslist[i].contractNum;      // 采购合同份数
                        }
                        if(res.data.data.fslist[i].contract_type == 2){
                            this.salenum = res.data.data.fslist[i].contractNum;         // 销售合同份数
                        }
                    }
                    for(let i=0;i<res.data.data.htllist.length;i++){
                        if(res.data.data.htllist[i].contract_type == 1){
                            this.chasingTonnage = res.data.data.htllist[i].sumNum;      // 采购合同吨数
                        }
                        if(res.data.data.htllist[i].contract_type == 2){
                            this.saleTonnage = res.data.data.htllist[i].sumNum;         // 销售合同吨数
                        }
                    }
                }
            })
        },
        // 采购下拉菜单
        settlemenchange(name,type){
            this.settlementtime = name;
            this.countCapital(type);
            this.countlistByPage(type);
        },
        // 销售下拉菜单
        salechange(name){
            this.saletime = name;
        },

        /* 
        税负率数据
        */
        burdenapi(type){
            this.loading = true;
            // let invoiceDatess;
            // if(type.length == 2){
            //     invoiceDatess = {
            //         startTime:type[0],
            //         endTime:type[1]
            //     };
            // }else{
            //     invoiceDatess = {
            //         startTime:this.invoice[0],
            //         endTime:this.invoice[1]
            //     };
            // }
            this.api.bill.countTaxRate(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.burdenData = [];
                    this.burdenTime = [];
                    let data = res.data.data;
                    var newatt = [];
                    for(let key in data){
                        if(key != 'ljsfl'){
                            let att = (key).replace(/-/g, "");
                            newatt.push(att*1);
                            newatt.sort(function(a, b){return a - b});
                        }else{
                            this.burdenCumulative = (data[key]*100).toFixed(2);
                        }
                    }
                    for(let i=0;i<newatt.length;i++){
                        let year = newatt[i].toString().substring(0,4);
                        let newdata =  newatt[i].toString().substring(4,6);
                        var newtime = year + '-' + newdata;
                        this.burdenTime.push(newtime);
                        let color;
                        if(data[newtime]*100 > 10){
                            color = "#EE6666";
                            this.credit = '良好';
                        }else{
                            color = "#FCDEDE";
                        }
                        let obj = {
                            value:(data[newtime]*100).toFixed(2),
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
            var burden = document.getElementById('burden');
            var myChartburden = echarts.init(burden);
            var optionmyChartburden;
            optionmyChartburden = {
                title: {
                    text: '税负率'
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
                        name: '时间',
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
                        // // height: 15,                              // 组件高度
                        // handleColor: '#CF1724',                     // h滑动图标的颜色
                        // fillerColor:"#F5C4C4",                      // 选中范围的填充颜色。
                        // borderColor:"#CF1724",                      // 边框颜色。
                        show: true,                                 // flase直接隐藏图形
                        xAxisIndex: [0],
                        left: '5%',                                 // 滚动条靠左侧的百分比
                        bottom: -5,
                        start: 0,                                   // 滚动条的起始位置
                        end: 70,                                    // 滚动条的截止位置（按比例分割你的柱状图x轴长度）

                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '百分比',
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
                        barWidth:'35%',
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
                                formatter: "警戒线"
                            },
                            data : [{
                                silent:true,             //鼠标悬停事件  true没有，false有
                                lineStyle:{               //警戒线的样式  ，虚实  颜色
                                    type:"dashed",
                                    color:"#CF1724"
                                },
                                name: '警戒线',
                                yAxis: 10
                            }]
                        },
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
                    },
                ],
            };
            optionmyChartburden && myChartburden.setOption(optionmyChartburden);
        },
        // burdenchange(time){
        //     this.burdenapi(time);
        // },

        /* 
        税额统计数据
        */
        countBillss(type){
            this.loading = true;
            this.api.bill.count(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.sumIncludMoney = '';
                    this.arrIncludMoney = '';
                    this.sumTaxAmount = 0;
                    this.arrTaxAmount = 0;
                    for(let i=0;i<res.data.data.length;i++){
                        if(res.data.data[i].type == 1){
                            this.sumIncludMoney = (res.data.data[i].sumIncludMoney) *1/10000;   // 进项含税金额  
                            this.sumTaxAmount =  res.data.data[i].sumTaxAmount *1/10000;        // 进项税额
                        }else if(res.data.data[i].type == 2){
                            this.arrIncludMoney = (res.data.data[i].sumIncludMoney) *1/10000;   // 销项含税金额
                            this.arrTaxAmount = res.data.data[i].sumTaxAmount *1/10000;         // 销项税额 
                        }
                    }
                }
            });
            type.companyId = sessionStorage.getItem('companyId');
            this.api.billRate.rateBill(type)
            .then(res=>{
                if(res.data.code == 200){
                    this.rateBill = (res.data.data.sumMoney*1/10000).toFixed(2)*1;
                }
            })
        },
        // dropdownchange(name,type){
        //     this.quantime = name;
        //     this.countBillss(type);
        // },

        /* 
        票据税额数据
        */
        countBill(type){
            this.loading = true;
            // let invoiceDatess;
            // if(type.length == 2){
            //     invoiceDatess = {
            //         startTime:type[0],
            //         endTime:type[1]
            //     };
            // }else{
            //     invoiceDatess = {
            //         startTime:this.invoice[0],
            //         endTime:this.invoice[1]
            //     };
            // }
            this.api.bill.count(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    if(res.data.data == '' || res.data.data.length <= 1){
                        this.sumTaxAmount = 0;
                        this.arrTaxAmount = 0;
                        this.chartDom();  // 税额
                    }else{
                        for(let i=0;i<res.data.data.length;i++){
                            if(res.data.data[i].type == 1){
                                this.sumTaxAmount =  res.data.data[i].sumTaxAmount *1/10000; // 进项税额
                                // this.sumTaxAmount = numdata.toFixed(2);
                            }else if(res.data.data[i].type == 2){
                                this.arrTaxAmount = res.data.data[i].sumTaxAmount *1/10000; // 销项税额 
                                // this.arrTaxAmount = newdata.toFixed(2)
                            }
                        }
                        this.chartDom();  // 税额
                    }
                }
            })
        },
        // invoicechange(val){
        //     this.countBill(val);
        // },
        chartDom(){
            var chartDom = document.getElementById('invoice');
            var myChart = echarts.init(chartDom);
            var option;
            var num;
            if(this.sumTaxAmount*1 == 0 || this.arrTaxAmount*1 == 0){
                num = 0;
            }else{
                let newnum = this.sumTaxAmount*1 / this.arrTaxAmount*1;
                num = (newnum*100).toFixed(2);
            }
            option = {
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
            // option = {
            //     title: {
            //         text: num,
            //         subtext: '进项占比',
            //         x: 'center',
            //         y: 'center',
            //         textStyle: {
            //             fontSize: 18,
            //             color: '#333'
            //         }, // 标题
            //         subtextStyle: {
            //             fontWeight:'600',
            //         fontSize: 14,
            //             color: '#333'
            //         }, // 副标题
            //     },
            //     color: ['#FAD337', '#FB7F5B'],
            //     series: [
            //     { 
            //         type: 'pie',
            //         radius: ['40%', '60%'], // 半径
            //         center: ['50%', '50%'], // 位置
            //         label: {
            //             normal: {
            //                 formatter: ' {b} \n {c} ',
            //             },
            //         },
            //         data: [
            //             {
            //                 value:mount,
            //                 name: '进项税额万(元)',
            //             },
            //             {
            //                 value:sumNumss,
            //                 name: '销项税额万(元)'
            //             },
            //         ],
            //     }]
            // };
            option && myChart.setOption(option);
        },

        /* 
        结算额数据
        */
        settCapital(type){
            this.loading = true;
            this.api.capital.setall(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.serienum = 0;
                    this.riesnum = 0;
                    this.dataArr = [];      // 时间
                    this.seriesData = [];   // 采购
                    this.riesData = [];     // 销售
                    let newdate = [];           // 组装采购数据
                    let newtwodate = [];        // 新采购数据
                    for(let i=0;i<res.data.data.cglb.length;i++){
                        let obj = {
                            date:(res.data.data.cglb[i].date).replace(/-/g, "")*1,
                            sumPrice:res.data.data.cglb[i].sumPrice
                        }
                        newdate.push(obj);
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumPrice:newdate[i].sumPrice
                        }
                        newtwodate.push(obj);
                    }
                    for(let i=0;i<newtwodate.length;i++){
                        this.dataArr.push(newtwodate[i].date);            // 时间
                        let num = newtwodate[i].sumPrice*1/10000;
                        this.seriesData.push(num.toFixed(2))                    // 采购数据
                        this.serienum += newtwodate[i].sumPrice *1/10000;
                    }
                    let xslbnewdate = [];                                       // 组装销售数据
                    let xslbnewtwodate = [];                                    // 新销售数据
                    for(let i=0;i<res.data.data.xslb.length;i++){
                        let obj = {
                            date:(res.data.data.xslb[i].date).replace(/-/g, "")*1,
                            sumPrice:res.data.data.xslb[i].sumPrice
                        }
                        xslbnewdate.push(obj);
                    }
                    xslbnewdate.sort(this.compare('date'));
                    for(let i=0;i<xslbnewdate.length;i++){
                        let year = xslbnewdate[i].date.toString().substring(0,4);
                        let newdatas =  xslbnewdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumPrice:xslbnewdate[i].sumPrice
                        }
                        xslbnewtwodate.push(obj);
                    }
                    for(let i=0;i<xslbnewtwodate.length;i++){
                        let nums = xslbnewtwodate[i].sumPrice*1/10000;
                        this.riesData.push(nums.toFixed(2))                   // 销售数据
                        this.riesnum += xslbnewtwodate[i].sumPrice*1/10000;
                    }
                    this.sett();
                }
            })
        },
        // 结算图表
        sett(){
            var settlement = document.getElementById('settlement');
            var myChartsettlement = echarts.init(settlement);
            var optionmyChartsettlement;
            optionmyChartsettlement = {
                title: {
                    text: '结算额'
                },
                grid:{
                    x:50,
                    y:75,
                    x2:70,
                    y2:20,
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
                // toolbox: {
                //     iconStyle:{
                //         borderColor:'#CF1724'
                //     },
                //     feature: {
                //         magicType: {
                //             show: true, 
                //             type: ['line', 'bar'],
                //             emphasis: {iconStyle:{textFill:'#CF1724'}}
                //         },
                //     }
                // },
                legend: {
                    data: ['采购结算额', '销售结算额'],
                    textStyle: {
                        color: "#333"
                    }
                },
                xAxis: [
                    {
                        type: 'category',
                        name: '时间',
                        data: this.dataArr,
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
                        max: 6000,
                        interval:1000,
                        axisLine: {
                            show: true,
                            lineStyle: {
                                color: '#333'
                            }
                        },
                    }
                ],
                series:[
                {
                    name: '采购结算额',
                    type: 'bar',
                    itemStyle: {
                        normal: {
                            color: '#72CE99'
                        }
                    },
                    data:this.seriesData
                },
                {
                    name: '销售结算额',
                    type: 'bar',
                    itemStyle: {
                        normal: {
                            color:'#58AFFF'
                        }
                    },
                    data: this.riesData,
                },
            ],
            };
            optionmyChartsettlement && myChartsettlement.setOption(optionmyChartsettlement);
        },
        // settChange(name,type){
        //     this.setttime = name;
        //     this.settCapital(type)
        // },

        /* 
        收付数据
        */
        // 处理时间排序函数
        compare(property){
            return function(a,b){
                var value1 = a[property];
                var value2 = b[property];
                return value1 - value2;
            }
        },
        paymentCapital(type){
            this.loading = true;
            // let invoiceDatess;
            // if(type.length == 2){
            //     invoiceDatess = {
            //         startTime:type[0],
            //         endTime:type[1]
            //     };
            // }else{
            //     invoiceDatess = {
            //         startTime:this.invoice[0],
            //         endTime:this.invoice[1]
            //     };
            // }
            this.api.chainCap.setall(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.datapayment = [];      // 时间
                    this.paymentData = [];      // 收款数据
                    this.paymentserData = [];   // 付款数据
                    let newdate = [];           // 组装付款数据
                    let newtwodate = [];        // 新付款数据
                    for(let i=0;i<res.data.data.fklb.length;i++){
                        let obj = {
                            date:(res.data.data.fklb[i].date).replace(/-/g, "")*1,
                            sumPrice:res.data.data.fklb[i].sumPrice
                        }
                        newdate.push(obj);
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        let length = newdate[i].date.toString().length
                        if(length > 6){
                            let daty = newdate[i].date.toString().substring(6,8);
                            var newtime = year + '-' + newdatas + '-' + daty;
                        }else{
                            var newtime = year + '-' + newdatas;
                        }
                        let obj = {
                            date:newtime,
                            sumPrice:newdate[i].sumPrice
                        }
                        newtwodate.push(obj);
                    }
                    for(let i=0;i<newtwodate.length;i++){
                        this.datapayment.push(newtwodate[i].date);          // 时间
                        let num = newtwodate[i].sumPrice*1/10000;
                        this.paymentserData.push(num.toFixed(2))            // 付款数据
                    }
                    let sklbnewdate = [];                                   // 组装收款数据
                    let sklbnewtwodate = [];                                // 新收款数据
                    for(let i=0;i<res.data.data.sklb.length;i++){
                        let obj = {
                            date:(res.data.data.sklb[i].date).replace(/-/g, "")*1,
                            sumPrice:res.data.data.sklb[i].sumPrice
                        }
                        sklbnewdate.push(obj);
                    }
                    sklbnewdate.sort(this.compare('date'));
                    for(let i=0;i<sklbnewdate.length;i++){
                        let year = sklbnewdate[i].date.toString().substring(0,4);
                        let newdatas =  sklbnewdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumPrice:sklbnewdate[i].sumPrice
                        }
                        sklbnewtwodate.push(obj);
                    }
                    for(let i=0;i<sklbnewtwodate.length;i++){
                        let num = sklbnewtwodate[i].sumPrice*1/10000;
                        this.paymentData.push(num.toFixed(2))                   // 收款数据
                    }
                    this.payment();
                }
            })
        },
        payment(){
            // 收付图表
            var chartDompayment = document.getElementById('payment');
            var myChartpayment = echarts.init(chartDompayment);
            var optionpayment;
            optionpayment = {
                title: {
                    text: '收付'
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
                    data: ['收款金额', '付款金额'],
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
                        max: 6000,
                        interval: 1000,
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
                        name: '收款金额',
                        type: 'bar',
                        itemStyle: {
                            normal: {
                                color: '#FAD790'
                            }
                        },
                        data:this.paymentData
                    },
                    {
                        name: '付款金额',
                        type: 'bar',
                        itemStyle: {
                            normal: {
                                color: '#EE6666'
                            }
                        },
                        data:this.paymentserData
                    }
                ]
            };
            optionpayment && myChartpayment.setOption(optionpayment);
        },
        // paymentChange(name,type){
        //     this.paymenttime = name;
        //     this.paymentCapital(type)
        // },
        
        /* 
        物流数据
        */
        logis(type){
            this.loading = true;
            // let invoiceDatess;
            // if(type.length == 2){
            //     invoiceDatess = {
            //         startTime:type[0],
            //         endTime:type[1]
            //     };
            // }else{
            //     invoiceDatess = {
            //         startTime:this.invoice[0],
            //         endTime:this.invoice[1]
            //     };
            // }
            this.api.logistics.setall(type)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                    this.Tonum = 0;
                    this.Selfnum = 0;
                    this.datalogistics = [];    // 时间
                    this.logisticsData = [];    // 自提
                    this.mountsData = [];       // 送到
                    let newdate = [];           // 组装自提数据
                    let newtwodate = [];        // 新自提数据
                    for(let i=0;i<res.data.data.ztlb.length;i++){
                        let obj = {
                            date:(res.data.data.ztlb[i].date).replace(/-/g, "")*1,
                            sumTonnage:res.data.data.ztlb[i].sumTonnage
                        }
                        newdate.push(obj);
                    }
                    newdate.sort(this.compare('date'));
                    for(let i=0;i<newdate.length;i++){
                        let year = newdate[i].date.toString().substring(0,4);
                        let newdatas =  newdate[i].date.toString().substring(4,6);
                        let length = newdate[i].date.toString().length
                        if(length > 6){
                            let daty = newdate[i].date.toString().substring(6,8);
                            var newtime = year + '-' + newdatas + '-' + daty;
                        }else{
                            var newtime = year + '-' + newdatas;
                        }
                        let obj = {
                            date:newtime,
                            sumTonnage:newdate[i].sumTonnage
                        }
                        newtwodate.push(obj);
                    }
                    for(let i=0;i<newtwodate.length;i++){
                        this.datalogistics.push(newtwodate[i].date);            // 时间
                        let num = newtwodate[i].sumTonnage*1;
                        this.logisticsData.push(num.toFixed(2))                 // 自提数据
                        this.Selfnum += newtwodate[i].sumTonnage *1;            // 自提总吨数
                        this.Selfnum = this.Selfnum.toFixed(2)*1;
                    }
                    let sdlbnewdate = [];                                       // 组装送到数据
                    let sdlbnewtwodate = [];                                    // 新送到数据
                    for(let i=0;i<res.data.data.sdlb.length;i++){
                        let obj = {
                            date:(res.data.data.sdlb[i].date).replace(/-/g, "")*1,
                            sumTonnage:res.data.data.sdlb[i].sumTonnage
                        }
                        sdlbnewdate.push(obj);
                    }
                    sdlbnewdate.sort(this.compare('date'));
                    for(let i=0;i<sdlbnewdate.length;i++){
                        let year = sdlbnewdate[i].date.toString().substring(0,4);
                        let newdatas =  sdlbnewdate[i].date.toString().substring(4,6);
                        var newtime = year + '-' + newdatas;
                        let obj = {
                            date:newtime,
                            sumTonnage:sdlbnewdate[i].sumTonnage
                        }
                        sdlbnewtwodate.push(obj);
                    }
                    for(let i=0;i<sdlbnewtwodate.length;i++){
                        let nums = sdlbnewtwodate[i].sumTonnage*1;
                        this.mountsData.push(nums.toFixed(2))                   // 送到数据
                        this.Tonum += nums;      // 送到总吨数
                        this.Tonum = this.Tonum.toFixed(2)*1;
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
            var Extraction = '自提 '+ this.Selfnum + ' 吨';     // 自提
            var goTo = '送到 ' + this.Tonum + ' 吨';         // 送到
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
                legend: {
                    data: [Extraction,goTo],
                    textStyle: {
                        color: "#333"
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
                        name: '吨',
                        // min: 0,
                        // max: 10000,
                        // interval: 2000,
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
                        name:Extraction,
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
                    },
                    {
                        name: goTo,
                        type: 'line',
                        smooth: true,
                        symbol: 'none',
                        itemStyle: {
                            normal: {
                                color: '#5887FF'
                            }
                        },
                        areaStyle: {
                            color: '#EAEBFF'
                        },
                        data:this.mountsData
                    }
                ]
            };
            optionlogistics && myChartlogistics.setOption(optionlogistics);
        },
        // logisticsChange(name,type){
        //     this.logisticstime = name;
        //     this.logis(type)
        // },

        /* 
        一键导入功能
        */
        // 取消导入
        Importble(){
            this.ImportVisible = false;
            this.fileForm.importnav = "";
            this.$refs.upload.clearFiles();
        },
        // 导入文件格式验证
        supbeforedemoUpload(file) {
            this.header.Authorization = sessionStorage.getItem('Authorization')
            this.fileForm.importnav = file.name
            let str = file.name;
            let index = str.lastIndexOf("\.");
            str = str.substring(index + 1,str.length)
            let types = [
            "xlsx",
            "xls",
            ];
            const isImg = types.indexOf(str) == -1;
            if (isImg) {
                this.$message.error("上传文件只能是 xlsx xls 格式,请重新上传!");
            }
            return !isImg;
        },
        // 导入功能方法
        Importdailog(){
            this.ImportVisible = true;
        },
        // 关闭导入弹窗
        ImporthandleClose() {
            this.ImportVisible = false;
            this.perstep = 0;
            this.active = 1;
            this.oneShow = true;
            this.twoShow = false;
            this.streeShow = false;
            this.pernav = "正在导入...";
            this.fileForm.importnav = "";
            this.$refs.upload.clearFiles();
        },
        // 第二步 第三步 导入结果
        ImportNext(formName){
            this.$refs[formName].validate((valid) => {
            if (valid) {
                this.perstep = 0
                this.active = 2;
                this.oneShow = false;
                this.twoShow = true;
                // 导入进度条
                let f = setInterval(()=>{
                    this.perstep += 3;
                    if(this.perstep == 99){
                    this.perstep = 99;
                    this.$refs.upload.submit();
                    clearInterval(f);
                    }
                },150) 
            }
            });
        },
        // 导入成功的钩子
        addFileLet(res){
            this.perstep = 0
            this.pernav = "导入完成";
            this.active = 3;
            this.oneShow = false;
            this.twoShow = false;
            if(res.code == 200){
                this.fourShow = true         // 成功
                this.streeShow = false;

                this.sucessShow = false;
                this.errorShow = true;
                this.sucessData = res.data.sucInfos;
                this.errorData = res.data.errs;
                this.cation = res.data.errInfos;
            }else{
                this.fourShow = false;
                this.streeShow = true;       // 失败
                this.cation = res.data;
            }
        },
        // 导入完成
        ImportNextstree(){
            this.ImportVisible = false;
            this.perstep = 0;
            this.active = 1;
            this.oneShow = true;
            this.twoShow = false;
            this.importnav = "";
            this.streeShow = false;
            this.fourShow = false;         // 成功
            this.pernav = "正在导入...";
            this.streeShow = false; 
            let invoiceDate = {
                startTime:this.screeninvoice[0],
                endTime:this.screeninvoice[1]
            };
            // 提示看板
            this.tskbAll();
            // 票据税额接口
            this.countBill(invoiceDate);
            // 税负率
            this.burdenapi(invoiceDate);
            // 收付
            this.paymentCapital(invoiceDate);
            // 物流
            this.logis(invoiceDate);
            // 结算额
            this.settCapital(invoiceDate);
            // 税额统计接口
            this.countBillss(invoiceDate);
            // 采购合同数据
            this.countCapital(5);
            // 采购结算量数据
            this.countlistByPage(5);
            this.current = 1;
            this.$refs.upload.clearFiles();
            this.fileForm.importnav = "";
        },

        // 一键删除导入文件
        numberClearFile(){
            if(this.importclearnumber.length == 36){
                this.$confirm('此操作将永久删除批次号对应的一键导入文件, 是否继续?', '提示', {
                  confirmButtonText: '确定',
                  cancelButtonText: '取消',
                  type: 'warning'
                }).then(() => {
                    this.api.listByPage.delOnce({im_batch_number:this.importclearnumber})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: 'success',
                                message: '删除成功!'
                            });
                        }
                    })
                })
            }else{
                this.$message({
                  type: 'warning',
                  message: '文件的批次号不能为空或批次号不正确!'
                });
            }
        }
    },
};
</script>
<style lang="less" scoped>
// 导入样式开始
.Timedialog{
    font-size: 18px;
    padding: 20px 20px 10px;
}
.Timedialog-box{
    padding: 30px 20px;
}
.importone{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  .importone-bottom{
    float: right;
    margin-right: -75px;
    margin-top: 52px;
  }
  .downBox{
    // display: inline-block;
    width: 100%;
    text-align: right;
    padding-right:5.5%;
    box-sizing: border-box;
  }
  a{

    font-size: 14px;
    color:#CF1724;
    text-decoration:none;
  }
  .importclear{
      width: 100%;
      padding-top:28px;
      box-sizing: border-box;
  }
  .importone-nav{
    padding-top:54px;
    box-sizing: border-box;
    span{
      font-size: 16px;
      line-height: 38px;
      color: #666666;
    }
    p{
      font-size: 14px;
      line-height: 30px;
      color: #666;
    }
  }
}
.importtwo{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  .importtwo-Box{
    width: 130px;
    margin: auto;
    text-align: center;
    span{
      line-height: 50px;
      color:#CF1724;
    }
  }
}
.importstree{
  width: 100%;
  height: 100%;
  overflow: hidden;
  padding:50px 75px 0;
  box-sizing: border-box;
  i{
    font-size: 18px;
    color: #CF1724;
  }
  span{
    font-size: 16px;
    color: #999;
    margin-left:5px;
    line-height: 32px;
  }
  p{
    font-size: 14px;
    color: #666;
    margin-left:30px;
  }
  .stree-buton{
    margin-left:24px;
  }
  .stree-ok{
    margin-top: 109px;
    .el-button{
      float: right;
      margin-right: -75px;
    }
  }
}

.importfour{
    width: 100%;
    height: 100%;
    overflow: hidden;
    padding:50px 75px 0;
    box-sizing: border-box;
    // 导入成功 失败返回数据开始样式
    .importfour-Box{
        width: auto;
        margin: auto;
        text-align: center;
        span{
        display: block;
        line-height: 50px;
        color:#999;
        }
        .err-box{
            width: 100%;
            .el-alert{
                float: left;
                width: 20%;
                margin: 1% 1% 1% 0;   
            }
        }
    }
    // 导入成功 失败返回数据结束样式
    .error-size{
        width: 100%;
        height: 200px;
        border: 1px solid #E4E4E4;
        border-radius: 4px;
        padding: 8px;
        box-sizing: border-box;
        margin-top:15px;
        overflow-y: scroll;
        text-align: left ;
        p{
            font-size: 14px;
            line-height: 24px;
            color: #999;
        }
    }
    .four-ok{
        margin-top: 30px;
        .el-button{
            float: right;
            margin-right: -75px;
        }
    }
}
// 导入样式结束
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
    // 物流地图
    .cityMap{
        width: 40%;
        height: 450px;
        margin:0 0.5% 0.5% 0.5%;
        float: left;
        background: #fff;
        #maps{
            width: 100%;
            height: 100%;
        }
    }
    .citynews{
        width: 58.5%;
        height: 450px;
        margin:0 0.5% 0.5% 0;
        float: left;
        background: #fff;
        padding: 16px 0 0 29px;
        box-sizing: border-box;
        .citynewsTitle{
            font-size: 18px;
            font-weight: 600;
            color: #333;
            
        }
        .citynewsImage{
            width: 100%;
            height: 94px;
            margin:2% 0;
            .img{
                width: 16%;
                height: 100%;
                float: left;
            }
            .citynewsBox{
                width: 80%;
                float: left;
                display: flex;
                flex-wrap: wrap;
                padding-left:3%;
                box-sizing: border-box;
                li{
                    width: 50%;
                    height: 33%;
                    line-height: 32px;
                    b{
                        font-size: 16px;
                        color: #333;
                    }
                    span{
                        font-size: 16px;
                        color: #333;
                    }
                }
            }
        }
        .citynewsEchartBox{
            width: 100%;
            height: 260px;
            // background: burlywood;
            #citynewsEchart{
                width: 100%;
                height: 100%;
            }
        }
    }
    // 采购 销售数据
    .Purchasing{
        float: left;
        width:34.5%;
        height:241px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;
        margin:0 0 0.5% 0.5%;
        padding: 16px 16px 16px 20px;
        box-sizing: border-box;
        .tity{
            font-size: 18px;
            font-weight: 600;
            color: #333;
            line-height: 26px;
        }
        span{
            line-height: 26px;
        }
        .Purchasing-time{
            width:auto;
            height: 26px;
            background: #F0F0F0;
            border-radius: 60px;
            float: right;
            outline: none;
            padding: 0 10px 2px;
            box-sizing: border-box;
            cursor: pointer;
            color: #333;
        }
        .Purchasing-conent{
            // width: 100%;
            // height: 85%;
            li{
                width: 50%;
                float: left;
                padding-top:30px;
                box-sizing: border-box;
                .styletop{
                    padding-top:30px;
                    h5{
                        font-size: 24px;
                        color: #333;
                        line-height: 56px; 
                    }
                    p{
                        font-size: 14px;
                        color: #828282;
                    }
                }
                img{
                    width: 32px;
                    height: 32px;
                    vertical-align: middle;
                }
                span{
                    font-size: 14px;
                    color: #828282;
                    vertical-align: middle;
                    margin-left: 8px;
                }
                h4{
                    font-size: 24px;
                    color: #333;
                    line-height: 32px;
                }
            }
        }
    }
    // 提示看板
    .Tips{
        float: left;
        width:29%;
        height:241px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;
        margin:0 0 0.5% 0.5%;
        padding: 16px 0 16px 20px;
        box-sizing: border-box;
        span{
            font-size: 18px;
            font-weight: 600;
            color: #333;
            line-height: 26px;
            vertical-align:middle;
        }
        img{
            width: 18px;
            height: 18px;
            margin-left: 8px;
            vertical-align:middle;
        }
        .tips-list{
            width: 100%;
            height: 165px;
            box-sizing: border-box;
            margin-top: 14px;
            li{
                padding:0 30px 0 10px;
                div{
                    .dian{
                        width: 6px;
                        height: 6px;
                        border-radius: 50%;
                        background-color:#CF1724;
                        float: left;
                        margin:11px 8px 0 0;
                    }
                    span{
                        width: 70%;
                        font-size: 14px;
                        line-height: 28px;
                        word-break:break-all;
                        color: #333;
                    }
                }
            }
        }
        .tips-box{
            width: 100%;
            height:165px;
            margin-top: 10px;
            text-align: center;
            padding-top:50px;
            box-sizing: border-box  ;
            p{
                font-size: 14px;
                color: #828282;
            }
        }
    }
    // 筛选器
    .screenBox{
        width: 99%;
        height:50px;
        float: left;
        background-color: #fff;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        margin:0 0 0.5% 0.5%;
        padding: 10px 40px 10px 20px;
        box-sizing: border-box;
        b{
            font-size: 18px;
            font-weight:800;
            line-height: 30px;
            color: #333;
        }
        .el-button{
            float: right;
        }
    }
    // 税额统计
    .quantity{
        float: left;
        width:99%;
        height:168px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;
        margin:0 0 0.5% 0.5%;
        padding: 16px 16px 16px 20px;
        box-sizing: border-box;
        .tity{
            font-size: 18px;
            font-weight: 600;
            color: #333;
            line-height: 26px;
        }
        span{
            line-height: 26px;
        }
        .quantity-time{
            width:auto;
            height: 26px;
            background: #F0F0F0;
            border-radius: 60px;
            float: right;
            outline: none;
            padding: 0 10px 2px;
            box-sizing: border-box;
            cursor: pointer;
            color: #333;
        }
        .quantity-conent{
            width: 100%;
            li{
                width: 20%;
                height: 100%;
                float: left;
                margin-top:30px;
                .dian{
                    width: 6px;
                    height: 6px;
                    margin: 25px 16px 25px 0;
                    float: left;
                    border-radius: 50%;
                }
                div{
                    h5{
                        font-size: 24px;
                        color: #333;
                        line-height: 35px
                    }
                    span{
                        font-size:14px;
                        color: #828282;
                    }
                }
                i{
                    font-size: 24px;
                    float: right;
                    margin-top:35px;
                    color: #828282;
                }
            }
        }
    }
    // 结算
    .settlement{
        float: left;
        width:99%;
        height:328px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff; 
        margin: 0 0 0 0.5%;
        #settlement{
            float: left;
            width: 80%;
            height: 100%;
            padding: 14px;
            box-sizing: border-box;
        }
        .settlement-right{
            float: left;
            width: 20%;
            height: 100%;
            padding: 16px;
            box-sizing: border-box;
            .quantity-time{
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
            }
            .dian{
                width: 6px;
                height: 6px;
                border-radius: 50%;
                float: left;
                margin:34px 8px 0 0;
            }
            .dian-size{
                float: left;
                font-size: 14px;
                color:rgba(0, 0, 0, 0.85);
                margin-top: 28px;
            }
            h5{
                font-size: 24px;
                line-height: 70px;
                color: #333;
            }

        }
    }
    // 收付
    .payment{
        float: left;
        width:99%;
        height:372px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;   
        margin:0.5%;
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
        width:99%;
        height:345px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;  
        margin:0 0.5% 0.3%;
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
    // 税负率
    .burden{
        float: left;
        width:69.2%;
        height:400px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
        background-color: #fff;  
        margin:0 0.5%;
        margin-right: 0;
        padding: 14px;
        box-sizing: border-box;
        position: relative;
        #burden{
            width:100%;
            height:100%;
        }
        .burden-time{
            position: absolute;
            right:14px;
            top: 14px;
            text-align: right;
            span{
                font-size: 14px;
                line-height: 32px;
                color: #333;
                margin-right: 17px;
            }
            
        }
    }
    // 发票
    .invoice{
        float: left;
        width:29.2%;
        height:400px;
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
            height: calc(160% - 38px);
            // padding-top:28px;
            // box-sizing: border-box;
            
        }
    }
}
// 内容结束
</style>