<template>
<!-- 注册申请页面 -->
    <div>
        <div class="login-concent">
            <div @click="Gologin" style="cursor: pointer;">
                <div class="img-lg">
                    <img src="../../src/assets/lg.svg" alt="">
                </div>
                <div class="img-size"><span>云管理</span></div>
            </div>
        </div>
        <div class="stepAll">
            <el-steps :active="active" align-center process-status="wait" finish-status="finish" v-if="stepShow" style="max-width:800px;margin:auto;">
                <el-step title="注册账号信息" icon="el-icon-ali1"></el-step>
                <el-step title="登记成功" icon="el-icon-ali2"></el-step>
                <el-step title="申请管理员" icon="el-icon-ali3" v-if="activeShow"></el-step>
                <el-step title="加入企业" icon="el-icon-ali3" v-if="activeShowtwo"></el-step>
            </el-steps>
        </div>
        <div class="stepBox"> 
            <div class="stepAllBox">
                <!-- 设置账号信息页面 -->
                <el-form v-if="active == 1" :model="loginForm" :rules="rule" ref="loginForm" label-width="auto" :hide-required-asterisk="false" class="concentBox">
                    <el-form-item prop="phone" label="手机号">
                    <el-input v-model="loginForm.phone" placeholder="请输入11位手机号码" autocomplete="off"></el-input>
                    </el-form-item>
                    <!-- 验证码开始 -->
                    <el-form-item prop="code" label="验证码" style="margin-bottom:1.5%;">
                    <el-input v-model="loginForm.code" placeholder="请输入验证码" autocomplete="off">
                        <template slot="suffix">
                        <span v-if="!isDisabled" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:16px;cursor: pointer;" ></span>
                        <span v-if="isDisabled" :class="color" v-html="codephone" style="margin-right:16px;cursor: pointer;"></span>
                        </template>
                    </el-input>
                    </el-form-item>
                    <!-- 验证码结束 -->
                    <el-form-item>
                        <el-checkbox v-model="loginForm.checked" style="float:left;">
                            <span style="font-size:12px;color:#333;font-weight: 400;cursor:pointer;">注册即代表您同意   
                                <span style="line-height: 18px;color:#CF1724; font-weight:500;cursor: pointer;" @click="isService()">服务协议</span> 和
                                <span style="line-height: 18px;color:#CF1724; font-weight:500;cursor: pointer;" @click="isService()">隐私政策</span>
                            </span>
                        </el-checkbox>
                    </el-form-item>
                    <el-form-item style="margin-bottom:0;margin-top:40px;">
                    <el-button type="primary" @click="submitFormLogin('loginForm')" style="width:100%;">下 一 步</el-button>
                    </el-form-item>
                </el-form>
                <!-- 登记成功页面 -->
                <div v-if="active ==2" class="concentBox">
                    <div class="secussIcon"><i class="el-icon-alidui"></i></div>
                    <h3>登记成功</h3>
                    <p>您已完成登记，请继续操作完成注册</p>
                    <div class="secussButton" style="display:flex;justify-content: space-between;">
                        <div class="elButton" v-if="AdjustForm" @click="administrators" style="margin-right:64px;">申请企业管理员</div>
                        <div class="elButton" @click="enterprise">加 入 企 业</div>
                    </div>
                </div>
                <!-- 已有管理员，加入企业 -->
                <div v-if="active ==5" class="concentBox">
                    <div class="secussIcon"><i class="el-icon-aligantanhao" style="color:#FB8138;"></i></div>
                    <h3>已有管理员</h3>
                    <p>企业已经申请过管理员，请联系管理员获取邀请码加入企业</p>
                    <div class="secussButton">
                        <div style="margin:auto;" class="elButton" @click="enterprise">加 入 企 业</div>
                    </div>
                </div>
                <!-- 正在审核 -->
                <div v-if="active ==6" class="concentBox">
                    <div class="secussIcon"><img src="../assets/kaifa.jpg" alt=""></div>
                    <h3>正在审核</h3>
                    <p>已的资料已提交成功，平台正在加速审核...</p>
                    <span style="color:#828282;"><span style="color:#CF1724;"> {{Jump}}s </span> 后自动跳转登录页</span>
                </div>
                <!-- 加入企业页面 -->
                <el-form v-if="active == 4" :model="businessForm" :rules="rulebusiness" ref="businessForm" label-width="auto" :hide-required-asterisk="false" class="concentBox">
                    <el-form-item prop="code" label="邀请码">
                        <el-input v-model="businessForm.code" placeholder="若未取得请联系企业管理员" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="username" label="本人姓名">
                        <el-input v-model="businessForm.username" placeholder="姓名必须2-4位中文" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="account" label="账号">
                        <el-input v-model="businessForm.account" placeholder="请输入长度为4-16位字母用于登录" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="password" label="登录密码">
                    <el-input type="password" v-model="businessForm.password" placeholder="长度为8-16位 包括数字和大小写字母" show-password autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item style="margin-bottom:0;margin-top:40px;">
                        <el-button type="primary" @click="submitFormBusiness('businessForm')" style="width:100%;">下 一 步</el-button>
                    </el-form-item>
                </el-form>
                <!-- 申请企业管理员页面 -->
                <el-form v-if="active == 3" :model="manageForm" :rules="managerule" ref="manageForm" label-width="auto" :hide-required-asterisk="false" class="concentBox">
                    <el-form-item prop="enterpriseName" label="企业名称">
                        <el-input v-model="manageForm.enterpriseName" placeholder="请输入企业全称" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="rateNumber" label="税号">
                        <el-input v-model="manageForm.rateNumber" placeholder="请输入税号" autocomplete="off"></el-input>
                    </el-form-item>

                    <el-form-item prop="legalPerson" label="法人">
                        <el-input v-model="manageForm.legalPerson" placeholder="姓名必须2-4位中文" autocomplete="off"></el-input>
                    </el-form-item>
                    <!-- <el-form-item prop="controllerPerson" label="实际控制人">
                        <el-input v-model="manageForm.controllerPerson" placeholder="请输入实际控制人" autocomplete="off"></el-input>
                    </el-form-item> -->
                    <el-form-item prop="contactsPerson" label="主要联系人">
                        <el-input v-model="manageForm.contactsPerson" placeholder="请输入主要联系人" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="adminName" label="管理员姓名">
                        <el-input v-model="manageForm.adminName" placeholder="姓名必须2-4位中文" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="account" label="账号">
                        <el-input v-model="manageForm.account" placeholder="请输入长度为4-16位字母用于账号" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="password" label="登录密码">
                        <el-input type="password" v-model="manageForm.password" placeholder="长度为8-16位 包括数字和大小写字母" show-password autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="qrUrl" label="营业执照">
                        <el-upload
                            class="img-uploader"
                            :action="http +'/saas-common/upload/upload'"
                            :headers="header"
                            :show-file-list="false"
                            name='file'
                            :on-success="handleQrSuccess"
                            :before-upload="beforelogoUpload">
                            <img v-if="manageForm.qrUrl" :src="manageForm.qrUrl" class="avatar"/>
                            <i v-else class="el-icon-plus avatar-uploader-icon"></i>
                        </el-upload>
                    </el-form-item>
                    <el-form-item label="企业类别">
                        <el-select v-model="manageForm.type" placeholder="请选择企业类别" style="width:100%;">
                            <el-option
                            v-for="item in typeOptions"
                            :key="item.id"
                            :label="item.name"
                            :value="item.id">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="经营地址">
                        <el-input v-model="manageForm.address" placeholder="请输入经营地地址" autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item label="注册地址">
                        <el-row>
                            <el-col :span="6">
                                <el-select @change="enterprises($event)" v-model="manageForm.provinceCode" placeholder="省" style="width:90%;margin-right:10%;">
                                    <el-option
                                    v-for="item in provinceOptions"
                                    :key="item.code"
                                    :label="item.name"
                                    :value="item.code">
                                    </el-option>
                                </el-select>
                            </el-col>
                            <el-col :span="6">
                                <el-select v-model="manageForm.cityCode" placeholder="市" style="width:90%;margin-right:10%;" @change="citys($event)">
                                    <el-option
                                    v-for="item in cityOptions"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.code">
                                    </el-option>
                                </el-select>
                            </el-col>
                            <el-col :span="6">
                                <el-select v-model="manageForm.areaCode" placeholder="县 / 区" style="width:90%;margin-right:10%;" @change="countys($event)">
                                    <el-option
                                    v-for="item in countyOptions"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.code">
                                    </el-option>
                                </el-select>
                            </el-col>
                            <el-col :span="6">
                                <el-select v-model="manageForm.countyCode" placeholder="乡" style="width:100%;" @change="villages($event)">
                                    <el-option
                                    v-for="item in villageOptions"
                                    :key="item.value"
                                    :label="item.name"
                                    :value="item.code">
                                    </el-option>
                                </el-select>
                            </el-col>
                        </el-row>
                    </el-form-item>
                    <el-form-item style="margin-bottom:0;margin-top:40px;">
                        <el-button type="primary" @click="submitFormmanage('manageForm')" style="width:100%;">下 一 步</el-button>
                    </el-form-item>
                </el-form>
            </div>
        </div>
        <div class="copyright-bottom">
            <div class="copyright"><span>Copyright © 山西海云启辰科技有限公司 版权所有</span><p><a target="_blank" href="https://beian.miit.gov.cn/">晋ICP备2022001434号</a></p></div>
        </div>
    </div>
</template>
<script>
export default {
  data() {
    return {
        AdjustForm:null,  
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
        http:'',
        Jump:5,
        active:1,
        activeShow:true,        // 管理员显示
        activeShowtwo:false,    // 企业隐藏
        stepShow:true,          // 步骤条显示与隐藏
        // 设置账号信息
        loginForm: {
            phone: '',
            password:'',
            code:'',            // 表单提交的验证码
            checked:false,      // 勾选同意协议
        },
        isDisabled: false,      // 控制按钮是否可以点击（false:可以点击，true:不可点击）
        color:'onPhone',
        codephone:'获取验证码',
        // 密码登录表单验证
        rule: {
            username: [
                { required: true, message: '请输入手机号', trigger: 'blur' },
                { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'blur'}
            ],
            code: [
                { required: true, message: '请输入验证码', trigger: 'blur' },
                { max: 6, message: '请输入6位验证码', trigger: 'blur' }
            ]
        },
        hasAdmin:'',        // 是否有管理员，true为有，false为没有        
        // 申请管理员表单
        manageForm:{
            userId:'',      // 企业id
            adminName:'',
            enterpriseName:'',
            account:'',
            rateNumber:'',
            address:'',
            legalPerson:'',
            controllerPerson:'',
            contactsPerson:'',
            id:'',
            addressone:'',
            provinceCode:'',
            cityCode:'',
            areaCode:'',
            countyCode:'',
            proName:'',
            areaInfo:'',
            areaName:'',
            countryName:'',
            type:'1',
            qrUrl:'',
            filelist:[],
        },
        header:{},
        
        // 行政区划
        provinceOptions:[],             // 省
        cityOptions: [],                // 市
        countyOptions: [],              // 县 区
        villageOptions: [],             // 乡
        // 申请管理员表单验证 
        managerule: {
            adminName: [
                { required: true, message: '请输入管理员姓名', trigger: 'blur' },
                { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
            ],
            enterpriseName: [
                { required: true, message: '请输入企业全称', trigger: 'blur' }
            ],
            account: [
                { required: true, message: '请输入账号', trigger: 'blur' },
                { pattern:/^[a-zA-Z]{4,16}$/, message: '请输入字母 4-16 位账号', trigger: 'blur'},
            ],
            rateNumber: [
                { required: true, message: '请输入税号', trigger: 'blur' },
                { pattern:/(^(?:(?![IOZSV])[\dA-Z]){2}\d{6}(?:(?![IOZSV])[\dA-Z]){10}$)|(^\d{15}$)/, message: '请输入正确的税号', trigger: 'blur'}
            ],
            address: [
                { required: true, message: '请输入经营地址', trigger: 'blur' }
            ],
            addressone: [
                { required: true, message: '请输入注册地址', trigger: 'blur' }
            ],
            legalPerson: [
                { required: true, message: '请输入法人姓名', trigger: 'blur' },
                { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
            ],
            controllerPerson: [
                { required: true, message: '请输入实际控制人姓名', trigger: 'blur' },
                { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
            ],
            contactsPerson: [
                { required: true, message: '请输入主要联系人姓名', trigger: 'blur' },
                { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
            ],
            password: [
                { required: true, message: '请输入 8 到 16 位数字和大小写字母', trigger: 'blur' },
                { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
            ],
            type: [
                { required: true, message: '请选择企业类别', trigger: 'blur' }
            ],
            qrUrl: [
                { required: true, message: '请上传营业执照', trigger: 'blur' }
            ],
        },
        // 加入企业表单
        businessForm:{
            id:'',
            code:'',
            username:'',
            account:''
        },
        //  加入企业表单验证
        rulebusiness: {
            code: [
                { required: true, message: '请输入邀请码', trigger: 'blur' }
            ],
            username: [
                { required: true, message: '请输入姓名', trigger: 'blur' },
                { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'blur'}
            ],
            account: [
                { required: true, message: '请输入账号', trigger: 'blur' },
                { pattern:/^[a-zA-Z]{4,16}$/, message: '请输入字母 4-16 位账号', trigger: 'blur'},
            ],
            password: [
                { required: true, message: '请输入 8 到 16 位数字和大小写字母', trigger: 'blur' },
                { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
            ]
        },
    };
  },
  created(){
    this.api.glconfig.list()
    .then(res=>{
        if(res.data.code == 200){
            this.AdjustForm = res.data.data[0].isOpenEtcRegist;
        }
    })
    sessionStorage.setItem('routeName',this.$route.name);
    this.http = this.$store.state.http;                 // IP地址
    this.header = {
        Signature:'328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f',
        AuthorizationId:'wxda0b1c468b83bc99'
    }
    // 企业地区
    // 省
    this.api.areaInfo.selByParentNoAuth({parentCode:0})
    .then(res=>{
        if(res.data.code == 200){
            this.provinceOptions = res.data.data;
        }
    })
    // if(this.manageForm.cityCode != null || this.manageForm.cityCode != ""){
    //     // 市
    //     this.api.areaInfo.all({parentCode:this.manageForm.provinceCode})
    //     .then(res=>{
    //         if(res.data.code == 200){
    //             this.cityOptions = res.data.data;
    //         }
    //     })
    //     // 县 区
    //     this.api.areaInfo.all({parentCode:this.manageForm.cityCode})
    //     .then(res=>{
    //         if(res.data.code == 200){
    //             this.countyOptions = res.data.data;
    //         }
    //     })
    //     // 乡
    //     this.api.areaInfo.all({parentCode:this.manageForm.areaCode})
    //     .then(res=>{
    //         if(res.data.code == 200){
    //             this.villageOptions = res.data.data;
    //         }
    //     })
    // }
  },
  mounted() {},
  methods: {
    // 用户服务协议  
    isService() {
      let routeData = this.$router.resolve({
        name: "Policy"
      });
      window.open(routeData.href, "_blank");
    },
    /* 
    第一步 设置账号信息
    */
    submitFormLogin(formName) {
      this.$refs[formName].validate((valid) => {
        if(valid){
            if(this.loginForm.checked == true){
                this.api.authAll.setAccountNoAuth(this.loginForm)
                .then(res=>{
                    if(res.data.code == 200){
                        // this.hasAdmin = res.data.data.hasAdmin;
                        this.manageForm.userId = res.data.data.userId;
                        this.businessForm.id = res.data.data.userId;
                        this.active = 2;
                    }
                })
            }else{
            this.$message({
                type: 'warning',
                message: "请同意勾选平台服务协议和隐私政策"
            });
          }
        }
      });
    },
    /* 
    第二步 登记成功
    */
    // 选择企业管理员
    administrators(){
        this.stepShow = true;
        this.activeShow = true;
        this.activeShowtwo = false;
        this.active = 3;        // 申请结果...填写申请信息
    },
    // 申请企业管理员
    submitFormmanage(formName) {
      this.$refs[formName].validate((valid) => {
        if(valid){ 
            this.api.enterData.addNoAuth(this.manageForm)
            .then(res=>{
                if(res.data.code == 200){
                    // 申请结果 true 已有管理员， false 申请管理员
                    if(res.data.data.hasAdmin == true){
                        this.stepShow = false;
                        this.active = 5;        // 申请结果...已有管理员
                    }else{
                        // 审核结果页
                        this.stepShow = false;
                        this.active = 6;        // 审核结果页
                        var that = this;
                        var f = setInterval(function(){
                            that.Jump--
                            if(that.Jump == 0){
                                clearInterval(f);
                                sessionStorage.removeItem('routeName')
                                that.$router.push({name:'Login'})
                            }
                        },1000)
                    }
                }
            })
        }
      });
    },
    // 选择省
    enterprises(id){
        let resultArr = this.provinceOptions.filter((item)=>{
            return item.code == id;
        });
        this.manageForm.proName = resultArr[0].name;
        this.api.areaInfo.selByPetNoAuth({parentCode:id})
        .then(res=>{
            if(res.data.code == 200){
                this.manageForm.cityCode = '';                 // 市
                this.manageForm.areaCode = '';                 // 县 区
                this.manageForm.countyCode = '';               // 乡
                this.cityOptions = res.data.data;
            }
        })
    },
    // 选择市
    citys(id){
        let resultArr = this.cityOptions.filter((item)=>{
            return item.code == id;
        });
        this.manageForm.cityName = resultArr[0].name;
        this.api.areaInfo.selByPetNoAuth({parentCode:id})
        .then(res=>{
            if(res.data.code == 200){
                this.manageForm.areaCode = [];                 // 县 区
                this.manageForm.countyCode = '';               // 乡
                this.countyOptions = res.data.data;
            }
        })
    },
    // 选择县 区
    countys(id){
        let resultArr = this.countyOptions.filter((item)=>{
            return item.code == id;
        });
        this.manageForm.areaName = resultArr[0].name;
        this.api.areaInfo.selByPetNoAuth({parentCode:id})
        .then(res=>{
            if(res.data.code == 200){
                this.manageForm.countyCode = '';               // 乡
                this.villageOptions = res.data.data;
            }
        })
    },
    // 选择乡
    villages(id){
        let resultArr = this.villageOptions.filter((item)=>{
            return item.code == id;
        });
        this.manageForm.countryName = resultArr[0].name;
        this.manageForm.addressone = "ok";
    },
    // 选择加入企业
    enterprise(){
        this.stepShow = true;
        this.activeShow = false;
        this.activeShowtwo = true;
        this.active = 4;        // 加入企业
    },
    // 加入企业表单提交
    submitFormBusiness(formName){
      this.$refs[formName].validate((valid) => {
        if(valid){
            this.api.authAll.joinCpyNoAuth(this.businessForm)
            .then(res=>{
                if(res.data.code == 200){
                    // 审核结果页
                    this.stepShow = false;
                    this.active = 6;        // 审核结果页
                    var that = this;
                    var f = setInterval(function(){
                        that.Jump--
                        if(that.Jump == 0){
                            clearInterval(f);
                            sessionStorage.removeItem('routeName')
                            that.$router.push({name:'Login'})
                        }
                    },1000)
                }
            })
        }
      });
    },
    // 营业执照上传
    handleQrSuccess(res, file) {
        if(res.code == 200){
            this.$message({
                message:'上传营业执照成功',
                type:'success',
            })
            console.log(res.data)
            this.manageForm.filelist = res.data;
            this.manageForm.qrUrl = this.http + '/saas-common/upload/imageShow?fileName=' + res.data[0].fileName + '&Signature=328b14094b22d98589a0b135efe853eae3cde83b53e8f7920265df1c88c6389f&AuthorizationId=wxda0b1c468b83bc99';
        }else{
            this.$message({
                message:'上传营业执照失败,请重新上传!',
                type:'error'
            })
        }
    },
    // 营业执照图片格式认证
    beforelogoUpload(file) {
        var types = [
            "image/jpeg",
            "image/png", 
            "image/jpg"
        ];
        const isImg = types.indexOf(file.type) == -1;
        const isLt3M = file.size / 1024 / 1024 < 3; 
        if (isImg) {
            this.$message.error("上传营业执照图片只能是 JPG PNG JPEG 格式!");
        }
        if (!isLt3M) {
            this.$message.error("上传营业执照图片大小不能超过 3MB!");
        }
        return !isImg && isLt3M;
    },
    // 点击获取手机验证码
    phoneClick(){
      if(this.loginForm.phone != ''){
        this.isDisabled = true;
        var time = 60;
        this.api.authAll.getPhoneCodeNoAuth({phone:this.loginForm.phone})
        .then(res=>{
            var that = this;
            if(res.data.code == 200){
                var f = setInterval(function(){
                var newTime = time - 1;
                that.codephone = '重新获取('+ newTime +')s'
                that.color="onokPhone"
                    time--
                    if(time == 0){
                    that.codephone="获取验证码"
                    that.color="onPhone"
                    clearInterval(f);
                    that.isDisabled = false;
                    }
                },1000)  
            }
        })
      }else{
        this.$message({
          type: 'warning',
          message: "请先输入手机号再获取验证码"
        });
      }
    },
    // 返回登录页面
    Gologin(){
        sessionStorage.removeItem('routeName')
        this.$router.push({name:'Login'})
    }
  },
};
</script>
<style lang="less" scoped>
// 上传图片样式开始
.img-uploader{
    width: 80px;
    height: 106px;
}
.img-uploader .el-upload {
    width: 80px;
    height: 106px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
}
.img-uploader .el-upload:hover {
    border-color: #409eff;
}
[class*=" avatar-uploader-icon"], [class^=avatar-uploader-icon]{
    border: 1px dashed #d9d9d9;
    border-radius: 4px;
    font-size: 28px;
    color: #8c939d;
    width: 80px;
    height: 106px;
    line-height: 106px;
    text-align: center;
}
.avatar {
    width: 80px;
    height: 106px;
    display: block;
}
// 上传图片样式结束
// 头部开始
.login-concent{
    width: 100%;
    height: 70px;
    padding:8px 5%;
    box-sizing: border-box;
    border-bottom: 1px solid #e4e4e4;
    line-height: 54px;
    overflow: hidden;
    position: fixed;
    top: 0;
    left: 0;
    background-color: #fff;
    z-index: 10;
    .img-lg{
        float: left;
        img{
            width:48px;
            height: 48px;
        }
    }
    .img-size{
        float: left;
        margin-left:24px;
        box-sizing: border-box;
        span{
            font-size: 24px;
            line-height: 54px;
            font-weight: 600;
            color: #333;
        }
    }
}
// 头部结束
.stepAll{
    width: 100%;
    height: 275px;
    padding-top:155px;
    box-sizing: border-box;
}
.stepBox{
    width: 100%;
    height:calc(100vh - 370px);
    box-sizing: border-box;
    overflow: hidden;
    overflow-y: scroll;
    margin: auto;
    .stepAllBox{
        max-width:562px;
        overflow: hidden;
        padding-bottom:80px;
        margin: auto;
    }
}
.concentBox {
    max-width:562px;
    height: 90%;
    padding: 12px 35px 35px;
    box-sizing: border-box;
    margin: auto;
    text-align: center;
    user-select:none;
    .secussIcon{
        width: 100%;
        i{
            font-size: 32px;
            line-height: 32px;
            color: #17D08D;
        }
    }
    h3{
        font-size: 30px;
        line-height: 85px;
        color: #333;
        font-weight: 600;
    }
    p{
        font-size: 18px;
        line-height:30px;
        color: #828282;
    }
    .secussButton{
        width: 100%;
        margin:64px auto 0;
        .elButton{
            width:100%;
            height: 46px;
            border: 1px solid #e4e4e4;
            border-radius: 4px;
            background-color: #f0f0f0;
            text-align: center;
            color: #333;
            font-size: 16px;
            line-height: 46px;
            cursor: pointer;
            
        }
        .elButton:hover{
            background-color:#D74653;
            color: #fff;
            border: 0;
        }
        .elButton:active{
            background-color:#CF1724;
        }
    }
}
.menu-list{
  width: 100%;
  text-align: center;
}
.totalTab {
  font-size: 24px;
  line-height: 88px;
  font-weight: 650;
  color: #333;
}
.strip {
  width: 24px;
  height: 3px;
  background-color: #ea222e;
  border-radius: 60px;
  margin: 5px 0 12px 55px;
}
.striptwo {
  width: 24px;
  height: 3px;
  background-color: #ea222e;
  border-radius: 60px;
  margin: 5px 0 12px 192px;
}
.text {
  width: 100%;
}
.identifybox{
  width: 113px;
  height: 40px;
  float: left;
  margin-left:8px;
  margin-bottom:20px;
  img{
    width: 100%;
    height: 100%;
  }
}
.operation{
  margin-top: 70px;
   div{
     width:135px;
     text-align: center;
     float: left;
   }
    span{
      font-size: 14px;
      font-weight: 600;
      color: #333;
      
    };
}
.onPhone{
  font-size: 12px;
  color: #EA222E;
}
.onokPhone{
  font-size: 12px;
  color: #BFBFBF;
}
// 版权开始
.copyright-bottom{
  width: 100%;
  height:60px;
  padding-top:10px;
  box-sizing: border-box;
  background-color:#F6F6F6;
  position: fixed;
  bottom: 0;
  left: 0;
  .copyright{
    width:auto;
    height: 27px;
    margin: auto;
    text-align: center;
    span{
      font-size: 12px;
      line-height: 24px;
      color: #828282;
    }
    p{
      font-size: 12px;
      color: #828282;
    }
    a{
        font-size: 12px;
        color: #828282;
    }
  }
}
// 版权结束
</style>