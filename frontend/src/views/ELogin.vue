<template>
    <div>
        <div class="casket">
            <div class="casketlog">
                <!-- 头部 -->
                <div class="login-concent">
                    <div class="img-lg">
                        <img src="../../src/assets/lg.svg" alt="">
                    </div>
                    <div class="img-size"><span>云管理</span></div>
                </div>
                <div class="conent">
                    <img src="../../src/assets/lg.svg" alt="">
                    <h3> 电子签章签署</h3>
                    <!-- <div class="totalTab"></div> -->
                    <div class="adBox">
                        <el-form :model="loginForm" :rules="rule" ref="loginForm" label-width="80px" label-position="left">
                            <el-form-item prop="phone" label="手机号">
                                <el-input v-model="loginForm.phone" placeholder="请输入手机号" prefix-icon="el-icon-user"></el-input>
                            </el-form-item>
                            <el-form-item prop="code" label="验证码">
                                <el-input v-model="loginForm.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入验证码">
                                <template slot="suffix">
                                    <span v-if="isDisabled == false" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:3px;cursor:pointer;"></span>
                                    <span v-if="isDisabled" :class="color" v-html="codephone" style="cursor: pointer;" ></span>
                                </template>
                                </el-input>
                            </el-form-item>
                            <div style="margin-top:10%;">
                                <el-button type="primary" @click.native.prevent="submitFormLogin('loginForm')" @keyup.enter.native="submitFormLogin('loginForm')" style="width:100%;hight:48px;">登 录</el-button>
                            </div>
                        </el-form>
                    </div>
                    <!-- tab选项卡 -->
                    <!-- <div class="totalTab"
                        :class="{ active: item.id === sel }"
                        v-for="item in tabs"
                        :key="item.id"
                        @click="select(item)" style="cursor: pointer;">
                        <span>{{ item.label }}</span>
                        <div class="stripnone" :style="item.width"
                        :class="{ strip: item.id === sel }">
                        </div>
                    </div>
                    <div class="adBox" v-if="sel === 1"> 
                        <el-form :model="AddData" :rules="rules" ref="AddData" label-width="80px" label-position="left">
                            <div style="clear:both;"></div>
                            <el-form-item label=" 参与人">
                                <el-select v-model="AddData.regType" @change="prologistics($event)" style="width:100%;">
                                    <el-option
                                    v-for="item in SelectcontractName"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <div v-if="logistics">     
                                <el-form-item prop="name" label="企业">
                                    <el-input v-model="AddData.name" placeholder="请输入企业名称"></el-input>
                                </el-form-item>   
                                <el-form-item prop="thirdPartyUserId" label="税号">
                                    <el-input v-model="AddData.thirdPartyUserId" placeholder="请输入营业执照编号（税号）"></el-input>
                                </el-form-item> 
                            </div> 
                            <el-form-item prop="username" label="姓名">
                                <el-input v-model="AddData.username" placeholder="请输入姓名"></el-input>
                            </el-form-item>
                            <el-form-item prop="phone" label="手机号">
                                <el-input v-model="AddData.phone" placeholder="请输入手机号"></el-input>
                            </el-form-item>
                            <el-form-item prop="code" label="验证码">
                                <el-input v-model="AddData.code" placeholder="请输入验证码">
                                    <template slot="suffix">
                                        <span v-if="isDisabledone == false" :class="colors" @click="phoneClicks" v-html="codephones" style="margin-right:3px;cursor:pointer;"></span>
                                        <span v-if="isDisabledone" :class="colors" v-html="codephones" style="cursor: pointer;" ></span>
                                    </template>
                                </el-input>
                            </el-form-item>
                            <div style="margin-top:10%;">
                                <el-button @click="AddDatasubmitForm('AddData')" type="primary" style="width:100%;">提 交</el-button>
                            </div>
                        </el-form>
                    </div> 
                    <div class="adBox" v-if="sel === 2"> 
                        <el-form :model="loginForm" :rules="rule" ref="loginForm" label-width="80px" label-position="left">
                            <el-form-item prop="phone" label="手机号">
                                <el-input v-model="loginForm.phone" placeholder="请输入手机号" prefix-icon="el-icon-user"></el-input>
                            </el-form-item>
                            <el-form-item prop="code" label="验证码">
                                <el-input v-model="loginForm.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入验证码">
                                <template slot="suffix">
                                    <span v-if="isDisabled == false" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:3px;cursor:pointer;"></span>
                                    <span v-if="isDisabled" :class="color" v-html="codephone" style="cursor: pointer;" ></span>
                                </template>
                                </el-input>
                            </el-form-item>
                            <div style="margin-top:10%;">
                                <el-button type="primary" @click.native.prevent="submitFormLogin('loginForm')" @keyup.enter.native="submitFormLogin('loginForm')" style="width:100%;hight:48px;">登 录</el-button>
                            </div>
                        </el-form>
                    </div>  -->
                </div>
                <!-- 底部 -->
                <div class="copyright-bottom">
                    <div class="copyright"><span>Copyright © 2021 云信通联能源科技有限公司 版权所有</span><p><a target="_blank" href="https://beian.miit.gov.cn/">晋ICP备2022001434号</a></p></div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>

export default {
    data(){
        return {
            sel:1,   // 选项卡默认显示tab
            tabs: [
                { label: "注 册", id: 1 ,width:"width:38px;"},
                { label: "登 录", id: 2 ,width:"width:38px;"},
            ],

            /*
            注册
            */
            logistics:true,                         // 企业默认显示
            // 参与人选择下拉框
            SelectcontractName:[
                { 
                    value: 1,
                    label: '个人'
                },
                { 
                    value: 2,
                    label: '企业'
                }
            ],  
            // 表单验证
            rules: { 
                name: [
                    { required: true, message: '请输入企业名称', trigger: 'blur' },
                ],
                thirdPartyUserId: [
                    { required: true, message: '请输入税号', trigger: 'blur' },
                    { pattern:/(^(?:(?![IOZSV])[\dA-Z]){2}\d{6}(?:(?![IOZSV])[\dA-Z]){10}$)|(^\d{15}$)/, message: '请输入正确的税号', trigger: 'blur'}
                ],
                username: [
                    { required: true, message: '请输入姓名', trigger: 'blur' },
                    { pattern: /^[\u4E00-\u9FA5]{2,4}$/, message: '姓名必须2-4位中文', trigger: 'change'}
                ],
                phone: [
                    { required: true, message: '请输入手机号', trigger: 'blur' },
                    { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'change'}
                ],
                code: [
                    { required: true, message: '请输入验证码', trigger: 'blur' },
                    { max: 6, message: '请输入6位验证码', trigger: 'change' }
                ]
            },
            // 新建表单数据
            AddData: {
                regType:2,
                name:'',
                thirdPartyUserId:'',
                username:'',
                phone:'',
                code:'',
            }, 
            colors:'onPhones',
            codephones:'获取验证码',
            isDisabledone: false, //控制按钮是否可以点击（false:可以点击，true:不可点击）
            /*
            登录
            */
            loginForm: {
                phone:'',
                code:'',        // 表单提交的验证码
            },
            color:'onPhone',
            codephone:'获取验证码',
            isDisabled: false, // 控制按钮是否可以点击（false:可以点击，true:不可点击）
            // identifyCodes:'',  // 图形验证码
            //密码登录表单验证
            rule: {
                phone: [
                    { required: true, message: '请输入手机号', trigger: 'blur' },
                    { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'blur'}
                ],
                code: [
                    { required: true, message: '请输入验证码', trigger: 'blur' },

                    
                    { max: 6, message: '请输入6位验证码', trigger: 'change' }
                ]
            },
        }
    },
    created(){
        var _self = this;
        document.onkeydown = function(e){
        if(window.event == undefined){
            var key = e.keyCode;
        }else{
            var key = window.event.keyCode;
        }
        if(key == 13){
            _self.submitFormLogin('loginForm');
        }
        }
    },
    methods:{
        // 选项卡
        select(item) {
            this.sel = item.id;
            if(this.sel == 1){
                this.AddData = {};
                this.AddData.regType = 2;
                this.logistics = true;          // 承运平台显示
            }else if(this.sel == 2){
                this.loginForm = {};
            }
        },
        // 参与人切换
        prologistics(id){
                // 企业
                if(id == 2){
                    this.AddData = {
                        regType:2,
                        name:'',
                        thirdPartyUserId:'',
                        username:'',
                        phone:'',
                        code:'',
                    }, 
                    this.logistics = true;          // 承运平台显示
                }
                // 个人
                if(id == 1){
                    this.AddData = {
                        regType:1,
                        name:'',
                        thirdPartyUserId:'',
                        username:'',
                        phone:'',
                        code:'',
                    }, 
                    this.logistics = false;         // 承运平台隐藏
                }
        },
        // 注册表单提交
        AddDatasubmitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    let obj;
                    if(this.AddData.regType == 2){
                        obj = {
                            regType:2,
                            username:this.AddData.username,
                            phone:this.AddData.phone,
                            code:this.AddData.code,
                            orgInfo:{
                                name:this.AddData.name,
                                thirdPartyUserId:this.AddData.thirdPartyUserId
                            },
                            userInfo:{
                                name:this.AddData.username,
                                thirdPartyUserId:this.AddData.phone
                            }
                        }
                    }else if(this.AddData.regType == 1){
                        obj = {
                            regType:1,
                            username:this.AddData.username,
                            phone:this.AddData.phone,
                            code:this.AddData.code,
                            userInfo:{
                                name:this.AddData.username,
                                thirdPartyUserId:this.AddData.phone
                            }
                        } 
                    }
                    this.api.authAll.regNoAuth(obj)
                    .then(res=>{
                        if(res.data.code == 200){
                            this.$message({
                                type: "success",
                                message: "提交注册成功,请登录！",
                                duration:1000,
                                onClose:()=>{
                                    this.sel = 2;
                                    this.loginForm.phone = this.AddData.phone;
                                }
                            });
                        }
                    })
                }
            });
        },

        // 用户名登录提交
        submitFormLogin(formName) {
        this.$refs[formName].validate((valid) => {
            if(valid){
                let param = new URLSearchParams();
                param.append('phone', this.loginForm.phone)
                param.append('code',this.loginForm.code)
                this.api.mobilelogin.mobilelogin(param)
                .then(res=>{
                    if(res.data.code == 200){
                        sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token)
                        sessionStorage.setItem('companyName',res.data.data.companyName)
                        sessionStorage.setItem('username',res.data.data.username)
                        this.$message({
                            type: "success",
                            message: "登录成功！",
                            duration:500,
                            onClose:()=>{
                                this.$router.push({name:"EsignIndex"})
                            }
                        });
                    }else{
                        this.isDisabledone = false;
                    }
                }) 
            }
        });
        },
        // 注册点击获取手机验证码
        phoneClicks(){
            if(this.AddData.phone != ''){
                this.isDisabledone = true;
                this.api.authAll.getPhoneCodeNoAuth({phone:this.AddData.phone})
                .then(res=>{
                    var that = this;
                    if(res.data.code == 200){
                        var time = 60;
                        var f = setInterval(function(){
                        let newTime = time - 1;
                        that.codephones = '重新获取('+ newTime +')s'
                        that.colors="onokPhones"
                            time--
                            if(time == 0){
                            that.codephones = "获取验证码"
                            that.colors="onPhones"
                            clearInterval(f);
                            that.isDisabledone = false;
                            }
                        },1000)  
                    }
                    if(res.data.code == 999){
                        that.isDisabledone = false;
                    }
                })
            }else{
                this.$message({
                type: 'warning',
                message: "请先输入手机号再获取验证码"
                });
                this.isDisabledone = false;
            }
        },
        // 点击登录获取手机验证码
        phoneClick(){
            if(this.loginForm.phone != ''){
                this.isDisabled = true;
                this.api.code.code({account:this.loginForm.phone})
                .then(res=>{
                var that = this;
                if(res.data.code == 200){
                    var time = 60;
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
                }else{
                    that.isDisabled = false;
                }
                if(res.data.code == 999){
                    that.isDisabled = false;
                }
                })
            }else{
                this.$message({
                type: 'warning',
                message: "请先输入手机号再获取验证码"
                });
                this.isDisabled = false;
            }
        }
    }
}
</script>
<style lang="less" scoped>
.casket{
    width: 100%;
}
.casketlog{
    width: 100%;
    // background: chocolate;
    margin: auto;
}
// 头部开始
.login-concent{
    width: 100%;
    height: 70px;
    padding:8px 5%;
    box-sizing: border-box;
    border-bottom: 1px solid #e4e4e4;
    line-height: 48px;
    overflow: hidden;
    // position: fixed;
    // top: 0;
    // left: 0;
    background-color: #fff;
    // z-index: 10;
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
        span{
            font-size: 24px;
            font-weight: 600;
            color: #333;
        }
    }
}
// 头部结束
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
// 中间内容开始
.conent{
    width: 40%;
    margin: auto;
    padding: 56px 25px;
    box-sizing: border-box;
    // background: olive;
    img{
        display: none;
    }
    h3{
        font-size: 26px;
        color: #333;
        text-align: center;
        font-weight: 600;
        line-height: 68px;
    }
}
.onPhone{
font-size: 12px;
color: #EA222E;
}
.onokPhone{
font-size: 12px;
color: #BFBFBF;
}
.onPhones{
font-size: 12px;
color: #EA222E;
}
.onokPhones{
font-size: 12px;
color: #BFBFBF;
}
// 中间内容结束
// 选项卡开始
.totalTab {
    width: 50%;
    float: left;
    font-size: 20px;
    line-height: 26px;
    font-weight: 650;
    // margin-top:8%;
    color: #333;
    text-align: center;
    padding: 5% 0;
    box-sizing: border-box;
}
.totalTab.active {
    color: #ea222e;
}
.stripnone{
    height: 3px;
    background-color: #ea222e;
    margin:4px auto 0;
    border-radius:25px;
    display: none;
}
.strip {
    display: block;
}
.adBox{
    height: 100%;
    margin-top:10%;
}
//选项卡结束
@media screen and (max-width:900px) and (min-width:750px) {
    .casketlog{
        width: 70%;
        // background: rgb(52, 25, 204);
        margin: auto;
    }
    .login-concent{
        display: none;
    }
    .copyright-bottom{
        display: none;
    }
    .conent{
        width: 100%;
        margin: auto;
        padding: 25px 25px 25px;
        box-sizing: border-box;
        img{
            width: 64px;
            height: 64px;
            display:block;
            margin: auto;
        }
        h3{
            display: block;
            font-size: 18px;
            color: #333;
            text-align: center;
            font-weight: 600;
            line-height: 68px;
        }
    }
}
@media screen and (max-width:750px) {
    .casketlog{
        width: 100%;
        margin: auto;
    }
    .login-concent{
        display: none;
    }
    .copyright-bottom{
        display: none;
    }
    .conent{
        width: 95%;
        margin: auto;
        padding: 25px 25px 25px;
        box-sizing: border-box;
        img{
            width: 64px;
            height: 64px;
            display:block;
            margin: auto;
        }
        h3{
            display: block;
            font-size: 18px;
            color: #333;
            text-align: center;
            font-weight: 600;
            line-height: 68px;
        }
    }
    .adBox{
        padding-top:10%;
        box-sizing: border-box;
    }
}
</style>
