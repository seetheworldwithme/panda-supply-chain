<template>
    <!-- 第一屏LED看板 -->
    <div class="FirstScreen">
        <div class="header">
            <p>煤炭大宗交易数字平台</p>
        </div>
        <div class="nav">
            <div class="section">
                <!-- 代码特效 -->
                <div class="codeEffects">
                    <div class="codeBox1"></div>
                    <div class="codeBox2"></div>
                    <div class="codeBox3"></div>
                    <div class="codeBox4"></div>
                    <div class="codeBox5"></div>
                </div>
                <!-- 粒子盒子 -->
                <div class="particleBox">
                    <vue-particles
                    color="#00FFFF"
                    :particleOpacity="0.3"
                    :particlesNumber="28"
                    shapeType="circle"
                    :particleSize="10"
                    linesColor="#fff"
                    :linesWidth="0"
                    :lineLinked="false"
                    :lineOpacity="0"
                    :linesDistance="500"
                    :moveSpeed="2"
                    :hoverEffect="false"
                    :clickEffect="false">
                    </vue-particles>
                </div>
            </div>
            <div class="concentBox">
                <div class="concentBoxloginForm">
                    <el-form :model="loginForm" :hide-required-asterisk="false">
                        <el-form-item prop="username">
                            <el-input type="text" v-model="loginForm.username" placeholder="请输入账号" prefix-icon="el-icon-user"></el-input>
                        </el-form-item>
                        <el-form-item prop="password">
                            <el-input type="password" v-model="loginForm.password" placeholder="请输入密码" show-password prefix-icon="el-icon-unlock"></el-input>
                        </el-form-item>
                        <el-form-item prop="code">
                            <el-input v-model="loginForm.code" prefix-icon="el-icon-mobile-phone" placeholder="请输入图形码">
                                <template slot="suffix">
                                    <div class="imgcolor" @click="phoneClick">
                                        <img :src="identifyCodes" alt="">
                                    </div>
                                </template>
                            </el-input>
                        </el-form-item>
                        <el-form-item style="margin-top:10%;">
                            <el-button type="primary" @click.native.prevent="submitFormLogin('loginForm')" @keyup.enter.native="submitFormLogin('loginForm')" style="width:100%;background-color:rgba(1, 57, 141, .6);border:0;color:#8ce9e9;">登 录</el-button>
                        </el-form-item>
                    </el-form>
                </div>
            </div>
        </div>
        <!-- <div class="footer">
            <div class="copyright"><span>Copyright © 2021 云信通联能源科技有限公司 版权所有</span><p> 赣ICP备2021002370号</p></div>
        </div> -->
        <!-- <ul class="FirstScreen-box">
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
            <li></li>
        </ul> -->
    </div>
</template>
<script>
export default{
    data(){
        return {
            // 用户名登录
            loginForm: {
                username: '',
                password:'',
                code:'',        // 表单提交的图形码
            },
            identifyCodes:'',  // 图形图形码
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
    mounted() {
        this.getImgCode();
    },
    methods:{
        //点击获取图形码
        phoneClick(){
            this.getImgCode();  // 调用图形码 
        },
        // 图形码
        getImgCode(){
            this.api.imglogin.imglogin()
            .then(res=>{
                if (res && res.data && res.data.size) {
                    const dataInfo = res.data
                    let reader = new window.FileReader()
                    // 使用readAsArrayBuffer读取文件, result属性中将包含一个 ArrayBuffer对象以表示所读取文件的数据
                    reader.readAsArrayBuffer(dataInfo)
                    var that = this;
                    reader.onload = function (e) {
                        const result = e.target.result
                        const contentType = result.type
                        // 生成blob图片,需要参数(字节数组, 文件类型)
                        const blob = new Blob([result], { type: contentType })
                        // 使用 Blob 创建一个指向类型化数组的URL, URL.createObjectURL是new Blob文件的方法,可以生成一个普通的url,可以直接使用,比如用在img.src上
                        const url = window.URL.createObjectURL(blob);
                        // console.log(url) // 产生一个类似 blob:d3958f5c-0777-0845-9dcf-2cb28783acaf 这样的URL字符串
                        that.identifyCodes = url;
                        console.log(url)
                    }
                }
            })
        },
        //用户名登录提交
        submitFormLogin() {
            if(this.loginForm.username == "" && this.loginForm.password == ""){
                this.$message({
                type: 'warning',
                message: "账号密码不能为空"
                });
            }else{
                let param = new URLSearchParams()
                param.append('username', this.loginForm.username)
                param.append('password', this.loginForm.password)
                param.append('code', this.loginForm.code)
                this.api.login.screenLogin(param)
                .then(res=>{
                    if(res.data.code == 200){
                        document.onkeydown = undefined;
                        // 登陆成功
                        sessionStorage.setItem('Authorization',"Bearer "+ res.data.data.token)
                        sessionStorage.setItem('xzqhCode',res.data.data.xzqhCode)
                        sessionStorage.setItem('companyId',res.data.data.companyId)
                        sessionStorage.setItem('accountTypes',4)
                        if(res.data.data.xzqhCode == '3201'){
                            this.$message({
                                type: 'success',
                                message: "登录成功",
                                duration:500,
                                onClose:()=>{
                                    this.$router.push({name:'NjMain'});
                                }
                            });
                        }else{
                            this.$message({
                                type: 'success',
                                message: "登录成功",
                                duration:500,
                                onClose:()=>{
                                    this.$router.push({name:'Dpmain'});
                                }
                            });
                        }

                    }else if(res.data.code == 2001){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }else if(res.data.code == 2003){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }else if(res.data.code == 3001){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }else if(res.data.code == 3002){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }else if(res.data.code == 3003){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }else if(res.data.code == 3004){
                        this.$message.error(res.data.msg)
                        this.getImgCode();
                    }
                })   
            }
        },
    },
}
</script>
<style scoped lang="less">
// 改变input里的字体颜色
/deep/input::-webkit-input-placeholder {
  color: #17a1e5;
  font-size: 15px;
}
// 改变input框背景颜色
/deep/.el-input__inner {
  background-color: transparent !important;
  border: 1px solid #1296db;
}
/deep/.el-input__inner:focus {
    border:1px solid #4D7CFE;
}
.FirstScreen{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    padding:1%;
    box-sizing: border-box;
    background-image: url(../assets/backimg.jpg);
    background-size: cover;
    background-repeat: no-repeat;
    // 头部
    .header{
        width: 100%;
        height: 66px;
        background-image: url(../assets/titles.png);
        background-size:100% 132px;
        background-repeat: no-repeat;
        text-align: center;
        position: relative;
        @font-face {
            font-family: aa;
            src: url(../assets/icon/FZMeiHei-M07.ttf);
        }
        p{
            font-family: aa;
            font-size: 30px;
            line-height: 58px;
            color: #fff;
            font-weight: 500;
        }
    }
   // 内容
    .nav{
        width: 100%;
        height: calc(100% - 126px);
        position: relative;
        .section{
            width: 100%;
            height: 65%;
            display: flex;
            flex-wrap:wrap;
            flex-direction: column-reverse;
            // 代码特效
            .codeEffects{
                width: 80%;
                height: 50%;
                position: fixed;
                top: 13%;
                left: 10%;
                z-index: 1;
                overflow: hidden;
                opacity: .3;
                .codeBox1{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 0;
                    background: url(../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
                @keyframes mymove
                {
                    0%  {top:-20%;}
                    10% {top:20%;}
                    20% {top:40%;}
                    30% {top:60%;}
                    40% {top:80%;}
                    50% {top:100%;}
                    60% {top:80%;}
                    70% {top:60%;}
                    80% {top:40%;}
                    90% {top:20%;}
                    100% {top:-20%;}
                }
                .codeBox2{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 20%;
                    background: url(../assets/code1.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove1 14s infinite alternate linear;
                }
                @keyframes mymove1
                {
                    0%  {bottom:-20%;}
                    10% {bottom:20%;}
                    20% {bottom:40%;}
                    30% {bottom:60%;}
                    40% {bottom:80%;}
                    50% {bottom:100%;}
                    60% {bottom:80%;}
                    70% {bottom:60%;}
                    80% {bottom:40%;}
                    90% {bottom:20%;}
                    100% {bottom:-20%;}
                }
                .codeBox3{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 40%;
                    background: url(../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
                .codeBox4{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 60%;
                    background: url(../assets/code1.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove1 14s infinite alternate linear;
                }
                .codeBox5{
                    width: 20%;
                    height: 100%;
                    position: absolute;
                    left: 80%;
                    background: url(../assets/code2.png) no-repeat top center;
                    background-size: 4% 100%;
                    animation:mymove 14s infinite alternate linear;
                }
            }
            // 粒子盒子
            .particleBox{
                width: 100%;
                height: 100%;
                position: absolute;
                top: 40%;
                left: 0;
                z-index: 2;
                opacity: 0.4;
            }
        }
        // 登录窗口
        .concentBox{
            width:28%;
            // height: 45%;
            overflow-y: scroll;
            background-color:rgba(0, 0, 0, 0.3);
            // border: 1px solid rgba(1, 57, 141, .9);
            box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .4);
            border-radius: 4px;
            box-sizing: border-box;
            padding-top:1.5%;
            box-sizing: border-box;
            position: absolute;
            top:50%;
            left:50%;
            transform:translate(-50%,-50%);
            z-index:30;
            .concentBoxloginForm{
                padding:6% 16% 5% 18%;
                box-sizing: border-box;
            }
            .imgcolor{
                width: 100%;
                height: 100%;
                
                cursor: pointer;
                filter: grayscale(100%);
                img{
                    
                    opacity: .4;
                    height: 100%;
                    margin-right:-9%;
                }
            }
        }
    }
    .footer{
        width: 100%;
        height:60px;
        padding-top:10px;
        box-sizing: border-box;
        .copyright{
            width:auto;
            height: 27px;
            margin: auto;
            text-align: center;
            opacity: .6;
            span{
                font-size: 12px;
                line-height: 24px;
                color: #00FFFF;
            }
            p{
                font-size: 12px;
                color: #00FFFF;
            }
        }
    }
}
// 测试九宫格
.FirstScreen-box{
    width: 100%;
    height: 100vh;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 10;
    li{
        width: 33.33%;
        height: 33.33%;
        border: 1px solid #ff6700;
        float: left;
        box-sizing: border-box;
    }
}
</style>
