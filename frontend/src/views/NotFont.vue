<template>
    <!-- 404 -->
    <div class="main">
        <div class="error">
            <img src="../assets/404.png" alt="">
            <p>页面丢失，请重新登录...</p>
            <h2><el-button @click="to" size="small">重 新 登 录</el-button></h2>
        </div>
    </div>

</template>
<script>
    export default{
        methods:{
            to(){
                sessionStorage.clear();
                this.$router.push({name:'Login'});
                window.location.reload(); // 整页刷新:路由表和按钮权限是按上一个账号裁剪过的,同标签页换账号会继承
            }
        }
    }
</script>

<style lang="less" scoped>
.main{
    width: 100%;
    height: 100vh;
    overflow: hidden;
    position: relative;
    .error{
        width:353px;
        height: 200px;
        position: absolute;
        top: -20%;
        bottom:0;
        left: 0;
        right: 0;
        margin: auto;
    }
    img{
        width: 353px;
        height: 200px;
    }
    p{
        font-size: 22px;
        // font-weight: 600;
        text-align: center;
        line-height: 68px;
        color: #333;
    }
    h2{
        text-align: center;
        color: #333;
        cursor: pointer;
    }
}
</style>