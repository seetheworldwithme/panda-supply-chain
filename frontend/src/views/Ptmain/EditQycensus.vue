<template>
    <!-- 编辑企业统计页面 -->
    <el-form :model="AddData" :rules="rules" ref="AddData" size="small" label-width="130px"  v-loading="uploading"
        element-loading-text="文件正在识别中，请勿关闭！"
        element-loading-spinner="el-icon-loading"
        element-loading-background="rgba(0, 0, 0, 0.8)">
        <!-- 头部开始 -->
        <el-row class="pur-top">
        <el-col :span="10"><span class="pur-size">{{ purTitle }}</span></el-col>
        <el-col :span="14" class="pur-right">
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
                    <el-collapse v-model="activeNames" accordion>
                        <el-collapse-item title="基础信息" name="1">
                            <!-- 基础信息 -->
                            <el-col :span="24">
                                <!-- <div class="see-top">基础信息</div> -->
                                <el-row>
                                    <el-col :span="2">
                                        <span style="color:red;">* </span>
                                        <label style="color:#666;font-size:14px;line-height:30px;">企业名称</label>
                                    </el-col>
                                    <el-form-item prop="enterpriseName">
                                        <el-col :span="7">
                                            <el-input v-model="AddData.enterpriseName" placeholder="请输入企业名称" @input="changesize"></el-input>
                                            <div class="bormp" v-if="listblon">
                                                <div class="bpmslist" v-for="(v,k) in listcm" :key="k" @click="bpmclick(v.Name)">{{v.Name}}</div>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2.5" >
                                        <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">关联企业</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <div class="details">
                                                <div class="details-nav">
                                                    <div class="details-tag" >
                                                        <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length == 1">
                                                            <el-tag style="margin:4px;"
                                                            size="small"
                                                            type="info">
                                                                <el-tooltip effect="dark" :content="'关联企业：' + tag.enterpriseName" placement="top">
                                                                    <span> {{tag.enterpriseName}}</span>
                                                                </el-tooltip>
                                                            </el-tag>
                                                        </div>
                                                        <div class="tag-box" v-for="(tag,k) in newtagesAll" :key="k" v-if="newtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                            <el-tag size="small" type="info" style="margin:0;">
                                                                <span>{{newtagesAll[0].enterpriseName}}</span>
                                                            </el-tag>
                                                            <b> . . .</b>
                                                        </div>
                                                    </div>
                                                    <div class="details-conent">
                                                        共 <span style="color:#CF1724;"> {{newtagesAll.length}}</span> 家
                                                    </div> 
                                                </div>
                                                <div class="details-size" @click="dialogSelect">
                                                    <span style="color:#333;font-size:14px;">关 联</span>
                                                </div>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>

                                <el-row>
                                    <el-col :span="2.5" >
                                        <label style="color:#666;font-size:14px;line-height:30px;;padding-left:12px;">关联用户</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <div class="details">
                                                <div class="details-nav">
                                                    <div class="details-tag" >
                                                        <div class="tag-box" v-for="(urtag,k) in urnewtagesAll" :key="k" v-if="urnewtagesAll.length == 1">
                                                            <el-tag style="margin:4px;"
                                                            size="small"
                                                            type="info">
                                                                <el-tooltip effect="dark" :content="'关联用户：' + urtag.phone" placement="top">
                                                                    <span> {{urtag.phone}}</span>
                                                                </el-tooltip>
                                                            </el-tag>
                                                        </div>
                                                        <div class="tag-box" v-for="(urtag,k) in urnewtagesAll" :key="k" v-if="urnewtagesAll.length >= 2" style="padding:4px;box-sizing: border-box;">
                                                            <el-tag size="small" type="info" style="width:auto;margin:0;">
                                                                <span>{{urnewtagesAll[0].phone}}</span>
                                                            </el-tag>
                                                            <b> . . .</b>
                                                        </div>
                                                    </div>
                                                    <div class="details-conent">
                                                        共 <span style="color:#CF1724;"> {{urnewtagesAll.length}}</span> 个
                                                    </div> 
                                                </div>
                                                <div class="details-size" @click="urdialogSelect">
                                                    <span style="color:#333;font-size:14px;">关 联</span>
                                                </div>
                                            </div>
                                        </el-col>
                                    </el-form-item>
                                </el-row>

                                <el-row>
                                    <el-col :span="2" >
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">注册时间</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-date-picker
                                            disabled
                                            v-model="AddData.regTime"
                                            type="date"
                                            placeholder="请选择注册时间"
                                            value-format="yyyy-MM-dd"
                                            style="width:100%;">
                                            </el-date-picker>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">税号</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.rateNumber" placeholder="请输入税号"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">法人姓名</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.legalPerson" placeholder="请输入法人姓名"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">注册资本</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.regMoney" placeholder="请输入注册资本"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">实缴资本</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.recCap" placeholder="请输入实缴资本"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">企业类型</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.econKind" placeholder="请输入企业类型"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">更新日期</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.updatedDate" placeholder="请输入更新日期"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">注册地址</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled v-model="AddData.address" placeholder="请输入注册地址"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">经营范围</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input disabled type="textarea" v-model="AddData.scope" :rows="4"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">省</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-select @change="enterprises($event)" v-model="AddData.provinceCode" placeholder="请选择省" style="width:100%;">
                                                <el-option
                                                v-for="item in provinceOptions"
                                                :key="item.code"
                                                :label="item.name"
                                                :value="item.code">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">市</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7"> 
                                            <el-select v-model="AddData.cityCode" placeholder="请选择市" style="width:100%;" @change="citys($event)">
                                                <el-option
                                                v-for="item in cityOptions"
                                                :key="item.value"
                                                :label="item.name"
                                                :value="item.code">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">县 / 区</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7"> 
                                            <el-select v-model="AddData.areaCode" placeholder="请选择县 / 区" style="width:100%;" @change="countys($event)">
                                                <el-option
                                                v-for="item in countyOptions"
                                                :key="item.value"
                                                :label="item.name"
                                                :value="item.code">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">乡</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-select v-model="AddData.countyCode" placeholder="请选择乡" style="width:100%;" @change="villages($event)">
                                                <el-option
                                                v-for="item in villageOptions"
                                                :key="item.value"
                                                :label="item.name"
                                                :value="item.code">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2" >
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">法人身份证</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input v-model="AddData.legIdentity" placeholder="请输入法人身份证号码"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2" >
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">上户时间</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7" >
                                            <el-date-picker
                                            v-model="AddData.shTime"
                                            type="month"
                                            placeholder="请选择上户时间"
                                            value-format="yyyy-MM"
                                            style="width:100%;">
                                            </el-date-picker>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">联系电话</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input v-model="AddData.phone" placeholder="请输入联系电话"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">业务员</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input v-model="AddData.businesser" placeholder="请输入业务员姓名"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">企业类别</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-select v-model="AddData.type" placeholder="请选择企业类别" style="width:100%;">
                                                <el-option
                                                v-for="item in typeOptions"
                                                :key="item.id"
                                                :label="item.name"
                                                :value="item.id">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <!-- 上传企业logo -->
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">*</span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">上传企业logo</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-upload
                                            :limit="1"
                                            :on-exceed="handleExceed"
                                            class="avatar-uploader"
                                            :class="noneBtnDealImg ? 'disUoloadSty':''"
                                            :action="http + '/saas-common/upload/batch'"
                                            :file-list="fileList"
                                            :headers="header"
                                            list-type="picture-card"
                                            :on-success="handleSuccess"
                                            name="file"
                                            :on-remove="handleRemove">
                                                <i class="el-icon-plus" ></i>
                                            </el-upload>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <div style="height: 30px;">
                                            <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">官网图片</label>
                                        </div>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-upload
                                                :limit="1"
                                                :action="http + '/saas-common/upload/batch'"
                                                :class="gwBtnDealImg ? 'disyeloadSty':''"
                                                :headers="header"
                                                :file-list="lisgwearr"
                                                name='file'
                                                list-type="picture-card"
                                                :on-success="addgwlist"
                                                :on-preview="gwyeardPrediv"
                                                :on-remove="gwRemove">
                                                <i class="el-icon-plus"></i>
                                            </el-upload>
                                            <el-dialog :visible.sync="gwyeVisible" append-to-body>
                                                <img width="100%" :src="gwImgUrl" alt="">
                                            </el-dialog>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <div style="height: 30px;">
                                            <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">官网视频</label>
                                        </div>
                                    </el-col>
                                    <el-form-item>
                                        <el-upload :file-list="vdoarr" :on-remove="imgdel" class="satar-uploader" :http-request="uploadFile" :data="data" action="">
                                            <el-button icon="el-icon-upload2">官网视频上传</el-button>

                                        </el-upload>
                                        <div style="margin-top:5px;">
                                            <el-button size="small" icon="el-icon-upload2" v-show="showVideo" @click="downidevo">下载</el-button>
                                            <!-- <el-button size="small" icon="el-icon-upload2" v-show="showVideo" @click="imgdel">删除</el-button> -->
                                        </div>

                                    </el-form-item>
                                </el-row>
                            </el-col>
                        </el-collapse-item>
                        <el-collapse-item title="资质信息" name="2">
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">精煤仓数量（个）</label>
                                </el-col>
                                <el-form-item prop="jcoalSilo">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.jcoalSilo" placeholder="请输入精煤仓数量"></el-input>
                                    </el-col> 
                                </el-form-item>  
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">精煤仓储量（吨）</label>
                                </el-col>
                                <el-form-item prop="jcoalSiloCapacity">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.jcoalSiloCapacity" placeholder="请输入精煤仓储容量"></el-input>
                                    </el-col>
                                </el-form-item> 
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">原煤仓数量（个）</label>
                                </el-col>
                                <el-form-item prop="ycoalSilo"> 
                                    <el-col :span="7">
                                        <el-input v-model="AddData.ycoalSilo" placeholder="请输入原煤仓数量"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">原煤仓储量（吨）</label>
                                </el-col>
                                <el-form-item prop="ycoalSiloCapacity">
                                    <el-col :span="7">
                                        <el-input v-model="AddData.ycoalSiloCapacity" placeholder="请输入原煤仓储容量"></el-input>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="2">
                                    <!-- <span style="color:red;">* </span> -->
                                    <label style="color:#666;font-size:14px;line-height:30px;margin-left: 12px;">检测设备仪器</label>
                                </el-col>
                                <el-form-item>
                                    <el-col :span="7">
                                        <el-select v-model="AddData.testEquipmentType" size="small" style="width: 100%;">
                                            <el-option
                                            v-for="item in Extended"
                                            :key="item.value"
                                            :label="item.dkey"
                                            :value="item.dvalue">
                                            </el-option>
                                        </el-select>
                                    </el-col>
                                </el-form-item>
                            </el-row>
                            <el-row>
                                <el-col :span="4">
                                    <el-col :span="24">
                                        <div style="height: 30px;">
                                            <label style="color:#666;font-size:14px;line-height:30px;">营业执照</label>
                                        </div>
                                    </el-col>
                                    <el-form-item label-width="0px">
                                        <el-col :span="24">
                                            <el-upload
                                                :limit="1"
                                                :action="http + '/saas-common/upload/batch'"
                                                :class="yeBtnDealImg ? 'disyeloadSty':''"
                                                :headers="header"
                                                :file-list="listyearr"
                                                name='file'
                                                list-type="picture-card"
                                                :on-success="addyelist"
                                                :on-preview="jcyeardPrediv"
                                                :on-remove="yeRemove">
                                                <i class="el-icon-plus"></i>
                                            </el-upload>
                                            <el-dialog :visible.sync="imgyeVisible" append-to-body>
                                                <img width="100%" :src="yeImgUrl" alt="">
                                            </el-dialog>
                                        </el-col>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="4">
                                    <el-col :span="24">
                                        <label style="color:#666;font-size:14px;line-height:30px;">开户许可证</label>
                                    </el-col>
                                    <el-form-item label-width="0px">
                                        <el-col :span="24">
                                            <el-upload
                                                :limit="1"
                                                :action="http + '/saas-common/upload/batch'"
                                                :class="mpBtnDealImg ? 'dismploadSty':''"
                                                :headers="header"
                                                :file-list="listmparr"
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
                                        </el-col>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="4">
                                    <el-col :span="24">
                                        <label style="color:#666;font-size:14px;line-height:30px;">法人身份证（正面）</label>
                                    </el-col>
                                    <el-form-item label-width="0px">
                                        <el-col :span="24">
                                            <el-upload
                                                :limit="1"
                                                :action="http + '/saas-common/upload/batch'"
                                                :class="jcBtnDealImg ? 'disjcloadSty':''"
                                                :headers="header"
                                                :file-list="listjcarr"
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
                                        </el-col>
                                    </el-form-item>
                                </el-col>
                                <el-col :span="4">
                                    <el-col :span="24">
                                        <label style="color:#666;font-size:14px;line-height:30px;">法人身份证（反面）</label>
                                    </el-col>
                                    <el-form-item label-width="0px">
                                        <el-col :span="24">
                                            <el-upload
                                                :limit="1"
                                                :action="http + '/saas-common/upload/batch'"
                                                :headers="header"
                                                :file-list="listsfzarr"
                                                :class="sfzBtnDealImg ? 'dissfzloadSty':''"
                                                name='file'
                                                list-type="picture-card"
                                                :on-success="addsfzlist"
                                                :on-preview="sfzPictureCardPrediv"
                                                :on-remove="sfzRemove">
                                                <i class="el-icon-plus"></i>
                                            </el-upload>
                                            <el-dialog :visible.sync="imgsfzogVisible" append-to-body>
                                                <img width="100%" :src="sfzogImageUrl" alt="">
                                            </el-dialog>
                                        </el-col>
                                    </el-form-item>
                                </el-col>
                            </el-row>
                        </el-collapse-item>
                        <el-collapse-item name="3">
                            <template slot="title">
                                <span style="margin-right: 8px;">成交后评星</span>
                                <el-rate
                                    v-model="pxvals" :max="1" disabled show-score text-color="#E2231A" :colors="['#E2231A', '#E2231A', '#E2231A']" score-template="{value}">
                                </el-rate>
                            </template>
                            <el-col :span="24">
                                <el-table
                                :data="tableDatas"
                                style="width: 100%">
                                    <el-table-column label="云仓仓号" prop="warehouseNb">
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
                                    <el-table-column label="挥发分(Vdaf)%≤" :show-overflow-tooltip="true"  >
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
                                    <el-table-column align="left">
                                        <template slot="header" slot-scope="scope">
                                            <span style="margin-right: 2px;">合格兑现率</span>
                                            <el-tooltip class="item" effect="dark" content="合同数量（发运数量）承诺兑现情况 ；结算单据综合评定勾选" placement="top">
                                                <i style="color: #333;" class="el-icon-warning"></i>
                                            </el-tooltip>
                                        </template>
                                        <template slot-scope="scope">
                                            <el-checkbox v-model="scope.row.isContractQualified"></el-checkbox>
                                        </template>
                                    </el-table-column>
                                    <el-table-column align="left">
                                        <template slot="header" slot-scope="scope">
                                            <span style="margin-right: 2px;">质量合格率</span>
                                            <el-tooltip class="item" effect="dark" content="根据业务执行质量合格率执行扣罚情况综合评定勾选" placement="top">
                                                <i style="color: #333;" class="el-icon-warning"></i>
                                            </el-tooltip>
                                        </template>
                                        <template slot-scope="scope">
                                            <el-checkbox v-model="scope.row.isQualityQualified"></el-checkbox>
                                        </template>
                                    </el-table-column>
                                </el-table>
                            </el-col>
                        </el-collapse-item>
                        <el-collapse-item title="配置信息" name="4">
                            <el-col :span="24">
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">网络货运</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-radio v-model="AddData.isNetwork" :label="true">是</el-radio>
                                            <el-radio v-model="AddData.isNetwork" :label="false">否</el-radio>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">小程序轮播标志</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-radio v-model="AddData.rotationFlag" :label="true">是</el-radio>
                                            <el-radio v-model="AddData.rotationFlag" :label="false">否</el-radio>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2.5" >
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">对公流程配置</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-select v-model="AddData.pubFlowId" style="width:100%;">
                                                <el-option
                                                v-for="item in AddData.pubFlows"
                                                :key="item.id"
                                                :label="item.planName"
                                                :value="item.id">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2.5" >
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">对私流程配置</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-select v-model="AddData.perFlowId" style="width:100%;">
                                                <el-option
                                                v-for="item in AddData.perFlows"
                                                :key="item.id"
                                                :label="item.planName"
                                                :value="item.id">
                                                </el-option>
                                            </el-select>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">签章免费份数</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input v-model="AddData.signFreeCount" placeholder="请输入签章免费份数"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                                <el-row>
                                    <el-col :span="2">
                                        <!-- <span style="color:red;">* </span> -->
                                        <label style="color:#666;font-size:14px;line-height:30px;margin-left:12px;">签章收费份数</label>
                                    </el-col>
                                    <el-form-item>
                                        <el-col :span="7">
                                            <el-input v-model="AddData.signChargeCount" placeholder="请输入签章收费份数"></el-input>
                                        </el-col>
                                    </el-form-item>
                                </el-row>
                            </el-col>
                        </el-collapse-item>
                    </el-collapse>
                </el-row>
            </div>
        </div>
        <!-- 内容结束 -->
        <!-- 供应商单位弹窗开始 -->
        <el-dialog
            :visible.sync="dialogSelectUnit"
            width="60%"
            :before-close="BuyerhandleClose">
            <template slot="title">
                <span>选择供应商</span>
            </template>
            <!-- 显示合同数据选择内容框 -->
            <div class="smbox">
                <el-tag
                size="small"
                type="info"
                v-for="(v,k) in commontags" :key="k"
                style="margin:5px 5px;color:#333;">
                <span v-if="v.enterpriseName">{{v.enterpriseName}}</span>
                </el-tag>
            </div>
            <!-- 供应商单位内容框 -->
            <!-- <el-input type="textarea" v-model="dialogtextarea"></el-input> -->
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
                            </div>
                            <el-table border ref="multipleTable" :data="CompanytableData" :header-cell-style="{ background: '#F6F6F6', color: '#333' }" style="width: 100%;" @selection-change="contractgetCurrentRow" :row-key="keyId">
                                <el-table-column type="selection" :reserve-selection="true"></el-table-column>
                                <el-table-column prop="id" v-if="show">
                                </el-table-column>
                                <el-table-column prop="enterpriseName" label="单位名称">
                                </el-table-column>
                                <el-table-column prop="createTime" label="创建时间">
                                </el-table-column>
                            </el-table>
                            <!-- 分页 -->
                            <div class="pur-pages">
                                <el-pagination
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
                </div>
            </div>
        </el-dialog>
        <!--供应商单位弹窗结束 -->

        <!-- 用户关联弹窗开始 -->
        <el-dialog
            :visible.sync="urdialogSelectUnit"
            width="60%"
            :before-close="urBuyerhandleClose">
            <template slot="title">
                <span>选择用户</span>
            </template>
            <!-- 显示用户关联数据内容框 -->
            <div class="smbox">
                <el-tag
                    size="small"
                    closable
                    @close="handleClose(k)"
                    type="info"
                    v-for="(v,k) in urcommontags" :key="k"
                    style="margin:5px 5px;color:#333;">
                    <span v-if="v.phone">{{v.phone}}</span>
                </el-tag>
                <!-- <el-input
                    class="input-new-tag"
                    v-if="inputVisible"
                    v-model="inputValue"
                    ref="saveTagInput"
                    size="small"
                    @keyup.enter.native="handleInputConfirm"
                    @blur="handleInputConfirm"
                    >
                </el-input> -->
                <!-- <el-button v-else class="button-new-tag" size="small" @click="showInput">+ New Tag</el-button> -->
            </div>
            <!-- 用户关联内容框 -->
            <div class="Wdialog">
                <div class="concentBox">
                    <div class="textbutton">
                        <!-- 用户关联应用按钮 -->
                        <el-button size="small" @click="ursalesUnitApplication">应用</el-button>
                    </div>
                    <div class="text">
                        <!-- 用户关联信息 -->
                        <div class="textAll">
                            <div class="textlast">
                                <el-input
                                    @keyup.enter.native="ursearchtwo"
                                    size="small"
                                    placeholder="请输入内容"
                                    prefix-icon="el-icon-search"
                                    v-model="urpurSearchtwo"
                                    style="width:200px;float:right;"
                                    >
                                </el-input>
                            </div>
                                <!-- @selection-change="urcontractgetCurrentRow"  -->

                            <el-table 
                                border 
                                ref="urmultipleTable" 
                                :data="urCompanytableData" 
                                :header-cell-style="{ background: '#F6F6F6', color: '#333' }" 
                                @select="selTag"
                                style="width: 100%;" 
                                :row-key="urkeyId">
                                <el-table-column type="selection" :reserve-selection="true"></el-table-column>
                                <el-table-column prop="id" v-if="urshow">
                                </el-table-column>
                                <el-table-column prop="phone" label="用户手机号">
                                </el-table-column>
                                <el-table-column
                                    prop="username"
                                    label="姓名"
                                    :show-overflow-tooltip="true">
                                    <template slot-scope="scope">
                                        <span v-if="scope.row.username">{{ scope.row.username }}</span>
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
                            </el-table>
                            <!-- 分页 -->
                            <div class="pur-pages">
                                <el-pagination
                                background
                                :page-size.sync="ursizetwo"
                                :current-page.sync="urcurrenttwo"
                                layout="prev, pager, next"
                                :total="urtotaltwo"
                                @current-change='changeCurrent'>
                                </el-pagination>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </el-dialog>
        <!--用户关联弹窗结束 -->
    </el-form>

</template>
<script>
// 引入上传文件方法
import { upload, uploadByPieces } from "../../upload.js";
export default {
    data() {
        return {

            // 视频上传
            data:{},
            src: '',
            showVideo: false,
            vdoarr:[],
            /*
            资质信息
            */
            uploading:false,   // 上传遮罩层
            // 官网图片
            gwImgUrl: '',
            gwyeVisible: false,
            lisgwearr:[],
            gwBtnDealImg:false,
            // 营业执照
            yeImgUrl: '',
            imgyeVisible: false,
            listyearr:[],
            yeBtnDealImg:false,
            // 开户许可证
            mpogImageUrl: '',
            imgmpogVisible: false,
            listmparr:[],
            mpBtnDealImg:false,
            // 身份证正面
            jcogImageUrl: '',
            imgjcogVisible: false,
            listjcarr:[],
            jcBtnDealImg:false,
            // 身份证反面
            sfzogImageUrl: '',
            imgsfzogVisible: false,
            listsfzarr:[],
            sfzBtnDealImg:false,
            // 手风琴
            activeNames:['2'],
            /*
            成交后评星
            */
            pxvals:null,
            tableDatas: [],
            /* 
            供应商单位数据列表
            */
            currenttwo:1,                   
            sizetwo:10,                      
            totaltwo:0,                     
            purSearchtwo:'',                               
            dialogSelectUnit:false,                 // 选择单位弹窗
            CompanytableData:[],                    // 单位表格数据
            show:false,                             // 选择单位表格隐藏 id
            commontags: [],                         // 合同集合回显数据
            newtagesAll:[],
            /* 
            用户列表数据列表
            */
            urcurrenttwo:1,                   
            ursizetwo:10,                      
            urtotaltwo:0,                     
            urpurSearchtwo:'',                               
            urdialogSelectUnit:false,                 // 选择用户弹窗
            urCompanytableData:[],                    // 用户表格数据
            urshow:false,                             // 选择用户表格隐藏 id
            urcommontags: [],                         // 用户集合回显数据
            urnewtagesAll:[],

            http:'',
            header:{},                              // 上传token
            imageUrl:'',
            dialogImageUrl: '',
            noneBtnDealImg:false, 
            // dialogVisible: false,
            /* 
            编辑销售收款单数据
            */
            typeOptions:[
                {
                    id:"1",
                    name:'一般纳税人'
                },
                {
                    id:"2",
                    name:'小规模纳税人'
                },
            ],
            purTitle: "",                           // 标题
            // 检测设备
            Extended:[
              {
                  dkey:'岩相监测',
                  dvalue:1
              },
              {
                  dkey:'小焦炉设备',
                  dvalue:2
              }
            ], 
            // 表单验证
            rules: {
                jcoalSilo: [
                    // { required: true, message: '请输入精煤仓数量', trigger: 'blur' },
                    { pattern: /^\d+$/, message: '数量格式不正确', trigger: 'blur'}
                ],  
                jcoalSiloCapacity: [
                    // { required: true, message: '请输入精煤仓储容量', trigger: 'blur' }
                    { pattern: /^\d+$/, message: '容量格式不正确', trigger: 'blur'}
                ],    
                ycoalSilo: [
                    // { required: true, message: '请输入原煤仓数量', trigger: 'blur' }
                    { pattern: /^\d+$/, message: '数量格式不正确', trigger: 'blur'}
                ],     
                ycoalSiloCapacity: [
                    // { required: true, message: '请输入原煤仓储容量', trigger: 'blur' }
                    { pattern: /^\d+$/, message: '容量格式不正确', trigger: 'blur'}
                ], 
                signFreeCount: [
                    { required: true, message: '请输入签章免费份数', trigger: 'blur' }
                ],
                signChargeCount: [
                    { required: true, message: '请输入签章收费份数', trigger: 'blur' }
                ],
                regTime: [
                    { required: true, message: '请输入注册时间', trigger: 'change' }
                ],
                legIdentity: [
                    { required: true, message: '请输入法人身份证号', trigger: 'blur' },
                    { pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/, message: '请输入正确的身份证号', trigger: 'blur'},
                ],
                enterpriseName: [
                    { required: true, message: '请输入企业名称', trigger: 'blur' }
                ],
                shTime: [
                    { required: true, message: '请输入上户时间', trigger: 'blur' }
                ],
                phone: [
                    { required: true, message: '请输入联系电话', trigger: 'blur' },
                    { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'blur'}
                ],
                regMoney: [
                    { required: true, message: '请输入注册资本', trigger: 'blur' }
                ],
                businesser: [
                    { required: true, message: '请输入业务员姓名', trigger: 'blur' },
                    { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
                ],
                legalPerson: [
                    { required: true, message: '请输入法人姓名', trigger: 'blur' },
                    { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
                ],
                rateNumber:[
                    { required: true, message: '请输入税号', trigger: 'blur' },
                    { pattern: /^[A-Z0-9]{15}$|^[A-Z0-9]{17}$|^[A-Z0-9]{18}$|^[A-Z0-9]{20}$/, message: '请输入正确的企业税号', trigger: 'blur'}
                ],
                isNetwork: [
                    { required: true, message: '请选择是否为网络货运', trigger: 'blur' }
                ],
                rotationFlag: [
                    { required: true, message: '请选择是否为小程序轮播标志', trigger: 'blur' }
                ],
            },
            // 行政区划
            provinceOptions:[],             // 省
            cityOptions: [],                // 市
            countyOptions: [],              // 县 区
            villageOptions: [],             // 乡
            // 编辑表单数据
            AddData: {
                carriageNameone:[],
                carriageNametwo:[],
                etcs:[],
                logo:null,
                carriageName:1,
            },
            // 企业认证
            listcm:[],   // 企业列表
            listblon:false,
            enterprise:{},  // 企业信息
            fileList:[],  // 企业logo
            inputVisible: false,
            inputValue: ''
        };
    },
    created() {
        this.header.Authorization = sessionStorage.getItem('Authorization');
        this.http  = this.$store.state.http;                  // IP地址
    },
    mounted () {
        // this.getVideo();
        // videojs("my-video",()=>{this.play();});
        // 接收主页单条数据
        this.AddData = JSON.parse(sessionStorage.getItem('queryAll'));
        if(!this.AddData){
            // 直接输入网址 / 后退进入时没有待编辑的记录:回列表页,否则下面取字段就抛错
            return this.$router.replace({name:'Qycensus'});
        }
        // 模块标题
        this.purTitle = this.AddData.enterpriseName;
        this.getpurcon(this.AddData.id);  // 查询评星数据
        this.pxvals = this.AddData.qualifiedScore;
        // 企业logo
        if(this.AddData.logo){
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.logo +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.fileList.push({url: url})
            this.noneBtnDealImg = this.fileList.length >= 1; // 上传完成后隐藏上传按钮
        }
        this.newtagesAll = this.AddData.etcs;  // 累计关联企业数
        this.urnewtagesAll = this.AddData.urcs;  // 累计关联用户数
        this.api.areaInfo.all({parentCode:0})
        .then(res=>{
            if(res.data.code == 200){
                this.provinceOptions = res.data.data;
            }
        })
        if(this.AddData.cityCode != null){
            // 市
            this.api.areaInfo.all({parentCode:this.AddData.provinceCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.cityOptions = res.data.data;
                }
            })
            // 县 区
            this.api.areaInfo.all({parentCode:this.AddData.cityCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.countyOptions = res.data.data;
                }
            })
            // 乡
            this.api.areaInfo.all({parentCode:this.AddData.areaCode})
            .then(res=>{
                if(res.data.code == 200){
                    this.villageOptions = res.data.data;
                }
            })
        }
        /*
        资质信息
        */
        if (this.AddData.advertisingVideoId) {
            let newobj = {
                name:'官网视频',
                response:{
                    code:200,
                    data:this.AddData.advertisingVideoId
                }
            }
            this.vdoarr.push(newobj);
            // let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.advertisingVideoId +
            // '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            // this.vdoarr.push({url: url,fileName:'官网视频'});
        };
        
        // 官网照片
        if (this.AddData.websiteImgId) {
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.websiteImgId +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.lisgwearr.push({url: url});
            this.gwBtnDealImg = this.lisgwearr.length >= 1;
        };
        // 营业执照
        if (this.AddData.busLicenseId) {
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.busLicenseId +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.listyearr.push({url: url});
            this.yeBtnDealImg = this.listyearr.length >= 1;
        };
        // 开户许可照
        if (this.AddData.accOpeningPermitId) {
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.accOpeningPermitId +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.listmparr.push({url: url});
            this.mpBtnDealImg = this.listmparr.length >= 1;
        };

        // 身份证正面
        if (this.AddData.legalFtCardId) {
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.legalFtCardId +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.listjcarr.push({url: url});
            this.jcBtnDealImg = this.listjcarr.length >= 1;
        };
        // 身份证反面
        if (this.AddData.legalBkCardId) {
            let url =  this.http + '/saas-common/upload/imageShow?fileName=' + this.AddData.legalBkCardId +
            '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99'
            this.listsfzarr.push({url: url});
            this.sfzBtnDealImg = this.listsfzarr.length >= 1;
        };
    },
    // props: ["data"],
    methods: {
      
        /*
        资质信息
        */

        async uploadFile({ data, file }) {
          // data是上传时附带的额外参数，file是文件
          let url = this.http + '/saas-common/upload/bigUpload'; //上传文件接口
          this.uploading = true;
          try {
            // 如果文件小于5MB，直接上传
            if (file.size < 5 * 1024 * 1024) {
                console.log('直接上传')
                let formData = new FormData();
                for (let key in data) {
                    formData.append(key, data[key]);
                }
                formData.append("file", file);
    
                const res = await upload(url,formData);
                this.uploading = false;
                return res;
            } else {
              // 如果文件大于等于5MB，分片上传
            //   console.log('分片上传')
              data.file = file;
              const res = await uploadByPieces(url,data);
              for(let i=0;i<res.length;i++){
                if(res[i].data.id){
                    this.showVideo = true;
                    this.AddData.advertisingVideoId = res[i].data.id;   // 官网视频id
                    this.src = this.http + '/saas-common/upload/getVideo?fileName=' + res[i].data.id + '&AuthorizationId=wxda0b1c468b83bc99&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f'
                }
              }
              if(res)
              this.uploading = false;
              return res;
            }
          } catch (e) {
            this.uploading = false;
            return e;
          }
        },
        // 删除合同文件
        imgdel(file, fileList){
            console.log(file, fileList)
            this.api.addfile.del({id:this.AddData.advertisingVideoId})
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: 'warning',
                        message: '删除文件成功!'
                    });
                    this.showVideo = false;
                    this.AddData.advertisingVideoId = '';
                }
            })
        },
        // 下载
        downidevo(){
            window.open(this.src)
        },
        // 上传官网照片
        addgwlist(response,file,fileList){
            if (response.code == 200) {
                this.AddData.websiteImgId = response.data[0].id;   // 官网图id
                this.AddData.websiteCedImgId = response.data[0].id;   // 官网图压缩id
                this.gwBtnDealImg = fileList.length >= 1;
                this.$message({
                    message: "上传官网照片成功！",
                    type: "success",
                });
            }
        },
        // 预览官网照片
        gwyeardPrediv(file) {
          this.gwImgUrl = file.url;
          this.gwyeVisible = true;
        },
        // 删除官网照片
        gwRemove(file, fileList) {
            this.api.addfile.del({id:this.AddData.websiteImgId})
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: 'warning',
                        message: '删除文件成功!'
                    });
                    this.gwBtnDealImg = fileList.length >= 1;
                    this.AddData.websiteImgId = '';   // 官网图id
                    this.lisgwearr = [];
                }
            })
            this.api.addfile.del({id:this.AddData.websiteCedImgId})
            .then(res=>{
                if(res.data.code == 200){
                    this.$message({
                        type: 'warning',
                        message: '删除文件成功!'
                    });
                    this.gwBtnDealImg = fileList.length >= 1;
                    this.AddData.websiteCedImgId = '';   // 官网图压缩id
                    this.lisgwearr = [];
                }
            })

        },
        // 上传营业执照
        addyelist(response,file,fileList){
            this.uploading = true;
            this.api.ller.getOcrResult({
                type:2,  // 1身份证 2营业执照 3许可证
                id:response.data[0].compressedId
            }).then(res=>{
                this.uploading = false;
                if(res.data.code == 200){
                    this.AddData.busLicenseId = response.data[0].id;
                    this.yeBtnDealImg = fileList.length >= 1;
                    this.$message({
                        message: '营业执照校验成功！',
                        type: 'success'
                    });
                }
                if(res.data.code == 701){
                    this.listyearr = [];
                    this.AddData.busLicenseId = '';
                    this.yeBtnDealImg = false;
                    if(res.data.data == 'LimitExceeded.TooLargeFileError'){
                        this.$message({
                            message: '上传的营业执照图片过大,请重新上传！',
                            type: 'warning'
                        });
                    } 
                    if(res.data.data == 'FailedOperation.NoBizLicense'){
                        this.$message({
                            message: '非营业执照,请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.OcrFailed'){
                        this.$message({
                            message: '营业执照识别失败，请上传清晰(完整)图像,请重新上传！',
                            type: 'warning'
                        });
                    }
                }
            })

        },
        // 预览营业执照
        jcyeardPrediv(file) {
          this.yeImgUrl = file.url;
          this.imgyeVisible = true;
        },
        // 删除营业执照
        yeRemove(file, fileList) {
            this.yeBtnDealImg = fileList.length >= 1;
            this.AddData.busLicenseId = '';
            this.listyearr = [];
        },

        // 上传开户许可证
        addmplist(response,file,fileList){
            this.uploading = true;
            this.api.ller.getOcrResult({
                type:3,  // 1身份证 2营业执照 3许可证
                id:response.data[0].compressedId
            }).then(res=>{
                this.uploading = false;
                if(res.data.code == 200) {
                    let permit = JSON.parse(res.data.data);
                    let permitTypeOne = permit.findIndex(el => el.detectedText == '开户许可证');
                    let accInfo = permit.findIndex(el => el.detectedText == '基本存款账户信息');
                    let legal = permit.findIndex(el => el.detectedText == '法定代表人:');
                    // if(permitTypeOne != -1 || (accInfo != -1&&legal != -1)){
                        this.AddData.accOpeningPermitId = response.data[0].id;
                        this.mpBtnDealImg = fileList.length >= 1;
                        this.$message({
                            message: '许可证校验成功！',
                            type: 'success'
                        });
                    // } else {
                    //     this.listmparr = [];
                    //     this.AddData.accOpeningPermitId = '';
                    //     this.mpBtnDealImg = false;
                    //     this.$message({
                    //         message: '请上传正确的许可证图片,请重新上传！',
                    //         type: 'warning'
                    //     });
                    // }
                }
                if(res.data.code == 701){
                    this.listmparr = [];
                    this.AddData.accOpeningPermitId = '';
                    this.mpBtnDealImg = false;
                    if(res.data.data == 'LimitExceeded.TooLargeFileError'){
                        this.$message({
                            message: '上传的许可证图片过大,请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.ImageNoText'){
                        this.$message({
                            message: '照片中未检测到文本,请重新上传！',
                            type: 'warning'
                        });
                    }
                }
            })

        },
        // 预览开户许可证
        mpPictureCardPrediv(file) {
          this.mpogImageUrl = file.url;
          this.imgmpogVisible = true;
        },
        // 删除开户许可证
        mpRemove(file, fileList) {
            this.mpBtnDealImg = fileList.length >= 1;
            this.AddData.accOpeningPermitId = '';
            this.listmparr = [];
        },

        // 上传身份证正面
        addjclist(response,file,fileList){
            this.uploading = true;
            // type 1、身份证  2、营业执照 3、开户许可证
            this.api.ller.getOcrResult({
                id:response.data[0].compressedId,
                type:1      
            }).then(res => {
                this.uploading = false;
                if(res.data.code == 200){
                    let IDcardFront = JSON.parse(res.data.data);
                    if(IDcardFront.Name != ''){
                        this.AddData.legalFtCardId = response.data[0].id;
                        this.jcBtnDealImg = fileList.length >= 1;
                        this.$message({
                            message: '身份证(正面)校验成功！',
                            type: 'success'
                        });
                    } else {
                        this.listjcarr = [];
                        this.AddData.legalFtCardId = '';
                        this.jcBtnDealImg = false;
                        this.$message({
                            message: '身份证(正面)校验失败，请重新上传！',
                            type: 'warning'
                        });
                    }
                }
                if (res.data.code == 701) {
                    this.listjcarr = [];
                    this.AddData.legalFtCardId = '';
                    this.jcBtnDealImg = false;

                    if(res.data.data == 'LimitExceeded.TooLargeFileError'){
                        this.$message({
                            message: '身份证(正面)上传的身份证图片过大，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.ImageNoIdCard'){
                        this.$message({
                            message: '图片未检测到身份证，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.IdCardInfoIllegal'){
                        this.$message({
                            message: '身份证信息不合法或图片模糊，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.ImageBlur'){
                        this.$message({
                            message: '身份证图片模糊，请重新上传清晰图像！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.IdCardTooSmall'){
                        this.$message({
                            message: '图片分辨率过小或身份证在原图中的占比过小，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.OcrFailed'){
                        this.$message({
                            message: '身份证识别失败，请重新上传清晰(完整)图像！',
                            type: 'warning'
                        });
                    }
                }
            })
        },
        // 预览身份证正面
        jcPictureCardPrediv(file) {
          this.jcogImageUrl = file.url;
          this.imgjcogVisible = true;
        },
        // 删除身份证正面
        jcRemove(file, fileList) {
            this.jcBtnDealImg = fileList.length >= 1;
            this.AddData.legalFtCardId = '';
            this.listjcarr = [];
        },

        // 上传身份证反面
        addsfzlist(response,file,fileList){
            this.uploading = true;
            // type 1、身份证  2、营业执照 3、开户许可证
            this.api.ller.getOcrResult({
                id:response.data[0].compressedId,
                type:1      
            }).then(res => {
                this.uploading = false;
                if(res.data.code == 200){
                    let IDcardBack = JSON.parse(res.data.data);
                    if(IDcardBack.Name == '' || IDcardBack.Authority != ''){
                        this.AddData.legalBkCardId = response.data[0].id;
                        this.sfzBtnDealImg = fileList.length >= 1;
                        this.$message({
                            message: '身份证(反面)校验成功！',
                            type: 'success'
                        });
                    } else {
                        this.listsfzarr = [];
                        this.AddData.legalBkCardId = '';
                        this.sfzBtnDealImg = false;
                        this.$message({
                            message: '身份证(反面)校验失败，请重新上传！',
                            type: 'warning'
                        });
                    }
                }
                if (res.data.code == 701) {
                    this.listsfzarr = [];
                    this.AddData.legalBkCardId = '';
                    this.sfzBtnDealImg = false;
                    if(res.data.data == 'LimitExceeded.TooLargeFileError'){
                        this.$message({
                            message: '上传的身份证图片过大，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.ImageNoIdCard'){
                        this.$message({
                            message: '照片未检测到身份证，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.IdCardInfoIllegal'){
                        this.$message({
                            message: '身份证信息不合法或图片模糊，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.ImageBlur'){
                        this.$message({
                            message: '身份证图片模糊，请重新上传清晰图像！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.IdCardTooSmall'){
                        this.$message({
                            message: '图片分辨率过小或身份证在原图中的占比过小，请重新上传！',
                            type: 'warning'
                        });
                    }
                    if(res.data.data == 'FailedOperation.OcrFailed'){
                        this.$message({
                            message: '身份证识别失败，请重新上传清晰(完整)图像！',
                            type: 'warning'
                        });
                    }
                }
            })
        },
        // 预览身份证反面
        sfzPictureCardPrediv(file) {
          this.sfzogImageUrl = file.url;
          this.imgsfzogVisible = true;
        },
        // 删除身份证反面
        sfzRemove(file, fileList) {
            this.sfzBtnDealImg = fileList.length >= 1;
            this.AddData.legalBkCardId = '';
            this.listsfzarr = [];
        },
        /*
        评星功能
        */
        getpurcon(id){
            // 渲染表格
            this.api.busrodect.all({
                effective:1,
                status:1,
                companyId:id,
                searchStr:''
            })
            .then(res=>{
              if(res.data.code == 200){
                    this.tableDatas = res.data.data;
              }
            })
        },
        /* 
        供应商单位的选择功能开始
        */
        // 合同编号选择表格
        contractgetCurrentRow(row){
            this.commontags = row;                         // 获取选中单位数据集合
        },
        keyId(row){
            return row.id;
        },
        // 供应商单位表格数据
        allsuppliertwo(){
            // 传参
            var listByPageDatatwo = {
                // partnerType:1,
                isExclude:true,
                searchStr:this.purSearchtwo,
                page:{
                    current:this.currenttwo,
                    size:this.sizetwo
                }
            }
            // 渲染表格
            this.api.enterData.all(listByPageDatatwo)
            .then(res=>{
                if(res.data.code == 200){
                    this.CompanytableData = res.data.data.records
                    this.totaltwo = res.data.data.total;
                    this.CompanytableData = res.data.data.records;
                    this.totaltwo = res.data.data.total;
                    this.$refs.multipleTable.clearSelection();
                    if (this.newtagesAll.length > 0) {
                        this.commontags = this.newtagesAll;
                        for (var i = 0; i < this.newtagesAll.length; i++) {
                            for (var j = 0; j < this.CompanytableData.length; j++) {
                                if (this.newtagesAll[i].id == this.CompanytableData[j].id) {//items里面包含有的都回显打勾                 
                                    this.$refs.multipleTable.toggleRowSelection(this.CompanytableData[j],true);
                                }
                            }
                        }
                    }

                }
            })
        },
        // 关联企业搜索功能
        searchtwo() {
            // this.CompanytableData = [];
            this.allsuppliertwo();
            this.currenttwo = 1;
        },
        // 关联企业弹窗
        dialogSelect(){
            this.allsuppliertwo();                              // 渲染表格
            this.dialogSelectUnit=true;                         // 开启弹窗
        },
        // 关联企业应用按钮事件
        salesUnitApplication(){
            
            this.newtagesAll = this.commontags;                 // 获取单位名称，回显到表单
            this.commontags = [];                               // 清空选中框回显内容
            this.$refs.multipleTable.clearSelection();          // 清空列表复选框
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
        },
        // 关联企业 关闭弹窗
        BuyerhandleClose() {
            this.commontags = [];                               // 清空选中框回显内容
            this.$refs.multipleTable.clearSelection();          // 清空列表复选框
            this.dialogSelectUnit = false;                      // 关闭整个弹窗
        },
        /* 
        供应商单位的选择功能结束
        */

        /* 
        用户关联的选择功能开始
        */
        // 用户关联选择表格
        urcontractgetCurrentRow(row){
        //    this.urcommontags = row;
        },
        selTag(selection,row){
            let arr = []
            let phone = row.phone
            this.urcommontags.forEach(item => {
                arr.push(item.phone)
            }) 
            if(!arr.includes(phone)){
                this.urcommontags.push(row)
            }else if(!selection.includes(row)){
                // @select 在取消勾选时同样触发:原先只增不减,取消的用户点「应用」后仍被关联
                this.urcommontags = this.urcommontags.filter(item => item.phone != phone)
            }
        },
        urkeyId(row){
            return row.id;
        },
        // 用户关联表格数据
        urallsuppliertwo(){
            this.urcurrenttwo = 1
            // 传参
            var listByPageDatatwo = {
                hasPhone:true,
                searchStr:this.urpurSearchtwo,
                page:{
                    current:this.urcurrenttwo,
                    size:this.ursizetwo
                },
                channel:2
            }
            // 渲染表格
            this.api.authAll.listByPage(listByPageDatatwo)
            .then(res=>{
                if(res.data.code == 200){
                    this.urCompanytableData = res.data.data.records;
                    this.urtotaltwo = res.data.data.total;
                    this.$refs.urmultipleTable.clearSelection();
                    if (this.urnewtagesAll.length > 0) {
                        this.urcommontags = this.urnewtagesAll;
                        // for (var i = 0; i < this.urnewtagesAll.length; i++) {
                        //     for (var j = 0; j < this.urCompanytableData.length; j++) {
                        //         if (this.urnewtagesAll[i].id == this.urCompanytableData[j].id) {//items里面包含有的都回显打勾                 
                        //             this.$refs.urmultipleTable.toggleRowSelection(this.urCompanytableData[j],true);
                        //         }
                        //     }
                        // }
                    }
                }
            })
        },
        changeCurrent(){
            // 传参
            var listByPageDatatwo = {
                hasPhone:true,
                searchStr:this.urpurSearchtwo,
                page:{
                    current:this.urcurrenttwo,
                    size:this.ursizetwo
                },
                channel:2
            }
            // 渲染表格
            this.api.authAll.listByPage(listByPageDatatwo)
            .then(res=>{
                if(res.data.code == 200){
                    this.urCompanytableData = res.data.data.records;
                    this.urtotaltwo = res.data.data.total;
                    if (this.urnewtagesAll.length > 0) {
                        this.urcommontags = this.urnewtagesAll;
                        // for (var i = 0; i < this.urnewtagesAll.length; i++) {
                        //     for (var j = 0; j < this.urCompanytableData.length; j++) {
                        //         if (this.urnewtagesAll[i].id == this.urCompanytableData[j].id) {//items里面包含有的都回显打勾                 
                        //             this.$refs.urmultipleTable.toggleRowSelection(this.urCompanytableData[j],true);
                        //         }
                        //     }
                        // }
                    }
                }
            })
        },
        // 用户关联搜索功能
        ursearchtwo() {
            this.urcurrenttwo = 1;
            this.changeCurrent();
        },
        // 关联用户弹窗
        urdialogSelect(){
            this.urallsuppliertwo();                              // 渲染表格
            this.urdialogSelectUnit=true;                         // 开启弹窗
        },
        // 关联用户应用按钮事件
        ursalesUnitApplication(){
            this.urnewtagesAll = this.urcommontags;               // 获取用户名称，回显到表单
            this.urcommontags = [];                               // 清空选中框回显内容
            this.urpurSearchtwo = ''
            this.$refs.urmultipleTable.clearSelection();          // 清空列表复选框
            this.urdialogSelectUnit = false;                      // 关闭整个弹窗
        },
        // 关联用户 关闭弹窗
        urBuyerhandleClose() {
            this.urcommontags = [];                               // 清空选中框回显内容
            this.$refs.urmultipleTable.clearSelection();          // 清空列表复选框
            this.urdialogSelectUnit = false;                      // 关闭整个弹窗
        },
        /* 
        用户的选择功能结束
        */


        // 查询企业
        changesize(){
            if(this.AddData.enterpriseName == ""){
                this.enterprise == {};
                this.listblon = false;
            }
            this.api.enterData.seletcname({etcName:this.AddData.enterpriseName})
            .then(res=>{
                if(res.data.code == 200 && res.data.data.length>0){
                    this.listcm = res.data.data;
                    this.listblon = true;
                }
            })
        },
        // 选中企业
        bpmclick(name){
            this.AddData.enterpriseName = name;
            this.listblon = false;
            this.api.enterData.seletcinfo({enterpriseName:name,isAutoAdd:false})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.regTime = res.data.data.regTime;
                    this.AddData.rateNumber = res.data.data.creditCode;
                    this.AddData.legalPerson = res.data.data.legalPerson;
                    this.AddData.regMoney = res.data.data.regMoney;
                    this.AddData.recCap = res.data.data.recCap;
                    this.AddData.econKind = res.data.data.econKind;
                    this.AddData.updatedDate = res.data.data.updatedDate;
                    this.AddData.address = res.data.data.address;
                    this.AddData.scope = res.data.data.scope;
                    this.AddData.imageUrl = res.data.data.imageUrl;
                    this.AddData.companyId = res.data.data.id;
                    if(res.data.data.provinceCode){
                        this.AddData.provinceCode = res.data.data.provinceCode;  // 省
                        this.enterprises(this.AddData.provinceCode); 

                        // if(res.data.data.cityCode){
                        //     this.AddData.cityCode = res.data.data.cityCode;          // 市
                        //     this.citys(this.AddData.cityCode);
                        //     if(res.data.data.areaCode){
                        //         this.AddData.areaCode = res.data.data.areaCode;          // 县/区
                        //         this.countys(this.AddData.areaCode);
                        //         if(res.data.data.countyCode){
                        //             this.AddData.countyCode = res.data.data.countyCode;      // 乡 
                        //             this.villages(this.AddData.countyCode);
                        //         }
                        //     }
                        // }
                    } 
                    
                }
            })
        },
        // 编辑企业信息表单提交
        AddDatasubmitForm(formName) {
            // if(!this.AddData.logo){
            //     return this.$message({
            //         type: "error",
            //         message: "请上传企业logo!",
            //         duration:800,  
            //     });
            // };
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.api.busrodect.updateBatch(this.tableDatas).then(res=>{});
                    let pxnum = 0;
                    this.tableDatas.map(item=>{
                        if(item.isContractQualified){
                            pxnum = pxnum + 0.5;
                        }
                        if(item.isQualityQualified){
                            pxnum = pxnum + 0.5;
                        }
                    })
                    if(pxnum >=1){
                        this.AddData.qualifiedScore = 1;
                    }else{
                        this.AddData.qualifiedScore = pxnum;
                    }
                    // 关联企业
                    this.AddData.etcs = [];
                    this.newtagesAll.map(item=>{
                        this.AddData.etcs.push({
                            groupId:this.AddData.id,
                            enterpriseName:item.enterpriseName,
                            id:item.id
                        })
                    })
                    // 关联用户
                    this.AddData.urcs = [];
                    this.urnewtagesAll.map(item=>{
                        this.AddData.urcs.push({
                            companyId:this.AddData.id,
                            phone:item.phone,
                            id:item.id
                        })
                    })
                    this.api.enterData.updateById(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "保存成功!",
                                duration:500,  
                                onClose(){
                                    sessionStorage.removeItem('queryAll');
                                    that.$router.push({name:"Qycensus"})
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
            this.$router.push({name:"Qycensus"})
        },
        // 选择省
        enterprises(id){
            let resultArr = this.provinceOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.proName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.cityCode = '';                 // 市
                    this.AddData.areaCode = '';                 // 县 区
                    this.AddData.countyCode = '';               // 乡
                    this.cityOptions = res.data.data;
                }
            })
        },
        // 选择市
        citys(id){
            let resultArr = this.cityOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.cityName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.areaCode = [];                 // 县 区
                    this.AddData.countyCode = '';               // 乡
                    this.countyOptions = res.data.data;
                }
            })
        },
        // 选择县 区
        countys(id){
            let resultArr = this.countyOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.areaName = resultArr[0].name;
            this.api.areaInfo.all({parentCode:id})
            .then(res=>{
                if(res.data.code == 200){
                    this.AddData.countyCode = '';               // 乡
                    this.villageOptions = res.data.data;
                }
            })
        },
        // 选择乡
        villages(id){
            let resultArr = this.villageOptions.filter((item)=>{
                return item.code == id;
            });
            this.AddData.countryName = resultArr[0].name;
        },
        handleRemove(file, fileList) {
            // this.imageUrl = ''
            this.noneBtnDealImg = fileList.length >= 1
            this.AddData.logo = '';
            this.fileList = [];
        },
        
        handleSuccess(response,file,fileList){
            this.AddData.logo = response.data[0].id
            // this.imageUrl = response.data[0].url;
            this.noneBtnDealImg = fileList.length >= 1
        },
        handleExceed(files, fileList, num) {
            this.$message.warning(`当前限制选择 1 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
        },
        handleClose(index){
            if(!this.urCompanytableData.includes(this.urcommontags[index])) {
                this.urcommontags.splice(index,1) 
                return
            }
            this.urCompanytableData.forEach((el,i) => {
                if(el.id == this.urcommontags[index].id){
                    console.log(el.phone,'phone');
                    setTimeout(() => {
                        this.$refs.urmultipleTable.toggleRowSelection(this.urCompanytableData[i],false);
                        this.urcommontags.splice(index,1)
                    },0)
                }
            })
        },
        // showInput() {
        //     this.inputVisible = true;
        //     this.$nextTick(()=> {
        //         this.$refs.saveTagInput.$refs.input.focus();
        //     });
        // },

        // handleInputConfirm() {
        //     let inputValue = this.inputValue;
        //     if (inputValue) {
        //         this.urcommontags.push({
        //             phone:inputValue
        //         })
        //     }
        //     this.inputVisible = false;
        //     this.inputValue = '';
        // }
    }
};
</script>
<style lang="less" scoped>
			// 上传视频样式开始
			.manyPhoto {
				width: 100%;
				height: 322rpx;
				margin: 24rpx 0 32rpx;
				box-sizing: border-box;

				image {
					width: 100%;
					height: 100%;
				}
			}

			.videoBox {
				width: 300px;
				height: 100rpx;
				margin: 24rpx 0 32rpx;
				box-sizing: border-box;

				video {
					width: 100%;
					height: 100%;

				}
			}

			.close-img {
				position: absolute;
				height: 36rpx !important;
				width: 36rpx !important;
				top: 16rpx;
				right: 16rpx;
				z-index:99;
			}

			// 上传视频样式结束
/deep/.avatar-uploader{
    .el-upload--picture-card{
        width: 80px !important;
        height: 80px !important;
        line-height: 90px !important;
    }
    .el-upload-list__item{
        width: 80px !important;
        height: 80px !important;
        line-height: 90px !important;
    }
}
/deep/.satar-uploader{
    .el-upload--picture-card{
        width: 200px !important;
        height: 30px !important;
        line-height: 30px !important;
    }
    .el-upload-list__item{
        width: 200px !important;
        height: 30px !important;
        line-height: 30px !important;
    }
    .el-upload-list__item{
        .el-icon-close{
            // display: none;
        }
    } 
}

/deep/.el-upload-list--picture-card .el-upload-list__item-status-label i{
    position: absolute;
    top:1px;
    left: 14px;
}
/deep/.el-form-item--small{
    .el-upload--picture-card{
        width: 236px;
        height: 140px;
        line-height: 150px;
    }
   .el-upload-list__item{
        width: 236px;
        height: 140px;
        line-height: 150px;
    }
}

// 关联企业弹窗样式开始
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
    
}
.text {
  width: 100%;
}
  // 分页
.pur-pages {
    width: 100%;
    height: 32px;
    margin: 42px 0 9px 9px;
}
// 关联企业弹窗样式结束

.bormp{
    width: 100%;
    min-height: 36px;
    border:1px solid #e4e4e4;
    border-top:0;
    padding: 6px 5px;
    box-sizing: border-box;
    .bpmslist{
        font-size: 14px;
        line-height: 24px;
        color: #333;
        box-sizing: border-box;
        cursor: pointer;
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
            width: 50%;
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

/deep/ .crsBanner {
    .el-form-item__label::after {
      content: '(最多1张)';
      display: block;
      font-size: 12px;
      color: #999;
      line-height: 12px;
    }
  }

/deep/ .allUpload {
    .el-form-item__content{
      // flex: none !important;
      display: flex;
    }
  }

/deep/ .el-upload-list__item{
  transition: none !important
}

/deep/ .disUoloadSty .el-upload--picture-card{
  display:none !important;   /* 上传logo按钮隐藏 */
}
/deep/ .disyeloadSty .el-upload--picture-card{
  display:none !important;   /* 上传营业执照按钮隐藏 */
}
/deep/ .dismploadSty .el-upload--picture-card{
  display:none !important;   /* 上传开户许可证按钮隐藏 */
}
/deep/ .disjcloadSty .el-upload--picture-card{
  display:none !important;   /* 上传身份证正面按钮隐藏 */
}
/deep/ .dissfzloadSty .el-upload--picture-card{
  display:none !important;   /* 上传身份证反面按钮隐藏 */
}


/deep/ .el-upload-list__item{
  width:300px;
  height:300px;
}

  .el-tag + .el-tag {
    margin-left: 10px;
  }
  .button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
  }
</style>