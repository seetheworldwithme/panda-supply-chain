<template>
  <!-- 煤种列表页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <el-button size="small" @click="opensjPur">上架云仓</el-button>
        <el-button size="small" @click="openxjPur">下架云仓</el-button>
        <el-button size="small" @click="Importdailog" v-if="$route.meta.import">一键导入煤种信息</el-button>
        <el-button size="small" @click="busrtdailog">一键配置经纪人</el-button>
        <el-button size="small" @click="exportlist">导出表格</el-button>
        <el-input
          style="width:14%;"
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch">
        </el-input>
      </el-col>
      <el-col :span="24">
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">企业是否核实</div>
          <el-select v-model="verifyStatus" size="small" @change="verifyChange">
              <el-option
              v-for="item in verifyOption"
              :key="item.value"
              :label="item.dkey"
              :value="item.dvalue">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">是否关联企业</div>
          <el-select v-model="caeName" size="small" @change="prtics">
              <el-option
              v-for="item in ssended"
              :key="item.value"
              :label="item.dkey"
              :value="item.dvalue">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">云仓上下架</div>
          <el-select v-model="carriageName" size="small" @change="prologistics">
              <el-option
              v-for="item in Extended"
              :key="item.value"
              :label="item.dkey"
              :value="item.dvalue">
              </el-option>
          </el-select>
        </el-col>
        <!-- 多功能筛选 -->
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">公司性质</div>
          <el-select v-model="carrnatur" placeholder="请选择公司性质" size="small" clearable @change="naturstics" style="width:100%;">
            <el-option
            v-for="item in natureall"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue">
            </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">煤种经纪人</div>
          <el-select v-model="brokernatur" placeholder="请选择经纪人" size="small" clearable @change="brokertics" style="width:100%;">
            <el-option
            v-for="item in brokerall"
            :key="item.id"
            :label="item.username"
            :value="item.id">
            </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">省</div>
          <el-select @change="enterprises" size="small" v-model="provinceCode" clearable placeholder="请选择省" style="width:100%;">
              <el-option
              v-for="item in provinceOptions"
              :key="item.code"
              :label="item.name"
              :value="item.code">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">市</div>
            <el-select v-model="cityCode" placeholder="请选择市" size="small" clearable style="width:100%;" @change="citys($event)">
                <el-option
                v-for="item in cityOptions"
                :key="item.value"
                :label="item.name"
                :value="item.code">
                </el-option>
            </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">区，县</div>
          <el-select v-model="areaCode" placeholder="请选择县 / 区" size="small" clearable style="width:100%;" @change="countys($event)">
              <el-option
              v-for="item in countyOptions"
              :key="item.value"
              :label="item.name"
              :value="item.code">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="2" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">是否平台认证</div>
          <el-select v-model="attestation" placeholder="请选择是否平台认证" size="small" clearable @change="attestationstics" style="width:100%;">
            <el-option
            v-for="item in attestationall"
            :key="item.value"
            :label="item.dkey"
            :value="item.dvalue">
            </el-option>
          </el-select>
        </el-col>
        <el-col :span="3" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">煤炭分类</div>
          <el-select ref="bigref" v-model="classify" clearable size="small" @change="changeClassify" placeholder="请选择煤炭分类">
              <el-option
              v-for="item in classifydata"
              :key="item.id"
              :label="item.productName"
              :value="item.id">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="3" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">煤炭品种</div>
          <el-select ref="proRef" v-model="breed" clearable size="small" @change="changeBreed" placeholder="请选择煤炭品种">
              <el-option
              v-for="item in breedmlns"
              :key="item.id"
              :label="item.productName"
              :value="item.id">
              </el-option>
          </el-select>
        </el-col>
        <el-col :span="4" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">创建煤种时间</div>
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
          style="width:100%;">
          </el-date-picker>
        </el-col>
        <el-col :span="4" style="padding-right:6px;box-sizing: border-box;">
          <div class="seaclass">更新煤种日期</div>
          <el-date-picker
          size="small"
          @change="renewTime"
          v-model="renewTimeArray"
          type="daterange"
          unlink-panels
          range-separator="~"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="yyyy-MM-dd"
          style="width:100%;">
          </el-date-picker>
        </el-col>
        <el-col :span="4" style="padding-right:6px;box-sizing:border-box;">
          <div class="seaclass">更新煤种价格日期（出厂价）</div>
          <el-date-picker
          size="small"
          @change="moenyTime"
          v-model="moenyTimeArray"
          type="daterange"
          unlink-panels
          range-separator="~"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="yyyy-MM-dd"
          style="width:100%;">
          </el-date-picker>
        </el-col>
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
            <!-- <el-table-column prop="id" label="煤种索引">
            </el-table-column> -->
            <el-table-column type="index" width="50" label="序号" fixed="left"></el-table-column>
            <el-table-column
              label="操作" width="80" fixed="left" v-if="$route.meta.edit">
              <template slot-scope="scope">
                <el-button v-if="scope.row.status == 0" @click="EditPurs(scope.row)" size="mini" type="warning" plain>编辑</el-button>
                <el-button v-else-if="scope.row.status == 2" @click="EditPurs(scope.row)" size="mini" type="danger" plain>编辑</el-button>
                <el-button v-else-if="scope.row.status == 3" @click="EditPurs(scope.row)" size="mini" type="info" plain>查看</el-button>
                <el-button v-else @click="EditPurs(scope.row)" size="mini" type="success" plain>编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column label="发布渠道" :show-overflow-tooltip="true" fixed="left">
              <template slot-scope="scope">
                <span v-if="scope.row.dataResource == 1">小程序</span>
                <span v-else-if="scope.row.dataResource == 2">官网</span>
                <span v-else-if="scope.row.dataResource == 3">APP</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column prop="warehouseNb" label="云仓仓号" fixed="left" :show-overflow-tooltip="true" width="130">
            </el-table-column>
            <el-table-column label="煤种访问人数" fixed="left" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span @click="acesPurs(scope.row)" style="color:#CF1724;cursor: pointer;">{{scope.row.vtPeopCount}} 人</span>
              </template>
            </el-table-column>
            <el-table-column prop="warseCode" label="云仓编号" :show-overflow-tooltip="true" width="80" fixed="left">
            </el-table-column>
            <el-table-column prop="enterpriseName" label="企业名称" fixed="left" :show-overflow-tooltip="true" width="140">
            </el-table-column>
            <el-table-column label="公司性质" :show-overflow-tooltip="true" fixed="left">
            <template slot-scope="scope">
              <span v-if="scope.row.wareIdentity == 1">洗煤厂</span>
              <span v-else-if="scope.row.wareIdentity == 2">贸易商</span>
              <span v-else-if="scope.row.wareIdentity == 3">焦化厂</span>
              <span v-else-if="scope.row.wareIdentity == 4">经纪人</span>
              <span v-else-if="scope.row.wareIdentity == 5">煤矿</span>
              <span v-else-if="scope.row.wareIdentity == 6">钢厂</span>
              <span v-else-if="scope.row.wareIdentity == 7">金融企业</span>
              <span v-else-if="scope.row.wareIdentity == 8">平台经纪人</span>
              <span v-else-if="scope.row.wareIdentity == 9">其他</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column prop="businessName" label="煤种经纪人" :show-overflow-tooltip="true" fixed="left">
          </el-table-column>
          <el-table-column prop="businessPhone" label="煤种经纪人电话" :show-overflow-tooltip="true" width="120" fixed="left">
          </el-table-column>
          <el-table-column label="是否上仓" :show-overflow-tooltip="true" fixed="left">
              <template slot-scope="scope">
                <div v-if="scope.row.status == 0">
                  <div class="statsBox" style="background: #FB8138;box-shadow: 0px 4px 4px #FB8138 ;"></div><span> 待审核</span>
                </div>
                <div v-if="scope.row.status == 1">
                  <div class="statsBox" style="background: #46C7A2;box-shadow: 0px 4px 4px #46C7A2;"></div><span> 已上架</span>
                </div>
                <div v-if="scope.row.status == 2">
                  <div class="statsBox" style="background: #CF1724;box-shadow: 0px 4px 4px #F7D0D3;"></div><span> 已拒绝</span>
                </div>
                <div v-if="scope.row.status == 3">
                  <div class="statsBox" style="background: #BFBFBF;box-shadow: 0px 4px 4px #BFBFBF;"></div><span> 已下架</span>
                </div>
              </template>
            </el-table-column>
          
          <el-table-column prop="areaBuserName" label="地区经纪人姓名" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="areaBuserPhone" label="地区经纪人电话" :show-overflow-tooltip="true" width="120">
          </el-table-column>
          <el-table-column prop="createPhone" label="发布人电话" :show-overflow-tooltip="true" width="120">
          </el-table-column>
          <el-table-column prop="coalRelName" label="联系人姓名（企业）" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="coalRelPhone" label="联系人电话（企业）" :show-overflow-tooltip="true" width="110">
          </el-table-column>
          <el-table-column prop="etcRelPhone" label="数据分配（电话）" :show-overflow-tooltip="true" width="110">
          </el-table-column>
            <el-table-column label="所属省" :show-overflow-tooltip="true" width="100" prop="proName"></el-table-column>
            <el-table-column label="所属市" :show-overflow-tooltip="true" width="100" prop="cityName"></el-table-column>
            <el-table-column label="所属区/县" :show-overflow-tooltip="true" width="100" prop="areaName"></el-table-column>
            <el-table-column label="核实信息">
              <template slot-scope="scope">
                <span style="color: #cf1724; cursor: pointer" @click="getVerifyInfo(scope.row)" v-if="scope.row.isVerify>0">详情</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="品名（一级分类）" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span  v-if="scope.row.fstName">{{scope.row.fstName}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column prop="secName" label="煤炭分类（二级分类）" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span  v-if="scope.row.secName">{{scope.row.secName}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column prop="thName" label="煤炭品种（三级分类）" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span  v-if="scope.row.thName">{{scope.row.thName}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="累计成交" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span v-if="scope.row.sumTonnageTraded">{{scope.row.sumTonnageTraded}} 吨</span>
                <span v-else>0 吨</span>
              </template>
            </el-table-column> -->
            <el-table-column label="平台可供应量" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span  v-if="scope.row.availableTonnage">{{scope.row.availableTonnage}} 吨</span>
                <span v-else>0 吨</span>
              </template>
            </el-table-column>
            <el-table-column label="现有库存" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span  v-if="scope.row.bleTonnage">{{scope.row.bleTonnage}} 吨</span>
                <span v-else>0 吨</span>
              </template>
            </el-table-column>
            <!-- <el-table-column prop="dataProvider" label="发布人" :show-overflow-tooltip="true" width="100">
            </el-table-column> -->
            <el-table-column label="出厂现汇含税价" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.exFactoryPrice">{{scope.row.exFactoryPrice}} 元 (出厂价)</span>
                <span v-if="scope.row.platePrice">{{scope.row.platePrice}} 元 (车板价)</span>
                <span v-if="!scope.row.platePrice && !scope.row.exFactoryPrice">--</span>
              </template>
            </el-table-column>
            <el-table-column label="煤种创建时间" :show-overflow-tooltip="true" width="200">
                <template slot-scope="scope">
                    <span v-if="scope.row.createTime">{{scope.row.createTime}}</span>
                </template>
            </el-table-column>
            <el-table-column label="煤种更新时间" :show-overflow-tooltip="true" width="200">
                <template slot-scope="scope">
                    <span v-if="scope.row.updateTime">{{scope.row.updateTime}}</span>
                </template>
            </el-table-column>
            <el-table-column label="价格更新时间" :show-overflow-tooltip="true" width="200">
                <template slot-scope="scope">
                    <p v-if="scope.row.fpriceUpdTime">出厂：{{scope.row.fpriceUpdTime}}</p>
                    <p v-if="scope.row.ppriceUpdTime">车板：{{scope.row.ppriceUpdTime}}</p>
                </template>
            </el-table-column>
            <el-table-column prop="downTime" label="煤种下架时间" :show-overflow-tooltip="true" width="160">
              <template slot-scope="scope">
                <span v-if="scope.row.downTime">{{scope.row.downTime}}</span>
                <span v-else>--</span>
              </template>
          </el-table-column>
            <el-table-column label="灰分(Ad)%≤" :show-overflow-tooltip="true">
                <template slot-scope="scope">
                    <span v-if="scope.row.maxAshContent">{{scope.row.maxAshContent}}</span>
                    <span v-else>--</span>
                </template>
            </el-table-column>
            <el-table-column label="水分(Mt)%≤" :show-overflow-tooltip="true">
            <template slot-scope="scope">
                <span v-if="scope.row.maxWaterContent">{{scope.row.maxWaterContent}}</span>
                <span v-else>--</span>
            </template>
            </el-table-column>
            <el-table-column label="硫(St,d)%≤" :show-overflow-tooltip="true">
            <template slot-scope="scope">
                <span v-if="scope.row.maxStd">{{scope.row.maxStd}}</span>
                <span v-else>--</span>
            </template>
            </el-table-column>
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
            <el-table-column label="哈氏可磨指数(HGI)≥" :show-overflow-tooltip="true">
            <template slot-scope="scope">
                <span v-if="scope.row.minHgi">{{scope.row.minHgi}}</span>
                <span v-else>--</span>
            </template>
            </el-table-column>
            <el-table-column label="岩相标准差S≤" :show-overflow-tooltip="true">
            <template slot-scope="scope">
                <span v-if="scope.row.lithofacies">{{scope.row.lithofacies}}</span>
                <span v-else>--</span>
            </template>
            </el-table-column>
            <el-table-column label="精煤回收率(%)" :show-overflow-tooltip="true">
            <template slot-scope="scope">
                <span v-if="scope.row.cleanCoalRecoveryRate">{{scope.row.cleanCoalRecoveryRate}}</span>
                <span v-else>--</span>
            </template>
            </el-table-column>
            <el-table-column label="热值≥" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                  <span v-if="scope.row.calorificValue">{{scope.row.calorificValue}}</span>
                  <span v-else>--</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="视频链接" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                  <span v-if="scope.row.videoList">{{http}}/saas-common/upload/getVideo?fileName={{scope.row.videoList}}&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f</span>
                  <span v-else>--</span>
              </template>
            </el-table-column> -->
            <el-table-column label="热门推荐" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span v-if="scope.row.isRecommende">开启</span>
                <span v-else>关闭</span>
              </template>
            </el-table-column>
            <el-table-column label="平台认证" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span v-if="scope.row.isPlatformCertification">已认证</span>
                <span v-else>未认证</span>
              </template>
            </el-table-column>
            <el-table-column label="平台核实" :show-overflow-tooltip="true" width="100">
              <template slot-scope="scope">
                <span v-if="scope.row.isRealCheck">已核实</span>
                <span v-else>未核实</span>
              </template>
            </el-table-column>
            <el-table-column label="备注" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.remark">{{scope.row.remark}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="提货地" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.remark">{{scope.row.takePlace}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="发站地" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.remark">{{scope.row.departureStation}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
            <el-table-column label="下架原因" :show-overflow-tooltip="true" width="140">
              <template slot-scope="scope">
                <span v-if="scope.row.delReason">{{scope.row.delReason}}</span>
                <span v-else>--</span>
              </template>
            </el-table-column>
			<!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="@/assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                <el-button size="small" @click="Importdailog" v-if="$route.meta.import">一键导入煤种信息</el-button>
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
      <!-- 编辑终端信息弹窗 -->
      <el-dialog center
        :closeOnClickModal="false"
        title="煤种信息详情"
        :visible.sync="editpaymentShow"
        width="98%"
        :before-close="editpaymentShowClose">
          <el-form style="width:95%;" :model="editData" :rules="editpaymentrules" ref="AddData" size="small" label-width="140px" v-loading="mzloading" element-loading-text = "更新中...">
            <el-row>
              <el-col :span="6">
                  <el-form-item label="身份选择">
                    <div class="labfl">
                      <div class="totalTab" :class="{ active: item.id === sel }" v-for="item in tabs" :key="item.id" @click="selectss(item)" style="cursor: pointer;">
                        <span>{{item.name}}</span>
                      </div>
                    </div>
                  </el-form-item>
                  <el-form-item label="煤种区别">
                    <div class="labfl">
                      <div class="totalTab" :class="{ active: v.id === selone }" v-for="v in tabones" :key="v.id" @click="seleones(v)" style="cursor: pointer;">{{v.name}}</div>
                    </div>
                  </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="煤炭分类" prop="bigProId">
                    <el-select ref="bigref" v-model="editData.bigProId" size="small" @change="changeHandler" placeholder="请选择煤炭分类">
                        <el-option
                        v-for="item in columnsdata"
                        :key="item.id"
                        :label="item.productName"
                        :value="item.id">
                        </el-option>
                    </el-select>
                  </el-form-item>
                  <el-form-item label="煤炭品种" prop="threeProId">
                    <el-select ref="proRef" v-model="editData.threeProId" size="small" @change="changepro" placeholder="请选择煤炭品种">
                        <el-option
                        v-for="item in bolumlns"
                        :key="item.id"
                        :label="item.productName"
                        :value="item.id">
                        </el-option>
                    </el-select>
                  </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="出厂现汇价" prop="exFactoryPrice">
                  <el-input v-model="editData.exFactoryPrice"></el-input>
                </el-form-item>
                <el-form-item label="车板价格" prop="platePrice">
                  <el-input v-model="editData.platePrice"></el-input>
                </el-form-item>
              </el-col>
              <el-col style="line-height:33px" :span="1">
                <el-tooltip style="margin-left:10px" class="item" effect="dark" :content="exPriceTip">
                  <i class="el-icon-warning"></i>
                </el-tooltip>
              </el-col>
              <el-col :span="4">
                <el-form-item label="平台可供应量">
                  <el-input v-model="editData.availableTonnage"></el-input>
                </el-form-item>
                <el-form-item label="现有库存">
                  <el-input v-model="editData.bleTonnage"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="关联企业">
                    <el-select ref="comref" v-model="editData.companyId" filterable size="small" placeholder="请选择关联企业" @change="compaypro" style="width:100%;">
                        <el-option
                        v-for="item in companydata"
                        :key="item.id"
                        :label="item.enterpriseName"
                        :value="item.id">
                        </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="关联标签"> 
                  <el-select v-model="editData.prags" multiple placeholder="请选择标签" style="width:100%;">
                    <el-option
                      v-for="item in prooptions"
                      :key="item.id"
                      :label="item.signName"
                      :value="item.id">
                    </el-option>
                  </el-select>
                </el-form-item>
                
              </el-col>
              <el-col :span="4">
                <el-form-item label="省" prop="provinceCode">
                  <el-select size="small" disabled v-model="editData.provinceCode" clearable placeholder="请选择省" style="width:100%;" @change="editerprise">
                      <el-option
                      v-for="item in provinceOptions"
                      :key="item.code"
                      :label="item.name"
                      :value="item.code">
                      </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="市" prop="cityCode">
                  <el-select v-model="editData.cityCode" disabled placeholder="请选择市" size="small" clearable style="width:100%;" @change="editcity">
                      <el-option
                      v-for="item in cityOptions"
                      :key="item.value"
                      :label="item.name"
                      :value="item.code">
                      </el-option>
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="1">
                <div style="height:1px"></div>
              </el-col>
              <el-col :span="4">
                <el-form-item label="县 / 区" prop="areaCode">
                    <el-select v-model="editData.areaCode" disabled placeholder="请选择县/区" size="small" clearable style="width:100%;" @change="editcounty">
                      <el-option
                      v-for="item in countyOptions"
                      :key="item.value"
                      :label="item.name"
                      :value="item.code"> 
                      </el-option>
                    </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="5">
                <el-form-item label="发布人手机号">
                  <el-input v-model="editData.createPhone" width="160"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="4">
                  <el-form-item label="企业名称" prop="enterpriseName">
                    <el-input v-model="editData.enterpriseName" disabled></el-input>
                  </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="云仓仓号" prop="warehouseNb">
                  <el-input v-model="editData.warehouseNb" disabled></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="云仓编号" prop="warseCode">
                  <el-input v-model="editData.warseCode"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="煤种编号">
                  <el-input v-model="editData.coalCode"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="煤种宣传语">
                  <el-input v-model="editData.propaganda"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="数据分配（电话）" prop="etcRelPhone">
                  <el-input v-model="editData.etcRelPhone"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="联系人姓名" prop="coalRelName">
                  <el-input v-model="editData.coalRelName"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="联系人电话" prop="coalRelPhone">
                  <el-input v-model="editData.coalRelPhone"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="煤种经纪人姓名" prop="businessName">
                   <el-select v-model="editData.businessName" placeholder="请选择煤种经纪人" @change="selectBroker">
                      <el-option
                        v-for="(item,index) in brokerList"
                        :key="index"
                        :label="item.username"
                        :value="index">
                      </el-option>
                    </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="煤种经纪人电话" prop="businessPhone">
                  <el-input v-model="editData.businessPhone" disabled></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="地区经纪人姓名" prop="areaBuserName">
                  <el-select v-model="editData.areaBuserName" placeholder="请选择地区经纪人" @change="selectAreaBroker">
                      <el-option
                        v-for="(item,index) in brokerList"
                        :key="index"
                        :label="item.username"
                        :value="index">
                      </el-option>
                    </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="地区经纪人电话" prop="areaBuserPhone">
                  <el-input v-model="editData.areaBuserPhone" disabled></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="4">
                <el-form-item label="热门推荐">
                  <el-switch v-model="editData.isRecommende"></el-switch>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="平台认证">
                  <el-switch v-model="editData.isPlatformCertification"></el-switch>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="平台核实">
                  <el-switch v-model="editData.isRealCheck"></el-switch>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="备注">
                  <el-input type="textarea" v-model="editData.remark"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="提货地">
                  <el-input type="textarea" v-model="editData.takePlace" disabled></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="发站地">
                  <el-input type="textarea" v-model="editData.departureStation" disabled></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="3">
                <el-form-item label="硫(St,d)%≤">
                  <el-input v-model="editData.maxStd"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="挥发分(Vdaf)%≤">
                  <el-input v-model="editData.maxVolatiles"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="粘结指数GR.I≥">
                  <el-input v-model="editData.minBondIndex"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="热态CSR≥">
                  <el-input v-model="editData.minCsr"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="热值">
                  <el-input v-model="editData.calorificValue"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="矸石(%)">
                  <el-input v-model="editData.gangue"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="煤泥">
                  <el-input v-model="editData.slurry"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="灰分(Ad)%≤">
                  <el-input v-model="editData.maxAshContent"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="水分(Mt)%≤">
                  <el-input v-model="editData.maxWaterContent"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="胶质层Ymm≥">
                  <el-input v-model="editData.minY"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="岩相标准差S">
                  <el-input v-model="editData.lithofacies"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="哈氏可磨指数(HGI)">
                  <el-input v-model="editData.minHgi"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="精煤回收率(%)">
                  <el-input v-model="editData.cleanCoalRecoveryRate"></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="3">
                <el-form-item label="中煤(%)">
                  <el-input v-model="editData.middlingCoal"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="12">
                <el-form-item label="检测报告">
                  <el-upload
                    :action="http + '/saas-common/upload/batch'"
                    :headers="header"
                    :file-list="jcimgList"
                    name='file'
                    list-type="picture-card"
                    :on-success="addjclist"
                    :on-preview="jcPictureCardPrediv"
                    :on-remove="jcRemove">
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="imgjcogVisible" append-to-body>
                    <img width="100%" :src="jcogImageUrl" alt="">
                  </el-dialog>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="煤棚照片">
                  <el-upload
                    :action="http + '/saas-common/upload/batch'"
                    :headers="header"
                    :file-list="mpimgList"
                    name='file'
                    list-type="picture-card"
                    :on-success="addmplist"
                    :on-preview="mpPictureCardPrediv"
                    :on-remove="mpRemove">
                    <i class="el-icon-plus"></i>
                  </el-upload>
                  <el-dialog :visible.sync="imgmpogVisible" append-to-body>
                    <img width="100%" :src="mpogImageUrl" alt="">
                  </el-dialog>
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label-width="0px">
              <div class="butcls">
                <template v-if="editData.status==2">
                  <el-button size="medium" type="primary" @click="ignoreClick(editData)" style="margin: auto;">重新审核</el-button>
                  <el-button size="medium" @click="editpaymentresetForm">关 闭</el-button>
                </template>
                <template v-else-if="editData.status==0">
                  <el-popover
                    v-loading="mzloading" element-loading-text = "更新中..."
                    placement="left"
                    title="拒绝原因"
                    width="200"
                    trigger="click">
                    <el-input type="textarea" :rows="3" placeholder="请输入拒绝原因" v-model="checkRemark"></el-input>
                    <el-button type="primary" size="medium" @click="jRefuse(editData.id)" style="width:100%;margin-top: 10px;">通知发布人</el-button>
                    <el-button slot="reference" size="medium">审核拒绝</el-button>
                  </el-popover>
                  <el-button size="medium" type="warning" @click="editpaymentsubmitm('AddData')" style="margin: auto;">审核通过</el-button>
                  <el-button type="primary" size="medium" @click="editpaymentresetForm">关 闭</el-button>
                </template>
                <template v-else-if="editData.status==1">
                  <el-button size="medium" type="warning" @click="editpaymentsubmitm('AddData')" style="margin: auto;">保 存</el-button>
                  <el-button  type="primary" size="medium" @click="editpaymentresetForm">关 闭</el-button>
                </template>
                <el-button v-else type="primary" size="medium" @click="editpaymentresetForm">关 闭</el-button>
              </div>
            </el-form-item>
          </el-form>
      </el-dialog>
      <!-- 导入弹窗功能 -->
      <el-dialog
        title="煤种信息导入"
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
                    :action="http + '/saas-warehouse/busProduct/importExcel'"
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
              <a :href="http + '/saas-common/upload/downById?fileName=YCPROMB.xlsx&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'">下载导入模板</a>
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
        <!-- 关联经纪人信息弹窗 -->
        <el-dialog
          :closeOnClickModal="false"
          title="关联经纪人"
          :visible.sync="baymentShow"
          width="25%"
          :before-close="ebymentShowClose">
            <el-form style="width:95%;" :model="associationData" :rules="bpaymentrules" ref="bData" size="small" label-width="140px">
              <el-form-item label="地区经纪人选择">
                  <el-select style="width:100%;" v-model="associationData.areaBuserName" size="small" @change="bgistics">
                    <el-option
                      v-for="item in bended"
                      :key="item.phone"
                      :label="item.username"
                      :value="item.phone">
                      </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="地区经纪人手机号">
                    <el-input v-model="associationData.areaBuserPhone"></el-input>
                </el-form-item>
                <el-form-item label="煤种经纪人选择">
                  <el-select style="width:100%;" v-model="associationData.businessName" size="small" @change="mzistics">
                    <el-option
                      v-for="item in bended"
                      :key="item.phone"
                      :label="item.username"
                      :value="item.phone">
                      </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="煤种经纪人手机号">
                    <el-input v-model="associationData.businessPhone"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="baymentsubmitm('bData')" style="float:right;margin:5% 0 0 5%;">确定</el-button>
                    <el-button @click="bmentresetForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>
        <!-- 关联云仓弹窗 -->
        <el-dialog
        v-loading="mzloading" element-loading-text = "更新中..."
          title="关联云仓号"
          :visible.sync="YuncangShow"
          width="25%">
            <el-form style="width:95%;" :model="YuncangData" :rules="Yuncangrules" ref="yData" size="small" label-width="140px">
              <el-form-item label="分配云仓号" prop="warehouseId">
                  <el-select filterable ref="ycref" style="width:100%;" v-model="YuncangData.warehouseId" size="small" @change="Yuncangstics">
                    <el-option
                      v-for="item in mnsdata"
                      :key="item.id"
                      :label="item.warehouseNb"
                      :value="item.id">
                      </el-option>
                  </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="Yuncangsubmitm('yData')" style="float:right;margin:5% 0 0 5%;">保存</el-button>
                    <el-button @click="YuncangForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>
        <!-- 下架原因弹框 -->
        <el-dialog
        title="下架原因"
        :visible.sync="dialoshelf"
        width="30%"
        :show-close="false">
        <el-input
          type="textarea" size="small"
          placeholder="请输入下架原因" v-model="purshelf">
        </el-input>
        <span slot="footer" class="dialog-footer">
          <el-button @click="shelfon">取消下架</el-button>
          <el-button type="primary" @click="shelfyes">确定下架</el-button>
        </span>
      </el-dialog>
      <!-- 抽屉 查看核实煤仓信息 -->
      <el-drawer :title="logo" :visible.sync="drawer" :with-header="false" size="700px">
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
                <img src="@/assets/zanwushuju.png" alt="">
              </div>
              <div style="font-size: 16px;color: #666;">暂无核实记录</div>
            </div>
        </div>
      </el-drawer>
      <!-- 煤种访问人数弹框开始 -->
      <el-dialog
        title="煤种访问人数"
        :visible.sync="acsdialsible"
        width="40%"
        :before-close="haacsClose">      
        <el-table :data="ascgridData" height="100%" border :header-cell-style="{ background: '#F6F6F6', color: '#333' }">
          <el-table-column property="userPhone" label="访问用户"></el-table-column>
          <el-table-column label="访问时间">
            <template slot-scope="scope">
              <span v-if="scope.row.createTime">{{ scope.row.createTime }}</span>
            </template>
          </el-table-column>
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
          @current-change='mzdata'
          @size-change='mzdata'>
          </el-pagination>
        </div>
      </el-dialog>
      <!-- 煤种访问人数弹框结束 -->
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            // 煤种访问人数
            acsdialsible:false,
            ascgridData:[],
            mzcurrent:1,              // 当前页数
            mzsize:10,                // 显示条数
            mztotal:0,                // 总条数
            mzid:null,
            targetUserId:null,
            // 下架弹框
            dialoshelf:false,
            purshelf:'',
            renewTimeArray:[],        // 更新时间搜索
            moenyTimeArray:[],        // 更新煤种价格时间搜索
            TimeArrays:[],            // 创建时间搜索  
            // 是否平台认证
            attestationall:[
              {
                  dkey:'全部',
                  dvalue:null
              },
              {
                  dkey:'平台认证',
                  dvalue:true
              },
              {
                  dkey:'平台未认证',
                  dvalue:false
              },
            ],
            attestation:null,
            // 煤炭区别
            classify:null,
            classifydata: [], // 煤炭分类二级大类
            breed:null,
            breedmlns:[],    // 煤炭品种三级分类
            // 行政区划
            provinceOptions:[],             // 省
            provinceCode:null,
            provinceName:null,
            cityOptions: [],                // 市
            cityCode:null,
            countyOptions: [],              // 县 区
            areaCode:null,
            // 经纪人
            brokerall:[],
            brokernatur:null,
            // 公司性质
            natureall:[
              {
                  dkey:'全部',
                  dvalue:null
              },
              {
                  dkey:'洗煤厂',
                  dvalue:1
              },
              {
                  dkey:'贸易商',
                  dvalue:2
              },
              {
                  dkey:'焦化厂',
                  dvalue:3
              },
              {
                  dkey:'经纪人',
                  dvalue:4
              },
              {
                  dkey:'煤矿',
                  dvalue:5
              },
              {
                  dkey:'钢厂',
                  dvalue:6
              },
              {
                  dkey:'金融企业',
                  dvalue:7
              },
              {
                  dkey:'平台经纪人',
                  dvalue:8
              },
              {
                  dkey:'其他',
                  dvalue:9
              },
            ],
            carrnatur:null,

          
            ssended:[
            {
                dkey:'全部企业',
                dvalue:null
              },
              {
                dkey:'关联企业',
                dvalue:true
              },
              {
                dkey:'未关联企业',
                dvalue:false
              },
            ],
            caeName:null,
            purSearch: "",              // 搜索  
            http:'',
            cloudhttp:'',
            dataAva:false,
            loading:false,
            current:1,                  // 当前页数
            size:100,                    // 显示条数
            total:0,                    // 总条数
            delmul:[],                  // 批量删除
            show:false,
            purTitle: "",               // 标题
            multipleSelection: [],      // 选中数据
            drawer: false,              // 查看抽屉
            logo: "",                   // 抽屉企业名称
            activeName: "first",        // 查看tab栏
            tableData: [],              // 表格数据
            // 查看...动态
            timeData:[], 
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
            // 编辑
            editpaymentShow:false,
            editData:{
              provinceCode:null,
              cityCode:null,
              areaCode:null,
            },
            editpaymentrules: {
                provinceCode: [
                  { required: true, message: '请输入省级名称', trigger: 'blur' }
                ],
                cityCode: [
                  { required: true, message: '请输入市级名称', trigger: 'blur' }
                ],
                areaCode: [
                  { required: true, message: '请输入区/县级名称', trigger: 'blur' }
                ],
                // etcRelPhone: [
                //   { required: true, message: '请输入数据分配（电话）', trigger: 'blur' }
                // ],
                coalRelName: [
                  { required: true, message: '请输入联系人姓名', trigger: 'blur' }
                ],
                coalRelPhone: [
                  { required: true, message: '请输入联系人电话', trigger: 'blur' }
                ],
                businessName: [
                  { required: true, message: '请输入煤种经纪人姓名', trigger: 'blur' }
                ],
                businessPhone: [
                  { required: true, message: '请输入煤种经纪人电话', trigger: 'blur' }
                ],
                areaBuserName: [
                  { required: true, message: '请输入地区经纪人姓名', trigger: 'blur' }
                ],
                areaBuserPhone: [
                  { required: true, message: '请输入地区经纪人电话', trigger: 'blur' }
                ],
                enterpriseName: [
                  { required: true, message: '请输入企业名称', trigger: 'blur' }
                ],
                coalCode: [
                  { required: true, message: '请输入煤种编号', trigger: 'blur' }
                ],
                warseCode: [
                  { required: true, message: '请输入云仓编号', trigger: 'blur' }
                ],
                threeProId: [
                  { required: true, message: '请选择煤炭品种', trigger: 'blur' }
                ],
                bigProId: [
                  { required: true, message: '请选择煤炭分类', trigger: 'blur' }
                ],
                availableTonnage: [
                  { required: true, message: '请输入可供应吨数（吨）', trigger: 'blur' }
                ],
                exFactoryPrice: [
                  { required: true, message: '请输入出厂现汇含税价格', trigger: 'blur' }, 
                ],
                bleTonnage: [
                  { required: true, message: '请输入现有库存吨数（吨）', trigger: 'blur' }
                ],
                prags: [
                  { required: true, message: '请关联标签', trigger: 'blur' }
                ],
            },
            Extended:[
              {
                dkey:'全部云仓',
                dvalue:null
              },
              {
                dkey:'待审核',
                dvalue:0
              },
              {
                dkey:'已上架',
                dvalue:1
              },
              {
                dkey:'已拒绝',
                dvalue:2
              },
              {
                dkey:'已下架',
                dvalue:3
              },

            ],            // 查询类型
            carriageName:null,

            // 关联经纪人
            baymentShow:false,
            associationData:{},
            bpaymentrules: {
                availableTonnage: [
                { required: true, message: '请输入可供应吨数（吨）', trigger: 'blur' }
                ],
                exFactoryPrice: [
                { required: true, message: '请输入出厂现汇含税价格', trigger: 'blur' }, 
                ],
                bleTonnage: [
                { required: true, message: '请输入现有库存吨数（吨）', trigger: 'blur' }
                ],
            },
            // 经纪人列表
            bended:[], 
            /* 
            编辑审核煤种
            */
           
            seeData:{},                 // 查看数据
            
            // 身份信息
            sel: 1,   // 选项卡默认显示tab
            tabs:[
              {id:1,name:'洗煤厂',wareldenShow:false},
              {id:2,name:'贸易商',wareldenShow:false},
              {id:3,name:'焦化厂',wareldenShow:false},
              {id:4,name:'经纪人',wareldenShow:false},
              {id:5,name:'煤矿',wareldenShow:false},
              {id:6,name:'钢厂',wareldenShow:false},
              {id:7,name:'金融企业',wareldenShow:false},
              {id:8,name:'平台经纪人',wareldenShow:true},
              {id:9,name:'其他',wareldenShow:true},
            ],
            // 煤种区别
            selone: 1,   // 选项卡默认显示tab
            tabones:[
              {id:1,name:'单一煤'},
              {id:2,name:'混洗煤'},
            ],
            columnsdata: [], // 煤炭分类二级大类
            bolumlns:[],    // 煤炭品种三级分类
            // 检测报告
            jcogImageUrl: '',
            imgjcogVisible: false,
            jcimgList:[],                // 检测报告文件列表
            listarr:[],
            // 煤棚照片
            mpogImageUrl: '',
            imgmpogVisible: false,
            mpimgList:[],               // 煤棚照片文件列表
            listtwo:[],
            mzloading:false,            // 全局遮罩
            checkRemark: '', // 拒绝原因
            /*
            云仓号分配
            */
            YuncangShow:false,
            mnsdata: [],        // 云仓号
            YuncangData:{},     // 分配云仓号表单
            Yuncangrules: {
              warehouseId: [
                  { required: true, message: '请选择云仓号', trigger: 'blur' }
              ],
            },
            
            companydata:[], // 关联企业列表
            prooptions:[], // 关联标签
            useIdnum:null,
            verifyOption:[
            {
                dkey:'全部企业',
                dvalue:null
              },
              {
                dkey:'已核实企业',
                dvalue:true
              },
              {
                dkey:'未核实企业',
                dvalue:false
              },
            ],
            verifyStatus:null,
            verifyList:[],   // 核实信息数组
            brokerList:[],  // 经纪人列表
            broker:'',       // 经纪人
            exPriceTip:'',     // 上一次出厂价提示
            plPriceTip:'',     // 上一次车板价提示
        };
    },
    created() {
      this.header.Authorization = sessionStorage.getItem('Authorization');
      this.api.authAll.list({ware_identity:8})
      .then(res=>{
        if(res.data.code == 200){
          this.brokerall = res.data.data;
        }
      })
      // 行政区划
      this.api.areaInfo.all({parentCode:0})
        .then(res=>{
            if(res.data.code == 200){
                this.provinceOptions = res.data.data;
            }
      })
      // 煤炭区别 煤炭二级分类
      this.api.Prodect.alllist({level:2})
      .then(res=>{
        if(res.data.code == 200){
          this.classifydata = res.data.data;
        }
      }) 
      // 获取经纪人列表
      this.getBrokerList()
    },
    mounted(){
      this.purTitle = this.$route.meta.title;
      this.logo = sessionStorage.getItem('companyName')     // 企业名称
      this.http  = this.$store.state.http;                  // IP地址
      this.cloudhttp = this.$store.state.cloudhttp;                  // IP地址
      this.getpurcon();
    },
    methods: {
        // 煤种访问人数
        acesPurs(data){
          this.mzid = data.id;
          this.targetUserId = data.createUserId;
          this.mzdata()
        },
        mzdata(){
          this.api.busprolog.selVisitedPop({
            proId:this.mzid,
            operateObj:4,
	          operateType:7,
            searchStr:'',
            targetUserId:this.targetUserId,
            page: {
              current: this.mzcurrent,
              size: this.mzsize
            },
          }).then(res=>{
            if(res.data.code == 200){
              this.ascgridData = res.data.data.records;
              this.mztotal = res.data.data.total;
              this.acsdialsible = true;
            }
          })
        },
        haacsClose(){
          this.ascgridData = [];
          this.acsdialsible = false;
        },
        // 弹框选择省
        editerprise(){
            if(this.editData.provinceCode == ''){
              this.editData.provinceCode = null;
            }else{
              let resultArr = this.provinceOptions.filter((item)=>{
                  return item.code == this.editData.provinceCode;
              });
              this.editData.proName = resultArr[0].name;
            }
            this.editData.cityCode = null;                 // 市
            this.editData.areaCode = null;                 // 县 区
            this.cityOptions = [];
            this.countyOptions = [];

            this.api.areaInfo.all({parentCode:this.editData.provinceCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 弹框选择市
        editcity(){            
            if(this.editData.cityCode == ''){
              this.editData.cityCode = null
            }else{
              let resultArr = this.cityOptions.filter((item)=>{
                  return item.code == this.editData.cityCode;
              });
              this.editData.cityName = resultArr[0].name;
            }
            this.editData.areaCode = null;                 // 县 区
            this.editData.cityOptions = [];
            this.countyOptions = [];               // 乡
            this.api.areaInfo.all({parentCode:this.editData.cityCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.countyOptions = res.data.data;
                }
            })
        },
        // 弹框选择县 区
        editcounty(){
          if(this.editData.areaCode != ''){
            let resultArr = this.countyOptions.filter((item)=>{
                return item.code == this.editData.areaCode;
            });
            this.editData.areaName = resultArr[0].name;
          }else{
            this.editData.areaCode = null;
          }
        },


        // 更新时间搜索
        renewTime(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 更新煤种价格时间搜索
        moenyTime(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 上架日期搜索
        Times() {
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 是否是平台认证
        attestationstics(e){
          // this.attestation = e;
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 煤炭区别二级选中
        changeClassify(){
          this.breedmlns = [];
          this.breed = null;
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
          this.boredlist(this.classify);
        },
        // 煤炭三级分类选中
        changeBreed(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 煤炭区别 煤炭品种三级分类
        boredlist(ids){
          this.api.Prodect.alllist({parent_id: ids})
          .then(res=>{
            if(res.data.code == 200){
              this.breedmlns = res.data.data;
            }
          })
        },
        // 选择省
        enterprises(){
            if(this.provinceCode == ''){
              this.provinceCode = null;
            }
            this.cityCode = null;                 // 市
            this.areaCode = null;                 // 县 区
            this.cityOptions = [];
            this.countyOptions = [];
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
            this.api.areaInfo.all({parentCode:this.provinceCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 选择市
        citys(){
            this.areaCode = null;                 // 县 区
            this.cityOptions = [];
            this.countyOptions = [];               // 乡
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
            this.api.areaInfo.all({parentCode:this.cityCode})
            .then(res=>{
                if(res.data.code == 200){

                    this.countyOptions = res.data.data;
                }
            })
        },
        // 选择县 区
        countys(){
            this.tableData = [];
            this.current = 1;
            this.getpurcon();
            // this.api.areaInfo.all({parentCode:this.areaCode})
            // .then(res=>{
            //     if(res.data.code == 200){
            //         this.villageOptions = res.data.data;
            //     }
            // })
        },
      // 经纪人
        brokertics(e){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        // 公司性质
        naturstics(e){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        /*
        编辑煤种信息
        */
        // 身份信息
        selectss(item){
          this.sel = item.id;
        },
        // 煤种区别
        seleones(item){
          this.selone = item.id;
        },
        // 煤炭二级分类
        changeHandler(){
          this.$nextTick(()=>{
            this.editData.bigProdata = this.$refs.bigref.selectedLabel;
            this.prolists(this.editData.bigProId);
          })
        },
        // 煤炭品种三级分类
        prolists(ids){
          this.api.Prodect.alllist({parent_id: ids})
          .then(res=>{
            if(res.data.code == 200){
              this.bolumlns = res.data.data;
            }
          })
        },
        changepro(){
          this.$nextTick(()=>{
            this.editData.Prodata = this.$refs.proRef.selectedLabel;
          })
        },
        // 删除检测报告
        jcRemove(file) {
          if(file.response.code == 200){
            this.api.addfile.del({id:file.response.data[0].id})
            .then(res => {
              if (res.data.code == 200) {
                this.listarr.map((item,key)=>{
                  if(item == file.response.data[0].id){
                    this.listarr.splice(key, 1);
                  }
                })
                this.$message({
                  type: 'success',
                  message: '删除检测报告成功!'
                });
              }
            })
          }
        },
        // 预览检测报告照片
        jcPictureCardPrediv(file) {
          this.jcogImageUrl = file.url;
          this.imgjcogVisible = true;
        },
        // 上传检测报告照片
        addjclist(res){
          if(res.code == 200){
            this.listarr.push(res.data[0].id);
          }
        },

        // 删除煤棚照片
        mpRemove(file,fileList) {
          // console.log(file)
          // console.log(fileList)
          // console.log(this.listtwo)
          if(file.response.code == 200){
            this.api.addfile.del({id:file.response.data[0].id})
            .then(res => {
              if (res.data.code == 200) {
                this.listtwo.map((item,key)=>{
                  if(item.id == file.response.data[0].id){ // listtwo 的元素是 {id, compressedId} 对象:原先拿对象跟 id 比永远不相等,照片删了 id 还留着
                    this.listtwo.splice(key, 1);
                  }
                })
                this.$message({
                  type: 'success',
                  message: '删除煤棚照片成功!'
                });
              }
            })
          }
        },
        // 预览煤棚照片
        mpPictureCardPrediv(file) {
          this.mpogImageUrl = file.url;
          this.imgmpogVisible = true;
        },
        // 上传煤棚照片
        addmplist(res){
          if(res.code == 200){
            this.listtwo.push({id:res.data[0].id,compressedId:res.data[0].compressedId});
          }
        },

        // 编辑终端信息打开
        EditPurs(data){
          // console.log(data,'==========data===============');
            this.jcimgList = [];                // 检测报告文件列表
            this.listarr = [];
            this.mpimgList = [];               // 煤棚照片文件列表
            this.listtwo = [];
          	// 煤种单条数据
						this.api.busrodect.one({id:data.id,checkType: 2})
            .then(res=>{
              if(res.data.code == 200){
                this.editData = res.data.data;
                this.sel = this.editData.wareIdentity;
                if(this.editData.status == 1){
                  this.exPriceTip = '上一次出厂价为：'+ this.editData.preExFactoryPrice
                } else {
                  this.exPriceTip = '上一次出厂价为：'+ this.editData.lastCheckFacPrice
                }
                    // console.log(this.editData,'==================');
                if(!this.editData.provinceCode){
                  if(this.editData.etcInfo){
                    this.editData.provinceCode = this.editData.etcInfo.province_code
                    this.editData.cityCode = this.editData.etcInfo.city_code
                    this.editData.areaCode = this.editData.etcInfo.area_code
                    this.editData.proName = this.editData.etcInfo.pro_name
                    this.editData.cityName = this.editData.etcInfo.city_name
                    this.editData.areaName = this.editData.etcInfo.area_name
                  }
                }
                if(this.editData.cityCode != null){
                    // 市
                    this.api.areaInfo.all({parentCode:this.editData.provinceCode})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.cityOptions = res.data.data;
                        }
                    })
                    // 县 区
                    this.api.areaInfo.all({parentCode:this.editData.cityCode})
                    .then(res=>{
                        if(res.data.code == 200){
                            this.countyOptions = res.data.data;
                        }
                    })
                }
                // 标签列表
                this.api.bustag.list()
                .then(res=>{
                  if(res.data.code == 200){
                    this.prooptions = res.data.data;
                    this.$forceUpdate();
                    // console.log(this.editData.tags)
                    if(this.editData.tags.length>0){
                      this.editData.tags.map(item=>{
                        this.editData.prags.push(item.id);
                      })
                    }

                  }
                })
                // 获取分类列表
                this.api.Prodect.alllist({level:2})
                .then(res=>{
                  if(res.data.code == 200){
                    this.columnsdata = res.data.data;
                  }
                }) 
                this.prolists(this.editData.bigProId);
                this.editpaymentShow = true;
                // 检测报告
                if (this.editData.imgIds) {
                  let arrimgid = this.editData.imgIds.split(',');
                  this.listarr = arrimgid;
                  this.listarr.map(item => {
                    let url = this.http +
                      '/saas-common/upload/imageShow?fileName=' + item +
                      '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
                    this.jcimgList.push({
                      url: url,
                      imgIds: item,
                      response:{
                          code:200,
                          data:[{id:item}]
                      }
                    })
                  })
                };
                // 煤棚照片
                this.api.addfile.list({association_id:this.editData.associationId,f_type:4,is_compressed_file:0})
                .then(res=>{
                  console.log(res)
                  console.log(this.editData.associationId);
                if(res.data.code == 200 && res.data.data.length != 0){
                    this.listtwo = res.data.data;
                    this.listtwo.map(item => {
                      let url = this.http +
                        '/saas-common/upload/imageShow?fileName=' + item.id +
                        '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
                      this.mpimgList.push({
                        url: url,
                        imgIds: item.id,
                        response:{
                          code:200,
                          data:[{id:item.id}]
                        }
                      })
                      console.log(this.mpimgList)
                    })
                }
                })
              }
            })
            // 企业列表
            this.api.enterData.list()
            .then(res=>{
              if(res.data.code == 200){
                this.companydata = res.data.data;
              }
            })

            this.api.ware.list()
            .then(res => {
              if (res.data.code == 200) {
                this.mnsdata = res.data.data;
                
              }
            })
        },
        // 关联企业
        compaypro(){
          this.$nextTick(()=>{
            this.editData.enterpriseName = this.$refs.comref.selectedLabel;
          })
        },
        // 编辑信息取消
        editpaymentShowClose(){
            this.editpaymentShow = false;
            this.editData = {};
            this.cityOptions = [];
            this.countyOptions = [];
        },
        // 取消编辑终端信息提交
        editpaymentresetForm() {
            this.editpaymentShow = false;
            this.editData = {};
            this.cityOptions = [];
            this.countyOptions = [];
        },
        // 审核拒绝
        jRefuse(cid) {
          let obj = {
            id: cid,
            status: 2, // 状态 1-通过 2-拒绝
            checkRemark: this.checkRemark,
          }
          // console.log(cid)
          this.mzloading = true;
          this.api.busrodect.edit(obj)
            .then(res => {
              if (res.data.code == 200) {
                  this.editData = {}
                  this.tableData = [];
                  this.current = 1;
                  this.getpurcon();
                  this.editpaymentShow = false;
                  this.$message({
                    type: 'success',
                    message: '煤种发布已拒绝！'
                  });
                  this.mzloading = false;
              }
            })
        },
        /*
        云仓号分配
        */
        // 监听云仓号列表下拉
        Yuncangstics(){
          this.$forceUpdate();
          this.$nextTick(()=>{
            this.YuncangData.warehouseNb = this.$refs.ycref.selectedLabel;
          })
        }, 
        // 云仓号分配取消
        YuncangForm() {
            this.YuncangShow = false;
            this.YuncangData = {};
        },
        // 分配云仓提交
        Yuncangsubmitm(formName){
          this.$refs[formName].validate((valid) => {
              if (valid) {
                  this.mzloading = true;
                  // console.log(this.listarr)
                    this.editData.proTags = [];
                    // console.log(this.editData.prags.length)
                    if(this.editData.prags){   
                      this.editData.prags.map(item=>{
                      this.editData.proTags.push({
                        proId:this.editData.id,
                        signId:item 
                      })
                    })
                  }
                  // console.log(this.YuncangData)
                  this.editData.warehouseId = this.YuncangData.warehouseId;
                  this.editData.warehouseNb = this.YuncangData.warehouseNb;
                  if(this.editData.status == 1){
                    this.editData.isChangeUpdTime = true;
                  }else if(this.editData.status == 0){
                    this.editData.isChangeUpdTime = true;
                  }
                  this.editData.status = 1;
                  this.editData.reportList = [];
                  this.listarr.map(item => {
                    this.editData.reportList.push({
                      id: item
                    })
                  })
                  this.editData.imgIds = this.listarr.join(',');
                  this.editData.imgList = this.listtwo;
                  if(this.listtwo.length > 0){ // 没有煤棚照片的煤种(导入 / 老数据)原先在这里抛错,遮罩不消
                      this.editData.coverId = this.listtwo[0].compressedId;
                  }
                  this.editData.wareIdentity = this.sel;
                  this.editData.user = {
                    id:this.editData.createUserId,
                    wareIdentity:this.sel
                  };
                  this.editData.coalType = this.selone;
                  this.api.busrodect.edit(this.editData)
                  .then(res=>{
                    if(res.data.code == 200){
                        this.$message({
                          type: 'success',
                          message: '煤种审核已通过!'
                        });
                        this.editData = {}
                        this.tableData = [];
                        this.current = 1;
                        this.getpurcon();
                        this.YuncangShow = false;
                        this.editpaymentShow = false;
                        this.mzloading = false;
                    }
                  })
              }
          });
        },
        // 分配云仓列表数据获取
        editpaymentsubmitm(formName) {
            // 相应煤种指标必填
            if(this.editData.bigProId == 3){
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxVolatiles == '' && !this.editData.maxWaterContent == '' && !this.editData.minBondIndex == '' && !this.editData.minY == '' && !this.editData.minCsr == '' && !this.editData.lithofacies == ''){

              }else{
                // 炼焦煤：硫分、灰分、挥发、水分、粘结、Y指、热态、岩相 为必填
                return this.$message({
                    type: 'warning',
                    message: '炼焦煤指标：硫、灰分、挥发、水分、粘结、Y指、热态、岩相 为必填!'
                });
              }
            }else if(this.editData.bigProId == 29){
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxWaterContent == '' && !this.editData.minCsr == ''){

              }else{
                // 冶金焦：硫、灰色、水分、强度为必填
                return this.$message({
                    type: 'warning',
                    message: '冶金焦指标：硫、灰分、水分、强度为必填！'
                });
              }

            }else if(this.editData.bigProId == 31){
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxVolatiles == '' && !this.editData.maxWaterContent == '' && !this.editData.calorificValue == ''){

              }else{
                // 化工煤：硫、灰分、挥发、水分、热值 为必填
                return this.$message({
                    type: 'warning',
                    message: '化工煤指标：硫、灰分、挥发、水分、热值为必填！'
                });
              }
            }else if(this.editData.bigProId == 4){
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxVolatiles == '' && !this.editData.calorificValue == '' && !this.editData.minHgi == ''){

              }else{
                // 喷吹煤：硫、灰分、挥发、热值、可磨 为必填
                return this.$message({
                    type: 'warning',
                    message: '喷吹煤指标：硫、灰分、挥发、热值、可磨 为必填！'
                });
              }
            }else if(this.editData.bigProId == 5){
              console.log(!this.editData.maxStd == '')
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxWaterContent == '' && !this.editData.minCsr == ''){

              }else{
                // 化工焦：硫、灰分、水分、强度为必填
                return this.$message({
                    type: 'warning',
                    message: '化工焦指标：硫、灰分、水分、强度为必填！'
                });
              }
            }else if(this.editData.bigProId == 6){
              if(!this.editData.maxStd == '' && !this.editData.maxAshContent == '' && !this.editData.maxVolatiles == '' && !this.editData.maxWaterContent == '' && !this.editData.calorificValue == ''){

              }else{
                // 动力煤：硫、灰分、挥发、水分、热值
                return this.$message({
                    type: 'warning',
                    message: '动力煤指标：硫、灰分、挥发、水分、热值为必填！'
                });
              }
            }
            this.useIdnum = this.editData.warehouseId;   
            this.$refs[formName].validate((valid) => {
              if (valid) {
                this.YuncangShow = true;
                if(this.useIdnum){
                  this.YuncangData.warehouseId = this.useIdnum;
                  this.YuncangData.warehouseNb = this.editData.warehouseNb;
                }
              }
            });
        },
        // 重新审核煤种
        ignoreClick(data){
          data.status = 0;
          this.mzloading = true;
          this.api.busrodect.edit(data)
            .then(res=>{
              if(res.data.code == 200){
                  this.editData = {}
                  this.tableData = [];
                  this.current = 1;
                  this.getpurcon();
                  this.editpaymentShow = false;
                  this.$message({
                  type: 'success',
                  message: '煤种进入待审核状态，快去审核吧!'
                  });
                  this.mzloading = false;
              }
            })
        },
        /*
        批量配置经纪人功能 
        */
        // 导出表格
        exportlist(){
          window.open(this.http + '/saas-warehouse/busProduct/exportExcel?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99')
        },
        busrtdailog(){
          if(this.multipleSelection == ""){
              this.$message({
                  type: 'warning',
                  message: '请先选择煤种再配置经纪人!'
              });
            }else{
              this.baymentShow = true;
              this.api.authAll.list({ware_identity:8})
              .then(res=>{
                if(res.data.code == 200){
                  this.bended = res.data.data;
                }
              })
            }
        },
        // 配置经纪人取消
        ebymentShowClose(){
            this.baymentShow = false;
            this.editData = {};
            this.getpurcon();
        },
        // 取消配置经纪人提交
        bmentresetForm() {
            this.baymentShow = false;
            this.editData = {};
            this.getpurcon();
        },
        // 配置经纪人提交
        baymentsubmitm(formName) {
            this.$refs[formName].validate((valid) => {
            if (valid) {
              let arr = [];
              this.multipleSelection.map(item=>{
                let obj = Object.assign({}, {id:item.id,isChangeUpdTime:true,updateTime:item.updateTime}, this.associationData);
                arr.push(obj)
              })
              this.api.busrodect.updateBatch(arr)
              .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                    type: 'success',
                    message: '配置经纪人成功!'
                    });
                    this.associationData = {}
                    this.tableData = [];
                    this.current = 1;
                    this.getpurcon();
                    this.baymentShow = false;
                }
              })
            }
            });
        },
        bgistics(e){
          this.associationData.areaBuserPhone = e;
          this.bended.map(item=>{
            if(item.phone == e){
              this.associationData.areaBuserId = item.id;
              this.associationData.areaBuserName = item.username;
            }
          })
        },
        mzistics(e){
          this.associationData.businessPhone = e;
          this.bended.map(item=>{
            if(item.phone == e){
              this.associationData.businessId = item.id;
              this.associationData.businessName = item.username;
            }
          })
        },
        /*
        批量上架
        */
        opensjPur(){
          if(this.multipleSelection == ""){
              this.$message({
                  type: 'warning',
                  message: '请选择要上架云仓的数据!'
              });
            }else{
              let bloane = true;
              this.multipleSelection.map(item=>{
                if(item.status != 3){
                  bloane = false;
                }
              })
              if(!bloane){
                this.$message({
                      type: 'warning',
                      message: '请批量选择已下架的云仓数据!'
                });
              }else{
                let arr = [];
                this.multipleSelection.map(item=>{
                  arr.push({
                    id:item.id,
                    status:1
                  })
                })
                this.api.busrodect.updateBatch(arr)
                .then(res=>{
                    if(res.data.code == 200){
                      this.$message({
                          type: 'success',
                          message: '上架云仓成功!'
                      });
                      this.tableData = [];
                      this.current = 1;
                      this.getpurcon();
                    }
                })
              }
          }
        }, 
        // 取消下架
        shelfon(){
          this.dialoshelf = false;
          this.purshelf = '';
        },
        // 确定下架
        shelfyes(){
          let arr = [];
          this.multipleSelection.map(item=>{
            arr.push({
              id:item.id,
              status:3,
              delReason:this.purshelf
            })
          })
          this.api.busrodect.updateBatch(arr)
          .then(res=>{
              if(res.data.code == 200){
                this.$message({
                    type: 'success',
                    message: '下架云仓成功!'
                });
                this.dialoshelf = false;
                this.tableData = [];
                this.current = 1;
                this.getpurcon();
              }
          })
        },
        // 批量下架 
        openxjPur(){
          if(this.multipleSelection == ""){
              this.$message({
                  type: 'warning',
                  message: '请选择要下架云仓的数据!'
              });
          }else{
            let bloane = true;
            this.multipleSelection.map(item=>{
              if(item.status != 1){
                bloane = false;
              }
            })
            if(!bloane){
              this.$message({
                    type: 'warning',
                    message: '请批量选择已上架的云仓数据!'
              });
            }else{
              this.dialoshelf = true;
              this.purshelf = '';
            }
          }
        },
        // 查询上下架搜索功能
        prologistics(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        prtics(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        verifyChange(){
          this.tableData = [];
          this.current = 1;
          this.getpurcon();
        },
        

        // 表格数据
        getpurcon(){
            this.loading = true;
            var listByPageData;
            if(this.carriageName!=null){
              listByPageData = {
                wareIdentity:this.carrnatur,
                businessId:this.brokernatur,
                provinceCode:this.provinceCode,
                cityCode:this.cityCode,
                areaCode:this.areaCode,
                bigProId:this.classify,
                threeProId:this.breed,
                isPlatformCertification:this.attestation,
                createStartTimeDate:this.TimeArrays != null?this.TimeArrays[0]:"",
                createEndTimeDate:this.TimeArrays != null?this.TimeArrays[1]:"",
                updateStartTimeDate:this.renewTimeArray != null?this.renewTimeArray[0]:"",
                updateEndTimeDate:this.renewTimeArray != null?this.renewTimeArray[1]:"",
                facPriceStartTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[0]:"",
                facPriceEndTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[1]:"",
                searchStr:this.purSearch,
                // isUpDown:true,
                status:this.carriageName,
                isRelationEtc:this.caeName,
                isRealCheck:this.verifyStatus, // 原先这两个分支不带企业核实筛选,未核实企业会混进结果
                page:{
                    current:this.current,
                    size:this.size
                },
                effective:1,    // 删除是 0  查询是 1
                publishType:1
              }
            }else if(this.caeName != null ){
              // 传参
              var listByPageData = {
                wareIdentity:this.carrnatur,
                businessId:this.brokernatur,
                provinceCode:this.provinceCode,
                cityCode:this.cityCode,
                areaCode:this.areaCode,
                bigProId:this.classify,
                threeProId:this.breed,
                isPlatformCertification:this.attestation,
                createStartTimeDate:this.TimeArrays != null?this.TimeArrays[0]:"",
                createEndTimeDate:this.TimeArrays != null?this.TimeArrays[1]:"",
                updateStartTimeDate:this.renewTimeArray != null?this.renewTimeArray[0]:"",
                updateEndTimeDate:this.renewTimeArray != null?this.renewTimeArray[1]:"",
                facPriceStartTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[0]:"",
                facPriceEndTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[1]:"",
                searchStr:this.purSearch,
                status:this.carriageName,
                isRelationEtc:this.caeName,
                isRealCheck:this.verifyStatus, // 原先这两个分支不带企业核实筛选,未核实企业会混进结果
                page:{
                    current:this.current,
                    size:this.size
                },
                effective:1,    // 删除是 0  查询是 1
                publishType:1
              }
            }else if(this.verifyStatus!=null){
              var listByPageData = {
                wareIdentity:this.carrnatur,
                businessId:this.brokernatur,
                provinceCode:this.provinceCode,
                cityCode:this.cityCode,
                areaCode:this.areaCode,
                bigProId:this.classify,
                threeProId:this.breed,
                isPlatformCertification:this.attestation,
                createStartTimeDate:this.TimeArrays != null?this.TimeArrays[0]:"",
                createEndTimeDate:this.TimeArrays != null?this.TimeArrays[1]:"",
                updateStartTimeDate:this.renewTimeArray != null?this.renewTimeArray[0]:"",
                updateEndTimeDate:this.renewTimeArray != null?this.renewTimeArray[1]:"",
                facPriceStartTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[0]:"",
                facPriceEndTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[1]:"",
                searchStr:this.purSearch,
                status:this.carriageName,
                isRelationEtc:this.caeName,
                isRealCheck:this.verifyStatus,
                page:{
                    current:this.current,
                    size:this.size
                },
                effective:1,    // 删除是 0  查询是 1
                publishType:1
              }
            }else {
              var listByPageData = {
                wareIdentity:this.carrnatur,
                businessId:this.brokernatur,
                provinceCode:this.provinceCode,
                cityCode:this.cityCode,
                areaCode:this.areaCode,
                bigProId:this.classify,
                threeProId:this.breed,
                isPlatformCertification:this.attestation,
                createStartTimeDate:this.TimeArrays != null?this.TimeArrays[0]:"",
                createEndTimeDate:this.TimeArrays != null?this.TimeArrays[1]:"",
                updateStartTimeDate:this.renewTimeArray != null?this.renewTimeArray[0]:"",
                updateEndTimeDate:this.renewTimeArray != null?this.renewTimeArray[1]:"",
                facPriceStartTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[0]:"",
                facPriceEndTimeDate:this.moenyTimeArray != null?this.moenyTimeArray[1]:"",
                searchStr:this.purSearch,
                status:this.carriageName,
                isRelationEtc:this.caeName,
                isRealCheck:this.verifyStatus,
                page:{
                    current:this.current,
                    size:this.size
                },
                effective:1,    // 删除是 0  查询是 1
                publishType:1
              }
            }
            // 渲染表格
            this.api.busrodect.all(listByPageData)
            .then(res=>{
              if(res.data.code == 200){
                  this.loading = false; 
                  if(res.data.data.records != ""){
                      // console.log(res.data.data)
                  this.tableData = res.data.data.records;
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
            this.fourShow = false; // 不清的话:导入成功后关闭再打开,结果面板会和步骤一同屏
            clearInterval(this.importTimer); // 关掉弹窗就不再继续导入
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
                  clearInterval(this.importTimer); // 关闭后 5 秒内重开再导入会叠两个定时器
                  this.importTimer = setInterval(()=>{
                      this.perstep += 3;
                      if(this.perstep == 99){
                          this.perstep = 99;
                          clearInterval(this.importTimer); // 先停定时器再提交:提交抛错(例如已离开页面)时定时器会一直跑
                          this.$refs.upload.submit();
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
            this.current = 1;
            this.getpurcon();
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
                  console.log(this.delmul,'====');
                  let newArr = []
                  this.delmul.forEach(item => {
                    let obj = {}
                    obj.id = item
                    obj.effective = 0
                    newArr.push(obj)
                  })
                  this.api.busrodect.updateBatch(newArr).then(res=>{
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
        // 获取核实信息
			getVerifyInfo(data) {
        this.drawer = true
				this.api.verifyInfo.all({
					searchStr: '',
					proId:data.id
				}).then(res => {
					if (res.data.code == 200) {
						this.verifyList = res.data.data
					}
				})
			},
      // 获取经纪人列表
      getBrokerList(){
        this.api.authAll.list({ware_identity:8}).then(res=>{
            if(res.data.code == 200){
             this.brokerList = res.data.data
            }
        })
      },
      // 选择煤种经纪人
      selectBroker(e){
        this.editData.businessName = this.brokerList[e].username
        this.editData.businessPhone = this.brokerList[e].phone
        this.editData.businessId = this.brokerList[e].id
      },
      // 选择地区经纪人
      selectAreaBroker(e){
        this.editData.areaBuserName = this.brokerList[e].username
        this.editData.areaBuserPhone = this.brokerList[e].phone
        this.editData.areaBuserId = this.brokerList[e].id
      }
    },
};
</script>
<style lang="less" scoped>
.mzpages{
  width: 100%;
  padding: 20px 20px;
  box-sizing: border-box;
  .el-pagination{
    float: right;
  }
}
.butcls{
  width: 300px;
  height: 88px;
  padding: 26px 0;
  margin: auto;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
}
// 身份选择样式开始
.labfl{
			display: flex;
			flex-wrap: wrap;
			.totalTab{
				padding: 6px;
				font-size: 14px;
				line-height: 20px;
				color: #333;
				display: inline-block;
				margin:0 8px 8px 0;
				background: #F9F9F9;
				border-radius:4px;
				text-align: center;
				border: 1px solid #F9F9F9;
				box-sizing: border-box;
			}
			.active{
				color: #E2231A;
				border: 1px solid #E2231A;
				background:#FFF3F2;
				box-sizing: border-box;
			}
		
		}
// 身份选择样式结束
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
  height: 200px;
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
  .seaclass{
    font-size: 14px;
    line-height: 28px;
    color: #333;
    margin:6px 0;
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
  height: calc(100vh - 292px);
  background-color: #fff;
  margin: 8px;
  box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
  padding: 8px;
  .pur-table{
    width:100%;
    height: 94%;

    // /deep/ .el-table__body-wrapper::-webkit-scrollbar-thumb {
    //     border-radius: 20px;
    //     width: 200px;
    // }
    /deep/ .el-table__body-wrapper::-webkit-scrollbar {
        width: 7px;
        height: 14px;
    }
    /deep/ .el-table__fixed{
      height: calc(100% - 14px) !important;
    }
    // /deep/ .el-table__body-wrapper{
    //    z-index: 2
    // }

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

// 查看核实信息抽屉
.draBox{
  height: 100%;
  padding: 40px 32px 0px;

  .draTitle{
    margin-bottom: 36px;
    font-size: 18px;
    color: #333;
    font-weight: 700;
  }

  .draContent{

  }

  .noData{
      width: 100%;
      height: calc(100% - 100px);
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      box-sizing: border-box;
  }
}

.contentItem{

  .timeNode{
    height: 14px;

    .itemIcon{
      width: 14px;
      height: 14px;
      float: left;
      margin-right: 10px;
      background: linear-gradient(90deg, #E2231A 0%, #CF2C2C 100%);
      border-radius: 50%;
      text-align: center;
      line-height: 14px;

      img{
        width: 10px;
        height: auto;
      }
    }
    .infoTime{
      float: left;
      font-size: 12px;
      color: #666;
    }
  }

  .itemInfo{
    padding-left: 6px;
    padding-bottom: 18px;

    .txtBox{
      height: auto;
      border-left: 2px solid #D9D9D9;
      padding: 6px 0 0 18px;

      .adressTxt{
        margin-bottom: 14px;

        img{
          width: 12px;
          height: 12px;
          margin-left: 4px;
        }
      }

      .comment{
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