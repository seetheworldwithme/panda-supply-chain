<template>
    <!-- 新建销售合同页面 -->
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
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">系统编码</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.systemCode" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- <el-row>
                            <el-col :span="2" >
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">项目类型</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-select v-model="AddData.contractResource" placeholder="请选择项目类型" style="width:100%;">
                                        <el-option
                                        v-for="item in projectName"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                        <el-row>
                            <el-col :span="2">
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">合同编号</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractCode" placeholder="请输入合同编号" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">合同名称</label>
                            </el-col>
                            <el-form-item prop="contractName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractName" placeholder="请输入合同名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签定时间</label>
                            </el-col>
                            <el-form-item prop="signingTime">
                                <el-col :span="7" >
                                    <el-date-picker
                                    @change="clickTime"
                                    v-model="AddData.signingTime"
                                    type="date"
                                    placeholder="请选择日期"
                                    value-format="yyyy-MM-dd"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
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
                                    placeholder="请选择月份"
                                    value-format="yyyy-MM"
                                    style="width:100%;">
                                    </el-date-picker>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">销售单位</label>
                            </el-col>
                            <el-form-item prop="salesUnit">
                                <el-col :span="7">
                                    <el-input v-model="AddData.salesUnit" placeholder="请输入销售单位全称">
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">购货单位</label>
                            </el-col>
                            <el-form-item prop="buyerUnit">
                                <el-col :span="7">
                                    <el-input v-model="AddData.buyerUnit" placeholder="请输入购货单位全称" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="dialogSelect">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">发货地</label>
                            </el-col>
                            <el-form-item prop="deliveryPlace">
                                <el-col :span="7">
                                    <el-input v-model="AddData.deliveryPlace" placeholder="请输入发货地全称" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="PlaceChoice">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">收货地</label>
                            </el-col>
                            <el-form-item prop="receiptPlace">
                                <el-col :span="7">
                                    <el-input v-model="AddData.receiptPlace" placeholder="请输入收货地全称" disabled>
                                        <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="Placeharvest">选 择</span>
                                    </el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                    <!-- 产品信息 -->
                    <el-col :span="24">
                        <div class="see-top">产品信息</div>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">产品大类选择</label>
                            </el-col>
                            <el-form-item prop="proBigType">
                                <el-col :span="2" v-for="(v,k) in AddproBigType" :key="k">
                                <el-radio v-model="AddData.proBigType" :label="v.dvalue * 1" @change="getValue()">{{v.dkey}}</el-radio>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">产品小类选择</label>
                            </el-col>
                            <el-form-item prop="proSmaType">
                                <el-col :span="7">
                                    <el-select filterable v-model="AddData.proSmaType" placeholder="请选择产品小类" style="width:100%;" @change="proSmaclick($event)">
                                        <el-option
                                        v-for="item in SelectproSmaType"
                                        :key="item.dvalue"
                                        :label="item.dkey"
                                        :value="item.dvalue * 1"
                                        >
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row v-if="pIDshow">
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">产品名称</label>
                            </el-col>
                            <el-form-item prop="proRemark">
                                <el-col :span="7">
                                    <el-input v-model="AddData.proRemark" placeholder="请输入产品名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">合同数量</label>
                            </el-col>
                            <el-form-item prop="contractNum">
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractNum" placeholder="吨"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">货品单价</label>
                            </el-col>
                            <el-form-item prop="unitPrice">
                                <el-input v-model="AddData.unitPrice" placeholder="元 / 吨" style="width:200px;" :maxlength="10"></el-input>
                                <el-radio v-model="AddData.taxIncluded" :label="true" style="margin-left:40px;">含税</el-radio>
                                <el-radio disabled v-model="AddData.taxIncluded" :label="false">不含税</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">总价金额</label>
                            </el-col>
                            <el-form-item prop="sumPrice">
                                <el-col :span="7">
                                    <el-input @focus="focusData" v-model="AddData.sumPrice" placeholder="元" @input="CardNumber" :maxlength="length"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">标的指标</label>
                            </el-col>
                            <el-form-item>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.sulfur" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">硫分St.d%</p>
                                        
                                    </el-col>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.volatilization" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">挥发分Vdaf%</p>
                                        
                                    </el-col>
                                    <el-col :span="2" style="margin-top:5px;">
                                        <el-input v-model="AddData.water" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">水分Mt%</p>
                                        
                                    </el-col>
                                    <el-col :span="2" v-for="(v, index) in AddData.list" :key="index" style="margin-top:5px;">
                                        <el-input v-model="v.dvalue" style="width:64px;"></el-input>
                                        <p style="margin-right:10px;">{{v.dkey}}</p>
                                    </el-col>
                                    <el-col :span="1">
                                        <el-tooltip effect="dark" content="标的指标需要在基础数据里设置" placement="bottom">
                                            <div class="AddPurIcon">?</div>
                                        </el-tooltip>
                                    </el-col>                                 
                            </el-form-item>
                        </el-row>                        
                        <el-row>
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">备注</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input type="textarea" v-model="AddData.remark" placeholder="请输入备注内容"></el-input>
                                </el-col>
                            </el-form-item> 
                        </el-row>
                    </el-col>
                    <!-- 扩展信息 -->
                    <el-col :span="24">
                        <div class="see-top">扩展信息</div>
                        <el-row v-for="(vl,k) in Extended.base_module" :key="k">
                            <el-col :span="2" >
                                <!-- <span style="color:red;">* </span> -->
                                <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">{{vl.dkey}}</label>
                            </el-col>
                            <el-form-item>
                                <!-- 价格类型 -->
                                <el-col :span="22" v-if="vl.childKey =='price_type'">
                                    <el-radio v-model="AddData.priceType" v-for="(val,kal) in Extended.price_type" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                </el-col>
                                <!-- 付款方式 -->
                                <el-col :span="22" v-if="vl.childKey =='payment_method'">
                                    <el-radio v-model="AddData.payMethod" v-for="(val,kal) in Extended.payment_method" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                </el-col>
                                <!-- 支付方式 -->
                                <el-col :span="22"  v-if="vl.childKey =='payment_type'">
                                    <el-radio v-model="AddData.payType" v-for="(val,kal) in Extended.payment_type" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                </el-col>
                                <!-- 运输模式 -->
                                <el-col :span="22"  v-if="vl.childKey =='transport_mode'">
                                    <el-radio v-model="AddData.transportMode" v-for="(val,kal) in Extended.transport_mode" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                </el-col>
                                <!-- 承运工具 -->
                                <el-col :span="22"  v-if="vl.childKey =='carrier'">
                                    <el-radio v-model="AddData.carrier" v-for="(val,kal) in Extended.carrier" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
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
                                    :limit="1"
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
                        <el-row v-if="AdjustForm">
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">在线支付</label>
                            </el-col>
                            <el-form-item prop="isPayOnline">
                                <el-radio v-model="AddData.isPayOnline" :label="true" @change="isPayOnlinelogistics($event)">开启</el-radio>
                                <el-radio v-model="AddData.isPayOnline" :label="false">关闭</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">电子签章</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.isEsign" :label="true" @change="isEsignlogistics($event)">开启</el-radio>
                                <el-radio v-model="AddData.isEsign" :label="false">关闭</el-radio>
                            </el-form-item>
                        </el-row>
                    </el-col>
                    <!-- 开启电子签章 -->
                    <div v-if="AddData.isEsign">
                        <!-- 发起方信息 -->
                        <el-col :span="24">
                            <div class="see-top">发起方信息
                                <span style="color:#CF1724;margin-left:3.5%;">* </span>
                                <span style="color:#666;font-size:12px;line-height:40px;">以下信息将用于电子签章服务注册，请谨慎填写！</span>
                            </div>
                            <el-row>
                                <el-col :span="2.5">
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">发起方</label>
                                </el-col>
                                <el-form-item prop="launchType">
                                    <el-col :span="7">  
                                        <el-radio v-model="AddData.launchType" :label="2">企业</el-radio>               
                                        <el-radio v-model="AddData.launchType" :label="0">个人</el-radio>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row v-if="AddData.launchType == 2">
                                <el-col :span="2.5">
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">发起方企业</label>
                                </el-col>
                                <el-form-item prop="etcName">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.etcName" placeholder="请输入企业名称" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">发起方姓名</label>
                                </el-col>
                                <el-form-item prop="participants">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.participants" placeholder="发起方姓名" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">发起方手机号</label>
                                </el-col>
                                <el-form-item prop="phone">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.phone" placeholder="请输入手机号" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </el-col>
                        <!-- 签署方信息 -->
                        <el-col :span="24">
                            <div class="see-top">签署方信息</div>
                            <el-row>
                                <el-col :span="2.5" >
                                    <span style="color:#CF1724;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">参与人</label>
                                </el-col>
                                <el-form-item prop="signRole">
                                    <el-col :span="7">  
                                        <el-radio v-model="AddData.signRole" :label="2" @change="prologistics($event)">企业</el-radio>               
                                        <el-radio v-model="AddData.signRole" :label="0" @change="prologistics($event)">个人</el-radio>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <!-- 企业 -->
                            <div v-if="AddData.signRole == 2" >
                                <!-- <div style="width:15px;height:20px;"></div> -->
                                <el-row v-for="(domain,index) in AddData.domains" :key="index">
                                    <el-col :span="2.5">
                                        <span style="color:#CF1724;">* </span>
                                        <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                                    </el-col>
                                    <el-form-item :prop="'domains.' + index + '.value'" :rules="{ required: true, message: '请输入' + domain.message + domain.id, trigger: 'blur' }">
                                        <el-col :span="7">
                                            <el-input v-model="domain.value" :placeholder="'请输入' + domain.label"></el-input>
                                        </el-col>
                                        <el-col :span="7" v-if="domain.labelone">
                                            <div class="removeDomainoneBox" @click.prevent="removeDomain(domain)">
                                                <i class="el-icon-remove"></i>
                                                <span> 删除</span>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                            </div>
                            <!-- 个人 -->
                            <div v-if="AddData.signRole == 0">
                                <el-row v-for="(domain,key) in AddData.domainsone" :key="key">
                                    <el-col :span="2.5">
                                        <span style="color:#CF1724;">* </span>
                                        <label style="color:#666;font-size:14px;line-height:30px;">{{domain.label}}</label>
                                    </el-col>
                                    <el-form-item :prop="'domainsone.' + key + '.value'" :rules="{ required: true, message: '请输入' + domain.message  + domain.id, trigger: 'blur' }">
                                        <el-col :span="7">
                                            <el-input v-model="domain.value" :placeholder="'请输入' + domain.label"></el-input>  
                                        </el-col>
                                        <el-col :span="2" v-if="domain.labelone">
                                            <div class="removeDomainoneBox" @click.prevent="removeDomainone(domain)">
                                                <i class="el-icon-remove"></i>
                                                <span> 删除</span>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                            </div>
                            <div>
                                <el-row>
                                    <el-form-item>
                                        <el-col :span="7" v-if="AddData.signRole == 2">
                                            <div class="zjcyfspan" @click="addDomain">
                                                <i class="el-icon-circle-plus"></i><span > 增加参与方</span>
                                            </div>
                                        </el-col>
                                        <el-col :span="7" v-if="AddData.signRole == 0">
                                            <div class="zjcyfspan" @click="addDomainone">
                                                <i class="el-icon-circle-plus"></i><span > 增加参与方</span>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                            </div>
                        </el-col>
                    </div>
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
                <span>选择客户单位</span>
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
                                <el-table-column prop="enterpriseName" label="单位名称" :render-header="renderHeaderMethods" :show-overflow-tooltip="true">
                                </el-table-column>
                                <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true">
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
                            <el-col :span="10">
                                <el-form-item label=" 单 位 ：" prop="enterpriseName">
                                    <!-- <el-autocomplete
                                    class="inline-input"
                                    v-model="CompanyNewData.enterpriseName"
                                    :fetch-suggestions="querySearch"
                                    placeholder="请输入单位名称"
                                    @select="handleSelect"
                                    style="width:90%;"
                                    ></el-autocomplete> -->
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

        <!-- 发货地弹窗开始 -->
        <!-- 发货地址选择 -->
        <el-dialog
            :visible.sync="PlaceChoiceUnit"
            width="60%"
            :before-close="PlaceChoicehandleClose">
            <template slot="title">
                <div style="font-size:18px;">地址选择</div>
            </template>
            <!-- 显示发货地地址选择内容框 -->
            <el-input type="textarea" v-model="PlaceChoicetarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
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
                        <span @click="PlaceChoiceCompanyshow">+ 新建</span>
                    </div>
                    <!-- 发货地表格 -->
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="PlaceChoicetableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="Placeradio" @change.native="PlaceChoicegetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="addressName" label="地址名称" :render-header="renderHeaderMethodsone" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="PlaceChoicepaginationshow"
                        background
                        :page-size.sync="PlaceChoicesize"
                        :current-page.sync="PlaceChoicecurrent"
                        layout="prev, pager, next"
                        :total="PlaceChoicetotal"
                        @current-change='PlaceChoiceallsupplier'
                        @size-change='PlaceChoiceallsupplier'>
                        </el-pagination>
                        <!-- 发货地应用按钮 -->
                        <el-button size="small" @click="deliveryPlaceApplication">应用</el-button> 
                    </div>
                    <!-- 新建发货地 -->
                    <el-form :model="PlaceChoiceNewData" :rules="Placerels" ref="PlaceChoiceNewData" label-width="80px" size="small" v-show="PlaceChoiceshow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="8">
                                <el-form-item label=" 地 址 ：" prop="addressName">
                                    <el-input v-model="PlaceChoiceNewData.addressName" placeholder="请输入地址名称" style="width:260px;"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="2" style="padding-left:25px;">
                                <el-button @click="PlaceChoiceNewShow" size="small">取消</el-button>
                            </el-col>
                            <el-col :span="2" style="padding-left:15px;">
                                <el-button size="small" @click="PlaceChoiceNewsubmitForm('PlaceChoiceNewData')">保存</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 发货地弹窗结束 -->

        <!-- 收货地弹窗开始 -->
        <!-- 收货地址选择 -->
        <el-dialog
            :visible.sync="PlaceharvestUnit"
            width="60%"
            :before-close="PlaceharvesthandleClose">
            <template slot="title">
                <div style="font-size:18px;">地址选择</div>
            </template>
            <!-- 显示收货地地址选择内容框 -->
            <el-input type="textarea" v-model="Placeharvesttarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="Placeharvestsearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="PlaceharvestpurSearch"
                        style="width:200px;float:right;"
                        >
                        </el-input>
                        <span @click="PlaceharvestCompanyshow">+ 新建</span>
                    </div>
                    <!-- 收货地表格 -->
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="PlaceharvesttableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="Placeharvestradio" @change.native="PlaceharvestgetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="addressName" label="地址名称" :render-header="renderHeaderMethodsone" :show-overflow-tooltip="true">
                        </el-table-column>
                        <el-table-column prop="createTime" label="创建时间" :show-overflow-tooltip="true">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="Placeharvestpaginationshow"
                        background
                        :page-size.sync="Placeharvestsize"
                        :current-page.sync="Placeharvestcurrent"
                        layout="prev, pager, next"
                        :total="Placeharvesttotal"
                        @current-change='Placeharvestallsupplier'
                        @size-change='Placeharvestallsupplier'>
                        </el-pagination>
                        <!-- 收货地应用按钮 -->
                        <el-button size="small" @click="PlaceharvestApplication">应用</el-button> 
                    </div>
                    <!-- 新建收货地 -->
                    <el-form :model="PlaceharvestNewData" :rules="Placeharvestrels" ref="PlaceharvestNewData" label-width="80px" size="small" v-show="Placeharvestshow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="8">
                                <el-form-item label=" 地 址 ：" prop="addressName">
                                    <el-input v-model="PlaceharvestNewData.addressName" placeholder="请输入地址名称" style="width:260px;"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="2" style="padding-left:25px;">
                                <el-button @click="PlaceharvestNewShow" size="small">取消</el-button>
                            </el-col>
                            <el-col :span="2" style="padding-left:15px;">
                                <el-button size="small" @click="PlaceharvestNewsubmitForm('PlaceharvestNewData')">保存</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 收货地弹窗结束 -->
    </el-form>
</template>
<script>
export default {
    data() {
        return {
            AdjustForm:null,
            /* 
            销售单位 购货单位数据列表
            */
            paginationshow:true, // 无数据隐藏分页(全部)
            current:1,  // 当前页数(全部)
            size:4,  // 显示条数(全部)
            totalone:0,  // 总条数(全部)
            purSearch:'',  // 选择单位搜索框(全部)
            paginationshowtwo:true, // 无数据隐藏分页(供应商)
            currenttwo:1,  // 当前页数(供应商)
            sizetwo:4,  // 显示条数(供应商)
            totaltwo:0,  // 总条数(供应商)
            purSearchtwo:'',  // 选择单位搜索框(供应商)
            paginationshowtrr:true, // 无数据隐藏分页(客户)
            currenttrr:1,  // 当前页数(客户)
            sizetrr:4,  // 显示条数(客户)
            totaltrr:0,  // 总条数(客户)
            purSearchtrr:'',  // 选择单位搜索框(客户)
            dialogSelectUnit:false,  // 选择单位弹窗
            dialogtextarea:'', // 显示单位内容框
            sel: 1,   // 选项卡默认显示tab
            tabs: [
                { label: "全部", id: 1 ,width:"width:35px;"},
                { label: "供应商", id: 2 ,width:"width:50px;"},
                { label: "客户", id: 3 ,width:"width:35px;"},
            ],
            CompanytableData:[],  // 选择单位表格数据
            radio:'',  // 选择单位单选框
            show:false,  // 选择单位表格隐藏 id
            CompanyDatashow:false,  //新建单位 显示隐藏
            // 新建单位表单
            CompanyNewData:{
                custType:'',
                partnerType:'',
                enterpriseName:'',
            },
            //  新建单位表单验证
            rels: {
                custType: [
                    { required: true, message: '请选择类型', trigger: 'blur' },
                ],
                partnerType: [
                    { required: true, message: '请选择合作方', trigger: 'blur' },
                ],
                enterpriseName: [
                    { required: true, message: '请输入单位名称', trigger: 'change' },
                ],
            },
            restaurants:[],
            salesUnitShow:false,  //销售单位应用按钮
            buyerUnitShow:false,   // 购货单位应用按钮
            
            /* 
            发货地数据列表
            */
            PlaceChoicepaginationshow:true, // 无数据隐藏分页
            PlaceChoicecurrent:1,  // 当前页数
            PlaceChoicesize:4,  // 每页显示条数
            PlaceChoicetotal:0,  // 总条数
            PlaceChoicepurSearch:'', // 搜索
            PlaceChoiceUnit:false,  // 弹窗
            PlaceChoicetarea:'', // 显示发货内容框
            PlaceChoiceshow:false,  //新建地址 显示隐藏
            // 新建发货地表单
            PlaceChoiceNewData:{
                addressName:'',
            },
            // 新建发货地表单验证
            Placerels: {
                addressName: [
                    { required: true, message: '请输入地址名称', trigger: 'blur' },
                ],
            },
            PlaceChoicetableData:[
            ],  // 选择发货地地址表格数据
            Placeradio:'',  // 发货地单选框

            /* 
            收货地数据列表
            */
            Placeharvestpaginationshow:true, // 无数据隐藏分页
            Placeharvestcurrent:1,  // 当前页数
            Placeharvestsize:4,  // 每页显示条数
            Placeharvesttotal:0,  // 总条数
            PlaceharvestpurSearch:'', // 搜索
            PlaceharvestUnit:false,  // 弹窗
            Placeharvesttarea:'', // 显示收货内容框
            Placeharvestshow:false,  //新建地址 显示隐藏
            // 新建收货地表单
            PlaceharvestNewData:{
                addressName:'',
            },
            // 新建收货地表单验证
            Placeharvestrels: {
                addressName: [
                    { required: true, message: '请输入地址名称', trigger: 'blur' },
                ],
            },
            PlaceharvesttableData:[],  // 选择收货地地址表格数据
            Placeharvestradio:'',  // 收货地单选框

            /* 
            编辑销售合同数据
            */
            length:'',
            http:'',
            header:{},                      // 上传请求头
            fileObj:{},                     // 附件文件格式验证
            enclosurefileList: [],          // 附件文件列表
            placecontent:{},  // 单位选中单条数据
            pIDshow:false,  // 产品名称显示隐藏
            purTitle: "", // 标题
            AddproBigType: [], // 产品大类选择
            SelectproSmaType: [],  // 小类选择下拉框
             // 项目类型选择下拉框
            // projectName:[
            //     { 
            //         value: 1,
            //         label: '供应链业务'
            //     },
            //     { 
            //         value: 2,
            //         label: '供应链金融'
            //     }
            // ], 
           // 表格数据
            AddData: {
                isPayOnline:false,
                systemCode: "CG.0000001202102220001",
                contractResource:"",
                contractCode: "",
                contractName: "",
                signingTime: "",
                salesUnit: "",
                buyerUnit: "",
                deliveryPlace: "",
                receiptPlace: "",
                proBigType: "",
                proSmaType: "",
                proRemark:"",
                contractNum: "",
                month:'',
                supplierId:'',
                unitPrice: "",
                taxIncluded:true,
                sumPrice: "",
                // 标的指标
                sulfur: "",
                volatilization: "",
                water: "",
                list:[],
                remark: "",
                priceType: "",
                payMethod: "",
                payType: "",
                transportMode: "",
                carrier: "",
                contractIds: '', // 合同 
                agreementIds: '', // 协议 字符串
                filelist:[],  // 合同 协议的集合id 

                isEsign:false,              // 电子签章
                participants:'',
                etcName:'',
                phone:'',
                launchType:2,               // 发起方类型
                domains: [],                // 企业
                domainsone:[],              // 个人
                signRole:2,                 // 参与方类型
                eqbFlow:{},
            }, 
            numbers:0,
            numbersone:0,
            // 编辑销售合同表单验证
            rules: {
                contractName: [
                    { required: true, message: '请输入合同名称', trigger: 'blur' }
                ],
                signingTime: [
                    { required: true, message: '请输入签订时间', trigger: 'blur' }
                ],
                salesUnit: [
                    { required: true, message: '请输入销售单位全称', trigger: 'blur' }
                ],
                contractCode: [
                    { required: true, message: '请输入合同编号', trigger: 'change' },
                    // { pattern: /(?=.*\d)(?=.*\W+)(?=.*[A-Z])(?!.*\n).*$/, message: '合同编号格式不正确', trigger: 'blur'}
                ],
                buyerUnit: [
                    { required: true, message: '请输入购货单位全称', trigger: 'change' }
                ],
                month: [
                    { required: true, message: '请选择所属月份', trigger: 'blur' }
                ],
                deliveryPlace: [
                    { required: true,message: '请输入发货地', trigger: 'change' }
                ],
                receiptPlace: [
                    { required: true, message: '请输入收货地', trigger: 'change' }
                ],
                contractNum: [
                    { required: true, message: '请输入合同数量', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '格式不正确', trigger: 'blur'}
                ],
                proBigType: [
                    { required: true, message: '请选择产品大类', trigger: 'blur' }
                ],
                proSmaType: [
                    { required: true, message: '请选择产品小类', trigger: 'blur' }
                ],
                proRemark: [
                    {required: true, message: '请选择产品名称', trigger: 'blur' }
                ],
                unitPrice: [
                    { required: true, message: '请输入货品单价', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '请输入正确的金额', trigger: 'blur'}
                ],
                sumPrice: [
                    { required: true, message: '请输入总价金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
            },
            Extended:'',  // 扩展信息菜单
        };
    },
    created() {
        this.api.glconfig.list()
        .then(res=>{
            if(res.data.code == 200){
                this.AdjustForm = res.data.data[0].isOpenPay;
            }
        })
        this.api.statistics.enterData({countyCode:0,searchStr:''})
        .then(res=>{
            if(res.data.code == 200){
                for(let i=0;i<res.data.data.length;i++){
                    this.restaurants.push({value:res.data.data[i].enterpriseName})
                }
            }
        })
    },
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')}
        this.http  = this.$store.state.http;                 // IP地址
        //接收主页单条数据
        this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));
        if(!this.AddData){
            // 直接输入网址 / 后退进入时没有待编辑的合同:回列表页,否则下面取 singers 就抛错
            return this.$router.replace({name:"SalesContract"});
        }
        if(!this.AddData.singers) this.AddData.singers = [];
        for(let i=0;i<this.AddData.singers.length;i++){
            if(!this.AddData.singers[i].isLaunchRole){
                this.AddData.launchType = this.AddData.singers[i].signRole;
                this.AddData.phone = this.AddData.singers[i].signPhone;
                this.AddData.etcName = this.AddData.singers[i].signEtcName;
                this.AddData.participants = this.AddData.singers[i].signName;
            }else{
                if(this.AddData.singers[i].signRole == 2){
                    this.AddData.signRole = this.AddData.singers[i].signRole;
                    this.numbers++; // 每个回显的签署方各用一个序号(原先共用同一个,再保存时会被合并)
                    this.AddData.domains.push(
                        {
                            id:this.numbers,
                            signRole:2,
                            label:'签署方姓名' + this.numbers,
                            value: this.AddData.singers[i].signName,
                            labelone:'删除',
                            type:'signName',
                            message:'签署方姓名'
                        },{
                            id:this.numbers,
                            signRole:2,
                            label:'签署方手机号' + this.numbers,
                            value: this.AddData.singers[i].signPhone,
                            labelone:'',
                            type:'signPhone',
                            message:'签署方手机号'
                        },{
                            id:this.numbers,
                            signRole:2,
                            label:'签署方企业名称' + this.numbers,
                            value: this.AddData.singers[i].signEtcName,
                            labelone:'',
                            type:'signEtcName',
                            message:'签署方企业名称'
                        }
                    );
                }else if(this.AddData.singers[i].signRole == 0){
                    // 个人签署方同样要回显:原先没有还原,再次保存时会被丢掉
                    this.AddData.domainsone.push(
                        {
                            id:++this.numbersone,
                            signRole:0,
                            label:'签署方姓名' + this.numbersone,
                            value: this.AddData.singers[i].signName,
                            labelone:'删除',
                            type:'signName',
                            message:'签署方姓名'
                        },{
                            id:this.numbersone,
                            signRole:0,
                            label:'签署方手机号' + this.numbersone,
                            value: this.AddData.singers[i].signPhone,
                            labelone:'',
                            type:'signPhone',
                            message:'签署方手机号'
                        }
                    );
                }
            }
        }
        if(this.AddData.custType == 1){
            this.placeblone = false;                        // 默认显示发货必填
            this.placebloneNo = true;                       // 发货非必填
            this.disabled = true;                           // 货品单位是否含税禁选
            this.disabledss = false;                        // 货品单位是否含税不禁选
            this.AddData.taxIncluded = false                // 选中不含税
        }else if(this.AddData.custType == 2){
            this.placeblone = true;                        // 默认显示发货必填
            this.placebloneNo = false;                       // 发货非必填
            this.disabled = false;                           // 货品单位是否含税禁选
            this.disabledss = false;                        // 货品单位是否含税不禁选
            this.AddData.taxIncluded = true                // 选中不含税
        }

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
                if(newarr[i].response.data[0].fileType == 1){
                    // 附件文件内容
                    this.enclosurefileList.push(newarr[i]);      
                }
            }
        })
        //基础维护数据  渲染扩展信息内容
        this.api.basics.all({template_id:1})
        .then(res=>{
            if(res.data.code == 200){
                var dictID =null;
                this.Extended = res.data.data;
                this.AddproBigType = res.data.data.pro_type; // 产品大类选择
                // 获取大类id
                for(let i=0;i<this.Extended.pro_type.length;i++){
                    if(this.Extended.pro_type[i].dvalue === this.AddData.proBigType * 1){
                        dictID = this.Extended.pro_type[i].dictId;
                    }
                }
                // 回显小类数据
                for(let k=0;k<this.Extended.pro__middle_type.length;k++){
                    var proArrss = this.Extended.pro__middle_type[k].dparentId.split(',');
                    for(let i=0;i<proArrss.length;i++){
                        if(proArrss[i] == ","){
                            this.splice(proArrss[i]);  
                        }
                        if(proArrss[i].indexOf(dictID) != -1){
                            this.SelectproSmaType.push(this.Extended.pro__middle_type[k]); // 产品小类
                            if(this.Extended.pro__middle_type[k].id == 46){
                                this.pIDshow = true; 
                            }else{
                                this.pIDshow = false; 
                            }
                        }
                    }
                }
            }
        })
    },
    methods: {
        // 在线支付提示
        isPayOnlinelogistics(){
            this.$confirm("开启在线支付后，合同将不可编辑，请谨慎操作！", "提示", {
                confirmButtonText: "开启",
                cancelButtonText: "关闭",
                type: "warning",
                closeOnPressEscape:false,
                closeOnClickModal:false,
                showClose:false,
            }).then(() => {
                this.AddData.isPayOnline = true;
            }).catch(() => {
                this.AddData.isPayOnline = false;
            });
        },
        /*
        电子签章
        */
        // 电子签章提示
        isEsignlogistics(){
            this.$confirm("开启电子签章后，合同将不可编辑，请谨慎操作！", "提示", {
                confirmButtonText: "开启",
                cancelButtonText: "关闭",
                type: "warning",
                closeOnPressEscape:false,
                closeOnClickModal:false,
                showClose:false,
            }).then(() => {
                this.AddData.isEsign = true;
            }).catch(() => {
                this.AddData.isEsign = false;
            });
        },
        // 企业动态增加input
        removeDomain(item) {
            // 签署方序号只增不减:回退后新加的签署方会与保留的重号,提交时被合并成一个
            var index = this.AddData.domains.indexOf(item);
            if (index !== -1) {
            this.AddData.domains.splice(index,3);
            }
        },
        addDomain() {  
            this.numbers++;
            this.AddData.domains.push(
                {
                    id:this.numbers,
                    signRole:2,
                    label:'签署方姓名' + this.numbers,
                    value: '',
                    labelone:'删除',
                    type:'signName',
                    message:'签署方姓名'
                },{
                    id:this.numbers,
                    signRole:2,
                    label:'签署方手机号' + this.numbers,
                    value: '',
                    labelone:'',
                    type:'signPhone',
                    message:'签署方手机号'
                },{
                    id:this.numbers,
                    signRole:2,
                    label:'签署方企业名称' + this.numbers,
                    value: '',
                    labelone:'',
                    type:'signEtcName',
                    message:'签署方企业名称'
                }
            );
        },
        // 个人动态增加input
        removeDomainone(item) {
            var index = this.AddData.domainsone.indexOf(item);
            if (index !== -1) {
            this.AddData.domainsone.splice(index,2);
            }
        },
        addDomainone() {
            // this.numbersone = this.numbersone - 1
            this.numbersone++;
            this.AddData.domainsone.push(
                {
                    id:this.numbersone,
                    signRole:0,
                    label:'签署方姓名' + this.numbersone,
                    value: '',
                    labelone:'删除',
                    type:'signName',
                    message:'签署方姓名'
                },{
                    id:this.numbersone,
                    signRole:0,
                    label:'签署方手机号' + this.numbersone,
                    value: '',
                    labelone:'',
                    type:'signPhone',
                    message:'签署方手机号'
                }
            );
        },
        // 参与人切换功能
        prologistics(id){
                if(id == 2){
                    this.AddData.domainsone = [];   // 个人 
                }
                if(id == 0){
                    this.AddData.domains = [];              // 企业
                }
        },

        // 总价金额计算方法
        focusData(){
            if(!this.AddData.unitPrice || !this.AddData.contractNum || this.AddData.unitPrice == 0 || this.AddData.contractNum == 0){
                this.AddData.sumPrice = 0;
            }else{
                this.AddData.sumPrice = ((this.AddData.unitPrice*1)*(this.AddData.contractNum*1)).toFixed(2)
            }
        },
        // 所属月份回显
        clickTime(){
            if(!this.AddData.signingTime) return; // 清空签定时间后取月份会抛错
            let strs = this.AddData.signingTime.substr(0,7);
            this.AddData.month = strs;
        },
        // 联想输入方法
        querySearch(queryString, cb) {
            var restaurants = this.restaurants;
            var results = queryString ? restaurants.filter(this.createFilter(queryString)) : restaurants;
            // 调用 callback 返回建议列表的数据
            cb(results);
        },
        createFilter(queryString) {
            return (restaurant) => {
            return (restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0);
            };
        },
        handleSelect(item) {
            this.CompanyNewData.enterpriseName = item.value;
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
        renderHeaderMethodsone(h,{column}){
            return h(
                'div',
                {
                    style:'display:flex;margin:auto;'
                },
                [
                    h('span',column.label),
                    h('promptMessages',{
                        props:{
                            messages:['相关操作在地址维护中操作']
                        }
                    })
                ]
            )
        },
        /* 
        编辑销售合同表单提交功能
        */ 
        // 总金额输入框限制
        CardNumber(val){
            if(val.indexOf('.') != -1){
                this.length = 12;
            }else{
                this.length = 10;
            }
        },
       // 大类目前的选中值
       getValue(){
           this.SelectproSmaType = [];
           var pro__middle_typeID = ""; 
            // 获取大类id
            for(let i=0;i<this.Extended.pro_type.length;i++){
                 this.SelectproSmaType = [];
                if(this.Extended.pro_type[i].dvalue == this.AddData.proBigType){
                    pro__middle_typeID = this.Extended.pro_type[i].dictId  // 获取产品大类id
                }
            }
            // 获取小类数据
            for(let k=0;k<this.Extended.pro__middle_type.length;k++){
                var proArrs = this.Extended.pro__middle_type[k].dparentId.split(',');
                for(let i=0;i<proArrs.length;i++){
                    if(proArrs[i] == ","){
                        this.splice(proArrs[i]);  
                    }
                    if(proArrs[i].indexOf(pro__middle_typeID) != -1){
                        this.AddData.proSmaType = '';  // 产品小类为空
                        this.SelectproSmaType.push(this.Extended.pro__middle_type[k]); // 产品小类选择
                    }
                }
            }
            // console.log(pro__middle_typeID)
            // console.log(this.SelectproSmaType)
            // console.log(this.Extended.pro__middle_type)
       },
       // 小类目前的选中值
       proSmaclick(pID){
            // console.log(pID); //获取的value
            // 选中小类其他显示产品名称
            if(pID == 20){
                this.pIDshow = true;
            }else{
                this.pIDshow = false;   
            }
            var proid = "";  // 下拉框选中value对应的id
            if(this.AddData.proSmaType == ''){
                this.AddData.list = [];
            }else{
                // 获取小类dictId
                for(let i=0;i<this.Extended.pro__middle_type.length;i++){
                    if(this.Extended.pro__middle_type[i].dvalue == pID){
                        this.AddData.list = [];
                        proid = this.Extended.pro__middle_type[i].dictId;
                    }
                }
                // console.log(proid)
                // 获取标的指标数据
                for(let k=0;k<this.Extended.pro__small_type.length;k++){
                    // console.log(this.Extended.pro__small_type)
                    var proArr = this.Extended.pro__small_type[k].dparentId.split(',');
                    for(let i=0;i<proArr.length;i++){
                        if(proArr[i] == ","){
                            this.splice(proArr[i]);  
                        }
                        if(proArr[i].indexOf(proid) != -1 && this.Extended.pro__small_type[k].dextend == true){
                            // 标的指标扩展字段
                            this.AddData.list = [];
                            this.AddData.list.push(
                                {
                                    dictId:this.Extended.pro__small_type[k].id,  // 字典id
                                    companyId:'', // 企业id
                                    dkey:this.Extended.pro__small_type[k].dkey,  // lab名称
                                    dvalue:'', // 内容
                                    dtype:'pro__small_type', // 扩展类型
                                }
                            ); // 产品小类选择
                        }
                    }
                }
            }
       },
        //编辑销售合同表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    let switchloan = null;
                    if(this.AddData.isEsign || this.AddData.isPayOnline){
                        switchloan = true;
                    }else if(this.AddData.isEsign && this.AddData.isPayOnline){
                        switchloan = true;
                    }else{
                        switchloan = false;
                    }
                    if(switchloan){
                        if(this.enclosurefileList.length != 0){
                            if(this.AddData.isEsign){
                                if(this.AddData.domains.length > 0 || this.AddData.domainsone.length > 0){
                                    if(this.AddData.launchType == 2){
                                        this.AddData.eqbFlow.orgInfo = {
                                            name:this.AddData.etcName,
                                        }
                                        this.AddData.eqbFlow.userInfo = {
                                            name:this.AddData.participants,
                                        }
                                    }else if(this.AddData.launchType == 0){
                                        this.AddData.eqbFlow.userInfo = {
                                            name:this.AddData.participants,
                                        }
                                    }
                                    var arrays = [];
                                    if(this.AddData.domains.length > 0){
                                        let a = {};
                                        for(let i of new Set(this.AddData.domains.map(d => d.id))){ // 只遍历还在的签署方序号
                                            a = {}
                                            for(let key in this.AddData.domains){
                                                if(this.AddData.domains[key].type == 'signName'){
                                                    this.AddData.domains[key].signName = this.AddData.domains[key].value
                                                }
                                                if(this.AddData.domains[key].type == 'signPhone'){
                                                    this.AddData.domains[key].signPhone = this.AddData.domains[key].value
                                                }
                                                if(this.AddData.domains[key].type == 'signEtcName'){
                                                    this.AddData.domains[key].signEtcName = this.AddData.domains[key].value
                                                }
                                                if(i === this.AddData.domains[key].id){
                                                    Object.assign(a,this.AddData.domains[key]);
                                                }
                                            }
                                            arrays.push(a);
                                        }
                                    }
                                    if(this.AddData.domainsone.length > 0){
                                        let b ={};
                                        for(let i of new Set(this.AddData.domainsone.map(d => d.id))){
                                            b = {}
                                            for(let key in this.AddData.domainsone){
                                                if(this.AddData.domainsone[key].type == 'signName'){
                                                    this.AddData.domainsone[key].signName = this.AddData.domainsone[key].value
                                                }
                                                if(this.AddData.domainsone[key].type == 'signPhone'){
                                                    this.AddData.domainsone[key].signPhone = this.AddData.domainsone[key].value
                                                }
                                                if(i === this.AddData.domainsone[key].id){
                                                    Object.assign(b,this.AddData.domainsone[key]);
                                                }
                                            }
                                            arrays.push(b)
                                        }
                                    }
                                    this.AddData.eqbFlow.singers = arrays;
                                    this.AddData.eqbFlow.contractName = this.enclosurefileList[0].name;
                                    let newobj = {
                                        signEtcName:this.AddData.etcName,
                                        signPhone:this.AddData.phone,
                                        signName:this.AddData.participants,
                                        signRole:this.AddData.launchType,
                                    }
                                    this.AddData.eqbFlow.singers.push(newobj);
                                    this.AddData.eqbFlow.topicName = this.AddData.contractName;
                                }else{
                                    this.$message({
                                        type: "warning",
                                        message: "签署方信息不能少于一个签署方!"
                                    });
                                    return; // 原先提示后仍继续提交,签章合同会没有签署方
                                }
                            }
                            var arr = [];
                            var cid;
                            cid = this.enclosurefileList.map((itme,index)=>{
                                arr.push({id:itme.response.data[0].id});
                                return itme.response.data[0].id;
                            }).join(',');
                            this.AddData.contractIds = cid;
                            this.AddData.agreementIds = '';
                            this.AddData.filelist = arr;
                            this.api.xsContract.edit(this.AddData)
                            .then(res=>{
                                if(res.data.code == 200){
                            this.delRemovedFiles();
                                    var that = this
                                    that.$message({
                                        type: "success",
                                        message: "编辑合同成功!",
                                        duration:500,  
                                        onClose(){
                                            sessionStorage.removeItem('queryAll');
                                            that.$router.push({name:"SalesContract"});
                                        }
                                    });
                                }
                            })
                        }else{
                            this.$message({
                                type: "warning",
                                message: "请上传合同文件!"
                            });
                        }
                    }else{
                        if(this.AddData.isEsign){
                            if(this.AddData.domains.length > 0 || this.AddData.domainsone.length > 0){
                                if(this.AddData.launchType == 2){
                                    this.AddData.eqbFlow.orgInfo = {
                                        name:this.AddData.etcName,
                                    }
                                    this.AddData.eqbFlow.userInfo = {
                                        name:this.AddData.participants,
                                    }
                                }else if(this.AddData.launchType == 0){
                                    this.AddData.eqbFlow.userInfo = {
                                        name:this.AddData.participants,
                                    }
                                }
                                var arrays = [];
                                if(this.AddData.domains.length > 0){
                                    let a = {};
                                    for(let i of new Set(this.AddData.domains.map(d => d.id))){ // 只遍历还在的签署方序号
                                        a = {}
                                        for(let key in this.AddData.domains){
                                            if(this.AddData.domains[key].type == 'signName'){
                                                this.AddData.domains[key].signName = this.AddData.domains[key].value
                                            }
                                            if(this.AddData.domains[key].type == 'signPhone'){
                                                this.AddData.domains[key].signPhone = this.AddData.domains[key].value
                                            }
                                            if(this.AddData.domains[key].type == 'signEtcName'){
                                                this.AddData.domains[key].signEtcName = this.AddData.domains[key].value
                                            }
                                            if(i === this.AddData.domains[key].id){
                                                Object.assign(a,this.AddData.domains[key]);
                                            }
                                        }
                                        arrays.push(a);
                                    }
                                }
                                if(this.AddData.domainsone.length > 0){
                                    let b ={};
                                    for(let i of new Set(this.AddData.domainsone.map(d => d.id))){
                                        b = {}
                                        for(let key in this.AddData.domainsone){
                                            if(this.AddData.domainsone[key].type == 'signName'){
                                                this.AddData.domainsone[key].signName = this.AddData.domainsone[key].value
                                            }
                                            if(this.AddData.domainsone[key].type == 'signPhone'){
                                                this.AddData.domainsone[key].signPhone = this.AddData.domainsone[key].value
                                            }
                                            if(i === this.AddData.domainsone[key].id){
                                                Object.assign(b,this.AddData.domainsone[key]);
                                            }
                                        }
                                        arrays.push(b)
                                    }
                                }
                                this.AddData.eqbFlow.singers = arrays;
                                this.AddData.eqbFlow.contractName = this.enclosurefileList[0].name;
                                let newobj = {
                                    signEtcName:this.AddData.etcName,
                                    signPhone:this.AddData.phone,
                                    signName:this.AddData.participants,
                                    signRole:this.AddData.launchType,
                                }
                                this.AddData.eqbFlow.singers.push(newobj);
                                this.AddData.eqbFlow.topicName = this.AddData.contractName;
                            }else{
                                this.$message({
                                    type: "warning",
                                    message: "签署方信息不能少于一个签署方!"
                                });
                            }
                        }
                        var arr = [];
                        var cid;
                        cid = this.enclosurefileList.map((itme,index)=>{
                            arr.push({id:itme.response.data[0].id});
                            return itme.response.data[0].id;
                        }).join(',');
                        this.AddData.contractIds = cid;
                        this.AddData.agreementIds = '';
                        this.AddData.filelist = arr;
                        this.api.xsContract.edit(this.AddData)
                        .then(res=>{
                            if(res.data.code == 200){
                        this.delRemovedFiles();
                                var that = this
                                that.$message({
                                    type: "success",
                                    message: "编辑合同成功!",
                                    duration:500,  
                                    onClose(){
                                        sessionStorage.removeItem('queryAll');
                                        that.$router.push({name:"SalesContract"});
                                    }
                                });
                            }
                        })
                    }
                }
            });
        },
        // 取消保存表单
        Addopen() {
            sessionStorage.removeItem('queryAll')
            this.$router.push({name:"SalesContract"})
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
        // 附件上传文件成功
        resfile(response, file, fileList){
            // 后端转 pdf 失败时回 code 200 + 空 data:拿不到附件号的文件不能留在列表里
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
        // 保存成功后才删编辑期间移除的回显附件(原先一移除就删,取消编辑后合同指向已删文件)
        delRemovedFiles(){
            (this.removedFileIds || []).forEach(id => this.api.addfile.del({id}));
            this.removedFileIds = [];
        },
        fileRemove(file,fileList){
            this.enclosurefileList = fileList; // 列表与界面一致:原先删除接口失败时保存仍会带上
            if(!this.fileObj[file.name]){
                // 回显的附件没有经过 before-upload:先记下,等保存成功再删
                this.removedFileIds = (this.removedFileIds || []).concat(file.response.data[0].id);
                return;
            }
            return  this.api.addfile.del({id:file.response.data[0].id})
                        .then(res=>{
                            if(res.data.code == 200){
                                this.$message({
                                    type: 'success',
                                    message: '删除文件成功!'
                                });
                                this.enclosurefileList = fileList;
                            }
                        })
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
            this.gid = row.id;
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
            this.AddData.supplierId = this.gid;                 // 客户单位id
            this.AddData.buyerUnit = this.dialogtextarea;       // 获取单位名称，回显到表单
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

        
        /* 
        发货地的选择功能
        */
        //发货地表格数据
        PlaceChoiceallsupplier(){
            //传参
            var listByPageDataone = {
                type:2,
                searchStr:this.PlaceChoicepurSearch,
                page:{
                    current:this.PlaceChoicecurrent,
                    size:this.PlaceChoicesize
                }
            }
            //渲染表格
            this.api.address.all(listByPageDataone)
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
        //发货地搜索功能
        PlaceChoicesearch() {
            if(this.PlaceChoicetableData != []){
                this.PlaceChoicetableData = [];
                this.PlaceChoiceallsupplier();
                this.PlaceChoicepaginationshow = true;
                this.PlaceChoicecurrent = 1
            }
        },
        //隐藏显示发货地地址新建
        PlaceChoiceCompanyshow(){
            this.PlaceChoiceshow = true
        },
        //删除地址
        PlaceChoiceDelete(data){
            this.$confirm('此操作将永久删除该文件, 是否继续?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.api.address.del(data.id)
                .then(res=>{
                    if(res.data.code == 200){
                        this.$message({
                            type: 'success',
                            message: '删除成功!'
                        });
                        this.PlaceChoicecurrent = 1
                        this.PlaceChoicetableData = [];
                        this.PlaceChoiceallsupplier();
                    }
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });          
            });
        },
        // 发货地选择表格
        PlaceChoicegetCurrentRow(row){ 
            //获取选中地址数据
            this.PlaceChoicetarea = row.addressName;
        },
        // 新建发货地提交表单
        PlaceChoiceNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.PlaceChoiceNewData.type = 2;
                    this.api.address.add(this.PlaceChoiceNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                            });
                            that.PlaceChoiceNewData = {};
                            that.PlaceChoiceshow = false;
                            that.PlaceChoiceallsupplier();
                            that.PlaceChoicepaginationshow = true;
                            
                        }
                    })
                }
            });
        },
        // 新建发货地取消提交表单
        PlaceChoiceNewShow() {
            this.PlaceChoiceNewData = {};
            this.PlaceChoiceshow = false
        },
        // 选择发货地弹窗
        PlaceChoice(){
            if(this.PlaceChoicetableData != []){
                this.PlaceChoicepaginationshow = true;
                this.PlaceChoiceallsupplier();
            }
            this.PlaceChoiceUnit=true; //开启弹窗
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.PlaceChoiceUnit == true){
                this.PlaceChoiceshow = false  
            }
        },
        // 选择发货地 关闭弹窗
        PlaceChoicehandleClose() {
            this.PlaceChoicepurSearch = "",
            this.Placeradio = "";  // 清空单选框
            this.PlaceChoicetarea = ""; // 清空显示地址内容框
            this.PlaceChoiceUnit = false; //关闭整个弹窗
        },
        // 发货地应用按钮 事件
        deliveryPlaceApplication(){
            this.AddData.deliveryPlace = this.PlaceChoicetarea; // 获取地址名称，回显到表单
            this.PlaceChoiceUnit = false; //关闭整个弹窗
            this.Placeradio = "";  // 清空单选框
            this.PlaceChoicetarea = ""; // 清空显示地址内容框
        },

        /* 
        收货地的选择功能
        */
        //收货地表格数据
        Placeharvestallsupplier(){
            //传参
            var listByPageDataone = {
                type:1,
                searchStr:this.PlaceharvestpurSearch,
                page:{
                    current:this.Placeharvestcurrent,
                    size:this.Placeharvestsize
                }
            }
            //渲染表格
            this.api.address.all(listByPageDataone)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.PlaceharvesttableData = res.data.data.records
                        this.Placeharvesttotal = res.data.data.total
                    }else{
                        this.Placeharvestpaginationshow = false;
                    }
                }
            })
        },
        //收货地搜索功能
        Placeharvestsearch() {
            if(this.PlaceharvesttableData != []){
                this.PlaceharvesttableData = [];
                this.Placeharvestallsupplier();
                this.Placeharvestpaginationshow = true;
                this.Placeharvestcurrent = 1
            }
        },
        //隐藏显示收货地地址新建
        PlaceharvestCompanyshow(){
            this.Placeharvestshow = true
        },
        // 收货地选择表格
        PlaceharvestgetCurrentRow(row){ 
            //获取选中地址数据
            this.Placeharvesttarea = row.addressName;
        },
        // 新建收货地提交表单
        PlaceharvestNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.PlaceharvestNewData.type = 1;
                    this.api.address.add(this.PlaceharvestNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                            });
                            that.Placeharvestshow = false;
                            that.Placeharvestallsupplier();
                            that.Placeharvestpaginationshow = true;
                            that.PlaceharvestNewData = {};
                        }
                    })
                }
            });
        },
        // 新建收发地取消提交表单
        PlaceharvestNewShow() {
            this.PlaceharvestNewData = {};
            this.Placeharvestshow = false
        },
        // 选择收货地弹窗
        Placeharvest(){
            if(this.PlaceharvesttableData != []){
                this.Placeharvestpaginationshow = true;
                this.Placeharvestallsupplier();
            }
            this.PlaceharvestUnit=true; //开启弹窗
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.PlaceharvestUnit == true){
                this.Placeharvestshow = false  
            }
        },
        // 选择收货地 关闭弹窗
        PlaceharvesthandleClose() {
            this.PlaceharvestpurSearch = "",
            this.Placeharvestradio = "";  // 清空单选框
            this.Placeharvesttarea = ""; // 清空显示地址内容框
            this.PlaceharvestUnit = false; //关闭整个弹窗
        },
        // 收货地应用按钮 事件
        PlaceharvestApplication(){
            this.AddData.receiptPlace = this.Placeharvesttarea; // 获取地址名称，回显到表单
            this.PlaceharvestUnit = false; //关闭整个弹窗
            this.Placeharvestradio = "";  // 清空单选框
            this.Placeharvesttarea = ""; // 清空显示地址内容框
        },
    },
    computed: {},
    watch: {}
};
</script>
<style lang="less" scoped>
// 电子签章增加 删除样式开始
.removeDomainoneBox{
    width: 60px;
    height: 32px;
    margin-left:4%;
    cursor: pointer;
    i{
        font-size: 16px;
        color: #CF1724;
        line-height: 32px;
        vertical-align: middle;
    }
    span{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
        font-weight: 600;
        vertical-align: middle;
    }
}
.zjcyfspan{
    width: 116px;
    height: 32px;
    text-align: center;
    cursor: pointer;
    border: 1px solid #CF1724;
    box-sizing: border-box;
    border-radius:5px;
    
    i{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
    }
    span{
        font-size: 14px;
        color: #CF1724;
        line-height: 32px;
        font-weight: 600;
    }
}
// 电子签章增加 删除样式结束
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
// 销售单位 购货单位 发 收货地弹窗样式开始
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
    // position: absolute;
    // top: 0;
    // right:0;
    width: 50%;
    height: 46px;
    padding: 9px 17px 5px 0;
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

// 销售单位 购货单位 发 收货地弹窗样式结束
</style>