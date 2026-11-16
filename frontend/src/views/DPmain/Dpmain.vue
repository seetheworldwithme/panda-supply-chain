<template>
    <!-- 第一屏LED看板222 -->  
    <div class="FirstScreen">
        
        <div class="header">
            <!-- <span>煤炭大宗数字交易平台</span> -->
            <div class="icon"><div v-if="accid != 4" class="echartsClick" size="mini"><span @click="echartsClick">返回平台</span></div><i v-if="accid == 4" class="iconfont el-icon-alituichu" @click="secede"></i></div>
        </div>
        <div class="nav">
            <div class="section">
                <!-- 粒子盒子 -->
                <div class="particleBox">
                    <vue-particles
                    color="#00FFFF"
                    :particleOpacity="0.4"
                    :particlesNumber="24"
                    shapeType="circle"
                    :particleSize="8"
                    linesColor="#00FFFF"
                    :linesWidth="1"
                    :lineLinked="true"
                    :lineOpacity=".3"
                    :linesDistance="350"
                    :moveSpeed="2"
                    :hoverEffect="false"
                    :clickEffect="false">
                    </vue-particles>
                </div>
            </div>
            <div class="concentBox">
                <el-carousel :interval="4000" type="card" height="600px" arrow="never" :autoplay="false">
                    <el-carousel-item v-for="(item,index) in itmes" :key="index" :label="item.name">
                        <img :src="item.img" alt="" @click="carouselClick(index)">
                    </el-carousel-item>
                </el-carousel>
            </div>
        </div>
    </div>
</template>
<script>
import imgone from "../../assets/cd1.jpg"
import imgtwo from "../../assets/cd2.jpg"
import imgtree from "../../assets/cd3.jpg"
import imgsix from "../../assets/rlhy.jpg"
import imgsenn from "../../assets/maps.jpg"
import imgsixs from "../../assets/imgsixs.png"

export default{
    data(){
        return {
            accid:'',
            itmes:[
                {
                    index:1,
                    name:'菜单一',
                    img:imgsix
                },
                {
                    index:2,
                    name:'菜单二',
                    img:imgone
                },
                {
                    index:3,
                    name:'菜单三',
                    img:imgtwo
                },
                {
                    index:4,
                    name:'菜单四',
                    img:imgtree
                },
                {
                    index:5,
                    name:'菜单五',
                    img:imgsenn
                },
                {
                    index:6,
                    name:'菜单六',
                    img:imgsixs
                },
            ],
        }
    },
    created(){
        this.accid = sessionStorage.getItem('accountTypes');
    },
    mounted() {

    },
    methods:{
        carouselClick(id){
            if(id == 0){
                this.$router.push({name:'Welcome'})
            }
            if(id == 1){
                this.$router.push({name:'Firstscreen'})
            }
            if(id == 2){
                this.$router.push({name:'Secondscreen'})
            }
            if(id == 3){
                this.$router.push({name:'Thirdscreen'})
            }
            if(id == 4){
                this.$router.push({name:'Maps'})
            }
            if(id == 5){
                this.$router.push({name:'warehouse'})
            }
        },
        // 退出
        secede(){
            sessionStorage.clear();
            this.$router.push({name:'DPLogin'})
        },
        echartsClick(){
            this.$router.push({name:'PtIndex'})
        }
    },
}
</script>
<style scoped lang="less">
// 轮播菜单
.el-carousel__item img {
    width: 100%;
    height: 100%;
}
.el-carousel__item:nth-child(2n) {
background-color: #99a9bf;
}
.el-carousel__item:nth-child(2n+1) {
background-color: #d3dce6;
}
.FirstScreen{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    padding:1%;
    box-sizing: border-box;
    background-image: url(../../assets/bg.png);
    background-size: cover;
    background-repeat: no-repeat;
    // 头部
    .header{
        width: 100%;
        text-align: left;
        position: relative;
        padding-left:2%;
        box-sizing: border-box;
        z-index: 99;
        @font-face {
            font-family: aa;
            src: url(../../assets/icon/FZMeiHei-M07.ttf);
        }
        span{
            font-family: aa;
            font-size: 48px;
            line-height: 128px;
            color: #00FFFF;
            font-weight: 500;
        }
        .icon{
            width: auto;
            float: right;
            padding-right:2%;
            box-sizing: border-box;
            .echartsClick{
                font-family: aa;
                width:100px;
                height: 128px;
                color: #00FFFF;
                float: left;
                text-align: center;
                cursor: pointer;
                span{
                    font-size: 18px;
                    cursor: 'pointer';
                    line-height: 128px;
                }
            }
            i{
                font-size: 16px;
                color: #00FFFF;
                line-height: 128px;
                cursor: pointer;
            }
        }
    }
   // 内容
    .nav{
        width: 100%;
        height: calc(100% - 66px);
        position: relative;
        .section{
            width: 100%;
            height: 65%;
            display: flex;
            flex-wrap:wrap;
            flex-direction: column-reverse;
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
        .concentBox{
            width: 100%;
            height: 90%;
            opacity: 0.7;
            position: absolute;
            top: 6%;
            left: 0;
            z-index: 30;
            padding-top:4%;
            box-sizing: border-box;
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
