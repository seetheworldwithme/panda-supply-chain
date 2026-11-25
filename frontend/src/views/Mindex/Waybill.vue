<template>
  <!-- 运单管理页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col
      >
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav"
      v-loading="loading"
      element-loading-text = "加载中...">
      <!-- 表格 -->
      <div class="pur-table">
        <el-table
          height="100%"
          border
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%;"
          @selection-change="purChange">
            <el-table-column type="selection"  fixed="left"></el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column v-if="$route.meta.edit"
              label="操作" width="80" fixed="left">
              <template slot-scope="scope">
                <el-button size="mini" @click="EditPurs(scope.row)">编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column label="系统编码" sortable :show-overflow-tooltip="true" width="260">
                <template slot-scope="scope">
                    <span style="color: #cd1724;cursor: pointer;"
                        @click="purSee(scope.$index, scope.row)"
                        >{{ scope.row.systemCode }}</span>
                </template>
            </el-table-column>
            <el-table-column v-if="$route.meta.edit" label="关联详情">
              <template slot-scope="scope">
                <span v-if="scope.row.carriageName == 1" style="color: #cd1724;cursor: pointer;" @click="pronPurs(scope.row.id)">磅单详情</span>
                <span v-if="scope.row.carriageName == 2" style="color: #cd1724;cursor: pointer;" @click="OnePurs(scope.row)">运单详情</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="订单编号" sortable :show-overflow-tooltip="true" width="260">
                <template slot-scope="scope">
                    <span v-if="scope.row.billNumber"
                        style="color: #cd1724;cursor: pointer;"
                        @click="purSee(scope.$index, scope.row)"
                        >{{ scope.row.billNumber }}</span>
                        <span v-else>---</span>
                </template>
            </el-table-column> -->
            <el-table-column label="订单昵称" :show-overflow-tooltip="true" >
              <template slot-scope="scope">
                <span v-if="scope.row.orderName">{{scope.row.orderName}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="订单编号" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.orderNumber">{{scope.row.orderNumber}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="订单创建时间" :show-overflow-tooltip="true" width="130">
              <template slot-scope="scope">
                <span v-if="scope.row.orderCreateTime">{{scope.row.orderCreateTime}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="订单开始时间" :show-overflow-tooltip="true" width="130">
              <template slot-scope="scope">
                <span v-if="scope.row.orderStartTime">{{scope.row.orderStartTime}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="订单结束时间" :show-overflow-tooltip="true" width="130">
              <template slot-scope="scope">
                <span v-if="scope.row.orderEndTime">{{scope.row.orderEndTime}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column prop="transMethod" label="运输模式" :show-overflow-tooltip="true">
              <template slot-scope="scope" v-if="scope.row.transMethod">
                {{scope.row.transMethod|transportmode}}
              </template>
            </el-table-column>
            <el-table-column prop="carriageName" label="承运单位" :show-overflow-tooltip="true" width="110">
              <template slot-scope="scope">
                <span v-if="scope.row.carriageName == 1">物流公司</span>
                <span v-else-if="scope.row.carriageName == 2">网络货运</span>
                <span v-else-if="scope.row.carriageName == 3">铁路运输</span>
              </template>
            </el-table-column>
            <el-table-column label="承运平台" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.carriagePlatform">{{scope.row.carriagePlatform|logisticsplatform}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column prop="logisticsCompany" label="物流单位" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.logisticsCompany">{{scope.row.logisticsCompany}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column prop="transType" label="承运工具" :show-overflow-tooltip="true">
              <template slot-scope="scope" v-if="scope.row.transType">
                {{scope.row.transType|carrier}}
              </template>
            </el-table-column>
            <!-- <el-table-column prop="contractNumber" label="合同编号" sortable :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="baseSalesUnit" label="销售单位" sortable :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="baseBuyerUnit" label="购货单位" sortable :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="baseDeliveryPlace" label="发货地" sortable :show-overflow-tooltip="true" width="260">
            </el-table-column> -->
            <!-- <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.subjectMatter">{{scope.row.subjectMatter}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column> -->
            <el-table-column prop="loadingAddress" label="装货地" :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <el-table-column prop="unloadingPoint" label="卸货地" :show-overflow-tooltip="true" width="260">
            </el-table-column>
            <!-- <el-table-column label="发货日期" sortable width="200" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <div v-if="scope.row.issStartDate && scope.row.issEndDate">
                    <span>{{scope.row.issStartDate}}</span><span> ~ </span><span>{{scope.row.issEndDate}}</span>
                </div>
                <div v-else>---</div>
              </template>
            </el-table-column>
            <el-table-column label="收货日期" sortable width="200" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <div v-if="scope.row.recStartDate && scope.row.recEndDate">
                  <span>{{scope.row.recStartDate}}</span><span> ~ </span><span>{{scope.row.recEndDate}}</span>
                </div>
                <div v-else>---</div>
              </template>
            </el-table-column> -->
            <el-table-column label="车数" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.vehicles">{{scope.row.vehicles}} 车</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="运费金额合计" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.freightAmount">{{scope.row.freightAmount|currency}} 元</span>
                <span v-if="scope.row.freightAmount==0 || scope.row.freightAmount==''">---</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="吨数合计" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span>{{scope.row.subjectCount}} 吨</span>
              </template>
            </el-table-column> -->
            <el-table-column label="原发吨数合计" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span>{{scope.row.originalTonnage}} 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="实收吨数合计" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span>{{scope.row.tonnagePaid}} 吨</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <span>合计：运费金额 {{newcontractNum|currency}} 元 , 原发吨数 {{yuansumPrice|currency}} 吨 , 实收吨数 {{newsumPrice|currency}} 吨 </span>
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
      <!-- 抽屉 查看功能 -->
      <el-drawer :title="logo" :visible.sync="drawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <el-col :span="24">
                  <div class="see-top">附件信息</div>
                  <dl class="see-file">
                    <dd>
                      <el-col :span="5">
                        <span>合同附件：</span>
                      </el-col>
                      <el-col :span="19">
                        <a v-for="(v,k) in filelist" :key="k" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)" @click="fileGo(v,seeData)"><p>{{v.name}}</p></a>
                      </el-col>
                    </dd>
                  </dl>
                </el-col>
                <el-col :span="24" v-if="seeData">
                    <div class="see-top">基础信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>系统编码：</label>
                            <span>{{seeData.systemCode}}</span>
                        </li>
                        <!-- <li>
                            <label>运单编号：</label>
                            <span v-if="seeData.billNumber">{{seeData.billNumber}}</span>
                            <span v-else>---</span>
                        </li> -->
                        <li>
                            <label>运输模式：</label>
                            <span>{{seeData.transMethod|transportmode}}</span>
                        </li>
                        <li>
                            <label>承运单位：</label>
                            <span v-if="seeData.carriageName == 1">物流公司</span>
                            <span v-if="seeData.carriageName == 2">网络货运</span>
                        </li>
                        <li>
                            <label>承运平台：</label>
                            <span v-if="seeData.carriagePlatform">{{seeData.carriagePlatform|logisticsplatform}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>物流单位：</label>
                            <span v-if="seeData.logisticsCompany">{{seeData.logisticsCompany}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>承运工具：</label>
                            <span>{{seeData.transType|carrier}}</span>
                        </li>
                        <li>
                            <label>标的物：</label>
                            <span v-if="seeData.subjectMatter">{{seeData.subjectMatter}}</span>
                            <span v-else>---</span>
                        </li>
                        <li>
                            <label>装货地点：</label>
                            <span>{{seeData.loadingAddress}}</span>
                        </li>
                        <li>
                            <label>卸货地点：</label>
                            <span>{{seeData.unloadingPoint}}</span>
                        </li>
                        <!-- <li v-if="seeData.issStartDate && seeData.issEndDate">
                            <label>发货日期：</label>
                            <span>{{seeData.issStartDate}}</span><span> ~ </span><span>{{seeData.issEndDate}}</span>
                        </li>
                        <li v-else>
                            <label>发货日期：</label>
                            <span>---</span>
                        </li>
                        <li v-if="seeData.recStartDate && seeData.recEndDate">
                            <label>收货日期：</label>
                            <span>{{seeData.recStartDate}}</span><span> ~ </span><span>{{seeData.recEndDate}}</span>
                        </li>
                        <li v-else>
                            <label>收货日期：</label>
                            <span>---</span>
                        </li> -->
                    </ul>
                </el-col>
                <el-col :span="24">
                    <div class="see-top">运单统计</div>
                    <ul class="see-lable">
                        <li>
                            <label>车数：</label>
                            <span>{{seeData.vehicles}} 车</span>
                        </li>
                        <li>
                            <label>运费金额合计：</label>
                            <span v-if="seeData.freightAmount">{{seeData.freightAmount|currency}} 元</span>
                            <span v-if="seeData.freightAmount== 0 || seeData.freightAmount==''">---</span>
                        </li>
                        <li>
                            <label>吨数合计：</label>
                            <span>{{seeData.subjectCount}} 吨</span>
                        </li>

                        <li>
                            <label>原发吨数合计：</label>
                            <span>{{seeData.originalTonnage}} 吨</span>
                        </li>
                        <li>
                            <label>实收吨数合计：</label>
                            <span>{{seeData.tonnagePaid}} 吨</span>
                        </li>
                    </ul>
                </el-col>
                <el-col :span="24">
                    <div class="see-top">合同关联信息</div>
                    <ul class="see-lable" v-for="(v,k) in tags" :key="k">
                        <li>
                            <label>合同编号：</label>
                            <span>{{v.systemCode}}</span>
                        </li>
                        <li>
                            <label>销售单位：</label>
                            <span>{{v.salesUnit}}</span>
                        </li>
                        <li>
                            <label>购货单位：</label>
                            <span>{{v.buyerUnit}}</span>
                        </li>
                        <li>
                            <label>发货地：</label>
                            <span>{{v.deliveryPlace}}</span>
                        </li>
                    </ul>
                </el-col>
              </el-row>
            </el-tab-pane>
            <!-- <el-tab-pane label="动态" name="second">
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
            </el-tab-pane> -->
            <!-- 团队模块 V.0.0.2 上线  -->
            <!-- <el-tab-pane label="团队" name="third">团队</el-tab-pane> -->
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer>
      <!-- 文件预览弹窗 -->
      <el-dialog
        title="文件预览"
        :visible.sync="ImgVisible"
        width="44%"
        :before-close="ImgClose">
            <div class="imgBox" v-if="itemshow == true" v-loading="loadingimgss" element-loading-text="文件加载中">
              <!-- <iframe :src="url" frameborder="0" width="800px" height="600px"></iframe> -->
              <div v-for="i in numPages" :key="i">
                  <pdf :src="src" :page="i"></pdf>
              </div>
            </div>  
            <viewer class="imgBox" :images="inmages" v-else v-loading="loadingimg" element-loading-text="文件加载中"> 
              <!-- // photo 一定要一个数组，否则报错 -->
                <img
                v-for="(src,index) in inmages"
                :src="src.filePath"
                :key="index">
            </viewer>
            <div class="imgbuttons">
              <el-button type="primary" icon="el-icon-printer" size="mini" @click="downPrint">打印</el-button>
              <el-button type="primary" icon="el-icon-download" style="float:right;" size="mini" @click="openfile">下载</el-button>
            </div>
      </el-dialog>
        <!-- 磅单详情弹窗开始 -->
        <el-dialog
            :visible.sync="pronChoicecont"
            width="60%"
            :before-close="pronberhandleClose">
            <template slot="title">
                <div style="font-size:18px;">磅单详情</div>
            </template>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <span style="margin-left:24px;color:#333;">累计净重 <b style="color:#EA222E;">{{PronNum}}</b> 吨</span>
                        <el-input
                        @keyup.enter.native="pronractsearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="pronSearch"
                        style="width:200px;">
                        </el-input>
                    </div>
                    <div class="concentBoxtable">
                        <!-- 磅单表格 -->
                        <el-table
                        border
                        height="100%"
                        ref="prontractTable"
                        :data="pronTableData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;">
                            <el-table-column prop="id" v-if="proshow">
                            </el-table-column>
                            <el-table-column prop="licensePlate" label="车牌号">
                            </el-table-column>
                            <el-table-column prop="coalType" label="煤种">
                            </el-table-column>
                            <el-table-column label="皮重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.tare}} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="毛重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.roughWeight}} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="净重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.netWeight}} 吨</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="puopages">
                        <el-pagination
                        background
                        :page-sizes="[50,100,300,500,]"
                        :page-size.sync="pronSize"
                        :current-page.sync="pronCurrent"
                        layout="total, sizes, prev, pager, next"
                        :total="pronTotal"
                        @current-change='pronNtractall'
                        @size-change='pronNtractall'>
                        </el-pagination>
                        <el-button size="small" @click="pronApplication">关闭</el-button>  
                    </div>
                </div>
            </div>
        </el-dialog>
        <!-- 磅单详情弹窗结束 -->
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
import pdf from 'vue-pdf';
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js'; // 加载中文的包
export default {
  components: {pdf},
    data() {
      return {
        // 合计
        newcontractNum:0,       // 合计金额
        yuansumPrice:0,         // 原发
        newsumPrice:0,          // 实收
        activewill:'1',
        tags:[],
        http:'',
        dataAva:false,
        loading:false,
        current:1,              // 当前页数
        size:50,                // 显示条数
        total:0,                // 总条数
        delmul:[],              // 批量删除
        show:false,
        purTitle: "",           // 标题
        purSearch: "",          // 搜索        
        multipleSelection: [],  // 选中数据
        drawer: false,          // 查看抽屉
        logo: "",               // 抽屉企业名称
        activeName: "first",    // 查看tab栏
        tableData: [],          // 表格数据
        seeData:{},             // 查看数据 
        filelist:[],            // 文件列表
        // 查看...动态
        timeData:[
            // {
            //   createtime:'2021-01-02',
            //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
            // }
        ],  
        // 合同文件预览功能
        ImgVisible:false,
        inmages:[],
        numPages:'',
        src:'',
        path:'',                        // pdf的地址，例如：/testFile.pdf
        geshi:'',
        itemshow:null,
        loadingimg:true,
        loadingimgss:true,
        oneData:'',                     // 单条数据

        /* 
        磅单关联数据列表
        */
        proshow:false,
        pronCurrent:1,                      // 当前页数
        pronSize:50,                        // 显示条数
        pronTotal:0,                        // 总条数
        pronChoicecont:false,               // 选择磅单弹窗
        pronSearch:'',                      // 磅单搜索框
        pronTableData:[],                   // 选择磅单表格数据  
        PronNum:0,                          // 弹框累计净重吨数
        pronId:null,
      };
    },
    created() {},
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName') // 企业名称
      this.http  = this.$store.state.http;                 // IP地址
      this.getpurcon();
    },
    methods: {
      /* 
      磅单详情功能
      */
      pronPurs(id){
          this.pronId = id;
          this.pronChoicecont = true;                     // 开启弹窗
          this.pronNtractall();
      },
      // 表格数据
      pronNtractall(){
          // 传参
          var listByPageData = {
              orderId:this.pronId,
              searchStr:this.pronSearch,
              page:{
                  current:this.pronCurrent,
                  size:this.pronSize
              }
          }
          // 渲染表格
          this.api.Weighing.all(listByPageData)
          .then(res=>{
              if(res.data.code == 200){
                  this.pronTableData = res.data.data.records;
                  let num = 0;
                  this.pronTableData.map((item) =>{
                      num += item.netWeight*1;
                  })
                  this.PronNum = num.toFixed(2);
                  this.pronTotal = res.data.data.total;
              }
          })
      },
      // 磅单详情搜索功能
      pronractsearch() {
          this.pronNtractall();
          this.pronCurrent = 1;
      },
      // 磅单详情关闭弹窗
      pronberhandleClose() {
          this.pronId = null;
          this.PronNum = 0;
          this.pronSearch = '';
          this.pronChoicecont = false;                    
      },
      // 磅单详情关闭
      pronApplication(){
        this.pronId = null;
        this.PronNum = 0;
        this.pronSearch = '';
        this.pronChoicecont = false;                 
      },

      /*
      合同文件预览
      */
      fileGo(obj){
        this.oneData = obj;
        this.ImgVisible = true;
        let faliename = this.oneData.fileName;
        this.geshi = faliename.substring(faliename.length - 4);
        let arrs = ['.pdf','.PDF','docx','DOCX','.doc','.DOC'];
        if(!arrs.includes(this.geshi)){
          this.itemshow = false;
          this.api.file.preFIle({file_id:this.oneData.id})
          .then(res=>{
              if(res.data.code == 200){
                this.loadingimg = false;
                this.inmages = res.data.data;
              }
          })
        }else{
          this.itemshow = true;
          this.src = pdf.createLoadingTask({
            url:this.http + '/saas-common/upload/previewStream?fileName=' + faliename + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f' ,     // 正式环境用这个！
            CMapReaderFactory
          })
          // 让所有页数一次性加载完，否则就只会加载第一页
          this.src.promise.then(pdf => {
            console.log(pdf)
              this.loadingimgss = false;
              this.numPages = pdf.numPages;
          }).catch(() => {})
        }
      },
      // 点击文件下载
      openfile(){
          let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + this.oneData.fileName;
          window.open(url)
      },
      // 打印文件
      downPrint(){
        let url =  this.http + '/saas-common/upload/previewStream?fileName=' + this.oneData.fileName + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f';     // 正式环境用这个！
        window.open(url)
      },
      // 关闭查看预览文件弹窗
      ImgClose() {
          this.ImgVisible = false;
          this.src = '';
          this.inmages = [];
          this.numPages = '';
          this.path = '';     // pdf的地址，例如：/testFile.pdf
          this.geshi = '';
          this.loadingimg = true;
          this.loadingimgss = true;
      },
      // 表格数据
      getpurcon(){
        this.loading = true;
        // 传参
        var listByPageData = {
          searchStr:this.purSearch,
          page:{
            current:this.current,
            size:this.size
          }
        }
        // 渲染表格
        this.api.logistics.all(listByPageData)
        .then(res=>{
          if(res.data.code == 200){
            this.loading = false; 
            if(res.data.data.records != ""){
              this.tableData = res.data.data.records
              this.total = res.data.data.total
            }else{
              this.dataAva = true;
            }
          }
        })
        // 合计
        this.newcontractNum = 0;
        this.yuansumPrice = 0;
        this.newsumPrice = 0;
        this.api.logistics.getSumPrice()
        .then(res=>{
          if(res.data.code == 200 && res.data.data != null){
            this.newcontractNum = res.data.data.freightAmount;
            this.yuansumPrice = res.data.data.originalTonnage;
            this.newsumPrice = res.data.data.tonnagePaid;
          }
        })
      },
      // 搜索功能
      search() {
        this.tableData = [];
        this.current = 1;
        this.getpurcon();
      },
      // 表格选中数据
      purChange(val) {
        this.multipleSelection = val;
        this.delmul = [];
        for(let i=0;i<this.multipleSelection.length;i++){
          this.delmul.push(this.multipleSelection[i].id) 
        }
      },
      // 新建运单信息
      addPurs(){
          this.$router.push({name:"AddWaybill"})
      },
      // 编辑运单信息
      EditPurs(row){
        this.api.logistics.one({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            // 发货日期回显
            if(res.data.data.issStartDate && res.data.data.issEndDate){
              res.data.data.Time = [res.data.data.issStartDate,res.data.data.issEndDate];
            }else{
              res.data.data.Time = [];
            }
            // 收货日期回显
            if(res.data.data.recStartDate && res.data.data.recEndDate){
              res.data.data.reTime = [res.data.data.recStartDate,res.data.data.recEndDate];
            }else{
              res.data.data.reTime = [];
            }
            res.data.data.ids = [];
            let purdata = res.data.data;
            purdata.id = row.id;
            sessionStorage.setItem('queryAlls',JSON.stringify(purdata))
            this.$router.push({name:"EditWaybill"})
          }
        })
      },
      // 查看
      purSee(index, row) {
        this.drawer = true;
        this.api.logistics.one({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            this.seeData = res.data.data;
          }
        })
        // 关联合同数据集合
        this.api.logistics.oneAll({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            this.tags = res.data.data;
          }
        })
        // 文件
        this.api.file.one({relationId:row.id})
        .then(res=>{
          this.filelist = res.data.data
        })
      },
      changeActive(e) {
          e.currentTarget.className = 'active';
      },
      removeActive(e) {
          e.currentTarget.className = '';
      },
      // 抽屉
      handleClick(tab) {
        // console.log(tab)
        // 判断tab栏选中状态，请求相应的接口
        if (tab.label == "基本信息") {
        }
        if (tab.label == "动态") {
          // console.log(2);
        }
      },
      // 删除功能
      openPur() {
          if(this.multipleSelection == ""){
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
              this.api.logistics.del(this.delmul)
              .then(res=>{
                if(res.data.code == 200){
                  this.$message({
                    type: 'success',
                    message: '删除成功!'
                  });
                  this.tableData = [];
                  this.current = 1;
                  this.getpurcon();
                  if(this.tableData == []){
                    this.dataAva = true;
                  }
                }
              })
            })
          }
      },
      // 查看运单详情
      OnePurs(data){
        sessionStorage.setItem('DetailsId',JSON.stringify(data.id))
        this.$router.push({name:"WaybillDetails"})
        // console.log(data)
      },
    },
};
</script>
<style lang="less" scoped>
// 磅单详情弹窗样式开始
.Wdialog{
    width: 100%;
    border: 1px solid #E0E0E0;
    margin-top: 16px;
}
// 选项卡
.concentBox {
    width: 100%; 
    position: relative;
    .concentBoxtable{
        height:330px;
    }
}
.textbutton{
    width:56px;
    height: 32px;
    position: absolute;
    right:25px;
    bottom: 10px;
}
.totalTab {
    width: 9%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin:9px 0 0 15px;
}
.textlast{
    width: 50%;
    height: 46px;
    padding: 9px 17px 5px 24px;
    float: right;
    span{
        font-size: 16px;
        line-height: 32px;
        color: #EA222E;
        font-weight: 600;
        cursor: pointer;
        float: right;
        margin-right:24px;
    }
    .el-input{
        float: right;
    }
}
.totalTab.active {
    color: #333;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:7px auto 16px;
    display: none;
}
.strip {
    display: block;
}
.text {
  width: 100%;
}
// 磅单分页
.puopages {
    width: 100%;
    height: 32px;
    margin: 148px 0 9px 9px;
    .el-pagination {
      float: left;
    }
    .el-button{
        float: right;
        margin-right: 24px;
    }
}
// 磅单表单
.CompanyData{
    width: 100%;
    padding: 8px 16px;
    position: absolute;
    left: 0;
    bottom:115px;
}
// 磅单详情弹窗样式结束

// 文件图片查看样式
.imgBox{
    width: 100%;
    height: 600px;
    overflow: hidden;
    overflow-y: scroll;
    position: relative;
    img{
      width: 100%;
    }
}
.imgbuttons{
  width:24%;
  height: 40px;
  position: absolute;
  top: 0;
  right: 0;
  bottom: -90%;
  left: 0;
  margin: auto;
  z-index: 10000;
}
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
  .pur-table{
    width:100%;
    height: 94%;
  } 
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
      left:0;
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
  // 分页
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
}
// 内容结束
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
  padding: 12px 0;
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
    a {
      span {
        color: #333;
        
      }
    }
  }
}
// 文件样式
.see-file{
  padding: 15px 0 ;
  overflow: hidden;
  span {
    font-size: 14px;
    color: #666;
  }
  dd {
    font-size: 14px;
    line-height: 34px;
    color: #333;
    a{
      text-decoration: none;
      color: #333;
    }
    .active{
      color: #EA222E;
    }
  }
}
// 查看抽屉样式结束
</style>