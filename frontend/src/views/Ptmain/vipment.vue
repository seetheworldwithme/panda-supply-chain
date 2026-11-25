<!-- 标签列表页 -->
<template>
  <!-- 云仓管理页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3">
        <span class="pur-size">{{ purTitle }}</span></el-col
      >
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
      <!-- 表格 -->
      <div class="pur-table">
        <el-table
          height="100%"
          border
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%"
        >
          <el-table-column prop="id" v-if="show"></el-table-column>
          <el-table-column
            label="操作"
            width="200"
            fixed="left"
            v-if="$route.meta.edit"
          >
            <template slot-scope="scope">
              <el-button
                size="mini"
                @click="addPurs(scope.row.memberType)"
                v-if="scope.row.memberType != 2"
                >权限配置</el-button>
              <el-button
                v-if="scope.row.memberType == 1"
                type="primary"
                @click="EditPurs(scope.row)"
                size="mini"
                >会员配置</el-button>
                <el-button
                v-if="scope.row.memberType == 3"
                type="primary"
                @click="seleEditPurs(scope.row)"
                size="mini"
                >会员配置</el-button>
            </template>
          </el-table-column>
          <el-table-column label="会员类型" prop="signName">
            <template slot-scope="scope">
              <span v-if="scope.row.memberType == 0">免费会员</span>
              <span v-if="scope.row.memberType == 1">充值会员</span>
              <span v-if="scope.row.memberType == 2">预付定金会员</span>
            </template>
          </el-table-column>
          <el-table-column label="年度会员" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.yearPrice"
                >{{ scope.row.yearPrice }} 元/年</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="月度会员" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.monthPrice"
                >{{ scope.row.monthPrice }} 元/月</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="预付保证金" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.advanceDeposit">{{
                scope.row.advanceDeposit
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="黑金劵" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.needScoreNum && scope.row.needMoney"
                >{{ scope.row.needScoreNum }} 积分 =
                {{ scope.row.needMoney }} 元</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="当月供应名额" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.gyMonthlyQuota"
                >{{ scope.row.gyMonthlyQuota }} 个</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="当月采购名额" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.cgMonthlyQuota"
                >{{ scope.row.cgMonthlyQuota }} 个</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="年度最终会员价格" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.finalYearPrice"
                >{{ scope.row.finalYearPrice }} 年/元</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="半年度最终会员价格" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.finalHalfYearPrice"
                >{{ scope.row.finalHalfYearPrice }} 半年/元</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="季度最终会员价格" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.finalQuarterYearPrice"
                >{{ scope.row.finalQuarterYearPrice }} 季度/元</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column label="月度最终会员价格" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.finalMonthPrice"
                >{{ scope.row.finalMonthPrice }} 月/元</span
              >
              <span v-else>--</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <!-- 权限配置弹窗 -->
      <el-dialog
        :closeOnClickModal="false"
        title="权限配置"
        :visible.sync="paymentShow"
        width="30%"
        :before-close="paymentShowClose"
      >
        <el-form
          style="width: 95%"
          :model="AddData"
          :rules="paymentrules"
          ref="AddData"
          size="small">
          <el-form-item>
            <div>煤种发布用户:</div>
            <el-checkbox v-model="AddData.isShowEtcPhone">企业联系人</el-checkbox>
            <el-checkbox v-model="AddData.isShowEtcName">企业名称</el-checkbox>
            <!-- <el-checkbox v-model="AddData.isShowWebsite">企业数字化官网</el-checkbox> -->
            <el-checkbox v-model="AddData.isSeeTestingReport">检测报告</el-checkbox>
            <el-checkbox v-model="AddData.isSeePriceTrend">价格走势</el-checkbox>
            <div>登录用户:</div>
            <el-checkbox v-model="AddData.isSeeWarehouseVisitor">我的云仓访客</el-checkbox>
            <el-checkbox v-model="AddData.isSeeDataBriefing">数据简报</el-checkbox>
            <el-checkbox v-model="AddData.isSeeCardVisitor">名片访客</el-checkbox>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="paymentsubmitm('AddData')"
              style="float: right; margin: 5% 0 0 2%"
              >确 定</el-button
            >
            <el-button
              @click="paymentresetForm"
              style="float: right; margin: 5% 0 0 2%"
              >取 消</el-button
            >
          </el-form-item>
        </el-form>
      </el-dialog>
      <!-- 会员配置弹窗 -->
      <el-dialog
        title="会员配置"
        :closeOnClickModal="false"
        :visible.sync="editpaymentShow"
        width="40%"
        :before-close="editpaymentShowClose">
        <el-form
          style="width: 95%"
          :model="editData"
          :rules="editpaymentrules"
          ref="editData"
          size="small"
          label-width="100px">
          <el-row>
            <el-col :span="8">
              <el-form-item label="年度会员">
                <el-input
                  v-model="editData.yearPrice"
                  placeholder="请输入价格"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="月度会员">
                <el-input
                  v-model="editData.monthPrice"
                  placeholder="请输入价格"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="预付保证金">
                <el-input
                  v-model="editData.advanceDeposit"
                  placeholder="请输入价格"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="黑金劵">
                <el-col :span="5">
                  <el-input
                    v-model="editData.needScoreNum"
                    placeholder="请输入黑金劵"
                  ></el-input>
                </el-col>
                <el-col :span="1" style="padding-left: 1%">
                  <span>--</span>
                </el-col>
                <el-col :span="5">
                  <el-input
                    v-model="editData.needMoney"
                    placeholder="请输入价格"
                  ></el-input>
                </el-col>
                <el-col :span="2" style="padding-left: 1%">
                  <span> 元</span>
                </el-col>
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="当月供应名额">
                <el-col :span="20">
                  <el-input
                    v-model="editData.gyMonthlyQuota"
                    placeholder="请输入数量"
                  ></el-input>
                </el-col>
                <el-col :span="4" style="padding-left: 1%">个</el-col>
              </el-form-item>
            </el-col>
            <el-col :span="11">
              <el-form-item label="当月采购名额">
                <el-col :span="20">
                  <el-input
                    v-model="editData.cgMonthlyQuota"
                    placeholder="请输入数量"
                  ></el-input>
                </el-col>
                <el-col :span="4" style="padding-left: 1%">个</el-col>
              </el-form-item>
            </el-col>
          </el-row>
          <!-- 轮播图上传 -->
          <el-row>
            <el-col :span="8">
              <el-col :span="24">
                <div style="height: 30px">
                  <label style="color: #666; font-size: 14px; line-height: 30px"
                    >供应会员</label>
                </div>
              </el-col>
              <el-form-item label-width="0px" prop="banner1">
                <el-col :span="24">
                  <el-upload
                    :limit="1"
                    :action="http + '/saas-common/upload/batch'"
                    :class="yeBtnDealImg ? 'disyeloadSty' : ''"
                    :headers="header"
                    :file-list="listyearr"
                    name="file"
                    list-type="picture-card"
                    :on-success="addyelist"
                    :on-preview="jcyeardPrediv"
                    :on-remove="yeRemove">
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="imgyeVisible" append-to-body>
                    <img width="100%" :src="yeImgUrl" alt="" />
                  </el-dialog>
                </el-col>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-col :span="24">
                <label style="color: #666; font-size: 14px; line-height: 30px"
                  >采购会员</label
                >
              </el-col>
              <el-form-item label-width="0px" prop="banner2">
                <el-col :span="24">
                  <el-upload
                    :limit="1"
                    :action="http + '/saas-common/upload/batch'"
                    :class="mpBtnDealImg ? 'dismploadSty' : ''"
                    :headers="header"
                    :file-list="listmparr"
                    name="file"
                    list-type="picture-card"
                    :on-success="addmplist"
                    :on-preview="mpPictureCardPrediv"
                    :on-remove="mpRemove"
                  >
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="imgmpogVisible" append-to-body>
                    <img width="100%" :src="mpogImageUrl" alt="" />
                  </el-dialog>
                </el-col>
              </el-form-item>
            </el-col>
            <el-col :span="8"> 
              <el-col :span="24">
                <label style="color: #666; font-size: 14px; line-height: 30px"
                  >代理销售</label
                >
              </el-col>
              <el-form-item label-width="0px" prop="banner3">
                <el-col :span="24">
                  <el-upload
                    :limit="1"
                    :action="http + '/saas-common/upload/batch'"
                    :class="jcBtnDealImg ? 'disjcloadSty' : ''"
                    :headers="header"
                    :file-list="listjcarr"
                    name="file"
                    list-type="picture-card"
                    :on-success="addjclist"
                    :on-preview="jcPictureCardPrediv"
                    :on-remove="jcRemove"
                  >
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="imgjcogVisible" append-to-body>
                    <img width="100%" :src="jcogImageUrl" alt="" />
                  </el-dialog>
                </el-col>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item>
            <el-button
              type="primary"
              @click="editpaymentsubmitm('editData')"
              style="float: right; margin: 5% 0 0 2%"
              >确定</el-button
            >
            <el-button
              @click="editpaymentresetForm"
              style="float: right; margin: 5% 0 0 2%"
              >取消</el-button
            >
          </el-form-item>
        </el-form>
      </el-dialog>


      <!-- 熊猫甄选会员配置弹窗 -->
      <el-dialog
        title="熊猫甄选会员配置"
        :closeOnClickModal="false"
        :visible.sync="selementShow"
        width="60%"
        :before-close="selseShowClose">
        <el-form
          :model="seleData"
          :rules="selerules"
          ref="seleData"
          size="small"
          label-width="160px">
          <el-row>
            <!-- 初始会员价格 -->
            <el-col :span="6">
              <el-form-item label="月度会员">
                <el-input v-model="seleData.monthPrice" placeholder="请输入价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="季度会员">
                <el-input v-model="seleData.quarterYearPrice" placeholder="请输入价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="半年度会员">
                <el-input v-model="seleData.halfYearPrice" placeholder="请输入价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="年度会员">
                <el-input v-model="seleData.yearPrice" placeholder="请输入价格"></el-input>
              </el-form-item>
            </el-col>
            <!-- 增值系数（%） -->
            <el-col :span="6">
              <el-form-item label="月度增值系数（%）">
                <el-input v-model="seleData.monthCoefficient" placeholder="请输入会员增值系数（%）"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="季度增值系数（%）">
                <el-input v-model="seleData.quarterYearCoefficient" placeholder="请输入会员增值系数（%）"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="半年度增值系数（%）">
                <el-input v-model="seleData.halfYearCoefficient" placeholder="请输入会员增值系数（%）"></el-input>
              </el-form-item>
            </el-col>
            <!-- 折扣（%） -->
            <el-col :span="6">
              <el-form-item label="月度折扣（%）">
                <el-input v-model="seleData.monthDiscount" placeholder="请输入会员折扣（%）"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="季度折扣（%）">
                <el-input v-model="seleData.quarterYearDiscount" placeholder="请输入会员折扣（%）"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="半年度折扣（%）">
                <el-input v-model="seleData.halfYearDiscount" placeholder="请输入会员折扣（%）"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="年度折扣（%）">
                <el-input v-model="seleData.yearDiscount" placeholder="请输入会员折扣（%）"></el-input>
              </el-form-item>
            </el-col>
            <!-- 最终价格 -->
            <el-col :span="6">
              <el-form-item label="月度最终价格">
                <el-input @focus="ydprice" v-model="seleData.finalMonthPrice" placeholder="请输入会员最终价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="季度最终价格">
                <el-input @focus="jdprice" v-model="seleData.finalQuarterYearPrice" placeholder="请输入会员最终价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="半年度最终价格">
                <el-input @focus="bndprice" v-model="seleData.finalHalfYearPrice" placeholder="请输入会员最终价格"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="年度最终价格">
                <el-input @focus="ndprice" v-model="seleData.finalYearPrice" placeholder="请输入会员最终价格"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <!-- 轮播图上传 -->
          <el-row>
            <el-col :span="8">
              <el-col :span="24">
                <div style="height: 30px">
                  <label style="color: #666; font-size: 14px; line-height: 30px"
                    >甄选会员</label
                  >
                </div>
              </el-col>
              <el-form-item label-width="0px" prop="banner4">
                <el-col :span="24">
                  <el-upload
                    :limit="1"
                    :action="http + '/saas-common/upload/batch'"
                    :class="seleDealImg ? 'disyeloadSty' : ''"
                    :headers="header"
                    :file-list="selelistyear"
                    name="file"
                    list-type="picture-card"
                    :on-success="addselelist"
                    :on-preview="seleyeardPrediv"
                    :on-remove="seleRemove">
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="seleVisible" append-to-body>
                    <img width="100%" :src="seleImgUrl" alt="" />
                  </el-dialog>
                </el-col>
              </el-form-item>
            </el-col>
            <el-col :span="16">
              <div>
                <h1 style="font-weight:600;margin-bottom:12px;">最终会员价计算规则:</h1>
                <h4>月度会员最终价 = (月度会员价格/12 + (月度会员价格/12)*月度增值系数)*月度折扣;</h4>
                <h4>季度会员最终价 = (季度会员价格/4 + (季度会员价格/4)*季度增值系数)*季度折扣;</h4>
                <h4>半年度会员最终价 = (半年度会员价格/2 + (半年度会员价格/2)*半年度增值系数)*半年度折扣;</h4>
                <h4>年度会员最终价 = 月年度会员价格*年度折扣;</h4>
              </div>
            </el-col>
          </el-row>
          <el-form-item>
            <el-button
              type="primary"
              @click="selesubmitm('seleData')"
              style="float: right; margin: 5% 0 0 2%"
              >确定</el-button>
            <el-button
              @click="seleresetForm"
              style="float: right; margin: 5% 0 0 2%"
              >取消</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
  <script>
import { number } from 'echarts';
export default {
  data() {
    return {
      // 甄选会员配置编辑
      selementShow: false,
      seleData: {
        banner4:'',
      },
      selerules: {
        banner4: [
          { required: true, message: "请上传banner4", trigger: "blur" },
        ],
      },
      // 甄选banner
      seleImgUrl: "",
      seleVisible: false,
      selelistyear: [],
      seleDealImg: false,


      show: false,
      purTitle: "", // 标题
      tableData: [], // 表格数据
      // 权限配置弹窗
      paymentShow: false,
      AddData: {},
      paymentrules: {
        signName: [
          { required: true, message: "请输入标签内容", trigger: "blur" },
        ],
      },
      radio: "1",
      checkList: [],
      // 会员配置编辑
      editpaymentShow: false,
      editData: {
        // banner1:'',
      },
      editpaymentrules: {
        signName: [
          { required: true, message: "请输入标签内容", trigger: "blur" },
        ],
        banner1: [
          { required: true, message: "请上传banner1", trigger: "blur" },
        ],
        banner2: [
          { required: true, message: "请上传banner2", trigger: "blur" },
        ],
        banner3: [
          { required: true, message: "请上传banner3", trigger: "blur" },
        ],
      },
      // banner
      http: "",
      header: {}, // 上传token
      // banner1
      yeImgUrl: "",
      imgyeVisible: false,
      listyearr: [],
      yeBtnDealImg: false,
      // banner2
      mpogImageUrl: "",
      imgmpogVisible: false,
      listmparr: [],
      mpBtnDealImg: false,
      // banner3
      jcogImageUrl: "",
      imgjcogVisible: false,
      listjcarr: [],
      jcBtnDealImg: false,
    };
  },
  created() {},
  mounted() {
    this.header.Authorization = sessionStorage.getItem("Authorization");
    this.http = this.$store.state.http; // IP地址
    this.purTitle = this.$route.meta.title;
    this.getpurcon();
  },
  methods: {
    /*
    最终会员价算法
    */
    // 年度最终会员计算
    ndprice(){
      if(this.seleData.yearDiscount && this.seleData.yearPrice){
        this.seleData.finalYearPrice = ((this.seleData.yearDiscount*this.seleData.yearPrice)*1/100).toFixed(0);
      }else{
        this.$message({
            type: "warning",
            message: "年度会员价格和折扣必填!",
          });
      }
    },
    // 半年度最终会员计算
    bndprice(){
      if(this.seleData.halfYearPrice && this.seleData.halfYearDiscount && this.seleData.halfYearCoefficient){
        let A = (this.seleData.halfYearPrice/2);
        let B = A*this.seleData.halfYearCoefficient/100;
        this.seleData.finalHalfYearPrice = ((A + B)*(this.seleData.halfYearDiscount/100)).toFixed(0);
        this.$forceUpdate();
        // console.log(((A + B)*(this.seleData.halfYearDiscount/100)).toFixed(0),this.seleData.finalHalfYearPrice)
      }else{
        this.$message({
            type: "warning",
            message: "半年度会员价格,半年度折扣,半年度增值系数必填!",
          });
      }
    },
    // 季度最终会员计算
    jdprice(){
      if(this.seleData.quarterYearPrice && this.seleData.quarterYearDiscount && this.seleData.quarterYearCoefficient){
        let A = (this.seleData.quarterYearPrice/4);
        let B = A*this.seleData.quarterYearCoefficient/100;
        this.seleData.finalQuarterYearPrice = ((A + B)*(this.seleData.quarterYearDiscount/100)).toFixed(0);
        this.$forceUpdate();
      }else{
        this.$message({
            type: "warning",
            message: "季度会员价格,季度折扣,季度增值系数必填!",
          });
      }
    },
    // 月度最终会员计算
    ydprice(){
      if(this.seleData.monthPrice && this.seleData.monthDiscount && this.seleData.monthCoefficient){
        let A = (this.seleData.monthPrice/12);
        let B = A*this.seleData.monthCoefficient/100;
        this.seleData.finalMonthPrice = ((A + B)*(this.seleData.monthDiscount/100)).toFixed(0);
        this.$forceUpdate();
      }else{
        this.$message({
            type: "warning",
            message: "月度会员价格,月度折扣,月度增值系数必填!",
          });
      }
    },
    /*
        上传功能
        */
    // 上传banner1
    addyelist(response, file, fileList) {
      if (response.code == 200) {
        this.editData.banner1 = response.data[0].id;
        this.yeBtnDealImg = fileList.length >= 1;
        this.$message({
          message: "banner1上传成功！",
          type: "success",
        });
      }
    },
    // 预览banner1
    jcyeardPrediv(file) {
      this.yeImgUrl = file.url;
      this.imgyeVisible = true;
    },
    // 保存成功后删除被移除的旧 banner 文件
    delRemovedFiles() {
      (this.removedFileIds || []).forEach((id) => this.api.addfile.del({ id }));
      this.removedFileIds = [];
    },
    // 删除banner1
    yeRemove(file, fileList) {
        // 先只从界面移除并记下旧文件,保存成功后再删:原先一点 × 就删服务器文件,随后点「取消」轮播图就裂了
        (this.removedFileIds = this.removedFileIds || []).push(this.editData.banner1);
                this.yeBtnDealImg = fileList.length >= 1;
                this.editData.banner1 = "";
                this.listyearr = [];
    },
    // 上传banner2
    addmplist(response, file, fileList) {
      if (response.code == 200) {
        this.editData.banner2 = response.data[0].id;
        this.mpBtnDealImg = fileList.length >= 1;
        this.$message({
          message: "banner2上传成功！",
          type: "success",
        });
      }
    },
    // 预览banner2
    mpPictureCardPrediv(file) {
      this.mpogImageUrl = file.url;
      this.imgmpogVisible = true;
    },
    // 删除banner2
    mpRemove(file, fileList) {
      // 先只从界面移除并记下旧文件,保存成功后再删:原先一点 × 就删服务器文件,随后点「取消」轮播图就裂了
      (this.removedFileIds = this.removedFileIds || []).push(this.editData.banner2);
                this.editData.banner2 = "";
                this.listmparr = [];
                this.mpBtnDealImg = fileList.length >= 1 ? true : false;
    },
    // 上传banner3
    addjclist(response, file, fileList) {
      if (response.code == 200) {
        this.editData.banner3 = response.data[0].id;
        this.jcBtnDealImg = fileList.length >= 1;
        this.$message({
          message: "banner3上传成功！",
          type: "success",
        });
      }
    },
    // 预览banner3
    jcPictureCardPrediv(file) {
      this.jcogImageUrl = file.url;
      this.imgjcogVisible = true;
    },
    // 删除banner3
    jcRemove(file, fileList) {
      // 先只从界面移除并记下旧文件,保存成功后再删:原先一点 × 就删服务器文件,随后点「取消」轮播图就裂了
      (this.removedFileIds = this.removedFileIds || []).push(this.editData.banner3);
                this.jcBtnDealImg = fileList.length >= 1;
                this.editData.banner3 = "";
                this.listjcarr = [];
    },
    // 上传banner4
    addselelist(response, file, fileList) {
      if (response.code == 200) {
        this.seleData.banner4 = response.data[0].id;
        this.seleDealImg = fileList.length >= 1;
        this.$message({
          message: "banner4上传成功！",
          type: "success",
        });
      }
    },
    // 预览banner4
    seleyeardPrediv(file) {
      this.seleImgUrl = file.url;
      this.seleVisible = true;
    },
    // 删除banner4
    seleRemove(file, fileList) {
        // 先只从界面移除并记下旧文件,保存成功后再删:原先一点 × 就删服务器文件,随后点「取消」轮播图就裂了
        (this.removedFileIds = this.removedFileIds || []).push(this.seleData.banner4);
                this.seleDealImg = fileList.length >= 1;
                this.seleData.banner4 = "";
                this.selelistyear = [];
    },
    // 表格数据
    getpurcon() {
      // 渲染表格
      this.api.busconfig.list().then((res) => {
        if (res.data.code == 200) {
          this.tableData = res.data.data;
        }
      });
    },

    // 权限配置信息弹窗
    addPurs(type) {
      this.api.busbenfig.list({ member_type: type }).then((res) => {
        if (res.data.code == 200) {
          this.paymentShow = true;
          this.AddData = res.data.data[0];
        }
      });
    },
    // 权限配置取消
    paymentShowClose() {
      this.paymentShow = false;
      this.getpurcon();
      this.AddData = {};
    },
    // 取消权限配置提交
    paymentresetForm() {
      this.paymentShow = false;
      this.getpurcon();
      this.AddData = {};
    },
    // 权限配置信息提交
    paymentsubmitm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.api.busbenfig.edit(this.AddData).then((res) => {
            if (res.data.code == 200) {
              this.$message({
                type: "success",
                message: "权限配置成功!",
              });
              this.AddData = {};
              this.tableData = [];
              this.getpurcon();
              this.paymentShow = false;
            }
          });
        }
      });
    },
    // 熊猫甄选会员配置
    seleEditPurs(data) {
      console.log(data)
      // banner4
      if (data.banner4) {
        let url =
          this.http +
          "/saas-common/upload/imageShow?fileName=" +
          data.banner4 +
          "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
        this.selelistyear = [];
        this.selelistyear.push({ url: url });
        this.seleDealImg = this.selelistyear.length >= 1;
      }
      this.selementShow = true;
      this.seleData = JSON.parse(JSON.stringify(data));
    },
    // 熊猫甄选会员配置信息取消
    selseShowClose() {
  this.removedFileIds = []; // 取消 / 关闭:不清的话下次保存别的配置会把仍被引用的 banner 删掉
      this.selementShow = false;
      this.getpurcon();
      this.seleData = {};
    },
    // 取消熊猫甄选会员配置信息提交
    seleresetForm() {
  this.removedFileIds = []; // 取消 / 关闭:不清的话下次保存别的配置会把仍被引用的 banner 删掉
      this.selementShow = false;
      this.getpurcon();
      this.seleData = {};
    },
    // 熊猫甄选会员配置信息提交
    selesubmitm(formName) {
      console.log(this.seleData)
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.api.busconfig.edit(this.seleData).then((res) => {
            if (res.data.code == 200) {
              this.delRemovedFiles();
              this.$message({
                type: "success",
                message: "熊猫甄选会员信息配置成功!",
              });
              this.seleData = {};
              this.tableData = [];
              this.getpurcon();
              this.selementShow = false;
            }
          });
        }
      });
    },


    // 会员配置信息
    EditPurs(data) {
      // banner1
      if (data.banner1) {
        let url =
          this.http +
          "/saas-common/upload/imageShow?fileName=" +
          data.banner1 +
          "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
        this.listyearr = [];
        this.listyearr.push({ url: url });
        this.yeBtnDealImg = this.listyearr.length >= 1;
      }
      // banner2
      if (data.banner2) {
        let url =
          this.http +
          "/saas-common/upload/imageShow?fileName=" +
          data.banner2 +
          "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
        this.listmparr = [];
        this.listmparr.push({ url: url });
        this.mpBtnDealImg = this.listmparr.length >= 1;
      }

      // banner3
      if (data.banner3) {
        let url =
          this.http +
          "/saas-common/upload/imageShow?fileName=" +
          data.banner3 +
          "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
        this.listjcarr = [];
        this.listjcarr.push({ url: url });
        this.jcBtnDealImg = this.listjcarr.length >= 1;
      }
      this.editpaymentShow = true;
      this.editData = JSON.parse(JSON.stringify(data));
    },
    // 会员配置信息取消
    editpaymentShowClose() {
  this.removedFileIds = []; // 取消 / 关闭:不清的话下次保存别的配置会把仍被引用的 banner 删掉
      this.editpaymentShow = false;
      this.getpurcon();
      this.editData = {};
    },
    // 取消会员配置信息提交
    editpaymentresetForm() {
  this.removedFileIds = []; // 取消 / 关闭:不清的话下次保存别的配置会把仍被引用的 banner 删掉
      this.editpaymentShow = false;
      this.getpurcon();
      this.editData = {};
    },
    // 会员配置信息提交
    editpaymentsubmitm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.api.busconfig.edit(this.editData).then((res) => {
            if (res.data.code == 200) {
              this.delRemovedFiles();
              this.$message({
                type: "success",
                message: "会员信息配置成功!",
              });
              this.editData = {};
              this.tableData = [];
              this.getpurcon();
              this.editpaymentShow = false;
            }else if(res.data.code == 602){
                // 原先这段挂在 if (valid) 上:表单校验不通过时 res 还不存在,直接抛 ReferenceError
                this.$message({
                    type: "warning",
                    message: res.data.msg,
                });
            }
          });
        }
      });
    },
  },
};
</script>
  <style lang="less" scoped>
/deep/ .disyeloadSty .el-upload--picture-card {
  display: none !important; /* 上传banner1按钮隐藏 */
}
/deep/ .dismploadSty .el-upload--picture-card {
  display: none !important; /* 上传banner2按钮隐藏 */
}
/deep/ .disjcloadSty .el-upload--picture-card {
  display: none !important; /* 上传banner3按钮隐藏 */
}
// 详情样式开始
.detilsBox {
  width: 100%;
  height: 100%;
  padding: 25px;
  box-sizing: border-box;
  h5 {
    font-size: 16px;
    line-height: 58px;
    color: #333333;
  }
  .detilconent {
    width: 100%;
    span {
      display: inline-block;
      width: 33.33%;
      font-size: 14px;
      line-height: 24px;
      color: #333333;
    }
  }
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
  .pur-table {
    width: 100%;
    height: 94%;
  }
  // 暂无数据样式开始
  .dataAvailable {
    width: 100%;
    height: 245px;
    position: relative;
    margin-top: 185px;
    margin-bottom: 279px;
    .dataAva {
      position: absolute;
      top: 75%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;
      img {
        width: 371px;
        height: 200px;
        margin: auto;
      }
      p {
        font-size: 16px;
        line-height: 24px;
        color: #333;
      }
      .el-button {
        margin-bottom: 100px;
      }
    }
  }
  // 暂无数据样式结束
  .pur-pages {
    width: 100%;
    height: 6%;
    padding: 6px 0 0 12px;
    box-sizing: border-box;
    .el-pagination {
      float: right;
    }
  }
}
// 内容结束
</style>