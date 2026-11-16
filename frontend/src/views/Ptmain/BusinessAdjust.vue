<template>
  <!-- 业务核算页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
        <el-col :span="12">
            <span class="pur-size">{{ purTitle }}</span>
            <el-select v-model="AdjustValue" placeholder="请选择" size="small" style="width:40%;" @change="AdjustChang" filterable>
                <el-option
                v-for="item in projectName"
                :key="item.value"
                :label="item.enterpriseName"
                :value="item.id">
                </el-option>
            </el-select>
        </el-col>
        <el-col :span="12" class="pur-right">
            <el-button type="primary" size="small" @click="AddopenPur" v-if="$route.meta.add">新增</el-button>
            <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
            <el-input
            @keyup.enter.native="search"
            size="small"
            placeholder="请输入内容"
            prefix-icon="el-icon-search"
            v-model="purSearch">
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
            <div class="addpurs" v-if="dialogVisibledata">
                <el-form class="AdjustForm" :model="AdjustForm" :rules="rules" ref="AdjustForm" size="small" label-width="auto">
                    <el-row>
                        <el-col :span="3">
                            <el-form-item prop="month" label="时间">
                                    <el-date-picker
                                    :editable="false"
                                    v-model="AdjustForm.month"
                                    type="month"
                                    placeholder="所属月份"
                                    value-format="yyyy-MM"
                                    style="width:110%;">
                                    </el-date-picker>
                            </el-form-item>
                        </el-col>
                        <el-col :span="4" style="margin-right:1%;">
                            <el-form-item prop="companyId" label="企业名称" style="width:100%;">
                                    <el-select size="small" v-model="AdjustForm.companyId" placeholder="请选择企业名称" disabled>
                                        <el-option
                                        v-for="item in projectName"
                                        :key="item.value"
                                        :label="item.enterpriseName"
                                        :value="item.id">
                                        </el-option>
                                    </el-select>
                            </el-form-item>
                        </el-col> 
                        <el-col :span="3">
                            <el-form-item prop="transCount" label="运输发票吨数">
                                <el-input size="small" v-model="AdjustForm.transCount" placeholder="单位(吨)"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="3">
                            <el-form-item prop="logisticsCount" label="开票数量">
                                <el-input size="small" v-model="AdjustForm.logisticsCount" placeholder="单位(吨)"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="3">
                            <el-form-item prop="busCount" label="开票金额">
                                <el-input size="small" v-model="AdjustForm.busCount" placeholder="单位(元)"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="3">
                            <el-form-item prop="taxPayment" label="缴税金额">
                                <el-input size="small" v-model="AdjustForm.taxPayment" placeholder="单位(元)"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="4" style="padding-left:1%;">
                            <el-button size="small" style="margin-left: 1%"  @click="AdjustresetForm">取 消</el-button>
                            <el-button type="primary" size="small" @click="AddDatasubmitForm('AdjustForm')">保 存</el-button>
                        </el-col>
                    </el-row>
                </el-form>
            </div>
            <el-table
            height="100%"
            border
            :data="tableData"
            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
            style="width: 100%;">
                <el-table-column
                label="操作" width="80" fixed="left" v-if="$route.meta.delete">
                <template slot-scope="scope">
                    <el-button size="mini" @click="delPurs(scope.row)">删除</el-button>
                </template>
                </el-table-column>
                <el-table-column prop="month" label="时间" :show-overflow-tooltip="true" width="200" sortable>
                </el-table-column>
                <el-table-column prop="enterpriseName" label="企业名称" :show-overflow-tooltip="true">
                </el-table-column> 
                <el-table-column prop="busCount" label="开票金额(元)" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <div v-if="scope.row.jejustShow">
                            <el-input size="small" v-focus v-model="newbusCount" @blur="CollectionloseFcous(scope.$index, scope.row)"></el-input>
                        </div>
                        <div v-else>
                            <span>{{scope.row.busCount}}</span>
                            <i v-if="$route.meta.edit" class="el-icon-edit" @click="Colleclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                        </div>
                    </template> 
                </el-table-column>
                <el-table-column prop="logisticsCount" label="开票数量(吨)" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <div v-if="scope.row.sljustShow">
                            <el-input size="small" v-focus v-model="newlogisticsCount" @blur="logistloseFcous(scope.$index, scope.row)"></el-input>
                        </div>
                        <div v-else>
                            <span>{{scope.row.logisticsCount}}</span>
                            <i v-if="$route.meta.edit" class="el-icon-edit" @click="logistColleclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                        </div>
                    </template> 
                </el-table-column>
                <el-table-column prop="transCount" label="运输发票吨数(吨)" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <div v-if="scope.row.ysjustShow">
                            <el-input size="small" v-focus v-model="newtransCount" @blur="transCountFcous(scope.$index, scope.row)"></el-input>
                        </div>
                        <div v-else>
                            <span>{{scope.row.transCount}}</span>
                            <i v-if="$route.meta.edit" class="el-icon-edit" @click="transCountColleclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                        </div>
                    </template> 
                </el-table-column>
                <el-table-column prop="taxPayment" label="缴税额(元)" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <div v-if="scope.row.jsljustShow">
                            <el-input size="small" v-focus v-model="newtaxPaymentCount" @blur="taxPaymentFcous(scope.$index, scope.row)"></el-input>
                        </div>
                        <div v-else>
                            <span>{{scope.row.taxPayment}}</span>
                            <i v-if="$route.meta.edit" class="el-icon-edit" @click="taxPaymentColleclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                        </div>
                    </template> 
                </el-table-column>
                <!-- 暂无数据 -->
                <div slot="empty" class="dataAvailable" v-if="dataAva">
                    <div class="dataAva">
                        <img src="../../assets/zanwushuju.png" alt="">
                        <p>暂无数据 !</p>
                        <!-- <el-button size="small" @click="Importdailog">导入</el-button>
                        <el-button size="small" @click="AddopenPur">新增</el-button> -->
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
        <!-- 新增弹窗 -->
        <!-- <el-dialog
        title="新增"
        :visible.sync="dialogVisibledata"
        width="40%">
            <div class="elTabBox">
                <el-form :model="AdjustForm" :rules="rules" ref="AdjustForm" size="small" label-width="120px">
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">时间</label>
                        </el-col>
                        <el-form-item prop="month">
                            <el-col :span="20">
                                <el-date-picker
                                v-model="AdjustForm.month"
                                type="month"
                                placeholder="请选择所属月份"
                                value-format="yyyy-MM"
                                style="width:100%;">
                                </el-date-picker>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6" >
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">企业名称</label>
                        </el-col>   
                        <el-form-item prop="companyId">
                            <el-col :span="20">
                                <el-select size="small" filterable v-model="AdjustForm.companyId" placeholder="请选择企业名称" style="width:100%;">
                                    <el-option
                                    v-for="item in projectName"
                                    :key="item.value"
                                    :label="item.enterpriseName"
                                    :value="item.id">
                                    </el-option>
                                </el-select>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">开票金额</label>
                        </el-col>
                        <el-form-item prop="busCount">
                            <el-col :span="20">
                                <el-input size="small" v-model="AdjustForm.busCount" placeholder="请输入开票金额金额(万元)"></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">开票数量</label>
                        </el-col>
                        <el-form-item prop="logisticsCount">
                            <el-col :span="20">
                                <el-input size="small" v-model="AdjustForm.logisticsCount" placeholder="请输入开票数量吨数(万吨)"></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <div class="adjusttop">
                        <el-button size="small" style="margin-left: 24px"  @click="dialogVisibledata = false">取 消</el-button>
                        <el-button type="primary" size="small" @click="AddDatasubmitForm('AdjustForm')">确 定</el-button>
                    </div>
                </el-form>
            </div>
        </el-dialog> -->
        <!-- 编辑弹窗 -->
        <!-- <el-dialog
        title="编辑"
        :visible.sync="AdjustVisibledata"
        width="40%">
            <div class="elTabBox">
                <el-form :model="editAdjustForm" :rules="editrules" ref="editAdjustForm" size="small" label-width="120px">
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">时间</label>
                        </el-col>
                        <el-form-item prop="month">
                            <el-col :span="20">
                                <el-date-picker
                                v-model="editAdjustForm.month"
                                type="month"
                                placeholder="请选择所属月份"
                                value-format="yyyy-MM"
                                style="width:100%;">
                                </el-date-picker>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6" >
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">企业名称</label>
                        </el-col>   
                        <el-form-item prop="companyId">
                            <el-col :span="20">
                                <el-select size="small" filterable v-model="editAdjustForm.companyId" placeholder="请选择企业名称" style="width:100%;">
                                    <el-option
                                    v-for="item in projectName"
                                    :key="item.value"
                                    :label="item.enterpriseName"
                                    :value="item.id">
                                    </el-option>
                                </el-select>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">开票数量</label>
                        </el-col>
                        <el-form-item prop="logisticsCount">
                            <el-col :span="20">
                                <el-input size="small" v-model="editAdjustForm.logisticsCount" placeholder="请输入开票数量(万吨)"></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="6">
                            <span style="color:red;">* </span>
                            <label style="color:#666;font-size:14px;line-height:40px;">开票金额</label>
                        </el-col>
                        <el-form-item prop="busCount">
                            <el-col :span="20">
                                <el-input size="small" v-model="editAdjustForm.busCount" placeholder="请输入开票金额(万元)"></el-input>
                            </el-col>
                        </el-form-item>
                    </el-row>
                    <div class="adjusttop">
                        <el-button size="small" style="margin-left: 24px"  @click="AdjustVisibledata = false">取 消</el-button>
                        <el-button type="primary" size="small" @click="editDatasubmitForm('editAdjustForm')">确 定</el-button>
                    </div>
                </el-form>
            </div>
        </el-dialog> -->
        <!-- 导入弹窗功能 -->
        <el-dialog
            title="业务核算导入"
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
                    <!-- importOnce    importExcelCg -->
                        <el-upload slot="append"
                        ref='upload'
                        :auto-upload="resfile"
                        :action="http + '/saas-common/chainBasicBussiness/importExcel'"
                        :headers="header"
                        name='file'
                        :show-file-list="false"
                        :on-success="addFileLet"
                        :on-change="supbeforedemoUpload"
                        :limit="1"><el-button size="mini">请选择文件</el-button>
                        </el-upload>
                    </el-input>
                </el-form-item>
                </el-col>
                <el-col :span="24">
                <a :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=YWHS.xlsx'">下载导入模板</a>
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
                <el-button size="mini" @click="ImportNextstree">完成</el-button>
            </div>
            </div>
            <!-- 步骤三...成功 -->
            <div class="importfour" v-show="fourShow">
            <div class="importfour-Box">
                <el-progress type="circle" :percentage="100" status="success"></el-progress>
                <span>数据完成导入，请在列表中点击查看</span>
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
            AdjustValue:27,                     // 头部下拉企业选中
            projectName:[],                     // 头部下拉框企业列表数据
            purSearch: "",                      // 搜索  
            http:'',
            dataAva:false,
            loading:false,
            current:1,                          // 当前页数
            size:50,                            // 显示条数
            total:0,                            // 总条数
            purTitle: "",                       // 标题
            tableData: [],                      // 表格数据
            // 导入
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
            /*
            新增
            */
            dialogVisibledata:false,
            AdjustForm:{},
            rules: {
                month: [
                    { required: true, message: '请选择所属月份', trigger: 'blur' }
                ],
                companyId: [
                    { required: true, message: '请选择企业名称', trigger: 'blur' }
                ],
                busCount: [
                    { required: true, message: '请输入开票金额', trigger: 'blur' },
                    // { pattern: /^[0-9]{1,3}(,[0-9]{3})*(.[0-9]{1,2})?$/, message: '金额格式不正确', trigger: 'blur'}
                ],
                transCount: [
                    { required: true, message: '请输入运输发票数', trigger: 'blur' },
                    // { pattern: /^[0-9]{1,3}(,[0-9]{3})*(.[0-9]{1,2})?$/, message: '发票数格式不正确', trigger: 'blur'}
                ],
                logisticsCount: [
                    { required: true, message: '请输入开票数量', trigger: 'change' },
                    // { pattern: /^[0-9]{1,3}(,[0-9]{3})*(.[0-9]{1,2})?$/, message: '开票数格式不正确', trigger: 'blur'}
                ],
                taxPayment: [
                    { required: true, message: '请输入缴税量', trigger: 'change' },
                    // { pattern: /^[0-9]{1,3}(,[0-9]{3})*(.[0-9]{1,2})?$/, message: '金额格式不正确', trigger: 'blur'}
                ],
            },
            /*
            编辑
            */
            newbusCount:'',         // 编辑开票金额
            newtransCount:'',       // 运输发票吨数
            newlogisticsCount:'',   // 编辑开票数量
            newtaxPaymentCount:'',  // 编辑缴税量
            editAdjustForm:{},
            // editrules: {
            //     month: [
            //         { required: true, message: '请选择所属月份', trigger: 'blur' }
            //     ],
            //     companyId: [
            //         { required: true, message: '请选择企业名称', trigger: 'blur' }
            //     ],
            //     busCount: [
            //         { required: true, message: '请输入开票金额(万元)', trigger: 'blur' },
            //     ],
            //     logisticsCount: [
            //         { required: true, message: '请输入开票数量(万吨)', trigger: 'change' }
            //     ],
            // },
        };
    },
    created() {},
    mounted(){
        this.purTitle = this.$route.meta.title;
        this.logo = sessionStorage.getItem('companyName') // 企业名称
        this.http  = this.$store.state.http;                 // IP地址
        // 企业列表
        this.api.adjust.selEcts({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
        .then(res=>{
            if(res.data.code == 200){
                this.projectName = res.data.data; // 企业名称
            }
        })
        this.getpurcon();
    },
    directives: {
        // 注册一个局部的自定义指令 v-focus
        focus: {
            // 指令的定义
            inserted: function (el) {
            // 聚焦元素
            el.querySelector('input').focus()
            }
        }
    },
    methods: {
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                companyId:this.AdjustValue,
                searchStr:this.purSearch,
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.adjust.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.loading = false; 
                    if(res.data.data.records != ""){
                        for(let i=0;i<res.data.data.records.length;i++){
                            res.data.data.records[i].jejustShow = false;
                            res.data.data.records[i].sljustShow = false;
                            res.data.data.records[i].jsljustShow = false;
                            res.data.data.records[i].ysjustShow = false;
                        }
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                    }else{
                        this.dataAva = true;
                    }
                }
            })
        },
        // 新增弹窗
        AddopenPur(){
            this.dialogVisibledata = true;
            this.api.adjust.selEcts({countyCode:sessionStorage.getItem('xzqhCode'),searchStr:''})
            .then(res=>{
                if(res.data.code == 200){
                    this.projectName = res.data.data; // 企业名称列表
                }
            })
            this.AdjustForm.companyId = this.AdjustValue;
        },
        // 新增表单
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    let arr = this.projectName.filter(item=>{
                        if(item.id == this.AdjustForm.companyId){
                            return item
                        }
                    })
                    this.AdjustForm.enterpriseName = arr[0].enterpriseName;
                    this.api.adjust.add(this.AdjustForm)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({  
                                type: "success",
                                message: "新建成功!"
                            });
                            this.dialogVisibledata = false;
                            this.AdjustForm = {};
                            this.getpurcon();
                        }
                    })
                }
            });
        },
        // 新增取消
        AdjustresetForm(){
            this.dialogVisibledata = false;
            this.AdjustForm = {};
        },
        // 编辑开票金额
        Colleclick(row){
            row.jejustShow = true;
            this.newbusCount = row.busCount;
        },
        // 编辑开票金额失去焦点事件
        CollectionloseFcous(index,val){
            if(this.newbusCount != ''){
                val.jejustShow = false;
                this.editAdjustForm = val;
                if(this.editAdjustForm.busCount != this.newbusCount *1){
                    this.editAdjustForm.busCount = this.newbusCount *1;
                    this.api.adjust.edit(this.editAdjustForm)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "编辑成功!"
                            });
                            this.editAdjustForm = {};
                            this.getpurcon();
                        }
                    })
                }
            }else{
                this.$message.warning('开票金额不能为空!');
            }
        },
        // 编辑开票数量
        logistColleclick(row){
            row.sljustShow = true;
            this.newlogisticsCount = row.logisticsCount;
        },
        // 编辑开票数量失去焦点事件
        logistloseFcous(index,val){
            if(this.newlogisticsCount != ''){
                val.sljustShow = false;
                this.editAdjustForm = val;
                if(this.editAdjustForm.logisticsCount != this.newlogisticsCount *1){
                    this.editAdjustForm.logisticsCount = this.newlogisticsCount *1;
                    this.api.adjust.edit(this.editAdjustForm)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "编辑成功!"
                            });
                            this.editAdjustForm = {};
                            this.getpurcon();
                        }
                    })
                }
            }else{
                this.$message.warning('开票数量不能为空!');
            }
        },
        // 编辑缴税量
        taxPaymentColleclick(row){
            row.jsljustShow = true;
            this.newtaxPaymentCount = row.taxPayment;
        },
        // 编辑缴税量失去焦点事件
        taxPaymentFcous(index,val){
            if(this.newtaxPaymentCount != ''){
                val.jsljustShow = false;
                this.editAdjustForm = val;
                if(this.editAdjustForm.taxPayment != this.newtaxPaymentCount *1){
                    this.editAdjustForm.taxPayment = this.newtaxPaymentCount *1;
                    this.api.adjust.edit(this.editAdjustForm)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "编辑成功!"
                            });
                            this.editAdjustForm = {};
                            this.getpurcon();
                        }
                    })
                }
            }else{
                this.$message.warning('缴税量不能为空!');
            }
        },
        // 编辑运输发票吨数失去焦点事件
        transCountFcous(index,val){
            if(this.newtransCount != ''){
                val.ysjustShow = false;
                this.editAdjustForm = val;
                if(this.editAdjustForm.transCount != this.newtransCount *1){
                    this.editAdjustForm.transCount = this.newtransCount *1;
                    this.api.adjust.edit(this.editAdjustForm)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "编辑成功!"
                            });
                            this.editAdjustForm = {};
                            this.getpurcon();
                        }
                    })
                }
            }else{
                this.$message.warning('运输发票吨数不能为空!');
            }
        },
        // 编辑运输发票吨数
        transCountColleclick(row){
            row.ysjustShow = true;
            this.newtransCount = row.transCount;
        },

        // 搜索功能
        search() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
            this.dialogVisibledata = false;
            this.AdjustForm = {};
        },
        // 搜索企业
        AdjustChang(){
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
            this.dialogVisibledata = false;
            this.AdjustForm = {};
        },
        // 删除
        delPurs(row){
            this.dialogVisibledata = false;
            this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
              this.api.adjust.del({id:row.id})
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
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });          
            });
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
            this.dialogVisibledata = false;
            this.AdjustForm = {};
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
            width: 19%;
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
// 新增弹窗样式开始
// ::v-deep .el-dialog__wrapper {
//     .el-dialog {
//         .el-dialog__body{
//         padding: 0px;
//         }
//         .el-dialog__header{
//         padding: 0px;
//         }
//     }
// }
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
    // height: 330px;
    padding: 0 120px 60px;
    box-sizing: border-box;
    position: relative;
    overflow: hidden;
}
// 新增弹窗样式结束
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
    display: inline-block;
    margin-right:24px;
  }
  .pur-right {
    .el-input {
      width: 30%;
      float: right;
      margin-right: 10px;
    }
    .el-button{
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
    position: relative;
    .addpurs{
        width:100% ;
        // height:53px;
        padding:8px 0 0 8px;
        box-sizing: border-box;
        background-color: #ffff;
        position: absolute;
        top:0px;
        left:0;
        z-index: 5;
        .AdjustForm{
            .el-form-item{
                float: left;
            }
        }
    }
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
      top:50%;
      left:50%;
      transform: translate(-50%, -50%);
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
</style>