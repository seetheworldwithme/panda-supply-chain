<template>
  <!-- 交易记录页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="addPurs">取消付款</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
        <el-date-picker
        size="small"
        @change="Times"
        v-model="TimeArrays"
        type="daterange"
        unlink-panels
        range-separator="~"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="yyyy-MM-dd"
        style="width:20%;float:right;margin-right:16px;">
        </el-date-picker>
        <!-- <el-date-picker
        @change="Times"
        size="small"
        v-model="keyupTime"
        type="date"
        value-format="yyyy-MM-dd"
        placeholder="选择日期">
        </el-date-picker> -->
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
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column type="selection" fixed="left"></el-table-column> 
            <el-table-column label="合同编号" sortable :show-overflow-tooltip="true" width="240">
                <template slot-scope="scope">
                    <div>
                        <span style="color: #cd1724;cursor: pointer;" @click="purSee(scope.row)">{{ scope.row.contractCode }}</span>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="交易信息" :show-overflow-tooltip="true" width="140">
                <template slot-scope="scope">
                    <div>
                        <span style="color: #cd1724;cursor: pointer;" @click="recordClick(scope.row)">详情</span>
                    </div>
                </template>
            </el-table-column>
            <el-table-column label="状态" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <div v-if="scope.row.status == 0">
                  <div class="statsBox" style="background: #CF1724;box-shadow: 0px 4px 4px #F7D0D3;"></div><span> 未付款</span>
                </div>
                <div v-if="scope.row.status == 1">
                  <div class="statsBox" style="background: #FB8138;box-shadow: 0px 4px 4px #FFE2D0;"></div><span> 已付款冻结中</span>
                </div>
                <div v-if="scope.row.status == 2">
                  <div class="statsBox"></div><span> 付款完成</span>
                </div>
                <div v-if="scope.row.status == 3">
                  <div class="statsBox" style="background: #FB8138;box-shadow: 0px 4px 4px #FFE2D0;"></div><span> 已取消付款</span>
                </div>
                <div v-if="scope.row.status == 4">
                  <div class="statsBox" style="background: #CF1724;box-shadow: 0px 4px 4px #F7D0D3;"></div><span> 付款失败</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="订单金额（元）" :show-overflow-tooltip="true" width="140" prop="freightAmount">
              <template slot-scope="scope">
                <span v-if="scope.row.payMoney">{{scope.row.payMoney|currency}} 元</span>
                <span v-else>0.00 元</span>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="操作时间" sortable :show-overflow-tooltip="true" width="180"> 
            </el-table-column>
            <el-table-column label="系统订单" prop="mchOrderNbr" sortable width="240" :show-overflow-tooltip="true">
            </el-table-column>
            <el-table-column prop="mchTransNbr" label="交易流水号" :show-overflow-tooltip="true" width="180">
            </el-table-column>
            <el-table-column label="支付方式" :show-overflow-tooltip="true" width="120"> 
              <template slot-scope="scope">
                <span v-if="scope.row.ordersubtype == 'GRT'">担保支付</span>
                <span v-if="scope.row.ordersubtype == 'DRT'">直接支付</span>
              </template>
            </el-table-column>
            <el-table-column prop="payeeUnitName" label="收款单位" :show-overflow-tooltip="true" width="240">
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template slot-scope="scope">
                <el-button v-if="scope.row.status == 1" size="mini" @click="gopayment(scope.row)">解 冻</el-button>
                <span v-else>--</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 抽屉 查看合同功能 -->
      <el-drawer :title="logo" :visible.sync="drawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <el-col :span="24">
                  <div class="see-top">基础信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>系统编码：</label>
                      <span>{{seeData.systemCode}}</span>
                    </li>
                    <li>
                      <label>项目类型：</label>
                      <span v-if="seeData.contractResource === 1">供应链业务</span>
                      <span v-if="seeData.contractResource === 2">供应链金融</span>
                    </li>
                    <li>
                      <label>合同编号：</label>
                      <span v-if="seeData.contractCode">{{seeData.contractCode}}</span>
                      <span v-else>---</span>
                    </li>
                    <li>
                      <label>合同名称：</label>
                      <span>{{seeData.contractName}}</span>
                    </li>
                    <li>
                      <label>签订时间：</label>
                      <span>{{seeData.signingTime}}</span>
                    </li>
                    <li>
                      <label>所属月份：</label>
                      <span v-if="seeData.month">{{seeData.month|timeDate}}</span>
                    </li>
                    <li>
                      <label>销售单位：</label>
                      <span>{{seeData.salesUnit}}</span>
                    </li>
                    <li>
                      <label>购货单位：</label>
                      <span>{{seeData.buyerUnit}}</span>
                    </li>
                    <li>
                      <label>发货地：</label>
                      <span>{{seeData.deliveryPlace}}</span>
                    </li>
                    <li>
                      <label>收货地：</label>
                      <span>{{seeData.receiptPlace}}</span>
                    </li>
                  </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">产品信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>产品大类：</label>
                      <span v-if="seeData.proBigType">{{seeData.proBigType|protype}}</span>
                    </li>
                    <li>
                      <label>产品小类：</label>
                      <span v-if="seeData.proSmaType">{{seeData.proSmaType|promiddletype}}</span>
                    </li>
                    <li>
                      <label>产品名称：</label>
                      <span v-if="seeData.proSmaType == 20">{{seeData.proRemark}}</span>
                      <span v-else>---</span>
                    </li>
                    <li>
                      <label>合同数量：</label>
                      <span>{{seeData.contractNum}} 吨</span>
                    </li>
                    <li>
                      <label>货品单价：</label>
                      <span>{{seeData.unitPrice}} 元 / 吨</span>
                    </li>
                    <li>
                      <label>总价金额：</label>
                      <span v-if="seeData.sumPrice">{{seeData.sumPrice|currency}} 元</span>
                    </li>
                    <li v-for="(v, index) in seeData.list" :key="index">
                      <label>标的指标：</label>
                      <div class="see-libox">
                        <span>硫分St.d% : </span>
                        <div>
                          <p v-if="seeData.sulfur">{{seeData.sulfur}} </p><p v-else>--- </p>
                        </div>
                        <span>挥发分Vdaf% : </span>
                        <div>
                          <span v-if="seeData.volatilization">{{seeData.volatilization}} </span><span v-else>--- </span>
                        </div>
                        <span>水分Mt% : </span>
                        <div>
                          <span v-if="seeData.water">{{seeData.water}} </span><span v-else>--- </span>
                        </div>
                        <span>{{v.dkey}} : </span>
                        <div>
                          <span v-if="v.dvalue">{{v.dvalue }} </span><span v-else>--- </span>
                        </div>
                      </div>
                    </li> 
                    <li style="clear:both;">
                      <label>备注：</label>
                      <span v-if="seeData.remark">{{seeData.remark}}</span>
                      <span v-else>---</span>
                    </li>
                  </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">扩展信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>价格类型：</label>
                      <span v-if="seeData.priceType">{{seeData.priceType|pricetype}}</span>
                    </li>
                    <li>
                      <label>付款方式：</label>
                      <span v-if="seeData.payMethod">{{seeData.payMethod|paymentmethod}}</span>
                    </li>
                    <li>
                      <label>支付方式：</label>
                      <span v-if="seeData.payType">{{seeData.payType|paymenttype}}</span>
                    </li>
                    <li>
                      <label>运输模式：</label>
                      <span v-if="seeData.transportMode">{{seeData.transportMode|transportmode}}</span>
                    </li>
                    <li>
                      <label>承运工具：</label>
                      <span v-if="seeData.carrier">{{seeData.carrier|carrier}}</span>
                    </li>
                  </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">扩展信息</div>
                  <dl class="see-file">
                    <dd>
                      <el-col :span="5">
                        <span>上传合同附件：</span>
                      </el-col>
                      <el-col :span="19">
                        <a v-for="(v,k) in filelist" :key="k"  :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"><p v-if="v.fileType== 1">{{v.name}}</p></a>
                      </el-col>
                    </dd>
                    <dd>
                      <el-col :span="5">
                        <span>上传补充协议：</span>
                      </el-col>
                      <el-col :span="19" style="padding-bottom: 48px;">
                        <a v-for="(v,k) in filelist" :key="k" :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"><p v-if="v.fileType== 2">{{v.name}}</p></a>
                      </el-col>
                    </dd>
                  </dl>
                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer>
      <!-- 查看交易记录功能 -->
      <el-drawer title="交易信息" :visible.sync="recorddrawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
                <el-row>
                <el-col :span="24">
                  <div class="see-top">交易信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>系统订单：</label>
                      <span>{{moenyData.mchOrderNbr}}</span>
                    </li>
                    <li>
                      <label>交易流水：</label>
                      <span>{{moenyData.mchTransNbr}}</span>
                    </li>
                    <li>
                      <label>订单金额：</label>
                      <span v-if="moenyData.orderAmt">{{moenyData.orderAmt|currency}} 元</span>
                    </li>
                  </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">付款方信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>付款单位：</label>
                      <span>{{moenyData.companyName}}</span>
                    </li>
                    <li>
                      <label>付款账号：</label>
                      <span>{{moenyData.payerAccNbr}}</span>
                    </li>
                    <li>
                      <label>付款账户：</label>
                      <span>{{moenyData.payerAccName}}</span>
                    </li>
                    <li>
                      <label>付款时间：</label>
                      <span v-if="moenyData.orderStrTime">{{moenyData.orderStrTime|orderStrTime}}</span>
                    </li>
                  </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">收款方信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>收款单位：</label>
                      <span>{{moenyData.payeeUnitName}}</span>
                    </li>
                    <li>
                      <label>收款账号：</label>
                      <span>{{moenyData.payeeAccNbr}}</span>
                    </li>
                    <li>
                      <label>收款账户：</label>
                      <span>{{moenyData.payeeAccName}}</span>
                    </li>
                  </ul>
                </el-col>
              </el-row>
          <div class="xia">《</div>
        </div>
      </el-drawer>
      <!-- 分页 -->
      <div class="pur-pages">
        <el-pagination
        background
        :page-sizes="[50, 100, 200, 300, 500]"
        :page-size.sync="size"
        :pager-count="5"
        :current-page.sync="current"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div>
      <!-- 解冻审批弹窗 -->
      <el-dialog
        title="解冻操作"
        :visible.sync="paymentShow"
        width="35%"
        :before-close="paymentShowClose">
          <el-form style="width:95%;" :model="AddData" :rules="paymentrules" ref="AddData" size="small" label-width="100px">
            <el-form-item label="操作人" prop="name">
              <el-input v-model="AddData.name" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="解冻金额" prop="transAmt">
              <el-input v-model="AddData.transAmt" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="系统订单" prop="mchOrderNbr">
              <el-input v-model="AddData.mchOrderNbr" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="交易流水号" prop="mchTransNbr">
              <el-input v-model="AddData.mchTransNbr" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item prop="code" label="验证码">
              <el-input v-model="AddData.code" placeholder="请输入验证码">
                <template slot="suffix">
                  <span v-if="!isDisabled" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:3px;cursor:pointer;"></span>
                  <span v-if="isDisabled" :class="color" v-html="codephone" style="cursor: pointer;" ></span>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="paymentsubmitm('AddData')" style="float:right;margin:5% 0 0 2%;">解 冻</el-button>
              <el-button @click="paymentresetForm" style="float:right;margin:5% 0 0 2%;">取 消</el-button>
            </el-form-item>
          </el-form>
      </el-dialog>
      <!-- 退款弹窗 -->
      <el-dialog
        title="退款信息"
        :visible.sync="refundShow"
        width="35%"
        :before-close="refundClose">
          <el-form style="width:95%;" :model="refundData" :rules="refundrules" ref="refundData" size="small" label-width="100px">
            <el-form-item label="操作人" prop="name">
              <el-input v-model="refundData.name" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="退款金额" prop="transAmt">
              <el-input v-model="refundData.transAmt" :disabled="true"></el-input>
            </el-form-item>
            <!-- <el-form-item label="退款账号" prop="orderAmt">
              <el-input v-model="refundData.orderAmt"></el-input>
            </el-form-item> -->
            <el-form-item prop="code" label="验证码">
              <el-input v-model="refundData.code" placeholder="请输入验证码">
                <template slot="suffix">
                  <span v-if="!refundisDisa" :class="refundcolor" @click="refundphoneClick" v-html="refundcodephone" style="margin-right:3px;cursor:pointer;"></span>
                  <span v-if="refundisDisa" :class="refundcolor" v-html="refundcodephone" style="cursor: pointer;" ></span>
                </template>
              </el-input>
            </el-form-item>
            <!-- <el-form-item label="退款原因" prop="payeeEtc">
              <el-input v-model="refundData.payeeEtc" :disabled="true"></el-input>
            </el-form-item> -->
            <el-form-item>
              <el-button type="primary" @click="refundsubmitm('refundData')" style="float:right;margin:5% 0 0 2%;">退 款</el-button>
              <el-button @click="refundesetForm" style="float:right;margin:5% 0 0 2%;">取 消</el-button>
            </el-form-item>
          </el-form>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            purSearch: "",          // 搜索  
            TimeArrays:[],           // 时间搜索 
            dataAva:false,
            loading:false,
            current:1,              // 当前页数
            size:200,               // 显示条数
            total:0,                // 总条数
            show:false,
            purTitle: "",           // 标题
            drawer: false,          // 查看合同抽屉
            recorddrawer:false,     // 查看交易记录
            logo: "",               // 抽屉企业名称
            activeName: "first",    // 查看tab栏
            tableData: [],          // 表格数据
            seeData:{},             // 查看合同数据
            filelist:[],            // 查看文件列表 
            // 查看...动态
            timeData:[
                // {
                //   createtime:'2021-01-02',
                //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ], 
            moenyData:{},   // 交易详情

            // 退款弹窗
            refundShow:false,
            refundData:{},
            refundrules: {
              code: [
                { required: true, message: '请输入验证码', trigger: 'blur' },
                { max: 6, message: '请输入6位验证码', trigger: 'blur' }
              ],
                // orderAmt: [
                //   { required: true, message: '请输入合同订单金额', trigger: 'change' },
                //   { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额格式不正确', trigger: 'blur'},
                // ],
            },
            // 验证码
            caleardis:null,
            refundisDisa: false, //控制按钮是否可以点击（false:可以点击，true:不可点击）
            refundcolor:'onPhone',
            refundcodephone:'获取验证码',
            delmul:null,              // 批量退款总金额
            refundList:[],            // 选中集合
            multipleSelection: [],    // 选中数据

            // 解冻审批弹窗
            paymentShow:false,
            AddData:{
              code:''
            },
            paymentrules: {
              code: [
                { required: true, message: '请输入验证码', trigger: 'blur' },
                { max: 6, message: '请输入6位验证码', trigger: 'blur' }
              ],
                // orderAmt: [
                //   { required: true, message: '请输入合同订单金额', trigger: 'change' },
                //   { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额格式不正确', trigger: 'blur'},
                // ],
            },
            // 验证码
            calearfin:null,
            isDisabled: false, //控制按钮是否可以点击（false:可以点击，true:不可点击）
            color:'onPhone',
            codephone:'获取验证码',
        };
    },
    created() {},
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName')     // 企业名称
      this.http  = this.$store.state.http;                  // IP地址
      this.getpurcon();
    },
    methods: {
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                transType:0,
                searchStr:this.purSearch,
                startTime:this.TimeArrays != null?this.TimeArrays[0]:"",
                endTime:this.TimeArrays != null?this.TimeArrays[1]:"",
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.Transinfo.all(listByPageData)
            .then(res=>{
              if(res.data.code == 200){
                  this.loading = false; 
                  if(res.data.data.records != ""){
                    this.tableData = res.data.data.records;
                    this.total = res.data.data.total;
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
        // 日期搜索
        Times() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 查看合同
        purSee(row) {
            this.drawer = true;
            this.activeName = 'first';
            // this.seeId = row.contractId;
            this.api.cgContract.one({id:row.contractId})
            .then(res=>{
                if(res.data.code == 200){
                this.seeData = res.data.data;
                }
            })
            // 文件
            this.api.file.one({relationId:row.contractId})
            .then(res=>{
                this.filelist = res.data.data;
            })
        },
        changeActive(e) {
            e.currentTarget.className = 'active';
        },
        removeActive(e) {
            e.currentTarget.className = '';
        },
        // 查看交易记录详情
        recordClick(row) {
            this.recorddrawer = true;
            // this.moenyData
            this.api.confirmPay.orderDetails({mchOrdNbr:row.mchOrderNbr,mchTransNbr:row.mchTransNbr,mchNbr:row.mchNbr})
            .then(res=>{
                if(res.data.code == 200){
                    this.moenyData = res.data.data.body.orderInfoQryRtnDto;
                    this.moenyData.payeeUnitName = row.payeeUnitName;
                    this.moenyData.companyName = row.payUnitName;
                }
            })
        },

        // 抽屉
        handleClick(tab) {
          // console.log(tab)
          // 判断tab栏选中状态，请求相应的接口
          if (tab.label == "基本信息") {
            // this.allone(this.seeId);
          }
          if (tab.label == "动态") {
            // console.log(2);
          }
        },

        // 解冻审批确认按钮
        gopayment(data){
          this.paymentShow = true;
          this.AddData.id = data.id;
          this.AddData.mchOrderNbr = data.mchOrderNbr;
          this.AddData.mchTransNbr = data.mchTransNbr;
          this.AddData.transAmt = data.payMoney;
          this.AddData.mchNbr = data.mchNbr;
          this.AddData.name = sessionStorage.getItem('username');
          this.codephone="获取验证码";
          this.color="onPhone";
          this.isDisabled = false;
          clearInterval(this.calearfin);
        },
        // 解冻审批取消
        paymentShowClose(){
          this.paymentShow = false;
          this.AddData = {};
          this.codephone="获取验证码";
          this.color="onPhone";
          this.isDisabled = false;
          clearInterval(this.calearfin);
        },
        // 解冻审批信息提交
        paymentsubmitm(formName) {
          this.$refs[formName].validate((valid) => {
            if (valid) {
              this.AddData.transType = 'CFM';
              this.api.confirmPay.GoConfirmPay(this.AddData)
              .then(res=>{
                if(res.data.code == 200){
                  this.$message({
                    type: 'success',
                    message: '操作成功，请稍后查看!'
                  });
                  this.tableData = [];
                  this.getpurcon();
                  this.current = 1;
                  this.paymentShow = false;
                  this.AddData = {};
                  this.codephone="获取验证码";
                  this.color="onPhone";
                  this.isDisabled = false;
                  clearInterval(this.calearfin);
                }
              })
            }
          });
        },
        // 取消解冻审批
        paymentresetForm() {
          this.paymentShow = false;
          this.AddData = {};
          this.codephone="获取验证码";
          this.color="onPhone";
          this.isDisabled = false;
          clearInterval(this.calearfin);
        },
        // 解冻审批点击获取手机验证码
        phoneClick(){
          this.isDisabled = true;
          let time = 60;
          this.api.code.refundMoney({templateId:'1358420'})
          .then(res=>{
            let that = this;
            if(res.data.code == 200){
              this.calearfin = setInterval(function(){
                let newTime = time - 1;
                that.codephone = '重新获取('+ newTime +')s'
                that.color="onokPhone"
                  time--
                  if(time == 0){
                    that.codephone="获取验证码"
                    that.color="onPhone"
                    clearInterval(this.calearfin);
                    that.isDisabled = false;
                  }
              },1000)  
            }
            if(res.data.code == 999){
              that.isDisabled = false;
            }
          })
        },

        // 退款按钮
        addPurs(){
          if(this.multipleSelection.length == 0){
            this.$message({
                type: 'warning',
                message: '请选择要取消付款的合同!'
            });
          }else{
            let switchs = false;
            for(let i=0;i<this.refundList.length;i++){
              if(this.refundList[i].status != 1){
                switchs = true;
              }
            }
            if(switchs){
                this.$message({
                  message: '发起失败，请选择已付款冻结的订单！',
                  type: 'warning'
                });
            }else{
              this.refundShow = true;
              this.refundData.name = sessionStorage.getItem('username');
              this.refundData.transAmt = this.delmul;
              switchs = false;
              this.refundcodephone="获取验证码"
              this.refundcolor="onPhone"
              clearInterval(this.caleardis);
              this.refundisDisa = false;
            }
          }
        },
        // 退款取消
        refundClose(){
          this.refundShow = false;
          this.refundData = {};
          this.refundcodephone="获取验证码"
          this.refundcolor="onPhone"
          clearInterval(this.caleardis);
          this.refundisDisa = false;
        },
        // 退款信息提交
        refundsubmitm(formName) {
          this.$refs[formName].validate((valid) => {
            if (valid) {
              if(this.refundData.orderAmt == 0){
                this.$message({
                  message: '付款金额必须大于零',
                  type: 'warning'
                });
              }else{
                console.log(this.refundList)
                this.api.confirmPay.cancelPay({refundList:this.refundList,code:this.refundData.code})
                .then(res=>{
                  if(res.data.code == 200){
                    this.$message({
                      type: 'success',
                      message: '退款成功，请稍后查看!'
                    });
                    this.tableData = [];
                    this.getpurcon();
                    this.current = 1;
                    this.refundShow = false;
                    this.refundData = {};
                    this.refundcodephone="获取验证码";
                    this.refundcolor="onPhone";
                    clearInterval(this.caleardis);
                    this.refundisDisa = false;
                  }
                })
              }
            }else{
              this.AddData.orderAmt = '';
            }
          });
        },
        // 取消退款
        refundesetForm() {
          this.refundShow = false;
          this.refundData = {};
          this.refundcodephone="获取验证码"
          this.refundcolor="onPhone"
          clearInterval(this.caleardis);
          this.refundisDisa = false;
        },
        // 退款点击获取手机验证码
        refundphoneClick(){
            this.refundisDisa = true;
            let time = 60;
            this.api.code.refundMoney({templateId:'1358422'})
            .then(res=>{
              let that = this;
              if(res.data.code == 200){
                this.caleardis = setInterval(function(){
                  let newTime = time - 1;
                  that.refundcodephone = '重新获取('+ newTime +')s'
                  that.refundcolor="onokPhone"
                    time--
                    if(time == 0){
                      that.refundcodephone="获取验证码"
                      that.refundcolor="onPhone"
                      clearInterval(this.caleardis);
                      that.refundisDisa = false;
                    }
                },1000)  
              }
              if(res.data.code == 999){
                that.refundisDisa = false;
              }
            })
        },
        // 表格选中数据
        purChange(val) {
          this.multipleSelection = val;
          this.delmul = null;
          this.refundList = [];
          var smul = 0;
          for(let i=0;i<this.multipleSelection.length;i++){
            smul += this.multipleSelection[i].payMoney*1 > 0? this.multipleSelection[i].payMoney*1:0;
            this.refundList.push({
              id:this.multipleSelection[i].id,
              mchOrderNbr:this.multipleSelection[i].mchOrderNbr,
              mchTransNbr:this.multipleSelection[i].mchTransNbr,
              transAmt:this.multipleSelection[i].payMoney,
              mchNbr:this.multipleSelection[i].mchNbr,
              status:this.multipleSelection[i].status,
              payerBankFlag:this.multipleSelection[i].payerBankFlag,
              name:sessionStorage.getItem('username'),
              transType:'CNL'
            })
          }
          this.delmul = smul;
        },
    },
};
</script>
<style lang="less" scoped>
// 订单状态样式
.statsBox{
  width: 10px;
  height: 10px;
  border-radius:50%;
  display: inline-block;
  background: #6FCF97;
  box-shadow: 0px 4px 4px #CBF9DE;
  span{
    color: #333;
  }
}
// 确认付款样式
.paymentsize{
  padding-top:6px;
  box-sizing:border-box;
  font-size: 12px;
  line-height: 22px;
  color: #333;
}
.onPhone{
  font-size: 12px;
  color: #EA222E;
}
.onokPhone{
  font-size: 12px;
  color: #BFBFBF;
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
  padding:24px 0 24px 24px;
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
  padding: 15px 0;
  label {
    display: inline-block;
    width: 18%;
    font-size: 14px;
    color: #666;
  }
  li {
    min-height: 40px;
    font-size: 14px;
    line-height: 40px;
    color: #333;
    label{
      float: left;
    }
    .see-libox{
      width: 76%;
      display: inline-block;
      span{
        color: #666;
      }
      div{
        display: inline-block;
        margin-right: 15px;
      }
    }
  }

}
// 文件样式
.see-file{
  padding: 15px 0 ;
  span {
    font-size: 14px;
    color: #666;
  }
  dd {
    display: flow-root;
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