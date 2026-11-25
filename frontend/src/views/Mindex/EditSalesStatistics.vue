<template>
    <!--编辑销项发票页面 -->
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
                        <div class="see-top">进项信息</div>
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
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">销售合同</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractCode" placeholder="请输入销售合同" disabled>
                                        <span slot="suffix" style="color:#333;cursor: pointer;padding-right:16px;" @click="contractNumberdialogSelect">关 联</span>
                                    </el-input>
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
                                <label style="color:#666;font-size:14px;line-height:30px;">所属月份</label>
                            </el-col>
                            <el-form-item prop="month">
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.month"
                                    type="month"
                                    placeholder="请选择所属月份"
                                    value-format="yyyy-MM"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发票代码</label>
                            </el-col>
                            <el-form-item prop="billCode">
                                <el-col :span="7">
                                    <el-input v-model="AddData.billCode" placeholder="请输入发票代码"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发票号码</label>
                            </el-col>
                            <el-form-item prop="billNumber">
                                <el-col :span="7">
                                    <el-input v-model="AddData.billNumber" placeholder="请输入发票号码"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">开票日期</label>
                            </el-col>
                            <el-form-item prop="openDate">
                                <el-col :span="7" >
                                    <el-date-picker
                                    v-model="AddData.openDate"
                                    type="date"
                                    placeholder="请选择日期"
                                    value-format="yyyy-MM-dd"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">税号</label>
                            </el-col>
                            <el-form-item prop="dutyParagraph">
                                <el-col :span="7">
                                    <el-input v-model="AddData.dutyParagraph" placeholder="请输入税号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">购方企业名称</label>
                            </el-col>
                            <el-form-item prop="companyName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.companyName" placeholder="请输购方企业名称全称" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:16px;cursor: pointer;" @click="dialogSelect">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">商品名称</label>
                            </el-col>
                            <el-form-item prop="goodsName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.goodsName" placeholder="请输入商品名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">单位</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.uit" placeholder="请输入单位"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">数量</label>
                            </el-col>
                            <el-form-item prop="count">
                                <el-col :span="7">
                                    <el-input v-model="AddData.count" placeholder="吨"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">单价(不含税)</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7" prop="unitPrice">
                                    <el-input v-model="AddData.unitPrice" placeholder="元 / 吨"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>                                                           
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">金额</label>
                            </el-col>
                            <el-form-item prop="baseMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.baseMoney" placeholder="元" @input="CardNumberbaseMoney" :maxlength="length"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">税率</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.taxRate" placeholder="请输入税率"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>   
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">税额</label>
                            </el-col>
                            <el-form-item prop="taxAmount">
                                <el-col :span="7">
                                    <el-input v-model="AddData.taxAmount" placeholder="元" @input="CardNumbertaxAmount" :maxlength="taxAmountlength"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>     
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">含税金额</label>
                            </el-col>
                            <el-form-item prop="includMoney">
                                <el-col :span="7">
                                    <el-input v-model="AddData.includMoney" placeholder="元" @input="CardNumberincludMoney" :maxlength="includMoneylength"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>   
                        <el-row>
                            <el-col :span="2">
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">规格型号</label>
                            </el-col>
                            <el-form-item prop="modle">
                                <el-col :span="7">
                                    <el-input v-model="AddData.modle" placeholder="请输入规格型号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row> 
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

        <!-- 客户单位弹窗开始 -->
        <el-dialog
            :visible.sync="dialogSelectUnit"
            width="60%"
            :before-close="BuyerhandleClose">
            <template slot="title">
                <span>选择客户</span>
            </template>
            <!-- 客户单位内容框 -->
            <el-input type="textarea" v-model="dialogtextarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textbutton">
                        <!-- 客户应用按钮 -->
                        <el-button size="small" @click="salesUnitApplication">应用</el-button>
                    </div>
                    <div class="text">
                        <!-- 客户单位信息 -->
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
                                <el-table-column prop="enterpriseName" label="单位名称" :render-header="renderHeaderMethods">
                                </el-table-column>
                                <el-table-column prop="createTime" label="创建时间">
                                </el-table-column>
                            </el-table>
                            <!-- 分页 -->
                            <div class="pur-pages">
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
        <!--客户单位弹窗结束 -->

        <!-- 销售合同关联弹窗开始 -->
        <el-dialog
            :visible.sync="PlaceChoicecontractNumber"
            width="60%"
            :before-close="contractNumberhandleClose">
            <template slot="title">
                <div style="font-size:18px;">关联合同</div>
            </template>
            <!-- 显示合同数据选择内容框 -->
            <el-input type="textarea" v-model="contracttarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="contractsearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="contractNumberSearch"
                        style="width:200px;">
                        </el-input>
                    </div>
                    <!-- 合同表格 -->
                    <el-table
                    border
                    :data="contractNumbertableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="contractradio" @change.native="contractgetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="contractCode" label="销售合同" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="contractName" label="合同名称" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="contractNum" label="合同数量" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column label="总价金额" :show-overflow-tooltip="true">
                            <template slot-scope="scope">
                                <span v-if="scope.row.sumPrice">{{scope.row.sumPrice|currency}} 元</span>
                            </template>
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="contractpaginationshow"
                        background
                        :page-size.sync="contractsize"
                        :current-page.sync="contractcurrent"
                        layout="prev, pager, next"
                        :total="contracttotal"
                        @current-change='contractall'
                        @size-change='contractall'>
                        </el-pagination>
                        <!-- 合同关联应用按钮 -->
                        <el-button size="small" @click="contractNumberApplication">应用</el-button>  
                    </div>
                </div>
            </div>
        </el-dialog>
        <!-- 销售合同关联弹窗结束 -->

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
                        v-show="PlaceChoicepaginationshow"
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
            PlaceChoicepaginationshow:true,         // 无数据隐藏分页
            PlaceChoicecurrent:1,                   // 当前页数
            PlaceChoicesize:4,                      // 每页显示条数
            PlaceChoicetotal:0,                     // 总条数
            PlaceChoicepurSearch:'',                // 搜索
            PlaceChoiceUnit:false,                  // 弹窗
            PlaceChoicetarea:'',                    // 显示销售结算内容框
            PlaceChoicetableData:[],                // 选择销售结算表格数据
            tags: [],                               // 销售结算回显数据
            moeny:0,                               // 合计金额

            /* 
            合同关联数据列表
            */
            contractpaginationshow:true,            // 无数据隐藏分页(合同)
            contractcurrent:1,                      // 当前页数
            contractsize:4,                         // 显示条数
            contracttotal:0,                        // 总条数
            PlaceChoicecontractNumber:false,        // 选择合同弹窗
            contracttarea:'',                       // 显示合同选中内容框
            contractNumberSearch:'',                // 合同搜索框
            contractNumbertableData:[],             // 选择合同表格数据
            contractradio:'',                       // 合同数据单选框

            /* 
            客户单位数据列表
            */
            paginationshowtwo:true,                 // 无数据隐藏分页(客户)
            currenttwo:1,                   
            sizetwo:4,                      
            totaltwo:0,                     
            purSearchtwo:'',                               
            dialogSelectUnit:false,                 // 选择单位弹窗
            dialogtextarea:'',                      // 显示单位内容框
            CompanytableData:[],                    // 选择单位表格数据
            radio:'',                               // 选择单位单选框
            show:false,                             // 选择单位表格隐藏 id
            CompanyDatashow:false,                  // 新建单位表单显示隐藏
            CompanyNewData:{},                      // 新建单位表单
            // 新建单位表单验证
            rels: {
                custType: [
                    { required: true, message: '请选择类型', trigger: 'blur' },
                ],
                enterpriseName: [
                    { required: true, message: '请输入单位名称', trigger: 'blur' },
                ],
            },

            /* 
           编辑销项发票数据
            */
            length:'',
            taxAmountlength:'',   
            includMoneylength:'',
            http:'',
            header:{},                      // 上传请求头
            fileObj:{},                     // 附件文件格式验证
            enclosurefileList: [],          // 附件文件列表
            purTitle: "", // 标题
            gid:'',
            placecontent:'',
            // 表单验证
            rules: {
                systemCode: [
                    { required: true, message: '请输入系统编码', trigger: 'change' }
                ],
                contractCode: [
                    { required: true, message: '请关联采购合同', trigger: 'change' }
                ],
                month: [
                    { required: true, message: '请选择所属月份', trigger: 'blur' }
                ],
                inType: [
                    { required: true, message: '请输入进项类型', trigger: 'blur' }
                ],
                billCode: [
                    { required: true, message: '请输入发票代码', trigger: 'blur' },
                    { pattern:/^[0-9]*$/, message: '发票代码格式不正确', trigger: 'blur'}
                ],
                billNumber: [
                    { required: true, message: '请输入发票号码', trigger: 'blur' },
                    { pattern:/^[0-9]*$/, message: '发票号码格式不正确', trigger: 'blur'}
                ],
                openDate: [
                    { required: true, message: '请选择开票日期', trigger: 'blur' }
                ],
                dutyParagraph: [
                    { required: true, message: '请输入销方税号', trigger: 'blur' },
                    { pattern:/^[A-Z0-9]+$/, message: '发票号码格式不正确', trigger: 'blur'}
                ],
                companyName: [
                    { required: true, message: '请输入销方企业名称', trigger: 'change' }
                ],
                taxAmount: [
                    { required: true, message: '请输入税额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '税额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                includMoney: [
                    { required: true, message: '请输入含税金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '含税金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                count: [
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '数量格式不正确', trigger: 'blur'}
                ],
                unitPrice: [
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '单价格式不正确', trigger: 'blur'}
                ],
                baseMoney: [
                    { required: true, message: '请输入金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
            },
            AddData: {
                id:'',
                baseMoney:'',
                billCode:'',
                billNumber:'',
                companyId:'',
                companyName:'',
                contractCode:'',
                contractId:'',
                count:'',
                dutyParagraph:'',
                goodsName:'',
                ids:[],
                month:'',
                includMoney:'',
                relationFlag:'',
                openDate:'',
                systemCode:'',
                taxAmount:'',
                taxRate:'',
                type:'',
                uit:'',
                unitPrice:'',
                modle:'',
                filelist:[]  // 合同 协议的集合id 
            }, // 新建表单数据
        };
    },
    created() {},
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 接收主页单条数据
        this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));
        if(!this.AddData){
            // 直接输入网址 / 后退进入时没有待编辑的记录:回列表页,否则下面取字段就抛错
            return this.$router.replace({name:'SalesStatistics'});
        }
        // 关联销售结算数据集合回显
        this.api.xxbill.oneAll({id:this.AddData.id})
        .then(res=>{
            this.tags = res.data.data;
            let num = 0; 
            for(let i=0;i<this.tags.length;i++){
                num += this.tags[i].sumPrice;
            }
            this.moeny = num.toFixed(2) * 1;
            var arr = [];
            for(let i=0;i<this.tags.length;i++){
                arr.push(this.tags[i].id)
            }
            this.AddData.ids = arr;
        })
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')}
        this.http  = this.$store.state.http;                 // IP地址
        // 回显文件列表
        this.api.file.one({relationId:JSON.parse(sessionStorage.getItem('queryAll')).id})
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
                // 附件文件内容
                this.enclosurefileList.push(newarr[i]); 
            }
        })
    },
    methods: {
        /* 
       编辑销项发票表单提交功能
        */ 
        // 金额输入框限制
        CardNumberbaseMoney(val) {
            if(val.indexOf('.') != -1){
                this.length = 12;
            }else{
                this.length = 10;
            }
        }, 
        // 税额输入框限制
        CardNumbertaxAmount(val) {
            if(val.indexOf('.') != -1){
                this.taxAmountlength = 12;
            }else{
                this.taxAmountlength = 10;
            }
        }, 
        // 含税金额输入框限制
        CardNumberincludMoney(val) {
            if(val.indexOf('.') != -1){
                this.includMoneylength = 12;
            }else{
                this.includMoneylength = 10;
            }
        }, 
        // 编辑销项发票表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    // 上传文件
                    var arr = [];
                    this.enclosurefileList.map((itme,index)=>{
                        arr.push({id:itme.response.data[0].id});
                    }).join(',')
                    this.AddData.filelist = arr;
                    this.AddData.contractId = this.placecontent.id      // 合同编号id
                    // 销售结算
                    if(this.AddData.ids.length == 0){
                        this.AddData.relationFlag = false;
                    }else{
                        this.AddData.relationFlag = true;
                    }
                    this.api.xxbill.edit(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                    this.delRemovedFiles();
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    sessionStorage.removeItem('queryAll');
                                    that.$router.push({name:"SalesStatistics"});
                                }
                            });
                        }
                    })
                }
            });
        },
        // 取消保存表单
        Addopen() {
            sessionStorage.removeItem('queryAll');
            this.$router.push({name:"SalesStatistics"});
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
        /* 
        合同关联的选择功能
        */
        // 表格数据
        contractall(){
            // 传参
            var listByPageData = {
                contractType:2,
                searchStr:this.contractNumberSearch,
                page:{
                current:this.contractcurrent,
                size:this.contractsize
                }
            }
            // 渲染表格
            this.api.listByPage.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                if(res.data.data.records != ""){
                    this.contractNumbertableData = res.data.data.records
                    this.contracttotal = res.data.data.total
                }else{
                    this.contractpaginationshow = false;
                }
                }
            })
        },
        // 合同关联搜索功能
        contractsearch() {
           this.contractcurrent = 1;
           this.contractall();
        },
        // 合同编号选择表格
        contractgetCurrentRow(row){
            this.placecontent = row;                    // 获取选中单位数据
            if( row.contractCode == ""){
                // 获取选中合同名称
                this.contracttarea = row.contractName;
            }else{
                // 获取选中合同编号数据
                this.contracttarea = row.contractCode;
            }
        },
        // 合同数据弹窗
        contractNumberdialogSelect(){
            this.contractall();
            this.PlaceChoicecontractNumber=true; //开启弹窗
        },
        // 选择合同关联关闭弹窗
        contractNumberhandleClose() {
            this.contractradio = "";                        // 清空单选框
            this.contracttarea = "";                        // 清空显示地址内容框
            this.PlaceChoicecontractNumber = false;         // 关闭整个弹窗
        },
        // 合同数据应用按钮 事件
        contractNumberApplication(){
            if(!this.placecontent || !this.placecontent.id){
                // 没选行就点「应用」:原先取 id 抛错,表单里留下空值
                return this.$message({ type: "warning", message: "请先选择一条合同" });
            }
            this.AddData.contractId = this.placecontent.id;
            this.AddData.contractCode = this.contracttarea                      // 销售合同
            this.AddData.companyName = this.placecontent.buyerUnit              // 销售合同的购货单位(销售方是本企业)
            this.PlaceChoicecontractNumber = false;                             // 关闭整个弹窗
            this.contractradio = "";                                            // 清空单选框
            this.contracttarea = "";                                            // 清空显示地址内容框
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
                    if(res.data.data.records != ""){
                        this.PlaceChoicetableData = res.data.data.records
                        this.PlaceChoicetotal = res.data.data.total
                    }else{
                        this.PlaceChoicepaginationshow = false;
                    }
                }
            })
        },
        // 销售结算搜索功能
        PlaceChoicesearch() {
            if(this.PlaceChoicetableData != []){
                this.PlaceChoicetableData = [];
                this.PlaceChoiceallsupplier();
                this.PlaceChoicepaginationshow = true;
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
                this.PlaceChoicepaginationshow = true;
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
            var arr = [];
            for(let i=0;i<this.tags.length;i++){
                arr.push(this.tags[i].id)
            }
            this.AddData.ids = arr;                     // 获取地址名称,回显到表单
            this.PlaceChoiceUnit = false;               // 关闭整个弹窗
        },

        /* 
        客户单位的选择功能
        */
        // 客户单位表格数据
        allsuppliertwo(){
            // 传参
            var listByPageDatatwo = {
                partnerType:2,
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
        // 客户单位搜索功能(客户)
        searchtwo() {
            this.CompanytableData = [];
            this.allsuppliertwo();
            this.currenttwo = 1;
        },
        // 客户单位 选择单位表格
        getCurrentRow(row){
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
                    this.CompanyNewData.partnerType = 2;
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
        // 关联客户单位弹窗
        dialogSelect(){
            this.allsuppliertwo();                              // 渲染表格
            this.dialogSelectUnit=true;                         // 开启弹窗
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.dialogSelectUnit == true){
                this.CompanyDatashow = false
            }
        },
        // 关联客户单位应用按钮事件
        salesUnitApplication(){
            this.AddData.companyName = this.dialogtextarea;    // 获取单位名称，回显到表单
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
            this.radio = "";                                    // 清空单选框
            this.dialogtextarea = "";                           // 清空显示单位内容框
        },
        // 关联客户单位 关闭弹窗
        BuyerhandleClose() {
            this.radio = "";                                    // 清空单选框
            this.dialogtextarea = "";                           // 清空显示单位内容框
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
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
// 关联合同弹窗样式开始
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
    padding: 8px 16px;
    position: absolute;
    left: 0;
    bottom:115px;
}
// 关联合同弹窗样式结束
</style>