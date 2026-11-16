<template>
<!-- 忘记密码页面 -->
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
            <el-steps :active="active" align-center process-status="wait" finish-status="finish" style="max-width:800px;margin:auto;">
                <el-step title="身份验证" icon="el-icon-ali1"></el-step>
                <el-step title="修改密码" icon="el-icon-ali2"></el-step>
                <el-step title="修改成功" icon="el-icon-ali3"></el-step>
            </el-steps>
        </div>
        <div class="stepBox">
            <div class="stepAllBox">
                <!-- 验证手机号页面 -->
                <el-form v-if="active == 1" :model="loginForm" :rules="rule" ref="loginForm" label-width="auto" :hide-required-asterisk="false" class="concentBox">
                    <el-form-item prop="phone" label="手机号">
                    <el-input v-model="loginForm.phone" placeholder="请输入11位手机号码" autocomplete="off"></el-input>
                    </el-form-item>
                    <!-- 验证码开始 -->
                    <el-form-item prop="code" label="验证码">
                        <el-input v-model="loginForm.code" placeholder="请输入验证码" autocomplete="off">
                            <template slot="suffix">
                            <span v-if="!isDisabled" :class="color" @click="phoneClick" v-html="codephone" style="margin-right:16px;cursor: pointer;"></span>
                            <span v-if="isDisabled" :class="color" v-html="codephone" style="margin-right:16px;cursor: pointer;"></span>
                            </template>
                        </el-input>
                    </el-form-item>
                    <!-- 验证码结束 -->
                    <el-form-item style="margin-bottom:0;margin-top:40px;">
                    <el-button type="primary" @click="submitFormLogin('loginForm')" style="width:100%;">下 一 步</el-button>
                    </el-form-item>
                </el-form>
                <!-- 修改密码页面 -->
                <el-form v-if="active == 2" :model="modifyForm" :rules="rules" ref="modifyForm" label-width="100px" :hide-required-asterisk="false" class="concentBox">
                    <el-form-item prop="account" label="账号">
                        <el-select v-model="modifyForm.id" placeholder="请选择账号" style="width:100%;">
                            <el-option
                            v-for="item in accountAll"
                            :key="item.value"
                            :label="item.account"
                            :value="item.id">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item prop="newPassword" label="输入新密码">
                        <el-input v-model="modifyForm.newPassword" placeholder="请输入长度为8-16位 包括数字和大小写字母" type="password" show-password autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item prop="checkPass" label="确认新密码">
                        <el-input v-model="modifyForm.checkPass" placeholder="请确认新的密码" type="password" show-password autocomplete="off"></el-input>
                    </el-form-item>
                    <el-form-item style="margin-bottom:0;margin-top:40px;">
                        <el-button type="primary" @click="submitFormModify('modifyForm')" style="width:100%;float:left;">下 一 步</el-button>
                    </el-form-item>
                </el-form>
                <!-- 修改成功，立即登录页面 -->
                <div v-if="active ==3" class="concentBox">
                    <div class="secussIcon"><i class="el-icon-alidui"></i></div>
                    <h3>密码修改成功</h3>
                    <p>请牢记您的新密码，<span style="color:#CF1724;">{{Jump}}s </span> 后跳转到登录页</p>
                    <div class="secussButton">
                        <div style="margin:auto;" class="elButton" @click="enterprise">立 即 登 录</div>
                    </div>
                </div>
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
        Jump:5,         // 秒数
        active:1,       // 步骤数
        // 第一步 身份验证
        loginForm: {
            phone:'',
            code:'', // 表单提交的验证码
        },
        userId:'',  // 验证通过的id
        isDisabled: false, //控制按钮是否可以点击（false:可以点击，true:不可点击）
        color:'onPhone',
        codephone:'获取验证码',
        // 身份验证表单提交
        rule: {
            phone: [
                { required: true, message: '请输入手机号', trigger: 'blur' },
                { pattern: /^(13[0-9]|14[01456879]|15[0-3,5-9]|16[2567]|17[0-8]|18[0-9]|19[0-3,5-9])\d{8}$/, message: '请输入正确的号码格式', trigger: 'blur'}
            ],
            code: [
                { required: true, message: '请输入验证码', trigger: 'blur' },
                { max: 6, message: '请输入6位验证码', trigger: 'blur' }
            ]
        },
        accountAll:[],      // 账号列表
        // 第二步 修改密码
        modifyForm: {
            id: '',
            newPassword: '',
            checkPass:''
        },
        // 修改密码表单提交
        rules: {
            id: [
                { required: true, message: '请输入账号', trigger: 'blur' }
            ],
            newPassword: [
                { required: true, message: '请输入 8 到 16 位数字和大小写字母', trigger: 'blur' },
                { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
            ],
            checkPass: [
                { required: true, message: '请输入 8 到 16 位数字和大小写字母', trigger: 'blur' },
                { pattern: /(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{8,16}$/, message: '请输入 8 到 16 位至少包含数字跟字母', trigger: 'blur'},
            ]
        },
    };
  },
  created(){},
  mounted() {},
  methods: {
    /* 
    第一步 验证手机号
    */
    submitFormLogin(formName) {
      this.$refs[formName].validate((valid) => {
        if(valid){
            this.api.authAll.valIdentity(this.loginForm)
            .then(res=>{
                if(res.data.code == 200){
                    this.accountAll = res.data.data;
                    this.active = 2;
                }
                if(res.data.code == 999){
                    Message({
                        message:'身份验证失败!',
                        type:'error',
                    })
                }
            })
        }
      });
    },
    /* 
    第二步 修改密码
    */
    submitFormModify(formName) {
      this.$refs[formName].validate((valid) => {
        if(valid){
          if(this.modifyForm.checkPass != this.modifyForm.newPassword){
            this.$message({
              message: '两次新密码不一致，请重新设置',
              type: 'warning'
            });
          }else{
            this.api.authAll.updPwd(this.modifyForm)
            .then(res=>{
                if(res.data.code == 200){
                    this.active = 3;
                    var that = this;
                    var f = setInterval(function(){
                        that.Jump--
                        if(that.Jump == 0){
                            clearInterval(f);
                            that.$router.push({name:'Login'})
                        }
                    },1000)
                }

            })
          }
        }
      });
    },
    /* 
    第三步 修改成功
    */
    // 返回立即登录页
    enterprise(){
        this.Jump = 0;
        this.$router.push({name:'Login'})
    },
    // 点击获取手机验证码
    phoneClick(){
      if(this.loginForm.phone != ''){
        this.isDisabled = true;
        var time = 60;
        this.api.authAll.valCode({phone:this.loginForm.phone})
        .then(res=>{
            var that = this;
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
        this.$router.push({name:'Login'})
    }
  },
};
</script>
<style lang="less" scoped>
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
        span{
            font-size: 24px;
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
        span{
            font-size: 18px;
            line-height:30px; 
        }
    }
    .secussButton{
        width: 80%;
        margin:30px auto 0;
        .elButton{
            width: 85%;
            height: 46px;
            border-radius: 4px;
            background-color: #CF1724;
            text-align: center;
            color: #fff;
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