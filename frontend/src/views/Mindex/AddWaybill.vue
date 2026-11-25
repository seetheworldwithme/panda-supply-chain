<template>
    <!-- 新建运单信息页面 -->
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
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">系统编码</label>
                            </el-col>
                            <el-form-item prop="systemCode">
                                <el-col :span="7">
                                    <el-input v-model="AddData.systemCode" disabled></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5" >
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">承运单位</label>
                            </el-col>
                            <el-form-item prop="carriageName">
                                <el-col :span="7">
                                    <el-select v-model="AddData.carriageName" style="width:100%;" @change="prologistics($event)">
                                        <el-option
                                        v-for="item in SelectcontractName"
                                        :key="item.value"
                                        :label="item.label"
                                        :value="item.value">
                                        </el-option>
                                    </el-select>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">标的物</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.subjectMatter" placeholder="标的物"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- 网络货运 -->
                        <div v-if="logistics">
                            <el-row>
                                <el-col :span="2">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">运单信息</label>
                                </el-col>
                                <el-form-item >
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in tagsone" :key="k" v-if="tagsone.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <span>{{tag.billNumber}}
                                                                <el-tooltip effect="dark" :content="tag.billNumber" placement="top">
                                                                    <span> {{tag.billNumber}}</span>
                                                                </el-tooltip>
                                                            </span>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in tagsone" :key="k" v-if="tagsone.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{tagsone[0].billNumber}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent">
                                                    共 <span style="color:#CF1724;"> {{tagsone.length}}</span> 条 
                                                </div> 
                                            </div>
                                            <div class="details-size" @click="PlaceChoiceone">
                                                <span>关 联</span>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>  
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">运输模式</label>
                                </el-col>
                                <el-form-item>
                                    <!-- 运输模式 -->
                                    <el-col :span="22">
                                        <el-radio v-model="AddData.transMethod" v-for="(val,kal) in Extended.transport_mode" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                    </el-col>
                                </el-form-item> 
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">承运平台</label>
                                </el-col>
                                <!-- 当选中网络货运显示单选框 -->
                                <el-form-item prop="carriagePlatform">
                                    <!-- 运输模式 -->
                                    <el-col :span="22">
                                        <el-radio v-model="AddData.carriagePlatform" v-for="(val,kal) in Extended.logistics_platform" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <!-- <el-row>
                                <el-col :span="5">
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">选择相应的网络货运平台进行数据外部导入对接</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="16">
                                        <div style="color:#CF1724;cursor: pointer;width:15%;">
                                            <i class="el-icon-download"></i><span> 网络货运平台导入</span>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row> -->
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">运输方式</label>
                                </el-col>
                                <el-form-item prop="transType">
                                    <!-- 运输方式 -->
                                    <el-col :span="22">
                                        <el-radio v-model="AddData.transType" v-for="(val,kal) in Extended.carrier" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                    </el-col>
                                </el-form-item> 
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">合同编号</label>
                                </el-col>
                                <el-form-item prop="ids">
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <span v-if="tag.sumPrice >= 0">{{tag.contractNum}} 吨 - 
                                                                <el-tooltip effect="dark" :content="'总价金额 ： ' + tag.sumPrice + ' 元'" placement="top">
                                                                    <span> {{tag.sumPrice|currency}} 元</span>
                                                                </el-tooltip>
                                                            </span>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span v-if="newtagesAll[0].sumPrice >= 0">{{newtagesAll[0].contractNum}} 吨 - {{newtagesAll[0].sumPrice|currency}} 元</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                            <div class="details-size" @click="contractNumberdialogSelect">
                                                <span style="color:#333;font-size:14px;">关 联</span>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>  
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">销售单位</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'销售单位：' + tag.salesUnit" placement="top">
                                                                <span> {{tag.salesUnit}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].salesUnit}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">购货单位</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'购货单位：' + tag.buyerUnit" placement="top">
                                                                <span> {{tag.buyerUnit}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].buyerUnit}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">发货地</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'发货地：' + tag.deliveryPlace" placement="top">
                                                                <span> {{tag.deliveryPlace}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].deliveryPlace}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">装货地</label>
                                </el-col>
                                <el-form-item prop="loadingAddress">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.loadingAddress" placeholder="装货地" disabled>
                                            <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="pretend">选 择</span>
                                        </el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">卸货地</label>
                                </el-col>
                                <el-form-item prop="unloadingPoint">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.unloadingPoint" placeholder="卸货地" disabled>
                                            <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="discharge">选 择</span>
                                        </el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">订单昵称</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.orderName" placeholder="请输入订单昵称"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">订单编号</label>
                                </el-col>
                                <el-form-item prop="orderNumber">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.orderNumber" placeholder="请输入订单编号"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">订单创建时间</label>
                                </el-col>
                                <el-form-item prop="orderCreateTime">
                                    <el-col :span="7">
                                        <el-date-picker
                                        style="width:100%;"
                                        v-model="AddData.orderCreateTime"
                                        type="datetime"
                                        value-format="yyyy-MM-dd HH:mm:ss"
                                        placeholder="选择订单创建时间">
                                        </el-date-picker>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">订单时间范围</label>
                                </el-col>
                                <el-form-item prop="oldStartTime">
                                    <el-col :span="7">
                                        <el-date-picker
                                        style="width:100%;"
                                        v-model="AddData.oldStartTime"
                                        type="datetimerange"
                                        value-format="yyyy-MM-dd HH:mm:ss"
                                        range-separator="至"
                                        start-placeholder="订单开始日期"
                                        end-placeholder="订单结束日期">
                                        </el-date-picker>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                        <!-- 物流公司 -->
                        <div v-if="bliLogistics">
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">物流公司</label>
                                </el-col>
                                <!-- 当选中物流公司显示输入框 -->
                                <el-form-item prop="logisticsCompany" key="1">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.logisticsCompany" placeholder="请输入物流公司名称"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">磅单</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in pronAll" :key="k" v-if="pronAll.length == 1">
                                                        <el-tag style="margin:4px;" size="small" type="info">
                                                                <el-tooltip effect="dark" :content="'净重 ： ' + tag.netWeight + ' 吨'" placement="top">
                                                                    <span>{{tag.licensePlate}} - {{tag.coalType}}</span>
                                                                </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in pronAll" :key="k" v-if="pronAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{pronAll[0].licensePlate}} - {{pronAll[0].coalType}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent">
                                                    共 <span style="color:#CF1724;"> {{pronAll.length}}</span> 条
                                                </div> 
                                            </div>
                                            <div class="details-size" @click="proNumberdialogSelect">
                                                <span style="color:#333;font-size:14px;">关 联</span>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :span="15">
                                        <span style="margin-left:12px;color:#333;">累计净重 <b style="color:#EA222E;">{{PronNum}}</b> 吨</span>
                                    </el-col>
                                </el-form-item>
                            </el-row> 
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">运输模式</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="22">
                                        <el-radio v-model="AddData.transMethod" v-for="(val,kal) in Extended.transport_mode" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                    </el-col>
                                </el-form-item> 
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">运输方式</label>
                                </el-col>
                                <el-form-item prop="transType" key="1">
                                    <el-col :span="22">
                                        <el-radio v-model="AddData.transType" v-for="(val,kal) in Extended.carrier" :key="kal" :label="val.dvalue">{{val.dkey}}</el-radio>
                                    </el-col>
                                </el-form-item> 
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">合同编号</label>
                                </el-col>
                                <el-form-item prop="ids" key="1">
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <span v-if="tag.sumPrice >= 0">{{tag.contractNum}} 吨 - 
                                                                <el-tooltip effect="dark" :content="'总价金额 ： ' + tag.sumPrice + ' 元'" placement="top">
                                                                    <span> {{tag.sumPrice|currency}} 元</span>
                                                                </el-tooltip>
                                                            </span>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span v-if="newtagesAll[0].sumPrice >= 0">{{newtagesAll[0].contractNum}} 吨 - {{newtagesAll[0].sumPrice|currency}} 元</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                            <div class="details-size" @click="contractNumberdialogSelect">
                                                <span style="color:#333;font-size:14px;">关 联</span>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>  
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">销售单位</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'销售单位：' + tag.salesUnit" placement="top">
                                                                <span> {{tag.salesUnit}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].salesUnit}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">购货单位</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'购货单位：' + tag.buyerUnit" placement="top">
                                                                <span> {{tag.buyerUnit}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].buyerUnit}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">发货地</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <div class="details">
                                            <div class="details-nav" style="width:100%;">
                                                <div class="details-tag" >
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                        <el-tag style="margin:4px;"
                                                        size="small"
                                                        type="info">
                                                            <el-tooltip effect="dark" :content="'发货地：' + tag.deliveryPlace" placement="top">
                                                                <span> {{tag.deliveryPlace}}</span>
                                                            </el-tooltip>
                                                        </el-tag>
                                                    </div>
                                                    <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                        <el-tag size="small" type="info" style="margin:0;">
                                                            <span>{{newtagesAll[0].deliveryPlace}}</span>
                                                        </el-tag>
                                                        <b> . . .</b>
                                                    </div>
                                                </div>
                                                <div class="details-conent" style="padding-right:10px;">
                                                    共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 条
                                                </div> 
                                            </div>
                                        </div>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">装货地</label>
                                </el-col>
                                <el-form-item key="1">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.loadingAddress" placeholder="装货地" disabled>
                                            <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="pretend">选 择</span>
                                        </el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">卸货地</label>
                                </el-col>
                                <el-form-item key="1">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.unloadingPoint" placeholder="卸货地" disabled>
                                            <span slot="suffix" style="color:#333;margin-right:18px;cursor: pointer;" @click="discharge">选 择</span>
                                        </el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">发货日期</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-date-picker
                                        v-model="AddData.Time"
                                        type="daterange"
                                        unlink-panels
                                        range-separator="至"
                                        start-placeholder="开始日期"
                                        end-placeholder="结束日期"
                                        value-format="yyyy-MM-dd"
                                        style="width:100%;">
                                        </el-date-picker>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5" >
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">收货日期</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-date-picker
                                        v-model="AddData.reTime"
                                        type="daterange"
                                        unlink-panels
                                        range-separator="至"
                                        start-placeholder="开始日期"
                                        end-placeholder="结束日期"
                                        value-format="yyyy-MM-dd"
                                        style="width:100%;">
                                        </el-date-picker>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                    </el-col>
                    <!-- 运单统计 -->
                    <el-col :span="24">
                        <div class="see-top">运单统计</div>
                        <!-- 网络货运 -->
                        <div v-if="logistics">
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">车数</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.vehicles" placeholder="车数" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">运费金额合计</label>
                                </el-col>
                                <el-form-item prop="freightAmount">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.freightAmount" placeholder="运费金额合计" @input="CardNumber" :maxlength="length" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <!-- <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">吨数合计</label>
                                </el-col>
                                <el-form-item prop="subjectCount">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.subjectCount" placeholder="吨数合计" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row> -->
                            <el-row>
                                <el-col :span="2.5">
                                    <span style="color:red;">* </span>
                                    <label style="color:#666;font-size:14px;line-height:30px;">原发吨数合计</label>
                                </el-col>
                                <el-form-item prop="originalTonnage">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.originalTonnage" placeholder="原发吨数合计" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收吨数合计</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.tonnagePaid" placeholder="实收吨数合计" disabled></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                        <!-- 物流公司 -->
                        <div v-if="bliLogistics">
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">车数</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.vehicles" placeholder="车数"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding:12px;">运费金额合计</label>
                                </el-col>
                                <el-form-item prop="freightAmount">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.freightAmount" placeholder="运费金额合计" @input="CardNumber" :maxlength="length"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">吨数合计</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.subjectCount" placeholder="吨数合计"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">原发吨数合计</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.originalTonnage" placeholder="原发吨数合计"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2.5">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;padding-left:12px;">实收吨数合计</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-input v-model="AddData.tonnagePaid" placeholder="实收吨数合计"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                        </div>
                    </el-col>
                    <!-- 附件信息 -->
                    <el-col :span="24">
                        <div class="see-top">附件信息</div>  
                        <!-- 上传合同附件 -->
                        <el-row>
                            <el-col :span="2.5" >
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

        <!-- 销售单位 购货单位弹窗开始...选择单位 -->
        <!-- <el-dialog
            :visible.sync="dialogSelectUnit"
            width="60%"
            :before-close="BuyerhandleClose">
            <template slot="title">
                <span>选择单位</span>
            </template>
            <el-input type="textarea" v-model="dialogtextarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textbutton">
                        <el-button size="small" v-show="salesUnitShow" @click="salesUnitApplication">应用</el-button>  
                        <el-button size="small" v-show="buyerUnitShow" @click="buyerUnitApplication">应用</el-button>
                    </div>
                    <div class="menu-list">
                        <div class="totalTab"
                        :class="{ active: item.id === sel }"
                        v-for="item in tabs"
                        :key="item.id"
                        @click="select(item)" style="cursor: pointer;">
                            <span>{{ item.label }}</span>
                            <div class="stripnone" :style="item.width"
                            :class="{ strip: item.id === sel }"></div>
                        </div>
                    </div>
                    <div class="text">
                        <div class="textAll" v-show="sel === 1">
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
                                <el-table-column prop="enterpriseName" label="单位名称" :render-header="renderHeaderMethods">
                                </el-table-column>
                                <el-table-column prop="createTime" label="创建时间">
                                </el-table-column>
                            </el-table>
                            <div class="pur-pages">
                                <el-pagination
                                v-show="paginationshow"
                                background
                                :page-size.sync="size"
                                :current-page.sync="current"
                                layout="prev, pager, next"
                                :total="totalone"
                                @current-change='allsupplierone'
                                @size-change='allsupplierone'>
                                </el-pagination>
                            </div>
                        </div>
                        <div class="textAll" v-show="sel === 2">
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
                        <div class="textAll" v-show="sel === 3">
                            <div class="textlast">
                                <el-input
                                @keyup.enter.native="searchtrr"
                                size="small"
                                placeholder="请输入内容"
                                prefix-icon="el-icon-search"
                                v-model="purSearchtrr"
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
                            <div class="pur-pages">
                                <el-pagination
                                v-show="paginationshowtrr"
                                background
                                :page-size.sync="sizetrr"
                                :current-page.sync="currenttrr"
                                layout="prev, pager, next"
                                :total="totaltrr"
                                @current-change='allsuppliertrr'
                                @size-change='allsuppliertrr'>
                                </el-pagination>
                            </div>
                        </div>
                    </div>
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
                            <el-col :span="6">
                                <el-form-item label=" 合作方 :" prop="partnerType" style="width:100%;">
                                    <el-select v-model="CompanyNewData.partnerType">
                                    <el-option label="供应商" value="1"></el-option>
                                    <el-option label="客户" value="2"></el-option>
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
        </el-dialog> -->
        <!-- 销售单位 购货单位弹窗结束...选择单位 -->

        <!-- 发货地址弹窗开始 -->
        <!-- <el-dialog
            :visible.sync="PlaceChoiceUnit"
            width="60%"
            :before-close="PlaceChoicehandleClose">
            <template slot="title">
                <div style="font-size:18px;">地址选择</div>
            </template>
            <el-input type="textarea" v-model="PlaceChoicetarea"></el-input>
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
                        <span @click="PlaceChoiceCompanyshow">+ 新建</span>
                    </div>
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
                        <el-table-column prop="addressName" label="地址名称" :render-header="renderHeaderMethodsone" >
                        </el-table-column>
                        <el-table-column prop="createTime" label="创建时间">
                        </el-table-column>
                    </el-table>
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
                        <el-button size="small" @click="deliveryPlaceApplication">应用</el-button> 
                    </div>
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
        </el-dialog> -->
        <!-- 发货地址弹窗结束 -->

        <!-- 合同关联弹窗开始 -->
        <el-dialog
            :visible.sync="PlaceChoicecontractNumber"
            width="60%"
            :before-close="contractNumberhandleClose">
            <template slot="title">
                <div style="font-size:18px;">关联合同</div>
            </template>
            <!-- 显示合同数据选择内容框 -->
            <div class="smbox">
                <el-tag
                size="small"
                type="info"
                v-for="(v,k) in commontags" :key="k"
                style="margin:5px 5px;color:#333;">
                <span v-if="v.contractNum">{{v.contractNum}} 吨 - {{v.sumPrice|currency}} 元</span>
                </el-tag>
            </div>
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
                    <div class="concentBoxtable">
                        <!-- 合同表格 -->
                        <el-table
                        border
                        height="100%"
                        ref="contractTable"
                        :data="contractNumbertableData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;"
                        @selection-change="contractgetCurrentRow"
                        :row-key="keyId">
                            <el-table-column type="selection" :reserve-selection="true"> </el-table-column>
                            <el-table-column prop="id" v-if="show">
                            </el-table-column>
                            <el-table-column prop="systemCode" label="合同编码">
                            </el-table-column>
                            <el-table-column prop="contractName" label="合同名称">
                            </el-table-column>
                            <el-table-column label="合同数量">
                                <template slot-scope="scope">
                                    <span>{{scope.row.contractNum}} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="总价金额">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.sumPrice != 0">{{scope.row.sumPrice|currency}} 元</span>
                                    <span v-else>0.00 元</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="contractpaginationshow"
                        background
                        :page-sizes="[50,100,300,500,]"
                        :page-size.sync="contractsize"
                        :current-page.sync="contractcurrent"
                        layout="total, sizes, prev, pager, next"
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
        <!-- 合同关联弹窗结束 -->

        <!-- 装货地址弹窗开始 -->
        <el-dialog
            :visible.sync="pretendUnit"
            width="60%"
            :before-close="pretendhandleClose">
            <template slot="title">
                <div style="font-size:18px;">地址选择</div>
            </template>
            <!-- 显示装货地址选择内容框 -->
            <el-input type="textarea" v-model="pretendtarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="pretendsearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="pretendpurSearch"
                        style="width:200px;float:right;"
                        >
                        </el-input>
                        <span @click="pretendCompanyshow">+ 新建</span>
                    </div>
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="pretendtableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="Placeradio" @change.native="pretendgetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="addressName" label="地址名称" :render-header="renderHeaderMethodsone" >
                        </el-table-column>
                        <el-table-column prop="createTime" label="创建时间">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="pretendpaginationshow"
                        background
                        :page-size.sync="pretendsize"
                        :current-page.sync="pretendcurrent"
                        layout="prev, pager, next"
                        :total="pretendtotal"
                        @current-change='pretendallsupplier'
                        @size-change='pretendallsupplier'>
                        </el-pagination>
                        <!-- 装货地应用按钮 -->
                        <el-button size="small" @click="pretendApplication">应用</el-button> 
                    </div>
                    <!-- 新建装货地 -->
                    <el-form :model="pretendNewData" :rules="pretendrels" ref="pretendNewData" label-width="80px" size="small" v-show="pretendshow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="8">
                                <el-form-item label=" 地 址 ：" prop="addressName">
                                    <el-input v-model="pretendNewData.addressName" placeholder="请输入地址名称" style="width:260px;"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="2" style="padding-left:25px;">
                                <el-button @click="pretendNewShow" size="small">取消</el-button>
                            </el-col>
                            <el-col :span="2" style="padding-left:15px;">
                                <el-button size="small" @click="pretendNewsubmitForm('pretendNewData')">保存</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 装货地址弹窗结束 -->

        <!-- 卸货地址弹窗开始 -->
        <el-dialog
            :visible.sync="dischargeUnit"
            width="60%"
            :before-close="dischargehandleClose">
            <template slot="title">
                <div style="font-size:18px;">地址选择</div>
            </template>
            <!-- 显示卸货地地址选择内容框 -->
            <el-input type="textarea" v-model="dischargetarea"></el-input>
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textlast">
                        <el-input
                        @keyup.enter.native="dischargesearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="dischargepurSearch"
                        style="width:200px;float:right;"
                        >
                        </el-input>
                        <span @click="dischargeCompanyshow">+ 新建</span>
                    </div>
                    <el-table
                    border
                    ref="PlacemultipleTable"
                    :data="dischargetableData"
                    :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                    style="width: 100%;">
                        <el-table-column width="33">
                            <template scope="scope">
                                <el-radio :label="scope.$index" v-model="dischargeradio" @change.native="dischargegetCurrentRow(scope.row)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</el-radio>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" v-if="show">
                        </el-table-column>
                        <el-table-column prop="addressName" label="地址名称" :render-header="renderHeaderMethodsone">
                        </el-table-column>
                        <el-table-column prop="createTime" label="创建时间">
                        </el-table-column>
                    </el-table>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        v-show="dischargepaginationshow"
                        background
                        :page-size.sync="dischargesize"
                        :current-page.sync="dischargecurrent"
                        layout="prev, pager, next"
                        :total="dischargetotal"
                        @current-change='dischargeallsupplier'
                        @size-change='dischargeallsupplier'>
                        </el-pagination>
                        <!-- 卸货地应用按钮 -->
                        <el-button size="small" @click="dischargeApplication">应用</el-button> 
                    </div>
                    <!-- 新建卸货地 -->
                    <el-form :model="dischargeNewData" :rules="dischargerels" ref="dischargeNewData" label-width="80px" size="small" v-show="dischargeshow">
                        <el-row :gutter="10" class="CompanyData">
                            <el-col :span="8">
                                <el-form-item label=" 地 址 ：" prop="addressName">
                                    <el-input v-model="dischargeNewData.addressName" placeholder="请输入地址名称" style="width:260px;"></el-input>
                                </el-form-item>
                            </el-col>
                            <el-col :span="2" style="padding-left:25px;">
                                <el-button @click="dischargeNewShow" size="small">取消</el-button>
                            </el-col>
                            <el-col :span="2" style="padding-left:15px;">
                                <el-button size="small" @click="dischargeNewsubmitForm('dischargeNewData')">保存</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                </div>
            </div>
        </el-dialog>
        <!-- 卸货地址弹窗结束 -->

        <!-- 磅单关联弹窗开始 -->
        <el-dialog
            :visible.sync="pronChoicecont"
            width="60%"
            :before-close="pronberhandleClose">
            <template slot="title">
                <div style="font-size:18px;">关联磅单</div>
            </template>
            <!-- 显示磅单数据选择内容框 -->
            <div class="smbox">
                <el-tag
                size="small"
                type="info"
                v-for="(v,k) in pronMontags" :key="k"
                style="margin:5px 5px;color:#333;">
                <span>{{v.licensePlate}} - {{v.coalType}}（{{v.netWeight}} 吨）</span>
                </el-tag>
            </div>
            <div class="Wdialog">
                <div class="concentBox">
                    <!-- 搜索框 -->
                    <div class="textlast">
                        <span style="margin-left:24px;color:#333;">累计净重 <b style="color:#EA222E;">{{PronNum}}</b> 吨</span>
                        <el-input
                        @keyup.enter.native="pronractsearch"
                        size="small"
                        placeholder="请输入内容"
                        prefix-icon="el-icon-search"
                        v-model="pronSearch"
                        style="width:200px;">
                        </el-input>
                    </div>
                    <div class="concentBoxtable">
                        <!-- 磅单表格 -->
                        <el-table
                        border
                        height="100%"
                        ref="prontractTable"
                        :data="pronTableData"
                        :header-cell-style="{ background: '#F6F6F6', color: '#333' }"
                        style="width: 100%;"
                        @selection-change="pronCurrentRow"
                        :row-key="proKeyId">
                            <el-table-column type="selection" :reserve-selection="true"> </el-table-column>
                            <el-table-column prop="id" v-if="proshow">
                            </el-table-column>
                            <el-table-column prop="licensePlate" label="车牌号">
                            </el-table-column>
                            <el-table-column prop="coalType" label="煤种">
                            </el-table-column>
                            <el-table-column label="皮重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.tare}} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="毛重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.roughWeight}} 吨</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="净重">
                                <template slot-scope="scope">
                                    <span>{{scope.row.netWeight}} 吨</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                    <!-- 分页 -->
                    <div class="pur-pages">
                        <el-pagination
                        background
                        :page-sizes="[50,100,300,500,]"
                        :page-size.sync="pronSize"
                        :current-page.sync="pronCurrent"
                        layout="total, sizes, prev, pager, next"
                        :total="pronTotal"
                        @current-change='pronNtractall'
                        @size-change='pronNtractall'>
                        </el-pagination>
                        <!-- 磅单关联应用按钮 -->
                        <el-button size="small" @click="pronApplication">应用</el-button>  
                    </div>
                </div>
            </div>
        </el-dialog>
        <!-- 磅单关联弹窗结束 -->
    </el-form>
</template>
<script>
export default {
    data() {
        return {
            /* 
            新建订单信息数据
            */
            tagsone: [],                            // 运单信息回显数据
            length:'',
            show:false,
            http:'',
            header:{},                              // 上传请求头
            fileObj:{},                             // 附件文件格式验证
            enclosurefileList: [],                  // 附件文件列表
            purTitle: "",                           // 标题
            logistics:true,                         // 承运平台默认显示
            bliLogistics:false,                     // 物流公司默认隐藏
            AddproBigType: [],                      // 产品大类选择
            Extended:'',                            // 扩展信息菜单
            SelectproSmaType: [],                   // 小类选择下拉框
            // placecontent:'', 
            // 承运单位选择下拉框
            SelectcontractName:[
                { 
                    value: 1,
                    label: '物流公司'
                },
                { 
                    value: 2,
                    label: '网络货运'
                },
            ],  
            // 表单验证
            rules: {
                orderName: [
                    { required: true, message: '请输入订单昵称', trigger: 'change' }
                ],
                orderCreateTime: [
                    { required: true, message: '请输入订单创建时间', trigger: 'change' }
                ],
                orderNumber: [
                    { required: true, message: '请输入订单编号', trigger: 'change' }
                ],
                oldStartTime: [
                    { required: true, message: '请输入订单时间范围', trigger: 'change' }
                ],
                systemCode: [
                    { required: true, message: '请输入系统编码', trigger: 'change' }
                ],
                carriageName: [
                    { required: true, message: '请选择承运单位', trigger: 'blur' }
                ],
                billNumber: [
                    { pattern: /(?=.*\d)(?=.*\W+)(?=.*[A-Z])(?!.*\n).*$/, message: '订单编号格式不正确', trigger: 'blur'}
                ],
                carriagePlatform: [
                    { required: true, message: '请选择承运平台', trigger: 'blur' }
                ],
                logisticsCompany: [
                    {required: true, message: '请输入物流公司名称', trigger: 'blur'},
                ],
                transType: [
                    { required: true, message: '请输入运输方式', trigger: 'blur' }
                ],
                contractNumber: [
                    { required: true, message: '请输入合同编号', trigger: 'change' }
                ],
                baseDeliveryPlace: [
                    { required: true, message: '请输入发货地', trigger: 'change' }
                ],
                ids: [
                    { required: true, message: '请关联合同编号', trigger: 'change' }
                ],
                // notBindOrder: [
                //     { required: true, message: '请关联磅单', trigger: 'change' }
                // ],
                subjectMatter: [
                    { required: true, message: '请输入标的物', trigger: 'blur' }
                ],
                subjectCount: [
                    { required: true, message: '请输入吨数', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '格式不正确', trigger: 'blur'}
                ],
                loadingAddress: [
                    { required: true, message: '请输入装货地', trigger: 'change' }
                ],
                freightAmount: [
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
                unloadingPoint: [
                    { required: true, message: '请输入卸货地', trigger: 'change' }
                ],
                Time: [
                    { required: true, message: '请输入发货日期', trigger: 'blur' }
                ],
                originalTonnage: [
                    { required: true, message: '请输入原发吨数', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '格式不正确', trigger: 'blur'}
                ],
                reTime: [
                    { required: true, message: '请输入收货日期', trigger: 'blur' }
                ],
                tonnagePaid: [
                    { required: true, message: '请输入实收吨数', trigger: 'blur' },
                    { pattern: /(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '格式不正确', trigger: 'blur'}
                ],
            },
            // 新建表单数据
            AddData: {
                systemCode:'',
                billNumber:'',
                transMethod:1,
                carriageName:2,
                carriagePlatform:'',
                logisticsCompany:'',
                transType:1,
                contractNumber:'',
                baseSalesUnit:'',
                baseBuyerUnit:'',
                baseDeliveryPlace:'',
                subjectMatter:'',
                subjectCount:0,
                logistics:[],
                ids:[],
                vehicles:'',
                loadingAddress:'',
                freightAmount:0,
                unloadingPoint:'',
                issuanceDate:'',
                tonnagePaid:0,
                originalTonnage:0,
                orderCreateTime:'', // 订单创建时间
                oldStartTime:[],    // 订单时间范围
                orderStartTime:'',  // 订单开始时间
                orderEndTime:'',    // 订单结束时间
                filelist:[]  // 合同 协议的集合id 
            }, 

            /* 
            合同关联数据列表
            */
            texttags:[],
            contractpaginationshow:true,            // 无数据隐藏分页(合同)
            contractcurrent:1,                      // 当前页数
            contractsize:50,                         // 显示条数
            contracttotal:0,                        // 总条数
            PlaceChoicecontractNumber:false,        // 选择合同弹窗
            contractNumberSearch:'',                // 合同搜索框
            contractNumbertableData:[],             // 选择合同表格数据
            // tags: [],                               // 合同集合回显数据
            moeny:0,                                // 合计金额
            commontags: [],                         // 合同集合回显数据
            newtagesAll:[],

            /* 
            装货地数据列表
            */
            pretendpaginationshow:true,             // 无数据隐藏分页
            pretendcurrent:1,                       // 当前页数
            pretendsize:4,                          // 每页显示条数
            pretendtotal:0,                         // 总条数
            pretendpurSearch:'',                    // 搜索
            pretendUnit:false,                      // 弹窗
            pretendtarea:'',                        // 显示内容框
            pretendshow:false,                      // 新建地址 显示隐藏
            pretendtableData:[],                    // 选择地址表格数据
            Placeradio:'',                          // 单选框
            // 新建表单
            pretendNewData:{
                addressName:'',
            },
            // 新建表单验证
            pretendrels: {
                addressName: [
                    { required: true, message: '请输入地址名称', trigger: 'blur' },
                ],
            },

            /* 
            卸货地数据列表
            */
            dischargepaginationshow:true,           // 无数据隐藏分页
            dischargecurrent:1,                     // 当前页数
            dischargesize:4,                        // 每页显示条数
            dischargetotal:0,                       // 总条数
            dischargepurSearch:'',                  // 搜索
            dischargeUnit:false,                    // 弹窗
            dischargetarea:'',                      // 显示收货内容框
            dischargeshow:false,                    // 新建地址 显示隐藏
            dischargetableData:[],                  // 选择卸货地地址表格数据
            dischargeradio:'',                      // 卸货地单选框
            // 新建卸货地表单
            dischargeNewData:{
                addressName:'',
            },
            // 新建卸货地表单验证
            dischargerels: {
                addressName: [
                    { required: true, message: '请输入地址名称', trigger: 'blur' },
                ],
            },

            /* 
            磅单关联数据列表
            */
            proshow:false,
            pronCurrent:1,                      // 当前页数
            pronSize:50,                        // 显示条数
            pronTotal:0,                        // 总条数
            pronChoicecont:false,               // 选择磅单弹窗
            pronSearch:'',                      // 磅单搜索框
            pronTableData:[],                   // 选择磅单表格数据
            pronMontags: [],                    // 磅单集合回显数据
            pronAll:[],   
            PronNum:0,                          // 弹框累计净重吨数 

            /* 
            销售单位 购货单位数据列表
            */
            // paginationshow:true,                    // 无数据隐藏分页(全部)
            // current:1,                              // 当前页数
            // size:4,                                 // 显示条数
            // totalone:0,                             // 总条数
            // purSearch:'',                           // 选择单位搜索框
            // paginationshowtwo:true,                 // 无数据隐藏分页(供应商)
            // currenttwo:1,                   
            // sizetwo:4,                      
            // totaltwo:0,                     
            // purSearchtwo:'',                
            // paginationshowtrr:true,                 // 无数据隐藏分页(客户)
            // currenttrr:1,                   
            // sizetrr:4,                      
            // totaltrr:0,                     
            // purSearchtrr:'',                
            // dialogSelectUnit:false,                 // 选择单位弹窗
            // dialogtextarea:'',                      // 显示单位内容框
            // sel: 1,                                 // 选项卡默认显示tab
            // tabs: [
            //     { label: "全部", id: 1 ,width:"width:35px;"},
            //     { label: "供应商", id: 2 ,width:"width:50px;"},
            //     { label: "客户", id: 3 ,width:"width:35px;"},
            // ],
            // CompanytableData:[],                    // 选择单位表格数据
            // radio:'',                               // 选择单位单选框
            // show:false,                             // 选择单位表格隐藏 id
            // CompanyDatashow:false,                  // 新建单位表单显示隐藏
            // salesUnitShow:false,                    // 销售单位应用按钮
            // buyerUnitShow:false,                    // 购货单位应用按钮
            // // 新建单位表单
            // CompanyNewData:{
            //     custType:'',
            //     partnerType:'',
            //     enterpriseName:'',
            // },
            // // 新建单位表单验证
            // rels: {
            //     custType: [
            //         { required: true, message: '请选择类型', trigger: 'blur' },
            //     ],
            //     partnerType: [
            //         { required: true, message: '请选择合作方', trigger: 'blur' },
            //     ],
            //     enterpriseName: [
            //         { required: true, message: '请输入单位名称', trigger: 'blur' },
            //     ],
            // },
            
            /* 
            发货地数据列表
            */
            // PlaceChoicepaginationshow:true,         // 无数据隐藏分页
            // PlaceChoicecurrent:1,                   // 当前页数
            // PlaceChoicesize:4,                      // 每页显示条数
            // PlaceChoicetotal:0,                     // 总条数
            // PlaceChoicepurSearch:'',                // 搜索
            // PlaceChoiceUnit:false,                  // 弹窗
            // PlaceChoicetarea:'',                    // 显示装货内容框
            // PlaceChoiceshow:false,                  // 新建地址表格显示隐藏
            // Placeradio:'',                          // 单选框
            // PlaceChoicetableData:[],                // 选择地址表格数据   
            // // 新建表单
            // PlaceChoiceNewData:{
            //     addressName:'',
            // },
            // // 新建表单验证
            // Placerels: {
            //     addressName: [
            //         { required: true, message: '请输入地址名称', trigger: 'blur' },
            //     ],
            // },
        };
    },
    created() {},
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 渲染系统编码
        this.api.logistics.code({codeType:1})
        .then(res=>{
            if(res.data.code == 200){
                this.AddData.systemCode = res.data.data;
            }
        });
        //基础维护数据
        this.api.basics.all({template_id:1})
        .then(res=>{
            if(res.data.code == 200){
                this.Extended = res.data.data;
                this.AddproBigType = res.data.data.pro_type;
            }
        });
        // 文件上传请求头
        this.header = {Authorization:sessionStorage.getItem('Authorization')};
        this.http  = this.$store.state.http;                 // IP地址
        // 回显关联信息
        if(JSON.parse(sessionStorage.getItem('queryAll')) != null){
            this.tagsone = JSON.parse(sessionStorage.getItem('queryAll')); 
        }else{
            this.tagsone = [];
        }
        // 运单统计
        if(this.tagsone.length > 0){
            let num = 0;
            let numone = 0;
            let numtwo = 0;
            let numstr = 0;
            // 车数
            this.AddData.vehicles = this.tagsone.length;
            for(let i=0;i<this.tagsone.length;i++){
                let obj = {id:this.tagsone[i].id};
                this.AddData.logistics.push(obj);
                num += this.tagsone[i].freightAmount*1;                 // 运费金额合计
                numone += this.tagsone[i].subjectCount*1;               // 吨数合计
                numtwo += this.tagsone[i].originalTonnage*1;            // 原发吨数合计
                numstr += this.tagsone[i].tonnagePaid*1;                // 实收吨数合计
            }
            this.AddData.freightAmount = num.toFixed(2)*1;              // 运费金额合计
            this.AddData.subjectCount = numone.toFixed(2)*1;            // 吨数合计
            this.AddData.originalTonnage = numtwo.toFixed(2)*1;         // 原发吨数合计
            this.AddData.tonnagePaid = numstr.toFixed(2)*1;             // 实收吨数合计
        }
    },
    methods: {
        // 关联运单列表信息按钮
        PlaceChoiceone(){
            sessionStorage.setItem('queryAll',JSON.stringify(this.tagsone));
            this.$router.push({name:'RelationBill'})
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
        新建订单信息表单提交功能
        */ 
        // 总金额输入框限制
        CardNumber(val){
            if(val.indexOf('.') != -1){
                this.length = 12;
            }else{
                this.length = 10;
            }
        },
       // 承运单位连级承运平台功能
       prologistics(id){
            // 物流公司 网络货运切换显示功能
            if(id == 2){
                this.newtagesAll = [];
                this.AddData.carriagePlatform = "";
                this.AddData.logisticsCompany = "";
                this.tagsone = [];
                this.tags = [];
                this.AddData.ids = [];
                this.AddData.vehicles = '';                  // 车数
                this.AddData.freightAmount = '';             // 运费金额
                this.AddData.subjectCount = '';              // 吨数
                this.AddData.originalTonnage = '';           // 原发吨数
                this.AddData.tonnagePaid = '';               // 实收吨数
                this.AddData.subjectMatter = '';             // 标的物
                this.AddData.loadingAddress = '';            // 装货地
                this.AddData.unloadingPoint = '';            // 卸货地
                this.AddData.Time = [];                      // 发货日期
                this.AddData.reTime = [];                    // 收货日期
                this.AddData.logistics = [];
                this.logistics = true;          // 承运平台显示
                this.bliLogistics = false;      // 物流公司隐藏
            }
            if(id == 1){
                this.newtagesAll = [];
                this.AddData.carriagePlatform = "";
                this.AddData.logisticsCompany = "";
                this.tagsone = [];
                this.tags = [];
                this.AddData.ids = [];
                this.AddData.vehicles = '';                  // 车数
                this.AddData.freightAmount = '';             // 运费金额
                this.AddData.subjectCount = '';              // 吨数
                this.AddData.originalTonnage = '';           // 原发吨数
                this.AddData.tonnagePaid = '';               // 实收吨数
                this.AddData.subjectMatter = '';             // 标的物
                this.AddData.loadingAddress = '';            // 装货地
                this.AddData.unloadingPoint = '';            // 卸货地
                this.AddData.logistics = [];
                this.logistics = false;         // 承运平台隐藏
                this.bliLogistics = true;       // 物流公司显示
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
                    pro__middle_typeID = this.Extended.pro_type[i].id                       // 获取产品大类id
                }
            }
            // 获取小类数据
            for(let k=0;k<this.Extended.pro__middle_type.length;k++){
                this.AddData.proSmaType = '';                                               // 产品小类为空
                if(this.Extended.pro__middle_type[k].dparentId == pro__middle_typeID){
                    this.SelectproSmaType.push(this.Extended.pro__middle_type[k]);          // 产品小类选择
                }
            }
            // console.log(this.SelectproSmaType)
            // console.log(pro__middle_typeID)
       },
       // 小类目前的选中值
       proSmaclick(pID){
            var proid = "";
            if(this.AddData.proSmaType == ''){
                this.AddData.list = [];
            }else{
                // 获取小类id
                for(let i=0;i<this.Extended.pro__middle_type.length;i++){
                    if(this.Extended.pro__middle_type[i].dvalue == pID){
                        this.AddData.list = [];
                        proid = this.Extended.pro__middle_type[i].id;
                    }
                }
                // console.log(proid)
                // 获取标的指标数据
                for(let k=0;k<this.Extended.pro__small_type.length;k++){
                    var proArr = this.Extended.pro__small_type[k].dparentId.split(',');
                    for(let i=0;i<proArr.length;i++){
                        if(proArr[i] == ","){
                            this.splice(proArr[i]);  
                        }
                        if(proArr[i].indexOf(proid) != -1 && this.Extended.pro__small_type[k].dextend == true){
                            // 标的指标扩展字段
                            this.AddData.list.push(
                                {
                                    dictId:this.Extended.pro__small_type[k].id,     // 字典id
                                    companyId:'',                                   // 企业id
                                    dkey:this.Extended.pro__small_type[k].dkey,     // lab名称
                                    dValue:'',                                      // name名称
                                    dType:'',                                       // 字典类型
                                }
                            );
                        }
                    }
                }
            }
       },
        // 新建订单信息表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    if(this.bliLogistics == true){
                        // 发货日期
                        if(this.AddData.Time){
                            this.AddData.issStartDate = this.AddData.Time[0];
                            this.AddData.issEndDate = this.AddData.Time[1];
                        }
                        // 收货日期
                        if(this.AddData.reTime){
                            this.AddData.recStartDate = this.AddData.reTime[0];
                            this.AddData.recEndDate = this.AddData.reTime[1];
                        }
                        // 文件上传
                        var arr = [];
                        this.enclosurefileList.map((itme,index)=>{
                            arr.push({id:itme.response.data[0].id});
                        }).join(',')
                        this.AddData.filelist = arr;
                        this.api.logistics.add(this.AddData)
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
                                        that.$router.push({name:"Waybill"})
                                    }
                                });
                            }
                        })
                    }else if(this.logistics){
                        if(this.AddData.logistics.length > 0){
                            // console.log(this.AddData);
                            this.AddData.orderStartTime = this.AddData.oldStartTime[0];
                            this.AddData.orderEndTime = this.AddData.oldStartTime[1];
                            // 文件上传 
                            var arr = [];
                            this.enclosurefileList.map((itme,index)=>{
                                arr.push({id:itme.response.data[0].id});
                            }).join(',')
                            this.AddData.filelist = arr;
                            this.api.logistics.add(this.AddData)
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
                                            that.$router.push({name:"Waybill"})
                                        }
                                    });
                                }
                            })
                        }else{
                            this.$message({
                                type: "warning",
                                message: "请关联运单信息!",
                            });
                        }

                    }

                }
            });
        },
        // 取消保存表单
        Addopen() {
            sessionStorage.removeItem('queryAll');
            this.$router.push({name:"Waybill"})
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
            //             sessionStorage.removeItem('queryAll');
            //             that.$router.push({name:"Waybill"})
            //         }
            //     });
            // })
            // .catch(() => {
            //     sessionStorage.removeItem('queryAll');
            //     this.$router.push({name:"Waybill"})
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
        发货地的选择功能
        */
        // // 发货地表格数据
        // PlaceChoiceallsupplier(){
        //     // 传参
        //     var listByPageDataone = {
        //         type:2,
        //         searchStr:this.PlaceChoicepurSearch,
        //         page:{
        //             current:this.PlaceChoicecurrent,
        //             size:this.PlaceChoicesize
        //         }
        //     }
        //     // 渲染表格
        //     this.api.address.all(listByPageDataone)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             if(res.data.data.records != ""){
        //                 this.PlaceChoicetableData = res.data.data.records
        //                 this.PlaceChoicetotal = res.data.data.total
        //             }else{
        //                 this.PlaceChoicepaginationshow = false;
        //             }
        //         }
        //     })
        // },
        // // 发货地搜索功能
        // PlaceChoicesearch() {
        //     if(this.PlaceChoicetableData != []){
        //         this.PlaceChoicetableData = [];
        //         this.PlaceChoiceallsupplier();
        //         this.PlaceChoicepaginationshow = true;
        //         this.PlaceChoicecurrent = 1
        //     }
        // },
        // // 隐藏显示发货地地址新建
        // PlaceChoiceCompanyshow(){
        //     this.PlaceChoiceshow = true
        // },
        // // 发货地选择表格
        // PlaceChoicegetCurrentRow(row){
        //     this.PlaceChoicetarea = row.addressName;        // 获取选中地址数据
        // },
        // // 新建发货地提交表单
        // PlaceChoiceNewsubmitForm(formName) {
        //     this.$refs[formName].validate((valid) => {
        //         if (valid) {
        //             this.PlaceChoiceNewData.type = 2;
        //             this.api.address.add(this.PlaceChoiceNewData)
        //             .then(res=>{
        //                 if(res.data.code == 200){
        //                     var that = this
        //                     that.$message({
        //                         type: "success",
        //                         message: "保存成功!",
        //                     });
        //                     that.PlaceChoiceNewData = {};
        //                     that.PlaceChoiceshow = false;
        //                     that.PlaceChoiceallsupplier();
        //                     that.PlaceChoicepaginationshow = true;
        //                 }
        //             })
        //         }
        //     });
        // },
        // // 新建发货地取消提交表单
        // PlaceChoiceNewShow() {
        //     this.PlaceChoiceNewData = {};
        //     this.PlaceChoiceshow = false
        // },
        // // 选择发货地弹窗
        // PlaceChoice(){
        //     if(this.PlaceChoicetableData != []){
        //         this.PlaceChoicepaginationshow = true;
        //         this.PlaceChoiceallsupplier();
        //     }
        //     this.PlaceChoiceUnit=true;        
        //     // 判断，如果弹窗为显示，新建表单则为隐藏
        //     if(this.PlaceChoiceUnit == true){
        //         this.PlaceChoiceshow = false  
        //     }
        // },
        // // 选择发货地 关闭弹窗
        // PlaceChoicehandleClose() {
        //     this.PlaceChoicepurSearch = "",
        //     this.Placeradio = "";               // 清空单选框
        //     this.PlaceChoicetarea = "";         // 清空显示地址内容框
        //     this.PlaceChoiceUnit = false;       // 关闭整个弹窗
        // },
        // // 发货地应用按钮 事件
        // deliveryPlaceApplication(){
        //     this.AddData.baseDeliveryPlace = this.PlaceChoicetarea;     // 获取地址名称,回显到表单
        //     this.PlaceChoiceUnit = false;                               // 关闭整个弹窗
        //     this.Placeradio = "";                                       // 清空单选框
        //     this.PlaceChoicetarea = "";                                 // 清空显示地址内容框
        // },

        /* 
        销售单位 购货单位的选择功能
        1. 选择销售单位显示销售单位的应用按钮
        2. 选择购货单位显示销售单位的应用按钮
        */
        // // 销售单位 购货单位表格数据(全部)
        // allsupplierone(){
        //     // 传参
        //     var listByPageDataone = {
        //         searchStr:this.purSearch,
        //         page:{
        //             current:this.current,
        //             size:this.size
        //         }
        //     }
        //     // 渲染表格
        //     this.api.supplier.all(listByPageDataone)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             if(res.data.data.records != ""){
        //                 this.CompanytableData = res.data.data.records
        //                 this.totalone = res.data.data.total
        //             }else{
        //                 this.paginationshow = false;
        //             }
        //         }
        //     })
        // },
        // // 销售 购货单位搜索功能(全部)
        // search() {
        //     this.CompanytableData = [];
        //     this.allsupplierone();
        //     this.current = 1;
        // },
        // // 销售单位 购货单位表格数据(供应商)
        // allsuppliertwo(){
        //     // 传参
        //     var listByPageDatatwo = {
        //         partnerType:1,
        //         searchStr:this.purSearchtwo,
        //         page:{
        //             current:this.currenttwo,
        //             size:this.sizetwo
        //         }
        //     }
        //     // 渲染表格
        //     this.api.supplier.all(listByPageDatatwo)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             if(res.data.data.records != ""){
        //                 this.CompanytableData = res.data.data.records
        //                 this.totaltwo = res.data.data.total
        //             }else{
        //                 this.paginationshowtwo = false;
        //             }
        //         }
        //     })
        // },
        // // 销售 购货单位搜索功能(供应商)
        // searchtwo() {
        //     this.CompanytableData = [];
        //     this.allsuppliertwo();
        //     this.currenttwo = 1;
        // },
        // // 销售单位 购货单位表格数据(客户)
        // allsuppliertrr(){
        //     // 传参
        //     var listByPageDatatrr = {
        //         partnerType:2,
        //         searchStr:this.purSearchtrr,
        //         page:{
        //             current:this.currenttrr,
        //             size:this.sizetrr
        //         }
        //     }
        //     // 渲染表格
        //     this.api.supplier.all(listByPageDatatrr)
        //     .then(res=>{
        //         if(res.data.code == 200){
        //             if(res.data.data.records != ""){
        //                 this.CompanytableData = res.data.data.records
        //                 this.totaltrr = res.data.data.total
        //             }else{
        //                 this.paginationshowtrr = false;
        //             }
        //         }
        //     })
        // },
        // // 销售 购货单位搜索功能(客户)
        // searchtrr() {
        //     this.CompanytableData = [];
        //     this.allsuppliertrr();
        //     this.currenttrr = 1;
        // },
        // // 选项卡
        // select(item) {
        //     this.sel = item.id;
        //     if(this.sel == 1){
        //         this.radio = "";
        //         this.current = 1;
        //         this.CompanytableData = [];
        //         this.allsupplierone();
        //     }else if(this.sel == 2){
        //         this.radio = "";
        //         this.currenttwo = 1;
        //         if(this.CompanytableData !=[]){
        //             this.CompanytableData = [];
        //             this.allsuppliertwo()
        //             this.paginationshowtwo = true;
        //         }else{
        //             this.allsuppliertwo()
        //         }
        //     }else if(this.sel == 3){
        //         this.radio = "";
        //         this.currentrr = 1;
        //         if(this.CompanytableData !=[]){
        //             this.CompanytableData = [];
        //             this.allsuppliertrr()
        //             this.paginationshowtrr = true;
        //         }else{
        //             this.allsuppliertrr()
        //         }
        //     }
        // },
        // // 销售 购货 选择单位表格
        // getCurrentRow(row){
        //     this.dialogtextarea = row.enterpriseName;   // 获取选中单位名称数据
        // },
        // // 隐藏显示新建
        // Companyshow(){
        //     this.CompanyDatashow = true
        // },
        // // 新建单位提交表单
        // CompanyNewsubmitForm(formName) {
        //     this.$refs[formName].validate((valid) => {
        //         if (valid) {
        //             this.api.supplier.add(this.CompanyNewData)
        //             .then(res=>{
        //                 if(res.data.code == 200){
        //                     var that = this
        //                     that.$message({
        //                         type: "success",
        //                         message: "保存成功!",
        //                         duration:500,  
        //                         onClose(){
        //                             that.sel = 1;
        //                             that.paginationshow = true;
        //                             that.allsupplierone();
        //                             that.CompanyDatashow = false;
        //                             that.CompanyNewData = {};
        //                         }
        //                     });
        //                 }
        //             })
        //         }
        //     });
        // },
        // // 新建单位取消提交表单
        // CompanyNewShow() {
        //     this.CompanyNewData = {};
        //     this.CompanyDatashow = false
        // },
        // // 选择销售单位弹窗
        // dialogSelect(){
        //     this.allsupplierone();          // 渲染表格
        //     this.dialogSelectUnit=true;     // 开启弹窗
        //     this.salesUnitShow = true;      // 显示销售单位应用按钮
        //     // 判断，如果弹窗为显示，新建表单则为隐藏
        //     if(this.dialogSelectUnit == true){
        //         this.CompanyDatashow = false
        //     }
        // },
        // // 销售单位应用按钮事件
        // salesUnitApplication(){
        //     this.AddData.baseSalesUnit = this.dialogtextarea;   // 获取单位名称，回显到表单
        //     this.dialogSelectUnit = false;                      // 关闭整个弹窗
        //     this.salesUnitShow = false;                         // 隐藏销售单位应用按钮
        //     this.radio = "";                                    // 清空单选框
        //     this.dialogtextarea = "";                           // 清空显示单位内容框
        // },
        // // 选择销货 购货单位 关闭弹窗
        // BuyerhandleClose() {
        //     this.radio = "";                // 清空单选框
        //     this.dialogtextarea = "";       // 清空显示单位内容框
        //     this.dialogSelectUnit = false;  // 关闭整个弹窗
        //     this.salesUnitShow = false;     // 销售单位应用按钮隐藏
        //     this.buyerUnitShow = false;     // 购货单位应用按钮隐藏
        // },
        // // 选择购货单位弹窗
        // Buyer(){
        //     this.dialogSelectUnit=true;     // 开启弹窗
        //     this.buyerUnitShow = true;      // 显示购货单位应用按钮
        //     // 判断，如果弹窗为显示，新建表单则为隐藏
        //     if(this.dialogSelectUnit == true){
        //         this.CompanyDatashow = false
        //     }
        // },
        // // 购货单位应用按钮 事件
        // buyerUnitApplication(){
        //     this.AddData.baseBuyerUnit = this.dialogtextarea;   // 获取单位名称，回显到表单
        //     this.dialogSelectUnit = false;                      // 关闭整个弹窗
        //     this.buyerUnitShow = false;                         // 隐藏销售单位应用按钮
        //     this.radio = "";                                    // 清空单选框
        //     this.dialogtextarea = "";                           // 清空显示单位内容框
        // },

        /* 
        磅单关联的选择功能
        */
        // 表格数据
        pronNtractall(){
            // 传参
            var listByPageData = {
                orderId:this.AddData.id?this.AddData.id:null,
                notBindOrder:true,
                searchStr:this.pronSearch,
                page:{
                    current:this.pronCurrent,
                    size:this.pronSize
                }
            }
            // 渲染表格
            this.api.Weighing.all(listByPageData)
            .then(res=>{
                if(res.data.code == 200){
                    this.pronTableData = res.data.data.records;
                    this.pronTotal = res.data.data.total;
                    this.$refs.prontractTable.clearSelection();
                    if (this.pronAll.length > 0) {
                        this.pronMontags = this.pronAll;
                        for (var i = 0; i < this.pronAll.length; i++) {
                            for (var j = 0; j < this.pronTableData.length; j++) {
                                if (this.pronAll[i].id == this.pronTableData[j].id) {   
                                    // items里面包含有的都回显打勾                 
                                    this.$refs.prontractTable.toggleRowSelection(this.pronTableData[j],true);
                                }
                            }
                        }
                    }
                }
            })
        },
        proKeyId(row){
            return row.id;
        },
        // 磅单关联搜索功能
        pronractsearch() {
           this.pronNtractall();
           this.pronCurrent = 1;
        },
        // 磅单车俩选择表格
        pronCurrentRow(row){
            this.pronMontags = row;                         // 获取选中单位数据集合
            let pronnumber = 0;
            row.map((item)=>{
                pronnumber += (item.netWeight*1)
            });
            this.PronNum = pronnumber.toFixed(2);
        },
        // 磅单数据弹窗
        proNumberdialogSelect(){
            this.pronChoicecont = true;                     // 开启弹窗
            this.pronNtractall();
        },
        // 选择磅单关联关闭弹窗
        pronberhandleClose() {
            this.pronSearch = '';
            this.pronChoicecont = false;                    // 关闭整个弹窗
        },
        // 磅单数据应用按钮事件
        pronApplication(){
            this.pronAll = this.pronMontags;                // 数据回显到关联框
            this.AddData.weighingList = [];
            for(let i=0;i<this.pronAll.length;i++){
                this.AddData.weighingList.push({id:this.pronAll[i].id})   // 组装数据,添加到提交表单
            }
            this.pronChoicecont = false;                    // 关闭整个弹窗
        },

        /* 
        合同关联的选择功能
        */
        // 表格数据
        contractall(){
            // 传参
            var listByPageData = {
                orderId:this.AddData.id?this.AddData.id:null,
                notBindOrder:true,
                contractType:'',
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
                    if (res.data.data.records.length > 0) {
                        this.contractNumbertableData = res.data.data.records;
                        this.contracttotal = res.data.data.total;
                        this.$refs.contractTable.clearSelection();
                        if (this.newtagesAll.length > 0) {
                            this.commontags = this.newtagesAll
                            for (var i = 0; i < this.newtagesAll.length; i++) {
                                for (var j = 0; j < this.contractNumbertableData.length; j++) {
                                    
                                    if (this.newtagesAll[i].id == this.contractNumbertableData[j].id) {//items里面包含有的都回显打勾                 
                                        this.$refs.contractTable.toggleRowSelection(this.contractNumbertableData[j],true);
                                    }
                                }
                            }
                        }
                    }else{
                        this.contractpaginationshow = false;
                    }
                }
            })
        },
        keyId(row){
            return row.id;
        },
        // 合同关联搜索功能
        contractsearch() {
           this.contractcurrent = 1;
           this.contractall();
        },
        // 合同编号选择表格
        contractgetCurrentRow(row){
            this.commontags = row;                         // 获取选中单位数据集合
        },
        // 合同数据弹窗
        contractNumberdialogSelect(){
            this.PlaceChoicecontractNumber = true; //开启弹窗
            this.contractall();
        },
        // 选择合同关联关闭弹窗
        contractNumberhandleClose() {
            this.contractNumberSearch = '';
            this.PlaceChoicecontractNumber = false;         // 关闭整个弹窗
        },
        // 合同数据应用按钮 事件
        contractNumberApplication(){
            this.newtagesAll = this.commontags;
            this.tags = this.commontags;
            this.AddData.ids = [];
            for(let i=0;i<this.newtagesAll.length;i++){
                this.AddData.ids.push(this.newtagesAll[i].id)   // 获取地址名称,回显到表单
            }
            this.PlaceChoicecontractNumber = false;         // 关闭整个弹窗
        },

        /* 
        装货地的选择功能
        */
        // 装货地表格数据
        pretendallsupplier(){
            // 传参
            var listByPageDataone = {
                type:3,
                searchStr:this.pretendpurSearch,
                page:{
                    current:this.pretendcurrent,
                    size:this.pretendsize
                }
            }
            // 渲染表格
            this.api.address.all(listByPageDataone)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.pretendtableData = res.data.data.records
                        this.pretendtotal = res.data.data.total
                    }else{
                        this.pretendpaginationshow = false;
                    }
                }
            })
        },
        // 装货地搜索功能
        pretendsearch() {
            if(this.pretendtableData != []){
                this.pretendtableData = [];
                this.pretendallsupplier();
                this.pretendpaginationshow = true;
                this.pretendcurrent = 1
            }
        },
        // 隐藏显示装货地址新建
        pretendCompanyshow(){
            this.pretendshow = true
        },
        // 装货地选择表格
        pretendgetCurrentRow(row){ 
            this.pretendtarea = row.addressName;       //获取选中地址数据
        },
        // 新建装货地提交表单
        pretendNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.pretendNewData.type = 3;
                    this.api.address.add(this.pretendNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                            });
                            that.pretendshow = false;
                            that.pretendallsupplier();
                            that.pretendpaginationshow = true;
                            that.pretendNewData = {};
                        }
                    })
                }
            });
        },
        // 新建装货地取消提交表单
        pretendNewShow() {
            this.pretendNewData = {};
            this.pretendshow = false
        },
        // 选择装货地弹窗
        pretend(){
            if(this.pretendtableData != []){
                this.pretendpaginationshow = true;
                this.pretendallsupplier();
            }
            this.pretendUnit=true;
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.pretendUnit == true){
                this.pretendshow = false  
            }
        },
        // 选择装货地 关闭弹窗
        pretendhandleClose() {
            this.pretendpurSearch = "",
            this.Placeradio = "";        // 清空单选框
            this.pretendtarea = "";        // 清空显示地址内容框
            this.pretendUnit = false;      // 关闭整个弹窗
        },
        // 装货地应用按钮 事件
        pretendApplication(){
            this.AddData.loadingAddress = this.pretendtarea;    // 获取地址名称，回显到表单
            this.pretendUnit = false;                           // 关闭整个弹窗
            this.Placeradio = "";                             // 清空单选框
            this.pretendtarea = "";                             // 清空显示地址内容框
        },

        /* 
        卸货地的选择功能
        */
        // 卸货地表格数据
        dischargeallsupplier(){
            // 传参
            var listByPageDataone = {
                type:4,
                searchStr:this.dischargepurSearch,
                page:{
                    current:this.dischargecurrent,
                    size:this.dischargesize
                }
            }
            // 渲染表格
            this.api.address.all(listByPageDataone)
            .then(res=>{
                if(res.data.code == 200){
                    if(res.data.data.records != ""){
                        this.dischargetableData = res.data.data.records
                        this.dischargetotal = res.data.data.total
                    }else{
                        this.dischargepaginationshow = false;
                    }
                }
            })
        },
        // 卸货地搜索功能
        dischargesearch() {
            if(this.dischargetableData != []){
                this.dischargetableData = [];
                this.dischargeallsupplier();
                this.dischargepaginationshow = true;
                this.dischargecurrent = 1
            }
        },
        // 隐藏显示卸货地址新建
        dischargeCompanyshow(){
            this.dischargeshow = true
        },
        // 卸货地选择表格
        dischargegetCurrentRow(row){ 
            this.dischargetarea = row.addressName;       // 获取选中地址数据
        },
        // 新建卸货地提交表单
        dischargeNewsubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.dischargeNewData.type = 4;
                    this.api.address.add(this.dischargeNewData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                            });
                            that.dischargeshow = false;
                            that.dischargeallsupplier();
                            that.dischargepaginationshow = true;
                            that.dischargeNewData = {};
                        }
                    })
                }
            });
        },
        // 新建卸货地取消提交表单
        dischargeNewShow() {
            this.dischargeNewData = {};
            this.dischargeshow = false
        },
        // 选择卸货地弹窗
        discharge(){
            if(this.dischargetableData != []){
                this.dischargepaginationshow = true;
                this.dischargeallsupplier();
            }
            this.dischargeUnit=true;
            // 判断，如果弹窗为显示，新建表单则为隐藏
            if(this.dischargeUnit == true){
                this.dischargeshow = false  
            }
        },
        // 选择卸货地 关闭弹窗
        dischargehandleClose() {
            this.dischargepurSearch = "",
            this.dischargeradio = "";        // 清空单选框
            this.dischargetarea = "";        // 清空显示地址内容框
            this.dischargeUnit = false;      // 关闭整个弹窗
        },
        // 卸货地应用按钮 事件
        dischargeApplication(){
            this.AddData.unloadingPoint = this.dischargetarea;      // 获取地址名称，回显到表单
            this.dischargeUnit = false;                             // 关闭整个弹窗
            this.dischargeradio = "";                               // 清空单选框
            this.dischargetarea = "";                               // 清空显示地址内容框
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
            width:50%;
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
// 销售单位 购货单位 发 装货地点弹窗样式开始
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
    .concentBoxtable{
        height:330px;
    }
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
    margin:9px 0 0 15px;
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


// 销售单位 购货单位 发 装货地点弹窗样式结束
</style>