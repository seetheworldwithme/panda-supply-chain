<!-- 租赁售卖页面 -->
<template>
  <div class="content">
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3">
        <span class="pur-size" style="margin-right: 8px">{{ purTitle }}</span>
        <span class="pur-size" style="font-size: 12px"
          >(页面用户访问量<i style="color: #cf1724; font-style: normal">{{
            accCount
          }}</i
          >)</span
        >
      </el-col>
      <el-col :span="21" class="pur-right">
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
        <!-- 筛选身份下拉框 -->
        <el-select v-model="carriageName" size="small" @change="prologistics">
          <el-option
            v-for="item in Extended"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">审核状态：</div>
        <el-select v-model="typeName" size="small" @change="changeType">
          <el-option
            v-for="item in goodsType"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">租赁售卖：</div>
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav" v-loading="loading" element-loading-text="加载中...">
      <!-- 表格 -->
      <div class="pur-table">
        <el-table
          height="100%"
          border
          :data="tableData"
          :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          style="width: 100%"
        >
          <el-table-column prop="id" v-if="show"> </el-table-column>
          <el-table-column label="操作" width="160" fixed="left">
            <template slot-scope="scope">
              <el-button @click="editInfo(scope.row.id)" size="mini"
                >编辑</el-button
              >
              <el-button size="mini" @click="delGoods(scope.row.id)"
                >删除</el-button
              >
            </template>
          </el-table-column>
          <el-table-column
            label="电话拨打量"
            :show-overflow-tooltip="true"
            width="120"
          >
            <template slot-scope="scope">
              <div
                @click="callDetails(scope.row.id)"
                style="cursor: pointer; color: #cf1724"
              >
                {{ scope.row.callCount }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            label="类型"
            :show-overflow-tooltip="true"
            width="120"
          >
            <template slot-scope="scope">
              <div v-if="scope.row.type == 1">租赁</div>
              <div v-if="scope.row.type == 2">售卖</div>
            </template>
          </el-table-column>
          <el-table-column
            label="企业名称"
            :show-overflow-tooltip="true"
            prop="enterpriseName"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="联系人"
            :show-overflow-tooltip="true"
            prop="coalRelName"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="联系电话"
            :show-overflow-tooltip="true"
            prop="coalRelPhone"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="报价"
            :show-overflow-tooltip="true"
            prop="offer"
            width="130"
          >
          </el-table-column>
          <el-table-column
            label="工艺"
            :show-overflow-tooltip="true"
            prop="coalWashingProcess"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="土地"
            :show-overflow-tooltip="true"
            prop="landSituation"
            width="150"
          >
          </el-table-column>
          <el-table-column
            label="状态"
            :show-overflow-tooltip="true"
            prop="coalPlantStatus"
            width="100"
          >
          </el-table-column>
          <el-table-column label="审核状态" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <div v-if="scope.row.status == 0">
                <div
                  class="statsBox"
                  style="background: #fb8138; box-shadow: 0px 4px 4px #fb8138"
                ></div>
                <span> 待审核</span>
              </div>
              <div v-if="scope.row.status == 1">
                <div
                  class="statsBox"
                  style="background: #46c7a2; box-shadow: 0px 4px 4px #46c7a2"
                ></div>
                <span> 已上架</span>
              </div>
              <div v-if="scope.row.status == 2">
                <div
                  class="statsBox"
                  style="background: #cf1724; box-shadow: 0px 4px 4px #f7d0d3"
                ></div>
                <span> 已拒绝</span>
              </div>
              <div v-if="scope.row.status == 3">
                <div
                  class="statsBox"
                  style="background: #bfbfbf; box-shadow: 0px 4px 4px #bfbfbf"
                ></div>
                <span> 已下架</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            label="最近站台"
            :show-overflow-tooltip="true"
            prop="nearestPlatform"
            width="140"
          >
          </el-table-column>
          <el-table-column
            label="煤厂产能（吨）"
            :show-overflow-tooltip="true"
            prop="coalPlantCapacity"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="精煤煤仓储能力（平方）"
            :show-overflow-tooltip="true"
            prop="cleanCoalStorage"
            width="180"
          >
          </el-table-column>
          <el-table-column
            label="原煤煤仓储能力（平方）"
            :show-overflow-tooltip="true"
            prop="rawCoalStorage"
            width="180"
          >
          </el-table-column>
          <el-table-column
            label="占地面积（亩）"
            :show-overflow-tooltip="true"
            prop="landArea"
            width="120"
          >
          </el-table-column>
          <el-table-column
            label="煤厂地区"
            :show-overflow-tooltip="true"
            prop="coalPlantArea"
            width="150"
          >
          </el-table-column>
          <el-table-column
            label="企业介绍"
            :show-overflow-tooltip="true"
            prop="cinfo"
            width="150"
          >
          </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable">
            <div class="dataAva">
              <img src="@/assets/zanwushuju.png" alt="" />
              <p>暂无数据！</p>
            </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <el-pagination
          style="float: right"
          background
          :page-sizes="[10, 30, 50, 100, 200]"
          :page-size.sync="size"
          :pager-count="5"
          :current-page.sync="current"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @current-change="getpurcon"
          @size-change="getpurcon"
        >
        </el-pagination>
      </div>
    </div>
    <!-- 编辑终端信息弹窗 -->
    <el-dialog
      center
      :closeOnClickModal="false"
      title="租赁售卖"
      :visible.sync="editShow"
      width="66%"
      :before-close="editShowClose"
    >
      <el-form
        ref="ruleForm"
        :model="form"
        :rules="rules"
        label-width="80px"
        label-position="top"
      >
        <el-row style="padding-left: 30px">
          <el-col :span="4" style="margin-right: 16px">
            <el-form-item label="类型" prop="type">
              <el-radio-group v-model="form.type">
                <el-radio :label="1">租赁</el-radio>
                <el-radio :label="2">售卖</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="4" style="margin-right: 16px">
            <el-form-item label="状态" prop="coalPlantStatus">
              <el-radio-group v-model="form.coalPlantStatus">
                <el-radio :label="1">运行</el-radio>
                <el-radio :label="2">停产</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="4" style="margin-right: 16px">
            <el-form-item label="土地" prop="landSituation">
              <el-radio-group v-model="form.landSituation">
                <el-radio :label="1">自有</el-radio>
                <el-radio :label="2">租赁</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="5">
            <el-form-item label="工艺" required>
              <el-checkbox-group @change="stateChange" v-model="checkList">
                <el-checkbox label="跳汰"></el-checkbox>
                <el-checkbox label="重介"></el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :span="6" style="margin-left: -21px">
            <el-form-item label="煤厂地址" prop="address">
              <el-input class="addressIpt" v-model="form.address" width="316"></el-input>
            </el-form-item>
          </el-col>
          <!-- <el-col :span="3">
            <el-form-item label="煤厂地区" prop="provinceCode">
              <el-select
                @change="enterprises"
                v-model="form.provinceCode"
                placeholder="请选择身份"
              >
                <el-option
                  v-for="item in provinceOptions"
                  :key="item.code"
                  :label="item.name"
                  :value="item.code"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="3">
            <el-form-item prop="cityCode" style="padding-top: 40px">
              <el-select
                @change="selectCity($event)"
                v-model="form.cityCode"
                placeholder="请选择城市"
              >
                <el-option
                  v-for="item in cityOptions"
                  :key="item.value"
                  :label="item.name"
                  :value="item.code"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="3">
            <el-form-item prop="areaCode" style="padding-top: 40px">
              <el-select
                @change="selectCountys($event)"
                v-model="form.areaCode"
                placeholder="请选择地区"
              >
                <el-option
                  v-for="item in countyOptions"
                  :key="item.value"
                  :label="item.name"
                  :value="item.code"
                >
                </el-option>
              </el-select>
            </el-form-item>
          </el-col> -->
        </el-row>
        <div class="box">
          <div class="iptBox">
            <el-form-item label="企业名称" prop="enterpriseName">
              <el-select class="selectName" @change="changeEnterpriseName" :remote-method="getEnterpriseName" v-model="form.enterpriseName" filterable remote placeholder="请选择">
                <el-option
                  v-for="item in enterpriseNameList"
                  :key="item.Name"
                  :label="item.Name"
                  :value="item.Name">
                </el-option>
              </el-select>
              <!-- <el-input v-model="form.enterpriseName" width="240"></el-input> -->
            </el-form-item>
            <el-form-item label="原煤仓储能力（平方）" prop="rawCoalStorage">
              <el-input v-model="form.rawCoalStorage"></el-input>
            </el-form-item>
            <el-form-item label="占地面积（亩）" prop="landArea">
              <el-input v-model="form.landArea"></el-input>
            </el-form-item>
            <el-form-item label="煤厂产能（吨）" prop="coalPlantCapacity">
              <el-input v-model="form.coalPlantCapacity"></el-input>
            </el-form-item>
            <el-form-item
              label="精煤煤仓储能力（平方）"
              prop="cleanCoalStorage"
            >
              <el-input v-model="form.cleanCoalStorage"></el-input>
            </el-form-item>
            <el-form-item label="最近站台" prop="nearestPlatform">
              <el-input v-model="form.nearestPlatform"></el-input>
            </el-form-item>
            <el-form-item label="报价（万元）" prop="offer">
              <el-input v-model="form.offer"></el-input>
            </el-form-item>
            <el-form-item label="联系人" prop="coalRelName">
              <el-input v-model="form.coalRelName"></el-input>
            </el-form-item>
            <el-form-item label="联系电话" prop="coalRelPhone">
              <el-input v-model="form.coalRelPhone"></el-input>
            </el-form-item>
          </div>
          <div class="txtBox">
            <el-form-item label="企业介绍" prop="cinfo">
              <el-input type="textarea" v-model="form.cinfo"></el-input>
            </el-form-item>
          </div>
        </div>
        <el-row>
          <el-col :span="12">
            <el-form-item
              label="煤厂证书（排污证、取水证等）"
              required
              style="padding-left: 30px"
            >
              <el-upload
                :action="http + '/saas-common/upload/batch'"
                :headers="header"
                :file-list="certificate"
                name="file"
                list-type="picture-card"
                :on-success="addCertList"
                :on-preview="certPreview"
                :before-remove="certBeforeRemove"
                :on-remove="certRemove"
              >
                <i class="el-icon-plus"></i>
              </el-upload>
              <el-dialog :visible.sync="certPreDialog">
                <img width="100%" :src="certUrl" alt="" />
              </el-dialog>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="洗煤厂照片" required>
              <el-upload
                :action="http + '/saas-common/upload/batch'"
                :headers="header"
                :file-list="imgList"
                name="file"
                list-type="picture-card"
                :on-success="addImgList"
                :on-preview="imgPreview"
                :before-remove="imgBeforeRemove"
                :on-remove="imgRemove"
              >
                <i class="el-icon-plus"></i>
              </el-upload>
              <el-dialog :visible.sync="imgPreDialog">
                <img width="100%" :src="imageUrl" alt="" />
              </el-dialog>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <div style="float: right">
            <el-button
              type="primary"
              @click="onSubmit('ruleForm')"
              class="dialogBtn"
              >通过</el-button
            >
            <el-button @click="refuse('ruleForm')" class="dialogBtn"
              >拒绝</el-button
            >
          </div>
        </el-form-item>
      </el-form>
    </el-dialog>
    <!-- 电话拨打详情抽屉 -->
    <el-drawer :visible.sync="isShow" :with-header="false" size="700px">
      <div class="draBox">
        <!-- 抽屉标题 -->
        <div class="draHeader">
          <div class="draTitle">电话拨打数据</div>
          <div @click="closeDrawer" class="draClose">
            <i class="el-icon-close"></i>
          </div>
        </div>
        <!-- 电话拨打列表 -->
        <div v-if="ascgridData.length != 0" class="draContent">
          <el-table
            :data="ascgridData"
            height="100%"
            border
            :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
          >
            <el-table-column
              label="访问用户"
              prop="callUserPhone"
            ></el-table-column>
            <el-table-column
              label="访问时间"
              prop="updateTime"
            ></el-table-column>
          </el-table>
          <!-- 分页 -->
          <div class="mzpages">
            <el-pagination
              background
              :page-sizes="[10, 30, 50, 100, 200]"
              :page-size.sync="mzsize"
              :pager-count="5"
              :current-page.sync="mzcurrent"
              layout="total, sizes, prev, pager, next, jumper"
              :total="mztotal"
              @current-change="mzdata"
              @size-change="mzdata"
            >
            </el-pagination>
          </div>
        </div>
        <div v-else class="noData"> 
          <div>
            <img src="../../assets/zanwushuju.png" alt="" />
          </div>
          <div style="font-size: 16px; color: #666">暂无数据!</div>
        </div>
      </div>
    </el-drawer>
    <!-- 内容结束 -->
  </div>
</template>
  <script>
export default {
  data() {
    return {
      http: null,
      loading: false,
      current: 1, // 当前页数
      size: 100, // 显示条数
      total: 0, // 总条数
      goodsStatus: null,
      goodsCategory: null,
      show: false,
      editShow: false,
      purTitle: "", // 标题
      purSearch: "", // 搜索
      tableData: [
        // {
        //   enterpriseName: "企业名称", // 企业名称
        //   coalRelName: "测试联系人",
        //   coalRelPhone: "15234141205",
        //   offer: 2000, // 报价
        //   coalWashingProcess: "工艺内容",
        //   landSituation: "土地土地",
        //   coalPlantStatus: "状态啊",
        //   nearestPlatform: "最近的站台是什么",
        //   coalPlantCapacity: 3000, // 煤厂产能(吨)
        //   cleanCoalStorage: 1000, // 精煤仓储能力（平方）
        //   rawCoalStorage: 2000, // 原煤仓储能力（平方）
        //   landArea: 4000, // 占地面积（亩）
        //   proName: "地区地区", // 省名称
        //   cityName: "", // 市名称
        //   areaName: "", // 区名称
        //   provinceCode: null, // 省级代码
        //   cityCode: null, // 市级代码
        //   areaCode: null, // 区级代码
        //   cinfo: "介绍介绍介绍", // 企业介绍
        //   status: 0, // 审核状态（0、未审核1、已通过2、已拒绝3、已下架）
        // },
      ], // 表格数据
      typeName: null,
      goodsType: [
        {
          dkey: "全部",
          dvalue: null,
        },
        {
          dkey: "租赁",
          dvalue: 1,
        },
        {
          dkey: "售卖",
          dvalue: 2,
        },
      ],
      Extended: [
        {
          dkey: "全部",
          dvalue: null,
        },
        {
          dkey: "待审核",
          dvalue: 0,
        },
        {
          dkey: "已上架",
          dvalue: 1,
        },
        {
          dkey: "已拒绝",
          dvalue: 2,
        },
        {
          dkey: "已下架",
          dvalue: 3,
        },
      ], // 查询类型
      carriageName: null,
      form: {},
      rules: {
        type: [{ required: true, message: "请选择类型", trigger: "change" }],
        coalPlantStatus: [
          { required: true, message: "请选择洗煤厂状态", trigger: "change" },
        ],
        landSituation: [
          { required: true, message: "请选择土地状态", trigger: "change" },
        ],
        address: [
          { required: true, message: "请填写煤厂地址", trigger: "change" },
        ],
        // provinceCode: [
        //   { required: true, message: "请选择省份", trigger: "change" },
        // ],
        // cityCode: [
        //   { required: true, message: "请选择城市", trigger: "change" },
        // ],
        // areaCode: [
        //   { required: true, message: "请选择地区", trigger: "change" },
        // ],
        enterpriseName: [
          { required: true, message: "请填写企业名称", trigger: "change" },
        ],
        rawCoalStorage: [
          {
            required: true,
            message: "请填写原煤仓储能力（平方）",
            trigger: "change",
          },
          {
            pattern: /^[1-9]\d*$/,
            message: "请填写大于0的整数",
          },
        ],
        landArea: [
          {
            required: true,
            message: "请填写占地面积（亩）",
            trigger: "change",
          },
          {
            pattern: /^[1-9]\d*$/,
            message: "请填写大于0的整数",
          },
        ],
        coalPlantCapacity: [
          {
            required: true,
            message: "煤厂产能（吨）",
            trigger: "change",
          },
          {
            pattern: /^[1-9]\d*$/,
            message: "请填写大于0的整数",
          },
        ],
        cleanCoalStorage: [
          {
            required: true,
            message: "请填写精煤煤仓储能力（平方）",
            trigger: "change",
          },
          {
            pattern: /^[1-9]\d*$/,
            message: "请填写大于0的整数",
          },
        ],
        nearestPlatform: [
          { required: true, message: "请填写最近站台", trigger: "change" },
        ],
        offer: [
          { required: true, message: "请输入报价", trigger: "change" },
          {
            pattern: /^[1-9]\d*$/,
            message: "请填写大于0的整数",
          },
        ],
        coalRelName: [
          { required: true, message: "请填写联系人", trigger: "change" },
        ],
        coalRelPhone: [
          { required: true, message: "请填写联系电话", trigger: "change" },
          {
            pattern:
              /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/,
            message: "请输入正确的手机格式",
          },
        ],
        cinfo: [
          { required: true, message: "请选择填写企业介绍", trigger: "change" },
        ],
      },
      certUrl: "",
      certPreDialog: false,
      imageUrl: "",
      imgPreDialog: false,
      header: {}, // 上传token
      isShow: false,
      mzcurrent: 1,
      mzsize: 10,
      ascgridData: [],
      mztotal: 0,
      accCount: 0,
      itemId: null,
      checkList: [],
      provinceOptions: [],
      cityOptions: [],
      countyOptions: [],
      certificate: [],
      imgList: [],
      enterpriseNameList:[]
    };
  },
  created() {
    this.http = this.$store.state.http; // IP地址
    this.header.Authorization = sessionStorage.getItem("Authorization");
    this.purTitle = this.$route.meta.title;
    this.getpurcon();
    this.getAccCount();
    this.getProvinceInfo();
  },
  mounted() {},
  methods: {
    // 新建标签取消
    paymentShowClose() {
      this.paymentShow = false;
      // this.getpurcon();
      this.AddData = []; 
    },
    // 表格数据
    getpurcon(data) {
      this.tableData = [];
      if (!data) this.current = 1; // 入参是翻页的页码:原先无条件置 1,第 2 页永远查不到
      this.loading = true;
      // 传参
      var listByPageData = {
        searchStr: this.purSearch,
        page: {
          current: this.current,
          size: this.size,
        },
        status: this.goodsStatus,
        type: this.goodsCategory,
      };
      // 渲染表格
      this.api.busEnterpriseLeasing.listByPage(listByPageData).then((res) => {
        this.loading = false;
        if (res.data.code == 200) {
          res.data.data.records.forEach((item) => {
            item.coalPlantArea = item.proName
              ? item.proName + item.cityName + item.areaName
              : "";
          });
          this.tableData = res.data.data.records;
          this.total = res.data.data.total;
        }
      });
    },
    // 搜索功能
    search() {
      this.tableData = [];
      this.current = 1;
      this.getpurcon();
    },
    // 查询上下架搜索功能
    prologistics(data) {
      this.tableData = [];
      this.goodsStatus = data;
      this.current = 1;
      this.getpurcon();
    },
    changeType(data) {
      this.tableData = [];
      this.goodsCategory = data;
      this.current = 1;
      this.getpurcon();
    },
    // 编辑
    async editInfo(id) {
      this.editShow = true;
      this.imgList = []; // 不清会残留上一条记录的证书图片,删它还会按旧下标删掉这条的附件
      this.checkList = [];
      try {
        const res = await this.api.busEnterpriseLeasing.getById({ id });
        if (res.data.code == 200) {
          this.form = res.data.data;
          if (this.form.coalWashingProcess) {
            this.checkList = this.form.coalWashingProcess.split(",");
          }
          this.form.cerList = [];
          this.form.photoList = [];
          this.api.areaInfo
            .all({ parentCode: this.form.provinceCode })
            .then((res) => {
              if (res.data.code == 200) {
                this.cityOptions = res.data.data;
              }
            });
          this.api.areaInfo
            .all({ parentCode: this.form.cityCode })
            .then((res) => {
              if (res.data.code == 200) {
                this.countyOptions = res.data.data;
              }
            });
          if (this.form.certAssociationId) {
            this.api.addfile
              .list({
                association_id: this.form.certAssociationId,
                is_compressed_file: 0,
              })
              .then((res) => {
                res.data.data.forEach((item) => {
                  let url =
                    this.http +
                    "/saas-common/upload/imageShow?fileName=" +
                    item.id +
                    "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
                  this.certificate.push({
                    url: url,
                  });
                  this.form.cerList.push({ id: item.id });
                });
              });
          }
          if (this.form.photoAssociationId) {
            this.api.addfile
              .list({
                association_id: this.form.photoAssociationId,
                is_compressed_file: 0,
              })
              .then((res) => {
                res.data.data.forEach((item) => {
                  let url =
                    this.http +
                    "/saas-common/upload/imageShow?fileName=" +
                    item.id +
                    "&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99";
                  this.imgList.push({
                    url: url,
                  });
                  this.form.photoList.push({ id: item.id });
                });
              });
          }
        }
      } catch (error) {
        console.log(error);
      }
    },
    editShowClose() {
      this.editShow = false;
    },
    onSubmit(formName) {
      if (this.checkList[0] == "") {
        return this.$message({
          message: "请选择至少一种工艺",
          type: "warning",
        });
      }
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.form.status = 1;
          this.api.busEnterpriseLeasing.updateById(this.form).then((res) => {
            if (res.data.code == 200) {
              this.imgList = [];
              this.certificate = [];
              this.checkList = [];
              this.editShow = false;
              this.getpurcon();
            }
          });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    refuse(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.form.status = 2;
          this.api.busEnterpriseLeasing.updateById(this.form).then((res) => {
            if (res.data.code == 200) {
              this.imgList = [];
              this.certificate = [];
              this.checkList = [];
              this.editShow = false;
              this.getpurcon();
            }
          });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    addCertList(response, file, fileList) {
      if (response.code == 200) {
        this.form.cerList.push({ id: Number(response.data[0].id) });
      }
    },
    addImgList(response, file, fileList) {
      if (response.code == 200) {
        this.form.photoList.push({ id: Number(response.data[0].id) });
      }
    },
    certBeforeRemove(file, fileList) {
      const index = fileList.findIndex((item) => item.url === file.url);
      this.form.cerList.splice(index, 1);
    },
    imgBeforeRemove(file, fileList) {
      const index = fileList.findIndex((item) => item.url === file.url);
      this.form.photoList.splice(index, 1);
    },
    certRemove(file, fileList) {},
    imgRemove(file, fileList) {
      // console.log(file, fileList);
    },
    certPreview(file) {
      this.certUrl = file.url;
      this.certPreDialog = true;
    },
    imgPreview(file) {
      this.imageUrl = file.url;
      this.imgPreDialog = true;
    },
    // 拨打详情
    callDetails(id) {
      this.itemId = id;
      this.api.phonelog
        .listByPage({
          leaseId: id,
          pType: 2,
          searchStr: "",
          page: {
            current: this.mzcurrent,
            size: this.mzsize,
          },
        })
        .then((res) => {
          if (res.data.code == 200) {
            this.ascgridData = res.data.data.records;
            this.mztotal = res.data.data.total;
            this.isShow = true;
          }
        });
    },
    async getAccCount() {
      let obj = { operateObj: 5, operateType: 7 };
      try {
        const res = await this.api.busprolog.selCountByUser(obj);
        if (res.data.code == 200) {
          this.accCount = res.data.data;
        }
      } catch (error) {
        console.log(error);
      }
    },
    closeDrawer() {
      this.isShow = false;
    },
    mzdata(data) {
      this.mzcurrent = data;
      this.callDetails(this.itemId);
    },
    stateChange(data) {
      this.form.coalWashingProcess = data.join(",");
    },
    getProvinceInfo() {
      this.api.areaInfo.all({ parentCode: 0 }).then((res) => {
        if (res.data.code == 200) {
          this.provinceOptions = res.data.data;
        }
      });
    },
    enterprises(data) {
      this.api.areaInfo
        .all({ parentCode: this.form.provinceCode })
        .then((res) => {
          if (res.data.code == 200) {
            this.cityOptions = res.data.data;
          }
        });
    },
    selectCity(e) {
      this.api.areaInfo.all({ parentCode: this.form.cityCode }).then((res) => {
        if (res.data.code == 200) {
          this.countyOptions = res.data.data;
        }
      });
    },
    selectCountys(e) {
      this.form.areaCode = e;
    },
    delGoods(id) {
      this.$confirm("此操作将永久删除该文件, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        this.api.busEnterpriseLeasing.del({id}).then((res) => {
          // 字段名拼错(cdoe):删除成功后既不提示也不刷新
          if (res.data.code == 200) {
            this.$message({
              type: "success",
              message: "删除成功!",
            });
            this.tableData = [];
            this.getpurcon();
          }
        });
      });
    },
    getEnterpriseName(query){
      this.api.enterData.seletcname({etcName:query})
				.then(res=>{
					if(res.data.code == 200 && res.data.data.length>0){
						this.enterpriseNameList = res.data.data;
					}
				})
    },
    changeEnterpriseName(name){
      this.api.enterData.seletcinfo({enterpriseName:name,isAutoAdd:true})
				.then(res=>{
					if(res.data.code == 200){
          console.log(res.data.data,'datatataa');
						// this.addForm.companyId = res.data.data.id;
            this.form.enterpriseName = res.data.data.enterpriseName
            this.form.address = res.data.data.address
            this.form.cinfo = res.data.data.scope
					}
				})
    }
  },
};
</script>
<style lang="less" scoped>
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
    .el-select {
      width: 11%;
      margin-right: 10px;
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
  }
}

.statsBox {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  background: #6fcf97;
  box-shadow: 0px 4px 4px #cbf9de;
  span {
    color: #333;
  }
}

/deep/ .el-form-item__content {
  .el-input {
    width: 240px !important;
  }

  .addressIpt{
    width: 316px !important;
  }

  .el-input--suffix {
    width: 128px !important;
    margin-right: 8px !important;

    .el-input__inner {
      width: 128px !important;
    }
  }

  .el-textarea__inner {
    width: 316px !important;
    height: 244px !important;
  }

  .selectName {
    width: 240px !important;

    .el-input--suffix {
      width: 240px !important;

      .el-input__inner {
        width: 240px !important;
      }
    }
  }
}

/deep/ .el-form-item {
  .el-form-item__label {
    padding: 0 !important;
  }
}

/deep/ .el-radio-group {
  .el-radio {
    margin-right: 10px !important;
  }
}

/deep/ .el-checkbox-group {
  .el-checkbox {
    margin-right: 10px !important;
  }
}

.box {
  width: 100%;
  display: flex;
  justify-content: space-between;

  .iptBox {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
  }

  .txtBox {
    width: 316px;
  }
}

.dialogBtn {
  width: 140px;
  height: 44px;
  margin-left: 33px;
}
.selTit {
  height: 32px;
  float: right;
  line-height: 32px;
  font-size: 14px;
}
.draBox {
  height: 100%;

  .draTitle {
    line-height: 60px;
    font-size: 18px;
    color: #333;
    font-weight: 700;
  }

  .draContent {
    box-sizing: border-box;
    height: calc(100% - 132px);
    padding: 20px;
  }

  .noData {
    width: 100%;
    height: calc(100% - 100px);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
  }
}
.draHeader {
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 0 32px;
  background-color: #f3f3f3;
  box-sizing: border-box;

  .draClose {
    cursor: pointer;

    i {
      line-height: 60px;
      font-size: 26px;
    }
  }
}
.mzpages {
  height: 32px;
  width: 100%;
  padding: 20px 0px;
  // box-sizing: border-box;
  .el-pagination {
    float: right;
  }
}
// 内容结束


</style>

    