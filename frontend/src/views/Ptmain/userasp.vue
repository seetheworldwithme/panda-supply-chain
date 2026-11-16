<template>
  <!-- 用户列表页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="2"
        ><span class="pur-size">{{ purTitle }}</span></el-col
      >
      <el-col :span="22" class="pur-right">
        <el-date-picker
          clearable
          size="small"
          @change="usrTie"
          v-model="usrTimArray"
          type="daterange"
          unlink-panels
          range-separator="~"
          start-placeholder="开始注册日期"
          end-placeholder="结束注册日期"
          value-format="yyyy-MM-dd"
          style="width: 20%; float: right">
        </el-date-picker>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch"
        >
        </el-input>
        <el-select v-model="carriageName" size="small" @change="prologistics">
          <el-option
            v-for="item in Extended"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">筛选平台用户：</div>
        <!-- 筛选有效用户下拉框 -->
        <el-select v-model="effective" size="small" @change="selectUser">
          <el-option
            v-for="item in effSelect"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue"
          >
          </el-option>
        </el-select>
        <div class="selTit">筛选用户：</div>
        <!-- 筛选身份下拉框 -->
        <el-select v-model="userIdentity" size="small" @change="selectIdentity">
          <el-option
            v-for="item in tabs"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          >
          </el-option>
        </el-select>
        <div class="selTit">筛选身份：</div>
        <!-- 筛选身份下拉框 -->
        <el-select v-model="hyvalue" size="small" @change="selectIhy">
          <el-option
            v-for="item in listhy"
            :key="item.dvalue"
            :label="item.dkey"
            :value="item.dvalue">
          </el-option>
        </el-select>
        <div class="selTit">筛选会员：</div>
        <el-popover
          placement="bottom"
          width="400"
          trigger="hover">
          <el-form :inline="true" :model="formInline">
            <el-form-item label="注册手机号">
              <el-input v-model="formInline.phone" placeholder="请输入手机号"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="small" @click="onSubmit">确认注册</el-button>
            </el-form-item>
          </el-form>
          <el-button slot="reference" size="small" type="primary">注册用户</el-button>
        </el-popover>
        
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
          style="width: 100%">
          <el-table-column prop="id" v-if="show"> </el-table-column>
          <el-table-column label="操作" width="100" fixed="left">
            <template slot-scope="scope">
              <el-dropdown trigger="click" placement="right">
                <el-button type="primary" size="small">操 作</el-button>
                <el-dropdown-menu slot="dropdown">  
                  <el-dropdown-item v-if="scope.row.wareIdentity == 8" @click.native="delPurs(scope.row)">取消平台经纪人</el-dropdown-item>
                  <el-dropdown-item v-else @click.native="EditPurs(scope.row)">设为平台经纪人</el-dropdown-item>
                  <el-dropdown-item @click.native="Editserc(scope.row,scope.$index)">编辑信息</el-dropdown-item>
                  <el-dropdown-item @click.native="uplserc(scope.row)" v-if="scope.row.accountType!=2 && scope.row.level!=2">设为管理员</el-dropdown-item>
                  <el-dropdown-item @click.native="nouplserc(scope.row)" v-else>取消管理员</el-dropdown-item>
                  <el-dropdown-item @click.native="delserc(scope.row)">注销账号</el-dropdown-item>
                  <el-dropdown-item @click.native="openVip(scope.row)" v-if="scope.row.memberType == 0">开通VIP会员</el-dropdown-item>
                  <el-dropdown-item @click.native="noVip(scope.row)" v-if="scope.row.memberType == 1">取消会员</el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
            </template>
          </el-table-column>
          <el-table-column
            prop="phone"
            label="手机号"
            :show-overflow-tooltip="true"
            width="140">
          <template slot-scope="scope">
              {{scope.row.phone?scope.row.phone:'用户未登录'}}
            </template>
          </el-table-column>
          <el-table-column prop="expirationTime" label="会员到期时间" :show-overflow-tooltip="true" width="160">
            <template slot-scope="scope">
              <span v-if="scope.row.expirationTime">{{ scope.row.expirationTime }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="wareIdentity"
            label="身份信息"
            :show-overflow-tooltip="true"
            width="140"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.wareIdentity == 1">洗煤厂</span>
              <span v-else-if="scope.row.wareIdentity == 2">贸易商</span>
              <span v-else-if="scope.row.wareIdentity == 3">焦化厂</span>
              <span v-else-if="scope.row.wareIdentity == 4">经纪人</span>
              <span v-else-if="scope.row.wareIdentity == 5">煤矿</span>
              <span v-else-if="scope.row.wareIdentity == 6">钢厂</span>
              <span v-else-if="scope.row.wareIdentity == 7">金融企业</span>
              <span v-else-if="scope.row.wareIdentity == 8">平台经纪人</span>
              <span v-else-if="scope.row.wareIdentity == 9">其它</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column prop="username" label="姓名" :show-overflow-tooltip="true" width="140">
            <template slot-scope="scope">
              <span v-if="scope.row.username">{{ scope.row.username }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="createTime"
            label="注册时间"
            :show-overflow-tooltip="true"
            width="160">
            <template slot-scope="scope">
              <span v-if="scope.row.createTime">{{
                scope.row.createTime
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <!-- <el-table-column
            label="兑换记录"
            :show-overflow-tooltip="true"
            width="80"
          >
            <template slot-scope="scope">
              <span
                style="color: #cf1724; cursor: pointer"
                @click="exchange(scope.row.id)"
                >详情</span
              >
            </template>
          </el-table-column> -->
          <!-- <el-table-column
            label="访问数据"
            :show-overflow-tooltip="true"
            width="80"
          >
            <template slot-scope="scope">
              <span
                style="color: #cf1724; cursor: pointer"
                @click="userchange(scope.row)"
                >详情</span
              >
            </template>
          </el-table-column> -->
          <el-table-column label="核实记录" :show-overflow-tooltip="true" width="80">
            <template slot-scope="scope">
              <span style="color: #cf1724; cursor: pointer" @click="verifychange(scope.row)" v-if="scope.row.wareIdentity==8">详情</span>
              <span v-else>--</span>
            </template>
          </el-table-column>  
          <el-table-column
            prop="followPerson"
            label="跟进人"
            :show-overflow-tooltip="true"
            width="140"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.followPerson">{{
                scope.row.followPerson
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="followStatus"
            label="跟进状态"
            :show-overflow-tooltip="true"
            width="140"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.followStatus">{{
                scope.row.followStatus
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <!-- <el-table-column prop="sendAddres" label="黑金劵" :show-overflow-tooltip="true" width="140">
              </el-table-column> -->
          <!-- <el-table-column prop="fstName" label="账号" :show-overflow-tooltip="true">
              </el-table-column> -->
          <el-table-column
            prop="nickname"
            label="昵称"
            :show-overflow-tooltip="true"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.nickname">{{ scope.row.nickname }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="enterpriseName"
            label="企业名称"
            :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.enterpriseName">{{ scope.row.enterpriseName }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="duty"
            label="职位"
            :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.duty">{{ scope.row.duty }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
            <div class="dataAva">
              <img src="@/assets/zanwushuju.png" alt="" />
              <p>暂无数据!</p>
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
          @current-change="getpurcon"
          @size-change="getpurcon"
        >
        </el-pagination>
      </div>
      <!-- 配置经纪人信息弹窗 -->
      <el-dialog
        :closeOnClickModal="false"
        title="配置经纪人信息"
        :visible.sync="editpaymentShow"
        width="30%"
        :before-close="editpaymentShowClose"
      >
        <el-form
          style="width: 95%"
          :model="editData"
          :rules="editpaymentrules"
          ref="AddData"
          size="small"
          label-width="140px"
        >
          <el-form-item label="经纪人" prop="username">
            <el-input v-model="editData.username"></el-input>
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="editData.phone"></el-input>
          </el-form-item>
          <el-form-item label="工作开始时间" prop="brokerStartDate">
            <el-date-picker
              style="width: 100%"
              v-model="editData.brokerStartDate"
              type="date"
              value-format="yyyy-MM-dd"
              placeholder="选择日期"
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item label="企业客服微信链接" prop="servicerUrl">
            <el-input v-model="editData.servicerUrl"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="editpaymentsubmitm('AddData')"
              style="float: right; margin: 5% 0 0 5%"
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
      <!-- 用户信息弹窗 -->
      <el-dialog
        :closeOnClickModal="false"
        title="用户信息"
        :visible.sync="editphow"
        width="30%"
        :before-close="editphowClose">
        <el-form
          style="width: 95%"
          :model="edituser"
          :rules="euserrules"
          ref="userData"
          size="small"
          label-width="140px">
          <el-form-item label="姓名">
            <el-input v-model="edituser.username"></el-input>
          </el-form-item>
          <el-form-item label="身份选择">
            <div class="labfl">
              <div
                class="totalTab"
                :class="{ active: item.id === sel }"
                v-for="item in tabs"
                :key="item.id"
                @click="selectss(item)"
                style="cursor: pointer">
                <span>{{ item.name }}</span>
              </div>
            </div>
          </el-form-item>
          <el-form-item label="跟进人">
            <el-input v-model="edituser.followPerson"></el-input>
          </el-form-item>
          <el-form-item label="跟进状态">
            <el-input v-model="edituser.followStatus"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="edituserbmitm('userData')"
              style="float: right; margin: 5% 0 0 5%"
              >确定</el-button
            >
            <el-button
              @click="edituserForm"
              style="float: right; margin: 5% 0 0 2%"
              >取消</el-button
            >
          </el-form-item>
        </el-form>
      </el-dialog>
      <!-- 兑换记录抽屉 -->
      <el-drawer title="黑金劵兑换记录" :visible.sync="drawer">
        <div class="blocktitle">
          当前账号累计拥有黑金劵 {{ seeData.totalScore }}
        </div>
        <div class="block">
          <el-timeline>
            <el-timeline-item
              :timestamp="v.createTime"
              placement="top"
              v-for="(v, k) in seeData.logs"
              :key="k"
            >
              <el-card>
                <h4 style="line-height: 32px">手机号：{{ v.userPhone }}</h4>
                <p>兑换类型： {{ v.content }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </div>
      </el-drawer>
      <!-- 访问数据抽屉 -->
      <el-drawer :visible.sync="druserer" :before-close="ImClose">
        <div class="usertitle">
          <span>访问数据</span>
          <div v-if="activeName!=1" style="display:inline-block;marginLeft:10px;">共 <span style="fontSize:16px;fontWeight:400;color:red;">{{seuserData.length}}</span> 条</div>
          <div class="usrscreen">
            <el-date-picker
              clearable
              size="small"
              @change="Times"
              v-model="TimeArrays"
              type="daterange"
              unlink-panels
              range-separator="~"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="yyyy-MM-dd"
              style="width: 70%; float: right"
            >
            </el-date-picker>
            <p>筛选查询</p>
          </div>
        </div>
        <div style="padding: 16px;box-sizing: border-box;">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="会员VIP页面" name="1">
              <el-table border :data="vipdata"
              :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
              style="width: 100%">
                <el-table-column prop="id" v-if="show"> </el-table-column>
                <el-table-column
                  label="最后一次访问时间"
                  :show-overflow-tooltip="true">
                  <template slot-scope="scope">
                    <span v-if="scope.row.lastVisTime">{{ scope.row.lastVisTime }} 次</span>
                    <span v-else>--</span>
                  </template>
                </el-table-column>
                <el-table-column
                  prop="callTimes"
                  label="拨打客服次数"
                  :show-overflow-tooltip="true"
                  width="120">
                </el-table-column>
                <el-table-column label="点击支付按钮次数" prop="clickTimes" :show-overflow-tooltip="true">
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="煤种详情页面" name="2">
              <!-- 列表页 -->
              <div class="usrclock" v-if="seuserData.length > 0">
                <div class="usrbox" v-for="(v, k) in seuserData" :key="k">
                  <div class="usrh1">{{ v.enterpriseName }}</div>
                  <div class="usrh2">{{ v.warehouseNb }}</div>
                  <div class="usrh3" v-if="v.createTime">
                    最后一次访问时间：{{ v.createTime }}
                  </div>
                  <div class="usrfrom">
                    <el-table
                      height="160"
                      border
                      :data="v.phoneLogs"
                      :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                      style="width: 100%">
                      <el-table-column prop="id" v-if="show"> </el-table-column>
                      <el-table-column
                        prop="username"
                        label="经纪人姓名"
                        :show-overflow-tooltip="true"
                        width="100">
                      </el-table-column>
                      <el-table-column
                        prop="phone"
                        label="经纪人电话"
                        :show-overflow-tooltip="true"
                        width="120"
                      >
                      </el-table-column>
                      <el-table-column
                        label="被拨打电话次数"
                        :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                          <el-popover placement="bottom" width="400" trigger="click">
                            <el-table :data="gridData">
                              <el-table-column
                                width="150"
                                prop="callUserPhone"
                                label="电话"
                              ></el-table-column>
                              <el-table-column
                                prop="createTime"
                                label="拨打时间"
                              ></el-table-column>
                              <!-- <el-table-column width="300" property="address" label="地址"></el-table-column> -->
                            </el-table>
                            <span
                              v-if="scope.row.callCount"
                              slot="reference"
                              style="color: #e2231a; cursor: pointer"
                              @click="girtusr(v.proId, scope.row.phone)"
                              >{{ scope.row.callCount }} 次</span
                            >
                            <span v-else>--</span>
                          </el-popover>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>
                </div>
              </div>
              <div class="usrifbox" v-else>
                <el-empty description="暂无访问数据！"></el-empty>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </el-drawer>
      <!-- 核实记录抽屉 -->
      <el-drawer :visible.sync="isShow" :with-header="false" size="700px">
        <div class="draBox">
          <!-- 抽屉标题 -->
          <div class="draTitle">经纪人核实记录</div>
          <!-- 核实内容 -->
           <div v-if="verifyList.length!=0" class="draContent">
              <!-- 核实项 -->
              <div class="contentItem" v-for="(item,index) in verifyList" :key="index">
                <!-- 核实时间节点 -->
                <div class="timeNode">
                  <div v-if="index==0" class="itemIcon">
                    <img src="@/assets/verifyIcon.png">
                  </div>
                  <div v-else class="itemIcon" style="background:#D9D9D9"></div>
                  <div class="infoTime">{{item.updateTime}} 核实人：{{item.verifyPeopleName}}</div>
                </div>
                <div class="itemInfo">
                  <!-- <div class="timeLine"></div> -->
                  <div class="txtBox">
                    <div class="adressTxt">{{item.verifyAddress}}<img src="@/assets/addressImg.png"></div>
                    <div class="comment">{{item.evaluate}}</div>
                  </div>
                </div>
              </div>
            </div>

          <div v-else class="noData">
            <div>
              <img src="../../assets/zanwushuju.png" alt="" />
            </div>
            <div style="font-size: 16px; color: #666">暂无核实记录</div>
          </div>
        </div>
      </el-drawer>
      <!-- VIP开通弹框 -->
      <el-dialog
        :closeOnClickModal="false"
        title="开通会员"
        :visible.sync="vipshow"
        width="30%"
        :before-close="vipClose">
        <el-form style="width: 95%" :model="vipusertime" size="small" label-width="140px">
          <el-form-item label="开通手机号">
            <el-input v-model="vipusertime.phone" disabled></el-input>
          </el-form-item>
          <el-form-item label="会员到期时间:" v-if="vipusertime.andtime && vipusertime.expirationTime">
            <!-- <el-date-picker
              clearable
              size="small"
              v-model="vipTimArray"
              type="daterange"
              unlink-panels
              range-separator="~"
              start-placeholder="开始注册日期"
              end-placeholder="结束注册日期"
              value-format="yyyy-MM-dd"
              style="width: 100%;">
            </el-date-picker> -->
            <el-date-picker
              clearable
              @change="changeTime"
              v-model="vipusertime.expirationTime"
              type="datetime"
              placeholder="选择日期时间"
              style="width: 100%;">
            </el-date-picker>
            <!-- <span>{{ vipusertime.andtime }} ~ {{ vipusertime.expirationTime }}</span> -->
            <!-- <el-input v-model="vipusertime.viptimes" disabled></el-input> -->
          </el-form-item>
          <el-form-item>
            <el-radio @input="viptime" v-model="vipusertime.activateType" label="1">年度会员</el-radio>
            <el-radio @input="viptime" v-model="vipusertime.activateType" label="2">月度会员</el-radio>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="editvipbmitm()" style="float: right; margin: 5% 0 0 5%" >确认开通</el-button>
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
      http: "",
      cloudhttp: "",
      purTitle: "", // 标题
      show: false,
      // 注册新用户
      formInline:{
        phone:'',
      },
      /*
      列表字段
      */
      usrTimArray: [], // 注册时间搜索
      purSearch: "", // 搜索
      dataAva: false,
      loading: false,
      current: 1, // 当前页数
      size: 200, // 显示条数
      total: 0, // 总条数
      tableData: [], // 表格数据

      drawer: false, // 查看兑换抽屉
      seeData: {}, // 查看兑换数据

      druserer: false, // 查看访问抽屉
      seuserData: [], // 查看访问数据
      TimeArrays: [], // 时间搜索
      usrlist: [],
      uid: null,
      gridData: [],
      activeName:'1',  // 数据访问tab
      vipdata:[],       // vip访问数据

      // 经纪人配置
      editpaymentShow: false,
      editData: {},
      editpaymentrules: {
        username: [
          { required: true, message: "请输入经纪人姓名", trigger: "blur" },
        ],
        phone: [
          { required: true, message: "请输入经纪人电话", trigger: "blur" },
        ],
        servicerUrl: [
          { required: true, message: "请输入客服链接", trigger: "blur" },
        ],
        brokerStartDate: [
          { required: true, message: "请选择工作开始时间", trigger: "blur" },
        ],
      },
      // 下拉筛选
      Extended: [
        {
          dkey: "全部用户",
          dvalue: null,
        },
        {
          dkey: "平台管理员",
          dvalue: 2,
        }
      ],
      carriageName: null,

      // 会员列表
      listhy:[
        {
          dkey:'全部用户',
          dvalue:null
        },
        {
          dkey:'到期会员',
          dvalue:0
        },
        {
          dkey:'供应会员',
          dvalue:1
        },
        {
          dkey:'甄选会员',
          dvalue:3
        },
      ],
      hyvalue:null,
      // 会员弹框
      vipTimArray:[],
      vipshow: false,
      vipusertime: {
        memberType:1,
        id:'',
        activateType:'1',
        phone:'',
        andtime:'',
        expirationTime:'',
      },
      // 用户编辑
      editphow: false,
      edituser: {},
      euserrules: {
        username: [{ required: true, message: "请输入姓名", trigger: "blur" }],
      },
      // 身份信息
      sel: 1, // 选项卡默认显示tab
      tabs: [
        { id: 1, name: "洗煤厂", wareldenShow: true },
        { id: 2, name: "贸易商", wareldenShow: true },
        { id: 3, name: "焦化厂", wareldenShow: true },
        { id: 4, name: "经纪人", wareldenShow: true },
        { id: 5, name: "煤矿", wareldenShow: true },
        { id: 6, name: "钢厂", wareldenShow: true },
        { id: 7, name: "金融企业", wareldenShow: true },
        { id: 8, name: "平台经纪人", wareldenShow: false },
        { id: 9, name: "其它", wareldenShow: false },
      ],
      isShow: false, // 查看核实记录
      verifyList: [], // 核实记录数组
       // 下拉筛选
      effSelect: [
        {
          dkey: "全部用户",
          dvalue: null,
        },
        {
          dkey: "有效用户",
          dvalue: true,
        }
      ],
      effective:null, // 有效用户
      editIndex:null,
      userIdentity:null, // 用户身份
    };
  },
  created() {},
  mounted() {
    this.purTitle = this.$route.meta.title;
    this.logo = sessionStorage.getItem("companyName"); // 企业名称
    this.http = this.$store.state.http; // IP地址
    this.cloudhttp = this.$store.state.cloudhttp; // IP地址
    this.getpurcon();
    this.vipusertime.activateType = '1';
    this.yeartime();
  },
  methods: {
    // 注册新用户
    onSubmit(){
      if(/^1\d{10}$/.test(this.formInline.phone)){
        this.api.authAll.add({
          channel:2,
          level:1,
          accountType:1,
          phone:this.formInline.phone,
          account:this.formInline.phone,
          password:this.formInline.phone,
        }).then(res=>{
          if(res.data.code == 200){
            this.$message({
                  message: '注册账号成功！',
                  type: 'success'
              });
            this.formInline = {};
          }
        })
      }else{
        this.$message({
            message: '请输入正确的手机号！',
            type: 'warning'
        });
        this.formInline = {};
      }
    },  
    // 年时间
    yeartime(){
      // 年度会员
      // 获取当前日期
      const today = new Date();
      // 获取当前年份
      const currentYear = today.getFullYear();
      // 获取当前月份（0-11，所以实际月份需要+1）
      const currentMonth = today.getMonth() + 1;
      // 获取当前日期（1-31）
      const currentDay = today.getDate();
      const currentHour= today.getHours() < 10 ? '0' + today.getHours() : today.getHours();
      const currentMinute= today.getMinutes() < 10 ? '0' + today.getMinutes() : today.getMinutes();
      const currentSecond= today.getSeconds() < 10 ? '0' + today.getSeconds() : today.getSeconds();

      // 当前日期的字符串表示
      const currentDateStr = `${currentYear}-${currentMonth.toString().padStart(2, '0')}-${currentDay.toString().padStart(2, '0')} ${currentHour}:${currentMinute}:${currentSecond}`;
      // 一年后的今天
      today.setFullYear(currentYear + 1); // 设置为一年后的今天
      // 一年后的年份
      const nextYear = today.getFullYear();
      // 一年后的月份
      const nextMonth = today.getMonth() + 1;
      // 一年后的日期
      const nextDay = today.getDate();
      const nextHour= today.getHours() < 10 ? '0' + today.getHours() : today.getHours();
      const nextMinute = today.getMinutes() < 10 ? '0' + today.getMinutes() : today.getMinutes();
      const nextSecond = today.getSeconds() < 10 ? '0' + today.getSeconds() : today.getSeconds();
      // 一年后今天的字符串表示
      const nextYearDateStr = `${nextYear}-${nextMonth.toString().padStart(2, '0')}-${nextDay.toString().padStart(2, '0')} ${nextHour}:${nextMinute}:${nextSecond}`;
      
      this.vipusertime.andtime = currentDateStr; // 开始时间
      // this.vipusertime.expirationTime = nextYearDateStr + ' 00:00:00'; ; // 结束时间
      this.vipusertime.expirationTime = nextYearDateStr
      this.vipTimArray = [currentDateStr,nextYearDateStr]
    },
    // 年度 月度会员选择
    viptime(e){
      this.vipusertime.activateType = e;
      if(e == 1){
        this.yeartime();
      }else{
        // 获取当前日期
        const today = new Date();
        // 月度会员
        const year = today.getFullYear(); // 当前年份
        const month = today.getMonth() + 1; // 当前月份（0-11），加1因为getMonth()返回0-11表示1-12月
        const day = today.getDate(); // 当前日（1-31）
        const hour= today.getHours() < 10 ? '0' + today.getHours() : today.getHours();
        const minute= today.getMinutes() < 10 ? '0' + today.getMinutes() : today.getMinutes();
        const second= today.getSeconds() < 10 ? '0' + today.getSeconds() : today.getSeconds();
        // 当前日期
        const currentDate = `${year}-${month.toString().padStart(2, '0')}-${day.toString().padStart(2, '0')} ${hour}:${minute}:${second}`;
        // 一个月后的今天
        today.setMonth(today.getMonth() + 1); // 在当前日期上加一个月
        const monthLaterDate = `${today.getFullYear()}-${(today.getMonth() + 1).toString().padStart(2, '0')}-${today.getDate().toString().padStart(2, '0')} ${hour}:${minute}:${second}`;
        console.log(monthLaterDate,'monthLaterDate');
        this.vipusertime.andtime = currentDate; // 开始时间
        // this.vipusertime.expirationTime = monthLaterDate + ' 00:00:00'; // 结束时间
        this.vipusertime.expirationTime = monthLaterDate
        this.vipTimArray = [currentDate,monthLaterDate]
        // console.log('当前日期:', currentDate);
        // console.log('一个月后的今天:', monthLaterDate);
      }

    },
    changeTime(time){
        // 年度会员
        // 获取当前日期
        const today = new Date(time);
        // 获取当前年份
        const currentYear = today.getFullYear();
        // 获取当前月份（0-11，所以实际月份需要+1）
        const currentMonth = today.getMonth() + 1;
        // 获取当前日期（1-31）
        const currentDay = today.getDate();
        const currentHour= today.getHours() < 10 ? '0' + today.getHours() : today.getHours();
        const currentMinute= today.getMinutes() < 10 ? '0' + today.getMinutes() : today.getMinutes();
        const currentSecond= today.getSeconds() < 10 ? '0' + today.getSeconds() : today.getSeconds();
        // 当前日期的字符串表示
        const currentDateStr = `${currentYear}-${currentMonth.toString().padStart(2, '0')}-${currentDay.toString().padStart(2, '0')} ${currentHour}:${currentMinute}:${currentSecond}`;
        this.vipusertime.expirationTime = currentDateStr
    },
    // 开通VIP会员
    openVip(data){
      // console.log(data);
      this.vipshow = true;
      this.vipusertime.id = data.id;
      this.vipusertime.phone = data.phone;
      this.vipusertime.activateType = '1';
      this.yeartime();
    },
    // 取消会员
    noVip(data){
      this.api.authAll.updateById({
        id:data.id,
        memberType:0
      })
      .then(res=>{
          if(res.data.code == 200){
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
            this.$message({
                message: '取消会员成功！',
                type: 'success'
            });
          }
      })
    },
    // 会员弹框取消
    vipClose(){
      this.vipshow = false;
      this.vipusertime.id = '';
      this.vipusertime.phone = '';
      this.vipusertime.activateType = '1';
      this.vipusertime.andtime = ''; // 开始时间
      this.vipusertime.expirationTime = ''; // 结束时间
    },
    // 确认开通会员
    editvipbmitm(){
      console.log(this.vipusertime.expirationTime,'vipTimArray');
      this.vipusertime.activateType = this.vipusertime.activateType*1;
      // this.vipusertime.expirationTime = this.vipTimArray[1]  + ' 00:00:00';
      this.api.authAll.updateById(this.vipusertime)
      .then(res=>{
        if(res.data.code == 200){
          this.vipshow = false;
          this.vipusertime = {
            memberType:1,
            id:'',
            activateType:'1',
            phone:'',
            andtime:'',
            expirationTime:'',
          }
          this.tableData = [];
          this.getpurcon();
          this.current = 1;
          this.$message({
              message: '开通会员成功！',
              type: 'success'
          });
        }
      })
    },
    // 选中筛选会员
    selectIhy(val){
      this.hyvalue = val
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 设为管理员
    uplserc(data){
      this.api.authAll.updateById({
        id:data.id,
        accountType:2,
        level:2
      }).then(res=>{
        if(res.data.code == 200){
          this.tableData = [];
          this.getpurcon();
          this.current = 1;
          this.$message({
              message: '设为熊猫煤仓管理员成功，请提醒用户退出重新登录！',
              type: 'success'
          });
        }
      })
    },
    // 取消管理员
    nouplserc(data){
      this.api.authAll.updateById({
        id:data.id,
        accountType:1,
        level:1
      }).then(res=>{
        if(res.data.code == 200){
          this.tableData = [];
          this.getpurcon();
          this.current = 1;
          this.$message({
              message: '已取消熊猫煤仓管理员！',
              type: 'success'
          });
        }
      })
    },
    // 注册日期搜搜
    usrTie() {
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 拨打次数列表
    girtusr(proId, targetUserPhone) {
      console.log(proId, targetUserPhone, this.uid);
      this.api.phonelog
        .list({
          call_user_id: this.uid,
          pro_id: proId,
          target_user_phone: targetUserPhone,
        })
        .then((res) => {
          if (res.data.code == 200) {
            this.gridData = res.data.data;
          }
          console.log(res);
        });
    },
    handleClick(tab) {
      if(tab.name == 1){
        this.vipdata = [];
        this.datavips();
      }else if(tab.name == 2){
        this.seuserData = [];
        this.datauser();
      }
    },
    // 访问数据抽屉
    userchange(data) {
      this.uid = data.id;
      this.druserer = true;
      this.vipdata = [];
      this.datavips();
    },
    // 煤种详情页访问数据
    datauser(){
      this.api.busprolog
      .selVisisted({
          userId: this.uid,
          startTime: this.TimeArrays != null ? this.TimeArrays[0] : "",
          endTime: this.TimeArrays != null ? this.TimeArrays[1] : "",
        })
      .then((res) => {
        if (res.data.code == 200) {
          this.seuserData = res.data.data;
        }
      });
    },
    datavips(){
      this.api.busprolog.selMemberData({userId:this.uid})
      .then(res=>{
        if (res.data.code == 200) {
          this.vipdata.push(res.data.data);
          console.log(this.vipdata)
        }
      })
    },
    ImClose() {
      this.druserer = false;
      this.seuserData = [];
      this.vipdata = [];
      this.TimeArrays = [];
      this.activeName = '1';
    },
    // 日期搜索
    Times() {
      if(this.activeName == 1){
        this.vipdata = [];
        this.datavips();
      }else if(this.activeName == 2){
        this.seuserData = [];
        this.datauser();
      }
    },

    // 身份信息
    selectss(item) {
      this.sel = item.id;
    },
    // 编辑用户
    Editserc(data,index) {
      // console.log(data,index)
      this.editphow = true;
      this.edituser = data;
      this.sel = data.wareIdentity;
      this.editIndex = index;
    },
    // 编辑用户取消
    editphowClose() {
      this.editphow = false;
      this.edituser = {};
    },
    // 取消编辑用户信息提交
    edituserForm() {
      this.editphow = false;
      this.edituser = {};
    },
    // 编辑用户提交
    edituserbmitm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.edituser.wareIdentity = this.sel;
          this.api.authAll.updateById(this.edituser).then((res) => {
            if (res.data.code == 200) {
              this.$message({
                type: "success",
                message: "编辑成功!",
              });
              // console.log(this.edituser,'this.edituser');
              this.edituser = {};
              this.tableData = [];
              this.getpurcon();
              this.current = 1;
              // this.tableData[this.editIndex] = {
              //   ...this.tableData[this.editIndex],
              //   username:this.edituser.username,
              //   wareIdentity:this.edituser.wareIdentity,
              //   followPerson:this.edituser.followPerson,
              //   followStatus:this.edituser.followStatus
              // }
              this.editphow = false;
            }
          });
        }
      });
    },

    // 取消经纪人
    delPurs(data) {
      // console.log(data);
      this.$confirm("此操作将取消经纪人配置, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.api.authAll.cancelBuser({ id: data.id }).then((res) => {
            if (res.data.code == 200) {
              this.$message({
                type: "success",
                message: "取消经纪人成功!",
              });
              this.tableData = [];
              this.getpurcon();
              this.current = 1;
            }
          });
        })
        .catch(() => {
          this.$message({
            type: "info",
            message: "已保留经纪人！",
          });
        });
    },
    // 兑换记录抽屉
    exchange(id) {
      this.api.buslog
        .exchangeRecords({ userId: id, searchStr: "" })
        .then((res) => {
          if (res.data.code == 200) {
            this.drawer = true;
            this.seeData = res.data.data;
          }
        });
    },

    // 编辑经纪人取消
    editpaymentShowClose() {
      this.editpaymentShow = false;
      this.editData = {};
    },
    // 取消编辑经纪人信息提交
    editpaymentresetForm() {
      this.editpaymentShow = false;
      this.editData = {};
    },
    // 编辑经纪人提交
    editpaymentsubmitm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.editData.wareIdentity = 8;
          this.api.authAll.updateById(this.editData).then((res) => {
            if (res.data.code == 200) {
              this.$message({
                type: "success",
                message: "设置经纪人成功!",
              });
              this.editData = {};
              this.tableData = [];
              this.getpurcon();
              this.current = 1;
              this.editpaymentShow = false;
            }
          });
        }
      });
    },
    // 设置经纪人按钮
    EditPurs(data) {
      // console.log(data.id);
      this.editpaymentShow = true;
      this.editData = JSON.parse(JSON.stringify(data));
      
    },

    // 表格数据
    getpurcon() {
      this.loading = true;
      // 传参
      var listByPageData = {
        memberType:this.hyvalue,
        searchStr: this.purSearch,
        level: this.carriageName == null ? null : this.carriageName,
        startTime: this.usrTimArray != null ? this.usrTimArray[0] : "",
        endTime: this.usrTimArray != null ? this.usrTimArray[1] : "",
        channel: 2,
        isCtiExDate:true,
        hasPhone:this.effective,
        page: {
          current: this.current,
          size: this.size,
        },
        wareIdentity:this.userIdentity
      };
      // 渲染表格
      this.api.authAll.listByPage(listByPageData).then((res) => {
        if (res.data.code == 200) {
          this.loading = false;
          if (res.data.data.records != "") {
            this.tableData = res.data.data.records;
            this.total = res.data.data.total;
          } else {
            this.dataAva = true;
          }
        }
      });
    },
    // 搜索功能
    search() {
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 查询经纪人
    prologistics() {
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 查看核实信息
    verifychange(data) {
      this.isShow = true;
      this.api.verifyInfo
        .all({
          searchStr: "",
          verifyPeopleId: data.id,
        })
        .then((res) => {
          if (res.data.code == 200) {
            this.verifyList = res.data.data;
          }
        });
    },
    selectUser(e){
      this.effective = e
      this.getpurcon();
    },
    // 注销账号
    delserc(data){
      this.$confirm('此操作将永久注销该账号并清空账号下所有数据，请谨慎操作！', '提示', {
        confirmButtonText: '注销',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.api.authAll.logOffUser({id:data.id})
        .then(res=>{
          if(res.data.code == 200){
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
            this.$message({
              type: 'success',
              message: '注销账号成功!'
            });
            
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消注销操作！'
        });          
      });
    },
    // 筛选身份
    selectIdentity(e){
      this.userIdentity = e
      this.getpurcon();
    },
  },
};
</script>
  <style lang="less" scoped>
// 身份选择样式开始
.labfl {
  display: flex;
  flex-wrap: wrap;
  .totalTab {
    padding: 6px;
    font-size: 14px;
    line-height: 20px;
    color: #333;
    display: inline-block;
    margin: 0 8px 8px 0;
    background: #f9f9f9;
    border-radius: 4px;
    text-align: center;
    border: 1px solid #f9f9f9;
    box-sizing: border-box;
  }
  .active {
    color: #e2231a;
    border: 1px solid #e2231a;
    background: #fff3f2;
    box-sizing: border-box;
  }
}
// 身份选择样式结束
// 访问数据样式开始
.usertitle {
  padding: 16px;
  box-sizing: border-box;
  span {
    font-size: 20px;
    line-height: 32px;
    color: #333;
    font-weight: 700;
  }
  .usrscreen {
    width: 70%;
    float: right;
    p {
      font-size: 16px;
      line-height: 32px;
      color: #333;
      float: right;
      margin-right: 6px;
    }
  }
}
.usrifbox {
  width: 100%;
  height: calc(100vh - 160px);
  padding-top: 80px;
  box-sizing: border-box;
}
.usrclock {
  width: 100%;
  height: calc(100vh - 160px);
  overflow: hidden;
  overflow-y: scroll;
  padding: 16px;
  box-sizing: border-box;

  .usrbox {
    width: 100%;
    height: 260px;
    background: #f7f7f7;
    border-radius: 8px;
    margin-bottom: 6px;
    padding: 18px 16px;
    box-sizing: border-box;
    display: flex;
    flex-wrap: wrap;
    .usrh1 {
      margin-right: 10px;
      font-size: 20px;
      line-height: 32px;
      color: #333;
    }
    .usrh2 {
      flex: 1;
      font-size: 20px;
      line-height: 32px;
      color: #333;
    }
    .usrh3 {
      width: 100%;
      font-size: 14px;
      line-height: 32px;
      color: #333;
    }
    .usrfrom {
      width: 100%;
      height: 160px;
      background: #ff6700;
    }
  }
}
// 访问数据样式结束
// 兑换记录样式开始
.blocktitle {
  font-size: 24px;
  line-height: 48px;
  font-weight: 600;
  color: #333;
  padding: 0 32px;
  box-sizing: border-box;
}
.block {
  width: 100%;
  padding: 32px;
  box-sizing: border-box;
}
// 兑换记录样式结束
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
    .el-select {
      width: 8%;
      margin-right: 10px;
      float: right;
    }
    .el-input {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .el-date-picker {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .el-button {
      margin-right: 10px;
      float: right;
    }
    .el-date-picker {
      width: 10%;
      float: right;
      margin-right: 10px;
    }
    .selTit{
      height:32px;
      float:right;
      line-height:32px;
      font-size: 14px;
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
      left: 25%;
      transform: translate(0, -50%);
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
  // 分页
  .pur-pages {
    width: 100%;
    height: 6%;
    padding: 6px 0 0 12px;
    box-sizing: border-box;
    span {
      font-size: 14px;
      line-height: 33px;
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
  .xia {
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
.see-file {
  padding: 15px 0;
  overflow: hidden;
  span {
    font-size: 14px;
    color: #666;
  }
  dd {
    font-size: 14px;
    line-height: 34px;
    color: #333;
    a {
      text-decoration: none;
      color: #333;
    }
    .active {
      color: #ea222e;
    }
  }
}
// 查看抽屉样式结束

// 查看核实信息抽屉
.draBox {
  height: 100%;
  padding: 40px 32px 0px;

  .draTitle {
    margin-left: 24px;
    margin-bottom: 18px;
    font-size: 18px;
    color: #333;
    font-weight: 700;
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

.company {
  margin-left: 24px;
  margin-bottom: 14px;
  font-size: 12px;
  color: #666;
}

.contentItem {
  .timeNode {
    height: 14px;

    .itemIcon {
      width: 14px;
      height: 14px;
      float: left;
      margin-right: 10px;
      background: linear-gradient(90deg, #e2231a 0%, #cf2c2c 100%);
      border-radius: 50%;
      text-align: center;
      line-height: 14px;

      img {
        width: 10px;
        height: auto;
      }
    }
    .infoTime {
      float: left;
      font-size: 12px;
      color: #666;
    }
  }

  .itemInfo {
    padding-left: 6px;
    padding-bottom: 18px;

    

    .txtBox {
      height: auto;
      border-left: 2px solid #D9D9D9;
      padding: 6px 0 0 18px;

      .adressTxt {
        margin-bottom: 14px;

        img {
          width: 12px;
          height: 12px;
          margin-left: 4px;
        }
      }

      .comment {
        width: 600px;
        font-size: 14px;
        font-weight: 700;
        color: #333;
        word-wrap: break-word;
      }
    }
  }
}
</style>