<template>
    <!-- 新建销售收款单页面 -->
    <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="130px">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="3"><span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="21" class="pur-right">
            <el-button size="small" @click="AddDatasubmitForm('AddData')">保存</el-button>
            <el-button size="small" style="margin-left: 24px"  @click="Addopen">取消</el-button>
        </el-col>
        </el-row>
        <!-- 头部结束 -->
        <!-- 内容开始 -->
        <div class="pur-nav">
            <!-- 表单 -->
            <div class="pur-table">
                <el-row>
                    <!-- 基础信息 -->
                    <el-col :span="24">
                        <div class="see-top">基础信息</div>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">系统编码</label>
                            </el-col>
                            <el-form-item prop="systemCode">
                                <el-col :span="7">
                                    <el-input v-model="AddData.systemCode" placeholder="请输入系统编码" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">销售结算</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <div class="details">
                                        <div class="details-nav">
                                            <div class="details-tag" >
                                                <div class="tag-box" v-for="(tag,k) in tags" :key="k" v-if="tags.length == 1">
                                                    <el-tag style="margin:4px;"
                                                    size="small"
                                                    type="info">
                                                        <span v-if="tag.sumPrice">{{tag.settlementCount}} 吨 - 
                                                            <el-tooltip effect="dark" :content="'结算总价 ： ' + tag.sumPrice + ' 元'" placement="top">
                                                                <span> {{tag.sumPrice|currency}} 元</span>
                                                            </el-tooltip>
                                                        </span>
                                                    </el-tag>
                                                </div>
                                                <div class="tag-box" v-for="(tag,k) in tags" :key="k" v-if="tags.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                    <el-tag size="small" type="info" style="margin:0;">
                                                        <span v-if="tags[0].sumPrice">{{tags[0].settlementCount}} 吨 - {{tags[0].sumPrice|currency}} 元</span>
                                                    </el-tag>
                                                    <b> . . .</b>
                                                </div>
                                            </div>
                                            <div class="details-conent">
                                                共 <span style="color:#CF1724;"> {{tags.length}}</span> 条 , 合计 : 
                                                    <el-tooltip class="item" effect="dark" :content="'合计 : ' + moeny + ' 元'" placement="top">
                                                        <span style="color:#CF1724;">{{moeny|currency}}</span>
                                                    </el-tooltip> 元
                                            </div> 
                                        </div>
                                        <div class="details-size" @click="PlaceChoice">
                                            <span>关 联</span>
                                        </div>
                                    </div>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">银行名称</label>
                            </el-col>
                            <el-form-item prop="bankName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.bankName" placeholder="请输入银行名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">账户名称</label>
                            </el-col>
                            <el-form-item prop="accountName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.accountName" placeholder="请输入账户名称全称" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">账户号码</label>
                            </el-col>
                            <el-form-item prop="account">
                                <el-col :span="7">
                                    <el-input v-model="AddData.account" placeholder="请输入账户号码" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:16px;cursor: pointer;" @click="dialogSelect">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">交易日期</label>
                            </el-col>
                            <el-form-item prop="transTime">
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.transTime"
                                    type="date"
                                    placeholder="请选择交易日期"
                                    value-format="yyyy-MM-dd"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">金额</label>
                            </el-col>
                            <el-form-item prop="money">
                                <el-col :span="7">
                                    <el-input v-model="AddData.money" placeholder="元" @input="CardNumberss" :maxlength="length"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">对方账户</label>
                            </el-col>
                            <el-form-item prop="oppositeName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.oppositeName" placeholder="请输入账户全称" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">对方账号</label>
                            </el-col>
                            <el-form-item prop="oppositeAccount">
                                <el-col :span="7">
                                    <el-input v-model="AddData.oppositeAccount" placeholder="请输入完整账号" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:16px;cursor: pointer;" @click="otherParty">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">摘要</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.absContent" placeholder="请输入摘要内容"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">用途</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.purpose" placeholder="请输入用途内容"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">用途</label>
                            </el-col>
                            <el-form-item prop="purpose">
                                <el-col :span="7">
                                    <el-select filterable v-model="AddData.purpose" placeholder="请选择用途" style="width:100%;">
                                        <el-option
                                        v-for="(item,k) in SelectproSmaType" :key="k"
                                        :label="item"
                                        :value="k"
                                        >
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                    </el-col>
                    <!-- 附件信息 -->
                    <el-col :span="24">
                        <div class="see-top">附件信息</div>
                        <!-- 上传合同附件 -->
                        <el-row>
                            <el-col :span="2" >
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">上传合同附件</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7"> 
                                    <el-upload
                                    :action="http + '/saas-common/upload/batch'"
                                    :file-list="enclosurefileList"
                                    :headers="header"
                                    :on-remove="fileRemove"
                                    :on-preview="openfile"
                                    :before-upload="beforedemoUpload"
                                    :on-success="resfile">
                                        <div class="AddbuttonUpload"><span>请选择附件</span></div>
                                    </el-upload>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row>
            </div>
        </div>
        <!-- 内容结束 -->
        <!-- 账户名称弹窗开始 -->
        <el-dialog
            :visible.sync="dialogSelectUnit"
            width="60%"
            :before-close="BuyerhandleClose">
            <template slot="title">
                <div style="font-size:18px;">账户选择</div>
            </template>
            <!-- 显示账户内容框 -->
            <el-input type="textarea" v-model="dialogtextarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox" style="padding-top:50px;">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="search"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="purSearch"
                        style="width:200px;float:right;"
                        >
                        </el-input>
                        <span @click="Companyshow">+ 新建</span>
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
                        <el-table-column prop="accountName" label="账户名称">
                        </el-table-column>
                        <el-table-column prop="accountNumber" label="账户号码">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-size.sync="contractsize"
                        :current-page.sync="contractcurrent"
                        layout="prev, pager, next"
                        :total="contracttotal"
                        @current-change='contractall'
                        @size-change='contractall'>
                        </el-pagination>
                        <!-- 账户信息应用按钮 -->
                        <el-button size="small" @click="accountNameApplication">应用</el-button>  
                    </div>
                    <!-- 新建账户信息 -->
                    <el-form :model="CompanyNewData" :rules="rels" ref="CompanyNewData" label-width="100px" size="small" v-show="CompanyDatashow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="10">
                                <el-form-item label="账户名称：" prop="accountName">
                                    <el-input v-model="CompanyNewData.accountName" placeholder="请输入账户名称"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="10">
                                <el-form-item label="账户号码：" prop="accountNumber">
                                    <el-input v-model="CompanyNewData.accountNumber" placeholder="请输入账户号码" @input="CardNumber"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="1" style="float:right;margin-right:1%;">
                                <el-button size="small" @click="CompanyNewsubmitForm('CompanyNewData')">保存</el-button>
                            </el-col>
                            <el-col :span="2" style="float:right;">
                                <el-button @click="CompanyNewShow" size="small">取消</el-button>
                            </el-col>

                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 账户名称弹窗结束 -->

        <!-- 对方账户信息弹窗开始 -->
        <el-dialog
            :visible.sync="otherPartyUnit"
            width="60%"
            :before-close="otherPartyhandleClose">
            <template slot="title">
                <div style="font-size:18px;">对方账户</div>
            </template>
            <!-- 显示 对方账户信息内容框 -->
            <el-input type="textarea" v-model="otherPartytarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="otherPartysearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="otherPartypurSearch"
                        style="width:200px;"
                        >
                        </el-input>
                        <span @click="otherPartyCompanyshow">+ 新建</span>
                    </div>
                    <!--  对方账户信息表格 -->
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="otherPartytableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="otherPartyPlaceradio" @change.native="otherPartygetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="accountName" label="对方账户">
                        </el-table-column>
                        <el-table-column prop="accountNumber" label="对方账号">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-size.sync="size"
                        :current-page.sync="current"
                        layout="prev, pager, next"
                        :total="total"
                        @current-change='tractall'
                        @size-change='tractall'>
                        </el-pagination>
                        <!-- 对方账户应用按钮 -->
                        <el-button size="small" @click="oppositeNameApplication">应用</el-button>  
                    </div>
                    <!-- 新建对方账户信息 -->
                    <el-form :model="otherPartyNewData" :rules="Placerels" ref="otherPartyNewData" label-width="80px" size="small" v-show="otherPartyshow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="8">
                                <el-form-item label="账户：" prop="accountName">
                                    <el-input v-model="otherPartyNewData.accountName" placeholder="请输入对方账户" style="width:260px;"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="10">
                                <el-form-item label="账号：" prop="accountNumber">
                                    <el-input v-model="otherPartyNewData.accountNumber" placeholder="请输入对方账号" style="width:260px;" @input="otherCardNumber"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="1" style="margin-right:12px;float:right;">
                                <el-button size="small" @click="otherPartyNewsubmitForm('otherPartyNewData')">保存</el-button>
                            </el-col>
                            <el-col :span="2" style="padding-left:25px;float:right;">
                                <el-button @click="otherPartyNewShow" size="small">取消</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 对方账户信息弹窗弹窗结束 -->

        <!-- 销售结算弹窗开始 -->
        <el-dialog
            :visible.sync="PlaceChoiceUnit"
            width="60%"
            :before-close="PlaceChoicehandleClose">
            <template slot="title">
                <div style="font-size:18px;">销售结算关联</div>
            </template>
            <!-- 显示销售结算选择内容框 -->
            <div class="smbox">
                <el-tag
                size="small"
                v-for="tag in tags"
                :key="tag.systemCode"
                type="info" style="margin:5px 5px;color:#333;">
                <span v-if="tag.sumPrice">{{tag.settlementCount}} 吨 - {{tag.sumPrice|currency}} 元</span>
                </el-tag>
            </div>
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="PlaceChoicesearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="PlaceChoicepurSearch"
                        style="width:200px;float:right;"
                        >
                        </el-input>
                        <!-- <span @click="PlaceChoiceCompanyshow">+ 新建</span> -->
                    </div>
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="PlaceChoicetableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;"
                    @selection-change="relationChange"
                    :summary-method="getSummaries"
                    show-summary
                    :row-key="keyId">
                        <el-table-column type="selection" :reserve-selection="true"> </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="systemCode" label="系统编号" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column label="结算时间" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span>{{scope.row.setStartTime}}</span><span> ~ </span><span>{{scope.row.setEndTime}}</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="结算数量" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span>{{scope.row.settlementCount}} 吨</span>
                            </template>
                        </el-table-column>
                        <el-table-column label="结算总价" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice">{{scope.row.sumPrice|currency}} 元</span>
                            </template>
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-size.sync="PlaceChoicesize"
                        :pager-count="5"
                        :current-page.sync="PlaceChoicecurrent"
                        layout="prev, pager, next"
                        :total="PlaceChoicetotal"
                        @current-change='PlaceChoiceallsupplier'
                        @size-change='PlaceChoiceallsupplier'>
                        </el-pagination>
                        <!-- 销售结算应用按钮 -->
                        <el-button size="small" @click="deliveryPlaceApplication">应用</el-button> 
                    </div>
                </div>
            </div>
        </el-dialog>
        <!-- 销售结算弹窗结束 -->
    </el-form>
</template>
<script>
export default {
    data() {
        return {
            /* 
            销售结算数据列表
            */
            PlaceChoicecurrent:1,                   // 当前页数
            PlaceChoicesize:4,                      // 每页显示条数
            PlaceChoicetotal:0,                     // 总条数
            PlaceChoicepurSearch:'',                // 搜索
            PlaceChoiceUnit:false,                  // 弹窗
            // PlaceChoicetarea:'',                 // 显示销售收款单内容框
            PlaceChoicetableData:[],                // 选择销售收款单表格数据
            tags: [],                               // 收款单回显数据
            moeny:0,                                // 合计金额

            /* 
            账户名称数据列表
            */
            contractcurrent:1,                      // 当前页数
            contractsize:4,                         // 显示条数
            contracttotal:0,                        // 总条数
            dialogSelectUnit:false,                 // 选择账户信息弹窗
            dialogtextarea:'',                      // 显示账户信息内容框
            dialogtext:'',                          // 显示账号信息
            purSearch:'',                           // 选择账户信息搜索框
            CompanytableData:[],                    // 选择账户信息表格数据
            radio:'',                               // 选择账户信息单选框
            show:false,                             // 选择账户信息表格隐藏 id
            CompanyDatashow:false,                  // 新建账户信息 显示隐藏
            // 新建账户表单
            CompanyNewData:{
                accountNumber:'',
                accountName:'',
            },
            // 新建账户表单验证
            rels: {
                accountNumber: [
                    { required: true, message: '请输入账户号码', trigger: 'blur' },
                ],
                accountName: [
                    { required: true, message: '请输入账户名称全称', trigger: 'blur' },
                ],
            },
            
            /* 
            对方账户信息数据列表
            */
            current:1,                              // 当前页数
            size:4,                                 // 显示条数
            total:0,                                // 总条数
            otherPartyUnit:false,                   // 对方账户信息弹窗
            otherPartytarea:'',                     // 对方账户信息内容框
            otherPartytext:'',                      // 对方账号信息内容框
            otherPartypurSearch:'',                 // 对方账户信息搜索框
            otherPartyshow:false,                   // 新建对方账户信息 显示隐藏
            otherPartyPlaceradio:'',                // 对方账户信息单选框
            // 新建对方账户信息表单
            otherPartyNewData:{
                accountName:'',
                accountNumber:''
            },
            // 新建对方账户信息表单验证
            Placerels: {
                accountName: [
                    { required: true, message: '请输入对方账户', trigger: 'blur' },
                ],
                accountNumber: [
                    { required: true, message: '请输入对方账号', trigger: 'blur' },
                ],
            },
            otherPartytableData:[],                 // 选择方账户信息表格数据
            
            /* 
            新建销售收款单数据
            */
            length:'',
            http:'',
            header:{},                      // 上传请求头
            fileObj:{},                     // 附件文件格式验证
            enclosurefileList: [],          // 附件文件列表
            purTitle: "",                   // 标题
            // SelectproSmaType:{},         // 用途
            Extended:'',                    // 扩展信息
            // 表单验证
            rules: {
                systemCode: [
                    { required: true, message: '请输入系统编码', trigger: 'change' }
                ],
                statementCode: [
                    { required: true, message: '请输入销售结算编号', trigger: 'change' }
                ],
                bankName:[
                    { required: true, message: '请输入银行名称', trigger: 'change' }
                ],
                accountName: [
                    { required: true, message: '请输入账户名称全称', trigger: 'change' }
                ],
                account: [
                    { required: true, message: '请输入账户号码', trigger: 'change' }
                ],
                transTime: [
                    { required: true, message: '请选择交易日期', trigger: 'blur' }
                ],
                money: [
                    { required: true, message: '请输入金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                oppositeName: [
                    { required: true, message: '请输入账户全称', trigger: 'change' }
                ],
                oppositeAccount: [
                    { required: true, message: '请输入完整账号', trigger: 'change' }
                ],
                absContent: [
                    { required: true, message: '请输入摘要内容', trigger: 'blur' }
                ],
                purpose: [
                    { required: true, message: '请输入用途', trigger: 'blur' }
                ],
            },
            // 新建表单数据
            AddData: {
                systemCode:'',
                sments:[],
                statementCode:'',
                accountName:'',
                relationFlag:'',
                account:'',
                transTime:'',
                money:'',
                oppositeName:'',
                oppositeAccount:'',
                absContent:'',
                purpose:'',
                filelist:[]  // 合同 协议的集合id 
            }, 
        };
    },
    created() {},
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 渲染系统编码
        this.api.xschainCap.code({codeType:2})
        .then(res=>{
            if(res.data.code == 200){
                this.AddData.systemCode = res.data.data
            }
        })
        //全部维护数据...用途
        // var att = sessionStorage.getItem('data'); 
        // var arr =JSON.parse(att);
        // this.SelectproSmaType = arr.purpose;
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')}
        this.http  = this.$store.state.http;                 // IP地址
    },
    methods: {
        /* 
        新建销售收款单表单提交功能
        */ 
        // 金额输入框限制
        CardNumberss(val) {
            if(val.indexOf('.') != -1){
                this.length = 12;
            }else{
                this.length = 10;
            }
        },
        // 新建销售收款单表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    // 文件上传
                    var arr = [];
                    this.enclosurefileList.map((itme,index)=>{
                        arr.push({id:itme.response.data[0].id});
                    }).join(',')
                    this.AddData.filelist = arr;
                    this.AddData.detype = 1;
                    this.AddData.sments = this.tags;
                    // 销售结算
                    if(this.AddData.sments.length == 0){
                        this.AddData.relationFlag = false;
                    }else{
                        this.AddData.relationFlag = true;
                    }
                    this.api.xschainCap.add(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                    this.delRemovedFiles();
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.$router.push({name:"Receipt"})
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            this.$router.push({name:"Receipt"})
            // this.$confirm("检测到表单已有内容填写，是否保存为草稿?", "提示", {
            //     confirmButtonText: "确定",
            //     cancelButtonText: "取消",
            //     type: "warning",
            // })
            // .then(() => {
            //     var that = this
            //     that.$message({
            //         type: "success",
            //         message: "保存草稿成功!",
            //         duration:1000,  
            //         onClose(){
            //             that.$router.push({name:"Receipt"})
            //         }
            //     });
            // })
            // .catch(() => {
            //     this.$router.push({name:"Receipt"})
            // });
        },
        // 附件文件格式验证
        beforedemoUpload(file) {
            let str = file.name;
            let index = str.lastIndexOf("\.");
            str = str.substring(index + 1,str.length)
            let types = [
                "png", 
                "jpg",
                "pdf",
                "doc",
                "docx",
                "xlsx",
                "xls"
            ];
            const status = types.indexOf(str) != -1;
            this.fileObj[file.name] = true;
            if(!status) {
                this.fileObj[file.name] = false;
                this.$message.error("上传文件只能是 PNG JPG PDF DOCX DOC XLSX XLS 格式,请重新上传!");
            }
            return status;
        },
        // 附件上传文件成功
        resfile(response, file, fileList){
            // HTTP 200 + code≠200 或 data 为空(例如后端转 pdf 失败)时拿不到附件号,文件不能留在列表里
            if(response.code == 200 && response.data && response.data.length){
                this.$message({
                    type: 'success',
                    message: '上传文件成功!'
                });
                this.enclosurefileList = fileList;
            }else{
                fileList.splice(fileList.indexOf(file), 1);
                this.$message.error(response.msg || '上传文件失败!');
            }
        },
        // 点击附件文件列表下载
        openfile(res){
            let url = this.http + '/saas-common/upload/downById?Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99&fileName=' + res.response.data[0].fileName;
            window.open(url)
        },
        // 删除附件文件
        // 保存成功后才删编辑期间移除的回显附件(原先一移除就删,取消编辑后记录指向已删文件)
        delRemovedFiles(){
            (this.removedFileIds || []).forEach(id => this.api.addfile.del({id}));
            this.removedFileIds = [];
        },
        fileRemove(file,fileList){
            this.enclosurefileList = fileList; // 列表与界面一致:原先删除接口失败时提交仍会带上
            if(!this.fileObj[file.name]){
                // 回显的附件没有经过 before-upload:先记下,等保存成功再删
                this.removedFileIds = (this.removedFileIds || []).concat(file.response.data[0].id);
                return;
            }
            return  this.api.addfile.del({id:file.response.data[0].id})
                        .then(res=>{
                            if(res.data.code == 200){
                                this.$message({
                                    type: 'warning',
                                    message: '删除文件成功!'
                                });
                                this.enclosurefileList = fileList;
                            }
                        })
        },

        /* 
        账户名称的选择功能
        */
        // 账户输入框分割符
        CardNumber(){
            this.CompanyNewData.accountNumber = this.CompanyNewData.accountNumber.replace(/\s/g, ' ').replace(/[^\d]/g, ' ').replace(/(\d{4})(?=\d)/g, '$1 ');
        },
        contractall(){
            this.radio = ""; // 翻页后清空单选框:新一页同下标的行会显示成上一页选中的样子
            // 传参
            var listByPageData = {
                ownerType:1,
                searchStr:this.purSearch,
                page:{
                current:this.contractcurrent,
                size:this.contractsize
                }
            }
            // 渲染表格
            this.api.Account.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.CompanytableData = res.data.data.records
                    this.contracttotal = res.data.data.total
                }
            })
        },
        // 账户名称搜索功能
        search() {
            this.contractcurrent = 1; // 原先误写成总条数 contracttotal,且是先请求后改页码:翻页后搜索仍带旧页码
            this.contractall();
        },
        // 账户名称 选择账户表格
        getCurrentRow(row){
            // 获取选中数据
            this.dialogtextarea = row.accountName;
            this.dialogtext = row.accountNumber
        },
        // 隐藏显示新建
        Companyshow(){
            this.CompanyDatashow = true
        },
        // 新建账户信息提交表单
        CompanyNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.CompanyNewData.ownerType = 1;
                    this.api.Account.add(this.CompanyNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.contractall();
                                    that.CompanyDatashow = false;
                                    that.CompanyNewData = {};
                                }
                            });
                        }
                    })
                } else {
                    let that = this
                    that.$message({
                        type: "error",
                        message: "保存失败,请重新选择!", 
                        onClose(){
                            that.CompanyDatashow = false;
                        }
                    });
                    return false;
                }
            });
        },
        // 新建账户信息取消提交表单
        CompanyNewShow() {
            this.CompanyNewData = {};
            this.CompanyDatashow = false
        },
        // 选择账户名称弹窗
        dialogSelect(){
            this.dialogSelectUnit=true;         // 开启弹窗
            this.contractall();
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.dialogSelectUnit == true){
                this.CompanyDatashow = false
            }
        },
        // 选择账户 关闭弹窗
        BuyerhandleClose() {
            this.radio = "";                    // 清空单选框
            this.dialogtextarea = "";           // 清空显示账户信息内容框
            this.dialogSelectUnit = false;      // 关闭整个弹窗
        },
        // 账户信息应用按钮 事件
        accountNameApplication(){
            this.AddData.accountName = this.dialogtextarea;         // 获取账户名称，回显到表单
            this.AddData.account = this.dialogtext;                 // 获取账号名称，回显到表单
            this.dialogSelectUnit = false;                          // 关闭整个弹窗
            this.radio = "";                                        // 清空单选框
            this.dialogtextarea = "";                               // 清空显示账户信息内容框
            this.dialogtext = "";                                   // 清空显示账号信息内容框
        },

        /* 
        对方账户信息的选择功能
        */
        // 对方账户输入框分割符
        otherCardNumber(){
            this.otherPartyNewData.accountNumber = this.otherPartyNewData.accountNumber.replace(/\s/g, ' ').replace(/[^\d]/g, ' ').replace(/(\d{4})(?=\d)/g, '$1 ');
        },
        tractall(){
            // 传参
            var listByPageData = {
                ownerType:2,
                searchStr:this.otherPartypurSearch,
                page:{
                    current:this.current,
                    size:this.size
                }
            }
            // 渲染表格
            this.api.Account.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.otherPartytableData = res.data.data.records
                    this.total = res.data.data.total
                }
            })
        },
        // 对方账户信息搜索功能
        otherPartysearch() {
            this.current = 1; // 原先误写成总条数 total,且是先请求后改页码
            this.tractall();
        },
        // 隐藏显示对方账户信息新建
        otherPartyCompanyshow(){
            this.otherPartyshow = true
        },
        // 对方账户信息选择表格
        otherPartygetCurrentRow(row){
            // 获取选中账户信息数据
            this.otherPartytarea = row.accountName;
            // 获取选中账号信息数据
            this.otherPartytext = row.accountNumber;
        },
        // 对方账户信息提交表单   
        otherPartyNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.otherPartyNewData.ownerType = 2;
                    this.api.Account.add(this.otherPartyNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    that.tractall();
                                    that.otherPartyshow = false;  
                                    that.otherPartyNewData = {};
                                }
                            });
                        }
                    })
                } else {
                    let that = this
                    that.$message({
                        type: "error",
                        message: "保存失败,请重新选择!", 
                        onClose(){
                            that.otherPartyshow = false;
                            that.otherPartyNewData = {};
                        }
                    });
                    return false;
                }
            });
        },
        // 新建对方账户信息取消提交表单
        otherPartyNewShow() {
            this.otherPartyNewData = {};
            this.otherPartyshow = false
        },
        // 选择对方账户信息弹窗
        otherParty(){
            this.otherPartyUnit=true;                               // 开启弹窗
            this.tractall();
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.otherPartyUnit == true){
                this.otherPartyshow = false
            }
        },
        // 对方账户信息关闭弹窗
        otherPartyhandleClose() {
            this.otherPartyPlaceradio = "";                         // 清空单选框
            this.otherPartytarea = "";                              // 清空显示对方账户信息内容框
            this.otherPartyUnit = false;                            // 关闭整个弹窗
        },
        // 对方账户信息应用按钮 事件
        oppositeNameApplication(){
            this.AddData.oppositeName = this.otherPartytarea;       // 获取对方账户信息，回显到表单
            this.AddData.oppositeAccount = this.otherPartytext;     // 获取对方账号信息，回显到表单
            this.otherPartyUnit = false;                            // 关闭整个弹窗
            this.otherPartyPlaceradio = "";                         // 清空单选框
            this.otherPartytarea = "";                              // 清空显示对方账户信息内容框
        },
 
        /* 
        销售结算的选择弹窗功能
        */
       // 合计行
        getSummaries(param) {
            const { columns, data } = param;
            const sums = [];
            columns.forEach((column, index) => {
                if (index === 1) {
                    sums[index] = '合计';
                    return;
                }
                if (index === 4) {
                    if(this.moeny == 0){
                        sums[index] = '0 元';
                    }else{
                        sums[index] = this.moeny + ' 元';
                    }
                    return;
                }
            });
            return sums;
        },
        // 销售结算表格数据
        PlaceChoiceallsupplier(){
            // 传参
            var listByPageDataone = {
                type:2,
                searchStr:this.PlaceChoicepurSearch,
                page:{
                    current:this.PlaceChoicecurrent,
                    size:this.PlaceChoicesize
                }
            }
            // 渲染表格
            this.api.capital.all(listByPageDataone)
            .then(res=>{
                if(res.data.code == 200){
                    this.PlaceChoicetableData = res.data.data.records
                    this.PlaceChoicetotal = res.data.data.total
                }
            })
        },
        // 销售结算搜索功能
        PlaceChoicesearch() {
            if(this.PlaceChoicetableData != []){
                this.PlaceChoicetableData = [];
                this.PlaceChoiceallsupplier();
                this.PlaceChoicecurrent = 1
            }
        },
        // 新建销售结算
        PlaceChoiceCompanyshow(){
            // 跳转新建销售结算
            this.$router.push({name:"AddSalesBill"})
        },
        // 销售结算选择表格
        relationChange(row){
            this.tags = row;
            let num = 0; 
            for(let i=0;i<this.tags.length;i++){
                num += this.tags[i].sumPrice;
            }
            this.moeny = num.toFixed(2) * 1;
        },
        keyId(row){
            return row.id;
        },
        // 选择销售结算弹窗
        PlaceChoice(){
            if(this.PlaceChoicetableData != []){
                this.PlaceChoiceallsupplier();
            }
            this.PlaceChoiceUnit=true;       
        },
        // 选销售结算 关闭弹窗
        PlaceChoicehandleClose() {
            this.PlaceChoicepurSearch = "",
            this.PlaceChoiceUnit = false;               // 关闭整个弹窗
        },
        // 销售结算应用按钮 事件
        deliveryPlaceApplication(){
            this.AddData.sments = this.tags;            // 获取地址名称,回显到表单
            this.PlaceChoiceUnit = false;               // 关闭整个弹窗
        },
    },
    computed: {},
    watch: {}
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
    .pur-size {
        font-size: 16px;
        line-height: 32px;
        color: #333;
        font-weight: 600;
    }
    .pur-right {
        .el-button {
            margin-left: 16px;
            float: right;
        }
    }
}
// 头部结束
// 表单内容开始
.smbox{
    width:100%;
    min-height:78px;
    border: 1px solid #DCDFE6;
    border-radius: 4px;
    padding:8px;
    box-sizing: border-box;
}
.details{
   width:100%;
   height:32px;
   border: 1px solid #DCDFE6;
   border-radius: 4px;
   .details-nav{
       background-color:#F5F7FA;
        width: calc(100% - 73px);
        height: 32px;
        float: left;
        overflow: hidden;
        .details-tag{
            width: 44%;
            height: 32px;
            display: inline-block;
            .tag-box{
                width: 100%;
                display:inline-block;
                height: 32px;
                color: #C4C4C4;
                .el-tag{
                    color:#333;
                    width: 80%;
                    overflow: hidden;           // 盒子溢出隐藏
                    text-overflow:ellipsis;     // 文字溢出显示省略号
                    white-space: nowrap;        // 文字不换行
                }
                b{
                    width: 20%;
                    display: inline-block;
                }
            }
        }
        .details-conent{
            width:56%;
            text-align: right;
            height: 32px;
            float: right;
            padding-right: 4px;
            box-sizing: border-box;
            overflow: hidden;           // 盒子溢出隐藏
            text-overflow:ellipsis;     // 文字溢出显示省略号
            white-space: nowrap;        // 文字不换行
        }
   }  
   .details-size{
       background-color: #F5F7FA;
        border-left: 1px solid #DCDFE6;
        height:100%;
        text-align: center;
        width:73px;
        box-sizing: border-box;
        color:#333;
        cursor: pointer;
        float: right;
        span{
            width: 100%;
            height: 32px;
            font-size: 13px;
        }
   }
}
.pur-nav {
    height: calc(100vh - 110px);
    background-color: #fff;
    margin: 0 60px;
    box-shadow: 2px 2px 10px rgba(102, 102, 102, 0.15);
    padding:33px 0 33px 33px;
    box-sizing: border-box;
    .pur-table {
        width: 100%;
        height: 100%;
        overflow-y: scroll;
        .see-top {
            height: 40px;
            font-weight: 600;
            font-size: 14px;
            line-height: 45px;
            border-bottom: 1px dashed #999999;
            margin-bottom:24px;
        }
        .see-lable {
            padding: 15px 0;
            .el-input{
                .input{
                    width: 430px;
                }
            }
        }
    }
}
.AddPurIcon{
    width: 16px;
    height: 16px;
    border-radius: 50%;
    border: 1px solid #828282;
    color: #828282;
    font-size:12px;
    text-align: center;
    line-height: 16px;
    margin-top: 14px;
    margin-left: -25px;
}
.AddbuttonUpload{
    width: 100px;
    height: 32px;
    background: #FFFFFF;
    border: 1px solid #E4E4E4;
    border-radius: 4px;
    font-size: 14px;
    color: #666;
    line-height: 32px;
    text-align: center;
}
.AddbuttonUpload:hover{
    background: #e4e4e4;
    color: #333;
}
// 表单内容结束
// 弹窗样式开始
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
    // position: relative;
}
.totalTab {
    width: 9%;
    float: left;
    font-size: 16px;
    line-height: 26px;
    font-weight: 650;
    color: #333;
    text-align: center;
    margin:9px 0 0 35px;
}
.textlast{
    // position: absolute;
    // top: 0;
    // right:0;
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
.totalTab.active {
    color: #333;
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
.pur-pages {
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
    padding: 8px 30px;
    position: absolute;
    left: 0;
    bottom:148px;
}

// 弹窗样式结束
</style>