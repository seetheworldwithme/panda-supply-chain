<template>
    <!-- 编辑配置合同页面 -->
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
                <!-- 基本信息配置 -->
                <el-row>
                    <el-col :span="24">
                        <div class="see-top">基本信息配置</div>
                        <!-- <el-divider content-position="center">5555</el-divider> -->
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">文件主题名称</label>
                            </el-col>
                            <el-form-item prop="planName">
                                <el-col :span="7">
                                    <el-input v-model="AddData.planName" placeholder="请输入文件主题名称"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">流程类型</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.flowType" :label="1">对公</el-radio>
                                <el-radio v-model="AddData.flowType" :label="2">对私</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">自动归档</label>
                            </el-col>
                            <el-form-item prop="autoarchive">
                                <el-radio v-model="AddData.autoarchive" :label="true">归档</el-radio>
                                <el-radio v-model="AddData.autoarchive" :label="false">不归档</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">自动开启流程</label>
                            </el-col>
                            <el-form-item prop="autoinitiate">
                                <el-radio v-model="AddData.autoinitiate" :label="true">开启</el-radio>
                                <el-radio v-model="AddData.autoinitiate" :label="false">关闭</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签署有效截止时间</label>
                            </el-col>
                            <el-form-item prop="signvalidity">
                                <el-col :span="7">
                                    <el-input v-model="AddData.signvalidity" placeholder="请输入签署有效截止时间,超过则无法继续签署,单位（毫秒）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">文件续签提醒</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractremind" placeholder="请输入提醒时间（单位为小时，不填为不提醒）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">文件有效截止日期</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.contractvalidity" placeholder="请输入提醒时间（单位为小时，不填为不失效）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    <!-- 任务信息配置 -->
                        <!-- 认证配置项 -->
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">个人认证方式</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.personavailableauthtypes">
                                    <el-checkbox label="PSN_BANK4_AUTHCODE">个人银行卡四要素认证</el-checkbox>
                                    <el-checkbox label="PSN_TELECOM_AUTHCODE">个人运营商三要素认证</el-checkbox>
                                    <el-checkbox label="PSN_FACEAUTH_BYURL">个人刷脸认证</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">指定意愿认证方式</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.willtypes">
                                    <el-checkbox label="CODE_SMS">短信验证码</el-checkbox>
                                    <el-checkbox label="FACE_ZHIMA_XY">支付宝刷脸</el-checkbox>
                                    <el-checkbox label="FACE_TECENT_CLOUD_H5">腾讯云刷脸</el-checkbox>
                                    <el-checkbox label="FACE_FACE_LIVENESS_RECOGNITION">e签宝刷脸</el-checkbox>
                                    <el-checkbox label="FACE_WE_CHAT_FACE">微信小程序刷脸</el-checkbox>
                                    <el-checkbox label="FACE_AUDIO_VIDEO_DUAL">智能视频认证</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">企业实名认证方式</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.orgavailableauthtypes">
                                    <el-checkbox label="ORG_BANK_TRANSFER">组织机构对公账户打款认证</el-checkbox>
                                    <el-checkbox label="ORG_ZM_AUTHORIZE">企业芝麻认证</el-checkbox>
                                    <el-checkbox label="ORG_LEGAL_AUTHORIZE">组织机构法定代表人授权书签署认证</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <!-- 通知配置项 -->
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">通知开发者地址</label>
                            </el-col>
                            <el-form-item prop="noticedeveloperurl">
                                <el-col :span="7">
                                    <el-input v-model="AddData.noticedeveloperurl" placeholder="请输入地址"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">通知方式</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.noticetype">
                                    <el-checkbox label="1">短信</el-checkbox>
                                    <el-checkbox label="2">邮件</el-checkbox>
                                    <el-checkbox label="3">支付宝</el-checkbox>
                                    <el-checkbox label="4">钉钉</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <!-- 签署配置项 -->
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签署重定向url地址</label>
                            </el-col>
                            <el-form-item prop="redirecturl">
                                <el-col :span="7">
                                    <el-input v-model="AddData.redirecturl" placeholder="请输入签署完成后重定向地址"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署平台</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.signplatform">
                                    <el-checkbox label="1">开放服务h5</el-checkbox>
                                    <el-checkbox label="2">支付宝签</el-checkbox>
                                    <el-checkbox label="3">微信小程序</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">一键落章</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.batchdropseal" :label="true">显示</el-radio>
                                <el-radio v-model="AddData.batchdropseal" :label="false">不显示</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署页提交倒计时</label>
                            </el-col>
                            <el-form-item>
                                <el-col :span="7">
                                    <el-input v-model="AddData.countdown" placeholder="单位秒，取值范围0-999"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署重定向跳转时间</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.redirectdelaytime" :label="0">0s</el-radio>
                                <el-radio v-model="AddData.redirectdelaytime" :label="1">1s</el-radio>
                                <el-radio v-model="AddData.redirectdelaytime" :label="2">2s</el-radio>
                                <el-radio v-model="AddData.redirectdelaytime" :label="3">3s</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">是否e签盾签署</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.mobileshieldway" :label="1">不使用e签盾</el-radio>
                                <el-radio v-model="AddData.mobileshieldway" :label="2">使用e签盾</el-radio>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row>
                <!-- 签署方信息 -->
                <el-row>
                    <el-col :span="24">
                        <div class="see-top">签署方信息配置</div>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">平台自动签署</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.platformsign" :label="true">平台方自动签署</el-radio>
                                <el-radio v-model="AddData.platformsign" :label="false">平台用户签署</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签署顺序</label>
                            </el-col>
                            <el-form-item prop="signorder">
                                <el-col :span="7">
                                    <el-input-number v-model="AddData.signorder" :min="1"></el-input-number>       
                                    <el-tooltip class="item" effect="dark" content="签署方签署顺序，默认1,且不小于1，顺序越小越先处理" placement="top" style="margin:0 5%;">
                                        <i class="el-icon-question"></i>
                                    </el-tooltip>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">强制阅读时间</label>
                            </el-col>
                            <el-form-item prop="forceReadTime">
                                <el-col :span="7">
                                    <el-input v-model="AddData.forceReadTime" placeholder="请输入签署强制阅读时间, 单位为秒，不传默认为0，最大999"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <!-- 签署方账号信息 -->
                        <!-- <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签署人真实信息</label>
                            </el-col>
                            <el-form-item prop="signerAccount">
                                <el-col :span="7">
                                    <el-input v-model="AddData.signerAccount" placeholder="请输入签署人真实信息，手机号或邮箱"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row> -->
                        <!-- 签署人基本信息 -->
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署人角色信息</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.actorindentitytype" :label="1">个人</el-radio>
                                <el-radio v-model="AddData.actorindentitytype" :label="2">企业</el-radio>
                                <el-radio v-model="AddData.actorindentitytype" :label="3">法人</el-radio>
                                <el-radio v-model="AddData.actorindentitytype" :label="4">经办人</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署区位置</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.assignedposbean" :label="true">固定签署区位置且无法移动</el-radio>
                                <el-radio v-model="AddData.assignedposbean" :label="false">不固定签署区位置且签署时可自由移动</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署自动执行</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.autoexecute" :label="true">自动签署</el-radio>
                                <el-radio v-model="AddData.autoexecute" :label="false">不自动签署</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">页码信息</label>
                            </el-col>
                            <el-form-item prop="pospage">
                                <el-col :span="7">
                                    <el-input v-model="AddData.pospage" placeholder="页码可以'-'分割指定页码范围,输入all代表全部页码"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">X坐标</label>
                            </el-col>
                            <el-form-item prop="posx">
                                <el-col :span="7">
                                    <el-input v-model="AddData.posx" placeholder="请输入x坐标，坐标为印章中心点"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">Y坐标</label>
                            </el-col>
                            <el-form-item prop="posy">
                                <el-col :span="7">
                                    <el-input v-model="AddData.posy" placeholder="请输入Y坐标，坐标为印章中心点"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">印章类型</label>
                            </el-col>
                            <el-form-item>
                                <el-checkbox-group v-model="AddData.sealtype">
                                    <el-checkbox label="1">手绘印章</el-checkbox>
                                    <el-checkbox label="2">模版印章</el-checkbox>
                                </el-checkbox-group>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签章日期字体大小</label>
                            </el-col>
                            <el-form-item prop="fontsize">
                                <el-col :span="7">
                                    <el-input v-model="AddData.fontsize" placeholder="请输入签章日期字体大小(px)"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签章日期格式</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.format" label="yyyy年MM月dd日">yyyy年MM月dd日</el-radio>
                                <el-radio v-model="AddData.format" label="yyyy-MM-dd">yyyy-MM-dd</el-radio>
                                <el-radio v-model="AddData.format" label="yyyy/MM/dd">yyyy/MM/dd</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">页码信息</label>
                            </el-col>
                            <el-form-item prop="qpospage">
                                <el-col :span="7">
                                    <el-input v-model="AddData.qpospage" placeholder="签署自动执行选择自动执行时，并且需要展示签署日期，则需要指定日期盖章页码 ，默认当前页"></el-input>
                                </el-col> 
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">X坐标</label>
                            </el-col>
                            <el-form-item prop="qposx">
                                <el-col :span="7">
                                    <el-input v-model="AddData.qposx" placeholder="请输入x坐标 ，坐标点为日期左下角"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">Y坐标</label>
                            </el-col>
                            <el-form-item prop="qposy">
                                <el-col :span="7">
                                    <el-input v-model="AddData.qposy" placeholder="请输入y坐标 ，坐标点为日期左下角"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">添加签署日期</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.signdatebeantype" :label="0">禁止</el-radio>
                                <el-radio v-model="AddData.signdatebeantype" :label="1">必须</el-radio>
                                <el-radio v-model="AddData.signdatebeantype" :label="2">不限制</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <label style="color:#666;font-size:14px;line-height:30px;">签署类型</label>
                            </el-col>
                            <el-form-item>
                                <el-radio v-model="AddData.signtype" :label="0">不限</el-radio>
                                <el-radio v-model="AddData.signtype" :label="1">单页签署</el-radio>
                                <el-radio v-model="AddData.signtype" :label="2">骑缝签署</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">签署区宽度</label>
                            </el-col>
                            <el-form-item prop="width">
                                <el-col :span="7">
                                    <el-input v-model="AddData.width" placeholder="请输入签署区宽度（默认印章的宽度）"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>

                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">第三方流水号</label>
                            </el-col>
                            <el-form-item prop="thirdorderno">
                                <el-col :span="7">
                                    <el-input v-model="AddData.thirdorderno" placeholder="请输入第三方流水号"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row>
                <!-- 待签署文件信息 -->
                <el-row>
                    <el-col :span="24">
                        <div class="see-top">待签署文件信息</div>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">是否加密</label>
                            </el-col>
                            <el-form-item prop="encryption">
                                <el-radio v-model="AddData.encryption" :label="0">不加密</el-radio>
                                <el-radio v-model="AddData.encryption" :label="1">加密</el-radio>
                            </el-form-item>
                        </el-row>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">文档密码</label>
                            </el-col>
                            <el-form-item prop="filepassword">
                                <el-col :span="7">
                                    <el-input v-model="AddData.filepassword" placeholder="请输入文档密码"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row>
                <!-- 流程信息配置 -->
                <!-- <el-row>
                    <el-col :span="24">
                        <div class="see-top">流程信息配置</div>
                        <el-row>
                            <el-col :span="2.5">
                                <span style="color:red;">* </span>
                                <label style="color:#666;font-size:14px;line-height:30px;">计费隔离码</label>
                            </el-col>
                            <el-form-item prop="billisolationcode">
                                <el-col :span="7">
                                    <el-input v-model="AddData.billisolationcode" placeholder="请输入计费隔离码"></el-input>
                                </el-col>
                            </el-form-item>
                        </el-row>
                    </el-col>
                </el-row> -->
            </div>
        </div>
        <!-- 内容结束 -->
    </el-form>
</template>
<script>
export default {
    data() {
        return {
            /* 
            编辑配置合同数据
            */
            purTitle:'',
            // 编辑配置合同表单验证
            rules: {
                contractName: [
                    { required: true, message: '请输入合同名称', trigger: 'blur' }
                ],
                contractCode: [
                    { required: true, message: '请输入合同编号', trigger: 'change' },
                    // { pattern: /(?=.*\d)(?=.*\W+)(?=.*[A-Z])(?!.*\n).*$/, message: '合同编号格式不正确', trigger: 'blur'}
                ],
                sumPrice: [
                    { required: true, message: '请输入总价金额', trigger: 'blur' },
                    { pattern:/(^[1-9]([0-9]+)?(\.[0-9]{1,2})?$)|(^(0){1}$)|(^[0-9]\.[0-9]([0-9])?$)/, message: '金额大小在 1 - 9999999999', trigger: 'blur'}
                ],
            },
            // 编辑数据表单
            AddData: {
                noticetype:[],
                personavailableauthtypes:[],
                willtypes:[],
                orgavailableauthtypes:[],
                signplatform:[],
                sealtype:[],
                signtype:'',
                signorder:1,
            },
        };
    },
    created() {},
    mounted () {
        // 模块标题
        this.purTitle = this.$route.meta.title;
        // 接收主页单条数据
        const queryAll = sessionStorage.getItem('queryAll');
        if (!queryAll) { this.$router.replace({name:'ContractDeploy'}); return } // 保存/取消后 queryAll 被清,后退到本页会整页报错
        this.AddData = JSON.parse(queryAll);
    },
    methods: {
        /* 
        编辑配置合同表单提交功能
        */ 
        //编辑配置合同表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    this.AddData.noticetype = this.AddData.noticetype.join(",");
                    this.AddData.personavailableauthtypes = this.AddData.personavailableauthtypes.toString();
                    this.AddData.willtypes = this.AddData.willtypes.toString();
                    this.AddData.orgavailableauthtypes = this.AddData.orgavailableauthtypes.toString();
                    this.AddData.signplatform = this.AddData.signplatform.toString();
                    this.AddData.sealtype = this.AddData.sealtype.toString();
                    this.api.eqbconfig.edit(this.AddData)
                    .then(res=>{
                        if(res.data.code == 200){
                            var that = this
                            that.$message({
                                type: "success",
                                message: "编辑成功!",
                                duration:500,  
                                onClose(){
                                    sessionStorage.removeItem('queryAll');
                                    that.$router.push({name:"ContractDeploy"});
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
            this.$router.push({name:"ContractDeploy"});
        },
    },
    computed: {},
    watch: {},
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
.tips{
    float: left;
    width: 20px;
    height: 20px;
    font-size: 14px;
    cursor: pointer;
    border-radius: 50%;
    border:1px solid #828282;
    text-align: center;
    line-height:20px;
    color: #828282;
    box-sizing: border-box;
    margin-top:15px;
    margin-left: 8px;
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

// 弹窗样式结束
</style>