<template>
  <!-- 核销应付单页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
    </el-row>
    <!-- 头部结束 -->
    <!-- 内容开始 -->
    <div class="pur-nav">
        <div class="receivable">
            <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="70px">
                <el-row>
                    <el-col :span="1.5">
                        <span style="color:red;">* </span>
                        <label style="color:#666;font-size:14px;line-height:32px;">供应商</label>
                    </el-col>
                    <el-form-item prop="supplierName">
                        <el-col :span="7">
                            <el-input v-model="AddData.supplierName" placeholder="请选择供应商全称" disabled style="width:100%;">
                                <span slot="suffix" style="color:#333;cursor: pointer;padding-right:16px;" @click="dialogSelect">关 联</span>
                            </el-input>
                        </el-col>
                        <el-col :span="15" class="receivable-conent">
                            <!-- <img src="@/assets/qian.svg" alt="">
                            <b v-if="amountPayable != 0">{{amountPayable|currency}}</b>
                            <b v-else> 0 </b>
                            <span>元</span>
                            <span style="margin:0 24px 0 32px;">采购结算单合计</span> -->
                            <el-button type="primary" size="medium " @click="AddDatasubmitForm('AddData')">核 销</el-button>
                            <span v-if="heceNum != 0" style="margin:0 0 0 32px;">应付单核销差额 {{heceNum|currency}} 元</span>
                            <span v-else style="margin:0 0 0 32px;">应付单核销差额 0.00 元</span>
                            <el-tooltip style="margin:10px 0 0 5px;" content="当结算实际核销金额大于收款单实际核销金额时 , 差额为正数 . 否则为负数 !" placement="top">
                                <i class="el-icon-alihelp "></i>
                            </el-tooltip>
                        </el-col>
                    </el-form-item>
                </el-row>
                <!-- 供应商单位弹窗开始 -->
                <el-dialog
                    :visible.sync="dialogSelectUnit"
                    width="60%"
                    :before-close="BuyerhandleClose">
                    <template slot="title">
                        <span>选择供应商单位</span>
                    </template>
                    <!-- 供应商单位内容框 -->
                    <el-input type="textarea" v-model="dialogtextarea"></el-input>
                    <div class="Wdialog">
                        <div class="concentBox">
                            <div class="textbutton">
                                <!-- 供应商应用按钮 -->
                                <el-button size="small" @click="salesUnitApplication">应用</el-button>
                            </div>
                            <div class="text">
                                <!-- 供应商单位信息 -->
                                <div class="textAll">
                                    <div class="textlast">
                                        <el-input
                                        @keyup.enter.native="searchtwo"
                                        size="small"
                                        placeholder="请输入内容"
                                        prefix-icon="el-icon-search"
                                        v-model="purSearchtwo"
                                        style="width:200px;float:right;"
                                        >
                                        </el-input>
                                        <!-- <span @click="Companyshow">+ 新建</span> -->
                                    </div>
                                    <el-table
                                    border
                                    ref="multipleTable"
                                    :data="CompanytableData"
                                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                                    style="width: 100%;">
                                        <el-table-column width="33">
                                            <template scope="scope">
                                                <el-radio :label="scope.$index" v-model="radio" @change.native="getCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                                            </template>
                                        </el-table-column>
                                        <el-table-column prop="id" v-if="show">
                                        </el-table-column>
                                        <el-table-column prop="enterpriseName" label="单位名称" :render-header="renderHeaderMethods">
                                        </el-table-column>
                                        <el-table-column prop="createTime" label="创建时间">
                                        </el-table-column>
                                    </el-table>
                                    <!-- 分页 -->
                                    <div class="pur-page">
                                        <el-pagination
                                        v-show="paginationshowtwo"
                                        background
                                        :page-size.sync="sizetwo"
                                        :current-page.sync="currenttwo"
                                        layout="prev, pager, next"
                                        :total="totaltwo"
                                        @current-change='allsuppliertwo'
                                        @size-change='allsuppliertwo'>
                                        </el-pagination>
                                    </div>
                                </div>
                            </div>
                            <!-- 新建单位 -->
                            <el-form :model="CompanyNewData" :rules="rels" ref="CompanyNewData" label-width="80px" size="small" v-show="CompanyDatashow">
                                <el-row :gutter="10" class="CompanyData">
                                    <el-col :span="6">
                                        <el-form-item label=" 类 型 :" prop="custType" style="width:100%;">
                                            <el-select v-model="CompanyNewData.custType">
                                            <el-option label="个人" value="1"></el-option>
                                            <el-option label="公司" value="2"></el-option>
                                            </el-select>
                                        </el-form-item>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-form-item label=" 单 位 ：" prop="enterpriseName">
                                            <el-input v-model="CompanyNewData.enterpriseName" placeholder="请输入单位名称"></el-input>
                                        </el-form-item>
                                    </el-col>
                                    <el-col :span="2" style="padding-left:55px;">
                                        <el-button @click="CompanyNewShow" size="small">取消</el-button>
                                    </el-col>
                                    <el-col :span="2" style="padding-left:40px;">
                                        <el-button size="small" @click="CompanyNewsubmitForm('CompanyNewData')">保存</el-button>
                                    </el-col>
                                </el-row>
                            </el-form>
                        </div>
                    </div>
                </el-dialog>
                <!--供应商单位弹窗结束 -->
            </el-form>
        </div>
        <!-- 采购结算列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
                <b>采购结算</b>
                <div class="sjhxjeMoeny">
                    <span v-if="settlementmoeny != 0">实际核销金额合计: {{settlementmoeny|currency}} 元</span>
                    <span v-else>实际核销金额合计: 0.00 元</span>
                </div>
                <div class="Advanced">
                    <el-date-picker
                    @change="ClickOppositeNameone"
                    size="small"
                    v-model="SearchoppositeNameone"
                    type="month"
                    value-format="yyyy-MM"
                    placeholder="所属月份">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="jshxpurSearchs"
                    size="small"
                    placeholder="请输入内容"
                    prefix-icon="el-icon-search"
                    v-model="jshxpurSearch">
                    </el-input>
                </div>
            </div>
            <!-- 表格 -->
            <div class="pur-table">
                <el-table
                    v-loading="settlementloading"
                    element-loading-text = "加载中..."
                    height="100%"
                    border
                    :data="settlementtableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;border-bottom:0;overflow:visible !important;"
                    @selection-change="settlementChange"
                    :row-key="settlementkeyId">
                    <el-table-column type="selection" width="40"> </el-table-column>
                    <el-table-column prop="id" v-if="show">
                    </el-table-column>
                    <el-table-column label="系统编码" sortable :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span
                                style="color: #cd1724;cursor: pointer;"
                                @click="settlementpurSee(scope.$index, scope.row)"
                                >{{ scope.row.systemCode }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="month" label="所属月份" :show-overflow-tooltip="true" width="110">
                    </el-table-column>
                    <el-table-column prop="goodsName" label="货物名称" sortable :show-overflow-tooltip="true" width="120">
                    </el-table-column>
                    <el-table-column label="可核销金额" sortable :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                        <span v-if="scope.row.khxje != null">{{scope.row.khxje|currency}} 元</span>
                        <span v-else-if="scope.row.sumPrice && scope.row.khxje == null">{{scope.row.sumPrice|currency}} 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="实际核销金额" sortable :show-overflow-tooltip="true" >
                        <template slot-scope="scope">
                            <div v-if="scope.row.sjhxblone">
                                <el-input size="small" v-model="newSjhxje" v-focus @blur="loseFcous(scope.$index, scope.row)"></el-input>
                            </div>
                            <div v-else>
                                <span v-if="scope.row.sjhxje != null">{{scope.row.sjhxje|currency}} 元</span>
                                <span v-else-if="scope.row.khxje !=null && scope.row.sjhxje == null">{{scope.row.khxje|currency}} 元</span>
                                <span v-else-if="scope.row.sumPrice && scope.row.sjhxje == null && scope.row.khxje == null">{{scope.row.sumPrice|currency}} 元</span>
                                <i class="el-icon-edit" @click="cellclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="pur-pages">
                <div style="float:right;width:75%;">
                    <el-pagination
                    background
                    :page-sizes="[10, 30, 50, 100, 200]"
                    :page-size.sync="settlementsize"
                    :pager-count="5"
                    :current-page.sync="settlementcurrent"
                    layout="total, sizes, prev, pager, next, jumper"
                    :total="settlementtotal"
                    @current-change='settlementgetpurcon'
                    @size-change='settlementgetpurcon'>
                    </el-pagination>
                </div>
            </div>
            <!-- 抽屉 查看功能 -->
            <el-drawer :title="logo" :visible.sync="settlementdrawer" size="700px">
                <div class="pur-tiao"></div>
                <div class="pur-drawerNav">
                <el-tabs v-model="settlementactiveName" @tab-click="settlementhandleClick">
                    <el-tab-pane label="基本信息" name="first">
                    <el-row>
                        <el-col :span="24">
                            <div class="see-top">基础信息</div>
                            <ul class="see-lable">
                                <li>
                                    <label>系统编码：</label>
                                    <span>{{settlementseeData.systemCode}}</span>
                                </li>
                                <li>
                                    <label>采购合同：</label>
                                    <span>{{settlementseeData.contractName}}</span>
                                </li>
                                <li>
                                    <label>供应商：</label>
                                    <span>{{settlementseeData.supplierName}}</span>
                                </li>
                                <li>
                                    <label>货物名称：</label>
                                    <span>{{settlementseeData.goodsName}}</span>
                                </li>
                                <li>
                                    <label>结算时间：</label>
                                    <span v-if="settlementseeData.setStartTime">{{settlementseeData.setStartTime|Time}}</span>
                                    <span> ~ </span>
                                    <span v-if="settlementseeData.setEndTime">{{settlementseeData.setEndTime|Time}}</span>
                                </li>
                                <li>
                                    <label>结算数量：</label>
                                    <span>{{settlementseeData.settlementCount}} 吨</span>
                                </li>
                                <li>
                                    <label>结算单价：</label>
                                    <span>{{settlementseeData.unitPrice}} 元 / 吨</span>
                                </li>
                                <li>
                                    <label>结算总价：</label>
                                    <span v-if="settlementseeData.sumPrice">{{settlementseeData.sumPrice|currency}} 元</span>
                                </li>
                                <li>
                                    <label>付款单列表：</label>
                                    <span v-if="settlementseeData.relDetailFlag == true" @click="settlementdetailsPurAll" style="color:#cf1724; cursor: pointer;">详情</span>
                                    <span v-else>---</span>
                                </li>
                                <li>
                                <label>进项发票列表：</label>
                                <span v-if="settlementseeData.relBillFlag == true" @click="settlementdetailsPuronessAll" style="color:#cf1724; cursor: pointer;">详情</span>
                                <span v-else>---</span>
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
                                <a v-for="(v,k) in settlementfilelist" :key="k"  :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="settlementchangeActive($event)" @mouseleave="settlementremoveActive($event)"><p>{{v.name}}</p></a>
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
                                    <el-timeline-item v-for="(v,i) in settlementtimeData" :key="i" :timestamp="v.createtime" placement="top">
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
            <!-- 采购付款单详情一级抽屉 -->
            <el-drawer
                size="350px"
                title="流水明细"
                :visible.sync="settlementdrawerss"
                :before-close="settlementdrawerssClose"
                :append-to-body="true">
                <div class="purdrawer">
                <p>流水信息</p>
                <div class="purdrawer-box">
                <dl class="purdrawer-nav">
                    <dd v-for="(v,k) in settlementpurdrawersize" :key="k">
                    <el-row>
                        <el-col :span="10"><span>{{v.transTime}}</span></el-col>
                        <el-col :span="10"><span v-if="v.money">{{v.money|currency}} 元</span></el-col>
                        <el-col :span="4"><el-button size="mini" @click="settlementdetailsdrawer(v.id)">查看</el-button></el-col>
                    </el-row>
                    </dd>
                </dl>
                </div>
                <el-row style="margin-top:30px;">
                    <el-col :span="9"><span>合计</span></el-col>
                    <el-col :span="10"><span style="margin-left:5px;" v-if="settlementtotalpurdrawersize">{{settlementtotalpurdrawersize|currency}} 元</span></el-col>
                </el-row>
                </div>
            </el-drawer>
            <!-- 采购付款单详情二级抽屉 -->
            <el-drawer
                size="300px"
                :append-to-body="true"
                :visible.sync="settlementinnerDrawer"
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
                    <p>{{settlementinnerDrawerData.systemCode}}</p>
                    </dd>
                    <dd>
                    <span>账户名称 :</span>
                    <p>{{settlementinnerDrawerData.accountName}}</p>
                    </dd>
                    <dd>
                    <span>账户号码 :</span>
                    <p>{{settlementinnerDrawerData.account}}</p>
                    </dd>
                    <dd>
                    <span>交易日期 :</span>
                    <p>{{settlementinnerDrawerData.transTime}}</p>
                    </dd>
                    <dd>
                    <span>金额 : </span>
                    <p v-if="settlementinnerDrawerData.money">{{settlementinnerDrawerData.money|currency}} 元</p>
                    </dd>
                    <dd>
                    <span>对方户名 :</span>
                    <p>{{settlementinnerDrawerData.oppositeName}}</p>
                    </dd>
                    <dd>
                    <span>对方账号 :</span>
                    <p>{{settlementinnerDrawerData.oppositeAccount}}</p>
                    </dd>
                    <dd>
                    <span>摘要 :</span>
                    <p>{{settlementinnerDrawerData.absContent}}</p>
                    </dd>
                    <dd>
                    <span>用途 :</span>
                    <p v-if="settlementinnerDrawerData.purpose">{{settlementinnerDrawerData.purpose}}</p>
                    <p v-else>---</p>
                    </dd>
                </dl>
            </el-drawer>
            <!-- 进项发票详情一级抽屉 -->
            <el-drawer
            size="350px"
            title="票据明细"
            :visible.sync="settlementdraweroness"
            :before-close="settlementdraweronessClose"
            :append-to-body="true">
            <div class="purdrawer">
                <p>进项信息</p>
                <div class="purdrawer-box">
                <dl class="purdrawer-nav">
                    <dd v-for="(v,k) in settlementpurdrawersizeoness" :key="k">
                    <el-row>
                        <el-col :span="10"><span>{{v.openDate}}</span></el-col>
                        <el-col :span="10"><span v-if="v.includMoney">{{v.includMoney|currency}} 元</span></el-col>
                        <el-col :span="4"><el-button size="mini" @click="settlementdetailsdraweroness(v.id)">查看</el-button></el-col>
                    </el-row>
                    </dd>
                </dl>
                </div>
                <el-row style="margin-top:30px;">
                <el-col :span="9"><span>合计</span></el-col>
                <el-col :span="10"><span style="margin-left:5px;" v-if="settlementtotalpurdrawersizeoness">{{settlementtotalpurdrawersizeoness|currency}} 元</span></el-col>
                </el-row>
            </div>
            </el-drawer>
            <!-- 进项发票详情二级抽屉 -->
            <el-drawer
            size="500px"
            :append-to-body="true"
            :visible.sync="settlementinnerDraweroness"
            :with-header="false"
            :modal="false"
            class="drawerstyle"
            style="
                margin-right: 350px !important;
                box-shadow: 2px 0 10px rgba(102, 102, 102, 0.15);">
                <div class="details-title">
                <p>详情</p>
                </div>
                <dl class="details-one">
                <dd>
                    <span>系统编号 :</span>
                    <p>{{settlementinnerDraweronessData.systemCode}}</p>
                </dd>
                <dd>
                    <span>发票代码 :</span>
                    <p>{{settlementinnerDraweronessData.billCode}}</p>
                </dd>
                <dd>
                    <span>发票号码:</span>
                    <p>{{settlementinnerDraweronessData.billNumber}}</p>
                </dd>
                <dd>
                    <span>销方税号 :</span>
                    <p>{{settlementinnerDraweronessData.dutyParagraph}}</p>
                </dd>
                <dd>
                    <span>商品名称 : </span>
                    <p v-if="settlementinnerDraweronessData.goodsName">{{settlementinnerDraweronessData.goodsName}}</p>
                    <p v-else>---</p>
                </dd>
                <dd>
                    <span>规格型号 :</span>
                    <p>{{settlementinnerDraweronessData.modle}}</p>
                </dd>
                <dd>
                    <span>单价(不含税) :</span>
                    <p v-if="settlementinnerDraweronessData.unitPrice">{{settlementinnerDraweronessData.unitPrice|currency}} 元 / 吨</p>
                    <p v-else>---</p>
                </dd>
                <dd>
                    <span>税率 :</span>
                    <p>{{settlementinnerDraweronessData.taxRate}}</p>
                </dd>
                </dl>
                <dl class="details-two">
                <dd>
                    <span>采购合同 :</span>
                    <p v-if="settlementinnerDraweronessData.contractName">{{settlementinnerDraweronessData.contractName}}</p>
                    <p v-else>---</p>
                </dd>
                <dd>
                    <span>开票日期 :</span>
                    <p>{{settlementinnerDraweronessData.openDate}}</p>
                </dd>
                <dd>
                    <span>销方企业名称 :</span>
                    <p>{{settlementinnerDraweronessData.companyName}}</p>
                </dd>
                <dd>
                    <span>单位 :</span>
                    <p v-if="settlementinnerDraweronessData.uit">{{settlementinnerDraweronessData.uit}}</p>
                    <p v-else>---</p>
                </dd>
                <dd>
                    <span>数量 :</span>
                    <p v-if="settlementinnerDraweronessData.count">{{settlementinnerDraweronessData.count}} 吨</p>
                    <p v-else>---</p>
                </dd>
                <dd>
                    <span>金额 :</span>
                    <p v-if="settlementinnerDraweronessData.baseMoney">{{settlementinnerDraweronessData.baseMoney|currency}} 元</p>
                </dd>
                <dd>
                    <span>税额 :</span>
                    <p v-if="settlementinnerDraweronessData.taxAmount">{{settlementinnerDraweronessData.taxAmount|currency}} 元</p>
                </dd>
                <dd>
                    <span>含税金额 :</span>
                    <p v-if="settlementinnerDraweronessData.includMoney">{{settlementinnerDraweronessData.includMoney|currency}} 元</p>
                    <p v-else>---</p>
                </dd>
                </dl>
            </el-drawer>
        </div>
        <div class="rece-right"><span>《</span></div>
        <!-- 采购付款单列表 -->
        <div class="rece">
            <!-- 高级搜索功能 -->
            <div class="pur-title">
                <b>采购付款单</b>
                <div class="sjhxjeMoeny">
                    <span v-if="contmoeny != 0">实际核销金额合计: {{contmoeny|currency}} 元</span>
                    <span v-else>实际核销金额合计: 0.00 元</span>
                </div>
                <div class="Advanced">
                    <el-select @change="keyupBanks" v-model="banks" filterable placeholder="请选择银行名称" size="small" :clearable='true'>
                        <el-option
                        v-for="item in optionsBanks"
                        :key="item.bank_name"
                        :label="item.bank_name"
                        :value="item.bank_name">
                        </el-option>
                    </el-select>
                    <el-date-picker
                    @change="DataTime"
                    size="small"
                    v-model="keyupTime"
                    type="date"
                    value-format="yyyy-MM-dd"
                    placeholder="选择日期">
                    </el-date-picker>
                    <el-input
                    @keyup.enter.native="ClickOppositeName"
                    size="small"
                    placeholder="搜索对方户名"
                    prefix-icon="el-icon-user"
                    v-model="SearchoppositeName"
                    >
                    </el-input>
                    <el-input
                    @keyup.enter.native="search"
                    size="small"
                    placeholder="搜索"
                    prefix-icon="el-icon-search"
                    v-model="purSearch"
                    >
                    </el-input>
                </div>
            </div>
            <!-- 表格 -->
            <div class="pur-table">
                <el-table
                v-loading="loading"
                element-loading-text = "加载中..."
                height="100%"
                border
                :data="tableData"
                :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                style="width: 100%;overflow:visible !important;"
                @selection-change="purChange">
                    <el-table-column type="selection" width="40"> </el-table-column>
                    <el-table-column prop="id" v-if="show">
                    </el-table-column>
                    <el-table-column prop="accountName" label="账户名称" sortable :show-overflow-tooltip="true" width="110">
                        <template slot-scope="scope">
                            <span
                                style="color: #cd1724;cursor: pointer;"
                                @click="purSee(scope.$index, scope.row)"
                                >{{ scope.row.accountName }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="transTime" label="交易日期" sortable width="110" :show-overflow-tooltip="true">
                    </el-table-column>
                    <el-table-column prop="oppositeName" label="对方户名" sortable :show-overflow-tooltip="true" width="110">
                    </el-table-column>
                    <el-table-column label="可核销金额" sortable :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="scope.row.khxje != null">{{scope.row.khxje|currency}} 元</span>
                            <span v-else-if="scope.row.money && scope.row.khxje == null">{{scope.row.money|currency}} 元</span>
                        </template>  
                    </el-table-column>
                    <el-table-column label="实际核销金额" sortable :show-overflow-tooltip="true">
                        <template slot-scope="scope">
                            <div v-if="scope.row.sjhxblone">
                                <el-input size="small" v-model="CollectionSjhxje" v-focus @blur="CollectionloseFcous(scope.$index, scope.row)"></el-input>
                            </div>
                            <div v-else>
                                <span v-if="scope.row.sjhxje != null">{{scope.row.sjhxje|currency}} 元</span>
                                <span v-else-if="scope.row.khxje !=null && scope.row.sjhxje == null">{{scope.row.khxje|currency}} 元</span>
                                <span v-else-if="scope.row.money && scope.row.sjhxje == null && scope.row.khxje == null">{{scope.row.money|currency}} 元</span>
                                <i class="el-icon-edit" @click="Colleclick(scope.row)" style="color:#CF1724;float:right;margin-right:10px;line-height: 23px;cursor:pointer"></i>
                            </div>
                        </template>  
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页 -->
            <div class="pur-pages" v-if="paginationshow">
                <div style="float:right;width:75%;">
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
            </div>
            <!-- 抽屉 查看功能 -->
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
                                    <label>采购结算：</label>
                                    <span v-if="seeData.relationFlag == true" @click="detailsPurAll" style="color:#cf1724; cursor: pointer;">详情</span>
                                    <span v-else>---</span>
                                </li>
                                <li>
                                    <label>账户名称：</label>
                                    <span>{{seeData.accountName}}</span>
                                </li>
                                <li>
                                    <label>账户号码：</label>
                                    <span>{{seeData.account}}</span>
                                </li>
                                <li>
                                    <label>交易日期：</label>
                                    <span>{{seeData.transTime}}</span>
                                </li>
                                <li>
                                    <label>金额：</label>
                                    <span v-if="seeData.money">{{seeData.money|currency}} 元</span>
                                </li>
                                <li>
                                    <label>对方户名：</label>
                                    <span>{{seeData.oppositeName}}</span>
                                </li>
                                <li>
                                    <label>对方账号：</label>
                                    <span>{{seeData.oppositeAccount}}</span>
                                </li>
                                <li>
                                    <label>摘要：</label>
                                    <span>{{seeData.absContent}}</span>
                                </li>
                                <li>
                                    <label>用途：</label>
                                    <span v-if="seeData.purpose">{{seeData.purpose}}</span>
                                    <span v-else>---</span>
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
                                <a v-for="(v,k) in filelist" :key="k" :href="http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + v.fileName" @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"><p>{{v.name}}</p></a>
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
                </el-tabs>
                <div class="xia">《</div>
                </div>
            </el-drawer>
            <!-- 采购付款单详情一级抽屉 -->
            <el-drawer
            size="350px"
            title="采购结算"
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
                <el-col :span="10"><span style="margin-left:5px;" v-if="totalpurdrawersize">{{totalpurdrawersize|currency}} 元</span></el-col>
                </el-row>
            </div>
            </el-drawer>
            <!-- 采购付款单详情二级抽屉 -->
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
                    <span>采购合同 :</span>
                    <p>{{innerDrawerData.contractName}}</p>
                </dd>
                <dd>
                    <span>供应商 :</span>
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
                    <p>{{innerDrawerData.unitPrice}} 元 / 吨</p>
                </dd>
                <dd>
                    <span>结算总价 :</span>
                    <p v-if="innerDrawerData.sumPrice">{{innerDrawerData.sumPrice|currency}} 元</p>
                </dd>
                </dl>
            </el-drawer>
        </div>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
    data() {
        return {
            http:'',
            show:false,
            purTitle: "",                                   // 标题
            logo: "",                                       // 抽屉标题

            /* 
            核销功能数据
            */
            heceNum:0,                                      // 核销差额
            amountPayable:0,                                // 应付金额       
            newSjhxje:'',                                   // 采购结算input
            CollectionSjhxje:'',                            // 采购付款单input
            // 表单验证
            rules: {
                supplierName: [ 
                    { required: true, message: '请选择客户名称', trigger: 'change' }
                ], 
            },
            // 新建核销表单数据
            AddData: {
                supplierName:'',
            },  

            /* 
            采购结算数据列表
            */
            settlementmoeny:0,                              // 结算合计金额
            // settlementdataAva:false,
            settlementloading:false,
            settlementcurrent:1,                            // 当前页数
            settlementsize:50,                              // 显示条数 
            settlementtotal:0,                              // 总条数
            jshxpurSearch:"",                               // 全部搜索
            SearchoppositeNameone:'',                       // 搜索所属月份
            settlementtableData: [],                        // 表格数据
            multiplesettlement: [],                         // 选中数据
            settlementdrawer: false,                        // 查看抽屉
            settlementdID:'',                               // 表格单条数据id
            settlementdrawerss: false,                      // 付款单详情抽屉开关
            settlementpurdrawersize:[],                     // 流水信息
            settlementtotalpurdrawersize:0,                 // 合计
            settlementinnerDrawer: false,                   // 付款单详情里面抽屉
            settlementinnerDrawerData:{},                   // 查看数据
            settlementdraweroness: false,                   // 进项发票详情抽屉开关
            settlementpurdrawersizeoness:[],                // 流水信息
            settlementtotalpurdrawersizeoness:0,            // 合计
            settlementinnerDraweroness: false,              // 进项发票详情里面抽屉
            settlementinnerDraweronessData:{},              // 查看数据
            settlementfilelist:[],                          // 文件列表
            settlementactiveName: "first",                  // 查看tab栏
            // 查看数据
            settlementseeData:{},  
            // 查看 动态
            settlementtimeData:[
                // {
                //   createtime:'2021-01-02',
                //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ], 

            /* 
            采购付款单数据列表
            */
            banks:'',                           
            optionsBanks:[],                    // 银行搜索
            purSearch:'',                       // 全部搜索
            keyupTime: '',                      // 时间搜索
            SearchoppositeName:'',              // 对方户名搜索
            contmoeny:0,                        // 付款单合计金额  
            paginationshow:true,                // 无数据隐藏分页
            // dataAva:false,
            loading:false,
            current:1,                          // 当前页数
            size:50,                            // 显示条数
            total:0,                            // 总条数
            tableData: [],                      // 表格数据
            multipleSelection: [],              // 选中数据
            drawer: false,                      // 查看抽屉
            dID:'',
            drawerss: false,                    // 采购付款单详情抽屉开关
            purdrawersize:[],                   // 流水信息
            totalpurdrawersize:0,               // 合计
            innerDrawer: false,                 // 采购付款单详情里面抽屉
            innerDrawerData:{},                 // 查看数据
            filelist:[],                        // 文件列表
            activeName: "first",                // 查看tab栏
            seeData:{},                         // 查看数据
            // 查看 动态
            timeData:[
                // {
                //   createtime:'2021-01-02',
                //   operation:'操作内容 ：史石平新建了合同 CG-21010801合同订单：CG-21010801'
                // }
            ], 

            /* 
            供应商弹窗数据列表
            */
            rowData:'',                         // 选中数据
            paginationshowtwo:true,             // 无数据隐藏分页(供应商)
            currenttwo:1,                   
            sizetwo:4,                      
            totaltwo:0,                     
            purSearchtwo:'',                               
            dialogSelectUnit:false,             // 选择单位弹窗
            dialogtextarea:'',                  // 显示单位内容框
            CompanytableData:[],                // 选择单位表格数据
            radio:'',                           // 选择单位单选框
            show:false,                         // 选择单位表格隐藏 id
            CompanyDatashow:false,              // 新建单位表单显示隐藏
            CompanyNewData:{},                  // 新建单位表单
            gid:'',                             // 表格id
            // 新建单位表单验证
            rels: {
                custType: [
                    { required: true, message: '请选择类型', trigger: 'blur' },
                ],
                enterpriseName: [
                    { required: true, message: '请输入单位名称', trigger: 'blur' },
                ],
            },
        };
    },
    created() {},
    mounted(){
        this.purTitle = this.$route.meta.title;                 // 标题
        this.logo = sessionStorage.getItem('companyName')       // 企业名称
        this.http  = this.$store.state.http;                    // IP地址
        this.settlementgetpurcon();                             // 采购结算
        this.getpurcon();                                       // 采购付款单
        // 获取银行搜索列表
        this.api.chainCap.banks()
        .then(res=>{
            if(res.data.code == 200){
                this.optionsBanks = res.data.data;
            }
        })
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
        // 结算实际核销修改金额事件
        cellclick(row){
            row.sjhxblone = true;
        },
        // 付款实际核销修改金额事件
        Colleclick(row){
            row.sjhxblone = true;
        },
        // 结算input失去焦点
        loseFcous(index,val){
            let eagstr = /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/;
            if(eagstr.test(this.newSjhxje)){
                if(val.khxje != null && val.khxje >= this.newSjhxje*1 > 0){
                    val.sjhxblone = false;
                    val.sjhxje = this.newSjhxje *1;
                    this.newSjhxje = '';
                    this.settlementmoeny = 0;
                    for(let i=0;i<this.multiplesettlement.length;i++){
                        this.settlementmoeny += this.multiplesettlement[i].sjhxje;
                    }
                }else if(val.khxje == null && val.sumPrice >= this.newSjhxje*1 > 0){
                    val.sjhxblone = false;
                    val.sjhxje = this.newSjhxje *1;
                    this.newSjhxje = '';
                    this.settlementmoeny = 0;
                    for(let i=0;i<this.multiplesettlement.length;i++){
                        this.settlementmoeny += this.multiplesettlement[i].sjhxje;
                    }
                }else{
                    this.$message.error('实际核销金额错误,原因：实际核销金额大于可核销金额');
                }
                // 合销差额
                this.heceNum = 0;
                this.heceNum = this.settlementmoeny - this.contmoeny;
            }else{
                val.sjhxblone = false;
            }
        },
        // 付款单input失去焦点
        CollectionloseFcous(index,val){
            let eagstr = /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/;
            if(eagstr.test(this.CollectionSjhxje)){
                if(val.khxje != null && val.khxje >= this.CollectionSjhxje * 1 > 0){
                    val.sjhxblone = false;
                    val.sjhxje = this.CollectionSjhxje *1;
                    this.CollectionSjhxje = '';
                    this.contmoeny = 0;
                    for(let i=0;i<this.multipleSelection.length;i++){
                        this.contmoeny += this.multipleSelection[i].sjhxje;
                    }
                }else if(val.khxje == null && val.money >= this.CollectionSjhxje*1 > 0){
                    val.sjhxblone = false;
                    val.sjhxje = this.CollectionSjhxje *1;
                    this.CollectionSjhxje = '';
                    this.contmoeny = 0;
                    for(let i=0;i<this.multipleSelection.length;i++){
                        this.contmoeny += this.multipleSelection[i].sjhxje;
                    }
                }else{
                    this.$message.error('实际核销金额错误,原因：实际核销金额大于可核销金额');
                }
                // 合销差额
                this.heceNum = 0;
                this.heceNum = this.settlementmoeny - this.contmoeny;
            }else{
                val.sjhxblone = false;
            }
        },
        // 核销表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    if(Math.abs(this.settlementmoeny) === Math.abs(this.contmoeny)){
                        let res;
                        if(this.settlementmoeny){
                            let num = this.settlementmoeny.toFixed(2);
                            let prev = num.toString().split('.')[0]; // 获取整数部分
                            let next = num.toString().split('.')[1];
                            res = prev.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') + '.' + next;
                        }else if(this.settlementmoeny ==0 || this.contmoeny == 0){
                            res = 0;
                        }
                        this.$confirm('核销金额共计：'+ res + ' 元', '提示', {
                        confirmButtonText: '确定',
                        cancelButtonText: '取消',
                        type: 'warning',
                        showClose:false,
                        }).then(() => {
                            let consts = [];
                            for(let i=0;i<this.multipleSelection.length;i++){
                                let constObj = {
                                    detailId:this.multipleSelection[i].id,
                                    detailMoney:this.multipleSelection[i].sjhxje
                                }
                                consts.push(constObj)
                            }
                            for(let i=0;i<this.multiplesettlement.length;i++){
                                let constObjs = {
                                    statementId:this.multiplesettlement[i].id,
                                    settMoney:this.multiplesettlement[i].sjhxje
                                }
                                consts.push(constObjs)
                            }
                            let DataObj = {
                                counts:consts,
                                decreaseType:2,
                                details:this.multipleSelection,
                                sments:this.multiplesettlement,
                                supplierId:this.gid
                            }
                            this.api.chainCap.writeoff(DataObj)
                            .then(res=>{
                                if(res.data.code == 200){
                                    var that = this;
                                    that.$message({
                                        type: "success",
                                        message: "核销成功!",
                                        duration:300,  
                                        onClose(){
                                            that.tableData = [];
                                            that.settlementtableData = [];
                                            that.getpurcon();
                                            that.current = 1;
                                            that.settlementgetpurcon();
                                            that.settlementcurrent = 1;
                                            that.api.supplier.one({id:that.gid})
                                            .then(res=>{
                                                if(res.data.code == 200){
                                                    that.amountPayable = res.data.data.amountPayable*1
                                                }
                                            })
                                        }
                                    });
                                }
                            })
                        }).catch(() => {
                            this.$message({
                                type: 'info',
                                message: '已取消核销'
                            });          
                        });
                    }else if(Math.abs(this.settlementmoeny) != Math.abs(this.contmoeny)){
                        this.$alert('实际核销金额不一致，请检查', '提示', {
                            confirmButtonText: '确定',
                            type: 'warning',
                            showClose:false,
                        });
                    }
                }
            });
        },

        /*
        采购付款单列表展示功能
        */
        // 表格数据
        getpurcon(){
            this.loading = true;
            // 传参
            var listByPageData = {
                detype:2,
                khxFlag:true,
                bankName:this.banks,
                transTime:this.keyupTime,
                oppositeName:this.SearchoppositeName,
                searchStr:this.purSearch,
                page:{
                current:this.current,
                size:this.size
                }
            }
            // 渲染表格
            this.api.chainCap.all(listByPageData)
            .then(res=>{
                this.loading = false;
                if(res.data.code == 200){
                if(res.data.data.records != ""){
                    for(let i=0;i<res.data.data.records.length;i++){
                        res.data.data.records[i].sjhxblone = false;
                    }
                    this.tableData = res.data.data.records;
                    this.total = res.data.data.total;
                }else{
                    // this.dataAva = true;
                    this.paginationshow = false;
                }
                }
            })
        },
        // 银行名称搜索
        keyupBanks() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 日期搜索
        DataTime() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 对方户名
        ClickOppositeName() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 全部搜索
        search() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 采购结算详情关闭按钮
        drawerssClose(){
            this.innerDrawer = false;
            this.drawerss = false;
        },
        // 表格查看采购结算详情
        detailsPur(data){
            this.purdrawersize = [];
            this.totalpurdrawersize = 0;
            this.api.chainCap.oneAll({id:data.id})
            .then(res=>{
                this.purdrawersize = res.data.data
                // 流水明细合计
                for(let i=0;i<this.purdrawersize.length;i++){
                    this.totalpurdrawersize += this.purdrawersize[i].sumPrice
                }
                let a = this.totalpurdrawersize;
                this.totalpurdrawersize = a.toFixed(2) * 1;
            })
            this.drawer = false;
            this.drawerss = true;
        },
        // 表格查看采购结算详情
        detailsPurAll(){
            this.purdrawersize = [];
            this.totalpurdrawersize = 0;
            this.api.chainCap.oneAll({id:this.dID})
            .then(res=>{
                this.purdrawersize = res.data.data
                // 流水明细合计
                for(let i=0;i<this.purdrawersize.length;i++){
                    this.totalpurdrawersize += this.purdrawersize[i].sumPrice
                }
                let a = this.totalpurdrawersize;
                this.totalpurdrawersize = a.toFixed(2) * 1;
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
        // 表格选中数据
        purChange(val) {
            this.multipleSelection = val;
            this.contmoeny = 0;
            for(let i=0;i<this.multipleSelection.length;i++){
                if(this.multipleSelection[i].sjhxje == null && this.multipleSelection[i].khxje == null){
                    this.contmoeny += this.multipleSelection[i].money;
                }else if(this.multipleSelection[i].sjhxje == null && this.multipleSelection[i].khxje != null){
                    this.contmoeny += this.multipleSelection[i].khxje;
                }else if(this.multipleSelection[i].sjhxje != null){
                    this.contmoeny += this.multipleSelection[i].sjhxje;
                }
            }
            // 合销差额
            this.heceNum = 0;
            this.heceNum = this.settlementmoeny - this.contmoeny;
        },
        // 查看
        purSee(index, row) {
            this.drawer = true;
            this.api.chainCap.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                    this.seeData = res.data.data;
                    this.dID = res.data.data.id;
                }
            })
            // 查看文件
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
            // console.log(tab.label);
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "动态") {
                // console.log(2);
            } else if (tab.label == "团队") {
                // console.log(3);
            }
        },


        /*
        采购结算列表展示功能
        */
        // 采购结算表格数据
        settlementgetpurcon(){
            this.settlementloading = true;
            // 传参
            var listByPageData = {
                type:1,
                month:this.SearchoppositeNameone,
                supplierId:this.rowData.id,
                searchStr:this.jshxpurSearch,
                khxFlag:true,
                page:{
                current:this.settlementcurrent,
                size:this.settlementsize
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageData)
            .then(res=>{
                this.settlementloading = false;
                if(res.data.code == 200){
                    for(let i=0;i<res.data.data.records.length;i++){
                        res.data.data.records[i].sjhxblone = false;
                    }
                    this.settlementtableData = res.data.data.records;
                    this.settlementtotal = res.data.data.total;
                }
            })
        },
        // 所属月份搜索
        ClickOppositeNameone() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },
        // 搜索全部功能
        jshxpurSearchs() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },
        // 采购付款单详情关闭按钮
        settlementdrawerssClose(){
            this.settlementinnerDrawer = false;
            this.settlementdrawerss = false;
        },
        // 表格查看采购付款单详情
        detailsPur(data){
            this.settlementpurdrawersize = [];
            this.settlementtotalpurdrawersize = 0;
            this.api.capital.oneAll({id:data.id})
            .then(res=>{
                this.settlementpurdrawersize = res.data.data
                // 流水明细合计
                for(let i=0;i<this.settlementpurdrawersize.length;i++){
                this.settlementtotalpurdrawersize += this.settlementpurdrawersize[i].money
                }
                let a = this.settlementtotalpurdrawersize;
                this.settlementtotalpurdrawersize = a.toFixed(2) * 1;
            })
            this.settlementdrawer = false;
            this.settlementdrawerss = true;
        },
        // 查看采购付款单
        settlementdetailsPurAll(){
            this.settlementpurdrawersize = [];
            this.settlementtotalpurdrawersize = 0;
            this.api.capital.oneAll({id:this.settlementdID})
            .then(res=>{
                this.settlementpurdrawersize = res.data.data
                // 流水明细合计
                for(let i=0;i<this.settlementpurdrawersize.length;i++){
                this.settlementtotalpurdrawersize += this.settlementpurdrawersize[i].money;
                }
                let a = this.settlementtotalpurdrawersize;
                this.settlementtotalpurdrawersize = a.toFixed(2) * 1;
            })
            this.settlementdrawer = false;
            this.settlementdrawerss = true;
        },
        // 二级抽屉数据回显
        settlementdetailsdrawer(index){
            for(let i=0;i<this.settlementpurdrawersize.length;i++){
                if(this.settlementpurdrawersize[i].id == index){
                this.settlementinnerDrawerData = this.settlementpurdrawersize[i]
                }
            }
            this.settlementinnerDrawer = true;
        },
        // 进项发票关闭按钮
        settlementdraweronessClose(){
            this.settlementdraweroness = false;
            this.settlementinnerDraweroness = false;
        },
        // 表格查看进项发票详情
        detailsPuroness(data){
            this.settlementpurdrawersizeoness = [];
            this.settlementtotalpurdrawersizeoness = 0;
            this.api.capital.twoAll({id:data.id})
            .then(res=>{
                this.settlementpurdrawersizeoness = res.data.data
                // 流水明细合计
                for(let i=0;i<this.settlementpurdrawersizeoness.length;i++){
                this.settlementtotalpurdrawersizeoness += this.settlementpurdrawersizeoness[i].includMoney
                }
                let a = this.settlementtotalpurdrawersizeoness;
                this.settlementtotalpurdrawersizeoness = a.toFixed(2) * 1;
            })
            this.settlementdrawer = false;
            this.settlementdraweroness = true;
        },
        // 表格查看进项发票详情
        settlementdetailsPuronessAll(){
            this.settlementpurdrawersizeoness = [];
            this.settlementtotalpurdrawersizeoness = 0;
            this.api.capital.twoAll({id:this.settlementdID})
            .then(res=>{
                this.settlementpurdrawersizeoness = res.data.data
                // 流水明细合计
                for(let i=0;i<this.settlementpurdrawersizeoness.length;i++){
                this.settlementtotalpurdrawersizeoness += this.settlementpurdrawersizeoness[i].includMoney
                }
                let a = this.settlementtotalpurdrawersizeoness;
                this.settlementtotalpurdrawersizeoness = a.toFixed(2) * 1;
            })
            this.settlementdrawer = false;
            this.settlementdraweroness = true;
        },
        // 二级抽屉数据回显
        settlementdetailsdraweroness(index){
            for(let i=0;i<this.settlementpurdrawersizeoness.length;i++){
                if(this.settlementpurdrawersizeoness[i].id == index){
                this.settlementinnerDraweronessData = this.settlementpurdrawersizeoness[i]
                }
            }
            this.settlementinnerDraweroness = true;
        },
        // 表格选中数据
        settlementChange(val) {
            this.multiplesettlement = val;
            this.settlementmoeny = 0;
            for(let i=0;i<this.multiplesettlement.length;i++){
                if(this.multiplesettlement[i].sjhxje == null && this.multiplesettlement[i].khxje == null){
                    this.settlementmoeny += this.multiplesettlement[i].sumPrice;
                }else if(this.multiplesettlement[i].sjhxje == null && this.multiplesettlement[i].khxje != null){
                    this.settlementmoeny += this.multiplesettlement[i].khxje;
                }else if(this.multiplesettlement[i].sjhxje != null){
                    this.settlementmoeny += this.multiplesettlement[i].sjhxje;
                }
            }
            // 合销差额
            this.heceNum = 0;
            this.heceNum = this.settlementmoeny - this.contmoeny;
        },
        settlementkeyId(row){
            return row.id;
        },
        // 查看
        settlementpurSee(index, row) {
            this.settlementdrawer = true;
            this.api.capital.one({id:row.id})
            .then(res=>{
                if(res.data.code == 200){
                this.settlementseeData = res.data.data;
                this.settlementdID = res.data.data.id;
                }
            })
            // 查看文件
            this.api.file.one({relationId:row.id})
            .then(res=>{
                this.settlementfilelist = res.data.data
            })
        },
        settlementchangeActive(e) {
            e.currentTarget.className = 'active';
        },
        settlementremoveActive(e) {
            e.currentTarget.className = '';
        },
        // 抽屉
        settlementhandleClick(tab) {
            // console.log(tab.label);
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "动态") {
                // console.log(2);
            } else if (tab.label == "团队") {
                // console.log(3);
            }
        },

        /* 
        供应商单位的选择功能
        */
        // 弹窗表单提示语
        renderHeaderMethods(h,{column}){
            return h(
                'div',
                {
                    style:'display:flex;margin:auto;'
                },
                [
                    h('span',column.label),
                    h('promptMessages',{
                        props:{
                            messages:['相关操作在合作方管理中操作']
                        }
                    })
                ]
            )
        },
        // 供应商单位表格数据
        allsuppliertwo(){
            // 传参
            var listByPageDatatwo = {
                partnerType:1,
                searchStr:this.purSearchtwo,
                page:{
                    current:this.currenttwo,
                    size:this.sizetwo
                }
            }
            // 渲染表格
            this.api.supplier.all(listByPageDatatwo)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.CompanytableData = res.data.data.records
                        this.totaltwo = res.data.data.total
                    }else{
                        this.paginationshowtwo = false;
                    }
                }
            })
        },
        // 供应商单位搜索功能(供应商)
        searchtwo() {
            this.CompanytableData = [];
            this.allsuppliertwo();
            this.currenttwo = 1;
        },
        // 供应商单位 选择单位表格
        getCurrentRow(row){
            this.gid = row.id;
            this.rowData = row;
            this.dialogtextarea = row.enterpriseName;   // 获取选中单位名称数据
        },
        // 隐藏显示新建
        Companyshow(){
            this.CompanyDatashow = true
        },
        // 新建单位提交表单
        CompanyNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.CompanyNewData.partnerType = 1;
                    this.api.supplier.add(this.CompanyNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.allsuppliertwo();
                                    that.CompanyDatashow = false;
                                    that.CompanyNewData = {};
                                }
                            });
                        }
                    })
                }
            });
        },
        // 新建单位取消提交表单
        CompanyNewShow() {
            this.CompanyNewData = {};
            this.CompanyDatashow = false
        },
        // 关联供应商单位弹窗
        dialogSelect(){
            this.allsuppliertwo();                              // 渲染表格
            this.dialogSelectUnit=true;                         // 开启弹窗
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.dialogSelectUnit == true){
                this.CompanyDatashow = false;
            }
        },
        // 关联供应商单位应用按钮事件
        salesUnitApplication(){
            this.amountPayable = this.rowData.amountPayable;  // 获取应付金额
            this.AddData.supplierName = this.dialogtextarea;        // 获取单位名称，回显到表单
            this.dialogSelectUnit = false;                          // 关闭整个弹窗
            this.radio = "";                                        // 清空单选框
            this.dialogtextarea = "";                               // 清空显示单位内容框
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },
        // 关联供应商单位 关闭弹窗
        BuyerhandleClose() {
            this.radio = "";                                    // 清空单选框
            this.dialogtextarea = "";                           // 清空显示单位内容框
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
        },
    },
};
</script>
<style lang="less" scoped>
// 付款单详情样式
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
// 子付款单详情样式
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
.details-one{
  width: 50%;
  height: 100%;
  background-color: #FBFBFB;
  padding: 20px 0 0 20px;
  font-size: 14px;
  float: left;
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
.details-two{
  float: left;
  width: 50%;
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
    .receivable{
        width:100%;
        height: 78px;
        // background-color: #f6f6f6;
        border: 1px solid #E4E4E4;
        padding:23px 35px;
        box-sizing: border-box;
        overflow: hidden;
        .receivable-conent{
            padding-left:32px;
            img{
                width: 16px;
                height: 16px;
                margin-bottom:2px;
            }
            b{
                font-size:28px;
                color: #333;
                margin:0 5px;
            }
            span{
                font-size:16px;
                color: #333;
            }
            .el-button{
               position: relative;
               bottom: 3px;
            }
        }
    }
    .rece{
        width: 49%;
        height: 87.5%;
        // height: auto;
        border: 1px solid #E4E4E4;
        margin-top:12px;
        float: left;
        overflow: hidden;
        box-sizing: border-box;
        .pur-title{
            width: 100%;
            height: 100px;
            padding:9px 24px 9px;
            box-sizing: border-box;
            position: relative;
            b{
                font-size: 18px;
                color: #333;
                line-height: 32px;
            }
            .sjhxjeMoeny{
                position: absolute;
                top: 0px;
                right: 32px;
                span{
                    font-size: 14px;
                    font-weight: 600;
                    color: #333;
                    line-height:50px;
                }  
            }
            .Advanced{
                width: 100%;
                height: 50px;
                padding:9px 0;
                box-sizing: border-box;
                .el-select{
                    float: left;
                    width: 24%;
                    margin-right: 1%;
                }
                .el-input {
                    width: 24%;
                    float: left;
                    margin-right: 1%;
                }
                .el-date-picker{
                    width: 25%;
                    float: left;
                    margin-right: 1%;
                }
            }
        }
        .pur-table{
            width: 100%;
            height: 75%;
        }
        .pur-pages {
            width: 100%;
            height: 42px;
            padding:10px 18px 0 24px;
            box-sizing: border-box; 
            .el-pagination {
                float: right;
            }
        }
    }
    .rece-right{
        width: 2%;
        height: 89%;
        float: left;
        transform: rotate(180deg);
        overflow: hidden;
        position: relative;
        box-sizing: border-box;
        span{
            position: absolute;
            right:20%;
            font-size: 24px;
            color: #999;
            font-weight: 300;
            line-height: 800px;
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
    }
    .active{
      color: #EA222E;
    }
  }
}
// 查看抽屉样式结束
// 关联弹窗样式开始
.el-dialog__body{
    padding: 0;
}
.Wdialog{
    width: 100%;
    border: 1px solid #E0E0E0;
    margin-top: 16px;
}
// 选项卡
.concentBox {
    width: 100%; 
    position: relative;
}
.textbutton{
    width:56px;
    height: 32px;
    position: absolute;
    right:25px;
    bottom: 10px;
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
  // 分页
.pur-page {
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
  // 新建表单
.CompanyData{
    width: 100%;
    padding: 8px 16px;
    position: absolute;
    left: 0;
    bottom:115px;
}
// 关联弹窗样式结束
</style>