<template>
<div class="logins">
  <div class="login-concent">
      <div class="img-lg">
        <img src="../../src/assets/lg.svg" alt="">
      </div>
      <div class="img-size"><span>云管理</span></div>
  </div>
  <div class="login-core">
    <!-- <img class="img" src="../../src/assets/banner.jpg" alt=""> -->
    <div class="concentBox">
      <div class="menu-list">云管理</div>
      <div class="concentBoxloginForm">
        <el-form :model="loginForm" :rules="rule" ref="loginForm" :hide-required-asterisk="false">
          <el-form-item prop="username">
            <el-input v-model="loginForm.username" placeholder="请输入账号" prefix-icon="el-icon-user"></el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input type="password" v-model="loginForm.password" placeholder="请输入密码" show-password prefix-icon="el-icon-unlock"></el-input>
          </el-form-item>
          <el-form-item prop="code">
            <el-input v-model="loginForm.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入验证码">
              <template slot="suffix">
                <span v-if="!isDisabled" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:3px;cursor:pointer;"></span>
                <span v-if="isDisabled" :class="color" v-html="codephone" style="cursor: pointer;" ></span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item style="margin-bottom:4%;">
            <el-button type="primary" @click.native.prevent="submitFormLogin('loginForm')" @keyup.enter.native="submitFormLogin('loginForm')" style="width:100%;hight:48px;">登录</el-button>
          </el-form-item>
          <el-checkbox v-model="loginForm.checked">
              <span style="font-size:12px;color:#333;font-weight: 400;cursor:pointer;">登录即代表您同意   
                <span style="line-height: 18px;color:#CF1724; font-weight:500;cursor: pointer;" @click="isService()">服务协议</span> 和
                <span style="line-height: 18px;color:#CF1724; font-weight:500;cursor: pointer;" @click="isService()">隐私政策</span>
              </span>
          </el-checkbox>
          <div class="operation">
            <!-- 注册申请 Register -->
            <div><router-link to="Register" style="color:#333;">注册申请</router-link></div>
            <div><router-link to="Forget" style="color:#333;">忘记密码</router-link></div>
          </div>
        </el-form>
      </div>
    </div>
  </div> 
  <div class="copyright-bottom">
    <div class="copyright"><span>Copyright © 山西海云启辰科技有限公司 版权所有</span><p><a target="_blank" href="https://beian.miit.gov.cn/">晋ICP备2022001434号</a></p></div>
  </div>
</div>

      <!-- <div style="clear: both"></div> -->
      <!-- <div class="text"> -->
        <!-- <div class="textAll" v-show="sel === 1"> -->
          <!-- <div class="strip"></div>
          <div>
              <el-form :model="phoneFormLogin" :rules="rules"  ref="phoneFormLogin" label-width="0px" :hide-required-asterisk="false">
                <el-form-item prop="phone">
                  <div style="margin-top: 15px;">
                    <el-input placeholder="请输入手机号码" prefix-icon="el-icon-user" v-model="phoneFormLogin.phone">
                    </el-input>
                  </div>
                </el-form-item>
                <el-form-item  prop="code">
                  <el-input v-model="phoneFormLogin.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入短信验证码">
                    <template slot="suffix"><span :class="color" @click="phoneClick" v-html="codephone" style="margin-right:16px;cursor: pointer;"></span></template>
                  </el-input>
                </el-form-item>
                <el-form-item style="margin-bottom:0;"> -->
                  <!-- <el-button type="primary" @keyup.enter.native="submitPhoneFormLogin('phoneFormLogin')" @click="submitPhoneFormLogin('phoneFormLogin')" style="width:100%;margin-top:18px;">登录</el-button> -->
                  <!-- <el-button type="primary" @click.native.prevent="submitPhoneFormLogin('phoneFormLogin')" @keyup.enter.native="submitPhoneFormLogin('phoneFormLogin')" style="width:100%;margin-top:18px;">登录</el-button>
                </el-form-item>
                <el-form-item prop="checkedPhone" style="margin-bottom:0;">
                    <el-checkbox-group v-model="phoneFormLogin.checkedPhone" style="height:35px;margin-bottom:-28px;">
                      <el-checkbox name="checkedPhone"><span style=" font-size:12px;line-height:80px;color: #333;font-weight: 400;">登录即代表您同意    -->
                        <!-- to="Service"<router-link to="Login" style="text-decoration:none;"><span style="line-height: 18px;color:#CF1724; font-weight: 500;">平台服务协议</span></router-link> 和
                        to="Policy"<router-link to="Login" style="text-decoration:none;"><span style="line-height: 18px;color:#CF1724; font-weight: 500;">隐私政策</span></router-link></span> -->
                      <!-- </el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
              </el-form>
          </div> -->
        <!-- </div> -->
        <!-- <div class="textAll" v-show="sel === 2"> -->
          <!-- <div class="striptwo"></div> -->
          <!-- <div>
              <el-form :model="loginForm" :rules="rule" ref="loginForm" label-width="0px" :hide-required-asterisk="false">
                <el-form-item prop="username">
                  <el-input v-model="loginForm.username" placeholder="请输入账号" prefix-icon="el-icon-user
                  "  style="width:270px;hight:48px;marign:auto;"></el-input>
                </el-form-item>
                <el-form-item prop="password">
                  <el-input type="password" v-model="loginForm.password" placeholder="请输入密码" show-password prefix-icon="el-icon-unlock"  style="width:270px;hight:48px;"></el-input>
                </el-form-item>
                <el-form-item prop="code">
                  <el-input v-model="loginForm.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入短信验证码">
                    <template slot="suffix">
                      <span :class="color" @click="phoneClick" v-html="codephone" style="margin-right:16px;cursor: pointer;" v-if="!isDisabled"></span>
                      <span :class="color" v-html="codephone" style="margin-right:16px;cursor: pointer;" v-if="isDisabled"></span>
                    </template>
                  </el-input>
                </el-form-item>
                <el-form-item style="margin-bottom:0;">
                  <el-button type="primary" @click.native.prevent="submitFormLogin('loginForm')" @keyup.enter.native="submitFormLogin('loginForm')" style="width:270px;hight:48px;">登录</el-button>
                </el-form-item>
                <el-form-item prop="checked" style="margin-bottom:0;">
                  <el-checkbox-group v-model="loginForm.checked" style="margin-bottom:-28px;height:35px;">
                    <el-checkbox name="checked"><span style="font-size:12px;line-height:80px;color:#333;font-weight: 400;cursor: pointer;">登录即代表您同意   
                      <router-link to="Login" style="text-decoration:none;"><span style="line-height: 18px;color:#CF1724; font-weight: 500;cursor: pointer;">平台服务协议</span></router-link> 和
                      <router-link to="Login" style="text-decoration:none;"><span style="line-height: 18px;color:#CF1724; font-weight: 500;cursor: pointer;">隐私政策</span></router-link></span>
                    </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
                <div class="operation">
                  <router-link to="Register" style="text-decoration:none;"><div><span>注册申请</span></div></router-link>
                  <router-link to="Forget" style="text-decoration:none;"><div><span>忘记密码</span></div></router-link>
                </div>
              </el-form>
          </div> -->
      <!-- </div> -->
</template>
<script>
export default {
  data() {
    return {
      // sel: 1,
      // tabs: [
      //   { label: "短信登录", id: 1 },
      //   { label: "密码登录", id: 2 },
      // ],
      //用户名登录
      loginForm: {
        username: '',
        password:'',
        code:'',        // 表单提交的验证码
        checked:true,  // 勾选同意协议
      },
      isDisabled: false, //控制按钮是否可以点击（false:可以点击，true:不可点击）
      color:'onPhone',
      codephone:'获取验证码',
      // identifyCodes:'',  // 图形验证码
      //密码登录表单验证
      rule: {
          username: [
            { required: true, message: '请输入账号', trigger: 'blur' },
            // { pattern:/^[a-zA-Z]{4,16}$/, message: '请输入 4 到 16 位账号', trigger: 'blur'},
          ],
          password: [
            { required: true, message: '请输入 8 到 16 位包括数字和大小写字母', trigger: 'blur' },
            // { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
          ],
          code: [
            { required: true, message: '请输入验证码', trigger: 'blur' },
            { max: 6, message: '请输入6位验证码', trigger: 'blur' }
          ]
      },
      //短信登录
      // phoneFormLogin: {
      //   phone:'',
      //   code:'',
      //   checkedPhone:[]
      // },
      // //验证码登录表单验证
      // rules: {
      //     phone: [
      //       { required: true, message: '请输入手机号', trigger: 'blur' },
      //       { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'blur'}
      //     ],
      //     code: [
      //       { required: true, message: '请输入验证码', trigger: 'blur' },
      //       { pattern:/^\d{6}$/, message: '请输入正确的验证码格式', trigger: 'blur'}
      //     ],
      //     checkedPhone: [
      //       {type: 'array',required: true, message: '请同意勾选平台服务协议和隐私政策', trigger: 'blur'}
      //     ]
      // },
    };
  },
  created(){
    sessionStorage.clear()
    var _self = this;
    document.onkeydown = function(e){
      if(window.event == undefined){
        var key = e.keyCode;
      }else{
        var key = window.event.keyCode;
      }
      if(key == 13){
        _self.submitFormLogin('loginForm');
        // _self.submitPhoneFormLogin('phoneFormLogin')
      }
    }
  },
  mounted() {},
  methods: {
    // select(item) {
    //   this.sel = item.id;
    //   if(this.sel == 2){
    //     this.getImgCode(); // 调用图形码
    //   }
    // },
    // 图形验证码
    // getImgCode(){
    //   this.api.imglogin.imglogin()
    //   .then(res=>{
    //     if (res && res.data && res.data.size) {
		//       const dataInfo = res.data
		//       let reader = new window.FileReader()
		//       // 使用readAsArrayBuffer读取文件, result属性中将包含一个 ArrayBuffer对象以表示所读取文件的数据
		//       reader.readAsArrayBuffer(dataInfo)
    //       var that = this;
		//       reader.onload = function (e) {
		//         const result = e.target.result
		//         const contentType = result.type
		//         // 生成blob图片,需要参数(字节数组, 文件类型)
		//         const blob = new Blob([result], { type: contentType })
		//         // 使用 Blob 创建一个指向类型化数组的URL, URL.createObjectURL是new Blob文件的方法,可以生成一个普通的url,可以直接使用,比如用在img.src上
		//         const url = window.URL.createObjectURL(blob);
		//         // console.log(url) // 产生一个类似 blob:d3958f5c-0777-0845-9dcf-2cb28783acaf 这样的URL字符串
    //         that.identifyCodes = url;
		//       }
		//     }
    //   })
    // },
    // 用户服务协议  
    isService() {
      let routeData = this.$router.resolve({
        name: "Policy"
      });
      window.open(routeData.href, "_blank");
    },
    //用户名登录提交
    submitFormLogin(formName) {
      this.$refs[formName].validate((valid) => {
        if(valid){
          if(this.loginForm.checked == true){
            let param = new URLSearchParams();
            param.append('username', this.loginForm.username)
            param.append('password', this.loginForm.password)
            param.append('code', this.loginForm.code)
            param.append('checked', this.loginForm.checked)
            this.api.login.login(param)
            .then(res=>{
              if(res.data.code == 200){
                document.onkeydown = undefined;
                if(res.data.data.level == 1){
                  // 企业路由刷新
                  this.$store.dispatch('asynSetRoutes');
                  sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token);
                  sessionStorage.setItem('companyName',res.data.data.companyName); 
                  sessionStorage.setItem('username',res.data.data.username);
                  sessionStorage.setItem('level',res.data.data.level);
                  sessionStorage.setItem('path','');
                  sessionStorage.setItem('pathID','');
                  sessionStorage.setItem('companyId',res.data.data.companyId);
                  sessionStorage.setItem('userId',res.data.data.userId);
                  // if(this.$store.state.newstate){
                    // if(this.$store.state.newstate){
                      this.$message({
                        type: 'success',
                        message: "登录成功!",
                        duration:1000,
                        onClose:()=>{
                          this.$router.push({name:"Main"});
                        }
                      })
                    // }else{
                    //   this.$message({
                    //     type: 'warning',
                    //     message: "管理员未分配权限,请尽快联系管理员!",
                    //     duration:1000,
                    //     onClose:()=>{
                    //       sessionStorage.clear();
                    //     }
                    //   })
                    // }
                  // }
                }
                if(res.data.data.level == 2){
                  // 平台路由刷新
                  this.$store.dispatch('asynSetplatformRoutes')
                  sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token)
                  sessionStorage.setItem('companyName',res.data.data.companyName)
                  sessionStorage.setItem('username',res.data.data.username)
                  sessionStorage.setItem('xzqhCode',res.data.data.xzqhCode)
                  sessionStorage.setItem('level',res.data.data.level)
                  sessionStorage.setItem('accountTypes',2);
                  // if(this.$store.state.newstate != null){
                    // if(this.$store.state.newstate){
                      this.$message({
                        type: 'success',
                        message: "登录成功!",
                        duration:1000,
                        onClose:()=>{
                          this.$router.push({name:"Main"});
                        }
                      })
                    // }else{
                    //   this.$message({
                    //     type: 'warning',
                    //     message: "管理员未分配权限,请尽快联系管理员!",
                    //     duration:1000,
                    //     onClose:()=>{
                    //       sessionStorage.clear();
                    //     }
                    //   })
                    // }
                  // }
                }
                if(res.data.data.level == 3 || res.data.data.level == 4 || res.data.data.level == 5 || res.data.data.level == 6){
                  // 政府路由刷新
                  this.$store.dispatch('asynSetgoverRoutes')
                  sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token)
                  sessionStorage.setItem('companyName',res.data.data.companyName)
                  sessionStorage.setItem('username',res.data.data.username)
                  sessionStorage.setItem('xzqhCode',res.data.data.xzqhCode)
                  sessionStorage.setItem('level',res.data.data.level)
                  this.$message({
                    type: 'success',
                    message: "登录成功",
                    duration:500,
                    onClose:()=>{
                      this.$router.push({name:"Xzmain"})
                    }
                  });
                }
              }else{
                this.$message.error(res.data.msg);
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
    //短信登录提交
    // submitPhoneFormLogin(formName) {
    //   this.$refs[formName].validate((valid) => {
    //     if (valid) {
    //       let param = new URLSearchParams()
    //       param.append('phone', this.phoneFormLogin.phone)
    //       param.append('code', this.phoneFormLogin.code)
    //       param.append('checkedPhone', this.phoneFormLogin.checkedPhone)
    //       this.api.mobilelogin.mobilelogin(param)
    //       .then(res=>{
    //         if(res.data.code == 200){
    //           document.onkeydown = undefined;
    //             sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token)
    //             sessionStorage.setItem('companyName',res.data.data.companyName)
    //             sessionStorage.setItem('username',res.data.data.username)
    //             this.$message({
    //                 type: 'success',
    //                 message: "登录成功",
    //                 duration:500,
    //                 onClose:()=>{
    //                   this.$router.push({name:"Main"})
    //                 }
    //             });
    //         }else if(res.data.code == 2001){
    //           this.$message.error(res.data.msg)
    //         }else if(res.data.code == 2003){
    //           this.$message.error(res.data.msg)
    //         }else if(res.data.code == 3001){
    //           this.$message.error(res.data.msg)
    //         }else if(res.data.code == 3002){
    //           this.$message.error(res.data.msg)
    //         }else if(res.data.code == 3003){
    //           this.$message.error(res.data.msg)
    //         }else if(res.data.code == 3004){
    //           this.$message.error(res.data.msg)
    //         }
    //       })
    //     }
    //   });
    // },
    //点击获取手机验证码
    phoneClick(){
      if(this.loginForm.username != ''){
        this.isDisabled = true;
        var time = 120;
        this.api.code.code({account:this.loginForm.username})
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
          if(res.data.code == 999){
            that.isDisabled = false;
          }
          if(res.data.code != 200){ // 账号不存在(2002)或请求异常后必须恢复按钮,否则只能刷新页面
            that.isDisabled = false;
          }
        })
      }else{
        this.$message({
          type: 'warning',
          message: "请先输入账号再获取验证码"
        });
        this.isDisabled = false;
      }
    }
  },
};
</script>
<style lang="less" scoped>
.logins{
  width: 100%;
  height: 100vh;
  // 头部开始
  .login-concent{
    // background-color: blue;
    height: 70px;
    padding:8px 5%;
    box-sizing: border-box;
    line-height: 48px;
    overflow: hidden;
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
  .login-core{
    background-color: burlywood;
    width: 100%;
    height: calc(100vh - 130px);
    background: url(../../src/assets/banner.jpg);
    background-size:cover;
    background-repeat: no-repeat;
    position: relative;
    // .img{
    //   width: 100%;
    //   height: calc(100vh - 130px);
    // }
    .concentBox{
      min-width: 335px;
      max-height: 85%;
      overflow-y: scroll;
      background-color: #fff;
      // opacity: 0.9;
      border-radius: 4px;
      box-sizing: border-box;
      position: absolute;
      right:4%;
      top:12%;
      z-index:5;
      .menu-list{
          padding-top:3%;
          width: 100%;
          text-align: center;
          font-size: 28px;
          line-height: 100px;
          font-weight: 600;
          color: #333;
      }
      .concentBoxloginForm{
        overflow: hidden;
        padding:0 10% 5%;
      }
      .operation{
        div{
            width:50%;
            text-align: center;
            float: left;
            font-size: 14px;
            font-weight: 600;
            line-height: 68px;
            color: #333;
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
    }
  }
  // 版权开始
  .copyright-bottom{
    width: 100%;
    height:60px;
    padding-top:10px;
    box-sizing: border-box;
    background-color: #fff;
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
}

// 版权结束
// .concentBox {
//   max-width: 340px;
//   max-height: 450px; 
//   background-color: rgba(255, 255, 255,.9);    
//   border-radius: 4px;
//   padding: 12px 35px 35px;
//   box-sizing: border-box;
//   position: absolute;
//   right:5%;
//   top:6%;
// }
// .strip {
//   width: 24px;
//   height: 3px;
//   background-color: #ea222e;
//   border-radius: 60px;
//   margin: 5px 0 12px 55px;
// }
// .striptwo {
//   width: 24px;
//   height: 3px;
//   background-color: #ea222e;
//   border-radius: 60px;
//   margin: 5px 0 12px 192px;
// }
// .text {
//   width: 100%;
// }

</style>