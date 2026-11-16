<template>
  <!-- 销项发票页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">{{ purTitle }}</span></el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <!-- <a :href="http + '/saas-bill/chainBill/exportExcel?type=2'"><el-button size="small">导出</el-button></a> -->
        <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
        <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
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
            <el-table-column type="selection"> </el-table-column>
            <el-table-column prop="id" v-if="show">
            </el-table-column>
            <el-table-column
              label="操作"
              width="100" fixed="left">
              <template slot-scope="scope">
                <el-button @click="Clicksee(scope.row)" size="mini" style="margin:0 0 5% 10px;">查看</el-button>
                <el-button @click="EditPurs(scope.row)" size="mini" v-if="$route.meta.edit">编辑</el-button>
              </template>
            </el-table-column>
            <el-table-column label="系统编码" sortable :show-overflow-tooltip="true" width="240">
                <template slot-scope="scope">
                    <span style="color: #cd1724;cursor: pointer;"
                        @click="purSee(scope.$index, scope.row)"
                        >{{ scope.row.systemCode }}</span>
                        <!-- <span v-else>{{scope.row.systemCode}}</span> -->
                </template>
            </el-table-column>
            <el-table-column prop="contractCode" label="销售合同" :show-overflow-tooltip="true" width="210">
              <template slot-scope="scope">
                  <span v-if="scope.row.contractCode">{{scope.row.contractCode}}</span>
                  <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="所属月份" sortable :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.month">{{scope.row.month|timeDate}}</span>
              </template>
            </el-table-column>
            <!-- <el-table-column label="销售结算" sortable width="120" :show-overflow-tooltip="true">
                <template slot-scope="scope">
                  <el-button v-if="scope.row.relationFlag == true" size="mini" type="text" @click="detailsPur(scope.row)">详情</el-button>
                  <span v-if="scope.row.relationFlag == false">---</span>
                </template>
            </el-table-column> -->
            <el-table-column prop="billCode" label="发票代码" :show-overflow-tooltip="true" width="160">
            </el-table-column>
            <el-table-column prop="billNumber" label="发票号码" :show-overflow-tooltip="true" width="160">
            </el-table-column>
            <el-table-column prop="openDate" label="开票日期" sortable :show-overflow-tooltip="true" width="120">
            </el-table-column>
            <el-table-column prop="dutyParagraph" label="税号" :show-overflow-tooltip="true" width="170">
            </el-table-column>
            <el-table-column prop="companyName" label="购方企业名称" :show-overflow-tooltip="true" width="230">
            </el-table-column>
            <el-table-column label="数量" :show-overflow-tooltip="true" width="120">
                <template slot-scope="scope">
                  <span v-if="scope.row.count">{{scope.row.count}} 吨</span>
                </template>
            </el-table-column>
            <el-table-column prop="unitPrice" label="单价" :show-overflow-tooltip="true" width="120">
                <template slot-scope="scope">
                  <span v-if="scope.row.unitPrice">{{scope.row.unitPrice}} 元 / 吨</span>
                </template>
            </el-table-column>
            <el-table-column label="金额" :show-overflow-tooltip="true" width="120">
                <template slot-scope="scope">
                  <span v-if="scope.row.baseMoney">{{scope.row.baseMoney|currency}} 元</span>
                </template>
            </el-table-column>
            <el-table-column label="税率" width="80" :show-overflow-tooltip="true">
              <template slot-scope="scope">
                <span v-if="scope.row.taxRate">{{scope.row.taxRate}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="税额" :show-overflow-tooltip="true" width="120">
                <template slot-scope="scope">
                  <span v-if="scope.row.taxAmount">{{scope.row.taxAmount|currency}} 元</span>
                </template>
            </el-table-column>
            <el-table-column label="含税金额" :show-overflow-tooltip="true" width="140">
                <template slot-scope="scope">
                  <span v-if="scope.row.includMoney">{{scope.row.includMoney|currency}} 元</span>
                </template>
            </el-table-column>
            <el-table-column label="商品名称" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.goodsName">{{scope.row.goodsName}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="单位" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.uit">{{scope.row.uit}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
            <el-table-column label="规格型号" :show-overflow-tooltip="true" width="120">
              <template slot-scope="scope">
                <span v-if="scope.row.modle">{{scope.row.modle}}</span>
                <span v-else>---</span>
              </template>
            </el-table-column>
          <!-- 暂无数据 -->
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以新建试试</p>
                  <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">新建</el-button>
                  <el-button size="small" @click="Importdailog" v-if="$route.meta.import">导入</el-button>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
        <span>合计：发票金额 {{newcontractNum|currency}} 元 , 税额 {{yuansumPrice|currency}} 元 , 含税金额 {{newsumPrice|currency}} 元</span>
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
                    <div class="see-top">销项信息</div>
                    <ul class="see-lable">
                        <li>
                            <label>系统编码：</label>
                            <span>{{seeData.systemCode}}</span>
                        </li>
                        <li>
                            <label>销售合同：</label>
                            <span v-if="seeData.contractCode">{{seeData.contractCode}}</span>
                            <span v-else>---</span>
                        </li>
                        <!-- <li>
                            <label>销售结算：</label>
                            <span v-if="seeData.relationFlag == true" @click="detailsPurAll" style="color:#cf1724; cursor: pointer;">详情</span>
                            <span v-else>---</span>
                        </li> -->
                        <li>
                            <label>发票代码：</label>
                            <span>{{seeData.billCode}}</span>
                        </li>
                        <li>
                            <label>发票号码：</label>
                            <span>{{seeData.billNumber}}</span>
                        </li>
                        <li>
                            <label>开票日期：</label>
                            <span>{{seeData.openDate}}</span>
                        </li>
                        <li>
                            <label>税号：</label>
                            <span>{{seeData.dutyParagraph}}</span>
                        </li>
                        <li>
                            <label>销方企业名称：</label>
                            <span>{{seeData.companyName}}</span>
                        </li>
                        <li>
                          <label>商品名称：</label>
                          <span v-if="seeData.goodsName">{{seeData.goodsName}}</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>单位：</label>
                          <span v-if="seeData.uit">{{seeData.uit}}</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>数量：</label>
                          <span v-if="seeData.count">{{seeData.count}} 吨</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>规格型号：</label>
                          <span v-if="seeData.modle">{{seeData.modle}}</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>单价：</label>
                          <span v-if="seeData.unitPrice">{{seeData.unitPrice|currency}} 元 / 吨</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>金额：</label>
                          <span v-if="seeData.baseMoney">{{seeData.baseMoney|currency}} 元</span>
                        </li>
                        <li>
                          <label>税率：</label>
                          <span v-if="seeData.taxRate">{{seeData.taxRate}}</span>
                          <span v-else>---</span>
                        </li>
                        <li>
                          <label>税额：</label>
                          <span v-if="seeData.taxAmount">{{seeData.taxAmount|currency}} 元</span>
                        </li>
                           <li>
                          <label>含税金额：</label>
                          <span v-if="seeData.includMoney">{{seeData.includMoney|currency}} 元</span>
                        </li>
                    </ul>
                </el-col>
                <el-col :span="24">
                  <div class="see-top">附件信息</div>
                  <dl class="see-file">
                    <dd>
                      <el-col :span="5">
                        <span>上传合同附件：</span>
                      </el-col>
                      <el-col :span="19">
                        <a v-for="(v,k) in filelist" :key="k" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)" @click="fileGo(v,seeData)"><p>{{v.name}}</p></a>
                      </el-col>
                    </dd>
                  </dl>
                </el-col>
              </el-row>
            </el-tab-pane>
            <!-- <el-tab-pane label="动态" name="second">
                <el-row>
                    <el-col :span="24">
                        <div class="see-top" style="margin-bottom:24px;">动态信息</div>
                        <el-timeline>
                            <el-timeline-item v-for="(v,i) in timeData" :key="i" :timestamp="v.createtime" placement="top">
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
    <!-- 销售结算详情一级抽屉 -->
    <el-drawer
      size="350px"
      title="销售结算"
      :visible.sync="drawerss"
      :before-close="drawerssClose"
      :append-to-body="true">
      <div class="purdrawer">
        <p>结算信息</p>
        <div class="purdrawer-box">
          <dl class="purdrawer-nav">
            <dd v-for="(v,k) in purdrawersize" :key="k">
              <el-row>
                <el-col :span="10"><span>{{v.settlementCount}} 吨</span></el-col>
                <el-col :span="10"><span v-if="v.sumPrice">{{v.sumPrice|currency}} 元</span></el-col>
                <el-col :span="4"><el-button size="mini" @click="detailsdrawer(v.id)">查看</el-button></el-col>
              </el-row>
            </dd>
          </dl>
        </div>
        <el-row style="margin-top:30px;">
          <el-col :span="9"><span>合计</span></el-col>
          <el-col :span="10"><span style="margin-left:5px;">{{totalpurdrawersize}} 元</span></el-col>
        </el-row>
      </div>
    </el-drawer>
    <!-- 销售结算详情二级抽屉 -->
    <el-drawer
      size="300px"
      :append-to-body="true"
      :visible.sync="innerDrawer"
      :with-header="false"
      :modal="false"
      class="drawerstyle"
      style="
        margin-right: 350px !important;
        box-shadow: 2px 0 10px rgba(102, 102, 102, 0.15);">
        <div class="details-title">
          <p>详情</p>
        </div>
        <dl class="details-nav">
          <dd>
            <span>系统编号 :</span>
            <p>{{innerDrawerData.systemCode}}</p>
          </dd>
          <dd>
            <span>销售合同 :</span>
            <p>{{innerDrawerData.contractName}}</p>
          </dd>
          <dd>
            <span>客户 :</span>
            <p>{{innerDrawerData.supplierName}}</p>
          </dd>
          <dd>
            <span>货物名称 :</span>
            <p>{{innerDrawerData.goodsName}}</p>
          </dd>
          <dd>
            <span>结算时间 : </span>
            <p>{{innerDrawerData.setStartTime}} ~ {{innerDrawerData.setEndTime}}</p>
          </dd>
          <dd>
            <span>结算数量 :</span>
            <p>{{innerDrawerData.settlementCount}} 吨</p>
          </dd>
          <dd>
            <span>结算单价 :</span>
            <p v-if="innerDrawerData.unitPrice">{{innerDrawerData.unitPrice|currency}} 元 / 吨</p>
          </dd>
          <dd>
            <span>结算总价 :</span>
            <p v-if="innerDrawerData.sumPrice">{{innerDrawerData.sumPrice|currency}} 元</p>
          </dd>
        </dl>
    </el-drawer>
      <!-- 导入弹窗功能 -->
      <el-dialog
        title="发票导入"
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
                    :action="http + '/saas-bill/chainBill/importExcelXx'"
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
              <a :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=XXFPTP.xlsx'">下载导入模板</a>
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
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
import pdf from 'vue-pdf';
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js'; // 加载中文的包
import slhy from '../../assets/slhy.svg';
import zst from '../../assets/zst.svg';
export default {
  components: {pdf},
  data() {
    return {
      // 展示图
      pathIDS:[],             // 四流合同集合
      zstimg:zst,
      slhyimg:slhy,
      seeId:'',               // 查看单条信息的id 
      searsData:{},           // 展示图数据

      // 合计
      newcontractNum:0,       // 金额
      yuansumPrice:0,         // 税额
      newsumPrice:0,          // 含税金额
      http:'',
      dataAva:false,
      loading:false,
      current:1,                          // 当前页数
      size:50,                            // 显示条数
      total:0,                            // 总条数
      delmul:[],                          // 批量删除
      show:false,
      purTitle: "",                       // 标题
      ImportVisible:false,                // 导入弹窗
      ImportVisible:false,                // 取消导入
      header:{},                          // 上传token
      filelist:[],            // 文件列表
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

      purSearch: "",                      // 搜索
      TimeArrays:[],           // 时间搜索
      tableData: [],                      // 表格数据
      multipleSelection: [],              // 选中数据
      drawer: false,                      // 查看抽屉
      dID:'',
      drawerss: false,                    // 销售结算详情抽屉开关
      purdrawersize:[],                   // 流水信息
      totalpurdrawersize:0,               // 合计
      innerDrawer: false,                 // 销售结算详情里面抽屉
      innerDrawerData:{},                 // 查看数据
      logo: "",                           // 抽屉标题
      activeName: "first",                // 查看tab栏
      seeData:{},                         // 查看数据
      // 查看 动态
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
    };
  },
  created() {},
  mounted(){
    this.purTitle = this.$route.meta.title;             // 标题
    this.logo = sessionStorage.getItem('companyName')   // 企业名称
    this.http  = this.$store.state.http;                 // IP地址
    this.getpurcon();
  },
  methods: {
      /*
      合同文件预览
      */
      fileGo(obj){
        console.log(obj)
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
          type:2,
          searchStr:this.purSearch,
          startTime:this.TimeArrays != null?this.TimeArrays[0]:"",
          endTime:this.TimeArrays != null?this.TimeArrays[1]:"",
          page:{
            current:this.current,
            size:this.size
          }
        }
        // 渲染表格
        this.api.xxbill.all(listByPageData)
        .then(res=>{
          this.loading = false;
          if(res.data.code == 200){
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
        this.yuansumPrice = 0;
        this.newsumPrice = 0;
        this.api.xxbill.amount({type:2})
        .then(res=>{
          if(res.data.code == 200 && res.data.data != null){
            this.newcontractNum = res.data.data.baseMoney;
            this.yuansumPrice = res.data.data.taxAmount;
            this.newsumPrice = res.data.data.sumMoney;
          }
        })
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
      // 销售结算详情关闭按钮
      drawerssClose(){
        this.innerDrawer = false;
        this.drawerss = false;
      },
      // 表格查看销售结算详情
      detailsPur(data){
        this.purdrawersize = [];
        this.totalpurdrawersize = 0;
        this.api.xxbill.oneAll({id:data.id})
        .then(res=>{
          if(res.data.code == 200){
            this.purdrawersize = res.data.data
            // 流水明细合计
            for(let i=0;i<this.purdrawersize.length;i++){
              this.totalpurdrawersize += this.purdrawersize[i].sumPrice
            }
            let a = this.totalpurdrawersize;
            this.totalpurdrawersize = a.toFixed(2) * 1;
          }
        })
        this.drawer = false;
        this.drawerss = true;
      },
      // 表格查看销售结算详情
      detailsPurAll(){
        this.purdrawersize = [];
        this.totalpurdrawersize = 0;
        this.api.xxbill.oneAll({id:this.dID})
        .then(res=>{
          if(res.data.code == 200){
            this.purdrawersize = res.data.data
            // 流水明细合计
            for(let i=0;i<this.purdrawersize.length;i++){
              this.totalpurdrawersize += this.purdrawersize[i].sumPrice
            }
            let a = this.totalpurdrawersize;
            this.totalpurdrawersize = a.toFixed(2) * 1;
          }
        })
        this.drawer = false;
        this.drawerss = true;
      },
      // 二级抽屉数据回显
      detailsdrawer(index){
        for(let i=0;i<this.purdrawersize.length;i++){
          if(this.purdrawersize[i].id == index){
            this.innerDrawerData = this.purdrawersize[i]
          }
        }
        this.innerDrawer = true;
      },
      //搜索功能
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
      // 表格选中数据
      purChange(val) {
        this.multipleSelection = val;
        this.delmul = [];
        for(let i=0;i<this.multipleSelection.length;i++){
          this.delmul.push(this.multipleSelection[i].id) 
        }
      },
      // 新建销项发票
      addPurs(){
          this.$router.push({name:"AddSalesStatistics"})
      },
      // 编辑销项发票
      EditPurs(row){
        this.api.bill.one({id:row.id})
        .then(res=>{
          if(res.data.code == 200){
            let purdata = res.data.data;
            purdata.id = row.id;
            sessionStorage.setItem('queryAll',JSON.stringify(purdata));
            this.$router.push({name:"EditSalesStatistics"})
          }
        })
      },
      // 查看基础信息函数
      allone(id){
        this.api.xxbill.one({id:id})
        .then(res=>{
          if(res.data.code == 200){
            this.seeData = res.data.data;
            this.dID = res.data.data.id;
          }
        })
        // 查看文件
        this.api.file.one({relationId:id})
        .then(res=>{
          this.filelist = res.data.data
        })
      },
      // 查看四流合一图函数
      sears(id){
        this.api.listByPage.sear({billId:id})
        .then(res=>{
          if(res.data.code == 200){
            this.searsData = res.data.data;
            let arr = res.data.data.contractIds;
            this.pathIDS = [];
            for(let key in arr){
              this.pathIDS.push(arr[key].contract_id)
            }
          }
        })
      },
      // 查看
      purSee(index, row) {
        this.drawer = true;
        this.activeName = 'first';
        this.seeId = row.id;
        this.allone(row.id)
      },
      // 查看展示图
      Clicksee(row){
        this.drawer = true;
        this.activeName = 'chart';
        this.seeId = row.id;
        this.sears(row.id)
      },
      // 跳转四流合一
      zstJump(){
        this.$router.push({name:"Relation",params: {paths:'Relation',id:this.pathIDS}});
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
          this.allone(this.seeId);
        }
        if (tab.label == "动态") {
          // console.log(2);
        }
        if (tab.label == "展示图") {
          this.sears(this.seeId);
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
                this.api.xxbill.del(this.delmul)
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
  },
};
</script>
<style lang="less" scoped>
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
// 采购结算详情样式
.purdrawer{
  width: 100%;
  padding:16px 0 16px 16px;
  box-sizing: border-box;
  border-top:1px solid #F5F5F5;
  p{
    font-size: 14px;
    color: #333;
    line-height:62px;
  }
  span{
    font-size: 14px;
    line-height: 42px;
    color: #333;
  }
  .purdrawer-box{
    width: 100%;
    max-height: 600px;
    .purdrawer-nav{
      max-height: 600px;
      overflow-x: hidden;
      overflow-y: scroll;
      border-bottom:1px solid #E4E4E4;
      padding-bottom:27px;
      margin-bottom: 17px;
      padding-right:16px;
      dd{
        height:42px;
        .el-button{
          margin-top:7px;
        }
      }
    }
  }
}
// 子采购央视详情样式
.details-title{
  background-color: #FBFBFB;
  height:75px;
  border-bottom:1px solid #F5F5F5;
  padding: 20px;
  box-sizing: border-box;
  p{
    font-size: 16px;
    color: #333;
  }
}
.details-nav{
  height: 100%;
  background-color: #FBFBFB;
  padding: 20px 0 0 20px;
  font-size: 14px;
  dd{
    span{
      color:#999;
    }
    p{
      display: block;
      min-height: 60px;
      color: #333;
      line-height: 60px;
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
    margin-bottom:265px;
    .dataAva{
      position: absolute;
      top:75%;
      left: 25%;
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
  padding: 15px 0;
  label {
    display: inline-block;
    width: 18%;
    font-size: 14px;
    color: #999;
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
      cursor: pointer;
    }
    .active{
      color: #EA222E;
    }
  }
}
// 查看抽屉样式结束
</style>