<template>
  <!-- 采购合同页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span>
      </el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <!-- 协同人模块 V.0.0.2 上线  -->
        <!-- <el-button size="small">删除协同人</el-button>
        <el-button size="small">添加协同人</el-button> -->
        <!-- <el-button size="small">导出</el-button> -->
        <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
        <el-button size="small" @click="addPurs" v-if="$route.meta.add">新建</el-button> 
        <el-button size="small" style="margin-left: 24px" @click="draftsAll">草稿箱</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch">
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
        <!-- 签署状态 -->
        <el-select v-model="signValue" clearable placeholder="请选择签署状态" size="small" @clear="signClear" @change="signchange">
          <el-option
            v-for="item in signStatusoptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
        <!-- 支付状态 -->
        <el-select v-model="PaymentValue" clearable placeholder="请选择支付状态" size="small" @clear="PaymentClear" @change="Paymentchange"> 
          <el-option
            v-for="item in PaymentStatusoptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
        <!-- 合同状态 -->
        <!-- <el-select v-model="ContractValue" clearable placeholder="请选择合同状态" size="small" @clear="ContractClear"> 
          <el-option
            v-for="item in ContractStatusoptions"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select> -->
      </el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav"
      v-loading="loading"
      element-loading-text = "加载中...">
      <div class="pur-table">
        <!-- 表格 -->    
        <el-table
        style="width: 100%"
        height="100%"
        border
        :data="tableData"
        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
        @selection-change="purChange"> 
          <el-table-column type="selection" fixed="left"></el-table-column>
          <el-table-column prop="id" v-if="show">
          </el-table-column>  
          <el-table-column
            label="操作" min-width="100" fixed="left">
            <template slot-scope="scope">
              <el-button @click="Clicksee(scope.row)" size="mini" style="margin:0 0 5% 10px;">查看</el-button>
              <!-- && !scope.row.isPayOnline -->
              <el-button size="mini" @click="EditPurs(scope.row)"  v-if="$route.meta.edit && scope.row.createUserId == userIds && !scope.row.isEsign && !scope.row.isPayOnline">编辑</el-button>
            </template>
          </el-table-column>
          <el-table-column label="合同编号" sortable :show-overflow-tooltip="true" width="160">
            <template slot-scope="scope">
              <div>
                <span v-if="scope.row.contractCode"
                  style="color: #cd1724;cursor: pointer;"
                  @click="purSee(scope.row)"
                  >{{ scope.row.contractCode }}</span
                >
                <span v-else>---</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column v-show="AdjustForm" label="在线支付" width="100" fixed="left">
            <template slot-scope="scope">
              <el-button size="mini" @click="gopayment(scope.row)">付 款</el-button>
            </template>
          </el-table-column>
          <el-table-column
            label="电子签章" min-width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.isEsign">开启</span>
              <span v-else>关闭</span>
            </template>
          </el-table-column>
          <el-table-column label="签署合同预览" :show-overflow-tooltip="true" width="120">
            <template slot-scope="scope">
              <!-- <span v-if="scope.row.flowStatus == -1">待确认的</span>
              <span v-else @click="confselSign(scope.row)" style="color:#CF1724;cursor: pointer;">签署预览</span> -->
              <span v-if="!scope.row.isEsign"> -- </span>
              <span v-else @click="confselSign(scope.row)" style="color:#CF1724;cursor: pointer;">签署预览</span>
            </template>
          </el-table-column>
          <el-table-column
            label="签署状态" min-width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.isEsignComplete">已完成</span>
              <span v-else>未完成</span>
            </template>
          </el-table-column>
          <el-table-column label="签署方信息" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="!scope.row.isEsign"> -- </span>
              <span v-else @click="seeone(scope.row)" style="color:#CF1724;cursor: pointer;">查看</span>
            </template>
          </el-table-column>
          <el-table-column label="系统编号" prop="systemCode" sortable width="260" :show-overflow-tooltip="true">
            <!-- <template slot-scope="scope">
              <span style="color: #cd1724;cursor: pointer;" @click="purSee(scope.row)">{{ scope.row.systemCode }}</span>
            </template> -->
          </el-table-column>
          <el-table-column prop="contractName" label="合同名称" :show-overflow-tooltip="true" width="140">
            <!-- <template slot-scope="scope">  
              <span
                style="color: #cd1724;cursor: pointer;"
                @click="purSee(scope.$index, scope.row)"
                >{{ scope.row.contractName}}</span>
            </template> -->
          </el-table-column>
          <!-- <el-table-column label="项目类型" :show-overflow-tooltip="true" width="120">
            <template slot-scope="scope">
              <span v-if="scope.row.contractResource === 1">供应链业务</span>
              <span v-else-if="scope.row.contractResource === 2">供应链金融</span>
            </template>
          </el-table-column> -->
          <el-table-column prop="signingTime" label="签订时间" sortable width="120" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column label="所属月份" sortable width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.month">{{scope.row.month|timeDate}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="salesUnit" label="销售单位" :show-overflow-tooltip="true" width="260">
          </el-table-column>
          <el-table-column prop="buyerUnit" label="购货单位" :show-overflow-tooltip="true" width="260">
          </el-table-column>
          <el-table-column prop="deliveryPlace" label="发货地" :show-overflow-tooltip="true" width="260">
          </el-table-column>
          <el-table-column prop="receiptPlace" label="收货地" :show-overflow-tooltip="true" width="260">
          </el-table-column>
          <el-table-column label="产品大类" :show-overflow-tooltip="true">
            <template slot-scope="scope" v-if="scope.row.proBigType">
              {{scope.row.proBigType|protype}}
            </template>
          </el-table-column>
          <el-table-column label="产品小类" :show-overflow-tooltip="true">
            <template slot-scope="scope" v-if="scope.row.proSmaType">
              {{scope.row.proSmaType|promiddletype}}
            </template>
          </el-table-column>
          <el-table-column label="产品名称" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.proSmaType == 20">{{scope.row.proRemark}}</span>
              <span v-else>---</span>
            </template>
          </el-table-column>
          <el-table-column prop="contractNum" label="合同数量" width="120" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column label="货品单价" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span>{{scope.row.unitPrice}} 元 / 吨</span>
            </template>
          </el-table-column>
          <el-table-column label="总价金额" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.sumPrice">{{scope.row.sumPrice|currency}} 元</span>
              <span v-else>---</span>
            </template>
          </el-table-column>
          <!-- <el-table-column label="标的指标" sortable width="250" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                  <div v-if="!scope.row.sulfur && !scope.row.volatilization && !scope.row.water">---</div>
                  <div v-else>
                      <div style="display: inline-block;"><span v-if="scope.row.sulfur">硫分St.d% : {{scope.row.sulfur}}</span><span v-else></span></div>
                      <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.volatilization">挥发分Vdaf% : {{scope.row.volatilization}}</span><span v-else></span></div>
                      <div style="display: inline-block;margin-left:15px;"><span v-if="scope.row.water">水分Mt% : {{scope.row.water}}</span><span v-else></span></div>
                  </div>
              </template>
          </el-table-column> -->
          <el-table-column label="价格类型" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.priceType">{{scope.row.priceType|pricetype}}</span>
            </template>
          </el-table-column>
          <el-table-column label="付款方式" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.payMethod">{{scope.row.payMethod|paymentmethod}}</span>
            </template>
          </el-table-column>
          <el-table-column label="支付方式" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.payType">{{scope.row.payType|paymenttype}}</span>
            </template>
          </el-table-column>
          <el-table-column label="运输模式" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.transportMode">{{scope.row.transportMode|transportmode}}</span>
            </template>
          </el-table-column>
          <el-table-column label="承运工具" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.carrier">{{scope.row.carrier|carrier}}</span>
            </template>
          </el-table-column>
          <el-table-column label="备注" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.remark">{{scope.row.remark}}</span>
              <span v-else>---</span>
            </template>
          </el-table-column>
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建合同试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
                <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <span>合计：合同数量 {{newcontractNum|currency}} 吨 , 总价金额 {{newsumPrice|currency}} 元</span>
        <el-pagination
        background
        :page-sizes="[10, 30, 50, 100, 200]"
        :pager-count="5"
        :page-size.sync="size"
        :current-page.sync="current"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        @current-change='getpurcon'
        @size-change='getpurcon'>
        </el-pagination>
      </div>
      <!-- 抽屉 查看功能 -->
      <el-drawer :title="logo" :visible.sync="drawer" size="700px" @close="closedrawerno">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <!-- 文件信息 -->
                <el-col :span="24">
                  <div class="see-top">文件信息</div>
                  <dl class="see-file">
                    <dd style="margin-bottom:10px;">
                      <el-col :span="5">
                        <span>合同附件：</span>
                      </el-col>
                      <el-col :span="19">
                        <p v-if="filelist.length == 0"> --- </p>
                        <p v-else v-for="(v,k) in filelist" :key="k"><span @mouseenter="changeActive($event)" @mouseleave="removeActive($event)" @click="fileGo(v,seeData)">{{v.name}}</span></p>
                        <!-- <a v-for="(v,k) in filelist" :key="k"  :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"><p v-if="v.fileType== 1">{{v.name}}</p></a> -->
                      </el-col>
                    </dd>
                    <dd v-if="seeData.createUserId == userIds">
                      <el-col :span="5">
                        <span>上传补充协议：</span>
                      </el-col>
                      <el-col :span="14"> 
                          <el-upload
                          :action="http + '/saas-common/upload/batch'"
                          :file-list="enclosurefileList"
                          :headers="header"
                          :on-remove="fileRemove"
                          :on-preview="openfiles"
                          :before-upload="beforedemoUpload"
                          :on-success="resfiles">
                              <div class="AddbuttonUpload"><span>请选择补充协议</span></div>
                          </el-upload>
                        <!-- <p v-for="(v,k) in filelist" :key="k"><span v-if="v.fileType== 2" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)" @click="fileGos(v,seeData)">{{v.name}}</span></p> -->
                      </el-col>
                    </dd>
                    <dd v-else>
                      <el-col :span="5">
                        <span>合同补充协议：</span>
                      </el-col>
                      <el-col :span="19">
                        <p v-if="filelist.length == 0"> --- </p>
                        <p v-else v-for="(v,kindex) in enclosurefileList" :key="kindex"><span @mouseenter="changeActive($event)" @mouseleave="removeActive($event)" @click="fileGo(v,seeData)">{{v.name}}</span></p>
                      </el-col>
                    </dd>
                  </dl>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">基础信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>系统编码：</label>
                      <span>{{seeData.systemCode}}</span>
                    </li>
                    <!-- <li>
                      <label>项目类型：</label>
                      <span v-if="seeData.contractResource === 1">供应链业务</span>
                      <span v-if="seeData.contractResource === 2">供应链金融</span>
                    </li> -->
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
                <el-col :span="24" style="padding-bottom: 48px;">
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
            <el-tab-pane label="展示图" name="chart">
              <el-row>
                <el-col :span="24">
                  <div class="chart-title" style="margin:24px 0;">单笔业务数据统计</div>
                  <ul class="chart-data">
                    <li>
                      <div>
                        <i class="el-icon-alihetong1" style="color: #2272EA;"></i>
                        <span>合同</span>
                      </div>
                      <div>
                        <span v-if="searsData.contractMoney"> {{searsData.contractMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.contractCount|ToFiexd}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliwuliu2" style="color: #3FBF9C;"></i>
                        <span>运单</span>
                      </div>
                      <div>
                        <span> {{searsData.loMoney}} </span>
                      </div>
                      <div>
                        <span> {{searsData.loCount|ToFiexd}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-alizijin" style="color:#F54D36;"></i>
                        <span>结算</span>
                      </div>
                      <div>
                        <span v-if="searsData.settMoney"> {{searsData.settMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.settCount|ToFiexd}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliticket-fill" style="color: #8092C0;"></i>
                        <span>票据</span>
                      </div>
                      <div>
                        <span v-if="searsData.billMoney"> {{searsData.billMoney|currency}} 元 </span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.billCount|ToFiexd}} 吨 </span>
                      </div>
                    </li>
                    <li>
                      <div>
                        <i class="el-icon-aliqianbao11" style="color:#F59D36;"></i>
                        <span>流水</span>
                      </div>
                      <div>
                        <span v-if="searsData.detMoney"> {{searsData.detMoney|currency}} 元</span>
                        <span v-else> 0.00 元</span>
                      </div>
                      <div>
                        <span> {{searsData.detCount}} </span>
                      </div>
                    </li>
                  </ul>
                  <div class="chart-title">单笔业务流程图</div>
                  <div class="chart-leftImg" @click="zstJump">
                    <img :src="zstimg" alt="">
                    <div class="htl">合同 ( {{searsData.contractNum}} ) </div>
                    <div class="js">结算 ( {{searsData.settNum}} ) </div>
                    <div class="wl">订单 ( {{searsData.orderNum}} ) </div>
                    <div class="zj">资金 ( {{searsData.detNum}} ) </div>
                    <div class="pjl">票据流 ( {{searsData.billNum}} ) </div>
                  </div>
                  <div class="chart-rightImg">
                    <img :src="slhyimg" alt="">
                  </div>

                </el-col>
              </el-row>
            </el-tab-pane>
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer>
      <!-- 导入弹窗功能 -->
      <el-dialog
        title="合同导入"
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
                    :action="http + '/saas-contract/chainCgContract/importAlone'"
                    :headers="header"
                    name='file'
                    :show-file-list="false"
                    :on-success="addFileLet"
                    :on-change="supbeforedemoUpload"
                    :limit="1"><el-button size="mini" style="color: #333;">请选择文件</el-button>
                    </el-upload>
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <a :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=HTMB.xlsx'">下载导入模板</a>
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
      <!-- 付款确认弹窗 -->
      <el-dialog
        title="付款信息"
        :visible.sync="paymentShow"
        width="35%"
        :before-close="paymentShowClose">
          <el-form style="width:95%;" :model="AddData" :rules="paymentrules" ref="AddData" size="small" label-width="100px">
            <el-form-item label="支付方式" prop="orderSubtype">
              <el-select v-model="AddData.orderSubtype" placeholder="请选择支付方式" style="width:100%;" :disabled="true">
                <el-option label="担保支付" value="GRT"></el-option>
                <el-option label="直接支付" value="DRT"></el-option>
              </el-select>
              <div class="paymentsize">
                <span style="color:#CF1724;">* </span>
                <span>担保支付采用付款方资金冻结方式。订单经办审批后，交易资金将被冻结在原付款账户，待买方客户确认收货之后，我行根据平台指令将从资金解付至收款商户账户。</span>
              </div>
            </el-form-item>
            <!-- <el-form-item label="系统订单" prop="name">
              <el-input v-model="AddData.name" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="交易流水号" prop="name">
              <el-input v-model="AddData.name" :disabled="true"></el-input>
            </el-form-item> -->
            <el-form-item label="订单金额" prop="orderAmt">
              <el-input v-model="AddData.orderAmt"></el-input>
            </el-form-item>
            <el-form-item label="收款单位" prop="payeeEtc">
              <el-input v-model="AddData.payeeEtc" :disabled="true"></el-input>
            </el-form-item>
            <!-- <el-form-item prop="code" label="验证码">
              <el-input v-model="AddData.code" placeholder="请输入验证码">
                <template slot="suffix">
                  <span v-if="!isDisabled" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:3px;cursor:pointer;"></span>
                  <span v-if="isDisabled" :class="color" v-html="codephone" style="cursor: pointer;" ></span>
                </template>
              </el-input>
            </el-form-item> -->
            <el-form-item>
              <el-button type="primary" @click="paymentsubmitm('AddData')" style="float:right;margin:5% 0 0 2%;">确定</el-button>
              <el-button @click="paymentresetForm" style="float:right;margin:5% 0 0 2%;">取消</el-button>
            </el-form-item>
          </el-form>
      </el-dialog>
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
      <!-- 草稿箱弹框 -->
      <el-dialog title="草稿箱" :visible.sync="dialogTableDrafts">
        <el-table style="cursor: pointer;" :data="DraftsDastas" border v-loading="loading" element-loading-text = "加载中..." highlight-current-row @current-change="DraftsChange">
          <el-table-column
          label="序号" type="index" width="50"></el-table-column>
          <el-table-column property="systemCode" label="系统编号" width="240"></el-table-column>
          <el-table-column property="contractCode" label="合同编号" width="160"></el-table-column>
          <el-table-column property="contractName" label="合同名称"></el-table-column>
          <el-table-column property="updateTime" label="修改时间" width="160"></el-table-column>
          <el-table-column
            label="操作" width="80">
            <template slot-scope="scope">
              <div style="color:#CF1724;cursor: pointer;" @click.stop="DraftsClicksee(scope.row)">
                <i class="el-icon-delete"></i><span> 删除</span>
              </div>
            </template>
          </el-table-column>
          <div slot="empty">
              <div class="dataAva">
                <p>暂无草稿数据！</p>
              </div>
          </div>
        </el-table> 
        <!-- 草稿箱分页 -->
        <div class="Draftspages">
          <el-pagination
          background
          :page-sizes="[10, 30, 50, 100, 200]"
          :pager-count="5"
          :page-size.sync="draftsSize"
          :current-page.sync="draftsCurrent"
          layout="total, sizes, prev, pager, next, jumper"
          :total="draftsTotal"
          @current-change='draftsPurcon'
          @size-change='draftsPurcon'>
          </el-pagination>
        </div>    
      </el-dialog>
      <!-- 签章文件预览弹窗 -->
      <el-drawer
        title="签署合同预览"
        :visible.sync="ImgcontactVisible"
        size="35%"
        :before-close="ImgconClose">
            <div class="imgBoxss" v-loading="urlloading" element-loading-text="文件加载中">
              <iframe :src="srcurl" frameborder="0" width="675px" height="99.5%"></iframe>
            </div>  
            <div class="imgbuttonsss" v-if="openoneData.isEsignComplete">
              <!-- <el-button type="primary" icon="el-icon-printer" size="mini" @click="downPrint">打印</el-button> -->
              <el-button type="primary" icon="el-icon-download" style="float:right;" size="mini" @click="openfilews">下载</el-button>
            </div>
      </el-drawer>
      <!-- 签署方信息 -->
      <el-dialog title="签署方信息" :visible.sync="dialogTable">
        <el-table :data="gridData" border>
          <el-table-column label="签署方类型" width="120" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.signRole == '0'">个人</span>
              <span v-if="scope.row.signRole == '2'">企业</span>
            </template>
          </el-table-column>
          <el-table-column label="签署方企业名称" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.signEtcName">{{scope.row.signEtcName}}</span>
              <span v-else> -- </span>
            </template>
          </el-table-column>
          <el-table-column prop="signName" label="签署方姓名" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="signPhone" label="签署方手机号" :show-overflow-tooltip="true">
          </el-table-column>
        </el-table>
      </el-dialog>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
import pdf from 'vue-pdf';
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js' // 加载中文的包
import slhy from '../../assets/slhy.svg';
import zst from '../../assets/zst.svg';
export default {
  components: {pdf},
  data() {
    return {
      AdjustForm:null,
      // 展示图
      pathIDS:[],             // 四流合同集合
      zstimg:zst,
      slhyimg:slhy,
      seeId:'',               // 查看单条信息的id 
      orderId:'',
      searsData:{},           // 展示图数据
      
      // 合计
      newcontractNum:0,
      newsumPrice:0,
      http:'',
      current:1,              // 当前页数
      size:50,                // 显示条数
      total:0,                // 总条数
      delmul:[],              // 批量删除
      show:false,
      purTitle: "",           // 模块标题
      purSearch: "",          // 搜索
      TimeArrays:[],          // 时间搜索
      // 支付状态
      PaymentStatusoptions:[
        {
         value:true,
         label:'已开启' 
        },
        {
         value:false,
         label:'未开启' 
        }
      ],
      PaymentValue:'',
      // 合同状态
      ContractStatusoptions:[
        {
         value:'0',
         label:'全部' 
        },
        {
         value:'1',
         label:'进行中' 
        },
        {
         value:'2',
         label:'已完成' 
        },
        {
         value:'3',
         label:'待确认' 
        },
      ],
      ContractValue:'',
      // 签署状态
      signStatusoptions:[
        {
         value:false,
         label:'未完成' 
        },
        {
         value:true,
         label:'已完成' 
        }, 
      ],
      signValue:'',
      tableData: [],          // 表格数据
      dataAva:false,
      loading:false,
      multipleSelection: [],  // 选中数据
      drawer: false,          // 查看抽屉
      logo: "",               // 抽屉标题
      activeName: "first",    // 查看tab栏
      seeData:{},             // 查看数据
      filelist:[],            // 文件列表 
      timeData:[
        // {
        //     createtime:'2021-01-02',
        //     operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
        // }
      ],                      // 查看 动态

      /*
        导入功能
      */
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

      // 付款确认弹窗
      paymentShow:false,
      AddData:{
        orderSubtype:'GRT',
        // code:''
      },
      paymentrules: {
          orderAmt: [
            { required: true, message: '请输入合同订单金额', trigger: 'change' },
            { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额格式不正确', trigger: 'blur'},
          ],
          //  /^(([1-9]\d*)(\.\d{1,2})?)$|(0\.0?([1-9]\d?))$/
          // code: [
            
          //   { required: true, message: '请输入验证码', trigger: 'blur' },
          //   { max: 6, message: '请输入6位验证码', trigger: 'blur' }
          // ]
      },
      // 验证码功能
      isDisabled: false,              // 控制按钮是否可以点击（false:可以点击，true:不可点击）
      color:'onPhone',
      codephone:'获取验证码',
      userIds:null,                   // 是否开启在线支付

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
      // header:{},                      // 上传请求头
      enclosurefileList: [],          // 协议文件列表
      fileObj:{},                     // 协议文件格式验证

      // 草稿箱功能
      draftsLoading:false,            // 加载动画
      dialogTableDrafts:false,        // 草稿箱弹窗
      DraftsDastas:[],                // 草稿箱数据
      draftsCurrent:1,                // 当前页数
      draftsSize:10,                  // 显示条数
      draftsTotal:0,                  // 总条数

      // 电子签章预览
      // 文件预览
      urlloading:true,
      ImgcontactVisible:false,
      srcurl:'',
      // openfileShow:false,
      openoneData:'',   // 预览单条信息
      // 签署方信息
      dialogTable:false,
      gridData:[],      // 签署方表格信息
    };
  },
  created() {
    this.api.glconfig.list()
    .then(res=>{
        if(res.data.code == 200){
            this.AdjustForm = res.data.data[0].isOpenPay;
        }
    })
  },
  mounted(){
    this.purTitle = this.$route.meta.title;
    this.logo = sessionStorage.getItem('companyName')     // 企业名称
    this.userIds = sessionStorage.getItem('userId');
    this.header = {Authorization:sessionStorage.getItem('Authorization')}
    this.http  = this.$store.state.http;                  // IP地址
    this.getpurcon();
  },
  methods: {
    /* 
    电子签章预览功能
    */
    // 查看签署方信息
    seeone(data){
      this.dialogTable = true
      this.api.eqbSingers.all({contractId:data.id})
      .then(res=>{
        if(res.data.code == 200){
          this.gridData = res.data.data;
        }
      })
    },
    // 查询签署人状态(合同预览)
    confselSign(data){
      this.openoneData = data;
      
      if(this.carriageName != 3){
          this.api.eqbflows.selSignUrl({eqbFlowId:data.eqbFlowId})
          .then(res=>{
              if(res.data.code == 200){
                this.ImgcontactVisible = true;
                this.urlloading = false;
                this.srcurl = res.data.data;
              }
          })
      }
    },
    // 关闭签署合同预览弹窗
    ImgconClose() {
        this.ImgcontactVisible = false;
        this.srcurl = '';
    },
    // 下载签署文件
    openfilews(){
      this.api.eqbflows.vDownFile({eqbFlowId:this.openoneData.eqbFlowId})
      .then(res=>{
        if(res.data.code == 200){
          window.open(res.data.data)
        }
      })
    },

    /* 
    草稿箱功能
    */
    draftsAll(){
      this.dialogTableDrafts = true;
      this.draftsPurcon();
    },
    // 草稿箱表格数据
    draftsPurcon(){
      this.draftsLoading = true;
      // 传参
      var listByPageData = {
        fileStatus:0,
        contractType:1,
        searchStr:'',
        page:{
          current:this.draftsCurrent,
          size:this.draftsSize
        }
      }
      // 渲染表格
      this.api.cgContract.all(listByPageData)
      .then(res=>{
        if(res.data.code == 200){
          this.draftsLoading = false;
          this.DraftsDastas = res.data.data.records;
          this.draftsTotal = res.data.data.total;
        }
      })
    },
    // 草稿箱选中行
    DraftsChange(val){
      this.api.cgContract.one({id:val.id})
      .then(res=>{
        if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = val.id;
          purdata.launchType = 2;
          purdata.signRole = 2;
          purdata.domains = [];
          purdata.domainsone = [];
          purdata.eqbFlow = {};
          sessionStorage.setItem('queryAll',JSON.stringify(purdata))
          this.$router.push({name:"AddPurchaseContract"})
        }
      })
    },
    // 删除草稿箱信息
    DraftsClicksee(data){
      this.api.cgContract.del([data.id])
      .then(res=>{
        if(res.data.code == 200){
          this.$message({
            type: 'success',
            message: '删除草稿成功!'
          });
          this.DraftsDastas = [];
          this.draftsPurcon();
          this.current = 1;
        }
      })
    },

    /*
    合同文件预览
    */
    fileGo(obj){
      console.log(obj)
      this.oneData = obj;
      this.ImgVisible = true;
      let faliename = this.oneData.response.data[0].fileName;
      this.geshi = faliename.substring(faliename.length - 4);
      let arrs = ['.pdf','.PDF','docx','DOCX','.doc','.DOC'];
      if(!arrs.includes(this.geshi)){
        this.itemshow = false;
        this.api.file.preFIle({file_id:this.oneData.response.data[0].id})
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
            this.loadingimgss = false;
            this.numPages = pdf.numPages;
        }).catch(() => {})
      }
    },
    // 点击文件下载
    openfile(){
        let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + this.oneData.response.data[0].fileName;
        window.open(url)
    },
    // 打印文件
    downPrint(){
      let url =  this.http + '/saas-common/upload/previewStream?fileName=' + this.oneData.response.data[0].fileName + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f';     // 正式环境用这个！
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
    // 协议文件格式验证
    beforedemoUpload(file) {
        let str = file.name;
        let index = str.lastIndexOf("\.");
        str = str.substring(index + 1,str.length)
        let types = [
            "png", 
            "jpg",
            "pdf",
            "docx",
            "doc",
            "jpeg",
            "PNG", 
            "JPG",
            "PDF",
            "DOCX",
            "DOC",
            "JPEG",
        ];
        const status = types.indexOf(str) != -1;
        this.fileObj[file.name] = true;
        if(!status) {
            this.fileObj[file.name] = false;
            this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC JPEG 格式,请重新上传!");
        }
        return status;
    },
    // 协议文件上传文件成功
    resfiles(response, file, fileList){
        if(response.code == 200){
            this.$message({
                type: 'success',
                message: '上传文件成功!'
            });
            this.enclosurefileList = fileList;
            sessionStorage.setItem('filelist',JSON.stringify(this.enclosurefileList));
        }
    },
    // 点击协议文件列表下载
    openfiles(res){
      this.fileGo(res)
    },
    // 删除附件文件
    fileRemove(file,fileList){
      this.api.addfile.del({id:file.response.data[0].id})
      .then(res=>{
          if(res.data.code == 200){
              this.$message({
                  type: 'warning',
                  message: '删除文件成功!'
              });
              this.enclosurefileList = fileList;
              sessionStorage.setItem('filelist',JSON.stringify(this.enclosurefileList))
          }
      })
    },
    // 关闭抽屉回调(保存文件)
    closedrawerno(){
      this.closedraws();
    },
    changeActive(e) {
        e.currentTarget.className = 'active';
    },
    removeActive(e) {
        e.currentTarget.className = '';
    },

    /*
    付款功能
    */
    // 付款确认按钮
    gopayment(data){
      // v-if="scope.row.isPayOnline && scope.row.isEsignComplete isEsign"
      if(data.isPayOnline){
        if((data.isEsign && data.isEsignComplete) || !data.isEsign){
          this.paymentShow = true;
          this.AddData.payeeEtc = data.salesUnit;
          this.AddData.contractId = data.id;
          this.AddData.contractCode = data.contractCode;
        }else{
          this.$message({
            message: '签署状态未完成，禁止支付！',
            type: 'warning'
          });
        }
      }else{
        this.$message({
          message: '未开启在线支付！',
          type: 'warning'
        });
      }
    },
    // 付款取消
    paymentShowClose(){
      this.paymentShow = false;
      this.AddData = {orderSubtype:'GRT'};
    },
    // 付款确认信息提交
    paymentsubmitm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          if(this.AddData.orderAmt == 0){
            this.$message({
              message: '付款金额必须大于零',
              type: 'warning'
            });
          }else{
            this.api.confirmPay.payment(this.AddData)
            .then(res=>{
              if(res.data.code == 200){
                this.paymentShow = false;
                this.AddData = {orderSubtype:'GRT'};
                window.open(res.data.data);
              }
            })
          }
        }else{
          this.AddData.orderAmt = '';
        }
      });
    },
    // 取消付款
    paymentresetForm() {
      this.paymentShow = false;
      this.AddData = {orderSubtype:'GRT'};
    },
    // 点击获取手机验证码付款
    phoneClick(){
      // if(this.loginForm.username != ''){
        this.api.code.code({account:'admin'})
        .then(res=>{
          var that = this;
          if(res.data.code == 200){
            that.isDisabled = true;
            var time = 60;
            var f = setInterval(function(){
              var newTime = time - 1;
              that.codephone = '重新获取('+ newTime +')s'
              that.color="onokPhone"
                time--
                if(time == 0){
                  that.codephone="获取验证码"
                  that.color="onPhone"
                  clearInterval(f);
                  that.isDisabled = false;
                }
            },1000)  
          }
          if(res.data.code == 999){
            that.isDisabled = false;
          }
        })
      // }else{
      //   this.$message({
      //     type: 'warning',
      //     message: "请先输入账号再获取验证码"
      //   });
      //   this.isDisabled = false;
      // }
    },

    // 采购表格数据
    getpurcon(){
      this.loading = true;
      // 传参  this.signValue = ""; this.PaymentValue = "";  
      var listByPageData = {
        fileStatus:1,
        contractType:1,
        searchStr:this.purSearch,
        startTime:this.TimeArrays != null?this.TimeArrays[0]:"",
        endTime:this.TimeArrays != null?this.TimeArrays[1]:"",
        isPayOnline:this.PaymentValue,
        isEsign:this.signValue,
        page:{
          current:this.current,
          size:this.size
        }
      }
      // 渲染表格
      this.api.cgContract.all(listByPageData)
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
      // 合计
      this.newcontractNum = 0;
      this.newsumPrice = 0;
      this.api.cgContract.getSumPrice({contractType:1})
      .then(res=>{
        if(res.data.code == 200 && res.data.data != null){
          this.newcontractNum = res.data.data.contractNum;
          this.newsumPrice = res.data.data.sumPrice;
        }
      })
    },
    // 搜索功能
    search() {
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 时间搜索
    Times() {
        this.tableData = [];
        this.getpurcon();
        this.current = 1;
    },
    // 支付搜索
    signchange(){
        this.tableData = [];
        this.getpurcon();
        this.current = 1;
    },
    // 支付清空
    PaymentClear(){
      this.PaymentValue = "";  
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },  
    // 合同搜索

    // 合同清空
    // ContractClear(){
    // },

    // 签署搜索
    Paymentchange(){
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 签署清空
    signClear(){
      this.signValue = "";
      this.tableData = [];
      this.getpurcon();
      this.current = 1;
    },
    // 表格选中数据
    purChange(val) {
      this.multipleSelection = val;
      this.delmul = [];
      for(let i=0;i<this.multipleSelection.length;i++){
        this.delmul.push(this.multipleSelection[i].id) 
      }
    },
    //新建采购合同
    addPurs(){
        this.$router.push({name:"AddPurchaseContract"})
    },
    // 编辑采购合同
    EditPurs(row){
      this.api.Transinfo.list({contract_id:row.id})
      .then(res=>{
        if(res.data.code == 200){
          if(res.data.data.length == 0){
            this.api.cgContract.one({id:row.id})
            .then(res=>{
              if(res.data.code == 200){
                let purdata = res.data.data;
                purdata.id = row.id;
                purdata.launchType = 2;
                purdata.signRole = 2;
                purdata.domains = [];
                purdata.domainsone = [];
                purdata.eqbFlow = {};
                sessionStorage.setItem('queryAll',JSON.stringify(purdata))
                this.$router.push({name:"EditPurchaseContract"})
              }
            })
          }else{
            this.$message({
              type: 'warning',
              message: '已付款合同不能进行编辑!'
            });
          }
        }
      })
    },
    // 查看采购合同
    purSee(row) {
      this.drawer = true;
      this.activeName = 'first';
      this.seeId = row.id;
      if(row.orderId == null){
        this.orderId = 0;
      }else{
        this.orderId = row.orderId;
      }
      this.allone(row.id)
    }, 
    // 查看四流合一图
    sears(id,rId){
      this.api.cgContract.sear({contractId:id,orderId:rId})
      .then(res=>{
        if(res.data.code == 200){
          this.searsData = res.data.data;
          let arr = res.data.data.contractId;
          this.pathIDS = [];
          if((typeof arr) != Array){
            this.pathIDS.push(arr)
          }else{
            for(let key in arr){
              this.pathIDS.push(key.contract_id)
            }
          }
        }
      })
      
    },
    // 查看基础信息
    allone(id){
      this.seeData = {};
      this.filelist = [];
      this.api.cgContract.one({id:id})
      .then(res=>{
        if(res.data.code == 200){
          this.seeData = res.data.data;
        }
      })
      // 文件
      this.api.file.one({relationId:id})
      .then(res=>{
            var newarr = [];
            var newobj;
            for(let i=0;i<res.data.data.length;i++){
                newobj = {
                    name:res.data.data[i].name,
                    response:{
                        code:200,
                        data:[res.data.data[i]]
                    }
                }
                newarr.push(newobj);
            }
            
            for(let i=0;i<newarr.length;i++){
                if(newarr[i].response.data[0].fileType == 1){
                    // 附件文件内容
                    this.filelist.push(newarr[i]);      
                }else if(newarr[i].response.data[0].fileType == 2){
                    // 补充协议内容
                    this.enclosurefileList.push(newarr[i]); 
                }
            }
            sessionStorage.setItem('filelist',JSON.stringify(this.enclosurefileList));
      })
    },
    // 查看展示图
    Clicksee(row){
      this.drawer = true;
      this.activeName = 'chart';
      this.seeId = row.id;
      if(row.orderId == null){
        this.orderId = 0;
      }else{
        this.orderId = row.orderId;
      }
      this.sears(row.id,this.orderId);
    },
    // 查看抽屉
    handleClick(tab) {
      // 判断tab栏选中状态，请求相应的接口
      if (tab.label == "基本信息") {
        this.seeData = {};
        this.filelist = [];
        this.enclosurefileList = [];
        this.allone(this.seeId,this.seeId);
      }
      if (tab.label == "动态") {
        // console.log(2);
      }
      if (tab.label == "展示图") {
        this.sears(this.seeId,this.orderId);
        this.closedraws();
      }
    },
    // 查看 编辑保存函数
    closedraws(){
      this.seeData = {};
      this.filelist = [];
      this.api.cgContract.one({id:this.seeId})
      .then(res=>{
        if(res.data.code == 200){
          this.seeData = res.data.data;
          // 文件
          this.api.file.one({relationId:this.seeId})
          .then(res=>{
                var newarr = [];
                var newobj;
                for(let i=0;i<res.data.data.length;i++){
                    newobj = {
                        name:res.data.data[i].name,
                        response:{
                            code:200,
                            data:[res.data.data[i]]
                        }
                    }
                    newarr.push(newobj);
                }
                for(let i=0;i<newarr.length;i++){
                    if(newarr[i].response.data[0].fileType == 1){
                        // 附件文件内容
                        this.filelist.push(newarr[i]);      
                    }else if(newarr[i].response.data[0].fileType == 2){
                        // 补充协议内容
                        this.enclosurefileList.push(newarr[i]); 
                    }
                }
                sessionStorage.setItem('filelist',JSON.stringify(this.enclosurefileList));
                let arr = [];
                let Array = null;
                let cid = null;
                let contract = this.filelist.length != 0?(this.filelist[0].response.data[0].id).toString():'';
                if(sessionStorage.getItem('filelist')){
                  Array = JSON.parse(sessionStorage.getItem('filelist'))
                }else{
                  Array = this.enclosurefileList;
                }
                cid = Array.map((itme)=>{
                    arr.push({id:itme.response.data[0].id});
                    return itme.response.data[0].id;
                }).join(',');
                this.seeData.filelist = arr;
                this.seeData.agreementIds = cid;
                this.seeData.contractIds = contract;
                this.seeData.eqbFlow = null;
                if(this.seeData.createUserId == this.userIds){
                  this.api.cgContract.edit(this.seeData)
                  .then(res=>{
                    if(res.data.code == 200){
                      this.seeData = {};
                      this.filelist = [];
                      this.enclosurefileList = [];
                      sessionStorage.removeItem('filelist')
                    }
                  })
                }
                this.filelist = [];
                this.enclosurefileList = [];
          })
        }
      });
    },
    // 删除功能
    openPur() {
        if(this.multipleSelection == ""){
          this.$message({
              type: 'warning',
              message: '请选择要删除的数据!'
          });
        }else{
          this.$confirm('此合同关联相关其他信息,删除后将不再关联,请谨慎操作 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
          }).then(() => {
            this.api.cgContract.del(this.delmul)
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
        }
    },
    // 跳转四流合一功能    @click="zstJump"
    zstJump(){
      this.$router.push({name:"Relation",params: {paths:'Relation',id:this.pathIDS}});
    },

    /* 
    导入功能方法
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
      this.getpurcon();
      this.current = 1;
      this.$refs.upload.clearFiles();
      this.fileForm.importnav = "";
    },
  },
};
</script>
<style lang="less" scoped>
// 签署文件查看样式开始
.imgBoxss{
    width: 100%;
    height:100%;
    overflow: hidden;
    overflow-y: scroll;
    img{
        width: 100%;
    }
}
.imgbuttonsss{
  width:10%;
  height: 40px;
  position: absolute;
  top: 0;
  right:55%;
  bottom:92%;
  left: 0;
  margin: auto;
  z-index: 10000;
}
// 签署文件查看样式结束
// 草稿箱样式开始
.Draftspages{
  width: 100%;
  height: 40px;
  padding-top: 10px;
  box-sizing: border-box;
  .el-pagination{
    float: right;
  }
}
// 上传补充协议样式
.AddbuttonUpload{
    width: 120px;
    height: 32px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 32px;
    text-align: center;
    float: left;
}
.AddbuttonUpload:hover{
    background: #e4e4e4;
    color: #333;
}
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
// 文件图片查看样式结束
// .AddbuttonUpload{
//     width: 100%;
//     height: 30px;
//     background: #FFFFFF;
//     border: 1px solid #E4E4E4;
//     border-radius: 4px;
//     font-size: 14px;
//     color: #666;
//     line-height: 30px;
//     text-align: left;
//     padding: 0 15px;
//     box-sizing: border-box;
// }
// .AddbuttonUpload:hover{
//     background: #e4e4e4;
//     color: #333;
// }
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
// 展示图样式开始
.chart-title{
  width: 160px;
  height: 30px;
  font-size: 18px;
  line-height: 30px;
  padding-left:18px;
  color: #fff;
  background-color:#FB8138;
  border-radius:0 15px 15px 0;
  margin: 24px 0;
}
.chart-leftImg{
  width: 50%;
  display: inline-block;
  padding-left: 20px;
  box-sizing: border-box;
  position: relative;
  cursor: pointer;
  img{
    width: 234px;
  }
  div{
    font-size:14px;
    line-height: 30px;
    color: #fff;
    text-align: center;
  }
  .htl{
    background-color:#2272EA;
    position:absolute;
    top:20%;
    left: 42%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .js{
    background-color:#F54D36;
    position: absolute;
    left:15%;
    top:61%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .wl{
    background-color:#48CAA6;
    position: absolute;
    left:69%;
    top:61%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .zj{
    background-color:#F59D36;
    position: absolute;
    left:12%;
    top:106%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
  .pjl{
    background-color:#8092C0;
    position: absolute;
    left:66%;
    top:106%;
    transform: translate(-50%, -50%);
    padding: 0 6px;
    border-radius:4px;
  }
}
.chart-rightImg{
  width: 50%;
  display: inline-block;
  padding-left:20px;
  box-sizing: border-box; 
  img{
    width:90%;
  }
}
.chart-data{
  li{
    overflow: hidden;
    width: 100%;
    padding-left:50px;
    box-sizing: border-box;
    font-weight: 600;
    div{
      float: left;
      width: 20%;
      font-size: 14px;
      color: #333;
      line-height: 30px;
      margin:0 40px 20px 0;
      i{
        font-size: 18px;
        line-height: 30px;
        vertical-align: middle;
        margin-right:6px;
      }
      span{
        font-size: 16px;
        line-height: 30px;
        color: #333;
        vertical-align: middle;
      }
    }
  }
}
// 展示图样式结束
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
    .el-select{
      width: 10%;
      float: right;
      margin-right:10px;
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
  position: relative;
  .pur-table{
    width:100%;
    height: 94%;
    // border: 1px solid #E4E4E4;
    // overflow: hidden;
    // overflow-y: scroll;
    // overflow-x: scroll;
    // .box{
    //   width:3600px;
    //   height: 100%;
    // }
  } 
  // 暂无数据样式
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
  // 分页样式
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
    p{
      text-decoration: none;
      color: #333;
      cursor: pointer;
    }
    .active{
      color: #EA222E;
    }
  }
}
// 查看抽屉样式结束
</style>