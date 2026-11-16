<template>
  <!-- 合同配置页面 -->
  <div>
    <!-- 头部开始 -->
    <el-row class="pur-top">
      <el-col :span="3"
        ><span class="pur-size">电子签章</span>
      </el-col>
      <el-col :span="21" class="pur-right">
        <el-button size="small" @click="openPur" v-if="$route.meta.delete">删除</el-button>
        <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">配置</el-button>
        <el-input
          @keyup.enter.native="search"
          size="small"
          placeholder="请输入内容"
          prefix-icon="el-icon-search"
          v-model="purSearch">
        </el-input>
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
          <el-table-column type="selection" fixed="left" width="40"></el-table-column>
          <el-table-column prop="id" v-if="show">
          </el-table-column>  
          <el-table-column
            label="操作" width="100" fixed="left">
            <template slot-scope="scope">
              <el-button @click="Clicksee(scope.row)" size="mini" style="margin:0 0 5% 10px;">查看</el-button>
              <el-button size="mini" @click="EditPurs(scope.row)" v-if="$route.meta.edit">编辑</el-button>
            </template>
          </el-table-column>
          <el-table-column prop="planName" label="文件主题名称" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column label="流程类型" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.flowType == 1">对公</span>
              <span v-if="scope.row.flowType == 2">对私</span>
            </template>
          </el-table-column>
          <el-table-column prop="filepassword" label="文档密码" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="autoarchive" label="是否自动归档" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.autoarchive == true">归档</span>
              <span v-else>不归档</span>
            </template>
          </el-table-column>
          <el-table-column prop="autoinitiate" label="是否自动开启流程" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.autoinitiate == true">开启</span>
              <span v-else>关闭</span>
            </template>
          </el-table-column>
          <el-table-column prop="contractremind" label="文件续签提醒" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.contractremind">{{scope.row.contractremind}}小时</span>
            </template>
          </el-table-column>
          <el-table-column prop="contractvalidity" label="文件有效截止日期" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.contractvalidity">{{scope.row.contractvalidity}}小时</span>
            </template>
          </el-table-column>
          <el-table-column prop="batchdropseal" label="是否一键落章" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.batchdropseal == true">显示</span>
              <span v-else>不显示</span>
            </template>
          </el-table-column>
          <el-table-column prop="countdown" label="签署页提交倒计时" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.countdown">{{scope.row.countdown}}秒</span>
            </template>
          </el-table-column>
          <!-- <el-table-column prop="billisolationcode" label="计费隔离码" :show-overflow-tooltip="true">
          </el-table-column> -->
          <el-table-column prop="mobileshieldway" label="是否使用e签盾签署" :show-overflow-tooltip="true" width="160">
            <template slot-scope="scope">
              <span v-if="scope.row.mobileshieldway == 1">不使用e签盾</span>
              <span v-else>使用e签盾</span>
            </template>
          </el-table-column>
          <el-table-column prop="noticedeveloperurl" label="通知开发者地址" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="noticetype" label="通知方式" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.noticetype" :key="key">
                <span v-if="item == 1">短信, </span>
                <span v-if="item == 2">邮件, </span>
                <span v-if="item == 3">支付宝, </span>
                <span v-if="item == 4">钉钉</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="个人认证方式" :show-overflow-tooltip="true" width="170">
            <template slot-scope="scope">
              <div v-for="(item,key) in scope.row.personavailableauthtypes" :key="key">
                <span v-if="item == 'PSN_BANK4_AUTHCODE'">个人银行卡四要素认证</span>
                <span v-if="item == 'PSN_TELECOM_AUTHCODE'">个人运营商三要素认证</span>
                <span v-if="item == 'PSN_FACEAUTH_BYURL'">个人刷脸认证</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="指定意愿认证方式" :show-overflow-tooltip="true" width="170">
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.willtypes" :key="key">
                <span v-if="item == 'CODE_SMS'">短信验证码, </span>
                <span v-if="item == 'FACE_ZHIMA_XY'">支付宝刷脸, </span>
                <span v-if="item == 'FACE_TECENT_CLOUD_H5'">腾讯云刷脸, </span>
                <span v-if="item == 'FACE_FACE_LIVENESS_RECOGNITION'">e签宝刷脸, </span>
                <span v-if="item == 'FACE_WE_CHAT_FACE'">微信小程序刷脸, </span>
                <span v-if="item == 'FACE_AUDIO_VIDEO_DUAL'">智能视频认证</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="企业实名认证方式" :show-overflow-tooltip="true" width="170">
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.orgavailableauthtypes" :key="key">
                <span v-if="item == 'ORG_BANK_TRANSFER'">组织机构对公账户打款认证, </span>
                <span v-if="item == 'ORG_ZM_AUTHORIZE'">企业芝麻认证, </span>
                <span v-if="item == 'ORG_LEGAL_AUTHORIZE'">组织机构法定代表人授权书签署认证</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="countdown" label="签署重定向跳转时间" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.redirectdelaytime == 0">0 秒</span>
              <span v-if="scope.row.redirectdelaytime == 1">1 秒</span>
              <span v-if="scope.row.redirectdelaytime == 2">2 秒</span>
              <span v-if="scope.row.redirectdelaytime == 3">3 秒</span>
            </template>
          </el-table-column>
          <el-table-column prop="redirecturl" label="签署重定向url地址" sortable :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column label="签署平台" :show-overflow-tooltip="true" >
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.signplatform" :key="key">
                <span v-if="item == '1'">开放服务h5, </span>
                <span v-if="item == '2'">支付宝签, </span>
                <span v-if="item == '3'">微信小程序</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="signvalidity" label="签署有效截止时间" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.signvalidity">{{scope.row.signvalidity}} 毫秒</span>
            </template>
          </el-table-column>
          <el-table-column prop="platformsign" label="平台自动签署" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.platformsign == true">平台方自动签署</span>
              <span v-else>对接平台的用户签署</span>
            </template>
          </el-table-column>
          <el-table-column prop="actorindentitytype" label="签署人角色信息" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.actorindentitytype" :key="key">
                <span v-if="item == 1">个人, </span>
                <span v-if="item == 2">企业, </span>
                <span v-if="item == 3">法人, </span>
                <span v-if="item == 4">经办人</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="assignedposbean" label="签署区位置" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.assignedposbean == true">固定签署区位置且无法移动</span>
              <span v-else>不固定签署区位置且签署时可自由移动</span>
            </template>
          </el-table-column>
          <el-table-column prop="autoexecute" label="签署自动执行" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.autoexecute == true">自动签署</span>
              <span v-else>不自动签署</span>
            </template>
          </el-table-column>
          <el-table-column prop="pospage" label="页码信息" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="posx" label="X坐标" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="posy" label="Y坐标" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column label="印章类型" :show-overflow-tooltip="true" >
            <template slot-scope="scope">
              <span v-for="(item,key) in scope.row.sealtype" :key="key">
                <span v-if="item == '1'">手绘印章, </span>
                <span v-if="item == '2'">模版印章</span>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="fontsize" label="签章字体大小" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="format" label="签章日期格式" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.format == 'yyyy年MM月dd日'">yyyy年MM月dd日</span>
              <span v-if="scope.row.format == 'yyyy-MM-dd'">yyyy-MM-dd</span>
              <span v-if="scope.row.format == 'yyyy/MM/dd'">yyyy/MM/dd</span>
            </template>
          </el-table-column>
          <el-table-column prop="qpospage" label="页码信息" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="qposx" label="X坐标" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="qposy" label="Y坐标" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="signdatebeantype" label="添加签署日期" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.signdatebeantype == 0">禁止</span>
              <span v-if="scope.row.signdatebeantype == 1">必须</span>
              <span v-if="scope.row.signdatebeantype == 2">不限制</span>
            </template>
          </el-table-column>
          <el-table-column label="签署类型" :show-overflow-tooltip="true" width="100">
            <template slot-scope="scope">
              <span v-if="scope.row.signtype == 0">不限</span>
              <span v-if="scope.row.signtype == 1">单页签署</span>
              <span v-if="scope.row.signtype == 2">骑缝签署</span>
            </template>
          </el-table-column>
          <el-table-column prop="width" label="签署区宽度" :show-overflow-tooltip="true">
            <template slot-scope="scope">
              <span v-if="scope.row.width">{{scope.row.width}} px</span>
            </template>
          </el-table-column>
          <el-table-column prop="signorder" label="签署顺序" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="forceReadTime" label="强制阅读时间(s)" :show-overflow-tooltip="true">
          </el-table-column>
          <el-table-column prop="thirdorderno" label="第三方流水号" :show-overflow-tooltip="true">
          </el-table-column>
          <div slot="empty" class="dataAvailable" v-if="dataAva">
              <div class="dataAva">
                <img src="../../assets/zanwushuju.png" alt="">
                <p>暂无数据，您可以配置合同试试</p>
                <el-button size="small" style="margin-left: 24px" @click="addPurs" v-if="$route.meta.add">配置</el-button>
              </div>
          </div>
        </el-table>
      </div>
      <!-- 分页 -->
      <div class="pur-pages">
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
      <el-drawer :title="logo" :visible.sync="drawer" size="700px">
        <div class="pur-tiao"></div>
        <div class="pur-drawerNav">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane label="基本信息" name="first">
              <el-row>
                <el-col :span="24">
                  <div class="see-top">基本信息配置</div>
                  <ul class="see-lable">
                    <li>
                      <label>文件主题名称：</label>
                      <span>{{seeData.planName}}</span>
                    </li>
                    <li>
                      <label>流程类型：</label>
                        <span v-if="seeData.flowType == 1">对公</span>
                        <span v-if="seeData.flowType == 2">对私</span>
                    </li>
                    <li>
                      <label>是否自动归档：</label>
                      <span v-if="seeData.autoarchive == true">归档</span>
                      <span v-else>不归档</span>
                    </li>
                    <li>
                      <label>是否自动开启流程：</label>
                      <span v-if="seeData.autoinitiate == true">开启</span>
                      <span v-else>关闭</span>
                    </li>
                    <li>
                      <label>文件续签提醒：</label>
                      <span>{{seeData.contractremind}} 小时</span>
                    </li>
                    <li>
                      <label>文件有效截止日期：</label>
                      <span>{{seeData.contractvalidity}} 小时</span>
                    </li>
                    <li>
                      <label>是否一键落章：</label>
                      <span v-if="seeData.batchdropseal == true">显示</span>
                      <span v-else>不显示</span>
                    </li>
                    <li>
                      <label>签署页提交倒计时：</label>
                      <span>{{seeData.countdown}} 秒</span>
                    </li>
                    <li>
                      <label>计费隔离码：</label>
                      <span>{{seeData.billisolationcode}}</span>
                    </li>
                    <li>
                      <label>是否使用e签盾签署：</label>
                      <span v-if="seeData.mobileshieldway == 1">不使用e签盾</span>
                      <span v-else>使用e签盾</span>
                    </li>
                    <li>
                      <label>通知开发者地址：</label>
                      <span>{{seeData.noticedeveloperurl}}</span>
                    </li>
                    <li>
                      <label>通知方式：</label>
                      <span v-for="(item,key) in seeData.noticetype" :key="key">
                        <span v-if="item == '1'">短信, </span>
                        <span v-if="item == '2'">邮件, </span>
                        <span v-if="item == '3'">支付宝, </span>
                        <span v-if="item == '4'">钉钉</span>
                      </span>
                    </li>
                    <li>
                      <label>个人认证方式：</label>
                      <span v-for="(item,key) in seeData.personavailableauthtypes" :key="key">
                        <span v-if="item == 'PSN_BANK4_AUTHCODE'">个人银行卡四要素认证, </span>
                        <span v-if="item == 'PSN_TELECOM_AUTHCODE'">个人运营商三要素认证, </span>
                        <span v-if="item == 'PSN_FACEAUTH_BYURL'">个人刷脸认证</span>
                      </span>
                    </li>
                    <li>
                      <label>指定意愿认证方式：</label>
                      <span v-for="(item,key) in seeData.willtypes" :key="key">
                        <span v-if="item == 'CODE_SMS'">短信验证码, </span>
                        <span v-if="item == 'FACE_ZHIMA_XY'">支付宝刷脸, </span>
                        <span v-if="item == 'FACE_TECENT_CLOUD_H5'">腾讯云刷脸, </span>
                        <span v-if="item == 'FACE_FACE_LIVENESS_RECOGNITION'">e签宝刷脸, </span>
                        <span v-if="item == 'FACE_WE_CHAT_FACE'">微信小程序刷脸, </span>
                        <span v-if="item == 'FACE_AUDIO_VIDEO_DUAL'">智能视频认证</span>
                      </span>
                    </li>
                    <li>
                      <label>企业实名认证方式：</label>
                      <span v-for="(item,key) in seeData.orgavailableauthtypes" :key="key">
                        <span v-if="item == 'ORG_BANK_TRANSFER'">组织机构对公账户打款认证, </span>
                        <span v-if="item == 'ORG_ZM_AUTHORIZE'">企业芝麻认证, </span>
                        <span v-if="item == 'ORG_LEGAL_AUTHORIZE'">组织机构法定代表人授权书签署认证</span>
                      </span>
                    </li>
                    <li>
                      <label>签署重定向跳转时间：</label>
                        <span v-if="seeData.redirectdelaytime == 0">0 秒</span>
                        <span v-if="seeData.redirectdelaytime == 1">1 秒</span>
                        <span v-if="seeData.redirectdelaytime == 2">2 秒</span>
                        <span v-if="seeData.redirectdelaytime == 3">3 秒</span>
                    </li>
                    <li>
                      <label>签署重定向url地址：</label>
                      <span>{{seeData.redirecturl}}</span>
                    </li>
                    <li>
                      <label>签署平台：</label>
                      <span v-for="(item,key) in seeData.signplatform" :key="key">
                        <span v-if="item == '1'">开放服务h5, </span>
                        <span v-if="item == '2'">支付宝签, </span>
                        <span v-if="item == '3'">微信小程序</span>
                      </span>
                    </li>
                    <li>
                      <label>签署有效截止时间：</label>
                      <span>{{seeData.signvalidity}} 毫秒</span>
                    </li>
                  </ul>

                  <div class="see-top">签署方信息配置</div>
                  <ul class="see-lable">
                    <li>
                      <label>平台自动签署：</label>
                      <span v-if="seeData.platformsign == true">平台方自动签署</span>
                      <span v-else>平台用户签署</span>
                    </li>
                    <li>
                      <label>签署人角色信息：</label>
                      <span v-for="(item,key) in seeData.actorindentitytype" :key="key">
                        <span v-if="item == '1'">个人, </span>
                        <span v-if="item == '2'">企业, </span>
                        <span v-if="item == '3'">法人</span>
                        <span v-if="item == '4'">经办人</span>
                      </span>
                    </li>
                    <li>
                      <label>签署区位置：</label>
                      <span v-if="seeData.assignedposbean == true">固定签署区位置且无法移动</span>
                      <span v-else>不固定签署区位置且签署时可自由移动</span>
                    </li>
                    <li>
                      <label>签署自动执行：</label>
                      <span v-if="seeData.autoexecute == true">自动签署</span>
                      <span v-else>不自动签署</span>
                    </li>
                    <li>
                      <label>页码信息：</label>
                      <span>{{seeData.pospage}}</span>
                    </li>
                    <li>
                      <label>X坐标：</label>
                      <span>{{seeData.posx}}</span>
                    </li>
                    <li>
                      <label>Y坐标：</label>
                      <span>{{seeData.posy}}</span>
                    </li>
                    <li>
                      <label>印章类型：</label>
                      <span v-for="(item,key) in seeData.sealtype" :key="key">
                        <span v-if="item == '1'">手绘印章, </span>
                        <span v-if="item == '2'">模版印章, </span>
                      </span>
                    </li>
                    <li>
                      <label>签章字体大小：</label>
                      <span>{{seeData.fontsize}}</span>
                    </li>
                    <li>
                      <label>签章日期格式：</label>
                        <span v-if="seeData.format == 'yyyy年MM月dd日'">yyyy年MM月dd日</span>
                        <span v-if="seeData.format == 'yyyy-MM-d'">yyyy-MM-d</span>
                        <span v-if="seeData.format == 'yyyy/MM/dd'">yyyy/MM/dd</span>
                    </li>
                    <li>
                      <label>页码信息：</label>
                      <span>{{seeData.qpospage}}</span>
                    </li>
                    <li>
                      <label>X坐标：</label>
                      <span>{{seeData.qposx}}</span>
                    </li>
                    <li>
                      <label>Y坐标：</label>
                      <span>{{seeData.qposy}}</span>
                    </li>
                    <li>
                      <label>添加签署日期：</label>
                        <span v-if="seeData.signdatebeantype == 0">禁止</span>
                        <span v-if="seeData.signdatebeantype == 1">必须</span>
                        <span v-if="seeData.signdatebeantype == 2">不限制</span>
                    </li>
                    <li>
                      <label>签署类型：</label>
                        <span v-if="seeData.signdatebeantype == 0">不限</span>
                        <span v-if="seeData.signdatebeantype == 1">单页签署</span>
                        <span v-if="seeData.signdatebeantype == 2">骑缝签署</span>
                    </li>
                    <li>
                      <label>签署区宽度：</label>
                      <span>{{seeData.width}} px</span>
                    </li>
                    <li>
                      <label>签署顺序：</label>
                      <span>{{seeData.signorder}}</span>
                    </li>
                    <li>
                      <label>签署强制阅读时间：</label>
                      <span>{{seeData.forceReadTime}} s</span>
                    </li>
                    <li>
                      <label>第三方流水号：</label>
                      <span>{{seeData.thirdorderno}}</span>
                    </li>
                  </ul>

                  <div class="see-top">待签署文件信息</div>
                  <ul class="see-lable">
                    <li>
                      <label>文档密码：</label>
                      <span>{{seeData.filepassword}}</span>
                    </li>
                    <li>
                      <label>是否加密：</label>
                      <span v-if="seeData.encryption == 0">不加密</span>
                      <span v-else>加密</span>
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
          </el-tabs>
          <div class="xia">《</div>
        </div>
      </el-drawer>
    </div>
    <!-- 内容结束 -->
  </div>
</template>
<script>
export default {
  data() {
    return {
      current:1,              // 当前页数
      size:50,                // 显示条数
      total:0,                // 总条数
      delmul:[],              // 批量删除
      show:false,
      purTitle: "",           // 模块标题
      purSearch: "",          // 搜索
      tableData: [],          // 表格数据
      dataAva:false,
      loading:false,
      multipleSelection: [],  // 选中数据

      // 查看
      drawer: false,          // 查看抽屉
      logo: "",               // 抽屉标题
      activeName: "first",    // 查看tab栏
      seeData:{},             // 查看数据
    };
  },
  created() {},
  mounted(){
    this.purTitle = this.$route.meta.title;
    this.logo = sessionStorage.getItem('companyName') // 企业名称
    this.getpurcon();
  },
  methods: {
    // 表格数据
    getpurcon(){
      this.loading = true;
      // 传参
      var listByPageData = {
        searchStr:this.purSearch,
        page:{
          current:this.current,
          size:this.size
        }
      }
      // 渲染表格
      this.api.eqbconfig.all(listByPageData)
      .then(res=>{
        if(res.data.code == 200){
          this.loading = false;
          if(res.data.data.records != ""){
            let arr = res.data.data.records;
            for(let key in arr){
              
              arr[key].noticetype = arr[key].noticetype.split(",");
              arr[key].personavailableauthtypes = arr[key].personavailableauthtypes.split(",");
              arr[key].willtypes = arr[key].willtypes.split(",");
              arr[key].orgavailableauthtypes = arr[key].orgavailableauthtypes.split(",");
              arr[key].signplatform = arr[key].signplatform.split(",");
              arr[key].sealtype = arr[key].sealtype.split(",");
            }
            this.tableData = arr;
            this.total = res.data.data.total;
          }else{
            this.dataAva = true;
          }
        }
      })
    },
    // 搜索功能
    search() {
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
    // 新建配置合同
    addPurs(){
        this.$router.push({name:"AddContractDeploy"})
    },
    // 编辑合同配置
    EditPurs(row){
      this.api.eqbconfig.one({id:row.id})
      .then(res=>{
        if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.id = row.id;
          purdata.noticetype = purdata.noticetype.split(",");
          purdata.personavailableauthtypes = purdata.personavailableauthtypes.split(",");
          purdata.willtypes = purdata.willtypes.split(",");
          purdata.orgavailableauthtypes = purdata.orgavailableauthtypes.split(",");
          purdata.signplatform = purdata.signplatform.split(",");
          purdata.sealtype = purdata.sealtype.split(",");
          sessionStorage.setItem('queryAll',JSON.stringify(purdata))
          this.$router.push({name:"EditContractDeploy"})
        }
      })
    },
    // 查看
    Clicksee(row) {
      this.drawer = true;
      this.activeName = 'first';
      this.api.eqbconfig.one({id:row.id})
      .then(res=>{
        if(res.data.code == 200){
          let purdata = res.data.data;
          purdata.noticetype = purdata.noticetype.split(",");
          purdata.personavailableauthtypes = purdata.personavailableauthtypes.split(",");
          purdata.willtypes = purdata.willtypes.split(",");
          purdata.orgavailableauthtypes = purdata.orgavailableauthtypes.split(",");
          purdata.signplatform = purdata.signplatform.split(",");
          purdata.sealtype = purdata.sealtype.split(",");
          this.seeData = purdata;
        }
      })
    },
    // 抽屉
    handleClick(tab) {
      // console.log(tab)
      // 判断tab栏选中状态，请求相应的接口
      if (tab.label == "基本信息") {
        this.allone(this.seeId,this.seeId);
      }
      if (tab.label == "动态") {
        // console.log(2);
      }
      if (tab.label == "展示图") {
        this.sears(this.seeId,this.orderId);
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
          this.$confirm('确定是否删除,请谨慎操作 是否继续?', '提示', {
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning'
          }).then(() => {
            this.api.eqbconfig.del(this.delmul)
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
  position: relative;
  .pur-table{
    width:100%;
    height: 94%;
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
      left:50%;
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
    width: 22%;
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
// 查看抽屉样式结束
</style>