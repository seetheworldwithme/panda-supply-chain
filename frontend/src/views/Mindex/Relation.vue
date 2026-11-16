<template>
    <!-- 四流合一单页面 -->
    <div>
        <!-- 头部开始 -->
        <!-- <el-row class="pur-top">
            <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
        </el-row> -->
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
            <!-- 合同信息列表 -->
            <div class="rece">
                <!-- 高级搜索功能 -->
                <div class="pur-title-header">
                    <b>合同 ({{ settlementtableData.length }}份) </b>
                    <div>
                        <el-input @keyup.enter.native="jshxpurSearchs" size="small" placeholder="搜索"
                            prefix-icon="el-icon-search" v-model="jshxpurSearch">
                        </el-input>
                        <el-select v-model="SelectValue" size="small" @change="Contractchange">
                            <el-option
                            v-for="item in SelectOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value">
                            </el-option>
                        </el-select>
                        <!-- <el-select v-model="SelectValue" size="small" @change="Contractchange">
                            <el-option label="全部合同" :value="0"></el-option>
                            <el-option label="采购合同" :value="1"></el-option>
                            <el-option label="销售合同" :value="2"></el-option>
                        </el-select> -->
                    </div>
                </div>
                <!-- 表格 -->
                <div class="pur-table">
                    <el-table v-loading="settlementloading" element-loading-text="加载中..." height="100%" border size="medium"
                        :data="settlementtableData" highlight-current-row @current-change="settlementhandleCurrentChange"
                        ref="mulsetTable" :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;cursor: pointer;">
                        <!-- <el-table-column type="index" width="90" label="序号" :index="indexMethod"
                            :show-overflow-tooltip="true">
                        </el-table-column> -->
                        <el-table-column prop="contractCode" label="合同编号" sortable :show-overflow-tooltip="true"
                            width="220">
                            <template slot-scope="scope">
                                <span style="color: #cd1724; cursor: pointer;" @click.stop="purSee(scope.row)">{{
                                    scope.row.contractCode }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="" label="电子签章" :show-overflow-tooltip="true" min-width="100">
                            <template slot-scope="scope">
                                <span v-if="scope.row.isEsign">开启</span>
                                <span v-else>关闭</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="签署合同预览" :show-overflow-tooltip="true" width="120">
                            <template slot-scope="scope">
                                <span v-if="!scope.row.isEsign"> -- </span>
                                <span v-else @click.stop="confselSign(scope.row)"
                                    style="color:#CF1724;cursor: pointer;">签署预览</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="签署状态" min-width="100">
                            <template slot-scope="scope">
                                <span v-if="scope.row.isEsignComplete">已完成</span>
                                <span v-else>未完成</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="签署方信息" :show-overflow-tooltip="true" width="100">
                            <template slot-scope="scope">
                                <span v-if="!scope.row.isEsign"> -- </span>
                                <span v-else @click.self="seeone(scope.row)"
                                    style="color:#CF1724;cursor: pointer;">查看</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="systemCode" label="系统编号" sortable width="260" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="contractName" label="合同名称" :show-overflow-tooltip="true" width="140">
                        </el-table-column>
                        <el-table-column prop="signingTime" label="签订时间" sortable width="120" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="month" label="所属月份" sortable width="120" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="salesUnit" label="销售单位" :show-overflow-tooltip="true" width="260">
                        </el-table-column>
                        <el-table-column prop="salesUnit" label="销售单位" :show-overflow-tooltip="true" width="260">
                        </el-table-column>
                        <el-table-column prop="buyerUnit" label="购货单位" :show-overflow-tooltip="true" width="260">
                        </el-table-column>
                        <el-table-column prop="deliveryPlace" label="发货地" :show-overflow-tooltip="true" width="260">
                        </el-table-column>
                        <el-table-column prop="receiptPlace" label="收货地" :show-overflow-tooltip="true" width="260">
                        </el-table-column>
                        <el-table-column prop="proBigType" label="产品大类" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope" v-if="scope.row.proBigType">
                                {{ scope.row.proBigType | protype }}
                            </template>
                        </el-table-column>
                        <el-table-column prop="proSmaType" label="产品小类" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope" v-if="scope.row.proSmaType">
                                {{ scope.row.proSmaType | promiddletype }}
                            </template>
                        </el-table-column>
                        <el-table-column label="产品名称" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.proSmaType == 20">{{ scope.row.proRemark }}</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="contractNum" label="合同数量" width="120" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column label="货品单价" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span>{{ scope.row.unitPrice }} 元 / 吨</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="总价金额" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice != null && scope.row.sumPrice != 0">{{ scope.row.sumPrice |
                                    currency }}
                                    元</span>
                                <span v-else-if="scope.row.sumPrice == 0">0.00 元</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="价格类型" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.priceType">{{ scope.row.priceType | pricetype }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="付款方式" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.payMethod">{{ scope.row.payMethod | paymentmethod }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="支付方式" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.payType">{{ scope.row.payType | paymenttype }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="运输模式" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.transportMode">{{ scope.row.transportMode | transportmode }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="承运工具" width="120" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.carrier">{{ scope.row.carrier | carrier }}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="备注" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.remark">{{ scope.row.remark }}</span>
                                <span v-else>---</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 分页 -->
                <div class="pur-pages">
                    <el-pagination background :page-sizes="[10, 30, 50, 100, 200]" :page-size.sync="settlementsize"
                        :pager-count="5" :current-page.sync="settlementcurrent"
                        layout="total, sizes, prev, pager, next, jumper" :total="settlementtotal"
                        @current-change='settlementgetpurcon' @size-change='settlementgetpurcon'>
                    </el-pagination>
                </div>
            </div>
            <!-- <div class="rece-right"><span>《</span></div> -->
            <!-- 结算 发票 流水 物流tab列表 -->
            <div class="rece" style="margin-top:24px;">
                <!-- 高级搜索功能 -->
                <div class="pur-title">
                    <!-- tab选项卡 -->
                    <div class="totalTab" :class="{ active: item.id === sel }" v-for="item in tabs" :key="item.id"
                        @click="select(item)" style="cursor: pointer;">
                        <span v-if="item.num === null">{{ item.label }}</span>
                        <span v-else>{{ item.label }} ({{ item.num }})</span>
                        <div class="stripnone" :style="item.width" :class="{ strip: item.id === sel }">
                        </div>
                    </div>
                </div>
                <!-- 结算 -->
                <div class="adBox" v-show="sel === 1">
                    <div class="Advanced">
                        <el-input @keyup.enter.native="search" size="small" placeholder="搜索" prefix-icon="el-icon-search"
                            v-model="purSearch">
                        </el-input>
                        <el-input @keyup.enter.native="ClickOppositeName" size="small" :placeholder="'搜索' + customer"
                            prefix-icon="el-icon-user" v-model="SearchoppositeName">
                        </el-input>
                        <el-date-picker size="small" @change="Time" v-model="TimeArray" type="daterange" unlink-panels
                            range-separator="~" start-placeholder="开始日期" end-placeholder="结束日期" value-format="yyyy-MM-dd"
                            style="width:30%;float:right;margin-right:2%;">
                        </el-date-picker>
                    </div>
                    <!-- 结算表格 -->
                    <div class="tab-table">
                        <el-table v-loading="loading" element-loading-text="加载中..." height="100%" border size="medium"
                            :data="tableData" :header-cell-style="{ background: '#F6F6F6', color: '#333' }">
                            <el-table-column prop="systemCode" label="合同编号" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click.stop="JspurSee(scope.$index, scope.row)">{{
                                            scope.row.systemCode }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="goodsName" label="货物名称" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column label="结算时间" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.setStartTime">{{ scope.row.setStartTime | Time }}</span>
                                    <span> ~ </span>
                                    <span v-if="scope.row.setEndTime">{{ scope.row.setEndTime | Time }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="supplierName" :label="customer" sortable :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="month" label="所属月份" sortable :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column label="结算单价" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.unitPrice">{{ scope.row.unitPrice | currency }} 元 / 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="settlementCount" label="结算数量" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span>{{ scope.row.settlementCount }} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="sumPrice" label="结算总价" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.sumPrice != 0">{{ scope.row.sumPrice | currency }} 元</span>
                                    <span v-else>0.00 元</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <div class="sjhxjeMoeny">
                            <span>合计: <span v-if="contmoeny != 0">结算总价金额累计 {{ contmoeny | currency }} 元</span><span v-else>结算总价金额累计 0.00
                                    元</span>，<span v-if="watertcountnum != 0">结算数量累计 {{ watertcountnum | currency }} 吨</span><span
                                    v-else>结算数量累计 0.00 吨</span></span>
                        </div>
                        <div style="float:right;width:auto;">
                            <el-pagination background :page-sizes="[10, 30, 50, 100, 200]" :page-size.sync="size"
                                :pager-count="5" :current-page.sync="current"
                                layout="total, sizes, prev, pager, next, jumper" :total="total" @current-change='getpurcon'
                                @size-change='getpurcon'>
                            </el-pagination>
                        </div>
                    </div>
                </div>
                <!-- 发票 -->
                <div class="adBox" v-show="sel === 2">
                    <div class="Advanced">
                        <el-input @keyup.enter.native="watersearch" size="small" placeholder="搜索"
                            prefix-icon="el-icon-search" v-model="waterpurSearch">
                        </el-input>
                        <el-input @keyup.enter.native="waterClickOppositeName" size="small"
                            :placeholder="'搜索' + Seller + '方企业名称'" prefix-icon="el-icon-user"
                            v-model="waterSearchoppositeName">
                        </el-input>
                        <el-date-picker @change="waterDataTime" size="small" v-model="waterkeyupTime" type="date"
                            value-format="yyyy-MM-dd" placeholder="选择日期">
                        </el-date-picker>
                    </div>
                    <!-- 票据表格 -->
                    <div class="tab-table">
                        <el-table v-loading="waterloading" element-loading-text="加载中..." height="100%" border size="medium"
                            :data="watertableData" :header-cell-style="{ background: '#F6F6F6', color: '#333' }">
                            <el-table-column label="系统编码" sortable :show-overflow-tooltip="true" width="240">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click="FppurSee(scope.$index, scope.row)">{{ scope.row.systemCode }}</span>
                                    <!-- <span v-else>{{scope.row.systemCode}}</span> -->
                                </template>
                            </el-table-column>
                            <el-table-column prop="companyName" :label="Seller + '方企业名称'" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="dutyParagraph" :label="Seller + '方税号'" :show-overflow-tooltip="true"
                                width="170">
                            </el-table-column>
                            <el-table-column prop="billCode" label="发票代码" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="billNumber" label="发票号码" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="openDate" label="开票日期" sortable :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column label="金额" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.baseMoney">{{ scope.row.baseMoney | currency }} 元</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="税率" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.taxRate">{{ scope.row.taxRate }}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="税额" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.taxAmount">{{ scope.row.taxAmount | currency }} 元</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="includMoney" label="含税金额" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.includMoney != 0">{{ scope.row.includMoney | currency }} 元</span>
                                    <span v-else>0.00 元</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="count" label="数量" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.count != null">{{ scope.row.count | currency }} 吨</span>
                                    <span v-else>0.00 吨</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <div class="sjhxjeMoeny">
                            <span>合计: <span v-if="watertcontmoeny != 0">含税金额累计 {{ watertcontmoeny | currency }} 元</span><span
                                    v-else>含税金额累计 0.00 元</span>，<span v-if="watertcount != 0">数量累计 {{ watertcount | currency }}
                                    吨</span><span v-else>数量累计 0.00 吨</span></span>
                        </div>
                        <div style="float:right;width:75%;">
                            <el-pagination background :page-sizes="[10, 30, 50, 100, 200]" :page-size.sync="watertsize"
                                :pager-count="5" :current-page.sync="watertcurrent"
                                layout="total, sizes, prev, pager, next, jumper" :total="waterttotal"
                                @current-change='watertgetpurcon' @size-change='watertgetpurcon'>
                            </el-pagination>
                        </div>
                    </div>
                </div>
                <!-- 流水 -->
                <div class="adBox" v-show="sel === 3">
                    <div class="Advanced">
                        <el-input @keyup.enter.native="newwatersearch" size="small" placeholder="搜索"
                            prefix-icon="el-icon-search" v-model="newwaterpurSearch">
                        </el-input>
                        <el-input @keyup.enter.native="newwaterClickOppositeName" size="small" placeholder="搜索对方账户"
                            prefix-icon="el-icon-user" v-model="newwaterSearchoppositeName">
                        </el-input>
                        <el-select @change="keyupBanks" v-model="lsbanks" filterable placeholder="搜索银行名称" size="small"
                            :clearable='true'>
                            <el-option v-for="item in optionsBanks" :key="item.bank_name" :label="item.bank_name"
                                :value="item.bank_name">
                            </el-option>
                        </el-select>
                        <el-date-picker @change="newwaterDataTime" size="small" v-model="newwaterkeyupTime" type="date"
                            value-format="yyyy-MM-dd" placeholder="选择日期">
                        </el-date-picker>
                    </div>
                    <!-- 收付款单表格 -->
                    <div class="tab-table">
                        <el-table v-loading="newwaterloading" element-loading-text="加载中..." height="100%" border
                            :data="newwatertableData" :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                            style="overflow:visible !important;">
                            <el-table-column label="系统编码" sortable :show-overflow-tooltip="true" width="240">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click="FkpurSee(scope.$index, scope.row)">{{ scope.row.systemCode }}</span>
                                    <!-- <span v-else>{{scope.row.systemCode}}</span> -->
                                </template>
                            </el-table-column>
                            <el-table-column prop="accountName" label="账户名称" :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <el-table-column prop="bankName" label="银行名称" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="account" label="账户号码" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="transTime" label="交易日期" sortable :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="oppositeName" label="对方账户" :show-overflow-tooltip="true">
                            </el-table-column>
                            <el-table-column prop="sjhxje" label="已核销金额" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.sjhxje != 0">{{ scope.row.sjhxje | currency }} 元</span>
                                    <span v-else>0.00 元</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="money" label="金额" sortable :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.money != null">{{ scope.row.money | currency }} 元</span>
                                    <span v-else>0.00 元</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <div class="sjhxjeMoeny">
                            <span>合计: <span v-if="newwatertcontmoeny != 0">已核销金额累计 {{ newwatertcontmoeny | currency }}
                                    元</span><span v-else>已核销金额累计 0.00 元</span>，<span v-if="newwatertcount != 0">金额累计{{
                                        newwatertcount | currency }}
                                    元</span><span v-else>金额累计 0.00 元</span></span>
                        </div>
                        <div style="float:right;width:75%;">
                            <el-pagination background :page-sizes="[10, 30, 50, 100, 200]" :page-size.sync="newwatertsize"
                                :pager-count="5" :current-page.sync="newwatertcurrent"
                                layout="total, sizes, prev, pager, next, jumper" :total="newwaterttotal"
                                @current-change='newwatertgetpurcon' @size-change='newwatertgetpurcon'>
                            </el-pagination>
                        </div>
                    </div>
                </div>
                <!-- 物流 -->
                <div class="adBox" v-show="sel === 4">
                    <div class="Advanced">
                        <el-input @keyup.enter.native="AllTonnageSearch" size="small" placeholder="搜索"
                            prefix-icon="el-icon-search" v-model="TonnagnewwaterpurSearch">
                        </el-input>
                        <el-input @keyup.enter.native="keyupTonnage" size="small" placeholder="搜索卸货地"
                            prefix-icon="el-icon-user" v-model="lsbanksTonnage">
                        </el-input>
                        <el-input @keyup.enter.native="TonnageOppositeName" size="small" placeholder="搜索装货地"
                            prefix-icon="el-icon-user" v-model="TonnageSearchopposite">
                        </el-input>
                        <el-date-picker @change="TonnageDataTime" size="small" v-model="TonnagekeyupTime" type="date"
                            value-format="yyyy-MM-dd" placeholder="选择日期">
                        </el-date-picker>
                    </div>
                    <!-- 物流表格 -->
                    <div class="tab-table">
                        <el-table v-loading="Tonnagnewwaterloading" element-loading-text="加载中..." height="100%" border
                            :data="TonnagwatertableData" :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                            style="overflow:visible !important;">
                            <el-table-column label="运单编码" sortable :show-overflow-tooltip="true" width="260">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click="WlpurSee(scope.$index, scope.row)">{{ scope.row.systemCode }}</span>
                                </template>
                            </el-table-column>

                            <el-table-column label="关联详情">
                            <template slot-scope="scope">
                                <span v-if="scope.row.carriageName == 1" style="color: #cd1724;cursor: pointer;" @click="bddialogSelect(scope.row)">磅单详情</span>
                                <span v-if="scope.row.carriageName == 2" style="color: #cd1724;cursor: pointer;" @click="dialogSelect(scope.row)">运单详情</span>
                            </template>
                            </el-table-column>


                            <!-- <el-table-column label="网货运单" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click="dialogSelect(scope.row.id)">查看</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="过磅单" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span style="color: #cd1724;cursor: pointer;"
                                        @click="bddialogSelect(scope.row.id)">查看</span>
                                </template>
                            </el-table-column> -->


                            <el-table-column label="订单昵称" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.orderName">{{ scope.row.orderName }}</span>
                                    <span v-else>--</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="订单编号" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.orderNumber">{{ scope.row.orderNumber }}</span>
                                    <span v-else>--</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="订单创建时间" :show-overflow-tooltip="true" width="130">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.orderCreateTime">{{ scope.row.orderCreateTime }}</span>
                                    <span v-else>--</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="订单开始时间" :show-overflow-tooltip="true" width="130">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.orderStartTime">{{ scope.row.orderStartTime }}</span>
                                    <span v-else>--</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="订单结束时间" :show-overflow-tooltip="true" width="130">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.orderEndTime">{{ scope.row.orderEndTime }}</span>
                                    <span v-else>--</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="transMethod" label="运输模式" :show-overflow-tooltip="true">
                                <template slot-scope="scope" v-if="scope.row.transMethod">
                                    {{ scope.row.transMethod | transportmode }}
                                </template>
                            </el-table-column>
                            <el-table-column prop="carriageName" label="承运单位" :show-overflow-tooltip="true" width="110">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.carriageName == 1">物流公司</span>
                                    <span v-else-if="scope.row.carriageName == 2">网络货运</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="承运平台" :show-overflow-tooltip="true" width="140">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.carriagePlatform">{{ scope.row.carriagePlatform |
                                        logisticsplatform }}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="logisticsCompany" label="物流单位" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.logisticsCompany">{{ scope.row.logisticsCompany }}</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column prop="transType" label="承运工具" :show-overflow-tooltip="true">
                                <template slot-scope="scope" v-if="scope.row.transType">
                                    {{ scope.row.transType | carrier }}
                                </template>
                            </el-table-column>
                            <!-- <el-table-column prop="contractNumber" label="合同编号" sortable :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <el-table-column prop="baseSalesUnit" label="销售单位" sortable :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <el-table-column prop="baseBuyerUnit" label="购货单位" sortable :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <el-table-column prop="baseDeliveryPlace" label="发货地" sortable :show-overflow-tooltip="true" width="260">
                            </el-table-column> -->
                                            <!-- <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.subjectMatter">{{scope.row.subjectMatter}}</span>
                                <span v-else>---</span>
                            </template>
                            </el-table-column> -->
                            <el-table-column prop="loadingAddress" label="装货地" :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <el-table-column prop="unloadingPoint" label="卸货地" :show-overflow-tooltip="true" width="260">
                            </el-table-column>
                            <!-- <el-table-column label="发货日期" sortable width="200" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <div v-if="scope.row.issStartDate && scope.row.issEndDate">
                                    <span>{{scope.row.issStartDate}}</span><span> ~ </span><span>{{scope.row.issEndDate}}</span>
                                </div>
                                <div v-else>---</div>
                            </template>
                            </el-table-column>
                            <el-table-column label="收货日期" sortable width="200" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <div v-if="scope.row.recStartDate && scope.row.recEndDate">
                                <span>{{scope.row.recStartDate}}</span><span> ~ </span><span>{{scope.row.recEndDate}}</span>
                                </div>
                                <div v-else>---</div>
                            </template>
                            </el-table-column> -->
                            <el-table-column label="车数" :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.vehicles">{{ scope.row.vehicles }} 车</span>
                                    <span v-else>---</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="运费金额合计" :show-overflow-tooltip="true" width="140">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.freightAmount">{{ scope.row.freightAmount | currency }} 元</span>
                                    <span v-if="scope.row.freightAmount == 0 || scope.row.freightAmount == ''">---</span>
                                </template>
                            </el-table-column>
                            <!-- <el-table-column label="吨数合计" :show-overflow-tooltip="true" width="100">
                            <template slot-scope="scope">
                                <span>{{scope.row.subjectCount}} 吨</span>
                            </template>
                            </el-table-column> -->
                            <el-table-column label="原发吨数合计" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span>{{ scope.row.originalTonnage }} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="实收吨数合计" :show-overflow-tooltip="true" width="120">
                                <template slot-scope="scope">
                                    <span>{{ scope.row.tonnagePaid }} 吨</span>
                                </template>
                            </el-table-column>
                            <!-- 暂无数据 -->
                            <!-- <div slot="empty" class="dataAvailable" v-if="dataAva">
                                <div class="dataAva">
                                    <img src="../../assets/zanwushuju.png" alt="">
                                    <p>暂无数据，您可以新建试试</p>
                                    <el-button size="small" style="margin-left: 24px" @click="addPurs"
                                        v-if="$route.meta.add">新建</el-button>
                                </div>
                            </div> -->
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <div class="sjhxjeMoeny">
                            <span>实收吨数合计: </span><span v-if="Tonnage != 0">{{ Tonnage | currency }} 吨</span><span v-else> 0.00
                                吨</span>
                        </div>
                        <div style="float:right;width:75%;">
                            <el-pagination background :page-sizes="[10, 30, 50, 100, 200]"
                                :page-size.sync="Tonnagnewwatertsize" :pager-count="5"
                                :current-page.sync="Tonnagnewwatertcurrent" layout="total, sizes, prev, pager, next, jumper"
                                :total="Tonnagnewwaterttotal" @current-change='Tonnagpurcon' @size-change='Tonnagpurcon'>
                            </el-pagination>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 抽屉 合同查看功能 -->
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
                                    <dd>
                                        <el-col :span="5">
                                            <span>合同附件：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <p v-if="filelist.length == 0"> --- </p>
                                            <a v-for="(v, k) in filelist" :key="k" @mouseenter="changeActive($event)"
                                                @mouseleave="removeActive($event)" @click="HtfileGo(v, seeData)">
                                                <p>{{ v.name }}</p>
                                            </a>
                                        </el-col>
                                    </dd>
                                    <dd>
                                        <el-col :span="5">
                                            <span>合同协议：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <p v-if="enclosurefileList.length == 0"> --- </p>
                                            <p v-else v-for="(v, kindex) in enclosurefileList" :key="kindex"><span
                                                    @mouseenter="changeActive($event)" @mouseleave="removeActive($event)"
                                                    @click="HtfileGo(v, seeData)">{{ v.name }}</span></p>
                                        </el-col>
                                    </dd>
                                </dl>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">基础信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>系统编码：</label>
                                        <span>{{ seeData.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>合同编号：</label>
                                        <span v-if="seeData.contractCode">{{ seeData.contractCode }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>合同名称：</label>
                                        <span>{{ seeData.contractName }}</span>
                                    </li>
                                    <li>
                                        <label>签订时间：</label>
                                        <span>{{ seeData.signingTime }}</span>
                                    </li>
                                    <li>
                                        <label>所属月份：</label>
                                        <span v-if="seeData.month">{{ seeData.month | timeDate }}</span>
                                    </li>
                                    <li>
                                        <label>销售单位：</label>
                                        <span>{{ seeData.salesUnit }}</span>
                                    </li>
                                    <li>
                                        <label>购货单位：</label>
                                        <span>{{ seeData.buyerUnit }}</span>
                                    </li>
                                    <li>
                                        <label>发货地：</label>
                                        <span>{{ seeData.deliveryPlace }}</span>
                                    </li>
                                    <li>
                                        <label>收货地：</label>
                                        <span>{{ seeData.receiptPlace }}</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">产品信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>产品大类：</label>
                                        <span v-if="seeData.proBigType">{{ seeData.proBigType | protype }}</span>
                                    </li>
                                    <li>
                                        <label>产品小类：</label>
                                        <span v-if="seeData.proSmaType">{{ seeData.proSmaType | promiddletype }}</span>
                                    </li>
                                    <li>
                                        <label>产品名称：</label>
                                        <span v-if="seeData.proSmaType == 20">{{ seeData.proRemark }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>合同数量：</label>
                                        <span>{{ seeData.contractNum }} 吨</span>
                                    </li>
                                    <li>
                                        <label>货品单价：</label>
                                        <span>{{ seeData.unitPrice }} 元 / 吨</span>
                                    </li>
                                    <li>
                                        <label>总价金额：</label>
                                        <span v-if="seeData.sumPrice">{{ seeData.sumPrice | currency }} 元</span>
                                    </li>
                                    <li v-for="(v, index) in seeData.list" :key="index">
                                        <label>标的指标：</label>
                                        <div class="see-libox">
                                            <span>硫分St.d% : </span>
                                            <div>
                                                <p v-if="seeData.sulfur">{{ seeData.sulfur }} </p>
                                                <p v-else>--- </p>
                                            </div>
                                            <span>挥发分Vdaf% : </span>
                                            <div>
                                                <span v-if="seeData.volatilization">{{ seeData.volatilization }}
                                                </span><span v-else>--- </span>
                                            </div>
                                            <span>水分Mt% : </span>
                                            <div>
                                                <span v-if="seeData.water">{{ seeData.water }} </span><span v-else>---
                                                </span>
                                            </div>
                                            <span>{{ v.dkey }} : </span>
                                            <div>
                                                <span v-if="v.dvalue">{{ v.dvalue }} </span><span v-else>--- </span>
                                            </div>
                                        </div>
                                    </li>
                                    <li style="clear:both;">
                                        <label>备注：</label>
                                        <span v-if="seeData.remark">{{ seeData.remark }}</span>
                                        <span v-else>---</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24" style="padding-bottom: 48px;">
                                <div class="see-top">扩展信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>价格类型：</label>
                                        <span v-if="seeData.priceType">{{ seeData.priceType | pricetype }}</span>
                                    </li>
                                    <li>
                                        <label>付款方式：</label>
                                        <span v-if="seeData.payMethod">{{ seeData.payMethod | paymentmethod }}</span>
                                    </li>
                                    <li>
                                        <label>支付方式：</label>
                                        <span v-if="seeData.payType">{{ seeData.payType | paymenttype }}</span>
                                    </li>
                                    <li>
                                        <label>运输模式：</label>
                                        <span v-if="seeData.transportMode">{{ seeData.transportMode | transportmode
                                        }}</span>
                                    </li>
                                    <li>
                                        <label>承运工具：</label>
                                        <span v-if="seeData.carrier">{{ seeData.carrier | carrier }}</span>
                                    </li>
                                </ul>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
                <div class="xia">《</div>
            </div>
        </el-drawer>
        <!-- 文件预览弹框 -->
        <el-dialog title="文件预览" :visible.sync="ImgVisible" width="44%" :before-close="ImgClose">
            <div class="imgBox" v-if="itemshow == true" v-loading="loadingimgss" element-loading-text="文件加载中">
                <!-- <iframe :src="url" frameborder="0" width="800px" height="600px"></iframe> -->
                <div v-for="i in numPages" :key="i">
                    <pdf :src="src" :page="i"></pdf>
                </div>
            </div>
            <viewer class="imgBox" :images="inmages" v-else v-loading="loadingimg" element-loading-text="文件加载中">
                <!-- // photo 一定要一个数组，否则报错 -->
                <img v-for="(src, index) in inmages" :src="src.filePath" :key="index">
            </viewer>
        </el-dialog>
        <!-- 结算 -->
        <el-drawer :title="logo" :visible.sync="jsdrawer" size="700px" :show-clos="false">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
                <el-tabs v-model="activeName" @tab-click="JshandleClick">
                    <el-tab-pane label="基本信息" name="first">
                        <el-row>
                            <el-col :span="24">
                                <div class="see-top">基础信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>系统编码：</label>
                                        <span>{{ jsseeData.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>采购合同：</label>
                                        <span>{{ jsseeData.contractName }}</span>
                                    </li>
                                    <li>
                                        <label>供应商：</label>
                                        <span>{{ jsseeData.supplierName }}</span>
                                    </li>
                                    <li>
                                        <label>所属月份：</label>
                                        <span>{{ jsseeData.month }}</span>
                                    </li>
                                    <li>
                                        <label>货物名称：</label>
                                        <span>{{ jsseeData.goodsName }}</span>
                                    </li>
                                    <li>
                                        <label>结算时间：</label>
                                        <span v-if="jsseeData.setStartTime">{{ jsseeData.setStartTime | Time }}</span>
                                        <span> ~ </span>
                                        <span v-if="jsseeData.setEndTime">{{ jsseeData.setEndTime | Time }}</span>
                                    </li>
                                    <li>
                                        <label>结算数量：</label>
                                        <span>{{ jsseeData.settlementCount }} 吨</span>
                                    </li>
                                    <li>
                                        <label>结算单价：</label>
                                        <span>{{ jsseeData.unitPrice }} 元 / 吨</span>
                                    </li>
                                    <li>
                                        <label>结算总价：</label>
                                        <span v-if="jsseeData.sumPrice">{{ jsseeData.sumPrice | currency }} 元</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">附件信息</div>
                                <dl class="see-file">
                                    <dd>
                                        <el-col :span="5">
                                            <span>文件附件：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <span v-if="filelist.length==0"> --- </span>
                                            <a v-else v-for="(v, k) in filelist" :key="k" @mouseenter="changeActive($event)"
                                                @mouseleave="removeActive($event)" @click="fileGo(v, jsseeData)">
                                                <p>{{ v.name }}</p>
                                            </a>
                                        </el-col>
                                    </dd>
                                </dl>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </el-drawer>
        <!-- 结算结束 -->
        <!-- 发票开始 -->
        <el-drawer :title="logo" :visible.sync="fpdrawer" size="700px" :show-clos="false">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
                <el-tabs v-model="activeName">
                    <el-tab-pane label="基本信息" name="first">
                        <el-row>
                            <el-col :span="24">
                                <div class="see-top">基础信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>系统编码：</label>
                                        <span>{{ fpseeData.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>采购合同：</label>
                                        <span>{{ fpseeData.contractCode }}</span>
                                    </li>
                                    <!-- <li>
                                        <label>采购结算：</label>
                                        <span v-if="seeData.relationFlag == true" @click="detailsPurAll" style="color:#cf1724; cursor: pointer;">详情</span>
                                        <span v-else>---</span>
                                    </li> -->
                                    <li>
                                        <label>所属月份：</label>
                                        <span v-if="fpseeData.month">{{ fpseeData.month | timeDate }}</span>
                                    </li>
                                    <li>
                                        <label>进项类型：</label>
                                        <span>{{ fpseeData.inType }}</span>
                                        <!-- <span v-if="seeData.inType == 1">货物进项</span>
                            <span v-if="seeData.inType == 2">运费进项</span>
                            <span v-if="seeData.inType == 3">其他进项</span> -->
                                    </li>
                                    <li>
                                        <label>发票代码：</label>
                                        <span>{{ fpseeData.billCode }}</span>
                                    </li>
                                    <li>
                                        <label>发票号码：</label>
                                        <span>{{ fpseeData.billNumber }}</span>
                                    </li>
                                    <li>
                                        <label>开票日期：</label>
                                        <span>{{ fpseeData.openDate }}</span>
                                    </li>
                                    <li>
                                        <label>销方税号：</label>
                                        <span>{{ fpseeData.dutyParagraph }}</span>
                                    </li>
                                    <li>
                                        <label>销方企业名称：</label>
                                        <span>{{ fpseeData.companyName }}</span>
                                    </li>
                                    <li>
                                        <label>商品名称：</label>
                                        <span v-if="fpseeData.goodsName">{{ fpseeData.goodsName }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>单位：</label>
                                        <span v-if="fpseeData.uit">{{ fpseeData.uit }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>数量：</label>
                                        <span v-if="fpseeData.count">{{ fpseeData.count }} 吨</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>单价：</label>
                                        <span v-if="fpseeData.unitPrice">{{ fpseeData.unitPrice | currency }} 元 / 吨</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>金额：</label>
                                        <span v-if="fpseeData.baseMoney">{{ fpseeData.baseMoney | currency }} 元</span>
                                    </li>
                                    <li>
                                        <label>税率：</label>
                                        <span v-if="fpseeData.taxRate">{{ fpseeData.taxRate }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>税额：</label>
                                        <span v-if="fpseeData.taxAmount">{{ fpseeData.taxAmount | currency }} 元</span>
                                    </li>
                                    <li>
                                        <label>含税金额：</label>
                                        <span v-if="fpseeData.includMoney">{{ fpseeData.includMoney | currency }} 元</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">附件信息</div>
                                <dl class="see-file">
                                    <dd>
                                        <el-col :span="5">
                                            <span>文件附件：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <span v-if="filelist.length==0"> --- </span>
                                            <a v-else v-for="(v, k) in filelist" :key="k" @mouseenter="changeActive($event)"
                                                @mouseleave="removeActive($event)" @click="fileGo(v, fpseeData)">
                                                <p>{{ v.name }}</p>
                                            </a>
                                        </el-col>
                                    </dd>
                                </dl>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </el-drawer>
        <!-- 发票结束 -->
        <!-- 付款单开始 -->
        <el-drawer :title="logo" :visible.sync="fkdrawer" size="700px" :show-clos="false">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
                <el-tabs v-model="activeName">
                    <el-tab-pane label="基本信息" name="first">
                        <el-row>
                            <el-col :span="24">
                                <div class="see-top">基础信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>系统编码：</label>
                                        <span>{{ fkseeData.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>账户名称：</label>
                                        <span>{{ fkseeData.accountName }}</span>
                                    </li>
                                    <li>
                                        <label>账户号码：</label>
                                        <span>{{ fkseeData.account }}</span>
                                    </li>
                                    <li>
                                        <label>交易日期：</label>
                                        <span>{{ fkseeData.transTime }}</span>
                                    </li>
                                    <li>
                                        <label>金额：</label>
                                        <span v-if="fkseeData.money">{{ fkseeData.money | currency }} 元</span>
                                    </li>
                                    <li>
                                        <label>对方户名：</label>
                                        <span>{{ fkseeData.oppositeName }}</span>
                                    </li>
                                    <li>
                                        <label>对方账号：</label>
                                        <span>{{ fkseeData.oppositeAccount }}</span>
                                    </li>
                                    <li>
                                        <label>摘要：</label>
                                        <span>{{ fkseeData.absContent }}</span>
                                    </li>
                                    <li>
                                        <label>用途：</label>
                                        <span v-if="fkseeData.purpose">{{ fkseeData.purpose }}</span>
                                        <span v-else>---</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">附件信息</div>
                                <dl class="see-file">
                                    <dd>
                                        <el-col :span="5">
                                            <span>文件附件：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <span v-if="filelist.length==0"> --- </span>
                                            <a v-else v-for="(v, k) in filelist" :key="k" @mouseenter="changeActive($event)"
                                                @mouseleave="removeActive($event)" @click="fileGo(v, fkseeData)">
                                                <p>{{ v.name }}</p>
                                            </a>
                                        </el-col>
                                    </dd>
                                </dl>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </el-drawer>
        <!-- 付款单结束 -->
        <!--物流开始 -->
        <el-drawer :title="logo" :visible.sync="wldrawer" size="700px" :show-clos="false">
            <div class="pur-tiao"></div>
            <div class="pur-drawerNav">
                <el-tabs v-model="activeName">
                    <el-tab-pane label="基本信息" name="first">
                        <el-row>
                            <el-col :span="24" v-if="wlseeData">
                                <div class="see-top">基础信息</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>系统编码：</label>
                                        <span>{{ wlseeData.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>运输模式：</label>
                                        <span>{{ wlseeData.transMethod | transportmode }}</span>
                                    </li>
                                    <li>
                                        <label>承运单位：</label>
                                        <span v-if="wlseeData.carriageName == 1">物流公司</span>
                                        <span v-if="wlseeData.carriageName == 2">网络货运</span>
                                    </li>
                                    <li>
                                        <label>承运平台：</label>
                                        <span v-if="wlseeData.carriagePlatform">{{ wlseeData.carriagePlatform |
                                            logisticsplatform }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>物流单位：</label>
                                        <span v-if="wlseeData.logisticsCompany">{{ wlseeData.logisticsCompany }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>承运工具：</label>
                                        <span>{{ wlseeData.transType | carrier }}</span>
                                    </li>
                                    <li>
                                        <label>标的物：</label>
                                        <span v-if="wlseeData.subjectMatter">{{ wlseeData.subjectMatter }}</span>
                                        <span v-else>---</span>
                                    </li>
                                    <li>
                                        <label>装货地点：</label>
                                        <span>{{ wlseeData.loadingAddress }}</span>
                                    </li>
                                    <li>
                                        <label>卸货地点：</label>
                                        <span>{{ wlseeData.unloadingPoint }}</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">运单统计</div>
                                <ul class="see-lable">
                                    <li>
                                        <label>车数：</label>
                                        <span>{{ wlseeData.vehicles }} 车</span>
                                    </li>
                                    <li>
                                        <label>运费金额合计：</label>
                                        <span v-if="wlseeData.freightAmount">{{ wlseeData.freightAmount | currency }}
                                            元</span>
                                        <span
                                            v-if="wlseeData.freightAmount == 0 || wlseeData.freightAmount == ''">---</span>
                                    </li>
                                    <li>
                                        <label>吨数合计：</label>
                                        <span>{{ wlseeData.subjectCount }} 吨</span>
                                    </li>

                                    <li>
                                        <label>原发吨数合计：</label>
                                        <span>{{ wlseeData.originalTonnage }} 吨</span>
                                    </li>
                                    <li>
                                        <label>实收吨数合计：</label>
                                        <span>{{ wlseeData.tonnagePaid }} 吨</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">合同关联信息</div>
                                <ul class="see-lable" v-for="(v, k) in tags" :key="k">
                                    <li>
                                        <label>合同编号：</label>
                                        <span>{{ v.systemCode }}</span>
                                    </li>
                                    <li>
                                        <label>销售单位：</label>
                                        <span>{{ v.salesUnit }}</span>
                                    </li>
                                    <li>
                                        <label>购货单位：</label>
                                        <span>{{ v.buyerUnit }}</span>
                                    </li>
                                    <li>
                                        <label>发货地：</label>
                                        <span>{{ v.deliveryPlace }}</span>
                                    </li>
                                </ul>
                            </el-col>
                            <el-col :span="24">
                                <div class="see-top">附件信息</div>
                                <dl class="see-file">
                                    <dd>
                                        <el-col :span="5">
                                            <span>文件附件：</span>
                                        </el-col>
                                        <el-col :span="19">
                                            <span v-if="filelist.length==0"> --- </span>
                                            <a v-else v-for="(v, k) in filelist" :key="k" @mouseenter="changeActive($event)"
                                                @mouseleave="removeActive($event)" @click="fileGo(v, wlseeData)">
                                                <p>{{ v.name }}</p>
                                            </a>
                                        </el-col>
                                    </dd>
                                </dl>
                            </el-col>
                        </el-row>
                    </el-tab-pane>
                </el-tabs>
            </div>
        </el-drawer>
        <!-- 物流结束 -->
        <!-- 运单弹窗开始 -->
        <el-dialog :visible.sync="dialogSelectUnit" width="60%" :before-close="BuyerhandleClose">
            <template slot="title">
                <span>网货运单</span>
            </template>
            <div class="Wdialog">
                <el-table border ref="multipleTable" :data="CompanytableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }" style="width: 100%;">
                    <el-table-column prop="billNumber" label="运单编号" sortable :show-overflow-tooltip="true"
                        width="200">
                    </el-table-column>
                    <el-table-column prop="carNumber" label="车号" :show-overflow-tooltip="true" width="100">
                            <template slot-scope="scope">
                                <span v-if="scope.row.carNumber">{{ scope.row.carNumber}}</span>
                                <span v-else>---</span>
                            </template>
                    </el-table-column>
                    <el-table-column prop="loadingAddress" label="装货地" :show-overflow-tooltip="true"
                        width="260">
                    </el-table-column>
                    <el-table-column prop="unloadingPoint" label="卸货地" :show-overflow-tooltip="true"
                        width="260">
                    </el-table-column>
                    <el-table-column prop="subjectMatter" label="标的物" :show-overflow-tooltip="true" width="100">
                    </el-table-column>
                    <el-table-column prop="freightUnit" label="运费单价" :show-overflow-tooltip="true" width="100">
                    </el-table-column>
                    <el-table-column label="运费金额" :show-overflow-tooltip="true" width="120"
                        prop="freightAmount">
                        <template slot-scope="scope">
                            <span v-if="scope.row.freightAmount">{{ scope.row.freightAmount | currency }} 元</span>
                            <span v-else>0.00 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="可开票吨数" :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="scope.row.invoicingTonnage">{{ scope.row.invoicingTonnage | currency }}
                                吨</span>
                            <span v-else>0.00 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="原发吨数" :show-overflow-tooltip="true" prop="originalTonnage">
                        <template slot-scope="scope">
                            <span v-if="scope.row.originalTonnage">{{ scope.row.originalTonnage }} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="实收吨数" :show-overflow-tooltip="true" prop="tonnagePaid">
                        <template slot-scope="scope">
                            <span v-if="scope.row.tonnagePaid">{{ scope.row.tonnagePaid }} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="应付金额" :show-overflow-tooltip="true" width="120" prop="payMoney">
                        <template slot-scope="scope">
                            <span v-if="scope.row.payMoney">{{ scope.row.payMoney | currency }} 元</span>
                            <span v-else>0.00 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="承运平台" :show-overflow-tooltip="true" width="120"
                        prop="carriagePlatform">
                        <template slot-scope="scope">
                            <span
                                v-if="scope.row.carriagePlatform">{{ scope.row.carriagePlatform | logisticsplatform }}</span>
                            <span v-else>---</span>
                        </template>
                    </el-table-column>
                    <!-- <el-table-column label="时间" :show-overflow-tooltip="true" width="180" prop="transTime">
                    </el-table-column>
                    <el-table-column prop="loadingLng" label="装货地经纬度" :show-overflow-tooltip="true" width="300">
                    </el-table-column>
                    <el-table-column prop="unloadingLng" label="卸货地经纬度" :show-overflow-tooltip="true"
                        width="300">
                    </el-table-column> -->
                    <el-table-column label="可开票金额" :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="scope.row.invoicingAmount">{{ scope.row.invoicingAmount | currency }}
                                元</span>
                            <span v-else>0.00 元</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="开票状态" :show-overflow-tooltip="true" width="120">
                        <template slot-scope="scope">
                            <span v-if="scope.row.invoicingStatus == 1">待开票</span>
                            <span v-if="scope.row.invoicingStatus == 2">已开票</span>
                            <span v-if="scope.row.invoicingStatus == 3">发票作废</span>
                            <span v-if="scope.row.invoicingStatus == 4">发票冲红</span>
                            <span v-else>---</span>
                        </template>
                    </el-table-column>



                </el-table>
                <div class="purpages">
                    <el-pagination v-show="paginationshowtwo" background :page-size.sync="sizetwo"
                        :current-page.sync="currenttwo" layout="total, sizes, prev, pager, next, jumper" :total="totaltwo"
                        @current-change='allsuppliertwo' @size-change='allsuppliertwo'>
                    </el-pagination>
                </div>
            </div>
        </el-dialog>
        <!--运单弹窗结束 -->
        <!-- 磅单弹窗开始 -->
        <el-dialog :visible.sync="bddialogSelectUnit" width="60%" :before-close="bdBuyerhandleClose">
            <template slot="title">
                <span>磅单详情</span>
            </template>
            <div class="Wdialog">
                <el-table border ref="multipleTable" :data="bdCompanytableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }" style="width: 100%;">
                    <el-table-column prop="licensePlate" label="车牌号">
                    </el-table-column>
                    <el-table-column prop="coalType" label="煤种">
                    </el-table-column>
                    <el-table-column label="皮重">
                        <template slot-scope="scope">
                            <span>{{ scope.row.tare }} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="毛重">
                        <template slot-scope="scope">
                            <span>{{ scope.row.roughWeight }} 吨</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="净重">
                        <template slot-scope="scope">
                            <span>{{ scope.row.netWeight }} 吨</span>
                        </template>
                    </el-table-column>
                </el-table>
                <div class="purpages">
                    <el-pagination  v-show="bdpaginationshowtwo" background :page-size.sync="bdsizetwo"
                        :current-page.sync="bdcurrenttwo" layout="total, sizes, prev, pager, next, jumper" :total="bdtotaltwo"
                        @current-change='bdallsuppliertwo' @size-change='bdallsuppliertwo'>
                    </el-pagination>
                </div>
            </div>
        </el-dialog>
        <!--  磅单弹窗结束-->
        <!-- 内容结束 -->
    </div>
</template>
<script>
import pdf from 'vue-pdf';
import CMapReaderFactory from 'vue-pdf/src/CMapReaderFactory.js'; // 加载中文的包
import zst from '../../assets/zst.svg';
import slhy from '../../assets/slhy.svg';
export default {
    components: { pdf },
    data() {
        return {
            /* 运单弹框 */
            Did: null,
            totaltwo: 0,  // 总条数(网络货运)
            currenttwo: 1,  // 当前页数(网络货运)
            sizetwo: 10,  // 显示条数(网络货运)
            paginationshowtwo: true, // 无数据隐藏分页(网络货运)
            CompanytableData: [],
            dialogSelectUnit: false,
            /* 磅单弹框 */
            pronId: null,
            bdtotaltwo: 0,  // 总条数(磅单)
            bdcurrenttwo: 1,  // 当前页数(磅单)
            bdsizetwo: 10,  // 显示条数(磅单)
            bdpaginationshowtwo: true, // 无数据隐藏分页(磅单)
            bddialogSelectUnit: false,  // 磅单弹窗
            bdCompanytableData: [],
            numPages: '',
            inmages: [],
            loadingimg: true,
            loadingimgss: true,
            itemshow: null,
            // 合同文件预览功能
            ImgVisible: false,
            slhyimg: slhy,
            zstimg: zst,
            searsData: {},           // 展示图数据
            header: {},                          // 上传token
            enclosurefileList: [],          // 协议文件列表
            userIds: null,// 是否开启在线支付
            seeData: {},// 查看数据
            filelist: [],// 文件列表 
            drawer: false,// 查看抽屉
            activeName: "first",    // 查看tab栏
            SelectValue: 0,  // 合同类型
            SelectOptions:[
                { 
                    value: 0,
                    label: '全部合同'
                },
                {
                    value: 1,
                    label: '采购合同'
                }, 
                {
                    value: 2,
                    label: '销售合同'
                },
            ],
            customer: '供应商 / 客户',
            Seller: '销 / 购',
            sel: 1,   // 选项卡默认显示tab
            tabs: [
                { label: "结算", id: 1, width: "width:32px;", num: null },
                { label: "发票", id: 2, width: "width:32px;", num: null },
                { label: "流水", id: 3, width: "width:32px;", num: null },
                { label: "物流", id: 4, width: "width:32px;", num: null },
            ],
            http: '',
            show: false,
            purTitle: "",                                   // 标题
            logo: "",                                       // 抽屉标题

            /* 
            合同信息数据列表
            */
            settlementloading: false,
            settlementcurrent: 1,                            // 当前页数
            settlementsize: 50,                              // 显示条数 
            settlementtotal: 0,                              // 总条数
            keyupTime: '',                                   // 时间搜索
            jshxpurSearch: "",                               // 全部搜索
            settlementtableData: [],                        // 表格数据
            multiplesettlement: null,                          // 选中单条数据

            /* 
            结算数据列表
            */
            // banks:'',                           // 银行搜索    
            purSearch: '',                       // 全部搜索
            SearchoppositeName: '',              // 供应商/客户搜索
            TimeArray: [],                       // 结算时间搜索
            contmoeny: 0,                        // 结算金额合计
            watertcountnum: 0,                   // 结算数量合计  
            loading: false,                      // 加载状态
            current: 1,                          // 当前页数
            size: 50,                            // 显示条数
            total: 0,                            // 总条数
            tableData: [],                       // 表格数据

            /* 
            发票数据列表
            */
            waterpurSearch: '',                       // 全部搜索
            waterkeyupTime: '',                      // 时间搜索
            waterSearchoppositeName: '',              // 对方户名搜索
            watertcontmoeny: 0,                       // 合计金额
            watertcount: 0,                           // 合计数量  
            waterloading: false,                      // 加载状态
            watertcurrent: 1,                         // 当前页数
            watertsize: 50,                           // 显示条数
            waterttotal: 0,                           // 总条数
            watertableData: [],                      // 表格数据

            /* 
            流水数据列表
            */
            newwaterpurSearch: '',                       // 全部搜索
            lsbanks: '',
            optionsBanks: [],                            // 银行搜索
            newwaterkeyupTime: '',                      // 时间搜索
            newwaterSearchoppositeName: '',              // 对方账户搜索
            newwatertcontmoeny: 0,                       // 已核销金额  
            newwatertcount: 0,                           // 合计金额  
            newwaterloading: false,                      // 加载状态
            newwatertcurrent: 1,                         // 当前页数
            newwatertsize: 50,                           // 显示条数
            newwaterttotal: 0,                           // 总条数
            newwatertableData: [],                      // 表格数据

            /* 
            物流数据列表
            */
            TonnagnewwaterpurSearch: '',                 // 全部搜索
            lsbanksTonnage: '',                          // 卸货地搜索              
            TonnagekeyupTime: '',                       // 时间搜索
            TonnageSearchopposite: '',                   // 装货地搜搜搜索
            Tonnage: 0,                                  // 实收吨数合计 
            Tonnagnewwaterloading: false,                // 加载状态
            Tonnagnewwatertcurrent: 1,                   // 当前页数
            Tonnagnewwatertsize: 50,                     // 显示条数
            Tonnagnewwaterttotal: 0,                     // 总条数
            TonnagwatertableData: [],                   // 表格数据  
            // /结算/ 
            jsseeId: "",
            jsactiveName: 'first',
            jsdrawer: false,
            dID: '',                             // 表格单条数据id
            jsseeData: {},                         // 查看数据
            //发票
            fpseeId: "",
            fpactiveName: 'first',
            fpdrawer: false,
            fpseeData: {},
            //付款单
            fkseeId: "",
            fkactiveName: 'first',
            fkdrawer: false,
            fkseeData: {},
            //订单
            wlseeId: "",
            wlactiveName: 'first',
            wldrawer: false,
            wlseeData: {},
            tags: [],
        };
    },
    created() {
        if (this.$route.params.paths == 'Relation') {
            sessionStorage.setItem('path', this.$route.params.paths)
            sessionStorage.setItem('pathID', JSON.stringify(this.$route.params.id))
            this.$router.go(0);
        }
        if (sessionStorage.getItem('path') == 'Relation') {
            // console.log('从四六和一进来')
            this.settlementloading = true;
            // 传参
            var idspath = [];
            if (JSON.parse(sessionStorage.getItem('pathID')).length == 0) {
                idspath.push(0);
            } else {
                idspath = JSON.parse(sessionStorage.getItem('pathID'))
            }
            var listByPageData = {
                fileStatus: 1,
                ids: idspath,
                searchStr: this.jshxpurSearch,
                createTime: this.keyupTime,
                // contractType:"",
                page: {
                    current: this.settlementcurrent,
                    size: this.settlementsize
                }
            }
            // 渲染表格
            this.api.listByPage.all(listByPageData)
            .then(res => {
                this.settlementloading = false;
                if (res.data.code == 200) {
                    this.settlementtableData = res.data.data.records;
                    // console.log('展示图数据'+res.data.data.records)
                    this.multiplesettlement = res.data.data.records[0];
                    this.SelectValue = this.multiplesettlement.contractType*1;
                    this.settlementtotal = res.data.data.total;
                    sessionStorage.setItem('path', '');
                    sessionStorage.setItem('pathID', '');
                }
            })
        } else {
            // console.log('正常进来')
            this.settlementgetpurcon();                             // 合同信息
        }
    },
    mounted() {
        this.tblehei = "610";
        this.purTitle = this.$route.meta.title;                 // 标题
        this.logo = sessionStorage.getItem('companyName')       // 企业名称
        this.http = this.$store.state.http;                    // IP地址
        // 获取银行搜索列表
        this.api.chainCap.banks()
            .then(res => {
                if (res.data.code == 200) {
                    this.optionsBanks = res.data.data;
                }
            })
    },
    methods: {
        //清空合同
        // ContractClear() {
        //     this.SelectValue = '0'
        //     this.tableData = [];
        //     this.tabs = [
        //         { label: "结算", id: 1, width: "width:32px;", num: null },
        //         { label: "发票", id: 2, width: "width:32px;", num: null },
        //         { label: "流水", id: 3, width: "width:32px;", num: null },
        //         { label: "物流", id: 4, width: "width:32px;", num: null },
        //     ]
        //     this.getpurcon();
        // },
        //选择合同类型
        Contractchange(val) {
            this.SelectValue = val;
            this.sel = 1;
            // console.log(this.multiplesettlement)
            // if(val==0){
            //     this.multiplesettlement.contractType = '';
            // }else{
            //     this.multiplesettlement.contractType = val
            // }
            // 合同清空
            this.settlementtableData = [];
            this.settlementgetpurcon();
            // 结算清空
            this.tableData = [];
            // this.getpurcon();
            // // 发票清空
            // this.watertableData = [];
            // // 流水清空
            // this.newwatertableData = [];  
            // // 物流清空
            // this.TonnagwatertableData = [];
            // tab栏还原
            this.tabs = [
                { label: "结算", id: 1, width: "width:32px;", num: null },
                { label: "发票", id: 2, width: "width:32px;", num: null },
                { label: "流水", id: 3, width: "width:32px;", num: null },
                { label: "物流", id: 4, width: "width:32px;", num: null },
            ]
            // 结算搜索框清空
            // this.TimeArray = null;
            // this.purSearch = "";
            // this.SearchoppositeName = "";
            // // 发票搜索框清空
            // // this.waterBanks = "";
            // this.waterpurSearch = "";
            // this.waterkeyupTime = "";
            // this.waterSearchoppositeName = "";
            // // 流水搜索框清空
            // this.lsbanks = "";
            // this.newwaterkeyupTime = "";
            // this.newwaterpurSearch = "";
            // this.newwaterSearchoppositeName = "";
            // // 物流搜索框清空
            // this.lsbanksTonnage = "";
            // this.TonnagekeyupTime = "";
            // this.TonnageSearchopposite = "";
            // this.TonnagnewwaterpurSearch = "";
            
        },
        // 合同信息选中条数
        settlementhandleCurrentChange(currentRow) {
            this.multiplesettlement = currentRow;
            // console.log('合同被选中' + currentRow.contractType)
            if (currentRow.contractType*1 == 1) {
                this.Cgconart(currentRow.id, currentRow.orderId)
                this.tabs = [
                    { label: "采购结算 ", id: 1, width: "width:64px;", num: null },
                    { label: "进项发票", id: 2, width: "width:64px;", num: null },
                    { label: "付款单", id: 3, width: "width:48px;", num: null },
                    { label: "物流订单", id: 4, width: "width:32px;", num: null },
                ]
                this.customer = "供应商";
                this.Seller = "销";
            } else if (currentRow.contractType*1 == 2) {
                this.XsContract(currentRow.id, currentRow.orderId)
                this.tabs = [
                    { label: "销售结算", id: 1, width: "width:64px;", num: null },
                    { label: "销项发票", id: 2, width: "width:64px;", num: null },
                    { label: "收款单", id: 3, width: "width:48px;", num: null },
                    { label: "物流订单", id: 4, width: "width:32px;", num: null },
                ]
                this.customer = "客户";
                this.Seller = "购"
            }
            this.sel = 1;
            this.getpurcon();                           // 结算列表
        },
        changeActive(e) {
            e.currentTarget.className = 'active';
        },
        //采购合同数据
        Cgconart(id, rId) {
            this.api.cgContract.sear({ contractId: id, orderId: rId })
                .then(res => {
                    let { code, data } = res.data;
                    if (code != 200) return
                    this.tabs[0].num = data.settNum
                    this.tabs[1].num = data.billNum
                    this.tabs[2].num = data.detNum
                    this.tabs[3].num = data.orderNum
                })
        },
        //销售合同数据
        XsContract(id, rId) {
            this.api.xsContract.sear({ contractId: id, orderId: rId })
                .then(res => {
                    let { code, data } = res.data;
                    if (code != 200) return
                    this.tabs[0].num = data.settNum
                    this.tabs[1].num = data.billNum
                    this.tabs[2].num = data.detNum
                    this.tabs[3].num = data.orderNum
                })
        },

        // 序号索引
        indexMethod(index) {
            return index + 1;
        },
        // 结算 进项 流水 物流选项卡
        select(item) {
            // 结算搜索框清空
            this.TimeArray = null;
            this.purSearch = "";
            this.SearchoppositeName = "";
            // 发票搜索框清空
            // this.waterBanks = "";
            this.waterpurSearch = "";
            this.waterkeyupTime = "";
            this.waterSearchoppositeName = "";
            // 流水搜索框清空
            this.lsbanks = "";
            this.newwaterkeyupTime = "";
            this.newwaterpurSearch = "";
            this.newwaterSearchoppositeName = "";
            // 物流搜索框清空
            this.lsbanksTonnage = "";
            this.TonnagekeyupTime = "";
            this.TonnageSearchopposite = "";
            this.TonnagnewwaterpurSearch = "";

            this.sel = item.id;
            if (this.sel == 1 && this.multiplesettlement) {
                if(this.multiplesettlement.systemCode){
                    this.getpurcon();                       // 结算列表
                }
            } else if (this.sel == 2 && this.multiplesettlement) {
                if(this.multiplesettlement.systemCode){
                    this.watertableData = [];
                    this.watertgetpurcon();                 // 票据列表
                }
            } else if (this.sel == 3 && this.multiplesettlement) {
                if(this.multiplesettlement.systemCode){
                    this.newwatertableData = [];
                    this.newwatertgetpurcon()               // 流水列表
                }
            } else if (this.sel == 4 && this.multiplesettlement) {
                if(this.multiplesettlement.systemCode){
                    this.TonnagwatertableData = [];
                    this.Tonnagpurcon()                     // 物流列表
                }

            }
        },

        /*
        合同信息列表展示功能
        */
        // 合同信息表格数据
        settlementgetpurcon() {
            this.settlementloading = true;
            // 传参
            var listByPageData = {
                fileStatus: 1,
                searchStr: this.jshxpurSearch,
                createTime: this.keyupTime,
                contractType: (this.SelectValue*1)==0?'':this.SelectValue,
                page: {
                    current: this.settlementcurrent,
                    size: this.settlementsize
                }
            }
            // console.log(listByPageData)
            // 渲染表格
            this.api.listByPage.all(listByPageData)
                .then(res => {
                    this.settlementloading = false;
                    if (res.data.code == 200) {
                        this.settlementtableData = res.data.data.records;
                        this.settlementtotal = res.data.data.total;
                    }
                })
        },
        // 搜索全部功能
        jshxpurSearchs() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },
        // 日期搜索
        DataTime() {
            this.settlementtableData = [];
            this.settlementgetpurcon();
            this.settlementcurrent = 1;
        },

        /*
        结算列表展示功能
        */
        // 表格数据
        getpurcon() {
            if(!this.multiplesettlement){
                return  this.$message({
                            type: "warning",
                            message: "提示：先选中合同再进行筛选!",
                            duration:2500,  
                        });
            }
            this.loading = true;
            if (this.TimeArray == null) {
                // 传参
                var listByPageData = {
                    contractId: this.multiplesettlement.id,
                    setStartTime: "",
                    setEndTime: "",
                    oppositeName: this.SearchoppositeName,
                    searchStr: this.purSearch,
                    page: {
                        current: this.current,
                        size: this.size
                    }
                }
            } else {
                // 传参
                var listByPageData = {
                    contractId: this.multiplesettlement.id,
                    setStartTime: this.TimeArray[0],
                    setEndTime: this.TimeArray[1],
                    oppositeName: this.SearchoppositeName,
                    searchStr: this.purSearch,
                    page: {
                        current: this.current,
                        size: this.size
                    }
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageData)
                .then(res => {
                    this.loading = false;
                    if (res.data.code == 200) {
                        this.tableData = res.data.data.records;
                        this.total = res.data.data.total;
                    }
                })
            this.api.capital.amount({ contractId: this.multiplesettlement.id })
                .then(res => {
                    this.contmoeny = 0;             // 合计金额
                    this.watertcountnum = 0;        // 合计数量
                    if (res.data.code == 200) {
                        this.contmoeny = res.data.data.sumMoney;
                        this.watertcountnum = res.data.data.settlementCount;
                    }
                })
        },
        // 供应商搜索
        ClickOppositeName() {
            this.tableData = [];
            this.getpurcon();
            this.current = 1;
        },
        // 结算时间搜索
        Time() {
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

        /*
        发票列表展示功能
        */
        // 表格数据
        watertgetpurcon() {
            if(!this.multiplesettlement){
                return  this.$message({
                            type: "warning",
                            message: "提示：先选中合同再进行筛选!",
                            duration:2500,  
                        });
            }
            this.waterloading = true;
            // 传参
            var listByPageData = {
                contractId: this.multiplesettlement.id,
                openDate: this.waterkeyupTime,
                companyName: this.waterSearchoppositeName,
                searchStr: this.waterpurSearch,
                page: {
                    current: this.watertcurrent,
                    size: this.watertsize
                }
            }
            // 渲染表格
            this.api.bill.bill(listByPageData)
                .then(res => {
                    this.waterloading = false;
                    if (res.data.code == 200) {
                        this.watertableData = res.data.data.records;
                        this.waterttotal = res.data.data.total;
                    }
                })
            this.api.bill.amount({ contractId: this.multiplesettlement.id })
                .then(res => {
                    this.watertcontmoeny = 0;   // 合计含税金额
                    this.watertcount = 0;       // 合计数量
                    if (res.data.code == 200) {
                        this.watertcontmoeny = res.data.data.sumMoney;
                        this.watertcount = res.data.data.count;
                    }
                })
        },
        // 日期搜索
        waterDataTime() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // 企业名称搜索
        waterClickOppositeName() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },
        // 全部搜索
        watersearch() {
            this.watertableData = [];
            this.watertgetpurcon();
            this.watertcurrent = 1;
        },

        /*
        流水列表展示功能
        */
        // 表格数据
        newwatertgetpurcon() {
            if(!this.multiplesettlement){
                return  this.$message({
                            type: "warning",
                            message: "提示：先选中合同再进行筛选!",
                            duration:2500,  
                        });
            }
            this.newwaterloading = true;
            // 传参 
            var listByPageData = {
                contractId: this.multiplesettlement.id,
                bankName: this.lsbanks,
                transTime: this.newwaterkeyupTime,
                oppositeName: this.newwaterSearchoppositeName,
                searchStr: this.newwaterpurSearch,
                page: {
                    current: this.newwatertcurrent,
                    size: this.newwatertsize
                }
            }
            // 渲染表格
            this.api.chainCap.getDetails(listByPageData)
                .then(res => {
                    this.newwaterloading = false;
                    if (res.data.code == 200) {
                        this.newwatertableData = res.data.data.records;
                        this.newwaterttotal = res.data.data.total;
                    }
                })
            this.api.chainCap.amount({ contractId: this.multiplesettlement.id })
                .then(res => {
                    this.newwatertcontmoeny = 0;   // 已核销金额
                    this.newwatertcount = 0;       // 金额
                    if (res.data.code == 200) {
                        this.newwatertcontmoeny = res.data.data.sumYhxje;
                        this.newwatertcount = res.data.data.sumMoney;
                    }
                })
        },
        // 银行名称搜索
        keyupBanks() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },
        // 日期搜索
        newwaterDataTime() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },
        // 企业名称搜索
        newwaterClickOppositeName() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },
        // 全部搜索
        newwatersearch() {
            this.newwatertableData = [];
            this.newwatertgetpurcon();
            this.newwatertcurrent = 1;
        },

        /*
        物流列表展示功能
        */
        // 表格数据
        Tonnagpurcon() {
            if(!this.multiplesettlement){
                return  this.$message({
                            type: "warning",
                            message: "提示：先选中合同再进行筛选!",
                            duration:2500,  
                        });
            }
            this.Tonnagnewwaterloading = true;
            // 传参 
            var listByPageData = {
                contractId: this.multiplesettlement.id,
                unloadingPoint: this.lsbanksTonnage,
                transTime: this.TonnagekeyupTime,
                loadingAddress: this.TonnageSearchopposite,
                searchStr: this.TonnagnewwaterpurSearch,
                page: {
                    current: this.Tonnagnewwatertcurrent,
                    size: this.Tonnagnewwatertsize
                }
            }
            this.api.logistics.selOrders(listByPageData)
            .then(res => {
                this.Tonnage = 0;   // 实收吨数合计
                this.Tonnagnewwaterloading = false;
                if (res.data.code == 200) {
                    this.TonnagwatertableData = res.data.data.records;
                    this.Tonnagnewwaterttotal = res.data.data.total;
                    if(res.data.data.records.length > 0){
                        this.Tonnage = res.data.data.records[0].tonnagePaid;               // 实收吨数合计
                    }
                }
            })
        },

        // 卸货地搜索
        keyupTonnage() {
            this.TonnagwatertableData = [];
            this.Tonnagpurcon();
            this.Tonnagnewwatertcurrent = 1;
        },
        // 日期搜索
        TonnageDataTime() {
            this.TonnagwatertableData = [];
            this.Tonnagpurcon();
            this.Tonnagnewwatertcurrent = 1;
        },
        // 装货地搜索
        TonnageOppositeName() {
            this.TonnagwatertableData = [];
            this.Tonnagpurcon();
            this.Tonnagnewwatertcurrent = 1;
        },
        // 全部搜索
        AllTonnageSearch() {
            this.TonnagwatertableData = [];
            this.Tonnagpurcon();
            this.Tonnagnewwatertcurrent = 1;
        },
        //查看合同
        purSee(row) {
            this.multiplesettlement = row
            this.drawer = true;
            this.activeName = 'first';
            this.seeId = row.id;
            if (row.orderId == null) {
                this.orderId = 0;
            } else {
                this.orderId = row.orderId;
            }
            this.allone(row.id)
        },
        allone(id) {
            this.seeData = {};
            this.filelist = [];
            this.api.cgContract.one({ id: id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.seeData = res.data.data;
                    }
                })
            // 文件
            this.api.file.one({ relationId: id })
                .then(res => {
                    var newarr = [];
                    var newobj;
                    for (let i = 0; i < res.data.data.length; i++) {
                        newobj = {
                            name: res.data.data[i].name,
                            response: {
                                code: 200,
                                data: [res.data.data[i]]
                            }
                        }
                        newarr.push(newobj);
                    }

                    for (let i = 0; i < newarr.length; i++) {
                        if (newarr[i].response.data[0].fileType == 1) {
                            // 附件文件内容
                            this.filelist.push(newarr[i]);
                        } else if (newarr[i].response.data[0].fileType == 2) {
                            // 补充协议内容
                            this.enclosurefileList.push(newarr[i]);
                        }
                    }
                    sessionStorage.setItem('filelist', JSON.stringify(this.enclosurefileList));
                })
        },
        // 查看抽屉
        handleClick(tab) {
            // 判断tab栏选中状态，请求相应的接口
            if (tab.label == "基本信息") {
                this.seeData = {};
                this.filelist = [];
                this.enclosurefileList = [];
                this.allone(this.seeId, this.seeId);
            }
            if (tab.label == "动态") {
            }
            if (tab.label == "展示图") {
                this.sears(this.seeId, this.orderId);
                this.closedraws();
            }
        },
        // 关闭抽屉回调(保存文件)
        closedrawerno() {
            this.closedraws();
        },
        removeActive(e) {
            e.currentTarget.className = '';
        },
        // 删除附件文件
        fileRemove(file, fileList) {
            this.api.addfile.del({ id: file.response.data[0].id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.$message({
                            type: 'warning',
                            message: '删除文件成功!'
                        });
                        this.enclosurefileList = fileList;
                        sessionStorage.setItem('filelist', JSON.stringify(this.enclosurefileList))
                    }
                })
        },
        // // 点击协议文件列表下载
        openfiles(res) {
            this.fileGo(res)
        },
        HtfileGo(obj) {
            this.oneData = obj;
            this.ImgVisible = true;
            let faliename = this.oneData.response.data[0].fileName;
            this.geshi = faliename.substring(faliename.length - 4);
            let arrs = ['.pdf', '.PDF', 'docx', 'DOCX', '.doc', '.DOC'];
            if (!arrs.includes(this.geshi)) {
                this.itemshow = false;
                this.api.file.preFIle({ file_id: this.oneData.id })
                    .then(res => {
                        if (res.data.code == 200) {
                            this.loadingimg = false;
                            this.inmages = res.data.data.slice(0, 2).slice(0, 1)
                        }
                    })
            } else {
                this.itemshow = true;
                this.src = pdf.createLoadingTask({
                    url: this.http + '/saas-common/upload/previewStream?fileName=' + faliename + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f',     // 正式环境用这个！
                    CMapReaderFactory
                })
                // 让所有页数一次性加载完，否则就只会加载第一页
                this.src.promise.then(pdf => {
                    this.loadingimgss = false;
                    this.numPages = pdf.numPages;
                }).catch(() => { })
            }
        },
        /*
           合同文件预览
          */
        fileGo(obj) {
            this.oneData = obj;
            this.ImgVisible = true;
            let faliename = this.oneData.fileName;
            this.geshi = faliename.substring(faliename.length - 4);
            let arrs = ['.pdf', '.PDF', 'docx', 'DOCX', '.doc', '.DOC'];
            if (!arrs.includes(this.geshi)) {
                this.itemshow = false;
                this.api.file.preFIle({ file_id: this.oneData.id })
                    .then(res => {
                        if (res.data.code == 200) {
                            this.loadingimg = false;
                            this.inmages = res.data.data.slice(0, 2).slice(0, 1)
                        }
                    })
            } else {
                this.itemshow = true;
                this.src = pdf.createLoadingTask({
                    url: this.http + '/saas-common/upload/previewStream?fileName=' + faliename + '&&AuthorizationId=wxda0b1c468b83bc99&&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f',     // 正式环境用这个！
                    CMapReaderFactory
                })
                // 让所有页数一次性加载完，否则就只会加载第一页
                this.src.promise.then(pdf => {
                    this.loadingimgss = false;
                    this.numPages = pdf.numPages;
                }).catch(() => { })
            }
        },
        // 协议文件格式验证
        beforedemoUpload(file) {
            let str = file.name;
            let index = str.lastIndexOf("\.");
            str = str.substring(index + 1, str.length)
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
            if (!status) {
                this.fileObj[file.name] = false;
                this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC JPEG 格式,请重新上传!");
            }
            return status;
        },
        // 协议文件上传文件成功
        resfiles(response, file, fileList) {
            if (response.code == 200) {
                this.$message({
                    type: 'success',
                    message: '上传文件成功!'
                });
                this.enclosurefileList = fileList;
                sessionStorage.setItem('filelist', JSON.stringify(this.enclosurefileList));
            }
        },
        // 跳转四流合一功能    @click="zstJump"
        zstJump() {
            this.$router.push({ name: "Relation", params: { paths: 'Relation', id: this.pathIDS } });
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
        // 查看 编辑保存函数
        closedraws() {
            this.seeData = {};
            this.filelist = [];
            this.api.cgContract.one({ id: this.seeId })
                .then(res => {
                    if (res.data.code == 200) {
                        this.seeData = res.data.data;
                        // 文件
                        this.api.file.one({ relationId: this.seeId })
                            .then(res => {
                                var newarr = [];
                                var newobj;
                                for (let i = 0; i < res.data.data.length; i++) {
                                    newobj = {
                                        name: res.data.data[i].name,
                                        response: {
                                            code: 200,
                                            data: [res.data.data[i]]
                                        }
                                    }
                                    newarr.push(newobj);
                                }
                                for (let i = 0; i < newarr.length; i++) {
                                    if (newarr[i].response.data[0].fileType == 1) {
                                        // 附件文件内容
                                        this.filelist.push(newarr[i]);
                                    } else if (newarr[i].response.data[0].fileType == 2) {
                                        // 补充协议内容
                                        this.enclosurefileList.push(newarr[i]);
                                    }
                                }
                                sessionStorage.setItem('filelist', JSON.stringify(this.enclosurefileList));
                                let arr = [];
                                let Array = null;
                                let cid = null;
                                let contract = this.filelist.length != 0 ? (this.filelist[0].response.data[0].id).toString() : '';
                                if (sessionStorage.getItem('filelist')) {
                                    Array = JSON.parse(sessionStorage.getItem('filelist'))
                                } else {
                                    Array = this.enclosurefileList;
                                }
                                cid = Array.map((itme) => {
                                    arr.push({ id: itme.response.data[0].id });
                                    return itme.response.data[0].id;
                                }).join(',');
                                this.seeData.filelist = arr;
                                this.seeData.agreementIds = cid;
                                this.seeData.contractIds = contract;
                                this.seeData.eqbFlow = null;
                                if (this.seeData.createUserId == this.userIds) {
                                    this.api.cgContract.edit(this.seeData)
                                        .then(res => {
                                            if (res.data.code == 200) {
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
        //结算抽屉
        JspurSee(index, row) {
            this.jsdrawer = true;
            this.jsactiveName = 'first';
            this.jsseeId = row.id;
            this.Jsallone(row.id)
        },
        //
        Jsallone(id) {
            this.api.cgcapital.one({ id: id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.jsseeData = res.data.data;
                        this.dID = res.data.data.id;
                    }
                })
            // 查看文件
            this.api.file.one({ relationId: id })
                .then(res => {
                    this.filelist = res.data.data
                })
        },
        //发票抽屉
        JshandleClick(tab) {
            if (tab.label == "基本信息") {
                this.allone(this.jsseeId)
            }
        },
        FppurSee(index, row) {
            this.fpdrawer = true;
            this.fpactiveName = 'first';
            this.fpseeId = row.id;
            this.Fpallone(row.id)
        },
        Fpallone(id) {
            this.api.jxbill.one({ id: id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.fpseeData = res.data.data;
                        this.dID = res.data.data.id;
                    }
                })
            // 查看文件
            this.api.file.one({ relationId: id })
                .then(res => {
                    this.filelist = res.data.data
                })
        },
        //付款抽屉
        FkpurSee(index, row) {
            this.fkdrawer = true;
            this.fkactiveName = 'first';
            this.fkseeId = row.id;
            this.Fkallone(row.id)
        },
        Fkallone(id) {
            this.api.cgchainCap.one({ id: id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.fkseeData = res.data.data;
                        this.dID = res.data.data.id;
                    }
                })
            // 查看文件
            this.api.file.one({ relationId: id })
                .then(res => {
                    this.filelist = res.data.data
                })
        },
        WlpurSee(index, row) {
            this.wldrawer = true;
            this.wlactiveName = 'first';
            this.wlseeId = row.id;
            this.api.logistics.one({ id: row.id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.wlseeData = res.data.data;
                    }
                })
            // 关联合同数据集合
            this.api.logistics.oneAll({ id: row.id })
                .then(res => {
                    if (res.data.code == 200) {
                        this.tags = res.data.data;
                    }
                })
            // 文件
            this.api.file.one({ relationId: row.id })
                .then(res => {
                    this.filelist = res.data.data
                })
        },
        /*
        物流--运单弹框
        */ 
        dialogSelect(data) {
            this.Did = data.id;
            this.allsuppliertwo();                              // 渲染表格
            this.dialogSelectUnit = true;                         // 开启弹窗
        },
        BuyerhandleClose() {
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
        },
        allsuppliertwo() {
            // 传参
            var listByPageDatatwo = {
                orderId: this.Did,
                partnerType: 2,
                searchStr: '',
                page: {
                    current: this.currenttwo,
                    size: this.sizetwo
                }
            }
            // 渲染表格
            this.api.waybill.all(listByPageDatatwo)
                .then(res => {
                    if (res.data.code == 200) {
                        if (res.data.data.records != "") {
                            this.CompanytableData = res.data.data.records
                            this.totaltwo = res.data.data.total
                        } else {
                            this.paginationshowtwo = false;
                        }
                    }
                })
        },
        /*
        物流--磅单弹框
        */
        bddialogSelect(data) {
            this.pronId = data.id;
            this.bdallsuppliertwo();                              // 渲染表格
            this.bddialogSelectUnit = true;                       // 开启弹窗
        },
        bdBuyerhandleClose() {
            this.bddialogSelectUnit = false;                      // 关闭整个弹窗
        },
        bdallsuppliertwo() {
            // 传参
            var listByPageData = {
                orderId: this.pronId,
                searchStr: '',
                page: {
                    current: this.bdcurrenttwo,
                    size: this.bdsizetwo
                }
            }
            // 渲染表格
            this.api.Weighing.all(listByPageData)
                .then(res => {
                    if (res.data.code == 200) {
                        if (res.data.data.records != "") {
                            this.bdCompanytableData = res.data.data.records;
                            this.bdtotaltwo = res.data.data.total;
                        } else {
                            this.bdpaginationshowtwo = false;
                        }
                    }
                })
        },
    },
};
</script>
<style lang="less" scoped>
// 磅单 运单详情弹窗样式开始
.Wdialog{
    width: 100%;
    // 磅单分页
    .purpages {
        width: 100%;
        height: 32px;
        margin: 28px 8px 0px 0px;
        .el-pagination {
        float: right;
        }
    }
}
// 磅单 运单详情弹窗样式结束
.imgBox {
    width: 100%;
    height: 600px;
    overflow: hidden;
    overflow-y: scroll;
    position: relative;

    img {
        width: 100%;
    }
}

// 选项卡开始
.totalTab {
    width: auto;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    margin-right: 3%;
    color: #333;
    text-align: center;
}

.totalTab.active {
    // width: 100px;
    color: #ea222e;
}

.stripnone {
    width: auto !important;
    height: 3px;
    background-color: #ea222e;
    margin: 4px auto 0;
    border-radius: 25px;
    display: none;
}

.strip {
    // width: auto;
    display: block;
}

.adBox {
    height: 100%;
    position: relative;
    .Advanced {
        width: 50%;
        height: 48px;
        // overflow: hidden;
        padding: 0 0 16px 32px;
        box-sizing: border-box;
        position: absolute;
        right: 0;
        top: -58px;

        .el-select {
            float: right;
            width: 24%;
            margin-right: 1%;
        }

        .el-input {
            width: 24%;
            float: right;
            margin-right: 1%;
        }
        .el-date-picker {
            width: 24%;
            float: right;
            margin-right: 1%;
        }
    }
    .tab-table {
        height: 88%;
    }


}

//选项卡结束
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
    height: calc(100vh - 94px);
    background-color: #fff;
    margin: 8px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding: 8px;

    .pur-title-header {
        width: 100%;
        // height:67px;
        overflow: hidden;
        padding: 16px 0 24px 32px;
        box-sizing: border-box;
        b {
            color: #333;
            font-size: 16px;
            line-height: 32px;
            float: left;
        }
        >div {
            float: right;
            .el-select {
                float: right;
                width: 40%;
                margin-right: 8px;
            }
            .el-input {
                width: 50%;
                float: right;
                margin-right:8px;
            }
        }
    }

    .rece {
        width: 100%;
        height: 45%;
        border: 1px solid #E4E4E4;
        // float: left;
        box-sizing: border-box;
        position: relative;

        .pur-title {
            width: 100%;
            // height:67px;
            overflow: hidden;
            padding: 16px 0 24px 32px;
            box-sizing: border-box;

            b {
                color: #333;
                font-size: 16px;
            }
        }

        .pur-table {
            width: 100%;
            height:calc(100% - 120px);

            >.el-table {
                height: 100%;
            }
        }

        .pur-pages {
            width: 100%;
            overflow: hidden;
            padding: 10px 18px 0 24px;
            box-sizing: border-box;
            .sjhxjeMoeny {
                float: left;
                // position: absolute;
                // top: -65px;
                // right: 32px;

                span {
                    font-size: 14px;
                    font-weight: 600;
                    color: #333;
                    line-height: 28px;
                }
            }
            .el-pagination {
                float: right;
            }
        }
    }

    .rece-right {
        width: 30px;
        margin: auto;
        line-height: 24px;
        transform: rotate(270deg);
        font-size: 24px;
        color: #999;
        font-weight: 300;
    }

    /deep/.el-table__body tr.current-row>td {
        background-color: #F7EDED !important;
    }
}

.see-top {
    height: 40px;
    font-weight: 500;
    font-size: 14px;
    line-height: 45px;
    border-bottom: 1px dashed #999999;
}

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

        label {
            float: left;
        }

        .see-libox {
            width: 76%;
            display: inline-block;

            span {
                color: #666;
            }

            div {
                display: inline-block;
                margin-right: 15px;
            }
        }
    }

}

// 文件样式
.see-file {
    padding: 15px 0;

    span {
        font-size: 14px;
        color: #666;
    }

    dd {
        display: flow-root;
        font-size: 14px;
        line-height: 34px;
        color: #333;

        p {
            text-decoration: none;
            color: #333;
            cursor: pointer;
        }

        .active {
            color: #EA222E;
        }
    }
}

// 上传补充协议样式
.AddbuttonUpload {
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

.AddbuttonUpload:hover {
    background: #e4e4e4;
    color: #333;
}

// 查看抽屉样式结束
// 内容结束</style>