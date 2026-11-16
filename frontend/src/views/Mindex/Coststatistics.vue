<template>
  <!-- 缴税明细页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <el-button size="small" @click="addPurs" v-if="$route.meta.add">新增</el-button>
        <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
        <!-- <el-input
        @keyup.enter.native="searchtwo"
        size="small"
        placeholder="卸货地搜索"
        prefix-icon="el-icon-search"
        v-model="purSearchtwo"
        >
        </el-input> -->
        <!-- <el-input
        @keyup.enter.native="searchone"
        size="small"
        placeholder="发货地搜索"
        prefix-icon="el-icon-search"
        v-model="purSearchone"
        >
        </el-input>
        <el-input
        @keyup.enter.native="searchtree"
        size="small"
        placeholder="运单编号"
        prefix-icon="el-icon-search"
        v-model="purSearchtree">
        </el-input>
        <el-date-picker
        @change="DataTime"
        size="small"
        v-model="keyupTime"
        type="date"
        value-format="yyyy-MM-dd"
        placeholder="运单创建时间">
        </el-date-picker>
        <el-select v-model="carriageName" size="small" clearable @change="prologistics" placeholder="网络货运平台搜索">
            <el-option
            v-for="item in Extended"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue">
            </el-option>
        </el-select> -->
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
            <el-table-column type="selection" fixed="left"> </el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column v-if="$route.meta.edit"
              label="操作"
              width="100" fixed="left">
              <template slot-scope="scope">
                <el-button @click="EditPurs(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="month" label="所属月份" sortable :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.month">{{scope.row.month|timeDate}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column prop="includMoney" label="销项含税金额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.includMoney">{{scope.row.includMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column prop="baseMoney" label="不含税金额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.baseMoney">{{scope.row.baseMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column prop="taxAmount" label="销项税额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.taxAmount">{{scope.row.taxAmount|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column prop="jxMoney" label="进项税额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.jxMoney">{{scope.row.jxMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="增值税额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.zzMoney">{{scope.row.zzMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="城建税" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.cjMoney">{{scope.row.cjMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="教育费附加税额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.jyMoney">{{scope.row.jyMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="地方教育费附加" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.dfjyMoney">{{scope.row.dfjyMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="印花税额" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.yhMoney">{{scope.row.yhMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="企业所得税" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.qyMoney">{{scope.row.qyMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="滞纳金" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.znjMoney">{{scope.row.znjMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column label="缴税合计" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.sumMoney">{{scope.row.sumMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" :show-overflow-tooltip="true" width="140">
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
                <el-button size="small" @click="addPurs" v-if="$route.meta.add">新增</el-button>
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
      <!-- 抽屉 查看功能 -->
      <!-- <el-drawer :title="logo" :visible.sync="drawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <el-col :span="24" v-if="seeData">
                    <div class="see-top">基础信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>系统编码：</label>
                            <span>{{seeData.systemCode}}</span>
                        </li>
                        <li>
                            <label>运单编号：</label>
                            <span v-if="seeData.billNumber">{{seeData.billNumber}}</span>
                            <span v-else>---</span>
                        </li>
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
                            <span>{{seeData.logisticsCompany|logisticsplatform}}</span>
                        </li>
                        <li>
                            <label>承运工具：</label>
                            <span>{{seeData.transType|carrier}}</span>
                        </li>
                        <li>
                            <label>标的物：</label>
                            <span>{{seeData.subjectMatter}}</span>
                        </li>
                        <li>
                            <label>数量：</label>
                            <span>{{seeData.subjectCount}} 吨</span>
                        </li>
                        <li>
                            <label>装货地点：</label>
                            <span>{{seeData.loadingAddress}}</span>
                        </li>
                        <li>
                            <label>运费金额：</label>
                            <span v-if="seeData.freightAmount">{{seeData.freightAmount|currency}} 元</span>
                            <span v-if="seeData.freightAmount== 0 || seeData.freightAmount==''">---</span>
                        </li>
                        <li>
                            <label>卸货地点：</label>
                            <span>{{seeData.unloadingPoint}}</span>
                        </li>
                    </ul>
                </el-col>
                <el-col :span="24">
                    <div class="see-top">磅单原发信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>发货日期：</label>
                            <span>{{seeData.issStartDate}}</span><span> ~ </span><span>{{seeData.issEndDate}}</span>
                        </li>
                        <li>
                            <label>原发吨数：</label>
                            <span>{{seeData.originalTonnage}} 吨</span>
                        </li>
                    </ul>
                </el-col>
                <el-col :span="24">
                    <div class="see-top">磅单实收信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>收货日期：</label>
                            <span>{{seeData.recStartDate}}</span><span> ~ </span><span>{{seeData.recEndDate}}</span>
                        </li>
                        <li>
                            <label>实收吨数：</label>
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
                <el-col :span="24">
                  <div class="see-top">附件信息</div>
                  <dl class="see-file">
                    <dd>
                      <el-col :span="5">
                        <span>合同附件：</span>
                      </el-col>
                      <el-col :span="19">
                        <a v-for="(v,k) in filelist" :key="k"  :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"><p>{{v.name}}</p></a>
                      </el-col>
                    </dd>
                  </dl>
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
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer> -->
      <!-- 导入弹窗功能 -->
      <el-dialog
        title="缴税明细导入"
        :visible.sync="ImportVisible"
        width="60%"
        :before-close="ImporthandleClose">
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
              <el-form-item prop="importnav">
                <el-input placeholder="仅限xlsx文件" v-model="fileForm.importnav" size="small">
                    <el-upload slot="append"
                    ref='upload'
                    :auto-upload="resfile"
                    :action="http + '/saas-bill/chainBillRate/importExcel'"
                    :headers="header"
                    name='file'
                    :show-file-list="false"
                    :on-success="addFileLet"
                    :on-change="supbeforedemoUpload"
                    :limit="1"><el-button size="mini" style="color:#333;">请选择文件</el-button>
                    </el-upload>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <a :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=QYJS.xlsx'">下载导入模板</a>
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
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            // Extended:[],             // 网络货运平台数据
            purSearch: "",              // 搜索  
            // purSearchtree:'',        // 运单编号搜索
            // keyupTime:'',            // 运单创建时间搜索 
            // carriageName:'',         // 网络货运平台搜索
            // purSearchone: "",        // 发货地搜索
            // purSearchtwo: "",        // 卸货地搜索
            // tags:[],
            http:'',
            dataAva:false,
            loading:false,
            current:1,                  // 当前页数
            size:50,                    // 显示条数
            total:0,                    // 总条数
            delmul:[],                  // 批量删除
            show:false,
            purTitle: "",               // 标题
            multipleSelection: [],      // 选中数据
            drawer: false,              // 查看抽屉
            logo: "",                   // 抽屉企业名称
            activeName: "first",        // 查看tab栏
            tableData: [],              // 表格数据
            seeData:{},                 // 查看数据
            filelist:[],                // 文件列表
            // 查看...动态
            timeData:[
                // {
                //   createtime:'2021-01-02',
                //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ], 
            ImportVisible:false,                // 导入弹窗
            ImportVisible:false,                // 取消导入
            header:{},                          // 上传token
            resfile:false,                      // 是否立即上传
            active:1,                           // 步骤条
            oneShow:true,                       // 第一步
            // 导入上传验证
            fileForm:{
                importnav:'',
            },
            rules: {
                importnav: [
                    { required: true, message: '请上传xlsx 或 xls文件', trigger: 'blur' },
                ],
            },
            twoShow:false,                      // 第二步
            perstep:0,                          // 进度%
            pernav:'正在导入...',
            streeShow:false,                    // 第三步...失败
            cation:'',                          // 导入错误提示
            fourShow:false,                     // 第三步...成功 
            sucessShow:false,           // 第三步...全部成功
            errorShow:false,            // 第三步...非全部成功
            sucessData:[],              // 正确数据
            errorData:[],               // 错误数据 
        };
    },
    created() {}, 
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName')     // 企业名称
      this.http  = this.$store.state.http;                  // IP地址
      // 基础维护数据
    //   this.api.basics.all({template_id:1})
    //   .then(res=>{
    //       if(res.data.code == 200){
    //           this.Extended = res.data.data.logistics_platform;
    //       }
    //   });
      this.getpurcon();
    },
    methods: {
        // 新建缴税明细
        addPurs(){
            this.$router.push({name:"AddCoststatistics"})
        },
        // 编辑缴税明细
        EditPurs(row){
          this.api.billRate.one({id:row.id})
          .then(res=>{
            if(res.data.code == 200){
              let purdata = res.data.data;
              purdata.id = row.id;
              sessionStorage.setItem('queryAll',JSON.stringify(purdata));
              this.$router.push({name:"EditCoststatistics"})
            }
          })
        },
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
            searchStr:this.purSearch,
            // billNumber:this.purSearchtree,
            // billData:this.keyupTime,
            // loadingAddress:this.purSearchone,
            // unloadingPoint:this.purSearchtwo,
            // carriagePlatform:this.carriageName,
            page:{
                current:this.current,
                size:this.size
            }
            }
            // 渲染表格
            this.api.billRate.all(listByPageData)
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
        },
        // 搜索功能
        search() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // // 发货地搜索功能
        // searchone() {
        //     this.tableData = [];
        //     this.getpurcon();
        //     this.current = 1;
        // },
        // // 卸货地搜索功能
        // searchtwo() {
        //     this.tableData = [];
        //     this.getpurcon();
        //     this.current = 1;
        // },
        // // 运单编号搜索功能
        // searchtree() {
        //     this.tableData = [];
        //     this.getpurcon();
        //     this.current = 1;
        // },
        // // 日期搜索
        // DataTime() {
        //     this.tableData = [];
        //     this.getpurcon();
        //     this.current = 1;
        // },
    //    // 承运平台搜索功能
    //    prologistics(){
    //       this.tableData = [];
    //       this.getpurcon();
    //       this.current = 1;
    //    },
        // 表格选中数据
        purChange(val) {
            this.multipleSelection = val;
            this.delmul = [];
            for(let i=0;i<this.multipleSelection.length;i++){
            this.delmul.push(this.multipleSelection[i].id) 
            }
        },
        // 查看
        // purSee(index, row) {
        //     this.drawer = true;
        //     this.api.billRate.one({id:row.id})
        //     .then(res=>{
        //     if(res.data.code == 200){
        //         this.seeData = res.data.data;
        //     }
        //     })
        //     // 关联合同数据集合
        //     this.api.logistics.oneAll({id:row.id})
        //     .then(res=>{
        //     if(res.data.code == 200){
        //         this.tags = res.data.data;
        //     }
        //     })
        //     // 查看文件
        //     this.api.file.one({relationId:row.id})
        //     .then(res=>{
        //     this.filelist = res.data.data
        //     })
        // },
        // changeActive(e) {
        //     e.currentTarget.className = 'active';
        // },
        // removeActive(e) {
        //     e.currentTarget.className = '';
        // },
        // 抽屉
        handleClick(tab) {
            // console.log(tab.label);
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "动态") {
                // console.log(2);
            } else if (tab.label == "团队") {
                // console.log(3);
            }
        },
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
            this.getpurcon();
            this.current = 1;
            this.$refs.upload.clearFiles();
            this.fileForm.importnav = "";
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
                this.api.billRate.del(this.delmul)
                .then(res=>{
                    if(res.data.code == 200){
                    this.$message({
                        type: 'success',
                        message: '删除成功!'
                    });
                    this.tableData = [];
                    this.getpurcon();
                    this.current = 1;
                    if(this.tableData == []){
                        this.dataAva = true;
                    }
                    }
                })
                })
            }
        },
    },
};
</script>
<style lang="less" scoped>
// 导入样式开始
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
  a{
    margin-left: 16%;
    font-size: 14px;
    color:#CF1724;
    text-decoration:none;
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
    .el-select{
      width: 11%;
      margin-right: 10px;
      float: right;
    }
    .el-input {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .el-button {
      margin-right: 10px;
      float: right;
      
    }
    .el-date-picker{
      width: 10%;
      float: right;
      margin-right: 10px;
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
      left:25%;
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